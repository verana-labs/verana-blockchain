package keeper

import (
	"fmt"

	"cosmossdk.io/collections"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/verana-labs/verana-blockchain/x/permission/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
		// state
		Permission        collections.Map[uint64, types.Permission]
		PermissionCounter collections.Item[uint64]
		PermissionSession collections.Map[string, types.PermissionSession]

		// external keeper
		credentialSchemaKeeper types.CredentialSchemaKeeper
		trustRegistryKeeper    types.TrustRegistryKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	credentialSchemaKeeper types.CredentialSchemaKeeper,
	trustRegistryKeeper types.TrustRegistryKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:                    cdc,
		storeService:           storeService,
		authority:              authority,
		logger:                 logger,
		Permission:             collections.NewMap(sb, types.PermissionKey, "permission", collections.Uint64Key, codec.CollValue[types.Permission](cdc)),
		PermissionCounter:      collections.NewItem(sb, types.PermissionCounterKey, "permission_counter", collections.Uint64Value),
		PermissionSession:      collections.NewMap(sb, types.PermissionSessionKey, "permission_session", collections.StringKey, codec.CollValue[types.PermissionSession](cdc)),
		credentialSchemaKeeper: credentialSchemaKeeper,
		trustRegistryKeeper:    trustRegistryKeeper,
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

func (k Keeper) GetPermissionByID(ctx sdk.Context, id uint64) (types.Permission, error) {
	return k.Permission.Get(ctx, id)
}

// CreatePermission creates a new permission and returns its ID
func (k Keeper) CreatePermission(ctx sdk.Context, perm types.Permission) (uint64, error) {
	id, err := k.getNextPermissionID(ctx)
	if err != nil {
		return 0, err
	}

	if err := k.Permission.Set(ctx, id, perm); err != nil {
		return 0, err
	}

	return id, nil
}

// getNextPermissionID gets the next available permission ID
func (k Keeper) getNextPermissionID(ctx sdk.Context) (uint64, error) {
	id, err := k.PermissionCounter.Get(ctx)
	if err != nil {
		id = 0
	}

	nextID := id + 1
	err = k.PermissionCounter.Set(ctx, nextID)
	if err != nil {
		return 0, fmt.Errorf("failed to set permission counter: %w", err)
	}

	return nextID, nil
}

func (k Keeper) UpdatePermission(ctx sdk.Context, perm types.Permission) error {
	return k.Permission.Set(ctx, perm.Id, perm)
}
