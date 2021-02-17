package app

import (
	"encoding/json"
	"fmt"
	"github.com/terra-project/mantle-sdk/committer"
	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/depsresolver"
	"github.com/terra-project/mantle-sdk/graph"
	"github.com/terra-project/mantle-sdk/graph/schemabuilders"
	"github.com/terra-project/mantle-sdk/indexer"
	"github.com/terra-project/mantle-sdk/querier"
	reg "github.com/terra-project/mantle-sdk/registry"
	"github.com/terra-project/mantle-sdk/types"
	"io/ioutil"
	"net/http"
	"time"
)

type RemoteMantle struct {
	db                   db.DB
	registry             *reg.Registry
	gqlInstance          *graph.RemoteGraphQLInstance
	depsResolverInstance depsresolver.DepsResolver
	committerInstance    committer.Committer
	indexerInstance      *indexer.IndexerBaseInstance
	baseMantleEndpoint   string
}

type RemoteSyncConfiguration struct {
	SyncUntil   uint64
	SyncFrom    uint64
	Reconnect   bool
	OnWSError   func(err error)
	OnInjectErr func(err error)
}

// NewRemoteMantle creates a mantle instance
// where there is no mantlemint
func NewRemoteMantle(
	db db.DB,
	baseMantleEndpoint string,
	indexers ...types.IndexerRegisterer,
) (mantleApp *RemoteMantle) {

	registry := reg.NewRegistry(indexers)
	remoteDepsResolver := depsresolver.NewRemoteDepsResolver()

	querierInstance := querier.NewQuerier(db, registry.KVIndexMap)

	// create gql instance w/ remote deps resolver and only the injected indexers
	gqlInstance := graph.NewRemoteGraphQLInstance(
		remoteDepsResolver,
		querierInstance,
		baseMantleEndpoint,
		schemabuilders.CreateRemoteModelSchemaBuilder(baseMantleEndpoint),
		schemabuilders.CreateModelSchemaBuilder(registry.KVIndexMap, registry.Models...),
	)

	// initializer committer
	committerInstance := committer.NewCommitter(db, registry.KVIndexMap)

	// initialize indexer
	indexerInstance := indexer.NewIndexerBaseInstance(
		registry.Indexers,
		registry.IndexerOutputs,
		gqlInstance.QueryInternal,
		gqlInstance.Commit,
	)

	// initialize deps resolver -- still required for inter-indexer sync
	depsResolverInstance := depsresolver.NewDepsResolver()

	// create remote mantle
	mantleApp = &RemoteMantle{
		db:                   db,
		registry:             &registry,
		gqlInstance:          gqlInstance,
		depsResolverInstance: depsResolverInstance,
		committerInstance:    committerInstance,
		indexerInstance:      indexerInstance,
		baseMantleEndpoint:   baseMantleEndpoint,
	}

	return
}

var LastSyncedHeightQuery = "query{LastSyncedHeight}"

// Sync starts mantle as a remote mode.
// In remote mode the chain data is not synced (i.e. no Inject happens)
// and only indexers are run
type LastSyncedHeightResponse struct {
	Data struct {
		LastSyncedHeight uint64 `json:"LastSyncedHeight"`
	} `json:"data"`
}

func (rmantle *RemoteMantle) Sync(config RemoteSyncConfiguration) {
	var lastKnownHeight uint64 = 0
	// listen to LastSyncHeight change
	// trigger indexer
	// TODO: refactor this into gql subscription & reconnect logic
	for {
		// poll every 200ms
		time.Sleep(200 * time.Millisecond)

		// get currentHeight
		currentHeight := getLastSyncedHeight(rmantle.baseMantleEndpoint)
		if currentHeight <= lastKnownHeight {
			continue
		}

		rmantle.indexerInstance.RunIndexerRound()
		indexerOutputs := rmantle.depsResolverInstance.GetState()

		// convert indexer outputs to slice
		var commitTargets = make([]interface{}, len(indexerOutputs))
		var i = 0
		for _, output := range indexerOutputs {
			commitTargets[i] = output
			i++
		}

		// set db critical zones
		rmantle.db.SetCriticalZone()

		// commit
		if commitErr := rmantle.committerInstance.Commit(currentHeight, commitTargets...); commitErr != nil {
			panic(commitErr)
		}

		// release db lock
		releaseErr := rmantle.db.ReleaseCriticalZone()
		if releaseErr != nil {
			panic(releaseErr)
		}
	}
}

func (rmantle *RemoteMantle) Server(port int) {
	go rmantle.gqlInstance.ServeHTTP(port)
}

func getLastSyncedHeight(baseMantleEndpoint string) uint64 {
	response, err := http.Get(fmt.Sprintf(
		"%s?query=%s",
		baseMantleEndpoint,
		LastSyncedHeightQuery,
	))

	if err != nil {
		panic(err)
	}

	bz, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	sync := LastSyncedHeightResponse{}
	if err := json.Unmarshal(bz, &sync); err != nil {
		panic(err)
	}

	return sync.Data.LastSyncedHeight
}
