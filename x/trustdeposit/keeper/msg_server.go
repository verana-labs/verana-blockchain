package keeper

import (
	"context"
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

	// Calculate claimable interest
	claimableInterest := (td.Share * params.TrustDepositShareValue) - td.Amount
	if claimableInterest <= 0 {
		return nil, fmt.Errorf("no claimable interest available")
	}

	// Calculate shares to reduce
	sharesToReduce := claimableInterest / params.TrustDepositShareValue

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

	// Calculate required minimum deposit
	requiredMinDeposit := td.Share * params.TrustDepositShareValue
	if requiredMinDeposit < td.Amount-msg.Claimed {
		return nil, fmt.Errorf("insufficient required minimum deposit")
	}

	// Calculate burn amount and transfer amount
	toBurn := (msg.Claimed * uint64(params.TrustDepositReclaimBurnRate)) / 100
	toTransfer := msg.Claimed - toBurn

	// Update trust deposit
	td.Claimable -= msg.Claimed
	td.Amount -= msg.Claimed
	td.Share -= msg.Claimed / params.TrustDepositShareValue

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
