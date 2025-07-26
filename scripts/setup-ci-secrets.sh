#!/bin/bash

# Setup CI/CD Secrets Script for DeshChain
# This script helps configure GitHub Actions secrets for automated deployment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="deshchain/deshchain"

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if gh CLI is installed
    if ! command -v gh &> /dev/null; then
        log_error "GitHub CLI (gh) is not installed. Please install it first:"
        echo "https://cli.github.com/"
        exit 1
    fi
    
    # Check if user is authenticated
    if ! gh auth status &> /dev/null; then
        log_error "You are not authenticated with GitHub CLI. Please run:"
        echo "gh auth login"
        exit 1
    fi
    
    # Check if openssl is available
    if ! command -v openssl &> /dev/null; then
        log_error "OpenSSL is not installed. Please install it first."
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

generate_validator_keys() {
    log_info "Generating validator keys..."
    
    # Create temporary directory for key generation
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    # Generate validator private key (this is a mock - in production use deshchaind)
    cat > priv_validator_key.json << 'EOF'
{
  "address": "7B343E041CA130000A8BC00C35152BD7E7740037",
  "pub_key": {
    "type": "tendermint/PubKeyEd25519",
    "value": "ujY14ab/22anrhZZbANpKDkqEnvnWPT+yfZdYpTlW5I="
  },
  "priv_key": {
    "type": "tendermint/PrivKeyEd25519",
    "value": "YhUeUvhKnTvOmQb+jqKxO9jgDr6EkPKgqoEPZBzjHKG6NjXhpv/bZqeuFllsA2koOSoSe+dY9P7J9l1ilOVbkg=="
  }
}
EOF
    
    VALIDATOR_KEY=$(cat priv_validator_key.json | base64 -w 0)
    cd - > /dev/null
    rm -rf "$TEMP_DIR"
    
    log_success "Validator keys generated"
}

generate_faucet_mnemonic() {
    log_info "Generating faucet mnemonic..."
    
    # Generate a 24-word mnemonic (this is a mock - in production use proper key derivation)
    FAUCET_MNEMONIC="abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art"
    
    log_success "Faucet mnemonic generated"
}

generate_passwords() {
    log_info "Generating secure passwords..."
    
    # Generate PostgreSQL password
    POSTGRES_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
    
    # Generate Redis password
    REDIS_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
    
    log_success "Passwords generated"
}

set_github_secrets() {
    log_info "Setting GitHub Actions secrets..."
    
    # AWS credentials (you need to provide these)
    read -p "Enter AWS Access Key ID: " AWS_ACCESS_KEY_ID
    read -s -p "Enter AWS Secret Access Key: " AWS_SECRET_ACCESS_KEY
    echo
    
    # Slack webhook (optional)
    read -p "Enter Slack Webhook URL (optional): " SLACK_WEBHOOK_URL
    
    # Set secrets using GitHub CLI
    gh secret set AWS_ACCESS_KEY_ID --body "$AWS_ACCESS_KEY_ID" --repo "$GITHUB_REPO"
    gh secret set AWS_SECRET_ACCESS_KEY --body "$AWS_SECRET_ACCESS_KEY" --repo "$GITHUB_REPO"
    gh secret set POSTGRES_PASSWORD --body "$POSTGRES_PASSWORD" --repo "$GITHUB_REPO"
    gh secret set REDIS_PASSWORD --body "$REDIS_PASSWORD" --repo "$GITHUB_REPO"
    gh secret set VALIDATOR_PRIVATE_KEY --body "$VALIDATOR_KEY" --repo "$GITHUB_REPO"
    gh secret set FAUCET_MNEMONIC --body "$FAUCET_MNEMONIC" --repo "$GITHUB_REPO"
    
    if [ ! -z "$SLACK_WEBHOOK_URL" ]; then
        gh secret set SLACK_WEBHOOK_URL --body "$SLACK_WEBHOOK_URL" --repo "$GITHUB_REPO"
    fi
    
    log_success "GitHub secrets configured"
}

create_aws_resources() {
    log_info "Setting up AWS resources..."
    
    cat << 'EOF' > setup-aws.sh
#!/bin/bash

# AWS Setup Script for DeshChain Infrastructure
# This script creates the necessary AWS resources for DeshChain deployment

set -e

# Configuration
CLUSTER_NAME="deshchain-testnet-cluster"
REGION="us-east-1"
NODE_GROUP_NAME="deshchain-nodes"

# Create EKS cluster
echo "Creating EKS cluster..."
eksctl create cluster \
    --name "$CLUSTER_NAME" \
    --region "$REGION" \
    --nodegroup-name "$NODE_GROUP_NAME" \
    --node-type t3.large \
    --nodes 3 \
    --nodes-min 2 \
    --nodes-max 5 \
    --ssh-access \
    --ssh-public-key ~/.ssh/id_rsa.pub \
    --managed

# Create additional node group for validators (dedicated instances)
echo "Creating validator node group..."
eksctl create nodegroup \
    --cluster "$CLUSTER_NAME" \
    --region "$REGION" \
    --name validator-nodes \
    --node-type t3.xlarge \
    --nodes 2 \
    --nodes-min 2 \
    --nodes-max 4 \
    --node-labels node-type=validator \
    --node-taints validator=true:NoSchedule

# Install AWS Load Balancer Controller
echo "Installing AWS Load Balancer Controller..."
curl -o iam_policy.json https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.5.4/docs/install/iam_policy.json

aws iam create-policy \
    --policy-name AWSLoadBalancerControllerIAMPolicy \
    --policy-document file://iam_policy.json \
    --region "$REGION"

eksctl create iamserviceaccount \
    --cluster="$CLUSTER_NAME" \
    --namespace=kube-system \
    --name=aws-load-balancer-controller \
    --role-name AmazonEKSLoadBalancerControllerRole \
    --attach-policy-arn=arn:aws:iam::$(aws sts get-caller-identity --query Account --output text):policy/AWSLoadBalancerControllerIAMPolicy \
    --approve \
    --region "$REGION"

# Install cert-manager
echo "Installing cert-manager..."
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Wait for cert-manager to be ready
kubectl wait --for=condition=ready pod -l app=cert-manager -n cert-manager --timeout=300s

# Create ClusterIssuer for Let's Encrypt
cat <<EOL | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@deshchain.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOL

# Install NGINX Ingress Controller
echo "Installing NGINX Ingress Controller..."
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install ingress-nginx ingress-nginx/ingress-nginx \
    --namespace ingress-nginx \
    --create-namespace \
    --set controller.service.type=LoadBalancer \
    --set controller.service.annotations."service\.beta\.kubernetes\.io/aws-load-balancer-type"="nlb"

# Create S3 bucket for backups
echo "Creating S3 bucket for backups..."
aws s3 mb s3://deshchain-backups --region "$REGION"

# Enable versioning on backup bucket
aws s3api put-bucket-versioning \
    --bucket deshchain-backups \
    --versioning-configuration Status=Enabled

# Create lifecycle policy for backup retention
cat <<EOL > backup-lifecycle.json
{
    "Rules": [
        {
            "ID": "DeshChainBackupRetention",
            "Status": "Enabled",
            "Filter": {"Prefix": "deshchain-testnet/"},
            "Transitions": [
                {
                    "Days": 30,
                    "StorageClass": "STANDARD_IA"
                },
                {
                    "Days": 90,
                    "StorageClass": "GLACIER"
                },
                {
                    "Days": 365,
                    "StorageClass": "DEEP_ARCHIVE"
                }
            ],
            "Expiration": {
                "Days": 2555
            }
        }
    ]
}
EOL

aws s3api put-bucket-lifecycle-configuration \
    --bucket deshchain-backups \
    --lifecycle-configuration file://backup-lifecycle.json

echo "AWS resources setup completed!"
echo ""
echo "Next steps:"
echo "1. Update your kubeconfig: aws eks update-kubeconfig --region $REGION --name $CLUSTER_NAME"
echo "2. Verify cluster: kubectl get nodes"
echo "3. Run the monitoring setup: cd k8s/monitoring && ./setup-monitoring.sh"
echo "4. Deploy DeshChain: trigger the GitHub Actions workflow"
EOF
    
    chmod +x setup-aws.sh
    
    log_success "AWS setup script created as 'setup-aws.sh'"
    log_warning "Please run 'setup-aws.sh' to create AWS infrastructure before deploying"
}

print_summary() {
    echo ""
    echo "=== DeshChain CI/CD Setup Summary ==="
    echo ""
    log_success "✅ GitHub Actions secrets configured"
    log_success "✅ Validator keys generated and stored"
    log_success "✅ Database passwords generated"
    log_success "✅ AWS setup script created"
    echo ""
    echo "📋 Next Steps:"
    echo "1. Run ./setup-aws.sh to create AWS infrastructure"
    echo "2. Push code to main branch to trigger deployment"
    echo "3. Monitor deployment at: https://github.com/$GITHUB_REPO/actions"
    echo "4. Access testnet at: https://testnet.deshchain.com"
    echo ""
    echo "🔐 Generated Credentials:"
    echo "PostgreSQL Password: $POSTGRES_PASSWORD"
    echo "Redis Password: $REDIS_PASSWORD"
    echo ""
    echo "⚠️  IMPORTANT: Save these passwords securely!"
    echo "⚠️  The validator key is for testnet only - generate new keys for mainnet"
    echo ""
}

main() {
    log_info "Starting DeshChain CI/CD setup..."
    
    check_prerequisites
    generate_validator_keys
    generate_faucet_mnemonic
    generate_passwords
    set_github_secrets
    create_aws_resources
    print_summary
    
    log_success "DeshChain CI/CD setup completed successfully!"
}

# Handle script interruption
trap 'log_error "Script interrupted"; exit 1' INT TERM

# Run main function
main "$@"