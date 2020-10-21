// test_test lol
package test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/mantle-sdk/app"
	"github.com/terra-project/mantle-sdk/db/badger"
	"testing"
)

func TestCreateBlock(t *testing.T) {

	func() {
		// create 2 accounts
		acc1 := NewAccount()
		acc2 := NewAccount()

		mantle := app.NewMantle(
			badger.NewBadgerDB(""), // memdb
			NewGenesis(
				NewGenesisAccount(acc1, Coins{
					{Denom: "uluna", Amount: sdk.NewInt(100000000000)},
				}),
			),
			nil, // for testing purposes we don't run any indexers
		)

		// create some sim blocks
		blocks := []*Block{
			NewBlock().WithTx(
				NewTx().WithMsg(
					NewMsgSend(
						acc1,
						acc2,
						Coins{
							{Denom: "uluna", Amount: sdk.NewInt(2000000)},
						},
					),
				).ToTx(),
			).ToBlock(),
			NewBlock().WithTx(
				NewTx().WithMsg(
					NewMsgSend(
						acc1,
						acc2,
						Coins{
							{Denom: "uluna", Amount: sdk.NewInt(2000000)},
						},
					),
				).ToTx(),
			).ToBlock(),
			NewBlock().WithTx(
				NewTx().WithMsg(
					NewMsgSend(
						acc1,
						acc2,
						Coins{
							{Denom: "uluna", Amount: sdk.NewInt(2000000)},
						},
					),
				).ToTx(),
			).ToBlock(),
		}

		for _, block := range blocks {
			mantle.Inject(&block)
		}

	}()
}
