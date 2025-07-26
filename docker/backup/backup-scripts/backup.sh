#!/bin/bash

# DeshChain Backup Script
# Performs comprehensive backup of DeshChain testnet data

set -e

# Configuration
BACKUP_DATE=$(date +%Y%m%d-%H%M%S)
BACKUP_DIR="/tmp/backup-$BACKUP_DATE"
CHAIN_ID="${CHAIN_ID:-deshchain-testnet-1}"
NAMESPACE="${KUBERNETES_NAMESPACE:-default}"
S3_BUCKET="${S3_BUCKET:-deshchain-backups}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

create_backup_directory() {
    log_info "Creating backup directory: $BACKUP_DIR"
    mkdir -p "$BACKUP_DIR"
}

backup_validator_data() {
    log_info "Backing up validator data..."
    
    # Get list of validator pods
    VALIDATOR_PODS=$(kubectl get pods -n "$NAMESPACE" -l app.kubernetes.io/component=validator -o jsonpath='{.items[*].metadata.name}')
    
    for pod in $VALIDATOR_PODS; do
        log_info "Backing up validator pod: $pod"
        
        # Create validator-specific backup directory
        mkdir -p "$BACKUP_DIR/validators/$pod"
        
        # Backup critical validator files
        kubectl exec -n "$NAMESPACE" "$pod" -- tar czf - \
            -C /deshchain/.deshchaind \
            config/genesis.json \
            config/node_key.json \
            config/priv_validator_key.json \
            data/priv_validator_state.json \
            > "$BACKUP_DIR/validators/$pod/validator-data.tar.gz"
        
        # Get validator info
        kubectl exec -n "$NAMESPACE" "$pod" -- deshchaind status > "$BACKUP_DIR/validators/$pod/status.json" || true
        
        log_success "Validator $pod backup completed"
    done
}

backup_postgresql() {
    log_info "Backing up PostgreSQL database..."
    
    # Get PostgreSQL pod
    POSTGRES_POD=$(kubectl get pods -n "$NAMESPACE" -l app.kubernetes.io/name=postgresql -o jsonpath='{.items[0].metadata.name}')
    
    if [ -n "$POSTGRES_POD" ]; then
        # Create database backup
        kubectl exec -n "$NAMESPACE" "$POSTGRES_POD" -- pg_dump \
            -U deshchain \
            -d deshchain_explorer \
            --no-password \
            --format=custom \
            --compress=9 \
            > "$BACKUP_DIR/postgresql-backup.dump"
        
        # Get database size info
        kubectl exec -n "$NAMESPACE" "$POSTGRES_POD" -- psql \
            -U deshchain \
            -d deshchain_explorer \
            -c "SELECT pg_size_pretty(pg_database_size('deshchain_explorer')) as size;" \
            --no-password -t > "$BACKUP_DIR/database-size.txt"
        
        log_success "PostgreSQL backup completed"
    else
        log_warning "PostgreSQL pod not found, skipping database backup"
    fi
}

backup_redis() {
    log_info "Backing up Redis data..."
    
    # Get Redis pod
    REDIS_POD=$(kubectl get pods -n "$NAMESPACE" -l app.kubernetes.io/name=redis -o jsonpath='{.items[0].metadata.name}')
    
    if [ -n "$REDIS_POD" ]; then
        # Trigger Redis save
        kubectl exec -n "$NAMESPACE" "$REDIS_POD" -- redis-cli BGSAVE
        
        # Wait for save to complete
        sleep 5
        
        # Copy Redis dump
        kubectl exec -n "$NAMESPACE" "$REDIS_POD" -- cat /data/dump.rdb > "$BACKUP_DIR/redis-dump.rdb"
        
        log_success "Redis backup completed"
    else
        log_warning "Redis pod not found, skipping Redis backup"
    fi
}

backup_configurations() {
    log_info "Backing up Kubernetes configurations..."
    
    mkdir -p "$BACKUP_DIR/k8s-configs"
    
    # Backup ConfigMaps
    kubectl get configmaps -n "$NAMESPACE" -o yaml > "$BACKUP_DIR/k8s-configs/configmaps.yaml"
    
    # Backup Secrets (names only, not values for security)
    kubectl get secrets -n "$NAMESPACE" -o jsonpath='{.items[*].metadata.name}' > "$BACKUP_DIR/k8s-configs/secret-names.txt"
    
    # Backup Services
    kubectl get services -n "$NAMESPACE" -o yaml > "$BACKUP_DIR/k8s-configs/services.yaml"
    
    # Backup Ingress
    kubectl get ingress -n "$NAMESPACE" -o yaml > "$BACKUP_DIR/k8s-configs/ingress.yaml"
    
    # Backup PersistentVolumeClaims
    kubectl get pvc -n "$NAMESPACE" -o yaml > "$BACKUP_DIR/k8s-configs/pvc.yaml"
    
    log_success "Kubernetes configurations backup completed"
}

create_metadata() {
    log_info "Creating backup metadata..."
    
    cat > "$BACKUP_DIR/backup-metadata.json" << EOF
{
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "backup_id": "$BACKUP_DATE",
    "chain_id": "$CHAIN_ID",
    "namespace": "$NAMESPACE",
    "backup_type": "full",
    "components": {
        "validators": true,
        "postgresql": $([ -f "$BACKUP_DIR/postgresql-backup.dump" ] && echo "true" || echo "false"),
        "redis": $([ -f "$BACKUP_DIR/redis-dump.rdb" ] && echo "true" || echo "false"),
        "configurations": true
    },
    "validator_count": $(ls -1 "$BACKUP_DIR/validators" 2>/dev/null | wc -l),
    "total_size_bytes": $(du -sb "$BACKUP_DIR" | cut -f1),
    "kubernetes_version": "$(kubectl version --client -o json | jq -r '.clientVersion.gitVersion')",
    "backup_version": "1.0.0"
}
EOF
    
    log_success "Backup metadata created"
}

compress_backup() {
    log_info "Compressing backup archive..."
    
    BACKUP_FILE="deshchain-backup-$BACKUP_DATE.tar.gz"
    
    # Create compressed archive
    tar czf "/tmp/$BACKUP_FILE" -C /tmp "$(basename "$BACKUP_DIR")"
    
    # Calculate checksums
    sha256sum "/tmp/$BACKUP_FILE" > "/tmp/$BACKUP_FILE.sha256"
    md5sum "/tmp/$BACKUP_FILE" > "/tmp/$BACKUP_FILE.md5"
    
    log_success "Backup compressed to $BACKUP_FILE"
}

upload_to_s3() {
    log_info "Uploading backup to S3..."
    
    if [ -z "$S3_BUCKET" ]; then
        log_warning "S3_BUCKET not set, skipping S3 upload"
        return
    fi
    
    BACKUP_FILE="deshchain-backup-$BACKUP_DATE.tar.gz"
    S3_PATH="s3://$S3_BUCKET/deshchain-testnet/$BACKUP_FILE"
    
    # Upload backup file
    aws s3 cp "/tmp/$BACKUP_FILE" "$S3_PATH" \
        --metadata backup-date="$BACKUP_DATE",chain-id="$CHAIN_ID",backup-type="full"
    
    # Upload checksums
    aws s3 cp "/tmp/$BACKUP_FILE.sha256" "s3://$S3_BUCKET/deshchain-testnet/$BACKUP_FILE.sha256"
    aws s3 cp "/tmp/$BACKUP_FILE.md5" "s3://$S3_BUCKET/deshchain-testnet/$BACKUP_FILE.md5"
    
    log_success "Backup uploaded to $S3_PATH"
}

cleanup_old_backups() {
    log_info "Cleaning up old backups..."
    
    if [ -z "$S3_BUCKET" ]; then
        log_warning "S3_BUCKET not set, skipping cleanup"
        return
    fi
    
    # List and delete old backups beyond retention period
    CUTOFF_DATE=$(date -d "$RETENTION_DAYS days ago" +%Y%m%d)
    
    aws s3 ls "s3://$S3_BUCKET/deshchain-testnet/" | \
        grep "deshchain-backup-" | \
        awk '{print $4}' | \
        while read -r backup_file; do
            # Extract date from filename
            backup_date=$(echo "$backup_file" | sed -n 's/.*deshchain-backup-\([0-9]\{8\}\)-.*/\1/p')
            
            if [ "$backup_date" -lt "$CUTOFF_DATE" ]; then
                log_info "Deleting old backup: $backup_file"
                aws s3 rm "s3://$S3_BUCKET/deshchain-testnet/$backup_file"
                aws s3 rm "s3://$S3_BUCKET/deshchain-testnet/$backup_file.sha256" 2>/dev/null || true
                aws s3 rm "s3://$S3_BUCKET/deshchain-testnet/$backup_file.md5" 2>/dev/null || true
            fi
        done
    
    log_success "Old backups cleanup completed"
}

cleanup_local_files() {
    log_info "Cleaning up local files..."
    
    # Remove backup directory
    rm -rf "$BACKUP_DIR"
    
    # Remove compressed files
    rm -f "/tmp/deshchain-backup-$BACKUP_DATE.tar.gz"
    rm -f "/tmp/deshchain-backup-$BACKUP_DATE.tar.gz.sha256"
    rm -f "/tmp/deshchain-backup-$BACKUP_DATE.tar.gz.md5"
    
    log_success "Local cleanup completed"
}

send_notification() {
    local status=$1
    local message=$2
    
    if [ -n "$SLACK_WEBHOOK_URL" ]; then
        local color="good"
        local emoji="✅"
        
        if [ "$status" != "success" ]; then
            color="danger"
            emoji="❌"
        fi
        
        curl -X POST -H 'Content-type: application/json' \
            --data "{
                \"attachments\": [{
                    \"color\": \"$color\",
                    \"title\": \"$emoji DeshChain Backup $status\",
                    \"text\": \"$message\",
                    \"fields\": [
                        {\"title\": \"Chain ID\", \"value\": \"$CHAIN_ID\", \"short\": true},
                        {\"title\": \"Backup Date\", \"value\": \"$BACKUP_DATE\", \"short\": true},
                        {\"title\": \"Namespace\", \"value\": \"$NAMESPACE\", \"short\": true}
                    ],
                    \"footer\": \"DeshChain Backup System\",
                    \"ts\": $(date +%s)
                }]
            }" \
            "$SLACK_WEBHOOK_URL" || log_warning "Failed to send Slack notification"
    fi
}

main() {
    log_info "Starting DeshChain backup process..."
    
    # Trap errors and cleanup
    trap 'log_error "Backup failed"; cleanup_local_files; send_notification "failed" "Backup process encountered an error"; exit 1' ERR
    
    create_backup_directory
    backup_validator_data
    backup_postgresql
    backup_redis
    backup_configurations
    create_metadata
    compress_backup
    upload_to_s3
    cleanup_old_backups
    cleanup_local_files
    
    local backup_size=$(du -sh /tmp/deshchain-backup-$BACKUP_DATE.tar.gz 2>/dev/null | cut -f1 || echo "unknown")
    send_notification "success" "Backup completed successfully. Size: $backup_size"
    
    log_success "DeshChain backup process completed successfully!"
}

# Run backup
main "$@"