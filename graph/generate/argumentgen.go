package generate

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/utils"
)

func GenerateArgument(kvindex *kvindex.KVIndex) graphql.FieldConfigArgument {
	args := graphql.FieldConfigArgument{}

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
