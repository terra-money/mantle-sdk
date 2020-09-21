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

func (block *TestkitBlock) WithHeight(height int64) *TestkitBlock {
	block.block.Header.Height = int64(height)
	return block
}

func (block *TestkitBlock) WithTime(t time.Time) *TestkitBlock {
	block.block.Header.Time = t.Format(time.RFC3339)
	return block
}

func (block *TestkitBlock) ToBlock() *Block {
	block.block.Header.ChainID = "tequila-0004"
	if block.block.Header.Time == "" {
		block.block.Header.Time = time.Now().Format(time.RFC3339)
	}

	if block.block.Header.Height == 0 {
		block.block.Header.Height = blockSequence
		blockSequence++
	}

	block.block.Data.Txs = block.txs
	return &block.block
}
