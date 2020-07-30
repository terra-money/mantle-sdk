package graph

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/terra-project/mantle/graph/depsresolver"
	"github.com/terra-project/mantle/querier"
	"github.com/terra-project/mantle/types"
	"github.com/terra-project/mantle/utils"
)

type SchemaBuilder func(fields *graphql.Fields) error

type GraphQLInstance struct {
	schema       graphql.Schema
	depsResolver depsresolver.DepsResolver
	querier      querier.Querier
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

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ContextHandler(
			server.prepareResolverContext(nil),
			w,
			r,
		)
	}))
	http.ListenAndServe(fmt.Sprintf(":%d", int(port)), nil)
}

func (server *GraphQLInstance) UpdateState(data interface{}) {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Struct {
		panic("Non struct type entity is provided to GraphQLInstance.UpdateState")
	}

	server.depsResolver.SetPredefinedState(data)
}

func (server *GraphQLInstance) ResolveQuery(
	gqlQuery string,
	variables types.GraphQLParams,
	dependencies []types.ModelType,
) *graphql.Result {

	log.Printf("[graphql] ResolveQuery\tq=%s,v=%v", gqlQuery, variables)

	params := graphql.Params{
		Schema:         server.schema,
		RequestString:  gqlQuery,
		VariableValues: variables,
		Context:        server.prepareResolverContext(dependencies),
	}

	// unresolved dependency are to be handled in resolver functions
	return graphql.Do(params)
}

// Commit persists indexer outputs in memory.
func (server *GraphQLInstance) Commit(entity interface{}) {
	server.depsResolver.Emit(entity)
}

func (server *GraphQLInstance) ExportStates() []interface{} {
	entities := make([]interface{}, 0)
	for _, entity := range server.depsResolver.GetState() {
		if entity == nil {
			continue
		}
		entities = append(entities, entity)
	}

	return entities
}

func (server *GraphQLInstance) prepareResolverContext(dependencies []types.ModelType) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.DepsResolverKey, server.depsResolver)
	ctx = context.WithValue(ctx, utils.QuerierKey, server.querier)

	// dependencies are only taken care of when running indexers
	if dependencies != nil {
		dependencyNames := make(utils.DependenciesKeyType)
		for _, dependencyModel := range dependencies {
			dependencyNames[dependencyModel.Name()] = true
		}
		ctx = context.WithValue(ctx, utils.DependenciesKey, dependencyNames)
	}

	return ctx
}

func buildSchema(schemabuilders ...SchemaBuilder) graphql.Schema {
	rootFields := &graphql.Fields{}

	for _, builder := range schemabuilders {
		err := builder(rootFields)

		if err != nil {
			panic(err)
		}
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootQuery",
			Fields: *rootFields,
		}),
		// Subscription: graphql.NewObject(graphql.ObjectConfig{
		// 	Name:   "RootSubscription",
		// 	Fields: *rootFields,
		// }),
	})

	if err != nil {
		panic(err)
	}

	return schema
}
