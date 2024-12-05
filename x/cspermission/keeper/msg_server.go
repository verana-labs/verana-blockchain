package keeper

import (
	"context"
	"cosmossdk.io/collections"
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

func (ms msgServer) validateSessionAccess(ctx sdk.Context, msg *types.MsgCreateOrUpdateCSPS) error {
	existingSession, err := ms.GetCSPS(ctx, msg.Id)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return err
	}

	if existingSession != nil {
		if existingSession.Controller != msg.Creator {
			return sdkerrors.ErrUnauthorized.Wrap("only session controller can update")
		}

		// Check for duplicate permission pairs
		for _, authz := range existingSession.SessionAuthz {
			if authz.ExecutorPermId == msg.ExecutorPermId &&
				authz.BeneficiaryPermId == msg.BeneficiaryPermId {
				return sdkerrors.ErrInvalidRequest.Wrap("permission pair already exists")
			}
		}
	}

	return nil
}

func (ms msgServer) validateExecutorPerm(ctx sdk.Context, msg *types.MsgCreateOrUpdateCSPS) (*types.CredentialSchemaPerm, error) {
	executorPerm, err := ms.CredentialSchemaPerm.Get(ctx, msg.ExecutorPermId)
	if err != nil {
		return nil, sdkerrors.ErrNotFound.Wrap("executor permission not found")
	}

	if executorPerm.Grantee != msg.Creator {
		return nil, sdkerrors.ErrUnauthorized.Wrap("only grantee can create/update session")
	}

	if executorPerm.CspType != types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER &&
		executorPerm.CspType != types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("executor must be ISSUER or VERIFIER")
	}

	return &executorPerm, nil
}

func (ms msgServer) validateBeneficiaryPerm(ctx sdk.Context, msg *types.MsgCreateOrUpdateCSPS, executorPerm *types.CredentialSchemaPerm) error {
	if executorPerm.CspType == types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER {
		if msg.BeneficiaryPermId != 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("beneficiary not allowed for ISSUER")
		}
		return nil
	}

	// VERIFIER type checks
	if msg.BeneficiaryPermId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("beneficiary required for VERIFIER")
	}

	beneficiaryPerm, err := ms.CredentialSchemaPerm.Get(ctx, msg.BeneficiaryPermId)
	if err != nil {
		return sdkerrors.ErrNotFound.Wrap("beneficiary permission not found")
	}

	if beneficiaryPerm.CspType != types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER {
		return sdkerrors.ErrInvalidRequest.Wrap("beneficiary must be ISSUER")
	}

	if beneficiaryPerm.SchemaId != executorPerm.SchemaId {
		return sdkerrors.ErrInvalidRequest.Wrap("schema IDs must match")
	}

	return nil
}

func (ms msgServer) createOrUpdateSession(ctx sdk.Context, msg *types.MsgCreateOrUpdateCSPS, executorPerm *types.CredentialSchemaPerm) error {
	session := &types.CredentialSchemaPermSession{
		Id:           msg.Id,
		Controller:   msg.Creator,
		UserAgentDid: msg.UserAgentDid,
	}

	existingSession, err := ms.GetCSPS(ctx, msg.Id)
	if err == nil {
		session = existingSession
	}

	// Add new authorization
	session.SessionAuthz = append(session.SessionAuthz, &types.SessionAuthz{
		ExecutorPermId:     msg.ExecutorPermId,
		BeneficiaryPermId:  msg.BeneficiaryPermId,
		WalletUserAgentDid: msg.WalletUserAgentDid,
	})

	return ms.CredentialSchemaPermSession.Set(ctx, msg.Id, *session)
}
