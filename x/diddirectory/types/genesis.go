package types

import "fmt"

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:         DefaultParams(),
		DidDirectories: []DIDDirectory{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate did directories
	seenDIDDirectories := make(map[string]bool)
	for _, tr := range gs.DidDirectories {
		// Check for duplicate DIDs
		if seenDIDDirectories[tr.Did] {
			return fmt.Errorf("duplicate DID Directory found in genesis state: %s", tr.Did)
		}
		seenDIDDirectories[tr.Did] = true
	}
	return nil
}
