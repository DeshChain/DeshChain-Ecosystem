#!/bin/bash

# DeshChain CPU Performance Profiler
# Advanced CPU profiling and optimization analysis

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m'

# Configuration
BINARY=${BINARY:-"deshchaind"}
OUTPUT_DIR="./cpu-profiling-$(date +%Y%m%d_%H%M%S)"
PROFILE_DURATION=${PROFILE_DURATION:-120}  # 2 minutes
SAMPLE_FREQUENCY=${SAMPLE_FREQUENCY:-99}   # 99 Hz
ENABLE_PERF=${ENABLE_PERF:-true}
ENABLE_PPROF=${ENABLE_PPROF:-true}

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
    echo -e "${PURPLE}[CPU-PROFILER] $1${NC}"
}

# Create output directory
mkdir -p "$OUTPUT_DIR"

highlight "Starting DeshChain CPU Performance Profiling"
log "Output directory: $OUTPUT_DIR"
log "Profile duration: ${PROFILE_DURATION}s"
log "Sample frequency: ${SAMPLE_FREQUENCY} Hz"

# Check dependencies
check_dependencies() {
    local missing_deps=""
    
    if [ "$ENABLE_PERF" = "true" ] && ! command -v perf &> /dev/null; then
        warning "perf not found - install linux-tools-generic for detailed profiling"
        ENABLE_PERF=false
        missing_deps="$missing_deps perf"
    fi
    
    if ! command -v go &> /dev/null; then
        warning "Go not found - pprof analysis will be limited"
        missing_deps="$missing_deps go"
    fi
    
    if ! command -v python3 &> /dev/null; then
        warning "Python3 not found - advanced analysis will be limited"
        missing_deps="$missing_deps python3"
    fi
    
    if [ -n "$missing_deps" ]; then
        warning "Missing optional dependencies:$missing_deps"
        info "Install them for enhanced profiling capabilities"
    fi
}

# Check if DeshChain is running
find_deshchain_process() {
    if ! pgrep -f "$BINARY" > /dev/null; then
        error "DeshChain binary '$BINARY' is not running"
        echo "Please start the DeshChain node and re-run this script"
        exit 1
    fi
    
    DESHCHAIN_PID=$(pgrep -f "$BINARY" | head -1)
    log "Found DeshChain process: PID $DESHCHAIN_PID"
    
    # Get process info
    cat > "$OUTPUT_DIR/process-info.txt" << EOF
DeshChain Process Information
============================
PID: $DESHCHAIN_PID
Binary: $BINARY
Start Time: $(ps -o lstart= -p $DESHCHAIN_PID 2>/dev/null || echo "unknown")
Command: $(ps -o cmd= -p $DESHCHAIN_PID 2>/dev/null || echo "unknown")
User: $(ps -o user= -p $DESHCHAIN_PID 2>/dev/null || echo "unknown")
Priority: $(ps -o ni= -p $DESHCHAIN_PID 2>/dev/null || echo "unknown")
Threads: $(ps -o nlwp= -p $DESHCHAIN_PID 2>/dev/null || echo "unknown")
EOF
}

# Collect system information
collect_system_info() {
    log "Collecting system information..."
    
    cat > "$OUTPUT_DIR/system-info.txt" << EOF
DeshChain CPU Profiling - System Information
============================================
Date: $(date)
Hostname: $(hostname)
Kernel: $(uname -r)
OS: $(lsb_release -d 2>/dev/null | cut -f2 || uname -o)
Architecture: $(uname -m)

CPU Information
===============
EOF
    
    if command -v lscpu &> /dev/null; then
        lscpu >> "$OUTPUT_DIR/system-info.txt"
    else
        cat /proc/cpuinfo >> "$OUTPUT_DIR/system-info.txt"
    fi
    
    cat >> "$OUTPUT_DIR/system-info.txt" << EOF

Load Average
============
$(uptime)

Current CPU Usage
=================
$(top -bn1 | grep "Cpu(s)" || echo "CPU usage information unavailable")

Memory Information
==================
$(free -h)

Process CPU Usage (Top 10)
==========================
$(ps aux --sort=-%cpu | head -11)
EOF
}

# Run perf profiling
run_perf_profiling() {
    if [ "$ENABLE_PERF" != "true" ]; then
        warning "Perf profiling disabled"
        return
    fi
    
    log "Starting perf CPU profiling..."
    
    local perf_output="$OUTPUT_DIR/perf-record.data"
    local perf_report="$OUTPUT_DIR/perf-report.txt"
    local perf_script="$OUTPUT_DIR/perf-script.txt"
    
    # Record performance data
    info "Recording perf data for ${PROFILE_DURATION}s..."
    if perf record -g -F $SAMPLE_FREQUENCY -p $DESHCHAIN_PID -o "$perf_output" -- sleep $PROFILE_DURATION 2>/dev/null; then
        log "✅ Perf recording completed"
        
        # Generate perf report
        info "Generating perf analysis reports..."
        perf report -i "$perf_output" --stdio > "$perf_report" 2>/dev/null || warning "Failed to generate perf report"
        
        # Generate perf script output for flame graphs
        perf script -i "$perf_output" > "$perf_script" 2>/dev/null || warning "Failed to generate perf script"
        
        # Try to generate flame graph if stackcollapse tools are available
        if command -v stackcollapse-perf.pl &> /dev/null && command -v flamegraph.pl &> /dev/null; then
            info "Generating flame graph..."
            local flamegraph="$OUTPUT_DIR/flamegraph.svg"
            stackcollapse-perf.pl "$perf_script" | flamegraph.pl > "$flamegraph" 2>/dev/null && \
                log "✅ Flame graph generated: $flamegraph" || \
                warning "Failed to generate flame graph"
        else
            info "Install FlameGraph tools for flame graph generation:"
            info "git clone https://github.com/brendangregg/FlameGraph.git"
        fi
        
        # Generate summary
        generate_perf_summary "$perf_report" "$OUTPUT_DIR/perf-summary.md"
        
    else
        error "Failed to run perf record - check permissions and perf installation"
        warning "Try: sudo sysctl -w kernel.perf_event_paranoid=1"
    fi
}

# Generate perf summary
generate_perf_summary() {
    local perf_report_file=$1
    local output_file=$2
    
    if [ ! -f "$perf_report_file" ]; then
        return
    fi
    
    log "Generating perf analysis summary..."
    
    cat > "$output_file" << EOF
# DeshChain CPU Performance Analysis (perf)

**Profile Duration:** ${PROFILE_DURATION}s  
**Sample Frequency:** ${SAMPLE_FREQUENCY} Hz  
**Process PID:** $DESHCHAIN_PID  

## Top CPU-Consuming Functions

\`\`\`
EOF
    
    # Extract top functions from perf report
    head -50 "$perf_report_file" | grep -E "^\s*[0-9]+\.[0-9]+%" | head -20 >> "$output_file" 2>/dev/null || true
    
    cat >> "$output_file" << EOF
\`\`\`

## Performance Hotspots

EOF
    
    # Analyze hotspots
    python3 -c "
import re

try:
    with open('$perf_report_file', 'r') as f:
        content = f.read()
    
    # Extract percentage lines
    lines = content.split('\n')
    hotspots = []
    
    for line in lines:
        match = re.match(r'\s*([0-9]+\.[0-9]+)%.*?([a-zA-Z_][a-zA-Z0-9_:]+)', line)
        if match:
            percent = float(match.group(1))
            function = match.group(2)
            if percent > 1.0:  # Only functions using > 1% CPU
                hotspots.append((percent, function))
    
    if hotspots:
        print('### Critical Hotspots (>5% CPU)')
        for percent, func in hotspots:
            if percent > 5.0:
                print(f'- **{func}**: {percent}% CPU usage')
        
        print()
        print('### Optimization Opportunities (1-5% CPU)')
        for percent, func in hotspots:
            if 1.0 <= percent <= 5.0:
                print(f'- {func}: {percent}% CPU usage')
    else:
        print('No significant hotspots identified.')

except Exception as e:
    print(f'Error analyzing perf data: {e}')
" >> "$output_file" 2>/dev/null || echo "### Analysis unavailable" >> "$output_file"
    
    cat >> "$output_file" << EOF

## Recommendations

1. **Focus on functions using >5% CPU time**
2. **Profile specific algorithms in hot functions**
3. **Consider algorithmic optimizations**
4. **Review memory access patterns**
5. **Implement caching for expensive operations**

---
*Generated from perf analysis*
EOF
}

# Run Go pprof profiling
run_pprof_profiling() {
    if [ "$ENABLE_PPROF" != "true" ]; then
        warning "pprof profiling disabled"
        return
    fi
    
    log "Starting Go pprof CPU profiling..."
    
    local pprof_output="$OUTPUT_DIR/cpu.prof"
    local pprof_report="$OUTPUT_DIR/pprof-report.txt"
    
    # Check if DeshChain exposes pprof endpoint
    local pprof_url="http://localhost:6060/debug/pprof/profile"
    
    info "Attempting to collect pprof CPU profile..."
    
    # Try to get profile from pprof endpoint
    if curl -s --max-time 5 "$pprof_url?seconds=$PROFILE_DURATION" -o "$pprof_output" 2>/dev/null; then
        log "✅ pprof CPU profile collected from HTTP endpoint"
        
        # Analyze the profile
        if command -v go &> /dev/null; then
            info "Generating pprof analysis..."
            
            # Generate text report
            go tool pprof -text "$pprof_output" > "$pprof_report" 2>/dev/null || warning "Failed to generate pprof text report"
            
            # Generate top functions
            go tool pprof -top "$pprof_output" > "$OUTPUT_DIR/pprof-top.txt" 2>/dev/null || warning "Failed to generate pprof top report"
            
            # Generate call tree
            go tool pprof -tree "$pprof_output" > "$OUTPUT_DIR/pprof-tree.txt" 2>/dev/null || warning "Failed to generate pprof tree report"
            
            # Try to generate SVG call graph
            if command -v dot &> /dev/null; then
                go tool pprof -svg "$pprof_output" > "$OUTPUT_DIR/pprof-callgraph.svg" 2>/dev/null && \
                    log "✅ pprof call graph generated: pprof-callgraph.svg" || \
                    warning "Failed to generate pprof SVG"
            fi
            
            generate_pprof_summary "$pprof_report" "$OUTPUT_DIR/pprof-summary.md"
        fi
    else
        warning "Could not collect pprof profile from $pprof_url"
        info "Ensure DeshChain is built with pprof enabled and endpoint is accessible"
        info "Add this to your main.go:"
        info "  import _ \"net/http/pprof\""
        info "  go func() { log.Println(http.ListenAndServe(\":6060\", nil)) }()"
    fi
}

# Generate pprof summary
generate_pprof_summary() {
    local pprof_report_file=$1
    local output_file=$2
    
    if [ ! -f "$pprof_report_file" ]; then
        return
    fi
    
    log "Generating pprof analysis summary..."
    
    cat > "$output_file" << EOF
# DeshChain CPU Performance Analysis (pprof)

**Profile Duration:** ${PROFILE_DURATION}s  
**Profile Type:** CPU  
**Process PID:** $DESHCHAIN_PID  

## Top CPU-Consuming Functions

\`\`\`
EOF
    
    # Extract top functions from pprof report
    head -30 "$pprof_report_file" >> "$output_file" 2>/dev/null || true
    
    cat >> "$output_file" << EOF
\`\`\`

## Analysis and Recommendations

EOF
    
    # Analyze pprof data
    python3 -c "
import re

try:
    with open('$pprof_report_file', 'r') as f:
        content = f.read()
    
    lines = content.split('\n')
    functions = []
    
    for line in lines:
        # Parse pprof output format
        if 'ms' in line and '%' in line:
            # Extract time and function name
            parts = line.strip().split()
            if len(parts) >= 3:
                try:
                    time_ms = float(parts[0].replace('ms', ''))
                    if time_ms > 10:  # Functions taking > 10ms
                        func_name = ' '.join(parts[2:])
                        functions.append((time_ms, func_name))
                except ValueError:
                    continue
    
    if functions:
        functions.sort(reverse=True)
        
        print('### Performance Critical Functions')
        print()
        
        high_impact = [f for f in functions if f[0] > 100]  # > 100ms
        if high_impact:
            print('**High Impact (>100ms):**')
            for time_ms, func in high_impact[:5]:
                print(f'- `{func}`: {time_ms}ms')
            print()
        
        medium_impact = [f for f in functions if 50 <= f[0] <= 100]
        if medium_impact:
            print('**Medium Impact (50-100ms):**')
            for time_ms, func in medium_impact[:5]:
                print(f'- `{func}`: {time_ms}ms')
            print()
        
        print('### Optimization Strategies')
        print()
        print('1. **Algorithm Optimization**: Review high-impact functions for algorithmic improvements')
        print('2. **Memory Access**: Optimize data structures and memory access patterns')
        print('3. **Parallelization**: Consider goroutines for CPU-intensive operations')
        print('4. **Caching**: Implement caching for expensive computations')
        print('5. **Profiling**: Use detailed profiling on specific hot functions')
    else:
        print('### Analysis')
        print('No significant CPU hotspots identified in this profile.')

except Exception as e:
    print(f'### Analysis Error')
    print(f'Could not parse pprof data: {e}')
" >> "$output_file" 2>/dev/null || echo "### Analysis unavailable" >> "$output_file"
    
    cat >> "$output_file" << EOF

## Next Steps

1. **Focus optimization efforts on functions with highest CPU time**
2. **Use \`go tool pprof\` for interactive analysis:**
   \`\`\`bash
   go tool pprof $pprof_output
   \`\`\`
3. **Generate call graphs for visual analysis:**
   \`\`\`bash
   go tool pprof -web $pprof_output
   \`\`\`
4. **Profile specific code paths under load**

---
*Generated from Go pprof analysis*
EOF
}

# Monitor CPU usage during profiling
monitor_cpu_usage() {
    local output_file="$OUTPUT_DIR/cpu-usage-timeline.csv"
    local duration=$1
    
    log "Monitoring CPU usage for ${duration}s..."
    
    echo "timestamp,cpu_percent,load_avg_1m,load_avg_5m,load_avg_15m,process_cpu,process_memory" > "$output_file"
    
    local samples=$((duration / 5))  # Sample every 5 seconds
    local count=0
    
    while [ $count -lt $samples ]; do
        local timestamp=$(date -Iseconds)
        local cpu_percent=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | sed 's/%us,//')
        local load_avg=$(uptime | awk -F'load average:' '{print $2}' | sed 's/ //g')
        local process_cpu=$(ps -p $DESHCHAIN_PID -o %cpu= 2>/dev/null | tr -d ' ')
        local process_memory=$(ps -p $DESHCHAIN_PID -o %mem= 2>/dev/null | tr -d ' ')
        
        # Parse load averages
        local load_1m=$(echo "$load_avg" | cut -d',' -f1)
        local load_5m=$(echo "$load_avg" | cut -d',' -f2)
        local load_15m=$(echo "$load_avg" | cut -d',' -f3)
        
        echo "$timestamp,$cpu_percent,$load_1m,$load_5m,$load_15m,$process_cpu,$process_memory" >> "$output_file"
        
        count=$((count + 1))
        sleep 5
    done
    
    log "CPU usage monitoring completed"
}

# Generate comprehensive report
generate_comprehensive_report() {
    local report_file="$OUTPUT_DIR/cpu-performance-report.md"
    
    log "Generating comprehensive CPU performance report..."
    
    cat > "$report_file" << EOF
# DeshChain CPU Performance Analysis Report

**Generated:** $(date)  
**Process PID:** $DESHCHAIN_PID  
**Profile Duration:** ${PROFILE_DURATION}s  
**Sample Frequency:** ${SAMPLE_FREQUENCY} Hz  

## Executive Summary

This report provides comprehensive CPU performance analysis for DeshChain, identifying bottlenecks and optimization opportunities.

### Key Findings

EOF
    
    # Add findings from analysis
    local critical_issues=0
    local optimization_opportunities=0
    
    if [ -f "$OUTPUT_DIR/perf-summary.md" ]; then
        critical_issues=$((critical_issues + $(grep -c "Critical Hotspots" "$OUTPUT_DIR/perf-summary.md" 2>/dev/null || echo 0)))
    fi
    
    if [ -f "$OUTPUT_DIR/pprof-summary.md" ]; then
        optimization_opportunities=$((optimization_opportunities + $(grep -c "High Impact" "$OUTPUT_DIR/pprof-summary.md" 2>/dev/null || echo 0)))
    fi
    
    cat >> "$report_file" << EOF
- **Critical Performance Issues:** $critical_issues identified
- **Optimization Opportunities:** $optimization_opportunities functions found
- **Profiling Methods Used:** $([ "$ENABLE_PERF" = "true" ] && echo "perf, " || echo "")$([ "$ENABLE_PPROF" = "true" ] && echo "pprof" || echo "")

## System Information

**Hardware:**
- CPU Cores: $(nproc)
- Architecture: $(uname -m)
- Load Average: $(uptime | awk -F'load average:' '{print $2}')

**Process Information:**
- Binary: $BINARY
- PID: $DESHCHAIN_PID
- Running Since: $(ps -o lstart= -p $DESHCHAIN_PID 2>/dev/null || echo "unknown")

## Performance Analysis

EOF
    
    # Include perf analysis if available
    if [ -f "$OUTPUT_DIR/perf-summary.md" ]; then
        echo "### Linux perf Analysis" >> "$report_file"
        echo "" >> "$report_file"
        tail -n +5 "$OUTPUT_DIR/perf-summary.md" >> "$report_file"
        echo "" >> "$report_file"
    fi
    
    # Include pprof analysis if available
    if [ -f "$OUTPUT_DIR/pprof-summary.md" ]; then
        echo "### Go pprof Analysis" >> "$report_file"
        echo "" >> "$report_file"
        tail -n +5 "$OUTPUT_DIR/pprof-summary.md" >> "$report_file"
        echo "" >> "$report_file"
    fi
    
    cat >> "$report_file" << EOF
## CPU Usage Timeline

$([ -f "$OUTPUT_DIR/cpu-usage-timeline.csv" ] && echo "See \`cpu-usage-timeline.csv\` for detailed CPU usage during profiling." || echo "CPU timeline data not available.")

## Generated Files

- 📊 **System Info:** \`system-info.txt\`
- 📈 **Process Info:** \`process-info.txt\`
$([ -f "$OUTPUT_DIR/perf-record.data" ] && echo "- 🔥 **Perf Data:** \`perf-record.data\`")
$([ -f "$OUTPUT_DIR/cpu.prof" ] && echo "- 📊 **pprof Profile:** \`cpu.prof\`")
$([ -f "$OUTPUT_DIR/flamegraph.svg" ] && echo "- 🔥 **Flame Graph:** \`flamegraph.svg\`")
$([ -f "$OUTPUT_DIR/pprof-callgraph.svg" ] && echo "- 📊 **Call Graph:** \`pprof-callgraph.svg\`")
- 📈 **CPU Timeline:** \`cpu-usage-timeline.csv\`

## Recommendations

### Immediate Actions
1. **Review high-impact functions** identified in the analysis
2. **Focus on functions using >5% CPU time**
3. **Implement performance monitoring** in production

### Development Optimizations
1. **Algorithm improvements** for CPU-intensive functions
2. **Memory access optimization** to reduce cache misses
3. **Parallelization** of suitable operations
4. **Caching strategies** for expensive computations

### Monitoring and Continuous Improvement
1. **Set up continuous profiling** in production
2. **Implement performance regression testing**
3. **Monitor CPU usage patterns** over time
4. **Regular performance reviews** and optimizations

---
*Generated by DeshChain CPU Performance Profiler*
EOF
}

# Main execution
check_dependencies
find_deshchain_process
collect_system_info

# Start CPU usage monitoring in background
monitor_cpu_usage $PROFILE_DURATION &
MONITOR_PID=$!

# Run profiling methods
if [ "$ENABLE_PERF" = "true" ]; then
    run_perf_profiling
fi

if [ "$ENABLE_PPROF" = "true" ]; then
    run_pprof_profiling
fi

# Wait for monitoring to complete
wait $MONITOR_PID 2>/dev/null || true

# Generate comprehensive report
generate_comprehensive_report

# Final summary
echo
highlight "DeshChain CPU Performance Profiling Completed!"
echo
echo -e "${BLUE}📁 Results Directory:${NC} $OUTPUT_DIR"
echo -e "${BLUE}📊 Comprehensive Report:${NC} $OUTPUT_DIR/cpu-performance-report.md"
echo -e "${BLUE}📈 CPU Timeline:${NC} $OUTPUT_DIR/cpu-usage-timeline.csv"

if [ -f "$OUTPUT_DIR/perf-record.data" ]; then
    echo -e "${BLUE}🔥 Perf Data:${NC} $OUTPUT_DIR/perf-record.data"
fi

if [ -f "$OUTPUT_DIR/cpu.prof" ]; then
    echo -e "${BLUE}📊 pprof Profile:${NC} $OUTPUT_DIR/cpu.prof"
    echo -e "${BLUE}🔧 Interactive Analysis:${NC} go tool pprof $OUTPUT_DIR/cpu.prof"
fi

if [ -f "$OUTPUT_DIR/flamegraph.svg" ]; then
    echo -e "${BLUE}🔥 Flame Graph:${NC} $OUTPUT_DIR/flamegraph.svg"
fi

echo
log "✅ CPU performance profiling completed successfully"