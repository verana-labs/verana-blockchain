package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
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

func (ms msgServer) CreateOrUpdateCSPS(goCtx context.Context, msg *types.MsgCreateOrUpdateCSPS) (*types.MsgCreateOrUpdateCSPSResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate session existence and controller
	if err := ms.validateSessionAccess(ctx, msg); err != nil {
		return nil, err
	}

	// Validate executor permission and type
	executorPerm, err := ms.validateExecutorPerm(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Validate beneficiary permission if required
	if err := ms.validateBeneficiaryPerm(ctx, msg, executorPerm); err != nil {
		return nil, err
	}

	// TODO: [MOD-CSPS-MSG-4-2-2] & [MOD-CSPS-MSG-4-2-3]
	// Calculate fees after validation module is ready
	// This will include:
	// 1. Building permission set recursively through validation chain
	// 2. Calculating beneficiary fees based on executor type
	// 3. Calculating trust deposits and rewards

	// TODO: Implement fee processing after validation module [MOD-CSP-MSG-4-3]
	// This will include:
	// 1. Transferring fees to grantees
	// 2. Increasing trust deposits
	// 3. Processing user agent rewards

	// Create or update session
	if err := ms.createOrUpdateSession(ctx, msg, executorPerm); err != nil {
		return nil, err
	}

	return &types.MsgCreateOrUpdateCSPSResponse{}, nil
}
