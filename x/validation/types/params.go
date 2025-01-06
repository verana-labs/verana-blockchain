package types

import (
	"errors"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// DefaultValidationTermRequestedTimeoutDays is the default timeout period in days
const DefaultValidationTermRequestedTimeoutDays uint64 = 7

// DefaultParams returns default validation parameters
func DefaultParams() Params {
	return Params{
		ValidationTermRequestedTimeoutDays: DefaultValidationTermRequestedTimeoutDays,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if p.ValidationTermRequestedTimeoutDays == 0 {
		return errors.New("validation term requested timeout days must be greater than 0")
	}
	return nil
}
