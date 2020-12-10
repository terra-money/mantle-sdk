package mantlemint

import (
	abcicli "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/state"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/terra-project/mantle-sdk/types"
)

type Mantlemint interface {
	Inject(*types.Block) (*types.BlockState, error)
	Init(*tmtypes.GenesisDoc) error
	GetCurrentHeight() int64
	GetCurrentBlock() *types.Block
	GetCurrentState() state.State
	SetBlockExecutor(executor MantlemintExecutor)
}

type MantlemintExecutor interface {
	ApplyBlock(state.State, tmtypes.BlockID, *types.Block) (state.State, int64, error)
	SetEventBus(publisher tmtypes.BlockEventPublisher)
}

type Middleware func(conn abcicli.Client) abcicli.Client
