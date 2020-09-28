package app

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/terra-project/mantle/committer"

	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/depsresolver"
	"github.com/terra-project/mantle/graph"
	"github.com/terra-project/mantle/graph/schemabuilders"
	"github.com/terra-project/mantle/indexer"
	"github.com/terra-project/mantle/querier"
	reg "github.com/terra-project/mantle/registry"
	"github.com/terra-project/mantle/subscriber"
	"github.com/terra-project/mantle/types"
	"github.com/terra-project/mantle/utils"
)

type Mantle struct {
	isSynced          bool
	app               *App
	lifecycle         *LifecycleContext
	gqlInstance       *graph.GraphQLInstance
	committerInstance committer.Committer
	indexerInstance   *indexer.IndexerBaseInstance
	db                db.DB
}

type SyncConfiguration struct {
	TendermintEndpoint string
	SyncUntil uint64
}

func NewMantle(
	db db.DB,
	genesis *utils.GenesisDoc,
	indexers ...types.IndexerRegisterer,
) *Mantle {
	// create new terra app with postgres-patched KVStore
	app := NewApp(db.GetCosmosAdapter(), genesis)

	// create an auxiliary terra app lifecycle
	lc := NewLifecycle(app, false)

	// gather outputs of indexer registry
	registry := reg.NewRegistry(indexers)

	// initialize gql
	depsResolverInstance := depsresolver.NewDepsResolver()
	querierInstance := querier.NewQuerier(db, registry.KVIndexMap)

	// instantiate gql
	gqlInstance := graph.NewGraphQLInstance(
		depsResolverInstance,
		querierInstance,
		schemabuilders.CreateABCIStubSchemaBuilder(app.GetApp()),
		schemabuilders.CreateModelSchemaBuilder(nil, reflect.TypeOf((*types.BaseState)(nil))),
		schemabuilders.CreateModelSchemaBuilder(registry.KVIndexMap, registry.Models...),
	)

	// initialize committer
	committerInstance := committer.NewCommitter(db, registry.KVIndexMap)

	// initialize indexer
	indexerInstance := indexer.NewIndexerBaseInstance(
		registry.Indexers,
		registry.IndexerOutputs,
		gqlInstance.QueryInternal,
		gqlInstance.Commit,
	)

	return &Mantle{
		isSynced:          false,
		app:               app,
		lifecycle:         lc,
		committerInstance: committerInstance,
		gqlInstance:       gqlInstance,
		indexerInstance:   indexerInstance,
		db:                db,
	}
}

func (mantle *Mantle) QuerySync(configuration SyncConfiguration, currentBlockHeight int64) {
	remoteBlock, err := subscriber.GetBlock(fmt.Sprintf("http://%s/block", configuration.TendermintEndpoint))

	if err != nil {
		panic(fmt.Errorf("error during mantle sync: remote head fetch failed. fromHeight=%d, (%s)", currentBlockHeight, err))
	}

	remoteHeight := remoteBlock.Header.Height
	syncingBlockHeight := currentBlockHeight
	tStart := time.Now()

	for syncingBlockHeight < remoteHeight {
		// stop sync if SyncUntil is given
		if configuration.SyncUntil != 0 && uint64(syncingBlockHeight) == configuration.SyncUntil {
			for{}
		}
		remoteBlock, err := subscriber.GetBlock(fmt.Sprintf("http://%s/block?height=%d", configuration.TendermintEndpoint, syncingBlockHeight+1))
		if err != nil {
			panic(fmt.Errorf("error during mantle sync: remote block(%d) fetch failed", syncingBlockHeight))
		}

		// run round
		mantle.Inject(remoteBlock)

		syncingBlockHeight++
	}

	dur := time.Now().Sub(tStart)

	if dur > time.Second {
		log.Printf("[mantle] QuerySync: %d to %d, Elapsed: %dms", currentBlockHeight, remoteHeight, dur.Milliseconds())
	}
}

func (mantle *Mantle) Sync(configuration SyncConfiguration) {
	mantle.QuerySync(configuration, 1)
	// subscribe to NewBlock event
	rpcSubscription := subscriber.NewRpcSubscription(fmt.Sprintf("ws://%s/websocket", configuration.TendermintEndpoint))
	blockChannel := rpcSubscription.Subscribe()

	for {
		select {
		case block := <-blockChannel:
			lastBlockHeight := mantle.app.GetApp().LastBlockHeight()

			// stop sync if SyncUntil is given
			if configuration.SyncUntil != 0 && uint64(lastBlockHeight) == configuration.SyncUntil {
				for{}
			}

			if block.Header.Height-lastBlockHeight != 1 {
				mantle.QuerySync(configuration, lastBlockHeight)
			} else {
				mantle.Inject(&block)
			}
		}
	}
}

func (mantle *Mantle) Server() {
	go mantle.gqlInstance.ServeHTTP(1337)
}

func (mantle *Mantle) Rebuild() {

}

func (mantle *Mantle) Start() {

}

func (mantle *Mantle) Inject(block *types.Block) types.BaseState {
	height := block.Header.Height

	// force compact every 100th block
	if height%100 == 0 {
		if compactErr := mantle.db.Compact(); compactErr != nil {
			panic(compactErr)
		}
	}

	tStart := time.Now()
	baseState := mantle.lifecycle.Inject(block)
	mantle.gqlInstance.UpdateState(baseState)
	mantle.indexerInstance.RunIndexerRound()

	// flush states to database
	// note that indexer outputs are committed __BEFORE__ IAVL
	// because reversing indexer outputs is trivial (i.e. overwrite them)
	// whereas IAVL reversal is a little tricky.
	exportedStates := mantle.gqlInstance.ExportStates()
	err := mantle.committerInstance.Commit(uint64(height), exportedStates...)

	mantle.lifecycle.Commit()

	defer func() {
		mantle.gqlInstance.Flush()
	}()

	if err != nil {
		panic(err)
	}
	tEnd := time.Now()

	log.Printf(
		"[mantle] Indexing finished for block(%d), committing %d indexer outputs processed in %dms",
		height,
		len(exportedStates),
		tEnd.Sub(tStart).Milliseconds(),
	)

	return baseState
}

func start(
	baseStateEvent chan types.BaseState,
	gql *graph.GraphQLInstance,
	indexer *indexer.IndexerBaseInstance,
) {
	log.Print("Starting mantle...")
	for {
		baseState := <-baseStateEvent
		fmt.Println(baseState)
	}
	// for {
	// 	baseState := <-baseStateEvent
	// 	subcontext := gql.CreateSubContext(baseState)
	// 	indexer.RunCollectors()
	// }
}
