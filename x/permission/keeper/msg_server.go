package keeper

import (
	"context"
	"fmt"
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

	return &types.MsgStartPermissionVPResponse{
		PermissionId: permID,
	}, nil
}

func (ms msgServer) RenewPermissionVP(goCtx context.Context, msg *types.MsgRenewPermissionVP) (*types.MsgRenewPermissionVPResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-PERM-MSG-2-2-2] Permission checks
	applicantPerm, err := ms.Keeper.GetPermission(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	// Verify creator is the grantee
	if applicantPerm.Grantee != msg.Creator {
		return nil, fmt.Errorf("creator is not the permission grantee")
	}

	// Get validator permission
	validatorPerm, err := ms.Keeper.GetPermission(ctx, applicantPerm.ValidatorPermId)
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
