package subscriber

import (
	types "github.com/terra-project/mantle/types"
)

type Subscriber interface {
	Subscribe() chan types.Block
	Close() error
}

type BlockEvent struct {
	Query string `json:"query"`
	Data  struct {
		Type  string `json:"type"`
		Value struct {
			Block            types.Block            `json:"block"`
			ResultBeginBlock types.ResultBeginBlock `json:"result_begin_block"`
			ResultEndBlock   types.ResultEndBlock   `json:"result_end_block"`
		} `json:"value"`
	} `json:"data"`
	Events types.Events `json:"events"`
}
