package rpc

import (
	"encoding/json"
	"fmt"
	client "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gorilla/mux"
	terra "github.com/terra-project/core/app"
	"github.com/terra-project/core/x/auth"
	authrest "github.com/terra-project/core/x/auth/client/rest"
	compatlocalclient "github.com/terra-project/mantle-compatibility/localclient"
	"github.com/terra-project/mantle-sdk/app/mantlemint"
	"github.com/terra-project/mantle-sdk/db/leveldb"
	"github.com/terra-project/mantle-sdk/testkit"
	"net/http"
	"sync"
)

type AccountsInitRequest struct {
	AccountName string `json:"account_name"`
	Mnemonic    string `json:"mnemonic"`
}

type ValidatorInitRequest struct {
	AccountName    string                       `json:"account_name"`
	SelfDelegation sdk.Coin                     `json:"self_delegation"`
	Commission     stakingtypes.CommissionRates `json:"commission"`
}

type AutomaticTxRequest struct {
	AccountName string      `json:"account_name"`
	Period      int         `json:"period"`
	Msgs        []sdk.Msg   `json:"msgs"`
	Fee         auth.StdFee `json:"fee"`
	Offset      int64       `json:"offset"`
	StartAt     int64       `json:"start_at"`
}

type AutomaticTxPauseRequest struct {
	AccountName string `json:"account_name"`
}

type AutomaticInjectRequest struct {
	ValidatorRounds []string `json:"validator_rounds"`
}

type Request struct {
	Genesis            json.RawMessage
	Accounts           []AccountsInitRequest   `json:"accounts"`
	Validators         []ValidatorInitRequest  `json:"validators"`
	AutomaticTxRequest []AutomaticTxRequest    `json:"auto_txs"`
	AutomaticInject    *AutomaticInjectRequest `json:"auto_inject"`
}

type Response struct {
	Identifier      string                     `json:"identifier"`
	ChainID         string                     `json:"chain_id"`
	Genesis         string                     `json:"genesis"`
	Accounts        []AddAccountResponse       `json:"accounts"`
	Validators      []CreateValidatorResponse  `json:"validators"`
	AutomaticTx     []AutomaticTxEntryResponse `json:"auto_txs"`
	AutomaticInject *AutomaticInjectResponse   `json:"auto_inject"`
}

func handleInitTestkit(ctx *TestkitRPCContext, router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := Request{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&request)
		if err != nil {
			panic(err)
		}

		response := initTestkit(ctx, request, router)

		w.WriteHeader(200)
		w.Write(MustMarshalJSON(response))
	}
}

func initTestkit(ctx *TestkitRPCContext, req Request, router *mux.Router) Response {
	// create genesis from file
	tg := testkit.NewGenesisFromFile(req.Genesis)

	// import all accounts
	accountResponse := make([]AddAccountResponse, len(req.Accounts))
	for i, accountRequest := range req.Accounts {
		info, seed := tg.ImportAccount(accountRequest.AccountName, accountRequest.Mnemonic)
		accountResponse[i] = AddAccountResponse{
			AccountName: accountRequest.AccountName,
			Address:     info.GetAddress().String(),
			Mnemonic:    seed,
		}
	}

	// create validators using accounts
	validatorResponse := make([]CreateValidatorResponse, len(req.Validators))
	for i, validatorRequest := range req.Validators {
		msg := tg.CreateValidator(
			validatorRequest.AccountName,
			validatorRequest.SelfDelegation,
			validatorRequest.Commission,
		)
		validatorResponse[i] = CreateValidatorResponse{
			Msg:              msg,
			ValidatorAddress: msg.ValidatorAddress,
			AccountName:      validatorRequest.AccountName,
		}
	}

	// seal genesis
	gendoc := tg.Seal()

	// create identifier for this network
	identifier := GenerateTestkitIdentifier(gendoc.ChainID)

	// create testkit context
	db := leveldb.NewLevelDB(identifier)
	tctx := testkit.NewTestkitContext(tg, db)

	// set automatic tx request
	autoTxResponse := make([]AutomaticTxEntryResponse, len(req.AutomaticTxRequest))
	for i, autoTxRequest := range req.AutomaticTxRequest {
		entry := testkit.NewAutomaticTxEntry(
			autoTxRequest.AccountName,
			autoTxRequest.Fee,
			autoTxRequest.Msgs,
			autoTxRequest.Period,
			autoTxRequest.Offset,
			autoTxRequest.StartAt,
		)

		tctx.AddAutomaticTxEntry(entry)

		autoTxResponse[i] = AutomaticTxEntryResponse{
			Msgs:        entry.Msgs,
			Fee:         entry.Fee,
			Period:      entry.Period,
			AccountName: entry.AccountName,
		}
	}

	// set automatic inject request if present
	var autoInjectResponse *AutomaticInjectResponse = nil
	if req.AutomaticInject != nil {
		tctx.SetAutomaticInjection(req.AutomaticInject.ValidatorRounds)
		autoInjectResponse = new(AutomaticInjectResponse)
		autoInjectResponse.ValidatorRounds = req.AutomaticInject.ValidatorRounds
	}

	// set testkit in rpc context
	ctx.SetTestkitContext(identifier, tctx)

	// set router
	terraApp := tctx.GetMantle().GetApp()
	globalQueryMutex := new(sync.Mutex)
	localClient := compatlocalclient.NewLocalClient(terraApp, globalQueryMutex)
	cliCtx := client.
		NewCLIContext().
		WithTrustNode(true).
		WithCodec(terra.MakeCodec()).
		WithClient(localClient)

	subrouter := router.PathPrefix(fmt.Sprintf("/%s", identifier)).Subrouter()
	terra.ModuleBasics.RegisterRESTRoutes(cliCtx, subrouter)
	authrest.RegisterRoutes(cliCtx, subrouter)

	gendocJson, _ := codec.MarshalJSON(gendoc)

	// replace mantlemint executer
	tctx.GetMantle().SetBlockExecutor(mantlemint.NewSimBlockExecutor)

	// return response
	return Response{
		Identifier:      identifier,
		ChainID:         gendoc.ChainID,
		Genesis:         string(gendocJson),
		Accounts:        accountResponse,
		Validators:      validatorResponse,
		AutomaticTx:     autoTxResponse,
		AutomaticInject: autoInjectResponse,
	}
}
