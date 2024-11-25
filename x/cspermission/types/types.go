package types

import (
	"regexp"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateCredentialSchemaPerm = "create_credential_schema_perm"

var (
	// Compiled regex patterns
	didRegex     = regexp.MustCompile(`^did:[a-zA-Z0-9]+:[a-zA-Z0-9._-]+$`)
	countryRegex = regexp.MustCompile(`^[A-Z]{2}$`)
)

var _ sdk.Msg = &MsgCreateCredentialSchemaPerm{}

// NewMsgCreateCredentialSchemaPerm creates a new MsgCreateCredentialSchemaPerm instance
func NewMsgCreateCredentialSchemaPerm(
	creator string,
	schemaId uint64,
	permType uint32,
	did string,
	grantee string,
	effectiveFrom time.Time,
	effectiveUntil *time.Time,
	country string,
	validationId uint64,
	validationFees uint64,
	issuanceFees uint64,
	verificationFees uint64,
) *MsgCreateCredentialSchemaPerm {
	return &MsgCreateCredentialSchemaPerm{
		Creator:          creator,
		SchemaId:         schemaId,
		CspType:          permType,
		Did:              did,
		Grantee:          grantee,
		EffectiveFrom:    effectiveFrom,
		EffectiveUntil:   effectiveUntil,
		Country:          country,
		ValidationId:     validationId,
		ValidationFees:   validationFees,
		IssuanceFees:     issuanceFees,
		VerificationFees: verificationFees,
	}
}

// Route implements sdk.Msg
func (msg *MsgCreateCredentialSchemaPerm) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg *MsgCreateCredentialSchemaPerm) Type() string {
	return TypeMsgCreateCredentialSchemaPerm
}

// GetSigners implements sdk.Msg
func (msg *MsgCreateCredentialSchemaPerm) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes implements sdk.Msg
func (msg *MsgCreateCredentialSchemaPerm) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgCreateCredentialSchemaPerm) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if !msg.EffectiveFrom.After(time.Now()) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "effective_from must be in the future")
	}

	// Validate mandatory parameters
	if msg.SchemaId == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "schema_id cannot be 0")
	}

	if msg.CspType < 1 || msg.CspType > 6 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "csp_type must be between 1 and 6")
	}

	permType := CredentialSchemaPermType(msg.CspType)
	switch permType {
	case CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER,
		CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER,
		CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_ISSUER_GRANTOR,
		CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_VERIFIER_GRANTOR,
		CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_TRUST_REGISTRY,
		CredentialSchemaPermType_CREDENTIAL_SCHEMA_PERM_TYPE_HOLDER:
		// Valid types
	default:
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid credential schema permission type")
	}

	// Validate DID
	if msg.Did == "" {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "did cannot be empty")
	}
	if !didRegex.MatchString(msg.Did) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid did format")
	}

	// Validate grantee address
	_, err = sdk.AccAddressFromBech32(msg.Grantee)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid grantee address (%s)", err)
	}

	// Validate effective dates
	if msg.EffectiveFrom.IsZero() {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "effective_from cannot be zero time")
	}

	if msg.EffectiveUntil != nil {
		if !msg.EffectiveUntil.After(msg.EffectiveFrom) {
			return errors.Wrap(sdkerrors.ErrInvalidRequest, "effective_until must be after effective_from")
		}
	}

	// Validate country code if present
	if msg.Country != "" && !countryRegex.MatchString(msg.Country) {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid country code format")
	}

	return nil
}

func (q *QueryListCSPRequest) ValidateRequest() error {
	if q.SchemaId == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("schema_id is required")
	}

	if q.Creator != "" {
		if _, err := sdk.AccAddressFromBech32(q.Creator); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrap(err.Error())
		}
	}

	if q.Grantee != "" {
		if _, err := sdk.AccAddressFromBech32(q.Grantee); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrap(err.Error())
		}
	}

	if q.Did != "" && !IsValidDID(q.Did) {
		return sdkerrors.ErrInvalidRequest.Wrap("invalid DID format")
	}

	if q.ResponseMaxSize == 0 {
		q.ResponseMaxSize = 64
	} else if q.ResponseMaxSize > 1024 {
		return sdkerrors.ErrInvalidRequest.Wrap("response_max_size must be between 1 and 1024")
	}

	return nil
}

func IsValidDID(did string) bool {
	// Basic DID validation regex
	// This is a simplified version and may need to be expanded based on specific DID method requirements
	didRegex := regexp.MustCompile(`^did:[a-zA-Z0-9]+:[a-zA-Z0-9._-]+$`)
	return didRegex.MatchString(did)
}
