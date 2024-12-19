package keeper

import (
	"cosmossdk.io/collections"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/verana-labs/verana-blockchain/x/validation/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
		Schema    collections.Schema

		// State
		Validation            collections.Map[uint64, types.Validation]
		Counter               collections.Map[string, uint64]
		ValidationByApplicant collections.Map[string, []uint64]
		ValidationByValidator collections.Map[uint64, []uint64]

		// External keepers
		csPermissionKeeper     types.CsPermissionKeeper
		credentialSchemaKeeper types.CredentialSchemaKeeper
		//trustDepositKeeper     types.TrustDepositKeeper // TODO: After TD module

	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	csPermissionKeeper types.CsPermissionKeeper,
	credentialSchemaKeeper types.CredentialSchemaKeeper,

) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}
	sb := collections.NewSchemaBuilder(storeService)

	return Keeper{
		cdc:                    cdc,
		storeService:           storeService,
		authority:              authority,
		logger:                 logger,
		Validation:             collections.NewMap(sb, types.ValidationKey, "validation", collections.Uint64Key, codec.CollValue[types.Validation](cdc)),
		Counter:                collections.NewMap(sb, types.CounterKey, "counter", collections.StringKey, collections.Uint64Value),
		csPermissionKeeper:     csPermissionKeeper,
		credentialSchemaKeeper: credentialSchemaKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetNextID(ctx sdk.Context) (uint64, error) {
	currentID, err := k.Counter.Get(ctx, "validation")
	if err != nil {
		currentID = 0
	}

	nextID := currentID + 1
	err = k.Counter.Set(ctx, "validation", nextID)
	if err != nil {
		return 0, fmt.Errorf("failed to set counter: %w", err)
	}

	return nextID, nil
}
