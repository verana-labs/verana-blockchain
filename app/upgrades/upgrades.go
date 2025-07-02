package upgrades

import (
	"github.com/verana-labs/verana-blockchain/app/upgrades/types"
	v5 "github.com/verana-labs/verana-blockchain/app/upgrades/v5"
)

var Upgrades = []types.Upgrade{
	v5.Upgrade,
}
