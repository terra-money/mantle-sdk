package testkit

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	TerraApp "github.com/terra-project/core/app"
	"github.com/terra-project/core/x/auth"
)

var localQuerier auth.NodeQuerier

// only use this for fixed testing, there is no fee estimation going on
func NewSignedTx(
	msgs []sdk.Msg,
	keybase keys.Keybase,
	signerAccountName string,
	chainId string,
	terraApp *TerraApp.TerraApp,
) auth.StdTx {
	signerKey, err := keybase.Get(signerAccountName)
	if err != nil {
		panic(err)
	}

	account := getOrCreateAccount(signerKey.GetAddress(), terraApp)

	txbldr := auth.NewTxBuilder(
		auth.DefaultTxEncoder(codec),
		account.GetAccountNumber(),
		account.GetSequence(),
		10000000,
		1.4,
		false,
		chainId,
		"",
		sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(1000000))),
		nil,
	).WithKeybase(keybase)

	signedTx, err := txbldr.SignStdTx(
		signerAccountName,
		defaultPassphrase,
		auth.NewStdTx(
			msgs,
			auth.NewStdFee(10000000, sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(1000000)))),
			nil,
			"",
		),
		true,
	)

	if err != nil {
		panic(err)
	}

	return signedTx
}

func getOrCreateAccount(
	address sdk.AccAddress,
	terraApp *TerraApp.TerraApp,
) exported.Account {
	// make this on the first request
	if localQuerier == nil {
		localQuerier = NewLocalQuerier(terraApp)
	}
	ar := auth.NewAccountRetriever(localQuerier)

	// check if account is known
	if accountExists := ar.EnsureExists(address); accountExists != nil {
		panic(accountExists)
	}

	acc, accErr := ar.GetAccount(address)
	if accErr != nil {
		panic(accErr)
	}

	return acc
}
