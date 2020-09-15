package utils

import (
	tmtypes "github.com/tendermint/tendermint/types"
)

type GenesisDoc = tmtypes.GenesisDoc

func GenesisDocFromFile(path string) *tmtypes.GenesisDoc {
	genesis, err := tmtypes.GenesisDocFromFile(path)
	if err != nil {
		panic(err)
	}

	return genesis
}

func GenesisDocFromJSON(blob []byte) *tmtypes.GenesisDoc {
	genesis, err := tmtypes.GenesisDocFromJSON(blob)
	if err != nil {
		panic(err)
	}

	return genesis
}
