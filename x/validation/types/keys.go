package types

const (
	// ModuleName defines the module name
	ModuleName = "validation"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_validation"
)

var (
	ParamsKey = []byte("p_validation")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
