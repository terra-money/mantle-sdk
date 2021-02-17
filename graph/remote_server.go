package graph

import (
	"context"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/terra-project/mantle-sdk/depsresolver"
	"github.com/terra-project/mantle-sdk/querier"
	"github.com/terra-project/mantle-sdk/types"
	"log"
	"net/http"
)

type RemoteGraphQLInstance struct {
	*GraphQLInstance
	baseMantleEndpoint string
}

func NewRemoteGraphQLInstance(
	depsResolver depsresolver.DepsResolver,
	querier querier.Querier,
	baseMantleEndpoint string,
	schemabuilders ...SchemaBuilder,
) *RemoteGraphQLInstance {
	return &RemoteGraphQLInstance{
		GraphQLInstance:    NewGraphQLInstance(depsResolver, querier, schemabuilders...),
		baseMantleEndpoint: baseMantleEndpoint,
	}
}

// TODO: reimplement me without using graphql-go/handler
func (server *RemoteGraphQLInstance) ServeHTTP(port int) {
	h := handler.New(&handler.Config{
		Schema: &server.schema,
		RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
			return server.depsResolver.GetState()
		},
		Pretty:     true,
		Playground: true,
	})

	http.Handle("/status", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ContextHandler(
			NewGraphContext().
				WithProxyResolverContext(server.baseMantleEndpoint).
				WithDepsResolver(server.depsResolver).
				WithQuerier(server.querier).
				WithImmediateResolveFlag(true).
				ToContext(),
			w,
			r,
		)
	}))
	http.ListenAndServe(fmt.Sprintf(":%d", int(port)), nil)
}

func (server *RemoteGraphQLInstance) QueryInternal(
	gqlQuery string,
	variables types.GraphQLParams,
	dependencies []types.Model,
) types.GraphQLResult {
	log.Printf("[graphql] Query\tq=%s,v=%v", gqlQuery, variables)
	params := graphql.Params{
		Schema:         server.schema,
		RequestString:  gqlQuery,
		VariableValues: variables,
		Context: NewGraphContext().
			WithProxyResolverContext(server.baseMantleEndpoint).
			WithDepsResolver(server.depsResolver).
			WithQuerier(server.querier).
			WithImmediateResolveFlag(false).
			WithDependencies(dependencies).
			ToContext(),
	}

	// unresolved dependency are to be handled in resolver functions
	return InternalGQLRun(params)
}
