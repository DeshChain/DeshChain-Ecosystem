#!/bin/bash

# DeshChain Stress Testing Script
# Tests various load scenarios to identify breaking points

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
CHAIN_ID=${CHAIN_ID:-"deshchain-testnet-1"}
NODE_ADDRESS=${NODE_ADDRESS:-"tcp://localhost:26657"}
RESULTS_DIR="./stress-test-results-$(date +%Y%m%d_%H%M%S)"
BINARY=${BINARY:-"./bin/deshchaind"}

log() {
    echo -e "${GREEN}[$(date '+%H:%M:%S')] $1${NC}"
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

# Create results directory
mkdir -p "$RESULTS_DIR"

log "Starting DeshChain stress testing..."
log "Results will be saved to: $RESULTS_DIR"

# Test scenarios
declare -a SCENARIOS=(
    "10:100:baseline"           # 10 workers, 100 tx each - baseline
    "25:100:moderate"           # 25 workers, 100 tx each - moderate load  
    "50:100:high"              # 50 workers, 100 tx each - high load
    "100:100:extreme"          # 100 workers, 100 tx each - extreme load
    "10:1000:burst"            # 10 workers, 1000 tx each - burst test
    "200:50:massive"           # 200 workers, 50 tx each - massive concurrency
)

# Test results tracking
declare -A RESULTS
FAILED_SCENARIOS=0
TOTAL_SCENARIOS=${#SCENARIOS[@]}

# Function to run individual stress test
run_stress_test() {
    local workers=$1
    local tx_per_worker=$2
    local scenario_name=$3
    
    log "Running scenario: $scenario_name ($workers workers, $tx_per_worker tx/worker)"
    
    local start_time=$(date +%s)
    local output_file="$RESULTS_DIR/${scenario_name}-results.json"
    local log_file="$RESULTS_DIR/${scenario_name}-log.txt"
    
    # Run the load test
    if go run scripts/load-testing/load-test.go \
        -chain-id="$CHAIN_ID" \
        -node="$NODE_ADDRESS" \
        -workers="$workers" \
        -tx-per-worker="$tx_per_worker" \
        -output="$output_file" \
        > "$log_file" 2>&1; then
        
        local end_time=$(date +%s)
        local duration=$((end_time - start_time))
        
        # Extract results from JSON
        if [ -f "$output_file" ]; then
            local success_rate=$(jq -r '.results.success_rate' "$output_file" 2>/dev/null || echo "0")
            local tps=$(jq -r '.results.tx_per_second' "$output_file" 2>/dev/null || echo "0")
            
            log "Scenario $scenario_name completed in ${duration}s - Success: ${success_rate}%, TPS: $tps"
            RESULTS["$scenario_name"]="SUCCESS:${success_rate}:${tps}:${duration}"
        else
            error "Results file not found for scenario $scenario_name"
            RESULTS["$scenario_name"]="FAILED:NO_RESULTS:0:$duration"
            ((FAILED_SCENARIOS++))
        fi
    else
        error "Scenario $scenario_name failed to execute"
        RESULTS["$scenario_name"]="FAILED:EXECUTION:0:0"
        ((FAILED_SCENARIOS++))
    fi
}

# Pre-test health check
log "Performing pre-test health check..."
if ! curl -sf "$NODE_ADDRESS/health" >/dev/null 2>&1; then
    warning "Node health check failed - continuing anyway"
fi

# Run all stress test scenarios
for scenario in "${SCENARIOS[@]}"; do
    IFS=':' read -r workers tx_per_worker name <<< "$scenario"
    
    run_stress_test "$workers" "$tx_per_worker" "$name"
    
    # Cool down period between tests
    log "Cooling down for 30 seconds..."
    sleep 30
done

# Generate comprehensive report
log "Generating stress test report..."

cat > "$RESULTS_DIR/stress-test-report.md" << EOF
# DeshChain Stress Test Report

**Date:** $(date)
**Chain ID:** $CHAIN_ID
**Node:** $NODE_ADDRESS
**Total Scenarios:** $TOTAL_SCENARIOS
**Failed Scenarios:** $FAILED_SCENARIOS

## Test Results Summary

| Scenario | Workers | Tx/Worker | Status | Success Rate | TPS | Duration |
|----------|---------|-----------|--------|--------------|-----|----------|
EOF

# Add results to report
for scenario in "${SCENARIOS[@]}"; do
    IFS=':' read -r workers tx_per_worker name <<< "$scenario"
    
    if [[ -n "${RESULTS[$name]}" ]]; then
        IFS=':' read -r status success_rate tps duration <<< "${RESULTS[$name]}"
        echo "| $name | $workers | $tx_per_worker | $status | ${success_rate}% | $tps | ${duration}s |" >> "$RESULTS_DIR/stress-test-report.md"
    else
        echo "| $name | $workers | $tx_per_worker | UNKNOWN | - | - | - |" >> "$RESULTS_DIR/stress-test-report.md"
    fi
done

cat >> "$RESULTS_DIR/stress-test-report.md" << EOF

## Performance Analysis

### Baseline Performance
EOF

# Analyze baseline performance
if [[ -n "${RESULTS[baseline]}" ]]; then
    IFS=':' read -r status success_rate tps duration <<< "${RESULTS[baseline]}"
    cat >> "$RESULTS_DIR/stress-test-report.md" << EOF
- **Success Rate:** ${success_rate}%
- **Throughput:** $tps TPS
- **Duration:** ${duration}s
- **Status:** $(if [ "$status" = "SUCCESS" ]; then echo "✅ Passed"; else echo "❌ Failed"; fi)

EOF
fi

cat >> "$RESULTS_DIR/stress-test-report.md" << EOF
### Scalability Analysis

EOF

# Analyze scalability
declare -A TPS_BY_SCENARIO
for scenario in "${SCENARIOS[@]}"; do
    IFS=':' read -r workers tx_per_worker name <<< "$scenario"
    if [[ -n "${RESULTS[$name]}" ]]; then
        IFS=':' read -r status success_rate tps duration <<< "${RESULTS[$name]}"
        if [ "$status" = "SUCCESS" ]; then
            TPS_BY_SCENARIO["$name"]="$tps"
        fi
    fi
done

# Find highest TPS
MAX_TPS=0
BEST_SCENARIO=""
for scenario in "${!TPS_BY_SCENARIO[@]}"; do
    tps="${TPS_BY_SCENARIO[$scenario]}"
    if (( $(echo "$tps > $MAX_TPS" | bc -l) )); then
        MAX_TPS="$tps"
        BEST_SCENARIO="$scenario"
    fi
done

cat >> "$RESULTS_DIR/stress-test-report.md" << EOF
- **Peak Performance:** $MAX_TPS TPS (scenario: $BEST_SCENARIO)
- **Successful Scenarios:** $((TOTAL_SCENARIOS - FAILED_SCENARIOS))/$TOTAL_SCENARIOS
- **Failure Rate:** $(( (FAILED_SCENARIOS * 100) / TOTAL_SCENARIOS ))%

### Breaking Point Analysis

EOF

# Identify breaking points
if [ $FAILED_SCENARIOS -gt 0 ]; then
    cat >> "$RESULTS_DIR/stress-test-report.md" << EOF
The following scenarios failed, indicating potential breaking points:

EOF
    for scenario in "${SCENARIOS[@]}"; do
        IFS=':' read -r workers tx_per_worker name <<< "$scenario"
        if [[ -n "${RESULTS[$name]}" ]]; then
            IFS=':' read -r status success_rate tps duration <<< "${RESULTS[$name]}"
            if [ "$status" != "SUCCESS" ]; then
                echo "- **$name:** $workers workers, $tx_per_worker tx/worker - $status" >> "$RESULTS_DIR/stress-test-report.md"
            fi
        fi
    done
else
    cat >> "$RESULTS_DIR/stress-test-report.md" << EOF
✅ All scenarios passed! The system handled all tested load levels successfully.

EOF
fi

cat >> "$RESULTS_DIR/stress-test-report.md" << EOF

## Recommendations

### Performance Optimization
1. **Target TPS:** Based on testing, target $MAX_TPS+ TPS for production
2. **Concurrency Limit:** $(if [ $FAILED_SCENARIOS -eq 0 ]; then echo "System handled 200+ concurrent workers successfully"; else echo "Review failed scenarios for concurrency limits"; fi)
3. **Memory Usage:** Monitor memory consumption during high-load scenarios

### Infrastructure Planning
1. **Validator Requirements:** Ensure validators can handle peak load
2. **Network Capacity:** Plan for sustained high-throughput operations
3. **Monitoring:** Implement real-time performance monitoring

### Production Readiness
$(if [ $FAILED_SCENARIOS -eq 0 ]; then
    echo "✅ **System appears ready for production load**"
else
    echo "⚠️ **Address failed scenarios before production deployment**"
fi)

---
*Report generated by DeshChain stress testing framework*
EOF

# Generate CSV for analysis
log "Generating CSV data for analysis..."
cat > "$RESULTS_DIR/stress-test-data.csv" << EOF
Scenario,Workers,TxPerWorker,Status,SuccessRate,TPS,Duration
EOF

for scenario in "${SCENARIOS[@]}"; do
    IFS=':' read -r workers tx_per_worker name <<< "$scenario"
    if [[ -n "${RESULTS[$name]}" ]]; then
        IFS=':' read -r status success_rate tps duration <<< "${RESULTS[$name]}"
        echo "$name,$workers,$tx_per_worker,$status,$success_rate,$tps,$duration" >> "$RESULTS_DIR/stress-test-data.csv"
    fi
done

# Performance benchmark comparison
log "Creating performance benchmark..."
cat > "$RESULTS_DIR/performance-benchmark.json" << EOF
{
  "benchmark_date": "$(date -Iseconds)",
  "chain_id": "$CHAIN_ID",
  "test_scenarios": $TOTAL_SCENARIOS,
  "failed_scenarios": $FAILED_SCENARIOS,
  "peak_tps": $MAX_TPS,
  "best_scenario": "$BEST_SCENARIO",
  "baseline_performance": {
EOF

if [[ -n "${RESULTS[baseline]}" ]]; then
    IFS=':' read -r status success_rate tps duration <<< "${RESULTS[baseline]}"
    cat >> "$RESULTS_DIR/performance-benchmark.json" << EOF
    "success_rate": $success_rate,
    "tps": $tps,
    "duration": $duration,
    "status": "$status"
EOF
fi

cat >> "$RESULTS_DIR/performance-benchmark.json" << EOF
  },
  "scenarios": {
EOF

first=true
for scenario in "${SCENARIOS[@]}"; do
    IFS=':' read -r workers tx_per_worker name <<< "$scenario"
    if [[ -n "${RESULTS[$name]}" ]]; then
        if [ "$first" = true ]; then
            first=false
        else
            echo "," >> "$RESULTS_DIR/performance-benchmark.json"
        fi
        
        IFS=':' read -r status success_rate tps duration <<< "${RESULTS[$name]}"
        cat >> "$RESULTS_DIR/performance-benchmark.json" << EOF
    "$name": {
      "workers": $workers,
      "tx_per_worker": $tx_per_worker,
      "status": "$status",
      "success_rate": $success_rate,
      "tps": $tps,
      "duration": $duration
    }EOF
    fi
done

cat >> "$RESULTS_DIR/performance-benchmark.json" << EOF

  }
}
EOF

# Final summary
echo
log "Stress testing completed!"
echo -e "${BLUE}Results Summary:${NC}"
echo -e "  Total Scenarios: $TOTAL_SCENARIOS"
echo -e "  Successful: $((TOTAL_SCENARIOS - FAILED_SCENARIOS))"
echo -e "  Failed: $FAILED_SCENARIOS"
echo -e "  Peak TPS: $MAX_TPS"
echo -e "  Best Scenario: $BEST_SCENARIO"
echo
echo -e "${BLUE}Files Generated:${NC}"
echo -e "  📊 Report: $RESULTS_DIR/stress-test-report.md"
echo -e "  📈 Data: $RESULTS_DIR/stress-test-data.csv" 
echo -e "  🎯 Benchmark: $RESULTS_DIR/performance-benchmark.json"
echo -e "  📁 All results: $RESULTS_DIR/"

# Exit with appropriate code
if [ $FAILED_SCENARIOS -gt 0 ]; then
    error "Some stress test scenarios failed. Review results before production deployment."
    exit 1
else
    log "All stress test scenarios passed! System is ready for production load."
    exit 0
fi