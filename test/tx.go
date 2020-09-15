package test

import (
	"github.com/terra-project/core/x/auth"
	"github.com/terra-project/mantle/types"
)

type (
	Tx = types.Tx
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

func (tx *TestkitTx) ToTx() Tx {
	return auth.NewStdTx(
		tx.msgs,
		auth.NewStdFee(
			0,
			Coins{},
		),
		nil,
		tx.memo,
	)
}
