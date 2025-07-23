#!/bin/bash

# DeshChain Complete Optimization Suite
# Runs all optimization tools in sequence for comprehensive analysis

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configuration
RESULTS_DIR="./complete-optimization-$(date +%Y%m%d_%H%M%S)"
ANALYSIS_DURATION=${ANALYSIS_DURATION:-600}  # 10 minutes
NODE_URL=${NODE_URL:-"http://localhost:26657"}
BINARY=${BINARY:-"deshchaind"}

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
    echo -e "${PURPLE}[OPTIMIZATION] $1${NC}"
}

section() {
    echo -e "${CYAN}$1${NC}"
}

# Create results directory
mkdir -p "$RESULTS_DIR"

highlight "Starting DeshChain Complete Optimization Analysis"
log "Results directory: $RESULTS_DIR"
log "Analysis duration: ${ANALYSIS_DURATION}s"
log "Node URL: $NODE_URL"

# Pre-flight checks
section "======================================================="
section "                  PRE-FLIGHT CHECKS                   "
section "======================================================="

# Check if DeshChain is running
if ! pgrep -f "$BINARY" > /dev/null; then
    error "DeshChain binary '$BINARY' is not running"
    echo "Please start the DeshChain node and re-run this script"
    exit 1
fi

DESHCHAIN_PID=$(pgrep -f "$BINARY" | head -1)
log "✅ DeshChain process found: PID $DESHCHAIN_PID"

# Check node connectivity
if ! curl -sf "$NODE_URL/health" >/dev/null 2>&1; then
    warning "⚠️ Node health check failed at $NODE_URL"
    info "Continuing with optimization analysis..."
else
    log "✅ Node is healthy and responding"
fi

# Check dependencies
MISSING_DEPS=""

if ! command -v go &> /dev/null; then
    warning "Go not found - some profiling features will be limited"
    MISSING_DEPS="$MISSING_DEPS go"
fi

if ! command -v python3 &> /dev/null; then
    warning "Python3 not found - blockchain optimization will be skipped"
    MISSING_DEPS="$MISSING_DEPS python3"
fi

if ! command -v perf &> /dev/null; then
    info "perf not found - install linux-tools-generic for enhanced CPU profiling"
fi

if [ -n "$MISSING_DEPS" ]; then
    warning "Missing optional dependencies:$MISSING_DEPS"
    info "Some optimization features may be limited"
fi

# System information
log "Collecting system information..."
cat > "$RESULTS_DIR/system-overview.txt" << EOF
DeshChain Complete Optimization Analysis
=======================================
Date: $(date)
Hostname: $(hostname)
OS: $(uname -a)
CPU: $(lscpu | grep "Model name" | cut -d: -f2 | sed 's/^ *//' || echo "unknown")
CPU Cores: $(nproc)
Memory: $(free -h | grep "Mem:" | awk '{print $2}')
Load Average: $(uptime | awk -F'load average:' '{print $2}')

DeshChain Process
================
Binary: $BINARY
PID: $DESHCHAIN_PID
Start Time: $(ps -o lstart= -p $DESHCHAIN_PID 2>/dev/null || echo "unknown")
Command: $(ps -o cmd= -p $DESHCHAIN_PID 2>/dev/null || echo "unknown")
Current CPU: $(ps -p $DESHCHAIN_PID -o %cpu= 2>/dev/null | tr -d ' ')%
Current Memory: $(ps -p $DESHCHAIN_PID -o %mem= 2>/dev/null | tr -d ' ')%

Analysis Configuration
=====================
Duration: ${ANALYSIS_DURATION}s
Results Directory: $RESULTS_DIR
Node URL: $NODE_URL
EOF

# Track optimization results
declare -A OPTIMIZATION_RESULTS
TOTAL_TOOLS=4
COMPLETED_TOOLS=0
FAILED_TOOLS=0

# 1. Memory Optimization Analysis
section "======================================================="
section "              MEMORY OPTIMIZATION                     "
section "======================================================="

highlight "Running memory optimization analysis..."

MEMORY_START_TIME=$(date +%s)
if DURATION=$((ANALYSIS_DURATION / 2)) ./scripts/optimization/memory-optimizer.sh 2>&1 | tee "$RESULTS_DIR/memory-optimizer.log"; then
    MEMORY_END_TIME=$(date +%s)
    MEMORY_DURATION=$((MEMORY_END_TIME - MEMORY_START_TIME))
    
    # Find and copy memory results
    MEMORY_RESULTS_DIR=$(find . -maxdepth 1 -name "memory-optimization-*" -type d 2>/dev/null | sort | tail -1)
    if [ -n "$MEMORY_RESULTS_DIR" ] && [ -d "$MEMORY_RESULTS_DIR" ]; then
        cp -r "$MEMORY_RESULTS_DIR"/* "$RESULTS_DIR/"
        rm -rf "$MEMORY_RESULTS_DIR"
        log "✅ Memory optimization completed in ${MEMORY_DURATION}s"
        OPTIMIZATION_RESULTS["memory"]="SUCCESS:${MEMORY_DURATION}"
        ((COMPLETED_TOOLS++))
    else
        warning "⚠️ Memory optimization results not found"
        OPTIMIZATION_RESULTS["memory"]="PARTIAL:${MEMORY_DURATION}"
        ((COMPLETED_TOOLS++))
    fi
else
    error "❌ Memory optimization failed"
    OPTIMIZATION_RESULTS["memory"]="FAILED:0"
    ((FAILED_TOOLS++))
fi

sleep 30  # Cool down between tools

# 2. CPU Performance Profiling
section "======================================================="
section "               CPU PERFORMANCE                        "
section "======================================================="

highlight "Running CPU performance profiling..."

CPU_START_TIME=$(date +%s)
if PROFILE_DURATION=$((ANALYSIS_DURATION / 3)) ./scripts/optimization/cpu-profiler.sh 2>&1 | tee "$RESULTS_DIR/cpu-profiler.log"; then
    CPU_END_TIME=$(date +%s)
    CPU_DURATION=$((CPU_END_TIME - CPU_START_TIME))
    
    # Find and copy CPU results
    CPU_RESULTS_DIR=$(find . -maxdepth 1 -name "cpu-profiling-*" -type d 2>/dev/null | sort | tail -1)
    if [ -n "$CPU_RESULTS_DIR" ] && [ -d "$CPU_RESULTS_DIR" ]; then
        # Rename files to avoid conflicts
        for file in "$CPU_RESULTS_DIR"/*; do
            if [ -f "$file" ]; then
                basename_file=$(basename "$file")
                cp "$file" "$RESULTS_DIR/cpu-${basename_file}"
            fi
        done
        rm -rf "$CPU_RESULTS_DIR"
        log "✅ CPU profiling completed in ${CPU_DURATION}s"
        OPTIMIZATION_RESULTS["cpu"]="SUCCESS:${CPU_DURATION}"
        ((COMPLETED_TOOLS++))
    else
        warning "⚠️ CPU profiling results not found"
        OPTIMIZATION_RESULTS["cpu"]="PARTIAL:${CPU_DURATION}"
        ((COMPLETED_TOOLS++))
    fi
else
    error "❌ CPU profiling failed"
    OPTIMIZATION_RESULTS["cpu"]="FAILED:0"
    ((FAILED_TOOLS++))
fi

sleep 30  # Cool down between tools

# 3. Go Application Profiling
section "======================================================="
section "            APPLICATION PROFILING                     "
section "======================================================="

highlight "Running Go application profiling..."

if command -v go &> /dev/null; then
    APP_START_TIME=$(date +%s)
    
    APP_PROFILE_DIR="$RESULTS_DIR/app-profiles"
    mkdir -p "$APP_PROFILE_DIR"
    
    if timeout $((ANALYSIS_DURATION / 2 + 60)) go run scripts/optimization/performance-profiler.go \
        -node="$NODE_URL" \
        -duration=$((ANALYSIS_DURATION / 2))s \
        -output="$APP_PROFILE_DIR" \
        -cpuprofile="$APP_PROFILE_DIR/app-cpu.prof" \
        -memprofile="$APP_PROFILE_DIR/app-mem.prof" \
        -blockprofile="$APP_PROFILE_DIR/app-block.prof" \
        -mutexprofile="$APP_PROFILE_DIR/app-mutex.prof" \
        2>&1 | tee "$RESULTS_DIR/app-profiler.log"; then
        
        APP_END_TIME=$(date +%s)
        APP_DURATION=$((APP_END_TIME - APP_START_TIME))
        log "✅ Application profiling completed in ${APP_DURATION}s"
        OPTIMIZATION_RESULTS["application"]="SUCCESS:${APP_DURATION}"
        ((COMPLETED_TOOLS++))
    else
        error "❌ Application profiling failed or timed out"
        OPTIMIZATION_RESULTS["application"]="FAILED:0"
        ((FAILED_TOOLS++))
    fi
else
    warning "⚠️ Go not available - skipping application profiling"
    OPTIMIZATION_RESULTS["application"]="SKIPPED:0"
fi

sleep 30  # Cool down between tools

# 4. Blockchain Optimization
section "======================================================="
section "           BLOCKCHAIN OPTIMIZATION                    "
section "======================================================="

highlight "Running blockchain-specific optimization..."

if command -v python3 &> /dev/null; then
    BLOCKCHAIN_START_TIME=$(date +%s)
    
    if timeout $((ANALYSIS_DURATION / 2 + 60)) python3 scripts/optimization/blockchain-optimizer.py \
        --node="$NODE_URL" \
        --duration=$((ANALYSIS_DURATION / 2)) \
        --interval=15 \
        --output="$RESULTS_DIR" \
        2>&1 | tee "$RESULTS_DIR/blockchain-optimizer.log"; then
        
        BLOCKCHAIN_END_TIME=$(date +%s)
        BLOCKCHAIN_DURATION=$((BLOCKCHAIN_END_TIME - BLOCKCHAIN_START_TIME))
        log "✅ Blockchain optimization completed in ${BLOCKCHAIN_DURATION}s"
        OPTIMIZATION_RESULTS["blockchain"]="SUCCESS:${BLOCKCHAIN_DURATION}"
        ((COMPLETED_TOOLS++))
    else
        error "❌ Blockchain optimization failed or timed out"
        OPTIMIZATION_RESULTS["blockchain"]="FAILED:0"
        ((FAILED_TOOLS++))
    fi
else
    warning "⚠️ Python3 not available - skipping blockchain optimization"
    OPTIMIZATION_RESULTS["blockchain"]="SKIPPED:0"
fi

# Generate comprehensive summary report
section "======================================================="
section "           GENERATING SUMMARY REPORT                  "
section "======================================================="

highlight "Generating comprehensive optimization summary..."

SUMMARY_FILE="$RESULTS_DIR/complete-optimization-summary.md"

cat > "$SUMMARY_FILE" << EOF
# DeshChain Complete Optimization Analysis Summary

**Generated:** $(date)  
**Analysis Duration:** ${ANALYSIS_DURATION}s  
**Process PID:** $DESHCHAIN_PID  
**Node URL:** $NODE_URL  

## Executive Summary

This report consolidates results from all DeshChain optimization tools, providing a comprehensive view of system performance and optimization opportunities.

### Analysis Results

- **Total Tools:** $TOTAL_TOOLS
- **Completed Successfully:** $COMPLETED_TOOLS
- **Failed:** $FAILED_TOOLS
- **Overall Success Rate:** $(( (COMPLETED_TOOLS * 100) / TOTAL_TOOLS ))%

## Tool Results

EOF

# Add results for each tool
for tool in memory cpu application blockchain; do
    if [[ -n "${OPTIMIZATION_RESULTS[$tool]}" ]]; then
        IFS=':' read -r status duration <<< "${OPTIMIZATION_RESULTS[$tool]}"
        case $status in
            "SUCCESS")
                echo "### ✅ ${tool^} Optimization" >> "$SUMMARY_FILE"
                echo "**Status:** Completed successfully in ${duration}s" >> "$SUMMARY_FILE"
                ;;
            "PARTIAL")
                echo "### ⚠️ ${tool^} Optimization" >> "$SUMMARY_FILE"
                echo "**Status:** Partially completed in ${duration}s" >> "$SUMMARY_FILE"
                ;;
            "FAILED")
                echo "### ❌ ${tool^} Optimization" >> "$SUMMARY_FILE"
                echo "**Status:** Failed" >> "$SUMMARY_FILE"
                ;;
            "SKIPPED")
                echo "### ⏭️ ${tool^} Optimization" >> "$SUMMARY_FILE"
                echo "**Status:** Skipped (dependencies not available)" >> "$SUMMARY_FILE"
                ;;
        esac
        echo "" >> "$SUMMARY_FILE"
    fi
done

cat >> "$SUMMARY_FILE" << EOF

## Critical Findings

EOF

# Analyze results and extract critical findings
CRITICAL_FINDINGS=0

# Check memory analysis
if [ -f "$RESULTS_DIR/memory-analysis.json" ]; then
    POTENTIAL_LEAK=$(python3 -c "
import json, sys
try:
    with open('$RESULTS_DIR/memory-analysis.json', 'r') as f:
        data = json.load(f)
    print(data.get('potential_leak', False))
except:
    print('False')
" 2>/dev/null || echo "False")
    
    if [ "$POTENTIAL_LEAK" = "True" ]; then
        echo "🚨 **MEMORY LEAK DETECTED** - Immediate attention required" >> "$SUMMARY_FILE"
        ((CRITICAL_FINDINGS++))
    fi
fi

# Check blockchain optimization
if [ -f "$RESULTS_DIR/blockchain-optimization-plan.json" ]; then
    HIGH_PRIORITY_COUNT=$(python3 -c "
import json, sys
try:
    with open('$RESULTS_DIR/blockchain-optimization-plan.json', 'r') as f:
        data = json.load(f)
    print(len(data.get('recommendations', {}).get('high_priority', [])))
except:
    print('0')
" 2>/dev/null || echo "0")
    
    if [ "$HIGH_PRIORITY_COUNT" -gt 0 ]; then
        echo "⚠️ **$HIGH_PRIORITY_COUNT High-Priority Blockchain Issues** - Review blockchain optimization report" >> "$SUMMARY_FILE"
        ((CRITICAL_FINDINGS++))
    fi
fi

if [ $CRITICAL_FINDINGS -eq 0 ]; then
    echo "✅ No critical performance issues identified" >> "$SUMMARY_FILE"
fi

cat >> "$SUMMARY_FILE" << EOF

## Generated Reports

### Core Analysis
- 📊 **System Overview:** \`system-overview.txt\`
- 📋 **This Summary:** \`complete-optimization-summary.md\`

### Memory Analysis
$([ -f "$RESULTS_DIR/memory-optimization-summary.md" ] && echo "- 🧠 **Memory Summary:** \`memory-optimization-summary.md\`")
$([ -f "$RESULTS_DIR/memory-recommendations.md" ] && echo "- 💡 **Memory Recommendations:** \`memory-recommendations.md\`")
$([ -f "$RESULTS_DIR/memory-usage.csv" ] && echo "- 📈 **Memory Data:** \`memory-usage.csv\`")

### CPU Analysis
$([ -f "$RESULTS_DIR/cpu-performance-report.md" ] && echo "- ⚡ **CPU Report:** \`cpu-performance-report.md\`")
$([ -f "$RESULTS_DIR/cpu.prof" ] && echo "- 📊 **CPU Profile:** \`cpu.prof\`")
$([ -f "$RESULTS_DIR/cpu-flamegraph.svg" ] && echo "- 🔥 **Flame Graph:** \`cpu-flamegraph.svg\`")

### Application Profiling
$([ -d "$RESULTS_DIR/app-profiles" ] && echo "- 🔧 **App Profiles:** \`app-profiles/\`")
$([ -f "$RESULTS_DIR/app-profiles/performance-analysis-report.md" ] && echo "- 📊 **App Analysis:** \`app-profiles/performance-analysis-report.md\`")

### Blockchain Optimization
$([ -f "$RESULTS_DIR/blockchain-optimization-report.md" ] && echo "- ⛓️ **Blockchain Report:** \`blockchain-optimization-report.md\`")
$([ -f "$RESULTS_DIR/blockchain-optimization-plan.json" ] && echo "- 📋 **Optimization Plan:** \`blockchain-optimization-plan.json\`")

## Next Steps

### Immediate Actions (Next 24 hours)
$([ $CRITICAL_FINDINGS -gt 0 ] && echo "1. **Address critical findings** identified above")
$([ -f "$RESULTS_DIR/memory-recommendations.md" ] && echo "2. **Review memory recommendations** for immediate optimizations")
$([ -f "$RESULTS_DIR/blockchain-optimization-plan.json" ] && echo "3. **Implement high-priority blockchain optimizations**")

### Short-term Actions (Next week)
1. **Implement performance optimizations** from all tool reports
2. **Set up continuous monitoring** based on identified metrics
3. **Create performance baselines** for tracking improvements
4. **Schedule regular optimization analysis**

### Long-term Actions (Next month)
1. **Implement comprehensive monitoring** dashboard
2. **Automate optimization analysis** in CI/CD pipeline
3. **Performance regression testing** for new releases
4. **Team training** on optimization tools and techniques

## Performance Monitoring Setup

Based on this analysis, set up monitoring for:
- Memory usage patterns and leak detection
- CPU hotspot monitoring
- Blockchain performance metrics (block time, TPS)
- Application profiling integration

---
*Generated by DeshChain Complete Optimization Suite*
EOF

# Create optimization checklist
CHECKLIST_FILE="$RESULTS_DIR/optimization-checklist.md"

cat > "$CHECKLIST_FILE" << EOF
# DeshChain Optimization Implementation Checklist

Use this checklist to track implementation of optimization recommendations.

## Memory Optimizations
$([ -f "$RESULTS_DIR/memory-recommendations.md" ] && echo "- [ ] Review memory recommendations")
- [ ] Implement object pooling for frequently allocated objects
- [ ] Set appropriate GOGC value
- [ ] Configure memory limits (GOMEMLIMIT)
- [ ] Set up memory leak monitoring

## CPU Optimizations
$([ -f "$RESULTS_DIR/cpu-performance-report.md" ] && echo "- [ ] Review CPU performance report")
- [ ] Optimize identified CPU hotspots
- [ ] Implement performance-critical algorithm improvements
- [ ] Add CPU usage monitoring
- [ ] Schedule regular CPU profiling

## Application Optimizations
$([ -d "$RESULTS_DIR/app-profiles" ] && echo "- [ ] Review application profiling results")
- [ ] Optimize high-impact functions
- [ ] Implement caching for expensive operations
- [ ] Add performance regression tests
- [ ] Integrate pprof endpoints for production

## Blockchain Optimizations
$([ -f "$RESULTS_DIR/blockchain-optimization-plan.json" ] && echo "- [ ] Review blockchain optimization plan")
- [ ] Implement consensus optimizations
- [ ] Optimize transaction processing
- [ ] Improve network connectivity
- [ ] Set up blockchain performance monitoring

## Monitoring Setup
- [ ] Set up Grafana dashboards for key metrics
- [ ] Configure alerting for performance thresholds
- [ ] Implement automated optimization analysis
- [ ] Create performance baseline documentation

## Validation
- [ ] Run load tests to validate improvements
- [ ] Compare before/after performance metrics
- [ ] Update performance documentation
- [ ] Schedule follow-up optimization analysis

---
**Implementation Priority:**
1. Address any critical findings first
2. Implement high-impact, low-effort optimizations
3. Set up monitoring and alerting
4. Plan and implement larger optimizations
5. Establish regular optimization review process
EOF

# Final summary
TOTAL_ANALYSIS_TIME=$(( $(date +%s) - $(date -d "$(head -1 "$RESULTS_DIR/system-overview.txt" | grep "Date:" | cut -d: -f2-)" +%s) ))

echo
section "======================================================="
section "              OPTIMIZATION COMPLETE                   "
section "======================================================="
echo

highlight "DeshChain Complete Optimization Analysis Finished!"
echo
echo -e "${CYAN}📊 Analysis Summary:${NC}"
echo -e "  Total Analysis Time: $((TOTAL_ANALYSIS_TIME / 60))m $((TOTAL_ANALYSIS_TIME % 60))s"
echo -e "  Tools Completed: $COMPLETED_TOOLS/$TOTAL_TOOLS"
echo -e "  Success Rate: $(( (COMPLETED_TOOLS * 100) / TOTAL_TOOLS ))%"
echo -e "  Critical Findings: $CRITICAL_FINDINGS"
echo
echo -e "${CYAN}📁 Generated Files:${NC}"
echo -e "  📋 Summary Report: $RESULTS_DIR/complete-optimization-summary.md"
echo -e "  ✅ Checklist: $RESULTS_DIR/optimization-checklist.md"
echo -e "  📊 System Overview: $RESULTS_DIR/system-overview.txt"
echo -e "  📁 All Results: $RESULTS_DIR/"
echo

if [ $CRITICAL_FINDINGS -gt 0 ]; then
    error "⚠️ $CRITICAL_FINDINGS critical issues found - review reports immediately"
    echo "   Priority: Address critical findings before production deployment"
    exit 1
elif [ $FAILED_TOOLS -gt 0 ]; then
    warning "⚠️ $FAILED_TOOLS optimization tools failed - check logs for details"
    echo "   Recommendation: Investigate failed tools and re-run analysis"
    exit 1
else
    log "✅ All optimization analysis completed successfully"
    echo "   Next: Review reports and implement recommended optimizations"
    exit 0
fi