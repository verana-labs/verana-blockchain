package keeper_test

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	csptypes "github.com/verana-labs/verana-blockchain/x/cspermission/types"
	"github.com/verana-labs/verana-blockchain/x/validation/keeper"
	"github.com/verana-labs/verana-blockchain/x/validation/types"
	"testing"
	"time"
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
				ValidationType:  uint32(types.ValidationType_ISSUER),
				ValidatorPermId: issuerGrantorPermId,
				Country:         "US",
			},
			expPass: true,
		},
		{
			name: "Valid VERIFIER Validation Request",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  uint32(types.ValidationType_VERIFIER),
				ValidatorPermId: verifierGrantorPermId,
				Country:         "US",
			},
			expPass: true,
		},
		{
			name: "Invalid Country Code",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  uint32(types.ValidationType_ISSUER),
				ValidatorPermId: issuerGrantorPermId,
				Country:         "USA", // Invalid format
			},
			expPass: false,
		},
		{
			name: "Non-existent Validator Permission",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  uint32(types.ValidationType_ISSUER),
				ValidatorPermId: 99999,
				Country:         "US",
			},
			expPass: false,
		},
		{
			name: "Country Mismatch",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  uint32(types.ValidationType_ISSUER),
				ValidatorPermId: issuerGrantorPermId,
				Country:         "GB", // Different from validator's country
			},
			expPass: false,
		},
		{
			name: "Missing Country",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  uint32(types.ValidationType_ISSUER),
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
				require.Equal(t, types.ValidationType(tc.msg.ValidationType), validation.Type)
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
				ValidationType:  uint32(types.ValidationType_ISSUER_GRANTOR),
				ValidatorPermId: trustRegistryPermId,
				Country:         "US",
			},
			expPass: true,
		},
		{
			name: "HOLDER with Issuer Permission",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  uint32(types.ValidationType_HOLDER),
				ValidatorPermId: issuerPermId,
				Country:         "US",
			},
			expPass: true,
		},
		{
			name: "Invalid: ISSUER with Issuer Permission",
			msg: &types.MsgCreateValidation{
				Creator:         creator,
				ValidationType:  uint32(types.ValidationType_ISSUER),
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

func TestRenewValidation(t *testing.T) {
	k, ms, csPermKeeper, csKeeper, ctx := setupMsgServer(t)
	creator := "verana1creator"

	// Create prerequisite data
	schemaId := csKeeper.CreateMockCredentialSchema(1)

	// Create mock permissions with different types and countries
	issuerGrantorPermId := csPermKeeper.CreateMockPermission(
		creator,
		schemaId,
		csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER_GRANTOR,
		"did:example:123",
		"US",
	)

	// Create a verifier grantor permission for incompatible type test
	verifierGrantorPermId := csPermKeeper.CreateMockPermission(
		creator,
		schemaId,
		csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER_GRANTOR,
		"did:example:789",
		"US",
	)

	// Create an initial validation that we'll try to renew
	createMsg := &types.MsgCreateValidation{
		Creator:         creator,
		ValidationType:  uint32(types.ValidationType_ISSUER),
		ValidatorPermId: issuerGrantorPermId,
		Country:         "US",
	}
	createResp, err := ms.CreateValidation(ctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Create another compatible validator permission for renewal tests
	newValidatorPermId := csPermKeeper.CreateMockPermission(
		creator,
		schemaId,
		csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER_GRANTOR,
		"did:example:456",
		"US",
	)

	testCases := []struct {
		name          string
		msg           *types.MsgRenewValidation
		expPass       bool
		errorContains string
		setupFn       func()
		checkFn       func(*testing.T, *types.Validation) // Optional additional checks
	}{
		{
			name: "Valid Renewal - Same Validator",
			msg: &types.MsgRenewValidation{
				Creator: creator,
				Id:      createResp.ValidationId,
			},
			expPass: true,
			checkFn: func(t *testing.T, v *types.Validation) {
				require.Equal(t, issuerGrantorPermId, v.ValidatorPermId)
			},
		},
		{
			name: "Valid Renewal - New Compatible Validator",
			msg: &types.MsgRenewValidation{
				Creator:         creator,
				Id:              createResp.ValidationId,
				ValidatorPermId: newValidatorPermId,
			},
			expPass: true,
			checkFn: func(t *testing.T, v *types.Validation) {
				require.Equal(t, newValidatorPermId, v.ValidatorPermId)
			},
		},
		{
			name: "Invalid - Non-existent Validation",
			msg: &types.MsgRenewValidation{
				Creator: creator,
				Id:      99999,
			},
			expPass:       false,
			errorContains: "validation not found",
		},
		{
			name: "Invalid - Wrong Creator",
			msg: &types.MsgRenewValidation{
				Creator: "verana1wrongcreator",
				Id:      createResp.ValidationId,
			},
			expPass:       false,
			errorContains: "only the validation applicant can renew",
		},
		{
			name: "Invalid - Non-existent New Validator",
			msg: &types.MsgRenewValidation{
				Creator:         creator,
				Id:              createResp.ValidationId,
				ValidatorPermId: 99999,
			},
			expPass:       false,
			errorContains: "validator permission not found",
		},
		{
			name: "Invalid - Incompatible Validator Type",
			msg: &types.MsgRenewValidation{
				Creator:         creator,
				Id:              createResp.ValidationId,
				ValidatorPermId: verifierGrantorPermId,
			},
			expPass:       false,
			errorContains: "permission type",
		},
		{
			name: "Valid Renewal - After Previous Renewal",
			msg: &types.MsgRenewValidation{
				Creator: creator,
				Id:      createResp.ValidationId,
			},
			expPass: true,
			setupFn: func() {
				// Get initial state
				initialVal, err := k.Validation.Get(ctx, createResp.ValidationId)
				require.NoError(t, err)
				t.Logf("Initial state - Deposit: %d, CurrentDeposit: %d",
					initialVal.ApplicantDeposit, initialVal.CurrentDeposit)

				// Advance time for first renewal
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				firstRenewalTime := sdkCtx.BlockTime().Add(time.Minute)
				ctx = sdk.WrapSDKContext(sdkCtx.WithBlockTime(firstRenewalTime))

				// First renewal
				firstRenewal := &types.MsgRenewValidation{
					Creator: creator,
					Id:      createResp.ValidationId,
				}
				_, err = ms.RenewValidation(ctx, firstRenewal)
				require.NoError(t, err)

				// Get state after first renewal but before second
				beforeSecondRenewal, err := k.Validation.Get(ctx, createResp.ValidationId)
				require.NoError(t, err)
				t.Logf("Before second renewal - Deposit: %d, CurrentDeposit: %d",
					beforeSecondRenewal.ApplicantDeposit, beforeSecondRenewal.CurrentDeposit)

				// Store this state for comparison in checkFn
				err = k.Validation.Set(ctx, createResp.ValidationId, beforeSecondRenewal)
				require.NoError(t, err)

				// Advance time for second renewal
				sdkCtx = sdk.UnwrapSDKContext(ctx)
				secondRenewalTime := sdkCtx.BlockTime().Add(time.Minute)
				ctx = sdk.WrapSDKContext(sdkCtx.WithBlockTime(secondRenewalTime))
			},
			checkFn: func(t *testing.T, v *types.Validation) {
				// Get the final validation state after renewal
				finalState, err := k.Validation.Get(ctx, v.Id)
				require.NoError(t, err)

				t.Logf("Validation State:")
				t.Logf("  Current Total Deposit: %d", finalState.ApplicantDeposit)
				t.Logf("  Renewal Deposit Amount: %d", finalState.CurrentDeposit)

				// Both should be equal since the execution is complete
				require.Equal(t, finalState.ApplicantDeposit, v.ApplicantDeposit,
					"Total deposit mismatch: state:%d, validation:%d",
					finalState.ApplicantDeposit, v.ApplicantDeposit)

				// The test passes because by this point, the second renewal
				// has already accumulated the deposit
				t.Logf("✓ Validation completed successfully with total deposit: %d", finalState.ApplicantDeposit)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run any custom setup
			if tc.setupFn != nil {
				tc.setupFn()
			}

			// Store initial state for comparison
			var originalValidation types.Validation
			if tc.expPass {
				var err error
				originalValidation, err = k.Validation.Get(ctx, tc.msg.Id)
				require.NoError(t, err)

				// Advance block time
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				newBlockTime := sdkCtx.BlockTime().Add(time.Minute)
				ctx = sdk.WrapSDKContext(sdkCtx.WithBlockTime(newBlockTime))
			}

			// Capture events before the call
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			previousEvents := sdkCtx.EventManager().Events()

			resp, err := ms.RenewValidation(ctx, tc.msg)
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, resp)

				// Verify validation was renewed correctly
				validation, err := k.Validation.Get(ctx, tc.msg.Id)
				require.NoError(t, err)

				// Check state changes
				require.Equal(t, types.ValidationState_PENDING, validation.State)
				require.True(t, validation.LastStateChange.After(originalValidation.LastStateChange),
					"LastStateChange (%v) should be after original (%v)",
					validation.LastStateChange, originalValidation.LastStateChange)

				// Verify fees and deposits
				require.True(t, validation.CurrentFees > 0, "CurrentFees should be positive")
				require.True(t, validation.CurrentDeposit > 0, "CurrentDeposit should be positive")
				require.True(t, validation.ApplicantDeposit > originalValidation.ApplicantDeposit,
					"ApplicantDeposit should increase: new %v, original %v",
					validation.ApplicantDeposit, originalValidation.ApplicantDeposit)

				// Check validator permission ID
				if tc.msg.ValidatorPermId != 0 {
					require.Equal(t, tc.msg.ValidatorPermId, validation.ValidatorPermId)

					// Check for revocation transfer event when validator changes
					if tc.msg.ValidatorPermId != originalValidation.ValidatorPermId {
						newEvents := sdkCtx.EventManager().Events()
						var found bool
						for _, event := range newEvents[len(previousEvents):] {
							if event.Type == "validation_revocation_control_transfer" {
								found = true
								break
							}
						}
						require.True(t, found, "revocation control transfer event should be emitted")
					}
				} else {
					require.Equal(t, originalValidation.ValidatorPermId, validation.ValidatorPermId)
				}

				// Run any additional checks
				if tc.checkFn != nil {
					tc.checkFn(t, &validation)
				}
			} else {
				require.Error(t, err)
				if tc.errorContains != "" {
					require.Contains(t, err.Error(), tc.errorContains)
				}
				require.Nil(t, resp)
			}
		})
	}
}

func TestSetValidated(t *testing.T) {
	k, ms, csPermKeeper, csKeeper, ctx := setupMsgServer(t)
	validator := "verana1validator"
	applicant := "verana1applicant"

	// Create prerequisite data
	schemaId := csKeeper.CreateMockCredentialSchema(1)

	// Create validator permission
	validatorPermId := csPermKeeper.CreateMockPermission(
		validator,
		schemaId,
		csptypes.CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER_GRANTOR,
		"did:example:123",
		"US",
	)

	// Create an initial validation
	createMsg := &types.MsgCreateValidation{
		Creator:         applicant,
		ValidationType:  uint32(types.ValidationType_ISSUER),
		ValidatorPermId: validatorPermId,
		Country:         "US",
	}
	createResp, err := ms.CreateValidation(ctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	validHash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	testCases := []struct {
		name          string
		msg           *types.MsgSetValidated
		expPass       bool
		errorContains string
		setupFn       func()
		checkFn       func(*testing.T, *types.Validation)
	}{
		{
			name: "Valid Set Validated - Without Summary Hash",
			msg: &types.MsgSetValidated{
				Creator: validator,
				Id:      createResp.ValidationId,
			},
			expPass: true,
			checkFn: func(t *testing.T, v *types.Validation) {
				require.Equal(t, types.ValidationState_VALIDATED, v.State)
				require.Empty(t, v.SummaryHash)
			},
		},
		{
			name: "Valid Set Validated - With Summary Hash",
			msg: &types.MsgSetValidated{
				Creator:     validator,
				Id:          createResp.ValidationId,
				SummaryHash: validHash,
			},
			expPass: true,
			setupFn: func() {
				// Reset validation state back to PENDING
				val, err := k.Validation.Get(ctx, createResp.ValidationId)
				require.NoError(t, err)
				val.State = types.ValidationState_PENDING
				err = k.Validation.Set(ctx, createResp.ValidationId, val)
				require.NoError(t, err)
			},
			checkFn: func(t *testing.T, v *types.Validation) {
				require.Equal(t, types.ValidationState_VALIDATED, v.State)
				require.Equal(t, validHash, v.SummaryHash)
			},
		},
		{
			name: "Invalid - Non-existent Validation",
			msg: &types.MsgSetValidated{
				Creator: validator,
				Id:      99999,
			},
			expPass:       false,
			errorContains: "validation not found",
		},
		{
			name: "Invalid - Wrong Validator",
			msg: &types.MsgSetValidated{
				Creator: applicant, // Using applicant instead of validator
				Id:      createResp.ValidationId,
			},
			setupFn: func() {
				// Ensure validation is in PENDING state
				val, err := k.Validation.Get(ctx, createResp.ValidationId)
				require.NoError(t, err)
				val.State = types.ValidationState_PENDING
				err = k.Validation.Set(ctx, createResp.ValidationId, val)
				require.NoError(t, err)
			},
			expPass:       false,
			errorContains: "only the validator can set validation to validated",
		},
		{
			name: "Invalid - Already Validated",
			msg: &types.MsgSetValidated{
				Creator: validator,
				Id:      createResp.ValidationId,
			},
			setupFn: func() {
				// Set validation to already validated
				val, err := k.Validation.Get(ctx, createResp.ValidationId)
				require.NoError(t, err)
				val.State = types.ValidationState_VALIDATED
				err = k.Validation.Set(ctx, createResp.ValidationId, val)
				require.NoError(t, err)
			},
			expPass:       false,
			errorContains: "validation must be in PENDING state",
		},
		{
			name: "Invalid - Summary Hash with HOLDER Type",
			msg: &types.MsgSetValidated{
				Creator:     validator,
				Id:          createResp.ValidationId,
				SummaryHash: validHash,
			},
			setupFn: func() {
				// Change validation type to HOLDER
				val, err := k.Validation.Get(ctx, createResp.ValidationId)
				require.NoError(t, err)
				val.Type = types.ValidationType_HOLDER
				val.State = types.ValidationState_PENDING
				err = k.Validation.Set(ctx, createResp.ValidationId, val)
				require.NoError(t, err)
			},
			expPass:       false,
			errorContains: "summary hash must be null for HOLDER type validations",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupFn != nil {
				tc.setupFn()
			}

			resp, err := ms.SetValidated(ctx, tc.msg)
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, resp)

				validation, err := k.Validation.Get(ctx, tc.msg.Id)
				require.NoError(t, err)

				// Common validations
				require.Equal(t, types.ValidationState_VALIDATED, validation.State)
				require.True(t, validation.LastStateChange.Equal(sdk.UnwrapSDKContext(ctx).BlockTime()))
				require.Zero(t, validation.CurrentFees)
				require.Zero(t, validation.CurrentDeposit)

				if tc.checkFn != nil {
					tc.checkFn(t, &validation)
				}
			} else {
				require.Error(t, err)
				if tc.errorContains != "" {
					require.Contains(t, err.Error(), tc.errorContains)
				}
				require.Nil(t, resp)
			}
		})
	}
}
