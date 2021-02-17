package schemabuilders

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/terra-project/mantle-sdk/graph"
	. "github.com/terra-project/mantle-sdk/graph/proxy_resolver"
	"github.com/terra-project/mantle-sdk/utils"
)

type (
	RemoteQueriesMap map[string]TypeDescriptor
	SubGraphRecFunc  func(nodeName string)
)

var (
	errInvalidContext = fmt.Errorf("context is not proxy resolver context")
	errInvalidSource  = fmt.Errorf("source is not proxy resolver context")
)

func CreateRemoteModelSchemaBuilder(baseMantleEndpoint string) graph.SchemaBuilder {
	return func(fields *graphql.Fields) error {
		schema := NewIntrospection(baseMantleEndpoint)

		// 1. go through all types, and create objects first
		remoteModelsMap := convertTypesToTypeMap(schema.Types)

		// 2. iterate through all FIELDS in RootQuery,
		//    and map out all things recursively
		rootQuery, ok := remoteModelsMap["RootQuery"]
		if !ok {
			return fmt.Errorf("remote mantle does not have root query")
		}

		rootQueryFields := rootQuery.Fields

		// iterate over queriable field, reconstruct query
		for _, queriableField := range rootQueryFields {
			name := queriableField.Name

			// query output object
			fieldOutput := reconstructFieldConfig(
				queriableField.Type,
				remoteModelsMap,
			)

			fieldArguments := reconstructFieldArgument(
				queriableField.Args,
				remoteModelsMap,
			)

			var deprecationReason string
			if queriableField.DeprecationReason != "" {
				deprecationReason = fmt.Sprintf("(proxy) %s", queriableField.DeprecationReason)
			}

			//
			(*fields)[name] = &graphql.Field{
				Name: name,
				Type: fieldOutput,
				Args: fieldArguments,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// recreate a full operation tree, make a request,
					// prop down the result.
					prc, ok := p.Context.Value(utils.ProxyResolverContextKey).(*ProxyResolverContext)
					if !ok {
						return nil, errInvalidContext
					}

					fmt.Println(p.Info.Path, p.Info.Operation)

					selection, isSelectionOk := p.Info.Path.Key.(string)
					if !isSelectionOk {
						selection = name
					}

					subgraph := createSubgraph(prc, p.Info.Operation.GetSelectionSet()).WithGraphQLVariables(p.Info.VariableValues)
					if source, err := subgraph.ResolveRoot(); err == nil {
						return source[selection], nil
					} else {
						return nil, err
					}
				},
				DeprecationReason: deprecationReason,
				Description:       fmt.Sprintf("(proxy) %s", queriableField.Description),
			}
		}

		return nil
	}
}

func convertTypesToTypeMap(types []TypeDescriptor) Definitions {
	m := make(map[string]TypeDescriptor)

	for _, t := range types {
		m[t.Name] = t
	}

	return m
}

func reconstructFieldConfig(
	queryType Type,
	remoteQueriesMap Definitions,
) graphql.Output {
	return GetGraphQLOutputType(queryType, remoteQueriesMap)
}

func reconstructFieldArgument(
	queryArguments []Argument,
	remoteQueriesMap Definitions,
) graphql.FieldConfigArgument {
	argumentConfig := graphql.FieldConfigArgument{}

	for _, queryArgument := range queryArguments {
		argumentConfig[queryArgument.Name] = &graphql.ArgumentConfig{
			Type:         GetGraphQLInputType(queryArgument.Type, remoteQueriesMap),
			DefaultValue: queryArgument.DefaultValue,
			Description:  queryArgument.Description,
		}
	}

	return argumentConfig
}

func createSubgraph(prc *ProxyResolverContext, selections *ast.SelectionSet) *ProxyResolverContext {
	// skip subgraph creation if no selections
	if selections == nil {
		return nil
	}

	for _, s := range selections.Selections {
		switch f := s.(type) {
		case *ast.Field:
			subgraphPrc := prc.CreateSubtree(f.Name.Value, astArgumentsToMap(f.Arguments))

			if f.Alias != nil {
				subgraphPrc.WithAlias(f.Alias.Value)
			}

			createSubgraph(subgraphPrc, f.SelectionSet)

			// TODO: implement me
			// case *ast.FragmentSpread:
			// 	fmt.Printf(f.Name.Value)
			// 	explainFullRequestTree(offset+1, f.GetSelectionSet())
			// case *ast.InlineFragment:
			// 	fmt.Printf(f.Kind)
			// 	explainFullRequestTree(offset+1, f.GetSelectionSet())
		}
	}

	return prc
}

func astArgumentsToMap(astArguments []*ast.Argument) map[string]ast.Value {
	argmap := make(map[string]ast.Value)
	for _, astArgument := range astArguments {
		argmap[astArgument.Name.Value] = astArgument.Value
	}

	return argmap
}
