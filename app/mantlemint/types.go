package mantlemint

import (
	abcicli "github.com/tendermint/tendermint/abci/client"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/terra-project/mantle-sdk/types"
)

type Mantlemint interface {
	Inject(*types.Block) (*types.BlockState, error)
	Init(*tmtypes.GenesisDoc) error
	GetCurrentHeight() int64
	GetCurrentBlock() *types.Block
}

type MantlemintContext map[string]interface{}
type MantlemintBlockFinalizer func(block *types.Block)
type Middleware func(conn abcicli.Client) abcicli.Client
