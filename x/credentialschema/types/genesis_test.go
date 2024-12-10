package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/verana-labs/verana-blockchain/x/credentialschema/types"
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
		//TODO: Fix by adding valid genesis state
		//{
		//	desc:     "valid genesis state",
		//	genState: &types.GenesisState{
		//
		//		// this line is used by starport scaffolding # types/genesis/validField
		//	},
		//	valid: true,
		//},
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
