package keeper

import (
	"errors"

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
		if authz.ExecutorPermId == msg.ExecutorPermId &&
			authz.BeneficiaryPermId == msg.BeneficiaryPermId &&
			authz.WalletAgentPermId == msg.WalletAgentPermId {
			return sdkerrors.ErrInvalidRequest.Wrap("authorization already exists")
		}
	}

	return nil
}

func (ms msgServer) validateExecutorPermission(ctx sdk.Context, msg *types.MsgCreateOrUpdatePermissionSession) (*types.Permission, error) {
	executorPerm, err := ms.Permission.Get(ctx, msg.ExecutorPermId)
	if err != nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("executor permission not found: %v", err)
	}

	// Executor must be ISSUER or VERIFIER
	if executorPerm.Type != types.PermissionType_PERMISSION_TYPE_ISSUER &&
		executorPerm.Type != types.PermissionType_PERMISSION_TYPE_VERIFIER {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("executor must be ISSUER or VERIFIER")
	}

	// Must be valid permission (not revoked/terminated)
	if executorPerm.Revoked != nil || executorPerm.Terminated != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("executor permission is revoked or terminated")
	}

	return &executorPerm, nil
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

//func (ms msgServer) calculateAndValidateFees(ctx sdk.Context, creator string, permSet PermissionSet, executorType types.PermissionType) (sdk.Coin, error) {
//	beneficiaryFees := math.NewInt(0)
//
//	// Calculate total beneficiary fees
//	for _, perm := range permSet {
//		if executorType == types.PermissionType_PERMISSION_TYPE_VERIFIER {
//			beneficiaryFees = beneficiaryFees.Add(math.NewInt(int64(perm.VerificationFees)))
//		} else {
//			beneficiaryFees = beneficiaryFees.Add(math.NewInt(int64(perm.IssuanceFees)))
//		}
//	}
//
//	// Get global variables
//	trustUnitPrice := ms.globalVariablesKeeper.GetTrustUnitPrice(ctx)
//	trustDepositRate := ms.globalVariablesKeeper.GetTrustDepositRate(ctx)
//	userAgentRewardRate := ms.globalVariablesKeeper.GetUserAgentRewardRate(ctx)
//	walletUserAgentRewardRate := ms.globalVariablesKeeper.GetWalletUserAgentRewardRate(ctx)
//
//	// Calculate total fees including trust deposit and rewards
//	totalFees := beneficiaryFees.Mul(trustUnitPrice)
//	trustFees := totalFees.Mul(trustDepositRate)
//	rewards := totalFees.Mul(userAgentRewardRate.Add(walletUserAgentRewardRate))
//
//	requiredAmount := sdk.NewCoin(ms.stakingKeeper.BondDenom(ctx), totalFees.Add(trustFees).Add(rewards))
//
//	// Validate sufficient balance
//	creatorAddr, _ := sdk.AccAddressFromBech32(creator)
//	if !ms.bankKeeper.HasBalance(ctx, creatorAddr, requiredAmount) {
//		return sdk.Coin{}, sdkerrors.ErrInsufficientFunds.Wrapf("insufficient funds: required %s", requiredAmount)
//	}
//
//	return requiredAmount, nil
//}

//func (ms msgServer) processFees(ctx sdk.Context, creator string, permSet PermissionSet, executorType types.PermissionType, totalFees sdk.Coin) error {
//	creatorAddr, _ := sdk.AccAddressFromBech32(creator)
//
//	for _, perm := range permSet {
//		fees := math.NewInt(0)
//		if executorType == types.PermissionType_PERMISSION_TYPE_VERIFIER {
//			fees = math.NewInt(int64(perm.VerificationFees))
//		} else {
//			fees = math.NewInt(int64(perm.IssuanceFees))
//		}
//
//		if fees.IsPositive() {
//			trustUnitPrice := ms.globalVariablesKeeper.GetTrustUnitPrice(ctx)
//			trustDepositRate := ms.globalVariablesKeeper.GetTrustDepositRate(ctx)
//
//			// Calculate direct fees (excluding trust deposit)
//			directFees := fees.Mul(trustUnitPrice).Mul(sdk.NewInt(1).Sub(trustDepositRate))
//
//			// Transfer direct fees to grantee
//			if err := ms.bankKeeper.SendCoins(
//				ctx,
//				creatorAddr,
//				sdk.AccAddressFromBech32(perm.Grantee),
//				sdk.NewCoins(sdk.NewCoin(ms.stakingKeeper.BondDenom(ctx), directFees)),
//			); err != nil {
//				return err
//			}
//
//			// TODO: Implement trust deposit increases when module is available
//		}
//	}
//
//	return nil
//}

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
		ExecutorPermId:    msg.ExecutorPermId,
		BeneficiaryPermId: msg.BeneficiaryPermId,
		WalletAgentPermId: msg.WalletAgentPermId,
	})

	return ms.PermissionSession.Set(ctx, msg.Id, *session)
}
