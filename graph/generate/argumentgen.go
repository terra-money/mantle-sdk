package generate

import (
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/terra-money/mantle-sdk/db/kvindex"
	"github.com/terra-money/mantle-sdk/utils"
)

func GenerateArgument(kvindex *kvindex.KVIndex) graphql.FieldConfigArgument {
	args := graphql.FieldConfigArgument{}

	// index "Height" exists for ALL models w/ kvindex.IsPrimary
	// by default, all entities have "Height (uint64)" index
	// since it's not part of KVIndex, add them here.
	if !kvindex.IsPrimaryKeyedModel() {
		args["Height"] = &graphql.ArgumentConfig{
			Description: fmt.Sprint("Height is the absolute block height when indexed"),
			Type:        utils.GetGraphQLType(reflect.Uint64),
		}
	}

	entries := kvindex.Entries()
	for indexName, indexEntry := range entries {
		entry := indexEntry
		args[indexName] = &graphql.ArgumentConfig{
			Description: fmt.Sprintf("index %s with type %s", indexName, entry.Type().Name()),
			Type:        utils.GetGraphQLType(entry.Type().Kind()),
		}
	}

	return args
}
