package keeper

import (
	"context"
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ListCSP(ctx context.Context, req *types.QueryListCSPRequest) (*types.QueryListCSPResponse, error) {
	if err := req.ValidateRequest(); err != nil {
		return nil, err
	}

	var perms []types.CredentialSchemaPerm
	err := k.CredentialSchemaPerm.Walk(ctx, nil, func(key uint64, perm types.CredentialSchemaPerm) (bool, error) {
		if perm.SchemaId != req.SchemaId {
			return false, nil
		}
		if req.Creator != "" && perm.CreatedBy != req.Creator {
			return false, nil
		}
		if req.Grantee != "" && perm.Grantee != req.Grantee {
			return false, nil
		}
		if req.Did != "" && perm.Did != req.Did {
			return false, nil
		}
		if req.Type != types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_UNSPECIFIED && perm.CspType != req.Type {
			return false, nil
		}

		perms = append(perms, perm)
		return len(perms) >= int(req.ResponseMaxSize), nil
	})

	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return &types.QueryListCSPResponse{
		Permissions: perms,
	}, nil
}
