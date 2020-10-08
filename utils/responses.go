package utils

import (
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func ConvertToABCIHeader(header *tmtypes.Header) abci.Header {
	var abciHeader = abci.Header{
		Version: abci.Version{
			Block: header.Version.Block.Uint64(),
			App:   header.Version.App.Uint64(),
		},
		ChainID: header.ChainID,
		Height:  header.Height,
		Time:    header.Time,
		LastBlockId: abci.BlockID{
			Hash: header.LastBlockID.Hash,
			PartsHeader: abci.PartSetHeader{
				Total: int32(header.LastBlockID.PartsHeader.Total),
				Hash:  header.LastBlockID.Hash,
			},
		},
		LastCommitHash:     header.LastCommitHash,
		DataHash:           header.DataHash,
		ValidatorsHash:     header.ValidatorsHash,
		NextValidatorsHash: header.NextValidatorsHash,
		ConsensusHash:      header.ConsensusHash,
		AppHash:            header.AppHash,
		LastResultsHash:    header.LastResultsHash,
		EvidenceHash:       header.EvidenceHash,
		ProposerAddress:    header.ProposerAddress,
	}

	return abciHeader
}
