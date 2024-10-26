package keeper

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
)

func (ms msgServer) validateIncreaseActiveGovernanceFrameworkVersionParams(ctx sdk.Context, msg *types.MsgIncreaseActiveGovernanceFrameworkVersion) error {
	if msg.TrId == 0 {
		return errors.New("trust registry ID is mandatory")
	}

	// Direct lookup by ID
	tr, err := ms.TrustRegistry.Get(ctx, msg.TrId)
	if err != nil {
		return fmt.Errorf("trust registry with ID %d does not exist: %w", msg.TrId, err)
	}

	if tr.Controller != msg.Creator {
		return errors.New("creator is not the controller of the trust registry")
	}

	nextVersion := tr.ActiveVersion + 1

	// Find GFV for next version
	var gfv types.GovernanceFrameworkVersion
	found := false
	err = ms.GFVersion.Walk(ctx, nil, func(id uint64, v types.GovernanceFrameworkVersion) (bool, error) {
		if v.TrId == msg.TrId && v.Version == nextVersion {
			gfv = v
			found = true
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("error checking versions: %w", err)
	}
	if !found {
		return fmt.Errorf("no governance framework version found for version %d", nextVersion)
	}

	// Check for document in trust registry's language
	var hasDefaultLanguageDoc bool
	err = ms.GFDocument.Walk(ctx, nil, func(id uint64, gfd types.GovernanceFrameworkDocument) (bool, error) {
		if gfd.GfvId == gfv.Id && gfd.Language == tr.Language {
			hasDefaultLanguageDoc = true
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("error checking documents: %w", err)
	}
	if !hasDefaultLanguageDoc {
		return errors.New("no document found for the default language of this version")
	}

	return nil
}

func (ms msgServer) executeIncreaseActiveGovernanceFrameworkVersion(ctx sdk.Context, msg *types.MsgIncreaseActiveGovernanceFrameworkVersion) error {
	// Direct lookup of trust registry by ID
	tr, err := ms.TrustRegistry.Get(ctx, msg.TrId)
	if err != nil {
		return fmt.Errorf("error finding trust registry: %w", err)
	}

	nextVersion := tr.ActiveVersion + 1

	// Find the GFV to activate
	var gfv types.GovernanceFrameworkVersion
	found := false
	err = ms.GFVersion.Walk(ctx, nil, func(id uint64, v types.GovernanceFrameworkVersion) (bool, error) {
		if v.TrId == msg.TrId && v.Version == nextVersion {
			gfv = v
			found = true
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("error finding version: %w", err)
	}
	if !found {
		return fmt.Errorf("governance framework version not found")
	}

	now := ctx.BlockTime()
	tr.ActiveVersion = nextVersion
	tr.Modified = now
	gfv.ActiveSince = now

	// Update trust registry using ID as key
	if err = ms.TrustRegistry.Set(ctx, tr.Id, tr); err != nil {
		return fmt.Errorf("failed to update trust registry: %w", err)
	}

	// Update version
	if err = ms.GFVersion.Set(ctx, gfv.Id, gfv); err != nil {
		return fmt.Errorf("failed to update governance framework version: %w", err)
	}

	return nil
}
