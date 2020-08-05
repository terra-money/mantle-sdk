package types

import (
	"time"

	"github.com/graphql-go/graphql"
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
	dependencies []ModelType,
) *graphql.Result
type GraphQLCommitter func(entity interface{}) error
