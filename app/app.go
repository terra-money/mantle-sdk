package app

import (
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/state"
	tmtypes "github.com/tendermint/tendermint/types"
	TerraApp "github.com/terra-project/core/app"
	compatapp "github.com/terra-project/mantle-compatibility/app"
	"github.com/terra-project/mantle-sdk/app/mantlemint"
	"github.com/terra-project/mantle-sdk/app/middlewares"
	"github.com/terra-project/mantle-sdk/utils"
	"log"
	"reflect"
	"time"

	"github.com/terra-project/mantle-sdk/committer"

	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/depsresolver"
	"github.com/terra-project/mantle-sdk/graph"
	"github.com/terra-project/mantle-sdk/graph/schemabuilders"
	"github.com/terra-project/mantle-sdk/indexer"
	"github.com/terra-project/mantle-sdk/querier"
	reg "github.com/terra-project/mantle-sdk/registry"
	"github.com/terra-project/mantle-sdk/subscriber"
	"github.com/terra-project/mantle-sdk/types"
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
	genesis *tmtypes.GenesisDoc,
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
		Block:              utils.ConvertBlockToRawBlock(block),
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
	log.Println("Local blockchain is behind, syncing previous blocks...")
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
		if _, err := mantle.Inject(remoteBlock); err != nil {
			panic(err)
		}

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

			log.Printf("lastBlockHeight=%v, remoteBlockHeight=%v\n", lastBlockHeight, block.Header.Height)

			// stop sync if SyncUntil is given
			if configuration.SyncUntil != 0 && uint64(lastBlockHeight) == configuration.SyncUntil {
				for {
				}
			}

			if block.Header.Height-lastBlockHeight != 1 {
				mantle.QuerySync(configuration, lastBlockHeight)
			} else {
				if _, err := mantle.Inject(&block); err != nil {
					panic(err)
				}
			}
		}
	}
}

func (mantle *Mantle) Server(port int) {
	go mantle.gqlInstance.ServeHTTP(port)
}

func (mantle *Mantle) Inject(block *types.Block) (*types.BlockState, error) {
	return mantle.mantlemint.Inject(block)
}

func (mantle *Mantle) ExportStates() map[string]interface{} {
	return mantle.depsResolverInstance.GetState()
}
