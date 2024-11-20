#!/bin/bash

set -e

# Function to log messages
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Detecting the OS
OS="$(uname)"
log "Detected OS: $OS"

# Function to handle sed compatibility between macOS and Linux
sed_inplace() {
    if [[ "$OS" == "Darwin" ]]; then
        sed -i '' "$1" "$2" # macOS version of sed needs '' for in-place editing
    else
        sed -i "$1" "$2"    # Linux version of sed
    fi
}

# Variables
CHAIN_ID="test-1"
MONIKER="testdev"
BINARY="veranad"
HOME_DIR="$HOME/.verana"
GENESIS_JSON_PATH="$HOME_DIR/config/genesis.json"
APP_TOML_PATH="$HOME_DIR/config/app.toml"
CONFIG_TOML_PATH="$HOME_DIR/config/config.toml"
VALIDATOR_NAME="cooluser"
VALIDATOR_AMOUNT="1000000000000000000000uvna"
GENTX_AMOUNT="1000000000uvna"

log "Starting Verana blockchain setup..."

# Ensure the binary is in the correct location
if [ ! -f "/usr/local/bin/$BINARY" ]; then
    log "Moving $BINARY to /usr/local/bin..."
    sudo mv ~/go/bin/$BINARY /usr/local/bin/
    if [ $? -ne 0 ]; then
        log "Error: Failed to move $BINARY to /usr/local/bin. Please check permissions."
        exit 1
    fi
fi

# Initialize the chain
log "Initializing the chain..."
$BINARY init $MONIKER --chain-id $CHAIN_ID
if [ $? -ne 0 ]; then
    log "Error: Failed to initialize the chain."
    exit 1
fi

# Add a validator key
log "Adding validator key..."
$BINARY keys add $VALIDATOR_NAME --keyring-backend test
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
$BINARY gentx $VALIDATOR_NAME $GENTX_AMOUNT --chain-id $CHAIN_ID --keyring-backend test
if [ $? -ne 0 ]; then
    log "Error: Failed to create genesis transaction."
    exit 1
fi

# Update minimum-gas-prices in app.toml
log "Updating minimum gas prices..."
sed_inplace 's/^minimum-gas-prices = ""/minimum-gas-prices = "0.25uvna"/' "$APP_TOML_PATH"
if [ $? -ne 0 ]; then
    log "Error: Failed to update minimum gas prices in app.toml."
    exit 1
fi

# Replace all occurrences of "stake" with "uvna" in genesis.json
log "Replacing 'stake' with 'uvna' in genesis.json..."
sed_inplace 's/stake/uvna/g' "$GENESIS_JSON_PATH"
if [ $? -ne 0 ]; then
    log "Error: Failed to replace 'stake' with 'uvna' in genesis.json."
    exit 1
fi

# Collect genesis transactions
log "Collecting genesis transactions..."
$BINARY collect-gentxs
if [ $? -ne 0 ]; then
    log "Error: Failed to collect genesis transactions."
    exit 1
fi

# Validate genesis file
log "Validating genesis file..."
$BINARY validate-genesis
if [ $? -ne 0 ]; then
    log "Error: Genesis file validation failed."
    exit 1
fi

# Update app.toml
log "Updating app.toml..."
sed_inplace 's/enable = false/enable = true/' "$APP_TOML_PATH"
sed_inplace 's/swagger = false/swagger = true/' "$APP_TOML_PATH"
sed_inplace 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "$APP_TOML_PATH"
if [ $? -ne 0 ]; then
    log "Error: Failed to update app.toml."
    exit 1
fi

# Update config.toml CORS settings
log "Updating CORS settings in config.toml..."
sed_inplace 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/' "$CONFIG_TOML_PATH"
if [ $? -ne 0 ]; then
    log "Error: Failed to update CORS settings in config.toml."
    exit 1
fi

# Start the chain
log "Starting the Verana blockchain..."
$BINARY start

log "Setup complete. If you encounter any issues, please check the logs above."