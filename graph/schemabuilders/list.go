package schemabuilders

import (
	"fmt"

	p "github.com/gertd/go-pluralize"
	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle/graph"
	"github.com/terra-project/mantle/querier"
	"github.com/terra-project/mantle/utils"
)

var pluralize = p.NewClient().Plural

func CreateListSchemaBuilder() graph.SchemaBuilder {
	return func(fields *graphql.Fields) error {
		for _, fieldConfig := range *fields {
			entityName := fieldConfig.Name
			entityType := fieldConfig.Type

			if fieldConfig.Args == nil {
				return fmt.Errorf("GraphQL resolver arguments are never set. Creating list field is disallowed: %s", entityName)
			}

			pluralizedName := pluralize(entityName)
			listField := graphql.Field{
				Name: pluralizedName,
				Args: fieldConfig.Args,
				Type: graphql.NewList(entityType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// for lists, the search is NEVER a hash seek.
					// grab data from kvstore and return as list
					args := p.Args
					q := p.Context.Value(utils.QuerierKey).(querier.Querier)

					searchedDocumentKeys := make(map[interface{}]bool)
					for indexKey, indexParam := range args {
						queryResolver, err := q.Build(entityName, indexKey, indexParam)
						if err != nil {
							return nil, err
						}

						resolvedDocumentKey, err := queryResolver.Resolve()
						if err != nil {
							return nil, err
						}

						// set this key to be found
						searchedDocumentKeys[resolvedDocumentKey] = true
					}

					// iterate again and get actual values
					entities := make([]interface{}, len(searchedDocumentKeys))
					for documentKey := range searchedDocumentKeys {
						entity, err := q.Get((documentKey).([]byte))
						if err != nil {
							return nil, fmt.Errorf("Document(%s) does not exist.", documentKey)
						}
						entities = append(entities, entity)
					}

					return entities, nil
				},
			}

			(*fields)[pluralizedName] = &listField
		}

		return nil
	}
}
