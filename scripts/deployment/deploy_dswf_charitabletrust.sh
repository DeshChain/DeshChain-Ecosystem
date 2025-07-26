#!/bin/bash

# DeshChain DSWF and CharitableTrust Module Deployment Script
# This script handles the deployment of DSWF and CharitableTrust modules
# Including initialization, parameter setting, and governance setup

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CHAIN_ID=${CHAIN_ID:-"deshchain-testnet"}
NODE_URL=${NODE_URL:-"http://localhost:26657"}
KEYRING_BACKEND=${KEYRING_BACKEND:-"test"}
GAS_PRICES=${GAS_PRICES:-"0.025unamo"}
GAS_ADJUSTMENT=${GAS_ADJUSTMENT:-"1.5"}

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to execute transaction
execute_tx() {
    local cmd=$1
    local desc=$2
    
    print_info "Executing: $desc"
    
    if eval "$cmd"; then
        print_success "$desc completed successfully"
    else
        print_error "$desc failed"
        exit 1
    fi
    
    # Wait for transaction to be processed
    sleep 6
}

# Check if deshchaind is available
if ! command -v deshchaind &> /dev/null; then
    print_error "deshchaind could not be found. Please install it first."
    exit 1
fi

print_info "Starting DSWF and CharitableTrust module deployment..."

# 1. Set up DSWF Parameters
print_info "Setting up DSWF module parameters..."

DSWF_PARAMS='{
    "min_fund_balance": {"denom": "unamo", "amount": "1000000000"},
    "max_allocation_percentage": "0.10",
    "min_liquidity_ratio": "0.20",
    "rebalancing_frequency": 90,
    "allocation_categories": [
        "infrastructure",
        "education",
        "healthcare",
        "technology",
        "agriculture",
        "emergency",
        "research",
        "social_welfare"
    ],
    "investment_horizon": 3650,
    "target_return_rate": "0.08",
    "max_risk_score": 5,
    "disbursement_batch_size": 10,
    "reporting_frequency": 30,
    "audit_requirement": true
}'

# Create param change proposal for DSWF
cat > /tmp/dswf_params_proposal.json <<EOF
{
    "title": "Initialize DSWF Module Parameters",
    "description": "Set initial parameters for DeshChain Sovereign Wealth Fund module",
    "changes": [
        {
            "subspace": "dswf",
            "key": "MinFundBalance",
            "value": "{\\"denom\\": \\"unamo\\", \\"amount\\": \\"1000000000\\"}"
        },
        {
            "subspace": "dswf",
            "key": "MaxAllocationPercentage",
            "value": "\\"0.10\\""
        },
        {
            "subspace": "dswf",
            "key": "MinLiquidityRatio",
            "value": "\\"0.20\\""
        },
        {
            "subspace": "dswf",
            "key": "RebalancingFrequency",
            "value": "90"
        },
        {
            "subspace": "dswf",
            "key": "InvestmentHorizon",
            "value": "3650"
        },
        {
            "subspace": "dswf",
            "key": "TargetReturnRate",
            "value": "\\"0.08\\""
        },
        {
            "subspace": "dswf",
            "key": "MaxRiskScore",
            "value": "5"
        }
    ]
}
EOF

# 2. Set up CharitableTrust Parameters
print_info "Setting up CharitableTrust module parameters..."

TRUST_PARAMS='{
    "enabled": true,
    "min_allocation_amount": {"denom": "unamo", "amount": "100000000"},
    "max_monthly_allocation_per_org": {"denom": "unamo", "amount": "100000000000"},
    "proposal_voting_period": 604800,
    "fraud_investigation_period": 30,
    "impact_report_frequency": 30,
    "distribution_categories": [
        "education",
        "healthcare",
        "rural_development",
        "women_empowerment",
        "emergency_relief",
        "skill_development",
        "environmental",
        "cultural_preservation"
    ]
}'

# 3. Initialize DSWF Governance
print_info "Setting up DSWF governance structure..."

# Get governance account address
GOV_ACCOUNT=${GOV_ACCOUNT:-$(deshchaind keys show gov -a --keyring-backend $KEYRING_BACKEND 2>/dev/null || echo "")}

if [ -z "$GOV_ACCOUNT" ]; then
    print_warning "Governance account not found. Please set GOV_ACCOUNT environment variable."
    GOV_ACCOUNT="desh1gov..."
fi

# Create initial fund managers
FUND_MANAGERS='[
    {
        "address": "desh1fundmanager1...",
        "name": "Investment Manager 1",
        "expertise": "Fixed Income and Treasury",
        "added_at": "2025-01-26T00:00:00Z"
    },
    {
        "address": "desh1fundmanager2...",
        "name": "Investment Manager 2",
        "expertise": "Equities and Growth Assets",
        "added_at": "2025-01-26T00:00:00Z"
    },
    {
        "address": "desh1fundmanager3...",
        "name": "Investment Manager 3",
        "expertise": "Alternative Investments",
        "added_at": "2025-01-26T00:00:00Z"
    }
]'

# 4. Initialize CharitableTrust Governance
print_info "Setting up CharitableTrust governance structure..."

# Create initial trustees
TRUSTEES='[
    {
        "address": "desh1trustee1...",
        "name": "Trustee 1",
        "role": "Chairman",
        "expertise": "Social Impact",
        "appointed_at": "2025-01-26T00:00:00Z",
        "term_end_date": "2027-01-26T00:00:00Z"
    },
    {
        "address": "desh1trustee2...",
        "name": "Trustee 2",
        "role": "Secretary",
        "expertise": "Financial Management",
        "appointed_at": "2025-01-26T00:00:00Z",
        "term_end_date": "2027-01-26T00:00:00Z"
    },
    {
        "address": "desh1trustee3...",
        "name": "Trustee 3",
        "role": "Member",
        "expertise": "NGO Operations",
        "appointed_at": "2025-01-26T00:00:00Z",
        "term_end_date": "2027-01-26T00:00:00Z"
    },
    {
        "address": "desh1trustee4...",
        "name": "Trustee 4",
        "role": "Member",
        "expertise": "Healthcare",
        "appointed_at": "2025-01-26T00:00:00Z",
        "term_end_date": "2027-01-26T00:00:00Z"
    },
    {
        "address": "desh1trustee5...",
        "name": "Trustee 5",
        "role": "Member",
        "expertise": "Education",
        "appointed_at": "2025-01-26T00:00:00Z",
        "term_end_date": "2027-01-26T00:00:00Z"
    },
    {
        "address": "desh1trustee6...",
        "name": "Trustee 6",
        "role": "Member",
        "expertise": "Rural Development",
        "appointed_at": "2025-01-26T00:00:00Z",
        "term_end_date": "2027-01-26T00:00:00Z"
    },
    {
        "address": "desh1trustee7...",
        "name": "Trustee 7",
        "role": "Member",
        "expertise": "Technology",
        "appointed_at": "2025-01-26T00:00:00Z",
        "term_end_date": "2027-01-26T00:00:00Z"
    }
]'

# 5. Create Genesis Configuration Files
print_info "Creating genesis configuration files..."

# DSWF Genesis
cat > /tmp/dswf_genesis.json <<EOF
{
    "params": $DSWF_PARAMS,
    "fund_governance": {
        "fund_managers": $FUND_MANAGERS,
        "required_signatures": 2,
        "approval_threshold": "0.67",
        "investment_committee": [
            "desh1committee1...",
            "desh1committee2...",
            "desh1committee3..."
        ],
        "risk_officer": "desh1risk...",
        "compliance_officer": "desh1compliance...",
        "audit_schedule": 180,
        "last_audit": "2025-01-26T00:00:00Z",
        "next_review": "2025-04-26T00:00:00Z"
    },
    "investment_portfolio": {
        "total_value": {"denom": "unamo", "amount": "0"},
        "liquid_assets": {"denom": "unamo", "amount": "0"},
        "invested_assets": {"denom": "unamo", "amount": "0"},
        "reserved_assets": {"denom": "unamo", "amount": "0"},
        "components": [],
        "total_returns": {"denom": "unamo", "amount": "0"},
        "annual_return_rate": "0",
        "risk_score": 3,
        "last_rebalanced": "2025-01-26T00:00:00Z"
    },
    "allocations": [],
    "monthly_reports": [],
    "allocation_count": 0
}
EOF

# CharitableTrust Genesis
cat > /tmp/charitabletrust_genesis.json <<EOF
{
    "params": $TRUST_PARAMS,
    "trust_governance": {
        "trustees": $TRUSTEES,
        "quorum": 4,
        "approval_threshold": "0.571",
        "advisory_committee": [],
        "transparency_officer": "desh1transparency...",
        "next_election": "2027-01-26T00:00:00Z"
    },
    "trust_fund_balance": {
        "total_balance": {"denom": "unamo", "amount": "0"},
        "allocated_amount": {"denom": "unamo", "amount": "0"},
        "available_amount": {"denom": "unamo", "amount": "0"},
        "total_distributed": {"denom": "unamo", "amount": "0"},
        "last_updated": "2025-01-26T00:00:00Z"
    },
    "allocations": [],
    "proposals": [],
    "impact_reports": [],
    "fraud_alerts": [],
    "allocation_count": 0,
    "proposal_count": 0,
    "report_count": 0,
    "alert_count": 0
}
EOF

# 6. Create deployment verification script
print_info "Creating deployment verification script..."

cat > /tmp/verify_deployment.sh <<'VERIFY_EOF'
#!/bin/bash

echo "Verifying DSWF and CharitableTrust deployment..."

# Check DSWF module
echo "1. Checking DSWF module status..."
deshchaind query dswf fund-status --output json

echo ""
echo "2. Checking DSWF governance..."
deshchaind query dswf governance --output json

echo ""
echo "3. Checking DSWF parameters..."
deshchaind query dswf params --output json

# Check CharitableTrust module
echo ""
echo "4. Checking CharitableTrust fund balance..."
deshchaind query charitabletrust trust-fund-balance --output json

echo ""
echo "5. Checking CharitableTrust governance..."
deshchaind query charitabletrust trust-governance --output json

echo ""
echo "6. Checking CharitableTrust parameters..."
deshchaind query charitabletrust params --output json

echo ""
echo "Deployment verification complete!"
VERIFY_EOF

chmod +x /tmp/verify_deployment.sh

# 7. Create operational setup script
print_info "Creating operational setup script..."

cat > /tmp/operational_setup.sh <<'OPS_EOF'
#!/bin/bash

# Operational setup for DSWF and CharitableTrust

echo "Setting up operational configurations..."

# 1. Create initial DSWF allocation proposal
echo "Creating sample DSWF allocation proposal..."
deshchaind tx dswf propose-allocation \
    "Rural Infrastructure Development" \
    "500000000unamo" \
    "infrastructure" \
    "desh1ruraldev..." \
    "Build 50 km rural roads in Maharashtra" \
    "Connect 20 villages to main highways" \
    "1.08" \
    "Low risk government project" \
    --proposers "desh1fundmanager1...,desh1fundmanager2..." \
    --from treasury \
    --gas-prices 0.025unamo \
    --gas-adjustment 1.5 \
    -y

sleep 6

# 2. Create initial CharitableTrust proposal
echo "Creating sample CharitableTrust allocation proposal..."
deshchaind tx charitabletrust create-proposal \
    "Q1 2025 Charity Distribution" \
    "Quarterly distribution to verified charitable organizations" \
    "10000000000unamo" \
    "Support education, healthcare, and rural development initiatives" \
    "Impact 100,000+ beneficiaries across India" \
    --allocations "1,Shiksha Foundation,2000000000unamo,Education for underprivileged,education;2,Arogya Trust,3000000000unamo,Rural healthcare camps,healthcare;3,Gram Vikas,5000000000unamo,Village infrastructure,rural_development" \
    --documents "https://ipfs.io/ipfs/Qm...proposal,https://ipfs.io/ipfs/Qm...budget" \
    --from trustee1 \
    --gas-prices 0.025unamo \
    --gas-adjustment 1.5 \
    -y

echo "Operational setup complete!"
OPS_EOF

chmod +x /tmp/operational_setup.sh

# 8. Summary
print_success "DSWF and CharitableTrust deployment scripts created successfully!"
print_info "Generated files:"
echo "  - /tmp/dswf_genesis.json"
echo "  - /tmp/charitabletrust_genesis.json"
echo "  - /tmp/dswf_params_proposal.json"
echo "  - /tmp/verify_deployment.sh"
echo "  - /tmp/operational_setup.sh"

print_info "Next steps:"
echo "1. Review and update the genesis files with actual addresses"
echo "2. Submit parameter change proposals through governance"
echo "3. Run verify_deployment.sh to check module status"
echo "4. Run operational_setup.sh to create initial allocations"

print_warning "Note: Replace placeholder addresses (desh1...) with actual addresses before deployment"