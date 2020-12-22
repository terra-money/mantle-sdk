package testkit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tm "github.com/tendermint/tendermint/types"
	"github.com/terra-project/core/x/auth"
	"github.com/terra-project/mantle-sdk/app"
	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/types"
	"log"
	"sync"
)

type TestkitContext struct {
	ID         string
	tg         *TestkitGenesis
	m          *sync.RWMutex
	validators []TestkitGenesisAccountToPrivValMap
	vc         *ValidatorContext
	mantle     *app.Mantle
	mempool    []types.StdTx
	db         db.DB

	autoTxs       []AutomaticTxEntry
	autoInjection *AutomaticInjection
}

func NewTestkitContext(
	tg *TestkitGenesis,
	db db.DB,
) *TestkitContext {
	if !tg.IsSealed() {
		panic("cannot create testkit context using unsealed genesis")
	}

	// create mantle from
	mantle := app.NewMantle(
		db,
		tg.GetGenesisDoc(),
	)

	mantle.Server(1337)

	validators := tg.GetValidators()

	var valpvs = make([]tm.PrivValidator, len(validators))
	for i, vm := range validators {
		valpvs[i] = vm.PrivValidator
	}

	return &TestkitContext{
		tg:         tg,
		m:          new(sync.RWMutex),
		validators: validators,
		vc:         NewValidatorContext(valpvs),
		mantle:     mantle,
		mempool:    []types.StdTx{},
		db:         db,
		// auto related
		autoTxs: make([]AutomaticTxEntry, 0),
		autoInjection: &AutomaticInjection{
			isEnabled:    false,
			lastProposer: 0,
			valRounds:    nil,
		},
	}
}

func (ctx *TestkitContext) GetMantle() *app.Mantle {
	return ctx.mantle
}

func (ctx *TestkitContext) ClearMempool() {
	ctx.m.Lock()
	ctx.mempool = []types.StdTx{}
	ctx.m.Unlock()
}

func (ctx *TestkitContext) AddToMempool(tx types.StdTx) (*types.BlockState, error) {
	ctx.m.Lock()
	ctx.mempool = append(ctx.mempool, tx)
	ctx.m.Unlock()

	// if auto injection is enabled, run injection as well
	if ctx.autoInjection.isEnabled {
		proposer := ctx.autoInjection.NextProposer()
		return ctx.Inject(ctx.PickProposerByAddress(proposer))
	}

	// if manual mode, return nil
	return nil, nil
}

func (ctx *TestkitContext) Inject(proposer tm.PrivValidator) (*types.BlockState, error) {
	lastState := ctx.mantle.GetLastState()

	// create block according to the last block
	nextBlock := NewBlock(lastState)

	// set all txs
	for _, tx := range ctx.mempool {
		nextBlock = nextBlock.WithTx(tx)
	}

	// any auto tx? put them in block
	// run them in parallel
	m := new(sync.WaitGroup)
	m.Add(len(ctx.autoTxs))

	txs := make([]auth.StdTx, len(ctx.autoTxs))
	for i, atx := range ctx.autoTxs {
		// skip if atx period is not met
		currentHeight := nextBlock.nextBlock.Header.Height
		atxStartedAt := atx.StartedAt

		// skip if startedAt is not met
		if currentHeight < atxStartedAt {
			m.Done()
			continue
		}

		// period check
		if (currentHeight-atxStartedAt)%int64(atx.Period) != 0 {
			m.Done()
			continue
		}

		// otherwise lets go
		index := i
		atxCopy := atx

		go func() {
			txs[index] = NewSignedTx(
				atxCopy.Msgs,
				atxCopy.Fee,
				atxCopy.AccountName,
				ctx.tg.GetKeybase(),
				ctx.tg.chainId,
				ctx.mantle.GetApp(),
			)
			m.Done()
		}()

	}

	m.Wait()

	for _, tx := range txs {
		if len(tx.Msgs) == 0 {
			continue
		}
		nextBlock.WithTx(tx)
	}

	// prep block for injection
	blockToInject := nextBlock.Finalize()

	log.Printf("[mantle/testkit/context] injecting block %d with %d txs", blockToInject.Height, len(blockToInject.Txs))

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

func (ctx *TestkitContext) PickProposerByIndex(index int) tm.PrivValidator {
	if index > len(ctx.validators) {
		panic("not enough validators")
	}
	return ctx.validators[index].PrivValidator
}

func (ctx *TestkitContext) PickProposerByAddress(address sdk.ValAddress) tm.PrivValidator {
	for _, val := range ctx.validators {
		if val.Account.Equals(address) {
			return val.PrivValidator
		}
	}

	panic("validator not found")
}
