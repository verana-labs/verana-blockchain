package keeper

import (
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
)

var _ types.QueryServer = Keeper{}
