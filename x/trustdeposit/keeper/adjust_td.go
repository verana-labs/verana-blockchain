package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/trustdeposit/types"
)

// AdjustTrustDeposit modifies the trust deposit for an account by the specified amount.
// If augend is positive, it increases the trust deposit.
// If augend is negative, it decreases the trust deposit and increases the claimable amount.
//
// The function follows the specification [MOD-TD-MSG-1] from the Verana blockchain specs.
//
// Parameters:
// - ctx: The SDK context
// - account: The account address as a Bech32 string
// - augend: The amount to adjust (positive for increase, negative for decrease)
//
// Returns:
// - error: If the operation fails
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

	// Get total shares and total deposits across all accounts
	totalShares, totalDeposits := k.GetTotalSharesAndDeposits(ctx)

	if err != nil {
		// If trust deposit doesn't exist and trying to decrease, abort
		if augend < 0 {
			return fmt.Errorf("cannot decrease non-existent trust deposit")
		}

		// Initialize new trust deposit
		td = types.TrustDeposit{
			Account:   account,
			Share:     0,
			Amount:    0,
			Claimable: 0,
		}
	}

	// Convert uint fields to int64 for calculations
	amount := int64(td.Amount)
	claimable := int64(td.Claimable)
	share := int64(td.Share)

	if augend > 0 {
		// Handle positive adjustment (increase)
		if claimable > 0 {
			if claimable >= augend {
				// Can cover from claimable amount
				claimable -= augend
			} else {
				// Need to transfer additional funds
				int64NeededDeposit := augend - claimable

				// Calculate new shares using dynamic share calculation
				// Shares per Token = TotalShares / TotalDeposits
				var newShares int64
				if totalDeposits > 0 {
					// Use dynamic ratio based on total shares and deposits
					newShares = (int64NeededDeposit * totalShares) / totalDeposits
				} else {
					// If no total deposits yet, use 1:1 ratio (first deposit)
					newShares = int64NeededDeposit
				}

				// Transfer tokens from account to module
				err := k.bankKeeper.SendCoinsFromAccountToModule(
					ctx,
					senderAcc,
					types.ModuleName,
					sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, int64NeededDeposit)),
				)
				if err != nil {
					return fmt.Errorf("failed to transfer tokens: %w", err)
				}

				amount += int64NeededDeposit
				share += newShares
				claimable = 0
			}
		} else {
			// No claimable amount, need to transfer full amount

			// Calculate new shares using dynamic share calculation
			// Shares per Token = TotalShares / TotalDeposits
			var newShares int64
			if totalDeposits > 0 {
				// Use dynamic ratio based on total shares and deposits
				newShares = (augend * totalShares) / totalDeposits
			} else {
				// If no total deposits yet, use 1:1 ratio (first deposit)
				newShares = augend
			}
			// Transfer tokens from account to module
			err := k.bankKeeper.SendCoinsFromAccountToModule(
				ctx,
				senderAcc,
				types.ModuleName,
				sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, augend)),
			)
			if err != nil {
				return fmt.Errorf("failed to transfer tokens: %w", err)
			}

			amount += augend
			share += newShares
		}
	} else { // augend < 0
		// Handle negative adjustment (decrease)
		// Get absolute value of augend for easier handling
		absAugend := -augend

		// if augend is negative and td.claimable - augend > td.deposit transaction MUST abort
		if claimable+absAugend > amount {
			return fmt.Errorf("claimable after adjustment would exceed deposit: %d > %d", claimable+absAugend, amount)
		}

		// Since augend is negative, we add absAugend to claimable
		// This implements "set td.claimable to td.claimable - augend" when augend is negative
		claimable += absAugend
	}

	// Convert back to uint for storage and ensure no negative values
	if amount < 0 {
		return fmt.Errorf("amount cannot be negative after adjustment: %d", amount)
	}
	if claimable < 0 {
		return fmt.Errorf("claimable amount cannot be negative after adjustment: %d", claimable)
	}
	if share < 0 {
		return fmt.Errorf("share cannot be negative after adjustment: %d", share)
	}

	td.Amount = uint64(amount)
	td.Claimable = uint64(claimable)
	td.Share = uint64(share)

	// Save updated trust deposit
	err = k.TrustDeposit.Set(ctx, account, td)
	if err != nil {
		return fmt.Errorf("failed to save trust deposit: %w", err)
	}

	return nil
}

// GetTotalSharesAndDeposits calculates the total shares and total deposits across all accounts
func (k Keeper) GetTotalSharesAndDeposits(ctx sdk.Context) (int64, int64) {
	var totalShares int64
	var totalDeposits int64

	// Use Walk function to iterate through all trust deposits
	_ = k.TrustDeposit.Walk(ctx, nil, func(key string, value types.TrustDeposit) (bool, error) {
		totalShares += int64(value.Share)
		totalDeposits += int64(value.Amount)
		return false, nil // Continue iteration
	})

	return totalShares, totalDeposits
}
