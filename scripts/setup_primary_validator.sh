set -e

# Function to log messages
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Detect if running on macOS or Linux
if [[ "$OSTYPE" == "darwin"* ]]; then
  SED_CMD="sed -i ''"
else
  SED_CMD="sed -i"
fi

# Variables
CHAIN_ID="vna-testnet-1"
MONIKER="validator1"
BINARY="veranad"
HOME_DIR="$HOME/.verana"
GENESIS_JSON_PATH="$HOME_DIR/config/genesis.json"
APP_TOML_PATH="$HOME_DIR/config/app.toml"
CONFIG_TOML_PATH="$HOME_DIR/config/config.toml"
VALIDATOR_NAME="cooluser"
VALIDATOR_AMOUNT="1000000000000000000000uvna"
GENTX_AMOUNT="1000000000uvna"
SEED_PHRASE_COOLUSER="pink glory help gown abstract eight nice crazy forward ketchup skill cheese"

# Default ports for primary validator
P2P_PORT="26656"
RPC_PORT="26657"
API_PORT="1317"
GRPC_PORT="9090"
GRPC_WEB_PORT="9091"

log "Starting Primary Validator setup..."

# Initialize the chain
log "Initializing the chain..."
$BINARY init $MONIKER --chain-id $CHAIN_ID
if [ $? -ne 0 ]; then
    log "Error: Failed to initialize the chain."
    exit 1
fi

# Add a validator key
log "Adding validator key..."
echo "$SEED_PHRASE_COOLUSER" | $BINARY keys add $VALIDATOR_NAME --recover --keyring-backend test
if [ $? -ne 0 ]; then
    log "Error: Failed to add validator key."
    exit 1
fi

# Add genesis account
log "Adding genesis account..."
$BINARY add-genesis-account $VALIDATOR_NAME $VALIDATOR_AMOUNT --keyring-backend test
if [ $? -ne 0 ]; then
    log "Error: Failed to add genesis account."
    exit 1
fi

# Create gentx
log "Creating genesis transaction..."
$BINARY gentx $VALIDATOR_NAME $GENTX_AMOUNT \
    --chain-id $CHAIN_ID \
    --moniker $MONIKER \
    --commission-rate "0.10" \
    --commission-max-rate "0.20" \
    --commission-max-change-rate "0.01" \
    --min-self-delegation "1" \
    --keyring-backend test
if [ $? -ne 0 ]; then
    log "Error: Failed to create genesis transaction."
    exit 1
fi

# Update minimum-gas-prices in app.toml
log "Updating minimum gas prices..."
$SED_CMD 's/^minimum-gas-prices = ""/minimum-gas-prices = "0.25uvna"/' "$APP_TOML_PATH"

# Configure ports in app.toml
$SED_CMD "s/:1317/:$API_PORT/" "$APP_TOML_PATH"
$SED_CMD "s/:9090/:$GRPC_PORT/" "$APP_TOML_PATH"
$SED_CMD "s/:9091/:$GRPC_WEB_PORT/" "$APP_TOML_PATH"

# Replace all occurrences of "stake" with "uvna" in genesis.json
log "Replacing 'stake' with 'uvna' in genesis.json..."
$SED_CMD 's/stake/uvna/g' "$GENESIS_JSON_PATH"

# Update governance params in genesis.json
log "Updating governance parameters in genesis.json..."
$SED_CMD 's/"max_deposit_period": ".*"/"max_deposit_period": "100s"/' "$GENESIS_JSON_PATH"
$SED_CMD 's/"voting_period": ".*"/"voting_period": "100s"/' "$GENESIS_JSON_PATH"
$SED_CMD 's/"expedited_voting_period": ".*"/"expedited_voting_period": "90s"/' "$GENESIS_JSON_PATH"
if [ $? -ne 0 ]; then
    log "Error: Failed to update governance parameters in genesis.json."
    exit 1
fi

# Configure ports in config.toml
$SED_CMD "s/:26656/:$P2P_PORT/" "$CONFIG_TOML_PATH"
$SED_CMD "s/:26657/:$RPC_PORT/" "$CONFIG_TOML_PATH"

# Enable API and CORS
log "Updating API and CORS settings..."
$SED_CMD 's/enable = false/enable = true/' "$APP_TOML_PATH"
$SED_CMD 's/swagger = false/swagger = true/' "$APP_TOML_PATH"
$SED_CMD 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "$APP_TOML_PATH"
$SED_CMD 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/' "$CONFIG_TOML_PATH"

# Collect genesis transactions
log "Collecting genesis transactions..."
$BINARY collect-gentxs

# Validate genesis file
log "Validating genesis file..."
$BINARY validate-genesis

# Save the genesis file for other validators
cp $GENESIS_JSON_PATH $HOME/genesis.json
log "Genesis file saved to $HOME/genesis.json"

# Get validator node ID
NODE_ID=$($BINARY tendermint show-node-id)
log "Primary Validator Node ID: $NODE_ID"
echo "Node ID: $NODE_ID" > $HOME/primary_validator_info.txt
echo "P2P Address: localhost:$P2P_PORT" >> $HOME/primary_validator_info.txt

# Start the chain
log "Starting the Primary Validator..."
$BINARY start

log "Primary Validator setup complete. If you encounter any issues, please check the logs above."
