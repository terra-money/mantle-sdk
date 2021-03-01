package types

import (
	"github.com/graphql-go/graphql"
)

type GraphQLParams map[string]interface{}
type GraphQLQuerier func(
	query string,
	variables GraphQLParams,
	dependencies []Model,
) *graphql.Result
type GraphQLCommitter func(entity interface{}) error
