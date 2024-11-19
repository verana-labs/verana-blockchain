package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	credentialschematypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
	trustregistrytypes "github.com/verana-labs/verana-blockchain/x/trustregistry/types"
)

func (ms msgServer) validatePermissions(ctx sdk.Context, msg *types.MsgCreateCredentialSchemaPerm, cs credentialschematypes.CredentialSchema, tr trustregistrytypes.TrustRegistry) error {
	permType := types.CredentialSchemaPermType(msg.CspType)

	switch permType {
	case types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY:
		return ms.validateTrustRegistryPerm(msg, tr)
	case types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER,
		types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER_GRANTOR:
		return ms.validateIssuerPerm(ctx, msg, cs)
	case types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER,
		types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER_GRANTOR:
		return ms.validateVerifierPerm(ctx, msg, cs)
	case types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_HOLDER:
		return ms.validateHolderPerm(ctx, msg, cs)
	default:
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid permission type")
	}
}

func (ms msgServer) validateTrustRegistryPerm(msg *types.MsgCreateCredentialSchemaPerm, tr trustregistrytypes.TrustRegistry) error {
	if msg.Creator != tr.Controller {
		return errors.Wrap(sdkerrors.ErrUnauthorized, "only trust registry controller can create TRUST_REGISTRY permissions")
	}
	if msg.Did != tr.Did {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "DID must match trust registry DID")
	}
	return nil
}

func (ms msgServer) validateIssuerPerm(ctx sdk.Context, msg *types.MsgCreateCredentialSchemaPerm, cs credentialschematypes.CredentialSchema) error {
	if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_OPEN {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "cannot create ISSUER permission when management mode is OPEN")
	}

	if msg.ValidationId == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "validation_id is required for ISSUER permissions")
	}

	return ms.validateValidationPermission(ctx, msg, cs, types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER)
}

func (ms msgServer) validateVerifierPerm(ctx sdk.Context, msg *types.MsgCreateCredentialSchemaPerm, cs credentialschematypes.CredentialSchema) error {
	if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_OPEN {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "cannot create VERIFIER permission when management mode is OPEN")
	}

	if msg.ValidationId == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "validation_id is required for VERIFIER permissions")
	}

	return ms.validateValidationPermission(ctx, msg, cs, types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER)
}

func (ms msgServer) validateHolderPerm(ctx sdk.Context, msg *types.MsgCreateCredentialSchemaPerm, cs credentialschematypes.CredentialSchema) error {
	if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_OPEN {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "cannot create HOLDER permission when management mode is OPEN")
	}

	if msg.ValidationId == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "validation_id is required for HOLDER permissions")
	}

	return ms.validateValidationPermission(ctx, msg, cs, types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_HOLDER)
}

func (ms msgServer) validateValidationPermission(ctx sdk.Context, msg *types.MsgCreateCredentialSchemaPerm, cs credentialschematypes.CredentialSchema, validatorKind types.CredentialSchemaPermType) error {
	// TODO: After Validation Module
	return nil
}

func (ms msgServer) checkOverlappingPermissions(ctx sdk.Context, msg *types.MsgCreateCredentialSchemaPerm) error {
	var overlappingPerms []types.CredentialSchemaPerm

	err := ms.CredentialSchemaPerm.Walk(ctx, nil, func(key uint64, perm types.CredentialSchemaPerm) (bool, error) {
		if isOverlapping(perm, msg) {
			overlappingPerms = append(overlappingPerms, perm)
		}
		return false, nil
	})

	if err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "error checking existing permissions")
	}

	for _, p := range overlappingPerms {
		if p.EffectiveUntil == nil {
			return errors.Wrap(sdkerrors.ErrInvalidRequest, "overlapping permission exists with no expiration")
		}
		if p.EffectiveUntil.After(msg.EffectiveFrom) {
			return errors.Wrap(sdkerrors.ErrInvalidRequest, "overlapping permission exists with later expiration")
		}
		if p.EffectiveFrom.Before(*msg.EffectiveUntil) {
			return errors.Wrap(sdkerrors.ErrInvalidRequest, "overlapping permission exists with earlier start")
		}
	}

	return nil
}

// Helper function to check if a permission overlaps with a new request
func isOverlapping(perm types.CredentialSchemaPerm, msg *types.MsgCreateCredentialSchemaPerm) bool {
	permType := types.CredentialSchemaPermType(msg.CspType)
	return perm.SchemaId == msg.SchemaId &&
		perm.CspType == permType &&
		perm.Country == msg.Country &&
		perm.Grantee == msg.Grantee &&
		perm.Revoked == nil &&
		perm.Terminated == nil
}

func (ms msgServer) createPermission(ctx sdk.Context, msg *types.MsgCreateCredentialSchemaPerm) error {
	id, err := ms.GetNextID(ctx, "csp")
	if err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "failed to generate ID")
	}

	permType := types.CredentialSchemaPermType(msg.CspType)

	perm := types.CredentialSchemaPerm{
		Id:               id,
		SchemaId:         msg.SchemaId,
		CspType:          permType,
		Did:              msg.Did,
		Grantee:          msg.Grantee,
		Created:          ctx.BlockTime(),
		CreatedBy:        msg.Creator,
		EffectiveFrom:    msg.EffectiveFrom,
		EffectiveUntil:   msg.EffectiveUntil,
		ValidationId:     msg.ValidationId,
		ValidationFees:   msg.ValidationFees,
		IssuanceFees:     msg.IssuanceFees,
		VerificationFees: msg.VerificationFees,
		Deposit:          0,
		Country:          msg.Country,
	}

	return ms.CredentialSchemaPerm.Set(ctx, id, perm)
}
