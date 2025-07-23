# DeshChain Load Testing Framework

Comprehensive load testing and performance validation framework for DeshChain blockchain.

## Overview

This framework provides production-ready tools for:
- **Load Testing**: Simulate real-world transaction loads
- **Stress Testing**: Identify system breaking points  
- **Performance Monitoring**: Real-time metrics and alerting
- **Benchmark Suite**: Comprehensive performance validation

## Tools

### 1. Load Test (`load-test.go`)
Core load testing tool that simulates concurrent transaction workloads.

```bash
# Basic load test
go run scripts/load-testing/load-test.go \
  -workers=10 \
  -tx-per-worker=100 \
  -chain-id="deshchain-testnet-1"

# Advanced configuration
go run scripts/load-testing/load-test.go \
  -workers=50 \
  -tx-per-worker=200 \
  -duration=300s \
  -node="tcp://localhost:26657" \
  -output="results.json"
```

**Parameters:**
- `--workers`: Number of concurrent workers (default: 10)
- `--tx-per-worker`: Transactions per worker (default: 100)
- `--duration`: Test duration (0 = use tx-per-worker)
- `--node`: Node RPC address (default: tcp://localhost:26657)
- `--chain-id`: Chain identifier (default: deshchain-testnet-1)
- `--output`: JSON results file path
- `--keyring-backend`: Keyring backend (default: test)

### 2. Stress Test (`stress-test.sh`)
Automated stress testing with multiple load scenarios.

```bash
# Run all stress test scenarios
./scripts/load-testing/stress-test.sh

# Custom configuration
CHAIN_ID="deshchain-mainnet-1" \
NODE_ADDRESS="tcp://validator.deshchain.network:26657" \
./scripts/load-testing/stress-test.sh
```

**Scenarios:**
- **Baseline**: 10 workers, 100 tx each
- **Moderate**: 25 workers, 100 tx each  
- **High**: 50 workers, 100 tx each
- **Extreme**: 100 workers, 100 tx each
- **Burst**: 10 workers, 1000 tx each
- **Massive**: 200 workers, 50 tx each

### 3. Performance Monitor (`performance-monitor.py`)
Real-time monitoring with alerting capabilities.

```bash
# Basic monitoring
python3 scripts/load-testing/performance-monitor.py \
  --node="http://localhost:26657" \
  --duration=3600 \
  --interval=30

# Advanced monitoring with alerts
python3 scripts/load-testing/performance-monitor.py \
  --node="http://localhost:26657" \
  --duration=0 \
  --interval=15 \
  --webhook="https://hooks.slack.com/services/YOUR/WEBHOOK/URL" \
  --max-block-time=15 \
  --min-tps=20 \
  --max-memory=85
```

**Parameters:**
- `--node`: Node RPC URL (default: http://localhost:26657)
- `--duration`: Monitoring duration in seconds (0 = infinite)
- `--interval`: Check interval in seconds (default: 30)
- `--output`: Metrics output file
- `--webhook`: Alert webhook URL
- `--max-block-time`: Max block time threshold (default: 10.0s)
- `--min-tps`: Min TPS threshold (default: 10.0)
- `--max-memory`: Max memory usage % (default: 80.0)
- `--max-cpu`: Max CPU usage % (default: 80.0)

### 4. Benchmark Suite (`benchmark-suite.sh`)
Comprehensive performance validation for production readiness.

```bash
# Full benchmark suite
./scripts/load-testing/benchmark-suite.sh

# Custom duration
BENCHMARK_DURATION=600 \
WARMUP_DURATION=120 \
./scripts/load-testing/benchmark-suite.sh
```

**Features:**
- System information collection
- Warmup phase
- Multiple benchmark scenarios
- Real-time performance monitoring
- Production readiness assessment
- Comprehensive reporting

## Prerequisites

### System Requirements
- **Go 1.21+**: For load testing tools
- **Python 3.8+**: For monitoring scripts
- **DeshChain Binary**: `deshchaind` in PATH
- **System Tools**: `curl`, `jq`, `bc`

### Python Dependencies
```bash
pip3 install aiohttp psutil
```

### Go Dependencies
Load testing tool will automatically download required dependencies.

## Quick Start

### 1. Basic Load Test
```bash
# Start with a simple load test
go run scripts/load-testing/load-test.go -workers=5 -tx-per-worker=50
```

### 2. Stress Testing
```bash
# Run comprehensive stress tests
./scripts/load-testing/stress-test.sh
```

### 3. Performance Monitoring
```bash
# Monitor for 1 hour
python3 scripts/load-testing/performance-monitor.py --duration=3600
```

### 4. Full Benchmark
```bash
# Complete production readiness assessment
./scripts/load-testing/benchmark-suite.sh
```

## Output Files

### Load Test Results (`results.json`)
```json
{
  "config": {
    "workers": 10,
    "tx_per_worker": 100,
    "test_duration": "2m30s",
    "tx_type": "bank-send"
  },
  "results": {
    "total_tx": 1000,
    "successful_tx": 950,
    "failed_tx": 50,
    "success_rate": 95.00,
    "total_duration": "2m25s",
    "tx_per_second": 6.52,
    "average_latency": "1.2s",
    "total_gas_used": 180000000,
    "total_fees": "950000namo",
    "error_count": 50
  }
}
```

### Stress Test Report (`stress-test-report.md`)
- Executive summary
- Scenario results table
- Performance analysis
- Breaking point identification
- Recommendations

### Performance Metrics (`performance-metrics.jsonl`)
Line-delimited JSON with timestamped metrics:
```json
{"timestamp":"2024-01-15T10:30:00Z","block_height":12345,"block_time":3.2,"tps":15.5,"memory_usage":45.2,"cpu_usage":32.1,"peer_count":8}
```

### Benchmark Report (`benchmark-report.md`)
- Comprehensive performance analysis
- Production readiness score
- Scalability assessment
- Infrastructure recommendations

## Performance Thresholds

### Production Targets
- **TPS**: 50+ sustained, 100+ peak
- **Block Time**: <6 seconds average
- **Success Rate**: >99%
- **Memory Usage**: <80%
- **CPU Usage**: <70%

### Alert Thresholds
- **Critical**: Block time >15s, Memory >90%, Disk >95%
- **Warning**: TPS <10, Memory >80%, CPU >80%
- **Info**: Peer count <5, Sync lag >10 blocks

## CI/CD Integration

### GitHub Actions
```yaml
- name: Load Testing
  run: |
    go run scripts/load-testing/load-test.go \
      -workers=20 \
      -tx-per-worker=100 \
      -output=ci-load-test.json
    
    # Check success rate
    SUCCESS_RATE=$(jq -r '.results.success_rate' ci-load-test.json)
    if (( $(echo "$SUCCESS_RATE < 95" | bc -l) )); then
      echo "Load test failed: Success rate $SUCCESS_RATE% < 95%"
      exit 1
    fi
```

### Production Deployment
```bash
# Pre-deployment validation
./scripts/load-testing/benchmark-suite.sh

# Check exit code
if [ $? -eq 0 ]; then
  echo "Production readiness validated"
  # Proceed with deployment
else
  echo "System not ready for production"
  exit 1
fi
```

## Troubleshooting

### Common Issues

**"Connection refused"**
- Ensure DeshChain node is running
- Check node address and port
- Verify firewall settings

**"Transaction failed"**
- Check account balances
- Verify chain ID
- Review node logs

**"High memory usage"**
- Reduce concurrent workers
- Check for memory leaks
- Monitor system resources

**"Low TPS"**
- Optimize transaction size
- Check network latency
- Review node configuration

### Debug Mode
```bash
# Enable verbose logging
export LOG_LEVEL=debug

# Run with detailed output
go run scripts/load-testing/load-test.go -workers=1 -tx-per-worker=10
```

## Advanced Usage

### Custom Transaction Types
Modify `load-test.go` to support different transaction types:
```go
// Add new transaction type
case "nft-mint":
    msg = &nfttypes.MsgMintNFT{...}
case "token-transfer":
    msg = &tokentypes.MsgTransfer{...}
```

### Custom Metrics
Extend `performance-monitor.py` with custom metrics:
```python
async def collect_custom_metrics(self):
    # Add custom metric collection
    custom_data = await self.fetch_custom_endpoint()
    return custom_data
```

### Load Profiles
Create custom load profiles for specific scenarios:
```bash
# DeFi load pattern
go run scripts/load-testing/load-test.go \
  -workers=30 \
  -tx-per-worker=500 \
  -tx-type="defi-swap"

# NFT minting pattern  
go run scripts/load-testing/load-test.go \
  -workers=10 \
  -tx-per-worker=1000 \
  -tx-type="nft-mint"
```

## Best Practices

### Test Environment
1. **Isolated Testing**: Use dedicated testnet
2. **Resource Monitoring**: Monitor system resources
3. **Baseline Establishment**: Run baseline tests first
4. **Gradual Scaling**: Increase load gradually

### Production Validation
1. **Pre-deployment**: Run full benchmark suite
2. **Canary Testing**: Start with small load
3. **Monitoring**: Continuous performance monitoring
4. **Rollback Plan**: Have rollback procedures ready

### Security Considerations
1. **Test Data**: Use test accounts only
2. **Network Isolation**: Isolate test networks
3. **Credential Management**: Secure test credentials
4. **Rate Limiting**: Respect API rate limits

---

For additional support or questions, please refer to the main DeshChain documentation or create an issue in the repository.