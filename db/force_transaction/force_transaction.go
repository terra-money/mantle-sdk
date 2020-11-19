package force_transaction

import (
	"github.com/terra-project/mantle-sdk/db"
)

type GlobalTransactionManager struct {
	db.DB
	session db.Batch
}

// WithGlobalTransactionManager wraps around db,
// and let all individual Set() operations go through a transaction session.
//
// This is only possible because mantle uses a single db for
// tendermint, cosmos and indexer outputs.
func WithGlobalTransactionManager(db db.DB) db.DBwithGlobalTransaction {
	return &GlobalTransactionManager{
		DB: db,
	}
}

func (ft *GlobalTransactionManager) SetGlobalTransactionBoundary() {
	ft.session = ft.Batch()
}


func (ft *GlobalTransactionManager)	FlushGlobalTransactionBoundary() error {
	defer func() {
		ft.session = nil
	}()

	return ft.session.Flush()
}

func (ft *GlobalTransactionManager) Set(key []byte, data []byte) error {
	if ft.session != nil {
		return ft.session.Set(key, data)
	} else {
		return ft.DB.Set(key, data)
	}
}

// Batch() always returns the currently set batch
func (ft *GlobalTransactionManager) Batch() db.Batch {
	if ft.session != nil {
		return ft.session
	} else {
		return ft.DB.Batch()
	}
}

func (ft *GlobalTransactionManager) Delete(key []byte) error {
	if ft.session != nil {
		return ft.session.Delete(key)
	} else {
		return ft.DB.Delete(key)
	}
}
