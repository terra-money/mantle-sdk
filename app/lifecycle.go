package app

import (
	"github.com/terra-project/mantle/types"
	"github.com/terra-project/mantle/utils"
)

type LifecycleContext struct {
	app                   *App
	transactionalAppState bool
	txDecoder             utils.TxDecoder
}

func NewLifecycle(
	app *App,
	transactionalAppState bool,
) *LifecycleContext {
	return &LifecycleContext{
		app:                   app,
		transactionalAppState: transactionalAppState,
		txDecoder:             utils.NewDecoder(),
	}
}

func (c *LifecycleContext) Start(
	blockChannel chan types.Block,
) chan types.BaseState {
	processedChannel := make(chan types.BaseState)
	go func() {
		for {
			block := <-blockChannel
			nextState := c.Inject(&block)
			processedChannel <- nextState
		}
	}()

	return processedChannel
}

func (c *LifecycleContext) Inject(block *types.Block) types.BaseState {
	// run begin blocker
	beginBlockerResponse := c.app.BeginBlocker(block)

	// run all txs
	deliverTxResponses := c.app.DeliverTxs(block.Data.Txs)

	// run end blocker
	endBlockerResponse := c.app.EndBlocker(block)

	// put together a primitive state
	txs := make([]types.Tx, 0)
	for _, tx := range block.Data.Txs {
		txdoc, err := c.txDecoder(tx)

		// TODO: handle me
		if err != nil {
			panic(err)
		}
		txs = append(txs, *txdoc)
	}

	primState := types.BaseState{
		Height:             block.Header.Height,
		BeginBlockResponse: beginBlockerResponse,
		EndBlockResponse:   endBlockerResponse,
		DeliverTxResponses: deliverTxResponses,
		Block:              *block,
		Txs:                txs,
	}

	return primState
}

func (c *LifecycleContext) Commit() []byte {
	commitResult := c.app.Commit(c.transactionalAppState)
	return commitResult.Data
}
