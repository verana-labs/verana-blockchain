package keeper

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
	"time"
)

func (ms msgServer) validateAddGovernanceFrameworkDocumentParams(ctx sdk.Context, msg *types.MsgAddGovernanceFrameworkDocument) error {
	//TODO: For language create a list of all the acceptable languages and then allow the given languages only
	if msg.Did == "" || msg.DocLanguage == "" || msg.DocUrl == "" || msg.DocHash == "" {
		return errors.New("missing mandatory parameter")
	}

	tr, err := ms.TrustRegistry.Get(ctx, msg.Did)
	if err != nil {
		return fmt.Errorf("trust registry with DID %s does not exist", msg.Did)
	}

	if tr.Controller != msg.Creator {
		return errors.New("creator is not the controller of the trust registry")
	}

	// Check if the version is valid
	var maxVersion int32
	err = ms.GFVersion.Walk(ctx, nil, func(key string, gfv types.GovernanceFrameworkVersion) (bool, error) {
		if gfv.TrDid == msg.Did && gfv.Version > maxVersion {
			maxVersion = gfv.Version
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("error checking versions: %w", err)
	}

	if msg.Version != maxVersion+1 && msg.Version != maxVersion {
		return fmt.Errorf("invalid version: must be %d or %d", maxVersion, maxVersion+1)
	}

	if msg.Version <= tr.ActiveVersion {
		return fmt.Errorf("version must be greater than the active version %d", tr.ActiveVersion)
	}

	if !isValidLanguageTag(msg.DocLanguage) {
		return errors.New("invalid language tag (must conform to rfc1766)")
	}

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
	var err error

	// Check if the version already exists
	err = ms.GFVersion.Walk(ctx, nil, func(key string, v types.GovernanceFrameworkVersion) (bool, error) {
		if v.TrDid == msg.Did && v.Version == msg.Version {
			gfv = v
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("error checking existing version: %w", err)
	}

	// If the version doesn't exist, create a new one
	if gfv.Id == "" {
		gfv = types.GovernanceFrameworkVersion{
			Id:          uuid.New().String(),
			TrDid:       msg.Did,
			Created:     now,
			Version:     msg.Version,
			ActiveSince: time.Time{}, // Set to zero time as it's not active yet
		}
		if err := ms.GFVersion.Set(ctx, gfv.Id, gfv); err != nil {
			return fmt.Errorf("failed to persist GovernanceFrameworkVersion: %w", err)
		}
	}

	// Create and persist the new GovernanceFrameworkDocument
	gfd := types.GovernanceFrameworkDocument{
		Id:       uuid.New().String(),
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
