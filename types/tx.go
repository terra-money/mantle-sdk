package types

import (
	"encoding/base64"
	"errors"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	TerraApp "github.com/terra-project/core/app"
)

type (
	TxDecoder = func(string) (*Tx, error)
)

var (
	decoder = NewDecoder()
)

type LazyTx struct {
	TxString string
}

func NewLazyTx(txstring string) LazyTx {
	return LazyTx{
		TxString: txstring,
	}
}

func (lazyTx LazyTx) Decode() *types.StdTx {
	tx, decodeErr := decoder(lazyTx.TxString)
	if decodeErr != nil {
		panic(decodeErr)
	}

	return tx
}

func NewDecoder() TxDecoder {
	codec := TerraApp.MakeCodec()
	return func(encodedTx string) (*Tx, error) {
		txBytes, err := base64.StdEncoding.DecodeString(encodedTx)
		if err != nil {
			return nil, err
		}

		if len(txBytes) == 0 {
			return &Tx{}, nil
		}

		tx := Tx{}
		err = codec.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, errors.New("error decoding transaction")
		}

		return &tx, nil
	}
}
