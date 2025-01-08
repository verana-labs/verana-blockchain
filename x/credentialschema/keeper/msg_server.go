package keeper

import (
	"context"
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/credentialschema/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (ms msgServer) CreateCredentialSchema(goCtx context.Context, msg *types.MsgCreateCredentialSchema) (*types.MsgCreateCredentialSchemaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Generate next ID
	nextID, err := ms.GetNextID(ctx, "cs")
	if err != nil {
		return nil, fmt.Errorf("failed to generate schema ID: %w", err)
	}

	// [MOD-CS-MSG-1-2-1] Basic checks
	if err := ms.validateCreateCredentialSchemaParams(ctx, msg); err != nil {
		return nil, err
	}

	// [MOD-CS-MSG-1-2-2] Fee checks
	//if err := ms.checkSufficientFees(ctx, msg.Creator); err != nil {
	// return nil, err
	//}

	// [MOD-CS-MSG-1-3] Execution
	if err := ms.executeCreateCredentialSchema(ctx, nextID, msg); err != nil {
		return nil, err
	}

	return &types.MsgCreateCredentialSchemaResponse{
		Id: nextID,
	}, nil
}

func (ms msgServer) UpdateCredentialSchema(goCtx context.Context, msg *types.MsgUpdateCredentialSchema) (*types.MsgUpdateCredentialSchemaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get credential schema
	cs, err := ms.CredentialSchema.Get(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("credential schema not found: %w", err)
	}

	// Check trust registry controller
	tr, err := ms.trustregistryKeeper.GetTrustRegistry(ctx, cs.TrId)
	if err != nil {
		return nil, fmt.Errorf("trust registry not found: %w", err)
	}
	if tr.Controller != msg.Creator {
		return nil, fmt.Errorf("only trust registry controller can update credential schema")
	}

	// Validate validity periods against params
	params := ms.GetParams(ctx)
	if err := ValidateValidityPeriods(
		params,
		msg.IssuerGrantorValidationValidityPeriod,
		msg.VerifierGrantorValidationValidityPeriod,
		msg.IssuerValidationValidityPeriod,
		msg.VerifierValidationValidityPeriod,
		msg.HolderValidationValidityPeriod,
	); err != nil {
		return nil, err
	}

	// Update credential schema
	cs.IssuerGrantorValidationValidityPeriod = msg.IssuerGrantorValidationValidityPeriod
	cs.VerifierGrantorValidationValidityPeriod = msg.VerifierGrantorValidationValidityPeriod
	cs.IssuerValidationValidityPeriod = msg.IssuerValidationValidityPeriod
	cs.VerifierValidationValidityPeriod = msg.VerifierValidationValidityPeriod
	cs.HolderValidationValidityPeriod = msg.HolderValidationValidityPeriod
	cs.Modified = ctx.BlockTime()

	if err := ms.CredentialSchema.Set(ctx, cs.Id, cs); err != nil {
		return nil, fmt.Errorf("failed to update credential schema: %w", err)
	}

	return &types.MsgUpdateCredentialSchemaResponse{}, nil
}

// ValidateValidityPeriods checks if all validity periods are within allowed ranges
func ValidateValidityPeriods(
	params types.Params,
	issuerGrantorPeriod,
	verifierGrantorPeriod,
	issuerPeriod,
	verifierPeriod,
	holderPeriod uint32,
) error {
	if issuerGrantorPeriod > params.CredentialSchemaIssuerGrantorValidationValidityPeriodMaxDays {
		return errors.New("issuer grantor validation validity period exceeds maximum allowed days")
	}
	if verifierGrantorPeriod > params.CredentialSchemaVerifierGrantorValidationValidityPeriodMaxDays {
		return errors.New("verifier grantor validation validity period exceeds maximum allowed days")
	}
	if issuerPeriod > params.CredentialSchemaIssuerValidationValidityPeriodMaxDays {
		return errors.New("issuer validation validity period exceeds maximum allowed days")
	}
	if verifierPeriod > params.CredentialSchemaVerifierValidationValidityPeriodMaxDays {
		return errors.New("verifier validation validity period exceeds maximum allowed days")
	}
	if holderPeriod > params.CredentialSchemaHolderValidationValidityPeriodMaxDays {
		return errors.New("holder validation validity period exceeds maximum allowed days")
	}
	return nil
}

func (ms msgServer) ArchiveCredentialSchema(goCtx context.Context, msg *types.MsgArchiveCredentialSchema) (*types.MsgArchiveCredentialSchemaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get credential schema
	cs, err := ms.CredentialSchema.Get(ctx, msg.Id)
	if err != nil {
		return nil, fmt.Errorf("credential schema not found: %w", err)
	}

	// Check trust registry controller
	tr, err := ms.trustregistryKeeper.GetTrustRegistry(ctx, cs.TrId)
	if err != nil {
		return nil, fmt.Errorf("trust registry not found: %w", err)
	}
	if tr.Controller != msg.Creator {
		return nil, fmt.Errorf("only trust registry controller can archive credential schema")
	}

	// Check archive state
	if msg.Archive {
		if cs.Archived != nil {
			return nil, fmt.Errorf("credential schema is already archived")
		}
	} else {
		if cs.Archived == nil {
			return nil, fmt.Errorf("credential schema is not archived")
		}
	}

	// Update archive state
	now := ctx.BlockTime()
	if msg.Archive {
		cs.Archived = &now
	} else {
		cs.Archived = nil
	}
	cs.Modified = now

	// Save updated credential schema
	if err := ms.CredentialSchema.Set(ctx, cs.Id, cs); err != nil {
		return nil, fmt.Errorf("failed to update credential schema: %w", err)
	}

	return &types.MsgArchiveCredentialSchemaResponse{}, nil
}
