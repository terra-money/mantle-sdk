package types

import (
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type (
	StdTx = auth.StdTx
)

// State houses all primitive data
type BaseState struct {
	Height             int64
	BeginBlockResponse abci.ResponseBeginBlock
	EndBlockResponse   abci.ResponseEndBlock
	DeliverTxResponses []abci.ResponseDeliverTx
	Block              Block
	Txs                []StdTx
}
