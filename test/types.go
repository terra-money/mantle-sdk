package test

import sdk "github.com/cosmos/cosmos-sdk/types"

type (
	Coin       = sdk.Coin
	Coins      = sdk.Coins
	Int        = sdk.Int
	Dec        = sdk.Dec
	DecCoin    = sdk.DecCoin
	AccAddress = sdk.AccAddress
	ValAddress = sdk.ValAddress
	Msg        = sdk.Msg
)

var (
	NewInt = sdk.NewInt
)
