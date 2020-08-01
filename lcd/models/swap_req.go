// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SwapReq swap req
//
// swagger:model SwapReq
type SwapReq struct {

	// ask denom
	// Required: true
	AskDenom *string `json:"ask_denom"`

	// base req
	// Required: true
	BaseReq *BaseReq `json:"base_req"`

	// offer coin
	// Required: true
	OfferCoin *Coin `json:"offer_coin"`

	// the `receiver` field can be skipped when the receiver is trader
	Receiver string `json:"receiver,omitempty"`
}

// Validate validates this swap req
func (m *SwapReq) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAskDenom(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBaseReq(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOfferCoin(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SwapReq) validateAskDenom(formats strfmt.Registry) error {

	if err := validate.Required("ask_denom", "body", m.AskDenom); err != nil {
		return err
	}

	return nil
}

func (m *SwapReq) validateBaseReq(formats strfmt.Registry) error {

	if err := validate.Required("base_req", "body", m.BaseReq); err != nil {
		return err
	}

	if m.BaseReq != nil {
		if err := m.BaseReq.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("base_req")
			}
			return err
		}
	}

	return nil
}

func (m *SwapReq) validateOfferCoin(formats strfmt.Registry) error {

	if err := validate.Required("offer_coin", "body", m.OfferCoin); err != nil {
		return err
	}

	if m.OfferCoin != nil {
		if err := m.OfferCoin.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("offer_coin")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SwapReq) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SwapReq) UnmarshalBinary(b []byte) error {
	var res SwapReq
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}