package keeper

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cstypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	csptypes "github.com/verana-labs/verana-blockchain/x/cspermission/types"
	"github.com/verana-labs/verana-blockchain/x/validation/types"
)

// validateRenewalBasics implements [MOD-V-MSG-2-2-1] basic checks
func (ms msgServer) validateRenewalBasics(ctx sdk.Context, msg *types.MsgRenewValidation) (*types.Validation, uint64, error) {
	if msg.Id == 0 {
		return nil, 0, errors.New("validation id is required")
	}

	val, err := ms.Validation.Get(ctx, msg.Id)
	if err != nil {
		return nil, 0, fmt.Errorf("validation not found: %w", err)
	}

	if val.Applicant != msg.Creator {
		return nil, 0, errors.New("only the validation applicant can renew it")
	}

	validatorPermID := msg.ValidatorPermId
	if validatorPermID == 0 {
		validatorPermID = val.ValidatorPermId
	}

	return &val, validatorPermID, nil
}

// validateRenewalPermissions implements [MOD-V-MSG-2-2-2] permission checks
func (ms msgServer) validateRenewalPermissions(ctx sdk.Context, msg *types.MsgRenewValidation, val *types.Validation, validatorPermID uint64) error {
	// Check if validator is being changed
	if validatorPermID != val.ValidatorPermId {
		// Get new validator permission
		perm, err := ms.csPermissionKeeper.GetCSPermission(ctx, validatorPermID)
		if err != nil {
			return fmt.Errorf("validator permission not found: %w", err)
		}

		// Get old validator permission for event emission and country check
		oldPerm, err := ms.csPermissionKeeper.GetCSPermission(ctx, val.ValidatorPermId)
		if err != nil {
			return fmt.Errorf("failed to get old validator permission: %w", err)
		}

		// Check country compatibility
		if oldPerm.Country != "" && perm.Country != oldPerm.Country {
			return errors.New("validator does not serve the specified country")
		}

		// Get credential schema for validation
		cs, err := ms.credentialSchemaKeeper.GetCredentialSchemaById(ctx, perm.SchemaId)
		if err != nil {
			return fmt.Errorf("failed to get credential schema: %w", err)
		}

		// Validate based on validation type and schema management mode
		switch val.Type {
		case types.ValidationType_ISSUER:
			if cs.IssuerPermManagementMode == cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION {
				if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER_GRANTOR {
					return errors.New("invalid validator permission type for issuer validation")
				}
			} else if cs.IssuerPermManagementMode == cstypes.CredentialSchemaPermManagementMode_TRUST_REGISTRY_VALIDATION {
				if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY {
					return errors.New("invalid validator permission type for issuer validation")
				}
			} else {
				return errors.New("invalid issuer permission management mode")
			}

		case types.ValidationType_ISSUER_GRANTOR:
			if cs.IssuerPermManagementMode == cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION {
				if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY {
					return errors.New("invalid validator permission type for issuer grantor validation")
				}
			} else {
				return errors.New("invalid issuer permission management mode for grantor")
			}

		case types.ValidationType_VERIFIER:
			if cs.VerifierPermManagementMode == cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION {
				if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER_GRANTOR {
					return errors.New("invalid validator permission type for verifier validation")
				}
			} else if cs.VerifierPermManagementMode == cstypes.CredentialSchemaPermManagementMode_TRUST_REGISTRY_VALIDATION {
				if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY {
					return errors.New("invalid validator permission type for verifier validation")
				}
			} else {
				return errors.New("invalid verifier permission management mode")
			}

		case types.ValidationType_VERIFIER_GRANTOR:
			if cs.VerifierPermManagementMode == cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION {
				if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY {
					return errors.New("invalid validator permission type for verifier grantor validation")
				}
			} else {
				return errors.New("invalid verifier permission management mode for grantor")
			}

		case types.ValidationType_HOLDER:
			if cs.VerifierPermManagementMode == cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION ||
				cs.VerifierPermManagementMode == cstypes.CredentialSchemaPermManagementMode_TRUST_REGISTRY_VALIDATION {
				if perm.CspType != csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER {
					return errors.New("invalid validator permission type for holder validation")
				}
			} else {
				return errors.New("invalid verifier permission management mode for holder")
			}
		}

		// Emit revocation control transfer event
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				"validation_revocation_control_transfer",
				sdk.NewAttribute("validation_id", fmt.Sprintf("%d", msg.Id)),
				sdk.NewAttribute("old_validator", oldPerm.Grantee),
				sdk.NewAttribute("new_validator", perm.Grantee),
			),
		)
	} else {
		// Verify existing validator permission is still valid
		_, err := ms.csPermissionKeeper.GetCSPermission(ctx, validatorPermID)
		if err != nil {
			return fmt.Errorf("validator permission not found: %w", err)
		}
	}

	return nil
}

func (ms msgServer) checkAndCalculateRenewalFees(ctx sdk.Context, validatorPermID uint64) (uint64, uint64, error) {
	perm, err := ms.csPermissionKeeper.GetCSPermission(ctx, validatorPermID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get validator permission: %w", err)
	}

	validationFees := perm.ValidationFees
	validationDeposit := validationFees

	return validationFees, validationDeposit, nil
}

func (ms msgServer) executeRenewalValidation(ctx sdk.Context, val *types.Validation, validatorPermID uint64, fees, deposit uint64) error {
	now := sdk.UnwrapSDKContext(ctx).BlockTime()

	// First update current values
	val.CurrentFees = fees
	val.CurrentDeposit = deposit

	// Then calculate and update total deposit
	val.ApplicantDeposit = val.ApplicantDeposit + val.CurrentDeposit

	// Update other state
	val.State = types.ValidationState_PENDING
	val.LastStateChange = now
	val.ValidatorPermId = validatorPermID

	// Save validation
	if err := ms.Validation.Set(ctx, val.Id, *val); err != nil {
		return fmt.Errorf("failed to save validation: %w", err)
	}

	return nil
}
