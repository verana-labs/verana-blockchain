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
  "$id": "https://json-schema.org/draft/2020-12/schema",
  "$vocabulary": {
    "https://json-schema.org/draft/2020-12/vocab/core": true,
    "https://json-schema.org/draft/2020-12/vocab/applicator": true,
    "https://json-schema.org/draft/2020-12/vocab/unevaluated": true,
    "https://json-schema.org/draft/2020-12/vocab/validation": true,
    "https://json-schema.org/draft/2020-12/vocab/meta-data": true,
    "https://json-schema.org/draft/2020-12/vocab/format-annotation": true,
    "https://json-schema.org/draft/2020-12/vocab/content": true
  },
  "$dynamicAnchor": "meta",
  "title": "Core and Validation specifications meta-schema",
  "allOf": [
    {
      "$ref": "https://json-schema.org/draft/2020-12/meta/core"
    },
    {
      "$ref": "https://json-schema.org/draft/2020-12/meta/applicator"
    },
    {
      "$ref": "https://json-schema.org/draft/2020-12/meta/unevaluated"
    },
    {
      "$ref": "https://json-schema.org/draft/2020-12/meta/validation"
    },
    {
      "$ref": "https://json-schema.org/draft/2020-12/meta/meta-data"
    },
    {
      "$ref": "https://json-schema.org/draft/2020-12/meta/format-annotation"
    },
    {
      "$ref": "https://json-schema.org/draft/2020-12/meta/content"
    }
  ],
  "type": [
    "object",
    "boolean"
  ],
  "properties": {
    "definitions": {
      "$comment": "While no longer an official keyword as it is replaced by $defs, this keyword is retained in the meta-schema to prevent incompatible extensions as it remains in common use.",
      "type": "object",
      "additionalProperties": {
        "$dynamicRef": "#meta"
      },
      "default": {}
    },
    "$defs": {
      "type": "object",
      "additionalProperties": {
        "$dynamicRef": "#meta"
      }
    },
    "$id": {
      "type": "string",
      "format": "uri-reference",
      "$comment": "Non-empty fragments not allowed.",
      "pattern": "^[^#]*#?$"
    },
    "$schema": {
      "type": "string",
      "format": "uri"
    },
    "$ref": {
      "type": "string",
      "format": "uri-reference"
    },
    "$anchor": {
      "type": "string",
      "pattern": "^[A-Za-z][-A-Za-z0-9.:_]*$"
    },
    "$dynamicRef": {
      "type": "string",
      "format": "uri-reference"
    },
    "$dynamicAnchor": {
      "type": "string",
      "pattern": "^[A-Za-z][-A-Za-z0-9.:_]*$"
    },
    "$vocabulary": {
      "type": "object",
      "propertyNames": {
        "type": "string",
        "format": "uri"
      },
      "additionalProperties": {
        "type": "boolean"
      }
    },
    "$comment": {
      "type": "string"
    },
    "title": {
      "type": "string"
    },
    "description": {
      "type": "string"
    },
    "default": true,
    "deprecated": {
      "type": "boolean",
      "default": false
    },
    "readOnly": {
      "type": "boolean",
      "default": false
    },
    "writeOnly": {
      "type": "boolean",
      "default": false
    },
    "examples": {
      "type": "array",
      "items": true
    },
    "multipleOf": {
      "type": "number",
      "exclusiveMinimum": 0
    },
    "maximum": {
      "type": "number"
    },
    "exclusiveMaximum": {
      "type": "number"
    },
    "minimum": {
      "type": "number"
    },
    "exclusiveMinimum": {
      "type": "number"
    },
    "maxLength": {
      "$ref": "#/$defs/nonNegativeInteger"
    },
    "minLength": {
      "$ref": "#/$defs/nonNegativeIntegerDefault0"
    },
    "pattern": {
      "type": "string",
      "format": "regex"
    },
    "maxItems": {
      "$ref": "#/$defs/nonNegativeInteger"
    },
    "minItems": {
      "$ref": "#/$defs/nonNegativeIntegerDefault0"
    },
    "uniqueItems": {
      "type": "boolean",
      "default": false
    },
    "maxContains": {
      "$ref": "#/$defs/nonNegativeInteger"
    },
    "minContains": {
      "$ref": "#/$defs/nonNegativeInteger",
      "default": 1
    },
    "maxProperties": {
      "$ref": "#/$defs/nonNegativeInteger"
    },
    "minProperties": {
      "$ref": "#/$defs/nonNegativeIntegerDefault0"
    },
    "required": {
      "$ref": "#/$defs/stringArray"
    },
    "dependentRequired": {
      "type": "object",
      "additionalProperties": {
        "$ref": "#/$defs/stringArray"
      }
    },
    "const": true,
    "enum": {
      "type": "array",
      "items": true
    },
    "type": {
      "anyOf": [
        {
          "$ref": "#/$defs/simpleTypes"
        },
        {
          "type": "array",
          "items": {
            "$ref": "#/$defs/simpleTypes"
          },
          "minItems": 1,
          "uniqueItems": true
        }
      ]
    },
    "format": {
      "type": "string"
    },
    "contentMediaType": {
      "type": "string"
    },
    "contentEncoding": {
      "type": "string"
    },
    "contentSchema": {
      "$dynamicRef": "#meta"
    },
    "properties": {
      "type": "object",
      "additionalProperties": {
        "$dynamicRef": "#meta"
      },
      "default": {}
    },
    "patternProperties": {
      "type": "object",
      "additionalProperties": {
        "$dynamicRef": "#meta"
      },
      "propertyNames": {
        "format": "regex"
      },
      "default": {}
    },
    "additionalProperties": {
      "$dynamicRef": "#meta"
    },
    "propertyNames": {
      "$dynamicRef": "#meta"
    },
    "unevaluatedProperties": {
      "$dynamicRef": "#meta"
    },
    "items": {
      "$dynamicRef": "#meta"
    },
    "additionalItems": {
      "$dynamicRef": "#meta",
      "deprecated": true
    },
    "contains": {
      "$dynamicRef": "#meta"
    },
    "allOf": {
      "$ref": "#/$defs/schemaArray"
    },
    "anyOf": {
      "$ref": "#/$defs/schemaArray"
    },
    "oneOf": {
      "$ref": "#/$defs/schemaArray"
    },
    "not": {
      "$dynamicRef": "#meta"
    }
  }
}`

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
	issuerPermManagementMode CredentialSchemaPermManagementMode,
	verifierPermManagementMode CredentialSchemaPermManagementMode,
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
	if err := validateJSONSchema(msg.JsonSchema); err != nil {
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

func validateJSONSchema(schemaJSON string) error {
	if schemaJSON == "" {
		return fmt.Errorf("json schema cannot be empty")
	}

	// Parse JSON
	var schemaDoc map[string]interface{}
	if err := json.Unmarshal([]byte(schemaJSON), &schemaDoc); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Load the meta-schema
	metaSchemaLoader := gojsonschema.NewStringLoader(jsonSchemaMetaSchema)
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)

	// Validate against meta-schema
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
	if msg.IssuerPermManagementMode == CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_UNSPECIFIED {
		return fmt.Errorf("issuer permission management mode must be specified")
	}

	// Check verifier permission management mode
	if msg.VerifierPermManagementMode == CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_UNSPECIFIED {
		return fmt.Errorf("verifier permission management mode must be specified")
	}

	// Define valid modes
	validModes := map[CredentialSchemaPermManagementMode]bool{
		CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_OPEN:                      true,
		CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION:        true,
		CredentialSchemaPermManagementMode_PERM_MANAGEMENT_MODE_TRUST_REGISTRY_VALIDATION: true,
	}

	if !validModes[msg.IssuerPermManagementMode] {
		return fmt.Errorf("invalid issuer permission management mode")
	}

	if !validModes[msg.VerifierPermManagementMode] {
		return fmt.Errorf("invalid verifier permission management mode")
	}

	return nil
}
