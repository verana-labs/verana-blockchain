package keeper_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/x/cspermission/keeper"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
)

func setupQueryServer(t testing.TB) (types.QueryServer, keeper.Keeper, types.MsgServer, *keepertest.MockTrustRegistryKeeper, *keepertest.MockCredentialSchemaKeeper, context.Context) {
	k, trKeeper, csKeeper, ctx := keepertest.CspermissionKeeper(t)
	return keeper.NewQueryServerImpl(k), k, keeper.NewMsgServerImpl(k), trKeeper, csKeeper, ctx
}

func TestGetCSP(t *testing.T) {
	qs, _, ms, trKeeper, csKeeper, ctx := setupQueryServer(t)
	creator := "verana1creator"

	// Create prerequisite data
	trId := trKeeper.CreateMockTrustRegistry(creator, "did:example:123")
	schemaId := csKeeper.CreateMockCredentialSchema(trId)

	// Create a test permission
	createMsg := &types.MsgCreateCredentialSchemaPerm{
		Creator:          creator,
		SchemaId:         schemaId,
		CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
		Did:              "did:example:123",
		Grantee:          "verana1grantee",
		EffectiveFrom:    time.Now().Add(time.Hour),
		ValidationId:     1,
		ValidationFees:   100,
		IssuanceFees:     200,
		VerificationFees: 300,
	}
	_, err := ms.CreateCredentialSchemaPerm(ctx, createMsg)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		request       *types.QueryGetCSPRequest
		expectedError bool
		check         func(*testing.T, *types.QueryGetCSPResponse)
	}{
		{
			name: "Valid Request",
			request: &types.QueryGetCSPRequest{
				Id: 1,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryGetCSPResponse) {
				require.NotNil(t, response.Permission)
				require.Equal(t, uint64(1), response.Permission.Id)
				require.Equal(t, "did:example:123", response.Permission.Did)
				require.Equal(t, schemaId, response.Permission.SchemaId)
			},
		},
		{
			name: "Zero ID",
			request: &types.QueryGetCSPRequest{
				Id: 0,
			},
			expectedError: true,
		},
		{
			name: "Non-existent ID",
			request: &types.QueryGetCSPRequest{
				Id: 99,
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := qs.GetCSP(ctx, tc.request)
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

func TestListCSP(t *testing.T) {
	qs, _, ms, trKeeper, csKeeper, ctx := setupQueryServer(t)
	creator := "cosmos1p3yh9uqdjghsdf4jqq38erd5l0kvtss0mq52yk"

	// Create prerequisite data
	trId := trKeeper.CreateMockTrustRegistry(creator, "did:example:123")
	schemaId := csKeeper.CreateMockCredentialSchema(trId)

	// Create multiple test permissions
	createMsg1 := &types.MsgCreateCredentialSchemaPerm{
		Creator:          creator,
		SchemaId:         schemaId,
		CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
		Did:              "did:example:123",
		Grantee:          "verana1grantee",
		EffectiveFrom:    time.Now().Add(time.Hour),
		ValidationId:     1,
		ValidationFees:   100,
		IssuanceFees:     200,
		VerificationFees: 300,
	}
	_, err := ms.CreateCredentialSchemaPerm(ctx, createMsg1)
	require.NoError(t, err)

	createMsg2 := &types.MsgCreateCredentialSchemaPerm{
		Creator:          creator,
		SchemaId:         schemaId,
		CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER),
		Did:              "did:example:456",
		Grantee:          "verana1grantee2",
		EffectiveFrom:    time.Now().Add(time.Hour),
		ValidationId:     1,
		ValidationFees:   100,
		IssuanceFees:     200,
		VerificationFees: 300,
	}
	_, err = ms.CreateCredentialSchemaPerm(ctx, createMsg2)
	require.NoError(t, err)

	testCases := []struct {
		name          string
		request       *types.QueryListCSPRequest
		expectedError bool
		check         func(*testing.T, *types.QueryListCSPResponse)
	}{
		{
			name: "List All for Schema",
			request: &types.QueryListCSPRequest{
				SchemaId:        schemaId,
				ResponseMaxSize: 10,
			},
			expectedError: false,
		},
		{
			name: "Filter by Creator",
			request: &types.QueryListCSPRequest{
				SchemaId:        schemaId,
				Creator:         creator,
				ResponseMaxSize: 10,
			},
			expectedError: false,
		},
		{
			name: "Filter by DID",
			request: &types.QueryListCSPRequest{
				SchemaId:        schemaId,
				Did:             "did:example:123",
				ResponseMaxSize: 10,
			},
			expectedError: false,
		},
		{
			name: "Invalid Schema ID",
			request: &types.QueryListCSPRequest{
				SchemaId: 0,
			},
			expectedError: true,
		},
		{
			name: "Filter by Type",
			request: &types.QueryListCSPRequest{
				SchemaId:        schemaId,
				Type:            types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER,
				ResponseMaxSize: 10,
			},
			expectedError: false,
		},
		{
			name: "Response Size Limit",
			request: &types.QueryListCSPRequest{
				SchemaId:        schemaId,
				ResponseMaxSize: 1,
			},
			expectedError: false,
		},
		{
			name: "Combined Filters",
			request: &types.QueryListCSPRequest{
				SchemaId:        schemaId,
				Creator:         creator,
				Did:             "did:example:123",
				Type:            types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER,
				ResponseMaxSize: 10,
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := qs.ListCSP(ctx, tc.request)
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

func TestIsAuthorizedIssuer(t *testing.T) {
	qs, _, ms, trKeeper, csKeeper, ctx := setupQueryServer(t)
	creator := "verana1creator"

	// Create prerequisite data
	trId := trKeeper.CreateMockTrustRegistry(creator, "did:example:123")
	schemaId := csKeeper.CreateMockCredentialSchema(trId)

	// Create test permission
	createMsg := &types.MsgCreateCredentialSchemaPerm{
		Creator:          creator,
		SchemaId:         schemaId,
		CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
		Did:              "did:example:123",
		Grantee:          "verana1grantee",
		EffectiveFrom:    time.Now().Add(time.Hour),
		ValidationId:     1,
		ValidationFees:   0,
		IssuanceFees:     0,
		VerificationFees: 0,
	}
	_, err := ms.CreateCredentialSchemaPerm(ctx, createMsg)
	require.NoError(t, err)

	testCases := []struct {
		name           string
		request        *types.QueryIsAuthorizedIssuerRequest
		expectedError  bool
		expectedStatus types.AuthorizationStatus
	}{
		{
			name: "Valid Authorization Check",
			request: &types.QueryIsAuthorizedIssuerRequest{
				IssuerDid:          "did:example:123",
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
				SchemaId:           schemaId,
				When:               timePtr(time.Now().Add(2 * time.Hour)),
			},
			expectedError:  false,
			expectedStatus: types.AuthorizationStatus_AUTHORIZED,
		},
		{
			name: "Invalid Schema ID",
			request: &types.QueryIsAuthorizedIssuerRequest{
				IssuerDid:          "did:example:123",
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
				SchemaId:           99999,
			},
			expectedError:  true,
			expectedStatus: types.AuthorizationStatus_UNSPECIFIED,
		},
		{
			name: "Permission Not Yet Effective",
			request: &types.QueryIsAuthorizedIssuerRequest{
				IssuerDid:          "did:example:123",
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
				SchemaId:           schemaId,
				When:               timePtr(time.Now()), // Before effective_from
			},
			expectedError:  false,
			expectedStatus: types.AuthorizationStatus_FORBIDDEN,
		},
		{
			name: "Non-matching Country Filter",
			request: &types.QueryIsAuthorizedIssuerRequest{
				IssuerDid:          "did:example:456",
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
				SchemaId:           schemaId,
				Country:            "GB",
				When:               timePtr(time.Now().Add(2 * time.Hour)),
			},
			expectedError:  false,
			expectedStatus: types.AuthorizationStatus_FORBIDDEN,
		},
		{
			name: "Invalid Country Code Format",
			request: &types.QueryIsAuthorizedIssuerRequest{
				IssuerDid:          "did:example:456",
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
				SchemaId:           schemaId,
				Country:            "USA", // Invalid format, should be 2 letters
				When:               timePtr(time.Now().Add(2 * time.Hour)),
			},
			expectedError:  true,
			expectedStatus: types.AuthorizationStatus_UNSPECIFIED,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := qs.IsAuthorizedIssuer(ctx, tc.request)
			if tc.expectedError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, response)
			fmt.Println("t", t)
			fmt.Println("tc.expectedStatus", tc.expectedStatus)
			fmt.Println("response.Status", response.Status)
			require.Equal(t, tc.expectedStatus, response.Status)
		})
	}
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
