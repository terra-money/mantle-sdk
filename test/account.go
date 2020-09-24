package test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var GlobalAccountNumber uint64 = 0

type TestAccount struct {
	name string
	info keys.Info
}

func (ta TestAccount) GetAddress() AccAddress {
	return AccAddressFromBech32(ta.info.GetAddress().String())
}

func NewAccount(name string) TestAccount {
	i, _, e := Keybase.CreateMnemonic(name, keys.English, "default", keys.Secp256k1)
	if e != nil {
		panic(e)
	}

	return TestAccount{
		info: i,
		name: name,
	}
}

func ImportAccount(name, mnemonic string) TestAccount {
	i, createAccErr := Keybase.CreateAccount(
		name,
		mnemonic,
		keys.DefaultBIP39Passphrase,
		"default",
		keys.CreateHDPath(0, 0).String(),
		keys.Secp256k1,
	)

	if createAccErr != nil {
		panic(createAccErr)
	}

	return TestAccount{
		info: i,
		name: name,
	}
}

func AccAddressFromBech32(addr string) AccAddress {
	account, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		panic(err)
	}
	return account
}
