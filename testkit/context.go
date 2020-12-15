package testkit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tm "github.com/tendermint/tendermint/types"
	"github.com/terra-project/mantle-sdk/app"
	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/types"
	"sync"
)

type TestkitContext struct {
	m          *sync.RWMutex
	validators []TestkitGenesisAccountToPrivValMap
	vc         *ValidatorContext
	mantle     *app.Mantle
	mempool    []types.StdTx
	db         db.DB
}

func NewTestkitContext(
	mantle *app.Mantle,
	db db.DB,
	validators []TestkitGenesisAccountToPrivValMap,
) *TestkitContext {

	var valpvs = make([]tm.PrivValidator, len(validators))
	for i, vm := range validators {
		valpvs[i] = vm.PrivValidator
	}

	return &TestkitContext{
		mantle:     mantle,
		mempool:    []types.StdTx{},
		m:          new(sync.RWMutex),
		validators: validators,
		vc:         NewValidatorContext(valpvs),
		db:         db,
	}
}

func (ctx *TestkitContext) ClearMempool() {
	ctx.m.Lock()
	ctx.mempool = []types.StdTx{}
	ctx.m.Unlock()
}

func (ctx *TestkitContext) AddToMempool(tx types.StdTx) {
	ctx.m.Lock()
	ctx.mempool = append(ctx.mempool, tx)
	ctx.m.Unlock()
}

func (ctx *TestkitContext) Inject(proposer tm.PrivValidator) (*types.BlockState, error) {
	lastState := ctx.mantle.GetLastState()

	// create block according to the last block
	nextBlock := NewBlock(lastState)

	// set all txs
	for _, tx := range ctx.mempool {
		nextBlock = nextBlock.WithTx(tx)
	}

	// prep block for injection
	blockToInject := nextBlock.Finalize()

	ctx.db.SetCriticalZone()

	// propose block
	proposedBlock := ctx.vc.Propose(proposer, ctx.mantle.GetLastState(), blockToInject)

	// let mantle inject; return blockState
	blockState, err := ctx.mantle.Inject(proposedBlock)

	ctx.ClearMempool()

	// force flush all batch, so state queries can be made
	ctx.db.ReleaseCriticalZone()

	return blockState, err
}

func (tc *TestkitContext) PickProposerByIndex(index int) tm.PrivValidator {
	if index > len(tc.validators) {
		panic("not enough validators")
	}
	return tc.validators[index].PrivValidator
}

func (tc *TestkitContext) PickProposerByAddress(address sdk.ValAddress) tm.PrivValidator {
	for _, val := range tc.validators {
		if val.Account.Equals(address) {
			return val.PrivValidator
		}
	}

	panic("validator not found")
}
