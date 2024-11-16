package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "cspermission"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_cspermission"

	RouterKey = ModuleName
)

var (
	ParamsKey               = []byte("p_cspermission")
	CredentialSchemaPermKey = collections.NewPrefix(1)
	CounterKey              = collections.NewPrefix(2)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
