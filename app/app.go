package app

import (
	"fmt"
	"log"

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
	isSynced bool
	app *App
	lifecycle *LifecycleContext
	gqlInstance *graph.GraphQLInstance
	indexerInstance *indexer.IndexerBaseInstance
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

	// create rpc subscriber
	//rpc := subscriber.NewRpcSubscription(
	//	rpcEndpoint,
	//	lcdEndpoint,
	//)

	// create subscriber -> app channel
	//blockEventChannel := rpc.Subscribe() // publishes to blockEventChannel
	//lifecycleEventChannel := lc.Start(blockEventChannel)

	// gather outputs of indexer registry
	registry := registry.NewRegistry(indexers)

	// initialize gql
	depsResolverInstance := depsresolver.NewDepsResolver()
	querierInstance := querier.NewQuerier(db, registry.KVIndexMap)
	gqlInstance := graph.NewGraphQLInstance(
		depsResolverInstance,
		querierInstance,
		schemabuilders.CreateABCIStubSchemaBuilder(app.GetApp()),
		schemabuilders.CreateModelSchemaBuilder(registry.Models...),
	)

	// initialize indexer
	indexerInstance := indexer.NewIndexerBaseInstance(
		registry.Indexers,
		gqlInstance.ResolveQuery,
		gqlInstance.Commit,
	)

	return &Mantle{
		isSynced: false,
		app: app,
		lifecycle: lc,
		gqlInstance: gqlInstance,
		indexerInstance: indexerInstance,
	}
}

func (mantle *Mantle) Sync() {
	for !mantle.isSynced {
		currentBlockHeight := mantle.app.GetApp().LastBlockHeight()

		mantle.isSynced = true
	}
}

func (mantle *Mantle) Start() {

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
