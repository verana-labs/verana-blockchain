package keeper_test

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cstypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/x/permission/keeper"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, *keepertest.MockCredentialSchemaKeeper, context.Context) {
	k, csKeeper, ctx := keepertest.PermissionKeeper(t)
	return k, keeper.NewMsgServerImpl(k), csKeeper, ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, _, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}

// Test for StartPermissionVP
func TestStartPermissionVP(t *testing.T) {
	k, ms, csKeeper, ctx := setupMsgServer(t)
	creator := sdk.AccAddress([]byte("test_creator")).String()

	// Create mock credential schema
	csKeeper.CreateMockCredentialSchema(1,
		cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION,
		cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION)

	// Create validator permission
	now := time.Now()
	validatorPerm := types.Permission{
		SchemaId:   1,
		Type:       3, // ISSUER_GRANTOR
		Grantee:    creator,
		Created:    &now,
		CreatedBy:  creator,
		Extended:   &now,
		ExtendedBy: creator,
		Modified:   &now,
		Country:    "US",
		VpState:    types.ValidationState_VALIDATION_STATE_VALIDATED,
	}
	validatorPermID, err := k.CreatePermission(sdk.UnwrapSDKContext(ctx), validatorPerm)
	require.NoError(t, err)

	testCases := []struct {
		name string
		msg  *types.MsgStartPermissionVP
		err  string
	}{
		{
			name: "Valid Permission VP Request",
			msg: &types.MsgStartPermissionVP{
				Creator:         creator,
				Type:            1, // ISSUER
				ValidatorPermId: validatorPermID,
				Country:         "US",
				Did:             "did:example:123",
			},
			err: "",
		},
		{
			name: "Non-existent Validator Permission",
			msg: &types.MsgStartPermissionVP{
				Creator:         creator,
				Type:            1,
				ValidatorPermId: 999,
				Country:         "US",
			},
			err: "validator permission not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.StartPermissionVP(ctx, tc.msg)
			if tc.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Greater(t, resp.PermissionId, uint64(0))

				// Verify created permission
				perm, err := k.GetPermission(sdk.UnwrapSDKContext(ctx), resp.PermissionId)
				require.NoError(t, err)
				require.Equal(t, tc.msg.Type, uint32(perm.Type))
				require.Equal(t, tc.msg.Creator, perm.Grantee)
				require.Equal(t, tc.msg.Country, perm.Country)
			}
		})
	}
}

func TestRenewPermissionVP(t *testing.T) {
	k, ms, csKeeper, ctx := setupMsgServer(t)
	creator := sdk.AccAddress([]byte("test_creator")).String()

	// Create mock credential schema
	csKeeper.CreateMockCredentialSchema(1,
		cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION,
		cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION)

	// Create validator permission
	now := time.Now()
	validatorPerm := types.Permission{
		SchemaId:   1,
		Type:       3, // ISSUER_GRANTOR
		Grantee:    creator,
		Created:    &now,
		CreatedBy:  creator,
		Extended:   &now,
		ExtendedBy: creator,
		Modified:   &now,
		Country:    "US",
		VpState:    types.ValidationState_VALIDATION_STATE_VALIDATED,
	}
	validatorPermID, err := k.CreatePermission(sdk.UnwrapSDKContext(ctx), validatorPerm)
	require.NoError(t, err)

	// Create applicant permission
	applicantPerm := types.Permission{
		SchemaId:        1,
		Type:            1, // ISSUER
		Grantee:         creator,
		Created:         &now,
		CreatedBy:       creator,
		Extended:        &now,
		ExtendedBy:      creator,
		Modified:        &now,
		Country:         "US",
		ValidatorPermId: validatorPermID,
		VpState:         types.ValidationState_VALIDATION_STATE_VALIDATED,
	}
	applicantPermID, err := k.CreatePermission(sdk.UnwrapSDKContext(ctx), applicantPerm)
	require.NoError(t, err)

	testCases := []struct {
		name string
		msg  *types.MsgRenewPermissionVP
		err  string
	}{
		{
			name: "Non-existent Permission",
			msg: &types.MsgRenewPermissionVP{
				Creator: creator,
				Id:      999,
			},
			err: "permission not found",
		},
		{
			name: "Wrong Creator",
			msg: &types.MsgRenewPermissionVP{
				Creator: sdk.AccAddress([]byte("wrong_creator")).String(),
				Id:      applicantPermID,
			},
			err: "creator is not the permission grantee",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.RenewPermissionVP(ctx, tc.msg)
			if tc.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify updated permission
				perm, err := k.GetPermission(sdk.UnwrapSDKContext(ctx), tc.msg.Id)
				require.NoError(t, err)
				require.Equal(t, types.ValidationState_VALIDATION_STATE_PENDING, perm.VpState)
				require.NotNil(t, perm.VpLastStateChange)
			}
		})
	}
}

func TestSetPermissionVPToValidated(t *testing.T) {
	k, ms, csKeeper, ctx := setupMsgServer(t)
	creator := sdk.AccAddress([]byte("test_creator")).String()
	validatorAddr := sdk.AccAddress([]byte("test_validator")).String()

	// Create mock credential schema with validation periods
	csKeeper.CreateMockCredentialSchema(1,
		cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION,
		cstypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION)

	// Create validator permission
	now := time.Now()
	validatorPerm := types.Permission{
		SchemaId:   1,
		Type:       3, // ISSUER_GRANTOR
		Grantee:    validatorAddr,
		Created:    &now,
		CreatedBy:  validatorAddr,
		Extended:   &now,
		ExtendedBy: validatorAddr,
		Modified:   &now,
		Country:    "US",
		VpState:    types.ValidationState_VALIDATION_STATE_VALIDATED,
	}
	validatorPermID, err := k.CreatePermission(sdk.UnwrapSDKContext(ctx), validatorPerm)
	require.NoError(t, err)

	// Create applicant permission
	applicantPerm := types.Permission{
		SchemaId:        1,
		Type:            1, // ISSUER
		Grantee:         creator,
		Created:         &now,
		CreatedBy:       creator,
		Extended:        &now,
		ExtendedBy:      creator,
		Modified:        &now,
		Country:         "US",
		ValidatorPermId: validatorPermID,
		VpState:         types.ValidationState_VALIDATION_STATE_PENDING,
	}
	applicantPermID, err := k.CreatePermission(sdk.UnwrapSDKContext(ctx), applicantPerm)
	require.NoError(t, err)

	testCases := []struct {
		name string
		msg  *types.MsgSetPermissionVPToValidated
		err  string
	}{
		{
			name: "Invalid Permission ID",
			msg: &types.MsgSetPermissionVPToValidated{
				Creator: validatorAddr,
				Id:      999,
			},
			err: "permission not found",
		},
		{
			name: "Wrong Validator",
			msg: &types.MsgSetPermissionVPToValidated{
				Creator: creator,
				Id:      applicantPermID,
			},
			err: "creator is not the validator",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.SetPermissionVPToValidated(ctx, tc.msg)
			if tc.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify updated permission
				perm, err := k.GetPermission(sdk.UnwrapSDKContext(ctx), tc.msg.Id)
				require.NoError(t, err)
				require.Equal(t, types.ValidationState_VALIDATION_STATE_VALIDATED, perm.VpState)
				require.Equal(t, tc.msg.ValidationFees, perm.ValidationFees)
				require.Equal(t, tc.msg.Country, perm.Country)
				require.NotNil(t, perm.EffectiveFrom)
				require.Equal(t, tc.msg.EffectiveUntil, perm.EffectiveUntil)
			}
		})
	}
}
