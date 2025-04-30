package keeper

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
)

func (ms msgServer) validateSessionAccess(ctx sdk.Context, msg *types.MsgCreateOrUpdatePermissionSession) error {
	existingSession, err := ms.PermissionSession.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil // New session case
		}
		return sdkerrors.ErrInvalidRequest.Wrapf("failed to get session: %v", err)
	}

	// Only session controller can update
	if existingSession.Controller != msg.Creator {
		return sdkerrors.ErrUnauthorized.Wrap("only session controller can update")
	}

	// Check for duplicate authorization
	for _, authz := range existingSession.Authz {
		if authz.ExecutorPermId == msg.IssuerPermId &&
			authz.BeneficiaryPermId == msg.VerifierPermId &&
			authz.WalletAgentPermId == msg.WalletAgentPermId {
			return sdkerrors.ErrInvalidRequest.Wrap("authorization already exists")
		}
	}

	return nil
}

func (ms msgServer) processFees(
	ctx sdk.Context,
	creator string,
	permSet []types.Permission,
	isVerifier bool,
	trustUnitPrice uint64,
	trustDepositRate math.LegacyDec,
) error {
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return fmt.Errorf("invalid creator address: %w", err)
	}

	// Process each permission's fees
	for _, perm := range permSet {
		var fees uint64
		if isVerifier {
			fees = perm.VerificationFees
		} else {
			fees = perm.IssuanceFees
		}

		if fees > 0 {
			// Calculate fees in denom
			feesInDenom := fees * trustUnitPrice

			// Calculate trust deposit amount
			trustDepositAmount := uint64(math.LegacyNewDec(int64(feesInDenom)).Mul(trustDepositRate).TruncateInt64())

			// Calculate direct fees (the portion that goes directly to the grantee)
			directFeesAmount := feesInDenom - trustDepositAmount

			// 1. Transfer direct fees from creator to permission grantee
			if directFeesAmount > 0 {
				granteeAddr, err := sdk.AccAddressFromBech32(perm.Grantee)
				if err != nil {
					return fmt.Errorf("invalid grantee address: %w", err)
				}

				err = ms.bankKeeper.SendCoins(
					ctx,
					creatorAddr,
					granteeAddr,
					sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, int64(directFeesAmount))),
				)
				if err != nil {
					return fmt.Errorf("failed to transfer direct fees: %w", err)
				}
			}

			// 2. Increase trust deposit for the grantee
			if trustDepositAmount > 0 {
				// First transfer funds from creator to module account
				err = ms.bankKeeper.SendCoinsFromAccountToModule(
					ctx,
					creatorAddr,
					types.ModuleName,
					sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, int64(trustDepositAmount))),
				)
				if err != nil {
					return fmt.Errorf("failed to transfer trust deposit to module: %w", err)
				}

				// Then adjust grantee's trust deposit
				err = ms.trustDeposit.AdjustTrustDeposit(
					ctx,
					perm.Grantee,
					int64(trustDepositAmount),
				)
				if err != nil {
					return fmt.Errorf("failed to adjust grantee trust deposit: %w", err)
				}
			}
		}
	}

	return nil
}

func (ms msgServer) createOrUpdateSession(ctx sdk.Context, msg *types.MsgCreateOrUpdatePermissionSession, now time.Time) error {
	session := &types.PermissionSession{
		Id:          msg.Id,
		Controller:  msg.Creator,
		AgentPermId: msg.AgentPermId,
		Modified:    &now,
	}

	existingSession, err := ms.PermissionSession.Get(ctx, msg.Id)
	if err == nil {
		// Update existing session
		session = &existingSession
		session.Modified = &now
	} else if errors.Is(err, collections.ErrNotFound) {
		// New session
		session.Created = &now
	} else {
		return err
	}

	// Add new authorization
	session.Authz = append(session.Authz, &types.SessionAuthz{
		ExecutorPermId:    msg.IssuerPermId,
		BeneficiaryPermId: msg.VerifierPermId,
		WalletAgentPermId: msg.WalletAgentPermId,
	})

	return ms.PermissionSession.Set(ctx, msg.Id, *session)
}

// findBeneficiaries gets the set of permissions that should receive fees
func (ms msgServer) findBeneficiaries(ctx sdk.Context, issuerPermId, verifierPermId uint64) ([]types.Permission, error) {
	var foundPerms []types.Permission

	// Helper function to check if a permission is already in the slice
	containsPerm := func(id uint64) bool {
		for _, p := range foundPerms {
			if p.Id == id {
				return true
			}
		}
		return false
	}

	// Process issuer permission hierarchy if provided
	if issuerPermId != 0 {
		issuerPerm, err := ms.Permission.Get(ctx, issuerPermId)
		if err != nil {
			return nil, fmt.Errorf("issuer permission not found: %w", err)
		}

		// Follow the validator chain up
		if issuerPerm.ValidatorPermId != 0 {
			currentPermID := issuerPerm.ValidatorPermId
			for currentPermID != 0 {
				currentPerm, err := ms.Permission.Get(ctx, currentPermID)
				if err != nil {
					return nil, fmt.Errorf("failed to get permission: %w", err)
				}

				// Add to set if valid and not already included
				if currentPerm.Revoked == nil && currentPerm.Terminated == nil && !containsPerm(currentPermID) {
					foundPerms = append(foundPerms, currentPerm)
				}

				// Move up
				currentPermID = currentPerm.ValidatorPermId
			}
		}
	}

	// Process verifier permission hierarchy if provided
	if verifierPermId != 0 {
		// First add issuer permission to the set if provided
		if issuerPermId != 0 {
			issuerPerm, err := ms.Permission.Get(ctx, issuerPermId)
			if err == nil && issuerPerm.Revoked == nil && issuerPerm.Terminated == nil && !containsPerm(issuerPermId) {
				foundPerms = append(foundPerms, issuerPerm)
			}
		}

		// Then process verifier's validator chain
		verifierPerm, err := ms.Permission.Get(ctx, verifierPermId)
		if err != nil {
			return nil, fmt.Errorf("verifier permission not found: %w", err)
		}

		if verifierPerm.ValidatorPermId != 0 {
			currentPermID := verifierPerm.ValidatorPermId
			for currentPermID != 0 {
				currentPerm, err := ms.Permission.Get(ctx, currentPermID)
				if err != nil {
					return nil, fmt.Errorf("failed to get permission: %w", err)
				}

				// Add to set if valid and not already included
				if currentPerm.Revoked == nil && currentPerm.Terminated == nil && !containsPerm(currentPermID) {
					foundPerms = append(foundPerms, currentPerm)
				}

				// Move up
				currentPermID = currentPerm.ValidatorPermId
			}
		}
	}

	return foundPerms, nil
}
