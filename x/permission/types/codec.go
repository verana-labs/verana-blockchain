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
	legacy.RegisterAminoMsg(cdc, &MsgStartPermissionVP{}, "/perm/v1/start-permission-vp")
	legacy.RegisterAminoMsg(cdc, &MsgRenewPermissionVP{}, "/perm/v1/renew-permission-vp")
	legacy.RegisterAminoMsg(cdc, &MsgSetPermissionVPToValidated{}, "/perm/v1/set-permission-vp-validated")
	legacy.RegisterAminoMsg(cdc, &MsgRequestPermissionVPTermination{}, "/perm/v1/request-vp-termination")
	legacy.RegisterAminoMsg(cdc, &MsgConfirmPermissionVPTermination{}, "/perm/v1/confirm-vp-termination")
	legacy.RegisterAminoMsg(cdc, &MsgCancelPermissionVPLastRequest{}, "/perm/v1/cancel-permission-vp-request")
	legacy.RegisterAminoMsg(cdc, &MsgCreateRootPermission{}, "/perm/v1/create-root-permission")
	legacy.RegisterAminoMsg(cdc, &MsgExtendPermission{}, "/perm/v1/extend-permission")
	legacy.RegisterAminoMsg(cdc, &MsgRevokePermission{}, "/perm/v1/revoke-permission")
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgStartPermissionVP{},
		&MsgRenewPermissionVP{},
		&MsgSetPermissionVPToValidated{},
		&MsgRequestPermissionVPTermination{},
		&MsgConfirmPermissionVPTermination{},
		&MsgCancelPermissionVPLastRequest{},
		&MsgCreateRootPermission{},
		&MsgExtendPermission{},
		&MsgRevokePermission{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
