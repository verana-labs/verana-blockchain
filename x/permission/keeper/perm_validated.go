package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cstypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
	"time"
)

func getValidityPeriod(permType uint32, cs cstypes.CredentialSchema) uint32 {
	switch permType {
	case 3: // ISSUER_GRANTOR
		return cs.IssuerGrantorValidationValidityPeriod
	case 4: // VERIFIER_GRANTOR
		return cs.VerifierGrantorValidationValidityPeriod
	case 1: // ISSUER
		return cs.IssuerValidationValidityPeriod
	case 2: // VERIFIER
		return cs.VerifierValidationValidityPeriod
	case 6: // HOLDER
		return cs.HolderValidationValidityPeriod
	default:
		return 0
	}
}

func calculateVPExp(currentVPExp *time.Time, validityPeriod uint64, now time.Time) *time.Time {
	if validityPeriod == 0 {
		return nil
	}

	var exp time.Time
	if currentVPExp == nil {
		exp = now.AddDate(0, 0, int(validityPeriod))
	} else {
		exp = currentVPExp.AddDate(0, 0, int(validityPeriod))
	}
	return &exp
}

func (ms msgServer) executeSetPermissionVPToValidated(ctx sdk.Context, perm types.Permission, msg *types.MsgSetPermissionVPToValidated, now time.Time, vpExp *time.Time) error {
	// Update permission
	perm.Modified = &now
	perm.VpState = types.ValidationState_VALIDATION_STATE_VALIDATED
	perm.VpLastStateChange = &now
	perm.VpCurrentFees = 0
	perm.VpCurrentDeposit = 0
	perm.VpSummaryDigestSri = msg.VpSummaryDigestSri
	perm.VpExp = vpExp
	perm.EffectiveUntil = msg.EffectiveUntil

	// Set initial values if not a renewal
	if perm.EffectiveFrom == nil {
		perm.ValidationFees = msg.ValidationFees
		perm.IssuanceFees = msg.IssuanceFees
		perm.VerificationFees = msg.VerificationFees
		perm.Country = msg.Country
		perm.EffectiveFrom = &now
	}

	// TODO: Handle fees and trust deposits after trust deposit module is ready
	// validatorTrustFees := perm.VpCurrentFees * (1 - GlobalVariables.TrustDepositRate)
	// validatorTrustDeposit := perm.VpCurrentFees - validatorTrustFees
	// Transfer fees from escrow to validator
	// Update validator deposit

	return ms.Keeper.UpdatePermission(ctx, perm)
}
