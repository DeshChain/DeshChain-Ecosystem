#!/bin/bash

# Production deployment script for DeshChain
# Target server: 139.162.191.142

set -e

echo "Starting DeshChain production deployment..."

# Server details
SERVER_IP="139.162.191.142"
SERVER_USER="root"
REMOTE_DIR="/opt/deshchain"

# Create deployment package
echo "Creating deployment package..."
mkdir -p deploy-package/{nginx/ssl,docker}

# Copy nginx configuration and SSL certificates
cp nginx/nginx.conf deploy-package/nginx/
cp nginx/ssl/cert.pem deploy-package/nginx/ssl/
cp nginx/ssl/key.pem deploy-package/nginx/ssl/

# Copy Docker configurations
cp docker-compose.21nodes.yml deploy-package/docker/
cp Dockerfile.mock deploy-package/docker/

# Create production docker-compose override
cat > deploy-package/docker/docker-compose.override.yml << 'EOF'
version: '3.8'

services:
  nginx:
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    restart: always

  landing-page:
    restart: always

  deshchain-explorer-frontend:
    restart: always

  deshchain-explorer-backend:
    restart: always

  deshchain-faucet:
    restart: always

  postgres:
    restart: always

  redis:
    restart: always
EOF

# Create deployment script for remote server
cat > deploy-package/deploy.sh << 'EOF'
#!/bin/bash
set -e

echo "Setting up DeshChain on production server..."

# Install Docker if not present
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com | sh
fi

# Install Docker Compose if not present
if ! command -v docker-compose &> /dev/null; then
    echo "Installing Docker Compose..."
    curl -L "https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
fi

# Create necessary directories
mkdir -p /opt/deshchain/{data,logs,nginx/ssl}

# Copy files to proper locations
cp nginx/nginx.conf /opt/deshchain/nginx/
cp nginx/ssl/* /opt/deshchain/nginx/ssl/
cp docker/* /opt/deshchain/

# Set proper permissions for SSL certificates
chmod 600 /opt/deshchain/nginx/ssl/key.pem
chmod 644 /opt/deshchain/nginx/ssl/cert.pem

# Navigate to deployment directory
cd /opt/deshchain

# Stop any existing containers
docker-compose -f docker-compose.21nodes.yml down || true

# Build and start services
echo "Building Docker images..."
docker build -f Dockerfile.mock -t deshchain:latest .

echo "Starting services..."
docker-compose -f docker-compose.21nodes.yml -f docker-compose.override.yml up -d

# Wait for services to start
sleep 10

# Check status
echo "Checking service status..."
docker-compose -f docker-compose.21nodes.yml ps

echo "Production deployment complete!"
echo "Services are available at:"
echo "  - Main site: https://deshchain.com"
echo "  - Explorer: https://explorer.deshchain.com"
echo "  - Testnet RPC: https://testnet.deshchain.com"
echo "  - Faucet: https://faucet.deshchain.com"
EOF

chmod +x deploy-package/deploy.sh

# Create README for deployment
cat > deploy-package/README.md << 'EOF'
# DeshChain Production Deployment

This package contains everything needed to deploy DeshChain to production.

## Contents
- `nginx/` - Nginx configuration with SSL certificates
- `docker/` - Docker configurations for 21-node testnet
- `deploy.sh` - Automated deployment script

## Deployment Steps
1. Transfer this package to the server
2. Run `./deploy.sh` on the server
3. Configure DNS records as specified

## DNS Configuration Required
- A record: deshchain.com -> 139.162.191.142
- A record: www.deshchain.com -> 139.162.191.142
- A record: explorer.deshchain.com -> 139.162.191.142
- A record: testnet.deshchain.com -> 139.162.191.142
- A record: faucet.deshchain.com -> 139.162.191.142

## SSL Certificates
CloudFlare Origin certificates are included and valid until 2040.
EOF

echo "Deployment package created successfully!"
echo ""
echo "To deploy to production server:"
echo "1. Transfer the deploy-package directory to $SERVER_IP"
echo "   scp -r deploy-package $SERVER_USER@$SERVER_IP:/tmp/"
echo ""
echo "2. SSH to the server and run deployment"
echo "   ssh $SERVER_USER@$SERVER_IP"
echo "   cd /tmp/deploy-package"
echo "   ./deploy.sh"
echo ""
echo "3. Configure DNS records as specified in README.md"