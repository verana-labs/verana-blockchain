package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"
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
