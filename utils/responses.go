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

	if err := json.Unmarshal(tmp["block"], &target); err != nil {
		panic("Error during UnmarshalBlockResponseFromLCD")
	}
}

func ConvertToABCIHeader(header *types.Header) abci.Header {
	var abciHeader = abci.Header{}

	time, err := time.Parse(time.RFC3339, header.Time)

	if err != nil {
		panic(err)
	}

	abciHeader.AppHash = []byte(header.AppHash)
	abciHeader.ChainID = header.ChainID
	abciHeader.Height = header.Height
	abciHeader.Time = time

	return abciHeader
}

func strToInt32(str string) int32 {
	data, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		panic(err)
	}

	return int32(data)
}
