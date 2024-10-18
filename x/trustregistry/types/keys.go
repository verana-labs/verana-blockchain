package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "trustregistry"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_trustregistry"
)

var (
	ParamsKey                      = collections.NewPrefix(0)
	TrustRegistryKey               = collections.NewPrefix(1)
	GovernanceFrameworkVersionKey  = collections.NewPrefix(2)
	GovernanceFrameworkDocumentKey = collections.NewPrefix(3)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
