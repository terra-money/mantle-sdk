package generate

import (
	"fmt"
	"reflect"
	"testing"
)

type TestEntity struct {
	CompositeField  CompositeField
	CompositeField2 Composite2
	TopLevelField   int
}

type CompositeField struct {
	Field1 string
	Field2 string
}

type Composite2 struct {
	NestedField NestedField
}

type NestedField struct {
	Field1 string
	Field2 string
}

func TestGenerateGraphResolver(t *testing.T) {
	name, rootFields := GenerateGraphResolver(reflect.TypeOf((*TestEntity)(nil)))

	if name != "TestEntity" {
		panic("name is not correct")
	}

	fmt.Println(rootFields)

	// schema, err := graphql.NewSchema(graphql.SchemaConfig{
	// 	Query: rootQuery,
	// })

	// if err != nil {
	// 	panic(err)
	// 	t.Fail()
	// }

	// h := handler.New(&handler.Config{
	// 	Schema:     &schema,
	// 	Pretty:     true,
	// 	Playground: true,
	// })

	// http.Handle("/", h)
	// http.ListenAndServe(":3030", nil)

}
