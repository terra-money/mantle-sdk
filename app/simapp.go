package app

import (
	compattypes "github.com/terra-project/mantle-compatibility/types"
	"github.com/terra-project/mantle-sdk/db/badger"
	"github.com/terra-project/mantle-sdk/types"
)

type (
	IndexerTracker func(queryRecord []QueryRecord, commitRecord []CommitRecord)
	QueryRecord    struct {
		Request   interface{}
		Variables types.GraphQLParams
	}
	CommitRecord interface{}
)

func NewSimMantle(genesis *compattypes.GenesisDoc, indexers ...types.IndexerRegisterer) *Mantle {
	return NewMantle(
		badger.NewBadgerDB(""),
		genesis,
		indexers...,
	)
}

// TrackIndexerResult is a HoC function which takes
// the original registerer and an additional tracker.
func TrackIndexerResult(
	registerer types.IndexerRegisterer,
	tracker IndexerTracker,
) types.IndexerRegisterer {
	var oIndexer types.Indexer
	var oModels []types.Model
	fauxRegister := func(indexer types.Indexer, models ...types.Model) {
		oIndexer = indexer
		oModels = models
	}

	// save oIndexer and oModels
	registerer(fauxRegister)

	var queryRecords []QueryRecord
	var commitRecords []CommitRecord

	// create a faux indexer, which calls the
	// original indexer w/ tracking history of query and commit calls
	var fauxIndexer types.Indexer = func(query types.Query, commit types.Commit) error {
		fauxQuery := func(request interface{}, variables types.GraphQLParams) error {
			record := QueryRecord{}
			record.Variables = variables

			queryErr := query(request, variables)

			record.Request = request
			queryRecords = append(queryRecords, record)
			return queryErr
		}

		fauxCommit := func(entity interface{}) error {
			commitRecords = append(commitRecords, entity)
			return commit(entity)
		}

		defer tracker(queryRecords, commitRecords)

		return oIndexer(fauxQuery, fauxCommit)
	}

	return func(register types.Register) {
		register(fauxIndexer, oModels...)
	}
}
