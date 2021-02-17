// Code generated by go-swagger; DO NOT EDIT.

package transactions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/terra-project/mantle-sdk/lcd/models"
)

// PostTxsEncodeReader is a Reader for the PostTxsEncode structure.
type PostTxsEncodeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostTxsEncodeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostTxsEncodeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPostTxsEncodeBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPostTxsEncodeInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPostTxsEncodeOK creates a PostTxsEncodeOK with default headers values
func NewPostTxsEncodeOK() *PostTxsEncodeOK {
	return &PostTxsEncodeOK{}
}

/*PostTxsEncodeOK handles this case with default header values.

The tx was successfully decoded and re-encoded
*/
type PostTxsEncodeOK struct {
	Payload *PostTxsEncodeOKBody
}

func (o *PostTxsEncodeOK) Error() string {
	return fmt.Sprintf("[POST /txs/encode][%d] postTxsEncodeOK  %+v", 200, o.Payload)
}

func (o *PostTxsEncodeOK) GetPayload() *PostTxsEncodeOKBody {
	return o.Payload
}

func (o *PostTxsEncodeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(PostTxsEncodeOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostTxsEncodeBadRequest creates a PostTxsEncodeBadRequest with default headers values
func NewPostTxsEncodeBadRequest() *PostTxsEncodeBadRequest {
	return &PostTxsEncodeBadRequest{}
}

/*PostTxsEncodeBadRequest handles this case with default header values.

The tx was malformed
*/
type PostTxsEncodeBadRequest struct {
}

func (o *PostTxsEncodeBadRequest) Error() string {
	return fmt.Sprintf("[POST /txs/encode][%d] postTxsEncodeBadRequest ", 400)
}

func (o *PostTxsEncodeBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPostTxsEncodeInternalServerError creates a PostTxsEncodeInternalServerError with default headers values
func NewPostTxsEncodeInternalServerError() *PostTxsEncodeInternalServerError {
	return &PostTxsEncodeInternalServerError{}
}

/*PostTxsEncodeInternalServerError handles this case with default header values.

Server proxy_resolver error
*/
type PostTxsEncodeInternalServerError struct {
}

func (o *PostTxsEncodeInternalServerError) Error() string {
	return fmt.Sprintf("[POST /txs/encode][%d] postTxsEncodeInternalServerError ", 500)
}

func (o *PostTxsEncodeInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*PostTxsEncodeBody post txs encode body
swagger:model PostTxsEncodeBody
*/
type PostTxsEncodeBody struct {

	// tx
	Tx *models.StdTx `json:"tx,omitempty"`
}

// Validate validates this post txs encode body
func (o *PostTxsEncodeBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateTx(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostTxsEncodeBody) validateTx(formats strfmt.Registry) error {

	if swag.IsZero(o.Tx) { // not required
		return nil
	}

	if o.Tx != nil {
		if err := o.Tx.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("tx" + "." + "tx")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostTxsEncodeBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostTxsEncodeBody) UnmarshalBinary(b []byte) error {
	var res PostTxsEncodeBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*PostTxsEncodeOKBody post txs encode o k body
swagger:model PostTxsEncodeOKBody
*/
type PostTxsEncodeOKBody struct {

	// tx
	Tx string `json:"tx,omitempty"`
}

// Validate validates this post txs encode o k body
func (o *PostTxsEncodeOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostTxsEncodeOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostTxsEncodeOKBody) UnmarshalBinary(b []byte) error {
	var res PostTxsEncodeOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
