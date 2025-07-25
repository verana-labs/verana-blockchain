package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	DefaultDidDirectoryTrustDeposit = uint64(5)
	DefaultDidDirectoryGracePeriod  = uint64(30)
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(trustDeposit, gracePeriod uint64) Params {
	return Params{
		DidDirectoryTrustDeposit: trustDeposit,
		DidDirectoryGracePeriod:  gracePeriod,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultDidDirectoryTrustDeposit,
		DefaultDidDirectoryGracePeriod,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
// Validate validates the set of params
func (p Params) Validate() error {
	if p.DidDirectoryTrustDeposit == 0 {
		return fmt.Errorf("did directory trust deposit must be positive")
	}
	if p.DidDirectoryGracePeriod == 0 {
		return fmt.Errorf("did directory grace period must be positive")
	}
	return nil
}
