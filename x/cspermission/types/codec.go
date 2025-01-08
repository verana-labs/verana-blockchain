package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateCredentialSchemaPerm{}, "/vpr/v1/csp/create-csp")
	legacy.RegisterAminoMsg(cdc, &MsgRevokeCredentialSchemaPerm{}, "/vpr/v1/csp/revoke-csp")
	legacy.RegisterAminoMsg(cdc, &MsgTerminateCredentialSchemaPerm{}, "/vpr/v1/csp/terminate-csp")
	legacy.RegisterAminoMsg(cdc, &MsgCreateOrUpdateCSPS{}, "/vpr/v1/csp/create-or-update-csps")
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgCreateCredentialSchemaPerm{},
		&MsgRevokeCredentialSchemaPerm{},
		&MsgTerminateCredentialSchemaPerm{},
		&MsgCreateOrUpdateCSPS{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
