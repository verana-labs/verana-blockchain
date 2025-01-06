package validation

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/verana-labs/verana-blockchain/x/validation/keeper"
	"github.com/verana-labs/verana-blockchain/x/validation/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	// Initialize all validations
	for _, validation := range genState.Validations {
		if err := k.Validation.Set(ctx, validation.Id, validation); err != nil {
			panic(fmt.Errorf("failed to set validation %d: %w", validation.Id, err))
		}
	}

	// Set the next validation ID
	if err := k.Counter.Set(ctx, "validation", genState.NextValidationId); err != nil {
		panic(fmt.Errorf("failed to set validation sequence: %w", err))
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	validations := make([]types.Validation, 0)

	// Get all validations
	k.Validation.Walk(ctx, nil, func(id uint64, validation types.Validation) (stop bool, err error) {
		validations = append(validations, validation)
		return false, nil
	})

	// Get next validation ID
	nextId, err := k.Counter.Get(ctx, "validation")
	if err != nil {
		panic(fmt.Errorf("failed to get validation sequence: %w", err))
	}

	return &types.GenesisState{
		Params:           k.GetParams(ctx),
		Validations:      validations,
		NextValidationId: nextId,
	}
}
