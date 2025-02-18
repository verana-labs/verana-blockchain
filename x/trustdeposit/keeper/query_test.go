package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/x/trustdeposit/types"
	"testing"
)

func TestGetTrustDeposit(t *testing.T) {
	keeper, ctx := keepertest.TrustdepositKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	// Create test account address
	testAddr := sdk.AccAddress([]byte("test_address")).String()

	// Test with non-existent trust deposit
	resp1, err := keeper.GetTrustDeposit(wctx, &types.QueryGetTrustDepositRequest{
		Account: testAddr,
	})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp1.TrustDeposit.Amount)
	require.Equal(t, uint64(0), resp1.TrustDeposit.Share)
	require.Equal(t, uint64(0), resp1.TrustDeposit.Claimable)

	// Create a trust deposit
	trustDeposit := types.TrustDeposit{
		Account:   testAddr,
		Share:     100,
		Amount:    1000,
		Claimable: 50,
	}
	err = keeper.TrustDeposit.Set(ctx, testAddr, trustDeposit)
	require.NoError(t, err)

	// Test with existing trust deposit
	resp2, err := keeper.GetTrustDeposit(wctx, &types.QueryGetTrustDepositRequest{
		Account: testAddr,
	})
	require.NoError(t, err)
	require.Equal(t, trustDeposit, resp2.TrustDeposit)

	// Test with invalid account address
	_, err = keeper.GetTrustDeposit(wctx, &types.QueryGetTrustDepositRequest{
		Account: "invalid_address",
	})
	require.Error(t, err)
}
