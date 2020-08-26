package generate

import (
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/terra-project/mantle/querier"
	"github.com/terra-project/mantle/utils"
)

func GenerateListGraphResolver(modelType reflect.Type, fieldConfig *graphql.Field) (*graphql.Field, error) {
	t := utils.GetType(modelType)
	entityName := t.Name()

	if fieldConfig.Args == nil {
		return nil, fmt.Errorf("GraphQL resolver arguments are never set. Creating list field is disallowed: %s", entityName)
	}

	// list objects have set of _range parameters defined
	rangeArgs := graphql.FieldConfigArgument{}
	for argName, arg := range fieldConfig.Args {
		rangeArgs[argName] = arg
		rangeArgs[fmt.Sprintf("%s_%s", argName, "range")] = &graphql.ArgumentConfig{
			Type:        graphql.NewList(arg.Type),
			Description: fmt.Sprintf("Ranged - %s", arg.Description),
		}
	}

	// if the output type is already a slice type,
	// don't make list of it.
	var outputType graphql.Output
	var outputName string
	if t.Kind() == reflect.Slice {
		outputName = fieldConfig.Name
		outputType = fieldConfig.Type
	} else {
		outputName = utils.Pluralize(fieldConfig.Name)
		outputType = graphql.NewList(fieldConfig.Type)
	}

	return &graphql.Field{
		Name: outputName,
		Args: rangeArgs,
		Type: outputType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			args := p.Args

			if args != nil {
				// query
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
				entities := make([]interface{}, 0)
				for documentKey := range intersection {
					doc, err := q.Get([]byte(documentKey))
					if err != nil {
						return nil, fmt.Errorf("Document(%s) does not exist.", documentKey)
					}

					docValue := reflect.New(t.Elem())
					if err := msgpack.Unmarshal(doc, docValue.Interface()); err != nil {
						return nil, err
					}

					entities = append(entities, docValue.Interface())
				}

				return entities, nil
			}

			// resolve current round

			return nil, nil

		},
	}, nil

}
