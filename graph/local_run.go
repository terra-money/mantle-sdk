package graph

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/terra-project/mantle/types"
	"github.com/vmihailenco/msgpack/v5"
)

// skip all internal fields
// msgpack will take care of them..
func InternalGQLRun(p graphql.Params) *types.GraphQLInternalResult {
	source := source.NewSource(&source.Source{
		Body: []byte(p.RequestString),
		Name: "GraphQL request",
	})

	AST, _ := parser.Parse(parser.ParseParams{Source: source})

	validationResult := graphql.ValidateDocument(&p.Schema, AST, nil)
	if !validationResult.IsValid {
		return InternalResultInvariant([]error{
			fmt.Errorf("nope"),
		})
	}

	fields := p.Schema.QueryType().Fields()
	resultMap := make(map[string][]byte)

	for _, definition := range AST.Definitions {
		switch definition := definition.(type) {
		case *ast.OperationDefinition:
			fieldName := definition.GetName()
			fieldConfig := fields[fieldName.Value]

			rp := graphql.ResolveParams{
				Source:  source,
				Args:    p.VariableValues,
				Info:    graphql.ResolveInfo{},
				Context: nil,
			}

			result, err := fieldConfig.Resolve(rp)
			if err != nil {
				return InternalResultInvariant([]error{
					err,
				})
			}

			pack, packErr := msgpack.Marshal(result)
			if packErr != nil {
				return InternalResultInvariant([]error{
					packErr,
				})
			}

			resultMap[fieldName.Value] = pack
		case *ast.FragmentDefinition:
			// noop yet
		}
	}

	return &types.GraphQLInternalResult{
		Data: resultMap,
	}
}

func InternalResultInvariant(errors []error) *types.GraphQLInternalResult {
	return &types.GraphQLInternalResult{
		Errors: errors,
	}
}
