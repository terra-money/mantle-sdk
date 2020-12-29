package rpc

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	tm "github.com/tendermint/tendermint/types"
	"github.com/terra-project/core/x/auth"
	"github.com/terra-project/mantle-sdk/testkit"
	"io/ioutil"
	"net/http"
)

func handleTxInject(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctxId, ok := vars["ctxId"]

		if !ok {
			panic("invalid ctxId")
		}

		var req struct {
			Tx   auth.StdTx `json:"tx" yaml:"tx"`
			Mode string     `json:"mode" yaml:"mode"`
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
			return
		}

		err = codec.UnmarshalJSON(body, &req)
		if err != nil {
			panic(err.Error())
			return
		}

		// log
		txHash := tm.Tx(codec.MustMarshalBinaryLengthPrefixed(req.Tx)).Hash()

		blockState, injectionErr := ctx.GetContext(ctxId).AddToMempool(req.Tx)

		if injectionErr != nil {
			panic(injectionErr)
		}

		txResult := blockState.ResponseDeliverTx[0]
		response := new(struct {
			Height    string              `json:"height"`
			TxHash    string              `json:"txhash"`
			RawLog    string              `json:"raw_log"`
			Log       sdk.ABCIMessageLogs `json:"logs"`
			GasWanted int64               `json:"string,gas_wanted"`
			GasUsed   int64               `json:"string,gas_used"`
			Code      uint32              `json:"code"`
			Codespace string              `json:"codespace"`
			Tx        auth.StdTx          `json:"tx"`
		})
		response.TxHash = fmt.Sprintf("%X", txHash)
		response.Height = fmt.Sprintf("%d", ctx.GetContext(ctxId).GetMantle().GetLastState().LastBlockHeight+1)
		response.RawLog = txResult.GetLog()
		response.GasWanted = txResult.GetGasWanted()
		response.GasUsed = txResult.GetGasUsed()
		response.Code = txResult.GetCode()
		response.Codespace = txResult.GetCodespace()
		response.Tx = req.Tx

		parsedLogs, _ := sdk.ParseABCILogs(txResult.Log)
		response.Log = parsedLogs

		w.WriteHeader(http.StatusOK)
		w.Write(MustMarshalJSON(response))
		return
	}
}

func handleBlockPropose(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctxId, ok := vars["ctxId"]

		if !ok {
			panic("invalid ctxId")
		}

		// inject!
		blockState, injectErr := ctx.GetContext(ctxId).Inject()
		if injectErr != nil {
			panic(injectErr.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(codec.MustMarshalJSON(blockState))
	}
}

func handleAutoTxRegister(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctxId, _ := vars["ctxId"]

		// read body
		autoTxRequest := AutomaticTxRequest{}
		bz, err := ioutil.ReadAll(r.Body)

		if err != nil {
			panic(err)
		}
		if err := codec.UnmarshalJSON(bz, &autoTxRequest); err != nil {
			panic(err)
		}

		entry := testkit.NewAutomaticTxEntry(
			autoTxRequest.AccountName,
			autoTxRequest.Fee,
			autoTxRequest.Msgs,
			autoTxRequest.Period,
			autoTxRequest.Offset,
			autoTxRequest.StartAt,
		)

		tctx := ctx.GetContext(ctxId)
		tctx.AddAutomaticTxEntry(entry)

		w.WriteHeader(200)
		w.Write(codec.MustMarshalJSON(entry))
	}
}

func handleAutoTxGet(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctxId, _ := vars["ctxId"]

		entries := ctx.GetContext(ctxId).GetAutomaticTxEntries()

		w.WriteHeader(200)
		w.Write(codec.MustMarshalJSON(entries))
	}
}

func handleAutoTxClearAll(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctxId, _ := vars["ctxId"]

		ctx.GetContext(ctxId).ClearAllAutomaticTxEntries()

		w.WriteHeader(200)
	}
}

func handleAutoTxClear(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctxId, _ := vars["ctxId"]
		atxId, ok := vars["atxId"]

		if !ok {
			fmt.Println("wtf??")
			panic("automatic tx id is not provided")
		}

		ctx.GetContext(ctxId).ClearAutomaticTxEntry(atxId)

		w.WriteHeader(200)
	}
}
