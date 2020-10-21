package types

import (
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type (
	StdTx             = auth.StdTx
	ResponseDeliverTx = abci.ResponseDeliverTx
	TxResult          struct {
		Result abci.ResponseDeliverTx
		Tx     tmtypes.Tx
	}
	BlockState struct {
		Height             int64
		ResponseBeginBlock abci.ResponseBeginBlock
		ResponseEndBlock   abci.ResponseEndBlock
		ResponseDeliverTx  []ResponseDeliverTx
		Block              Block
	}
)
