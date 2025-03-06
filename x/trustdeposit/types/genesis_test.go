package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/verana-labs/verana-blockchain/x/trustdeposit/types"
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
					TrustDepositReclaimBurnRate: uint32(60), // 60%
					TrustDepositShareValue:      uint64(1),  // Initial value: 1
					TrustDepositRate:            uint32(20), // 20%
					WalletUserAgentRewardRate:   uint32(20), // 20%
					UserAgentRewardRate:         uint32(20), // 20%
				},
				// Add other genesis state fields as needed
			},
			valid: true,
		},
		{
			desc: "invalid trust deposit reclaim burn rate",
			genState: &types.GenesisState{
				Params: types.Params{
					TrustDepositReclaimBurnRate: uint32(101), // Invalid: > 100%
					TrustDepositShareValue:      uint64(1),
					TrustDepositRate:            uint32(20),
					WalletUserAgentRewardRate:   uint32(20),
					UserAgentRewardRate:         uint32(20),
				},
			},
			valid: false,
		},
		{
			desc: "invalid trust deposit share value",
			genState: &types.GenesisState{
				Params: types.Params{
					TrustDepositReclaimBurnRate: uint32(60),
					TrustDepositShareValue:      uint64(0), // Invalid: cannot be 0
					TrustDepositRate:            uint32(20),
					WalletUserAgentRewardRate:   uint32(20),
					UserAgentRewardRate:         uint32(20),
				},
			},
			valid: false,
		},
		{
			desc: "invalid trust deposit rate",
			genState: &types.GenesisState{
				Params: types.Params{
					TrustDepositReclaimBurnRate: uint32(60),
					TrustDepositShareValue:      uint64(1),
					TrustDepositRate:            uint32(101), // Invalid: > 100%
					WalletUserAgentRewardRate:   uint32(20),
					UserAgentRewardRate:         uint32(20),
				},
			},
			valid: false,
		},
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
