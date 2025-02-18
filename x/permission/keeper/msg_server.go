package keeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// StartPermissionVP handles the MsgStartPermissionVP message
func (ms msgServer) StartPermissionVP(goCtx context.Context, msg *types.MsgStartPermissionVP) (*types.MsgStartPermissionVPResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-PERM-MSG-1-2-2] Permission checks
	validatorPerm, err := ms.validatePermissionChecks(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("permission validation failed: %w", err)
	}

	// [MOD-PERM-MSG-1-2-3] Fee checks
	fees, deposit, err := ms.validateAndCalculateFees(ctx, msg.Creator, validatorPerm)
	if err != nil {
		return nil, fmt.Errorf("fee validation failed: %w", err)
	}

	// [MOD-PERM-MSG-1-3] Execute the permission VP creation
	permID, err := ms.executeStartPermissionVP(ctx, msg, validatorPerm, fees, deposit)
	if err != nil {
		return nil, fmt.Errorf("failed to execute permission VP: %w", err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStartPermissionVP,
			sdk.NewAttribute(types.AttributeKeyPermissionID, strconv.FormatUint(permID, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyFees, strconv.FormatUint(fees, 10)),
			sdk.NewAttribute(types.AttributeKeyDeposit, strconv.FormatUint(deposit, 10)),
		),
	})

	return &types.MsgStartPermissionVPResponse{
		PermissionId: permID,
	}, nil
}

func (ms msgServer) RenewPermissionVP(goCtx context.Context, msg *types.MsgRenewPermissionVP) (*types.MsgRenewPermissionVPResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-PERM-MSG-2-2-2] Permission checks
	applicantPerm, err := ms.Keeper.GetPermissionByID(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	// Verify creator is the grantee
	if applicantPerm.Grantee != msg.Creator {
		return nil, fmt.Errorf("creator is not the permission grantee")
	}

	// Get validator permission
	validatorPerm, err := ms.Keeper.GetPermissionByID(ctx, applicantPerm.ValidatorPermId)
	if err != nil {
		return nil, fmt.Errorf("validator permission not found: %w", err)
	}

	// [MOD-PERM-MSG-2-2-3] Fee checks
	validationFees, validationDeposit, err := ms.validateAndCalculateFees(ctx, msg.Creator, validatorPerm)
	if err != nil {
		return nil, fmt.Errorf("fee validation failed: %w", err)
	}

	// [MOD-PERM-MSG-2-3] Execution
	if err := ms.executeRenewPermissionVP(ctx, applicantPerm, validationFees, validationDeposit); err != nil {
		return nil, fmt.Errorf("failed to execute permission VP renewal: %w", err)
	}

	return &types.MsgRenewPermissionVPResponse{}, nil
}

func (ms msgServer) executeRenewPermissionVP(ctx sdk.Context, perm types.Permission, fees, deposit uint64) error {
	// TODO: After trustdeposit module
	// Increment trust deposit
	//if err := ms.trustDepositKeeper.IncreaseTrustDeposit(ctx, perm.Grantee, deposit); err != nil {
	//    return fmt.Errorf("failed to increase trust deposit: %w", err)
	//}

	// Send validation fees to escrow if greater than 0
	//if fees > 0 {
	//    if err := ms.transferToEscrow(ctx, perm.Grantee, fees); err != nil {
	//        return fmt.Errorf("failed to transfer fees to escrow: %w", err)
	//    }
	//}

	now := ctx.BlockTime()

	// Update permission
	perm.VpState = types.ValidationState_VALIDATION_STATE_PENDING
	perm.VpLastStateChange = &now
	perm.Deposit += deposit
	perm.VpCurrentFees = fees
	perm.VpCurrentDeposit = deposit
	perm.Modified = &now

	// Store updated permission
	return ms.Keeper.UpdatePermission(ctx, perm)
}

func (ms msgServer) SetPermissionVPToValidated(goCtx context.Context, msg *types.MsgSetPermissionVPToValidated) (*types.MsgSetPermissionVPToValidatedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	now := ctx.BlockTime()

	// [MOD-PERM-MSG-3-2-1] Basic checks
	applicantPerm, err := ms.Keeper.GetPermissionByID(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	// Check renewal-specific constraints
	if applicantPerm.EffectiveFrom != nil {
		if msg.ValidationFees != applicantPerm.ValidationFees {
			return nil, fmt.Errorf("validation fees cannot be changed during renewal")
		}
		if msg.IssuanceFees != applicantPerm.IssuanceFees {
			return nil, fmt.Errorf("issuance fees cannot be changed during renewal")
		}
		if msg.VerificationFees != applicantPerm.VerificationFees {
			return nil, fmt.Errorf("verification fees cannot be changed during renewal")
		}
		if msg.Country != applicantPerm.Country {
			return nil, fmt.Errorf("country cannot be changed during renewal")
		}
	}

	// Check summary digest SRI
	if applicantPerm.Type == types.PermissionType_PERMISSION_TYPE_HOLDER && msg.VpSummaryDigestSri != "" {
		return nil, fmt.Errorf("vp_summary_digest_sri must be null for HOLDER type")
	}

	// [MOD-PERM-MSG-3-2-2] Validator permission checks
	validatorPerm, err := ms.Keeper.GetPermissionByID(ctx, applicantPerm.ValidatorPermId)
	if err != nil {
		return nil, fmt.Errorf("validator permission not found: %w", err)
	}

	if validatorPerm.Grantee != msg.Creator {
		return nil, fmt.Errorf("creator is not the validator")
	}

	// Get validation period and calculate expiration
	cs, err := ms.credentialSchemaKeeper.GetCredentialSchemaById(ctx, applicantPerm.SchemaId)
	if err != nil {
		return nil, fmt.Errorf("credential schema not found: %w", err)
	}

	validityPeriod := getValidityPeriod(uint32(applicantPerm.Type), cs)
	vpExp := calculateVPExp(applicantPerm.VpExp, uint64(validityPeriod), now)

	// Check effective_until if provided
	if msg.EffectiveUntil != nil {
		if applicantPerm.EffectiveUntil == nil {
			if !msg.EffectiveUntil.After(now) {
				return nil, fmt.Errorf("effective_until must be after current time")
			}
			if vpExp != nil && msg.EffectiveUntil.After(*vpExp) {
				return nil, fmt.Errorf("effective_until cannot be after validation expiration")
			}
		} else {
			if !msg.EffectiveUntil.After(*applicantPerm.EffectiveUntil) {
				return nil, fmt.Errorf("effective_until must be after current effective_until")
			}
			if vpExp != nil && msg.EffectiveUntil.After(*vpExp) {
				return nil, fmt.Errorf("effective_until cannot be after validation expiration")
			}
		}
	} else {
		msg.EffectiveUntil = vpExp
	}

	// [MOD-PERM-MSG-3-3] Execution
	if err := ms.executeSetPermissionVPToValidated(ctx, applicantPerm, msg, now, vpExp); err != nil {
		return nil, fmt.Errorf("failed to execute set to validated: %w", err)
	}

	return &types.MsgSetPermissionVPToValidatedResponse{}, nil
}

func (ms msgServer) RequestPermissionVPTermination(goCtx context.Context, msg *types.MsgRequestPermissionVPTermination) (*types.MsgRequestPermissionVPTerminationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	now := ctx.BlockTime()

	// [MOD-PERM-MSG-4-2-1] Basic checks
	applicantPerm, err := ms.Keeper.GetPermissionByID(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	if applicantPerm.VpState != types.ValidationState_VALIDATION_STATE_VALIDATED {
		return nil, fmt.Errorf("permission must be in VALIDATED state")
	}

	// Check termination authorization
	if applicantPerm.VpExp != nil && now.After(*applicantPerm.VpExp) {
		// VP has expired - either party can terminate
		validatorPerm, err := ms.Keeper.GetPermissionByID(ctx, applicantPerm.ValidatorPermId)
		if err != nil {
			return nil, fmt.Errorf("validator permission not found: %w", err)
		}
		if msg.Creator != applicantPerm.Grantee && msg.Creator != validatorPerm.Grantee {
			return nil, fmt.Errorf("only grantee or validator can terminate expired VP")
		}
	} else {
		// VP not expired - only grantee can terminate
		if msg.Creator != applicantPerm.Grantee {
			return nil, fmt.Errorf("only grantee can terminate active VP")
		}
	}

	// [MOD-PERM-MSG-4-3] Execution
	err = ms.executeRequestPermissionVPTermination(ctx, applicantPerm, msg.Creator, now)
	if err != nil {
		return nil, fmt.Errorf("failed to execute termination request: %w", err)
	}

	return &types.MsgRequestPermissionVPTerminationResponse{}, nil
}

func (ms msgServer) executeRequestPermissionVPTermination(ctx sdk.Context, perm types.Permission, terminator string, now time.Time) error {
	// Update basic fields
	perm.Modified = &now
	perm.VpTermRequested = &now
	perm.VpLastStateChange = &now

	// Set state based on conditions
	if perm.Type != types.PermissionType_PERMISSION_TYPE_HOLDER && // not HOLDER
		(perm.VpExp != nil && now.After(*perm.VpExp)) { // expired
		// Immediate termination
		perm.VpState = types.ValidationState_VALIDATION_STATE_TERMINATED
		perm.Terminated = &now
		perm.TerminatedBy = terminator

		// Handle deposits
		if err := ms.handleTerminationDeposits(ctx, &perm); err != nil {
			return fmt.Errorf("failed to handle termination deposits: %w", err)
		}
	} else {
		// Request termination
		perm.VpState = types.ValidationState_VALIDATION_STATE_TERMINATION_REQUESTED
	}

	return ms.Keeper.UpdatePermission(ctx, perm)
}

func (ms msgServer) handleTerminationDeposits(ctx sdk.Context, perm *types.Permission) error {
	// TODO: After trust deposit module is ready
	// if perm.Deposit > 0 {
	//     if err := ms.trustDepositKeeper.DecreaseTrustDeposit(ctx, perm.Grantee, perm.Deposit); err != nil {
	//         return err
	//     }
	//     perm.Deposit = 0
	// }
	//
	// if perm.ValidatorDeposit > 0 {
	//     validatorPerm, err := ms.Keeper.GetPermission(ctx, perm.ValidatorPermId)
	//     if err != nil {
	//         return err
	//     }
	//     if err := ms.trustDepositKeeper.DecreaseTrustDeposit(ctx, validatorPerm.Grantee, perm.ValidatorDeposit); err != nil {
	//         return err
	//     }
	//     perm.ValidatorDeposit = 0
	// }

	return nil
}

func (ms msgServer) ConfirmPermissionVPTermination(goCtx context.Context, msg *types.MsgConfirmPermissionVPTermination) (*types.MsgConfirmPermissionVPTerminationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	now := ctx.BlockTime()

	// Load applicant permission
	applicantPerm, err := ms.Keeper.GetPermissionByID(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	// Check permission state
	if applicantPerm.VpState != types.ValidationState_VALIDATION_STATE_TERMINATION_REQUESTED {
		return nil, fmt.Errorf("permission must be in TERMINATION_REQUESTED state")
	}

	// [MOD-PERM-MSG-5-2-2] Permission checks
	validatorPerm, err := ms.Keeper.GetPermissionByID(ctx, applicantPerm.ValidatorPermId)
	if err != nil {
		return nil, fmt.Errorf("validator permission not found: %w", err)
	}

	// Calculate timeout
	termRequestTimeout := applicantPerm.VpTermRequested.AddDate(0, 0, int(ms.Keeper.GetParams(ctx).ValidationTermRequestedTimeoutDays))
	timeoutReached := now.After(termRequestTimeout)

	// Check authorization
	if !timeoutReached {
		// Before timeout: only validator can confirm
		if msg.Creator != validatorPerm.Grantee {
			return nil, fmt.Errorf("only validator can confirm termination before timeout")
		}
	} else {
		// After timeout: either validator or applicant can confirm
		if msg.Creator != validatorPerm.Grantee && msg.Creator != applicantPerm.Grantee {
			return nil, fmt.Errorf("only validator or applicant can confirm termination after timeout")
		}
	}

	// [MOD-PERM-MSG-5-3] Execution
	if err := ms.executeConfirmPermissionVPTermination(ctx, applicantPerm, validatorPerm, msg.Creator, now); err != nil {
		return nil, fmt.Errorf("failed to execute termination confirmation: %w", err)
	}

	return &types.MsgConfirmPermissionVPTerminationResponse{}, nil
}

func (ms msgServer) executeConfirmPermissionVPTermination(ctx sdk.Context, applicantPerm types.Permission, validatorPerm types.Permission, confirmer string, now time.Time) error {
	// Update basic fields
	applicantPerm.Modified = &now
	applicantPerm.VpState = types.ValidationState_VALIDATION_STATE_TERMINATED
	applicantPerm.VpLastStateChange = &now
	applicantPerm.Terminated = &now
	applicantPerm.TerminatedBy = confirmer

	// Handle deposits based on who confirmed
	if applicantPerm.Deposit > 0 {
		// TODO: After trust deposit module implementation
		// if err := ms.trustDepositKeeper.DecreaseTrustDeposit(ctx, applicantPerm.Grantee, applicantPerm.Deposit); err != nil {
		//     return fmt.Errorf("failed to decrease applicant trust deposit: %w", err)
		// }
		applicantPerm.Deposit = 0
	}

	// Only return validator deposit if validator confirmed
	if confirmer == validatorPerm.Grantee && applicantPerm.VpValidatorDeposit > 0 {
		// TODO: After trust deposit module implementation
		// if err := ms.trustDepositKeeper.DecreaseTrustDeposit(ctx, validatorPerm.Grantee, applicantPerm.ValidatorDeposit); err != nil {
		//     return fmt.Errorf("failed to decrease validator trust deposit: %w", err)
		// }
		applicantPerm.VpValidatorDeposit = 0
	}

	// Persist changes
	return ms.Keeper.UpdatePermission(ctx, applicantPerm)
}

func (ms msgServer) CancelPermissionVPLastRequest(goCtx context.Context, msg *types.MsgCancelPermissionVPLastRequest) (*types.MsgCancelPermissionVPLastRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Load applicant permission
	applicantPerm, err := ms.Keeper.GetPermissionByID(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	// Check if creator is the grantee
	if applicantPerm.Grantee != msg.Creator {
		return nil, fmt.Errorf("creator is not the permission grantee")
	}

	// Check permission state
	if applicantPerm.VpState != types.ValidationState_VALIDATION_STATE_PENDING {
		return nil, fmt.Errorf("permission must be in PENDING state")
	}

	// [MOD-PERM-MSG-6-3] Execution
	if err := ms.executeCancelPermissionVPLastRequest(ctx, applicantPerm); err != nil {
		return nil, fmt.Errorf("failed to execute VP cancellation: %w", err)
	}

	return &types.MsgCancelPermissionVPLastRequestResponse{}, nil
}

func (ms msgServer) executeCancelPermissionVPLastRequest(ctx sdk.Context, perm types.Permission) error {
	now := ctx.BlockTime()

	// Update basic fields
	perm.Modified = &now
	perm.VpLastStateChange = &now

	// Set state based on vp_exp
	if perm.VpExp == nil {
		perm.VpState = types.ValidationState_VALIDATION_STATE_TERMINATED
	} else {
		perm.VpState = types.ValidationState_VALIDATION_STATE_VALIDATED
	}

	// Handle current fees if any
	if perm.VpCurrentFees > 0 {
		// TODO: After bank module integration
		// Transfer fees back from escrow
		// if err := ms.bankKeeper.SendCoinsFromModuleToAccount(
		//     ctx,
		//     types.ModuleName,
		//     sdk.AccAddress(perm.Grantee),
		//     sdk.NewCoins(sdk.NewCoin(ms.Keeper.GetParams(ctx).FeeDenom, sdk.NewInt(int64(perm.VpCurrentFees)))),
		// ); err != nil {
		//     return fmt.Errorf("failed to refund fees: %w", err)
		// }
		perm.VpCurrentFees = 0
	}

	// Handle current deposit if any
	if perm.VpCurrentDeposit > 0 {
		// TODO: After trust deposit module integration
		// if err := ms.trustDepositKeeper.DecreaseTrustDeposit(ctx, perm.Grantee, perm.VpCurrentDeposit); err != nil {
		//     return fmt.Errorf("failed to decrease trust deposit: %w", err)
		// }
		perm.VpCurrentDeposit = 0
	}

	// Persist changes
	return ms.Keeper.UpdatePermission(ctx, perm)
}

func (ms msgServer) CreateRootPermission(goCtx context.Context, msg *types.MsgCreateRootPermission) (*types.MsgCreateRootPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	now := ctx.BlockTime()

	// Check credential schema exists
	_, err := ms.credentialSchemaKeeper.GetCredentialSchemaById(ctx, msg.SchemaId)
	if err != nil {
		return nil, fmt.Errorf("credential schema not found: %w", err)
	}

	// [MOD-PERM-MSG-7-2-2] Permission checks
	if err := ms.validateCreateRootPermissionAuthority(ctx, msg); err != nil {
		return nil, err
	}

	// [MOD-PERM-MSG-7-3] Execution
	id, err := ms.executeCreateRootPermission(ctx, msg, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create root permission: %w", err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateRootPermission,
			sdk.NewAttribute(types.AttributeKeyRootPermissionID, strconv.FormatUint(id, 10)),
			sdk.NewAttribute(types.AttributeKeySchemaID, strconv.FormatUint(msg.SchemaId, 10)),
			sdk.NewAttribute(types.AttributeKeyTimestamp, now.String()),
		),
	})

	return &types.MsgCreateRootPermissionResponse{
		Id: id,
	}, nil
}

func (ms msgServer) validateCreateRootPermissionAuthority(ctx sdk.Context, msg *types.MsgCreateRootPermission) error {
	// Get credential schema
	cs, err := ms.credentialSchemaKeeper.GetCredentialSchemaById(ctx, msg.SchemaId)
	if err != nil {
		return fmt.Errorf("credential schema not found: %w", err)
	}

	// Load trust registry
	tr, err := ms.trustRegistryKeeper.GetTrustRegistry(ctx, cs.TrId)
	if err != nil {
		return fmt.Errorf("trust registry not found: %w", err)
	}

	// Check if creator is the controller
	if tr.Controller != msg.Creator {
		return fmt.Errorf("creator is not the trust registry controller")
	}

	return nil
}

func (ms msgServer) executeCreateRootPermission(ctx sdk.Context, msg *types.MsgCreateRootPermission, now time.Time) (uint64, error) {
	// Create new permission
	perm := types.Permission{
		SchemaId:         msg.SchemaId,
		Type:             types.PermissionType_PERMISSION_TYPE_TRUST_REGISTRY,
		Did:              msg.Did,
		Grantee:          msg.Creator,
		Created:          &now,
		CreatedBy:        msg.Creator,
		Modified:         &now,
		EffectiveFrom:    msg.EffectiveFrom,
		EffectiveUntil:   msg.EffectiveUntil,
		Country:          msg.Country,
		ValidationFees:   msg.ValidationFees,
		IssuanceFees:     msg.IssuanceFees,
		VerificationFees: msg.VerificationFees,
		Deposit:          0,
		VpState:          types.ValidationState_VALIDATION_STATE_VALIDATED,
	}

	// Store the permission
	id, err := ms.Keeper.CreatePermission(ctx, perm)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission: %w", err)
	}

	return id, nil
}

func (ms msgServer) ExtendPermission(goCtx context.Context, msg *types.MsgExtendPermission) (*types.MsgExtendPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	now := ctx.BlockTime()

	// [MOD-PERM-MSG-8-2-1] Basic checks
	applicantPerm, err := ms.Keeper.GetPermissionByID(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	// Check effective_until is after current effective_until
	if applicantPerm.EffectiveUntil != nil && !msg.EffectiveUntil.After(*applicantPerm.EffectiveUntil) {
		return nil, fmt.Errorf("effective_until must be after current effective_until")
	}

	// Check effective_until is before or equal to vp_exp
	if applicantPerm.VpExp != nil && msg.EffectiveUntil.After(*applicantPerm.VpExp) {
		return nil, fmt.Errorf("effective_until cannot be after validation expiration")
	}

	// [MOD-PERM-MSG-8-2-2] Validator permission checks
	if err := ms.validateExtendPermissionAuthority(ctx, msg.Creator, applicantPerm); err != nil {
		return nil, err
	}

	// [MOD-PERM-MSG-8-3] Execution
	if err := ms.executeExtendPermission(ctx, applicantPerm, msg.Creator, msg.EffectiveUntil, now); err != nil {
		return nil, fmt.Errorf("failed to extend permission: %w", err)
	}

	return &types.MsgExtendPermissionResponse{}, nil
}

func (ms msgServer) validateExtendPermissionAuthority(ctx sdk.Context, creator string, perm types.Permission) error {
	if perm.ValidatorPermId == 0 {
		// For TRUST_REGISTRY type, creator must be the grantee
		if perm.Type == types.PermissionType_PERMISSION_TYPE_TRUST_REGISTRY {
			if perm.Grantee != creator {
				return fmt.Errorf("creator is not the permission grantee")
			}
		} else {
			return fmt.Errorf("invalid permission type for extension")
		}
	} else {
		// For other types, creator must be the validator
		validatorPerm, err := ms.Keeper.GetPermissionByID(ctx, perm.ValidatorPermId)
		if err != nil {
			return fmt.Errorf("validator permission not found: %w", err)
		}
		if validatorPerm.Grantee != creator {
			return fmt.Errorf("creator is not the validator")
		}
	}
	return nil
}

func (ms msgServer) executeExtendPermission(ctx sdk.Context, perm types.Permission, creator string, effectiveUntil *time.Time, now time.Time) error {
	perm.EffectiveUntil = effectiveUntil
	perm.Extended = &now
	perm.Modified = &now
	perm.ExtendedBy = creator

	return ms.Keeper.UpdatePermission(ctx, perm)
}

// RevokePermission handles the MsgRevokePermission message
func (ms msgServer) RevokePermission(goCtx context.Context, msg *types.MsgRevokePermission) (*types.MsgRevokePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	now := ctx.BlockTime()

	// [MOD-PERM-MSG-9-2-1] Basic checks
	applicantPerm, err := ms.Keeper.GetPermissionByID(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	// [MOD-PERM-MSG-9-2-2] Validator permission checks
	validatorPerm, err := ms.Keeper.GetPermissionByID(ctx, applicantPerm.ValidatorPermId)
	if err != nil {
		return nil, fmt.Errorf("validator permission not found: %w", err)
	}

	if validatorPerm.Grantee != msg.Creator {
		return nil, fmt.Errorf("creator is not the validator")
	}

	// [MOD-PERM-MSG-9-3] Execution
	if err := ms.executeRevokePermission(ctx, applicantPerm, msg.Creator, now); err != nil {
		return nil, fmt.Errorf("failed to revoke permission: %w", err)
	}

	return &types.MsgRevokePermissionResponse{}, nil
}

func (ms msgServer) executeRevokePermission(ctx sdk.Context, perm types.Permission, creator string, now time.Time) error {
	perm.Revoked = &now
	perm.Modified = &now
	perm.RevokedBy = creator

	return ms.Keeper.UpdatePermission(ctx, perm)
}

type PermissionSet []types.Permission

func (ps PermissionSet) contains(id uint64) bool {
	for _, perm := range ps {
		if perm.Id == id {
			return true
		}
	}
	return false
}

func (ms msgServer) CreateOrUpdatePermissionSession(goCtx context.Context, msg *types.MsgCreateOrUpdatePermissionSession) (*types.MsgCreateOrUpdatePermissionSessionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	now := ctx.BlockTime()

	// Validate session access
	if err := ms.validateSessionAccess(ctx, msg); err != nil {
		return nil, err
	}

	// Validate executor permission
	executorPerm, err := ms.validateExecutorPermission(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Validate agent permission
	if err := ms.validateAgentPermission(ctx, msg); err != nil {
		return nil, err
	}

	// Validate wallet agent permission if provided
	if err := ms.validateWalletAgentPermission(ctx, msg); err != nil {
		return nil, err
	}

	// Build permission set and calculate fees
	_, err = ms.buildPermissionSet(ctx, executorPerm, msg.BeneficiaryPermId)
	if err != nil {
		return nil, fmt.Errorf("failed to build permission set: %w", err)
	}

	// Calculate and validate fees
	//fees, err := ms.calculateAndValidateFees(ctx, msg.Creator, permSet, executorPerm.Type)
	//if err != nil {
	//	return nil, err
	//}

	// Process fees and create/update session
	//if err := ms.processFees(ctx, msg.Creator, permSet, executorPerm.Type, fees); err != nil {
	//	return nil, err
	//}

	// Create or update session
	if err := ms.createOrUpdateSession(ctx, msg, now); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateOrUpdatePermissionSession,
			sdk.NewAttribute(types.AttributeKeySessionID, msg.Id),
			sdk.NewAttribute(types.AttributeKeyExecutorPermID, strconv.FormatUint(msg.ExecutorPermId, 10)),
			sdk.NewAttribute(types.AttributeKeyBeneficiaryPermID, strconv.FormatUint(msg.BeneficiaryPermId, 10)),
			sdk.NewAttribute(types.AttributeKeyAgentPermID, strconv.FormatUint(msg.AgentPermId, 10)),
			sdk.NewAttribute(types.AttributeKeyTimestamp, now.String()),
		),
	})

	return &types.MsgCreateOrUpdatePermissionSessionResponse{
		Id: msg.Id,
	}, nil
}
