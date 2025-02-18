package keeper

import (
	"context"
	errors2 "errors"
	"fmt"
	"regexp"
	"sort"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	credentialschematypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ListPermissions(goCtx context.Context, req *types.QueryListPermissionsRequest) (*types.QueryListPermissionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-PERM-QRY-1-2] Checks
	// Validate response_max_size
	if req.ResponseMaxSize == 0 {
		req.ResponseMaxSize = 64 // Default value
	}
	if req.ResponseMaxSize < 1 || req.ResponseMaxSize > 1024 {
		return nil, status.Error(codes.InvalidArgument, "response_max_size must be between 1 and 1,024")
	}

	var permissions []types.Permission

	// [MOD-PERM-QRY-1-3] Execution
	// Collect all matching permissions
	err := k.Permission.Walk(ctx, nil, func(key uint64, perm types.Permission) (bool, error) {
		// Apply modified_after filter if provided
		if req.ModifiedAfter != nil && !perm.Modified.After(*req.ModifiedAfter) {
			return false, nil
		}

		permissions = append(permissions, perm)
		return len(permissions) >= int(req.ResponseMaxSize), nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Sort by modified time ascending
	sort.Slice(permissions, func(i, j int) bool {
		return permissions[i].Modified.Before(*permissions[j].Modified)
	})

	return &types.QueryListPermissionsResponse{
		Permissions: permissions,
	}, nil
}

func (k Keeper) GetPermission(goCtx context.Context, req *types.QueryGetPermissionRequest) (*types.QueryGetPermissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-PERM-QRY-2-2] Checks
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "permission ID cannot be 0")
	}

	// [MOD-PERM-QRY-2-3] Execution
	permission, err := k.Permission.Get(ctx, req.Id)
	if err != nil {
		if errors2.Is(collections.ErrNotFound, err) {
			return nil, status.Error(codes.NotFound, "permission not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get permission: %v", err))
	}

	return &types.QueryGetPermissionResponse{
		Permission: permission,
	}, nil
}

func (k Keeper) GetPermissionSession(ctx context.Context, req *types.QueryGetPermissionSessionRequest) (*types.QueryGetPermissionSessionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "session ID is required")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	session, err := k.PermissionSession.Get(sdkCtx, req.Id)
	if err != nil {
		if errors2.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "session not found")
		}
		return nil, status.Error(codes.Internal, "failed to get session")
	}

	return &types.QueryGetPermissionSessionResponse{
		Session: &session,
	}, nil
}

func (k Keeper) ListPermissionSessions(ctx context.Context, req *types.QueryListPermissionSessionsRequest) (*types.QueryListPermissionSessionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// Validate response_max_size
	if req.ResponseMaxSize == 0 {
		req.ResponseMaxSize = 64 // Default value
	}
	if req.ResponseMaxSize < 1 || req.ResponseMaxSize > 1024 {
		return nil, status.Error(codes.InvalidArgument, "response_max_size must be between 1 and 1,024")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var sessions []types.PermissionSession

	err := k.PermissionSession.Walk(sdkCtx, nil, func(key string, session types.PermissionSession) (bool, error) {
		// Apply modified_after filter if provided
		if req.ModifiedAfter != nil && !session.Modified.After(*req.ModifiedAfter) {
			return false, nil
		}

		sessions = append(sessions, session)
		return len(sessions) >= int(req.ResponseMaxSize), nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list sessions")
	}

	// Sort by modified time ascending
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Modified.Before(*sessions[j].Modified)
	})

	return &types.QueryListPermissionSessionsResponse{
		Sessions: sessions,
	}, nil
}

// IsAuthorizedIssuer implements the Is Authorized Issuer query
func (k Keeper) IsAuthorizedIssuer(ctx context.Context, req *types.QueryIsAuthorizedIssuerRequest) (*types.QueryIsAuthorizedIssuerResponse, error) {
	if err := validateIsAuthorizedIssuerRequest(req); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Load credential schema
	cs, err := k.credentialSchemaKeeper.GetCredentialSchemaById(sdkCtx, req.SchemaId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load credential schema")
	}

	// Check if issuance is open
	if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_OPEN {
		return &types.QueryIsAuthorizedIssuerResponse{
			Result: types.AuthorizationResult_AUTHORIZATION_RESULT_AUTHORIZED,
		}, nil
	}

	// Define time point
	timePoint := req.When
	if timePoint == nil {
		now := sdkCtx.BlockTime()
		timePoint = &now
	}

	// Find matching issuer permission
	var issuerPerm *types.Permission
	err = k.Permission.Walk(sdkCtx, nil, func(id uint64, perm types.Permission) (bool, error) {
		if perm.Did == req.IssuerDid &&
			perm.Type == types.PermissionType_PERMISSION_TYPE_ISSUER &&
			perm.SchemaId == req.SchemaId {

			// Check country match
			if req.Country != "" {
				if perm.Country != "" && perm.Country != req.Country {
					return false, nil
				}
			}

			// Check time validity
			if perm.EffectiveFrom != nil && timePoint.Before(*perm.EffectiveFrom) {
				return false, nil
			}
			if perm.EffectiveUntil != nil && timePoint.After(*perm.EffectiveUntil) {
				return false, nil
			}
			if perm.Revoked != nil && timePoint.After(*perm.Revoked) {
				return false, nil
			}
			if perm.Terminated != nil && timePoint.After(*perm.Terminated) {
				return false, nil
			}

			issuerPerm = &perm
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to search permissions")
	}

	if issuerPerm == nil {
		return &types.QueryIsAuthorizedIssuerResponse{
			Result: types.AuthorizationResult_AUTHORIZATION_RESULT_FORBIDDEN,
		}, nil
	}

	// Calculate permission set and fees
	permSet, err := k.buildPermissionSet(sdkCtx, issuerPerm)
	if err != nil {
		return nil, err
	}

	trustFees := uint64(0)
	for _, perm := range permSet {
		trustFees += perm.IssuanceFees
	}

	// If no fees required, return authorized
	if trustFees == 0 {
		return &types.QueryIsAuthorizedIssuerResponse{
			Result: types.AuthorizationResult_AUTHORIZATION_RESULT_AUTHORIZED,
		}, nil
	}

	// Session check
	//if req.SessionId == 0 {
	//	return &types.QueryIsAuthorizedIssuerResponse{
	//		Result: types.AuthorizationResult_AUTHORIZATION_RESULT_SESSION_REQUIRED,
	//	}, nil
	//}

	// Load session
	session, err := k.PermissionSession.Get(sdkCtx, req.SessionId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load session")
	}

	// Check session authorization exactly as specified
	//TODO: after specs update
	//if session.AgentDid != req.AgentDid {
	//	return &types.QueryIsAuthorizedIssuerResponse{
	//		Result: types.AuthorizationResult_AUTHORIZATION_RESULT_SESSION_REQUIRED,
	//	}, nil
	//}

	// Check if authz contains (issuer_perm.id, null, wallet_agent_did)
	for _, authz := range session.Authz {
		if authz.ExecutorPermId == issuerPerm.Id &&
			authz.BeneficiaryPermId == 0 &&
			authz.WalletAgentPermId == 0 {
			return &types.QueryIsAuthorizedIssuerResponse{
				Result: types.AuthorizationResult_AUTHORIZATION_RESULT_AUTHORIZED,
			}, nil
		}
	}

	return &types.QueryIsAuthorizedIssuerResponse{
		Result: types.AuthorizationResult_AUTHORIZATION_RESULT_SESSION_REQUIRED,
	}, nil
}

func (k Keeper) buildPermissionSet(ctx sdk.Context, perm *types.Permission) (PermissionSet, error) {
	permSet := make(PermissionSet, 0)
	currentPerm := perm

	// Process ancestors of executor perm
	for currentPerm.ValidatorPermId != 0 {
		validatorPerm, err := k.Permission.Get(ctx, currentPerm.ValidatorPermId)
		if err != nil {
			return nil, errors.Wrapf(sdkerrors.ErrNotFound, "validator permission not found: %d", currentPerm.ValidatorPermId)
		}

		// Add to set if not revoked and not terminated
		if validatorPerm.Revoked == nil && validatorPerm.Terminated == nil {
			permSet.add(validatorPerm)
		}

		currentPerm = &validatorPerm
	}

	return permSet, nil
}

func validateIsAuthorizedIssuerRequest(req *types.QueryIsAuthorizedIssuerRequest) error {
	if req == nil {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "empty request")
	}

	if !isValidDID(req.IssuerDid) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid issuer DID")
	}

	if !isValidDID(req.AgentDid) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid agent DID")
	}

	if !isValidDID(req.WalletAgentDid) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid wallet agent DID")
	}

	if req.SchemaId == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "schema ID required")
	}

	if req.Country != "" && !isValidCountryCode(req.Country) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid country code")
	}

	return nil
}

func isValidDID(did string) bool {
	if did == "" {
		return false
	}
	match, _ := regexp.MatchString(`^did:[a-zA-Z0-9]+:[a-zA-Z0-9._-]+$`, did)
	return match
}

func isValidCountryCode(code string) bool {
	if code == "" {
		return false
	}
	match, _ := regexp.MatchString(`^[A-Z]{2}$`, code)
	return match
}

// add adds a permission to the set if it doesn't already exist
func (ps *PermissionSet) add(perm types.Permission) {
	if !ps.contains(perm.Id) {
		*ps = append(*ps, perm)
	}
}

func (k Keeper) IsAuthorizedVerifier(ctx context.Context, req *types.QueryIsAuthorizedVerifierRequest) (*types.QueryIsAuthorizedVerifierResponse, error) {
	if err := validateIsAuthorizedVerifierRequest(req); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	timePoint := req.When
	if timePoint == nil {
		now := sdkCtx.BlockTime()
		timePoint = &now
	}

	// Load credential schema
	cs, err := k.credentialSchemaKeeper.GetCredentialSchemaById(sdkCtx, req.SchemaId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load credential schema")
	}

	// Check if verification is open
	if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_OPEN {
		return &types.QueryIsAuthorizedVerifierResponse{
			Result: types.AuthorizationResult_AUTHORIZATION_RESULT_AUTHORIZED,
		}, nil
	}

	// Find matching verifier permission
	var verifierPerm *types.Permission
	err = k.Permission.Walk(sdkCtx, nil, func(id uint64, perm types.Permission) (bool, error) {
		if perm.Did == req.VerifierDid &&
			perm.Type == types.PermissionType_PERMISSION_TYPE_VERIFIER &&
			perm.SchemaId == req.SchemaId {

			// Check country match
			if req.Country != "" {
				if perm.Country != "" && perm.Country != req.Country {
					return false, nil
				}
			}

			// Check time validity
			if !isPermissionValidAtTime(perm, timePoint) {
				return false, nil
			}

			verifierPerm = &perm
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to search verifier permissions")
	}

	if verifierPerm == nil {
		return &types.QueryIsAuthorizedVerifierResponse{
			Result: types.AuthorizationResult_AUTHORIZATION_RESULT_FORBIDDEN,
		}, nil
	}

	// Find matching issuer permission
	var issuerPerm *types.Permission
	err = k.Permission.Walk(sdkCtx, nil, func(id uint64, perm types.Permission) (bool, error) {
		if perm.Did == req.IssuerDid &&
			perm.Type == types.PermissionType_PERMISSION_TYPE_ISSUER &&
			perm.SchemaId == req.SchemaId {

			// Check country match
			if req.Country != "" {
				if perm.Country != "" && perm.Country != req.Country {
					return false, nil
				}
			}

			// Check time validity
			if !isPermissionValidAtTime(perm, timePoint) {
				return false, nil
			}

			issuerPerm = &perm
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to search issuer permissions")
	}

	if issuerPerm == nil {
		return &types.QueryIsAuthorizedVerifierResponse{
			Result: types.AuthorizationResult_AUTHORIZATION_RESULT_FORBIDDEN,
		}, nil
	}

	// Calculate permission set and fees
	permSet, err := k.buildPermissionSet2(sdkCtx, verifierPerm, issuerPerm)
	if err != nil {
		return nil, err
	}

	// Calculate total fees
	trustFees := uint64(0)
	for _, perm := range permSet {
		trustFees += perm.VerificationFees
	}

	// If no fees required, return authorized
	if trustFees == 0 {
		return &types.QueryIsAuthorizedVerifierResponse{
			Result: types.AuthorizationResult_AUTHORIZATION_RESULT_AUTHORIZED,
		}, nil
	}

	// Session required if fees exist but no session provided
	//if req.SessionId == 0 {
	//	return &types.QueryIsAuthorizedVerifierResponse{
	//		Result: types.AuthorizationResult_AUTHORIZATION_RESULT_SESSION_REQUIRED,
	//	}, nil
	//}

	// Verify session
	session, err := k.PermissionSession.Get(sdkCtx, req.SessionId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}

	for _, authz := range session.Authz {
		if authz.ExecutorPermId == verifierPerm.Id &&
			authz.BeneficiaryPermId == issuerPerm.Id {
			return &types.QueryIsAuthorizedVerifierResponse{
				Result: types.AuthorizationResult_AUTHORIZATION_RESULT_AUTHORIZED,
			}, nil
		}
	}

	return &types.QueryIsAuthorizedVerifierResponse{
		Result: types.AuthorizationResult_AUTHORIZATION_RESULT_SESSION_REQUIRED,
	}, nil
}

func (k Keeper) buildPermissionSet2(ctx sdk.Context, verifierPerm, issuerPerm *types.Permission) (PermissionSet, error) {
	permSet := make(PermissionSet, 0)

	// Process verifier ancestors
	currentPerm := verifierPerm
	for currentPerm.ValidatorPermId != 0 {
		validatorPerm, err := k.Permission.Get(ctx, currentPerm.ValidatorPermId)
		if err != nil {
			return nil, errors.Wrapf(sdkerrors.ErrNotFound, "validator permission not found: %d", currentPerm.ValidatorPermId)
		}

		if validatorPerm.Revoked == nil && validatorPerm.Terminated == nil {
			permSet.add(validatorPerm)
		}

		currentPerm = &validatorPerm
	}

	// Process issuer ancestors
	currentPerm = issuerPerm
	for currentPerm.ValidatorPermId != 0 {
		validatorPerm, err := k.Permission.Get(ctx, currentPerm.ValidatorPermId)
		if err != nil {
			return nil, errors.Wrapf(sdkerrors.ErrNotFound, "validator permission not found: %d", currentPerm.ValidatorPermId)
		}

		if validatorPerm.Revoked == nil && validatorPerm.Terminated == nil {
			permSet.add(validatorPerm)
		}

		currentPerm = &validatorPerm
	}

	return permSet, nil
}

func validateIsAuthorizedVerifierRequest(req *types.QueryIsAuthorizedVerifierRequest) error {
	if req == nil {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "empty request")
	}

	if !isValidDID(req.VerifierDid) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid verifier DID")
	}

	if !isValidDID(req.IssuerDid) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid issuer DID")
	}

	if !isValidDID(req.AgentDid) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid agent DID")
	}

	if !isValidDID(req.WalletAgentDid) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid wallet agent DID")
	}

	if req.SchemaId == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "schema ID required")
	}

	if req.Country != "" && !isValidCountryCode(req.Country) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid country code")
	}

	return nil
}

func isPermissionValidAtTime(perm types.Permission, timePoint *time.Time) bool {
	if perm.EffectiveFrom != nil && timePoint.Before(*perm.EffectiveFrom) {
		return false
	}
	if perm.EffectiveUntil != nil && timePoint.After(*perm.EffectiveUntil) {
		return false
	}
	if perm.Revoked != nil && timePoint.After(*perm.Revoked) {
		return false
	}
	if perm.Terminated != nil && timePoint.After(*perm.Terminated) {
		return false
	}
	return true
}
