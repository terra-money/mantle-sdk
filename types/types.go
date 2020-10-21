package types

import (
	tmtypes "github.com/tendermint/tendermint/types"
)

// TODO: move this part to mantle-compatibility
type (
	Tx         = tmtypes.Tx
	GenesisDoc = tmtypes.GenesisDoc
	Block      = tmtypes.Block
	Header     = tmtypes.Header
)
