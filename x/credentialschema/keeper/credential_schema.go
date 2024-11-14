package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/credentialschema/types"
)

func (ms msgServer) validateCreateCredentialSchemaParams(ctx sdk.Context, msg *types.MsgCreateCredentialSchema) error {
	params := ms.GetParams(ctx)

	// Validate trust registry ownership
	tr, err := ms.trustregistryKeeper.GetTrustRegistry(ctx, msg.TrId)
	if err != nil {
		return fmt.Errorf("trust registry not found: %w", err)
	}
	if tr.Controller != msg.Creator {
		return fmt.Errorf("creator is not the controller of the trust registry")
	}

	// Check schema size
	if uint64(len(msg.JsonSchema)) > params.CredentialSchemaSchemaMaxSize {
		return fmt.Errorf("schema size exceeds maximum allowed size of %d bytes", params.CredentialSchemaSchemaMaxSize)
	}

	// Validate validity periods against params
	if err := validateValidityPeriodsWithParams(msg, params); err != nil {
		return fmt.Errorf("invalid validity period: %w", err)
	}

	return nil
}

func validateValidityPeriodsWithParams(msg *types.MsgCreateCredentialSchema, params types.Params) error {
	if msg.IssuerGrantorValidationValidityPeriod > params.CredentialSchemaIssuerGrantorValidationValidityPeriodMaxDays {
		return fmt.Errorf("issuer grantor validation validity period exceeds maximum of %d days",
			params.CredentialSchemaIssuerGrantorValidationValidityPeriodMaxDays)
	}

	if msg.VerifierGrantorValidationValidityPeriod > params.CredentialSchemaVerifierGrantorValidationValidityPeriodMaxDays {
		return fmt.Errorf("verifier grantor validation validity period exceeds maximum of %d days",
			params.CredentialSchemaVerifierGrantorValidationValidityPeriodMaxDays)
	}

	if msg.IssuerValidationValidityPeriod > params.CredentialSchemaIssuerValidationValidityPeriodMaxDays {
		return fmt.Errorf("issuer validation validity period exceeds maximum of %d days",
			params.CredentialSchemaIssuerValidationValidityPeriodMaxDays)
	}

	if msg.VerifierValidationValidityPeriod > params.CredentialSchemaVerifierValidationValidityPeriodMaxDays {
		return fmt.Errorf("verifier validation validity period exceeds maximum of %d days",
			params.CredentialSchemaVerifierValidationValidityPeriodMaxDays)
	}

	if msg.HolderValidationValidityPeriod > params.CredentialSchemaHolderValidationValidityPeriodMaxDays {
		return fmt.Errorf("holder validation validity period exceeds maximum of %d days",
			params.CredentialSchemaHolderValidationValidityPeriodMaxDays)
	}

	return nil
}

func (ms msgServer) executeCreateCredentialSchema(ctx sdk.Context, schemaID uint64, msg *types.MsgCreateCredentialSchema) error {
	// Get params using the getter method
	params := ms.GetParams(ctx)

	// Calculate trust deposit
	trustDeposit := params.CredentialSchemaTrustDeposit

	// Create the credential schema
	credentialSchema := types.CredentialSchema{
		Id:                                      schemaID, // Use the generated ID
		TrId:                                    msg.TrId,
		Created:                                 ctx.BlockTime(),
		Deposit:                                 trustDeposit,
		JsonSchema:                              msg.JsonSchema,
		IssuerGrantorValidationValidityPeriod:   msg.IssuerGrantorValidationValidityPeriod,
		VerifierGrantorValidationValidityPeriod: msg.VerifierGrantorValidationValidityPeriod,
		IssuerValidationValidityPeriod:          msg.IssuerValidationValidityPeriod,
		VerifierValidationValidityPeriod:        msg.VerifierValidationValidityPeriod,
		HolderValidationValidityPeriod:          msg.HolderValidationValidityPeriod,
		IssuerPermManagementMode:                types.CredentialSchemaPermManagementMode(msg.IssuerPermManagementMode),
		VerifierPermManagementMode:              types.CredentialSchemaPermManagementMode(msg.VerifierPermManagementMode),
	}

	// TODO:Handle trust deposit

	// Persist the credential schema using keeper method
	if err := ms.SetCredentialSchema(ctx, credentialSchema); err != nil {
		return fmt.Errorf("failed to persist credential schema: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateCredentialSchema,
			sdk.NewAttribute(types.AttributeKeyId, fmt.Sprintf("%d", schemaID)),
			sdk.NewAttribute(types.AttributeKeyTrId, fmt.Sprintf("%d", msg.TrId)),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return nil
}
