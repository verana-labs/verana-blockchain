package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	credentialschematypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
	trustregistrytypes "github.com/verana-labs/verana-blockchain/x/trustregistry/types"
)

func (ms msgServer) validateRevokePermissions(ctx sdk.Context, creator string, csp *types.CredentialSchemaPerm, cs credentialschematypes.CredentialSchema, tr trustregistrytypes.TrustRegistry) error {
	switch csp.CspType {
	case types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY:
		if creator != tr.Controller {
			return errors.Wrap(sdkerrors.ErrUnauthorized, "only trust registry controller can revoke TRUST_REGISTRY permissions")
		}

	case types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER,
		types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER_GRANTOR:
		return ms.validateIssuerRevoke(ctx, creator, csp, cs)

	case types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER,
		types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER_GRANTOR:
		return ms.validateVerifierRevoke(ctx, creator, csp, cs)

	case types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_HOLDER:
		return ms.validateHolderRevoke(ctx, creator, csp, cs)
	}

	return nil
}

func (ms msgServer) validateIssuerRevoke(ctx sdk.Context, creator string, csp *types.CredentialSchemaPerm, cs credentialschematypes.CredentialSchema) error {
	if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_OPEN {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "cannot revoke ISSUER permission when management mode is OPEN")
	}

	return ms.validateRevokeValidation(ctx, creator, csp, cs, types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER)
}

func (ms msgServer) validateVerifierRevoke(ctx sdk.Context, creator string, csp *types.CredentialSchemaPerm, cs credentialschematypes.CredentialSchema) error {
	if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_OPEN {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "cannot revoke VERIFIER permission when management mode is OPEN")
	}

	return ms.validateRevokeValidation(ctx, creator, csp, cs, types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER)
}

func (ms msgServer) validateHolderRevoke(ctx sdk.Context, creator string, csp *types.CredentialSchemaPerm, cs credentialschematypes.CredentialSchema) error {
	if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_OPEN {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "cannot revoke HOLDER permission when management mode is OPEN")
	}

	return ms.validateRevokeValidation(ctx, creator, csp, cs, types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_HOLDER)
}

func (ms msgServer) validateRevokeValidation(ctx sdk.Context, creator string, csp *types.CredentialSchemaPerm, cs credentialschematypes.CredentialSchema, validatorKind types.CredentialSchemaPermType) error {
	// TODO: After validation module
	return nil
}
