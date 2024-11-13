package trustregistry

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/verana-labs/verana-blockchain/x/trustregistry/keeper"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
	// Initialize trust registries
	for _, tr := range genState.TrustRegistries {
		// Set trust registry
		if err := k.TrustRegistry.Set(ctx, tr.Id, tr); err != nil {
			panic(fmt.Sprintf("failed to set trust registry: %s", err))
		}

		// Set DID index
		if err := k.TrustRegistryDIDIndex.Set(ctx, tr.Did, tr.Id); err != nil {
			panic(fmt.Sprintf("failed to set DID index: %s", err))
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	// Export all trust registries
	var trustRegistries []types.TrustRegistry
	err := k.TrustRegistry.Walk(ctx, nil, func(key uint64, tr types.TrustRegistry) (bool, error) {
		trustRegistries = append(trustRegistries, tr)
		return false, nil
	})
	if err != nil {
		panic(fmt.Sprintf("failed to export trust registries: %s", err))
	}

	genesis.TrustRegistries = trustRegistries

	return genesis
}
