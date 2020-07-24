package utils

import (
	"encoding/base64"
	"errors"

	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	TerraApp "github.com/terra-project/core/app"
)

type (
	Tx        = auth.StdTx
	TxDecoder = func(string) (*Tx, error)
)

func NewDecoder() TxDecoder {
	codec := TerraApp.MakeCodec()
	return func(encodedTx string) (*Tx, error) {
		txBytes, err := base64.StdEncoding.DecodeString(encodedTx)
		if err != nil {
			return nil, err
		}

		if len(txBytes) == 0 {
			return nil, errors.New("txBytes are empty")
		}

		tx := Tx{}
		err = codec.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, errors.New("error decoding transaction")
		}

		return &tx, nil
	}
}
