package testkit

import (
	"github.com/tendermint/tendermint/state"
	tm "github.com/tendermint/tendermint/types"
	"github.com/terra-project/mantle-sdk/types"
	"time"
)

type ValidatorContext struct {
	validators     []tm.PrivValidator
	lastCommitSigs []tm.CommitSig
}

func NewValidatorContext(
	validators []tm.PrivValidator,
) *ValidatorContext {
	return &ValidatorContext{
		validators:     validators,
		lastCommitSigs: nil,
	}
}

func (vc *ValidatorContext) Propose(
	proposer tm.PrivValidator,
	state state.State,
	block *types.Block,
) *types.Block {
	pubkey, err := proposer.GetPubKey()
	if err != nil {
		panic(err)
	}

	// assign proposer address
	block.Header.ProposerAddress = pubkey.Address()

	// now
	then := block.Time.Add(1 * time.Second)

	// in testkit, there isn't really a voting process going on
	// since we have access to all privValidator keys,
	// we can forge any vote states and make it true

	// cdcEncode(h.Version),
	// 	cdcEncode(h.ChainID),
	// 	cdcEncode(h.Height),
	// 	cdcEncode(h.Time),
	// 	cdcEncode(h.LastBlockID),         <-----
	// 	cdcEncode(h.LastCommitHash),      <-----
	// 	cdcEncode(h.DataHash),            <-----
	// 	cdcEncode(h.ValidatorsHash),
	// 	cdcEncode(h.NextValidatorsHash),
	// 	cdcEncode(h.ConsensusHash),
	// 	cdcEncode(h.AppHash),             <----- d/c as mantle replaces them
	// 	cdcEncode(h.LastResultsHash),
	// 	cdcEncode(h.EvidenceHash),
	// 	cdcEncode(h.ProposerAddress),     <----- fabricated

	// fill last commit
	block.LastCommit = &tm.Commit{
		Height:     block.Header.Height - 1,
		Round:      1,
		BlockID:    state.LastBlockID,
		Signatures: vc.lastCommitSigs,
	}

	// make current block id
	currentBlockId := tm.BlockID{
		Hash:        block.Hash(),
		PartsHeader: block.MakePartSet(tm.BlockPartSizeBytes).Header(),
	}

	signatures := []tm.CommitSig{}
	for valIdx, _ := range state.Validators.Validators {
		vote, voteErr := makeVote(
			state.Copy(),
			block.Header.Height-1,
			currentBlockId,
			vc.validators[valIdx],
		)
		vote.Round = 1

		if voteErr != nil {
			panic(voteErr)
		}

		vote.Timestamp = then
		vote.Type = 0x20 // fix type to Commit

		signatures = append(signatures, vote.CommitSig())
	}

	// make signatures for them

	if block.Height > 1 {
		block.LastCommitHash = block.LastCommit.Hash()
		block.LastBlockID = state.LastBlockID
	}

	vc.lastCommitSigs = signatures

	return block
}

func makeVote(
	state state.State,
	height int64,
	blockId tm.BlockID,
	privVal tm.PrivValidator,
) (*tm.Vote, error) {
	return tm.MakeVote(
		height,
		blockId,
		state.Validators,
		privVal,
		state.ChainID,
		time.Now(),
	)
}
