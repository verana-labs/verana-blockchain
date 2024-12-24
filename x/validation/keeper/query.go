package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/validation/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"
)

var _ types.QueryServer = queryServer{}

func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

func (qs queryServer) ListValidations(ctx context.Context, req *types.QueryListValidationsRequest) (*types.QueryListValidationsResponse, error) {
	if err := validateListValidationsRequest(req); err != nil {
		return nil, err
	}

	// Default response size if not specified
	if req.ResponseMaxSize == 0 {
		req.ResponseMaxSize = 64
	}

	var validations []types.Validation
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Collect matching validations
	err := qs.k.Validation.Walk(sdkCtx, nil, func(key uint64, val types.Validation) (bool, error) {
		// Apply filters
		if !matchesFilters(val, req) {
			return false, nil
		}

		validations = append(validations, val)
		return len(validations) >= int(req.ResponseMaxSize), nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Sort based on exp_before parameter
	if req.ExpBefore != nil {
		sort.Slice(validations, func(i, j int) bool {
			return validations[i].Exp.Before(*validations[j].Exp)
		})
	} else {
		sort.Slice(validations, func(i, j int) bool {
			return validations[i].LastStateChange.Before(validations[j].LastStateChange)
		})
	}

	return &types.QueryListValidationsResponse{
		Validations: validations,
	}, nil
}

func validateListValidationsRequest(req *types.QueryListValidationsRequest) error {
	// Check response_max_size bounds
	if req.ResponseMaxSize < 1 || req.ResponseMaxSize > 1024 {
		return status.Error(codes.InvalidArgument, "response_max_size must be between 1 and 1,024")
	}

	// Check that at least one of controller or validator_perm_id is specified
	if req.Controller == "" && req.ValidatorPermId == 0 {
		return status.Error(codes.InvalidArgument, "either controller or validator_perm_id must be specified")
	}

	return nil
}

func matchesFilters(val types.Validation, req *types.QueryListValidationsRequest) bool {
	// Check controller filter
	if req.Controller != "" && val.Applicant != req.Controller {
		return false
	}

	// Check validator permission ID filter
	if req.ValidatorPermId != 0 && val.ValidatorPermId != req.ValidatorPermId {
		return false
	}

	// Check type filter
	if req.Type != types.ValidationType_TYPE_UNSPECIFIED && val.Type != req.Type {
		return false
	}

	// Check state filter
	if req.State != types.ValidationState_STATE_UNSPECIFIED && val.State != req.State {
		return false
	}

	// Check expiration filter
	if req.ExpBefore != nil && (val.Exp.IsZero() || !val.Exp.Before(*req.ExpBefore)) {
		return false
	}

	return true
}

func (qs queryServer) GetValidation(ctx context.Context, req *types.QueryGetValidationRequest) (*types.QueryGetValidationResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "validation id is required")
	}

	validation, err := qs.k.Validation.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "validation not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetValidationResponse{
		Validation: validation,
	}, nil
}
