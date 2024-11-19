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
				overlapFrom := baseTime.Add(4 * time.Hour)  // Starts before first perm
				overlapUntil := baseTime.Add(7 * time.Hour) // Ends after first perm
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
				nonOverlapFrom := baseTime.Add(2 * time.Hour)   // Starts at +2h
				nonOverlapUntil := baseTime.Add(13 * time.Hour) // Ends at +13h (before first perm starts at +5h)

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
