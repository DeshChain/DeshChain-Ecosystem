#!/bin/bash
# Copyright 2024 DeshChain Foundation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# DeshChain Testnet Initialization Script
# This script sets up a new testnet environment

set -e

# Configuration
CHAIN_ID="deshchain-testnet-1"
DENOM="namo"
TOTAL_SUPPLY="1428627660000000"  # 1.42B tokens with 6 decimals
GENESIS_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
BINARY="deshchaind"
HOME_DIR="$HOME/.deshchain"
CONFIG_DIR="$HOME_DIR/config"
GENESIS_FILE="$CONFIG_DIR/genesis.json"

# Validator information
VALIDATOR_MONIKER=${VALIDATOR_MONIKER:-"testnet-validator"}
VALIDATOR_CHAIN_ID=${VALIDATOR_CHAIN_ID:-$CHAIN_ID}

# Network configuration
P2P_PORT=${P2P_PORT:-26656}
RPC_PORT=${RPC_PORT:-26657}
API_PORT=${API_PORT:-1317}
GRPC_PORT=${GRPC_PORT:-9090}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[ERROR] $1${NC}"
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

# Check if binary exists
check_binary() {
    if ! command -v $BINARY &> /dev/null; then
        error "$BINARY binary not found. Please install DeshChain first."
    fi
    
    log "Found $BINARY binary: $(which $BINARY)"
    log "Version: $($BINARY version)"
}

# Clean previous data
clean_data() {
    if [ -d "$HOME_DIR" ]; then
        warning "Existing DeshChain data found at $HOME_DIR"
        read -p "Do you want to remove it? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            log "Removing existing data..."
            rm -rf "$HOME_DIR"
        else
            error "Cannot continue with existing data. Please backup and remove manually."
        fi
    fi
}

# Initialize chain
init_chain() {
    log "Initializing chain with ID: $CHAIN_ID"
    
    $BINARY init "$VALIDATOR_MONIKER" \
        --chain-id="$CHAIN_ID" \
        --home="$HOME_DIR"
    
    log "Chain initialized successfully"
}

# Create validator key
create_keys() {
    log "Creating validator key..."
    
    # Create validator key
    if ! $BINARY keys show validator --home="$HOME_DIR" &>/dev/null; then
        $BINARY keys add validator \
            --home="$HOME_DIR" \
            --keyring-backend=test
    else
        log "Validator key already exists"
    fi
    
    # Get validator address
    VALIDATOR_ADDR=$($BINARY keys show validator -a --home="$HOME_DIR" --keyring-backend=test)
    log "Validator address: $VALIDATOR_ADDR"
    
    # Create additional test accounts
    for account in alice bob charlie; do
        if ! $BINARY keys show $account --home="$HOME_DIR" &>/dev/null; then
            log "Creating test account: $account"
            $BINARY keys add $account \
                --home="$HOME_DIR" \
                --keyring-backend=test
        fi
    done
}

# Configure genesis
configure_genesis() {
    log "Configuring genesis file..."
    
    # Get validator address and pubkey
    VALIDATOR_ADDR=$($BINARY keys show validator -a --home="$HOME_DIR" --keyring-backend=test)
    VALIDATOR_PUBKEY=$($BINARY tendermint show-validator --home="$HOME_DIR")
    
    # Add genesis accounts with allocations
    log "Adding genesis accounts..."
    
    # Validator account (10% of total supply)
    VALIDATOR_AMOUNT="142862766000000$DENOM"
    $BINARY genesis add-genesis-account $VALIDATOR_ADDR $VALIDATOR_AMOUNT --home="$HOME_DIR"
    
    # Test accounts
    TEST_AMOUNT="10000000000$DENOM"  # 10M tokens each
    for account in alice bob charlie; do
        ACCOUNT_ADDR=$($BINARY keys show $account -a --home="$HOME_DIR" --keyring-backend=test)
        $BINARY genesis add-genesis-account $ACCOUNT_ADDR $TEST_AMOUNT --home="$HOME_DIR"
    done
    
    # Community fund (15% of total supply)
    COMMUNITY_AMOUNT="214294149000000$DENOM"
    COMMUNITY_ADDR="desh1community000000000000000000000000000000"
    $BINARY genesis add-genesis-account $COMMUNITY_ADDR $COMMUNITY_AMOUNT --home="$HOME_DIR"
    
    # Development fund (15% of total supply)
    DEVELOPMENT_AMOUNT="214294149000000$DENOM"
    DEVELOPMENT_ADDR="desh1development00000000000000000000000000000"
    $BINARY genesis add-genesis-account $DEVELOPMENT_ADDR $DEVELOPMENT_AMOUNT --home="$HOME_DIR"
    
    # Create genesis transaction
    log "Creating genesis transaction..."
    SELF_DELEGATION="100000000000$DENOM"  # 100M tokens self-delegation
    $BINARY genesis gentx validator $SELF_DELEGATION \
        --chain-id="$CHAIN_ID" \
        --home="$HOME_DIR" \
        --keyring-backend=test \
        --commission-rate="0.10" \
        --commission-max-rate="0.20" \
        --commission-max-change-rate="0.01"
    
    # Collect genesis transactions
    $BINARY genesis collect-gentxs --home="$HOME_DIR"
    
    # Validate genesis
    log "Validating genesis..."
    $BINARY genesis validate-genesis --home="$HOME_DIR"
}

# Configure node settings
configure_node() {
    log "Configuring node settings..."
    
    # Update config.toml
    CONFIG_TOML="$CONFIG_DIR/config.toml"
    
    # Set chain ID
    sed -i "s/chain_id = \"\"/chain_id = \"$CHAIN_ID\"/" "$CONFIG_TOML"
    
    # Set ports
    sed -i "s/laddr = \"tcp:\/\/127.0.0.1:26656\"/laddr = \"tcp:\/\/0.0.0.0:$P2P_PORT\"/" "$CONFIG_TOML"
    sed -i "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/127.0.0.1:$RPC_PORT\"/" "$CONFIG_TOML"
    
    # Enable prometheus metrics
    sed -i 's/prometheus = false/prometheus = true/' "$CONFIG_TOML"
    
    # Set timeouts for faster blocks
    sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/' "$CONFIG_TOML"
    sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "500ms"/' "$CONFIG_TOML"
    sed -i 's/timeout_prevote = "1s"/timeout_prevote = "1s"/' "$CONFIG_TOML"
    sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "500ms"/' "$CONFIG_TOML"
    sed -i 's/timeout_precommit = "1s"/timeout_precommit = "1s"/' "$CONFIG_TOML"
    sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "500ms"/' "$CONFIG_TOML"
    sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/' "$CONFIG_TOML"
    
    # Update app.toml
    APP_TOML="$CONFIG_DIR/app.toml"
    
    # Set minimum gas prices
    sed -i "s/minimum-gas-prices = \"\"/minimum-gas-prices = \"0.0001$DENOM\"/" "$APP_TOML"
    
    # Enable API
    sed -i '/\[api\]/,/\[/{s/enable = false/enable = true/}' "$APP_TOML"
    sed -i "s/address = \"tcp:\/\/localhost:1317\"/address = \"tcp:\/\/0.0.0.0:$API_PORT\"/" "$APP_TOML"
    
    # Enable gRPC
    sed -i '/\[grpc\]/,/\[/{s/enable = false/enable = true/}' "$APP_TOML"
    sed -i "s/address = \"localhost:9090\"/address = \"0.0.0.0:$GRPC_PORT\"/" "$APP_TOML"
    
    # Configure pruning
    sed -i 's/pruning = "default"/pruning = "custom"/' "$APP_TOML"
    sed -i 's/pruning-keep-recent = "0"/pruning-keep-recent = "100"/' "$APP_TOML"
    sed -i 's/pruning-interval = "0"/pruning-interval = "10"/' "$APP_TOML"
    
    # Enable snapshots
    sed -i 's/snapshot-interval = 0/snapshot-interval = 1000/' "$APP_TOML"
    sed -i 's/snapshot-keep-recent = 2/snapshot-keep-recent = 5/' "$APP_TOML"
}

# Create systemd service
create_service() {
    log "Creating systemd service..."
    
    SERVICE_FILE="/etc/systemd/system/deshchaind.service"
    
    sudo tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=DeshChain Testnet Node
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=$USER
ExecStart=$BINARY start --home=$HOME_DIR
Restart=on-failure
RestartSec=3
LimitNOFILE=65536
Environment="DAEMON_NAME=deshchaind"
Environment="DAEMON_HOME=$HOME_DIR"

[Install]
WantedBy=multi-user.target
EOF

    sudo systemctl daemon-reload
    sudo systemctl enable deshchaind
    
    log "Systemd service created and enabled"
}

# Create useful scripts
create_scripts() {
    log "Creating utility scripts..."
    
    SCRIPTS_DIR="$HOME_DIR/scripts"
    mkdir -p "$SCRIPTS_DIR"
    
    # Node status script
    cat > "$SCRIPTS_DIR/status.sh" <<'EOF'
#!/bin/bash
echo "=== DeshChain Testnet Status ==="
echo "Node Status:"
deshchaind status --home=$HOME/.deshchain | jq '.SyncInfo'
echo
echo "Validator Info:"
deshchaind query staking validator $(deshchaind keys show validator --bech val -a --home=$HOME/.deshchain --keyring-backend=test) --home=$HOME/.deshchain
echo
echo "Account Balances:"
for account in validator alice bob charlie; do
    addr=$(deshchaind keys show $account -a --home=$HOME/.deshchain --keyring-backend=test 2>/dev/null || echo "N/A")
    if [ "$addr" != "N/A" ]; then
        balance=$(deshchaind query bank balances $addr --home=$HOME/.deshchain -o json | jq -r '.balances[0].amount // "0"')
        echo "$account: $balance namo"
    fi
done
EOF

    # Transaction script
    cat > "$SCRIPTS_DIR/send.sh" <<'EOF'
#!/bin/bash
if [ $# -ne 3 ]; then
    echo "Usage: $0 <from_account> <to_account> <amount>"
    echo "Example: $0 validator alice 1000000namo"
    exit 1
fi

FROM=$1
TO_ADDR=$(deshchaind keys show $2 -a --home=$HOME/.deshchain --keyring-backend=test)
AMOUNT=$3

deshchaind tx bank send $(deshchaind keys show $FROM -a --home=$HOME/.deshchain --keyring-backend=test) $TO_ADDR $AMOUNT \
    --chain-id=deshchain-testnet-1 \
    --home=$HOME/.deshchain \
    --keyring-backend=test \
    --gas=auto \
    --gas-adjustment=1.5 \
    --fees=1000namo \
    --yes
EOF

    # Delegate script
    cat > "$SCRIPTS_DIR/delegate.sh" <<'EOF'
#!/bin/bash
if [ $# -ne 3 ]; then
    echo "Usage: $0 <from_account> <validator> <amount>"
    echo "Example: $0 alice validator 1000000namo"
    exit 1
fi

FROM=$1
VALIDATOR_ADDR=$(deshchaind keys show $2 --bech val -a --home=$HOME/.deshchain --keyring-backend=test)
AMOUNT=$3

deshchaind tx staking delegate $VALIDATOR_ADDR $AMOUNT \
    --from=$FROM \
    --chain-id=deshchain-testnet-1 \
    --home=$HOME/.deshchain \
    --keyring-backend=test \
    --gas=auto \
    --gas-adjustment=1.5 \
    --fees=1000namo \
    --yes
EOF

    # Make scripts executable
    chmod +x "$SCRIPTS_DIR"/*.sh
    
    log "Utility scripts created in $SCRIPTS_DIR"
}

# Generate configuration summary
generate_summary() {
    log "Generating configuration summary..."
    
    SUMMARY_FILE="$HOME_DIR/testnet-summary.md"
    
    cat > "$SUMMARY_FILE" <<EOF
# DeshChain Testnet Configuration Summary

## Chain Information
- **Chain ID**: $CHAIN_ID
- **Genesis Time**: $GENESIS_TIME
- **Denomination**: $DENOM
- **Total Supply**: $TOTAL_SUPPLY

## Network Ports
- **P2P**: $P2P_PORT
- **RPC**: $RPC_PORT
- **API**: $API_PORT
- **gRPC**: $GRPC_PORT

## Validator Information
- **Moniker**: $VALIDATOR_MONIKER
- **Address**: $VALIDATOR_ADDR
- **Self-Delegation**: 100,000,000 NAMO

## Test Accounts
| Account | Address | Initial Balance |
|---------|---------|-----------------|
| alice   | $(deshchaind keys show alice -a --home="$HOME_DIR" --keyring-backend=test) | 10,000,000 NAMO |
| bob     | $(deshchaind keys show bob -a --home="$HOME_DIR" --keyring-backend=test) | 10,000,000 NAMO |
| charlie | $(deshchaind keys show charlie -a --home="$HOME_DIR" --keyring-backend=test) | 10,000,000 NAMO |

## Quick Commands

### Start Node
\`\`\`bash
sudo systemctl start deshchaind
sudo systemctl status deshchaind
\`\`\`

### Check Status
\`\`\`bash
$HOME_DIR/scripts/status.sh
\`\`\`

### Send Tokens
\`\`\`bash
$HOME_DIR/scripts/send.sh validator alice 1000000namo
\`\`\`

### Delegate Tokens
\`\`\`bash
$HOME_DIR/scripts/delegate.sh alice validator 5000000namo
\`\`\`

### View Logs
\`\`\`bash
sudo journalctl -u deshchaind -f
\`\`\`

## API Endpoints
- **RPC**: http://localhost:$RPC_PORT
- **API**: http://localhost:$API_PORT
- **gRPC**: localhost:$GRPC_PORT

## Configuration Files
- **Node Config**: $CONFIG_DIR/config.toml
- **App Config**: $CONFIG_DIR/app.toml
- **Genesis**: $CONFIG_DIR/genesis.json

## Useful Links
- **Explorer**: https://testnet-explorer.deshchain.network
- **Faucet**: https://testnet-faucet.deshchain.network
- **Documentation**: https://docs.deshchain.network

---
Generated on: $(date)
EOF

    log "Configuration summary saved to $SUMMARY_FILE"
}

# Main execution
main() {
    log "Starting DeshChain testnet initialization..."
    
    # Pre-flight checks
    check_binary
    
    # Setup process
    clean_data
    init_chain
    create_keys
    configure_genesis
    configure_node
    create_service
    create_scripts
    generate_summary
    
    log "Testnet initialization completed successfully!"
    echo
    info "Configuration summary: $HOME_DIR/testnet-summary.md"
    info "Utility scripts: $HOME_DIR/scripts/"
    echo
    info "To start your testnet node:"
    echo "  sudo systemctl start deshchaind"
    echo
    info "To check node status:"
    echo "  $HOME_DIR/scripts/status.sh"
    echo
    info "To view logs:"
    echo "  sudo journalctl -u deshchaind -f"
    echo
    warning "This is a testnet configuration with test keys. DO NOT use in production!"
}

# Script options
case "${1:-}" in
    --help|-h)
        echo "DeshChain Testnet Initialization Script"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h          Show this help message"
        echo "  --clean-only        Only clean existing data"
        echo "  --no-service        Skip systemd service creation"
        echo
        echo "Environment Variables:"
        echo "  VALIDATOR_MONIKER   Validator name (default: testnet-validator)"
        echo "  P2P_PORT           P2P port (default: 26656)"
        echo "  RPC_PORT           RPC port (default: 26657)"
        echo "  API_PORT           API port (default: 1317)"
        echo "  GRPC_PORT          gRPC port (default: 9090)"
        exit 0
        ;;
    --clean-only)
        check_binary
        clean_data
        log "Data cleaned successfully"
        exit 0
        ;;
    --no-service)
        log "Skipping systemd service creation"
        CREATE_SERVICE=false
        ;;
esac

# Run main function
main

log "Testnet initialization script completed!"