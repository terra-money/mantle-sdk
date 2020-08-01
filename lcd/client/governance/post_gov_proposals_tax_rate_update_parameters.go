// Code generated by go-swagger; DO NOT EDIT.

package governance

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

// NewPostGovProposalsTaxRateUpdateParams creates a new PostGovProposalsTaxRateUpdateParams object
// with the default values initialized.
func NewPostGovProposalsTaxRateUpdateParams() *PostGovProposalsTaxRateUpdateParams {
	var ()
	return &PostGovProposalsTaxRateUpdateParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostGovProposalsTaxRateUpdateParamsWithTimeout creates a new PostGovProposalsTaxRateUpdateParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostGovProposalsTaxRateUpdateParamsWithTimeout(timeout time.Duration) *PostGovProposalsTaxRateUpdateParams {
	var ()
	return &PostGovProposalsTaxRateUpdateParams{

		timeout: timeout,
	}
}

// NewPostGovProposalsTaxRateUpdateParamsWithContext creates a new PostGovProposalsTaxRateUpdateParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostGovProposalsTaxRateUpdateParamsWithContext(ctx context.Context) *PostGovProposalsTaxRateUpdateParams {
	var ()
	return &PostGovProposalsTaxRateUpdateParams{

		Context: ctx,
	}
}

// NewPostGovProposalsTaxRateUpdateParamsWithHTTPClient creates a new PostGovProposalsTaxRateUpdateParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostGovProposalsTaxRateUpdateParamsWithHTTPClient(client *http.Client) *PostGovProposalsTaxRateUpdateParams {
	var ()
	return &PostGovProposalsTaxRateUpdateParams{
		HTTPClient: client,
	}
}

/*PostGovProposalsTaxRateUpdateParams contains all the parameters to send to the API endpoint
for the post gov proposals tax rate update operation typically these are written to a http.Request
*/
type PostGovProposalsTaxRateUpdateParams struct {

	/*PostProposalBody
	  The tax rate update body that contains new tax rate info

	*/
	PostProposalBody PostGovProposalsTaxRateUpdateBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post gov proposals tax rate update params
func (o *PostGovProposalsTaxRateUpdateParams) WithTimeout(timeout time.Duration) *PostGovProposalsTaxRateUpdateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post gov proposals tax rate update params
func (o *PostGovProposalsTaxRateUpdateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post gov proposals tax rate update params
func (o *PostGovProposalsTaxRateUpdateParams) WithContext(ctx context.Context) *PostGovProposalsTaxRateUpdateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post gov proposals tax rate update params
func (o *PostGovProposalsTaxRateUpdateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post gov proposals tax rate update params
func (o *PostGovProposalsTaxRateUpdateParams) WithHTTPClient(client *http.Client) *PostGovProposalsTaxRateUpdateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post gov proposals tax rate update params
func (o *PostGovProposalsTaxRateUpdateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPostProposalBody adds the postProposalBody to the post gov proposals tax rate update params
func (o *PostGovProposalsTaxRateUpdateParams) WithPostProposalBody(postProposalBody PostGovProposalsTaxRateUpdateBody) *PostGovProposalsTaxRateUpdateParams {
	o.SetPostProposalBody(postProposalBody)
	return o
}

// SetPostProposalBody adds the postProposalBody to the post gov proposals tax rate update params
func (o *PostGovProposalsTaxRateUpdateParams) SetPostProposalBody(postProposalBody PostGovProposalsTaxRateUpdateBody) {
	o.PostProposalBody = postProposalBody
}

// WriteToRequest writes these params to a swagger request
func (o *PostGovProposalsTaxRateUpdateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.PostProposalBody); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}