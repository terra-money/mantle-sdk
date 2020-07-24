package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	types "github.com/terra-project/mantle/types"
)

func UnmarshalBlockResponseFromLCD(blockResponse []byte, target *types.Block) {
	tmp := make(map[string]json.RawMessage)

	if err := json.Unmarshal(blockResponse, &tmp); err != nil {
		panic(fmt.Sprintf("Error while converting block: %s", err))
	}

	if err := json.Unmarshal([]byte(tmp["block"]), &target); err != nil {
		panic("Error during UnmarshalBlockResponseFromLCD")
	}
}

func ConvertToABCIHeader(header *types.Header) abci.Header {
	var abciHeader = abci.Header{}

	t, err := strconv.Atoi(header.Time)

	if err != nil {
		fmt.Println(header)
		panic("Error in time")
	}

	abciHeader.AppHash = []byte(header.AppHash)
	abciHeader.ChainID = header.ChainID
	abciHeader.ConsensusHash = []byte(header.ConsensusHash)
	abciHeader.DataHash = []byte(header.DataHash)
	abciHeader.EvidenceHash = []byte(header.EvidenceHash)
	abciHeader.Height = header.Height
	abciHeader.LastBlockId.Hash = []byte(header.LastBlockID.Hash)
	abciHeader.LastBlockId.PartsHeader.Hash = []byte(header.LastBlockID.Parts.Hash)
	abciHeader.LastBlockId.PartsHeader.Total = strToInt32(header.LastBlockID.Parts.Total)
	abciHeader.LastCommitHash = []byte(header.LastCommitHash)
	abciHeader.LastResultsHash = []byte(header.LastResultsHash)
	abciHeader.NextValidatorsHash = []byte(header.NextValidatorsHash)
	// abciHeader.NumTxs = header.NumTxs // removed as of tendermint 0.33.0
	abciHeader.ProposerAddress = []byte(header.ProposerAddress)
	abciHeader.Time = time.Unix(0, int64(t))
	// abciHeader.TotalTxs = header.TotalTxs // remove as of tendermint 0.33.0
	abciHeader.ValidatorsHash = []byte(header.ValidatorsHash)
	abciHeader.Version.App = header.Version.App
	abciHeader.Version.Block = header.Version.Block

	return abciHeader
}

func strToInt32(str string) int32 {
	data, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		panic(err)
	}

	return int32(data)
}
