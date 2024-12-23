package keeper

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	csptypes "github.com/verana-labs/verana-blockchain/x/cspermission/types"
	"github.com/verana-labs/verana-blockchain/x/validation/keeper"
	"github.com/verana-labs/verana-blockchain/x/validation/types"
)

// MockCsPermissionKeeper mocks the credential schema permission keeper
type MockCsPermissionKeeper struct {
	permissions map[uint64]csptypes.CredentialSchemaPerm
}

func NewMockCsPermissionKeeper() *MockCsPermissionKeeper {
	return &MockCsPermissionKeeper{
		permissions: make(map[uint64]csptypes.CredentialSchemaPerm),
	}
}

func (m *MockCsPermissionKeeper) GetCSPermission(ctx sdk.Context, id uint64) (*csptypes.CredentialSchemaPerm, error) {
	if perm, ok := m.permissions[id]; ok {
		return &perm, nil
	}
	return nil, csptypes.ErrPermNotFound
}

func (m *MockCsPermissionKeeper) CreateMockPermission(
	creator string,
	schemaId uint64,
	permType csptypes.CredentialSchemaPermType,
	did string,
	country string,
) uint64 {
	id := uint64(len(m.permissions) + 1)
	now := time.Now()
	effectiveFrom := now.Add(time.Hour)

	m.permissions[id] = csptypes.CredentialSchemaPerm{
		Id:               id,
		SchemaId:         schemaId,
		CspType:          permType,
		Did:              did,
		Country:          country,
		Grantee:          creator,
		Created:          now,
		EffectiveFrom:    effectiveFrom,
		ValidationFees:   100,
		IssuanceFees:     200,
		VerificationFees: 300,
	}
	return id
}

// MockCredentialSchemaKeeper as you already have

func ValidationKeeper(t testing.TB) (keeper.Keeper, *MockCsPermissionKeeper, *MockCredentialSchemaKeeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	csPermKeeper := NewMockCsPermissionKeeper()
	csKeeper := NewMockCredentialSchemaKeeper()

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		csPermKeeper,
		csKeeper,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, csPermKeeper, csKeeper, ctx
}
