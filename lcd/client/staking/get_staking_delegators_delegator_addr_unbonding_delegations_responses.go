// Code generated by go-swagger; DO NOT EDIT.

package staking

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/terra-project/mantle/lcd/models"
)

// GetStakingDelegatorsDelegatorAddrUnbondingDelegationsReader is a Reader for the GetStakingDelegatorsDelegatorAddrUnbondingDelegations structure.
type GetStakingDelegatorsDelegatorAddrUnbondingDelegationsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK creates a GetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK with default headers values
func NewGetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK() *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK {
	return &GetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK{}
}

/*GetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK handles this case with default header values.

OK
*/
type GetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK struct {
	Payload []*models.UnbondingDelegation
}

func (o *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK) Error() string {
	return fmt.Sprintf("[GET /staking/delegators/{delegatorAddr}/unbonding_delegations][%d] getStakingDelegatorsDelegatorAddrUnbondingDelegationsOK  %+v", 200, o.Payload)
}

func (o *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK) GetPayload() []*models.UnbondingDelegation {
	return o.Payload
}

func (o *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest creates a GetStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest with default headers values
func NewGetStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest() *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest {
	return &GetStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest{}
}

/*GetStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest handles this case with default header values.

Invalid delegator address
*/
type GetStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest struct {
}

func (o *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest) Error() string {
	return fmt.Sprintf("[GET /staking/delegators/{delegatorAddr}/unbonding_delegations][%d] getStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest ", 400)
}

func (o *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError creates a GetStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError with default headers values
func NewGetStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError() *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError {
	return &GetStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError{}
}

/*GetStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError handles this case with default header values.

Internal Server Error
*/
type GetStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError struct {
}

func (o *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /staking/delegators/{delegatorAddr}/unbonding_delegations][%d] getStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError ", 500)
}

func (o *GetStakingDelegatorsDelegatorAddrUnbondingDelegationsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}