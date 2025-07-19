# DeshChain Node Setup Guide

Complete guide for setting up and running DeshChain nodes (validators and full nodes).

## Table of Contents

1. [System Requirements](#system-requirements)
2. [Hardware Recommendations](#hardware-recommendations)
3. [Software Prerequisites](#software-prerequisites)
4. [Node Installation](#node-installation)
5. [Configuration](#configuration)
6. [Network Setup](#network-setup)
7. [Security Hardening](#security-hardening)
8. [Monitoring and Maintenance](#monitoring-and-maintenance)
9. [Troubleshooting](#troubleshooting)

## System Requirements

### Minimum Requirements (Testnet)
- **CPU**: 4 cores (Intel Core i5 or equivalent)
- **RAM**: 8 GB
- **Storage**: 100 GB SSD
- **Network**: 100 Mbps broadband
- **OS**: Ubuntu 20.04+ / CentOS 8+ / macOS 12+

### Recommended Requirements (Mainnet)
- **CPU**: 8 cores (Intel Core i7/AMD Ryzen 7 or equivalent)
- **RAM**: 16 GB
- **Storage**: 500 GB NVMe SSD
- **Network**: 1 Gbps dedicated
- **OS**: Ubuntu 22.04 LTS

### Validator Requirements (Mainnet)
- **CPU**: 16 cores (Intel Xeon/AMD EPYC)
- **RAM**: 32 GB
- **Storage**: 1 TB NVMe SSD
- **Network**: 1 Gbps dedicated with redundancy
- **OS**: Ubuntu 22.04 LTS
- **Backup**: Secondary server for failover

## Hardware Recommendations

### Cloud Providers
#### Recommended Configurations

**AWS EC2:**
- **Testnet**: t3.large (2 vCPU, 8 GB RAM) + 100 GB gp3 SSD
- **Mainnet**: c5.2xlarge (8 vCPU, 16 GB RAM) + 500 GB gp3 SSD
- **Validator**: c5.4xlarge (16 vCPU, 32 GB RAM) + 1 TB gp3 SSD

**Google Cloud:**
- **Testnet**: e2-standard-4 (4 vCPU, 16 GB RAM) + 100 GB SSD
- **Mainnet**: c2-standard-8 (8 vCPU, 32 GB RAM) + 500 GB SSD
- **Validator**: c2-standard-16 (16 vCPU, 64 GB RAM) + 1 TB SSD

**Digital Ocean:**
- **Testnet**: s-4vcpu-8gb + 100 GB SSD
- **Mainnet**: c-8 + 500 GB SSD
- **Validator**: c-16 + 1 TB SSD

### Bare Metal Providers
- **Hetzner**: AX41 (AMD Ryzen 5, 64 GB RAM, 512 GB NVMe)
- **OVH**: Rise-1 (Intel Xeon, 32 GB RAM, 512 GB NVMe)
- **Vultr**: Bare Metal (Intel Xeon, 32 GB RAM, 240 GB SSD)

## Software Prerequisites

### Operating System Setup

#### Ubuntu 22.04 LTS (Recommended)
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install essential packages
sudo apt install -y build-essential git curl wget jq lz4 unzip

# Install Go 1.21+
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### CentOS 8/RHEL 8
```bash
# Update system
sudo dnf update -y

# Install essential packages
sudo dnf install -y gcc gcc-c++ make git curl wget jq lz4 unzip

# Install Go 1.21+
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### macOS
```bash
# Install Homebrew if not installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install dependencies
brew install git curl wget jq lz4

# Install Go 1.21+
brew install go@1.21
echo 'export PATH="/opt/homebrew/opt/go@1.21/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Docker Installation (Optional but Recommended)

#### Ubuntu/Debian
```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### CentOS/RHEL
```bash
# Install Docker
sudo dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo dnf install -y docker-ce docker-ce-cli containerd.io

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker

# Add user to docker group
sudo usermod -aG docker $USER
```

## Node Installation

### Method 1: Binary Installation (Recommended)

#### Download Latest Release
```bash
# Create directory
mkdir -p $HOME/deshchain
cd $HOME/deshchain

# Download latest binary
LATEST_VERSION=$(curl -s https://api.github.com/repos/deshchain/deshchain/releases/latest | jq -r .tag_name)
wget https://github.com/deshchain/deshchain/releases/download/${LATEST_VERSION}/deshchaind-linux-amd64

# Make executable and install
chmod +x deshchaind-linux-amd64
sudo mv deshchaind-linux-amd64 /usr/local/bin/deshchaind

# Verify installation
deshchaind version
```

### Method 2: Build from Source

#### Clone and Build
```bash
# Clone repository
git clone https://github.com/deshchain/deshchain.git
cd deshchain

# Checkout latest stable version
git checkout $(git describe --tags --abbrev=0)

# Build binary
make build

# Install binary
sudo cp build/deshchaind /usr/local/bin/

# Verify installation
deshchaind version
```

### Method 3: Docker Installation

#### Pull Docker Image
```bash
# Pull latest image
docker pull deshchain/deshchaind:latest

# Create alias for convenience
echo 'alias deshchaind="docker run --rm -it -v ~/.deshchain:/root/.deshchain deshchain/deshchaind:latest"' >> ~/.bashrc
source ~/.bashrc
```

## Configuration

### Initialize Node

#### Create Node Directory
```bash
# Initialize node
deshchaind init "your-node-name" --chain-id deshchain-1

# Create necessary directories
mkdir -p ~/.deshchain/cosmovisor/genesis/bin
mkdir -p ~/.deshchain/cosmovisor/upgrades
```

#### Download Genesis File

**Mainnet:**
```bash
# Download genesis file
wget https://raw.githubusercontent.com/deshchain/networks/main/mainnet/genesis.json -O ~/.deshchain/config/genesis.json

# Verify genesis file
deshchaind validate-genesis ~/.deshchain/config/genesis.json
```

**Testnet:**
```bash
# Download testnet genesis file
wget https://raw.githubusercontent.com/deshchain/networks/main/testnet/genesis.json -O ~/.deshchain/config/genesis.json

# Verify genesis file
deshchaind validate-genesis ~/.deshchain/config/genesis.json
```

### Configure Node Settings

#### Update config.toml
```bash
# Edit config file
nano ~/.deshchain/config/config.toml
```

**Key Settings:**
```toml
# P2P Configuration
[p2p]
# Seeds for initial peer discovery
seeds = "seed1.deshchain.network:26656,seed2.deshchain.network:26656"

# Persistent peers
persistent_peers = "peer1@node1.deshchain.network:26656,peer2@node2.deshchain.network:26656"

# Maximum number of inbound peers
max_num_inbound_peers = 40

# Maximum number of outbound peers
max_num_outbound_peers = 10

# RPC Configuration
[rpc]
# TCP or UNIX socket address for the RPC server to listen on
laddr = "tcp://127.0.0.1:26657"

# Maximum number of simultaneous connections
max_open_connections = 900

# Consensus Configuration
[consensus]
# How long we wait for a proposal block before prevoting nil
timeout_propose = "3s"
timeout_propose_delta = "500ms"
timeout_prevote = "1s"
timeout_prevote_delta = "500ms"
timeout_precommit = "1s"
timeout_precommit_delta = "500ms"
timeout_commit = "1s"

# Mempool Configuration
[mempool]
# Size of the mempool
size = 5000
# Maximum number of transactions in the cache
cache_size = 10000
```

#### Update app.toml
```bash
# Edit app configuration
nano ~/.deshchain/config/app.toml
```

**Key Settings:**
```toml
# State Sync Configuration
[state-sync]
# State sync snapshots allow other nodes to rapidly join the network
snapshot-interval = 1000
snapshot-keep-recent = 2

# API Configuration
[api]
# Enable defines if the API server should be enabled
enable = true
# Address defines the API server to listen on
address = "tcp://0.0.0.0:1317"

# gRPC Configuration
[grpc]
# Enable defines if the gRPC server should be enabled
enable = true
# Address defines the gRPC server address to bind to
address = "0.0.0.0:9090"

# Pruning Configuration
pruning = "default"
pruning-keep-recent = "100"
pruning-keep-every = "0"
pruning-interval = "10"

# Minimum gas prices
minimum-gas-prices = "0.0001namo"
```

### Set Up Cosmovisor (Recommended for Automatic Upgrades)

#### Install Cosmovisor
```bash
# Install cosmovisor
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest

# Create cosmovisor directories
mkdir -p ~/.deshchain/cosmovisor/genesis/bin
mkdir -p ~/.deshchain/cosmovisor/upgrades

# Copy binary to cosmovisor
cp $(which deshchaind) ~/.deshchain/cosmovisor/genesis/bin/
```

#### Configure Cosmovisor Environment
```bash
# Add to ~/.bashrc or ~/.zshrc
echo 'export DAEMON_NAME=deshchaind' >> ~/.bashrc
echo 'export DAEMON_HOME=$HOME/.deshchain' >> ~/.bashrc
echo 'export DAEMON_RESTART_AFTER_UPGRADE=true' >> ~/.bashrc
echo 'export DAEMON_ALLOW_DOWNLOAD_BINARIES=false' >> ~/.bashrc
echo 'export UNSAFE_SKIP_BACKUP=false' >> ~/.bashrc
source ~/.bashrc
```

## Network Setup

### Firewall Configuration

#### UFW (Ubuntu/Debian)
```bash
# Enable firewall
sudo ufw enable

# Allow SSH
sudo ufw allow 22/tcp

# Allow P2P port
sudo ufw allow 26656/tcp

# Allow RPC port (if needed for public access)
sudo ufw allow 26657/tcp

# Allow gRPC port (if needed)
sudo ufw allow 9090/tcp

# Allow API port (if needed)
sudo ufw allow 1317/tcp

# Check status
sudo ufw status
```

#### FirewallD (CentOS/RHEL)
```bash
# Start and enable firewalld
sudo systemctl start firewalld
sudo systemctl enable firewalld

# Allow SSH
sudo firewall-cmd --permanent --add-service=ssh

# Allow P2P port
sudo firewall-cmd --permanent --add-port=26656/tcp

# Allow RPC port (if needed)
sudo firewall-cmd --permanent --add-port=26657/tcp

# Reload firewall
sudo firewall-cmd --reload
```

### Load Balancer Setup (For Validators)

#### Nginx Configuration
```bash
# Install nginx
sudo apt install -y nginx

# Create configuration file
sudo nano /etc/nginx/sites-available/deshchain
```

**Nginx Config:**
```nginx
upstream deshchain_rpc {
    server 127.0.0.1:26657;
}

upstream deshchain_api {
    server 127.0.0.1:1317;
}

server {
    listen 80;
    server_name your-validator.com;

    location /rpc/ {
        proxy_pass http://deshchain_rpc/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api/ {
        proxy_pass http://deshchain_api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/deshchain /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

## Security Hardening

### SSH Security

#### Disable Root Login and Password Authentication
```bash
# Edit SSH config
sudo nano /etc/ssh/sshd_config
```

**Key Settings:**
```
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
Port 2222  # Change default port
MaxAuthTries 3
ClientAliveInterval 300
ClientAliveCountMax 2
```

```bash
# Restart SSH service
sudo systemctl restart ssh
```

### Fail2Ban Setup
```bash
# Install fail2ban
sudo apt install -y fail2ban

# Create custom jail
sudo nano /etc/fail2ban/jail.local
```

**Fail2Ban Config:**
```ini
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 3

[sshd]
enabled = true
port = 2222
logpath = /var/log/auth.log

[nginx-http-auth]
enabled = true
port = http,https
logpath = /var/log/nginx/error.log
```

```bash
# Start and enable fail2ban
sudo systemctl start fail2ban
sudo systemctl enable fail2ban
```

### Key Management

#### Hardware Security Module (HSM) Setup
For validators, consider using an HSM for key security:

```bash
# Install YubiHSM tools (example)
sudo apt install -y yubihsm-shell

# Configure validator key with HSM
deshchaind keys add validator --recover --keyring-backend file
```

#### Key Backup
```bash
# Backup validator key (KEEP SECURE!)
cp ~/.deshchain/config/priv_validator_key.json ~/validator_key_backup.json

# Backup node key
cp ~/.deshchain/config/node_key.json ~/node_key_backup.json

# Store backups in secure, offline location
```

### System Hardening

#### Automatic Updates
```bash
# Install unattended-upgrades
sudo apt install -y unattended-upgrades

# Configure automatic updates
sudo dpkg-reconfigure unattended-upgrades
```

#### Disable Unnecessary Services
```bash
# List all services
systemctl list-unit-files --type=service

# Disable unnecessary services (example)
sudo systemctl disable apache2
sudo systemctl disable bluetooth
sudo systemctl disable cups
```

## Service Setup

### Systemd Service Configuration

#### Create Service File
```bash
sudo nano /etc/systemd/system/deshchaind.service
```

**Service Configuration:**
```ini
[Unit]
Description=DeshChain Node
After=network-online.target

[Service]
User=deshchain
ExecStart=/usr/local/bin/deshchaind start --home /home/deshchain/.deshchain
Restart=on-failure
RestartSec=3
LimitNOFILE=65535
Environment=DAEMON_NAME=deshchaind
Environment=DAEMON_HOME=/home/deshchain/.deshchain
Environment=DAEMON_RESTART_AFTER_UPGRADE=true

[Install]
WantedBy=multi-user.target
```

#### For Cosmovisor
```bash
sudo nano /etc/systemd/system/cosmovisor.service
```

**Cosmovisor Service:**
```ini
[Unit]
Description=Cosmovisor DeshChain
After=network-online.target

[Service]
User=deshchain
ExecStart=/home/deshchain/go/bin/cosmovisor run start --home /home/deshchain/.deshchain
Restart=on-failure
RestartSec=3
LimitNOFILE=65535
Environment=DAEMON_NAME=deshchaind
Environment=DAEMON_HOME=/home/deshchain/.deshchain
Environment=DAEMON_RESTART_AFTER_UPGRADE=true
Environment=DAEMON_ALLOW_DOWNLOAD_BINARIES=false
Environment=UNSAFE_SKIP_BACKUP=false

[Install]
WantedBy=multi-user.target
```

#### Start Services
```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable and start service
sudo systemctl enable deshchaind  # or cosmovisor
sudo systemctl start deshchaind   # or cosmovisor

# Check status
sudo systemctl status deshchaind
```

## Monitoring and Maintenance

### Log Management

#### Check Logs
```bash
# View real-time logs
sudo journalctl -u deshchaind -f

# View logs from specific time
sudo journalctl -u deshchaind --since "2024-01-01 00:00:00"

# Export logs
sudo journalctl -u deshchaind --since yesterday > deshchain_logs.txt
```

#### Log Rotation
```bash
# Configure log rotation
sudo nano /etc/logrotate.d/deshchain
```

**Log Rotation Config:**
```
/var/log/deshchain/*.log {
    daily
    missingok
    rotate 7
    compress
    notifempty
    create 0644 deshchain deshchain
    postrotate
        systemctl reload deshchaind
    endscript
}
```

### Performance Monitoring

#### Install Monitoring Tools
```bash
# Install htop, iotop, and nethogs
sudo apt install -y htop iotop nethogs

# Install node_exporter for Prometheus
wget https://github.com/prometheus/node_exporter/releases/latest/download/node_exporter-*.linux-amd64.tar.gz
tar xvfz node_exporter-*.linux-amd64.tar.gz
sudo mv node_exporter-*/node_exporter /usr/local/bin/
sudo useradd -rs /bin/false node_exporter
```

#### Create Node Exporter Service
```bash
sudo nano /etc/systemd/system/node_exporter.service
```

**Node Exporter Service:**
```ini
[Unit]
Description=Node Exporter
After=network.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter

[Install]
WantedBy=multi-user.target
```

```bash
# Start node exporter
sudo systemctl enable node_exporter
sudo systemctl start node_exporter
```

### Health Checks

#### Node Status Script
```bash
# Create health check script
nano ~/check_node_health.sh
```

**Health Check Script:**
```bash
#!/bin/bash

# Check if node is running
if systemctl is-active --quiet deshchaind; then
    echo "✅ DeshChain node is running"
else
    echo "❌ DeshChain node is not running"
    exit 1
fi

# Check sync status
SYNC_STATUS=$(deshchaind status 2>&1 | jq -r '.SyncInfo.catching_up')
if [ "$SYNC_STATUS" == "false" ]; then
    echo "✅ Node is synced"
else
    echo "⏳ Node is syncing"
fi

# Check latest block
LATEST_BLOCK=$(deshchaind status 2>&1 | jq -r '.SyncInfo.latest_block_height')
echo "📊 Latest block: $LATEST_BLOCK"

# Check peers
PEERS=$(deshchaind status 2>&1 | jq -r '.NodeInfo.other.tx_index')
echo "🔗 Connected peers: $PEERS"

# Check disk usage
DISK_USAGE=$(df -h ~/.deshchain | tail -1 | awk '{print $5}')
echo "💾 Disk usage: $DISK_USAGE"
```

```bash
# Make executable
chmod +x ~/check_node_health.sh

# Add to crontab for regular checks
crontab -e
# Add: */5 * * * * /home/deshchain/check_node_health.sh >> /var/log/node_health.log
```

### Backup and Recovery

#### Database Backup
```bash
# Stop node
sudo systemctl stop deshchaind

# Create backup directory
mkdir -p ~/backups/$(date +%Y%m%d)

# Backup data directory
tar -czf ~/backups/$(date +%Y%m%d)/deshchain_data.tar.gz -C ~/.deshchain data

# Backup configuration
tar -czf ~/backups/$(date +%Y%m%d)/deshchain_config.tar.gz -C ~/.deshchain config

# Start node
sudo systemctl start deshchaind
```

#### Automated Backup Script
```bash
# Create backup script
nano ~/backup_node.sh
```

**Backup Script:**
```bash
#!/bin/bash

BACKUP_DIR="$HOME/backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "Stopping DeshChain node..."
sudo systemctl stop deshchaind

echo "Creating backup..."
tar -czf "$BACKUP_DIR/deshchain_data.tar.gz" -C ~/.deshchain data
tar -czf "$BACKUP_DIR/deshchain_config.tar.gz" -C ~/.deshchain config

echo "Starting DeshChain node..."
sudo systemctl start deshchaind

echo "Backup completed: $BACKUP_DIR"

# Clean old backups (keep last 7 days)
find ~/backups -type d -mtime +7 -exec rm -rf {} +
```

```bash
# Make executable and schedule
chmod +x ~/backup_node.sh

# Add to crontab for daily backup
crontab -e
# Add: 0 2 * * * /home/deshchain/backup_node.sh
```

## Troubleshooting

### Common Issues

#### Node Won't Start
```bash
# Check logs for errors
sudo journalctl -u deshchaind -n 100

# Check configuration
deshchaind validate-genesis ~/.deshchain/config/genesis.json

# Reset node data (CAUTION: This will delete all local data)
deshchaind unsafe-reset-all --home ~/.deshchain
```

#### Sync Issues
```bash
# Check sync status
deshchaind status | jq '.SyncInfo'

# Add more peers
deshchaind config config.toml p2p.persistent_peers "peer1@ip:port,peer2@ip:port"

# Use state sync for faster sync
# Edit config.toml [statesync] section
```

#### High Memory Usage
```bash
# Check memory usage
free -h
ps aux | grep deshchaind

# Restart node to clear memory
sudo systemctl restart deshchaind

# Consider increasing swap
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

#### Connection Issues
```bash
# Check firewall
sudo ufw status

# Check port binding
netstat -tulpn | grep deshchaind

# Test connectivity
telnet seed1.deshchain.network 26656
```

### Emergency Procedures

#### Node Recovery
```bash
# Stop node
sudo systemctl stop deshchaind

# Restore from backup
tar -xzf ~/backups/latest/deshchain_data.tar.gz -C ~/.deshchain/
tar -xzf ~/backups/latest/deshchain_config.tar.gz -C ~/.deshchain/

# Start node
sudo systemctl start deshchaind
```

#### Key Recovery
```bash
# Recover validator key
deshchaind keys add validator --recover --keyring-backend file
# Enter mnemonic when prompted

# Verify key
deshchaind keys show validator --keyring-backend file
```

### Performance Optimization

#### Database Optimization
```bash
# Compact database
deshchaind compact-db --home ~/.deshchain

# Prune old data
deshchaind prune --home ~/.deshchain
```

#### System Optimization
```bash
# Increase file limits
echo "* soft nofile 65535" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65535" | sudo tee -a /etc/security/limits.conf

# Optimize network settings
echo "net.core.rmem_max = 16777216" | sudo tee -a /etc/sysctl.conf
echo "net.core.wmem_max = 16777216" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

## Support and Resources

### Official Resources
- **Documentation**: https://docs.deshchain.network
- **GitHub**: https://github.com/deshchain/deshchain
- **Discord**: https://discord.gg/deshchain
- **Telegram**: https://t.me/deshchain

### Community Resources
- **Node Operators Group**: https://t.me/deshchain_validators
- **Technical Support**: support@deshchain.network
- **Bug Reports**: https://github.com/deshchain/deshchain/issues

### Emergency Contacts
- **Critical Issues**: emergency@deshchain.network
- **Security Issues**: security@deshchain.network
- **24/7 Support**: +1-XXX-XXX-XXXX (for validators)

---

*This guide is regularly updated. Check the official documentation for the latest version.*