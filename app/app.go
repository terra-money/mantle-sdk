package app

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/terra-project/mantle/committer"

	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/graph"
	"github.com/terra-project/mantle/graph/depsresolver"
	"github.com/terra-project/mantle/graph/schemabuilders"
	"github.com/terra-project/mantle/indexer"
	"github.com/terra-project/mantle/querier"
	"github.com/terra-project/mantle/registry"
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
}

func NewMantle(
	db db.DB,
	transactionalAppState bool,
	genesisPath string,
	indexers []types.IndexerRegisterer,
) *Mantle {
	genesis := utils.GenesisDocFromFile(genesisPath)

	// create new terra app with postgres-patched KVStore
	app := NewApp(db.GetCosmosAdapter(), genesis)

	// create an auxiliary terra app lifecycle
	lc := NewLifecycle(app, transactionalAppState)

	// gather outputs of indexer registry
	registry := registry.NewRegistry(indexers)

	// initialize gql
	depsResolverInstance := depsresolver.NewDepsResolver()
	querierInstance := querier.NewQuerier(db, registry.KVIndexMap)
	gqlInstance := graph.NewGraphQLInstance(
		depsResolverInstance,
		querierInstance,
		schemabuilders.CreateABCIStubSchemaBuilder(app.GetApp()),
		schemabuilders.CreateModelSchemaBuilder(reflect.TypeOf((*types.BaseState)(nil))),
		schemabuilders.CreateModelSchemaBuilder(registry.Models...),
		schemabuilders.CreateListSchemaBuilder(),
	)

	// initialize committer
	committerInstance := committer.NewCommitter(db, registry.KVIndexMap)

	// initialize indexer
	indexerInstance := indexer.NewIndexerBaseInstance(
		registry.Indexers,
		registry.IndexerOutputs,
		gqlInstance.ResolveQuery,
		gqlInstance.Commit,
	)

	return &Mantle{
		isSynced:          false,
		app:               app,
		lifecycle:         lc,
		committerInstance: committerInstance,
		gqlInstance:       gqlInstance,
		indexerInstance:   indexerInstance,
	}
}

func (mantle *Mantle) Sync() {
	currentBlockHeight := mantle.app.GetApp().LastBlockHeight()
	remoteBlock, err := subscriber.GetBlockLCD("https://tequila-fcd.terra.dev/blocks/latest")
	if err != nil {
		panic(fmt.Errorf("error during mantle sync: remote head fetch failed. fromHeight=%d, (%s)", currentBlockHeight, err))
	}

	remoteHeight := remoteBlock.Header.Height

	if remoteHeight <= currentBlockHeight {
		log.Printf("[mantle] Sync unnecessary, remoteHeight=%d, currentBlockHeight=%d", remoteHeight, currentBlockHeight)
		return
	}

	syncingBlockHeight := currentBlockHeight + 1
	for syncingBlockHeight < remoteHeight {
		remoteBlock, err := subscriber.GetBlockLCD(fmt.Sprintf("https://tequila-fcd.terra.dev/blocks/%d", syncingBlockHeight))
		if err != nil {
			panic(fmt.Errorf("error during mantle sync: remote block(%d) fetch failed", syncingBlockHeight))
		}

		log.Printf("[mantle] Syncing block(%d)", remoteBlock.Header.Height)

		// run round
		mantle.round(remoteBlock)

		syncingBlockHeight++
	}
}

func (mantle *Mantle) Server() {
	go mantle.gqlInstance.ServeHTTP(1337)
}

func (mantle *Mantle) Rebuild() {

}

func (mantle *Mantle) Start() {

}

func (mantle *Mantle) round(block *types.Block) {
	height := block.Header.Height

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
