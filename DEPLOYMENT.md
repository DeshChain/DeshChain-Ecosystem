# DeshChain Production Deployment Guide

This guide provides comprehensive instructions for deploying DeshChain to production environments.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- 16GB+ RAM (32GB+ recommended for mainnet)
- 1TB+ SSD storage
- Ubuntu 20.04+ or CentOS 8+

### Build & Deploy

```bash
# 1. Build the binary
make build

# 2. Initialize mainnet node
./scripts/init-mainnet.sh

# 3. Start the node
make start-mainnet

# 4. Or use Docker
docker-compose up -d deshchain
```

## 📋 Deployment Options

### Option 1: Native Binary Deployment

#### Mainnet
```bash
# Initialize mainnet environment
./scripts/init-mainnet.sh

# Start the node
systemctl start deshchaind

# Monitor logs
journalctl -u deshchaind -f
```

#### Testnet
```bash
# Initialize testnet environment
./scripts/init-testnet.sh

# Start testnet
make start-testnet
```

### Option 2: Docker Deployment

#### Single Node
```bash
# Start full node
docker-compose up -d deshchain

# With monitoring
docker-compose --profile monitoring up -d

# Validator setup
docker-compose --profile validator up -d deshchain-validator
```

#### Multi-Node Setup
```bash
# Start complete infrastructure
docker-compose --profile monitoring --profile validator --profile backup up -d
```

### Option 3: Kubernetes Deployment

```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n deshchain
```

## 🏗️ Architecture Overview

### Core Components
- **DeshChain Node**: Main blockchain node
- **Tax Module**: 2.5% transaction tax collection
- **Revenue Module**: Platform fee distribution 
- **Donation Module**: NGO payment system
- **Governance Module**: Phased governance with founder protections
- **Sikkebaaz**: Memecoin launchpad
- **GramSuraksha**: Blockchain pension system

### Infrastructure Components
- **Prometheus**: Metrics collection
- **Grafana**: Monitoring dashboards
- **AlertManager**: Alert notifications
- **Nginx**: Load balancer/proxy
- **Backup Service**: Automated backups

## 🔧 Configuration

### Environment Variables

#### Required
```bash
export CHAIN_ID="deshchain-mainnet-1"
export MONIKER="your-node-name"
export KEYRING_BACKEND="file"
```

#### Optional
```bash
export ENABLE_API="true"
export ENABLE_GRPC="true" 
export PROMETHEUS_ENABLED="true"
export BACKUP_ENABLED="true"
```

### Network Configuration

#### Ports
- **26656**: P2P (must be open)
- **26657**: RPC (internal only)
- **1317**: REST API (optional)
- **9090**: gRPC (internal only)
- **26660**: Prometheus metrics

#### Firewall Rules
```bash
# Allow P2P
ufw allow 26656/tcp

# Allow SSH (customize port)
ufw allow 22/tcp

# Deny all other inbound
ufw default deny incoming
ufw default allow outgoing
ufw enable
```

## 📊 Monitoring & Observability

### Health Checks
```bash
# Basic health check
curl http://localhost:26657/health

# Node status
deshchaind status | jq

# Custom health script
./scripts/health_check.sh
```

### Metrics
Access Grafana dashboard at `http://localhost:3000`
- Default credentials: admin/admin123
- DeshChain Overview dashboard pre-configured
- System metrics, blockchain metrics, custom alerts

### Alerting
Configured alerts for:
- Node down/unresponsive
- High memory/CPU usage
- Disk space low
- Sync issues
- Validator missing blocks
- Low peer connections

## 🔐 Security

### Security Checklist
- [ ] Firewall configured
- [ ] SSH key-based authentication
- [ ] Regular security updates
- [ ] Validator keys secured
- [ ] Backup encryption enabled
- [ ] Network monitoring active

### Security Auditing
```bash
# Run security audit
./scripts/security/audit.sh

# Network security scan
./scripts/security/network-scan.sh localhost
```

### Key Management
```bash
# Backup validator keys
cp ~/.deshchain/config/priv_validator_key.json /secure/backup/

# Create new keys
deshchaind keys add mykey --keyring-backend file

# Export keys
deshchaind keys export mykey --keyring-backend file
```

## 💾 Backup & Recovery

### Automated Backup
```bash
# Enable backup service
docker-compose --profile backup up -d backup

# Manual backup
./scripts/backup.sh
```

### Backup Contents
- Validator private keys
- Node configuration
- Blockchain state (optional)
- Application data

### Recovery Process
```bash
# Stop node
systemctl stop deshchaind

# Restore from backup
tar -xzf backup.tar.gz -C ~/.deshchain/

# Restart node
systemctl start deshchaind
```

## 🚀 Scaling & Performance

### Hardware Recommendations

#### Minimum (Testnet)
- 4 CPU cores
- 16GB RAM
- 500GB SSD
- 100Mbps internet

#### Recommended (Mainnet)
- 8+ CPU cores
- 32GB+ RAM
- 2TB+ NVMe SSD
- 1Gbps internet

#### Enterprise (Validator)
- 16+ CPU cores
- 64GB+ RAM
- 4TB+ NVMe SSD
- Redundant internet connections

### Performance Tuning
```bash
# Increase file descriptor limits
ulimit -n 65536

# Optimize disk I/O
echo 'vm.dirty_ratio = 5' >> /etc/sysctl.conf
echo 'vm.dirty_background_ratio = 2' >> /etc/sysctl.conf

# Network optimizations
echo 'net.core.rmem_max = 134217728' >> /etc/sysctl.conf
echo 'net.core.wmem_max = 134217728' >> /etc/sysctl.conf
```

## 🔄 Upgrades & Maintenance

### Cosmos SDK Upgrades
DeshChain uses Cosmovisor for automated upgrades:

```bash
# Prepare upgrade
cosmovisor add-upgrade <upgrade-name> /path/to/new/binary

# Upgrade happens automatically at specified block height
```

### Manual Upgrades
```bash
# Stop node
systemctl stop deshchaind

# Backup current state
./scripts/backup.sh

# Install new version
make install

# Restart node
systemctl start deshchaind
```

### Maintenance Windows
- Regular system updates: Monthly
- Security patches: As needed
- Major upgrades: Quarterly
- Backup verification: Weekly

## 🆘 Troubleshooting

### Common Issues

#### Node Not Syncing
```bash
# Check peers
deshchaind status | jq '.SyncInfo'

# Add peers
deshchaind config set p2p.persistent_peers "peer1@ip:26656,peer2@ip:26656"

# Reset node (last resort)
deshchaind unsafe-reset-all
```

#### High Memory Usage
```bash
# Enable pruning
sed -i 's/pruning = "default"/pruning = "custom"/' ~/.deshchain/config/app.toml
sed -i 's/pruning-keep-recent = "0"/pruning-keep-recent = "100"/' ~/.deshchain/config/app.toml
```

#### Validator Issues
```bash
# Check validator status
deshchaind query staking validator $(deshchaind keys show validator --bech val -a)

# Check missed blocks
deshchaind query slashing signing-info $(deshchaind tendermint show-validator)
```

### Support Channels
- **Documentation**: https://docs.deshchain.network
- **Discord**: https://discord.gg/deshchain
- **Telegram**: https://t.me/deshchain
- **GitHub Issues**: https://github.com/deshchain/deshchain/issues

## 📝 Production Checklist

### Pre-Launch
- [ ] Security audit completed
- [ ] Load testing performed
- [ ] Backup/recovery tested
- [ ] Monitoring configured
- [ ] Incident response plan ready
- [ ] Team trained on operations

### Launch Day
- [ ] Genesis file distributed
- [ ] Network peers coordinated
- [ ] Launch sequence executed
- [ ] Initial validation performed
- [ ] Community notified

### Post-Launch
- [ ] 24/7 monitoring active
- [ ] Performance metrics collected
- [ ] User feedback gathered
- [ ] Issue tracking system updated
- [ ] Documentation maintained

## 🎯 Success Metrics

### Technical KPIs
- **Uptime**: >99.9%
- **Block Time**: ~3 seconds
- **Transaction Throughput**: 1000+ TPS
- **Sync Time**: <24 hours
- **Peer Connections**: 10+ stable peers

### Business KPIs
- **Active Users**: Monthly growth
- **Transaction Volume**: Daily value
- **Validator Participation**: >67% bonded tokens
- **Governance Participation**: >40% voting power
- **NGO Donations**: Monthly distribution tracking

---

**🎉 Congratulations!** You're now ready to deploy DeshChain to production. Remember to follow security best practices and maintain regular backups.

For additional support, reach out to the DeshChain community through our official channels.