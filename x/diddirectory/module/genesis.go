package diddirectory

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/verana-labs/verana-blockchain/x/diddirectory/keeper"
	"github.com/verana-labs/verana-blockchain/x/diddirectory/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
	// Initialize did directories
	for _, dd := range genState.DidDirectories {
		// Set did directory
		if err := k.DIDDirectory.Set(ctx, dd.Did, dd); err != nil {
			panic(fmt.Sprintf("failed to set DID directory: %s", err))
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	// Export all did directories
	var didDirectories []types.DIDDirectory
	err := k.DIDDirectory.Walk(ctx, nil, func(key string, tr types.DIDDirectory) (bool, error) {
		didDirectories = append(didDirectories, tr)
		return false, nil
	})
	if err != nil {
		panic(fmt.Sprintf("failed to export DID directory: %s", err))
	}

	genesis.DidDirectories = didDirectories
	return genesis
}
