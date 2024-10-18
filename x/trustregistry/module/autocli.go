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
					Use:       "get-trust-registry [did]",
					Short:     "Get the trust registry information for a given DID",
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
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Get the current module parameters",
				},
				{
					RpcMethod: "ListTrustRegistries",
					Use:       "list-trust-registries [flags]",
					Short:     "List Trust Registries",
					Long:      "List Trust Registries with optional filtering and pagination",
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"modified": {
							Name:         "modified",
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
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "CreateTrustRegistry",
					Use:       "create-trust-registry [did] [aka] [language] [doc-url] [doc-hash]",
					Short:     "Create a new trust registry",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "did"},
						{ProtoField: "aka"},
						{ProtoField: "language"},
						{ProtoField: "doc_url"},
						{ProtoField: "doc_hash"},
					},
				},
				{
					RpcMethod: "AddGovernanceFrameworkDocument",
					Use:       "add-governance-framework-document [did] [doc-language] [doc-url] [doc-hash] [version]",
					Short:     "Add a new governance framework document to an existing trust registry",
					Long:      "Add a new governance framework document to an existing trust registry. The version must be either equal to the highest existing version or exactly one more.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "did"},
						{ProtoField: "doc_language"},
						{ProtoField: "doc_url"},
						{ProtoField: "doc_hash"},
						{ProtoField: "version"},
					},
				},
				{
					RpcMethod: "IncreaseActiveGovernanceFrameworkVersion",
					Use:       "increase-active-gf-version [did]",
					Short:     "Increase the active governance framework version for a trust registry",
					Long:      "Increase the active governance framework version for a trust registry. This can only be done by the controller of the trust registry.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "did"},
					},
				},
			},
		},
	}
}
