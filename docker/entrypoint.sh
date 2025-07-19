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

set -e

# DeshChain Docker Entrypoint Script
# Handles initialization and startup of DeshChain node

# Default values
CHAIN_ID=${CHAIN_ID:-"deshchain-1"}
MONIKER=${MONIKER:-"docker-node"}
PERSISTENT_PEERS=${PERSISTENT_PEERS:-""}
SEEDS=${SEEDS:-""}
MINIMUM_GAS_PRICES=${MINIMUM_GAS_PRICES:-"0.0001namo"}
ENABLE_API=${ENABLE_API:-"true"}
ENABLE_GRPC=${ENABLE_GRPC:-"true"}
PRUNING=${PRUNING:-"default"}
SNAPSHOT_INTERVAL=${SNAPSHOT_INTERVAL:-"1000"}
GENESIS_URL=${GENESIS_URL:-""}

# Function to log messages
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') [DESHCHAIN] $1"
}

# Function to initialize node
init_node() {
    log "Initializing DeshChain node..."
    
    # Check if node is already initialized
    if [ ! -f "$DAEMON_HOME/config/config.toml" ]; then
        log "Node not initialized. Running init..."
        deshchaind init "$MONIKER" --chain-id "$CHAIN_ID" --home "$DAEMON_HOME"
    else
        log "Node already initialized."
    fi
}

# Function to download genesis file
download_genesis() {
    if [ -n "$GENESIS_URL" ] && [ ! -f "$DAEMON_HOME/config/genesis.json" ]; then
        log "Downloading genesis file from $GENESIS_URL"
        curl -L "$GENESIS_URL" -o "$DAEMON_HOME/config/genesis.json"
    elif [ ! -f "$DAEMON_HOME/config/genesis.json" ]; then
        log "Warning: No genesis file found and no GENESIS_URL provided"
        log "You may need to manually copy the genesis file"
    fi
}

# Function to configure node
configure_node() {
    log "Configuring DeshChain node..."
    
    local config_file="$DAEMON_HOME/config/config.toml"
    local app_file="$DAEMON_HOME/config/app.toml"
    
    # Configure config.toml
    if [ -n "$PERSISTENT_PEERS" ]; then
        log "Setting persistent peers: $PERSISTENT_PEERS"
        sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEERS\"/" "$config_file"
    fi
    
    if [ -n "$SEEDS" ]; then
        log "Setting seeds: $SEEDS"
        sed -i "s/seeds = \"\"/seeds = \"$SEEDS\"/" "$config_file"
    fi
    
    # Enable prometheus metrics
    sed -i 's/prometheus = false/prometheus = true/' "$config_file"
    
    # Configure app.toml
    log "Setting minimum gas prices: $MINIMUM_GAS_PRICES"
    sed -i "s/minimum-gas-prices = \"\"/minimum-gas-prices = \"$MINIMUM_GAS_PRICES\"/" "$app_file"
    
    # API configuration
    if [ "$ENABLE_API" = "true" ]; then
        log "Enabling API server"
        sed -i 's/enable = false/enable = true/' "$app_file"
        sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/' "$app_file"
    fi
    
    # gRPC configuration
    if [ "$ENABLE_GRPC" = "true" ]; then
        log "Enabling gRPC server"
        sed -i '/\[grpc\]/,/\[/ s/enable = false/enable = true/' "$app_file"
        sed -i 's/address = "localhost:9090"/address = "0.0.0.0:9090"/' "$app_file"
    fi
    
    # Pruning configuration
    log "Setting pruning strategy: $PRUNING"
    sed -i "s/pruning = \"default\"/pruning = \"$PRUNING\"/" "$app_file"
    
    # Snapshot configuration
    log "Setting snapshot interval: $SNAPSHOT_INTERVAL"
    sed -i "s/snapshot-interval = 0/snapshot-interval = $SNAPSHOT_INTERVAL/" "$app_file"
}

# Function to setup cosmovisor
setup_cosmovisor() {
    if [ "$USE_COSMOVISOR" = "true" ]; then
        log "Setting up Cosmovisor..."
        
        # Create cosmovisor directories
        mkdir -p "$DAEMON_HOME/cosmovisor/genesis/bin"
        mkdir -p "$DAEMON_HOME/cosmovisor/upgrades"
        
        # Copy binary to cosmovisor
        cp /usr/local/bin/deshchaind "$DAEMON_HOME/cosmovisor/genesis/bin/"
        
        # Install cosmovisor if not present
        if ! command -v cosmovisor &> /dev/null; then
            log "Installing Cosmovisor..."
            go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
        fi
    fi
}

# Function to reset node data (unsafe)
reset_node() {
    if [ "$UNSAFE_RESET" = "true" ]; then
        log "WARNING: Performing unsafe reset of node data"
        deshchaind unsafe-reset-all --home "$DAEMON_HOME"
    fi
}

# Function to restore from snapshot
restore_snapshot() {
    if [ -n "$SNAPSHOT_URL" ]; then
        log "Restoring from snapshot: $SNAPSHOT_URL"
        
        # Stop any running processes
        pkill -f deshchaind || true
        
        # Reset node data
        deshchaind unsafe-reset-all --home "$DAEMON_HOME" --keep-addr-book
        
        # Download and extract snapshot
        cd "$DAEMON_HOME"
        curl -L "$SNAPSHOT_URL" | tar -xzf -
        
        log "Snapshot restoration completed"
    fi
}

# Function to validate configuration
validate_config() {
    log "Validating configuration..."
    
    # Check if genesis file exists and is valid
    if [ -f "$DAEMON_HOME/config/genesis.json" ]; then
        if ! deshchaind validate-genesis "$DAEMON_HOME/config/genesis.json" --home "$DAEMON_HOME"; then
            log "ERROR: Invalid genesis file"
            exit 1
        fi
        log "Genesis file is valid"
    fi
    
    # Check if config files exist
    if [ ! -f "$DAEMON_HOME/config/config.toml" ] || [ ! -f "$DAEMON_HOME/config/app.toml" ]; then
        log "ERROR: Configuration files missing"
        exit 1
    fi
    
    log "Configuration validation passed"
}

# Function to wait for network
wait_for_network() {
    if [ -n "$WAIT_FOR_SYNC" ] && [ "$WAIT_FOR_SYNC" = "true" ]; then
        log "Waiting for network connectivity..."
        
        # Extract first peer IP for connectivity test
        if [ -n "$PERSISTENT_PEERS" ]; then
            PEER_IP=$(echo "$PERSISTENT_PEERS" | cut -d'@' -f2 | cut -d':' -f1 | head -n1)
            while ! nc -z "$PEER_IP" 26656; do
                log "Waiting for peer connectivity to $PEER_IP:26656..."
                sleep 5
            done
            log "Network connectivity established"
        fi
    fi
}

# Function to backup data
backup_data() {
    if [ "$BACKUP_ON_START" = "true" ]; then
        local backup_dir="/backup/$(date +%Y%m%d_%H%M%S)"
        log "Creating backup at $backup_dir"
        
        mkdir -p "$backup_dir"
        cp -r "$DAEMON_HOME/config" "$backup_dir/"
        
        if [ -d "$DAEMON_HOME/data" ]; then
            tar -czf "$backup_dir/data.tar.gz" -C "$DAEMON_HOME" data
        fi
        
        log "Backup completed"
    fi
}

# Function to set up monitoring
setup_monitoring() {
    if [ "$ENABLE_MONITORING" = "true" ]; then
        log "Setting up monitoring..."
        
        # Create monitoring script
        cat > /usr/local/bin/health_check.sh << 'EOF'
#!/bin/bash
# Health check script for DeshChain node

# Check if node is running
if ! pgrep -f deshchaind > /dev/null; then
    echo "Node is not running"
    exit 1
fi

# Check if node is responding
if ! curl -f http://localhost:26657/health > /dev/null 2>&1; then
    echo "Node is not responding"
    exit 1
fi

# Check sync status
SYNC_STATUS=$(curl -s http://localhost:26657/status | jq -r '.result.sync_info.catching_up')
if [ "$SYNC_STATUS" = "true" ]; then
    echo "Node is syncing"
else
    echo "Node is synced"
fi

exit 0
EOF
        chmod +x /usr/local/bin/health_check.sh
        
        log "Monitoring setup completed"
    fi
}

# Function to run pre-start hooks
run_pre_start_hooks() {
    if [ -d "/hooks/pre-start" ]; then
        log "Running pre-start hooks..."
        for hook in /hooks/pre-start/*; do
            if [ -x "$hook" ]; then
                log "Running hook: $hook"
                "$hook"
            fi
        done
    fi
}

# Function to run post-start hooks
run_post_start_hooks() {
    if [ -d "/hooks/post-start" ]; then
        log "Running post-start hooks..."
        for hook in /hooks/post-start/*; do
            if [ -x "$hook" ]; then
                log "Running hook: $hook"
                "$hook" &
            fi
        done
    fi
}

# Function to handle signals
handle_signal() {
    log "Received signal, shutting down gracefully..."
    
    # Kill background processes
    pkill -P $$
    
    # Stop node if running
    if pgrep -f deshchaind > /dev/null; then
        pkill -TERM -f deshchaind
        sleep 10
        pkill -KILL -f deshchaind 2>/dev/null || true
    fi
    
    log "Shutdown completed"
    exit 0
}

# Main execution function
main() {
    log "Starting DeshChain node container..."
    log "Chain ID: $CHAIN_ID"
    log "Moniker: $MONIKER"
    log "Home directory: $DAEMON_HOME"
    
    # Set up signal handlers
    trap handle_signal SIGTERM SIGINT
    
    # Ensure we're running as the correct user
    if [ "$(id -u)" = "0" ]; then
        log "Running as root, switching to deshchain user..."
        chown -R deshchain:deshchain /home/deshchain
        exec su-exec deshchain "$0" "$@"
    fi
    
    # Run pre-start hooks
    run_pre_start_hooks
    
    # Setup functions
    init_node
    download_genesis
    configure_node
    setup_cosmovisor
    setup_monitoring
    
    # Optional functions
    reset_node
    restore_snapshot
    backup_data
    
    # Validation
    validate_config
    wait_for_network
    
    # Start the node
    log "Starting DeshChain node..."
    
    if [ "$USE_COSMOVISOR" = "true" ]; then
        log "Using Cosmovisor for node management"
        exec cosmovisor run "$@"
    else
        log "Starting node directly"
        
        # Run post-start hooks in background
        run_post_start_hooks
        
        # Start the node
        exec deshchaind "$@" --home "$DAEMON_HOME"
    fi
}

# Script entry point
if [ "${1#-}" != "$1" ] || [ -z "$1" ]; then
    # If first argument starts with dash or is empty, assume we want to run deshchaind
    set -- deshchaind start "$@"
fi

# If running deshchaind start, use our main function
if [ "$1" = "deshchaind" ] && [ "$2" = "start" ]; then
    shift 2
    main start "$@"
else
    # Otherwise, execute the command directly
    exec "$@"
fi