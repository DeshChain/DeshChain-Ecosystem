# DeshChain Performance Optimization Tools

Comprehensive suite of performance optimization tools for DeshChain blockchain, providing deep insights into system performance and actionable optimization recommendations.

## Overview

This toolkit provides four specialized optimization tools:
- **Performance Profiler**: Deep application profiling with Go pprof integration
- **Memory Optimizer**: Memory usage analysis and leak detection
- **CPU Profiler**: CPU performance analysis with perf and pprof
- **Blockchain Optimizer**: Blockchain-specific performance optimization

## Tools

### 1. Performance Profiler (`performance-profiler.go`)
Comprehensive Go application profiling with automatic bottleneck detection.

```bash
# Basic profiling
go run scripts/optimization/performance-profiler.go \
  -duration=5m \
  -output=./profiles

# Advanced profiling with HTTP endpoint
go run scripts/optimization/performance-profiler.go \
  -node="http://localhost:26657" \
  -duration=10m \
  -output=./detailed-profiles \
  -pprof=true \
  -pprof-port=6060
```

**Features:**
- CPU, memory, block, and mutex profiling
- Automatic bottleneck detection
- Performance recommendations
- HTTP pprof server integration
- Comprehensive analysis reports

### 2. Memory Optimizer (`memory-optimizer.sh`)
Advanced memory usage analysis with leak detection and optimization recommendations.

```bash
# Basic memory analysis
./scripts/optimization/memory-optimizer.sh

# Extended analysis
DURATION=600 SAMPLE_INTERVAL=5 ./scripts/optimization/memory-optimizer.sh
```

**Features:**
- Real-time memory usage tracking
- Memory leak detection
- Growth pattern analysis
- System memory optimization checks
- Detailed recommendations

### 3. CPU Profiler (`cpu-profiler.sh`)
Comprehensive CPU performance analysis using Linux perf and Go pprof.

```bash
# Basic CPU profiling
./scripts/optimization/cpu-profiler.sh

# Extended profiling with custom duration
PROFILE_DURATION=300 SAMPLE_FREQUENCY=97 ./scripts/optimization/cpu-profiler.sh
```

**Features:**
- Linux perf integration
- Go pprof CPU profiling
- Flame graph generation
- Call graph analysis
- Hotspot identification

### 4. Blockchain Optimizer (`blockchain-optimizer.py`)
Blockchain-specific performance optimization with consensus and throughput analysis.

```bash
# Basic blockchain optimization
python3 scripts/optimization/blockchain-optimizer.py \
  --node="http://localhost:26657" \
  --duration=600

# Comprehensive analysis
python3 scripts/optimization/blockchain-optimizer.py \
  --node="http://localhost:26657" \
  --duration=1800 \
  --interval=15 \
  --output=./blockchain-analysis
```

**Features:**
- Block time analysis
- Throughput optimization
- Network health assessment
- Consensus performance
- Implementation roadmap

## Quick Start

### Prerequisites

**System Requirements:**
- Linux system (Ubuntu 20.04+ recommended)
- Go 1.21+
- Python 3.8+
- DeshChain node running

**Optional Tools:**
```bash
# For enhanced profiling
sudo apt-get install linux-tools-generic
pip3 install aiohttp psutil

# For flame graphs
git clone https://github.com/brendangregg/FlameGraph.git
export PATH=$PATH:$(pwd)/FlameGraph
```

### Basic Performance Analysis

```bash
# 1. Start with memory analysis
./scripts/optimization/memory-optimizer.sh

# 2. Run CPU profiling
./scripts/optimization/cpu-profiler.sh

# 3. Analyze blockchain performance
python3 scripts/optimization/blockchain-optimizer.py

# 4. Deep application profiling
go run scripts/optimization/performance-profiler.go -duration=5m
```

### Comprehensive Analysis

```bash
# Run all optimization tools in sequence
./scripts/optimization/run-all-optimizations.sh
```

## Output Files

### Performance Profiler
- `performance-analysis-report.md` - Comprehensive analysis report
- `performance-analysis-results.json` - Structured results data
- `cpu.prof`, `mem.prof`, `block.prof`, `mutex.prof` - Go pprof profiles

### Memory Optimizer
- `memory-optimization-summary.md` - Executive summary
- `memory-recommendations.md` - Detailed recommendations
- `memory-usage.csv` - Raw memory usage data
- `memory-analysis.json` - Analysis results

### CPU Profiler
- `cpu-performance-report.md` - Comprehensive CPU analysis
- `perf-record.data` - Linux perf data
- `cpu.prof` - Go pprof CPU profile
- `flamegraph.svg` - CPU flame graph
- `cpu-usage-timeline.csv` - CPU usage over time

### Blockchain Optimizer
- `blockchain-optimization-report.md` - Full optimization report
- `blockchain-optimization-plan.json` - Structured optimization plan
- `blockchain-metrics.json` - Historical metrics data

## Performance Optimization Workflow

### 1. Initial Assessment
```bash
# Quick performance check
go run scripts/optimization/performance-profiler.go -duration=2m

# Check for memory issues
DURATION=300 ./scripts/optimization/memory-optimizer.sh

# Assess blockchain performance
python3 scripts/optimization/blockchain-optimizer.py --duration=300
```

### 2. Deep Analysis
```bash
# Extended profiling for production issues
go run scripts/optimization/performance-profiler.go \
  -duration=30m \
  -node="http://localhost:26657" \
  -pprof=true

# Detailed CPU analysis with flame graphs
PROFILE_DURATION=600 ./scripts/optimization/cpu-profiler.sh

# Comprehensive blockchain analysis
python3 scripts/optimization/blockchain-optimizer.py \
  --duration=3600 \
  --interval=10
```

### 3. Implementation
1. **Review Reports**: Analyze generated markdown reports
2. **Prioritize Issues**: Focus on high-priority recommendations
3. **Implement Fixes**: Apply suggested optimizations
4. **Validate Improvements**: Re-run tools to measure impact

## Integration

### Makefile Integration
```bash
# Add to production monitoring
make optimize-performance    # Run all optimization tools
make profile-memory         # Memory optimization only
make profile-cpu           # CPU profiling only
make optimize-blockchain   # Blockchain optimization only
```

### CI/CD Integration
```yaml
- name: Performance Optimization Check
  run: |
    # Quick performance validation
    timeout 300 python3 scripts/optimization/blockchain-optimizer.py \
      --duration=120 \
      --output=ci-optimization
    
    # Check for critical issues
    CRITICAL_ISSUES=$(jq -r '.recommendations.high_priority | length' \
      ci-optimization/blockchain-optimization-plan.json)
    
    if [ "$CRITICAL_ISSUES" -gt 0 ]; then
      echo "Critical performance issues detected: $CRITICAL_ISSUES"
      exit 1
    fi
```

### Production Monitoring
```bash
# Automated daily optimization checks
0 2 * * * /path/to/scripts/optimization/blockchain-optimizer.py \
  --duration=1800 \
  --output=/var/log/deshchain/daily-optimization

# Weekly comprehensive analysis
0 3 * * 0 /path/to/scripts/optimization/run-all-optimizations.sh
```

## Configuration

### Environment Variables
```bash
# Global settings
export DESHCHAIN_OPTIMIZATION_DIR="/var/log/deshchain/optimization"
export DESHCHAIN_PROFILE_DURATION="600"
export DESHCHAIN_SAMPLE_INTERVAL="10"

# Memory optimizer
export MEMORY_SAMPLE_INTERVAL="5"
export MEMORY_ANALYSIS_DURATION="300"

# CPU profiler
export CPU_SAMPLE_FREQUENCY="99"
export ENABLE_FLAME_GRAPHS="true"

# Blockchain optimizer
export BLOCKCHAIN_ANALYSIS_INTERVAL="30"
export BLOCKCHAIN_NODE_URL="http://localhost:26657"
```

### Tool-Specific Configuration

#### Performance Profiler
```bash
go run scripts/optimization/performance-profiler.go \
  -node="http://localhost:26657" \
  -duration=10m \
  -output=./profiles \
  -cpuprofile=cpu.prof \
  -memprofile=mem.prof \
  -blockprofile=block.prof \
  -mutexprofile=mutex.prof \
  -pprof=true \
  -pprof-port=6060
```

#### Memory Optimizer
```bash
BINARY="deshchaind" \
DURATION=600 \
SAMPLE_INTERVAL=5 \
./scripts/optimization/memory-optimizer.sh
```

#### CPU Profiler
```bash
BINARY="deshchaind" \
PROFILE_DURATION=300 \
SAMPLE_FREQUENCY=97 \
ENABLE_PERF=true \
ENABLE_PPROF=true \
./scripts/optimization/cpu-profiler.sh
```

#### Blockchain Optimizer
```bash
python3 scripts/optimization/blockchain-optimizer.py \
  --node="http://localhost:26657" \
  --duration=1800 \
  --interval=15 \
  --output=./blockchain-analysis
```

## Troubleshooting

### Common Issues

**"Permission denied" for perf**
```bash
# Enable perf for non-root users
sudo sysctl -w kernel.perf_event_paranoid=1
# Or run with sudo
sudo ./scripts/optimization/cpu-profiler.sh
```

**"Process not found"**
```bash
# Ensure DeshChain is running
pgrep -f deshchaind
# Start DeshChain if needed
./bin/deshchaind start
```

**"Python dependencies missing"**
```bash
# Install required packages
pip3 install aiohttp psutil
```

**"Go tool pprof not found"**
```bash
# Ensure Go is installed and in PATH
go version
which go
```

### Performance Tips

1. **Run During Load**: Profile during actual usage for realistic results
2. **Multiple Samples**: Run multiple times to identify consistent patterns
3. **Baseline Comparison**: Establish baseline metrics before optimization
4. **Incremental Changes**: Apply one optimization at a time to measure impact
5. **Continuous Monitoring**: Set up regular optimization checks

## Advanced Usage

### Custom Profiling Scenarios

**High-Load Analysis**
```bash
# Start load testing
make test-stress &
LOAD_PID=$!

# Profile during load
go run scripts/optimization/performance-profiler.go -duration=10m

# Stop load testing
kill $LOAD_PID
```

**Memory Leak Investigation**
```bash
# Extended memory monitoring
DURATION=3600 SAMPLE_INTERVAL=1 ./scripts/optimization/memory-optimizer.sh

# Look for sustained growth patterns
python3 -c "
import json
with open('memory-optimization-*/memory-analysis.json') as f:
    data = json.load(f)
    if data['potential_leak']:
        print('LEAK DETECTED: Growth rate:', data['growth_rate_percent'], '%')
"
```

**Production Optimization**
```bash
# Non-intrusive production profiling
go run scripts/optimization/performance-profiler.go \
  -duration=5m \
  -output=./prod-profiles \
  -pprof=false  # Disable HTTP server

# Low-frequency CPU sampling
SAMPLE_FREQUENCY=19 ./scripts/optimization/cpu-profiler.sh
```

## Best Practices

### Optimization Workflow
1. **Baseline Measurement**: Always establish baseline performance
2. **Systematic Analysis**: Use all tools for comprehensive view
3. **Prioritized Implementation**: Address high-priority issues first
4. **Validation Testing**: Verify improvements with benchmarks
5. **Continuous Monitoring**: Regular optimization checks

### Production Considerations
1. **Non-Intrusive Profiling**: Use minimal overhead settings
2. **Off-Peak Analysis**: Run intensive profiling during low-traffic periods
3. **Automated Monitoring**: Set up continuous optimization tracking
4. **Alert Integration**: Connect optimization tools to monitoring systems
5. **Documentation**: Maintain optimization logs and decisions

---

For detailed analysis of specific performance issues, refer to the individual tool documentation and generated reports.