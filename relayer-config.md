# Setting Up an IBC Relayer between Verana and Osmosis Testnet

This guide walks through setting up the Go IBC Relayer to connect a custom Verana blockchain with the Osmosis testnet.

## Prerequisites
- Go 1.19+ installed
- Basic understanding of Cosmos SDK and IBC
- Access to a Verana node
- Access to the internet for connecting to Osmosis testnet

## 1. Initial Setup

Create a working directory and clone the relayer repository:

```bash
# Create a directory for your relayer work
mkdir relay-go-test
cd relay-go-test

# Clone the Go relayer repository
git clone https://github.com/cosmos/relayer.git
cd relayer

# Check out a stable version
git checkout v2.4.2

# Build and install the relayer
make install
```

Verify the installation:

```bash
# Check if the relayer is properly installed
rly -h

# Initialize the relayer configuration
rly config init

# View the default config
rly config show
```

## 2. Adding Chains

Add the Osmosis testnet chain:

```bash
# Add Osmosis testnet using the built-in registry
rly chains add osmosistestnet --testnet

# Verify the chain was added
rly chains list
# Output: 1: osmo-test-5 -> type(cosmos) key(✘) bal(✘) path(✘)
```

Create a configuration file for the Verana chain:

```bash
# Create a JSON config file for Verana
nano verana.json
```

Add the following content to `verana.json`:

```json
{
  "type": "cosmos",
  "value": {
    "key": "default",
    "chain-id": "vna-devnet-1",
    "rpc-addr": "http://node1.devnet.verana.network:26657",
    "account-prefix": "verana",
    "keyring-backend": "test",
    "gas-adjustment": 1.5,
    "gas-prices": "25uvna",
    "debug": true,
    "timeout": "20s",
    "output-format": "json",
    "sign-mode": "direct"
  }
}
```

Add the Verana chain:

```bash
# Add Verana chain using the config file
rly chains add --file verana.json verana

# Verify both chains are added
rly chains list
# Output:
# 1: vna-devnet-1         -> type(cosmos) key(✘) bal(✘) path(✘)
# 2: osmo-test-5          -> type(cosmos) key(✘) bal(✘) path(✘)
```

## 3. Key Management

Generate keys for both chains:

```bash
# Add a key for the Osmosis testnet
rly keys add osmosistestnet osmouser
# Output will include a mnemonic and address - SAVE THESE SECURELY

# Add a key for Verana
rly keys add verana veranauser
# Output will include a mnemonic and address - SAVE THESE SECURELY
```

## 4. Funding the Relayer Accounts

Fund both accounts with tokens:

```bash
# For Verana chain (using an existing funded account)
veranad tx bank send cooluser verana1k2qxlps95nfeamcgkxtlawp6s5rfmd0usyy267 7000000uvna \
  --from cooluser \
  --chain-id vna-devnet-1 \
  --keyring-backend test \
  -y \
  --node http://node1.devnet.verana.network:26657 \
  --fees 600000uvna

# For Osmosis testnet, use a faucet:
# Visit https://testnet.ping.pub/osmosis/faucet
```

## 5. Update Configuration

Edit the relayer configuration to use your keys:

```bash
# Display current configuration
rly config show

# Update keys in the configuration file
# Navigate to ~/.relayer/config/config.yaml and update:
# - For Osmosis: key: osmouser
# - For Verana: key: veranauser
```

Verify account balances:

```bash
# Check Verana balance
rly q balance verana
# Expected output: address {verana1k2qxlps95nfeamcgkxtlawp6s5rfmd0usyy267} balance {7000000uvna}

# Check Osmosis balance
rly q balance osmosistestnet
# Should show the balance received from faucet
```

## 6. Creating a Path

Create a directory for paths and add the path configuration:

```bash
# Create paths directory
mkdir -p ~/.relayer/paths

# Create a path JSON file
nano ~/.relayer/paths/verana-osmosistestnet.json
```

Add this content to the path file:

```json
{
  "src": {
    "chain-id": "osmo-test-5",
    "client-id": "",
    "connection-id": ""
  },
  "dst": {
    "chain-id": "vna-devnet-1",
    "client-id": "",
    "connection-id": ""
  },
  "src-channel-filter": {
    "rule": "",
    "channel-list": []
  }
}
```

Add the path to your configuration:

```bash
# Add paths from directory
rly paths add-dir ~/.relayer/paths

# Verify chains show path availability
rly chains list
# Expected output:
# 1: osmo-test-5          -> type(cosmos) key(✔) bal(✔) path(✔)
# 2: vna-devnet-1         -> type(cosmos) key(✔) bal(✔) path(✔)

# Check path status
rly paths list
# Expected output:
# 0: verana-osmosistestnet -> chns(✔) clnts(✘) conn(✘) (osmo-test-5<>vna-devnet-1)
```

## 7. Creating IBC Clients

Create clients on both chains:

```bash
# Create client on Verana chain (uses Osmosis data)
rly tx client verana osmosistestnet verana-osmosistestnet --override

# Create client on Osmosis chain (uses Verana data)
rly tx client osmosistestnet verana verana-osmosistestnet --override

# Note: The --override flag is crucial when working with custom chains
# that have specific fee requirements
```

## 8. Starting the Relayer

Link the chains and start the relayer:

```bash
# Link the chains and start relaying packets
rly tx link-then-start verana-osmosistestnet

# Verify successful connection
rly paths list
# Expected output when successful:
# 0: verana-osmosistestnet -> chns(✔) clnts(✔) conn(✔) (osmo-test-5<>vna-devnet-1)

# IBC transfer direct
rly tx transfer verana osmosistestnet 1000000uvna osmo1h87hfk67n7ytppj2a3xleh9m454eqytnpp3njk channel-0

# IBC transfer via veranad client
veranad tx ibc-transfer transfer transfer channel-1 osmo1h87hfk67n7ytppj2a3xleh9m454eqytnpp3njk 2000000uvna --from verana-rly --keyring-backend test --chain-id vna-devnet-1 --node http://node1.devnet.verana.network:26657 --fees 600000uvna

```

## Resources

- [Go Relayer Documentation](https://github.com/cosmos/relayer)
- [IBC Protocol Documentation](https://ibcprotocol.org/)
- [Cosmos SDK Documentation](https://docs.cosmos.network/)