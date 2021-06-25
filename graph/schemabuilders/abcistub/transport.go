package abcistub

import (
	"encoding/json"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"net/http"

	client "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	terra "github.com/terra-project/core/app"
	compatlocalclient "github.com/terra-project/mantle-compatibility/localclient"
)

type RoundTripper struct {
	mux   *mux.Router
	cache *lru.Cache
}

func NewRoundTripper(mux *mux.Router, cache *lru.Cache) *RoundTripper {
	return &RoundTripper{
		mux:   mux,
		cache: cache,
	}
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cached, ok := rt.cache.Get(req.URL.String())
	if ok {
		next := cached.(ReaderCloser).Clone()
		return &http.Response{
			StatusCode: 200,
			Body:       next,
		}, nil
	}

	// create a stub response
	out := &ResponseBuffer{}

	// relay to lcd...
	rt.mux.ServeHTTP(out, req)

	ret := make(map[string]interface{})
	json.Unmarshal(out.Finalize(), &ret)

	if err, exists := ret["error"]; exists {
		return nil, errors.New(err.(string))
	}

	// hacky fix: if wasm, result is always string
	// stringify all result

	if req.URL.Path[:5] == "/wasm" {
		if _, exists := ret["result"]; exists {
			resultBytes, resultBytesErr := json.Marshal(ret["result"])
			if resultBytesErr != nil {
				return nil, fmt.Errorf("wasm result stringify failed, err=%s", resultBytesErr)
			}

			ret["result"] = string(resultBytes)
		}

		nextBody, nextBodyErr := json.Marshal(&ret)
		if nextBodyErr != nil {
			return nil, fmt.Errorf("wasm response stringify failed, err=%s", nextBodyErr)
		}

		_, nextBodyWriteErr := out.Write(nextBody)
		if nextBodyWriteErr != nil {
			return nil, fmt.Errorf("wasm response write failed, err=%s", nextBodyWriteErr)
		}
	}

	body := out.Body()

	rt.cache.Add(req.URL.String(), body)

	return &http.Response{
		StatusCode: 200,
		Body:       body.Clone(),
	}, nil
}

func NewABCIStubTransport(localClient compatlocalclient.LocalClient, cache *lru.Cache) (*httptransport.Runtime, error) {
	router := mux.NewRouter().SkipClean(true)

	viper.Set(flags.FlagTrustNode, true)

	f := httptransport.New("", "/", nil)
	f.Transport = NewRoundTripper(router, cache)

	ctx := client.
		NewCLIContext().
		WithTrustNode(true).
		WithCodec(terra.MakeCodec()).
		WithClient(localClient)

	// register all REST routes
	terra.ModuleBasics.RegisterRESTRoutes(ctx, router)

	return f, nil
}
