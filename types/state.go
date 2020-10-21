package types

import (
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type (
	StdTx              = auth.StdTx
	ResponseDeliverTx  = abci.ResponseDeliverTx
	ResponseBeginBlock = abci.ResponseBeginBlock
	ResponseEndBlock   = abci.ResponseEndBlock
	BlockState         struct {
		Height             int64
		ResponseBeginBlock abci.ResponseBeginBlock
		ResponseEndBlock   abci.ResponseEndBlock
		ResponseDeliverTx  []ResponseDeliverTx
		Block              Block
	}
)
