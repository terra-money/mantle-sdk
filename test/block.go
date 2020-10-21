package test

import (
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/terra-project/core/app"
	"github.com/terra-project/mantle-sdk/types"
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
		txs   []tmtypes.Tx
	}
)

func NewBlock() *TestkitBlock {
	return &TestkitBlock{}
}

func (block *TestkitBlock) WithTx(tx Tx) *TestkitBlock {
	block.txs = append(block.txs, codec.MustMarshalBinaryLengthPrefixed(tx))
	return block
}

func (block *TestkitBlock) WithHeight(height int64) *TestkitBlock {
	block.block.Header.Height = height
	return block
}

func (block *TestkitBlock) WithTime(t time.Time) *TestkitBlock {
	block.block.Header.Time = t
	return block
}

func (block *TestkitBlock) ToBlock() *Block {
	block.block.Header.ChainID = "mantle-test"
	if block.block.Header.Time.IsZero() {
		block.block.Header.Time = time.Now()
	}

	if block.block.Header.Height == 0 {
		block.block.Header.Height = blockSequence
		blockSequence++
	}

	block.block.Data.Txs = block.txs

	return &block.block
}
