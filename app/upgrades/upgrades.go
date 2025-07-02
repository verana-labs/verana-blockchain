package upgrades

import (
	"github.com/verana-labs/verana-blockchain/app/upgrades/types"
	v4 "github.com/verana-labs/verana-blockchain/app/upgrades/v4"
)

var Upgrades = []types.Upgrade{
	v4.Upgrade,
}
