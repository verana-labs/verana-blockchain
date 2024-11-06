package keeper

import (
	"github.com/verana-labs/verana-blockchain/x/credentialschema/types"
)

var _ types.QueryServer = Keeper{}
