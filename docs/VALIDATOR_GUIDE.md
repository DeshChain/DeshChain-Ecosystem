# DeshChain Validator Guide

Complete guide for setting up and operating DeshChain validators with economic incentives and geographic bonuses.

## Table of Contents

1. [Validator Overview](#validator-overview)
2. [Economic Model](#economic-model)
3. [Geographic Incentives](#geographic-incentives)
4. [Requirements](#requirements)
5. [Validator Setup](#validator-setup)
6. [Operations](#operations)
7. [Security](#security)
8. [Monitoring](#monitoring)
9. [Troubleshooting](#troubleshooting)
10. [Best Practices](#best-practices)

## Validator Overview

DeshChain validators secure the network through Tendermint consensus and earn rewards through multiple streams:
- Block rewards and transaction fees
- Geographic location bonuses (India preference)
- Performance bonuses
- MEV (Maximal Extractable Value) distribution
- Long-term validator incentives

### Validator Responsibilities
- **Block Production**: Propose and validate blocks
- **Network Security**: Maintain 99.9%+ uptime
- **Governance Participation**: Vote on proposals
- **Community Engagement**: Support ecosystem growth
- **Cultural Preservation**: Promote Indian values

### Validator Set Limits
- **Maximum Validators**: 150 active validators
- **Minimum Self-Delegation**: 10,000 NAMO
- **Validator Commission**: 0-100% (recommended 5-15%)
- **Unbonding Period**: 21 days

## Economic Model

### Revenue Streams

#### 1. Block Rewards (40% of total emissions)
```
Base APY: 12-18% (varies with staking ratio)
India Bonus: +2% APY
Performance Bonus: +1% APY (top 25% validators)
Total Potential: Up to 21% APY
```

#### 2. Transaction Fees (Variable)
- **2.5% tax on all transactions**
- **50% to validators** (1.25% of transaction volume)
- **25% to development** (0.625%)
- **15% to operations** (0.375%)
- **10% to founder royalty** (0.25%)

#### 3. DEX Trading Fees
- **0.25% trading fee on Money Order DEX**
- **25% to validators** (0.0625% of trading volume)
- **50% to liquidity providers**
- **25% to platform

#### 4. Sikkebaaz Launchpad Fees
- **5% fee on token launches**
- **20% to validators** (1% of launchpad volume)
- **30% to platform**
- **50% to anti-dump mechanisms**

#### 5. MEV Distribution
- **60% to validators** (based on performance)
- **25% to delegators**
- **15% to protocol development**

### Performance-Based Rewards

#### Performance Metrics
1. **Uptime Score** (40% weight)
   - 99.9%+ uptime: 100 points
   - 99.5-99.9%: 80 points
   - 99.0-99.5%: 60 points
   - <99.0%: 0 points

2. **Block Production** (25% weight)
   - Missed blocks penalty: -10 points per missed block
   - Fast block times: +5 points per block under 1s

3. **Governance Participation** (20% weight)
   - Vote on all proposals: 100 points
   - Miss votes: -20 points per missed vote

4. **Network Contribution** (15% weight)
   - Full node operations: +10 points
   - API endpoints: +15 points
   - Archive node: +20 points
   - Relayer operations: +25 points

#### Bonus Calculation
```go
// Performance bonus calculation
func CalculatePerformanceBonus(validator Validator) sdk.Dec {
    baseReward := GetBaseReward(validator)
    performanceScore := CalculatePerformanceScore(validator)
    
    // Top 25% validators get performance bonus
    if performanceScore >= 350 { // out of 400 points
        bonusMultiplier := sdk.NewDecWithPrec(10, 2) // 10% bonus
        return baseReward.Mul(bonusMultiplier)
    }
    
    return sdk.ZeroDec()
}
```

## Geographic Incentives

### India-Based Validator Bonuses

#### Location Verification
Validators must provide proof of Indian location:
1. **IP Address Verification**: Automatic detection
2. **Government ID**: Aadhaar/PAN verification
3. **Address Proof**: Electricity/phone bill
4. **Bank Account**: Indian bank account verification

#### Tier-Based Incentives

**Tier 1 Cities** (+1.5% APY):
- Mumbai, Delhi, Bangalore, Hyderabad, Chennai, Kolkata, Pune, Ahmedabad

**Tier 2 Cities** (+2.0% APY):
- Jaipur, Lucknow, Kanpur, Nagpur, Indore, Patna, Vadodara, Coimbatore

**Tier 3 Cities** (+2.5% APY):
- All other cities and towns

#### Rural Incentives (+3.0% APY)
Special incentives for validators in rural areas:
- Population < 50,000
- Additional infrastructure support
- Technical assistance programs

#### State-Specific Programs
- **Digital India States**: Additional 0.5% bonus
- **Startup Hub States**: Technology support
- **Border States**: Strategic importance bonus

### Verification Process

#### KYC Requirements
```yaml
Required Documents:
  - Aadhaar Card (mandatory)
  - PAN Card (mandatory)
  - Bank Account Details
  - Address Proof (recent)
  - Business Registration (if applicable)
  
Verification Timeline: 3-5 business days
Re-verification: Annual
```

#### Geographic Monitoring
- **Continuous IP monitoring**
- **Random geo-location checks**
- **Community reporting system**
- **Penalty for false claims**: Loss of geographic bonus for 6 months

## Requirements

### Hardware Requirements

#### Minimum (Testnet)
- **CPU**: 4 cores, 2.5GHz+
- **RAM**: 8 GB
- **Storage**: 200 GB SSD
- **Network**: 100 Mbps, 99.9% uptime
- **Bandwidth**: 1 TB/month

#### Recommended (Mainnet)
- **CPU**: 8 cores, 3.0GHz+
- **RAM**: 32 GB
- **Storage**: 1 TB NVMe SSD
- **Network**: 1 Gbps, 99.99% uptime
- **Bandwidth**: Unlimited

#### Professional Validator
- **CPU**: 16+ cores, Intel Xeon/AMD EPYC
- **RAM**: 64 GB+
- **Storage**: 2 TB+ NVMe SSD (RAID 1)
- **Network**: 10 Gbps with redundancy
- **Backup**: Hot standby server
- **Monitoring**: 24/7 monitoring setup

### Software Requirements
- **OS**: Ubuntu 22.04 LTS (recommended)
- **Go**: 1.21+
- **DeshChain**: Latest version
- **Cosmovisor**: For automatic upgrades
- **Monitoring**: Prometheus + Grafana
- **Backup**: Automated backup solution

### Financial Requirements
- **Minimum Self-Delegation**: 10,000 NAMO (~$500-1000)
- **Recommended Stake**: 100,000+ NAMO for better ranking
- **Operating Costs**: $200-500/month
- **Security Deposit**: Additional 5,000 NAMO for slashing protection

## Validator Setup

### Pre-Setup Checklist
- [ ] Hardware/VPS provisioned
- [ ] Domain name configured (recommended)
- [ ] SSL certificates obtained
- [ ] Firewall configured
- [ ] Monitoring setup planned
- [ ] Backup strategy defined

### Step 1: Server Preparation

#### Initial Server Setup
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Create validator user
sudo adduser deshvalidator
sudo usermod -aG sudo deshvalidator
sudo su - deshvalidator

# Setup firewall
sudo ufw enable
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 26656/tcp   # P2P
sudo ufw allow 26657/tcp   # RPC (optional, for monitoring)
```

#### Install Dependencies
```bash
# Install Go
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install essential tools
sudo apt install -y build-essential git curl wget jq lz4 unzip htop
```

### Step 2: Install DeshChain

#### Download Binary
```bash
# Create directory
mkdir -p $HOME/deshchain
cd $HOME/deshchain

# Download latest release
LATEST_VERSION=$(curl -s https://api.github.com/repos/deshchain/deshchain/releases/latest | jq -r .tag_name)
wget https://github.com/deshchain/deshchain/releases/download/${LATEST_VERSION}/deshchaind-linux-amd64

# Install binary
chmod +x deshchaind-linux-amd64
sudo mv deshchaind-linux-amd64 /usr/local/bin/deshchaind

# Verify installation
deshchaind version
```

#### Initialize Node
```bash
# Set variables
MONIKER="your-validator-name"
CHAIN_ID="deshchain-1"

# Initialize node
deshchaind init $MONIKER --chain-id $CHAIN_ID

# Download genesis file
wget https://raw.githubusercontent.com/deshchain/networks/main/mainnet/genesis.json -O ~/.deshchain/config/genesis.json

# Verify genesis
deshchaind validate-genesis ~/.deshchain/config/genesis.json
```

### Step 3: Configure Node

#### Update Configuration Files
```bash
# Set minimum gas prices
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.0001namo"/' ~/.deshchain/config/app.toml

# Set peers
PEERS="peer1@node1.deshchain.network:26656,peer2@node2.deshchain.network:26656"
sed -i "s/persistent_peers = \"\"/persistent_peers = \"$PEERS\"/" ~/.deshchain/config/config.toml

# Set seeds
SEEDS="seed1@seed1.deshchain.network:26656,seed2@seed2.deshchain.network:26656"
sed -i "s/seeds = \"\"/seeds = \"$SEEDS\"/" ~/.deshchain/config/config.toml

# Enable state sync (for faster sync)
sed -i 's/enable = false/enable = true/' ~/.deshchain/config/config.toml
```

#### Setup Cosmovisor
```bash
# Install cosmovisor
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest

# Create directories
mkdir -p ~/.deshchain/cosmovisor/genesis/bin
mkdir -p ~/.deshchain/cosmovisor/upgrades

# Copy binary
cp $(which deshchaind) ~/.deshchain/cosmovisor/genesis/bin/

# Set environment variables
echo 'export DAEMON_NAME=deshchaind' >> ~/.bashrc
echo 'export DAEMON_HOME=$HOME/.deshchain' >> ~/.bashrc
echo 'export DAEMON_RESTART_AFTER_UPGRADE=true' >> ~/.bashrc
source ~/.bashrc
```

### Step 4: Create Validator

#### Generate Keys
```bash
# Create validator key
deshchaind keys add validator --keyring-backend file

# IMPORTANT: Backup your mnemonic phrase securely!
# Create backup of validator key
cp ~/.deshchain/config/priv_validator_key.json ~/validator_key_backup.json
```

#### Fund Validator
```bash
# Get validator address
VALIDATOR_ADDR=$(deshchaind keys show validator -a --keyring-backend file)
echo "Validator address: $VALIDATOR_ADDR"

# Send minimum 10,000 NAMO to this address
# You can buy NAMO on supported exchanges or participate in initial distribution
```

#### Start Node and Sync
```bash
# Create systemd service
sudo tee /etc/systemd/system/deshchaind.service > /dev/null <<EOF
[Unit]
Description=DeshChain Validator
After=network-online.target

[Service]
User=deshvalidator
ExecStart=/home/deshvalidator/go/bin/cosmovisor run start --home /home/deshvalidator/.deshchain
Restart=on-failure
RestartSec=3
LimitNOFILE=65535
Environment=DAEMON_NAME=deshchaind
Environment=DAEMON_HOME=/home/deshvalidator/.deshchain
Environment=DAEMON_RESTART_AFTER_UPGRADE=true

[Install]
WantedBy=multi-user.target
EOF

# Start service
sudo systemctl daemon-reload
sudo systemctl enable deshchaind
sudo systemctl start deshchaind

# Check logs
sudo journalctl -u deshchaind -f
```

#### Wait for Sync
```bash
# Check sync status
deshchaind status | jq '.SyncInfo'

# Node is synced when "catching_up" is false
```

#### Create Validator Transaction
```bash
# Set variables
VALIDATOR_NAME="Your Validator Name"
VALIDATOR_DESCRIPTION="Description of your validator"
VALIDATOR_WEBSITE="https://yourwebsite.com"
VALIDATOR_CONTACT="contact@yourvalidator.com"
COMMISSION_RATE="0.10"  # 10% commission
MIN_SELF_DELEGATION="10000000000"  # 10,000 NAMO (with 6 decimals)

# Create validator
deshchaind tx staking create-validator \
  --amount=$MIN_SELF_DELEGATION"namo" \
  --pubkey=$(deshchaind tendermint show-validator) \
  --moniker="$VALIDATOR_NAME" \
  --chain-id=deshchain-1 \
  --commission-rate=$COMMISSION_RATE \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation=$MIN_SELF_DELEGATION \
  --gas="auto" \
  --gas-adjustment=1.5 \
  --gas-prices="0.0001namo" \
  --from=validator \
  --keyring-backend=file \
  --details="$VALIDATOR_DESCRIPTION" \
  --website="$VALIDATOR_WEBSITE" \
  --security-contact="$VALIDATOR_CONTACT"
```

## Operations

### Daily Operations

#### Check Node Status
```bash
# Node status
deshchaind status | jq

# Validator info
deshchaind query staking validator $(deshchaind keys show validator --bech val -a --keyring-backend file)

# Check if validator is in active set
deshchaind query staking validators --limit 1000 | jq '.validators[] | select(.description.moniker=="'$MONIKER'")'
```

#### Monitor Performance
```bash
# Check uptime
uptime

# Check disk usage
df -h

# Check memory usage
free -h

# Check network connections
netstat -tulpn | grep deshchaind
```

### Delegation Management

#### Check Delegations
```bash
# Check self-delegation
deshchaind query staking delegations $(deshchaind keys show validator -a --keyring-backend file)

# Check all delegations to your validator
VALIDATOR_ADDR=$(deshchaind keys show validator --bech val -a --keyring-backend file)
deshchaind query staking delegations-to $VALIDATOR_ADDR
```

#### Withdraw Rewards
```bash
# Withdraw commission and self-delegation rewards
deshchaind tx distribution withdraw-rewards $VALIDATOR_ADDR \
  --commission \
  --from=validator \
  --chain-id=deshchain-1 \
  --gas="auto" \
  --gas-adjustment=1.5 \
  --gas-prices="0.0001namo" \
  --keyring-backend=file
```

### Governance Participation

#### List Active Proposals
```bash
# List all proposals
deshchaind query gov proposals

# Get proposal details
deshchaind query gov proposal <proposal-id>
```

#### Vote on Proposals
```bash
# Vote on proposal (yes/no/abstain/no_with_veto)
deshchaind tx gov vote <proposal-id> yes \
  --from=validator \
  --chain-id=deshchain-1 \
  --gas="auto" \
  --gas-adjustment=1.5 \
  --gas-prices="0.0001namo" \
  --keyring-backend=file
```

### Validator Maintenance

#### Update Commission
```bash
# Update commission rate (within limits)
deshchaind tx staking edit-validator \
  --commission-rate="0.15" \
  --from=validator \
  --chain-id=deshchain-1 \
  --gas="auto" \
  --gas-adjustment=1.5 \
  --gas-prices="0.0001namo" \
  --keyring-backend=file
```

#### Update Validator Info
```bash
# Update validator metadata
deshchaind tx staking edit-validator \
  --details="Updated description" \
  --website="https://newwebsite.com" \
  --security-contact="newsecurity@validator.com" \
  --from=validator \
  --chain-id=deshchain-1 \
  --gas="auto" \
  --gas-adjustment=1.5 \
  --gas-prices="0.0001namo" \
  --keyring-backend=file
```

## Security

### Key Management

#### Hardware Security Module (HSM)
For production validators, use HSM for key security:

```bash
# Install HSM tools (example for YubiHSM)
sudo apt install -y yubihsm-shell

# Configure HSM integration
# Follow manufacturer's documentation for setup
```

#### Key Backup Strategy
1. **Offline Storage**: Store mnemonic in secure, offline location
2. **Multiple Copies**: Create 3 copies in different physical locations
3. **Encryption**: Encrypt backup files with strong passwords
4. **Regular Testing**: Test restore process quarterly

#### Key Rotation
```bash
# Generate new node key (for P2P communications)
deshchaind unsafe-reset-all --keep-addr-book

# Restart node to use new key
sudo systemctl restart deshchaind
```

### Server Security

#### SSH Hardening
```bash
# Create SSH key pair (on local machine)
ssh-keygen -t ed25519 -C "validator@deshchain"

# Copy public key to server
ssh-copy-id -i ~/.ssh/id_ed25519.pub deshvalidator@your-server-ip

# Edit SSH config on server
sudo nano /etc/ssh/sshd_config
```

**SSH Configuration:**
```
Port 2222
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
MaxAuthTries 3
ClientAliveInterval 300
AllowUsers deshvalidator
```

```bash
# Restart SSH service
sudo systemctl restart ssh
```

#### Firewall Rules
```bash
# Advanced firewall rules
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 2222/tcp      # SSH (custom port)
sudo ufw allow 26656/tcp     # P2P
sudo ufw limit 2222/tcp      # Rate limit SSH

# For monitoring (optional)
sudo ufw allow from trusted-ip to any port 26657  # RPC
sudo ufw allow from trusted-ip to any port 9090   # gRPC
```

#### Intrusion Detection
```bash
# Install fail2ban
sudo apt install -y fail2ban

# Configure jail
sudo nano /etc/fail2ban/jail.local
```

**Fail2Ban Configuration:**
```ini
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 3

[sshd]
enabled = true
port = 2222
logpath = /var/log/auth.log
maxretry = 3
```

### Slashing Protection

#### Tmkms (Tendermint Key Management System)
```bash
# Install tmkms
cargo install tmkms --features=yubihsm

# Initialize tmkms
tmkms init ~/.tmkms

# Configure tmkms
nano ~/.tmkms/tmkms.toml
```

**Tmkms Configuration:**
```toml
[[chain]]
id = "deshchain-1"
key_format = { type = "bech32", account_key_prefix = "deshchainpub", consensus_key_prefix = "deshchainvalconspub" }
state_file = "~/.tmkms/state/deshchain-1_priv_validator_state.json"

[[validator]]
addr = "tcp://127.0.0.1:61278"
chain_id = "deshchain-1"
reconnect = true
secret_key = "~/.tmkms/secrets/deshchain-1_consensus_key"

[[providers.yubihsm]]
adapter = { type = "usb" }
auth = { key = 1, password_file = "~/.tmkms/secrets/yubihsm_password" }
keys = [{ chain_ids = ["deshchain-1"], key = 1 }]
```

#### Double-Sign Protection
```bash
# Monitor double-sign attempts
grep "duplicate vote" ~/.deshchain/logs/deshchain.log

# Implement monitoring alerts
# Use external monitoring service for double-sign detection
```

## Monitoring

### Prometheus Setup

#### Install Prometheus
```bash
# Create prometheus user
sudo useradd --no-create-home --shell /bin/false prometheus

# Download Prometheus
wget https://github.com/prometheus/prometheus/releases/latest/download/prometheus-*.linux-amd64.tar.gz
tar xvfz prometheus-*.linux-amd64.tar.gz

# Install Prometheus
sudo mv prometheus-*/prometheus /usr/local/bin/
sudo mv prometheus-*/promtool /usr/local/bin/

# Create directories
sudo mkdir /etc/prometheus
sudo mkdir /var/lib/prometheus
sudo chown prometheus:prometheus /var/lib/prometheus
```

#### Configure Prometheus
```bash
# Create configuration
sudo nano /etc/prometheus/prometheus.yml
```

**Prometheus Configuration:**
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'deshchain-validator'
    static_configs:
      - targets: ['localhost:26660']  # Tendermint metrics
      - targets: ['localhost:1317']   # Cosmos SDK metrics

  - job_name: 'node-exporter'
    static_configs:
      - targets: ['localhost:9100']

  - job_name: 'process-exporter'
    static_configs:
      - targets: ['localhost:9256']
```

#### Create Prometheus Service
```bash
sudo nano /etc/systemd/system/prometheus.service
```

**Service Configuration:**
```ini
[Unit]
Description=Prometheus
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \
    --config.file /etc/prometheus/prometheus.yml \
    --storage.tsdb.path /var/lib/prometheus/ \
    --web.console.templates=/etc/prometheus/consoles \
    --web.console.libraries=/etc/prometheus/console_libraries \
    --web.listen-address=0.0.0.0:9090

[Install]
WantedBy=multi-user.target
```

### Grafana Dashboard

#### Install Grafana
```bash
# Install Grafana
wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -
echo "deb https://packages.grafana.com/oss/deb stable main" | sudo tee -a /etc/apt/sources.list.d/grafana.list
sudo apt update && sudo apt install -y grafana

# Start Grafana
sudo systemctl enable grafana-server
sudo systemctl start grafana-server
```

#### Configure Grafana
1. Access Grafana: `http://your-server-ip:3000`
2. Login: admin/admin (change password)
3. Add Prometheus data source: `http://localhost:9090`
4. Import DeshChain dashboard (ID: coming soon)

### Alert Manager

#### Install Alert Manager
```bash
# Download AlertManager
wget https://github.com/prometheus/alertmanager/releases/latest/download/alertmanager-*.linux-amd64.tar.gz
tar xvfz alertmanager-*.linux-amd64.tar.gz

# Install AlertManager
sudo mv alertmanager-*/alertmanager /usr/local/bin/
sudo mv alertmanager-*/amtool /usr/local/bin/
```

#### Configure Alerts
```bash
# Create alert rules
sudo nano /etc/prometheus/alert_rules.yml
```

**Alert Rules:**
```yaml
groups:
- name: deshchain_alerts
  rules:
  - alert: ValidatorDown
    expr: up{job="deshchain-validator"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Validator is down"
      description: "DeshChain validator has been down for more than 1 minute"

  - alert: HighMemoryUsage
    expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes > 0.9
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High memory usage"

  - alert: DiskSpaceLow
    expr: node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"} < 0.1
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Disk space low"
```

### Monitoring Scripts

#### Health Check Script
```bash
# Create monitoring script
nano ~/monitor_validator.sh
```

**Monitoring Script:**
```bash
#!/bin/bash

VALIDATOR_ADDR=$(deshchaind keys show validator --bech val -a --keyring-backend file)
TELEGRAM_TOKEN="your_telegram_bot_token"
CHAT_ID="your_chat_id"

# Check if node is running
if ! pgrep -f deshchaind > /dev/null; then
    echo "❌ DeshChain node is not running!"
    curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_TOKEN/sendMessage" \
        -d "chat_id=$CHAT_ID&text=❌ DeshChain validator node is DOWN!"
    exit 1
fi

# Check sync status
SYNC_STATUS=$(deshchaind status 2>&1 | jq -r '.SyncInfo.catching_up')
if [ "$SYNC_STATUS" == "true" ]; then
    echo "⏳ Node is syncing..."
else
    echo "✅ Node is synced"
fi

# Check validator status
VALIDATOR_INFO=$(deshchaind query staking validator $VALIDATOR_ADDR -o json 2>/dev/null)
if [ $? -eq 0 ]; then
    STATUS=$(echo $VALIDATOR_INFO | jq -r '.status')
    JAILED=$(echo $VALIDATOR_INFO | jq -r '.jailed')
    
    if [ "$JAILED" == "true" ]; then
        echo "🚨 VALIDATOR IS JAILED!"
        curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_TOKEN/sendMessage" \
            -d "chat_id=$CHAT_ID&text=🚨 DeshChain validator is JAILED!"
    elif [ "$STATUS" == "BOND_STATUS_BONDED" ]; then
        echo "✅ Validator is active and bonded"
    else
        echo "⚠️ Validator is not in active set"
    fi
fi

# Check missed blocks
MISSED_BLOCKS=$(deshchaind query slashing signing-info $(deshchaind tendermint show-validator) | jq -r '.missed_blocks_counter')
echo "📊 Missed blocks: $MISSED_BLOCKS"

if [ "$MISSED_BLOCKS" -gt 50 ]; then
    echo "⚠️ High number of missed blocks!"
    curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_TOKEN/sendMessage" \
        -d "chat_id=$CHAT_ID&text=⚠️ DeshChain validator missed $MISSED_BLOCKS blocks"
fi
```

```bash
# Make executable and schedule
chmod +x ~/monitor_validator.sh

# Add to crontab
crontab -e
# Add: */5 * * * * /home/deshvalidator/monitor_validator.sh
```

## Troubleshooting

### Common Issues

#### Node Won't Start
```bash
# Check logs
sudo journalctl -u deshchaind -n 100

# Common fixes:
# 1. Check disk space
df -h

# 2. Check permissions
sudo chown -R deshvalidator:deshvalidator ~/.deshchain

# 3. Reset corrupted state
deshchaind unsafe-reset-all --keep-addr-book
```

#### Validator Jailed
```bash
# Check why jailed
deshchaind query slashing signing-info $(deshchaind tendermint show-validator)

# Unjail validator
deshchaind tx slashing unjail \
  --from=validator \
  --chain-id=deshchain-1 \
  --gas="auto" \
  --gas-adjustment=1.5 \
  --gas-prices="0.0001namo" \
  --keyring-backend=file
```

#### Sync Issues
```bash
# Reset and fast sync
deshchaind unsafe-reset-all --keep-addr-book

# Use state sync for faster sync
# Edit config.toml [statesync] section
```

#### High Memory Usage
```bash
# Clear system cache
sudo sync && echo 3 | sudo tee /proc/sys/vm/drop_caches

# Restart node
sudo systemctl restart deshchaind

# Consider upgrading RAM
```

### Performance Optimization

#### Database Tuning
```bash
# Compact database
deshchaind compact-db

# Enable pruning
# Edit app.toml
pruning = "custom"
pruning-keep-recent = "100"
pruning-keep-every = "0"
pruning-interval = "10"
```

#### Network Optimization
```bash
# Optimize TCP settings
echo 'net.core.rmem_max = 16777216' | sudo tee -a /etc/sysctl.conf
echo 'net.core.wmem_max = 16777216' | sudo tee -a /etc/sysctl.conf
echo 'net.ipv4.tcp_rmem = 4096 65536 16777216' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

## Best Practices

### Operational Excellence

#### Infrastructure
1. **Redundancy**: Setup backup validator (sentry nodes)
2. **Monitoring**: 24/7 monitoring with alerts
3. **Automation**: Automate routine tasks
4. **Documentation**: Maintain operational runbooks
5. **Testing**: Regular disaster recovery testing

#### Security
1. **Key Management**: Use HSM for production
2. **Access Control**: Limit SSH access
3. **Updates**: Keep system updated
4. **Backups**: Regular encrypted backups
5. **Auditing**: Regular security audits

#### Performance
1. **Hardware**: Use recommended specs or better
2. **Network**: High-speed, low-latency connection
3. **Optimization**: Regular performance tuning
4. **Capacity Planning**: Monitor growth trends
5. **Scaling**: Plan for delegation growth

### Community Engagement

#### Marketing Your Validator
1. **Website**: Professional validator website
2. **Social Media**: Active presence on Twitter, Telegram
3. **Content**: Educational content about DeshChain
4. **Community**: Participate in governance discussions
5. **Support**: Help new delegators

#### Building Trust
1. **Transparency**: Publish regular updates
2. **Reliability**: Maintain high uptime
3. **Communication**: Responsive to delegator concerns
4. **Fair Commission**: Competitive but sustainable rates
5. **Long-term Commitment**: Demonstrate dedication

### Economic Optimization

#### Commission Strategy
- **Initial Rate**: Start with 5-10% to attract delegators
- **Adjustment Timeline**: Change gradually (max 1% per day)
- **Competitive Analysis**: Monitor other validators
- **Value Proposition**: Justify commission with services

#### Delegation Growth
1. **Self-Delegation**: Maintain meaningful self-stake
2. **Marketing**: Active promotion to potential delegators
3. **Services**: Offer additional services (RPC, APIs)
4. **Partnerships**: Collaborate with other validators
5. **Performance**: Maintain top-tier performance

### Geographic Advantage

#### Maximizing India Bonuses
1. **Location Verification**: Complete all KYC requirements
2. **Tier Optimization**: Consider location for maximum bonus
3. **Rural Opportunities**: Explore rural deployment
4. **Community Building**: Engage local blockchain community
5. **Cultural Promotion**: Actively promote Indian values

#### Compliance
1. **Regular Updates**: Keep location verification current
2. **Documentation**: Maintain all required documents
3. **Reporting**: Respond promptly to verification requests
4. **Community Standards**: Follow all community guidelines
5. **Legal Compliance**: Ensure local regulatory compliance

---

## Support and Resources

### Official Resources
- **Validator Portal**: https://validators.deshchain.network
- **Documentation**: https://docs.deshchain.network/validators
- **GitHub**: https://github.com/deshchain/deshchain
- **Discord**: https://discord.gg/deshchain-validators

### Validator Community
- **Telegram**: https://t.me/deshchain_validators
- **Technical Support**: validators@deshchain.network
- **Emergency Contact**: +91-XXXX-XXXXXX (24/7 for validators)

### Training and Certification
- **Validator Academy**: Monthly training sessions
- **Certification Program**: Official validator certification
- **Best Practices Workshops**: Quarterly workshops
- **Mentorship Program**: Pairing new with experienced validators

---

*This guide is updated regularly. Check the official documentation for the latest version and join the validator community for ongoing support.*