package test

import (
	"encoding/base64"
	"github.com/terra-project/core/app"
	"github.com/terra-project/mantle/types"
	"time"
)

var (
	codec               = app.MakeCodec()
	blockSequence int64 = 2
)

type (
	Block        = types.Block
	TestkitBlock struct {
		block Block
		txs   []string
	}
)

func NewBlock() *TestkitBlock {
	return &TestkitBlock{}
}

func (block *TestkitBlock) WithTx(tx Tx) *TestkitBlock {
	block.txs = append(block.txs, base64.StdEncoding.EncodeToString(codec.MustMarshalBinaryLengthPrefixed(tx)))
	return block
}

func (block *TestkitBlock) ToBlock() *Block {
	block.block.Header.Time = time.Now().Format(time.RFC3339)
	block.block.Header.Height = blockSequence
	block.block.Data.Txs = block.txs

	blockSequence++

	return &block.block
}
