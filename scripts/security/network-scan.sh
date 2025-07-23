#!/bin/bash

# DeshChain Network Security Scanner
# Scans running DeshChain node for network vulnerabilities

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuration
HOST=${1:-localhost}
REPORT_DIR="./network-scan-$(date +%Y%m%d_%H%M%S)"

log() {
    echo -e "${GREEN}[$(date '+%H:%M:%S')] $1${NC}"
}

warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

# Create report directory
mkdir -p "$REPORT_DIR"

log "Starting network security scan for $HOST"

# DeshChain standard ports
PORTS=(
    "26656:P2P"
    "26657:RPC"
    "26658:ABCI"
    "1317:API"
    "9090:gRPC"
    "9091:gRPC-Web"
    "26660:Prometheus"
)

# Port scan
log "Scanning DeshChain ports..."
echo "# DeshChain Network Security Scan Report" > "$REPORT_DIR/network-report.md"
echo "**Host:** $HOST" >> "$REPORT_DIR/network-report.md"
echo "**Date:** $(date)" >> "$REPORT_DIR/network-report.md"
echo "" >> "$REPORT_DIR/network-report.md"
echo "## Port Scan Results" >> "$REPORT_DIR/network-report.md"

OPEN_PORTS=0
SECURITY_ISSUES=0

for port_info in "${PORTS[@]}"; do
    port=$(echo "$port_info" | cut -d: -f1)
    service=$(echo "$port_info" | cut -d: -f2)
    
    log "Checking port $port ($service)..."
    
    if timeout 5 bash -c "</dev/tcp/$HOST/$port" 2>/dev/null; then
        log "Port $port ($service) is OPEN"
        echo "- **$port ($service):** ✅ OPEN" >> "$REPORT_DIR/network-report.md"
        ((OPEN_PORTS++))
        
        # Security checks for open ports
        case $port in
            26657|1317|9090)
                if [ "$HOST" != "localhost" ] && [ "$HOST" != "127.0.0.1" ]; then
                    warning "Port $port ($service) is exposed externally - potential security risk"
                    ((SECURITY_ISSUES++))
                    echo "  - ⚠️ **Security Risk:** Exposed externally" >> "$REPORT_DIR/network-report.md"
                fi
                ;;
        esac
    else
        log "Port $port ($service) is CLOSED"
        echo "- **$port ($service):** ❌ CLOSED" >> "$REPORT_DIR/network-report.md"
    fi
done

echo "" >> "$REPORT_DIR/network-report.md"
echo "**Summary:** $OPEN_PORTS ports open, $SECURITY_ISSUES security issues" >> "$REPORT_DIR/network-report.md"

# RPC endpoint security checks
if timeout 5 bash -c "</dev/tcp/$HOST/26657" 2>/dev/null; then
    log "Testing RPC endpoint security..."
    echo "" >> "$REPORT_DIR/network-report.md"
    echo "## RPC Security Analysis" >> "$REPORT_DIR/network-report.md"
    
    # Check if RPC allows dangerous methods
    DANGEROUS_METHODS=("unsafe_reset_all" "unsafe_flush_mempool" "dial_seeds" "dial_peers")
    
    for method in "${DANGEROUS_METHODS[@]}"; do
        response=$(curl -s -X POST -H "Content-Type: application/json" \
            -d "{\"jsonrpc\":\"2.0\",\"method\":\"$method\",\"params\":[],\"id\":1}" \
            "http://$HOST:26657" 2>/dev/null || echo "")
        
        if echo "$response" | grep -q "error.*method not found"; then
            log "Dangerous method $method is disabled ✅"
            echo "- **$method:** ✅ Disabled" >> "$REPORT_DIR/network-report.md"
        elif echo "$response" | grep -q "result"; then
            warning "Dangerous method $method is ENABLED ⚠️"
            echo "- **$method:** ⚠️ **ENABLED (Security Risk)**" >> "$REPORT_DIR/network-report.md"
            ((SECURITY_ISSUES++))
        else
            log "Method $method status unknown"
            echo "- **$method:** ❓ Unknown" >> "$REPORT_DIR/network-report.md"
        fi
    done
    
    # Check CORS headers
    cors_response=$(curl -s -I -H "Origin: http://evil.com" "http://$HOST:26657/health" 2>/dev/null || echo "")
    if echo "$cors_response" | grep -qi "access-control-allow-origin: \*"; then
        warning "CORS is configured to allow all origins - potential security risk"
        echo "- **CORS:** ⚠️ **Allows all origins (Security Risk)**" >> "$REPORT_DIR/network-report.md"
        ((SECURITY_ISSUES++))
    else
        log "CORS appears to be properly configured"
        echo "- **CORS:** ✅ Properly configured" >> "$REPORT_DIR/network-report.md"
    fi
fi

# API endpoint security checks
if timeout 5 bash -c "</dev/tcp/$HOST/1317" 2>/dev/null; then
    log "Testing API endpoint security..."
    echo "" >> "$REPORT_DIR/network-report.md"
    echo "## API Security Analysis" >> "$REPORT_DIR/network-report.md"
    
    # Check for sensitive endpoints
    SENSITIVE_ENDPOINTS=("/cosmos/base/node/v1beta1/config" "/cosmos/auth/v1beta1/accounts")
    
    for endpoint in "${SENSITIVE_ENDPOINTS[@]}"; do
        response=$(curl -s "http://$HOST:1317$endpoint" 2>/dev/null || echo "")
        if echo "$response" | grep -q "error\|not found"; then
            log "Sensitive endpoint $endpoint is protected ✅"
            echo "- **$endpoint:** ✅ Protected" >> "$REPORT_DIR/network-report.md"
        else
            warning "Sensitive endpoint $endpoint may be exposed"
            echo "- **$endpoint:** ⚠️ **May be exposed**" >> "$REPORT_DIR/network-report.md"
            ((SECURITY_ISSUES++))
        fi
    done
    
    # Check rate limiting
    log "Testing rate limiting..."
    start_time=$(date +%s)
    for i in {1..10}; do
        curl -s "http://$HOST:1317/cosmos/base/tendermint/v1beta1/node_info" >/dev/null 2>&1
    done
    end_time=$(date +%s)
    duration=$((end_time - start_time))
    
    if [ $duration -lt 2 ]; then
        warning "No rate limiting detected - rapid requests succeeded"
        echo "- **Rate Limiting:** ⚠️ **Not detected (Security Risk)**" >> "$REPORT_DIR/network-report.md"
        ((SECURITY_ISSUES++))
    else
        log "Rate limiting appears to be in effect"
        echo "- **Rate Limiting:** ✅ Detected" >> "$REPORT_DIR/network-report.md"
    fi
fi

# TLS/SSL checks
log "Checking TLS/SSL configuration..."
echo "" >> "$REPORT_DIR/network-report.md"
echo "## TLS/SSL Analysis" >> "$REPORT_DIR/network-report.md"

if command -v openssl >/dev/null 2>&1; then
    # Check if HTTPS is available
    if timeout 5 openssl s_client -connect "$HOST:26657" -verify_return_error </dev/null 2>/dev/null; then
        log "TLS is configured for RPC"
        echo "- **RPC TLS:** ✅ Configured" >> "$REPORT_DIR/network-report.md"
        
        # Check certificate details
        cert_info=$(timeout 5 openssl s_client -connect "$HOST:26657" -servername "$HOST" </dev/null 2>/dev/null | openssl x509 -text -noout 2>/dev/null || echo "")
        if echo "$cert_info" | grep -q "CN="; then
            log "Certificate details available"
        fi
    else
        warning "TLS not configured for RPC - traffic is unencrypted"
        echo "- **RPC TLS:** ⚠️ **Not configured (Security Risk)**" >> "$REPORT_DIR/network-report.md"
        ((SECURITY_ISSUES++))
    fi
else
    warning "OpenSSL not available, skipping TLS checks"
    echo "- **TLS Check:** ❓ OpenSSL not available" >> "$REPORT_DIR/network-report.md"
fi

# Firewall detection
log "Checking firewall configuration..."
echo "" >> "$REPORT_DIR/network-report.md"
echo "## Firewall Analysis" >> "$REPORT_DIR/network-report.md"

# Try to connect to high-numbered ports (should be blocked)
test_ports=(8080 8443 9999)
blocked_ports=0

for test_port in "${test_ports[@]}"; do
    if ! timeout 2 bash -c "</dev/tcp/$HOST/$test_port" 2>/dev/null; then
        ((blocked_ports++))
    fi
done

if [ $blocked_ports -eq ${#test_ports[@]} ]; then
    log "Firewall appears to be active (test ports blocked)"
    echo "- **Firewall Status:** ✅ Active (test ports blocked)" >> "$REPORT_DIR/network-report.md"
else
    warning "Firewall may not be properly configured"
    echo "- **Firewall Status:** ⚠️ **May not be active (Security Risk)**" >> "$REPORT_DIR/network-report.md"
    ((SECURITY_ISSUES++))
fi

# Generate final report
echo "" >> "$REPORT_DIR/network-report.md"
echo "## Security Summary" >> "$REPORT_DIR/network-report.md"
echo "- **Open Ports:** $OPEN_PORTS" >> "$REPORT_DIR/network-report.md"
echo "- **Security Issues:** $SECURITY_ISSUES" >> "$REPORT_DIR/network-report.md"

if [ $SECURITY_ISSUES -eq 0 ]; then
    SECURITY_RATING="A (Excellent)"
    RATING_COLOR="$GREEN"
elif [ $SECURITY_ISSUES -le 2 ]; then
    SECURITY_RATING="B (Good)"
    RATING_COLOR="$YELLOW"
else
    SECURITY_RATING="C (Needs Improvement)"
    RATING_COLOR="$RED"
fi

echo "- **Security Rating:** $SECURITY_RATING" >> "$REPORT_DIR/network-report.md"

echo "" >> "$REPORT_DIR/network-report.md"
echo "## Recommendations" >> "$REPORT_DIR/network-report.md"
echo "1. Only expose necessary ports to the internet" >> "$REPORT_DIR/network-report.md"
echo "2. Use a firewall to restrict access to sensitive ports" >> "$REPORT_DIR/network-report.md"
echo "3. Enable TLS/SSL for all HTTP endpoints" >> "$REPORT_DIR/network-report.md"
echo "4. Implement rate limiting on API endpoints" >> "$REPORT_DIR/network-report.md"
echo "5. Regularly monitor for unauthorized access attempts" >> "$REPORT_DIR/network-report.md"

# Display results
echo
log "Network security scan completed!"
echo -e "${RATING_COLOR}Security Rating: $SECURITY_RATING${NC}"
echo "Open Ports: $OPEN_PORTS"
echo "Security Issues: $SECURITY_ISSUES"
echo "Report saved to: $REPORT_DIR/network-report.md"

# Exit with error if critical issues found
if [ $SECURITY_ISSUES -gt 3 ]; then
    error "Critical network security issues detected!"
    exit 1
fi

log "Network scan completed successfully"