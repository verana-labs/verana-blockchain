package keeper

import (
	"cosmossdk.io/math"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	credentialschematypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
)

func (ms msgServer) validatePermissionChecks(ctx sdk.Context, msg *types.MsgStartPermissionVP) (types.Permission, error) {
	// Load validator permission
	validatorPerm, err := ms.Keeper.GetPermissionByID(ctx, msg.ValidatorPermId)
	if err != nil {
		return types.Permission{}, fmt.Errorf("validator permission not found: %w", err)
	}

	// Check if validator permission is valid
	if err := IsValidPermission(validatorPerm, msg.Country, ctx.BlockTime()); err != nil {
		return types.Permission{}, fmt.Errorf("validator permission is not valid: %w", err)
	}

	// Check country compatibility
	if validatorPerm.Country != "" && validatorPerm.Country != msg.Country {
		return types.Permission{}, fmt.Errorf("validator permission country mismatch")
	}

	// Load credential schema
	cs, err := ms.credentialSchemaKeeper.GetCredentialSchemaById(ctx, validatorPerm.SchemaId)
	if err != nil {
		return types.Permission{}, fmt.Errorf("credential schema not found: %w", err)
	}

	// Validate permission type combinations
	if err := validatePermissionTypeCombination(types.PermissionType(msg.Type), validatorPerm.Type, cs); err != nil {
		return types.Permission{}, err
	}

	return validatorPerm, nil
}

func (ms msgServer) validateAndCalculateFees(ctx sdk.Context, creator string, validatorPerm types.Permission) (uint64, uint64, error) {
	// Get global variables
	trustUnitPrice := ms.trustRegistryKeeper.GetTrustUnitPrice(ctx)
	trustDepositRate := ms.trustDeposit.GetTrustDepositRate(ctx)

	validationFeesInDenom := validatorPerm.ValidationFees * trustUnitPrice
	validationTrustDepositInDenom := ms.Keeper.validationTrustDepositInDenomAmount(validationFeesInDenom, trustDepositRate)

	return validationFeesInDenom, validationTrustDepositInDenom, nil
}

func (k Keeper) validationTrustDepositInDenomAmount(validationFeesInDenom uint64, trustDepositRate math.LegacyDec) uint64 {
	validationFeesInDenomDec := math.LegacyNewDec(int64(validationFeesInDenom))
	validationTrustDepositInDenom := validationFeesInDenomDec.Mul(trustDepositRate)
	return validationTrustDepositInDenom.TruncateInt().Uint64()
}

func (ms msgServer) executeStartPermissionVP(ctx sdk.Context, msg *types.MsgStartPermissionVP, validatorPerm types.Permission, fees, deposit uint64) (uint64, error) {
	// Increment trust deposit if deposit is greater than 0
	if deposit > 0 {
		if err := ms.trustDeposit.AdjustTrustDeposit(ctx, msg.Creator, int64(deposit)); err != nil {
			return 0, fmt.Errorf("failed to increase trust deposit: %w", err)
		}
	}

	// Send validation fees to escrow account if greater than 0
	if fees > 0 {
		// Get sender address
		senderAddr, err := sdk.AccAddressFromBech32(msg.Creator)
		if err != nil {
			return 0, fmt.Errorf("invalid creator address: %w", err)
		}

		// Transfer fees to module escrow account
		err = ms.bankKeeper.SendCoinsFromAccountToModule(
			ctx,
			senderAddr,
			types.ModuleName, // Using module name as the escrow account
			sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, int64(fees))),
		)
		if err != nil {
			return 0, fmt.Errorf("failed to transfer validation fees to escrow: %w", err)
		}
	}

	// Create new permission entry
	now := ctx.BlockTime()
	newPerm := types.Permission{
		Type:              types.PermissionType(msg.Type),
		SchemaId:          validatorPerm.SchemaId,
		Did:               msg.Did,
		Grantee:           msg.Creator,
		Created:           &now,
		CreatedBy:         msg.Creator,
		Modified:          &now,
		ValidationFees:    0,
		IssuanceFees:      0,
		VerificationFees:  0,
		Deposit:           deposit,
		Country:           msg.Country,
		ValidatorPermId:   msg.ValidatorPermId,
		VpState:           types.ValidationState_VALIDATION_STATE_PENDING,
		VpLastStateChange: &now,
		VpCurrentFees:     fees,
		VpCurrentDeposit:  deposit,
	}

	// Store the permission
	id, err := ms.Keeper.CreatePermission(ctx, newPerm)
	if err != nil {
		return 0, fmt.Errorf("failed to create permission: %w", err)
	}

	return id, nil
}

func validatePermissionTypeCombination(requestedType, validatorType types.PermissionType, cs credentialschematypes.CredentialSchema) error {
	switch requestedType {
	case types.PermissionType_PERMISSION_TYPE_ISSUER:
		if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION {
			if validatorType != types.PermissionType_PERMISSION_TYPE_ISSUER_GRANTOR {
				return fmt.Errorf("issuer permission requires ISSUER_GRANTOR validator")
			}
		} else if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_ECOSYSTEM {
			if validatorType != types.PermissionType_PERMISSION_TYPE_ECOSYSTEM {
				return fmt.Errorf("issuer permission requires ECOSYSTEM validator")
			}
		} else if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_OPEN {
			// Mode is OPEN which means anyone can issue credential of this schema
			// But formal permission creation is still needed when payment is required
			// Check if validator has the correct type for fee collection
			if validatorType != types.PermissionType_PERMISSION_TYPE_ECOSYSTEM {
				return fmt.Errorf("open issuance still requires ECOSYSTEM validator for fee collection")
			}
		} else {
			return fmt.Errorf("issuer permission not supported with current schema settings")
		}

	case types.PermissionType_PERMISSION_TYPE_ISSUER_GRANTOR:
		if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION {
			if validatorType != types.PermissionType_PERMISSION_TYPE_ECOSYSTEM {
				return fmt.Errorf("issuer grantor permission requires ECOSYSTEM validator")
			}
		} else {
			return fmt.Errorf("issuer grantor permission not supported with current schema settings")
		}

	case types.PermissionType_PERMISSION_TYPE_VERIFIER:
		if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION {
			if validatorType != types.PermissionType_PERMISSION_TYPE_VERIFIER_GRANTOR {
				return fmt.Errorf("verifier permission requires VERIFIER_GRANTOR validator")
			}
		} else if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_ECOSYSTEM {
			if validatorType != types.PermissionType_PERMISSION_TYPE_ECOSYSTEM {
				return fmt.Errorf("verifier permission requires ECOSYSTEM validator")
			}
		} else if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_OPEN {
			// Mode is OPEN which means anyone can verify credentials of this schema
			// This doesn't imply no payment is necessary - formal permission might be
			// required when payment is needed
			// Check if validator has the correct type for fee collection
			if validatorType != types.PermissionType_PERMISSION_TYPE_ECOSYSTEM {
				return fmt.Errorf("open verification still requires ECOSYSTEM validator for fee collection")
			}
		} else {
			return fmt.Errorf("verifier permission not supported with current schema settings")
		}

	case types.PermissionType_PERMISSION_TYPE_VERIFIER_GRANTOR:
		if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION {
			if validatorType != types.PermissionType_PERMISSION_TYPE_ECOSYSTEM {
				return fmt.Errorf("verifier grantor permission requires ECOSYSTEM validator")
			}
		} else {
			return fmt.Errorf("verifier grantor permission not supported with current schema settings")
		}

	case types.PermissionType_PERMISSION_TYPE_HOLDER:
		if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION ||
			cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_ECOSYSTEM {
			if validatorType != types.PermissionType_PERMISSION_TYPE_ISSUER {
				return fmt.Errorf("holder permission requires ISSUER validator")
			}
		} else if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_OPEN {
			// Even in OPEN mode, holder permissions might require validation from an ISSUER
			if validatorType != types.PermissionType_PERMISSION_TYPE_ISSUER {
				return fmt.Errorf("holder permission requires ISSUER validator even in OPEN verification mode")
			}
		} else {
			return fmt.Errorf("holder permission not supported with current schema settings")
		}
	}

	return nil
}
