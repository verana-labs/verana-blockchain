package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ListCredentialSchemas(goCtx context.Context, req *types.QueryListCredentialSchemasRequest) (*types.QueryListCredentialSchemasResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate response_max_size
	maxSize := uint32(64) // default
	if req.ResponseMaxSize > 0 {
		if req.ResponseMaxSize > 1024 {
			return nil, status.Error(codes.InvalidArgument, "response_max_size must be between 1 and 1024")
		}
		maxSize = req.ResponseMaxSize
	}

	var schemas []types.CredentialSchema
	var err error

	err = k.CredentialSchema.Walk(ctx, nil, func(key uint64, schema types.CredentialSchema) (bool, error) {
		// Filter by trust registry if specified
		if req.TrId != 0 && schema.TrId != req.TrId {
			return false, nil
		}

		// Filter by created_after if specified
		if req.CreatedAfter != nil && schema.Created.Before(*req.CreatedAfter) {
			return false, nil
		}

		schemas = append(schemas, schema)

		// Stop if we've reached max size
		return len(schemas) >= int(maxSize), nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListCredentialSchemasResponse{
		Schemas: schemas,
	}, nil
}

func (k Keeper) GetCredentialSchema(goCtx context.Context, req *types.QueryGetCredentialSchemaRequest) (*types.QueryGetCredentialSchemaResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	schema, err := k.CredentialSchema.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "credential schema not found")
	}

	return &types.QueryGetCredentialSchemaResponse{
		Schema: schema,
	}, nil
}

func (k Keeper) RenderJsonSchema(goCtx context.Context, req *types.QueryRenderJsonSchemaRequest) (*types.QueryRenderJsonSchemaResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	schema, err := k.CredentialSchema.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "credential schema not found")
	}

	return &types.QueryRenderJsonSchemaResponse{
		Schema: schema.JsonSchema,
	}, nil
}
