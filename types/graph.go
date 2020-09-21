package types

import (
	"time"
)

type BaseFilter struct {
	Height     uint64
	Datetime   time.Time
	ParentHash string
}

type GraphQLQueryRaw interface{}
type GraphQLParams map[string]interface{}
type GraphQLQuerier func(
	query string,
	variables GraphQLParams,
	dependencies []Model,
) GraphQLResult
type GraphQLCommitter func(entity interface{}) error

// graphql results
type GraphQLResult interface {
	HasErrors() bool
}

type GraphQLInternalResult struct {
	Data   map[string][]byte
	Errors []error
}

func (ir *GraphQLInternalResult) HasErrors() bool {
	return len(ir.Errors) > 0
}
