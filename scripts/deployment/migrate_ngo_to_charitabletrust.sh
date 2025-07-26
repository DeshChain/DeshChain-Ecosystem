#!/bin/bash

# Migration Script: NGO/Charity to CharitableTrust Module
# This script migrates existing NGO/charity data to the new CharitableTrust module

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
CHAIN_ID=${CHAIN_ID:-"deshchain-mainnet"}
NODE_URL=${NODE_URL:-"http://localhost:26657"}
KEYRING_BACKEND=${KEYRING_BACKEND:-"test"}
MIGRATION_HEIGHT=${MIGRATION_HEIGHT:-"0"} # Set to specific height for upgrade

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Create migration state directory
MIGRATION_DIR="/tmp/charitabletrust_migration_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$MIGRATION_DIR"

log_info "Starting NGO to CharitableTrust migration..."
log_info "Migration data will be stored in: $MIGRATION_DIR"

# Step 1: Export current NGO/charity data
log_info "Step 1: Exporting current NGO/charity data..."

# Query all registered NGOs from donation module
deshchaind query donation organizations --output json > "$MIGRATION_DIR/ngo_organizations.json" || {
    log_error "Failed to export NGO organizations"
    exit 1
}

# Query all donation records
deshchaind query donation donations --output json > "$MIGRATION_DIR/donation_records.json" || {
    log_error "Failed to export donation records"
    exit 1
}

# Query revenue distribution history for NGOs
deshchaind query revenue distributions --output json > "$MIGRATION_DIR/revenue_distributions.json" || {
    log_error "Failed to export revenue distributions"
    exit 1
}

log_success "Exported current NGO/charity data"

# Step 2: Transform data to CharitableTrust format
log_info "Step 2: Transforming data to CharitableTrust format..."

cat > "$MIGRATION_DIR/transform_data.py" <<'PYTHON_EOF'
#!/usr/bin/env python3
import json
import sys
from datetime import datetime, timedelta

def transform_organizations(ngo_data):
    """Transform NGO organizations to CharitableTrust format"""
    transformed = []
    
    for org in ngo_data.get('organizations', []):
        # Map NGO categories to CharitableTrust categories
        category_mapping = {
            'Education': 'education',
            'Healthcare': 'healthcare',
            'Rural Development': 'rural_development',
            'Women Empowerment': 'women_empowerment',
            'Emergency Relief': 'emergency_relief',
            'Environment': 'environmental',
            'Culture': 'cultural_preservation',
            'Skills': 'skill_development'
        }
        
        transformed_org = {
            'id': org['id'],
            'wallet_id': org['wallet_id'],
            'name': org['name'],
            'description': org.get('description', ''),
            'category': category_mapping.get(org.get('category', 'Other'), 'social_welfare'),
            'registration_number': org.get('registration_number', ''),
            'tax_exemption_certificate': org.get('tax_exemption_certificate', ''),
            'verified': org.get('verified', False),
            'verification_date': org.get('verification_date', ''),
            'verification_documents': org.get('verification_documents', []),
            'impact_metrics': {
                'total_beneficiaries': org.get('beneficiaries_served', 0),
                'projects_completed': org.get('projects_completed', 0),
                'funds_utilized': org.get('total_funds_received', '0unamo'),
                'active_projects': org.get('active_projects', 0)
            },
            'contact_info': {
                'address': org.get('address', ''),
                'email': org.get('email', ''),
                'phone': org.get('phone', ''),
                'website': org.get('website', '')
            },
            'bank_details': org.get('bank_details', {}),
            'created_at': org.get('created_at', datetime.now().isoformat()),
            'status': 'active' if org.get('active', True) else 'inactive'
        }
        
        transformed.append(transformed_org)
    
    return transformed

def create_historical_allocations(donation_data, revenue_data):
    """Create historical allocation records from donation and revenue data"""
    allocations = []
    allocation_id = 1
    
    # Process direct donations
    for donation in donation_data.get('donations', []):
        allocation = {
            'id': allocation_id,
            'charitable_org_wallet_id': donation['recipient_org_id'],
            'organization_name': donation.get('recipient_name', ''),
            'amount': donation['amount'],
            'purpose': donation.get('purpose', 'General donation'),
            'category': donation.get('category', 'general'),
            'proposal_id': 0,  # No proposal for direct donations
            'approved_by': ['migration_system'],
            'allocated_at': donation['timestamp'],
            'expected_impact': donation.get('expected_impact', ''),
            'monitoring': {
                'reporting_frequency': 90,
                'required_reports': ['impact', 'financial'],
                'kpis': ['beneficiaries_reached', 'funds_utilized'],
                'monitoring_duration': 180,
                'site_visits_required': False,
                'financial_audit_required': False
            },
            'status': 'completed',
            'distribution': {
                'tx_hash': donation.get('tx_hash', ''),
                'distributed_at': donation['timestamp'],
                'distributed_by': 'migration_system'
            }
        }
        allocations.append(allocation)
        allocation_id += 1
    
    # Process revenue distributions to NGOs
    for distribution in revenue_data.get('distributions', []):
        if distribution.get('recipient_type') == 'ngo':
            allocation = {
                'id': allocation_id,
                'charitable_org_wallet_id': distribution['recipient_id'],
                'organization_name': distribution.get('recipient_name', ''),
                'amount': distribution['amount'],
                'purpose': 'Revenue share distribution',
                'category': 'revenue_share',
                'proposal_id': 0,
                'approved_by': ['revenue_system'],
                'allocated_at': distribution['timestamp'],
                'expected_impact': 'Platform revenue sharing',
                'monitoring': {
                    'reporting_frequency': 30,
                    'required_reports': ['financial'],
                    'kpis': ['funds_utilized'],
                    'monitoring_duration': 90,
                    'site_visits_required': False,
                    'financial_audit_required': True
                },
                'status': 'completed',
                'distribution': {
                    'tx_hash': distribution.get('tx_hash', ''),
                    'distributed_at': distribution['timestamp'],
                    'distributed_by': 'revenue_system'
                }
            }
            allocations.append(allocation)
            allocation_id += 1
    
    return allocations

def calculate_trust_fund_balance(allocations):
    """Calculate the trust fund balance from historical data"""
    total_distributed = 0
    
    for allocation in allocations:
        amount_str = allocation['amount']
        # Extract numeric value from amount string (e.g., "1000000unamo" -> 1000000)
        amount_value = int(''.join(filter(str.isdigit, amount_str.split('unamo')[0])))
        total_distributed += amount_value
    
    return {
        'total_balance': {'denom': 'unamo', 'amount': '0'},  # Will be set from actual balance
        'allocated_amount': {'denom': 'unamo', 'amount': '0'},
        'available_amount': {'denom': 'unamo', 'amount': '0'},
        'total_distributed': {'denom': 'unamo', 'amount': str(total_distributed)},
        'last_updated': datetime.now().isoformat()
    }

def main():
    # Load data
    with open('ngo_organizations.json', 'r') as f:
        ngo_data = json.load(f)
    
    with open('donation_records.json', 'r') as f:
        donation_data = json.load(f)
    
    with open('revenue_distributions.json', 'r') as f:
        revenue_data = json.load(f)
    
    # Transform data
    organizations = transform_organizations(ngo_data)
    allocations = create_historical_allocations(donation_data, revenue_data)
    trust_fund_balance = calculate_trust_fund_balance(allocations)
    
    # Create migration state
    migration_state = {
        'organizations': organizations,
        'allocations': allocations,
        'trust_fund_balance': trust_fund_balance,
        'migration_metadata': {
            'source_module': 'donation',
            'target_module': 'charitabletrust',
            'migration_date': datetime.now().isoformat(),
            'total_organizations': len(organizations),
            'total_allocations': len(allocations),
            'total_distributed': trust_fund_balance['total_distributed']
        }
    }
    
    # Save migration state
    with open('migration_state.json', 'w') as f:
        json.dump(migration_state, f, indent=2)
    
    print(f"Migration state created: {len(organizations)} organizations, {len(allocations)} allocations")

if __name__ == '__main__':
    main()
PYTHON_EOF

cd "$MIGRATION_DIR"
python3 transform_data.py || {
    log_error "Failed to transform data"
    exit 1
}

log_success "Data transformation complete"

# Step 3: Create migration governance proposal
log_info "Step 3: Creating migration governance proposal..."

cat > "$MIGRATION_DIR/migration_proposal.json" <<EOF
{
    "title": "Migrate NGO Module to CharitableTrust Module",
    "description": "This proposal migrates all existing NGO and charity data from the legacy donation module to the new CharitableTrust module. This includes all registered organizations, historical donations, and revenue distributions. The migration ensures continuity of charitable operations while providing enhanced governance, transparency, and impact tracking features.",
    "changes": [
        {
            "subspace": "donation",
            "key": "Enabled",
            "value": "false"
        },
        {
            "subspace": "charitabletrust", 
            "key": "Enabled",
            "value": "true"
        }
    ],
    "migration_plan": {
        "phase1": "Disable donation module and enable CharitableTrust module",
        "phase2": "Migrate all organization data with verification status preserved",
        "phase3": "Create historical allocation records from donations and revenue distributions",
        "phase4": "Update revenue module to use CharitableTrust for distributions",
        "phase5": "Verify migration integrity and activate CharitableTrust operations"
    },
    "deposit": "10000000000unamo"
}
EOF

log_success "Migration proposal created"

# Step 4: Create migration execution script
log_info "Step 4: Creating migration execution script..."

cat > "$MIGRATION_DIR/execute_migration.sh" <<'EXEC_EOF'
#!/bin/bash

# Execute the migration after proposal passes

set -e

echo "Executing CharitableTrust migration..."

# 1. Submit the migration proposal
echo "Submitting migration proposal..."
PROPOSAL_ID=$(deshchaind tx gov submit-proposal param-change migration_proposal.json \
    --from validator \
    --gas-prices 0.025unamo \
    --gas-adjustment 1.5 \
    -y \
    --output json | jq -r '.logs[0].events[] | select(.type=="submit_proposal") | .attributes[] | select(.key=="proposal_id") | .value')

echo "Proposal ID: $PROPOSAL_ID"

# 2. Vote on the proposal (for testnet)
echo "Voting on proposal..."
deshchaind tx gov vote $PROPOSAL_ID yes \
    --from validator \
    --gas-prices 0.025unamo \
    -y

# 3. Wait for proposal to pass
echo "Waiting for proposal to pass..."
sleep 30

# 4. After proposal passes, import migration state
echo "Importing migration state..."
deshchaind tx charitabletrust import-migration migration_state.json \
    --from migration-authority \
    --gas-prices 0.025unamo \
    --gas-adjustment 2.0 \
    -y

echo "Migration execution complete!"
EXEC_EOF

chmod +x "$MIGRATION_DIR/execute_migration.sh"

# Step 5: Create verification script
log_info "Step 5: Creating migration verification script..."

cat > "$MIGRATION_DIR/verify_migration.sh" <<'VERIFY_EOF'
#!/bin/bash

# Verify the migration was successful

echo "Verifying CharitableTrust migration..."

# 1. Check if donation module is disabled
echo "1. Checking donation module status..."
DONATION_STATUS=$(deshchaind query donation params --output json | jq -r '.enabled')
if [ "$DONATION_STATUS" = "false" ]; then
    echo "✓ Donation module is disabled"
else
    echo "✗ Donation module is still enabled!"
fi

# 2. Check if CharitableTrust module is enabled
echo -e "\n2. Checking CharitableTrust module status..."
TRUST_STATUS=$(deshchaind query charitabletrust params --output json | jq -r '.enabled')
if [ "$TRUST_STATUS" = "true" ]; then
    echo "✓ CharitableTrust module is enabled"
else
    echo "✗ CharitableTrust module is not enabled!"
fi

# 3. Verify organization count
echo -e "\n3. Verifying organization migration..."
ORIGINAL_COUNT=$(cat ngo_organizations.json | jq '.organizations | length')
MIGRATED_COUNT=$(deshchaind query charitabletrust allocations --output json | jq '.allocations | map(select(.charitable_org_wallet_id)) | unique_by(.charitable_org_wallet_id) | length')
echo "Original organizations: $ORIGINAL_COUNT"
echo "Migrated organizations: $MIGRATED_COUNT"

# 4. Verify allocation count
echo -e "\n4. Verifying allocation migration..."
ALLOCATION_COUNT=$(deshchaind query charitabletrust allocations --output json | jq '.allocations | length')
echo "Total allocations migrated: $ALLOCATION_COUNT"

# 5. Verify trust fund balance
echo -e "\n5. Verifying trust fund balance..."
deshchaind query charitabletrust trust-fund-balance --output json | jq '.'

# 6. Check revenue module integration
echo -e "\n6. Checking revenue module integration..."
REVENUE_CONFIG=$(deshchaind query revenue params --output json | jq -r '.distribution_targets.charitable_trust')
if [ ! -z "$REVENUE_CONFIG" ]; then
    echo "✓ Revenue module configured for CharitableTrust"
else
    echo "✗ Revenue module not updated!"
fi

echo -e "\nMigration verification complete!"
VERIFY_EOF

chmod +x "$MIGRATION_DIR/verify_migration.sh"

# Step 6: Create rollback script (in case of issues)
log_info "Step 6: Creating rollback script..."

cat > "$MIGRATION_DIR/rollback_migration.sh" <<'ROLLBACK_EOF'
#!/bin/bash

# Rollback migration in case of issues

echo "Rolling back CharitableTrust migration..."

# 1. Create rollback proposal
cat > rollback_proposal.json <<EOF
{
    "title": "Rollback CharitableTrust Migration",
    "description": "Emergency rollback of CharitableTrust migration due to issues",
    "changes": [
        {
            "subspace": "charitabletrust",
            "key": "Enabled", 
            "value": "false"
        },
        {
            "subspace": "donation",
            "key": "Enabled",
            "value": "true"
        }
    ],
    "deposit": "10000000000unamo"
}
EOF

# 2. Submit rollback proposal
ROLLBACK_ID=$(deshchaind tx gov submit-proposal param-change rollback_proposal.json \
    --from validator \
    --gas-prices 0.025unamo \
    --gas-adjustment 1.5 \
    -y \
    --output json | jq -r '.logs[0].events[] | select(.type=="submit_proposal") | .attributes[] | select(.key=="proposal_id") | .value')

echo "Rollback proposal ID: $ROLLBACK_ID"
echo "Please vote on this proposal to complete rollback"

# 3. Export current state for recovery
deshchaind query charitabletrust allocations --output json > charitabletrust_state_backup.json
deshchaind query donation organizations --output json > donation_state_current.json

echo "State backed up for recovery"
echo "Rollback initiated - manual intervention required"
ROLLBACK_EOF

chmod +x "$MIGRATION_DIR/rollback_migration.sh"

# Summary
log_success "Migration preparation complete!"
log_info "Migration directory: $MIGRATION_DIR"
log_info "Generated files:"
echo "  - migration_state.json: Transformed data ready for import"
echo "  - migration_proposal.json: Governance proposal for migration"
echo "  - execute_migration.sh: Script to execute the migration"
echo "  - verify_migration.sh: Script to verify migration success"
echo "  - rollback_migration.sh: Emergency rollback script"

log_warning "IMPORTANT: Review migration_state.json before proceeding!"
log_info "Next steps:"
echo "1. Review the migration state and proposal"
echo "2. Run execute_migration.sh to start the migration"
echo "3. Run verify_migration.sh after migration completes"
echo "4. Keep rollback_migration.sh ready in case of issues"