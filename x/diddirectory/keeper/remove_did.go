package keeper

import (
	"errors"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/diddirectory/types"
)

func (ms msgServer) validateRemoveDIDParams(ctx sdk.Context, msg *types.MsgRemoveDID) error {
	if msg.Did == "" {
		return errors.New("DID is required")
	}

	if !isValidDID(msg.Did) {
		return errors.New("invalid DID syntax")
	}

	// Get DID entry
	didEntry, err := ms.DIDDirectory.Get(ctx, msg.Did)
	if err != nil {
		return fmt.Errorf("DID not found: %w", err)
	}

	// Get grace period
	params := ms.GetParams(ctx)
	gracePeriod := time.Duration(params.DidDirectoryGracePeriod) * 24 * time.Hour

	// Check authorization
	now := ctx.BlockTime()
	if now.Before(didEntry.Exp.Add(gracePeriod)) {
		// Before grace period: only controller can remove
		if msg.Creator != didEntry.Controller {
			return errors.New("only the controller can remove this DID before grace period")
		}
	}

	return nil
}

func (ms msgServer) executeRemoveDID(ctx sdk.Context, msg *types.MsgRemoveDID) error {
	_, err := ms.DIDDirectory.Get(ctx, msg.Did)
	if err != nil {
		return fmt.Errorf("error retrieving DID: %w", err)
	}

	// Release trust deposit to controller

	// Remove DID entry
	if err = ms.DIDDirectory.Remove(ctx, msg.Did); err != nil {
		return fmt.Errorf("failed to remove DID: %w", err)
	}

	return nil
}
