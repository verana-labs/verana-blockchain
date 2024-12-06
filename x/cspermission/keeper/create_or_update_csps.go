package keeper

import (
	"cosmossdk.io/collections"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
)

func (ms msgServer) validateSessionAccess(ctx sdk.Context, msg *types.MsgCreateOrUpdateCSPS) error {
	existingSession, err := ms.GetCSPSession(ctx, msg.Id)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return err
	}

	if existingSession != nil {
		if existingSession.Controller != msg.Creator {
			return sdkerrors.ErrUnauthorized.Wrap("only session controller can update")
		}

		// Check for duplicate permission pairs
		for _, authz := range existingSession.SessionAuthz {
			if authz.ExecutorPermId == msg.ExecutorPermId &&
				authz.BeneficiaryPermId == msg.BeneficiaryPermId {
				return sdkerrors.ErrInvalidRequest.Wrap("permission pair already exists")
			}
		}
	}

	return nil
}

func (ms msgServer) validateExecutorPerm(ctx sdk.Context, msg *types.MsgCreateOrUpdateCSPS) (*types.CredentialSchemaPerm, error) {
	executorPerm, err := ms.CredentialSchemaPerm.Get(ctx, msg.ExecutorPermId)
	if err != nil {
		return nil, sdkerrors.ErrNotFound.Wrap("executor permission not found")
	}

	if executorPerm.Grantee != msg.Creator {
		return nil, sdkerrors.ErrUnauthorized.Wrap("only grantee can create/update session")
	}

	if executorPerm.CspType != types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER &&
		executorPerm.CspType != types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("executor must be ISSUER or VERIFIER")
	}

	return &executorPerm, nil
}

func (ms msgServer) validateBeneficiaryPerm(ctx sdk.Context, msg *types.MsgCreateOrUpdateCSPS, executorPerm *types.CredentialSchemaPerm) error {
	if executorPerm.CspType == types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER {
		if msg.BeneficiaryPermId != 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("beneficiary not allowed for ISSUER")
		}
		return nil
	}

	// VERIFIER type checks
	if msg.BeneficiaryPermId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("beneficiary required for VERIFIER")
	}

	beneficiaryPerm, err := ms.CredentialSchemaPerm.Get(ctx, msg.BeneficiaryPermId)
	if err != nil {
		return sdkerrors.ErrNotFound.Wrap("beneficiary permission not found")
	}

	if beneficiaryPerm.CspType != types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER {
		return sdkerrors.ErrInvalidRequest.Wrap("beneficiary must be ISSUER")
	}

	if beneficiaryPerm.SchemaId != executorPerm.SchemaId {
		return sdkerrors.ErrInvalidRequest.Wrap("schema IDs must match")
	}

	return nil
}

func (ms msgServer) createOrUpdateSession(ctx sdk.Context, msg *types.MsgCreateOrUpdateCSPS, executorPerm *types.CredentialSchemaPerm) error {
	session := &types.CredentialSchemaPermSession{
		Id:           msg.Id,
		Controller:   msg.Creator,
		UserAgentDid: msg.UserAgentDid,
	}

	existingSession, err := ms.GetCSPSession(ctx, msg.Id)
	if err == nil {
		session = existingSession
	}

	// Add new authorization
	session.SessionAuthz = append(session.SessionAuthz, &types.SessionAuthz{
		ExecutorPermId:     msg.ExecutorPermId,
		BeneficiaryPermId:  msg.BeneficiaryPermId,
		WalletUserAgentDid: msg.WalletUserAgentDid,
	})

	return ms.CredentialSchemaPermSession.Set(ctx, msg.Id, *session)
}
