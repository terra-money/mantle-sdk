package main

import (
	"github.com/terra-project/mantle-sdk/testkit/rpc"
)

func main() {
	rpc.StartRPCServer(11317, new(rpc.TestkitRPCContext))
}
