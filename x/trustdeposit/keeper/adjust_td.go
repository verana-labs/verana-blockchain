package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/trustdeposit/types"
)

func (k Keeper) AdjustTrustDeposit(ctx sdk.Context, account string, augend int64) error {
	// Basic validation
	senderAcc, err := sdk.AccAddressFromBech32(account)
	if err != nil {
		return fmt.Errorf("invalid account address: %w", err)
	}
	if account == "" {
		return fmt.Errorf("account cannot be empty")
	}
	if augend == 0 {
		return fmt.Errorf("augend must be non-zero")
	}

	// Load existing trust deposit if it exists
	td, err := k.TrustDeposit.Get(ctx, account)
	if err != nil {
		// If trust deposit doesn't exist and trying to decrease, abort
		if augend < 0 {
			return fmt.Errorf("cannot decrease non-existent trust deposit")
		}

		// Initialize new trust deposit
		td = types.TrustDeposit{
			Account: account,
			Share:   0,
			Amount:  0,
		}
	}

	// Get module params for share value calculation
	params := k.GetParams(ctx)

	if augend > 0 {
		// Handle positive adjustment (increase)
		neededDeposit := uint64(augend)
		if td.Claimable > 0 {
			if td.Claimable >= uint64(augend) {
				// Can cover from claimable amount
				td.Claimable -= uint64(augend)
			} else {
				// Need to transfer additional funds
				neededDeposit = uint64(augend) - td.Claimable

				// Calculate new shares for the additional deposit
				newShares := (neededDeposit * 100) / params.TrustDepositShareValue

				// Transfer tokens from account to module
				if err := k.bankKeeper.SendCoinsFromAccountToModule(
					ctx,
					senderAcc,
					types.ModuleName,
					sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, int64(neededDeposit))),
				); err != nil {
					return fmt.Errorf("failed to transfer tokens: %w", err)
				}

				td.Amount += neededDeposit
				td.Share += newShares
				td.Claimable = 0
			}
		} else {
			// No claimable amount, need to transfer full amount
			newShares := (uint64(augend) * 100) / params.TrustDepositShareValue

			// Transfer tokens from account to module
			if err := k.bankKeeper.SendCoinsFromAccountToModule(
				ctx,
				senderAcc,
				types.ModuleName,
				sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, augend)),
			); err != nil {
				return fmt.Errorf("failed to transfer tokens: %w", err)
			}

			td.Amount += uint64(augend)
			td.Share += newShares
		}
	} else {
		// Handle negative adjustment (decrease)
		// Ensure not trying to decrease more than deposit
		if uint64(-augend) > td.Amount {
			return fmt.Errorf("cannot decrease more than deposited amount")
		}
		td.Claimable = td.Claimable - uint64(-augend)
	}

	// Save updated trust deposit
	if err := k.TrustDeposit.Set(ctx, account, td); err != nil {
		return fmt.Errorf("failed to save trust deposit: %w", err)
	}

	return nil
}
