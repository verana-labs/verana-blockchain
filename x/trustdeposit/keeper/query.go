package keeper

import (
	"github.com/verana-labs/verana-blockchain/x/trustdeposit/types"
)

var _ types.QueryServer = Keeper{}
