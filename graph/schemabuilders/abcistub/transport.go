package abcistub

import (
	"encoding/json"
	"errors"
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
	mux *mux.Router
}

func NewRoundTripper(mux *mux.Router) *RoundTripper {
	return &RoundTripper{
		mux: mux,
	}
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// create a stub response
	out := &ResponseBuffer{}

	// relay to lcd...
	rt.mux.ServeHTTP(out, req)

	ret := make(map[string]interface{})
	json.Unmarshal(out.Finalize(), &ret)

	if err, exists := ret["error"]; exists {
		return nil, errors.New(err.(string))
	}

	return &http.Response{
		StatusCode: 200,
		Body:       out.Body(),
	}, nil
}

func NewABCIStubTransport(localClient compatlocalclient.LocalClient) (*httptransport.Runtime, error) {
	router := mux.NewRouter().SkipClean(true)

	viper.Set(flags.FlagTrustNode, true)

	f := httptransport.New("", "/", nil)
	f.Transport = NewRoundTripper(router)

	ctx := client.
		NewCLIContext().
		WithTrustNode(true).
		WithCodec(terra.MakeCodec()).
		WithClient(localClient)

	// register all REST routes
	terra.ModuleBasics.RegisterRESTRoutes(ctx, router)

	return f, nil
}
