package generate

import (
	"math/big"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// here defined are the list of known cosmos-sdk scalars.
// when generating gql query or schema, we need to check if each
// leaf value implements these interfaces,
// and if so, the returned gql scalar type __must__ be used.

func IsCosmosScalar(leafType reflect.Type) (*graphql.Scalar, bool) {
	for _, test := range lists {
		if reflect.PtrTo(leafType).Implements(test.Type) {
			return test.Scalar, true
		}
	}

	return nil, false
}

type ScalarGeneratorPair struct {
	Type   reflect.Type
	Scalar *graphql.Scalar
}

var lists = []ScalarGeneratorPair{
	{
		Type: reflect.TypeOf((*interface{ BigInt() *big.Int })(nil)).Elem(),
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
}
