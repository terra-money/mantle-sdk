package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/terra-project/mantle-sdk/depsresolver"
	"github.com/terra-project/mantle-sdk/graph/proxy_resolver"
	"github.com/terra-project/mantle-sdk/querier"
	"github.com/terra-project/mantle-sdk/types"
	"github.com/terra-project/mantle-sdk/utils"
)

type graphContext struct {
	ctx context.Context
}

func NewGraphContext() graphContext {
	return graphContext{ctx: context.Background()}
}

func (g graphContext) WithDepsResolver(resolver depsresolver.DepsResolver) graphContext {
	return graphContext{
		ctx: context.WithValue(g.ctx, utils.DepsResolverKey, resolver),
	}
}

func (g graphContext) WithQuerier(querier querier.Querier) graphContext {
	return graphContext{
		ctx: context.WithValue(g.ctx, utils.QuerierKey, querier),
	}
}

func (g graphContext) WithImmediateResolveFlag(resolveImmediately bool) graphContext {
	return graphContext{
		ctx: context.WithValue(g.ctx, utils.ImmediateResolveFlagKey, resolveImmediately),
	}
}

func (g graphContext) WithDependencies(dependencies []types.Model) graphContext {
	if dependencies != nil {
		dependencyNames := make(utils.DependenciesKeyType)
		for _, dependencyModel := range dependencies {
			dependencyNames[dependencyModel.Name()] = true
		}

		return graphContext{
			ctx: context.WithValue(g.ctx, utils.DependenciesKey, dependencyNames),
		}
	}

	return g
}

func (g graphContext) WithProxyResolverContext(baseMantleEndpoint string) graphContext {
	prc := proxy_resolver.NewProxyResolverContext(func(query []byte) (map[string]interface{}, error) {
		response := CreateRemoteMantleRequest(baseMantleEndpoint, query)
		gqlResponse := new(struct {
			Data   map[string]interface{} `json:"data"`
			Errors *[]struct {
				Message   string `json:"message"`
				Locations []struct {
					Line   uint64 `json:"line"`
					Column uint64 `json:"column"`
				} `json:"locations"`
			} `json:"errors"`
		})

		if err := json.Unmarshal(response, gqlResponse); err != nil {
			panic(err)
		}

		if gqlResponse.Errors != nil {
			errString, err := json.Marshal(gqlResponse.Errors)
			if err != nil {
				panic("should not happen")
			}

			return nil, fmt.Errorf(string(errString))
		}

		return gqlResponse.Data, nil // TODO: error
	})

	return graphContext{
		ctx: context.WithValue(g.ctx, utils.ProxyResolverContextKey, prc),
	}
}

func (g graphContext) ToContext() context.Context {
	return g.ctx
}
