package proxy_resolver

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle-sdk/graph/graph_types"
)

func GetGraphQLOutputType(outputType Type, definitions Definitions) graphql.Output {
	switch outputType.Kind {
	// graphql-native graph_types + custom graph_types
	case "SCALAR":
		switch outputType.Name {
		// take care of all graphql-native graph_types
		case "Int":
			return graphql.Int
		case "Float":
			return graphql.Float
		case "String":
			return graphql.String
		case "Boolean":
			return graphql.Boolean
		case "ID":
			return graphql.ID
		case "DateTime":
			return graphql.DateTime
		case "Buffer":
			return graphql.String
		// check cosmos-scalar map
		default:
			outputName, ok := outputType.Name.(string)
			if !ok {
				panic(errNoName(outputType.Name))
			}
			cosmosScalar := graph_types.GetCosmosScalarByName(outputName)

			// if name is unknown, mantle can't handle it. panic here
			if cosmosScalar == nil {
				fmt.Printf("warning: unknown scalar %v, coercing to string\n", outputName)
				return graphql.String
			}

			return cosmosScalar
		}

	// get the matching OBJECT type of definitions, reconstruct them into graphql type
	case "OBJECT":
		objectName, ok := outputType.Name.(string)
		if !ok {
			panic(errNoName(outputType.Name))
		}

		if objectCache := getTypeFromCache(objectName); objectCache != nil {
			return objectCache
		}

		definition := definitions[objectName]
		subselectionFields := graphql.Fields{}
		subselections := definition.Fields

		for _, selection := range subselections {
			selectionName := selection.Name
			selectionType := GetGraphQLOutputType(selection.Type, definitions)

			subselectionFields[selectionName] = &graphql.Field{
				Name: selectionName,
				Type: selectionType,
				Args: nil, // TODO: enable this as we allow subgraph arguments
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, isSourceValue := p.Source.(map[string]interface{})
					if !isSourceValue {
						return nil, errInvalidSource(p.Info.Path.AsArray())
					}

					alias, isPathString := p.Info.Path.Key.(string)
					if !isPathString {
						alias = selectionName
					}

					return source[alias], nil
				},
				DeprecationReason: selection.DeprecationReason,
				Description:       selection.Description,
			}
		}

		object := graphql.NewObject(graphql.ObjectConfig{
			Name:        definition.Name,
			Interfaces:  nil, // TODO: fix me
			Fields:      subselectionFields,
			IsTypeOf:    nil, // TODO: fix me
			Description: definition.Description,
		})

		setTypeInCache(definition.Name, object)

		return object

	case "LIST":
		ofType := outputType.OfType
		intermediateOutputType := Type{
			Kind:   ofType.Kind,
			Name:   ofType.Name,
			OfType: ofType.OfType,
		}
		return graphql.NewList(GetGraphQLOutputType(intermediateOutputType, definitions))

	// unknown type, panic
	default:
		// TODO: not sure what to do here
		outputName, _ := outputType.Name.(string)
		panic(errUnknownArgumentType(outputName))
	}
}
