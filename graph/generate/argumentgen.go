package generate

import (
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/utils"
)

func GenerateArgument(kvindex *kvindex.KVIndex) graphql.FieldConfigArgument {
	args := graphql.FieldConfigArgument{}

	// by default, all entities have "Height (uint64)" index
	// since it's not part of KVIndex, add them here.
	args["Height"] = &graphql.ArgumentConfig{
		Description: fmt.Sprint("Height is the absolute block height when indexed"),
		Type: utils.GetGraphQLType(reflect.Uint64),
	}

	entries := kvindex.GetEntries()
	for indexName, indexEntry := range entries {
		entry := indexEntry.GetEntry()
		args[indexName] = &graphql.ArgumentConfig{
			Description: fmt.Sprintf("index %s with type %s", indexName, entry.Type.String()),
			Type:        utils.GetGraphQLType(entry.Type),
		}
	}

	return args
}
