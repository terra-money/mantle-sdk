package rpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	terra "github.com/terra-project/core/app"
	"github.com/terra-project/mantle-compatibility/types"
)

var codec = terra.MakeCodec()

type TestkitCreatorFunc func(r *mux.Router, ctx *TestkitRPCContext)

// AddAccountResponse is RPC response type for genesis AddAccount
type AddAccountResponse struct {
	AccountName string `json:"account_name"`
	Address     string `json:"address"`
	Mnemonic    string `json:"mnemonic"`
}

// CreateValidatorResponse is RPC response type for genesis CreateValidator
type CreateValidatorResponse struct {
	Msg              types.MsgCreateValidator `json:"Msg"`
	ValidatorAddress sdk.ValAddress           `json:"validator_address,string"`
	AccountName      string                   `json:"account_name"`
}

type AutomaticTxEntryResponse AutomaticTxRequest
type AutomaticInjectResponse AutomaticInjectRequest

type TxResultLog struct {
	MsgIdx int64       `json:"msg_index,string"`
	Log    string      `json:"log"`
	Events []sdk.Event `json:"events"`
}
