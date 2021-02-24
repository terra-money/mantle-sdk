package server

import (
	"fmt"
	client "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	terra "github.com/terra-project/core/app"
	compatlocalclient "github.com/terra-project/mantle-compatibility/localclient"
	"net/http"
	"sync"
)

type MantleLCDProxy struct{}

func NewMantleLCDServer() *MantleLCDProxy {
	return &MantleLCDProxy{}
}

func (lcd *MantleLCDProxy) Server(port int, app *terra.TerraApp) {
	router := mux.NewRouter().SkipClean(true)
	m := &sync.Mutex{}

	localClient := compatlocalclient.NewLocalClient(app, m)

	ctx := client.
		NewCLIContext().
		WithTrustNode(true).
		WithCodec(terra.MakeCodec()).
		WithClient(localClient)

	terra.ModuleBasics.RegisterRESTRoutes(ctx, router)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
