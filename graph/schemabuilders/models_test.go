package schemabuilders

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

type TestModel1 struct {
	Foo string
	Bar struct {
		Hello string `mantle:"index=hello"`
	}
}

type TestModel2 struct {
	Foo uint64
	Bar struct {
		Nested struct {
			Timestamp uint64 `mantle:"index=timestamp"`
		}
	}
}

func TestCreateModelSchemaBuilder(t *testing.T) {
	builder := CreateModelSchemaBuilder(
		reflect.TypeOf((*TestModel1)(nil)),
		reflect.TypeOf((*TestModel2)(nil)),
	)

	fields := graphql.Fields{}
	builder(&fields)

	// can't really test beyond this
	// actual acceptance shall be test in server_test
	assert.Equal(t, len(fields), 2)

	// test whether indexes are registered as graphql args
	index1, index1Exists := fields["TestModel1"].Args["hello"]
	assert.Equal(t, index1Exists, true)
	assert.Equal(t, index1.Type, graphql.String)

	index2, index2Exists := fields["TestModel2"].Args["timestamp"]
	assert.Equal(t, index2Exists, true)
	assert.Equal(t, index2.Type, graphql.Int)
}
