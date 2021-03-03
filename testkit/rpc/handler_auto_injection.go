package rpc

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type AutoInjectionRequest struct {
	validatorRounds []string `json:"validator_rounds"`
}

func handleSetAutoInjection(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctxId, ok := vars["ctxId"]

		request := AutoInjectionRequest{}
		dec := json.NewDecoder(r.Body)
		dec.Decode(&request)

		if !ok {
			panic("invalid ctxId")
		}

		// inject!
		ctx.GetContext(ctxId).SetAutomaticInjection(request.validatorRounds)

		w.WriteHeader(http.StatusOK)
	}
}

func handleDisableAutoInjection(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctxId, ok := vars["ctxId"]

		if !ok {
			panic("invalid ctxId")
		}

		ctx.GetContext(ctxId).DeleteAutomaticInjection()

		w.WriteHeader(http.StatusOK)
	}
}
