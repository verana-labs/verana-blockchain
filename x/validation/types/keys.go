package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "validation"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_validation"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	ParamsKey = []byte("p_validation")

	// ValidationKey stores validation entries by ID
	ValidationKey = collections.NewPrefix(1)

	// CounterKey for generating validation IDs
	CounterKey = collections.NewPrefix(2)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
