package test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
	"math/rand"
)

// super simple acc creation
// since mantle doesn't care about signatures,
// accounts just have to match in its length (AddrLen = 20bytes)
func NewAccount() AccAddress {
	config := sdk.GetConfig()
	config.SetCoinType(core.CoinType)
	config.SetFullFundraiserPath(core.FullFundraiserPath)
	config.SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)

	acc := make(AccAddress, 20)
	_, err := rand.Read(acc)
	if err != nil {
		panic(err)
	}

	return acc
}

func AccountFromBech32(addr string) AccAddress {
	config := sdk.GetConfig()
	config.SetCoinType(core.CoinType)
	config.SetFullFundraiserPath(core.FullFundraiserPath)
	config.SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)
	account, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		panic(err)
	}
	return account
}
