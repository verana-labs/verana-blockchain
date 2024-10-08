package veranablockchain_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/testutil/nullify"
	veranablockchain "github.com/verana-labs/verana-blockchain/x/veranablockchain/module"
	"github.com/verana-labs/verana-blockchain/x/veranablockchain/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.VeranablockchainKeeper(t)
	veranablockchain.InitGenesis(ctx, k, genesisState)
	got := veranablockchain.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
