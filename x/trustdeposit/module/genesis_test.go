package trustdeposit_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/testutil/nullify"
	trustdeposit "github.com/verana-labs/verana-blockchain/x/trustdeposit/module"
	"github.com/verana-labs/verana-blockchain/x/trustdeposit/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TrustdepositKeeper(t)
	trustdeposit.InitGenesis(ctx, k, genesisState)
	got := trustdeposit.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
