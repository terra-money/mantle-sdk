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
	"github.com/terra-project/mantle-sdk/graph/generate"
	"github.com/terra-project/mantle-sdk/utils"
	"log"
	"os"
	"os/signal"
	"reflect"
	"runtime/debug"
	"sync"
	"syscall"
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
	app                  *TerraApp.TerraApp
	registry             *reg.Registry
	mantlemint           mantlemint.Mantlemint
	gqlInstance          *graph.GraphQLInstance
	depsResolverInstance depsresolver.DepsResolver
	committerInstance    committer.Committer
	indexerInstance      *indexer.IndexerBaseInstance
	db                   db.DB
	m                    *sync.Mutex
}

type SyncConfiguration struct {
	TendermintEndpoint string
	SyncUntil          uint64
	Reconnect          bool
	OnWSError          func(err error)
	OnInjectError      func(err error)
}

var (
	errInvalidBlock = fmt.Errorf("invalid block")
)

func NewMantle(
	db db.DB,
	genesis *tmtypes.GenesisDoc,
	indexers ...types.IndexerRegisterer,
) (mantleApp *Mantle) {
	// wrap db w/ force global_transaction manager

	// create new terra app with postgres-patched KVStore
	tmdb := db.GetCosmosAdapter()
	terraApp := compatapp.NewTerraApp(tmdb)

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
		schemabuilders.CreateMantleStateSchemaBuilder(nil, nil),
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
		m:                    new(sync.Mutex),
	}

	// create a signal handler
	sigChannel := make(chan os.Signal, 1)
	mantleApp.gracefulShutdownOnSignal(
		sigChannel,
		func() { db.Purge(false) },
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	mantleApp.gracefulShutdownOnSignal(
		sigChannel,
		func() { db.Purge(true) },
		syscall.SIGILL,
		syscall.SIGABRT,
		syscall.SIGKILL,
		syscall.SIGHUP,
	)

	// initialize within transaction boundary
	mantleApp.db.SetCriticalZone()

	// initialize chain
	if initErr := mantleApp.mantlemint.Init(genesis); initErr != nil {
		panic(initErr)
	}

	if releaseErr := mantleApp.db.ReleaseCriticalZone(); releaseErr != nil {
		panic(releaseErr)
	}

	return mantleApp
}

func (mantle *Mantle) GetApp() *TerraApp.TerraApp {
	return mantle.app
}

func (mantle *Mantle) SetBlockExecutor(be mantlemint.MantlemintExecutorCreator) {
	mantle.mantlemint.SetBlockExecutor(be)
}

func (mantle *Mantle) indexerLifecycle(responses state.ABCIResponses) {
	var block = mantle.mantlemint.GetCurrentBlock()
	var height = block.Header.Height

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
}

func (mantle *Mantle) querySync(configuration SyncConfiguration, currentBlockHeight int64) {
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
			if configuration.OnInjectError != nil {
				configuration.OnInjectError(err)
			} else {
				panic(err)
			}
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
	rpcSubscription, connRefused := subscriber.NewRpcSubscription(
		fmt.Sprintf("ws://%s/websocket", configuration.TendermintEndpoint),
		configuration.OnWSError,
	)

	// connRefused here is most likely triggered by ECONNREFUSED
	// in case reconnect flag is set, try reestablish the connection after 5 seconds.
	if connRefused != nil {
		if configuration.Reconnect {
			select {
			case <-time.NewTimer(5 * time.Second).C:
				mantle.Sync(configuration)
			}
			return

		} else {
			panic(connRefused)
		}
	}

	blockChannel := rpcSubscription.Subscribe(configuration.Reconnect)

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
				log.Printf("lastBlockHeight=%v, remoteBlockHeight=%v\n", lastBlockHeight, block.Header.Height)
				mantle.querySync(configuration, lastBlockHeight)
			} else {
				if _, err := mantle.Inject(&block); err != nil {
					// if OnInjectError is set,
					// relay injection error to the caller
					if configuration.OnInjectError != nil {
						configuration.OnInjectError(err)
					} else {
						panic(err)
					}
				}
			}
		}
	}
}

func (mantle *Mantle) Server(port int) {
	go mantle.gqlInstance.ServeHTTP(port)
}

func (mantle *Mantle) Inject(block *types.Block) (*types.BlockState, error) {
	// handle any inevitable panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("!! mantle panicked w/ message: %s", r)
			debug.PrintStack()
			log.Print("[mantle] panic during inject, attempting graceful shutdown")

			// if mantle reaches this point, there was a panic during injection.
			// in such case db access is all gone, it is safe to NOT get a lock.
			// but doing it, just in case :)
			mantle.m.Lock() // never unlock
			mantle.db.Purge(true)
			log.Print("[mantle] shutdown done")
			os.Exit(0)
		}
	}()

	mantle.m.Lock()
	defer mantle.m.Unlock()

	// set global global_transaction boundary for
	// tendermint, cosmos, mantle
	mantle.db.SetCriticalZone()

	// time
	tStart := time.Now()

	// inject this block
	blockState, injectErr := mantle.mantlemint.Inject(block)

	// flush to db after injection is done
	// never call this in defer, as in panic cases we need to be able to revert this commit
	mantle.db.ReleaseCriticalZone()

	// if injection was successful,
	// flush all to disk
	if injectErr != nil {
		return blockState, injectErr
	}

	mantle.db.ReleaseCriticalZone()

	// time end
	tEnd := time.Now()

	log.Printf(
		"[mantle] Indexing finished for block(%d), processed in %dms",
		block.Header.Height,
		tEnd.Sub(tStart).Milliseconds(),
	)

	return blockState, injectErr
}

func (mantle *Mantle) LocalQuery(query interface{}, variables types.GraphQLParams) error {
	qs := generate.GenerateQuery(query, variables)
	res := mantle.gqlInstance.QueryInternal(
		qs,
		variables,
		nil,
	)

	return graph.UnmarshalInternalQueryResult(res, query)
}

func (mantle *Mantle) ExportStates() map[string]interface{} {
	return mantle.depsResolverInstance.GetState()
}

func (mantle *Mantle) GetLastState() state.State {
	return mantle.mantlemint.GetCurrentState()
}

func (mantle *Mantle) GetLastHeight() int64 {
	return mantle.mantlemint.GetCurrentHeight()
}

func (mantle *Mantle) GetLastBlock() *types.Block {
	return mantle.mantlemint.GetCurrentBlock()
}

func (mantle *Mantle) gracefulShutdownOnSignal(
	sig chan os.Signal,
	callback func(),
	signalTypes ...os.Signal,
) {
	// handle
	signal.Notify(
		sig,
		signalTypes...,
	)

	go func() {
		received := <-sig
		log.Printf("[mantle] received %v, cleanup...", received.String())
		// wait until inject is cleared
		log.Printf("[mantle] attempting graceful shutdown...")
		mantle.m.Lock() // never unlock
		callback()
		log.Printf("[mantle] shutdown done")
		os.Exit(0)
	}()
}
