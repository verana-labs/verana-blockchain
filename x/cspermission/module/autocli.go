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
  [schema-id]        : ID of the credential schema (uint64)
  [csp-type]         : Permission type (ISSUER, VERIFIER, ISSUER_GRANTOR, VERIFIER_GRANTOR, TRUST_REGISTRY, HOLDER)
  [did]              : DID of the grantee service
  [grantee]          : Account address of the grantee
  [effective-from]   : Start date (format: yyyyMMddHHmm)
  [validation-fees]  : Fees for validation process
  [issuance-fees]    : Fees for credential issuance
  [verification-fees]: Fees for credential verification

Optional Flags:
  --effective-until  : End date (format: yyyyMMddHHmm)
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
							Usage:        "End date (format: yyyyMMddHHmm)",
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
