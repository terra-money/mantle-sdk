// Code generated by go-swagger; DO NOT EDIT.

package tendermint_rpc

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetValidatorsetsLatestParams creates a new GetValidatorsetsLatestParams object
// with the default values initialized.
func NewGetValidatorsetsLatestParams() *GetValidatorsetsLatestParams {

	return &GetValidatorsetsLatestParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetValidatorsetsLatestParamsWithTimeout creates a new GetValidatorsetsLatestParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetValidatorsetsLatestParamsWithTimeout(timeout time.Duration) *GetValidatorsetsLatestParams {

	return &GetValidatorsetsLatestParams{

		timeout: timeout,
	}
}

// NewGetValidatorsetsLatestParamsWithContext creates a new GetValidatorsetsLatestParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetValidatorsetsLatestParamsWithContext(ctx context.Context) *GetValidatorsetsLatestParams {

	return &GetValidatorsetsLatestParams{

		Context: ctx,
	}
}

// NewGetValidatorsetsLatestParamsWithHTTPClient creates a new GetValidatorsetsLatestParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetValidatorsetsLatestParamsWithHTTPClient(client *http.Client) *GetValidatorsetsLatestParams {

	return &GetValidatorsetsLatestParams{
		HTTPClient: client,
	}
}

/*GetValidatorsetsLatestParams contains all the parameters to send to the API endpoint
for the get validatorsets latest operation typically these are written to a http.Request
*/
type GetValidatorsetsLatestParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get validatorsets latest params
func (o *GetValidatorsetsLatestParams) WithTimeout(timeout time.Duration) *GetValidatorsetsLatestParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get validatorsets latest params
func (o *GetValidatorsetsLatestParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get validatorsets latest params
func (o *GetValidatorsetsLatestParams) WithContext(ctx context.Context) *GetValidatorsetsLatestParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get validatorsets latest params
func (o *GetValidatorsetsLatestParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get validatorsets latest params
func (o *GetValidatorsetsLatestParams) WithHTTPClient(client *http.Client) *GetValidatorsetsLatestParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get validatorsets latest params
func (o *GetValidatorsetsLatestParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetValidatorsetsLatestParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}