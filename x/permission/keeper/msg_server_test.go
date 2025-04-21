package keeper_test

import (
	"context"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cstypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"

	"github.com/stretchr/testify/require"

	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/x/permission/keeper"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
)

//func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, *keepertest.MockCredentialSchemaKeeper, context.Context) {
//	k, csKeeper, ctx := keepertest.PermissionKeeper(t)
//	return k, keeper.NewMsgServerImpl(k), csKeeper, ctx
//}

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, *keepertest.MockCredentialSchemaKeeper, *keepertest.MockTrustRegistryKeeper, context.Context) {
	k, csKeeper, trkKeeper, ctx := keepertest.PermissionKeeper(t)
	return k, keeper.NewMsgServerImpl(k), csKeeper, trkKeeper, ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, _, _, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}

// Test for StartPermissionVP
func TestStartPermissionVP(t *testing.T) {
	k, ms, csKeeper, trkKeeper, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	creator := sdk.AccAddress([]byte("test_creator")).String()
	validDid := "did:example:123456789abcdefghi"

	// First create a trust registry for our credential schema
	trID := trkKeeper.CreateMockTrustRegistry(creator, validDid)

	// Create mock credential schema with specific permission management modes
	csKeeper.UpdateMockCredentialSchema(1, trID,
		cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION,
		cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION)

	// Create validator permission (ISSUER_GRANTOR)
	now := time.Now()
	// This should be VALIDATED as it's a prerequisite
	validatorPerm := types.Permission{
		SchemaId:   1,
		Type:       types.PermissionType_PERMISSION_TYPE_ISSUER_GRANTOR,
		Grantee:    creator,
		Created:    &now,
		CreatedBy:  creator,
		Extended:   &now,
		ExtendedBy: creator,
		Modified:   &now,
		Country:    "US",
		VpState:    types.ValidationState_VALIDATION_STATE_VALIDATED, // validator must be validated
	}
	validatorPermID, err := k.CreatePermission(sdkCtx, validatorPerm)
	require.NoError(t, err)

	// Create another validator permission (VERIFIER_GRANTOR with different country)
	verifierGrantorPerm := types.Permission{
		SchemaId:   1,
		Type:       types.PermissionType_PERMISSION_TYPE_VERIFIER_GRANTOR,
		Grantee:    creator,
		Created:    &now,
		CreatedBy:  creator,
		Extended:   &now,
		ExtendedBy: creator,
		Modified:   &now,
		Country:    "FR", // Different country
		VpState:    types.ValidationState_VALIDATION_STATE_VALIDATED,
	}
	verifierGrantorPermID, err := k.CreatePermission(sdkCtx, verifierGrantorPerm)
	require.NoError(t, err)

	testCases := []struct {
		name string
		msg  *types.MsgStartPermissionVP
		err  string
	}{
		{
			name: "Valid ISSUER Permission Request",
			msg: &types.MsgStartPermissionVP{
				Creator:         creator,
				Type:            uint32(types.PermissionType_PERMISSION_TYPE_ISSUER),
				ValidatorPermId: validatorPermID,
				Country:         "US",
				Did:             validDid,
			},
			err: "",
		},
		{
			name: "Non-existent Validator Permission",
			msg: &types.MsgStartPermissionVP{
				Creator:         creator,
				Type:            uint32(types.PermissionType_PERMISSION_TYPE_ISSUER),
				ValidatorPermId: 999,
				Country:         "US",
				Did:             validDid,
			},
			err: "validator permission not found",
		},
		//{
		//	name: "Country Mismatch",
		//	msg: &types.MsgStartPermissionVP{
		//		Creator:         creator,
		//		Type:            uint32(types.PermissionType_PERMISSION_TYPE_ISSUER),
		//		ValidatorPermId: validatorPermID,
		//		Country:         "FR", // Different from validator's country
		//		Did:             validDid,
		//	},
		//	err: "permission validation failed: validator permission is not valid: permission country mismatch: permission has US, requested FR does not contain validator permission country mismatch",
		//},
		{
			name: "Invalid Permission Type Combination - ISSUER with wrong validator",
			msg: &types.MsgStartPermissionVP{
				Creator:         creator,
				Type:            uint32(types.PermissionType_PERMISSION_TYPE_ISSUER),
				ValidatorPermId: verifierGrantorPermID, // Wrong validator type
				Country:         "FR",
				Did:             validDid,
			},
			err: "issuer permission requires ISSUER_GRANTOR validator",
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
				perm, err := k.GetPermissionByID(sdkCtx, resp.PermissionId)
				require.NoError(t, err)
				require.Equal(t, tc.msg.Type, uint32(perm.Type))
				require.Equal(t, tc.msg.Creator, perm.Grantee)
				require.Equal(t, tc.msg.Country, perm.Country)
				require.Equal(t, tc.msg.ValidatorPermId, perm.ValidatorPermId)
				require.Equal(t, types.ValidationState_VALIDATION_STATE_PENDING, perm.VpState)
				require.NotNil(t, perm.Created)
				require.NotNil(t, perm.Modified)
				require.NotNil(t, perm.VpLastStateChange)
			}
		})
	}
}

func TestRenewPermissionVP(t *testing.T) {
	k, ms, csKeeper, _, ctx := setupMsgServer(t)
	creator := sdk.AccAddress([]byte("test_creator")).String()

	// Create mock credential schema
	csKeeper.CreateMockCredentialSchema(1,
		cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION,
		cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION)

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
				perm, err := k.GetPermissionByID(sdk.UnwrapSDKContext(ctx), tc.msg.Id)
				require.NoError(t, err)
				require.Equal(t, types.ValidationState_VALIDATION_STATE_PENDING, perm.VpState)
				require.NotNil(t, perm.VpLastStateChange)
			}
		})
	}
}

func TestSetPermissionVPToValidated(t *testing.T) {
	k, ms, csKeeper, _, ctx := setupMsgServer(t)
	creator := sdk.AccAddress([]byte("test_creator")).String()
	validatorAddr := sdk.AccAddress([]byte("test_validator")).String()

	// Create mock credential schema with validation periods
	csKeeper.CreateMockCredentialSchema(1,
		cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION,
		cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION)

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
				perm, err := k.GetPermissionByID(sdk.UnwrapSDKContext(ctx), tc.msg.Id)
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

func TestMsgServerCreateRootPermission(t *testing.T) {
	k, ms, mockCsKeeper, trkKeeper, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	creator := sdk.AccAddress([]byte("test_creator")).String()
	validDid := "did:example:123456789abcdefghi"

	// First create a trust registry and store its ID
	trID := trkKeeper.CreateMockTrustRegistry(creator, validDid)

	// Create mock credential schema with specific permission management modes and trust registry ID
	mockCsKeeper.UpdateMockCredentialSchema(1,
		trID, // Set the trust registry ID
		cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION,
		cstypes.CredentialSchemaPermManagementMode_GRANTOR_VALIDATION)

	now := time.Now()
	futureTime := now.Add(24 * time.Hour)

	testCases := []struct {
		name    string
		msg     *types.MsgCreateRootPermission
		isValid bool
	}{
		{
			name: "Valid Create Root Permission",
			msg: &types.MsgCreateRootPermission{
				Creator:          creator,
				SchemaId:         1,
				Did:              validDid,
				ValidationFees:   100,
				IssuanceFees:     50,
				VerificationFees: 25,
				Country:          "US",
				EffectiveFrom:    &now,
				EffectiveUntil:   &futureTime,
			},
			isValid: true,
		},
		{
			name: "Non-existent Schema ID",
			msg: &types.MsgCreateRootPermission{
				Creator:          creator,
				SchemaId:         999,
				Did:              validDid,
				ValidationFees:   100,
				IssuanceFees:     50,
				VerificationFees: 25,
			},
			isValid: false,
		},
		{
			name: "Wrong Creator (Not Trust Registry Controller)",
			msg: &types.MsgCreateRootPermission{
				Creator:          sdk.AccAddress([]byte("wrong_creator")).String(),
				SchemaId:         1,
				Did:              validDid,
				ValidationFees:   100,
				IssuanceFees:     50,
				VerificationFees: 25,
			},
			isValid: false,
		},
	}

	var expectedID uint64 = 1 // Track expected auto-generated ID

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := ms.CreateRootPermission(ctx, tc.msg)
			if tc.isValid {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify ID was auto-generated correctly
				require.Equal(t, expectedID, resp.Id)

				// Get the created permission
				perm, err := k.GetPermissionByID(sdkCtx, resp.Id)
				require.NoError(t, err)

				// Verify all fields are set correctly
				require.Equal(t, tc.msg.SchemaId, perm.SchemaId)
				require.Equal(t, tc.msg.Did, perm.Did)
				require.Equal(t, tc.msg.Creator, perm.Grantee)
				require.Equal(t, types.PermissionType_PERMISSION_TYPE_TRUST_REGISTRY, perm.Type)
				require.Equal(t, tc.msg.ValidationFees, perm.ValidationFees)
				require.Equal(t, tc.msg.IssuanceFees, perm.IssuanceFees)
				require.Equal(t, tc.msg.VerificationFees, perm.VerificationFees)
				require.Equal(t, tc.msg.Country, perm.Country)

				// Verify time fields if set
				if tc.msg.EffectiveFrom != nil {
					require.Equal(t, tc.msg.EffectiveFrom.Unix(), perm.EffectiveFrom.Unix())
				}
				if tc.msg.EffectiveUntil != nil {
					require.Equal(t, tc.msg.EffectiveUntil.Unix(), perm.EffectiveUntil.Unix())
				}

				// Verify auto-populated fields
				require.NotNil(t, perm.Created)
				require.NotNil(t, perm.Modified)
				require.Equal(t, tc.msg.Creator, perm.CreatedBy)

				expectedID++ // Increment expected ID for next valid creation
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
			}
		})
	}
}
