package keeper

import (
	"context"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	credentialschematypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
	"time"
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

func (k Keeper) GetCSP(ctx context.Context, req *types.QueryGetCSPRequest) (*types.QueryGetCSPResponse, error) {
	if req.Id == 0 {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "id must be provided")
	}

	perm, err := k.CredentialSchemaPerm.Get(ctx, req.Id)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrKeyNotFound, "credential schema permission not found: %d", req.Id)
	}

	return &types.QueryGetCSPResponse{
		Permission: perm,
	}, nil
}

func (k Keeper) IsAuthorizedIssuer(ctx context.Context, req *types.QueryIsAuthorizedIssuerRequest) (*types.QueryIsAuthorizedIssuerResponse, error) {
	// Validate request
	if err := validateIsAuthorizedIssuerRequest(req); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Get the credential schema
	cs, err := k.credentialSchemaKeeper.GetCredentialSchemaById(sdkCtx, req.SchemaId)
	if err != nil {
		return nil, errors.Wrapf(sdkerrors.ErrNotFound, "credential schema not found: %d", req.SchemaId)
	}

	// If issuer mode is OPEN, return authorized immediately
	if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_OPEN {
		return &types.QueryIsAuthorizedIssuerResponse{
			Status: types.AuthorizationStatus_AUTHORIZED,
		}, nil
	}

	// Set check time
	checkTime := time.Now()
	if req.When != nil {
		checkTime = *req.When
	}

	// Find matching issuer permission
	var issuerPerm *types.CredentialSchemaPerm
	err = k.CredentialSchemaPerm.Walk(ctx, nil, func(key uint64, perm types.CredentialSchemaPerm) (bool, error) {
		if perm.Did == req.IssuerDid &&
			perm.SchemaId == req.SchemaId &&
			perm.CspType == types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER &&
			(req.Country == "" || perm.Country == "" || perm.Country == req.Country) {

			// Check time validity
			if isPermissionValidAtTime(&perm, checkTime) {
				issuerPerm = &perm
				return true, nil
			}
		}
		return false, nil
	})

	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "error checking permissions")
	}

	// If no valid permission found, return forbidden
	if issuerPerm == nil {
		return &types.QueryIsAuthorizedIssuerResponse{
			Status: types.AuthorizationStatus_FORBIDDEN,
		}, nil
	}

	// TODO: Implement fee calculation after validation module is ready
	// This would use MOD-CSPS-MSG-1-2-2 and MOD-CSPS-MSG-1-2-3
	trustFees := uint64(0) // Placeholder until validation module is ready

	// If no fees required, return authorized
	if trustFees == 0 {
		return &types.QueryIsAuthorizedIssuerResponse{
			Status: types.AuthorizationStatus_AUTHORIZED,
		}, nil
	}

	// If fees required but no session provided, require session
	if req.SessionId == 0 {
		return &types.QueryIsAuthorizedIssuerResponse{
			Status: types.AuthorizationStatus_SESSION_REQUIRED,
		}, nil
	}

	// TODO: Check session after validation module is ready
	// This would verify session.user_agent_did and session.session_authz[]

	return &types.QueryIsAuthorizedIssuerResponse{
		Status: types.AuthorizationStatus_SESSION_REQUIRED,
	}, nil
}

func validateIsAuthorizedIssuerRequest(req *types.QueryIsAuthorizedIssuerRequest) error {
	if req.IssuerDid == "" || !types.IsValidDID(req.IssuerDid) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid issuer DID")
	}

	if req.UserAgentDid == "" || !types.IsValidDID(req.UserAgentDid) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid user agent DID")
	}

	if req.WalletUserAgentDid == "" || !types.IsValidDID(req.WalletUserAgentDid) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid wallet user agent DID")
	}

	if req.SchemaId == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "schema ID is required")
	}

	if req.Country != "" && !isValidCountryCode(req.Country) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid country code")
	}

	return nil
}

func isPermissionValidAtTime(perm *types.CredentialSchemaPerm, checkTime time.Time) bool {
	// Check if time is within effective period
	if checkTime.Before(perm.EffectiveFrom) {
		return false
	}
	if perm.EffectiveUntil != nil && checkTime.After(*perm.EffectiveUntil) {
		return false
	}

	// Check if not revoked at check time
	if perm.Revoked != nil && !checkTime.After(*perm.Revoked) {
		return false
	}

	// Check if not terminated at check time
	if perm.Terminated != nil && !checkTime.After(*perm.Terminated) {
		return false
	}

	return true
}

func isValidCountryCode(code string) bool {
	// Basic check for ISO 3166-1 alpha-2 code
	if len(code) != 2 {
		return false
	}
	for _, c := range code {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}
