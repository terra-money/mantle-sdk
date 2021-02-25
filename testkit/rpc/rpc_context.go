package rpc

import (
	"github.com/terra-project/mantle-sdk/testkit"
)

type TestkitRPCContext struct {
	ctxs           map[string]*testkit.TestkitContext
	lastMantlePort int
}

func NewTestkitRPCContext() *TestkitRPCContext {
	return &TestkitRPCContext{
		ctxs:           make(map[string]*testkit.TestkitContext),
		lastMantlePort: 51337,
	}
}

func (rpcCtx *TestkitRPCContext) SetTestkitContext(
	identifier string,
	testkitContext *testkit.TestkitContext,
) {
	rpcCtx.ctxs[identifier] = testkitContext
}

func (rpcCtx *TestkitRPCContext) GetNextMantlePort() int {
	rpcCtx.lastMantlePort = rpcCtx.lastMantlePort + 1
	return rpcCtx.lastMantlePort
}

func (rpcCtx *TestkitRPCContext) GetContext(ctxId string) *testkit.TestkitContext {
	return rpcCtx.ctxs[ctxId]
}
