package keeper

import (
	"github.com/verana-labs/verana-blockchain/x/veranablockchain/types"
)

var _ types.QueryServer = Keeper{}
