package keeper

import (
	"fmt"

	"cosmossdk.io/collections"

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
		trustRegistryKeeper types.TrustRegistryKeeper

		// State management
		Schema collections.Schema
		//Params           collections.Item[types.Params]
		CredentialSchema collections.Map[uint64, types.CredentialSchema]
		Counter          collections.Map[string, uint64]
		trustDeposit     types.TrustDepositKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	bankKeeper types.BankKeeper,
	trustRegistryKeeper types.TrustRegistryKeeper,
	trustDeposit types.TrustDepositKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		cdc:                 cdc,
		storeService:        storeService,
		logger:              logger,
		authority:           authority,
		bankKeeper:          bankKeeper,
		trustRegistryKeeper: trustRegistryKeeper,

		// Initialize collections
		CredentialSchema: collections.NewMap(
			sb,
			types.CredentialSchemaKey,
			"credential_schema",
			collections.Uint64Key,
			codec.CollValue[types.CredentialSchema](cdc),
		),
		Counter: collections.NewMap(
			sb,
			types.CounterKey,
			"counter",
			collections.StringKey,
			collections.Uint64Value,
		),
		trustDeposit: trustDeposit,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetCredentialSchemaById GetCredentialSchema returns a credential schema by ID
func (k Keeper) GetCredentialSchemaById(ctx sdk.Context, id uint64) (types.CredentialSchema, error) {
	return k.CredentialSchema.Get(ctx, id)
}

// SetCredentialSchema sets a credential schema
func (k Keeper) SetCredentialSchema(ctx sdk.Context, schema types.CredentialSchema) error {
	return k.CredentialSchema.Set(ctx, schema.Id, schema)
}

// DeleteCredentialSchema deletes a credential schema
func (k Keeper) DeleteCredentialSchema(ctx sdk.Context, id uint64) error {
	return k.CredentialSchema.Remove(ctx, id)
}

// IterateCredentialSchemas iterates over all credential schemas
func (k Keeper) IterateCredentialSchemas(ctx sdk.Context, fn func(schema types.CredentialSchema) (stop bool)) error {
	return k.CredentialSchema.Walk(ctx, nil, func(key uint64, value types.CredentialSchema) (bool, error) {
		return fn(value), nil
	})
}

func (k Keeper) GetNextID(ctx sdk.Context, entityType string) (uint64, error) {
	currentID, err := k.Counter.Get(ctx, entityType)
	if err != nil {
		currentID = 0
	}

	nextID := currentID + 1
	err = k.Counter.Set(ctx, entityType, nextID)
	if err != nil {
		return 0, fmt.Errorf("failed to set counter: %w", err)
	}

	return nextID, nil
}
