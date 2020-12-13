package rpc

import (
	"github.com/gorilla/mux"
)

func RegisterTestkitRPC(
	r *mux.Router,
	ctx *TestkitRPCContext,
) {
	genesis := r.PathPrefix("/genesis").Subrouter()
	genesis.Use(FailIfGenesisNotInitializedMiddleware(ctx))
	genesis.HandleFunc("/add-account/{accountName}", handleAddAccount(ctx)).Methods("POST")
	genesis.HandleFunc("/create-validator/{accountName}", handleCreateValidator(ctx)).Methods("POST")
	genesis.HandleFunc("/accounts", handleGetAccounts(ctx)).Methods("GET")
	genesis.HandleFunc("/validators", handleGetValidators(ctx)).Methods("GET")

	// /genesis/seal is the testkit starter
	// after this route is handled, no further POST actions to genesis can be done
	genesis.HandleFunc("/seal", handleSeal(ctx, r)).Methods("POST")

	// testkit related primitives
	r.HandleFunc("/genesis/{chainId}", handleInitGenesis(ctx)).Methods("POST")

	// async broadcast only
	r.HandleFunc("/txs", handleTxInject(ctx)).Methods("POST")
	r.HandleFunc("/inject/{validatorAddress}", handleBlockPropose(ctx)).Methods("POST")
}
