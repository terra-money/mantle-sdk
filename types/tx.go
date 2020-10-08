package types

import (
	"fmt"
	terra "github.com/terra-project/core/app"
)

//
var (
	codec = terra.MakeCodec()
)

func TxDecoder(txbytes []byte) (*StdTx, error) {
	var tx = StdTx{}

	if len(txbytes) == 0 {
		return nil, fmt.Errorf("Tx bytes are empty")
	}

	// StdTx.Msg is an interface. The concrete types
	// are registered by MakeTxCodec
	err := codec.UnmarshalBinaryLengthPrefixed(txbytes, &tx)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}
