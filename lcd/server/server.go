package server

import (
	"fmt"
	client "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	lru "github.com/hashicorp/golang-lru"
	compatlocalclient "github.com/terra-money/mantle-compatibility/localclient"
	abcistub "github.com/terra-money/mantle-sdk/graph/schemabuilders/abcistub"
	terra "github.com/terra-project/core/app"
	"net/http"
	"sync"
)

type MantleLCDProxy struct {
	queryMtx *sync.Mutex
	cache    *lru.Cache
}

func NewMantleLCDServer(queryMtx *sync.Mutex, cache *lru.Cache) *MantleLCDProxy {
	return &MantleLCDProxy{
		queryMtx: queryMtx,
		cache:    cache,
	}
}

func (lcd *MantleLCDProxy) Server(port int, app *terra.TerraApp) {
	router := mux.NewRouter().SkipClean(true)
	localClient := compatlocalclient.NewLocalClient(app, lcd.queryMtx)

	ctx := client.
		NewCLIContext().
		WithTrustNode(true).
		WithCodec(terra.MakeCodec()).
		WithClient(localClient)

	terra.ModuleBasics.RegisterRESTRoutes(ctx, router)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
