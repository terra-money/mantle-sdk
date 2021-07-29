package middlewares

import (
	abcicli "github.com/tendermint/tendermint/abci/client"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/terra-money/mantle-sdk/app/mantlemint"
)

type CacheMiddlewareCallback func(lastKnownHeight int64)
type CacheMiddleware struct {
	abcicli.Client
	callback        CacheMiddlewareCallback
	LastKnownHeight int64
}

// CacheMiddleware invalidates cache at each commit
func NewCacheMiddleware(
	callback CacheMiddlewareCallback,
) mantlemint.Middleware {
	return func(client abcicli.Client) abcicli.Client {
		return &CacheMiddleware{
			Client:          client,
			callback:        callback,
			LastKnownHeight: -1,
		}
	}
}

func (middleware *CacheMiddleware) BeginBlockSync(req abci.RequestBeginBlock) (*abci.ResponseBeginBlock, error) {
	middleware.LastKnownHeight = req.Header.Height
	return middleware.Client.BeginBlockSync(req)
}

func (middleware *CacheMiddleware) BeginBlockAsync(req abci.RequestBeginBlock) *abcicli.ReqRes {
	middleware.LastKnownHeight = req.Header.Height
	return middleware.Client.BeginBlockAsync(req)
}

func (middleware *CacheMiddleware) DeliverTxSync(req abci.RequestDeliverTx) (*abci.ResponseDeliverTx, error) {
	return middleware.Client.DeliverTxSync(req)
}

func (middleware *CacheMiddleware) DeliverTxAsync(req abci.RequestDeliverTx) *abcicli.ReqRes {
	return middleware.Client.DeliverTxAsync(req)
}

func (middleware *CacheMiddleware) EndBlockSync(req abci.RequestEndBlock) (*abci.ResponseEndBlock, error) {
	return middleware.Client.EndBlockSync(req)
}

func (middleware *CacheMiddleware) EndBlockAsync(req abci.RequestEndBlock) *abcicli.ReqRes {
	return middleware.Client.EndBlockAsync(req)
}

func (middleware *CacheMiddleware) CommitSync() (*abci.ResponseCommit, error) {
	middleware.callback(middleware.LastKnownHeight)
	return middleware.Client.CommitSync()
}

func (middleware *CacheMiddleware) CommitAsync() *abcicli.ReqRes {
	middleware.callback(middleware.LastKnownHeight)
	return middleware.Client.CommitAsync()
}
