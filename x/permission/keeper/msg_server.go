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
