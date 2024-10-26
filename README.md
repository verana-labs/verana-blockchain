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
