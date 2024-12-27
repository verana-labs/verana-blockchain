package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/validation/types"
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

func (ms msgServer) CreateValidation(goCtx context.Context, msg *types.MsgCreateValidation) (*types.MsgCreateValidationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-V-MSG-1-2-2] Create New Validation permission checks
	if err := ms.validatePermissions(ctx, msg); err != nil {
		return nil, fmt.Errorf("permission check failed: %w", err)
	}

	// [MOD-V-MSG-1-2-3] Create New Validation fee checks
	fees, deposit, err := ms.checkAndCalculateFees(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("fee check failed: %w", err)
	}

	// [MOD-V-MSG-1-3] Create New Validation execution
	validation, err := ms.executeCreateValidation(ctx, msg, fees, deposit)
	if err != nil {
		return nil, fmt.Errorf("failed to create validation: %w", err)
	}

	return &types.MsgCreateValidationResponse{
		ValidationId: validation.Id,
	}, nil
}

func (ms msgServer) RenewValidation(goCtx context.Context, msg *types.MsgRenewValidation) (*types.MsgRenewValidationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-V-MSG-2-2-1] Basic checks and load validation
	val, validatorPermID, err := ms.validateRenewalBasics(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("basic validation check failed: %w", err)
	}

	// [MOD-V-MSG-2-2-2] Revocation and permission checks
	if err := ms.validateRenewalPermissions(ctx, msg, val, validatorPermID); err != nil {
		return nil, fmt.Errorf("permission check failed: %w", err)
	}

	// [MOD-V-MSG-2-2-3] Fee checks
	fees, deposit, err := ms.checkAndCalculateRenewalFees(ctx, validatorPermID)
	if err != nil {
		return nil, fmt.Errorf("fee check failed: %w", err)
	}

	// [MOD-V-MSG-2-3] Execute renewal
	if err := ms.executeRenewalValidation(ctx, val, validatorPermID, fees, deposit); err != nil {
		return nil, fmt.Errorf("failed to execute renewal: %w", err)
	}

	return &types.MsgRenewValidationResponse{
		ValidationId: msg.Id,
	}, nil
}

func (ms msgServer) SetValidated(goCtx context.Context, msg *types.MsgSetValidated) (*types.MsgSetValidatedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-V-MSG-3-2-1] Basic checks
	val, err := ms.Validation.Get(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("validation not found: %w", err)
	}

	if val.State != types.ValidationState_PENDING {
		return nil, fmt.Errorf("validation must be in PENDING state")
	}

	// Validate summary hash if provided
	if msg.SummaryHash != "" {
		if val.Type == types.ValidationType_HOLDER {
			return nil, fmt.Errorf("summary hash must be null for HOLDER type validations")
		}
	}

	// [MOD-V-MSG-3-2-2] Validator permission checks
	perm, err := ms.csPermissionKeeper.GetCSPermission(ctx, val.ValidatorPermId)
	if err != nil {
		return nil, fmt.Errorf("validator permission not found: %w", err)
	}

	if perm.Grantee != msg.Creator {
		return nil, fmt.Errorf("only the validator can set validation to validated")
	}

	// [MOD-V-MSG-3-3] Execute validation
	now := ctx.BlockTime()

	// Update validation state
	val.State = types.ValidationState_VALIDATED
	val.LastStateChange = now
	val.SummaryHash = msg.SummaryHash

	// TODO: Handle fees and deposits validatorTrustFees & validatorTrustDeposit

	// TODO: Transfer fees from escrow (module accounts) to validator
	// TODO: Update trust deposits using TrustDeposit module

	// Update validation
	val.CurrentFees = 0
	val.CurrentDeposit = 0

	// Set expiration based on validation type
	// TODO: Get validity periods from params and set val.Exp accordingly

	if err := ms.Validation.Set(ctx, val.Id, val); err != nil {
		return nil, fmt.Errorf("failed to update validation: %w", err)
	}

	return &types.MsgSetValidatedResponse{
		ValidationId: msg.Id,
	}, nil
}
