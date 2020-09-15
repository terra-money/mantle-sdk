package app

//
// import (
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	core "github.com/terra-project/core/types"
// 	"github.com/terra-project/mantle/db/badger"
// 	"github.com/terra-project/mantle/sim"
// 	"github.com/terra-project/mantle/types"
// 	"testing"
// )
//
// func initConfig() {
// 	config := sdk.GetConfig()
// 	config.SetCoinType(core.CoinType)
// 	config.SetFullFundraiserPath(core.FullFundraiserPath)
// 	config.SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
// 	config.SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)
// 	config.SetBech32PrefixForConsensusNode(core.Bech32PrefixConsAddr, core.Bech32PrefixConsPub)
// 	// don't seal
// 	// config.Seal()
// }
// func TestApp(t *testing.T) {
// 	initConfig()
//
// 	genesis := sim.CreateGenesis([]sim.GenesisAccount{
// 		sim.CreateGenesisAccount(
// 			"terra1a6nw09knat6sqxnz0xf0pt9mks35jfy0gpgwcy",
// 			sdk.Coins{
// 				{
// 					Denom:  "uluna",
// 					Amount: sdk.NewInt(10000000),
// 				},
// 			},
// 		),
// 	})
//
// 	app := NewMantle(
// 		badger.NewBadgerDB(""),
// 		false,
// 		genesis,
// 		[]types.IndexerRegisterer{},
// 	)
//
// 	app.Rebuild()
//
// }
