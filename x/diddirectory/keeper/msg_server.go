package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/diddirectory/types"
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

func (ms msgServer) AddDID(goCtx context.Context, msg *types.MsgAddDID) (*types.MsgAddDIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic parameter validation
	if err := ms.validateAddDIDParams(ctx, msg); err != nil {
		return nil, err
	}

	// Fee checks
	if err := ms.checkSufficientFees(ctx, msg.Creator, msg.Years); err != nil {
		return nil, err
	}

	// Execute the addition
	if err := ms.executeAddDID(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgAddDIDResponse{}, nil
}
