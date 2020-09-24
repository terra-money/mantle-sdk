package generate

import (
	"fmt"
	"github.com/terra-project/mantle/depsresolver"
	"github.com/terra-project/mantle/serdes"
	"reflect"
	"sort"

	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle/querier"
	"github.com/terra-project/mantle/utils"
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
	rangeArgs["limit"] = &graphql.ArgumentConfig{
		Type:        graphql.Int,
		Description: "Limit",
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

			if args != nil && len(args) > 0 {
				// query
				q := p.Context.Value(utils.QuerierKey).(querier.Querier)

				// must of been map[[]byte]bool but golang doesn't like that
				// for byte comparison purposes, string is fine
				intersectionSets := make([]map[string]bool, 0)

				for indexKey, indexParam := range FilterArgs(args, ReservedArgKeys) {
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
				var currentSortOrder = defaultOrder
				desigatedSortOrder, exists := p.Args["sort"]
				if exists {
					currentSortOrder = stringOrderToConst[desigatedSortOrder.(string)]
				}
				switch currentSortOrder {
				case DESC:
					sort.Sort(sort.Reverse(sort.StringSlice(sortedIntersection)))
				case ASC:
					sort.Sort(sort.StringSlice(sortedIntersection))
				}

				// iterate again and get actual values
				entities := make([]interface{}, 0)
				var count = 0
				var limit = p.Args["limit"]
				for _, documentKey := range sortedIntersection {
					if limit != nil && count > p.Args["limit"].(int) {
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

			// if resolve immediately flag is set, don't await for dependency.
			// in case of list resolver, data must be fetched from database
			// with default entry size of 100
			if p.Context.Value(utils.ImmediateResolveFlagKey).(bool) {
				q := p.Context.Value(utils.QuerierKey).(querier.Querier)
				resolver, resolverErr := q.Build(entityName, "", nil)
				if resolverErr != nil {
					return nil, fmt.Errorf("resolver build failed, entityName=%s, err=%s", entityName, resolverErr)
				}

				it, itErr := resolver.Resolve()
				if itErr != nil {
					return nil, fmt.Errorf("resolver iteration failed, entityName=%s, err=%s", entityName, itErr)
				}

				var documentKeys = make([][]byte, defaultLimit)
				var i = 0

				for it.Valid() && i < defaultLimit {
					if i >= defaultLimit {
						break
					}
					documentKeys[i] = it.Key()
					it.Next()
					i++
				}

				it.Close()

				entities := make([]interface{}, 0)

				for _, documentKey := range documentKeys {
					if len(documentKey) == 0 {
						break
					}
					doc, err := q.Get(documentKey)
					if err != nil {
						return nil, fmt.Errorf("document(%s) does not exist.", documentKey)
					}

					docValue := reflect.New(t.Elem())
					if err := serdes.Deserialize(t, doc, docValue.Interface()); err != nil {
						return nil, fmt.Errorf("msgunpack failed, entityName=%s, err=%s", entityName, err)
					}

					entities = append(entities, docValue.Interface())
				}

				return entities, nil
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
