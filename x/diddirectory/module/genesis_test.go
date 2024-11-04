package diddirectory_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/testutil/nullify"
	diddirectory "github.com/verana-labs/verana-blockchain/x/diddirectory/module"
	"github.com/verana-labs/verana-blockchain/x/diddirectory/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.DiddirectoryKeeper(t)
	diddirectory.InitGenesis(ctx, k, genesisState)
	got := diddirectory.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
