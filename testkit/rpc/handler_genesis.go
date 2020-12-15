package rpc

import (
	"encoding/json"
	"fmt"
	client "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/mantle-sdk/app/mantlemint"
	"math/rand"
	"time"

	"github.com/gorilla/mux"
	terra "github.com/terra-project/core/app"
	authrest "github.com/terra-project/core/x/auth/client/rest"
	"github.com/terra-project/mantle-compatibility/localclient"
	"github.com/terra-project/mantle-compatibility/types"
	"github.com/terra-project/mantle-sdk/app"
	"github.com/terra-project/mantle-sdk/db/leveldb"
	"github.com/terra-project/mantle-sdk/testkit"
	"net/http"
	"sync"
)

func handleInitGenesis(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ctx.tg != nil {
			panic("there is an ongoing genesis, please purge testkit before setting up a new genesis")
		}

		vars := mux.Vars(r)
		chainId, ok := vars["chainId"]
		if !ok {
			panic("chainId is not correct or not supplied")
		}
		ctx.tg = testkit.NewTestkitGenesis(chainId)

		response := new(struct {
			ChainID string `json:"chainId"`
		})

		response.ChainID = chainId
		w.WriteHeader(http.StatusOK)
		w.Write(MustMarshalJSON(response))
	}
}

func handleGetAccounts(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(codec.MustMarshalJSON(ctx.tg.GetAccounts()))

		return
	}
}

func handleGetValidators(ctx *TestkitRPCContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validators := ctx.tg.GetValidators()

		response := make([]string, len(validators))
		for idx, val := range validators {
			response[idx] = val.Account.String()
		}

		w.WriteHeader(http.StatusOK)
		w.Write(codec.MustMarshalJSON(response))

		return
	}
}

func handleAddAccount(ctx *TestkitRPCContext) http.HandlerFunc {
	type AddAccountResponse struct {
		AccountName string `json:"account_name"`
		Address     string `json:"address"`
		Mnemonic    string `json:"mnemonic"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if ctx.tg.IsSealed() {
			panic("can't add account after genesis seal")
		}

		if ctx.tg == nil {
			panic("testkit genesis not initialized")
		}

		vars := mux.Vars(r)
		accountName, ok := vars["accountName"]
		if !ok {
			panic("account name is not set")
		}

		i, seed := ctx.tg.AddAccount(accountName)

		response := AddAccountResponse{
			AccountName: accountName,
			Address:     i.GetAddress().String(),
			Mnemonic:    seed,
		}

		responseBody, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)

		return
	}
}

func handleCreateValidator(ctx *TestkitRPCContext) http.HandlerFunc {
	type CreateValidatorResponse struct {
		Msg              types.MsgCreateValidator `json:"Msg"`
		ValidatorAddress sdk.ValAddress           `json:"validator_address,string"`
		AccountName      string                   `json:"account_name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if ctx.tg == nil {
			panic("testkit genesis not initialized")
		}

		vars := mux.Vars(r)
		accountName, ok := vars["accountName"]
		if !ok {
			panic("account name is not set")
			return
		}

		msg := ctx.tg.CreateValidator(
			accountName,
			sdk.NewCoin("uluna", sdk.NewInt(1000000)),
			testkit.ZeroCommission,
		)

		response := CreateValidatorResponse{
			Msg:              msg,
			ValidatorAddress: msg.ValidatorAddress,
			AccountName:      accountName,
		}

		w.WriteHeader(http.StatusOK)
		w.Write(codec.MustMarshalJSON(response))

		return
	}
}

func handleSeal(ctx *TestkitRPCContext, router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ctx.tg.IsSealed() {
			panic("genesis already sealed")
		}

		if len(ctx.tg.GetValidators()) == 0 {
			panic("at least 1 validator must be present")
		}

		// create random shasum for this specific testkit context
		rand.Seed(time.Now().Unix())
		identifier := fmt.Sprintf("%d", rand.Intn(int(time.Now().Unix())))

		// create db
		db := leveldb.NewLevelDB(identifier)

		// create mantle
		mantle := app.NewMantle(
			db,
			ctx.tg.Seal(),
		)

		// setup for simulation
		mantle.SetBlockExecutor(mantlemint.NewSimBlockExecutor)
		mantle.Server(1337)

		// create local client for lcd works
		localClient := localclient.NewLocalClient(
			mantle.GetApp(),
			new(sync.Mutex),
		)
		appCtx := client.
			NewCLIContext().
			WithTrustNode(true).
			WithCodec(terra.MakeCodec()).
			WithClient(localClient)

		// register lcd routes
		terra.ModuleBasics.RegisterRESTRoutes(appCtx, router)

		// register auth routes
		authrest.RegisterRoutes(appCtx, router)

		// put all deps to ctx
		ctx.mantle = mantle
		ctx.tc = testkit.NewTestkitContext(
			mantle,
			db,
			ctx.tg.GetValidators(),
		)

		response := new(struct {
			Identifier string                   `json:"identifier"`
			Accounts   []testkit.TestkitAccount `json:"accounts"`
			Validators []sdk.ValAddress         `json:"validators"`
		})

		response.Identifier = identifier
		response.Accounts = ctx.tg.GetAccounts()
		response.Validators = ctx.tg.GetValidatorAddrs()

		w.WriteHeader(http.StatusOK)
		w.Write(codec.MustMarshalJSON(response))

		return
	}
}
