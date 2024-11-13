package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:          DefaultParams(),
		TrustRegistries: []TrustRegistry{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate trust registries
	seenDIDs := make(map[uint64]bool)
	for _, tr := range gs.TrustRegistries {
		// Check for duplicate DIDs
		if seenDIDs[tr.Id] {
			return fmt.Errorf("duplicate ID found in genesis state: %s", tr.Did)
		}
		seenDIDs[tr.Id] = true
	}

	return nil
}
