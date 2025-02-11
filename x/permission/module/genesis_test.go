package permission_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/testutil/nullify"
	permission "github.com/verana-labs/verana-blockchain/x/permission/module"
	"github.com/verana-labs/verana-blockchain/x/permission/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, _, _, ctx := keepertest.PermissionKeeper(t)
	permission.InitGenesis(ctx, k, genesisState)
	got := permission.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
