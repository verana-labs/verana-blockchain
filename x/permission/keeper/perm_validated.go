package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cstypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
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

	// Handle fees and trust deposits
	if perm.VpCurrentFees > 0 {
		// Load validator permission
		validatorPerm, err := ms.Keeper.GetPermissionByID(ctx, perm.ValidatorPermId)
		if err != nil {
			return fmt.Errorf("failed to get validator permission: %w", err)
		}

		// Get validator address
		validatorAddr, err := sdk.AccAddressFromBech32(validatorPerm.Grantee)
		if err != nil {
			return fmt.Errorf("invalid validator address: %w", err)
		}

		// Get trust deposit rate - assuming this returns a uint32 value representing a percentage (e.g., 20 for 20%)
		trustDepositRate := ms.trustDeposit.GetTrustDepositRate(ctx)
		validatorTrustDeposit := ms.Keeper.validatorTrustDepositAmount(perm.VpCurrentFees, trustDepositRate)

		// Calculate validator's direct fee portion (excluding trust deposit)
		validatorTrustFees := perm.VpCurrentFees - validatorTrustDeposit

		// Transfer direct fees from module escrow to validator
		if validatorTrustFees > 0 {
			err = ms.bankKeeper.SendCoinsFromModuleToAccount(
				ctx,
				types.ModuleName, // Module escrow account
				validatorAddr,    // Validator account
				sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, int64(validatorTrustFees))),
			)
			if err != nil {
				return fmt.Errorf("failed to transfer fees to validator: %w", err)
			}
		}

		// Increase validator's trust deposit
		if validatorTrustDeposit > 0 {
			err = ms.trustDeposit.AdjustTrustDeposit(
				ctx,
				validatorPerm.Grantee,
				int64(validatorTrustDeposit),
			)
			if err != nil {
				return fmt.Errorf("failed to adjust validator trust deposit: %w", err)
			}

			// Update validator deposit in applicant permission
			perm.VpValidatorDeposit += validatorTrustDeposit
		}
	}

	// Set current fees and deposit to zero after processing
	perm.VpCurrentFees = 0
	perm.VpCurrentDeposit = 0

	return ms.Keeper.UpdatePermission(ctx, perm)
}

func (k Keeper) validatorTrustDepositAmount(vpCurrentFees uint64, trustDepositRate math.LegacyDec) uint64 {
	vpCurrentFeesDec := math.LegacyNewDec(int64(vpCurrentFees))
	validatorTrustDeposit := vpCurrentFeesDec.Mul(trustDepositRate)
	return validatorTrustDeposit.TruncateInt().Uint64()
}
