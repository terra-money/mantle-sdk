package types

type IndexerQuerier func(request interface{}, variables GraphQLParams) error
type IndexerCommitter func(entity interface{})
type Indexer func(query Query, commit Commit)
type IndexerRegisterer func(Register)

// some aliases
type Query = IndexerQuerier
type Commit = IndexerCommitter

type RegisterIndexer func(
	indexer Indexer,
	models ...interface{},
)
