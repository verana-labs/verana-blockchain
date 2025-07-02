package v4

import (
	store "cosmossdk.io/store/types"
	"github.com/verana-labs/verana-blockchain/app/upgrades/types"
)

const UpgradeName = "v0.4"

var Upgrade = types.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
