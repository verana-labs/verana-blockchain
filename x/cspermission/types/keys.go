package types

const (
	// ModuleName defines the module name
	ModuleName = "cspermission"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_cspermission"
)

var (
	ParamsKey = []byte("p_cspermission")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
