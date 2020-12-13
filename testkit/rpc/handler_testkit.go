package rpc

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	tm "github.com/tendermint/tendermint/types"
	"github.com/terra-project/core/x/auth"
	"io/ioutil"
	"log"
	"net/http"
)

func handleTxInject(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ctx.tc == nil {
			panic("testkit is not started yet")
			return
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

		ctx.tc.AddToMempool(req.Tx)

		response := new(struct {
			Height string `json:'height'`
			TxHash string `json:"txhash"`
		})

		txHash := tm.Tx(codec.MustMarshalBinaryLengthPrefixed(req.Tx)).Hash()
		log.Printf("[mantle/testkit-rpc/tx] mempool tx (%X)\n", txHash)

		response.TxHash = fmt.Sprintf("%X", txHash)
		response.Height = fmt.Sprintf("%d", ctx.mantle.GetLastState().LastBlockHeight+1)

		w.WriteHeader(http.StatusOK)
		w.Write(MustMarshalJSON(response))
		return
	}
}

func handleBlockPropose(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ctx.tc == nil {
			panic("testkit is not started yet")
			return
		}

		vars := mux.Vars(r)
		valAddr, ok := vars["validatorAddress"]
		if !ok {
			panic("validatorAddress not given")
			return
		}

		val, err := sdk.ValAddressFromBech32(valAddr)
		if err != nil {
			panic("invalid validator address")
		}

		validatorPriv := ctx.tc.PickProposerByAddress(val)

		// inject!
		blockState, injectErr := ctx.tc.Inject(validatorPriv)
		if injectErr != nil {
			panic(injectErr.Error())
			return
		}

		response := codec.MustMarshalJSON(blockState)
		w.WriteHeader(http.StatusOK)
		w.Write(response)

		return
	}
}
