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
	r.HandleFunc("/{ctxId}/inject/{validatorAddress}", handleBlockPropose(ctx)).Methods("POST")
	r.HandleFunc("/{ctxId}/register_auto_tx", handleAutoTxRegister(ctx)).Methods("POST")
	r.HandleFunc("/{ctxId}/register_auto_tx_pause", handleAutoTxPauseRegister(ctx)).Methods("POST")
}
