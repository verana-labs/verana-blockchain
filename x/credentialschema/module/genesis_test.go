package credentialschema_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/verana-labs/verana-blockchain/testutil/keeper"
	"github.com/verana-labs/verana-blockchain/testutil/nullify"
	credentialschema "github.com/verana-labs/verana-blockchain/x/credentialschema/module"
	"github.com/verana-labs/verana-blockchain/x/credentialschema/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, _, ctx := keepertest.CredentialschemaKeeper(t)
	credentialschema.InitGenesis(ctx, k, genesisState)
	got := credentialschema.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
