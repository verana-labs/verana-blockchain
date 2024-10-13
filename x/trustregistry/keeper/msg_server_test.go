package keeper_test

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"

	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/keeper"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, context.Context) {
	k, ctx := keepertest.TrustregistryKeeper(t)
	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}

func TestMsgServerCreateTrustRegistry(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	creator := sdk.AccAddress([]byte("test_creator")).String()
	validDid := "did:example:123456789abcdefghi"

	testCases := []struct {
		name    string
		msg     *types.MsgCreateTrustRegistry
		isValid bool
	}{
		{
			name: "Valid Create Trust Registry",
			msg: &types.MsgCreateTrustRegistry{
				Creator:  creator,
				Did:      validDid,
				Aka:      "http://example.com",
				Language: "en",
				DocUrl:   "http://example.com/doc",
				DocHash:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
			isValid: true,
		},
		{
			name: "Invalid DID",
			msg: &types.MsgCreateTrustRegistry{
				Creator:  creator,
				Did:      "invalid-did",
				Language: "en",
				DocUrl:   "http://example.com/doc",
				DocHash:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
			isValid: false,
		},
		{
			name: "Missing Language",
			msg: &types.MsgCreateTrustRegistry{
				Creator: creator,
				Did:     validDid,
				DocUrl:  "http://example.com/doc",
				DocHash: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.CreateTrustRegistry(ctx, tc.msg)
			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, resp)
				// Check if the trust registry was actually created
				tr, err := k.TrustRegistry.Get(ctx, tc.msg.Did)
				require.NoError(t, err)
				require.Equal(t, tc.msg.Did, tr.Did)
				require.Equal(t, tc.msg.Creator, tr.Controller)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}
