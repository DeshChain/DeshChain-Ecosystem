#!/bin/bash

# DeshChain Security Audit Script
# Comprehensive security analysis for DeshChain blockchain

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Logging functions
log() {
    echo -e "${GREEN}[$(date '+%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

# Configuration
REPORT_DIR="./security-audit-$(date +%Y%m%d_%H%M%S)"
BINARY="deshchaind"

# Create report directory
mkdir -p "$REPORT_DIR"

# Initialize report
cat > "$REPORT_DIR/audit-report.md" << EOF
# DeshChain Security Audit Report
**Date:** $(date)
**Version:** $(git describe --tags 2>/dev/null || echo "unknown")
**Commit:** $(git rev-parse HEAD 2>/dev/null || echo "unknown")

## Executive Summary
This report contains the results of automated security analysis of DeshChain.

EOF

log "Starting DeshChain security audit..."
log "Report directory: $REPORT_DIR"

# 1. Code Analysis with gosec
log "Running gosec security scanner..."
if command -v gosec >/dev/null 2>&1; then
    gosec -fmt json -out "$REPORT_DIR/gosec-results.json" ./... 2>/dev/null || true
    gosec -fmt text ./... > "$REPORT_DIR/gosec-results.txt" 2>/dev/null || true
    
    # Parse results
    if [ -f "$REPORT_DIR/gosec-results.json" ]; then
        GOSEC_ISSUES=$(jq '.Stats.found // 0' "$REPORT_DIR/gosec-results.json" 2>/dev/null || echo "0")
        if [ "$GOSEC_ISSUES" -gt 0 ]; then
            warning "Found $GOSEC_ISSUES security issues with gosec"
        else
            log "No security issues found with gosec"
        fi
    fi
    
    cat >> "$REPORT_DIR/audit-report.md" << EOF

## Code Security Analysis (gosec)
- **Issues Found:** $GOSEC_ISSUES
- **Detailed Report:** See gosec-results.txt

EOF
else
    warning "gosec not installed, skipping code analysis"
fi

# 2. Dependency vulnerability scan
log "Scanning dependencies for vulnerabilities..."
if command -v govulncheck >/dev/null 2>&1; then
    govulncheck ./... > "$REPORT_DIR/vulnerability-scan.txt" 2>&1 || true
    
    if grep -q "No vulnerabilities found" "$REPORT_DIR/vulnerability-scan.txt"; then
        log "No known vulnerabilities in dependencies"
        VULN_STATUS="✅ Clean"
    else
        warning "Potential vulnerabilities found in dependencies"
        VULN_STATUS="⚠️ Issues found"
    fi
    
    cat >> "$REPORT_DIR/audit-report.md" << EOF

## Dependency Vulnerability Scan
- **Status:** $VULN_STATUS
- **Detailed Report:** See vulnerability-scan.txt

EOF
else
    warning "govulncheck not installed, skipping vulnerability scan"
fi

# 3. Binary analysis
log "Analyzing binary security features..."
if [ -f "bin/$BINARY" ]; then
    # Check if binary is statically linked
    if ldd "bin/$BINARY" >/dev/null 2>&1; then
        STATIC_LINKED="❌ Dynamically linked"
        warning "Binary is dynamically linked"
    else
        STATIC_LINKED="✅ Statically linked"
        log "Binary is statically linked"
    fi
    
    # Check for debug symbols
    if objdump -h "bin/$BINARY" 2>/dev/null | grep -q "debug"; then
        DEBUG_SYMBOLS="❌ Debug symbols present"
        warning "Binary contains debug symbols"
    else
        DEBUG_SYMBOLS="✅ Debug symbols stripped"
        log "Debug symbols stripped from binary"
    fi
    
    # Check for stack protection
    if readelf -s "bin/$BINARY" 2>/dev/null | grep -q "__stack_chk_fail"; then
        STACK_PROTECTION="✅ Stack protection enabled"
        log "Stack protection is enabled"
    else
        STACK_PROTECTION="❌ Stack protection not detected"
        warning "Stack protection not detected"
    fi
    
    cat >> "$REPORT_DIR/audit-report.md" << EOF

## Binary Security Analysis
- **Static Linking:** $STATIC_LINKED
- **Debug Symbols:** $DEBUG_SYMBOLS
- **Stack Protection:** $STACK_PROTECTION

EOF
else
    warning "Binary not found at bin/$BINARY, skipping binary analysis"
fi

# 4. Configuration security
log "Analyzing configuration security..."
CONFIG_ISSUES=0

# Check for default ports
if grep -r "26656\|26657\|1317\|9090" config/ 2>/dev/null | grep -v "#"; then
    warning "Default ports detected in configuration"
    ((CONFIG_ISSUES++))
fi

# Check for localhost bindings
if grep -r "localhost\|127.0.0.1" config/ 2>/dev/null | grep -v "#"; then
    info "Localhost bindings found (good for security)"
else
    warning "No localhost bindings found - services may be exposed"
    ((CONFIG_ISSUES++))
fi

# Check for empty passwords or keys
if grep -ri "password.*=.*\"\"\|key.*=.*\"\"" config/ 2>/dev/null; then
    error "Empty passwords or keys found in configuration"
    ((CONFIG_ISSUES++))
fi

cat >> "$REPORT_DIR/audit-report.md" << EOF

## Configuration Security
- **Issues Found:** $CONFIG_ISSUES
- **Default Ports:** $(if grep -r "26656\|26657\|1317\|9090" config/ 2>/dev/null >/dev/null; then echo "⚠️ Detected"; else echo "✅ Not detected"; fi)
- **Localhost Bindings:** $(if grep -r "localhost\|127.0.0.1" config/ 2>/dev/null >/dev/null; then echo "✅ Present"; else echo "❌ Not found"; fi)

EOF

# 5. Cryptographic analysis
log "Analyzing cryptographic implementations..."
CRYPTO_ISSUES=0

# Check for weak random number generation
if grep -r "math/rand" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/"; then
    error "Weak random number generation detected (math/rand)"
    ((CRYPTO_ISSUES++))
fi

# Check for hardcoded cryptographic keys or seeds
if grep -ri "BEGIN.*PRIVATE\|BEGIN.*RSA\|private.*key.*=\|seed.*=" . --include="*.go" | grep -v "_test.go" | head -5; then
    warning "Potential hardcoded cryptographic material found"
    ((CRYPTO_ISSUES++))
fi

# Check for deprecated crypto functions
if grep -r "md5\|sha1\|des\|rc4" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" | head -5; then
    warning "Deprecated cryptographic functions detected"
    ((CRYPTO_ISSUES++))
fi

cat >> "$REPORT_DIR/audit-report.md" << EOF

## Cryptographic Security
- **Issues Found:** $CRYPTO_ISSUES
- **Weak RNG:** $(if grep -r "math/rand" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then echo "⚠️ Detected"; else echo "✅ Not detected"; fi)
- **Hardcoded Keys:** $(if grep -ri "BEGIN.*PRIVATE\|BEGIN.*RSA" . --include="*.go" | grep -v "_test.go" >/dev/null; then echo "⚠️ Possible"; else echo "✅ Not detected"; fi)

EOF

# 6. Access control analysis
log "Analyzing access controls..."
ACCESS_ISSUES=0

# Check for admin/root checks
if ! grep -r "IsAdmin\|IsRoot\|CheckPermission" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then
    warning "No obvious access control checks found"
    ((ACCESS_ISSUES++))
fi

# Check for authentication mechanisms
if ! grep -r "Authenticate\|ValidateToken\|CheckAuth" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then
    info "Authentication mechanisms not detected (may be handled by Cosmos SDK)"
fi

cat >> "$REPORT_DIR/audit-report.md" << EOF

## Access Control Analysis
- **Issues Found:** $ACCESS_ISSUES
- **Access Checks:** $(if grep -r "IsAdmin\|IsRoot\|CheckPermission" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then echo "✅ Present"; else echo "⚠️ Not obvious"; fi)

EOF

# 7. Input validation analysis
log "Analyzing input validation..."
INPUT_ISSUES=0

# Check for SQL injection potential (though we use KV store)
if grep -r "Exec\|Query.*+\|sprintf.*%" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" | head -5; then
    warning "Potential injection vulnerabilities detected"
    ((INPUT_ISSUES++))
fi

# Check for input sanitization
if ! grep -r "Validate\|Sanitize\|Clean" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then
    warning "No obvious input validation functions found"
    ((INPUT_ISSUES++))
fi

cat >> "$REPORT_DIR/audit-report.md" << EOF

## Input Validation Analysis
- **Issues Found:** $INPUT_ISSUES
- **Validation Functions:** $(if grep -r "Validate\|Sanitize" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then echo "✅ Present"; else echo "⚠️ Not obvious"; fi)

EOF

# 8. Network security
log "Analyzing network security..."
NETWORK_ISSUES=0

# Check for TLS configuration
if ! grep -r "TLS\|Certificate\|tls.Config" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then
    warning "No TLS configuration found"
    ((NETWORK_ISSUES++))
fi

# Check for rate limiting
if ! grep -r "RateLimit\|Throttle\|Limit" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then
    warning "No rate limiting detected"
    ((NETWORK_ISSUES++))
fi

cat >> "$REPORT_DIR/audit-report.md" << EOF

## Network Security Analysis
- **Issues Found:** $NETWORK_ISSUES
- **TLS Configuration:** $(if grep -r "TLS\|Certificate" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then echo "✅ Present"; else echo "⚠️ Not detected"; fi)
- **Rate Limiting:** $(if grep -r "RateLimit\|Throttle" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then echo "✅ Present"; else echo "⚠️ Not detected"; fi)

EOF

# 9. Error handling analysis
log "Analyzing error handling..."
ERROR_ISSUES=0

# Check for information disclosure in errors
if grep -r "fmt.Printf.*err\|log.Printf.*err\|panic(err)" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" | head -5; then
    warning "Potential information disclosure in error messages"
    ((ERROR_ISSUES++))
fi

# Check for proper error handling
if grep -r "_, err :=.*; err != nil" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then
    log "Proper error handling patterns detected"
else
    warning "No obvious error handling patterns found"
    ((ERROR_ISSUES++))
fi

cat >> "$REPORT_DIR/audit-report.md" << EOF

## Error Handling Analysis
- **Issues Found:** $ERROR_ISSUES
- **Information Disclosure:** $(if grep -r "fmt.Printf.*err\|log.Printf.*err" . --include="*.go" | grep -v "_test.go" | grep -v "vendor/" >/dev/null; then echo "⚠️ Possible"; else echo "✅ Not detected"; fi)

EOF

# 10. Generate summary
TOTAL_ISSUES=$((CONFIG_ISSUES + CRYPTO_ISSUES + ACCESS_ISSUES + INPUT_ISSUES + NETWORK_ISSUES + ERROR_ISSUES))

if [ "$TOTAL_ISSUES" -eq 0 ]; then
    OVERALL_SCORE="A (Excellent)"
    OVERALL_COLOR="$GREEN"
elif [ "$TOTAL_ISSUES" -le 3 ]; then
    OVERALL_SCORE="B (Good)"
    OVERALL_COLOR="$YELLOW"
elif [ "$TOTAL_ISSUES" -le 6 ]; then
    OVERALL_SCORE="C (Fair)"
    OVERALL_COLOR="$YELLOW"
else
    OVERALL_SCORE="D (Needs Improvement)"
    OVERALL_COLOR="$RED"
fi

cat >> "$REPORT_DIR/audit-report.md" << EOF

## Summary
- **Overall Security Score:** $OVERALL_SCORE
- **Total Issues Found:** $TOTAL_ISSUES
- **Recommendations:** See detailed sections above

## Next Steps
1. Review and fix all identified security issues
2. Implement additional security measures as recommended
3. Consider professional security audit for production deployment
4. Set up continuous security monitoring

---
*This report was generated automatically. Manual review by security experts is recommended.*
EOF

# Display summary
echo
log "Security audit completed!"
echo -e "${OVERALL_COLOR}Overall Security Score: $OVERALL_SCORE${NC}"
echo -e "Total Issues Found: $TOTAL_ISSUES"
echo -e "Detailed report: $REPORT_DIR/audit-report.md"

# Create issues summary for CI
if [ -n "$GITHUB_ACTIONS" ]; then
    echo "security_score=$OVERALL_SCORE" >> "$GITHUB_OUTPUT"
    echo "total_issues=$TOTAL_ISSUES" >> "$GITHUB_OUTPUT"
    echo "report_path=$REPORT_DIR" >> "$GITHUB_OUTPUT"
fi

# Exit with error if critical issues found
if [ "$TOTAL_ISSUES" -gt 6 ]; then
    error "Critical security issues detected. Review required before deployment."
    exit 1
fi

log "Security audit completed successfully"