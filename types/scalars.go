package types

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
)

var MantleScalars = []graphql.Output{
	Order,
}

type JSONScalar string
type JSONScalarCustomMarshaler func(originalData interface{}, data []byte) []byte

func NewJSONScalar(data interface{}, customMarshaller JSONScalarCustomMarshaler) JSONScalar {
	bz, bzErr := json.Marshal(data)
	if bzErr != nil {
		panic(bzErr)
	}

	if customMarshaller != nil {
		bz = customMarshaller(data, bz)
	}

	return JSONScalar(bz)
}
func (scalar *JSONScalar) IsJsonScalar() {}

var Order = graphql.NewEnum(graphql.EnumConfig{
	Name: "Order",
	Values: graphql.EnumValueConfigMap{
		"ASC": &graphql.EnumValueConfig{
			Value: "ASC",
			//DeprecationReason: "",
			Description: "Ascending order",
		},
		"DESC": &graphql.EnumValueConfig{
			Value: "DESC",
			//DeprecationReason: "",
			Description: "Ascending order",
		},
	},
	Description: "",
})
