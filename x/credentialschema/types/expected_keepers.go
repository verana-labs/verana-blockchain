package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	trustregistrytypes "github.com/verana-labs/verana-blockchain/x/trustregistry/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// TrustRegistryKeeper defines the expected trust registry keeper
type TrustRegistryKeeper interface {
	GetTrustRegistry(ctx sdk.Context, id uint64) (trustregistrytypes.TrustRegistry, error)
}
