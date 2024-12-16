package keeper

import (
	credentialschematypes "github.com/verana-labs/verana-blockchain/x/credentialschema/types"
	"testing"

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

	"github.com/verana-labs/verana-blockchain/x/cspermission/keeper"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
)

type MockCredentialSchemaKeeper struct {
	credentialschemas map[uint64]credentialschematypes.CredentialSchema
}

func NewMockCredentialSchemaKeeper() *MockCredentialSchemaKeeper {
	return &MockCredentialSchemaKeeper{
		credentialschemas: make(map[uint64]credentialschematypes.CredentialSchema),
	}
}

func (m *MockCredentialSchemaKeeper) CreateMockCredentialSchema(trId uint64) uint64 {
	// Generate next ID based on map length
	id := uint64(len(m.credentialschemas) + 1)

	// Create mock credential schema
	m.credentialschemas[id] = credentialschematypes.CredentialSchema{
		Id:                         id,
		TrId:                       trId,
		IssuerPermManagementMode:   credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION,
		VerifierPermManagementMode: credentialschematypes.CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION,
	}

	return id
}

func (m *MockCredentialSchemaKeeper) GetCredentialSchemaById(ctx sdk.Context, id uint64) (credentialschematypes.CredentialSchema, error) {
	return m.GetCredentialSchema(ctx, id)
}

func (m *MockCredentialSchemaKeeper) GetCredentialSchema(ctx sdk.Context, id uint64) (credentialschematypes.CredentialSchema, error) {
	if schema, ok := m.credentialschemas[id]; ok {
		return schema, nil
	}
	return credentialschematypes.CredentialSchema{}, credentialschematypes.ErrCredentialSchemaNotFound
}

func CspermissionKeeper(t testing.TB) (keeper.Keeper, *MockTrustRegistryKeeper, *MockCredentialSchemaKeeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	trustRegistryKeeper := NewMockTrustRegistryKeeper()
	credentialSchemaKeeper := NewMockCredentialSchemaKeeper()

	k := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		trustRegistryKeeper,
		credentialSchemaKeeper,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, trustRegistryKeeper, credentialSchemaKeeper, ctx
}
