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
