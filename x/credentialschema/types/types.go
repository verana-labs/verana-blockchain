package types

import (
	"cosmossdk.io/errors"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/xeipuuv/gojsonschema"
)

// Official meta-schema for Draft 2020-12
const jsonSchemaMetaSchema = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.com/meta-schema/credential-schema",
  "title": "Credential Schema Meta-Schema",
  "type": "object",
  "required": ["$id", "$schema", "type", "properties"],
  "properties": {
    "$id": {
      "type": "string",
      "format": "uri-reference",
      "pattern": "^/dtr/v1/cs/js/\\d+$",
      "description": "$id must be a URI matching the rendering URL format"
    },
    "$schema": {
      "type": "string",
      "enum": ["https://json-schema.org/draft/2020-12/schema"],
      "description": "$schema must be the Draft 2020-12 URI"
    },
    "type": {
      "type": "string",
      "enum": ["object"],
      "description": "The root type must be 'object'"
    },
    "properties": {
      "type": "object",
      "additionalProperties": {
        "type": "object",
        "properties": {
          "type": {
            "type": "string",
            "enum": ["string", "number", "integer", "boolean", "object", "array"],
            "description": "The type of each property"
          },
          "description": {
            "type": "string"
          },
          "default": {
            "type": ["string", "number", "integer", "boolean", "object", "array", "null"]
          }
        },
        "required": ["type"]
      }
    },
    "required": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "description": "List of required properties"
    },
    "additionalProperties": {
      "type": "boolean",
      "default": true
    },
    "$defs": {
      "type": "object",
      "additionalProperties": {
        "type": "object"
      },
      "description": "Optional definitions for reusable schema components"
    }
  },
  "additionalProperties": false,
  "examples": [
    {
      "$schema": "https://json-schema.org/draft/2020-12/schema",
      "$id": "/dtr/v1/cs/js/1",
      "type": "object",
      "properties": {
        "name": {
          "type": "string",
          "description": "Name of the entity"
        }
      },
      "required": ["name"],
      "additionalProperties": false
    }
  ]
}
`

const TypeMsgCreateCredentialSchema = "create_credential_schema"

var _ sdk.Msg = &MsgCreateCredentialSchema{}

// NewMsgCreateCredentialSchema creates a new MsgCreateCredentialSchema instance
func NewMsgCreateCredentialSchema(
	creator string,
	id uint64,
	trId uint64,
	jsonSchema string,
	issuerGrantorValidationValidityPeriod uint32,
	verifierGrantorValidationValidityPeriod uint32,
	issuerValidationValidityPeriod uint32,
	verifierValidationValidityPeriod uint32,
	holderValidationValidityPeriod uint32,
	issuerPermManagementMode uint32,
	verifierPermManagementMode uint32,
) *MsgCreateCredentialSchema {
	return &MsgCreateCredentialSchema{
		Creator:                                 creator,
		Id:                                      id,
		TrId:                                    trId,
		JsonSchema:                              jsonSchema,
		IssuerGrantorValidationValidityPeriod:   issuerGrantorValidationValidityPeriod,
		VerifierGrantorValidationValidityPeriod: verifierGrantorValidationValidityPeriod,
		IssuerValidationValidityPeriod:          issuerValidationValidityPeriod,
		VerifierValidationValidityPeriod:        verifierValidationValidityPeriod,
		HolderValidationValidityPeriod:          holderValidationValidityPeriod,
		IssuerPermManagementMode:                issuerPermManagementMode,
		VerifierPermManagementMode:              verifierPermManagementMode,
	}
}

// Route implements sdk.Msg
func (msg *MsgCreateCredentialSchema) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (msg *MsgCreateCredentialSchema) Type() string {
	return TypeMsgCreateCredentialSchema
}

// GetSigners implements sdk.Msg
func (msg *MsgCreateCredentialSchema) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes implements sdk.Msg
func (msg *MsgCreateCredentialSchema) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements sdk.Msg
func (msg *MsgCreateCredentialSchema) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Check mandatory parameters
	if msg.Id == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "id cannot be 0")
	}

	if msg.TrId == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "tr_id cannot be 0")
	}

	// Validate JSON Schema
	if err := validateJSONSchema(msg.JsonSchema, msg.Id); err != nil {
		return errors.Wrapf(ErrInvalidJSONSchema, err.Error())
	}

	// Validate validity periods (must be >= 0)
	if err := validateValidityPeriods(msg); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Validate permission management modes
	if err := validatePermManagementModes(msg); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

func validateJSONSchema(schemaJSON string, id uint64) error {
	if schemaJSON == "" {
		return fmt.Errorf("json schema cannot be empty")
	}

	// Parse JSON
	var schemaDoc map[string]interface{}
	if err := json.Unmarshal([]byte(schemaJSON), &schemaDoc); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	if len(schemaJSON) > int(DefaultCredentialSchemaSchemaMaxSize) {
		return fmt.Errorf("json schema exceeds maximum size of %d bytes", DefaultCredentialSchemaSchemaMaxSize)
	}

	// Check for $id field and print it for debugging
	schemaId, ok := schemaDoc["$id"].(string)
	if !ok {
		return fmt.Errorf("$id must be a string")
	}

	expectedUrl := fmt.Sprintf("/dtr/v1/cs/js/%d", id) // Adjust URL pattern as per your API
	if schemaId != expectedUrl {
		return fmt.Errorf("$id must match the schema query URL pattern: %s", expectedUrl)
	}

	// Load the meta-schema and validate
	metaSchemaLoader := gojsonschema.NewStringLoader(jsonSchemaMetaSchema)
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)
	result, err := gojsonschema.Validate(metaSchemaLoader, schemaLoader)
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}

	if !result.Valid() {
		var errMsgs []string
		for _, err := range result.Errors() {
			errMsgs = append(errMsgs, err.String())
		}
		return fmt.Errorf("invalid JSON schema: %v", errMsgs)
	}

	// Check required fields
	requiredFields := []string{"$schema", "$id", "type"}
	for _, field := range requiredFields {
		if _, ok := schemaDoc[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate type is 'object'
	if schemaType, ok := schemaDoc["type"].(string); !ok || schemaType != "object" {
		return fmt.Errorf("root schema type must be 'object'")
	}

	// Validate properties exist
	properties, ok := schemaDoc["properties"].(map[string]interface{})
	if !ok || len(properties) == 0 {
		return fmt.Errorf("schema must define non-empty properties")
	}
	return nil
}

func validateValidityPeriods(msg *MsgCreateCredentialSchema) error {
	if msg.IssuerGrantorValidationValidityPeriod < 0 {
		return fmt.Errorf("issuer grantor validation validity period cannot be negative")
	}
	if msg.VerifierGrantorValidationValidityPeriod < 0 {
		return fmt.Errorf("verifier grantor validation validity period cannot be negative")
	}
	// Add maximum value checks
	if msg.IssuerGrantorValidationValidityPeriod > 0 &&
		msg.IssuerGrantorValidationValidityPeriod > DefaultCredentialSchemaIssuerGrantorValidationValidityPeriodMaxDays {
		return fmt.Errorf("issuer grantor validation validity period exceeds maximum allowed days")
	}

	if msg.VerifierGrantorValidationValidityPeriod > 0 &&
		msg.VerifierGrantorValidationValidityPeriod > DefaultCredentialSchemaVerifierGrantorValidationValidityPeriodMaxDays {
		return fmt.Errorf("verifier grantor validation validity period exceeds maximum allowed days")
	}

	if msg.IssuerValidationValidityPeriod < 0 {
		return fmt.Errorf("issuer validation validity period cannot be negative")
	}
	if msg.VerifierValidationValidityPeriod < 0 {
		return fmt.Errorf("verifier validation validity period cannot be negative")
	}
	if msg.HolderValidationValidityPeriod < 0 {
		return fmt.Errorf("holder validation validity period cannot be negative")
	}
	return nil
}

func validatePermManagementModes(msg *MsgCreateCredentialSchema) error {
	// Check issuer permission management mode
	if msg.IssuerPermManagementMode == 0 {
		return fmt.Errorf("issuer permission management mode must be specified")
	}
	if msg.IssuerPermManagementMode > 3 {
		return fmt.Errorf("invalid issuer permission management mode: must be between 1 and 3")
	}

	// Check verifier permission management mode
	if msg.VerifierPermManagementMode == 0 {
		return fmt.Errorf("verifier permission management mode must be specified")
	}
	if msg.VerifierPermManagementMode > 3 {
		return fmt.Errorf("invalid verifier permission management mode: must be between 1 and 3")
	}

	return nil
}
