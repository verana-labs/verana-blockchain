package keeper

import (
	"github.com/verana-labs/verana-blockchain/x/validation/types"
)

var _ types.QueryServer = Keeper{}
