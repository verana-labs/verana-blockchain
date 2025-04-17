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

func (ms msgServer) validateAgentPermission(ctx sdk.Context, msg *types.MsgCreateOrUpdatePermissionSession) error {
	agentPerm, err := ms.Permission.Get(ctx, msg.AgentPermId)
	if err != nil {
		return sdkerrors.ErrNotFound.Wrap("agent permission not found")
	}

	if agentPerm.Type != types.PermissionType_PERMISSION_TYPE_HOLDER {
		return sdkerrors.ErrInvalidRequest.Wrap("agent permission must be HOLDER type")
	}

	if agentPerm.Revoked != nil || agentPerm.Terminated != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("agent permission is revoked or terminated")
	}

	return nil
}

func (ms msgServer) validateWalletAgentPermission(ctx sdk.Context, msg *types.MsgCreateOrUpdatePermissionSession) error {
	if msg.WalletAgentPermId == 0 {
		return nil // Optional field
	}

	walletAgentPerm, err := ms.Permission.Get(ctx, msg.WalletAgentPermId)
	if err != nil {
		return sdkerrors.ErrNotFound.Wrap("wallet agent permission not found")
	}

	if walletAgentPerm.Type != types.PermissionType_PERMISSION_TYPE_HOLDER {
		return sdkerrors.ErrInvalidRequest.Wrap("wallet agent permission must be HOLDER type")
	}

	if walletAgentPerm.Revoked != nil || walletAgentPerm.Terminated != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("wallet agent permission is revoked or terminated")
	}

	return nil
}

func (ms msgServer) buildPermissionSet(ctx sdk.Context, executorPerm *types.Permission, beneficiaryPermId uint64) (PermissionSet, error) {
	var permSet PermissionSet

	// Process executor perm ancestors
	currentPerm := executorPerm
	for currentPerm.ValidatorPermId != 0 {
		validatorPerm, err := ms.Permission.Get(ctx, currentPerm.ValidatorPermId)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to get validator permission: %v", err)
		}

		// Add if not revoked and not terminated and not already in set
		if validatorPerm.Revoked == nil && validatorPerm.Terminated == nil && !permSet.contains(validatorPerm.Id) {
			permSet = append(permSet, validatorPerm)
		}

		currentPerm = &validatorPerm
	}

	// For VERIFIER type, process beneficiary permission chain
	if executorPerm.Type == types.PermissionType_PERMISSION_TYPE_VERIFIER {
		if beneficiaryPermId == 0 {
			return nil, sdkerrors.ErrInvalidRequest.Wrap("beneficiary permission required for VERIFIER")
		}

		beneficiaryPerm, err := ms.Permission.Get(ctx, beneficiaryPermId)
		if err != nil {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to get beneficiary permission: %v", err)
		}

		// Validate beneficiary is ISSUER
		if beneficiaryPerm.Type != types.PermissionType_PERMISSION_TYPE_ISSUER {
			return nil, sdkerrors.ErrInvalidRequest.Wrap("beneficiary must be ISSUER")
		}

		// Add beneficiary if valid
		if beneficiaryPerm.Revoked == nil && beneficiaryPerm.Terminated == nil && !permSet.contains(beneficiaryPerm.Id) {
			permSet = append(permSet, beneficiaryPerm)
		}

		// Process beneficiary ancestors
		currentPerm = &beneficiaryPerm
		for currentPerm.ValidatorPermId != 0 {
			validatorPerm, err := ms.Permission.Get(ctx, currentPerm.ValidatorPermId)
			if err != nil {
				return nil, sdkerrors.ErrInvalidRequest.Wrapf("failed to get validator permission: %v", err)
			}

			if validatorPerm.Revoked == nil && validatorPerm.Terminated == nil && !permSet.contains(validatorPerm.Id) {
				permSet = append(permSet, validatorPerm)
			}

			currentPerm = &validatorPerm
		}
	}

	return permSet, nil
}

func (ms msgServer) calculateAndValidateFees(ctx sdk.Context, creator string, permSet PermissionSet, executorType types.PermissionType) (sdk.Coin, error) {
	beneficiaryFees := math.NewInt(0)

	// Calculate total beneficiary fees
	for _, perm := range permSet {
		if executorType == types.PermissionType_PERMISSION_TYPE_VERIFIER {
			beneficiaryFees = beneficiaryFees.Add(math.NewInt(int64(perm.VerificationFees)))
		} else {
			beneficiaryFees = beneficiaryFees.Add(math.NewInt(int64(perm.IssuanceFees)))
		}
	}

	// Get global variables
	trustUnitPrice := ms.trustRegistryKeeper.GetTrustUnitPrice(ctx)
	trustDepositRate := ms.trustDeposit.GetTrustDepositRate(ctx)
	userAgentRewardRate := ms.trustDeposit.GetUserAgentRewardRate(ctx)
	walletUserAgentRewardRate := ms.trustDeposit.GetWalletUserAgentRewardRate(ctx)

	// Calculate total fees including trust deposit and rewards
	totalFees := beneficiaryFees.Mul(math.NewInt(int64(trustUnitPrice)))
	trustFees := ms.Keeper.trustFeesAmount(totalFees.Int64(), trustDepositRate)

	rewardRateSum := userAgentRewardRate.Add(walletUserAgentRewardRate)
	rewards := ms.Keeper.rewardsAmount(totalFees.Int64(), rewardRateSum)

	requiredAmount := sdk.NewCoin(types.BondDenom, totalFees.Add(trustFees).Add(rewards))

	// Validate sufficient balance
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("invalid creator address: %w", err)
	}

	if !ms.bankKeeper.HasBalance(ctx, creatorAddr, requiredAmount) {
		return sdk.Coin{}, sdkerrors.ErrInsufficientFunds.Wrapf("insufficient funds: required %s", requiredAmount)
	}

	return requiredAmount, nil
}

func (k Keeper) trustFeesAmount(totalFees int64, trustDepositRate math.LegacyDec) math.Int {
	totalFeesDec := math.LegacyNewDec(totalFees)
	trustFees := totalFeesDec.Mul(trustDepositRate)
	return trustFees.TruncateInt()
}

func (k Keeper) rewardsAmount(totalFees int64, rewardRateSum math.LegacyDec) math.Int {
	totalFeesDec := math.LegacyNewDec(totalFees)
	rewards := totalFeesDec.Mul(rewardRateSum)
	return rewards.TruncateInt()
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

func (k Keeper) trustDepositAmount(feesInDenom uint64, trustDepositRate math.LegacyDec) uint64 {
	feesInDenomDec := math.LegacyNewDec(int64(feesInDenom))
	trustDeposit := feesInDenomDec.Mul(trustDepositRate)
	return trustDeposit.TruncateInt().Uint64()
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
