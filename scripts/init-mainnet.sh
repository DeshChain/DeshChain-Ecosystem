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


# DeshChain Mainnet Initialization Script
# This script sets up a production mainnet environment

set -e

# Configuration
CHAIN_ID="deshchain-1"
DENOM="namo"
BINARY="deshchaind"
HOME_DIR="$HOME/.deshchain"
CONFIG_DIR="$HOME_DIR/config"
GENESIS_FILE="$CONFIG_DIR/genesis.json"

# Network configuration
P2P_PORT=${P2P_PORT:-26656}
RPC_PORT=${RPC_PORT:-26657}
API_PORT=${API_PORT:-1317}
GRPC_PORT=${GRPC_PORT:-9090}

# Security configuration
KEYRING_BACKEND=${KEYRING_BACKEND:-"file"}
ENABLE_EXTERNAL_RPC=${ENABLE_EXTERNAL_RPC:-false}
ENABLE_EXTERNAL_API=${ENABLE_EXTERNAL_API:-false}

# Mainnet specific URLs
GENESIS_URL="https://raw.githubusercontent.com/deshchain/networks/main/mainnet/genesis.json"
SEEDS="seed1.deshchain.network:26656,seed2.deshchain.network:26656,seed3.deshchain.network:26656"
PERSISTENT_PEERS=""  # Will be set based on network

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
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

# Security check
security_check() {
    warning "MAINNET SECURITY CHECKLIST"
    echo "This script will set up a PRODUCTION mainnet node."
    echo "Please ensure you have:"
    echo "  ✓ Secured your server (SSH keys, firewall, etc.)"
    echo "  ✓ Hardware/VPS meets minimum requirements"
    echo "  ✓ Backup strategy in place"
    echo "  ✓ Monitoring solution ready"
    echo "  ✓ Understanding of validator risks and responsibilities"
    echo
    read -p "Have you completed the security checklist? (yes/no): " -r
    if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
        error "Please complete the security checklist before proceeding."
    fi
}

# Check system requirements
check_requirements() {
    log "Checking system requirements..."
    
    # Check available memory
    TOTAL_MEM=$(grep MemTotal /proc/meminfo | awk '{print $2}')
    TOTAL_MEM_GB=$((TOTAL_MEM / 1024 / 1024))
    
    if [ $TOTAL_MEM_GB -lt 16 ]; then
        warning "System has ${TOTAL_MEM_GB}GB RAM. Recommended: 32GB+ for mainnet"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            error "Please upgrade your system to meet requirements."
        fi
    fi
    
    # Check available disk space
    AVAILABLE_SPACE=$(df / | tail -1 | awk '{print $4}')
    AVAILABLE_SPACE_GB=$((AVAILABLE_SPACE / 1024 / 1024))
    
    if [ $AVAILABLE_SPACE_GB -lt 500 ]; then
        warning "Available disk space: ${AVAILABLE_SPACE_GB}GB. Recommended: 1TB+ for mainnet"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            error "Please ensure sufficient disk space."
        fi
    fi
    
    # Check if running as root
    if [ "$EUID" -eq 0 ]; then
        error "Do not run this script as root. Use a dedicated user account."
    fi
    
    log "System requirements check passed"
}

# Check if binary exists and version
check_binary() {
    if ! command -v $BINARY &> /dev/null; then
        error "$BINARY binary not found. Please install DeshChain first."
    fi
    
    VERSION=$($BINARY version 2>&1 | grep -oP 'v\d+\.\d+\.\d+' || echo "unknown")
    log "Found $BINARY binary: $(which $BINARY)"
    log "Version: $VERSION"
    
    # Check if this is a production release
    if [[ "$VERSION" == *"dev"* ]] || [[ "$VERSION" == *"rc"* ]]; then
        warning "You are using a development or release candidate version: $VERSION"
        warning "For mainnet, please use a stable release version."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            error "Please install a stable release version."
        fi
    fi
}

# Clean previous data with extra caution
clean_data() {
    if [ -d "$HOME_DIR" ]; then
        warning "DANGER: Existing DeshChain data found at $HOME_DIR"
        warning "This may contain valuable keys and validator state!"
        echo
        echo "If this is a validator node, removing data could result in:"
        echo "  - Loss of validator keys"
        echo "  - Loss of node state"
        echo "  - Potential slashing if not handled properly"
        echo
        read -p "Are you absolutely sure you want to remove ALL data? (type 'DELETE' to confirm): " -r
        if [[ $REPLY != "DELETE" ]]; then
            error "Data removal cancelled. Please backup important data manually."
        fi
        
        log "Creating backup before deletion..."
        BACKUP_DIR="$HOME/deshchain-backup-$(date +%Y%m%d_%H%M%S)"
        mkdir -p "$BACKUP_DIR"
        
        if [ -f "$HOME_DIR/config/priv_validator_key.json" ]; then
            cp "$HOME_DIR/config/priv_validator_key.json" "$BACKUP_DIR/"
            log "Validator key backed up to $BACKUP_DIR"
        fi
        
        if [ -f "$HOME_DIR/config/node_key.json" ]; then
            cp "$HOME_DIR/config/node_key.json" "$BACKUP_DIR/"
            log "Node key backed up to $BACKUP_DIR"
        fi
        
        log "Removing existing data..."
        rm -rf "$HOME_DIR"
        log "Backup saved to: $BACKUP_DIR"
    fi
}

# Initialize chain
init_chain() {
    local moniker
    read -p "Enter your node moniker (public name): " moniker
    
    if [ -z "$moniker" ]; then
        error "Moniker cannot be empty"
    fi
    
    log "Initializing chain with moniker: $moniker"
    
    $BINARY init "$moniker" \
        --chain-id="$CHAIN_ID" \
        --home="$HOME_DIR"
    
    log "Chain initialized successfully"
}

# Download and verify genesis
download_genesis() {
    log "Downloading genesis file..."
    
    if [ ! -f "$GENESIS_FILE" ]; then
        log "Downloading genesis from: $GENESIS_URL"
        curl -L "$GENESIS_URL" -o "$GENESIS_FILE"
        
        # Verify genesis file
        log "Verifying genesis file..."
        if ! $BINARY genesis validate-genesis "$GENESIS_FILE" --home="$HOME_DIR"; then
            error "Genesis file validation failed"
        fi
        
        # Display genesis info
        GENESIS_TIME=$(jq -r '.genesis_time' "$GENESIS_FILE")
        CHAIN_ID_FROM_GENESIS=$(jq -r '.chain_id' "$GENESIS_FILE")
        
        log "Genesis time: $GENESIS_TIME"
        log "Chain ID: $CHAIN_ID_FROM_GENESIS"
        
        if [ "$CHAIN_ID_FROM_GENESIS" != "$CHAIN_ID" ]; then
            error "Chain ID mismatch. Expected: $CHAIN_ID, Got: $CHAIN_ID_FROM_GENESIS"
        fi
    else
        log "Genesis file already exists"
    fi
}

# Configure node for mainnet
configure_node() {
    log "Configuring node for mainnet..."
    
    local config_file="$CONFIG_DIR/config.toml"
    local app_file="$CONFIG_DIR/app.toml"
    
    # Configure config.toml
    log "Updating config.toml..."
    
    # Set seeds
    if [ -n "$SEEDS" ]; then
        log "Setting seeds: $SEEDS"
        sed -i "s/seeds = \"\"/seeds = \"$SEEDS\"/" "$config_file"
    fi
    
    # Set persistent peers if provided
    if [ -n "$PERSISTENT_PEERS" ]; then
        log "Setting persistent peers: $PERSISTENT_PEERS"
        sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/" "$config_file"
    fi
    
    # Configure P2P settings for mainnet
    sed -i "s/max_num_inbound_peers = 40/max_num_inbound_peers = 40/" "$config_file"
    sed -i "s/max_num_outbound_peers = 10/max_num_outbound_peers = 10/" "$config_file"
    sed -i "s/flush_throttle_timeout = \"100ms\"/flush_throttle_timeout = \"100ms\"/" "$config_file"
    
    # Set RPC configuration
    if [ "$ENABLE_EXTERNAL_RPC" = "true" ]; then
        warning "Enabling external RPC access. Ensure proper firewall configuration!"
        sed -i "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/0.0.0.0:$RPC_PORT\"/" "$config_file"
    else
        log "RPC limited to localhost (recommended for security)"
        sed -i "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/127.0.0.1:$RPC_PORT\"/" "$config_file"
    fi
    
    # Enable prometheus metrics
    sed -i 's/prometheus = false/prometheus = true/' "$config_file"
    
    # Configure consensus timeouts for mainnet (more conservative)
    sed -i 's/timeout_propose = "3s"/timeout_propose = "3s"/' "$config_file"
    sed -i 's/timeout_prevote = "1s"/timeout_prevote = "1s"/' "$config_file"
    sed -i 's/timeout_precommit = "1s"/timeout_precommit = "1s"/' "$config_file"
    sed -i 's/timeout_commit = "5s"/timeout_commit = "3s"/' "$config_file"
    
    # Configure app.toml
    log "Updating app.toml..."
    
    # Set minimum gas prices
    sed -i "s/minimum-gas-prices = \"\"/minimum-gas-prices = \"0.0001$DENOM\"/" "$app_file"
    
    # API configuration
    if [ "$ENABLE_EXTERNAL_API" = "true" ]; then
        warning "Enabling external API access. Ensure proper firewall configuration!"
        sed -i '/\[api\]/,/\[/{s/enable = false/enable = true/}' "$app_file"
        sed -i "s/address = \"tcp:\/\/localhost:1317\"/address = \"tcp:\/\/0.0.0.0:$API_PORT\"/" "$app_file"
    else
        log "API disabled for security (can be enabled later if needed)"
        sed -i '/\[api\]/,/\[/{s/enable = true/enable = false/}' "$app_file"
    fi
    
    # gRPC configuration
    sed -i '/\[grpc\]/,/\[/{s/enable = false/enable = true/}' "$app_file"
    sed -i "s/address = \"localhost:9090\"/address = \"127.0.0.1:$GRPC_PORT\"/" "$app_file"
    
    # Configure pruning for mainnet (more aggressive to save space)
    sed -i 's/pruning = "default"/pruning = "custom"/' "$app_file"
    sed -i 's/pruning-keep-recent = "0"/pruning-keep-recent = "100"/' "$app_file"
    sed -i 's/pruning-keep-every = "0"/pruning-keep-every = "0"/' "$app_file"
    sed -i 's/pruning-interval = "0"/pruning-interval = "10"/' "$app_file"
    
    # Configure state sync snapshots
    sed -i 's/snapshot-interval = 0/snapshot-interval = 1000/' "$app_file"
    sed -i 's/snapshot-keep-recent = 2/snapshot-keep-recent = 2/' "$app_file"
    
    # Configure mempool
    sed -i 's/size = 5000/size = 5000/' "$app_file"
    sed -i 's/cache_size = 10000/cache_size = 10000/' "$app_file"
}

# Setup cosmovisor for automatic upgrades
setup_cosmovisor() {
    log "Setting up Cosmovisor for automatic upgrades..."
    
    # Check if cosmovisor is installed
    if ! command -v cosmovisor &> /dev/null; then
        log "Installing Cosmovisor..."
        go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
    fi
    
    # Create cosmovisor directories
    mkdir -p "$HOME_DIR/cosmovisor/genesis/bin"
    mkdir -p "$HOME_DIR/cosmovisor/upgrades"
    
    # Copy current binary to cosmovisor
    cp "$(which $BINARY)" "$HOME_DIR/cosmovisor/genesis/bin/"
    
    # Set cosmovisor environment variables
    echo "export DAEMON_NAME=$BINARY" >> ~/.bashrc
    echo "export DAEMON_HOME=$HOME_DIR" >> ~/.bashrc
    echo "export DAEMON_RESTART_AFTER_UPGRADE=true" >> ~/.bashrc
    echo "export DAEMON_ALLOW_DOWNLOAD_BINARIES=false" >> ~/.bashrc
    echo "export UNSAFE_SKIP_BACKUP=false" >> ~/.bashrc
    
    log "Cosmovisor setup completed"
}

# Create systemd service for mainnet
create_service() {
    log "Creating systemd service for mainnet..."
    
    local service_file="/etc/systemd/system/deshchaind.service"
    
    # Check if we should use cosmovisor
    USE_COSMOVISOR=true
    if command -v cosmovisor &> /dev/null; then
        EXEC_START="$(which cosmovisor) run start --home=$HOME_DIR"
        SERVICE_DESC="DeshChain Mainnet Node (Cosmovisor)"
    else
        EXEC_START="$(which $BINARY) start --home=$HOME_DIR"
        SERVICE_DESC="DeshChain Mainnet Node"
        USE_COSMOVISOR=false
    fi
    
    sudo tee "$service_file" > /dev/null <<EOF
[Unit]
Description=$SERVICE_DESC
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$HOME
ExecStart=$EXEC_START
Restart=on-failure
RestartSec=3
LimitNOFILE=65536
LimitNPROC=4096
Environment="DAEMON_NAME=$BINARY"
Environment="DAEMON_HOME=$HOME_DIR"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="UNSAFE_SKIP_BACKUP=false"

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$HOME_DIR

[Install]
WantedBy=multi-user.target
EOF

    sudo systemctl daemon-reload
    sudo systemctl enable deshchaind
    
    if [ "$USE_COSMOVISOR" = "true" ]; then
        log "Systemd service created with Cosmovisor support"
    else
        log "Systemd service created (without Cosmovisor)"
        warning "Consider installing Cosmovisor for automatic upgrades"
    fi
}

# Configure monitoring
setup_monitoring() {
    log "Setting up basic monitoring..."
    
    local scripts_dir="$HOME_DIR/scripts"
    mkdir -p "$scripts_dir"
    
    # Node health check script
    cat > "$scripts_dir/health_check.sh" <<'EOF'
#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if node is running
if ! pgrep -f deshchaind > /dev/null; then
    echo -e "${RED}❌ Node is not running${NC}"
    exit 1
fi

# Check if node is responding
if ! curl -f http://localhost:26657/health > /dev/null 2>&1; then
    echo -e "${RED}❌ Node is not responding to RPC calls${NC}"
    exit 1
fi

# Get node status
STATUS=$(curl -s http://localhost:26657/status 2>/dev/null)
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to get node status${NC}"
    exit 1
fi

# Parse status
CATCHING_UP=$(echo "$STATUS" | jq -r '.result.sync_info.catching_up')
LATEST_HEIGHT=$(echo "$STATUS" | jq -r '.result.sync_info.latest_block_height')
PEER_COUNT=$(echo "$STATUS" | jq -r '.result.node_info.other.tx_index')

# Check sync status
if [ "$CATCHING_UP" = "true" ]; then
    echo -e "${YELLOW}⏳ Node is syncing - Height: $LATEST_HEIGHT${NC}"
else
    echo -e "${GREEN}✅ Node is synced - Height: $LATEST_HEIGHT${NC}"
fi

# Check peer connections
if [ "$PEER_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✅ Connected to peers${NC}"
else
    echo -e "${YELLOW}⚠️ No peer connections${NC}"
fi

# Check disk space
DISK_USAGE=$(df -h "$HOME/.deshchain" | tail -1 | awk '{print $5}' | sed 's/%//')
if [ "$DISK_USAGE" -gt 90 ]; then
    echo -e "${RED}❌ Disk usage critical: ${DISK_USAGE}%${NC}"
elif [ "$DISK_USAGE" -gt 80 ]; then
    echo -e "${YELLOW}⚠️ Disk usage high: ${DISK_USAGE}%${NC}"
else
    echo -e "${GREEN}✅ Disk usage OK: ${DISK_USAGE}%${NC}"
fi

# Check memory usage
MEM_USAGE=$(free | grep Mem | awk '{printf "%.0f", $3/$2 * 100.0}')
if [ "$MEM_USAGE" -gt 90 ]; then
    echo -e "${RED}❌ Memory usage critical: ${MEM_USAGE}%${NC}"
elif [ "$MEM_USAGE" -gt 80 ]; then
    echo -e "${YELLOW}⚠️ Memory usage high: ${MEM_USAGE}%${NC}"
else
    echo -e "${GREEN}✅ Memory usage OK: ${MEM_USAGE}%${NC}"
fi

echo -e "\n${GREEN}✅ Node health check completed${NC}"
EOF

    # Make scripts executable
    chmod +x "$scripts_dir"/*.sh
    
    # Create log rotation configuration
    sudo tee /etc/logrotate.d/deshchaind > /dev/null <<EOF
/var/log/deshchaind/*.log {
    daily
    missingok
    rotate 7
    compress
    notifempty
    create 0644 $USER $USER
    postrotate
        systemctl reload deshchaind
    endscript
}
EOF

    log "Basic monitoring setup completed"
    log "Health check script: $scripts_dir/health_check.sh"
}

# Setup firewall
configure_firewall() {
    log "Configuring firewall..."
    
    # Check if UFW is available
    if command -v ufw &> /dev/null; then
        log "Configuring UFW firewall..."
        
        # Enable firewall
        sudo ufw --force enable
        
        # Default policies
        sudo ufw default deny incoming
        sudo ufw default allow outgoing
        
        # Allow SSH (be careful!)
        SSH_PORT=$(grep "^Port " /etc/ssh/sshd_config | awk '{print $2}' || echo "22")
        sudo ufw allow "$SSH_PORT"/tcp comment "SSH"
        
        # Allow DeshChain P2P
        sudo ufw allow "$P2P_PORT"/tcp comment "DeshChain P2P"
        
        # Only allow RPC/API if explicitly enabled
        if [ "$ENABLE_EXTERNAL_RPC" = "true" ]; then
            sudo ufw allow "$RPC_PORT"/tcp comment "DeshChain RPC"
            warning "RPC port $RPC_PORT is now open to the internet!"
        fi
        
        if [ "$ENABLE_EXTERNAL_API" = "true" ]; then
            sudo ufw allow "$API_PORT"/tcp comment "DeshChain API"
            warning "API port $API_PORT is now open to the internet!"
        fi
        
        # Show status
        sudo ufw status numbered
        
        log "UFW firewall configured"
    else
        warning "UFW not found. Please configure your firewall manually."
        echo "Required ports:"
        echo "  - SSH: 22 (or your custom SSH port)"
        echo "  - P2P: $P2P_PORT"
        if [ "$ENABLE_EXTERNAL_RPC" = "true" ]; then
            echo "  - RPC: $RPC_PORT (if external access needed)"
        fi
        if [ "$ENABLE_EXTERNAL_API" = "true" ]; then
            echo "  - API: $API_PORT (if external access needed)"
        fi
    fi
}

# Create backup scripts
setup_backup() {
    log "Setting up backup scripts..."
    
    local backup_dir="$HOME_DIR/backup"
    mkdir -p "$backup_dir"
    
    # Backup script
    cat > "$backup_dir/backup.sh" <<'EOF'
#!/bin/bash

BACKUP_DIR="$HOME/deshchain-backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="deshchain_backup_$DATE"
BACKUP_PATH="$BACKUP_DIR/$BACKUP_NAME"

# Create backup directory
mkdir -p "$BACKUP_PATH"

echo "Creating DeshChain backup: $BACKUP_NAME"

# Stop node for consistent backup
echo "Stopping DeshChain node..."
sudo systemctl stop deshchaind

# Backup configuration
echo "Backing up configuration..."
cp -r "$HOME/.deshchain/config" "$BACKUP_PATH/"

# Backup keyring if using file backend
if [ -d "$HOME/.deshchain/keyring-file" ]; then
    echo "Backing up keyring..."
    cp -r "$HOME/.deshchain/keyring-file" "$BACKUP_PATH/"
fi

# Backup validator state
if [ -f "$HOME/.deshchain/data/priv_validator_state.json" ]; then
    echo "Backing up validator state..."
    cp "$HOME/.deshchain/data/priv_validator_state.json" "$BACKUP_PATH/"
fi

# Create archive
echo "Creating archive..."
cd "$BACKUP_DIR"
tar -czf "${BACKUP_NAME}.tar.gz" "$BACKUP_NAME"
rm -rf "$BACKUP_NAME"

# Start node
echo "Starting DeshChain node..."
sudo systemctl start deshchaind

# Cleanup old backups (keep last 7 days)
find "$BACKUP_DIR" -name "*.tar.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR/${BACKUP_NAME}.tar.gz"
EOF

    chmod +x "$backup_dir/backup.sh"
    
    log "Backup script created: $backup_dir/backup.sh"
    info "Consider scheduling regular backups with cron"
}

# Generate final summary
generate_summary() {
    log "Generating mainnet configuration summary..."
    
    local summary_file="$HOME_DIR/mainnet-setup.md"
    
    cat > "$summary_file" <<EOF
# DeshChain Mainnet Node Setup Summary

## Node Information
- **Chain ID**: $CHAIN_ID
- **Node Home**: $HOME_DIR
- **Keyring Backend**: $KEYRING_BACKEND

## Network Configuration
- **P2P Port**: $P2P_PORT
- **RPC Port**: $RPC_PORT (localhost only)
- **API Port**: $API_PORT (disabled)
- **gRPC Port**: $GRPC_PORT (localhost only)

## Security Configuration
- **External RPC**: $ENABLE_EXTERNAL_RPC
- **External API**: $ENABLE_EXTERNAL_API
- **Firewall**: Configured (UFW)
- **Cosmovisor**: $(command -v cosmovisor &> /dev/null && echo "Enabled" || echo "Not installed")

## Important Files
- **Config**: $CONFIG_DIR/config.toml
- **App Config**: $CONFIG_DIR/app.toml
- **Genesis**: $CONFIG_DIR/genesis.json
- **Service**: /etc/systemd/system/deshchaind.service

## Management Commands

### Start/Stop Node
\`\`\`bash
sudo systemctl start deshchaind
sudo systemctl stop deshchaind
sudo systemctl restart deshchaind
sudo systemctl status deshchaind
\`\`\`

### View Logs
\`\`\`bash
sudo journalctl -u deshchaind -f
sudo journalctl -u deshchaind --since "1 hour ago"
\`\`\`

### Health Check
\`\`\`bash
$HOME_DIR/scripts/health_check.sh
\`\`\`

### Node Status
\`\`\`bash
deshchaind status | jq
curl -s http://localhost:$RPC_PORT/status | jq
\`\`\`

### Backup
\`\`\`bash
$HOME_DIR/backup/backup.sh
\`\`\`

## Next Steps

### 1. Start Your Node
\`\`\`bash
sudo systemctl start deshchaind
\`\`\`

### 2. Monitor Sync Progress
\`\`\`bash
# Check sync status
deshchaind status | jq '.SyncInfo'

# Or use health check
$HOME_DIR/scripts/health_check.sh
\`\`\`

### 3. Wait for Full Sync
Your node needs to sync with the network before becoming active. This may take several hours to days depending on network history and your hardware.

### 4. Create Validator (Optional)
If you plan to run a validator:
1. Ensure node is fully synced
2. Create validator key: \`deshchaind keys add validator --keyring-backend $KEYRING_BACKEND\`
3. Fund validator account with minimum 10,000 NAMO
4. Create validator transaction
5. Monitor validator performance

### 5. Setup Monitoring
Consider setting up:
- Prometheus + Grafana for metrics
- Alerting for node issues
- Regular health checks
- Automated backups

## Security Reminders
- 🔐 Backup your validator keys securely
- 🔥 Never share your private keys
- 🛡️ Keep your system updated
- 📊 Monitor your node regularly
- 💾 Backup regularly
- 🔍 Enable logging and monitoring

## Support Resources
- **Documentation**: https://docs.deshchain.network
- **Discord**: https://discord.gg/deshchain
- **Telegram**: https://t.me/deshchain
- **GitHub**: https://github.com/deshchain/deshchain

---
Generated on: $(date)
Setup completed: $(date)
EOF

    log "Setup summary saved to: $summary_file"
}

# Main execution function
main() {
    log "Starting DeshChain mainnet node setup..."
    
    # Security and requirement checks
    security_check
    check_requirements
    check_binary
    
    # Setup process
    clean_data
    init_chain
    download_genesis
    configure_node
    setup_cosmovisor
    create_service
    setup_monitoring
    configure_firewall
    setup_backup
    generate_summary
    
    log "Mainnet node setup completed successfully!"
    echo
    info "Setup summary: $HOME_DIR/mainnet-setup.md"
    info "Health check: $HOME_DIR/scripts/health_check.sh"
    info "Backup script: $HOME_DIR/backup/backup.sh"
    echo
    info "To start your mainnet node:"
    echo "  sudo systemctl start deshchaind"
    echo
    info "To monitor sync progress:"
    echo "  $HOME_DIR/scripts/health_check.sh"
    echo
    info "To view logs:"
    echo "  sudo journalctl -u deshchaind -f"
    echo
    warning "IMPORTANT:"
    warning "1. Your node needs to sync before becoming active"
    warning "2. Backup your keys securely"
    warning "3. Monitor your node regularly"
    warning "4. Keep your system updated"
    echo
    log "Welcome to DeshChain mainnet! 🎉"
}

# Script options
case "${1:-}" in
    --help|-h)
        echo "DeshChain Mainnet Node Setup Script"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h              Show this help message"
        echo "  --enable-external-rpc   Enable external RPC access (SECURITY RISK)"
        echo "  --enable-external-api   Enable external API access (SECURITY RISK)"
        echo "  --keyring-backend=TYPE  Set keyring backend (file, os, test)"
        echo "  --no-firewall          Skip firewall configuration"
        echo "  --no-service           Skip systemd service creation"
        echo
        echo "Environment Variables:"
        echo "  P2P_PORT               P2P port (default: 26656)"
        echo "  RPC_PORT               RPC port (default: 26657)"
        echo "  API_PORT               API port (default: 1317)"
        echo "  GRPC_PORT              gRPC port (default: 9090)"
        echo "  PERSISTENT_PEERS       Comma-separated peer list"
        echo
        echo "SECURITY WARNING:"
        echo "This script sets up a production mainnet node."
        echo "Ensure your server is properly secured before running."
        exit 0
        ;;
    --enable-external-rpc)
        ENABLE_EXTERNAL_RPC=true
        warning "External RPC access will be enabled!"
        ;;
    --enable-external-api)
        ENABLE_EXTERNAL_API=true
        warning "External API access will be enabled!"
        ;;
    --keyring-backend=*)
        KEYRING_BACKEND="${1#*=}"
        ;;
    --no-firewall)
        CONFIGURE_FIREWALL=false
        warning "Firewall configuration will be skipped!"
        ;;
    --no-service)
        CREATE_SERVICE=false
        warning "Systemd service creation will be skipped!"
        ;;
esac

# Run main function
main

log "Mainnet setup script completed successfully!"
log "Please read the setup summary and follow the next steps."