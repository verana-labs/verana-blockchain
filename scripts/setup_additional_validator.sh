set -e

# Function to log messages
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Check if validator number is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <validator-number> (e.g., 2 for second validator)"
    exit 1
fi

VALIDATOR_NUM=$1

# Variables
CHAIN_ID="test-1"
MONIKER="validator$VALIDATOR_NUM"
BINARY="veranad"
HOME_DIR="$HOME/.verana$VALIDATOR_NUM"
APP_TOML_PATH="$HOME_DIR/config/app.toml"
CONFIG_TOML_PATH="$HOME_DIR/config/config.toml"
VALIDATOR_NAME="validator$VALIDATOR_NUM"
STAKE_AMOUNT="1000000000uvna"  # Equal voting power with first validator

# Calculate ports (increment from default ports)
P2P_PORT=$((26656 + ($VALIDATOR_NUM - 1) * 100))
RPC_PORT=$((26657 + ($VALIDATOR_NUM - 1) * 100))
API_PORT=$((1317 + ($VALIDATOR_NUM - 1) * 100))
GRPC_PORT=$((9090 + ($VALIDATOR_NUM - 1) * 100))
GRPC_WEB_PORT=$((9091 + ($VALIDATOR_NUM - 1) * 100))

log "Starting Validator $VALIDATOR_NUM setup..."

# Initialize the chain with custom home
log "Initializing the chain..."
$BINARY init $MONIKER --chain-id $CHAIN_ID --home $HOME_DIR

# Add a validator key
log "Adding validator key..."
$BINARY keys add $VALIDATOR_NAME --keyring-backend test --home $HOME_DIR

# Store the validator address
VALIDATOR_ADDRESS=$($BINARY keys show $VALIDATOR_NAME -a --keyring-backend test --home $HOME_DIR)
log "Validator address: $VALIDATOR_ADDRESS"

# Copy genesis file from primary validator
log "Copying genesis file..."
cp $HOME/genesis.json $HOME_DIR/config/genesis.json

# Update minimum-gas-prices in app.toml
log "Updating minimum gas prices..."
sed -i '' 's/^minimum-gas-prices = ""/minimum-gas-prices = "0.25uvna"/' "$APP_TOML_PATH"

# Configure ports in app.toml
sed -i '' "s/:1317/:$API_PORT/" "$APP_TOML_PATH"
sed -i '' "s/:9090/:$GRPC_PORT/" "$APP_TOML_PATH"
sed -i '' "s/:9091/:$GRPC_WEB_PORT/" "$APP_TOML_PATH"

# Configure ports in config.toml
sed -i '' "s/:26656/:$P2P_PORT/" "$CONFIG_TOML_PATH"
sed -i '' "s/:26657/:$RPC_PORT/" "$CONFIG_TOML_PATH"

# Enable API and CORS
log "Updating API and CORS settings..."
sed -i '' 's/enable = false/enable = true/' "$APP_TOML_PATH"
sed -i '' 's/swagger = false/swagger = true/' "$APP_TOML_PATH"
sed -i '' 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' "$APP_TOML_PATH"
sed -i '' 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/' "$CONFIG_TOML_PATH"

# Read primary validator info
PRIMARY_NODE_ID=$(grep "Node ID:" $HOME/primary_validator_info.txt | cut -d' ' -f3)
PRIMARY_P2P_ADDR=$(grep "P2P Address:" $HOME/primary_validator_info.txt | cut -d' ' -f3)

# Configure persistent peers
log "Configuring persistent peers..."
sed -i '' "s/persistent_peers = \"\"/persistent_peers = \"$PRIMARY_NODE_ID@$PRIMARY_P2P_ADDR\"/" "$CONFIG_TOML_PATH"

# Get node ID
NODE_ID=$($BINARY tendermint show-node-id --home $HOME_DIR)
log "Validator $VALIDATOR_NUM Node ID: $NODE_ID"

# Request tokens from the primary validator
log "Requesting tokens from primary validator..."
$BINARY tx bank send \
    cooluser \
    $VALIDATOR_ADDRESS \
    1500000000000uvna \
    --chain-id=$CHAIN_ID \
    --keyring-backend=test \
    --home ~/.verana \
    --fees 800000uvna \
    --gas 800000 \
    --gas-adjustment 1.3 \
    -y

# Wait for tokens to arrive
log "Waiting for tokens to arrive..."
sleep 10

# Create validator.json
log "Creating validator configuration..."
cat > "$HOME_DIR/validator.json" << EOF
{
  "pubkey": $($BINARY tendermint show-validator --home $HOME_DIR),
  "amount": "$STAKE_AMOUNT",
  "moniker": "$MONIKER",
  "identity": "",
  "website": "",
  "security": "",
  "details": "Validator $VALIDATOR_NUM",
  "commission-rate": "0.1",
  "commission-max-rate": "0.2",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1"
}
EOF

# Create validator transaction
log "Creating validator..."
$BINARY tx staking create-validator "$HOME_DIR/validator.json" \
    --from=$VALIDATOR_NAME \
    --chain-id=$CHAIN_ID \
    --keyring-backend=test \
    --home $HOME_DIR \
    --fees 800000uvna \
    --gas 800000 \
    --gas-adjustment 1.3 \
    -y

log "Validator $VALIDATOR_NUM setup complete!"
log "Starting validator with visible logs..."

# Start the chain with visible logs (no background process)
exec $BINARY start --home $HOME_DIR