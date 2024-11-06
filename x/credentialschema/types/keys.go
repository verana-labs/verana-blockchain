package types

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "credentialschema"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_credentialschema"

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	ParamsKey                 = []byte("p_credentialschema")
	CredentialSchemaKeyPrefix = []byte{0x01} // prefix for credential schema
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetCredentialSchemaKey returns the store key to retrieve a CredentialSchema from the index fields
func GetCredentialSchemaKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append(CredentialSchemaKeyPrefix, bz...)
}
