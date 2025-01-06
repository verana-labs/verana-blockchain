# Verana Blockchain Trust Registry Module

This README provides instructions for setting up the Verana blockchain and interacting with the Trust Registry module.

## Setting Up the Chain

1. Clone the repository:
   ```bash
   git clone https://github.com/verana-labs/verana-blockchain.git
   cd verana-blockchain
   ```

2. Run the setup script:
   ```bash
   ./scripts/setup_verana.sh
   ```

   This script initializes the chain and starts the node.

## Interacting with the Trust Registry Module

### Using CLI

1. Create a Trust Registry:
   ```bash
   veranad tx trustregistry create-trust-registry \
   did:example:123456789abcdefghi \
   "http://example-aka.com" \
   en \
   https://example.com/governance-framework.pdf \
   e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

   Note: This operation returns a Trust Registry ID that you'll need for subsequent operations.

2. Add Governance Framework Document:
   ```bash
   veranad tx trustregistry add-governance-framework-document \
   <tr_id> \
   en \
   https://example.com/updated-governance-framework.pdf \
   e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 \
   2 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

   Note: Replace `<tr_id>` with the actual Trust Registry ID.

3. Increase Active Governance Framework Version:
   ```bash
   veranad tx trustregistry increase-active-gf-version \
   <tr_id> \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

### Queries

1. Query Trust Registry by ID:
   ```bash
   veranad q trustregistry get-trust-registry <tr_id> \
   --active-gf-only \
   --preferred-language en \
   --output json
   ```

2. Query Trust Registry by DID:
   ```bash
   veranad q trustregistry get-trust-registry-by-did did:example:123456789abcdefghi \
   --active-gf-only \
   --preferred-language en \
   --output json
   ```

3. List Trust Registries:
   ```bash
   veranad q trustregistry list-trust-registries \
   --controller <account_address> \
   --modified-after "2023-01-01T00:00:00Z" \
   --active-gf-only \
   --preferred-language en \
   --response-max-size 100 \
   --output json
   ```

### Query Parameters

- `active-gf-only`: If true, returns only the current active version's data
- `preferred-language`: Return documents in this language when available
- `modified-after`: Filter registries modified after this timestamp
- `response-max-size`: Limit number of results (1-1024, default 64)

## Running Tests

To run the test suite for the Trust Registry module:

1. Run all tests:
   ```bash
   make test
   ```

2. Run tests with coverage:
   ```bash
   make test-coverage
   ```

   View the coverage report:
   ```bash
   open coverage.html
   ```

Note: Replace `cooluser`, chain ID, gas prices, and other parameters according to your setup.


## Multi-Validator Setup and Testing

### Setting Up Multiple Validators

1. Clean up any existing data:
   ```bash
   rm -rf ~/.verana ~/.verana2
   ```

2. Start Primary Validator (Terminal 1):
   ```bash
   ./scripts/setup_primary_validator.sh
   ```

3. Start Second Validator (Terminal 2):
   ```bash
   ./scripts/setup_additional_validator.sh 2
   ```

### Testing with Multiple Validators

1. Create Trust Registry through Primary Validator:
   ```bash
   # Terminal 1 (Primary Validator)
   veranad tx trustregistry create-trust-registry \
   did:example:123456789abcdefghi \
   "http://example-aka.com" \
   en \
   https://example.com/governance-framework.pdf \
   e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna \
   --home ~/.verana
   ```

2. Create Trust Registry through Secondary Validator:
   ```bash
   # Terminal 2 (Secondary Validator)
   veranad tx trustregistry create-trust-registry \
   did:example:456789abcdefghi \
   "http://example2-aka.com" \
   es \
   https://example2.com/governance-framework.pdf \
   e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna \
   --node tcp://localhost:26757
   ```
   #### Without --node flag, it defaults to 26657 (primary validator's RPC)
   #### To interact with secondary validator, need to specify --node tcp://localhost:26757

### Querying Transactions and Blocks

1. Query Transaction by Height:
   ```bash
   # Can be executed on either validator
   veranad q txs --query "tx.height=57"
   ```

2. Query Validators:
   ```bash
   # Check validator set
   veranad q tendermint-validator-set --home ~/.verana
   ```

## Interacting with the DID Directory Module

### Using CLI

1. Add a DID:
   ```bash
   veranad tx diddirectory add-did \
   did:example:123456789abcdefghi \
   5 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```
   Note: The second parameter (5) is the number of years for DID registration (1-31 years, defaults to 1 if not specified)

2. Renew a DID:
   ```bash
   veranad tx diddirectory renew-did \
   did:example:123456789abcdefghi \
   2 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```
   Note: The second parameter (2) is the number of additional years to extend the registration

3. Remove a DID:
   ```bash
   veranad tx diddirectory remove-did \
   did:example:123456789abcdefghi \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```
   Note: Only the controller can remove before grace period. Anyone can remove after grace period.

4. Touch a DID:
   ```bash
   veranad tx diddirectory touch-did \
   did:example:123456789abcdefghi \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```
   Note: Updates the last modified time to trigger reindexing

### Queries

1. List DIDs:
   ```bash
   veranad q diddirectory list-dids \
   --account <controller_address> \
   --changed "2024-01-01T00:00:00Z" \
   --expired=false \
   --over-grace=false \
   --max-results 64 \
   --output json
   ```

2. Get DID Details:
   ```bash
   veranad q diddirectory get-did \
   did:example:123456789abcdefghi \
   --output json
   ```

### Query Parameters

- `account`: Filter DIDs by controller account address
- `changed-after`: Filter DIDs modified after timestamp
- `expired`: Show expired DIDs
- `over-grace`: Show DIDs that are past grace period
- `max-results`: Maximum number of results to return (1-1024, default 64)

## Interacting with the Credential Schema Module

### Using CLI

1. Create a Credential Schema:
   ```bash
   echo '{
       "$schema": "https://json-schema.org/draft/2020-12/schema",
       "$id": "/dtr/v1/cs/js/1",
       "type": "object",
       "$defs": {},
       "properties": {
           "name": {
               "type": "string"
           }
       },
       "required": ["name"],
       "additionalProperties": false
   }' > schema.json

   veranad tx credentialschema create-credential-schema \
   1 \
   "$(cat schema.json)" \
   365 \
   365 \
   180 \
   180 \
   180 \
   2 \
   2 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

### Queries

1. List Credential Schemas:
   ```bash
   veranad q credentialschema list-schemas \
   --tr_id 1 \
   --created_after "2024-01-01T00:00:00Z" \
   --response_max_size 100 \
   --output json
   ```

2. Get Credential Schema:
   ```bash
   veranad q credentialschema get 1 \
   --output json
   ```

3. Get JSON Schema Definition:
   ```bash
   veranad q credentialschema schema 1 \
   --output json
   ```

### Query Parameters

- `tr_id`: Filter schemas by trust registry ID
- `created_after`: Show schemas created after this datetime (RFC3339 format)
- `response_max_size`: Maximum number of results (1-1024, default 64)

Note:
- The issuer and verifier mode values are:
   - 1: OPEN
   - 2: GRANTOR_VALIDATION
   - 3: TRUST_REGISTRY_VALIDATION
- A trust registry must exist before creating a credential schema
- The schema creator must be the controller of the referenced trust registry

## Interacting with the Credential Schema Permission Module

### Using CLI

1. Create a Credential Schema Permission:
```bash
   veranad tx cspermission create-credential-schema-perm \
   1 \
   1 \
   "did:example:123" \
   verana1mda3hc2z8jnmk86zkvm9wlfgfmxwg2msf2a3ka \
   "2024-03-16T15:00:00Z" \
   100 \
   200 \
   300 \
   --effective-until "2025-03-16T15:00:00Z" \
   --country US \
   --validation-id 123 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
```

The permission types are:
- 1 = ISSUER
- 2 = VERIFIER
- 3 = ISSUER_GRANTOR
- 4 = VERIFIER_GRANTOR
- 5 = TRUST_REGISTRY
- 6 = HOLDER

## Interacting with the Validation Module

### Prerequisites

1. Create and fund validator key:
   ```bash
   # Create validator key
   veranad keys add validator --keyring-backend test

   # Fund validator account (using test_user which should already have funds)
   veranad tx bank send \
   cooluser \
   $(veranad keys show validator -a --keyring-backend test) \
   1000000000uvna \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas auto \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

2. First create a Trust Registry (will need the ID):
   ```bash
   veranad tx trustregistry create-trust-registry \
   did:example:123456789abcdefghi \
   "http://example-aka.com" \
   en \
   https://example.com/governance-framework.pdf \
   e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas auto \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

3. Create a Credential Schema (using Trust Registry ID from step 1):
   ```bash
   echo '{
       "$schema": "https://json-schema.org/draft/2020-12/schema",
       "$id": "/dtr/v1/cs/js/1",
       "type": "object",
       "properties": {
           "name": {
               "type": "string"
           }
       },
       "required": ["name"],
       "additionalProperties": false
   }' > schema.json

   veranad tx credentialschema create-credential-schema \
   1 \
   "$(cat schema.json)" \
   365 \
   365 \
   180 \
   180 \
   180 \
   2 \
   2 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

4. Create necessary Permissions (e.g., ISSUER_GRANTOR permission for validating ISSUER requests):
   ```bash
   veranad tx cspermission create-credential-schema-perm \
   1 \
   3 \
   "did:example:123" \
   $(veranad keys show validator -a --keyring-backend test) \
   "2024-12-29T15:00:00Z" \
   100 \
   200 \
   300 \
   --effective-until "2025-03-16T15:00:00Z" \
   --country US \
   --validation-id 123 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

### Transactions

1. Create a Validation:
   ```bash
   veranad tx validation create-validation \
   3 \
   1 \
   US \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

2. Renew a Validation:
   ```bash
   veranad tx validation renew-validation \
   1 \
   1 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

   Note: Starting a renewal process with a different validator implicitly transfers revocation control of existing permissions to the new validator.

3. Set Validation to Validated:
   ```bash
   veranad tx validation set-validated \
   1 \
   e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 \
   --from validator \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

   Note:
   - Only the validator can set a validation to VALIDATED state
   - The summary hash parameter is optional and must be null for HOLDER type validations

4. Request Validation Termination:
   ```bash
   veranad tx validation request-termination \
   1 \
   --from cooluser \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```

   Note:
   - Only the validation applicant can request termination
   - Validation must be in VALIDATED state
   - After requesting termination, the validation will be set to TERMINATION_REQUESTED state
   - A separate confirmation step will be required to complete the termination

5. Confirm Validation Termination:
   ```bash
   veranad tx validation confirm-termination 1 \
   --from validator \
   --keyring-backend test \
   --chain-id test-1 \
   --gas 800000 \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna
   ```   

### Queries

1. List Validations:
   ```bash
   veranad q validation list-validations \
   --controller $(veranad keys show cooluser -a --keyring-backend test) \
   --validator-perm-id 1 \
   --type 3 \
   --state 4 \
   --response-max-size 64 \
   --output json
   ```

2. Get Validation by ID:
   ```bash
   veranad q validation get-validation <validation_id> \
   --output json
   ```

### Query Parameters

- `controller`: Filter by controller account address
- `validator-perm-id`: Filter by validator permission ID
- `type`: Filter by validation type (ISSUER_GRANTOR, VERIFIER_GRANTOR, ISSUER, VERIFIER, HOLDER)
- `state`: Filter by validation state (PENDING, VALIDATED, TERMINATED)
- `response-max-size`: Maximum number of results (1-1024, default 64)
- `exp-before`: Filter validations expiring before timestamp (RFC3339 format)

### Validation Types
- `ISSUER`: For becoming an issuer of credentials
- `VERIFIER`: For becoming a verifier of credentials
- `ISSUER_GRANTOR`: For becoming an issuer grantor
- `VERIFIER_GRANTOR`: For becoming a verifier grantor
- `HOLDER`: For getting issued a credential

### Permission Types
- 1 = ISSUER
- 2 = VERIFIER
- 3 = ISSUER_GRANTOR
- 4 = VERIFIER_GRANTOR
- 5 = TRUST_REGISTRY
- 6 = HOLDER

Note: The validation flow depends on the credential schema's permission management modes. For example, to become an ISSUER when the schema's issuer_perm_management_mode is GRANTOR_VALIDATION, you need to get validated by an ISSUER_GRANTOR.