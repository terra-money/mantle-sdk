package middlewares

import (
	abcicli "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/state"
	"github.com/terra-project/mantle-sdk/app/mantlemint"
)

type IndexerMiddlewareCallback func(responses state.ABCIResponses)
type IndexerMiddleware struct {
	abcicli.Client
	callback  IndexerMiddlewareCallback
	responses *state.ABCIResponses
}

func NewIndexerMiddleware(
	callback IndexerMiddlewareCallback,
) mantlemint.Middleware {
	return func(client abcicli.Client) abcicli.Client {
		return &IndexerMiddleware{
			Client:    client,
			callback:  callback,
			responses: nil,
		}
	}
}

func (middleware *IndexerMiddleware) BeginBlockSync(req abci.RequestBeginBlock) (*abci.ResponseBeginBlock, error) {
	middleware.responses = newABCIResponses()

	res, err := middleware.Client.BeginBlockSync(req)
	if err != nil {
		return nil, err
	}

	middleware.responses.BeginBlock = res
	return res, err
}

func (middleware *IndexerMiddleware) BeginBlockAsync(req abci.RequestBeginBlock) *abcicli.ReqRes {
	invariant()
	return nil
}

func (middleware *IndexerMiddleware) DeliverTxSync(req abci.RequestDeliverTx) (*abci.ResponseDeliverTx, error) {
	invariant()
	return nil, nil
}

func (middleware *IndexerMiddleware) DeliverTxAsync(req abci.RequestDeliverTx) *abcicli.ReqRes {
	res := middleware.Client.DeliverTxAsync(req)
	if r, ok := res.Response.Value.(*abci.Response_DeliverTx); ok {
		middleware.responses.DeliverTxs = append(middleware.responses.DeliverTxs, r.DeliverTx)
	}

	return res
}

func (middleware *IndexerMiddleware) EndBlockSync(req abci.RequestEndBlock) (*abci.ResponseEndBlock, error) {
	res, err := middleware.Client.EndBlockSync(req)
	if err != nil {
		return nil, err
	}

	middleware.responses.EndBlock = res
	return res, err
}

func (middleware *IndexerMiddleware) EndBlockAsync(req abci.RequestEndBlock) *abcicli.ReqRes {
	invariant()
	return nil
}

func (middleware *IndexerMiddleware) CommitSync() (*types.ResponseCommit, error) {
	// set response to flush after running this function
	defer func() {
		middleware.responses = nil
	}()

	// run indexer before client commit
	if middleware.responses != nil {
		middleware.callback(*middleware.responses)
	}

	// commit
	return middleware.Client.CommitSync()
}

func (middleware *IndexerMiddleware) CommitAsync() *abcicli.ReqRes {
	// not implemented
	invariant()
	return nil
}

//
func newABCIResponses() *state.ABCIResponses {
	return &state.ABCIResponses{
		DeliverTxs: make([]*abci.ResponseDeliverTx, 0),
		EndBlock:   nil,
		BeginBlock: nil,
	}
}

// private methods
var (
	errNotImplemented = "not implemented"
)

func invariant() {
	panic(errNotImplemented)
}
