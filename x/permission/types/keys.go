package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "permission"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_permission"
)

var (
	ParamsKey            = []byte("p_permission")
	PermissionKey        = collections.NewPrefix(0)
	PermissionCounterKey = collections.NewPrefix(1)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
