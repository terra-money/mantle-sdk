package proxy_resolver

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle-sdk/graph/graph_types"
)

func GetGraphQLInputType(argType *Input, definitions Definitions) graphql.Input {

	switch argType.Kind {
	// graphql-native graph_types + custom graph_types
	case "SCALAR":
		switch argType.Name {
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
		// check cosmos-scalar map
		default:
			cosmosScalar := graph_types.GetCosmosScalarByName(argType.Name)

			// if name is unknown, mantle can't handle it. panic here
			if cosmosScalar == nil {
				panic(errUnknownCustomScalar(argType.Name))
			}

			return cosmosScalar
		}

	// get the matching OBJECT type of definitions, reconstruct them into graphql type
	case "OBJECT":
		panic(fmt.Errorf("object argument is not supported yet"))

	case "LIST":
		ofType := argType.OfType
		return graphql.NewList(GetGraphQLInputType(ofType, definitions))

	case "ENUM":
		underlyingType := mustGetUnderlyingInputType(argType, definitions)

		if pervasive := graph_types.GetMantlePervasiveByName(underlyingType.Name); pervasive != nil {
			return pervasive
		}

		enumValues := underlyingType.EnumValues
		enumValueConfigMap := graphql.EnumValueConfigMap{}

		if enumCache := getTypeFromCache(underlyingType.Name); enumCache != nil {
			return enumCache
		}

		for _, enumValue := range enumValues {
			enumValueConfigMap[enumValue.Name] = &graphql.EnumValueConfig{
				Value:             enumValue.Name,
				DeprecationReason: enumValue.DeprecationReason,
				Description:       enumValue.Description,
			}
		}

		enum := graphql.NewEnum(graphql.EnumConfig{
			Name:        underlyingType.Name,
			Values:      enumValueConfigMap,
			Description: underlyingType.Description,
		})

		setTypeInCache(underlyingType.Name, enum)

		return enum

	// unknown type, panic
	default:
		panic(errUnknownArgumentType(argType.Name))
	}
}

func mustGetUnderlyingInputType(argType *Input, definitions Definitions) TypeDescriptor {
	underlyingType, ok := definitions[argType.Name]

	if !ok {
		panic(fmt.Errorf("underlying type %s does not exist in definitions map", argType.Name))
	}

	return underlyingType
}
