package cspermission

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/verana-labs/verana-blockchain/x/cspermission/keeper"
	"github.com/verana-labs/verana-blockchain/x/cspermission/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module params
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	// Initialize credential schema permissions
	for _, perm := range genState.CredentialSchemaPerms {
		if err := k.CredentialSchemaPerm.Set(ctx, perm.Id, perm); err != nil {
			panic(fmt.Sprintf("failed to initialize credential schema permission: %v", err))
		}
	}

	// Initialize credential schema permission sessions
	for _, session := range genState.CredentialSchemaPermSessions {
		if err := k.CredentialSchemaPermSession.Set(ctx, session.Id, session); err != nil {
			panic(fmt.Sprintf("failed to initialize credential schema permission session: %v", err))
		}
	}

	// Set the counter for next CSP ID
	if err := k.Counter.Set(ctx, "counter", genState.NextCspId); err != nil {
		panic(fmt.Sprintf("failed to set CSP counter: %v", err))
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export all credential schema permissions
	var perms []types.CredentialSchemaPerm
	if err := k.CredentialSchemaPerm.Walk(ctx, nil, func(key uint64, perm types.CredentialSchemaPerm) (bool, error) {
		perms = append(perms, perm)
		return false, nil
	}); err != nil {
		panic(fmt.Sprintf("failed to export credential schema permissions: %v", err))
	}
	genesis.CredentialSchemaPerms = perms

	// Export all credential schema permission sessions
	var sessions []types.CredentialSchemaPermSession
	if err := k.CredentialSchemaPermSession.Walk(ctx, nil, func(key string, session types.CredentialSchemaPermSession) (bool, error) {
		sessions = append(sessions, session)
		return false, nil
	}); err != nil {
		panic(fmt.Sprintf("failed to export credential schema permission sessions: %v", err))
	}
	genesis.CredentialSchemaPermSessions = sessions

	// Export the next CSP ID
	nextID, err := k.GetNextID(ctx, "csp")
	if err != nil {
		nextID = 1 // Default to 1 if counter doesn't exist
	}
	genesis.NextCspId = nextID

	return genesis
}
