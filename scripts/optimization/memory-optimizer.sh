#!/bin/bash

# DeshChain Memory Optimization Tool
# Analyzes and optimizes memory usage patterns

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
OUTPUT_DIR="./memory-optimization-$(date +%Y%m%d_%H%M%S)"
DURATION=${DURATION:-300}  # 5 minutes
SAMPLE_INTERVAL=${SAMPLE_INTERVAL:-10}  # 10 seconds

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
    echo -e "${PURPLE}[MEMORY-OPT] $1${NC}"
}

# Create output directory
mkdir -p "$OUTPUT_DIR"

highlight "Starting DeshChain Memory Optimization Analysis"
log "Output directory: $OUTPUT_DIR"
log "Analysis duration: ${DURATION}s"
log "Sample interval: ${SAMPLE_INTERVAL}s"

# Check if binary exists and is running
if ! pgrep -f "$BINARY" > /dev/null; then
    warning "DeshChain binary '$BINARY' is not running"
    echo "Please start the DeshChain node and re-run this script"
    exit 1
fi

DESHCHAIN_PID=$(pgrep -f "$BINARY" | head -1)
log "Found DeshChain process: PID $DESHCHAIN_PID"

# System information
log "Collecting system information..."
cat > "$OUTPUT_DIR/system-info.txt" << EOF
DeshChain Memory Optimization Analysis
=====================================
Date: $(date)
Hostname: $(hostname)
OS: $(uname -a)
CPU: $(lscpu | grep "Model name" | cut -d: -f2 | sed 's/^ *//' || echo "unknown")
CPU Cores: $(nproc)
Total Memory: $(free -h | grep "Mem:" | awk '{print $2}')
Available Memory: $(free -h | grep "Mem:" | awk '{print $7}')
Swap: $(free -h | grep "Swap:" | awk '{print $2}')

Process Information
==================
Binary: $BINARY
PID: $DESHCHAIN_PID
Start Time: $(ps -o lstart= -p $DESHCHAIN_PID)
Command Line: $(ps -o cmd= -p $DESHCHAIN_PID)
EOF

# Memory analysis functions
collect_memory_stats() {
    local timestamp=$1
    local output_file=$2
    
    # Process memory info
    if [ -f "/proc/$DESHCHAIN_PID/status" ]; then
        local vmrss=$(grep "VmRSS:" /proc/$DESHCHAIN_PID/status | awk '{print $2}')
        local vmsize=$(grep "VmSize:" /proc/$DESHCHAIN_PID/status | awk '{print $2}')
        local vmpeak=$(grep "VmPeak:" /proc/$DESHCHAIN_PID/status | awk '{print $2}')
        local vmhwm=$(grep "VmHWM:" /proc/$DESHCHAIN_PID/status | awk '{print $2}')
        
        # System memory
        local mem_total=$(grep "MemTotal:" /proc/meminfo | awk '{print $2}')
        local mem_free=$(grep "MemFree:" /proc/meminfo | awk '{print $2}')
        local mem_available=$(grep "MemAvailable:" /proc/meminfo | awk '{print $2}')
        local mem_cached=$(grep "Cached:" /proc/meminfo | awk '{print $2}')
        local mem_buffers=$(grep "Buffers:" /proc/meminfo | awk '{print $2}')
        
        # Calculate percentages
        local rss_percent=$(( vmrss * 100 / mem_total ))
        local available_percent=$(( mem_available * 100 / mem_total ))
        
        echo "$timestamp,$vmrss,$vmsize,$vmpeak,$vmhwm,$mem_total,$mem_free,$mem_available,$mem_cached,$mem_buffers,$rss_percent,$available_percent" >> "$output_file"
    fi
}

analyze_memory_growth() {
    local data_file=$1
    local output_file=$2
    
    log "Analyzing memory growth patterns..."
    
    python3 -c "
import csv
import json
from datetime import datetime

data = []
with open('$data_file', 'r') as f:
    reader = csv.reader(f)
    for row in reader:
        if len(row) >= 12:
            data.append({
                'timestamp': row[0],
                'rss_kb': int(row[1]),
                'size_kb': int(row[2]),
                'peak_kb': int(row[3]),
                'hwm_kb': int(row[4]),
                'rss_percent': float(row[10])
            })

if len(data) < 2:
    print('Insufficient data for analysis')
    exit(1)

# Calculate growth rates
initial_rss = data[0]['rss_kb']
final_rss = data[-1]['rss_kb']
max_rss = max(d['rss_kb'] for d in data)
min_rss = min(d['rss_kb'] for d in data)

growth_rate = (final_rss - initial_rss) / initial_rss * 100 if initial_rss > 0 else 0
volatility = (max_rss - min_rss) / min_rss * 100 if min_rss > 0 else 0

# Detect memory leaks (sustained growth)
leak_threshold = 5  # 5% growth is concerning
is_leaking = growth_rate > leak_threshold

# Calculate average memory usage
avg_rss = sum(d['rss_kb'] for d in data) / len(data)
avg_percent = sum(d['rss_percent'] for d in data) / len(data)

analysis = {
    'initial_memory_kb': initial_rss,
    'final_memory_kb': final_rss,
    'max_memory_kb': max_rss,
    'min_memory_kb': min_rss,
    'average_memory_kb': int(avg_rss),
    'average_percent': round(avg_percent, 2),
    'growth_rate_percent': round(growth_rate, 2),
    'volatility_percent': round(volatility, 2),
    'potential_leak': is_leaking,
    'samples': len(data),
    'duration_seconds': len(data) * $SAMPLE_INTERVAL
}

with open('$output_file', 'w') as f:
    json.dump(analysis, f, indent=2)

print(f'Memory analysis completed: {len(data)} samples')
print(f'Growth rate: {growth_rate:.2f}%')
print(f'Potential leak: {is_leaking}')
" 2>/dev/null || warning "Python analysis failed - ensure Python 3 is installed"
}

generate_memory_recommendations() {
    local analysis_file=$1
    local output_file=$2
    
    log "Generating memory optimization recommendations..."
    
    if [ ! -f "$analysis_file" ]; then
        warning "Analysis file not found, generating generic recommendations"
        cat > "$output_file" << EOF
# Memory Optimization Recommendations

## Generic Recommendations

1. **Enable Memory Profiling**
   - Use pprof for detailed memory analysis
   - Implement periodic memory snapshots
   - Monitor heap growth patterns

2. **Garbage Collection Tuning**
   - Adjust GOGC environment variable
   - Consider using GOMEMLIMIT for Go 1.19+
   - Monitor GC frequency and duration

3. **Memory Pool Implementation**
   - Implement object pools for frequently allocated objects
   - Use sync.Pool for short-lived objects
   - Consider buffer pools for I/O operations

4. **Data Structure Optimization**
   - Use memory-efficient data structures
   - Implement proper cleanup for maps and slices
   - Consider using byte pools for network operations

5. **Caching Strategy**
   - Implement LRU caches with size limits
   - Use TTL for cached entries
   - Monitor cache hit rates and memory usage
EOF
        return
    fi
    
    python3 -c "
import json

with open('$analysis_file', 'r') as f:
    analysis = json.load(f)

recommendations = []

# Memory leak detection
if analysis.get('potential_leak', False):
    recommendations.extend([
        '🚨 **MEMORY LEAK DETECTED** - Growth rate: {:.2f}%'.format(analysis['growth_rate_percent']),
        'Implement immediate memory profiling to identify leak sources',
        'Review goroutine lifecycle management',
        'Check for unclosed resources (files, connections, channels)',
        'Implement proper cleanup in defer statements'
    ])

# High memory usage
if analysis.get('average_percent', 0) > 80:
    recommendations.extend([
        '⚠️ **HIGH MEMORY USAGE** - Average: {:.1f}%'.format(analysis['average_percent']),
        'Consider implementing memory limits',
        'Optimize data structures and algorithms',
        'Implement data compression where applicable',
        'Review caching strategies and implement TTL'
    ])

# High volatility
if analysis.get('volatility_percent', 0) > 50:
    recommendations.extend([
        '📊 **HIGH MEMORY VOLATILITY** - Variance: {:.1f}%'.format(analysis['volatility_percent']),
        'Implement object pooling to reduce allocations',
        'Review batch processing logic',
        'Consider using buffered operations',
        'Implement backpressure mechanisms'
    ])

# Large memory footprint
if analysis.get('average_memory_kb', 0) > 1024 * 1024:  # > 1GB
    recommendations.extend([
        '💾 **LARGE MEMORY FOOTPRINT** - Average: {:.0f} MB'.format(analysis['average_memory_kb'] / 1024),
        'Implement memory-mapped files for large datasets',
        'Consider database offloading for large state',
        'Implement data pagination and lazy loading',
        'Review data retention policies'
    ])

# General recommendations
general_recommendations = [
    '## General Optimization Strategies',
    '',
    '### Immediate Actions',
    '1. Enable Go memory profiling: \`go tool pprof heap\`',
    '2. Implement memory monitoring dashboards',
    '3. Set up memory usage alerts',
    '4. Review and optimize hot paths',
    '',
    '### Code Optimizations', 
    '1. Use \`sync.Pool\` for frequently allocated objects',
    '2. Implement proper cleanup with \`defer\` statements',
    '3. Use byte slices instead of strings for mutable data',
    '4. Implement connection pooling for external resources',
    '',
    '### Configuration Tuning',
    '1. Tune GOGC environment variable (default: 100)',
    '2. Set GOMEMLIMIT for Go 1.19+ (soft memory limit)',
    '3. Configure appropriate heap size limits',
    '4. Implement graceful degradation under memory pressure',
    '',
    '### Monitoring and Alerting',
    '1. Set up memory usage thresholds (warning: 70%, critical: 85%)',
    '2. Monitor GC frequency and pause times',
    '3. Track memory growth trends over time',
    '4. Implement automatic memory dumps on high usage'
]

# Write recommendations
with open('$output_file', 'w') as f:
    f.write('# DeshChain Memory Optimization Recommendations\\n\\n')
    f.write('**Analysis Summary:**\\n')
    f.write(f'- Average Memory Usage: {analysis.get(\"average_memory_kb\", 0) / 1024:.1f} MB ({analysis.get(\"average_percent\", 0):.1f}%)\\n')
    f.write(f'- Memory Growth Rate: {analysis.get(\"growth_rate_percent\", 0):.2f}%\\n')
    f.write(f'- Memory Volatility: {analysis.get(\"volatility_percent\", 0):.1f}%\\n')
    f.write(f'- Potential Leak: {\"Yes\" if analysis.get(\"potential_leak\", False) else \"No\"}\\n\\n')
    
    if recommendations:
        f.write('## Critical Issues\\n\\n')
        for i, rec in enumerate(recommendations, 1):
            f.write(f'{i}. {rec}\\n')
        f.write('\\n')
    
    for line in general_recommendations:
        f.write(line + '\\n')

print('Memory recommendations generated successfully')
" 2>/dev/null || warning "Failed to generate detailed recommendations"
}

create_memory_dashboard() {
    local data_file=$1
    local output_file=$2
    
    log "Creating memory usage dashboard..."
    
    python3 -c "
import csv
import json
from datetime import datetime, timedelta

# Read memory data
data = []
start_time = None

try:
    with open('$data_file', 'r') as f:
        reader = csv.reader(f)
        for row in reader:
            if len(row) >= 12:
                timestamp = row[0]
                if start_time is None:
                    start_time = datetime.fromisoformat(timestamp.replace('Z', '+00:00'))
                
                current_time = datetime.fromisoformat(timestamp.replace('Z', '+00:00'))
                elapsed = (current_time - start_time).total_seconds()
                
                data.append({
                    'elapsed_seconds': elapsed,
                    'rss_mb': int(row[1]) / 1024,
                    'size_mb': int(row[2]) / 1024,
                    'rss_percent': float(row[10]),
                    'available_percent': float(row[11])
                })
except Exception as e:
    print(f'Error reading data: {e}')
    exit(1)

# Generate simple ASCII chart
chart_width = 60
chart_height = 20

if data:
    max_rss = max(d['rss_mb'] for d in data)
    min_rss = min(d['rss_mb'] for d in data)
    
    print('Memory Usage Over Time (RSS in MB)')
    print('=' * 70)
    print(f'Max: {max_rss:.1f} MB, Min: {min_rss:.1f} MB')
    print()
    
    # Create ASCII chart
    for i in range(chart_height, 0, -1):
        threshold = min_rss + (max_rss - min_rss) * i / chart_height
        line = f'{threshold:6.1f} |'
        
        for d in data[::max(1, len(data)//chart_width)]:
            if d['rss_mb'] >= threshold:
                line += '*'
            else:
                line += ' '
        print(line)
    
    print(' ' * 8 + '-' * chart_width)
    print(' ' * 8 + f'0s{\" \" * (chart_width-10)}{data[-1][\"elapsed_seconds\"]:.0f}s')

# Save dashboard data
dashboard_data = {
    'title': 'DeshChain Memory Usage Dashboard',
    'timestamp': datetime.now().isoformat(),
    'duration_seconds': data[-1]['elapsed_seconds'] if data else 0,
    'samples': len(data),
    'memory_stats': {
        'max_rss_mb': max(d['rss_mb'] for d in data) if data else 0,
        'min_rss_mb': min(d['rss_mb'] for d in data) if data else 0,
        'avg_rss_mb': sum(d['rss_mb'] for d in data) / len(data) if data else 0,
        'max_percent': max(d['rss_percent'] for d in data) if data else 0,
        'avg_percent': sum(d['rss_percent'] for d in data) / len(data) if data else 0
    },
    'timeline_data': data
}

with open('$output_file', 'w') as f:
    json.dump(dashboard_data, f, indent=2)

print(f'\\nDashboard data saved to: $output_file')
" 2>/dev/null || warning "Failed to create memory dashboard"
}

# Memory optimization checks
run_memory_checks() {
    local output_file="$OUTPUT_DIR/memory-checks.txt"
    
    log "Running memory optimization checks..."
    
    cat > "$output_file" << EOF
DeshChain Memory Optimization Checks
===================================

1. Go Runtime Configuration
EOF
    
    # Check Go environment variables
    if [ -n "$GOGC" ]; then
        echo "   ✅ GOGC is set: $GOGC" >> "$output_file"
    else
        echo "   ⚠️ GOGC not set (using default: 100)" >> "$output_file"
    fi
    
    if [ -n "$GOMEMLIMIT" ]; then
        echo "   ✅ GOMEMLIMIT is set: $GOMEMLIMIT" >> "$output_file"
    else
        echo "   ⚠️ GOMEMLIMIT not set (consider setting for Go 1.19+)" >> "$output_file"
    fi
    
    cat >> "$output_file" << EOF

2. System Memory Configuration
EOF
    
    # Check swap configuration
    local swap_total=$(free | grep "Swap:" | awk '{print $2}')
    if [ "$swap_total" -gt 0 ]; then
        echo "   ✅ Swap is configured: $(free -h | grep "Swap:" | awk '{print $2}')" >> "$output_file"
    else
        echo "   ⚠️ No swap configured (consider adding swap for memory pressure relief)" >> "$output_file"
    fi
    
    # Check overcommit settings
    local overcommit=$(cat /proc/sys/vm/overcommit_memory 2>/dev/null || echo "unknown")
    echo "   📊 Memory overcommit setting: $overcommit" >> "$output_file"
    
    cat >> "$output_file" << EOF

3. Process Memory Limits
EOF
    
    # Check ulimits
    local mem_limit=$(ulimit -m 2>/dev/null || echo "unlimited")
    local virtual_limit=$(ulimit -v 2>/dev/null || echo "unlimited")
    
    echo "   📊 Memory limit (RSS): $mem_limit" >> "$output_file"
    echo "   📊 Virtual memory limit: $virtual_limit" >> "$output_file"
    
    # Check if running in container
    if [ -f "/.dockerenv" ] || grep -q "docker\|lxc" /proc/1/cgroup 2>/dev/null; then
        echo "   🐳 Running in container (check container memory limits)" >> "$output_file"
    else
        echo "   🖥️ Running on bare metal/VM" >> "$output_file"
    fi
    
    cat >> "$output_file" << EOF

4. Memory Optimization Recommendations
   - Consider setting GOGC to 50-200 based on workload
   - Set GOMEMLIMIT to 80-90% of available memory for Go 1.19+
   - Monitor memory growth patterns regularly
   - Implement memory profiling in production
   - Use object pools for frequently allocated objects
   - Configure appropriate swap space (1-2x RAM)
EOF
}

# Main execution
log "Starting memory data collection..."

# Initialize CSV file
MEMORY_DATA_FILE="$OUTPUT_DIR/memory-usage.csv"
echo "timestamp,rss_kb,size_kb,peak_kb,hwm_kb,total_kb,free_kb,available_kb,cached_kb,buffers_kb,rss_percent,available_percent" > "$MEMORY_DATA_FILE"

# Collect memory samples
SAMPLES=0
MAX_SAMPLES=$((DURATION / SAMPLE_INTERVAL))

log "Collecting $MAX_SAMPLES samples over ${DURATION}s..."

while [ $SAMPLES -lt $MAX_SAMPLES ]; do
    TIMESTAMP=$(date -Iseconds)
    collect_memory_stats "$TIMESTAMP" "$MEMORY_DATA_FILE"
    
    SAMPLES=$((SAMPLES + 1))
    PROGRESS=$((SAMPLES * 100 / MAX_SAMPLES))
    
    if [ $((SAMPLES % 6)) -eq 0 ]; then  # Update every minute
        info "Progress: $PROGRESS% ($SAMPLES/$MAX_SAMPLES samples)"
    fi
    
    sleep $SAMPLE_INTERVAL
done

log "Memory data collection completed"

# Analyze collected data
ANALYSIS_FILE="$OUTPUT_DIR/memory-analysis.json"
analyze_memory_growth "$MEMORY_DATA_FILE" "$ANALYSIS_FILE"

# Generate recommendations
RECOMMENDATIONS_FILE="$OUTPUT_DIR/memory-recommendations.md"
generate_memory_recommendations "$ANALYSIS_FILE" "$RECOMMENDATIONS_FILE"

# Create dashboard
DASHBOARD_FILE="$OUTPUT_DIR/memory-dashboard.json"
create_memory_dashboard "$MEMORY_DATA_FILE" "$DASHBOARD_FILE"

# Run memory checks
run_memory_checks

# Generate summary report
SUMMARY_FILE="$OUTPUT_DIR/memory-optimization-summary.md"
log "Generating summary report..."

cat > "$SUMMARY_FILE" << EOF
# DeshChain Memory Optimization Summary

**Analysis Date:** $(date)  
**Process PID:** $DESHCHAIN_PID  
**Duration:** ${DURATION}s  
**Samples:** $SAMPLES  

## Quick Summary

EOF

if [ -f "$ANALYSIS_FILE" ]; then
    python3 -c "
import json
with open('$ANALYSIS_FILE', 'r') as f:
    analysis = json.load(f)

print(f'- **Average Memory Usage:** {analysis.get(\"average_memory_kb\", 0) / 1024:.1f} MB')
print(f'- **Memory Growth Rate:** {analysis.get(\"growth_rate_percent\", 0):.2f}%')
print(f'- **Memory Volatility:** {analysis.get(\"volatility_percent\", 0):.1f}%')
print(f'- **Potential Memory Leak:** {\"⚠️ Yes\" if analysis.get(\"potential_leak\", False) else \"✅ No\"}')
" >> "$SUMMARY_FILE" 2>/dev/null || echo "- Analysis data unavailable" >> "$SUMMARY_FILE"
fi

cat >> "$SUMMARY_FILE" << EOF

## Generated Files

- 📊 **Raw Data:** \`memory-usage.csv\`
- 📈 **Analysis:** \`memory-analysis.json\`
- 📋 **Recommendations:** \`memory-recommendations.md\`
- 🎛️ **Dashboard:** \`memory-dashboard.json\`
- ✅ **System Checks:** \`memory-checks.txt\`

## Next Steps

1. Review the recommendations in \`memory-recommendations.md\`
2. Implement suggested optimizations
3. Set up continuous memory monitoring
4. Schedule regular memory analysis

---
*Generated by DeshChain Memory Optimizer*
EOF

# Final summary
echo
highlight "DeshChain Memory Optimization Analysis Completed!"
echo
echo -e "${BLUE}📁 Results Directory:${NC} $OUTPUT_DIR"
echo -e "${BLUE}📊 Summary Report:${NC} $OUTPUT_DIR/memory-optimization-summary.md"
echo -e "${BLUE}📈 Recommendations:${NC} $OUTPUT_DIR/memory-recommendations.md"
echo -e "${BLUE}🎛️ Dashboard Data:${NC} $OUTPUT_DIR/memory-dashboard.json"
echo

# Check for critical issues
if [ -f "$ANALYSIS_FILE" ]; then
    POTENTIAL_LEAK=$(python3 -c "import json; f=open('$ANALYSIS_FILE'); print(json.load(f).get('potential_leak', False))" 2>/dev/null || echo "False")
    if [ "$POTENTIAL_LEAK" = "True" ]; then
        warning "🚨 POTENTIAL MEMORY LEAK DETECTED!"
        echo "   Review the recommendations immediately"
        exit 1
    fi
fi

log "✅ Memory optimization analysis completed successfully"