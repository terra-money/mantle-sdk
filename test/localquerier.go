package test

import (
	"github.com/terra-project/core/x/auth"
	"github.com/terra-project/mantle-compatibility/localclient"
	"github.com/terra-project/mantle/app"
)

type LocalQuerier struct {
	lc localclient.LocalClient
}

func NewLocalQuerier() auth.NodeQuerier {
	return &LocalQuerier{lc: localclient.NewLocalClient(app.GlobalTerraApp)}
}

func (lc *LocalQuerier) QueryWithData(path string, data []byte) ([]byte, int64, error) {
	res, err := lc.lc.ABCIQuery(path, data)
	if err != nil {
		return nil, 0, err
	}

	return res.Response.Value, res.Response.Height, nil
}
