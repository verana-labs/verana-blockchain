package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
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

func (ms msgServer) CreateTrustRegistry(goCtx context.Context, msg *types.MsgCreateTrustRegistry) (*types.MsgCreateTrustRegistryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-TR-MSG-1-2-1] Create New Trust Registry basic checks
	if err := ms.validateCreateTrustRegistryParams(ctx, msg); err != nil {
		return nil, err
	}

	// [MOD-TR-MSG-1-2-2] Create New Trust Registry fee checks
	if err := ms.checkSufficientFees(ctx, msg.Creator); err != nil {
		return nil, err
	}

	// [MOD-TR-MSG-1-3] Create New Trust Registry execution
	now := ctx.BlockTime()
	tr, gfv, gfd, err := ms.createTrustRegistryEntries(ctx, msg, now)
	if err != nil {
		return nil, err
	}

	if err := ms.persistEntries(ctx, tr, gfv, gfd); err != nil {
		return nil, err
	}

	return &types.MsgCreateTrustRegistryResponse{}, nil
}

func (ms msgServer) AddGovernanceFrameworkDocument(goCtx context.Context, msg *types.MsgAddGovernanceFrameworkDocument) (*types.MsgAddGovernanceFrameworkDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := ms.validateAddGovernanceFrameworkDocumentParams(ctx, msg); err != nil {
		return nil, err
	}

	if err := ms.checkSufficientFees(ctx, msg.Creator); err != nil {
		return nil, err
	}

	if err := ms.executeAddGovernanceFrameworkDocument(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgAddGovernanceFrameworkDocumentResponse{}, nil
}

func (ms msgServer) IncreaseActiveGovernanceFrameworkVersion(goCtx context.Context, msg *types.MsgIncreaseActiveGovernanceFrameworkVersion) (*types.MsgIncreaseActiveGovernanceFrameworkVersionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate parameters
	if err := ms.validateIncreaseActiveGovernanceFrameworkVersionParams(ctx, msg); err != nil {
		return nil, err
	}

	// Check fees
	if err := ms.checkSufficientFees(ctx, msg.Creator); err != nil {
		return nil, err
	}

	// Execute the increase
	if err := ms.executeIncreaseActiveGovernanceFrameworkVersion(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgIncreaseActiveGovernanceFrameworkVersionResponse{}, nil
}

func (ms msgServer) UpdateTrustRegistry(goCtx context.Context, msg *types.MsgUpdateTrustRegistry) (*types.MsgUpdateTrustRegistryResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Get trust registry
    tr, err := ms.TrustRegistry.Get(ctx, msg.Id)
    if err != nil {
        return nil, fmt.Errorf("trust registry not found: %w", err)
    }

    // Check controller
    if tr.Controller != msg.Creator {
        return nil, fmt.Errorf("only trust registry controller can update trust registry")
    }

    // Update fields
    tr.Did = msg.Did
    tr.Aka = msg.Aka
    tr.Modified = ctx.BlockTime()

    // Save updated trust registry
    if err := ms.TrustRegistry.Set(ctx, tr.Id, tr); err != nil {
        return nil, fmt.Errorf("failed to update trust registry: %w", err)
    }

    return &types.MsgUpdateTrustRegistryResponse{}, nil
} 