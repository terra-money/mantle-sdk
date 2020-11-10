package generate

import (
	"fmt"
	"github.com/terra-project/mantle-sdk/constants"
	"github.com/terra-project/mantle-sdk/depsresolver"
	"github.com/terra-project/mantle-sdk/graph"
	"github.com/terra-project/mantle-sdk/serdes"
	"github.com/terra-project/mantle-sdk/types"
	"reflect"
	"sort"

	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle-sdk/querier"
	"github.com/terra-project/mantle-sdk/utils"
)

func GenerateListGraphResolver(modelType reflect.Type, fieldConfig *graphql.Field) (*graphql.Field, error) {
	t := utils.GetType(modelType)
	entityName := t.Name()

	// list objects have set of _range parameters defined
	rangeArgs := graphql.FieldConfigArgument{}
	for argName, arg := range fieldConfig.Args {
		rangeArgs[argName] = arg
		rangeArgs[fmt.Sprintf("%s_%s", argName, "range")] = &graphql.ArgumentConfig{
			Type:        graphql.NewList(arg.Type),
			Description: fmt.Sprintf("Ranged - %s", arg.Description),
		}
	}

	// set limit argument
	rangeArgs["Limit"] = &graphql.ArgumentConfig{
		Type:        graphql.Int,
		Description: "Limit",
	}

	// order scalars
	rangeArgs["Order"] = &graphql.ArgumentConfig{
		Type:        types.Order,
		Description: "Sort order",
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
			filteredArgs := FilterArgs(args, ReservedArgKeys)

			var limit = constants.DefaultLimit
			if customLimit, customLimitExists := args["Limit"]; customLimitExists {
				limit = customLimit.(int)
			}

			if p.Context.Value(utils.ImmediateResolveFlagKey).(bool) || args != nil && len(filteredArgs) > 0 {
				// query
				q := p.Context.Value(utils.QuerierKey).(querier.Querier)

				// must of been map[[]byte]bool but golang doesn't like that
				// for byte comparison purposes, string is fine
				intersectionSets := make([]map[string]bool, 0)

				// only do sequential indexer when no specific arguments is given
				if len(filteredArgs) == 0 {
					queryResolver, err := q.Build(entityName, "", nil)
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
						// stop before limit is filled
						if len(keysHashMap) >= limit {
							break
						}

						keysHashMap[string(it.Key())] = true
						it.Next()
					}

					it.Close()

					intersectionSets = append(intersectionSets, keysHashMap)
				}

				// run them in thunk
				pctx := graph.CreateParallel(len(filteredArgs))

				for indexKey, indexParam := range filteredArgs {
					pctx.Run(func() (interface{}, error) {
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

						return keysHashMap, nil
					})
					//
					// queryResolver, err := q.Build(entityName, indexKey, indexParam)
					// if err != nil {
					// 	return nil, err
					// }
					//
					// it, err := queryResolver.Resolve()
					// if err != nil {
					// 	return nil, err
					// }
					//
					// // for every key found, mark them found
					// keysHashMap := make(map[string]bool)
					//
					// for it.Valid() {
					// 	keysHashMap[string(it.Key())] = true
					// 	it.Next()
					// }
					//
					// it.Close()
					//
					// intersectionSets = append(intersectionSets, keysHashMap)
				}

				intersectionSetResults := pctx.Sync()

				// if at least 1 error exists,
				// cut this resolver and return error
				if pctx.HasErrors() {
					for _, r := range intersectionSetResults {
						return nil, r.Error
					}
				}

				for _, r := range intersectionSetResults {
					intersectionSets = append(intersectionSets, r.Result.(map[string]bool))
				}

				// if intersectionSets was never populated,
				// we couldn't find anything. return nil
				if len(intersectionSets) == 0 {
					return nil, nil
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

				// order intersection
				var sortedIntersection = make([]string, len(intersection))
				var i = 0
				for key, _ := range intersection {
					sortedIntersection[i] = key
					i++
				}

				var currentSortOrder = constants.GetOrder(args)
				switch currentSortOrder {
				case constants.ASC:
					sort.Sort(sort.StringSlice(sortedIntersection))
				case constants.DESC:
					sort.Sort(sort.Reverse(sort.StringSlice(sortedIntersection)))
				}

				// iterate again and get actual values
				entities := make([]interface{}, 0)
				var count = 0
				for _, documentKey := range sortedIntersection {
					if limit != 0 && count >= limit {
						break
					}

					doc, err := q.Get([]byte(documentKey))
					if err != nil {
						return nil, fmt.Errorf("Document(%s) does not exist.", documentKey)
					}

					docValue := reflect.New(t.Elem())

					if err := serdes.Deserialize(t, doc, docValue.Interface()); err != nil {
						return nil, fmt.Errorf("document unpack failed, path=%s, err=%s", p.Info.Path.AsArray(), err)
					}

					entities = append(entities, docValue.Interface())
					count++
				}

				return entities, nil
			}

			depsResolver, depsResolverNotCleared := p.Context.Value(utils.DepsResolverKey).(depsresolver.DepsResolver)
			if !depsResolverNotCleared {
				panic(fmt.Sprintf("DepsResolver is either cleared or not set, in ResolveFunc for %s", entityName))
			}

			// resolve current round
			var dependencies = p.Context.Value(utils.DependenciesKey).(utils.DependenciesKeyType)
			var isAwaitedDependency = p.Args["Height"] == nil
			var isSelfReferencing = dependencies[entityName] == true

			if isAwaitedDependency && isSelfReferencing {
				return nil, fmt.Errorf("Self reference is disallowed. entityName=%s", entityName)
			}

			return depsResolver.Resolve(t), nil
		},
	}, nil

}
