package scalars

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var scalars = make(map[reflect.Type]ScalarGeneratorPair)

// here defined are the list of known cosmos-sdk scalars.
// when generating gql query or schema, we need to check if each
// leaf value implements these interfaces,
// and if so, the returned gql scalar type __must__ be used.
func IsCosmosScalar(leafType reflect.Type) (*graphql.Scalar, bool) {
	for _, test := range lists {
		if test.Check(leafType, test.Scalar) {
			return test.Scalar, true
		}
	}

	return nil, false
}

// GetCosmosScalarByName returns cosmos scalar by name
func GetCosmosScalarByName(scalarName string) *graphql.Scalar {
	for _, scalar := range lists {
		if scalar.Scalar.Name() == scalarName {
			return scalar.Scalar
		}
	}

	// errorneous case where scalar could not be found
	return nil
}

type ScalarGeneratorPair struct {
	Check  func(target reflect.Type, scalar *graphql.Scalar) bool
	Scalar *graphql.Scalar
}

var lists = []ScalarGeneratorPair{
	// bigint
	{
		Check: func(target reflect.Type, scalar *graphql.Scalar) bool {
			if target == reflect.TypeOf((*types.Int)(nil)).Elem() {
				return true
			} else {
				return false
			}
		},
		Scalar: graphql.NewScalar(graphql.ScalarConfig{
			Name:        "BigInt",
			Description: "BigInt scalar type represents cosmos-sdk specific big int implementation",
			Serialize: func(value interface{}) interface{} {
				return value
			},
			ParseValue: func(value interface{}) interface{} {
				return value
			},
			ParseLiteral: func(valueAST ast.Value) interface{} {
				return valueAST.GetValue()
			},
		}),
	},
	// StdTx/Msg
	{
		Check: func(target reflect.Type, scalar *graphql.Scalar) bool {
			t := reflect.TypeOf((*types.Msg)(nil)).Elem()
			if target.Implements(t) {
				return true
			}
			return false
		},
		Scalar: graphql.NewScalar(graphql.ScalarConfig{
			Name:        "Msg",
			Description: "Standard cosmos-sdk tx messages",
			Serialize: func(value interface{}) interface{} {
				msgInJsonBytes, _ := json.Marshal(value)
				return string(msgInJsonBytes)
			},
			ParseValue: func(value interface{}) interface{} {
				return value
			},
			ParseLiteral: func(valueAST ast.Value) interface{} {
				return valueAST.GetValue()
			},
		}),
	},
	// byte buffer
	{
		Check: func(target reflect.Type, scalar *graphql.Scalar) bool {
			if target == reflect.TypeOf(([]byte)(nil)) {
				return true
			}
			return false
		},
		Scalar: graphql.NewScalar(graphql.ScalarConfig{
			Name:        "Buffer",
			Description: "[]byte serialized as string",
			Serialize: func(value interface{}) interface{} {
				return string(value.([]byte))
			},
			ParseValue: func(value interface{}) interface{} {
				return value
			},
			ParseLiteral: func(valueAST ast.Value) interface{} {
				return valueAST.GetValue()
			},
		}),
	},
	// time.Time
	{
		Check: func(target reflect.Type, scalar *graphql.Scalar) bool {
			if target == reflect.TypeOf((*time.Time)(nil)).Elem() {
				return true
			}
			return false
		},
		Scalar: graphql.NewScalar(graphql.ScalarConfig{
			Name:        "Time",
			Description: "golang's time.Time",
			Serialize: func(value interface{}) interface{} {
				return value.(time.Time).String()
			},
			ParseValue: func(value interface{}) interface{} {
				return value.(time.Time)
			},
			ParseLiteral: func(valueAST ast.Value) interface{} {
				return valueAST.GetValue()
			},
		}),
	},
}
