package keeper

import (
	"context"
	"cosmossdk.io/collections"
	errors2 "errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"sort"
	"time"
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

func (k Keeper) FindPermissionsWithDID(goCtx context.Context, req *types.QueryFindPermissionsWithDIDRequest) (*types.QueryFindPermissionsWithDIDResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-PERM-QRY-3-2] Checks
	if req.Did == "" {
		return nil, status.Error(codes.InvalidArgument, "DID is required")
	}
	if !isValidDID(req.Did) {
		return nil, status.Error(codes.InvalidArgument, "invalid DID format")
	}

	// Check type - convert uint32 to PermissionType
	if req.Type == 0 {
		return nil, status.Error(codes.InvalidArgument, "permission type is required")
	}

	// Validate permission type value is in range
	permType := types.PermissionType(req.Type)
	if permType < types.PermissionType_PERMISSION_TYPE_ISSUER ||
		permType > types.PermissionType_PERMISSION_TYPE_HOLDER {
		return nil, status.Error(codes.InvalidArgument,
			fmt.Sprintf("invalid permission type value: %d, must be between 1 and 6", req.Type))
	}

	// Check schema ID
	if req.SchemaId == 0 {
		return nil, status.Error(codes.InvalidArgument, "schema ID is required")
	}

	// Check schema exists
	_, err := k.credentialSchemaKeeper.GetCredentialSchemaById(ctx, req.SchemaId)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("credential schema not found: %v", err))
	}

	// Check country code if provided
	if req.Country != "" && !isValidCountryCode(req.Country) {
		return nil, status.Error(codes.InvalidArgument, "invalid country code format")
	}

	// [MOD-PERM-QRY-3-3] Execution
	var foundPerms []types.Permission

	// TODO: If index is implemented, use it here to get permission IDs by schema and hash
	// For now, we'll scan all permissions

	err = k.Permission.Walk(ctx, nil, func(id uint64, perm types.Permission) (bool, error) {
		// Filter by schema ID
		if perm.SchemaId != req.SchemaId {
			return false, nil
		}

		// Filter by DID and type
		if perm.Did != req.Did || perm.Type != permType {
			return false, nil
		}

		// Filter by country
		if req.Country != "" && perm.Country != "" && perm.Country != req.Country {
			return false, nil
		}

		// If "when" is not specified, add all matching permissions
		if req.When == nil {
			foundPerms = append(foundPerms, perm)
			return false, nil
		}

		// Filter by time validity
		if isPermissionValidAtTime(perm, *req.When) {
			foundPerms = append(foundPerms, perm)
		}

		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to query permissions: %v", err))
	}

	return &types.QueryFindPermissionsWithDIDResponse{
		Permissions: foundPerms,
	}, nil
}

// Helper function to check if a permission is valid at a specific time
func isPermissionValidAtTime(perm types.Permission, when time.Time) bool {
	// Check effective_from
	if perm.EffectiveFrom != nil && when.Before(*perm.EffectiveFrom) {
		return false
	}

	// Check effective_until
	if perm.EffectiveUntil != nil && !when.Before(*perm.EffectiveUntil) {
		return false
	}

	// Check revoked
	if perm.Revoked != nil && !when.Before(*perm.Revoked) {
		return false
	}

	// Check terminated
	if perm.Terminated != nil && !when.Before(*perm.Terminated) {
		return false
	}

	return true
}

func isValidDID(did string) bool {
	// Basic DID validation regex
	// This is a simplified version and may need to be expanded based on specific DID method requirements
	didRegex := regexp.MustCompile(`^did:[a-zA-Z0-9]+:[a-zA-Z0-9._-]+$`)
	return didRegex.MatchString(did)
}

func isValidCountryCode(code string) bool {
	// Basic check for ISO 3166-1 alpha-2 format
	match, _ := regexp.MatchString(`^[A-Z]{2}$`, code)
	return match
}

func (k Keeper) FindBeneficiaries(goCtx context.Context, req *types.QueryFindBeneficiariesRequest) (*types.QueryFindBeneficiariesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// [MOD-PERM-QRY-4-2] Checks
	// At least one of issuer_perm_id or verifier_perm_id must be provided
	if req.IssuerPermId == 0 && req.VerifierPermId == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one of issuer_perm_id or verifier_perm_id must be provided")
	}

	var issuerPerm, verifierPerm *types.Permission

	// Load issuer permission if specified
	if req.IssuerPermId != 0 {
		perm, err := k.Permission.Get(ctx, req.IssuerPermId)
		if err != nil {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("issuer permission not found: %v", err))
		}
		issuerPerm = &perm
	}

	// Load verifier permission if specified
	if req.VerifierPermId != 0 {
		perm, err := k.Permission.Get(ctx, req.VerifierPermId)
		if err != nil {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("verifier permission not found: %v", err))
		}
		verifierPerm = &perm
	}

	// [MOD-PERM-QRY-4-3] Execution
	// Use a map to implement the set functionality
	foundPermMap := make(map[uint64]types.Permission)

	// Process issuer permission hierarchy
	if issuerPerm != nil {
		// Start with the validator of issuer_perm
		if issuerPerm.ValidatorPermId != 0 {
			currentPermID := issuerPerm.ValidatorPermId

			// Traverse up the validator chain
			for currentPermID != 0 {
				currentPerm, err := k.Permission.Get(ctx, currentPermID)
				if err != nil {
					return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get permission: %v", err))
				}

				// Add to set if not revoked or terminated
				if currentPerm.Revoked == nil && currentPerm.Terminated == nil {
					foundPermMap[currentPermID] = currentPerm
				}

				// Move up to the next validator
				currentPermID = currentPerm.ValidatorPermId
			}
		}
	}

	// Process verifier permission hierarchy
	if verifierPerm != nil {
		// First add issuer_perm to the set if it exists
		if issuerPerm != nil && issuerPerm.Revoked == nil && issuerPerm.Terminated == nil {
			foundPermMap[req.IssuerPermId] = *issuerPerm
		}

		// Start with the validator of verifier_perm
		if verifierPerm.ValidatorPermId != 0 {
			currentPermID := verifierPerm.ValidatorPermId

			// Traverse up the validator chain
			for currentPermID != 0 {
				currentPerm, err := k.Permission.Get(ctx, currentPermID)
				if err != nil {
					return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get permission: %v", err))
				}

				// Add to set if not revoked or terminated
				if currentPerm.Revoked == nil && currentPerm.Terminated == nil {
					foundPermMap[currentPermID] = currentPerm
				}

				// Move up to the next validator
				currentPermID = currentPerm.ValidatorPermId
			}
		}
	}

	// Convert map to array
	permissions := make([]types.Permission, 0, len(foundPermMap))
	for _, perm := range foundPermMap {
		permissions = append(permissions, perm)
	}

	return &types.QueryFindBeneficiariesResponse{
		Permissions: permissions,
	}, nil
}
