package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
)

func (msg MsgStartPermissionVP) ValidateBasic() error {
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

func (msg MsgRenewPermissionVP) ValidateBasic() error {
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
