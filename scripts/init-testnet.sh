#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== DeshChain Testnet Initialization ===${NC}"

# Configuration
CHAIN_ID="deshchain-testnet-1"
MONIKER="DeshChain Testnet Node"
DENOM="namo"
HOME_DIR="$HOME/.deshchain-testnet"
BINARY="deshchaind"

# Cleanup
echo -e "${YELLOW}Cleaning up existing testnet data...${NC}"
rm -rf $HOME_DIR

# Initialize chain
echo -e "${GREEN}Initializing testnet chain...${NC}"
$BINARY init "$MONIKER" --chain-id $CHAIN_ID --home $HOME_DIR

# Create test accounts
echo -e "${GREEN}Creating test accounts...${NC}"
$BINARY keys add alice --keyring-backend test --home $HOME_DIR
$BINARY keys add bob --keyring-backend test --home $HOME_DIR
$BINARY keys add validator --keyring-backend test --home $HOME_DIR
$BINARY keys add founder --keyring-backend test --home $HOME_DIR

# Add genesis accounts with test tokens
echo -e "${GREEN}Adding genesis accounts...${NC}"
$BINARY genesis add-genesis-account alice 50000000000000$DENOM --keyring-backend test --home $HOME_DIR
$BINARY genesis add-genesis-account bob 50000000000000$DENOM --keyring-backend test --home $HOME_DIR
$BINARY genesis add-genesis-account validator 10000000000000$DENOM --keyring-backend test --home $HOME_DIR
$BINARY genesis add-genesis-account founder 10000000000000$DENOM --keyring-backend test --home $HOME_DIR

# Add test NGOs
echo -e "${GREEN}Adding test NGO accounts...${NC}"
$BINARY genesis add-genesis-account desh1ngo000000000000000000000000000000000001 1000000000000$DENOM --home $HOME_DIR
$BINARY genesis add-genesis-account desh1ngo000000000000000000000000000000000002 1000000000000$DENOM --home $HOME_DIR
$BINARY genesis add-genesis-account desh1ngo000000000000000000000000000000000003 1000000000000$DENOM --home $HOME_DIR

# Create genesis transaction
echo -e "${GREEN}Creating genesis transaction...${NC}"
$BINARY genesis gentx validator 5000000000000$DENOM \
  --chain-id $CHAIN_ID \
  --moniker="$MONIKER" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1000000" \
  --keyring-backend test \
  --home $HOME_DIR

# Collect genesis transactions
echo -e "${GREEN}Collecting genesis transactions...${NC}"
$BINARY genesis collect-gentxs --home $HOME_DIR

# Update testnet configuration
echo -e "${GREEN}Updating testnet configuration...${NC}"

# Update app.toml for testing
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0'$DENOM'"/' $HOME_DIR/config/app.toml
sed -i 's/enable = false/enable = true/' $HOME_DIR/config/app.toml
sed -i 's/swagger = false/swagger = true/' $HOME_DIR/config/app.toml

# Enable API and gRPC for testing
sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/' $HOME_DIR/config/app.toml
sed -i 's/address = "localhost:9090"/address = "0.0.0.0:9090"/' $HOME_DIR/config/app.toml

# Update config.toml for faster blocks in testnet
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/' $HOME_DIR/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/' $HOME_DIR/config/config.toml

# Enable Prometheus for testing
sed -i 's/prometheus = false/prometheus = true/' $HOME_DIR/config/config.toml

# Validate genesis
echo -e "${GREEN}Validating genesis...${NC}"
$BINARY genesis validate-genesis --home $HOME_DIR

echo -e "${GREEN}=== Testnet initialization complete! ===${NC}"
echo -e ""
echo -e "${YELLOW}Test Accounts:${NC}"
echo -e "Alice: $($BINARY keys show alice -a --keyring-backend test --home $HOME_DIR)"
echo -e "Bob: $($BINARY keys show bob -a --keyring-backend test --home $HOME_DIR)"
echo -e "Validator: $($BINARY keys show validator -a --keyring-backend test --home $HOME_DIR)"
echo -e "Founder: $($BINARY keys show founder -a --keyring-backend test --home $HOME_DIR)"
echo -e ""
echo -e "${YELLOW}Test NGO Addresses:${NC}"
echo -e "NGO 1: desh1ngo000000000000000000000000000000000001"
echo -e "NGO 2: desh1ngo000000000000000000000000000000000002"
echo -e "NGO 3: desh1ngo000000000000000000000000000000000003"
echo -e ""
echo -e "To start the testnet: ${GREEN}make start-testnet${NC}"
echo -e "To reset testnet: ${GREEN}make init-testnet${NC}"