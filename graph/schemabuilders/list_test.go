package schemabuilders

import (
	"testing"

	p "github.com/gertd/go-pluralize"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func TestListSchemaBuilder(t *testing.T) {
	var pluralize = p.NewClient().Plural
	schemabuilder := CreateListSchemaBuilder()

	// failing test
	// list object creation is disallowed when Args are empty
	// why? becuase lists are basically for range search,
	// the original (singular) object must have some sort of args (indexes)
	func() {
		testFields := &graphql.Fields{
			"FailingTest": &graphql.Field{
				Name:    "FailingTest",
				Type:    graphql.Boolean,
				Resolve: graphql.DefaultResolveFn,
			},
		}

		err := schemabuilder(testFields)
		assert.NotNil(t, err)
	}()

	// passing test
	// note that full test with resolver is done in server_test
	func() {
		fieldType := graphql.Int
		testFields := &graphql.Fields{
			"Test": &graphql.Field{
				Name: "Test",
				Args: graphql.FieldConfigArgument{
					"someIndex": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Type:    fieldType,
				Resolve: graphql.DefaultResolveFn,
			},
		}


		err := schemabuilder(testFields)

		assert.Nil(t, err)

		pluralFieldName := pluralize("Test")
		listTarget, ok := (*testFields)[pluralFieldName]

		assert.Equal(t, ok, true, "check if pluralized field exists")
		assert.Equal(t, graphql.NewList(fieldType), listTarget.Type, "list field type is slice of original")
		_, ok = listTarget.Args["someIndex"]
		assert.Equal(t, ok, true, "singular argument preserved")
		rangeArgs, ok := listTarget.Args["someIndex_range"]
		assert.Equal(t, ok, true, "ranged argument exists")
		assert.Equal(t, rangeArgs.Type, graphql.NewList(graphql.Int), "ranged argument type is slice of original")
	}()
}
