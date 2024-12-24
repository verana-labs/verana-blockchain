package keeper_test

import (
	"context"
	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/x/validation/keeper"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	csptypes "github.com/verana-labs/verana-blockchain/x/cspermission/types"
	"github.com/verana-labs/verana-blockchain/x/validation/types"
)

func setupQueryServer(t testing.TB) (types.QueryServer, *keeper.Keeper, types.MsgServer, *keepertest.MockCsPermissionKeeper, *keepertest.MockCredentialSchemaKeeper, context.Context) {
	k, csPermKeeper, csKeeper, ctx := keepertest.ValidationKeeper(t)
	return keeper.NewQueryServerImpl(k), &k, keeper.NewMsgServerImpl(k), csPermKeeper, csKeeper, ctx
}

func TestListValidations(t *testing.T) {
	qs, k, ms, csPermKeeper, csKeeper, ctx := setupQueryServer(t)
	creator := sdk.AccAddress([]byte("test_creator")).String()

	// Create prerequisite data
	schemaId := csKeeper.CreateMockCredentialSchema(1)

	// Create different permission types
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

	// Helper function to get appropriate permission ID based on validation type
	getPermId := func(vType types.ValidationType) uint64 {
		switch vType {
		case types.ValidationType_ISSUER:
			return issuerGrantorPermId
		case types.ValidationType_VERIFIER:
			return verifierGrantorPermId
		default:
			return issuerGrantorPermId
		}
	}

	// Helper function to create validations
	createValidation := func(applicant string, vType types.ValidationType, state types.ValidationState, expTime *time.Time) uint64 {
		msg := &types.MsgCreateValidation{
			Creator:         applicant,
			ValidationType:  vType,
			ValidatorPermId: getPermId(vType),
			Country:         "US",
		}
		resp, err := ms.CreateValidation(sdk.UnwrapSDKContext(ctx), msg)
		require.NoError(t, err)

		if state != types.ValidationState_PENDING {
			validation, err := k.Validation.Get(sdk.UnwrapSDKContext(ctx), resp.ValidationId)
			require.NoError(t, err)
			validation.State = state
			validation.Exp = expTime
			err = k.Validation.Set(sdk.UnwrapSDKContext(ctx), resp.ValidationId, validation)
			require.NoError(t, err)
		}

		return resp.ValidationId
	}

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	yesterday := now.Add(-24 * time.Hour)

	// Create test validations
	_ = createValidation(creator, types.ValidationType_ISSUER, types.ValidationState_PENDING, &tomorrow)
	_ = createValidation(creator, types.ValidationType_VERIFIER, types.ValidationState_VALIDATED, &yesterday)
	_ = createValidation(sdk.AccAddress([]byte("other_creator")).String(), types.ValidationType_ISSUER, types.ValidationState_PENDING, nil)

	testCases := []struct {
		name          string
		request       *types.QueryListValidationsRequest
		expectedError bool
		check         func(*testing.T, *types.QueryListValidationsResponse)
	}{
		{
			name: "List All By Controller",
			request: &types.QueryListValidationsRequest{
				Controller:      creator,
				ResponseMaxSize: 10,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryListValidationsResponse) {
				require.Len(t, response.Validations, 2)
				for _, v := range response.Validations {
					require.Equal(t, creator, v.Applicant)
				}
			},
		},
		{
			name: "Filter By Type",
			request: &types.QueryListValidationsRequest{
				ValidatorPermId: issuerGrantorPermId,
				Type:            types.ValidationType_ISSUER,
				ResponseMaxSize: 10,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryListValidationsResponse) {
				require.Len(t, response.Validations, 2)
				for _, v := range response.Validations {
					require.Equal(t, types.ValidationType_ISSUER, v.Type)
				}
			},
		},
		{
			name: "Filter By State",
			request: &types.QueryListValidationsRequest{
				ValidatorPermId: issuerGrantorPermId,
				State:           types.ValidationState_PENDING,
				ResponseMaxSize: 10,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryListValidationsResponse) {
				require.Len(t, response.Validations, 2)
				for _, v := range response.Validations {
					require.Equal(t, types.ValidationState_PENDING, v.State)
				}
			},
		},
		{
			name: "Response Size Limit",
			request: &types.QueryListValidationsRequest{
				ValidatorPermId: issuerGrantorPermId,
				ResponseMaxSize: 1,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryListValidationsResponse) {
				require.Len(t, response.Validations, 1)
			},
		},
		{
			name: "Invalid Response Size",
			request: &types.QueryListValidationsRequest{
				ValidatorPermId: issuerGrantorPermId,
				ResponseMaxSize: 1025,
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := qs.ListValidations(sdk.UnwrapSDKContext(ctx), tc.request)
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

func TestGetValidation(t *testing.T) {
	qs, _, ms, csPermKeeper, csKeeper, ctx := setupQueryServer(t)
	creator := sdk.AccAddress([]byte("test_creator")).String()

	// Create prerequisite data
	schemaId := csKeeper.CreateMockCredentialSchema(1)

	issuerGrantorPermId := csPermKeeper.CreateMockPermission(
		creator,
		schemaId,
		csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER_GRANTOR,
		"did:example:123",
		"US",
	)

	// Create a test validation
	msg := &types.MsgCreateValidation{
		Creator:         creator,
		ValidationType:  types.ValidationType_ISSUER,
		ValidatorPermId: issuerGrantorPermId,
		Country:         "US",
	}
	createResp, err := ms.CreateValidation(sdk.UnwrapSDKContext(ctx), msg)
	require.NoError(t, err)
	validationId := createResp.ValidationId

	testCases := []struct {
		name          string
		request       *types.QueryGetValidationRequest
		expectedError bool
		check         func(*testing.T, *types.QueryGetValidationResponse)
	}{
		{
			name: "Valid Request",
			request: &types.QueryGetValidationRequest{
				Id: validationId,
			},
			expectedError: false,
			check: func(t *testing.T, response *types.QueryGetValidationResponse) {
				require.NotNil(t, response.Validation)
				require.Equal(t, validationId, response.Validation.Id)
				require.Equal(t, creator, response.Validation.Applicant)
				require.Equal(t, types.ValidationType_ISSUER, response.Validation.Type)
				require.Equal(t, issuerGrantorPermId, response.Validation.ValidatorPermId)
				require.Equal(t, types.ValidationState_PENDING, response.Validation.State)
			},
		},
		{
			name: "Zero ID",
			request: &types.QueryGetValidationRequest{
				Id: 0,
			},
			expectedError: true,
		},
		{
			name: "Non-existent ID",
			request: &types.QueryGetValidationRequest{
				Id: 99999,
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := qs.GetValidation(sdk.UnwrapSDKContext(ctx), tc.request)
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
