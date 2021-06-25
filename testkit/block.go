package testkit

import (
	"github.com/tendermint/tendermint/state"
	"github.com/terra-money/mantle-sdk/types"
	"time"
)

type Block struct {
	nextBlock *types.Block
	lastState *state.State
}

func NewBlock(lastState state.State) *Block {
	var blockTime time.Time
	var height int64

	// if genesis
	if lastState.LastBlockHeight == 0 {
		blockTime = lastState.LastBlockTime
		height = 1
	} else {
		blockTime = lastState.LastBlockTime.Add(6 * time.Second)
		height = lastState.LastBlockHeight + 1
	}

	nextBlock := types.Block{}
	nextBlock.Header = types.Header{
		Version:            lastState.Version.Consensus,
		ChainID:            lastState.ChainID,
		Height:             height,
		Time:               blockTime,
		LastBlockID:        lastState.LastBlockID,
		LastCommitHash:     lastState.LastBlockID.Hash,
		ValidatorsHash:     lastState.Validators.Hash(),
		NextValidatorsHash: lastState.NextValidators.Hash(),
		ConsensusHash:      lastState.ConsensusParams.Hash(),
		AppHash:            lastState.AppHash,
		LastResultsHash:    lastState.LastResultsHash,
		ProposerAddress:    nil,
		EvidenceHash:       nil, // can't simulate double signing
		DataHash:           nil, // ??
	}

	return &Block{
		nextBlock: &nextBlock,
		lastState: &lastState,
	}
}

func (block *Block) WithTx(tx types.StdTx) *Block {
	block.nextBlock.Data.Txs = append(block.nextBlock.Data.Txs, codec.MustMarshalBinaryLengthPrefixed(tx))
	return block
}

func (block *Block) WithTime(t time.Time) *Block {
	block.nextBlock.Header.Time = t
	return block
}

// finalize block header
func (block *Block) Finalize() *types.Block {
	block.nextBlock.DataHash = block.nextBlock.Data.Hash()
	return block.nextBlock
}
