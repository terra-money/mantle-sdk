package testkit

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/terra-project/core/x/auth"
)

var localQuerier auth.NodeQuerier

// only use this for fixed testing, there is no fee estimation going on
func NewTxBldr(
	fee auth.StdFee,
	keybase keys.Keybase,
	chainId string,
	acc exported.Account,
) auth.TxBuilder {
	txbldr := auth.NewTxBuilder(
		auth.DefaultTxEncoder(codec),
		acc.GetAccountNumber(),
		acc.GetSequence(),
		10000000,
		1.4,
		false,
		chainId,
		"",
		fee.Amount,
		nil,
	).WithKeybase(keybase)

	return txbldr
}

func SignTxBldr(txbldr auth.TxBuilder, signerAccountName string, msgs []sdk.Msg, fee auth.StdFee) auth.StdTx {
	signedTx, err := txbldr.SignStdTx(
		signerAccountName,
		defaultPassphrase,
		auth.NewStdTx(
			msgs,
			fee,
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
