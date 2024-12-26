package keeper

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cstypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	csptypes "github.com/verana-labs/verana-blockchain/x/cspermission/types"
	"github.com/verana-labs/verana-blockchain/x/validation/types"
)

// validatePermissions implements [MOD-V-MSG-1-2-2] permission checks
func (ms msgServer) validatePermissions(ctx sdk.Context, msg *types.MsgCreateValidation) error {
	perm, err := ms.csPermissionKeeper.GetCSPermission(ctx, msg.ValidatorPermId)
	if err != nil {
		return fmt.Errorf("validator permission not found: %w", err)
	}

	// Check country compatibility
	if perm.Country != "" && perm.Country != msg.Country {
		return errors.New("validator does not serve the specified country")
	}

	cs, err := ms.credentialSchemaKeeper.GetCredentialSchemaById(ctx, perm.SchemaId)
	if err != nil {
		return fmt.Errorf("failed to get credential schema: %w", err)
	}

	switch msg.ValidationType {
	case uint32(types.ValidationType_ISSUER):
		// For debugging
		fmt.Printf("IssuerPermManagementMode Type: %T\n", cs.IssuerPermManagementMode)
		fmt.Printf("PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION Type: %T\n",
			cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION)

		if cs.IssuerPermManagementMode == cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION {
			if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER_GRANTOR {
				return errors.New("invalid validator permission type for issuer validation")
			}
		} else if cs.IssuerPermManagementMode == cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_TRUST_REGISTRY_VALIDATION {
			if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY {
				return errors.New("invalid validator permission type for issuer validation")
			}
		} else {
			return errors.New("invalid issuer permission management mode")
		}

	case uint32(types.ValidationType_ISSUER_GRANTOR):
		if cs.IssuerPermManagementMode == cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION {
			if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY {
				return errors.New("invalid validator permission type for issuer grantor validation")
			}
		} else {
			return errors.New("invalid issuer permission management mode for grantor")
		}

	case uint32(types.ValidationType_VERIFIER):
		if cs.VerifierPermManagementMode == cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION {
			if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER_GRANTOR {
				return errors.New("invalid validator permission type for verifier validation")
			}
		} else if cs.VerifierPermManagementMode == cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_TRUST_REGISTRY_VALIDATION {
			if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY {
				return errors.New("invalid validator permission type for verifier validation")
			}
		} else {
			return errors.New("invalid verifier permission management mode")
		}

	case uint32(types.ValidationType_VERIFIER_GRANTOR):
		if cs.VerifierPermManagementMode == cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION {
			if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY {
				return errors.New("invalid validator permission type for verifier grantor validation")
			}
		} else {
			return errors.New("invalid verifier permission management mode for grantor")
		}

	case uint32(types.ValidationType_HOLDER):
		if cs.VerifierPermManagementMode == cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION ||
			cs.VerifierPermManagementMode == cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_TRUST_REGISTRY_VALIDATION {
			if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER {
				return errors.New("invalid validator permission type for holder validation")
			}
		} else {
			return errors.New("invalid verifier permission management mode for holder")
		}

	default:
		return errors.New("invalid validation type")
	}

	return nil
}

// checkAndCalculateFees implements [MOD-V-MSG-1-2-3] fee checks
func (ms msgServer) checkAndCalculateFees(ctx sdk.Context, msg *types.MsgCreateValidation) (uint64, uint64, error) {
	perm, err := ms.csPermissionKeeper.GetCSPermission(ctx, msg.ValidatorPermId)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get validator permission: %w", err)
	}

	//TODO: update calculations after TD Module
	// Calculate fees and deposit
	validationFees := perm.ValidationFees
	validationDeposit := validationFees

	//TODO: send validation fees to validation module account

	return validationFees, validationDeposit, nil
}

// executeCreateValidation implements [MOD-V-MSG-1-3] execution
func (ms msgServer) executeCreateValidation(ctx sdk.Context, msg *types.MsgCreateValidation, fees, deposit uint64) (*types.Validation, error) {
	// Generate new validation ID
	id, err := ms.Keeper.GetNextID(ctx)
	if err != nil {
		return nil, err
	}

	now := ctx.BlockTime()

	// Create validation entry
	validation := &types.Validation{
		Id:                id,
		Applicant:         msg.Creator,
		Type:              types.ValidationType(msg.ValidationType),
		Created:           now,
		ValidatorPermId:   msg.ValidatorPermId,
		State:             types.ValidationState_PENDING,
		LastStateChange:   now,
		ApplicantDeposit:  deposit,
		ValidatorDeposits: []types.ValidatorDeposit{},
		CurrentFees:       fees,
		CurrentDeposit:    deposit,
	}

	// TODO: Handle deposits and fees IncreaseTrustDeposit

	// Save validation
	if err := ms.Validation.Set(ctx, id, *validation); err != nil {
		return nil, fmt.Errorf("failed to save validation: %w", err)
	}

	return validation, nil
}
