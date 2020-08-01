// Code generated by go-swagger; DO NOT EDIT.

package treasury

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new treasury API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for treasury API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientService is the interface for Client methods
type ClientService interface {
	GetTreasuryParameters(params *GetTreasuryParametersParams) (*GetTreasuryParametersOK, error)

	GetTreasuryRewardWeight(params *GetTreasuryRewardWeightParams) (*GetTreasuryRewardWeightOK, error)

	GetTreasurySeigniorageProceeds(params *GetTreasurySeigniorageProceedsParams) (*GetTreasurySeigniorageProceedsOK, error)

	GetTreasuryTaxCapDenom(params *GetTreasuryTaxCapDenomParams) (*GetTreasuryTaxCapDenomOK, error)

	GetTreasuryTaxProceeds(params *GetTreasuryTaxProceedsParams) (*GetTreasuryTaxProceedsOK, error)

	GetTreasuryTaxRate(params *GetTreasuryTaxRateParams) (*GetTreasuryTaxRateOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  GetTreasuryParameters gets treasury module params
*/
func (a *Client) GetTreasuryParameters(params *GetTreasuryParametersParams) (*GetTreasuryParametersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetTreasuryParametersParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetTreasuryParameters",
		Method:             "GET",
		PathPattern:        "/treasury/parameters",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetTreasuryParametersReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetTreasuryParametersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetTreasuryParameters: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetTreasuryRewardWeight gets current reward weight
*/
func (a *Client) GetTreasuryRewardWeight(params *GetTreasuryRewardWeightParams) (*GetTreasuryRewardWeightOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetTreasuryRewardWeightParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetTreasuryRewardWeight",
		Method:             "GET",
		PathPattern:        "/treasury/reward_weight",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetTreasuryRewardWeightReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetTreasuryRewardWeightOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetTreasuryRewardWeight: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetTreasurySeigniorageProceeds retrieves the size of the seigniorage pool
*/
func (a *Client) GetTreasurySeigniorageProceeds(params *GetTreasurySeigniorageProceedsParams) (*GetTreasurySeigniorageProceedsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetTreasurySeigniorageProceedsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetTreasurySeigniorageProceeds",
		Method:             "GET",
		PathPattern:        "/treasury/seigniorage_proceeds",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetTreasurySeigniorageProceedsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetTreasurySeigniorageProceedsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetTreasurySeigniorageProceeds: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetTreasuryTaxCapDenom gets tax cap of the denom
*/
func (a *Client) GetTreasuryTaxCapDenom(params *GetTreasuryTaxCapDenomParams) (*GetTreasuryTaxCapDenomOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetTreasuryTaxCapDenomParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetTreasuryTaxCapDenom",
		Method:             "GET",
		PathPattern:        "/treasury/tax_cap/{denom}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetTreasuryTaxCapDenomReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetTreasuryTaxCapDenomOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetTreasuryTaxCapDenom: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetTreasuryTaxProceeds gets current tax proceeds
*/
func (a *Client) GetTreasuryTaxProceeds(params *GetTreasuryTaxProceedsParams) (*GetTreasuryTaxProceedsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetTreasuryTaxProceedsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetTreasuryTaxProceeds",
		Method:             "GET",
		PathPattern:        "/treasury/tax_proceeds",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetTreasuryTaxProceedsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetTreasuryTaxProceedsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetTreasuryTaxProceeds: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetTreasuryTaxRate gets current tax rate
*/
func (a *Client) GetTreasuryTaxRate(params *GetTreasuryTaxRateParams) (*GetTreasuryTaxRateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetTreasuryTaxRateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetTreasuryTaxRate",
		Method:             "GET",
		PathPattern:        "/treasury/tax_rate",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetTreasuryTaxRateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetTreasuryTaxRateOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetTreasuryTaxRate: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}