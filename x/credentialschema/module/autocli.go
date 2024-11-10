package credentialschema

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/verana-labs/verana-blockchain/api/veranablockchain/credentialschema"
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
					RpcMethod: "ListCredentialSchemas",
					Use:       "list-schemas",
					Short:     "List credential schemas with optional filters",
					Long: `List credential schemas with optional filters.
Example:
$ veranad query credentialschema list-schemas
$ veranad query credentialschema list-schemas --tr_id 1 --created_after 2024-01-01T00:00:00Z --response_max_size 100`,
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"tr_id": {
							Name:         "tr_id",
							Usage:        "Filter by trust registry ID",
							DefaultValue: "0",
						},
						"created_after": {
							Name:         "created_after",
							Usage:        "Show schemas created after this datetime (RFC3339 format)",
							DefaultValue: "",
						},
						"response_max_size": {
							Name:         "response_max_size",
							Usage:        "Maximum number of results (1-1024, default 64)",
							DefaultValue: "64",
						},
					},
				},
				{
					RpcMethod: "GetCredentialSchema",
					Use:       "get [id]",
					Short:     "Get a credential schema by ID",
					Long: `Get a credential schema by its ID.

Example:
$ veranad query credentialschema get 1`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
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
					RpcMethod: "CreateCredentialSchema",
					Use:       "create-credential-schema [id] [tr-id] [json-schema] [issuer-grantor-period] [verifier-grantor-period] [issuer-period] [verifier-period] [holder-period] [issuer-mode] [verifier-mode]",
					Short:     "Create a new credential schema",
					Long: `Create a new credential schema with the specified parameters:
- id: unique identifier for the schema
- tr-id: trust registry ID
- json-schema: path to JSON schema file or JSON string
- issuer-grantor-period: validation period for issuer grantors (days)
- verifier-grantor-period: validation period for verifier grantors (days)
- issuer-period: validation period for issuers (days)
- verifier-period: validation period for verifiers (days)
- holder-period: validation period for holders (days)
- issuer-mode: permission management mode for issuers (0=UNSPECIFIED, 1=OPEN, 2=GRANTOR_VALIDATION, 3=TRUST_REGISTRY_VALIDATION)
- verifier-mode: permission management mode for verifiers (0=UNSPECIFIED, 1=OPEN, 2=GRANTOR_VALIDATION, 3=TRUST_REGISTRY_VALIDATION)

Example:
$ veranad tx credentialschema create-credential-schema 1 1 schema.json 365 365 180 180 180 2 2`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
						{
							ProtoField: "tr_id",
						},
						{
							ProtoField: "json_schema",
						},
						{
							ProtoField: "issuer_grantor_validation_validity_period",
						},
						{
							ProtoField: "verifier_grantor_validation_validity_period",
						},
						{
							ProtoField: "issuer_validation_validity_period",
						},
						{
							ProtoField: "verifier_validation_validity_period",
						},
						{
							ProtoField: "holder_validation_validity_period",
						},
						{
							ProtoField: "issuer_perm_management_mode",
						},
						{
							ProtoField: "verifier_perm_management_mode",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
