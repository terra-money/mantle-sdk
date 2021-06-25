package mantlemint

import (
	abcicli "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/mock"
	"github.com/tendermint/tendermint/state"
	tmdb "github.com/tendermint/tm-db"
	TerraApp "github.com/terra-money/core/app"
	"io/ioutil"
	"sync"
)

func NewMantlemintExecutor(
	db tmdb.DB,
	AppConn abcicli.Client,
) *state.BlockExecutor {
	return state.NewBlockExecutor(
		db,
		log.NewTMLogger(ioutil.Discard),
		AppConn,
		mock.Mempool{},           // no mempool
		state.MockEvidencePool{}, // no evidence pool
	)
}

func NewMantlemintSimulationExecutor(
	db tmdb.DB,
	AppConn abcicli.Client,
) *state.BlockExecutor {
	return state.NewBlockExecutor(
		db,
		log.NewTMLogger(ioutil.Discard),
		AppConn,
		mock.Mempool{},           // no mempool
		state.MockEvidencePool{}, // no evidence pool
	)
}

func NewMantleAppConn(
	terraApp *TerraApp.TerraApp,
) abcicli.Client {
	mtx := new(sync.Mutex)
	return abcicli.NewLocalClient(mtx, terraApp)
}
