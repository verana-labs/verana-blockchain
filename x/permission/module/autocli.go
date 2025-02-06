package permission

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/verana-labs/verana-blockchain/api/veranablockchain/permission"
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
					RpcMethod: "StartPermissionVP",
					Use:       "start-permission-vp [type] [validator-perm-id] [country]",
					Short:     "Start a new permission validation process",
					Long: `Start a new permission validation process with the specified parameters:
- type: Permission type (0=Unspecified, 1=Issuer, 2=Verifier, 3=IssuerGrantor, 4=VerifierGrantor, 5=TrustRegistry, 6=Holder)
- validator-perm-id: ID of the validator permission
- country: ISO 3166-1 alpha-2 country code`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "type",
						},
						{
							ProtoField: "validator_perm_id",
						},
						{
							ProtoField: "country",
						},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"did": {
							Name:         "did",
							Usage:        "Optional DID for this permission",
							DefaultValue: "",
						},
					},
				},
				{
					RpcMethod: "RenewPermissionVP",
					Use:       "renew-permission-vp [id]",
					Short:     "Renew a permission validation process",
					Long: `Renew a permission validation process for an existing permission:
- id: ID of the permission to renew`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				{
					RpcMethod: "SetPermissionVPToValidated",
					Use:       "set-permission-vp-validated [id]",
					Short:     "Set permission validation process to validated state",
					Long: `Set a permission validation process to validated state with optional parameters:
- id: ID of the permission to validate
- effective-until: Optional timestamp until when this permission is effective (RFC3339 format)
- validation-fees: Optional validation fees
- issuance-fees: Optional issuance fees
- verification-fees: Optional verification fees
- country: Optional country code (ISO 3166-1 alpha-2)
- vp-summary-digest-sri: Optional digest SRI of validation information`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"effective_until": {
							Name:         "effective-until",
							Usage:        "Timestamp until when this permission is effective (RFC3339)",
							DefaultValue: "",
						},
						"validation_fees": {
							Name:         "validation-fees",
							Usage:        "Validation fees",
							DefaultValue: "0",
						},
						"issuance_fees": {
							Name:         "issuance-fees",
							Usage:        "Issuance fees",
							DefaultValue: "0",
						},
						"verification_fees": {
							Name:         "verification-fees",
							Usage:        "Verification fees",
							DefaultValue: "0",
						},
						"country": {
							Name:         "country",
							Usage:        "Country code (ISO 3166-1 alpha-2)",
							DefaultValue: "",
						},
						"vp_summary_digest_sri": {
							Name:         "vp-summary-digest-sri",
							Usage:        "Digest SRI of validation information",
							DefaultValue: "",
						},
					},
				},
				{
					RpcMethod: "RequestPermissionVPTermination",
					Use:       "request-permission-vp-termination [id]",
					Short:     "Request termination of a permission validation process",
					Long: `Request termination of a permission validation process:
- id: ID of the permission validation process to terminate
Note: For expired VPs, either the grantee or validator can request termination.
For active VPs, only the grantee can request termination unless it's a HOLDER type.`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
