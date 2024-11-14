package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/credentialschema/types"
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

func (ms msgServer) CreateCredentialSchema(goCtx context.Context, msg *types.MsgCreateCredentialSchema) (*types.MsgCreateCredentialSchemaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Generate next ID
	nextID, err := ms.GetNextID(ctx, "cs")
	if err != nil {
		return nil, fmt.Errorf("failed to generate schema ID: %w", err)
	}

	// [MOD-CS-MSG-1-2-1] Basic checks
	if err := ms.validateCreateCredentialSchemaParams(ctx, msg); err != nil {
		return nil, err
	}

	// [MOD-CS-MSG-1-2-2] Fee checks
	//if err := ms.checkSufficientFees(ctx, msg.Creator); err != nil {
	// return nil, err
	//}

	// [MOD-CS-MSG-1-3] Execution
	if err := ms.executeCreateCredentialSchema(ctx, nextID, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreateCredentialSchemaResponse{
		Id: nextID,
	}, nil
}
