package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
		Schema    collections.Schema

		// State
		CredentialSchemaPerm        collections.Map[uint64, types.CredentialSchemaPerm]
		Counter                     collections.Map[string, uint64]
		CredentialSchemaPermSession collections.Map[string, types.CredentialSchemaPermSession]

		// External keepers
		trustRegistryKeeper    types.TrustRegistryKeeper
		credentialSchemaKeeper types.CredentialSchemaKeeper
		//validationKeeper       types.ValidationKeeper // TODO: After validation module
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	trustRegistryKeeper types.TrustRegistryKeeper,
	credentialSchemaKeeper types.CredentialSchemaKeeper,
	//validationKeeper types.ValidationKeeper,

) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}
	sb := collections.NewSchemaBuilder(storeService)
	return Keeper{
		cdc:                         cdc,
		storeService:                storeService,
		authority:                   authority,
		logger:                      logger,
		CredentialSchemaPerm:        collections.NewMap(sb, types.CredentialSchemaPermKey, "credential_schema_perm", collections.Uint64Key, codec.CollValue[types.CredentialSchemaPerm](cdc)),
		Counter:                     collections.NewMap(sb, types.CounterKey, "counter", collections.StringKey, collections.Uint64Value),
		CredentialSchemaPermSession: collections.NewMap(sb, types.CredentialSchemaPermSessionKey, "credential_schema_perm_session", collections.StringKey, codec.CollValue[types.CredentialSchemaPermSession](cdc)),
		trustRegistryKeeper:         trustRegistryKeeper,
		credentialSchemaKeeper:      credentialSchemaKeeper,
		//validationKeeper:       validationKeeper,
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

func (k Keeper) GetCSPSession(ctx sdk.Context, id string) (*types.CredentialSchemaPermSession, error) {
	csps, err := k.CredentialSchemaPermSession.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &csps, nil
}
