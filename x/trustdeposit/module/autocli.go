package trustdeposit

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/verana-labs/verana-blockchain/api/veranablockchain/trustdeposit"
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
					RpcMethod: "GetTrustDeposit",
					Use:       "get-trust-deposit [account]",
					Short:     "Query trust deposit for an account",
					Long:      "Get the trust deposit information for a given account address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "account",
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
					RpcMethod: "ReclaimTrustDepositInterests",
					Use:       "reclaim-interests",
					Short:     "Reclaim earned interest from trust deposits",
					Long:      "Reclaim any available interest earned from trust deposits. The interest is calculated based on share value and current deposit amount.",
				},
				{
					RpcMethod: "ReclaimTrustDeposit",
					Use:       "reclaim-deposit [amount]",
					Short:     "Reclaim trust deposit",
					Long:      "Reclaim a specified amount from your claimable trust deposit balance. Note that a portion will be burned according to the reclaim burn rate.",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{
							ProtoField: "claimed",
						},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
