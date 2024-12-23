package keeper_test

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	csptypes "github.com/verana-labs/verana-blockchain/x/cspermission/types"
	"github.com/verana-labs/verana-blockchain/x/validation/keeper"
	"github.com/verana-labs/verana-blockchain/x/validation/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, *keepertest.MockCsPermissionKeeper, *keepertest.MockCredentialSchemaKeeper, context.Context) {
	k, csPermKeeper, csKeeper, ctx := keepertest.ValidationKeeper(t)
	return k, keeper.NewMsgServerImpl(k), csPermKeeper, csKeeper, ctx
}

func TestCreateValidation(t *testing.T) {
	k, ms, csPermKeeper, csKeeper, ctx := setupMsgServer(t)
	creator := "verana1creator"

	// Create prerequisite data
	schemaId := csKeeper.CreateMockCredentialSchema(1) // Create mock credential schema

	// Create mock permissions
	issuerGrantorPermId := csPermKeeper.CreateMockPermission(
		creator,
		schemaId,
		csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER_GRANTOR,
		"did:example:123",
		"US",
	)

	verifierGrantorPermId := csPermKeeper.CreateMockPermission(
		creator,
		schemaId,
		csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER_GRANTOR,
		"did:example:123",
		"US",
	)

	testCases := []struct {
		name    string
		msg     *types.MsgCreateValidation
		expPass bool
	}{
		{
			name: "Valid ISSUER Validation Request",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  types.ValidationType_ISSUER,
				ValidatorPermId: issuerGrantorPermId,
				Country:         "US",
			},
			expPass: true,
		},
		{
			name: "Valid VERIFIER Validation Request",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  types.ValidationType_VERIFIER,
				ValidatorPermId: verifierGrantorPermId,
				Country:         "US",
			},
			expPass: true,
		},
		{
			name: "Invalid Country Code",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  types.ValidationType_ISSUER,
				ValidatorPermId: issuerGrantorPermId,
				Country:         "USA", // Invalid format
			},
			expPass: false,
		},
		{
			name: "Non-existent Validator Permission",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  types.ValidationType_ISSUER,
				ValidatorPermId: 99999,
				Country:         "US",
			},
			expPass: false,
		},
		{
			name: "Country Mismatch",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  types.ValidationType_ISSUER,
				ValidatorPermId: issuerGrantorPermId,
				Country:         "GB", // Different from validator's country
			},
			expPass: false,
		},
		{
			name: "Missing Country",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  types.ValidationType_ISSUER,
				ValidatorPermId: issuerGrantorPermId,
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			perm, err := csPermKeeper.GetCSPermission(sdk.UnwrapSDKContext(ctx), tc.msg.ValidatorPermId)
			if err == nil && perm != nil {
				cs, _ := csKeeper.GetCredentialSchemaById(sdk.UnwrapSDKContext(ctx), perm.SchemaId)
				t.Logf("Permission Type: %v", perm.CspType)
				t.Logf("Schema Management Mode: %v", cs.IssuerPermManagementMode)
			}

			//cs, _ := csKeeper.GetCredentialSchemaById(sdk.UnwrapSDKContext(ctx), perm.SchemaId)

			resp, err := ms.CreateValidation(sdk.UnwrapSDKContext(ctx), tc.msg)
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify validation was created
				validation, err := k.Validation.Get(sdk.UnwrapSDKContext(ctx), resp.ValidationId)
				require.NoError(t, err)
				require.Equal(t, tc.msg.Creator, validation.Applicant)
				require.Equal(t, tc.msg.ValidationType, validation.Type)
				require.Equal(t, tc.msg.ValidatorPermId, validation.ValidatorPermId)
				require.Equal(t, types.ValidationState_PENDING, validation.State)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}

func TestValidationTypePermissionMatching(t *testing.T) {
	_, ms, csPermKeeper, csKeeper, ctx := setupMsgServer(t)
	creator := "verana1creator"

	schemaId := csKeeper.CreateMockCredentialSchema(1)

	// Create permissions of different types
	trustRegistryPermId := csPermKeeper.CreateMockPermission(
		creator,
		schemaId,
		csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY,
		"did:example:123",
		"US",
	)

	issuerPermId := csPermKeeper.CreateMockPermission(
		creator,
		schemaId,
		csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER,
		"did:example:123",
		"US",
	)

	testCases := []struct {
		name    string
		msg     *types.MsgCreateValidation
		expPass bool
	}{
		{
			name: "ISSUER_GRANTOR with Trust Registry Permission",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  types.ValidationType_ISSUER_GRANTOR,
				ValidatorPermId: trustRegistryPermId,
				Country:         "US",
			},
			expPass: true,
		},
		{
			name: "HOLDER with Issuer Permission",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  types.ValidationType_HOLDER,
				ValidatorPermId: issuerPermId,
				Country:         "US",
			},
			expPass: true,
		},
		{
			name: "Invalid: ISSUER with Issuer Permission",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  types.ValidationType_ISSUER,
				ValidatorPermId: issuerPermId,
				Country:         "US",
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.CreateValidation(ctx, tc.msg)
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
