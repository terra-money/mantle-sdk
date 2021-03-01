package proxy_resolver

func noop() {

}

// package proxy_resolver
//
// import (
// 	"context"
// 	"fmt"
// 	"github.com/graphql-go/graphql"
// 	"github.com/stretchr/testify/assert"
// 	"testing"
// )
//
// func TestProxyResolverContext(t *testing.T) {
// 	var responseCb ProxyResolverResponseCallback = func(query []byte) map[string]interface{} {
// 		assert.Equal(
// 			t,
// 			"query{l1(TxHash:\"8DEF870E21E3A619F5E3EAF29E7452BC5044BA61C8D0FBF414B408363B0D6B2D\",){l11,l12,},l2,l3{l31,},}",
// 			string(query),
// 		)
//
// 		return map[string]interface{}{
// 			"l1": map[string]interface{}{
// 				"l11": "l11response",
// 				"l12": "l12response",
// 			},
// 			"l2": "l2response",
// 			"l3": map[string]interface{}{
// 				"l31": "l31response",
// 			},
// 		}
// 	}
// 	prc := NewProxyResolverContext(responseCb)
//
// 	l1 := prc.CreateSubtree("l1", map[string]interface{}{
// 		"TxHash": "8DEF870E21E3A619F5E3EAF29E7452BC5044BA61C8D0FBF414B408363B0D6B2D",
// 	})
// 	_ = prc.CreateSubtree("l2", nil)
// 	l3 := prc.CreateSubtree("l3", nil)
//
// 	_ = l1.CreateSubtree("l11", nil)
// 	_ = l1.CreateSubtree("l12", nil)
// 	l31 := l3.CreateSubtree("l31", nil)
//
// 	// test query generation
//
// 	l31data, l31Error := l31.Resolve()
//
// 	l3ResponseDeserialized, ok := l31data.(string)
// 	assert.Nil(t, l31Error)
// 	assert.True(t, ok)
// 	assert.Equal(t, "l31response", l3ResponseDeserialized)
//
// }
//
// func TestRecursiveProxyResolverInGraphQL(t *testing.T) {
// 	// test struct
// 	// query {
// 	// 	l1 {
// 	// 		l11,
// 	// 		l12,
// 	// 		l13 {
// 	// 			l131
// 	// 			l132
// 	// 		},
// 	// 	},
// 	// 	l2,
// 	// 	l3,
// 	// }
//
// 	// create matching graphql resolver
// 	fields := graphql.Fields{
// 		"l1": &graphql.Field{
// 			Name: "l1",
// 			Args: graphql.FieldConfigArgument{
// 				"arg": &graphql.ArgumentConfig{
// 					Type:         graphql.Int,
// 					DefaultValue: 1,
// 				},
// 			},
// 			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 				// is root, create prc
// 				// is root
// 				prc := p.Context.Value(ProxyResolverContextKey).(*ProxyResolverContext)
// 				return prc.CreateSubtree("l1", p.Args), nil
// 			},
// 			Type: graphql.NewObject(graphql.ObjectConfig{
// 				Name: "l1",
// 				Fields: graphql.Fields{
// 					"l11": &graphql.Field{
// 						Name: "l11",
// 						Type: graphql.String,
// 						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 							prc := registerSubtreeIfSourceIsPRC(p.Source, "l11", p.Args)
// 							return registerLeafResolver(prc)
// 						},
// 					},
// 					"l12": &graphql.Field{
// 						Name: "l12",
// 						Type: graphql.String,
// 						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 							prc := registerSubtreeIfSourceIsPRC(p.Source, "l12", p.Args)
// 							return registerLeafResolver(prc)
// 						},
// 					},
// 					"l13": &graphql.Field{
// 						Name: "l13",
// 						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 							return registerSubtreeIfSourceIsPRC(p.Source, "l13", p.Args), nil
// 						},
// 						Type: graphql.NewObject(graphql.ObjectConfig{
// 							Name: "l13",
// 							Fields: graphql.Fields{
// 								"l131": &graphql.Field{
// 									Name: "l131",
// 									Type: graphql.String,
// 									Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 										prc := registerSubtreeIfSourceIsPRC(p.Source, "l131", p.Args)
// 										return registerLeafResolver(prc)
// 									},
// 								},
// 								"l132": &graphql.Field{
// 									Name: "l132",
// 									Type: graphql.String,
// 									Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 										prc := registerSubtreeIfSourceIsPRC(p.Source, "l132", p.Args)
// 										return registerLeafResolver(prc)
// 									},
// 								},
// 							},
// 						}),
// 					},
// 				},
// 			}),
// 		},
// 		"l2": &graphql.Field{
// 			Name: "l2",
// 			Type: graphql.Int,
// 			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 				// is root, create prc
// 				// is root
// 				prc := p.Context.Value(ProxyResolverContextKey).(*ProxyResolverContext)
// 				return registerLeafResolver(prc)
// 			},
// 		},
// 		"l3": &graphql.Field{
// 			Name: "l3",
// 			Type: graphql.String,
// 			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 				// is root, create prc
// 				// is root
// 				rootPrc := p.Context.Value(ProxyResolverContextKey).(*ProxyResolverContext)
// 				prc := rootPrc.CreateSubtree("l3", p.Args)
// 				return registerLeafResolver(prc)
// 			},
// 		},
// 		"l4": &graphql.Field{
// 			Name: "l4",
// 			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 				// is root, create prc
// 				// is root
// 				rootPrc := p.Context.Value(ProxyResolverContextKey).(*ProxyResolverContext)
// 				return rootPrc.CreateSubtree("l4", p.Args), nil
// 			},
// 			Type: graphql.NewObject(graphql.ObjectConfig{
// 				Name: "l4",
// 				Fields: graphql.Fields{
// 					"l41": &graphql.Field{
// 						Name: "l41",
// 						Type: graphql.String,
// 						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 							prc := registerSubtreeIfSourceIsPRC(p.Source, "l41", p.Args)
// 							return registerLeafResolver(prc)
// 						},
// 					},
// 					"l42": &graphql.Field{
// 						Name: "l42",
// 						Type: graphql.String,
// 						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 							prc := registerSubtreeIfSourceIsPRC(p.Source, "l42", p.Args)
// 							return registerLeafResolver(prc)
// 						},
// 					},
// 					"l43": &graphql.Field{
// 						Name: "l43",
// 						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 							return registerSubtreeIfSourceIsPRC(p.Source, "l43", p.Args), nil
// 						},
// 						Type: graphql.NewObject(graphql.ObjectConfig{
// 							Name: "l43",
// 							Fields: graphql.Fields{
// 								"l431": &graphql.Field{
// 									Name: "l431",
// 									Type: graphql.String,
// 									Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 										prc := registerSubtreeIfSourceIsPRC(p.Source, "l431", p.Args)
// 										return registerLeafResolver(prc)
// 									},
// 								},
// 								"l432": &graphql.Field{
// 									Name: "l432",
// 									Type: graphql.String,
// 									Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 										prc := registerSubtreeIfSourceIsPRC(p.Source, "l432", p.Args)
// 										return registerLeafResolver(prc)
// 									},
// 								},
// 							},
// 						}),
// 					},
// 				},
// 			}),
// 		},
// 	}
//
// 	// handle http request
// 	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
// 		Query: graphql.NewObject(graphql.ObjectConfig{
// 			Name:   "rootQuery",
// 			Fields: fields,
// 		}),
// 	})
//
// 	// make stub response. use this to check against the resolved response
// 	// NOTE: since leaf node resolvers are run in thunks,
// 	//		 we don't have guarantee over the execution order.
// 	//		 hence requestString would NEVER be the same as the input request string
// 	//		 check normality w/ result
// 	stubResponse := map[string]interface{}{
// 		"l1": map[string]interface{}{
// 			"l12": "l12",
// 			"l13": map[string]interface{}{
// 				"l131": "l131",
// 				"l132": "l132",
// 			},
// 		},
// 		"l3": "l3",
// 		"l4": map[string]interface{}{
// 			"l43": map[string]interface{}{
// 				"l431": "l431",
// 				"l432": "l432", // dangling
// 			},
// 		},
// 	}
// 	rootPrc := NewProxyResolverContext(func(query []byte) map[string]interface{} {
// 		return stubResponse
// 	})
//
// 	result := graphql.Do(graphql.Params{
// 		Schema: schema,
//
// 		// prop-down root PRC for remote mantle sync job
// 		Context: context.WithValue(context.Background(), ProxyResolverContextKey, rootPrc),
//
// 		// request string as defined for this tc
// 		RequestString: "query{l1(arg:15125,){l12,l13{l132,l131,},},l3,l4{l43{l431,},},}",
// 	})
//
// 	assert.ObjectsAreEqual(stubResponse, result.Data)
// }
//
// func registerSubtreeIfSourceIsPRC(source interface{}, subtreeName string, arguments map[string]interface{}) *ProxyResolverContext {
// 	sourcePRC, ok := source.(*ProxyResolverContext)
// 	if !ok {
// 		panic(fmt.Errorf("%s: source is not is not prc", subtreeName))
// 	}
//
// 	return sourcePRC.CreateSubtree(subtreeName, arguments)
// }
//
// func registerLeafResolver(prc *ProxyResolverContext) (func() (interface{}, error), error) {
// 	return func() (interface{}, error) {
// 		return prc.Resolve()
// 	}, nil
// }
