package types

import (
	"github.com/graphql-go/graphql"
)

var MantleScalars = []graphql.Output{
	Order,
}

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
