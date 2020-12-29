package rpc

import (
	"github.com/gorilla/mux"
)

func RegisterTestkitRPC(
	r *mux.Router,
	ctx *TestkitRPCContext,
) {
	// testkit related primitives
	r.HandleFunc("/init", handleInitTestkit(ctx, r)).Methods("POST")

	// async broadcast only
	r.HandleFunc("/{ctxId}/txs", handleTxInject(ctx)).Methods("POST")

	// manual injection
	r.HandleFunc("/{ctxId}/inject", handleBlockPropose(ctx)).Methods("POST")
	r.HandleFunc("/{ctxId}/automatic_tx", handleAutoTxGet(ctx)).Methods("GET")
	r.HandleFunc("/{ctxId}/automatic_tx", handleAutoTxRegister(ctx)).Methods("POST")
	r.HandleFunc("/{ctxId}/automatic_tx", handleAutoTxClearAll(ctx)).Methods("DELETE")      // register
	r.HandleFunc("/{ctxId}/automatic_tx/{atxId}", handleAutoTxClear(ctx)).Methods("DELETE") // delete
}
