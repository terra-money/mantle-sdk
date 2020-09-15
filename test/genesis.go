package test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/terra-project/core/app"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/genaccounts"
	"time"
)

type (
	GenesisDoc     = tmtypes.GenesisDoc
	GenesisAccount = genaccounts.GenesisAccount
)

func NewGenesis(genesisAccounts ...GenesisAccount) *GenesisDoc {
	codec := app.MakeCodec()
	appStates := app.ModuleBasics.DefaultGenesis()
	appStates["accounts"] = codec.MustMarshalJSON(genesisAccounts)
	appStatesJson, err := codec.MarshalJSON(appStates)

	if err != nil {
		panic(fmt.Errorf("could not initialize appstate, %v", err))
	}

	gendoc := &GenesisDoc{
		ChainID:     "mantle-sim",
		Validators:  nil,
		AppState:    appStatesJson,
		GenesisTime: time.Now(),
	}

	if gendocErr := gendoc.ValidateAndComplete(); gendocErr != nil {
		panic(gendocErr)
	}

	return gendoc
}

var sequence uint64 = 0

func NewGenesisAccount(
	address AccAddress,
	coins Coins,
) GenesisAccount {
	config := sdk.GetConfig()
	config.SetCoinType(core.CoinType)
	config.SetFullFundraiserPath(core.FullFundraiserPath)
	config.SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)

	return genaccounts.NewGenesisAccountRaw(
		address,
		coins,
		nil,
		0,
		0,
		nil,
		"",
		"",
	)
}
