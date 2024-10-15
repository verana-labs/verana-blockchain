package keeper

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
)

func (ms msgServer) validateIncreaseActiveGovernanceFrameworkVersionParams(ctx sdk.Context, msg *types.MsgIncreaseActiveGovernanceFrameworkVersion) error {
	if msg.Did == "" {
		return errors.New("DID is mandatory")
	}

	tr, err := ms.TrustRegistry.Get(ctx, msg.Did)
	if err != nil {
		return fmt.Errorf("trust registry with DID %s does not exist", msg.Did)
	}

	if tr.Controller != msg.Creator {
		return errors.New("creator is not the controller of the trust registry")
	}

	nextVersion := tr.ActiveVersion + 1
	var gfv types.GovernanceFrameworkVersion
	err = ms.GFVersion.Walk(ctx, nil, func(key string, v types.GovernanceFrameworkVersion) (bool, error) {
		if v.TrDid == msg.Did && v.Version == nextVersion {
			gfv = v
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("error checking versions: %w", err)
	}
	if gfv.Id == "" {
		return fmt.Errorf("no governance framework version found for version %d", nextVersion)
	}

	var gfdFound bool
	err = ms.GFDocument.Walk(ctx, nil, func(key string, gfd types.GovernanceFrameworkDocument) (bool, error) {
		if gfd.GfvId == gfv.Id && gfd.Language == tr.Language {
			gfdFound = true
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("error checking documents: %w", err)
	}
	if !gfdFound {
		return errors.New("no document found for the default language of this version")
	}

	return nil
}

func (ms msgServer) executeIncreaseActiveGovernanceFrameworkVersion(ctx sdk.Context, msg *types.MsgIncreaseActiveGovernanceFrameworkVersion) error {
	tr, err := ms.TrustRegistry.Get(ctx, msg.Did)
	if err != nil {
		return fmt.Errorf("failed to get trust registry: %w", err)
	}

	nextVersion := tr.ActiveVersion + 1
	var gfv types.GovernanceFrameworkVersion
	err = ms.GFVersion.Walk(ctx, nil, func(key string, v types.GovernanceFrameworkVersion) (bool, error) {
		if v.TrDid == msg.Did && v.Version == nextVersion {
			gfv = v
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("error checking versions: %w", err)
	}
	if gfv.Id == "" {
		return fmt.Errorf("no governance framework version found for version %d", nextVersion)
	}

	now := ctx.BlockTime()
	tr.ActiveVersion = nextVersion
	tr.Modified = now
	gfv.ActiveSince = now

	if err := ms.TrustRegistry.Set(ctx, tr.Did, tr); err != nil {
		return fmt.Errorf("failed to update trust registry: %w", err)
	}

	if err := ms.GFVersion.Set(ctx, gfv.Id, gfv); err != nil {
		return fmt.Errorf("failed to update governance framework version: %w", err)
	}

	return nil
}
