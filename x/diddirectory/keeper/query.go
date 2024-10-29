package keeper

import (
	"github.com/verana-labs/verana-blockchain/x/diddirectory/types"
)

var _ types.QueryServer = Keeper{}
