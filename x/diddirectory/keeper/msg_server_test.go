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

func TestMsgServerRenewDID(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	creator := sdk.AccAddress([]byte("test_creator")).String()
	wrongCreator := sdk.AccAddress([]byte("wrong_creator")).String()
	validDID := "did:example:123456789abcdefghi"

	// Set default params for testing
	params := types.DefaultParams()
	require.NoError(t, k.SetParams(ctx, params))

	// First create a DID
	createMsg := &types.MsgAddDID{
		Creator: creator,
		Did:     validDID,
		Years:   1,
	}
	_, err := ms.AddDID(ctx, createMsg)
	require.NoError(t, err)

	// Get initial DID entry for later comparison
	initialEntry, err := k.DIDDirectory.Get(ctx, validDID)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		msg     *types.MsgRenewDID
		isValid bool
	}{
		{
			name: "Empty DID",
			msg: &types.MsgRenewDID{
				Creator: creator,
				Did:     "",
				Years:   1,
			},
			isValid: false,
		},
		{
			name: "Invalid DID Format",
			msg: &types.MsgRenewDID{
				Creator: creator,
				Did:     "invalid-did",
				Years:   1,
			},
			isValid: false,
		},
		{
			name: "Years Too High",
			msg: &types.MsgRenewDID{
				Creator: creator,
				Did:     validDID,
				Years:   32,
			},
			isValid: false,
		},
		{
			name: "Wrong Controller",
			msg: &types.MsgRenewDID{
				Creator: wrongCreator,
				Did:     validDID,
				Years:   1,
			},
			isValid: false,
		},
		{
			name: "Non-existent DID",
			msg: &types.MsgRenewDID{
				Creator: creator,
				Did:     "did:example:nonexistent",
				Years:   1,
			},
			isValid: false,
		},
		{
			name: "Valid Renewal - Default Years",
			msg: &types.MsgRenewDID{
				Creator: creator,
				Did:     validDID,
				Years:   0, // Should default to 1
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.RenewDID(ctx, tc.msg)

			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Get updated DID entry
				storedDID, err := k.DIDDirectory.Get(ctx, tc.msg.Did)
				require.NoError(t, err)

				// Check years and deposit calculations
				years := tc.msg.Years
				if years == 0 {
					years = 1
				}
				expectedDeposit := initialEntry.Deposit + int64(params.DidDirectoryTrustDeposit*uint64(years))
				require.Equal(t, expectedDeposit, storedDID.Deposit)

				// Verify expiration is extended by years
				expectedExp := initialEntry.Exp.AddDate(int(years), 0, 0)
				require.Equal(t, expectedExp, storedDID.Exp)

				// Store the updated values for next test case
				initialEntry = storedDID

			} else {
				require.Error(t, err)
				require.Nil(t, resp)

				if tc.msg.Did == validDID {
					// Verify DID wasn't modified for invalid attempts
					currentDID, err := k.DIDDirectory.Get(ctx, tc.msg.Did)
					require.NoError(t, err)
					require.Equal(t, initialEntry, currentDID)
				}
			}
		})
	}
}

func TestMsgServerRemoveDID(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	creator := sdk.AccAddress([]byte("test_creator")).String()
	wrongCreator := sdk.AccAddress([]byte("wrong_creator")).String()
	validDID := "did:example:123456789abcdefghi"
	validDID2 := "did:example:987654321abcdefghi"

	// Set default params for testing
	params := types.DefaultParams()
	require.NoError(t, k.SetParams(ctx, params))

	// First create a DID
	createMsg := &types.MsgAddDID{
		Creator: creator,
		Did:     validDID,
		Years:   1,
	}
	_, err := ms.AddDID(ctx, createMsg)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		msg     *types.MsgRemoveDID
		setup   func(*sdk.Context)
		isValid bool
	}{
		{
			name: "Empty DID",
			msg: &types.MsgRemoveDID{
				Creator: creator,
				Did:     "",
			},
			isValid: false,
		},
		{
			name: "Invalid DID Format",
			msg: &types.MsgRemoveDID{
				Creator: creator,
				Did:     "invalid-did",
			},
			isValid: false,
		},
		{
			name: "Non-existent DID",
			msg: &types.MsgRemoveDID{
				Creator: creator,
				Did:     "did:example:nonexistent",
			},
			isValid: false,
		},
		{
			name: "Wrong Creator Before Grace Period",
			msg: &types.MsgRemoveDID{
				Creator: wrongCreator,
				Did:     validDID,
			},
			isValid: false,
		},
		{
			name: "Anyone Can Remove After Grace Period",
			msg: &types.MsgRemoveDID{
				Creator: wrongCreator,
				Did:     validDID,
			},
			setup: func(ctx *sdk.Context) {
				futureTime := time.Now().AddDate(2, 0, 0)
				*ctx = ctx.WithBlockTime(futureTime)
			},
			isValid: true,
		},
		{
			name: "Valid Removal By Controller",
			msg: &types.MsgRemoveDID{
				Creator: creator,
				Did:     validDID2,
			},
			setup: func(ctx *sdk.Context) {
				// Create a new DID for controller removal test
				createMsg := &types.MsgAddDID{
					Creator: creator,
					Did:     validDID2,
					Years:   1,
				}
				_, err := ms.AddDID(*ctx, createMsg)
				require.NoError(t, err)
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(&ctx)
			}

			resp, err := ms.RemoveDID(ctx, tc.msg)

			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify DID was removed
				_, err := k.DIDDirectory.Get(ctx, tc.msg.Did)
				require.Error(t, err) // Should error as DID no longer exists
			} else {
				require.Error(t, err)
				require.Nil(t, resp)

				if tc.msg.Did == validDID {
					// Verify DID still exists for invalid removal attempts
					_, err := k.DIDDirectory.Get(ctx, tc.msg.Did)
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestMsgServerTouchDID(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	creator := sdk.AccAddress([]byte("test_creator")).String()
	otherCreator := sdk.AccAddress([]byte("other_creator")).String()
	validDID := "did:example:123456789abcdefghi"

	// Set default params for testing
	params := types.DefaultParams()
	require.NoError(t, k.SetParams(ctx, params))

	cleanup := func() {
		_ = k.DIDDirectory.Remove(ctx, validDID)
	}

	testCases := []struct {
		name       string
		msg        *types.MsgTouchDID
		beforeTest func()
		isValid    bool
	}{
		{
			name: "Empty DID",
			msg: &types.MsgTouchDID{
				Creator: creator,
				Did:     "",
			},
			isValid: false,
		},
		{
			name: "Invalid DID Format",
			msg: &types.MsgTouchDID{
				Creator: creator,
				Did:     "invalid-did",
			},
			isValid: false,
		},
		{
			name: "Non-existent DID",
			msg: &types.MsgTouchDID{
				Creator: creator,
				Did:     "did:example:nonexistent",
			},
			isValid: false,
		},
		{
			name: "Valid Touch By Original Creator",
			msg: &types.MsgTouchDID{
				Creator: creator,
				Did:     validDID,
			},
			beforeTest: func() {
				createMsg := &types.MsgAddDID{
					Creator: creator,
					Did:     validDID,
					Years:   1,
				}
				_, err := ms.AddDID(ctx, createMsg)
				require.NoError(t, err)
			},
			isValid: true,
		},
		{
			name: "Valid Touch By Other Creator",
			msg: &types.MsgTouchDID{
				Creator: otherCreator,
				Did:     validDID,
			},
			beforeTest: func() {
				createMsg := &types.MsgAddDID{
					Creator: creator,
					Did:     validDID,
					Years:   1,
				}
				_, err := ms.AddDID(ctx, createMsg)
				require.NoError(t, err)
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer cleanup()

			if tc.beforeTest != nil {
				tc.beforeTest()
			}

			// Store initial state if DID exists
			var initialEntry types.DIDDirectory
			var hasInitial bool
			if tc.msg.Did != "" {
				entry, err := k.DIDDirectory.Get(ctx, tc.msg.Did)
				if err == nil {
					initialEntry = entry
					hasInitial = true
				}
			}

			// Move time forward before touching
			newTime := ctx.BlockTime().Add(10 * time.Second)
			ctx = ctx.WithBlockTime(newTime)

			resp, err := ms.TouchDID(ctx, tc.msg)

			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify DID was updated
				updatedEntry, err := k.DIDDirectory.Get(ctx, tc.msg.Did)
				require.NoError(t, err)

				// Modified time should be greater than initial time
				require.True(t, updatedEntry.Modified.After(initialEntry.Modified))

				// Modified time should match new block time
				require.Equal(t, newTime, updatedEntry.Modified)

				// Other fields should remain unchanged
				initialEntry.Modified = updatedEntry.Modified // Set equal for comparison
				require.Equal(t, initialEntry, updatedEntry)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)

				if hasInitial {
					// Verify DID wasn't modified for invalid attempts
					currentEntry, err := k.DIDDirectory.Get(ctx, tc.msg.Did)
					require.NoError(t, err)
					require.Equal(t, initialEntry, currentEntry)
				}
			}
		})
	}
}
