package indexer

import (
	"fmt"
	"testing"

	"github.com/terra-project/mantle-sdk/graph"
	"github.com/terra-project/mantle-sdk/types"
)

type BaseState struct {
	Foo string
	Bar string
	Idx int
}

func createGraphQLInstance() *graph.GraphQLInstance {
	var baseState = BaseState{
		Foo: "foo",
		Bar: "bar",
		Idx: 1,
	}

	// create gql instance for query resolution
	gqlInstance := graph.NewGraphQLInstance(baseState)
	gqlInstance.UpdateState("BaseState", baseState)

	return gqlInstance
}

func DBG(t *testing.T) {
	gqlInstance := createGraphQLInstance()
	gqlInstance.ServeHTTP()
}

func TestRunIndexerRound(t *testing.T) {
	gqlInstance := createGraphQLInstance()

	baseInstance := NewIndexerBaseInstance(
		[]types.Indexer{
			testIndexer1,
			testIndexer2,
		},
		gqlInstance.ResolveQuery,
		gqlInstance.Commit,
	)

	baseInstance.RunIndexerRound()

	states := gqlInstance.ExportStates()
	fmt.Println(states)
}

//////////////////////////////////////////
type TestEntity1 struct {
	Field1 string
	Field2 int
}
type TestEntity2 struct {
	Field1 string
	Field2 TestEntity2SubStruct
}
type TestEntity2SubStruct struct {
	Data string
}

/// TestIndexer1
type Test1Query struct {
	BaseState BaseState
}

func testIndexer1(query types.Query, commit types.Commit) {
	request := Test1Query{}
	query(&request, nil)

	entity := TestEntity1{
		Field1: request.BaseState.Foo,
		Field2: request.BaseState.Idx,
	}
	commit(entity)
}

// TestIndexer2
type TestQuery2 struct {
	BaseState BaseState
}

func testIndexer2(query types.Query, commit types.Commit) {
	request := TestQuery2{}
	query(&request, nil)

	entity := TestEntity2{
		Field1: request.BaseState.Foo,
		Field2: TestEntity2SubStruct{
			Data: request.BaseState.Bar,
		},
	}
	commit(entity)
}
