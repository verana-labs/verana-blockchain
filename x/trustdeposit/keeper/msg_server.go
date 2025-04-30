package keeper

import (
	"context"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/trustdeposit/types"
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

func (ms msgServer) ReclaimTrustDepositInterests(goCtx context.Context, msg *types.MsgReclaimTrustDepositInterests) (*types.MsgReclaimTrustDepositInterestsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get account running the method
	account := msg.Creator

	// Load TrustDeposit entry
	td, err := ms.Keeper.TrustDeposit.Get(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("trust deposit not found for account: %s", account)
	}

	// Get module params for share value calculation
	params := ms.Keeper.GetParams(ctx)

	// Calculate claimable interest using decimal math
	depositAmount := ms.Keeper.ShareToAmount(td.Share, params.TrustDepositShareValue)

	// Guards against underflow
	if depositAmount <= td.Amount {
		return nil, fmt.Errorf("no claimable interest available")
	}

	claimableInterest := depositAmount - td.Amount

	// Calculate shares to reduce using decimal math
	sharesToReduce := ms.Keeper.AmountToShare(claimableInterest, params.TrustDepositShareValue)

	// Update trust deposit shares
	td.Share -= sharesToReduce

	// Transfer claimable interest from module to account
	coins := sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, int64(claimableInterest)))
	if err := ms.Keeper.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		sdk.AccAddress(account),
		coins,
	); err != nil {
		return nil, fmt.Errorf("failed to transfer interest: %w", err)
	}

	// Save updated trust deposit
	if err := ms.Keeper.TrustDeposit.Set(ctx, account, td); err != nil {
		return nil, fmt.Errorf("failed to update trust deposit: %w", err)
	}

	return &types.MsgReclaimTrustDepositInterestsResponse{
		ClaimedAmount: claimableInterest,
	}, nil
}

// ShareToAmount converts share value to amount using decimal math
func (k Keeper) ShareToAmount(share uint64, shareValue math.LegacyDec) uint64 {
	shareDec := math.LegacyNewDec(int64(share))
	amountDec := shareDec.Mul(shareValue)
	return amountDec.TruncateInt().Uint64()
}

// AmountToShare converts amount to share value using decimal math
func (k Keeper) AmountToShare(amount uint64, shareValue math.LegacyDec) uint64 {
	amountDec := math.LegacyNewDec(int64(amount))
	if shareValue.IsZero() {
		return 0 // Prevent division by zero
	}
	shareDec := amountDec.Quo(shareValue)
	return shareDec.TruncateInt().Uint64()
}

func (ms msgServer) ReclaimTrustDeposit(goCtx context.Context, msg *types.MsgReclaimTrustDeposit) (*types.MsgReclaimTrustDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validations
	if msg.Claimed == 0 {
		return nil, fmt.Errorf("claimed amount must be greater than 0")
	}

	// Get account running the method
	account := msg.Creator

	// Load TrustDeposit entry
	td, err := ms.Keeper.TrustDeposit.Get(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("trust deposit not found for account: %s", account)
	}

	// Check if claimed amount is within claimable balance
	if td.Claimable < msg.Claimed {
		return nil, fmt.Errorf("claimed amount exceeds claimable balance")
	}

	// Get module params for calculations
	params := ms.Keeper.GetParams(ctx)

	// Calculate required minimum deposit using decimal math
	requiredMinDeposit := ms.Keeper.ShareToAmount(td.Share, params.TrustDepositShareValue)

	if td.Amount < msg.Claimed {
		return nil, fmt.Errorf("amount less than claimed")
	}

	if requiredMinDeposit < (td.Amount - msg.Claimed) {
		return nil, fmt.Errorf("insufficient required minimum deposit")
	}

	// Calculate burn amount and transfer amount using decimal math
	toBurn := ms.Keeper.CalculateBurnAmount(msg.Claimed, params.TrustDepositReclaimBurnRate)
	toTransfer := msg.Claimed - toBurn

	// Calculate share reduction using decimal math
	shareReduction := ms.Keeper.AmountToShare(msg.Claimed, params.TrustDepositShareValue)

	// Update trust deposit
	td.Claimable -= msg.Claimed
	td.Amount -= msg.Claimed
	td.Share -= shareReduction

	// Transfer claimable amount minus burn to the account
	if toTransfer > 0 {
		transferCoins := sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, int64(toTransfer)))
		if err := ms.Keeper.bankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			sdk.AccAddress(account),
			transferCoins,
		); err != nil {
			return nil, fmt.Errorf("failed to transfer coins: %w", err)
		}
	}

	// Burn the calculated amount
	if toBurn > 0 {
		burnCoins := sdk.NewCoins(sdk.NewInt64Coin(types.BondDenom, int64(toBurn)))
		if err := ms.Keeper.bankKeeper.BurnCoins(
			ctx,
			types.ModuleName,
			burnCoins,
		); err != nil {
			return nil, fmt.Errorf("failed to burn coins: %w", err)
		}
	}

	// Save updated trust deposit
	if err := ms.Keeper.TrustDeposit.Set(ctx, account, td); err != nil {
		return nil, fmt.Errorf("failed to update trust deposit: %w", err)
	}

	return &types.MsgReclaimTrustDepositResponse{
		BurnedAmount:  toBurn,
		ClaimedAmount: toTransfer,
	}, nil
}

// CalculateBurnAmount applies burn rate to claimed amount using decimal math
func (k Keeper) CalculateBurnAmount(claimed uint64, burnRate math.LegacyDec) uint64 {
	claimedDec := math.LegacyNewDec(int64(claimed))
	burnAmountDec := claimedDec.Mul(burnRate)
	return burnAmountDec.TruncateInt().Uint64()
}
