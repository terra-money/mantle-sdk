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
	// why? becuase lists are basically for @range search,
	// the original (singular) object must have some sort of args (indexes)
	func() {
		testFields := &graphql.Fields{
			"Test": &graphql.Field{
				Name:    "Test",
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
		testFields := &graphql.Fields{
			"Test": &graphql.Field{
				Name: "Test",
				Args: graphql.FieldConfigArgument{
					"someIndex": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Type:    graphql.Boolean,
				Resolve: graphql.DefaultResolveFn,
			},
		}

		err := schemabuilder(testFields)
		assert.Nil(t, err)

		pluralFieldName := pluralize("Test")
		_, ok := (*testFields)[pluralFieldName]

		assert.Equal(t, ok, true)
	}()
}
