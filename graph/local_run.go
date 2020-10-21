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
		var errors []error
		for _, e := range validationResult.Errors {
			errors = append(errors, fmt.Errorf(e.Error()))
		}
		return InternalResultInvariant(errors)
	}

	fields := p.Schema.QueryType().Fields()
	resultMap := make(map[string][]byte)

	//
	querySelections := AST.Definitions[0].(*ast.OperationDefinition).SelectionSet.Selections

	for _, selection := range querySelections {
		switch selection := selection.(type) {
		case *ast.Field:
			fieldName := selection.Name
			fieldConfig := fields[fieldName.Value]

			var variables = make(map[string]interface{})

			// map arguments to query arguments
			if len(selection.Arguments) > 0 {
				for _, argument := range selection.Arguments {
					argumentName := argument.Name.Value
					switch argumentValue := argument.Value.GetValue().(type) {
					case ast.Variable:
						variables[argumentName] = p.VariableValues[argumentValue.Name.Value]
					default:
						variables[argumentName] = argumentValue
					}
				}
			}

			//args := getArgumentValues(fieldDef.Args, fieldAST.Arguments, eCtx.VariableValues)

			rp := graphql.ResolveParams{
				Source:  source,
				Args:    variables,
				Info:    graphql.ResolveInfo{},
				Context: p.Context,
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

//
//// Prepares an object map of argument values given a list of argument
//// definitions and list of argument AST nodes.
//func getArgumentValues(
//	argDefs []*graphql.Argument, argASTs []*ast.Argument,
//	variableValues map[string]interface{}) map[string]interface{} {
//
//	argASTMap := map[string]*ast.Argument{}
//	for _, argAST := range argASTs {
//		if argAST.Name != nil {
//			argASTMap[argAST.Name.Value] = argAST
//		}
//	}
//	results := map[string]interface{}{}
//	for _, argDef := range argDefs {
//		var (
//			tmp   interface{}
//			value ast.Value
//		)
//		if tmpValue, ok := argASTMap[argDef.PrivateName]; ok {
//			value = tmpValue.Value
//		}
//		if tmp = valueFromAST(value, argDef.Type, variableValues); isNullish(tmp) {
//			tmp = argDef.DefaultValue
//		}
//		if !isNullish(tmp) {
//			results[argDef.PrivateName] = tmp
//		}
//	}
//	return results
//}
