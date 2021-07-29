package testkit

import (
	"fmt"
	TerraApp "github.com/terra-project/core/app"
	"github.com/terra-project/core/x/auth"
	"github.com/terra-money/mantle-compatibility/localclient"
	"sync"
)

type LocalQuerier struct {
	lc localclient.LocalClient
}

func NewLocalQuerier(app *TerraApp.TerraApp) auth.NodeQuerier {
	return &LocalQuerier{
		lc: localclient.NewLocalClient(app, new(sync.Mutex)),
	}
}

func (lc *LocalQuerier) QueryWithData(path string, data []byte) ([]byte, int64, error) {
	fmt.Println(path, string(data))
	res, err := lc.lc.ABCIQuery(path, data)
	if err != nil {
		return nil, 0, err
	}

	return res.Response.Value, res.Response.Height, nil
}
