#!/bin/bash

# DeshChain 21-Node Testnet Initialization Script
set -e

echo "🚀 Initializing DeshChain 21-Node Testnet..."

# Configuration
CHAIN_ID="deshchain-testnet-1"
DENOM="unamo"
STAKE_DENOM="stake"
VALIDATORS=21
INITIAL_COINS="100000000000${DENOM},100000000000${STAKE_DENOM}"

# Node names - Indian cultural icons
NODE_NAMES=(
    "Bharatmata"
    "Gandhiji"
    "Netaji"
    "Sardar"
    "Bhagat"
    "Azad"
    "Tilak"
    "Tagore"
    "Vivekananda"
    "APJ"
    "Raman"
    "Ambedkar"
    "Savarkar"
    "Shivaji"
    "Ashoka"
    "Chandragupta"
    "Akbar"
    "RaniLaxmibai"
    "Kabir"
    "Tulsidas"
    "Mirabai"
)

# Clean up previous testnet data
echo "🧹 Cleaning up previous testnet data..."
rm -rf testnet/

# Create testnet directories
echo "📁 Creating testnet directories..."
for i in $(seq 0 $((VALIDATORS-1))); do
    mkdir -p testnet/node${i}
done

# Build the binary if not exists
if [ ! -f "build/deshchaind" ]; then
    echo "🔨 Building DeshChain binary..."
    make build
fi

# Initialize first node
echo "🏗️ Initializing first validator node..."
./build/deshchaind init ${NODE_NAMES[0]} --chain-id $CHAIN_ID --home testnet/node0

# Create initial genesis
echo "📜 Creating genesis file..."
./build/deshchaind genesis add-genesis-account validator0 $INITIAL_COINS --home testnet/node0
./build/deshchaind genesis gentx validator0 50000000000$STAKE_DENOM \
    --chain-id $CHAIN_ID \
    --moniker=${NODE_NAMES[0]} \
    --commission-rate=0.10 \
    --commission-max-rate=0.20 \
    --commission-max-change-rate=0.01 \
    --home testnet/node0

# Copy genesis to other nodes and initialize them
for i in $(seq 1 $((VALIDATORS-1))); do
    echo "🏗️ Initializing validator node${i} - ${NODE_NAMES[$i]}..."
    
    # Initialize node
    ./build/deshchaind init ${NODE_NAMES[$i]} --chain-id $CHAIN_ID --home testnet/node${i}
    
    # Copy genesis from node0
    cp testnet/node0/config/genesis.json testnet/node${i}/config/genesis.json
    
    # Create account and gentx for this validator
    ./build/deshchaind genesis add-genesis-account validator${i} $INITIAL_COINS --home testnet/node${i}
    ./build/deshchaind genesis gentx validator${i} 50000000000$STAKE_DENOM \
        --chain-id $CHAIN_ID \
        --moniker=${NODE_NAMES[$i]} \
        --commission-rate=0.10 \
        --commission-max-rate=0.20 \
        --commission-max-change-rate=0.01 \
        --home testnet/node${i}
done

# Collect all gentxs
echo "📋 Collecting genesis transactions..."
mkdir -p testnet/node0/config/gentx/
for i in $(seq 1 $((VALIDATORS-1))); do
    cp testnet/node${i}/config/gentx/* testnet/node0/config/gentx/
done

# Collect gentxs and update genesis
./build/deshchaind genesis collect-gentxs --home testnet/node0

# Copy final genesis to all nodes
echo "📤 Distributing final genesis file..."
for i in $(seq 1 $((VALIDATORS-1))); do
    cp testnet/node0/config/genesis.json testnet/node${i}/config/genesis.json
done

# Configure persistent peers
echo "🔗 Configuring peer connections..."
PEERS=""
for i in $(seq 0 $((VALIDATORS-1))); do
    NODE_ID=$(./build/deshchaind tendermint show-node-id --home testnet/node${i})
    IP="172.20.0.$((10+i))"
    if [ "$i" -eq 0 ]; then
        PEERS="${NODE_ID}@${IP}:26656"
    else
        PEERS="${PEERS},${NODE_ID}@${IP}:26656"
    fi
done

# Update config for each node
for i in $(seq 0 $((VALIDATORS-1))); do
    echo "⚙️ Configuring node${i}..."
    
    # Update config.toml
    sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PEERS\"/" testnet/node${i}/config/config.toml
    sed -i "s/prometheus = false/prometheus = true/" testnet/node${i}/config/config.toml
    sed -i "s/cors_allowed_origins = \[\]/cors_allowed_origins = \[\"*\"\]/" testnet/node${i}/config/config.toml
    
    # Update app.toml
    sed -i "s/enable = false/enable = true/" testnet/node${i}/config/app.toml
    sed -i "s/swagger = false/swagger = true/" testnet/node${i}/config/app.toml
    
    # Set different ports for each node to avoid conflicts
    P2P_PORT=$((26656 + i*10))
    RPC_PORT=$((26657 + i*10))
    API_PORT=$((1317 + i*10))
    GRPC_PORT=$((9090 + i*10))
    
    sed -i "s/:26656/:$P2P_PORT/" testnet/node${i}/config/config.toml
    sed -i "s/:26657/:$RPC_PORT/" testnet/node${i}/config/config.toml
    sed -i "s/:1317/:$API_PORT/" testnet/node${i}/config/app.toml
    sed -i "s/:9090/:$GRPC_PORT/" testnet/node${i}/config/app.toml
done

# Create nginx configuration
echo "🌐 Creating nginx load balancer configuration..."
cat > testnet/nginx.conf << EOF
events {
    worker_connections 1024;
}

http {
    upstream rpc_backend {
        least_conn;
        server node0:26657;
        server node1:26657;
        server node2:26657;
        server node3:26657;
        server node4:26657;
        server node5:26657;
        server node6:26657;
        server node7:26657;
        server node8:26657;
        server node9:26657;
        server node10:26657;
        server node11:26657;
        server node12:26657;
        server node13:26657;
        server node14:26657;
        server node15:26657;
        server node16:26657;
        server node17:26657;
        server node18:26657;
        server node19:26657;
        server node20:26657;
    }

    upstream api_backend {
        least_conn;
        server node0:1317;
        server node1:1317;
        server node2:1317;
        server node3:1317;
        server node4:1317;
        server node5:1317;
        server node6:1317;
        server node7:1317;
        server node8:1317;
        server node9:1317;
        server node10:1317;
        server node11:1317;
        server node12:1317;
        server node13:1317;
        server node14:1317;
        server node15:1317;
        server node16:1317;
        server node17:1317;
        server node18:1317;
        server node19:1317;
        server node20:1317;
    }

    server {
        listen 80;
        server_name localhost;

        location /rpc {
            proxy_pass http://rpc_backend;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        }

        location /api {
            proxy_pass http://api_backend;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        }
    }
}
EOF

# Create faucet account
echo "💰 Creating faucet account..."
FAUCET_MNEMONIC="witness effort dose make crucial vote nature glove observe dilemma alpha invite lady wage fall shaft stock melody birth check refuse emotion fiscal cruise"
echo $FAUCET_MNEMONIC | ./build/deshchaind keys add faucet --recover --home testnet/node0 --keyring-backend test
FAUCET_ADDRESS=$(./build/deshchaind keys show faucet -a --home testnet/node0 --keyring-backend test)
./build/deshchaind genesis add-genesis-account $FAUCET_ADDRESS 1000000000000000${DENOM} --home testnet/node0

# Copy updated genesis to all nodes
for i in $(seq 1 $((VALIDATORS-1))); do
    cp testnet/node0/config/genesis.json testnet/node${i}/config/genesis.json
done

echo "✅ 21-Node testnet initialization complete!"
echo ""
echo "📊 Summary:"
echo "- Chain ID: $CHAIN_ID"
echo "- Number of validators: $VALIDATORS"
echo "- Token denomination: $DENOM"
echo "- Faucet address: $FAUCET_ADDRESS"
echo ""
echo "🚀 To start the testnet, run:"
echo "docker-compose -f docker-compose.21nodes.yml up -d"
echo ""
echo "📍 Access points:"
echo "- Landing Page: http://localhost:3004"
echo "- Explorer: http://localhost:3000"
echo "- Faucet: http://localhost:4000"
echo "- RPC: http://localhost/rpc"
echo "- API: http://localhost/api"
echo "- Grafana: http://localhost:3005 (admin/deshchain123)"
echo "- Prometheus: http://localhost:9091"