package testkit

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/terra-project/core/x/bank"
	"github.com/terra-money/mantle-sdk/app"
	"github.com/terra-money/mantle-sdk/app/mantlemint"
	"github.com/terra-money/mantle-sdk/db/leveldb"
	"github.com/terra-money/mantle-sdk/types"
	"testing"
)

func TestTeskit(t *testing.T) {
	tg := NewTestkitGenesis("testkit")
	a1, _ := tg.AddAccount("test1")
	a2, _ := tg.AddAccount("test2")

	// create validator
	tg.CreateValidator(
		"test1",
		sdk.NewCoin("uluna", sdk.NewInt(1000000)),
		ZeroCommission,
	)

	fmt.Println(a1.GetAddress().String())
	fmt.Println(a2.GetAddress().String())

	// finalize genesis
	gendoc := tg.Seal()

	// create mantle
	db := leveldb.NewLevelDB("test")
	mantle := app.NewMantle(db, gendoc)

	mantle.SetBlockExecutor(mantlemint.NewSimBlockExecutor(db.GetCosmosAdapter(), mantle.GetApp()))

	//
	tk := NewTestkitContext(
		mantle,
		db,
		tg.GetValidators(),
	)

	// inject genesis block
	var blockState *types.BlockState
	var err error
	tk.Inject(tk.PickProposerByIndex(0))

	// add to testkit mempool
	tk.AddToMempool(NewSignedTx(
		[]sdk.Msg{bank.NewMsgSend(
			a1.GetAddress(),
			a2.GetAddress(),
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(10000000))),
		)},
		tg.GetKeybase(),
		"test1",
		"testkit",
		mantle.GetApp(),
	))

	//
	blockState, err = tk.Inject(tk.PickProposerByIndex(0))
	if err != nil {
		panic(err)
	}
	fmt.Println(blockState)

	// add to testkit mempool
	tk.AddToMempool(NewSignedTx(
		[]sdk.Msg{bank.NewMsgSend(
			a1.GetAddress(),
			a2.GetAddress(),
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(10000000))),
		)},
		tg.GetKeybase(),
		"test1",
		"testkit",
		mantle.GetApp(),
	))

	//
	blockState, err = tk.Inject(tk.PickProposerByIndex(0))
	if err != nil {
		panic(err)
	}
	fmt.Println(blockState)

	// add to testkit mempool
	tk.AddToMempool(NewSignedTx(
		[]sdk.Msg{bank.NewMsgSend(
			a1.GetAddress(),
			a2.GetAddress(),
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(10000000))),
		)},
		tg.GetKeybase(),
		"test1",
		"testkit",
		mantle.GetApp(),
	))

	//
	blockState, err = tk.Inject(tk.PickProposerByIndex(0))
	if err != nil {
		panic(err)
	}
	fmt.Println(blockState)

	// add to testkit mempool
	tk.AddToMempool(NewSignedTx(
		[]sdk.Msg{bank.NewMsgSend(
			a1.GetAddress(),
			a2.GetAddress(),
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(10000000))),
		)},
		tg.GetKeybase(),
		"test1",
		"testkit",
		mantle.GetApp(),
	))

	//
	blockState, err = tk.Inject(tk.PickProposerByIndex(0))
	if err != nil {
		panic(err)
	}
	fmt.Println(blockState)

	// add to testkit mempool
	tk.AddToMempool(NewSignedTx(
		[]sdk.Msg{bank.NewMsgSend(
			a1.GetAddress(),
			a2.GetAddress(),
			sdk.NewCoins(sdk.NewCoin("uluna", sdk.NewInt(10000000))),
		)},
		tg.GetKeybase(),
		"test1",
		"testkit",
		mantle.GetApp(),
	))

	//
	blockState, err = tk.Inject(tk.PickProposerByIndex(0))
	if err != nil {
		panic(err)
	}
	fmt.Println(blockState)

}
