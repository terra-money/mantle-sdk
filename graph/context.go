package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/terra-project/mantle-sdk/depsresolver"
	"github.com/terra-project/mantle-sdk/graph/proxy_resolver"
	"github.com/terra-project/mantle-sdk/querier"
	"github.com/terra-project/mantle-sdk/types"
	"github.com/terra-project/mantle-sdk/utils"
	"log"
	"time"
)

type graphContext struct {
	ctx    context.Context
	height uint64
}

func NewGraphContext() graphContext {
	return graphContext{ctx: context.Background()}
}

const heightKey = "heightKey"

func (g graphContext) WithHeight(height uint64) graphContext {
	return graphContext{
		ctx: context.WithValue(g.ctx, heightKey, height),
	}
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
	var remoteQuerier proxy_resolver.ProxyResolverResponseCallback
	remoteQuerier = func(query []byte) (map[string]interface{}, error) {
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

		lastSyncedHeight, isLastSyncedHeightOk := gqlResponse.Data["__referenceHeight"].(float64)
		if !isLastSyncedHeightOk {
			// should not happen, but just in case
			return nil, fmt.Errorf("invalid LastSyncedHeight from mantle")
		}

		referenceHeight, ok := g.ctx.Value(heightKey).(uint64)
		if !ok {
			return nil, fmt.Errorf("reference height not set")
		}

		var lastSyncedHeightUint64 = uint64(lastSyncedHeight)

		// in such case, retry
		if lastSyncedHeightUint64 < referenceHeight {
			log.Printf("[proxyResolver] invalid height lastSyncedHeight %d != referenceHeight %d, retrying...", uint64(lastSyncedHeight), referenceHeight)
			time.Sleep(200 * time.Millisecond)
			return remoteQuerier(query)
		}

		if gqlResponse.Errors != nil {
			errString, err := json.Marshal(gqlResponse.Errors)
			if err != nil {
				return nil, errors.Wrap(err, "invalid gql errors")
			}

			return nil, errors.Wrap(fmt.Errorf("%s", errString), "base mantle query errored")
		}

		return gqlResponse.Data, nil // TODO: error
	}

	prc := proxy_resolver.NewProxyResolverContext(remoteQuerier)

	// always create LastSyncedHeight, so state snapshot is guaranteed to be at the same height
	prc.CreateSubtree("LastSyncedHeight", nil).WithAlias("__referenceHeight")

	return graphContext{
		ctx: context.WithValue(g.ctx, utils.ProxyResolverContextKey, prc),
	}
}

func (g graphContext) ToContext() context.Context {
	return g.ctx
}
