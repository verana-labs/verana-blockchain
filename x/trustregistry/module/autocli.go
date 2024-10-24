package trustregistry

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	modulev1 "github.com/verana-labs/verana-blockchain/api/veranablockchain/trustregistry"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "GetTrustRegistry",
					Use:       "get-trust-registry [tr_id]",
					Short:     "Get trust registry information by ID",
					Long:      "Get the trust registry information for a given trust registry ID, with options to filter by active governance framework and preferred language",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "tr_id"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"active_gf_only": {
							Name:         "active-gf-only",
							DefaultValue: "false",
							Usage:        "If true, include only current governance framework data",
						},
						"preferred_language": {
							Name:         "preferred-language",
							DefaultValue: "",
							Usage:        "Preferred language for the returned documents",
						},
					},
				},
				{
					RpcMethod: "GetTrustRegistryWithDID",
					Use:       "get-trust-registry-by-did [did]",
					Short:     "Get trust registry information by DID",
					Long:      "Get the trust registry information for a given DID, with options to filter by active governance framework and preferred language",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "did"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"active_gf_only": {
							Name:         "active-gf-only",
							DefaultValue: "false",
							Usage:        "If true, include only current governance framework data",
						},
						"preferred_language": {
							Name:         "preferred-language",
							DefaultValue: "",
							Usage:        "Preferred language for the returned documents",
						},
					},
				},
				{
					RpcMethod: "ListTrustRegistries",
					Use:       "list-trust-registries",
					Short:     "List Trust Registries",
					Long:      "List Trust Registries with optional filtering and pagination. Results are ordered by modified time ascending.",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"controller": {
							Name:         "controller",
							Usage:        "Filter by controller account address",
							DefaultValue: "",
						},
						"modified_after": {
							Name:         "modified-after",
							Usage:        "Filter by modified time (RFC3339 format)",
							DefaultValue: "",
						},
						"active_gf_only": {
							Name:         "active-gf-only",
							Usage:        "Include only current governance framework data",
							DefaultValue: "false",
						},
						"preferred_language": {
							Name:         "preferred-language",
							Usage:        "Preferred language for returned documents",
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
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Get the current module parameters",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "CreateTrustRegistry",
					Use:       "create-trust-registry [did] [language] [doc-url] [doc-hash]",
					Short:     "Create a new trust registry",
					Long:      "Create a new trust registry with the specified DID, language, and initial governance framework document",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "did"},
						{ProtoField: "language"},
						{ProtoField: "doc_url"},
						{ProtoField: "doc_hash"},
					},
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"aka": {
							Name:         "aka",
							Usage:        "Optional additional URI for the trust registry",
							DefaultValue: "",
						},
					},
				},
				{
					RpcMethod: "AddGovernanceFrameworkDocument",
					Use:       "add-governance-framework-document [tr_id] [doc-language] [doc-url] [doc-hash] [version]",
					Short:     "Add a new governance framework document",
					Long:      "Add a new governance framework document to an existing trust registry. The version must be either equal to the highest existing version or exactly one more than the highest version.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "tr_id"},
						{ProtoField: "doc_language"},
						{ProtoField: "doc_url"},
						{ProtoField: "doc_hash"},
						{ProtoField: "version"},
					},
				},
				{
					RpcMethod: "IncreaseActiveGovernanceFrameworkVersion",
					Use:       "increase-active-gf-version [tr_id]",
					Short:     "Increase the active governance framework version",
					Long:      "Increase the active governance framework version for a trust registry. This can only be done by the controller of the trust registry and requires a document in the trust registry's default language for the new version.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "tr_id"},
					},
				},
			},
		},
	}
}
