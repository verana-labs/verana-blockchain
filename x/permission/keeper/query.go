package keeper

import (
	"github.com/verana-labs/verana-blockchain/x/permission/types"
)

var _ types.QueryServer = Keeper{}
