package utils

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/terra-project/mantle-sdk/types"
)

// ConvertBlockHeaderToTMHeader is to convert BlockHeader to its correct format
func ConvertBlockHeaderToTMHeader(block json.RawMessage) tmtypes.Block {
	bb := tmtypes.Block{}
	if err := codec.Cdc.UnmarshalJSON(block, &bb); err != nil {
		panic(err)
	}

	return bb

}

func ConvertBlockToRawBlock(block *types.Block) types.RawBlock {
	return types.RawBlock{
		Header:     block.Header,
		Data:       block.Data,
		Evidence:   block.Evidence,
		LastCommit: *block.LastCommit,
	}
}
