package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
	"time"
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

func (ms msgServer) CreateCredentialSchemaPerm(goCtx context.Context, msg *types.MsgCreateCredentialSchemaPerm) (*types.MsgCreateCredentialSchemaPermResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the credential schema
	cs, err := ms.credentialSchemaKeeper.GetCredentialSchemaById(ctx, msg.SchemaId)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrNotFound, "credential schema not found: %d", msg.SchemaId)
	}

	// Get the trust registry
	tr, err := ms.trustRegistryKeeper.GetTrustRegistry(ctx, cs.TrId)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrNotFound, "trust registry not found: %d", cs.TrId)
	}

	if !msg.EffectiveFrom.After(time.Now()) {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "effective_from must be in the future")
	}

	// Validate permissions based on type
	if err := ms.validatePermissions(ctx, msg, cs, tr); err != nil {
		return nil, err
	}

	// Check for overlapping permissions
	if err := ms.checkOverlappingPermissions(ctx, msg); err != nil {
		return nil, err
	}

	// Create the permission
	if err := ms.createPermission(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreateCredentialSchemaPermResponse{}, nil
}
