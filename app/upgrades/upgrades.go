package upgrades

import (
	"github.com/verana-labs/verana-blockchain/app/upgrades/types"
	v3 "github.com/verana-labs/verana-blockchain/app/upgrades/v3"
)

var Upgrades = []types.Upgrade{
	v3.Upgrade,
}
