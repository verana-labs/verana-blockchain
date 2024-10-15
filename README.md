# Verana Blockchain Trust Registry Module

This README provides instructions for setting up the Verana blockchain and interacting with the Trust Registry module.

## Setting Up the Chain

1. Clone the repository:
   ```
   git clone https://github.com/verana-labs/verana-blockchain.git
   cd verana-blockchain
   ```

2. Run the setup script:
   ```
   ./scripts/setup_verana.sh
   ```

   This script initializes the chain and starts the node.

## Interacting with the Trust Registry Module

### Using CLI

1. Create a Trust Registry:
   ```
   veranad tx trustregistry create-trust-registry \
   did:example:123456789abcdefghi \
   http://example-aka.com \
   en \
   https://example.com/governance-framework.pdf \
   e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 \
   --from cooluser --keyring-backend test \
   --chain-id test-1 \
   --gas auto \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna \
   --keyring-backend test
   ```

2. Add Governance Framework Document:
   ```
   veranad tx trustregistry add-governance-framework-document \
   did:example:123456789abcdefghi \
   en \
   https://example.com/updated-governance-framework.pdf \
   e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 \
   2 \
   --from cooluser --keyring-backend test \
   --chain-id test-1 \
   --gas auto \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna \
   --keyring-backend test
   ```

3. Increase Active Governance Framework Version:
   ```
   veranad tx trustregistry increase-active-gf-version \
   did:example:123456789abcdefghi \
   --from cooluser --keyring-backend test \
   --chain-id test-1 \
   --gas auto \
   --gas-adjustment 1.3 \
   --gas-prices 1.1uvna \
   --keyring-backend test
   ```

4. Query Trust Registry:
   ```
   veranad q trustregistry get-trust-registry did:example:123456789abcdefghi \     
    --active-gf-only \
    --preferred-language en \
    --output json
   ```

### Using curl

For each transaction, you need to first create and sign the transaction, then broadcast it using curl.
