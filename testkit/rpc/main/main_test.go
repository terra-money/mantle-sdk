package main

import (
	"github.com/terra-project/mantle-sdk/testkit/rpc"
	"testing"
)

func TestApp(t *testing.T) {
	rpc.StartRPCServer(11317, rpc.NewTestkitRPCContext())
}
