package types

const (
	// ModuleName defines the module name
	ModuleName = "veranablockchain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_veranablockchain"
)

var (
	ParamsKey = []byte("p_veranablockchain")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
