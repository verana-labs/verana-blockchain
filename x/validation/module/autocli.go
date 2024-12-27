package validation

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/verana-labs/verana-blockchain/api/veranablockchain/validation"
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
					RpcMethod: "ListValidations",
					Use:       "list-validations",
					Short:     "List validations with optional filters",
					Long:      "List validations with optional filters: controller, validator permission ID, type, state, and expiration",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"controller": {
							Name:         "controller",
							Usage:        "Filter by controller account address",
							DefaultValue: "",
						},
						"validator_perm_id": {
							Name:         "validator-perm-id",
							Usage:        "Filter by validator permission ID",
							DefaultValue: "0",
						},
						"type": {
							Name:         "type",
							Usage:        "Filter by validation type (ISSUER_GRANTOR, VERIFIER_GRANTOR, ISSUER, VERIFIER, HOLDER)",
							DefaultValue: "",
						},
						"state": {
							Name:         "state",
							Usage:        "Filter by validation state (PENDING, VALIDATED, TERMINATED)",
							DefaultValue: "",
						},
						"response_max_size": {
							Name:         "response-max-size",
							Usage:        "Maximum number of results (1-1024)",
							DefaultValue: "64",
						},
						"exp_before": {
							Name:         "exp-before",
							Usage:        "Filter by expiration before timestamp (RFC3339 format)",
							DefaultValue: "",
						},
					},
				},
				{
					RpcMethod: "GetValidation",
					Use:       "get-validation [id]",
					Short:     "Get validation by ID",
					Long:      "Get the validation information for a given validation ID",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
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
					RpcMethod: "CreateValidation",
					Use:       "create-validation [validation-type] [validator-perm-id] [country]",
					Short:     "Create a new validation entry",
					Long: `Create a new validation entry with the specified parameters:
- type: ISSUER_GRANTOR, VERIFIER_GRANTOR, ISSUER, VERIFIER, HOLDER
- validator-perm-id: ID of the validator's permission
- country: Alpha-2 country code (ISO 3166)`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "validation_type",
						},
						{
							ProtoField: "validator_perm_id",
						},
						{
							ProtoField: "country",
						},
					},
				},
				{
					RpcMethod: "RenewValidation",
					Use:       "renew-validation [id] [validator-perm-id]",
					Short:     "Renew an existing validation",
					Long: `Renew an existing validation with optional new validator:
- id: ID of the validation to renew
- validator-perm-id: Optional new validator permission ID (if not specified, uses current validator)`,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
						{
							ProtoField: "validator_perm_id",
							Optional:   true,
						},
					},
				},
				{
					RpcMethod: "SetValidated",
					Use:       "set-validated [id] [summary-hash]",
					Short:     "Set a validation to VALIDATED state",
					Long:      "Set a validation to VALIDATED state. Only the validator can execute this method. Optional summary hash can be provided for non-HOLDER validations.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "id",
						},
						{
							ProtoField: "summary_hash",
							Optional:   true,
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
