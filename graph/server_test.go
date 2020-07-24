package graph_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terra-project/mantle/db/badger"
	"github.com/terra-project/mantle/db/kvindex"
	. "github.com/terra-project/mantle/graph"
	"github.com/terra-project/mantle/graph/depsresolver"
	"github.com/terra-project/mantle/graph/schemabuilders"
	"github.com/terra-project/mantle/querier"
	testutil "github.com/terra-project/mantle/test"
)

type TestChild1 struct {
	TestChild1Field1 int
	TestChild1Field2 string
}

type TestChild2 struct {
	TestChild2Field1 int
	TestChild2Field2 TestChild2Field2
}

type TestChild2Field2 struct {
	TestChild2Field2Field1 string `mantle:"index=test"`
}

type TestStruct struct {
	TestChild1 TestChild1
	TestChild2 TestChild2
	TestField  int
}

type TestQuery struct {
	TestStruct TestStruct
}

var initialState TestStruct = TestStruct{
	TestChild1: TestChild1{
		TestChild1Field1: 1,
		TestChild1Field2: "1",
	},
	TestChild2: TestChild2{
		TestChild2Field1: 1,
		TestChild2Field2: TestChild2Field2{
			TestChild2Field2Field1: "Hello World",
		},
	},
	TestField: 1,
}

func createTestServer(
	initialState interface{},
) *GraphQLInstance {
	app := testutil.CreateTestApp()
	db := badger.NewBadgerDB("")
	kvi := kvindex.NewKVIndex(reflect.TypeOf(initialState))
	kvindexMap := kvindex.NewKVIndexMap(kvi)
	depsresolver := depsresolver.NewDepsResolver()
	depsresolver.SetPredefinedState(initialState)

	querier := querier.NewQuerier(db, kvindexMap)

	instance := NewGraphQLInstance(
		depsresolver,
		querier,
		schemabuilders.CreateABCIStubSchemaBuilder(app.GetApp()),
		schemabuilders.CreateModelSchemaBuilder(reflect.TypeOf(TestStruct{})),
		schemabuilders.CreateListSchemaBuilder(),
	)

	// commit initial state so other test indexers can see it
	instance.Commit(initialState)

	return instance
}

func TestServeHTTP(t *testing.T) {
	server := createTestServer(initialState)
	server.ServeHTTP(3030)
}

func TestResolveQuery(t *testing.T) {
	server := createTestServer(initialState)
	qs := `
		query {
			TestStruct {
				TestChild1 {
					TestChild1Field1
					TestChild1Field2
				}
				TestChild2 {
					TestChild2Field1
					TestChild2Field2 {
						TestChild2Field2Field1
					}
				}
				TestField
			}
		}
	`

	result := server.ResolveQuery(qs, nil)

	if result.HasErrors() {
		for _, err := range result.Errors {
			fmt.Println(err)
		}
		t.FailNow()
	}

	res := TestQuery{}
	err := UnmarshalGraphQLResult(result, &res)
	if err != nil {
		t.Error(err)
	}

	// assert
	assert.Equal(t, initialState.TestChild1.TestChild1Field1, res.TestStruct.TestChild1.TestChild1Field1)
	assert.Equal(t, initialState.TestChild1.TestChild1Field2, res.TestStruct.TestChild1.TestChild1Field2)
	assert.Equal(t, initialState.TestChild2.TestChild2Field1, res.TestStruct.TestChild2.TestChild2Field1)
	assert.Equal(t, initialState.TestChild2.TestChild2Field2.TestChild2Field2Field1, res.TestStruct.TestChild2.TestChild2Field2.TestChild2Field2Field1)

	// test after updating state
	state2 := TestStruct{
		TestChild1: TestChild1{
			TestChild1Field1: 100,
			TestChild1Field2: "100",
		},
		TestChild2: TestChild2{
			TestChild2Field1: 100,
			TestChild2Field2: TestChild2Field2{
				TestChild2Field2Field1: "World Hello",
			},
		},
		TestField: 1,
	}
	server.UpdateState(state2)

	res2 := TestQuery{}
	result2 := server.ResolveQuery(qs, nil)
	if err = UnmarshalGraphQLResult(result2, &res2); err != nil {
		t.Error(err)
	}

	assert.Equal(t, state2.TestChild1.TestChild1Field1, res2.TestStruct.TestChild1.TestChild1Field1)
	assert.Equal(t, state2.TestChild1.TestChild1Field2, res2.TestStruct.TestChild1.TestChild1Field2)
	assert.Equal(t, state2.TestChild2.TestChild2Field1, res2.TestStruct.TestChild2.TestChild2Field1)
	assert.Equal(t, state2.TestChild2.TestChild2Field2.TestChild2Field2Field1, res2.TestStruct.TestChild2.TestChild2Field2.TestChild2Field2Field1)

	// state export
	exported := server.ExportStates()
	if exported == nil {
		t.Errorf("TestStruct is not stored in current session")
	}

	if _, ok := exported["TestStruct"].(TestStruct); !ok {
		t.Errorf("Exported state type assertion failed")
	}

	assert.Equal(t, state2.TestChild1.TestChild1Field1, exported["TestStruct"].(TestStruct).TestChild1.TestChild1Field1)
	assert.Equal(t, state2.TestChild1.TestChild1Field2, exported["TestStruct"].(TestStruct).TestChild1.TestChild1Field2)
	assert.Equal(t, state2.TestChild2.TestChild2Field1, exported["TestStruct"].(TestStruct).TestChild2.TestChild2Field1)
	assert.Equal(t, state2.TestChild2.TestChild2Field2.TestChild2Field2Field1, exported["TestStruct"].(TestStruct).TestChild2.TestChild2Field2.TestChild2Field2Field1)
}

func BenchmarkCycle(b *testing.B) {
	server := createTestServer(initialState)
	qs := `
		query {
			TestStruct {
				TestChild1 {
					TestChild1Field1
					TestChild1Field2
				}
				TestChild2 {
					TestChild2Field1
					TestChild2Field2 {
						TestChild2Field2Field1
					}
				}
				TestField
			}
		}
	`

	for i := 0; i < b.N; i++ {
		result := server.ResolveQuery(qs, nil)
		res := TestQuery{}
		UnmarshalGraphQLResult(result, &res)
	}
}

func TestQuerierRouter(t *testing.T) {

	// to initialize module manager
	server := createTestServer(initialState)

	server.ServeHTTP(3030)

}
