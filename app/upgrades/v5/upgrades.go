package v5

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/verana-labs/verana-blockchain/app/upgrades/types"
)

// migrateStore copies all KV pairs from the old store key to the new one
func migrateStore(ctx sdk.Context, app interface {
	GetKey(string) *storetypes.KVStoreKey
}, oldKey, newKey string) {
	oldStore := ctx.KVStore(app.GetKey(oldKey))
	newStore := ctx.KVStore(app.GetKey(newKey))

	iterator := oldStore.Iterator(nil, nil) // iterate all keys
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()
		newStore.Set(key, value)
	}
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ types.BaseAppParamManager,
	keepers types.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(context context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(context)

		// Migrate all renamed module stores
		if app, ok := keepers.(interface {
			GetKey(string) *storetypes.KVStoreKey
		}); ok {
			migrateStore(ctx, app, "credentialschema", "cs")
			migrateStore(ctx, app, "diddirectory", "dd")
			migrateStore(ctx, app, "permission", "perm")
			migrateStore(ctx, app, "trustdeposit", "td")
			migrateStore(ctx, app, "trustregistry", "tr")
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
