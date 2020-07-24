package test

import (
	"fmt"

	dbm "github.com/tendermint/tm-db"
	"github.com/terra-project/mantle/app"
	"github.com/terra-project/mantle/utils"
)

func CreateTestDB() dbm.DB {
	return dbm.NewMemDB()
}

func CreateTestApp() *app.App {
	db := CreateTestDB()
	genesis := utils.GenesisDocFromFile(fmt.Sprintf("%s/%s", utils.Rootdir(), "./test/genesis.json"))
	return app.NewApp(db, genesis)
}
