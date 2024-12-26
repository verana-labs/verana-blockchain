package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"regexp"
)

const TypeMsgCreateValidation = "create_validation"
const TypeMsgRenewValidation = "renew_validation"

var _ sdk.Msg = &MsgCreateValidation{}
var _ sdk.Msg = &MsgRenewValidation{}

// NewMsgCreateValidation creates a new MsgCreateValidation instance
func NewMsgCreateValidation(
	creator string,
	validationType ValidationType,
	validatorPermId uint64,
	country string,
) *MsgCreateValidation {
	return &MsgCreateValidation{
		Creator:         creator,
		ValidationType:  uint32(validationType),
		ValidatorPermId: validatorPermId,
		Country:         country,
	}
}

// Route implements sdk.Msg
func (msg *MsgCreateValidation) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg *MsgCreateValidation) Type() string {
	return TypeMsgCreateValidation
}

// GetSigners implements sdk.Msg
func (msg *MsgCreateValidation) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic implements [MOD-V-MSG-1-2-1] basic checks
func (msg *MsgCreateValidation) ValidateBasic() error {
	// Check mandatory parameters
	if msg.Creator == "" {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "creator is required")
	}
	if msg.ValidatorPermId == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "validator permission id is required")
	}
	if msg.Country == "" {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "country is required")
	}

	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate country code (ISO 3166-1 alpha-2)
	if !isValidCountryCode(msg.Country) {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid country code format: must be ISO 3166-1 alpha-2")
	}

	//Validate validation type
	if msg.ValidationType == uint32(ValidationType_TYPE_UNSPECIFIED) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "validation type must be specified")
	}
	if !isValidValidationType(ValidationType(msg.ValidationType)) {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid validation type: %d", msg.ValidationType)
	}

	return nil
}

// isValidCountryCode validates ISO 3166-1 alpha-2 country codes
func isValidCountryCode(country string) bool {
	pattern := regexp.MustCompile(`^[A-Z]{2}$`)
	return pattern.MatchString(country)
}

// isValidValidationType checks if the validation type is valid
func isValidValidationType(vType ValidationType) bool {
	// Check if type is within valid range (excluding UNSPECIFIED)
	return vType > ValidationType_TYPE_UNSPECIFIED && vType <= ValidationType_HOLDER
}

func NewMsgRenewValidation(
	creator string,
	id uint64,
	validatorPermId uint64,
) *MsgRenewValidation {
	return &MsgRenewValidation{
		Creator:         creator,
		Id:              id,
		ValidatorPermId: validatorPermId,
	}
}

// Route implements sdk.Msg
func (msg *MsgRenewValidation) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg *MsgRenewValidation) Type() string {
	return TypeMsgRenewValidation
}

// GetSigners implements sdk.Msg
func (msg *MsgRenewValidation) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgRenewValidation) ValidateBasic() error {
	if msg.Creator == "" {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "creator address is required")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Id == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "validation id is required")
	}

	return nil
}
