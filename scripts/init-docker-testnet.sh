#!/bin/bash

# DeshChain Docker-based 21-Node Testnet Initialization
set -e

echo "🚀 Initializing DeshChain Docker 21-Node Testnet..."

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

# Use deshchaind from Docker
DESHCHAIND="docker run --rm -v $(pwd)/testnet:/testnet deshchain:latest"

# Initialize first node
echo "🏗️ Initializing first validator node..."
$DESHCHAIND init ${NODE_NAMES[0]} --chain-id $CHAIN_ID --home /testnet/node0

# Create validator keys for all nodes
echo "🔑 Creating validator keys..."
for i in $(seq 0 $((VALIDATORS-1))); do
    echo "Creating keys for validator${i}..."
    docker run --rm -v $(pwd)/testnet:/testnet deshchain:latest keys add validator${i} --home /testnet/node${i} --keyring-backend test
done

# Add genesis accounts
echo "💰 Adding genesis accounts..."
for i in $(seq 0 $((VALIDATORS-1))); do
    ADDR=$(docker run --rm -v $(pwd)/testnet:/testnet deshchain:latest keys show validator${i} -a --home /testnet/node${i} --keyring-backend test)
    $DESHCHAIND add-genesis-account $ADDR $INITIAL_COINS --home /testnet/node0
done

# Create genesis transactions
echo "📋 Creating genesis transactions..."
for i in $(seq 0 $((VALIDATORS-1))); do
    echo "Creating gentx for node${i}..."
    if [ $i -eq 0 ]; then
        $DESHCHAIND gentx validator${i} 50000000000$STAKE_DENOM \
            --chain-id $CHAIN_ID \
            --moniker=${NODE_NAMES[$i]} \
            --commission-rate=0.10 \
            --commission-max-rate=0.20 \
            --commission-max-change-rate=0.01 \
            --home /testnet/node${i} \
            --keyring-backend test
    else
        # Initialize other nodes
        $DESHCHAIND init ${NODE_NAMES[$i]} --chain-id $CHAIN_ID --home /testnet/node${i}
        # Copy genesis from node0
        cp testnet/node0/config/genesis.json testnet/node${i}/config/genesis.json
        # Create gentx
        $DESHCHAIND gentx validator${i} 50000000000$STAKE_DENOM \
            --chain-id $CHAIN_ID \
            --moniker=${NODE_NAMES[$i]} \
            --commission-rate=0.10 \
            --commission-max-rate=0.20 \
            --commission-max-change-rate=0.01 \
            --home /testnet/node${i} \
            --keyring-backend test
    fi
done

# Collect all gentxs
echo "📋 Collecting genesis transactions..."
mkdir -p testnet/node0/config/gentx/
for i in $(seq 1 $((VALIDATORS-1))); do
    cp testnet/node${i}/config/gentx/* testnet/node0/config/gentx/
done

# Collect gentxs
$DESHCHAIND collect-gentxs --home /testnet/node0

# Copy final genesis to all nodes
echo "📤 Distributing final genesis file..."
for i in $(seq 1 $((VALIDATORS-1))); do
    cp testnet/node0/config/genesis.json testnet/node${i}/config/genesis.json
done

# Configure persistent peers
echo "🔗 Configuring peer connections..."
PEERS=""
for i in $(seq 0 $((VALIDATORS-1))); do
    NODE_ID=$(docker run --rm -v $(pwd)/testnet:/testnet deshchain:latest tendermint show-node-id --home /testnet/node${i})
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
done

# Create nginx configuration
echo "🌐 Creating nginx load balancer configuration..."
cat > testnet/nginx.conf << 'EOF'
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
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }

        location /api {
            proxy_pass http://api_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
    }
}
EOF

# Create faucet account
echo "💰 Creating faucet account..."
FAUCET_MNEMONIC="witness effort dose make crucial vote nature glove observe dilemma alpha invite lady wage fall shaft stock melody birth check refuse emotion fiscal cruise"
echo $FAUCET_MNEMONIC | docker run --rm -i -v $(pwd)/testnet:/testnet deshchain:latest keys add faucet --recover --home /testnet/node0 --keyring-backend test
FAUCET_ADDRESS=$(docker run --rm -v $(pwd)/testnet:/testnet deshchain:latest keys show faucet -a --home /testnet/node0 --keyring-backend test)
$DESHCHAIND add-genesis-account $FAUCET_ADDRESS 1000000000000000${DENOM} --home /testnet/node0

# Copy updated genesis to all nodes
for i in $(seq 1 $((VALIDATORS-1))); do
    cp testnet/node0/config/genesis.json testnet/node${i}/config/genesis.json
done

echo "✅ Docker-based 21-Node testnet initialization complete!"
echo ""
echo "📊 Summary:"
echo "- Chain ID: $CHAIN_ID"
echo "- Number of validators: $VALIDATORS"
echo "- Token denomination: $DENOM"
echo "- Faucet address: $FAUCET_ADDRESS"
echo ""
echo "🚀 To start the testnet, run:"
echo "docker-compose -f docker-compose.21nodes.yml up -d"