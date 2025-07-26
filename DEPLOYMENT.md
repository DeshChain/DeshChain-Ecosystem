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
- **Identity Module**: W3C DID-compliant identity management with biometric authentication
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
- **IPFS Node**: Decentralized storage for identity documents and NFT metadata
- **Redis**: Identity caching layer for high-performance verification
- **PostgreSQL**: Identity analytics and audit data storage

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

#### Identity Module Configuration
```bash
# Identity system configuration
export IDENTITY_ENABLED="true"
export BIOMETRIC_ENABLED="true"
export KYC_PROVIDER="india_stack"
export PRIVACY_LEVEL="advanced"

# India Stack integration
export AADHAAR_API_KEY="your-aadhaar-api-key"
export DIGILOCKER_CLIENT_ID="your-digilocker-client-id"
export DIGILOCKER_CLIENT_SECRET="your-digilocker-client-secret"
export UPI_API_ENDPOINT="https://api.upi.npci.org.in"

# Storage configuration
export IPFS_ENDPOINT="http://localhost:5001"
export REDIS_URL="redis://localhost:6379"
export POSTGRES_URL="postgresql://username:password@localhost:5432/deshchain_identity"

# Security configuration
export BIOMETRIC_ENCRYPTION_KEY="your-256-bit-encryption-key"
export ZK_PROVING_KEY_PATH="/etc/deshchain/zk-proving-keys"
export HSM_ENABLED="false"  # Set to true for production
export HSM_SLOT_ID="0"
```

### Network Configuration

#### Ports
- **26656**: P2P (must be open)
- **26657**: RPC (internal only)
- **1317**: REST API (optional)
- **9090**: gRPC (internal only)
- **26660**: Prometheus metrics
- **5001**: IPFS API (internal only)
- **4001**: IPFS Swarm (must be open for IPFS network)
- **6379**: Redis (internal only)
- **5432**: PostgreSQL (internal only)
- **8080**: Identity service API (internal only)

#### Firewall Rules
```bash
# Allow P2P
ufw allow 26656/tcp

# Allow IPFS Swarm for decentralized storage
ufw allow 4001/tcp

# Allow SSH (customize port)
ufw allow 22/tcp

# Internal services (only from specific IPs if multi-server setup)
# ufw allow from 10.0.0.0/8 to any port 6379  # Redis
# ufw allow from 10.0.0.0/8 to any port 5432  # PostgreSQL
# ufw allow from 10.0.0.0/8 to any port 5001  # IPFS API

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

# Identity module health check
curl http://localhost:8080/health
curl http://localhost:8080/identity/v1/health

# Infrastructure health checks
curl http://localhost:5001/api/v0/version  # IPFS
redis-cli ping  # Redis
pg_isready -h localhost -p 5432  # PostgreSQL

# Custom health script
./scripts/health_check.sh

# Identity-specific health check
./scripts/identity_health_check.sh
```

### Metrics
Access Grafana dashboard at `http://localhost:3000`
- Default credentials: admin/admin123
- DeshChain Overview dashboard pre-configured
- System metrics, blockchain metrics, custom alerts
- Identity module metrics dashboard
- Biometric authentication success rates
- KYC verification performance
- Privacy compliance metrics

### Alerting
Configured alerts for:
- Node down/unresponsive
- High memory/CPU usage
- Disk space low
- Sync issues
- Validator missing blocks
- Low peer connections
- Identity service unavailable
- Biometric authentication failures exceeding threshold
- KYC verification service disruption
- IPFS node disconnection
- Redis cache performance degradation
- PostgreSQL connection issues
- Privacy compliance violations
- Suspicious identity activity patterns

## 🔐 Security

### Security Checklist
- [ ] Firewall configured
- [ ] SSH key-based authentication
- [ ] Regular security updates
- [ ] Validator keys secured
- [ ] Backup encryption enabled
- [ ] Network monitoring active
- [ ] Identity module security configured
- [ ] Biometric encryption keys secured
- [ ] ZK proving keys properly installed
- [ ] HSM integration configured (production)
- [ ] India Stack API credentials secured
- [ ] IPFS node security hardened
- [ ] Redis authentication configured
- [ ] PostgreSQL access controls implemented
- [ ] Privacy compliance monitoring active

### Security Auditing
```bash
# Run security audit
./scripts/security/audit.sh

# Network security scan
./scripts/security/network-scan.sh localhost

# Identity module security audit
./scripts/security/identity-audit.sh

# Biometric system security check
./scripts/security/biometric-security-check.sh

# Privacy compliance audit
./scripts/security/privacy-audit.sh

# India Stack integration security check
./scripts/security/india-stack-security.sh
```

### Key Management
```bash
# Backup validator keys
cp ~/.deshchain/config/priv_validator_key.json /secure/backup/

# Create new keys
deshchaind keys add mykey --keyring-backend file

# Export keys
deshchaind keys export mykey --keyring-backend file

# Identity module key management
# Backup biometric encryption keys
cp /etc/deshchain/biometric-keys/* /secure/backup/identity/

# Backup ZK proving keys
cp -r /etc/deshchain/zk-proving-keys /secure/backup/identity/

# Backup DID registry keys
deshchaind tx identity backup-keys --from=admin --gas=auto

# Generate new identity service keys
./scripts/identity/generate-service-keys.sh

# Rotate biometric encryption keys (scheduled maintenance)
./scripts/identity/rotate-biometric-keys.sh
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
- Identity module configuration
- Biometric encryption keys
- ZK proving keys and circuits
- DID registry snapshots
- Identity analytics database
- IPFS private keys and configuration
- Redis cache snapshots (optional)
- PostgreSQL database dumps

### Recovery Process
```bash
# Stop all services
systemctl stop deshchaind
systemctl stop ipfs
systemctl stop redis-server
systemctl stop postgresql

# Restore main node data
tar -xzf backup.tar.gz -C ~/.deshchain/

# Restore identity module data
tar -xzf identity-backup.tar.gz -C /etc/deshchain/

# Restore IPFS data
tar -xzf ipfs-backup.tar.gz -C ~/.ipfs/

# Restore PostgreSQL database
pg_restore -h localhost -U deshchain_user -d deshchain_identity identity_backup.dump

# Restore Redis data (if backed up)
redis-cli --rdb redis-backup.rdb

# Start services in order
systemctl start postgresql
systemctl start redis-server
systemctl start ipfs
systemctl start deshchaind

# Verify identity module functionality
./scripts/identity/verify-recovery.sh
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

#### Enterprise (Validator with Identity Services)
- 16+ CPU cores
- 64GB+ RAM
- 4TB+ NVMe SSD
- Redundant internet connections
- Hardware Security Module (HSM) for identity keys
- Additional 2TB SSD for IPFS storage
- Dedicated GPU for biometric processing (optional)
- Backup identity verification hardware

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

# Identity module performance tuning
# Redis optimization
echo 'vm.overcommit_memory = 1' >> /etc/sysctl.conf
echo 'net.core.somaxconn = 65535' >> /etc/sysctl.conf

# PostgreSQL optimization
echo 'kernel.shmmax = 17179869184' >> /etc/sysctl.conf  # 16GB
echo 'kernel.shmall = 4194304' >> /etc/sysctl.conf

# IPFS optimization
echo 'fs.file-max = 2097152' >> /etc/sysctl.conf

# Biometric processing optimization (if using GPU)
# nvidia-smi -pm 1  # Enable persistence mode
# nvidia-smi -acp 0  # Disable auto boost
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

## 🆔 Identity Module Setup

### Prerequisites for Identity Module
```bash
# Install required dependencies
apt-get update
apt-get install -y postgresql-14 redis-server ipfs

# Install biometric libraries (if using biometric authentication)
apt-get install -y libfprint-2 libopencv-dev

# Install ZK proof dependencies
wget https://github.com/iden3/circom/releases/download/v2.1.5/circom-linux-amd64
chmod +x circom-linux-amd64
mv circom-linux-amd64 /usr/local/bin/circom

# Install snarkjs for ZK proof generation
npm install -g snarkjs
```

### Identity Module Configuration
```bash
# Create identity configuration directory
mkdir -p /etc/deshchain/identity
cd /etc/deshchain/identity

# Generate ZK proving keys
./scripts/identity/setup-zk-keys.sh

# Configure PostgreSQL for identity analytics
sudo -u postgres createdb deshchain_identity
sudo -u postgres createuser deshchain_user
sudo -u postgres psql -c "ALTER USER deshchain_user PASSWORD 'secure_password';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE deshchain_identity TO deshchain_user;"

# Run identity database migrations
psql -h localhost -U deshchain_user -d deshchain_identity -f ./migrations/identity_schema.sql

# Configure Redis for identity caching
redis-cli CONFIG SET requirepass "redis_secure_password"
redis-cli CONFIG SET maxmemory 4gb
redis-cli CONFIG SET maxmemory-policy allkeys-lru

# Initialize IPFS for identity documents
ipfs init
ipfs config Addresses.API /ip4/127.0.0.1/tcp/5001
ipfs config Addresses.Gateway /ip4/127.0.0.1/tcp/8080
ipfs daemon &

# Configure India Stack integration
./scripts/identity/setup-india-stack.sh
```

### Identity Module Verification
```bash
# Test identity module functionality
deshchaind tx identity create-did --from=admin --gas=auto

# Test biometric registration
curl -X POST http://localhost:8080/identity/v1/biometric/register \
  -H "Content-Type: application/json" \
  -d '{"did":"did:desh:test123","biometric_type":"fingerprint","template":"base64_template"}'

# Test KYC integration
curl -X POST http://localhost:8080/identity/v1/kyc/verify \
  -H "Content-Type: application/json" \
  -d '{"did":"did:desh:test123","document_type":"aadhaar","document_number":"xxxx-xxxx-1234"}'

# Test credential issuance
deshchaind tx identity issue-credential \
  --issuer="did:desh:issuer123" \
  --subject="did:desh:subject456" \
  --credential-type="KYCCredential" \
  --from=admin \
  --gas=auto

# Verify identity services are healthy
./scripts/identity/health_check.sh
```

### Production Identity Setup Checklist
- [ ] PostgreSQL database secured with authentication
- [ ] Redis configured with password and memory limits
- [ ] IPFS node secured and properly configured
- [ ] ZK proving keys generated and backed up
- [ ] Biometric encryption keys generated and secured
- [ ] India Stack API credentials configured and tested
- [ ] HSM integration configured (for production)
- [ ] Identity analytics dashboard configured
- [ ] Privacy compliance monitoring enabled
- [ ] Audit trail collection verified
- [ ] Cross-module identity sharing tested
- [ ] Identity backup and recovery procedures tested

### Identity Module Scaling
```bash
# Scale identity services for high throughput
# Configure Redis cluster for distributed caching
redis-cli --cluster create 127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002 \
  127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 --cluster-replicas 1

# Configure PostgreSQL read replicas
pg_basebackup -h master-server -D /var/lib/postgresql/14/replica -U replication -W

# Configure IPFS cluster for distributed storage
ipfs-cluster-service init
ipfs-cluster-service daemon &

# Load balance identity API endpoints
# Configure nginx upstream for identity services
```

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

#### Identity Module Issues

##### Identity Service Not Responding
```bash
# Check identity service status
systemctl status deshchain-identity
curl http://localhost:8080/health

# Check logs
journalctl -u deshchain-identity -f

# Restart identity service
systemctl restart deshchain-identity
```

##### Biometric Authentication Failures
```bash
# Check biometric service logs
tail -f /var/log/deshchain/biometric.log

# Verify biometric hardware
lsusb | grep -i "fingerprint\|biometric"

# Test biometric template matching
./scripts/identity/test-biometric-matching.sh

# Clear corrupted biometric templates
deshchaind tx identity clear-biometric --did="did:desh:user123" --from=admin
```

##### KYC Verification Failures
```bash
# Check India Stack connectivity
curl -I https://api.uidai.gov.in/

# Test Aadhaar API connectivity
./scripts/identity/test-aadhaar-connection.sh

# Check DigiLocker integration
curl -X GET "https://api.digilocker.gov.in/public/oauth2/1/authorize" \
  -H "Authorization: Bearer $DIGILOCKER_TOKEN"

# Verify KYC credentials
deshchaind query identity kyc-status did:desh:user123
```

##### IPFS Storage Issues
```bash
# Check IPFS daemon status
ipfs swarm peers | wc -l  # Should show connected peers

# Check IPFS storage
ipfs repo stat

# Fix IPFS connectivity
ipfs swarm connect /ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ

# Garbage collect old data
ipfs repo gc
```

##### Redis Cache Performance Issues
```bash
# Check Redis memory usage
redis-cli info memory

# Check Redis connection
redis-cli ping

# Monitor Redis performance
redis-cli monitor

# Clear Redis cache if needed
redis-cli flushdb
```

##### PostgreSQL Database Issues
```bash
# Check PostgreSQL status
pg_isready -h localhost -p 5432

# Check database size
psql -h localhost -U deshchain_user -d deshchain_identity \
  -c "SELECT pg_size_pretty(pg_database_size('deshchain_identity'));"

# Check slow queries
psql -h localhost -U deshchain_user -d deshchain_identity \
  -c "SELECT query, mean_exec_time FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;"

# Vacuum and analyze tables
psql -h localhost -U deshchain_user -d deshchain_identity -c "VACUUM ANALYZE;"
```

##### ZK Proof Generation Failures
```bash
# Verify ZK proving keys
ls -la /etc/deshchain/zk-proving-keys/

# Test ZK proof generation
./scripts/identity/test-zk-proof.sh

# Regenerate ZK keys if corrupted
./scripts/identity/regenerate-zk-keys.sh

# Check circom installation
circom --version
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
- **Identity Response Time**: <50ms (95th percentile)
- **Biometric Authentication Success**: >99.5%
- **KYC Verification Time**: <30 seconds
- **ZK Proof Generation**: <200ms

### Business KPIs
- **Active Users**: Monthly growth
- **Transaction Volume**: Daily value
- **Validator Participation**: >67% bonded tokens
- **Governance Participation**: >40% voting power
- **NGO Donations**: Monthly distribution tracking
- **Identity Adoption**: DIDs created per month
- **Privacy Compliance**: 100% GDPR/DPDP adherence
- **Cross-Module Identity Usage**: Integration success rate

---

**🎉 Congratulations!** You're now ready to deploy DeshChain to production. Remember to follow security best practices and maintain regular backups.

For additional support, reach out to the DeshChain community through our official channels.