package server

import (
	"fmt"
	client "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	terra "github.com/terra-money/core/app"
	compatlocalclient "github.com/terra-money/mantle-compatibility/localclient"
	"net/http"
	"sync"
)

type MantleLCDProxy struct {
	queryMtx *sync.Mutex
}

func NewMantleLCDServer(queryMtx *sync.Mutex) *MantleLCDProxy {
	return &MantleLCDProxy{
		queryMtx: queryMtx,
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
