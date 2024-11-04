package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

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
		{
			name: "Invalid Language Format",
			msg: &types.MsgCreateTrustRegistry{
				Creator:  creator,
				Did:      validDid,
				Language: "invalid-language-format",
				DocUrl:   "http://example.com/doc",
				DocHash:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
			isValid: false,
		},
		{
			name: "Invalid AKA URI",
			msg: &types.MsgCreateTrustRegistry{
				Creator:  creator,
				Did:      validDid,
				Aka:      "invalid-uri",
				Language: "en",
				DocUrl:   "http://example.com/doc",
				DocHash:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
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

				// Get ID from DID index
				id, err := k.TrustRegistryDIDIndex.Get(ctx, tc.msg.Did)
				require.NoError(t, err)

				// Get trust registry using ID
				tr, err := k.TrustRegistry.Get(ctx, id)
				require.NoError(t, err)
				require.Equal(t, tc.msg.Did, tr.Did)
				require.Equal(t, tc.msg.Creator, tr.Controller)
				require.Equal(t, int32(1), tr.ActiveVersion)
				require.Equal(t, tc.msg.Language, tr.Language)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}

func TestMsgServerAddGovernanceFrameworkDocument(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	creator := sdk.AccAddress([]byte("test_creator")).String()
	validDid := "did:example:123456789abcdefghi"

	// First, create a trust registry
	createMsg := &types.MsgCreateTrustRegistry{
		Creator:  creator,
		Did:      validDid,
		Language: "en",
		DocUrl:   "http://example.com/doc",
		DocHash:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}
	_, err := ms.CreateTrustRegistry(ctx, createMsg)
	require.NoError(t, err)

	// Get trust registry ID
	trID, err := k.TrustRegistryDIDIndex.Get(ctx, validDid)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		setupFunc func() // Additional setup for test case
		msg       *types.MsgAddGovernanceFrameworkDocument
		isValid   bool
	}{
		{
			name: "Valid Add Document with Next Version",
			msg: &types.MsgAddGovernanceFrameworkDocument{
				Creator:     creator,
				TrId:        trID,
				DocLanguage: "en",
				DocUrl:      "http://example.com/doc2",
				DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				Version:     2, // Exactly maxVersion + 1
			},
			isValid: true,
		},
		{
			name: "Valid Add Document to Same Version with Different Language",
			msg: &types.MsgAddGovernanceFrameworkDocument{
				Creator:     creator,
				TrId:        trID,
				DocLanguage: "fr",
				DocUrl:      "http://example.com/doc2-fr",
				DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				Version:     2, // Same version, different language
			},
			isValid: true,
		},
		{
			name: "Valid Add Next Version",
			msg: &types.MsgAddGovernanceFrameworkDocument{
				Creator:     creator,
				TrId:        trID,
				DocLanguage: "en",
				DocUrl:      "http://example.com/doc3",
				DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				Version:     3, // Exactly maxVersion + 1
			},
			isValid: true,
		},
		{
			name: "Invalid Version (Less than Active Version)",
			msg: &types.MsgAddGovernanceFrameworkDocument{
				Creator:     creator,
				TrId:        trID,
				DocLanguage: "en",
				DocUrl:      "http://example.com/doc-old",
				DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				Version:     1,
			},
			isValid: false,
		},
		{
			name: "Invalid Trust Registry ID",
			msg: &types.MsgAddGovernanceFrameworkDocument{
				Creator:     creator,
				TrId:        99999,
				DocLanguage: "en",
				DocUrl:      "http://example.com/doc2",
				DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				Version:     2,
			},
			isValid: false,
		},
		{
			name: "Invalid Language Format",
			msg: &types.MsgAddGovernanceFrameworkDocument{
				Creator:     creator,
				TrId:        trID,
				DocLanguage: "invalid-language",
				DocUrl:      "http://example.com/doc2",
				DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				Version:     2,
			},
			isValid: false,
		},
		{
			name: "Wrong Controller",
			msg: &types.MsgAddGovernanceFrameworkDocument{
				Creator:     "wrong-controller",
				TrId:        trID,
				DocLanguage: "en",
				DocUrl:      "http://example.com/doc2",
				DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				Version:     2,
			},
			isValid: false,
		},
		{
			name: "Invalid Version (Skipping Version)",
			setupFunc: func() {
				// Add version 3 document first
				msg := &types.MsgAddGovernanceFrameworkDocument{
					Creator:     creator,
					TrId:        trID,
					DocLanguage: "en",
					DocUrl:      "http://example.com/doc3",
					DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
					Version:     3,
				}
				_, err := ms.AddGovernanceFrameworkDocument(ctx, msg)
				require.NoError(t, err)
			},
			msg: &types.MsgAddGovernanceFrameworkDocument{
				Creator:     creator,
				TrId:        trID,
				DocLanguage: "en",
				DocUrl:      "http://example.com/doc5",
				DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				Version:     5, // Invalid: should be 4 (maxVersion + 1)
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFunc != nil {
				tc.setupFunc()
			}

			resp, err := ms.AddGovernanceFrameworkDocument(ctx, tc.msg)
			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify document was added
				var found bool
				err = k.GFDocument.Walk(ctx, nil, func(id uint64, gfd types.GovernanceFrameworkDocument) (bool, error) {
					if gfd.Language == tc.msg.DocLanguage && gfd.Url == tc.msg.DocUrl {
						found = true
						return true, nil
					}
					return false, nil
				})
				require.NoError(t, err)
				require.True(t, found)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}

func TestMsgServerIncreaseActiveGovernanceFrameworkVersion(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)

	creator := sdk.AccAddress([]byte("test_creator")).String()
	validDid := "did:example:123456789abcdefghi"

	// Create initial trust registry
	createMsg := &types.MsgCreateTrustRegistry{
		Creator:  creator,
		Did:      validDid,
		Language: "en",
		DocUrl:   "http://example.com/doc",
		DocHash:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}
	_, err := ms.CreateTrustRegistry(ctx, createMsg)
	require.NoError(t, err)

	// Get trust registry ID
	trID, err := k.TrustRegistryDIDIndex.Get(ctx, validDid)
	require.NoError(t, err)

	// Add version 2 documents
	addGFDocMsg := &types.MsgAddGovernanceFrameworkDocument{
		Creator:     creator,
		TrId:        trID,
		DocLanguage: "es", // First add Spanish version
		DocUrl:      "http://example.com/doc2-es",
		DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		Version:     2,
	}
	_, err = ms.AddGovernanceFrameworkDocument(ctx, addGFDocMsg)
	require.NoError(t, err)

	// Test cases for version increase
	testCases := []struct {
		name      string
		setupFunc func() // Additional setup for test case
		msg       *types.MsgIncreaseActiveGovernanceFrameworkVersion
		isValid   bool
	}{
		{
			name: "Cannot Increase Version - Missing Default Language Document",
			msg: &types.MsgIncreaseActiveGovernanceFrameworkVersion{
				Creator: creator,
				TrId:    trID,
			},
			isValid: false,
		},
		{
			name: "Valid Version Increase",
			setupFunc: func() {
				// Add English (default language) document for version 2
				msg := &types.MsgAddGovernanceFrameworkDocument{
					Creator:     creator,
					TrId:        trID,
					DocLanguage: "en",
					DocUrl:      "http://example.com/doc2-en",
					DocHash:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
					Version:     2,
				}
				_, err := ms.AddGovernanceFrameworkDocument(ctx, msg)
				require.NoError(t, err)
			},
			msg: &types.MsgIncreaseActiveGovernanceFrameworkVersion{
				Creator: creator,
				TrId:    trID,
			},
			isValid: true,
		},
		{
			name: "Wrong Controller",
			msg: &types.MsgIncreaseActiveGovernanceFrameworkVersion{
				Creator: "wrong-controller",
				TrId:    trID,
			},
			isValid: false,
		},
		{
			name: "Non-existent Trust Registry",
			msg: &types.MsgIncreaseActiveGovernanceFrameworkVersion{
				Creator: creator,
				TrId:    99999,
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFunc != nil {
				tc.setupFunc()
			}

			resp, err := ms.IncreaseActiveGovernanceFrameworkVersion(ctx, tc.msg)
			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify version increase
				tr, err := k.TrustRegistry.Get(ctx, tc.msg.TrId)
				require.NoError(t, err)
				require.Equal(t, int32(2), tr.ActiveVersion)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}
