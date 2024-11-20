package cspermission

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/verana-labs/verana-blockchain/api/veranablockchain/cspermission"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListCSP",
					Use:       "list-csp [schema-id]",
					Short:     "List credential schema permissions",
					Long: `List credential schema permissions filtered by various parameters.
Mandatory:
  schema-id: ID of the credential schema

Optional flags:
  --creator: Filter by creator address
  --grantee: Filter by grantee address
  --did: Filter by grantee DID
  --type: Filter by permission type (1-6)
  --response-max-size: Maximum number of results (1-1024, default 64)`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "schema_id",
						},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"creator": {
							Name:         "creator",
							Usage:        "Filter by creator address",
							DefaultValue: "",
						},
						"grantee": {
							Name:         "grantee",
							Usage:        "Filter by grantee address",
							DefaultValue: "",
						},
						"did": {
							Name:         "did",
							Usage:        "Filter by grantee DID",
							DefaultValue: "",
						},
						"type": {
							Name:         "type",
							Usage:        "Filter by permission type",
							DefaultValue: "0",
						},
						"response_max_size": {
							Name:         "response-max-size",
							Usage:        "Maximum number of results",
							DefaultValue: "64",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "CreateCredentialSchemaPerm",
					Use:       "create-credential-schema-perm [schema-id] [csp-type] [did] [grantee] [effective-from] [validation-fees] [issuance-fees] [verification-fees]",
					Short:     "Create a new credential schema permission",
					Long: `Create a new credential schema permission with the specified parameters.

Parameters:
  [schema-id]         : ID of the credential schema (uint64)
  [csp-type]         : Permission type:
                       1 = ISSUER
                       2 = VERIFIER
                       3 = ISSUER_GRANTOR
                       4 = VERIFIER_GRANTOR
                       5 = TRUST_REGISTRY
                       6 = HOLDER
  [did]              : DID of the grantee service
  [grantee]          : Account address of the grantee
  [effective-from]   : Start date (RFC3339 format, e.g., 2024-03-16T15:00:00Z)
  [validation-fees]  : Fees for validation process
  [issuance-fees]   : Fees for credential issuance
  [verification-fees]: Fees for credential verification

Optional Flags:
  --effective-until  : End date (RFC3339 format, e.g., 2025-03-16T15:00:00Z)
  --country         : ISO 3166-1 alpha-2 country code
  --validation-id   : ID of the validation entry required for non-trust-registry controllers`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "schema_id",
						},
						{
							ProtoField: "csp_type",
						},
						{
							ProtoField: "did",
						},
						{
							ProtoField: "grantee",
						},
						{
							ProtoField: "effective_from",
						},
						{
							ProtoField: "validation_fees",
						},
						{
							ProtoField: "issuance_fees",
						},
						{
							ProtoField: "verification_fees",
						},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"effective_until": {
							Name:         "effective-until",
							DefaultValue: "",
							Usage:        "End date (RFC3339 format)",
						},
						"country": {
							Name:         "country",
							DefaultValue: "",
							Usage:        "ISO 3166-1 alpha-2 country code",
						},
						"validation_id": {
							Name:         "validation-id",
							DefaultValue: "0",
							Usage:        "ID of the validation entry required for non-trust-registry controllers",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
