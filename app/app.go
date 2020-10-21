package app

import (
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/state"
	TerraApp "github.com/terra-project/core/app"
	compatapp "github.com/terra-project/mantle-compatibility/app"
	"github.com/terra-project/mantle/app/mantlemint"
	"github.com/terra-project/mantle/app/middlewares"
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
	isSynced             bool
	app                  *TerraApp.TerraApp
	registry             *reg.Registry
	mantlemint           mantlemint.Mantlemint
	gqlInstance          *graph.GraphQLInstance
	depsResolverInstance depsresolver.DepsResolver
	committerInstance    committer.Committer
	indexerInstance      *indexer.IndexerBaseInstance
	db                   db.DB
}

type SyncConfiguration struct {
	TendermintEndpoint string
	SyncUntil          uint64
}

var (
	GlobalTerraApp *TerraApp.TerraApp
)

func NewMantle(
	db db.DB,
	genesis *utils.GenesisDoc,
	indexers ...types.IndexerRegisterer,
) (mantleApp *Mantle) {
	// create new terra app with postgres-patched KVStore
	tmdb := db.GetCosmosAdapter()
	terraApp := compatapp.NewTerraApp(tmdb)
	GlobalTerraApp = terraApp

	// gather outputs of indexer registry
	registry := reg.NewRegistry(indexers)

	// initialize gql
	depsResolverInstance := depsresolver.NewDepsResolver()
	querierInstance := querier.NewQuerier(db, registry.KVIndexMap)

	// instantiate gql
	gqlInstance := graph.NewGraphQLInstance(
		depsResolverInstance,
		querierInstance,
		schemabuilders.CreateABCIStubSchemaBuilder(terraApp),
		schemabuilders.CreateModelSchemaBuilder(nil, reflect.TypeOf((*types.BlockState)(nil))),
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

	mantleApp = &Mantle{
		isSynced: false,
		app:      terraApp,
		registry: nil,
		mantlemint: mantlemint.NewMantlemint(
			tmdb,
			terraApp,

			// in order to prevent indexer output to be in disparity
			// w/ tendermint state, indexers MUST run before commit.
			// use middlewares to run indexers (in CommitSync/CommitAsync)
			middlewares.NewIndexerMiddleware(func(responses state.ABCIResponses) {
				mantleApp.indexerLifecycle(responses)
			}),
		),
		gqlInstance:          gqlInstance,
		depsResolverInstance: depsResolverInstance,
		committerInstance:    committerInstance,
		indexerInstance:      indexerInstance,
		db:                   db,
	}

	// initialize chain
	if initErr := mantleApp.mantlemint.Init(genesis); initErr != nil {
		panic(initErr)
	}

	return mantleApp
}

func (mantle *Mantle) indexerLifecycle(responses state.ABCIResponses) {
	var block = mantle.mantlemint.GetCurrentBlock()
	var height = block.Header.Height

	tStart := time.Now()

	// create blockState
	var deliverTxsCopy = make([]abci.ResponseDeliverTx, len(responses.DeliverTxs))
	for i, deliverTx := range responses.DeliverTxs {
		deliverTxsCopy[i] = *deliverTx
	}
	blockState := types.BlockState{
		Height:             block.Height,
		ResponseBeginBlock: *responses.BeginBlock,
		ResponseEndBlock:   *responses.EndBlock,
		ResponseDeliverTx:  deliverTxsCopy,
		Block:              *block,
	}

	// set BlockState in depsResolver
	mantle.depsResolverInstance.SetPredefinedState(blockState)

	// RunIndexerRound panics when an indexer fails
	mantle.indexerInstance.RunIndexerRound()

	// flush states to database
	// note that indexer outputs are committed __BEFORE__ IAVL
	// because reversing indexer outputs is trivial (i.e. overwrite them)
	// whereas IAVL reversal is a little tricky.
	indexerOutputs := mantle.depsResolverInstance.GetState()
	defer mantle.depsResolverInstance.Dispose()

	// convert indexer outputs to slice
	var commitTargets = make([]interface{}, len(indexerOutputs))
	var i = 0
	for _, output := range indexerOutputs {
		commitTargets[i] = output
		i++
	}

	if commitErr := mantle.committerInstance.Commit(uint64(height), commitTargets...); commitErr != nil {
		panic(commitErr)
	}

	tEnd := time.Now()

	log.Printf(
		"[mantle] Indexing finished for block(%d), processed in %dms",
		height,
		tEnd.Sub(tStart).Milliseconds(),
	)
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
			for {
			}
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
	// subscribe to NewBlock event
	rpcSubscription := subscriber.NewRpcSubscription(fmt.Sprintf("ws://%s/websocket", configuration.TendermintEndpoint))
	blockChannel := rpcSubscription.Subscribe()

	for {
		select {
		case block := <-blockChannel:
			lastBlockHeight := mantle.app.LastBlockHeight()

			// stop sync if SyncUntil is given
			if configuration.SyncUntil != 0 && uint64(lastBlockHeight) == configuration.SyncUntil {
				for {
				}
			}

			if block.Header.Height-lastBlockHeight != 1 {
				mantle.QuerySync(configuration, lastBlockHeight)
			} else {
				mantle.Inject(&block)
			}
		}
	}
}

func (mantle *Mantle) Server(port int) {
	go mantle.gqlInstance.ServeHTTP(port)
}

func (mantle *Mantle) Inject(block *types.Block) error {
	// inject
	if err := mantle.mantlemint.Inject(block); err != nil {
		return err
	}

	return nil
}

func (mantle *Mantle) ExportStates() map[string]interface{} {
	return mantle.depsResolverInstance.GetState()
}
