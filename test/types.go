package test

import (
	compat "github.com/terra-project/mantle-compatibility/types"
)

type (
	Coin           = compat.Coin
	Coins          = compat.Coins
	Int            = compat.Int
	Dec            = compat.Dec
	DecCoin        = compat.DecCoin
	AccAddress     = compat.AccAddress
	ValAddress     = compat.ValAddress
	Msg            = compat.Msg
	GenesisDoc     = compat.GenesisDoc
	GenesisAccount = compat.GenesisAccount
	HexBytes       = compat.HexBytes
)

var (
	NewInt = compat.NewInt
)
