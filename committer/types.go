package committer

type Committer interface {
	Commit(height uint64, entities ...interface{}) error
}
