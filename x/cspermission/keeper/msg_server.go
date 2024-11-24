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

func (ms msgServer) RevokeCredentialSchemaPerm(ctx context.Context, msg *types.MsgRevokeCredentialSchemaPerm) (*types.MsgRevokeCredentialSchemaPermResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	csp, err := ms.CredentialSchemaPerm.Get(sdkCtx, msg.Id)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrNotFound, "permission not found: %d", msg.Id)
	}

	if csp.Revoked != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "permission is already revoked")
	}

	cs, err := ms.credentialSchemaKeeper.GetCredentialSchemaById(sdkCtx, csp.SchemaId)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrNotFound, "credential schema not found: %d", csp.SchemaId)
	}

	tr, err := ms.trustRegistryKeeper.GetTrustRegistry(sdkCtx, cs.TrId)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrNotFound, "trust registry not found: %d", cs.TrId)
	}

	if err := ms.validateRevokePermissions(sdkCtx, msg.Creator, &csp, cs, tr); err != nil {
		return nil, err
	}

	revokedTime := sdkCtx.BlockTime()
	csp.Revoked = &revokedTime
	csp.RevokedBy = msg.Creator

	if csp.Deposit > 0 {
		// Handle deposit decrease - implement after trust deposit module
		csp.Deposit = 0
	}

	if err := ms.CredentialSchemaPerm.Set(sdkCtx, msg.Id, csp); err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "failed to update permission")
	}

	return &types.MsgRevokeCredentialSchemaPermResponse{}, nil
}

func (ms msgServer) TerminateCredentialSchemaPerm(ctx context.Context, msg *types.MsgTerminateCredentialSchemaPerm) (*types.MsgTerminateCredentialSchemaPermResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Get the permission
	csp, err := ms.CredentialSchemaPerm.Get(sdkCtx, msg.Id)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrNotFound, "permission not found: %d", msg.Id)
	}

	// Check if already terminated
	if csp.Terminated != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "permission is already terminated")
	}

	// Check grantee is the one terminating
	if csp.Grantee != msg.Creator {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "only grantee can terminate permission")
	}

	// TODO: Check validation state if validation exists

	// Set termination details
	terminatedTime := sdkCtx.BlockTime()
	csp.Terminated = &terminatedTime
	csp.TerminatedBy = msg.Creator

	// Handle deposit
	if csp.Deposit > 0 {
		// TODO: Implement trust deposit decrease
		csp.Deposit = 0
	}

	// Update permission
	if err := ms.CredentialSchemaPerm.Set(sdkCtx, msg.Id, csp); err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "failed to update permission")
	}

	return &types.MsgTerminateCredentialSchemaPermResponse{}, nil
}
