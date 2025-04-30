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
				{
					RpcMethod: "ListPermissions",
					Use:       "list-permissions",
					Short:     "List all permissions",
					Long:      "List all permissions with optional filtering by modified time and pagination",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"modified_after": {
							Name:         "modified-after",
							Usage:        "Filter by modified time (RFC3339 format)",
							DefaultValue: "",
						},
						"response_max_size": {
							Name:         "response-max-size",
							Usage:        "Maximum number of results to return (1-1024)",
							DefaultValue: "64",
						},
					},
				},
				{
					RpcMethod: "GetPermission",
					Use:       "get-permission [id]",
					Short:     "Get permission by ID",
					Long:      "Get detailed information about a permission by its ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				{
					RpcMethod: "GetPermissionSession",
					Use:       "get-permission-session [id]",
					Short:     "Get permission session by ID",
					Long:      "Get details about a specific permission session by its ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				{
					RpcMethod: "ListPermissionSessions",
					Use:       "list-permission-sessions",
					Short:     "List permission sessions",
					Long:      "List all permission sessions with optional filtering and pagination",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"modified_after": {
							Name:         "modified-after",
							Usage:        "Filter by modified time (RFC3339 format)",
							DefaultValue: "",
						},
						"response_max_size": {
							Name:         "response-max-size",
							Usage:        "Maximum number of results to return (1-1024)",
							DefaultValue: "64",
						},
					},
				},
				{
					RpcMethod: "FindPermissionsWithDID",
					Use:       "find-permissions-with-did [did] [type] [schema-id]",
					Short:     "Find permissions with DID",
					Long:      "Find permissions matching the specified DID, type, and schema ID with optional filtering by country and timestamp",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "did"},
						{ProtoField: "type"},
						{ProtoField: "schema_id"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"country": {
							Name:         "country",
							DefaultValue: "",
							Usage:        "Filter by country code (ISO 3166-1 alpha-2)",
						},
						"when": {
							Name:         "when",
							DefaultValue: "",
							Usage:        "Filter by validity at specified timestamp (RFC3339 format)",
						},
					},
				},
				{
					RpcMethod: "FindBeneficiaries",
					Use:       "find-beneficiaries",
					Short:     "Find beneficiary permissions in the permission tree",
					Long:      "Find beneficiary permissions by traversing the permission tree for issuer and/or verifier permissions",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"issuer_perm_id": {
							Name:         "issuer-perm-id",
							DefaultValue: "0",
							Usage:        "ID of the issuer permission",
						},
						"verifier_perm_id": {
							Name:         "verifier-perm-id",
							DefaultValue: "0",
							Usage:        "ID of the verifier permission",
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
					Use:       "request-vp-termination [id]",
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
				// Add to the RpcCommandOptions array in the Tx ServiceCommandDescriptor:
				{
					RpcMethod: "ConfirmPermissionVPTermination",
					Use:       "confirm-vp-termination [id]",
					Short:     "Confirm the termination of a permission VP",
					Long:      "Confirm the termination of a permission VP. Can be called by the validator, or by the grantee after the timeout period has elapsed.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				{
					RpcMethod: "CancelPermissionVPLastRequest",
					Use:       "cancel-permission-vp-request [id]",
					Short:     "Cancel a pending permission VP request",
					Long:      "Cancel a pending permission VP request. Can only be executed by the permission grantee and only when the permission is in PENDING state.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				{
					RpcMethod: "CreateRootPermission",
					Use:       "create-root-permission [schema-id] [did] [validation-fees] [issuance-fees] [verification-fees]",
					Short:     "Create a new root permission for a credential schema",
					Long:      "Create a new root permission for a credential schema. Can only be executed by the trust registry controller.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "schema_id",
						},
						{
							ProtoField: "did",
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
						"country": {
							Name:         "country",
							DefaultValue: "",
							Usage:        "Optional country code (ISO 3166-1 alpha-2)",
						},
						"effective_from": {
							Name:         "effective-from",
							DefaultValue: "",
							Usage:        "Optional timestamp (RFC3339) from when the permission is effective",
						},
						"effective_until": {
							Name:         "effective-until",
							DefaultValue: "",
							Usage:        "Optional timestamp (RFC3339) until when the permission is effective",
						},
					},
				},
				{
					RpcMethod: "ExtendPermission",
					Use:       "extend-permission [id] [effective-until]",
					Short:     "Extend a permission's effective duration",
					Long:      "Extend a permission's effective duration. Can only be executed by the validator of the permission.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
						{
							ProtoField: "effective_until",
						},
					},
				},
				{
					RpcMethod: "RevokePermission",
					Use:       "revoke-permission [id]",
					Short:     "Revoke a permission",
					Long:      "Revoke a permission. Can only be executed by the validator of the permission.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
					},
				},
				{
					RpcMethod: "CreateOrUpdatePermissionSession",
					Use:       "create-or-update-perm-session [id] [agent-perm-id]",
					Short:     "Create or update a permission session",
					Long: `Create or update a permission session with the specified parameters:
- id: UUID of the session
- agent-perm-id: ID of the agent permission (HOLDER)
Optional parameters:
- issuer-perm-id: ID of the issuer permission
- verifier-perm-id: ID of the verifier permission
- wallet-agent-perm-id: ID of the wallet agent permission if different from agent

At least one of issuer-perm-id or verifier-perm-id must be provided.`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
						{
							ProtoField: "agent_perm_id",
						},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"issuer_perm_id": {
							Name:         "issuer-perm-id",
							Usage:        "ID of the issuer permission",
							DefaultValue: "0",
						},
						"verifier_perm_id": {
							Name:         "verifier-perm-id",
							Usage:        "ID of the verifier permission",
							DefaultValue: "0",
						},
						"wallet_agent_perm_id": {
							Name:         "wallet-agent-perm-id",
							Usage:        "ID of the wallet agent permission if different from agent",
							DefaultValue: "0",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
