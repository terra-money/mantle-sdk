package test

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/terra-project/core/x/auth"
	authutilsterra "github.com/terra-project/core/x/auth/client/utils"
	compatlocalclient "github.com/terra-project/mantle-compatibility/localclient"
	"github.com/terra-project/mantle-sdk/app"
	"github.com/terra-project/mantle-sdk/types"
)

// dragons ahead
type (
	Tx = types.StdTx
)

type TestkitTx struct {
	msgs []Msg
	memo string
}

func NewTx() *TestkitTx {
	return &TestkitTx{}
}

func (tx *TestkitTx) WithMsg(msg Msg) *TestkitTx {
	tx.msgs = append(tx.msgs, msg)
	return tx
}

func (tx *TestkitTx) WithMemo(memo string) *TestkitTx {
	tx.memo = memo
	return tx
}

func (tx *TestkitTx) ToTx(signer TestAccount) Tx {
	signerAccount := getOrCreateAccount(signer.GetAddress())
	encoder := auth.DefaultTxEncoder(codec)
	txbldr := auth.NewTxBuilder(
		encoder,
		signerAccount.GetAccountNumber(),
		signerAccount.GetSequence(),
		0,
		0.0,
		false,
		"mantle-test",
		"",
		Coins{
			{Denom: "uusd", Amount: NewInt(100000000)},
		},
		nil,
	).WithKeybase(Keybase)

	emptyTx := auth.NewStdTx(
		tx.msgs,
		auth.StdFee{},
		nil,
		"",
	)
	emptyTx.Signatures = []auth.StdSignature{}

	localClient := compatlocalclient.NewLocalClient(app.GlobalTerraApp, nil)
	ctx := context.NewCLIContext().WithClient(localClient).WithCodec(codec)
	fees, gas, gaserr := authutilsterra.ComputeFeesWithStdTx(
		ctx,
		emptyTx,
		1.2,
		[]DecCoin{
			DecCoin{Denom: "uluna", Amount: types2.NewDecWithPrec(15, 6)},
		},
	)

	if gaserr != nil {
		panic(gaserr)
	}

	emptyTx.Fee = auth.StdFee{
		Amount: fees,
		Gas:    gas + 20000,
	}

	signedTx, signedTxErr := txbldr.SignStdTx(
		signer.info.GetName(),
		"default",
		emptyTx,
		false,
	)

	if signedTxErr != nil {
		panic(signedTxErr)
	}

	setAccountSequence(signerAccount)

	return signedTx
}

// must use rest clients, as store is not set in deliverState
// during this phase
func getOrCreateAccount(address AccAddress) exported.Account {
	querier := NewLocalQuerier()
	ar := auth.NewAccountRetriever(querier)

	acc, accErr := ar.GetAccount(address)
	// account was never set
	if acc == nil {
		panic("account was never set")
	} else {
		if GlobalAccountNumber < acc.GetAccountNumber() {
			GlobalAccountNumber = acc.GetAccountNumber()
		}
	}
	if accErr != nil {
		acc.SetAccountNumber(GlobalAccountNumber + 1)
	}

	return acc
}

func setAccountSequence(signerAccount exported.Account) {
	signerAccount.SetSequence(signerAccount.GetSequence() + 1)
}
