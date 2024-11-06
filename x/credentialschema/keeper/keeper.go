package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/verana-labs/verana-blockchain/x/credentialschema/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority           string
		bankKeeper          types.BankKeeper
		trustregistryKeeper types.TrustRegistryKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	bankKeeper types.BankKeeper,
	trustregistryKeeper types.TrustRegistryKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:                 cdc,
		storeService:        storeService,
		authority:           authority,
		logger:              logger,
		bankKeeper:          bankKeeper,
		trustregistryKeeper: trustregistryKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetCredentialSchema returns a credential schema by ID
func (k Keeper) GetCredentialSchema(ctx sdk.Context, id uint64) (types.CredentialSchema, error) {
	kvStore := k.storeService.OpenKVStore(ctx)

	var credentialSchema types.CredentialSchema
	bz, err := kvStore.Get(types.GetCredentialSchemaKey(id))
	if err != nil {
		return credentialSchema, err
	}
	if bz == nil {
		return credentialSchema, types.ErrCredentialSchemaNotFound
	}

	k.cdc.MustUnmarshal(bz, &credentialSchema)
	return credentialSchema, nil
}

// SetCredentialSchema sets a credential schema
func (k Keeper) SetCredentialSchema(ctx sdk.Context, credentialSchema types.CredentialSchema) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&credentialSchema)

	return kvStore.Set(types.GetCredentialSchemaKey(credentialSchema.Id), bz)
}

// DeleteCredentialSchema deletes a credential schema
func (k Keeper) DeleteCredentialSchema(ctx sdk.Context, id uint64) error {
	kvStore := k.storeService.OpenKVStore(ctx)
	return kvStore.Delete(types.GetCredentialSchemaKey(id))
}

// IterateCredentialSchemas iterates over all credential schemas
func (k Keeper) IterateCredentialSchemas(ctx sdk.Context, fn func(schema types.CredentialSchema) (stop bool)) {
	kvStore := k.storeService.OpenKVStore(ctx)
	iterator, err := kvStore.Iterator(types.CredentialSchemaKeyPrefix, nil)
	if err != nil {
		panic(err)
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var schema types.CredentialSchema
		k.cdc.MustUnmarshal(iterator.Value(), &schema)
		if stop := fn(schema); stop {
			break
		}
	}
}
