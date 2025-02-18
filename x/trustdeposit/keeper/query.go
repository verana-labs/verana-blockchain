package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/trustdeposit/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) GetTrustDeposit(goCtx context.Context, req *types.QueryGetTrustDepositRequest) (*types.QueryGetTrustDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate account address
	if _, err := sdk.AccAddressFromBech32(req.Account); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid account address: %s", err))
	}

	// Get trust deposit for account
	trustDeposit, err := k.TrustDeposit.Get(ctx, req.Account)
	if err != nil {
		// If not found, return zero values
		trustDeposit = types.TrustDeposit{
			Account:   req.Account,
			Share:     0,
			Amount:    0,
			Claimable: 0,
		}
	}

	return &types.QueryGetTrustDepositResponse{
		TrustDeposit: trustDeposit,
	}, nil
}
