#!/bin/bash

# DeshChain Benchmark Suite
# Comprehensive performance benchmarking for production readiness

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m'

# Configuration
CHAIN_ID=${CHAIN_ID:-"deshchain-testnet-1"}
NODE_ADDRESS=${NODE_ADDRESS:-"tcp://localhost:26657"}
RESULTS_DIR="./benchmark-$(date +%Y%m%d_%H%M%S)"
BENCHMARK_DURATION=${BENCHMARK_DURATION:-300}  # 5 minutes
WARMUP_DURATION=${WARMUP_DURATION:-60}         # 1 minute

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

highlight() {
    echo -e "${PURPLE}[BENCHMARK] $1${NC}"
}

# Create results directory
mkdir -p "$RESULTS_DIR"

highlight "Starting DeshChain Comprehensive Benchmark Suite"
log "Results directory: $RESULTS_DIR"
log "Benchmark duration: ${BENCHMARK_DURATION}s"
log "Warmup duration: ${WARMUP_DURATION}s"

# System information
log "Collecting system information..."
cat > "$RESULTS_DIR/system-info.txt" << EOF
DeshChain Benchmark System Information
=====================================
Date: $(date)
Hostname: $(hostname)
OS: $(uname -a)
CPU: $(lscpu | grep "Model name" | cut -d: -f2 | sed 's/^ *//')
CPU Cores: $(nproc)
Memory: $(free -h | grep "Mem:" | awk '{print $2}')
Disk: $(df -h / | tail -1 | awk '{print $2}')
Network: $(ip route get 8.8.8.8 | grep -oP 'src \K\S+')

Chain Configuration
==================
Chain ID: $CHAIN_ID
Node Address: $NODE_ADDRESS
Binary: $(which deshchaind || echo "not found")
Version: $(deshchaind version 2>/dev/null || echo "unknown")
EOF

# Check prerequisites
log "Checking prerequisites..."
MISSING_DEPS=""

if ! command -v deshchaind &> /dev/null; then
    warning "deshchaind binary not found in PATH"
    MISSING_DEPS="$MISSING_DEPS deshchaind"
fi

if ! command -v go &> /dev/null; then
    warning "Go not found - load testing may fail"
    MISSING_DEPS="$MISSING_DEPS go"
fi

if ! command -v python3 &> /dev/null; then
    warning "Python3 not found - monitoring may fail"
    MISSING_DEPS="$MISSING_DEPS python3"
fi

if [ -n "$MISSING_DEPS" ]; then
    error "Missing dependencies:$MISSING_DEPS"
    echo "Please install missing dependencies and re-run the benchmark"
    exit 1
fi

# Install Python dependencies if needed
if ! python3 -c "import aiohttp, psutil" 2>/dev/null; then
    log "Installing Python dependencies..."
    pip3 install aiohttp psutil 2>/dev/null || warning "Failed to install Python dependencies"
fi

# Pre-benchmark health check
log "Performing pre-benchmark health check..."
HEALTH_STATUS="UNKNOWN"
if curl -sf "${NODE_ADDRESS}/health" >/dev/null 2>&1; then
    HEALTH_STATUS="HEALTHY"
    log "✅ Node is healthy and responding"
else
    HEALTH_STATUS="UNHEALTHY"
    warning "❌ Node health check failed - continuing anyway"
fi

echo "Health Status: $HEALTH_STATUS" >> "$RESULTS_DIR/system-info.txt"

# Start performance monitoring in background
log "Starting performance monitoring..."
python3 scripts/load-testing/performance-monitor.py \
    --node="$NODE_ADDRESS" \
    --duration=$((BENCHMARK_DURATION + WARMUP_DURATION + 120)) \
    --interval=10 \
    --output="$RESULTS_DIR/performance-metrics.jsonl" \
    --max-block-time=15 \
    --min-tps=5 \
    > "$RESULTS_DIR/monitor.log" 2>&1 &
MONITOR_PID=$!

log "Performance monitor started (PID: $MONITOR_PID)"

# Benchmark test scenarios
declare -a BENCHMARKS=(
    "warmup:5:50:Warmup test"
    "baseline:10:100:Baseline performance"
    "moderate:25:200:Moderate load"
    "high:50:300:High throughput"
    "burst:20:500:Burst transactions"
    "sustained:30:1000:Sustained load"
    "extreme:100:200:Extreme concurrency"
)

# Results tracking
declare -A BENCHMARK_RESULTS
TOTAL_BENCHMARKS=${#BENCHMARKS[@]}
FAILED_BENCHMARKS=0

# Function to run individual benchmark
run_benchmark() {
    local name=$1
    local workers=$2
    local tx_per_worker=$3
    local description=$4
    local is_warmup=${5:-false}
    
    if [ "$is_warmup" = "true" ]; then
        info "Running warmup: $description"
    else
        highlight "Running benchmark: $description"
    fi
    
    local start_time=$(date +%s)
    local output_file="$RESULTS_DIR/${name}-results.json"
    local log_file="$RESULTS_DIR/${name}-log.txt"
    
    # Add delay before benchmark
    if [ "$is_warmup" != "true" ]; then
        log "Preparing for benchmark (30s cooldown)..."
        sleep 30
    fi
    
    # Run the load test
    local exit_code=0
    timeout $((BENCHMARK_DURATION + 60)) go run scripts/load-testing/load-test.go \
        -chain-id="$CHAIN_ID" \
        -node="$NODE_ADDRESS" \
        -workers="$workers" \
        -tx-per-worker="$tx_per_worker" \
        -output="$output_file" \
        > "$log_file" 2>&1 || exit_code=$?
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    if [ $exit_code -eq 0 ] && [ -f "$output_file" ]; then
        # Extract results
        local success_rate=$(jq -r '.results.success_rate' "$output_file" 2>/dev/null || echo "0")
        local tps=$(jq -r '.results.tx_per_second' "$output_file" 2>/dev/null || echo "0")
        local avg_latency=$(jq -r '.results.average_latency' "$output_file" 2>/dev/null || echo "0ms")
        local total_tx=$(jq -r '.results.total_tx' "$output_file" 2>/dev/null || echo "0")
        
        if [ "$is_warmup" != "true" ]; then
            log "✅ Benchmark completed - Success: ${success_rate}%, TPS: $tps, Latency: $avg_latency"
            BENCHMARK_RESULTS["$name"]="SUCCESS:${success_rate}:${tps}:${avg_latency}:${total_tx}:${duration}"
        else
            log "✅ Warmup completed - Success: ${success_rate}%, TPS: $tps"
        fi
    else
        if [ "$is_warmup" != "true" ]; then
            error "❌ Benchmark failed"
            BENCHMARK_RESULTS["$name"]="FAILED:0:0:0ms:0:$duration"
            ((FAILED_BENCHMARKS++))
        else
            warning "⚠️ Warmup had issues but continuing"
        fi
    fi
}

# Run warmup
log "Starting warmup phase..."
run_benchmark "warmup" 5 50 "System warmup" true
sleep 30

# Run benchmarks
log "Starting benchmark phase..."
for benchmark in "${BENCHMARKS[@]}"; do
    IFS=':' read -r name workers tx_per_worker description <<< "$benchmark"
    
    # Skip warmup in main loop
    if [ "$name" = "warmup" ]; then
        continue
    fi
    
    run_benchmark "$name" "$workers" "$tx_per_worker" "$description"
done

# Stop monitoring
log "Stopping performance monitor..."
kill $MONITOR_PID 2>/dev/null || true
wait $MONITOR_PID 2>/dev/null || true

# Generate comprehensive report
log "Generating benchmark report..."

cat > "$RESULTS_DIR/benchmark-report.md" << EOF
# DeshChain Comprehensive Benchmark Report

**Date:** $(date)  
**Duration:** ${BENCHMARK_DURATION}s per benchmark  
**System:** $(hostname) - $(lscpu | grep "Model name" | cut -d: -f2 | sed 's/^ *//')  
**Chain ID:** $CHAIN_ID  
**Node:** $NODE_ADDRESS  

## Executive Summary

- **Total Benchmarks:** $TOTAL_BENCHMARKS
- **Successful:** $((TOTAL_BENCHMARKS - FAILED_BENCHMARKS))
- **Failed:** $FAILED_BENCHMARKS
- **Success Rate:** $(( (TOTAL_BENCHMARKS - FAILED_BENCHMARKS) * 100 / TOTAL_BENCHMARKS ))%

## Benchmark Results

| Test | Workers | Tx/Worker | Success Rate | TPS | Avg Latency | Total Tx | Duration |
|------|---------|-----------|--------------|-----|-------------|----------|----------|
EOF

# Add benchmark results to table
for benchmark in "${BENCHMARKS[@]}"; do
    IFS=':' read -r name workers tx_per_worker description <<< "$benchmark"
    
    if [ "$name" = "warmup" ]; then
        continue
    fi
    
    if [[ -n "${BENCHMARK_RESULTS[$name]}" ]]; then
        IFS=':' read -r status success_rate tps latency total_tx duration <<< "${BENCHMARK_RESULTS[$name]}"
        echo "| $description | $workers | $tx_per_worker | ${success_rate}% | $tps | $latency | $total_tx | ${duration}s |" >> "$RESULTS_DIR/benchmark-report.md"
    else
        echo "| $description | $workers | $tx_per_worker | - | - | - | - | - |" >> "$RESULTS_DIR/benchmark-report.md"
    fi
done

cat >> "$RESULTS_DIR/benchmark-report.md" << EOF

## Performance Analysis

### Peak Performance
EOF

# Find peak performance
MAX_TPS=0
BEST_BENCHMARK=""
for benchmark in "${BENCHMARKS[@]}"; do
    IFS=':' read -r name workers tx_per_worker description <<< "$benchmark"
    
    if [ "$name" = "warmup" ] || [[ ! -n "${BENCHMARK_RESULTS[$name]}" ]]; then
        continue
    fi
    
    IFS=':' read -r status success_rate tps latency total_tx duration <<< "${BENCHMARK_RESULTS[$name]}"
    if [ "$status" = "SUCCESS" ] && (( $(echo "$tps > $MAX_TPS" | bc -l 2>/dev/null || echo "0") )); then
        MAX_TPS="$tps"
        BEST_BENCHMARK="$description"
    fi
done

cat >> "$RESULTS_DIR/benchmark-report.md" << EOF
- **Peak TPS:** $MAX_TPS
- **Best Scenario:** $BEST_BENCHMARK
- **Performance Rating:** $(if [ "$FAILED_BENCHMARKS" -eq 0 ]; then echo "🟢 Excellent"; elif [ "$FAILED_BENCHMARKS" -le 2 ]; then echo "🟡 Good"; else echo "🔴 Needs Improvement"; fi)

### Scalability Assessment

EOF

# Scalability analysis
if [ "$FAILED_BENCHMARKS" -eq 0 ]; then
    cat >> "$RESULTS_DIR/benchmark-report.md" << EOF
✅ **All benchmarks passed** - System demonstrates excellent scalability across all tested scenarios.

**Scalability Characteristics:**
- Handles high concurrency (100+ workers)
- Maintains performance under sustained load
- Successfully processes burst transactions
- Stable performance across different load patterns

EOF
else
    cat >> "$RESULTS_DIR/benchmark-report.md" << EOF
⚠️ **Some benchmarks failed** - System shows scalability limitations under certain conditions.

**Performance Limitations Identified:**
- Review failed scenarios for bottlenecks
- Consider infrastructure scaling
- Optimize for high-concurrency scenarios

EOF
fi

cat >> "$RESULTS_DIR/benchmark-report.md" << EOF
## Resource Utilization

See \`performance-metrics.jsonl\` for detailed resource usage during benchmarks.

## Production Readiness Assessment

EOF

# Production readiness
READINESS_SCORE=0
if [ "$FAILED_BENCHMARKS" -eq 0 ]; then
    READINESS_SCORE=$((READINESS_SCORE + 40))
fi

if (( $(echo "$MAX_TPS >= 50" | bc -l 2>/dev/null || echo "0") )); then
    READINESS_SCORE=$((READINESS_SCORE + 30))
fi

if [ "$HEALTH_STATUS" = "HEALTHY" ]; then
    READINESS_SCORE=$((READINESS_SCORE + 20))
fi

READINESS_SCORE=$((READINESS_SCORE + 10))  # Base score

cat >> "$RESULTS_DIR/benchmark-report.md" << EOF
**Production Readiness Score: $READINESS_SCORE/100**

$(if [ "$READINESS_SCORE" -ge 80 ]; then
    echo "🟢 **READY FOR PRODUCTION** - System meets all performance requirements"
elif [ "$READINESS_SCORE" -ge 60 ]; then
    echo "🟡 **CONDITIONAL READINESS** - Address identified issues before production"
else
    echo "🔴 **NOT READY** - Significant performance issues require resolution"
fi)

### Recommendations

1. **Performance Tuning**
   - Target sustained TPS: $MAX_TPS+
   - Monitor resource usage during peak loads
   - Implement performance alerting

2. **Infrastructure Planning**
   - Plan for 2x peak tested load
   - Implement horizontal scaling capabilities
   - Set up comprehensive monitoring

3. **Production Deployment**
   - Gradual rollout with monitoring
   - Load balancing across validators
   - Real-time performance dashboards

## Appendix

### Test Configuration
- Chain ID: $CHAIN_ID
- Node Address: $NODE_ADDRESS
- Benchmark Duration: ${BENCHMARK_DURATION}s per test
- Warmup Duration: ${WARMUP_DURATION}s

### Files Generated
- \`benchmark-report.md\` - This comprehensive report
- \`performance-metrics.jsonl\` - Raw performance data
- \`*-results.json\` - Individual benchmark results
- \`*-log.txt\` - Detailed logs for each test
- \`system-info.txt\` - System configuration details

---
*Report generated by DeshChain Benchmark Suite v1.0*
EOF

# Generate summary JSON
cat > "$RESULTS_DIR/benchmark-summary.json" << EOF
{
  "benchmark_date": "$(date -Iseconds)",
  "chain_id": "$CHAIN_ID",
  "node_address": "$NODE_ADDRESS",
  "system_info": {
    "hostname": "$(hostname)",
    "cpu_cores": $(nproc),
    "memory": "$(free -h | grep "Mem:" | awk '{print $2}')",
    "os": "$(uname -s)"
  },
  "test_configuration": {
    "benchmark_duration": $BENCHMARK_DURATION,
    "warmup_duration": $WARMUP_DURATION,
    "total_benchmarks": $TOTAL_BENCHMARKS
  },
  "results_summary": {
    "successful_benchmarks": $((TOTAL_BENCHMARKS - FAILED_BENCHMARKS)),
    "failed_benchmarks": $FAILED_BENCHMARKS,
    "success_rate": $(( (TOTAL_BENCHMARKS - FAILED_BENCHMARKS) * 100 / TOTAL_BENCHMARKS )),
    "peak_tps": $MAX_TPS,
    "best_scenario": "$BEST_BENCHMARK",
    "production_readiness_score": $READINESS_SCORE
  },
  "benchmark_details": {
EOF

# Add individual benchmark results
first=true
for benchmark in "${BENCHMARKS[@]}"; do
    IFS=':' read -r name workers tx_per_worker description <<< "$benchmark"
    
    if [ "$name" = "warmup" ]; then
        continue
    fi
    
    if [ "$first" = true ]; then
        first=false
    else
        echo "," >> "$RESULTS_DIR/benchmark-summary.json"
    fi
    
    if [[ -n "${BENCHMARK_RESULTS[$name]}" ]]; then
        IFS=':' read -r status success_rate tps latency total_tx duration <<< "${BENCHMARK_RESULTS[$name]}"
        cat >> "$RESULTS_DIR/benchmark-summary.json" << EOF
    "$name": {
      "description": "$description",
      "workers": $workers,
      "tx_per_worker": $tx_per_worker,
      "status": "$status",
      "success_rate": $success_rate,
      "tps": $tps,
      "average_latency": "$latency",
      "total_transactions": $total_tx,
      "duration": $duration
    }EOF
    fi
done

cat >> "$RESULTS_DIR/benchmark-summary.json" << EOF

  }
}
EOF

# Final summary
echo
highlight "DeshChain Benchmark Suite Completed!"
echo
echo -e "${BLUE}📊 Results Summary:${NC}"
echo -e "  📈 Peak TPS: $MAX_TPS"
echo -e "  ✅ Successful Tests: $((TOTAL_BENCHMARKS - FAILED_BENCHMARKS))/$TOTAL_BENCHMARKS"
echo -e "  🎯 Production Readiness: $READINESS_SCORE/100"
echo -e "  🏆 Best Scenario: $BEST_BENCHMARK"
echo
echo -e "${BLUE}📁 Generated Files:${NC}"
echo -e "  📋 Report: $RESULTS_DIR/benchmark-report.md"
echo -e "  📊 Summary: $RESULTS_DIR/benchmark-summary.json"
echo -e "  📈 Metrics: $RESULTS_DIR/performance-metrics.jsonl"
echo -e "  📂 All Results: $RESULTS_DIR/"
echo

# Exit with appropriate code based on production readiness
if [ "$READINESS_SCORE" -ge 80 ]; then
    log "🚀 System is ready for production deployment!"
    exit 0
elif [ "$READINESS_SCORE" -ge 60 ]; then
    warning "⚠️ System needs improvements before production deployment"
    exit 1
else
    error "❌ System requires significant work before production"
    exit 2
fi