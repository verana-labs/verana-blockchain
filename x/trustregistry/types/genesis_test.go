package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params: types.Params{
					TrustUnitPrice:            types.DefaultTrustUnitPrice,
					TrustRegistryTrustDeposit: types.DefaultTrustRegistryTrustDeposit,
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "invalid trust unit price",
			genState: &types.GenesisState{
				Params: types.Params{
					TrustUnitPrice:            0, // Invalid value
					TrustRegistryTrustDeposit: types.DefaultTrustRegistryTrustDeposit,
				},
			},
			valid: false,
		},
		{
			desc: "invalid trust registry deposit",
			genState: &types.GenesisState{
				Params: types.Params{
					TrustUnitPrice:            types.DefaultTrustUnitPrice,
					TrustRegistryTrustDeposit: 0, // Invalid value
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
