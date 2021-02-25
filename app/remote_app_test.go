package app

import (
	"github.com/terra-project/mantle-sdk/db/leveldb"
	"testing"
)

func TestRemoteApp(t *testing.T) {
	memdb := leveldb.NewLevelDB("remoteapp-test")
	remoteMantle := NewRemoteMantle(
		memdb,
		"https://tequila-mantle.terra.dev",
	)

	remoteMantle.Server(1337)
	remoteMantle.Sync(RemoteSyncConfiguration{})
}
