package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/x/credentialschema/keeper"
	"github.com/verana-labs/verana-blockchain/x/credentialschema/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, *keepertest.MockTrustRegistryKeeper, context.Context) {
	k, mockTrk, ctx := keepertest.CredentialschemaKeeper(t)
	return k, keeper.NewMsgServerImpl(k), mockTrk, ctx
}

func TestMsgServerCreateCredentialSchema(t *testing.T) {
	k, ms, mockTrk, ctx := setupMsgServer(t)

	creator := sdk.AccAddress([]byte("test_creator")).String()
	validDid := "did:example:123456789abcdefghi"

	// First create a trust registry
	trID := mockTrk.CreateMockTrustRegistry(creator, validDid)

	validJsonSchema := `{
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "/dtr/v1/cs/js/1",
        "type": "object",
        "$defs": {},
        "properties": {
            "name": {
                "type": "string"
            }
        },
        "required": ["name"],
        "additionalProperties": false
    }`

	testCases := []struct {
		name    string
		msg     *types.MsgCreateCredentialSchema
		isValid bool
	}{
		{
			name: "Valid Create Credential Schema",
			msg: &types.MsgCreateCredentialSchema{
				Creator:                                 creator,
				Id:                                      1,
				TrId:                                    trID, // Use the ID from trust registry response
				JsonSchema:                              validJsonSchema,
				IssuerGrantorValidationValidityPeriod:   365,
				VerifierGrantorValidationValidityPeriod: 365,
				IssuerValidationValidityPeriod:          180,
				VerifierValidationValidityPeriod:        180,
				HolderValidationValidityPeriod:          180,
				IssuerPermManagementMode:                2,
				VerifierPermManagementMode:              2,
			},
			isValid: true,
		},
		{
			name: "Non-existent Trust Registry",
			msg: &types.MsgCreateCredentialSchema{
				Creator:                                 creator,
				Id:                                      2,
				TrId:                                    999, // Non-existent trust registry
				JsonSchema:                              validJsonSchema,
				IssuerGrantorValidationValidityPeriod:   365,
				VerifierGrantorValidationValidityPeriod: 365,
				IssuerValidationValidityPeriod:          180,
				VerifierValidationValidityPeriod:        180,
				HolderValidationValidityPeriod:          180,
				IssuerPermManagementMode:                2,
				VerifierPermManagementMode:              2,
			},
			isValid: false,
		},
		{
			name: "Wrong Trust Registry Controller",
			msg: &types.MsgCreateCredentialSchema{
				Creator:                                 sdk.AccAddress([]byte("wrong_creator")).String(),
				Id:                                      3,
				TrId:                                    trID,
				JsonSchema:                              validJsonSchema,
				IssuerGrantorValidationValidityPeriod:   365,
				VerifierGrantorValidationValidityPeriod: 365,
				IssuerValidationValidityPeriod:          180,
				VerifierValidationValidityPeriod:        180,
				HolderValidationValidityPeriod:          180,
				IssuerPermManagementMode:                2,
				VerifierPermManagementMode:              2,
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.CreateCredentialSchema(ctx, tc.msg)
			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify schema was created
				schema, err := k.CredentialSchema.Get(ctx, tc.msg.Id)
				require.NoError(t, err)
				require.Equal(t, tc.msg.JsonSchema, schema.JsonSchema)
				require.Equal(t, tc.msg.IssuerPermManagementMode, uint32(schema.IssuerPermManagementMode))
				require.Equal(t, tc.msg.VerifierPermManagementMode, uint32(schema.VerifierPermManagementMode))
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}
