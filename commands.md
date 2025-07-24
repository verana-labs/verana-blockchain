# Trust Registry Module CLI Commands

## Module Overview

```bash
veranad tx tr            
Transactions commands for the tr module

Usage:
  veranad tx tr [flags]
  veranad tx tr [command]

Available Commands:
  add-governance-framework-document Add a governance framework document
  archive-trust-registry            Archive or unarchive a trust registry
  create-trust-registry             Create a new trust registry
  increase-active-gf-version        Increase the active governance framework version
  update-params                     Execute the UpdateParams RPC method
  update-trust-registry             Update a trust registry
```

## Transaction Commands

### 1. Create Trust Registry

Creates a new trust registry with governance framework documents.

**Syntax:**
```bash
veranad tx tr create-trust-registry <did> <language> <doc-url> <doc-digest-sri> [aka] --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<did>`: Decentralized Identifier (DID) - must follow DID specification
- `<language>`: ISO 639-1 language code (e.g., en, fr, es)
- `<doc-url>`: URL to the governance framework document
- `<doc-digest-sri>`: SHA-384 hash with SRI format prefix
- `[aka]`: Optional - Also Known As URL

**Examples:**

Basic creation:
```bash
veranad tx tr create-trust-registry did:example:123456789abcdefghi en https://example.com/doc sha384-MzNNbQTWCSUSi0bbz7dbua+RcENv7C6FvlmYJ1Y+I727HsPOHdzwELMYO9Mz68M26 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

With AKA (Also Known As):
```bash
veranad tx tr create-trust-registry did:example:123456789abcdefghi en https://example.com/doc sha384-MzNNbQTWCSUSi0bbz7dbua+RcENv7C6FvlmYJ1Y+I727HsPOHdzwELMYO9Mz68M26 --aka http://example.com --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 2. Add Governance Framework Document

Adds a new governance framework document to an existing trust registry.

**Syntax:**
```bash
veranad tx tr add-governance-framework-document <trust-registry-id> <doc-language> <doc-url> <doc-digest-sri> <version> --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<trust-registry-id>`: Numeric ID of the trust registry
- `<doc-language>`: ISO 639-1 language code
- `<doc-url>`: URL to the governance framework document
- `<doc-digest-sri>`: SHA-384 hash with SRI format prefix
- `<version>`: Version number (must be sequential)

**Examples:**

Add document for next version:
```bash
veranad tx tr add-governance-framework-document 1 en https://example.com/doc2 sha384-MzNNbQTWCSUSi0bbz7dbua+RcENv7C6FvlmYJ1Y+I727HsPOHdzwELMYO9Mz68M26 2 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Add document in different language for same version:
```bash
veranad tx tr add-governance-framework-document 1 fr https://example.com/doc2-fr sha384-MzNNbQTWCSUSi0bbz7dbua+RcENv7C6FvlmYJ1Y+I727HsPOHdzwELMYO9Mz68M26 2 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Add document for version 3:
```bash
veranad tx tr add-governance-framework-document 1 es https://example.com/doc3-es sha384-MzNNbQTWCSUSi0bbz7dbua+RcENv7C6FvlmYJ1Y+I727HsPOHdzwELMYO9Mz68M26 3 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 3. Increase Active Governance Framework Version

Increases the active version of the governance framework. Requires that a document exists in the default language for the target version.

**Syntax:**
```bash
veranad tx tr increase-active-gf-version <trust-registry-id> --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<trust-registry-id>`: Numeric ID of the trust registry

**Example:**
```bash
veranad tx tr increase-active-gf-version 1 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

**Note:** This command will fail if there's no document in the default language for the next version.

---

### 4. Update Trust Registry

Updates the DID and/or AKA fields of an existing trust registry.

**Syntax:**
```bash
veranad tx tr update-trust-registry <trust-registry-id> <new-did> [new-aka] --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<trust-registry-id>`: Numeric ID of the trust registry
- `<new-did>`: New DID for the trust registry
- `[new-aka]`: Optional - New AKA URL (use empty string to clear)

**Examples:**

Update DID and AKA:
```bash
veranad tx tr update-trust-registry 1 did:example:newdid --aka http://new.example.com --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 5. Archive Trust Registry

Archives or unarchives a trust registry.

**Syntax:**
```bash
veranad tx tr archive-trust-registry <trust-registry-id> <archive-flag> --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<trust-registry-id>`: Numeric ID of the trust registry
- `<archive-flag>`: Boolean value (`true` to archive, `false` to unarchive)

**Examples:**

Archive a trust registry:
```bash
veranad tx tr archive-trust-registry 1 true --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Unarchive a trust registry:
```bash
veranad tx tr archive-trust-registry 1 false --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---


## Parameter Validation Rules

### DID Format
- Must follow DID specification
- Example: `did:example:123456789abcdefghi`

### Language Codes
- Must be valid ISO 639-1 language codes
- Examples: `en`, `fr`, `es`, `de`, `zh`

### Document Digest SRI
- Must use SHA-384 hash with SRI format
- Format: `sha384-<base64-encoded-hash>`
- Example: `sha384-MzNNbQTWCSUSi0bbz7dbua+RcENv7C6FvlmYJ1Y+I727HsPOHdzwELMYO9Mz68M26`

### Version Rules
- Versions must be sequential
- Cannot skip versions when adding documents
- Must have default language document before increasing active version

### Access Control
- Only the trust registry controller (creator) can perform updates
- Wrong controller will result in transaction failure

## Common Error Scenarios

1. **Invalid Version**: Attempting to add a document with a version that skips numbers
2. **Missing Default Language**: Trying to increase version without default language document
3. **Wrong Controller**: Non-controller attempting to modify trust registry
4. **Already Archived/Unarchived**: Attempting to archive an already archived registry
5. **Invalid Language Format**: Using non-ISO 639-1 language codes
6. **Non-existent Trust Registry**: Using invalid trust registry ID

## Transaction Fees

All transactions require gas fees. Use `--gas auto` for automatic gas estimation or specify a specific gas limit. Fee examples:
- `--fees 50000uvna` (50,000 micro-VNA)
- `--gas auto`


######################################################################################################################################

# Credential Schema Module CLI Commands

This document provides comprehensive CLI commands for the Credential Schema (cs) module in the Verana blockchain.

## Module Overview

```bash
veranad tx cs
Transactions commands for the cs module

Usage:
  veranad tx cs [flags]
  veranad tx cs [command]

Available Commands:
  archive                  Archive or unarchive a credential schema
  create-credential-schema Create a new credential schema
  update                   Update a credential schema's validity periods
```

## Transaction Commands

### 1. Create Credential Schema

Creates a new credential schema linked to a trust registry.

**Syntax:**
```bash
veranad tx cs create-credential-schema <trust-registry-id> <json-schema> <issuer-grantor-validity> <verifier-grantor-validity> <issuer-validity> <verifier-validity> <holder-validity> <issuer-perm-mode> <verifier-perm-mode> --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<trust-registry-id>`: Numeric ID of the trust registry (must exist and caller must be controller)
- `<json-schema>`: JSON schema definition (properly escaped JSON string)
- `<issuer-grantor-validity>`: Issuer grantor validation validity period in days
- `<verifier-grantor-validity>`: Verifier grantor validation validity period in days
- `<issuer-validity>`: Issuer validation validity period in days
- `<verifier-validity>`: Verifier validation validity period in days
- `<holder-validity>`: Holder validation validity period in days
- `<issuer-perm-mode>`: Issuer permission management mode (integer)
- `<verifier-perm-mode>`: Verifier permission management mode (integer)

**Example:**

Basic credential schema creation:
```bash
veranad tx cs create-credential-schema 1 '{"$schema":"https://json-schema.org/draft/2020-12/schema","$id":"/vpr/v1/cs/js/1","type":"object","properties":{"name":{"type":"string"}},"required":["name"],"additionalProperties":false}' 365 365 180 180 180 2 2 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

**Note:** The JSON schema must be properly escaped when passed as a command line argument. For complex schemas, consider using a file:

```bash
# Save schema to file first
cat > schema.json << 'EOF'
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "$id": "/vpr/v1/cs/js/1",
    "type": "object",
    "$defs": {},
    "properties": {
        "name": {
            "type": "string"
        },
        "email": {
            "type": "string",
            "format": "email"
        }
    },
    "required": ["name"],
    "additionalProperties": false
}
EOF

# Use in command (you'll need to escape or quote properly)
veranad tx cs create-credential-schema 1 "$(cat schema.json)" 365 365 180 180 180 2 2 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 2. Update Credential Schema

Updates the validity periods of an existing credential schema.

**Syntax:**
```bash
veranad tx cs update <credential-schema-id> <issuer-grantor-validity> <verifier-grantor-validity> <issuer-validity> <verifier-validity> <holder-validity> --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<credential-schema-id>`: Numeric ID of the credential schema
- `<issuer-grantor-validity>`: New issuer grantor validation validity period in days
- `<verifier-grantor-validity>`: New verifier grantor validation validity period in days
- `<issuer-validity>`: New issuer validation validity period in days
- `<verifier-validity>`: New verifier validation validity period in days
- `<holder-validity>`: New holder validation validity period in days

**Examples:**

Update validity periods:
```bash
veranad tx cs update 1 365 365 180 180 180 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Extend all validity periods:
```bash
veranad tx cs update 1 730 730 365 365 365 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Shorter validity periods for testing:
```bash
veranad tx cs update 1 30 30 7 7 7 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 3. Archive Credential Schema

Archives or unarchives a credential schema.

**Syntax:**
```bash
veranad tx cs archive <credential-schema-id> <archive-flag> --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<credential-schema-id>`: Numeric ID of the credential schema
- `<archive-flag>`: Boolean value (`true` to archive, `false` to unarchive)

**Examples:**

Archive a credential schema:
```bash
veranad tx cs archive 1 true --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Unarchive a credential schema:
```bash
veranad tx cs archive 1 false --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

## Parameter Details

### JSON Schema Requirements
- Must be valid JSON Schema (Draft 2020-12 recommended)
- Should include `$schema`, `$id`, `type`, and `properties` fields
- Must define required fields and additionalProperties behavior
- Common pattern: `{"$schema": "https://json-schema.org/draft/2020-12/schema", "$id": "/path/to/schema", "type": "object", "properties": {...}, "required": [...], "additionalProperties": false}`

### Validity Periods
- Measured in days
- Must not exceed system maximum limits
- Common values:
    - **Grantor periods**: 365-730 days (1-2 years)
    - **Validation periods**: 30-365 days (1 month to 1 year)
    - **Testing periods**: 1-7 days for development

### Permission Management Modes
- Integer values representing different permission models
- Mode `2` appears to be a standard mode in examples
- Specific mode meanings depend on system configuration

### Trust Registry Requirements
- Credential schema must be linked to an existing trust registry
- Only the trust registry controller can create/modify credential schemas
- Trust registry must not be archived

## Validation Rules

### Access Control
- **Create**: Only trust registry controller can create schemas
- **Update**: Only trust registry controller can update schemas
- **Archive**: Only trust registry controller can archive/unarchive schemas

### Business Logic
- Cannot archive an already archived schema
- Cannot unarchive a schema that's not archived
- Validity periods cannot exceed system maximums
- JSON schema must be valid and parseable

## Transaction Fees

All transactions require gas fees. Use `--gas auto` for automatic gas estimation:
- `--fees 50000uvna` (50,000 micro-VNA)
- `--gas auto`

For complex JSON schemas, gas consumption may be higher due to storage requirements.



######################################################################################################################################

# Permission Module CLI Commands

This document provides comprehensive CLI commands for the Permission (perm) module in the Verana blockchain.

## Module Overview

```bash
veranad tx perm
Transactions commands for the perm module

Usage:
  veranad tx perm [flags]
  veranad tx perm [command]

Available Commands:
  cancel-perm-vp-request           Cancel a pending perm VP request
  confirm-vp-termination          Confirm the termination of a perm VP
  create-or-update-perm-session   Create or update a perm session
  create-perm                     Create a new perm for open schemas
  create-root-perm                Create a new root perm for a credential schema
  extend-perm                     Extend a perm's effective duration
  renew-perm-vp                   Renew a perm validation process
  repay-perm-slashed-td           Repay a slashed perm's trust deposit
  request-vp-termination          Request termination of a perm validation process
  revoke-perm                     Revoke a perm
  set-perm-vp-validated           Set perm validation process to validated state
  slash-perm-td                   Slash a perm's trust deposit
  start-perm-vp                   Start a new perm validation process
```

## Transaction Commands

### 1. Start Permission Validation Process

Initiates a new permission validation process for credential schema permissions requiring grantor validation.

**Syntax:**
```bash
veranad tx perm start-perm-vp <permission-type> <validator-perm-id> <country> <did> --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<permission-type>`: Integer representing permission type (1=ISSUER, 2=VERIFIER, 4=HOLDER)
- `<validator-perm-id>`: ID of the validator permission that will validate this request
- `<country>`: ISO 3166-1 alpha-2 country code (e.g., "US", "FR", "DE")
- `<did>`: Decentralized identifier of the requestor

**Examples:**

Start ISSUER permission validation:
```bash
veranad tx perm start-perm-vp 1 3 "US" --did "did:example:123456789abcdefghi" --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Start VERIFIER permission validation:
```bash
veranad tx perm start-perm-vp 2 3 "FR" "did:example:987654321abcdefghi" --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 2. Set Permission VP to Validated

Sets a pending permission validation process to validated state. Only the designated validator can execute this command.

**Syntax:**
```bash
veranad tx perm set-perm-vp-validated <permission-id> <validation-fees> <issuance-fees> <verification-fees> <country> --effective-until <timestamp> --vp-summary-digest-sri <digest> --from <validator> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<permission-id>`: ID of the permission to validate
- `<validation-fees>`: Validation fees amount in smallest unit
- `<issuance-fees>`: Issuance fees amount in smallest unit
- `<verification-fees>`: Verification fees amount in smallest unit
- `<country>`: ISO 3166-1 alpha-2 country code

**Flags:**
- `--effective-until`: Effective until timestamp in RFC3339 format
- `--vp-summary-digest-sri`: VP summary digest SRI (required for non-HOLDER types)

**Examples:**

Validate an ISSUER permission:
```bash
veranad tx perm set-perm-vp-validated 456 10 5 3 "US" --effective-until "2024-12-31T23:59:59Z" --vp-summary-digest-sri "sha384-validDigest123" --from validator --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Validate a HOLDER permission (no digest required):
```bash
veranad tx perm set-perm-vp-validated 789 10 5 3 "US" --effective-until "2024-12-31T23:59:59Z" --from validator --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 3. Create Root Permission

Creates a new root (ECOSYSTEM) permission for a credential schema. Only the trust registry controller can execute this command.

**Syntax:**
```bash
veranad tx perm create-root-perm <schema-id> <did> <validation-fees> <issuance-fees> <verification-fees> <country> --effective-from <timestamp> --effective-until <timestamp> --from <controller> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<schema-id>`: Credential schema ID
- `<did>`: Decentralized identifier of the trust registry controller
- `<validation-fees>`: Validation fees amount in smallest unit
- `<issuance-fees>`: Issuance fees amount in smallest unit
- `<verification-fees>`: Verification fees amount in smallest unit
- `<country>`: ISO 3166-1 alpha-2 country code

**Flags:**
- `--effective-from`: Effective from timestamp in RFC3339 format (optional)
- `--effective-until`: Effective until timestamp in RFC3339 format (optional)

**Examples:**

Create root permission with time bounds:
```bash
veranad tx perm create-root-perm 1 "did:example:trustregistry123" 100 50 25 --country "US" --effective-from "2026-01-01T00:00:00Z" --effective-until "2027-01-01T00:00:00Z" --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Create root permission without time bounds:
```bash
veranad tx perm create-root-perm 1 "did:example:trustregistry123" 100 50 25 --country "US" --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 4. Create Permission (Open Schemas)

Creates a new permission for credential schemas with OPEN permission management mode.

**Syntax:**
```bash
veranad tx perm create-perm <schema-id> <permission-type> <did> <country> <verification-fees> --effective-from <timestamp> --effective-until <timestamp> --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<schema-id>`: Credential schema ID (must have OPEN permission mode)
- `<permission-type>`: Permission type (1=ISSUER, 2=VERIFIER)
- `<did>`: Decentralized identifier
- `<country>`: ISO 3166-1 alpha-2 country code
- `<verification-fees>`: Verification fees amount in smallest unit

**Flags:**
- `--effective-from`: Effective from timestamp in RFC3339 format (optional)
- `--effective-until`: Effective until timestamp in RFC3339 format (optional)

**Examples:**

Create ISSUER permission for open schema:
```bash
veranad tx perm create-perm 1 1 "did:example:issuer123456789" "US" 100 --effective-from "2024-01-01T00:00:00Z" --effective-until "2025-01-01T00:00:00Z" --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Create VERIFIER permission for open schema:
```bash
veranad tx perm create-perm 1 2 "did:example:verifier987654321" "FR" 50 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 5. Renew Permission Validation Process

Renews an existing permission validation process, setting it back to PENDING state.

**Syntax:**
```bash
veranad tx perm renew-perm-vp <permission-id> --from <grantee> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<permission-id>`: ID of the permission to renew

**Examples:**

Renew permission validation:
```bash
veranad tx perm renew-perm-vp 456 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 6. Request VP Termination

Requests termination of a permission validation process.

**Syntax:**
```bash
veranad tx perm request-vp-termination <permission-id> --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<permission-id>`: ID of the permission to terminate

**Examples:**

Request termination (by grantee):
```bash
veranad tx perm request-vp-termination 456 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Request termination of expired VP (by validator):
```bash
veranad tx perm request-vp-termination 789 --from validator --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 7. Confirm VP Termination

Confirms the termination of a permission validation process.

**Syntax:**
```bash
veranad tx perm confirm-vp-termination <permission-id> --from <user> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<permission-id>`: ID of the permission to confirm termination

**Examples:**

Confirm termination by validator (before timeout):
```bash
veranad tx perm confirm-vp-termination 456 --from validator --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Confirm termination by applicant (after timeout):
```bash
veranad tx perm confirm-vp-termination 789 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 8. Cancel Permission VP Request

Cancels a pending permission VP request, returning fees and deposits.

**Syntax:**
```bash
veranad tx perm cancel-perm-vp-request <permission-id> --from <grantee> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<permission-id>`: ID of the pending permission request to cancel

**Examples:**

Cancel pending VP request:
```bash
veranad tx perm cancel-perm-vp-request 456 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 9. Extend Permission

Extends a permission's effective duration.

**Syntax:**
```bash
veranad tx perm extend-perm <permission-id> <effective-until> --from <validator> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<permission-id>`: ID of the permission to extend
- `<effective-until>`: New effective until timestamp in RFC3339 format

**Examples:**

Extend permission by validator:
```bash
veranad tx perm extend-perm 456 "2025-12-31T23:59:59Z" --from validator --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Extend ecosystem permission by trust registry controller:
```bash
veranad tx perm extend-perm 123 "2026-01-01T00:00:00Z" --from trustregistry --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 10. Revoke Permission

Revokes a permission permanently.

**Syntax:**
```bash
veranad tx perm revoke-perm <permission-id> --from <validator> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<permission-id>`: ID of the permission to revoke

**Examples:**

Revoke permission:
```bash
veranad tx perm revoke-perm 456 --from validator --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 11. Create or Update Permission Session

Creates or updates a permission session for coordinating permissions.

**Syntax:**
```bash
veranad tx perm create-or-update-perm-session <session-id> <agent-perm-id> <wallet-agent-perm-id> --issuer-perm-id <id> --verifier-perm-id <id> --from <controller> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<session-id>`: UUID for the session
- `<agent-perm-id>`: ID of the agent permission (must be HOLDER type)
- `<wallet-agent-perm-id>`: ID of the wallet agent permission (must be HOLDER type)

**Flags:**
- `--issuer-perm-id`: ID of the issuer permission (optional)
- `--verifier-perm-id`: ID of the verifier permission (optional)

**Examples:**

Create session with issuer:
```bash
veranad tx perm create-or-update-perm-session "550e8400-e29b-41d4-a716-446655440000" 789 890 --issuer-perm-id 123 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Create session with verifier:
```bash
veranad tx perm create-or-update-perm-session "550e8400-e29b-41d4-a716-446655440001" 789 890 --verifier-perm-id 234 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Create session with both issuer and verifier:
```bash
veranad tx perm create-or-update-perm-session "550e8400-e29b-41d4-a716-446655440002" 789 890 --issuer-perm-id 123 --verifier-perm-id 234 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 12. Slash Permission Trust Deposit

Slashes a permission's trust deposit as a penalty.

**Syntax:**
```bash
veranad tx perm slash-perm-td <permission-id> <amount> --from <authority> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<permission-id>`: ID of the permission to slash
- `<amount>`: Amount to slash from the trust deposit

**Examples:**

Slash deposit by validator:
```bash
veranad tx perm slash-perm-td 456 500 --from validator --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

Slash deposit by ecosystem controller:
```bash
veranad tx perm slash-perm-td 456 300 --from ecosystem --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

---

### 13. Repay Permission Slashed Trust Deposit

Repays a previously slashed trust deposit.

**Syntax:**
```bash
veranad tx perm repay-perm-slashed-td <permission-id> --from <grantee> --chain-id <chain-id> --keyring-backend test --fees <amount> --gas auto
```

**Parameters:**
- `<permission-id>`: ID of the permission with slashed deposit to repay

**Examples:**

Repay slashed deposit:
```bash
veranad tx perm repay-perm-slashed-td 456 --from cooluser --chain-id vna-testnet-1 --keyring-backend test --fees 50000uvna --gas auto
```

## Query Commands

### Get Permission
```bash
veranad query perm get-permission <permission-id>
```

### List Permissions
```bash
veranad query perm list-permissions --response-max-size <count>
```

### Find Permissions with DID
```bash
veranad query perm find-permissions-with-did <did> <type> <schema-id> <country>
```

### Get Permission Session
```bash
veranad query perm get-permission-session <session-id>
```

### List Permission Sessions
```bash
veranad query perm list-permission-sessions --response-max-size <count>
```

### Find Beneficiaries
```bash
veranad query perm find-beneficiaries --issuer-perm-id <id> --verifier-perm-id <id>
```

## Parameter Details

### Permission Types
- **1**: ISSUER - Can issue credentials
- **2**: VERIFIER - Can verify credentials
- **3**: ISSUER_GRANTOR - Can validate issuer permission requests
- **4**: HOLDER - Can hold credentials
- **5**: VERIFIER_GRANTOR - Can validate verifier permission requests
- **6**: ECOSYSTEM - Root permission for trust registry controllers

### Validation States
- **1**: PENDING - Awaiting validation
- **2**: VALIDATED - Active and validated
- **3**: TERMINATED - Permanently terminated
- **4**: TERMINATION_REQUESTED - Termination requested, awaiting confirmation

### Country Codes
- Must be valid ISO 3166-1 alpha-2 codes
- Examples: "US", "FR", "DE", "GB", "JP", "AU"
- Case sensitive (uppercase required)

### DID Format
- Must follow W3C DID specification
- Common pattern: `did:method:identifier`
- Examples:
  - `did:example:123456789abcdefghi`
  - `did:key:z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK`

### Fee Amounts
- Expressed in smallest unit (e.g., uvna for micro-VNA)
- Common ranges:
  - **Validation fees**: 10-100 uvna
  - **Issuance fees**: 5-50 uvna
  - **Verification fees**: 3-30 uvna

### Session IDs
- Must be valid UUIDs (version 4 recommended)
- Format: `xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx`
- Generate using: `uuidgen` (Linux/macOS) or online UUID generators

## Validation Rules

### Access Control
- **Start VP**: Any user with valid validator permission
- **Set Validated**: Only designated validator can validate
- **Create Root**: Only trust registry controller
- **Create Open**: Any user for schemas with OPEN mode
- **Renew VP**: Only permission grantee
- **Terminate**: Grantee or validator (depending on conditions)
- **Extend**: Only validator or ecosystem controller
- **Revoke**: Only validator
- **Slash**: Only validator or ecosystem controller
- **Repay**: Only permission grantee

### Business Logic
- Validator permission must exist and be validated
- Country codes must match between permission and validator (in some cases)
- Effective dates must be logical (until > from)
- Cannot extend beyond validation expiration
- Cannot slash more than available deposit
- Sessions require HOLDER type permissions for agents

### Permission Dependencies
- ISSUER/VERIFIER require ISSUER_GRANTOR/VERIFIER_GRANTOR validators
- HOLDER permissions require ISSUER validators
- Validator permissions must be VALIDATED state
- Permission sessions require all referenced permissions to exist and be active

## Transaction Fees

All transactions require gas fees. Use `--gas auto` for automatic gas estimation:
- Standard fee: `--fees 50000uvna` (50,000 micro-VNA)
- Complex operations may require higher fees
- Use `--gas-prices 0.025uvna` for manual gas price setting

For operations involving large data (e.g., sessions with many permissions), consider increasing gas limits manually.