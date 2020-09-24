package main

import (
	"log"
	"net/http"

	"github.com/terra-project/mantle/utils"

	_ "net/http/pprof"

	"github.com/terra-project/mantle/app"
	"github.com/terra-project/mantle/db/badger"
)

func main() {
	badgerdb := badger.NewBadgerDB("mantle-db")
	defer badgerdb.Close()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	mantle := app.NewMantle(
		badgerdb,
		utils.GenesisDocFromFile("./genesis.json"),
	)

	mantle.Server()
	mantle.Sync(app.SyncConfiguration{
		TendermintEndpoint: "localhost:26657",
	})
}
