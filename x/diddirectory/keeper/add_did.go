package keeper

import (
	"errors"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/diddirectory/types"
)

func (ms msgServer) validateAddDIDParams(ctx sdk.Context, msg *types.MsgAddDID) error {
	if msg.Did == "" {
		return errors.New("DID is required")
	}

	// Validate DID format
	if !isValidDID(msg.Did) {
		return errors.New("invalid DID syntax")
	}

	// Check if DID already exists
	_, err := ms.DIDDirectory.Get(ctx, msg.Did)
	if err == nil {
		return errors.New("DID already exists")
	}

	// Validate years (1-31)
	years := msg.Years
	if years == 0 {
		years = 1
	}
	if years > 31 {
		return errors.New("years must be between 1 and 31")
	}

	return nil
}

func isValidDID(did string) bool {
	// Basic DID validation regex
	// This is a simplified version and may need to be expanded based on specific DID method requirements
	didRegex := regexp.MustCompile(`^did:[a-zA-Z0-9]+:[a-zA-Z0-9._-]+$`)
	return didRegex.MatchString(did)
}

func (ms msgServer) checkSufficientFees(_ sdk.Context, _ string, _ uint32) error {
	return nil
}

func (ms msgServer) executeAddDID(ctx sdk.Context, msg *types.MsgAddDID) error {
	params := ms.GetParams(ctx)

	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	years := msg.Years
	if years == 0 {
		years = 1
	}

	now := ctx.BlockTime()
	expiration := now.AddDate(int(years), 0, 0)

	// Create DID entry
	didEntry := types.DIDDirectory{
		Did:        msg.Did,
		Controller: msg.Creator,
		Created:    now,
		Modified:   now,
		Exp:        expiration,
		Deposit:    int64(params.DidDirectoryTrustDeposit * uint64(years)),
	}

	// Lock trust deposit

	// Lock removal gas in escrow

	// Store the DID entry
	if err = ms.DIDDirectory.Set(ctx, msg.Did, didEntry); err != nil {
		return fmt.Errorf("failed to store DID: %w", err)
	}

	return nil
}
