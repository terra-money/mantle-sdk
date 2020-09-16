package test

import (
	"github.com/terra-project/mantle-compatibility/genesis"
)

func NewGenesis(genesisAccounts ...GenesisAccount) *GenesisDoc {
	return genesis.NewGenesis(genesisAccounts...)
}

var sequence uint64 = 0

func NewGenesisAccount(
	address AccAddress,
	coins Coins,
) GenesisAccount {
	return genesis.NewGenesisAccount(address, coins)
}
