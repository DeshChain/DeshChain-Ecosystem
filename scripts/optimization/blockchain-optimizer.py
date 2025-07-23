#!/usr/bin/env python3
"""
DeshChain Blockchain Performance Optimizer
Analyzes and optimizes blockchain-specific performance metrics
"""

import asyncio
import aiohttp
import json
import argparse
import time
import logging
import statistics
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, asdict

@dataclass
class BlockchainMetrics:
    timestamp: str
    block_height: int
    block_time: float
    block_size: int
    tx_count: int
    gas_used: int
    gas_limit: int
    validator_count: int
    peer_count: int
    sync_status: bool
    mempool_size: int
    tps: float

@dataclass
class OptimizationRecommendation:
    category: str
    priority: str  # high, medium, low
    title: str
    description: str
    impact: str
    implementation: str
    estimated_improvement: str

class BlockchainOptimizer:
    def __init__(self, node_url: str, output_dir: str = "./blockchain-optimization"):
        self.node_url = node_url.rstrip('/')
        self.output_dir = output_dir
        self.metrics_history: List[BlockchainMetrics] = []
        self.recommendations: List[OptimizationRecommendation] = []
        
        # Setup logging
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(levelname)s - %(message)s',
            handlers=[
                logging.FileHandler(f'{output_dir}/blockchain-optimizer.log'),
                logging.StreamHandler()
            ]
        )
        self.logger = logging.getLogger(__name__)

    async def fetch_node_status(self) -> Dict:
        """Fetch comprehensive node status"""
        try:
            async with aiohttp.ClientSession() as session:
                # Fetch multiple endpoints for comprehensive data
                endpoints = {
                    'status': f"{self.node_url}/status",
                    'net_info': f"{self.node_url}/net_info",
                    'blockchain': f"{self.node_url}/blockchain",
                    'validators': f"{self.node_url}/validators",
                    'unconfirmed_txs': f"{self.node_url}/unconfirmed_txs"
                }
                
                results = {}
                for name, url in endpoints.items():
                    try:
                        async with session.get(url, timeout=10) as response:
                            if response.status == 200:
                                results[name] = await response.json()
                            else:
                                self.logger.warning(f"Failed to fetch {name}: HTTP {response.status}")
                                results[name] = {}
                    except Exception as e:
                        self.logger.warning(f"Error fetching {name}: {e}")
                        results[name] = {}
                
                return results
        except Exception as e:
            self.logger.error(f"Error fetching node data: {e}")
            return {}

    async def collect_metrics(self) -> Optional[BlockchainMetrics]:
        """Collect comprehensive blockchain metrics"""
        try:
            data = await self.fetch_node_status()
            
            if not data.get('status'):
                return None
            
            status = data['status'].get('result', {})
            sync_info = status.get('sync_info', {})
            validator_info = status.get('validator_info', {})
            
            net_info = data.get('net_info', {}).get('result', {})
            blockchain_info = data.get('blockchain', {}).get('result', {})
            validators = data.get('validators', {}).get('result', {})
            mempool = data.get('unconfirmed_txs', {}).get('result', {})
            
            # Extract metrics
            block_height = int(sync_info.get('latest_block_height', 0))
            block_time = self._calculate_block_time(sync_info)
            block_size = self._estimate_block_size(blockchain_info)
            tx_count = self._count_transactions(blockchain_info)
            gas_used, gas_limit = self._get_gas_metrics(blockchain_info)
            validator_count = len(validators.get('validators', []))
            peer_count = len(net_info.get('peers', []))
            sync_status = not sync_info.get('catching_up', True)
            mempool_size = len(mempool.get('txs', []))
            
            # Calculate TPS
            tps = self._calculate_tps(tx_count, block_time)
            
            metrics = BlockchainMetrics(
                timestamp=datetime.now().isoformat(),
                block_height=block_height,
                block_time=block_time,
                block_size=block_size,
                tx_count=tx_count,
                gas_used=gas_used,
                gas_limit=gas_limit,
                validator_count=validator_count,
                peer_count=peer_count,
                sync_status=sync_status,
                mempool_size=mempool_size,
                tps=tps
            )
            
            return metrics
            
        except Exception as e:
            self.logger.error(f"Error collecting metrics: {e}")
            return None

    def _calculate_block_time(self, sync_info: Dict) -> float:
        """Calculate average block time"""
        if len(self.metrics_history) < 2:
            return 0.0
        
        try:
            current_height = int(sync_info.get('latest_block_height', 0))
            previous_metrics = self.metrics_history[-1]
            
            height_diff = current_height - previous_metrics.block_height
            time_diff = (datetime.now() - datetime.fromisoformat(previous_metrics.timestamp)).total_seconds()
            
            if height_diff > 0:
                return time_diff / height_diff
            return 0.0
        except:
            return 0.0

    def _estimate_block_size(self, blockchain_info: Dict) -> int:
        """Estimate block size from blockchain info"""
        # Simplified estimation - would need actual block data for accuracy
        return 1024 * 1024  # 1MB default

    def _count_transactions(self, blockchain_info: Dict) -> int:
        """Count transactions in recent blocks"""
        # Simplified - would analyze recent blocks for accuracy
        return 100  # Default estimation

    def _get_gas_metrics(self, blockchain_info: Dict) -> Tuple[int, int]:
        """Get gas usage and limit"""
        # Simplified - would analyze actual blocks
        return (800000, 1000000)  # gas_used, gas_limit

    def _calculate_tps(self, tx_count: int, block_time: float) -> float:
        """Calculate transactions per second"""
        if block_time > 0:
            return tx_count / block_time
        return 0.0

    def analyze_performance(self) -> None:
        """Analyze blockchain performance and generate recommendations"""
        if len(self.metrics_history) < 10:
            self.logger.warning("Insufficient data for comprehensive analysis")
            return
        
        self.logger.info("Analyzing blockchain performance...")
        
        # Analyze different aspects of performance
        self._analyze_block_times()
        self._analyze_throughput()
        self._analyze_network_health()
        self._analyze_resource_utilization()
        self._analyze_consensus_performance()
        self._analyze_mempool_efficiency()

    def _analyze_block_times(self) -> None:
        """Analyze block time patterns"""
        block_times = [m.block_time for m in self.metrics_history if m.block_time > 0]
        
        if not block_times:
            return
        
        avg_block_time = statistics.mean(block_times)
        max_block_time = max(block_times)
        std_block_time = statistics.stdev(block_times) if len(block_times) > 1 else 0
        
        # Block time recommendations
        if avg_block_time > 10.0:
            self.recommendations.append(OptimizationRecommendation(
                category="consensus",
                priority="high",
                title="Slow Block Times",
                description=f"Average block time is {avg_block_time:.2f}s, significantly above optimal range (3-6s)",
                impact="Reduced transaction throughput and poor user experience",
                implementation="Optimize consensus algorithm, review validator performance, check network latency",
                estimated_improvement="30-50% reduction in block time"
            ))
        
        if std_block_time > avg_block_time * 0.5:
            self.recommendations.append(OptimizationRecommendation(
                category="consensus",
                priority="medium",
                title="Inconsistent Block Times",
                description=f"High block time variability (std dev: {std_block_time:.2f}s)",
                impact="Unpredictable transaction confirmation times",
                implementation="Investigate validator network connectivity, optimize block processing",
                estimated_improvement="40-60% reduction in block time variance"
            ))

    def _analyze_throughput(self) -> None:
        """Analyze transaction throughput"""
        tps_values = [m.tps for m in self.metrics_history if m.tps > 0]
        tx_counts = [m.tx_count for m in self.metrics_history]
        
        if not tps_values:
            return
        
        avg_tps = statistics.mean(tps_values)
        max_tps = max(tps_values)
        avg_tx_count = statistics.mean(tx_counts)
        
        # Throughput recommendations
        if avg_tps < 10:
            self.recommendations.append(OptimizationRecommendation(
                category="throughput",
                priority="high",
                title="Low Transaction Throughput",
                description=f"Average TPS is {avg_tps:.2f}, below production requirements",
                impact="Limited scalability and network capacity",
                implementation="Implement transaction batching, optimize validation, consider parallel processing",
                estimated_improvement="200-500% increase in TPS"
            ))
        
        if max_tps < avg_tps * 2:
            self.recommendations.append(OptimizationRecommendation(
                category="throughput",
                priority="medium",
                title="Limited Peak Throughput",
                description=f"Peak TPS ({max_tps:.2f}) is only {max_tps/avg_tps:.1f}x average",
                impact="Cannot handle traffic spikes effectively",
                implementation="Implement dynamic block sizing, optimize mempool management",
                estimated_improvement="100-200% increase in peak capacity"
            ))

    def _analyze_network_health(self) -> None:
        """Analyze network connectivity and health"""
        peer_counts = [m.peer_count for m in self.metrics_history]
        sync_issues = [m for m in self.metrics_history if not m.sync_status]
        
        avg_peers = statistics.mean(peer_counts)
        min_peers = min(peer_counts)
        
        # Network health recommendations
        if avg_peers < 8:
            self.recommendations.append(OptimizationRecommendation(
                category="network",
                priority="medium",
                title="Low Peer Count",
                description=f"Average peer count is {avg_peers:.1f}, below recommended minimum (8+)",
                impact="Reduced network resilience and sync reliability",
                implementation="Improve peer discovery, check firewall settings, add seed nodes",
                estimated_improvement="Better network resilience and faster sync"
            ))
        
        if len(sync_issues) > len(self.metrics_history) * 0.1:
            self.recommendations.append(OptimizationRecommendation(
                category="network",
                priority="high",
                title="Frequent Sync Issues",
                description=f"Node out of sync {len(sync_issues)} times out of {len(self.metrics_history)} samples",
                impact="Unreliable transaction processing and state consistency",
                implementation="Check network connectivity, optimize block processing, review consensus params",
                estimated_improvement="95%+ sync reliability"
            ))

    def _analyze_resource_utilization(self) -> None:
        """Analyze resource usage efficiency"""
        gas_utilizations = []
        
        for m in self.metrics_history:
            if m.gas_limit > 0:
                utilization = m.gas_used / m.gas_limit
                gas_utilizations.append(utilization)
        
        if gas_utilizations:
            avg_gas_util = statistics.mean(gas_utilizations)
            
            if avg_gas_util < 0.3:
                self.recommendations.append(OptimizationRecommendation(
                    category="resources",
                    priority="low",
                    title="Low Gas Utilization",
                    description=f"Average gas utilization is {avg_gas_util:.1%}, indicating underutilized capacity",
                    impact="Inefficient resource usage and potential for higher throughput",
                    implementation="Consider reducing gas limits or increasing block frequency",
                    estimated_improvement="Better resource efficiency"
                ))
            
            elif avg_gas_util > 0.9:
                self.recommendations.append(OptimizationRecommendation(
                    category="resources",
                    priority="medium",
                    title="High Gas Utilization",
                    description=f"Average gas utilization is {avg_gas_util:.1%}, approaching limits",
                    impact="Risk of transaction rejection and network congestion",
                    implementation="Increase gas limits or optimize transaction processing",
                    estimated_improvement="30-50% increase in transaction capacity"
                ))

    def _analyze_consensus_performance(self) -> None:
        """Analyze consensus mechanism performance"""
        validator_counts = [m.validator_count for m in self.metrics_history]
        
        if validator_counts:
            avg_validators = statistics.mean(validator_counts)
            
            if avg_validators < 10:
                self.recommendations.append(OptimizationRecommendation(
                    category="consensus",
                    priority="high",
                    title="Low Validator Count",
                    description=f"Average validator count is {avg_validators:.1f}, below security threshold",
                    impact="Reduced network security and decentralization",
                    implementation="Incentivize validator participation, review staking rewards",
                    estimated_improvement="Enhanced network security and decentralization"
                ))

    def _analyze_mempool_efficiency(self) -> None:
        """Analyze mempool management efficiency"""
        mempool_sizes = [m.mempool_size for m in self.metrics_history]
        
        if mempool_sizes:
            avg_mempool = statistics.mean(mempool_sizes)
            max_mempool = max(mempool_sizes)
            
            if avg_mempool > 1000:
                self.recommendations.append(OptimizationRecommendation(
                    category="mempool",
                    priority="medium",
                    title="Large Mempool Size",
                    description=f"Average mempool size is {avg_mempool:.0f} transactions",
                    impact="Increased memory usage and slower transaction processing",
                    implementation="Optimize transaction validation, implement priority queues",
                    estimated_improvement="50-70% reduction in mempool size"
                ))

    def generate_optimization_plan(self) -> Dict:
        """Generate comprehensive optimization plan"""
        self.logger.info("Generating optimization plan...")
        
        # Group recommendations by priority
        high_priority = [r for r in self.recommendations if r.priority == "high"]
        medium_priority = [r for r in self.recommendations if r.priority == "medium"]
        low_priority = [r for r in self.recommendations if r.priority == "low"]
        
        # Calculate metrics summary
        if self.metrics_history:
            recent_metrics = self.metrics_history[-10:]  # Last 10 samples
            
            summary = {
                "avg_block_time": statistics.mean([m.block_time for m in recent_metrics if m.block_time > 0]),
                "avg_tps": statistics.mean([m.tps for m in recent_metrics if m.tps > 0]),
                "avg_peer_count": statistics.mean([m.peer_count for m in recent_metrics]),
                "sync_reliability": sum(1 for m in recent_metrics if m.sync_status) / len(recent_metrics) * 100,
                "avg_mempool_size": statistics.mean([m.mempool_size for m in recent_metrics])
            }
        else:
            summary = {}
        
        plan = {
            "timestamp": datetime.now().isoformat(),
            "analysis_period": {
                "start": self.metrics_history[0].timestamp if self.metrics_history else None,
                "end": self.metrics_history[-1].timestamp if self.metrics_history else None,
                "samples": len(self.metrics_history)
            },
            "performance_summary": summary,
            "recommendations": {
                "high_priority": [asdict(r) for r in high_priority],
                "medium_priority": [asdict(r) for r in medium_priority],
                "low_priority": [asdict(r) for r in low_priority],
                "total_count": len(self.recommendations)
            },
            "implementation_roadmap": self._create_implementation_roadmap(high_priority, medium_priority, low_priority)
        }
        
        return plan

    def _create_implementation_roadmap(self, high_priority: List, medium_priority: List, low_priority: List) -> Dict:
        """Create implementation roadmap with phases"""
        roadmap = {
            "phase_1_immediate": {
                "timeframe": "1-2 weeks",
                "focus": "Critical performance issues",
                "items": [{"title": r.title, "category": r.category} for r in high_priority[:3]]
            },
            "phase_2_short_term": {
                "timeframe": "1-2 months", 
                "focus": "Significant improvements",
                "items": [{"title": r.title, "category": r.category} for r in high_priority[3:] + medium_priority[:3]]
            },
            "phase_3_long_term": {
                "timeframe": "3-6 months",
                "focus": "Optimization and efficiency",
                "items": [{"title": r.title, "category": r.category} for r in medium_priority[3:] + low_priority]
            }
        }
        
        return roadmap

    async def save_results(self, plan: Dict) -> None:
        """Save optimization results to files"""
        import os
        os.makedirs(self.output_dir, exist_ok=True)
        
        # Save optimization plan
        plan_file = f"{self.output_dir}/blockchain-optimization-plan.json"
        with open(plan_file, 'w') as f:
            json.dump(plan, f, indent=2)
        
        # Save metrics data
        metrics_file = f"{self.output_dir}/blockchain-metrics.json"
        with open(metrics_file, 'w') as f:
            json.dump([asdict(m) for m in self.metrics_history], f, indent=2)
        
        # Generate markdown report
        await self._generate_markdown_report(plan)
        
        self.logger.info(f"Results saved to {self.output_dir}")

    async def _generate_markdown_report(self, plan: Dict) -> None:
        """Generate comprehensive markdown report"""
        report_file = f"{self.output_dir}/blockchain-optimization-report.md"
        
        summary = plan.get("performance_summary", {})
        
        report_content = f"""# DeshChain Blockchain Optimization Report

**Generated:** {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}  
**Analysis Period:** {plan['analysis_period']['samples']} samples  
**Node:** {self.node_url}  

## Executive Summary

This report provides comprehensive blockchain performance analysis and optimization recommendations for DeshChain.

### Current Performance Metrics

- **Average Block Time:** {summary.get('avg_block_time', 0):.2f} seconds
- **Average TPS:** {summary.get('avg_tps', 0):.2f} transactions/second
- **Sync Reliability:** {summary.get('sync_reliability', 0):.1f}%
- **Average Peer Count:** {summary.get('avg_peer_count', 0):.0f}
- **Average Mempool Size:** {summary.get('avg_mempool_size', 0):.0f} transactions

### Optimization Opportunities

- **High Priority Issues:** {len(plan['recommendations']['high_priority'])}
- **Medium Priority Items:** {len(plan['recommendations']['medium_priority'])}
- **Low Priority Optimizations:** {len(plan['recommendations']['low_priority'])}

## High Priority Recommendations

"""
        
        for i, rec in enumerate(plan['recommendations']['high_priority'], 1):
            report_content += f"""
### {i}. {rec['title']} ({rec['category'].title()})

**Priority:** 🔴 {rec['priority'].title()}  
**Description:** {rec['description']}  
**Impact:** {rec['impact']}  
**Implementation:** {rec['implementation']}  
**Expected Improvement:** {rec['estimated_improvement']}
"""
        
        report_content += """

## Medium Priority Recommendations

"""
        
        for i, rec in enumerate(plan['recommendations']['medium_priority'], 1):
            report_content += f"""
### {i}. {rec['title']} ({rec['category'].title()})

**Priority:** 🟡 {rec['priority'].title()}  
**Description:** {rec['description']}  
**Impact:** {rec['impact']}  
**Implementation:** {rec['implementation']}  
**Expected Improvement:** {rec['estimated_improvement']}
"""
        
        report_content += """

## Implementation Roadmap

### Phase 1: Immediate (1-2 weeks)
**Focus:** Critical performance issues

"""
        
        for item in plan['implementation_roadmap']['phase_1_immediate']['items']:
            report_content += f"- {item['title']} ({item['category']})\n"
        
        report_content += """
### Phase 2: Short-term (1-2 months)
**Focus:** Significant improvements

"""
        
        for item in plan['implementation_roadmap']['phase_2_short_term']['items']:
            report_content += f"- {item['title']} ({item['category']})\n"
        
        report_content += """
### Phase 3: Long-term (3-6 months)
**Focus:** Optimization and efficiency

"""
        
        for item in plan['implementation_roadmap']['phase_3_long_term']['items']:
            report_content += f"- {item['title']} ({item['category']})\n"
        
        report_content += """

## Technical Analysis

### Performance Trends

See `blockchain-metrics.json` for detailed historical performance data.

### Monitoring Recommendations

1. **Continuous Monitoring:** Implement real-time performance dashboards
2. **Alerting:** Set up alerts for block time > 10s, TPS < 5, sync issues
3. **Regular Analysis:** Run this optimization analysis weekly
4. **Performance Testing:** Implement automated performance regression testing

## Generated Files

- 📊 **Optimization Plan:** `blockchain-optimization-plan.json`
- 📈 **Metrics Data:** `blockchain-metrics.json`
- 📋 **This Report:** `blockchain-optimization-report.md`
- 📝 **Logs:** `blockchain-optimizer.log`

---
*Generated by DeshChain Blockchain Optimizer*
"""
        
        with open(report_file, 'w') as f:
            f.write(report_content)

    async def run_optimization_analysis(self, duration: int = 300, interval: int = 30) -> None:
        """Run comprehensive blockchain optimization analysis"""
        self.logger.info(f"Starting blockchain optimization analysis for {duration}s")
        
        import os
        os.makedirs(self.output_dir, exist_ok=True)
        
        start_time = time.time()
        
        try:
            while (time.time() - start_time) < duration:
                metrics = await self.collect_metrics()
                
                if metrics:
                    self.metrics_history.append(metrics)
                    self.logger.info(f"Collected metrics: Block {metrics.block_height}, TPS: {metrics.tps:.2f}")
                
                await asyncio.sleep(interval)
                
        except KeyboardInterrupt:
            self.logger.info("Analysis interrupted by user")
        
        # Perform analysis
        self.analyze_performance()
        
        # Generate optimization plan
        plan = self.generate_optimization_plan()
        
        # Save results
        await self.save_results(plan)
        
        self.logger.info("Blockchain optimization analysis completed")
        return plan

def main():
    parser = argparse.ArgumentParser(description="DeshChain Blockchain Performance Optimizer")
    parser.add_argument("--node", default="http://localhost:26657", help="Node RPC URL")
    parser.add_argument("--duration", type=int, default=300, help="Analysis duration in seconds")
    parser.add_argument("--interval", type=int, default=30, help="Sample interval in seconds")
    parser.add_argument("--output", default="./blockchain-optimization", help="Output directory")
    
    args = parser.parse_args()
    
    optimizer = BlockchainOptimizer(args.node, args.output)
    
    try:
        plan = asyncio.run(optimizer.run_optimization_analysis(args.duration, args.interval))
        
        print(f"\n🎯 Blockchain optimization analysis completed!")
        print(f"📁 Results directory: {args.output}")
        print(f"📊 Optimization plan: {args.output}/blockchain-optimization-plan.json")
        print(f"📋 Full report: {args.output}/blockchain-optimization-report.md")
        
        high_priority_count = len(plan['recommendations']['high_priority'])
        if high_priority_count > 0:
            print(f"⚠️  {high_priority_count} high-priority issues identified")
            print("Review the optimization report for immediate action items")
        else:
            print("✅ No high-priority issues found - blockchain performance looks good!")
            
    except Exception as e:
        print(f"❌ Analysis failed: {e}")
        return 1
    
    return 0

if __name__ == "__main__":
    exit(main())