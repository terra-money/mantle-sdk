// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Block block
//
// swagger:model Block
type Block struct {

	// evidence
	Evidence []string `json:"evidence"`

	// header
	Header *BlockHeader `json:"header,omitempty"`

	// last commit
	LastCommit *BlockLastCommit `json:"last_commit,omitempty"`

	// txs
	Txs []string `json:"txs"`
}

// Validate validates this block
func (m *Block) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHeader(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLastCommit(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Block) validateHeader(formats strfmt.Registry) error {

	if swag.IsZero(m.Header) { // not required
		return nil
	}

	if m.Header != nil {
		if err := m.Header.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("header")
			}
			return err
		}
	}

	return nil
}

func (m *Block) validateLastCommit(formats strfmt.Registry) error {

	if swag.IsZero(m.LastCommit) { // not required
		return nil
	}

	if m.LastCommit != nil {
		if err := m.LastCommit.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("last_commit")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Block) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Block) UnmarshalBinary(b []byte) error {
	var res Block
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// BlockLastCommit block last commit
//
// swagger:model BlockLastCommit
type BlockLastCommit struct {

	// block id
	BlockID *BlockID `json:"block_id,omitempty"`

	// precommits
	Precommits []*BlockLastCommitPrecommitsItems0 `json:"precommits"`
}

// Validate validates this block last commit
func (m *BlockLastCommit) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBlockID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePrecommits(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BlockLastCommit) validateBlockID(formats strfmt.Registry) error {

	if swag.IsZero(m.BlockID) { // not required
		return nil
	}

	if m.BlockID != nil {
		if err := m.BlockID.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("last_commit" + "." + "block_id")
			}
			return err
		}
	}

	return nil
}

func (m *BlockLastCommit) validatePrecommits(formats strfmt.Registry) error {

	if swag.IsZero(m.Precommits) { // not required
		return nil
	}

	for i := 0; i < len(m.Precommits); i++ {
		if swag.IsZero(m.Precommits[i]) { // not required
			continue
		}

		if m.Precommits[i] != nil {
			if err := m.Precommits[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("last_commit" + "." + "precommits" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *BlockLastCommit) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BlockLastCommit) UnmarshalBinary(b []byte) error {
	var res BlockLastCommit
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// BlockLastCommitPrecommitsItems0 block last commit precommits items0
//
// swagger:model BlockLastCommitPrecommitsItems0
type BlockLastCommitPrecommitsItems0 struct {

	// block id
	BlockID *BlockID `json:"block_id,omitempty"`

	// height
	Height string `json:"height,omitempty"`

	// round
	Round string `json:"round,omitempty"`

	// signature
	Signature string `json:"signature,omitempty"`

	// timestamp
	Timestamp string `json:"timestamp,omitempty"`

	// type
	Type float64 `json:"type,omitempty"`

	// validator address
	ValidatorAddress string `json:"validator_address,omitempty"`

	// validator index
	ValidatorIndex string `json:"validator_index,omitempty"`
}

// Validate validates this block last commit precommits items0
func (m *BlockLastCommitPrecommitsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBlockID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BlockLastCommitPrecommitsItems0) validateBlockID(formats strfmt.Registry) error {

	if swag.IsZero(m.BlockID) { // not required
		return nil
	}

	if m.BlockID != nil {
		if err := m.BlockID.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("block_id")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *BlockLastCommitPrecommitsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BlockLastCommitPrecommitsItems0) UnmarshalBinary(b []byte) error {
	var res BlockLastCommitPrecommitsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}