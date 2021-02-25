package graph

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/terra-project/mantle-sdk/depsresolver"
	"github.com/terra-project/mantle-sdk/querier"
	"github.com/terra-project/mantle-sdk/types"
)

type SchemaBuilder func(fields *graphql.Fields) error

type GraphQLInstance struct {
	schema             graphql.Schema
	depsResolver       depsresolver.DepsResolver
	querier            querier.Querier
	baseMantleEndpoint string
}

func NewGraphQLInstance(
	depsResolver depsresolver.DepsResolver,
	querier querier.Querier,
	schemabuilders ...SchemaBuilder,
) *GraphQLInstance {
	return &GraphQLInstance{
		depsResolver: depsResolver,
		querier:      querier,
		schema:       buildSchema(schemabuilders...),
	}
}

// TODO: reimplement me without using graphql-go/handler
func (server *GraphQLInstance) ServeHTTP(port int) {
	h := handler.New(&handler.Config{
		Schema: &server.schema,
		RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
			return server.depsResolver.GetState()
		},
		Pretty:     true,
		Playground: true,
	})

	mux := http.NewServeMux()
	mux.Handle("/status", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ContextHandler(
			NewGraphContext().
				WithDepsResolver(server.depsResolver).
				WithQuerier(server.querier).
				WithImmediateResolveFlag(true).
				ToContext(),
			w,
			r,
		)
	}))

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func (server *GraphQLInstance) QueryInternal(
	gqlQuery string,
	variables types.GraphQLParams,
	dependencies []types.Model,
) *graphql.Result {
	log.Printf("[graphql] Query\tq=%s,v=%v", gqlQuery, variables)
	params := graphql.Params{
		Schema:         server.schema,
		RequestString:  gqlQuery,
		VariableValues: variables,
		Context: NewGraphContext().
			WithDepsResolver(server.depsResolver).
			WithQuerier(server.querier).
			WithImmediateResolveFlag(false).
			WithDependencies(dependencies).
			ToContext(),
	}

	// unresolved dependency are to be handled in resolver functions
	return InternalGQLRun(params)
}

func (server *GraphQLInstance) UpdateState(data interface{}) {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Struct {
		panic("Non struct type entity is provided to GraphQLInstance.UpdateState")
	}

	server.depsResolver.SetPredefinedState(data)
}

// Commit persists indexer outputs in memory.
func (server *GraphQLInstance) Commit(entity interface{}) error {
	return server.depsResolver.Emit(entity)
}
