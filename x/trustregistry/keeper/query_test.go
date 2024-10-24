package keeper_test

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/keeper"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
	"testing"
)

func setupTestData(t *testing.T) (keeper.Keeper, types.QueryServer, context.Context, uint64) {
	k, ctx := keepertest.TrustregistryKeeper(t)
	qs := keeper.NewQueryServerImpl(k)
	ms := keeper.NewMsgServerImpl(k)

	// Create a trust registry
	creator := sdk.AccAddress([]byte("test_creator")).String()
	createMsg := &types.MsgCreateTrustRegistry{
		Creator:  creator,
		Did:      "did:example:123",
		Language: "en",
		DocUrl:   "http://example.com/doc1",
		DocHash:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}
	_, err := ms.CreateTrustRegistry(ctx, createMsg)
	require.NoError(t, err)

	// Get the trust registry ID
	trID, err := k.TrustRegistryDIDIndex.Get(ctx, "did:example:123")
	require.NoError(t, err)

	// Add documents in different languages for version 2
	addDocMsg := &types.MsgAddGovernanceFrameworkDocument{
		Creator:     creator,
		TrId:        trID,
		DocLanguage: "en",
		DocUrl:      "http://example.com/doc2-en",
		DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		Version:     2,
	}
	_, err = ms.AddGovernanceFrameworkDocument(ctx, addDocMsg)
	require.NoError(t, err)

	addDocMsg.DocLanguage = "es"
	addDocMsg.DocUrl = "http://example.com/doc2-es"
	_, err = ms.AddGovernanceFrameworkDocument(ctx, addDocMsg)
	require.NoError(t, err)

	return k, qs, ctx, trID
}

func TestGetTrustRegistry(t *testing.T) {
	_, qs, ctx, trID := setupTestData(t)

	testCases := []struct {
		name          string
		request       *types.QueryGetTrustRegistryRequest
		expectedError bool
		check         func(*testing.T, *types.QueryGetTrustRegistryResponse)
	}{
		{
			name: "Valid Request - All Documents",
			request: &types.QueryGetTrustRegistryRequest{
				TrId:         trID,
				ActiveGfOnly: false,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryGetTrustRegistryResponse) {
				require.NotNil(t, response.TrustRegistry)
				require.Equal(t, trID, response.TrustRegistry.Id)
				require.Len(t, response.Versions, 2)  // Version 1 and 2
				require.Len(t, response.Documents, 3) // 1 doc for v1, 2 docs for v2
			},
		},
		{
			name: "Valid Request - Active Only",
			request: &types.QueryGetTrustRegistryRequest{
				TrId:         trID,
				ActiveGfOnly: true,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryGetTrustRegistryResponse) {
				require.NotNil(t, response.TrustRegistry)
				require.Len(t, response.Versions, 1)
				require.Equal(t, int32(1), response.Versions[0].Version)
				require.Len(t, response.Documents, 1)
			},
		},
		{
			name: "Valid Request - Preferred Language",
			request: &types.QueryGetTrustRegistryRequest{
				TrId:              trID,
				PreferredLanguage: "es",
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryGetTrustRegistryResponse) {
				require.NotNil(t, response.TrustRegistry)
				for _, doc := range response.Documents {
					if doc.GfvId == 2 { // For version 2
						require.Equal(t, "es", doc.Language)
					}
				}
			},
		},
		{
			name: "Invalid Trust Registry ID",
			request: &types.QueryGetTrustRegistryRequest{
				TrId: 99999,
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := qs.GetTrustRegistry(ctx, tc.request)
			if tc.expectedError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, response)
			if tc.check != nil {
				tc.check(t, response)
			}
		})
	}
}

func TestGetTrustRegistryWithDID(t *testing.T) {
	_, qs, ctx, _ := setupTestData(t)

	testCases := []struct {
		name          string
		request       *types.QueryGetTrustRegistryWithDIDRequest
		expectedError bool
		check         func(*testing.T, *types.QueryGetTrustRegistryResponse)
	}{
		{
			name: "Valid Request",
			request: &types.QueryGetTrustRegistryWithDIDRequest{
				Did:          "did:example:123",
				ActiveGfOnly: false,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryGetTrustRegistryResponse) {
				require.NotNil(t, response.TrustRegistry)
				require.Equal(t, "did:example:123", response.TrustRegistry.Did)
			},
		},
		{
			name: "Invalid DID",
			request: &types.QueryGetTrustRegistryWithDIDRequest{
				Did: "invalid-did",
			},
			expectedError: true,
		},
		{
			name: "Non-existent DID",
			request: &types.QueryGetTrustRegistryWithDIDRequest{
				Did: "did:example:nonexistent",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := qs.GetTrustRegistryWithDID(ctx, tc.request)
			if tc.expectedError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, response)
			if tc.check != nil {
				tc.check(t, response)
			}
		})
	}
}

func TestListTrustRegistries(t *testing.T) {
	k, qs, ctx, _ := setupTestData(t)

	// Create additional trust registry for testing
	ms := keeper.NewMsgServerImpl(k)
	createMsg := &types.MsgCreateTrustRegistry{
		Creator:  "another_creator",
		Did:      "did:example:456",
		Language: "fr",
		DocUrl:   "http://example.com/doc-fr",
		DocHash:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}
	_, err := ms.CreateTrustRegistry(ctx, createMsg)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		request       *types.QueryListTrustRegistriesRequest
		expectedError bool
		check         func(*testing.T, *types.QueryListTrustRegistriesResponse)
	}{
		{
			name: "List All",
			request: &types.QueryListTrustRegistriesRequest{
				ResponseMaxSize: 10,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryListTrustRegistriesResponse) {
				require.Len(t, response.TrustRegistries, 2)
			},
		},
		{
			name: "Filter by Controller",
			request: &types.QueryListTrustRegistriesRequest{
				Controller:      "another_creator",
				ResponseMaxSize: 10,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryListTrustRegistriesResponse) {
				require.Len(t, response.TrustRegistries, 1)
				require.Equal(t, "another_creator", response.TrustRegistries[0].Controller)
			},
		},
		{
			name: "Invalid Response Max Size",
			request: &types.QueryListTrustRegistriesRequest{
				ResponseMaxSize: 1025, // More than maximum allowed
			},
			expectedError: true,
		},
		{
			name: "Default Response Max Size",
			request: &types.QueryListTrustRegistriesRequest{
				ResponseMaxSize: 10, // More than 2
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryListTrustRegistriesResponse) {
				require.Len(t, response.TrustRegistries, 2)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := qs.ListTrustRegistries(ctx, tc.request)
			if tc.expectedError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, response)
			if tc.check != nil {
				tc.check(t, response)
			}
		})
	}
}

func TestParams(t *testing.T) {
	_, qs, ctx, _ := setupTestData(t)

	response, err := qs.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Params)
}
