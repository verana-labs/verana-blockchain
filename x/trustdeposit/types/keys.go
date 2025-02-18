package types

const (
	// ModuleName defines the module name
	ModuleName = "trustdeposit"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_trustdeposit"
)

var (
	ParamsKey = []byte("p_trustdeposit")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
