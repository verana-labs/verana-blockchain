package upgrades

import (
	"github.com/verana-labs/verana-blockchain/app/upgrades/types"
	v2 "github.com/verana-labs/verana-blockchain/app/upgrades/v2"
)

var Upgrades = []types.Upgrade{
	v2.Upgrade,
}
