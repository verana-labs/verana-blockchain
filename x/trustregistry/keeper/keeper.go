package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
		// state management
		Schema        collections.Schema
		Params        collections.Item[types.Params]
		TrustRegistry collections.Map[string, types.TrustRegistry]
		GFVersion     collections.Map[string, types.GovernanceFrameworkVersion]
		GFDocument    collections.Map[string, types.GovernanceFrameworkDocument]
		// module references
		//bankKeeper trustregistry.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc: cdc,
		//addressCodec:  addressCodec,
		storeService:  storeService,
		authority:     authority,
		logger:        logger,
		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		TrustRegistry: collections.NewMap(sb, types.TrustRegistryKey, "trust_registry", collections.StringKey, codec.CollValue[types.TrustRegistry](cdc)),
		GFVersion:     collections.NewMap(sb, types.GovernanceFrameworkVersionKey, "gf_version", collections.StringKey, codec.CollValue[types.GovernanceFrameworkVersion](cdc)),
		GFDocument:    collections.NewMap(sb, types.GovernanceFrameworkDocumentKey, "gf_document", collections.StringKey, codec.CollValue[types.GovernanceFrameworkDocument](cdc)),
		//bankKeeper:    bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k

}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
