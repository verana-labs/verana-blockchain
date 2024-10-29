package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/x/diddirectory/types"
)

func TestListDIDs(t *testing.T) {
	keeper, ctx := keepertest.DiddirectoryKeeper(t)

	// Set up test data
	now := time.Now()
	ctx = ctx.WithBlockTime(now)

	// Create test DIDs with different attributes
	dids := []struct {
		did        string
		controller string
		modified   time.Time
		expired    bool
		overGrace  bool
	}{
		{
			did:        "did:example:1",
			controller: "cosmos1controller1",
			modified:   now.Add(-1 * time.Hour),
			expired:    false,
			overGrace:  false,
		},
		{
			did:        "did:example:2",
			controller: "cosmos1controller2",
			modified:   now.Add(-2 * time.Hour),
			expired:    true,
			overGrace:  false,
		},
		{
			did:        "did:example:3",
			controller: "cosmos1controller1",
			modified:   now.Add(-3 * time.Hour),
			expired:    true,
			overGrace:  true,
		},
	}

	params := types.DefaultParams()
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	// Store test DIDs
	for _, d := range dids {
		expTime := now.Add(24 * time.Hour) // Future expiry for non-expired
		if d.expired {
			if d.overGrace {
				expTime = now.AddDate(0, 0, -int(params.DidDirectoryGracePeriod)-1) // Past grace period
			} else {
				expTime = now.AddDate(0, 0, -1) // Just expired
			}
		}

		didEntry := types.DIDDirectory{
			Did:        d.did,
			Controller: d.controller,
			Created:    d.modified,
			Modified:   d.modified,
			Exp:        expTime,
			Deposit:    5,
		}
		err = keeper.DIDDirectory.Set(ctx, d.did, didEntry)
		require.NoError(t, err)
	}

	testCases := []struct {
		name      string
		req       *types.QueryListDIDsRequest
		expected  int
		expectErr bool
	}{
		{
			name: "List All DIDs",
			req: &types.QueryListDIDsRequest{
				ResponseMaxSize: 10,
			},
			expected: 3,
		},
		{
			name: "Filter by Controller",
			req: &types.QueryListDIDsRequest{
				Account:         "cosmos1controller1",
				ResponseMaxSize: 10,
			},
			expected: 2,
		},
		{
			name: "Filter by Changed Time",
			req: &types.QueryListDIDsRequest{
				Changed:         &now,
				ResponseMaxSize: 10,
			},
			expected: 0,
		},
		{
			name: "Filter Expired",
			req: &types.QueryListDIDsRequest{
				Expired:         true,
				OverGrace:       false,
				ResponseMaxSize: 10,
			},
			expected: 2, // Should include both expired DIDs
		},
		{
			name: "Filter Over Grace",
			req: &types.QueryListDIDsRequest{
				Expired:         true,
				OverGrace:       true,
				ResponseMaxSize: 10,
			},
			expected: 1, // Should include only the over-grace DID
		},
		{
			name: "Invalid Response Size",
			req: &types.QueryListDIDsRequest{
				ResponseMaxSize: 1025,
			},
			expectErr: true,
		},
		{
			name:     "List with Default Response Size",
			req:      &types.QueryListDIDsRequest{},
			expected: 3,
		},
		{
			name: "List with Small Response Size",
			req: &types.QueryListDIDsRequest{
				ResponseMaxSize: 2,
			},
			expected: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := keeper.ListDIDs(ctx, tc.req)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, response)
			require.Len(t, response.Dids, tc.expected)

			// Verify sorting by modified time
			for i := 1; i < len(response.Dids); i++ {
				require.True(t, response.Dids[i-1].Modified.Before(response.Dids[i].Modified))
			}
		})
	}
}
