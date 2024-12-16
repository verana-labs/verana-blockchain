package credentialschema

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/verana-labs/verana-blockchain/x/credentialschema/keeper"
	"github.com/verana-labs/verana-blockchain/x/credentialschema/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
	// Initialize Credential Schemas
	for _, cs := range genState.CredentialSchemas {
		// Set credential schema
		if err := k.CredentialSchema.Set(ctx, cs.Id, cs); err != nil {
			panic(fmt.Sprintf("failed to set Credential Schema: %s", err))
		}
		//TODO: Add Incremental id
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	// Export all credential schema
	var credentialSchemas []types.CredentialSchema
	err := k.CredentialSchema.Walk(ctx, nil, func(key uint64, cs types.CredentialSchema) (bool, error) {
		credentialSchemas = append(credentialSchemas, cs)
		return false, nil
	})
	if err != nil {
		panic(fmt.Sprintf("failed to export Credential Schema: %s", err))
	}

	genesis.CredentialSchemas = credentialSchemas

	return genesis
}
