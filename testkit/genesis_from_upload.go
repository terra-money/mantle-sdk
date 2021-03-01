package testkit

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/types"
	core "github.com/terra-project/core/types"
)

func NewGenesisFromFile(genesisBlob json.RawMessage) *TestkitGenesis {
	config := sdk.GetConfig()
	config.SetCoinType(core.CoinType)
	config.SetFullFundraiserPath(core.FullFundraiserPath)
	config.SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(core.Bech32PrefixConsAddr, core.Bech32PrefixConsPub)

	gendoc, err := types.GenesisDocFromJSON(genesisBlob)
	if err != nil {
		panic(err)
	}

	// create testkit genesis, empty for now
	tg := NewTestkitGenesis(gendoc.ChainID)
	tg.genesis = gendoc // use default genesis

	return tg
}
