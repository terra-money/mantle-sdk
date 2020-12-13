package rpc

import (
	"github.com/gorilla/mux"
	terra "github.com/terra-project/core/app"
	"github.com/terra-project/mantle-sdk/app"
	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/testkit"
)

var codec = terra.MakeCodec()

type TestkitRPCContext struct {
	mantle *app.Mantle
	tc     *testkit.TestkitContext
	tg     *testkit.TestkitGenesis
	db     db.DB
}

type TestkitCreatorFunc func(r *mux.Router, ctx *TestkitRPCContext)
