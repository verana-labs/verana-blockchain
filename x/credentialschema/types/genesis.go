package types

import "fmt"

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:            DefaultParams(),
		CredentialSchemas: []CredentialSchema{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate credential schemas
	seenCredentialSchemas := make(map[uint64]bool)
	for _, cs := range gs.CredentialSchemas {
		// Check for duplicate CSs
		if seenCredentialSchemas[cs.Id] {
			return fmt.Errorf("duplicate Credential Schema found in genesis state: %d", cs.Id)
		}
		seenCredentialSchemas[cs.Id] = true
	}
	return nil
}
