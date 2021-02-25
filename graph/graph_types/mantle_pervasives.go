package graph_types

import (
	"github.com/graphql-go/graphql"
)

type JSONScalar string // TODO: must go
var mantlePervasives = map[string]graphql.Output{
	"Order": Order,
}

func GetMantlePervasiveByName(typeName string) graphql.Output {
	return mantlePervasives[typeName]
}

var Order = graphql.NewEnum(graphql.EnumConfig{
	Name: "Order",
	Values: graphql.EnumValueConfigMap{
		"ASC": &graphql.EnumValueConfig{
			Value:       "ASC",
			Description: "Ascending order",
		},
		"DESC": &graphql.EnumValueConfig{
			Value:       "DESC",
			Description: "Ascending order",
		},
	},
	Description: "",
})
