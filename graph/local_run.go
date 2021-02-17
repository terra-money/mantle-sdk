package graph

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	"github.com/terra-project/mantle-sdk/serdes"
	"github.com/terra-project/mantle-sdk/types"
	"sync"
)

type InternalRunResult struct {
	Data interface{}
	Err  error
}

// skip all proxy_resolver fields
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
		return internalResultInvariant(errors)
	}

	fields := p.Schema.QueryType().Fields()
	gqlRunResult := types.GraphQLInternalResult{
		Data:   make(map[string][]byte),
		Errors: nil,
	}

	// parse out root
	querySelections := AST.Definitions[0].(*ast.OperationDefinition).SelectionSet.Selections

	// prepare for a parallel resolve
	// some resolvers might return a thunk, and we need to
	// wait for them in goroutines to avoid congestion
	wg := sync.WaitGroup{}
	wg.Add(len(querySelections))

	//
	sync := sync.RWMutex{}

	// go
	for _, selection := range querySelections {
		go func() {
			defer wg.Done()

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
						case *ast.Name:
							variables[argumentName] = p.VariableValues[argumentValue.Value]
						case ast.Variable:
							variables[argumentName] = p.VariableValues[argumentValue.Name.Value]
						default:
							variables[argumentName] = argumentValue
						}
					}
				}

				rp := graphql.ResolveParams{
					Source:  source,
					Args:    variables,
					Info:    graphql.ResolveInfo{},
					Context: p.Context,
				}

				// mutex control
				sync.Lock()
				defer sync.Unlock()

				result, err := unthunkResult(fieldConfig.Resolve(rp))

				if err != nil {
					gqlRunResult.Errors = []error{err}
					return
				}

				pack, packErr := serdes.Serialize(nil, result)
				if packErr != nil {
					gqlRunResult.Errors = []error{err}
					return
				}

				// save
				gqlRunResult.Data[fieldName.Value] = pack
				return

			case *ast.FragmentDefinition:
				// noop yet
			}
		}()
	}

	wg.Wait()

	return &gqlRunResult
}

func unthunkResult(result interface{}, err error) (interface{}, error) {
	thunkFunc, isThunk := result.(func() (interface{}, error))

	// resolved data might be a thunk, run this thunk and wait for the result
	// in this way, thunk resolve becomes series but that's fine for proxy_resolver indexing mechanism
	if !isThunk {
		return result, err
	}

	return thunkFunc()
}

func internalResultInvariant(errors []error) *types.GraphQLInternalResult {
	return &types.GraphQLInternalResult{
		Errors: errors,
	}
}
