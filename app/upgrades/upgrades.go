package upgrades

import (
	"github.com/verana-labs/verana-blockchain/app/upgrades/types"
	v6 "github.com/verana-labs/verana-blockchain/app/upgrades/v6"
)

var Upgrades = []types.Upgrade{
	v6.Upgrade,
}
