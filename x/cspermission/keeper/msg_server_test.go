package keeper_test

import (
	"context"
	"testing"
	"time"

	_ "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	_ "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"github.com/verana-labs/verana-blockchain/x/cspermission/keeper"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, *keepertest.MockTrustRegistryKeeper, *keepertest.MockCredentialSchemaKeeper, context.Context) {
	k, trKeeper, csKeeper, ctx := keepertest.CspermissionKeeper(t)
	return k, keeper.NewMsgServerImpl(k), trKeeper, csKeeper, ctx
}

func TestCreateCredentialSchemaPerm(t *testing.T) {
	_, ms, trKeeper, csKeeper, ctx := setupMsgServer(t)
	creator := "verana1creator"

	// Create mock trust registry
	trId := trKeeper.CreateMockTrustRegistry(creator, "did:example:123")

	// Create mock credential schema
	schemaId := csKeeper.CreateMockCredentialSchema(trId)

	testCases := []struct {
		name    string
		msg     *types.MsgCreateCredentialSchemaPerm
		setup   func()
		expPass bool
	}{
		{
			name: "Valid Trust Registry Permission",
			msg: &types.MsgCreateCredentialSchemaPerm{
				Creator:          creator,
				SchemaId:         schemaId,
				CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY),
				Did:              "did:example:123",
				Grantee:          "verana1grantee",
				EffectiveFrom:    time.Now().Add(time.Hour).UTC(),
				ValidationFees:   100,
				IssuanceFees:     200,
				VerificationFees: 300,
			},
			expPass: true,
		},
		{
			name: "Invalid Schema ID",
			msg: &types.MsgCreateCredentialSchemaPerm{
				Creator:          creator,
				SchemaId:         9999,
				CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY),
				Did:              "did:example:123",
				Grantee:          "verana1grantee",
				EffectiveFrom:    time.Now().Add(time.Hour).UTC(),
				ValidationFees:   100,
				IssuanceFees:     200,
				VerificationFees: 300,
			},
			expPass: false,
		},
		{
			name: "Invalid Permission Type",
			msg: &types.MsgCreateCredentialSchemaPerm{
				Creator:          creator,
				SchemaId:         schemaId,
				CspType:          99,
				Did:              "did:example:123",
				Grantee:          "verana1grantee",
				EffectiveFrom:    time.Now().Add(time.Hour).UTC(),
				ValidationFees:   100,
				IssuanceFees:     200,
				VerificationFees: 300,
			},
			expPass: false,
		},
		{
			name: "Non-Controller Creating Trust Registry Permission",
			msg: &types.MsgCreateCredentialSchemaPerm{
				Creator:          "verana1notcontroller",
				SchemaId:         schemaId,
				CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY),
				Did:              "did:example:123",
				Grantee:          "verana1grantee",
				EffectiveFrom:    time.Now().Add(time.Hour).UTC(),
				ValidationFees:   100,
				IssuanceFees:     200,
				VerificationFees: 300,
			},
			expPass: false,
		},
		{
			name: "Past Effective Date",
			msg: &types.MsgCreateCredentialSchemaPerm{
				Creator:          creator,
				SchemaId:         schemaId,
				CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY),
				Did:              "did:example:123",
				Grantee:          "verana1grantee",
				EffectiveFrom:    time.Now().Add(-time.Hour).UTC(),
				ValidationFees:   100,
				IssuanceFees:     200,
				VerificationFees: 300,
			},
			expPass: false,
		},
		{
			name: "Invalid DID Mismatch",
			msg: &types.MsgCreateCredentialSchemaPerm{
				Creator:          creator,
				SchemaId:         schemaId,
				CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY),
				Did:              "did:example:wrongdid",
				Grantee:          "verana1grantee",
				EffectiveFrom:    time.Now().Add(time.Hour).UTC(),
				ValidationFees:   100,
				IssuanceFees:     200,
				VerificationFees: 300,
			},
			expPass: false,
		},
		{
			name: "Valid Issuer Permission with Validation",
			msg: &types.MsgCreateCredentialSchemaPerm{
				Creator:          creator,
				SchemaId:         schemaId,
				CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
				Did:              "did:example:123",
				Grantee:          "verana1grantee",
				EffectiveFrom:    time.Now().Add(time.Hour).UTC(),
				ValidationId:     1,
				ValidationFees:   100,
				IssuanceFees:     200,
				VerificationFees: 300,
			},
			expPass: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}

			resp, err := ms.CreateCredentialSchemaPerm(ctx, tc.msg)
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, resp)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}

func TestCreateCredentialSchemaPermWithOverlap(t *testing.T) {
	_, ms, trKeeper, csKeeper, ctx := setupMsgServer(t)
	creator := "verana1creator"

	trId := trKeeper.CreateMockTrustRegistry(creator, "did:example:123")
	schemaId := csKeeper.CreateMockCredentialSchema(trId)

	// Create initial permission
	baseTime := time.Now().UTC()
	effectiveFrom1 := baseTime.Add(5 * time.Hour)  // First perm starts at +5h
	effectiveUntil1 := baseTime.Add(6 * time.Hour) // First perm ends at +6h

	msg1 := &types.MsgCreateCredentialSchemaPerm{
		Creator:          creator,
		SchemaId:         schemaId,
		CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
		Did:              "did:example:123",
		Grantee:          "verana1grantee",
		EffectiveFrom:    effectiveFrom1,
		EffectiveUntil:   &effectiveUntil1,
		ValidationId:     1,
		ValidationFees:   100,
		IssuanceFees:     200,
		VerificationFees: 300,
	}

	resp, err := ms.CreateCredentialSchemaPerm(ctx, msg1)
	require.NoError(t, err)
	require.NotNil(t, resp)

	testCases := []struct {
		name    string
		msg     *types.MsgCreateCredentialSchemaPerm
		expPass bool
	}{
		{
			name: "Overlapping Time Period",
			msg: func() *types.MsgCreateCredentialSchemaPerm {
				overlapFrom := baseTime.Add(5 * time.Hour)  // During first permission's period
				overlapUntil := baseTime.Add(7 * time.Hour) // After first permission ends
				return &types.MsgCreateCredentialSchemaPerm{
					Creator:          creator,
					SchemaId:         schemaId,
					CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
					Did:              "did:example:123",
					Grantee:          "verana1grantee",
					EffectiveFrom:    overlapFrom,
					EffectiveUntil:   &overlapUntil,
					ValidationId:     1,
					ValidationFees:   100,
					IssuanceFees:     200,
					VerificationFees: 300,
				}
			}(),
			expPass: false,
		},
		{
			name: "Non-overlapping Time Period",
			msg: func() *types.MsgCreateCredentialSchemaPerm {
				// Per spec, new permission must end before existing permission starts
				nonOverlapFrom := baseTime.Add(12 * time.Hour)  // Starts at +12h
				nonOverlapUntil := baseTime.Add(13 * time.Hour) // Ends at +13h

				t.Logf("\nBase time: %v", baseTime)
				t.Logf("First permission: %v to %v", effectiveFrom1, effectiveUntil1)
				t.Logf("New permission: %v to %v", nonOverlapFrom, nonOverlapUntil)

				return &types.MsgCreateCredentialSchemaPerm{
					Creator:          creator,
					SchemaId:         schemaId,
					CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
					Did:              "did:example:123",
					Grantee:          "verana1grantee",
					EffectiveFrom:    nonOverlapFrom,
					EffectiveUntil:   &nonOverlapUntil,
					ValidationId:     1,
					ValidationFees:   100,
					IssuanceFees:     200,
					VerificationFees: 300,
				}
			}(),
			expPass: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("\nTesting: %s\nExisting Permission: %v to %v\nNew Permission: %v to %v",
				tc.name,
				effectiveFrom1, effectiveUntil1,
				tc.msg.EffectiveFrom, tc.msg.EffectiveUntil)

			resp, err := ms.CreateCredentialSchemaPerm(ctx, tc.msg)
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, resp)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}

func TestRevokeCredentialSchemaPerm(t *testing.T) {
	k, ms, trKeeper, csKeeper, ctx := setupMsgServer(t)
	creator := "verana1creator"

	// Create mock trust registry
	trId := trKeeper.CreateMockTrustRegistry(creator, "did:example:123")

	// Create mock credential schema
	schemaId := csKeeper.CreateMockCredentialSchema(trId)

	// First create a permission that we'll try to revoke
	baseTime := time.Now().UTC()
	effectiveFrom := baseTime.Add(time.Hour)
	effectiveUntil := baseTime.Add(2 * time.Hour)

	createMsg := &types.MsgCreateCredentialSchemaPerm{
		Creator:          creator,
		SchemaId:         schemaId,
		CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY),
		Did:              "did:example:123",
		Grantee:          "verana1grantee",
		EffectiveFrom:    effectiveFrom,
		EffectiveUntil:   &effectiveUntil,
		ValidationFees:   100,
		IssuanceFees:     200,
		VerificationFees: 300,
	}

	resp, err := ms.CreateCredentialSchemaPerm(ctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, resp)

	testCases := []struct {
		name    string
		msg     *types.MsgRevokeCredentialSchemaPerm
		expPass bool
	}{
		{
			name: "Valid Revocation By Controller",
			msg: &types.MsgRevokeCredentialSchemaPerm{
				Creator: creator,
				Id:      1, // First created permission
			},
			expPass: true,
		},
		{
			name: "Non-existent Permission ID",
			msg: &types.MsgRevokeCredentialSchemaPerm{
				Creator: creator,
				Id:      99,
			},
			expPass: false,
		},
		{
			name: "Unauthorized Revocation Attempt",
			msg: &types.MsgRevokeCredentialSchemaPerm{
				Creator: "verana1unauthorized",
				Id:      1,
			},
			expPass: false,
		},
		{
			name: "Already Revoked Permission",
			msg: &types.MsgRevokeCredentialSchemaPerm{
				Creator: creator,
				Id:      1, // Try to revoke the same permission again
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.RevokeCredentialSchemaPerm(ctx, tc.msg)
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify the permission was actually revoked
				perm, err := k.CredentialSchemaPerm.Get(ctx, tc.msg.Id)
				require.NoError(t, err)
				require.NotNil(t, perm.Revoked)
				require.Equal(t, tc.msg.Creator, perm.RevokedBy)
				require.Zero(t, perm.Deposit)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}

// Test revocation for different permission types
func TestRevokeCredentialSchemaPermTypes(t *testing.T) {
	k, ms, trKeeper, csKeeper, ctx := setupMsgServer(t)
	creator := "verana1creator"

	trId := trKeeper.CreateMockTrustRegistry(creator, "did:example:123")
	schemaId := csKeeper.CreateMockCredentialSchema(trId)

	baseTime := time.Now().UTC()
	effectiveFrom := baseTime.Add(time.Hour)
	effectiveUntil := baseTime.Add(2 * time.Hour)

	// Keep track of permission IDs
	var currentId uint64 = 1

	testCases := []struct {
		name          string
		setupMsg      *types.MsgCreateCredentialSchemaPerm
		revokeCreator string
		expPass       bool
	}{
		{
			name: "Revoke ISSUER Permission",
			setupMsg: &types.MsgCreateCredentialSchemaPerm{
				Creator:          creator,
				SchemaId:         schemaId,
				CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
				Did:              "did:example:123",
				Grantee:          "verana1grantee",
				EffectiveFrom:    effectiveFrom,
				EffectiveUntil:   &effectiveUntil,
				ValidationId:     1,
				ValidationFees:   100,
				IssuanceFees:     200,
				VerificationFees: 300,
			},
			revokeCreator: creator,
			expPass:       true,
		},
		{
			name: "Revoke VERIFIER Permission",
			setupMsg: &types.MsgCreateCredentialSchemaPerm{
				Creator:          creator,
				SchemaId:         schemaId,
				CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER),
				Did:              "did:example:123",
				Grantee:          "verana1grantee",
				EffectiveFrom:    effectiveFrom,
				EffectiveUntil:   &effectiveUntil,
				ValidationId:     1,
				ValidationFees:   100,
				IssuanceFees:     200,
				VerificationFees: 300,
			},
			revokeCreator: creator,
			expPass:       true,
		},
		{
			name: "Revoke HOLDER Permission",
			setupMsg: &types.MsgCreateCredentialSchemaPerm{
				Creator:          creator,
				SchemaId:         schemaId,
				CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_HOLDER),
				Did:              "did:example:123",
				Grantee:          "verana1grantee",
				EffectiveFrom:    effectiveFrom,
				EffectiveUntil:   &effectiveUntil,
				ValidationId:     1,
				ValidationFees:   100,
				IssuanceFees:     200,
				VerificationFees: 300,
			},
			revokeCreator: creator,
			expPass:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create permission
			resp, err := ms.CreateCredentialSchemaPerm(ctx, tc.setupMsg)
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Try to revoke it using current ID
			revokeMsg := &types.MsgRevokeCredentialSchemaPerm{
				Creator: tc.revokeCreator,
				Id:      currentId,
			}

			revokeResp, err := ms.RevokeCredentialSchemaPerm(ctx, revokeMsg)
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, revokeResp)

				// Verify revocation
				perm, err := k.CredentialSchemaPerm.Get(ctx, currentId)
				require.NoError(t, err)
				require.NotNil(t, perm.Revoked)
				require.Equal(t, tc.revokeCreator, perm.RevokedBy)
			} else {
				require.Error(t, err)
				require.Nil(t, revokeResp)
			}

			// Increment ID for next test case
			currentId++
		})
	}
}

func TestTerminateCSP(t *testing.T) {
	k, ms, trKeeper, csKeeper, ctx := setupMsgServer(t)
	creator := "verana1creator"
	grantee := "verana1grantee"

	// Create mock trust registry and schema
	trId := trKeeper.CreateMockTrustRegistry(creator, "did:example:123")
	schemaId := csKeeper.CreateMockCredentialSchema(trId)

	// Create a permission to terminate
	baseTime := time.Now().UTC()
	effectiveFrom := baseTime.Add(time.Hour)
	effectiveUntil := baseTime.Add(2 * time.Hour)

	createMsg := &types.MsgCreateCredentialSchemaPerm{
		Creator:          creator,
		SchemaId:         schemaId,
		CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
		Did:              "did:example:123",
		Grantee:          grantee,
		EffectiveFrom:    effectiveFrom,
		EffectiveUntil:   &effectiveUntil,
		ValidationId:     1,
		ValidationFees:   100,
		IssuanceFees:     200,
		VerificationFees: 300,
	}

	resp, err := ms.CreateCredentialSchemaPerm(ctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, resp)

	testCases := []struct {
		name    string
		msg     *types.MsgTerminateCredentialSchemaPerm
		setup   func()
		expPass bool
	}{
		{
			name: "Valid Termination By Grantee",
			msg: &types.MsgTerminateCredentialSchemaPerm{
				Creator: grantee,
				Id:      1,
			},
			expPass: true,
		},
		{
			name: "Non-existent Permission ID",
			msg: &types.MsgTerminateCredentialSchemaPerm{
				Creator: grantee,
				Id:      99,
			},
			expPass: false,
		},
		{
			name: "Unauthorized Termination Attempt",
			msg: &types.MsgTerminateCredentialSchemaPerm{
				Creator: "verana1unauthorized",
				Id:      1,
			},
			expPass: false,
		},
		{
			name: "Already Terminated Permission",
			msg: &types.MsgTerminateCredentialSchemaPerm{
				Creator: grantee,
				Id:      1,
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}

			resp, err := ms.TerminateCredentialSchemaPerm(ctx, tc.msg)
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify termination
				perm, err := k.CredentialSchemaPerm.Get(ctx, tc.msg.Id)
				require.NoError(t, err)
				require.NotNil(t, perm.Terminated)
				require.Equal(t, tc.msg.Creator, perm.TerminatedBy)
				require.Zero(t, perm.Deposit)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}

func TestCreateOrUpdateCSPS(t *testing.T) {
	_, ms, trKeeper, csKeeper, ctx := setupMsgServer(t)
	creator := "verana1creator"

	// Create prerequisite data
	trId := trKeeper.CreateMockTrustRegistry(creator, "did:example:123")
	schemaId := csKeeper.CreateMockCredentialSchema(trId)

	// Create test permission for executor
	baseTime := time.Now().UTC()
	effectiveUntil := baseTime.Add(2 * time.Hour)

	createPermMsg := &types.MsgCreateCredentialSchemaPerm{
		Creator:          creator,
		SchemaId:         schemaId,
		CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
		Did:              "did:example:123",
		Grantee:          "verana1grantee",
		EffectiveFrom:    baseTime.Add(time.Hour),
		EffectiveUntil:   &effectiveUntil, // Set expiry time
		ValidationId:     1,
		ValidationFees:   0,
		IssuanceFees:     0,
		VerificationFees: 0,
	}

	_, err := ms.CreateCredentialSchemaPerm(ctx, createPermMsg)
	require.NoError(t, err)

	// Create another permission for update test
	effectiveUntil2 := baseTime.Add(4 * time.Hour)

	createPermMsg2 := &types.MsgCreateCredentialSchemaPerm{
		Creator:          creator,
		SchemaId:         schemaId,
		CspType:          uint32(types.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER),
		Did:              "did:example:234",
		Grantee:          "verana1grantee",
		EffectiveFrom:    baseTime.Add(3 * time.Hour), // Start after first permission ends
		EffectiveUntil:   &effectiveUntil2,
		ValidationId:     1,
		ValidationFees:   0,
		IssuanceFees:     0,
		VerificationFees: 0,
	}
	_, err = ms.CreateCredentialSchemaPerm(ctx, createPermMsg2)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		msg     *types.MsgCreateOrUpdateCSPS
		setup   func() // Additional setup for specific test cases
		expPass bool
	}{
		{
			name: "Valid ISSUER CSPS Creation",
			msg: &types.MsgCreateOrUpdateCSPS{
				Creator:            "verana1grantee",
				Id:                 "123e4567-e89b-12d3-a456-426614174000",
				ExecutorPermId:     1,
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
			},

			expPass: true,
		},
		{
			name: "Invalid UUID Format",
			msg: &types.MsgCreateOrUpdateCSPS{
				Id:                 "invalid-uuid",
				ExecutorPermId:     1,
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
			},
			expPass: false,
		},
		{
			name: "Non-existent Executor Permission",
			msg: &types.MsgCreateOrUpdateCSPS{
				Id:                 "123e4567-e89b-12d3-a456-426614174001",
				ExecutorPermId:     99,
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
			},
			expPass: false,
		},
		{
			name: "Invalid Beneficiary for ISSUER",
			msg: &types.MsgCreateOrUpdateCSPS{
				Id:                 "123e4567-e89b-12d3-a456-426614174002",
				ExecutorPermId:     1,
				BeneficiaryPermId:  2, // ISSUER type should not have beneficiary
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
			},
			expPass: false,
		},
		{
			name: "Missing User Agent DID",
			msg: &types.MsgCreateOrUpdateCSPS{
				Id:                 "123e4567-e89b-12d3-a456-426614174003",
				ExecutorPermId:     1,
				WalletUserAgentDid: "did:example:wallet",
			},
			expPass: false,
		},
		{
			name: "Missing Wallet User Agent DID",
			msg: &types.MsgCreateOrUpdateCSPS{
				Id:             "123e4567-e89b-12d3-a456-426614174004",
				ExecutorPermId: 1,
				UserAgentDid:   "did:example:agent",
			},
			expPass: false,
		},
		{
			name: "Duplicate Permission Pair",
			setup: func() {
				initialMsg := &types.MsgCreateOrUpdateCSPS{
					Creator:            "verana1grantee",
					Id:                 "123e4567-e89b-12d3-a456-426614174004",
					ExecutorPermId:     1,
					UserAgentDid:       "did:example:agent",
					WalletUserAgentDid: "did:example:wallet",
				}
				_, err := ms.CreateOrUpdateCSPS(ctx, initialMsg)
				require.NoError(t, err)
			},
			msg: &types.MsgCreateOrUpdateCSPS{
				Creator:            "verana1grantee",
				Id:                 "123e4567-e89b-12d3-a456-426614174004",
				ExecutorPermId:     1, // Same executor_perm_id
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
			},
			expPass: false,
		},
		{
			name: "Valid Update Existing Session",
			setup: func() {

				// Create initial session
				initialMsg := &types.MsgCreateOrUpdateCSPS{
					Creator:            "verana1grantee",
					Id:                 "123e4567-e89b-12d3-a456-426614174005",
					ExecutorPermId:     1,
					UserAgentDid:       "did:example:agent",
					WalletUserAgentDid: "did:example:wallet",
				}
				_, err = ms.CreateOrUpdateCSPS(ctx, initialMsg)
				require.NoError(t, err)
			},
			msg: &types.MsgCreateOrUpdateCSPS{
				Creator:            "verana1grantee",
				Id:                 "123e4567-e89b-12d3-a456-426614174005",
				ExecutorPermId:     2, // Different executor permission
				UserAgentDid:       "did:example:agent",
				WalletUserAgentDid: "did:example:wallet",
			},
			expPass: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}

			resp, err := ms.CreateOrUpdateCSPS(ctx, tc.msg)
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, resp)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}
