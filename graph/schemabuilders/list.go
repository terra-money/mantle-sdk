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
		listFields := map[string]*graphql.Field{}

		for _, fieldConfig := range *fields {
			entityName := fieldConfig.Name
			entityType := fieldConfig.Type

			if fieldConfig.Args == nil {
				return fmt.Errorf("GraphQL resolver arguments are never set. Creating list field is disallowed: %s", entityName)
			}

			// list objects have set of _range parameters defined
			rangeArgs := graphql.FieldConfigArgument{}
			for argName, arg := range fieldConfig.Args {
				rangeArgs[argName] = arg
				rangeArgs[fmt.Sprintf("%s_%s", argName, "range")] = &graphql.ArgumentConfig{
					Type: graphql.NewList(arg.Type),
					Description: fmt.Sprintf("Ranged - %s", arg.Description),
				}
			}

			pluralizedName := pluralize(entityName)
			listField := graphql.Field{
				Name: pluralizedName,
				Args: rangeArgs,
				Type: graphql.NewList(entityType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// for lists, the search is NEVER a hash seek.
					// grab data from kvstore and return as list
					args := p.Args
					q := p.Context.Value(utils.QuerierKey).(querier.Querier)

					// must of been map[[]byte]bool but golang doesn't like that
					// for byte comparison purposes, string is fine
					intersectionSets := make([]map[string]bool, 0)

					for indexKey, indexParam := range args {
						queryResolver, err := q.Build(entityName, indexKey, indexParam)
						if err != nil {
							return nil, err
						}

						it, err := queryResolver.Resolve()
						if err != nil {
							return nil, err
						}

						// for every key found, mark them found
						keysHashMap := make(map[string]bool)
						for it.Valid() {
							keysHashMap[string(it.Key())] = true
							it.Next()
						}

						it.Close()

						intersectionSets = append(intersectionSets, keysHashMap)
					}

					// find intersections
					intersection := intersectionSets[0]
					for _, set := range intersectionSets[1:] {
						nextIntersection := map[string]bool{}
						for key, _ := range set {
							if _, ok := intersection[key]; ok {
								nextIntersection[key] = true
							}
						}

						intersection = nextIntersection
					}

					// iterate again and get actual values
					entities := make([]interface{}, len(intersection))
					for documentKey := range intersection {
						entity, err := q.Get([]byte(documentKey))
						if err != nil {
							return nil, fmt.Errorf("Document(%s) does not exist.", documentKey)
						}
						entities = append(entities, entity)
					}

					return entities, nil
				},
			}

			listFields[pluralizedName] = &listField
		}

		// add
		for key, field := range listFields {
			(*fields)[key] = field
		}

		return nil
	}
}
