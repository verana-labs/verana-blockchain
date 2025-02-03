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
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
