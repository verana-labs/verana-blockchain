package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"regexp"
)

func (msg *MsgUpdateTrustRegistry) ValidateBasic() error {
	if msg.Creator == "" {
		return fmt.Errorf("creator address is required")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid creator address: %s", err)
	}

	if msg.Id == 0 {
		return fmt.Errorf("trust registry id is required")
	}

	if msg.Did == "" {
		return fmt.Errorf("did is required")
	}

	if !isValidDID(msg.Did) {
		return fmt.Errorf("invalid did")
	}

	return nil
}

func isValidDID(did string) bool {
	// Basic DID validation regex
	// This is a simplified version and may need to be expanded based on specific DID method requirements
	didRegex := regexp.MustCompile(`^did:[a-zA-Z0-9]+:[a-zA-Z0-9._-]+$`)
	return didRegex.MatchString(did)
}

func (msg *MsgArchiveTrustRegistry) ValidateBasic() error {
	if msg.Creator == "" {
		return fmt.Errorf("creator address is required")
	}

	if msg.Id == 0 {
		return fmt.Errorf("trust registry id is required")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid creator address: %s", err)
	}

	return nil
}