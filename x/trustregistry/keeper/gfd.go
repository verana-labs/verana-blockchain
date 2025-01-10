package keeper

import (
	"errors"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
)

func (ms msgServer) validateAddGovernanceFrameworkDocumentParams(ctx sdk.Context, msg *types.MsgAddGovernanceFrameworkDocument) error {
	// Check mandatory parameters
	if msg.Id == 0 || msg.DocLanguage == "" || msg.DocUrl == "" || msg.DocHash == "" {
		return errors.New("missing mandatory parameter")
	}

	// Direct lookup of trust registry by ID
	tr, err := ms.TrustRegistry.Get(ctx, msg.Id)
	if err != nil {
		return fmt.Errorf("trust registry with ID %d does not exist: %w", msg.Id, err)
	}

	// Check controller
	if tr.Controller != msg.Creator {
		return errors.New("creator is not the controller of the trust registry")
	}

	// Check version validity
	var maxVersion int32
	var hasVersion bool
	err = ms.GFVersion.Walk(ctx, nil, func(id uint64, gfv types.GovernanceFrameworkVersion) (bool, error) {
		if gfv.TrId == msg.Id {
			if gfv.Version == msg.Version {
				hasVersion = true
			}
			if gfv.Version > maxVersion {
				maxVersion = gfv.Version
			}
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("error checking versions: %w", err)
	}

	// Validate version according to spec
	if !hasVersion && msg.Version != maxVersion+1 {
		return fmt.Errorf("invalid version: must be %d or %d", maxVersion, maxVersion+1)
	}

	if msg.Version <= tr.ActiveVersion {
		return fmt.Errorf("version must be greater than the active version %d", tr.ActiveVersion)
	}

	// Validate language tag
	if !isValidLanguageTag(msg.DocLanguage) {
		return errors.New("invalid language tag (must conform to rfc1766)")
	}

	// Validate URL and hash
	if !isValidURL(msg.DocUrl) {
		return errors.New("invalid document URL")
	}

	if !isValidHash(msg.DocHash) {
		return errors.New("invalid document hash")
	}

	return nil
}

func (ms msgServer) executeAddGovernanceFrameworkDocument(ctx sdk.Context, msg *types.MsgAddGovernanceFrameworkDocument) error {
	now := ctx.BlockTime()

	var gfv types.GovernanceFrameworkVersion
	var gfvExists bool

	// Check if version exists
	err := ms.GFVersion.Walk(ctx, nil, func(id uint64, v types.GovernanceFrameworkVersion) (bool, error) {
		if v.TrId == msg.Id && v.Version == msg.Version {
			gfv = v
			gfvExists = true
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("error checking existing version: %w", err)
	}

	// Create new version if needed
	if !gfvExists {
		gfvID, err := ms.Keeper.GetNextID(ctx, "gfv")
		if err != nil {
			return fmt.Errorf("failed to generate GFV ID: %w", err)
		}

		gfv = types.GovernanceFrameworkVersion{
			Id:          gfvID,
			TrId:        msg.Id,
			Created:     now,
			Version:     msg.Version,
			ActiveSince: time.Time{}, // Zero time as per spec - not active yet
		}
		if err := ms.GFVersion.Set(ctx, gfv.Id, gfv); err != nil {
			return fmt.Errorf("failed to persist GovernanceFrameworkVersion: %w", err)
		}
	}

	// Create new document
	gfdID, err := ms.Keeper.GetNextID(ctx, "gfd")
	if err != nil {
		return fmt.Errorf("failed to generate GFD ID: %w", err)
	}

	gfd := types.GovernanceFrameworkDocument{
		Id:       gfdID,
		GfvId:    gfv.Id,
		Created:  now,
		Language: msg.DocLanguage,
		Url:      msg.DocUrl,
		Hash:     msg.DocHash,
	}

	if err := ms.GFDocument.Set(ctx, gfd.Id, gfd); err != nil {
		return fmt.Errorf("failed to persist GovernanceFrameworkDocument: %w", err)
	}

	return nil
}
