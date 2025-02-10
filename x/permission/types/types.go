package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
	"time"
)

func (msg *MsgStartPermissionVP) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	if msg.ValidatorPermId == 0 {
		return fmt.Errorf("validator permission ID cannot be 0")
	}

	if msg.Type == 0 || msg.Type > 6 {
		return fmt.Errorf("permission type must be between 1 and 6")
	}

	if msg.Country == "" {
		return fmt.Errorf("country must be specified")
	}

	if !isValidCountryCode(msg.Country) {
		return fmt.Errorf("invalid country code format")
	}

	if msg.Did != "" && !isValidDID(msg.Did) {
		return fmt.Errorf("invalid DID format")
	}

	return nil
}

func isValidCountryCode(code string) bool {
	// Basic check for ISO 3166-1 alpha-2 format
	match, _ := regexp.MatchString(`^[A-Z]{2}$`, code)
	return match
}

func isValidDID(did string) bool {
	// Basic DID validation regex
	match, _ := regexp.MatchString(`^did:[a-zA-Z0-9]+:[a-zA-Z0-9._-]+$`, did)
	return match
}

func (msg *MsgRenewPermissionVP) ValidateBasic() error {
	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Validate permission ID
	if msg.Id == 0 {
		return fmt.Errorf("permission ID cannot be 0")
	}

	return nil
}

// ValidateBasic for MsgSetPermissionVPToValidated
func (msg *MsgSetPermissionVPToValidated) ValidateBasic() error {
	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Validate permission ID
	if msg.Id == 0 {
		return fmt.Errorf("permission ID cannot be 0")
	}

	// Validate fees are non-negative
	if msg.ValidationFees < 0 {
		return fmt.Errorf("validation fees cannot be negative")
	}
	if msg.IssuanceFees < 0 {
		return fmt.Errorf("issuance fees cannot be negative")
	}
	if msg.VerificationFees < 0 {
		return fmt.Errorf("verification fees cannot be negative")
	}

	// Validate country code if provided
	if msg.Country != "" && !isValidCountryCode(msg.Country) {
		return fmt.Errorf("invalid country code format")
	}

	// Validate effective until if provided
	if msg.EffectiveUntil != nil && msg.EffectiveUntil.Before(time.Now()) {
		return fmt.Errorf("effective until must be in the future")
	}

	return nil
}

// ValidateBasic for MsgRequestPermissionVPTermination
func (msg *MsgRequestPermissionVPTermination) ValidateBasic() error {
	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Validate permission ID
	if msg.Id == 0 {
		return fmt.Errorf("permission ID cannot be 0")
	}

	return nil
}

// ValidateBasic for MsgConfirmPermissionVPTermination
func (msg *MsgConfirmPermissionVPTermination) ValidateBasic() error {
	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Validate permission ID
	if msg.Id == 0 {
		return fmt.Errorf("permission ID cannot be 0")
	}

	return nil
}

// ValidateBasic for MsgConfirmPermissionVPTermination
func (msg *MsgCancelPermissionVPLastRequest) ValidateBasic() error {
	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Validate permission ID
	if msg.Id == 0 {
		return fmt.Errorf("permission ID cannot be 0")
	}

	return nil
}

func (msg *MsgCreateRootPermission) ValidateBasic() error {
	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Validate schema ID
	if msg.SchemaId == 0 {
		return fmt.Errorf("schema ID cannot be 0")
	}

	// Validate DID
	if msg.Did == "" {
		return fmt.Errorf("DID is required")
	}
	if !isValidDID(msg.Did) {
		return fmt.Errorf("invalid DID format")
	}

	// Validate fees are non-negative
	if msg.ValidationFees < 0 {
		return fmt.Errorf("validation fees cannot be negative")
	}
	if msg.IssuanceFees < 0 {
		return fmt.Errorf("issuance fees cannot be negative")
	}
	if msg.VerificationFees < 0 {
		return fmt.Errorf("verification fees cannot be negative")
	}

	// Validate country code if present
	if msg.Country != "" && !isValidCountryCode(msg.Country) {
		return fmt.Errorf("invalid country code format")
	}

	// Validate effective dates if present
	now := time.Now()
	if msg.EffectiveFrom != nil {
		if !msg.EffectiveFrom.After(now) {
			return fmt.Errorf("effective_from must be in the future")
		}

		if msg.EffectiveUntil != nil {
			if !msg.EffectiveUntil.After(*msg.EffectiveFrom) {
				return fmt.Errorf("effective_until must be after effective_from")
			}
		}
	}

	return nil
}

func (msg *MsgExtendPermission) ValidateBasic() error {
	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Validate permission ID
	if msg.Id == 0 {
		return fmt.Errorf("permission ID cannot be 0")
	}

	// Validate effective_until is in the future
	if msg.EffectiveUntil != nil && !msg.EffectiveUntil.After(time.Now()) {
		return fmt.Errorf("effective_until must be in the future")
	}

	return nil
}

func (msg *MsgRevokePermission) ValidateBasic() error {
	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Validate permission ID
	if msg.Id == 0 {
		return fmt.Errorf("permission ID cannot be 0")
	}

	return nil
}
