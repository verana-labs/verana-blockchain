package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/x/diddirectory/keeper"
	"github.com/verana-labs/verana-blockchain/x/diddirectory/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, sdk.Context) {
	k, ctx := keepertest.DiddirectoryKeeper(t)
	// Set block time directly on the SDK context
	ctx = ctx.WithBlockTime(time.Now())
	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}

func TestMsgServerAddDID(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	creator := sdk.AccAddress([]byte("test_creator")).String()
	validDID := "did:example:123456789abcdefghi"

	testCases := []struct {
		name    string
		msg     *types.MsgAddDID
		isValid bool
	}{
		{
			name: "Valid Add DID - Default Years",
			msg: &types.MsgAddDID{
				Creator: creator,
				Did:     validDID,
				Years:   0, // Should default to 1
			},
			isValid: true,
		},
		{
			name: "Valid Add DID - Multiple Years",
			msg: &types.MsgAddDID{
				Creator: creator,
				Did:     validDID + "2",
				Years:   5,
			},
			isValid: true,
		},
		{
			name: "Empty DID",
			msg: &types.MsgAddDID{
				Creator: creator,
				Did:     "",
				Years:   1,
			},
			isValid: false,
		},
		{
			name: "Invalid DID Format",
			msg: &types.MsgAddDID{
				Creator: creator,
				Did:     "invalid-did",
				Years:   1,
			},
			isValid: false,
		},
		{
			name: "Years Too High",
			msg: &types.MsgAddDID{
				Creator: creator,
				Did:     validDID + "3",
				Years:   32,
			},
			isValid: false,
		},
		{
			name: "Duplicate DID",
			msg: &types.MsgAddDID{
				Creator: creator,
				Did:     validDID, // Same as first test case
				Years:   1,
			},
			isValid: false,
		},
	}

	// Set default params for testing
	params := types.DefaultParams()
	require.NoError(t, k.SetParams(ctx, params))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.AddDID(ctx, tc.msg)

			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify DID was stored
				storedDID, err := k.DIDDirectory.Get(ctx, tc.msg.Did)
				require.NoError(t, err)
				require.Equal(t, tc.msg.Did, storedDID.Did)
				require.Equal(t, tc.msg.Creator, storedDID.Controller)

				// Check years and expiration
				years := tc.msg.Years
				if years == 0 {
					years = 1
				}
				expectedDeposit := int64(params.DidDirectoryTrustDeposit * uint64(years))
				require.Equal(t, expectedDeposit, storedDID.Deposit)

				// Verify timestamps
				require.False(t, storedDID.Created.IsZero())
				require.False(t, storedDID.Modified.IsZero())
				require.False(t, storedDID.Exp.IsZero())

				// Verify expiration is years from creation
				expectedExp := storedDID.Created.AddDate(int(years), 0, 0)
				require.Equal(t, expectedExp, storedDID.Exp)

			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}
