package keeper

import (
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
	// TODO: After trust deposit module
	// Get global variables
	//trustUnitPrice := ms.Keeper.GetTrustUnitPrice(ctx)
	//trustDepositRate := ms.Keeper.GetTrustDepositRate(ctx)

	// Calculate validation fees
	//validationFeesInDenom := validatorPerm.ValidationFees * trustUnitPrice
	//validationTrustDepositInDenom := validationFeesInDenom * trustDepositRate

	// Check if account has sufficient balance
	//creatorAddr, err := sdk.AccAddressFromBech32(creator)
	//if err != nil {
	//	return 0, 0, fmt.Errorf("invalid creator address: %w", err)
	//}
	//
	//balance := ms.bankKeeper.GetBalance(ctx, creatorAddr, ms.Keeper.GetParams(ctx).FeeDenom)
	//requiredAmount := validationFeesInDenom + validationTrustDepositInDenom
	//
	//if balance.Amount.LT(sdk.NewInt(int64(requiredAmount))) {
	//	return 0, 0, fmt.Errorf("insufficient funds: required %d, got %d", requiredAmount, balance.Amount.Int64())
	//}

	return 0, 0, nil
}

func (ms msgServer) executeStartPermissionVP(ctx sdk.Context, msg *types.MsgStartPermissionVP, validatorPerm types.Permission, fees, deposit uint64) (uint64, error) {
	// TODO: After trustdeposit module
	// Increment trust deposit
	//if err := ms.trustDepositKeeper.IncreaseTrustDeposit(ctx, msg.Creator, deposit); err != nil {
	//	return 0, fmt.Errorf("failed to increase trust deposit: %w", err)
	//}

	// Send validation fees to escrow if greater than 0
	//if fees > 0 {
	//	if err := ms.transferToEscrow(ctx, msg.Creator, fees); err != nil {
	//		return 0, fmt.Errorf("failed to transfer fees to escrow: %w", err)
	//	}
	//}

	// Create new permission entry
	now := ctx.BlockTime()
	newPerm := types.Permission{
		SchemaId:          validatorPerm.SchemaId,
		Type:              types.PermissionType(msg.Type),
		Did:               msg.Did,
		Grantee:           msg.Creator,
		Created:           &now,
		CreatedBy:         msg.Creator,
		Extended:          &now,
		ExtendedBy:        msg.Creator,
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
		if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION &&
			validatorType != types.PermissionType_PERMISSION_TYPE_ISSUER_GRANTOR {
			return fmt.Errorf("issuer permission requires ISSUER_GRANTOR validator")
		}
		if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_TRUST_REGISTRY_VALIDATION &&
			validatorType != types.PermissionType_PERMISSION_TYPE_TRUST_REGISTRY {
			return fmt.Errorf("issuer permission requires TRUST_REGISTRY validator")
		}

	case types.PermissionType_PERMISSION_TYPE_ISSUER_GRANTOR:
		if cs.IssuerPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION &&
			validatorType != types.PermissionType_PERMISSION_TYPE_TRUST_REGISTRY {
			return fmt.Errorf("issuer grantor permission requires TRUST_REGISTRY validator")
		}

	case types.PermissionType_PERMISSION_TYPE_VERIFIER:
		if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION &&
			validatorType != types.PermissionType_PERMISSION_TYPE_VERIFIER_GRANTOR {
			return fmt.Errorf("verifier permission requires VERIFIER_GRANTOR validator")
		}
		if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_TRUST_REGISTRY_VALIDATION &&
			validatorType != types.PermissionType_PERMISSION_TYPE_TRUST_REGISTRY {
			return fmt.Errorf("verifier permission requires TRUST_REGISTRY validator")
		}

	case types.PermissionType_PERMISSION_TYPE_VERIFIER_GRANTOR:
		if cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION &&
			validatorType != types.PermissionType_PERMISSION_TYPE_TRUST_REGISTRY {
			return fmt.Errorf("verifier grantor permission requires TRUST_REGISTRY validator")
		}

	case types.PermissionType_PERMISSION_TYPE_HOLDER:
		if (cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION ||
			cs.VerifierPermManagementMode == credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_TRUST_REGISTRY_VALIDATION) &&
			validatorType != types.PermissionType_PERMISSION_TYPE_ISSUER {
			return fmt.Errorf("holder permission requires ISSUER validator")
		}
	}

	return nil
}
