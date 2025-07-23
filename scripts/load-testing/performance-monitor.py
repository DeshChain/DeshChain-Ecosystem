#!/usr/bin/env python3
"""
DeshChain Performance Monitoring Script
Real-time monitoring and alerting for blockchain performance
"""

import asyncio
import aiohttp
import json
import time
import argparse
import logging
import psutil
import subprocess
from datetime import datetime, timedelta
from typing import Dict, List, Optional
from dataclasses import dataclass, asdict

@dataclass
class PerformanceMetrics:
    timestamp: str
    block_height: int
    block_time: float
    tx_count: int
    tps: float
    memory_usage: float
    cpu_usage: float
    disk_usage: float
    peer_count: int
    sync_status: bool
    gas_usage: int
    validator_voting_power: int

@dataclass
class AlertThresholds:
    max_block_time: float = 10.0  # seconds
    min_tps: float = 10.0
    max_memory_usage: float = 80.0  # percentage
    max_cpu_usage: float = 80.0  # percentage
    max_disk_usage: float = 85.0  # percentage
    min_peer_count: int = 3

class DeshChainMonitor:
    def __init__(self, node_url: str, alert_webhook: Optional[str] = None):
        self.node_url = node_url.rstrip('/')
        self.alert_webhook = alert_webhook
        self.metrics_history: List[PerformanceMetrics] = []
        self.thresholds = AlertThresholds()
        self.alerts_sent = set()
        
        # Setup logging
        logging.basicConfig(
            level=logging.INFO,
            format='%(asctime)s - %(levelname)s - %(message)s',
            handlers=[
                logging.FileHandler(f'performance-monitor-{datetime.now().strftime("%Y%m%d")}.log'),
                logging.StreamHandler()
            ]
        )
        self.logger = logging.getLogger(__name__)

    async def fetch_node_status(self) -> Dict:
        """Fetch node status from RPC endpoint"""
        try:
            async with aiohttp.ClientSession() as session:
                async with session.get(f"{self.node_url}/status") as response:
                    if response.status == 200:
                        return await response.json()
                    else:
                        self.logger.error(f"Failed to fetch status: HTTP {response.status}")
                        return {}
        except Exception as e:
            self.logger.error(f"Error fetching node status: {e}")
            return {}

    async def fetch_net_info(self) -> Dict:
        """Fetch network info from RPC endpoint"""
        try:
            async with aiohttp.ClientSession() as session:
                async with session.get(f"{self.node_url}/net_info") as response:
                    if response.status == 200:
                        return await response.json()
                    else:
                        return {}
        except Exception as e:
            self.logger.error(f"Error fetching net info: {e}")
            return {}

    async def fetch_blockchain_info(self) -> Dict:
        """Fetch blockchain info"""
        try:
            async with aiohttp.ClientSession() as session:
                async with session.get(f"{self.node_url}/blockchain") as response:
                    if response.status == 200:
                        return await response.json()
                    else:
                        return {}
        except Exception as e:
            self.logger.error(f"Error fetching blockchain info: {e}")
            return {}

    def get_system_metrics(self) -> Dict:
        """Get system resource metrics"""
        try:
            return {
                'memory_usage': psutil.virtual_memory().percent,
                'cpu_usage': psutil.cpu_percent(interval=1),
                'disk_usage': psutil.disk_usage('/').percent
            }
        except Exception as e:
            self.logger.error(f"Error getting system metrics: {e}")
            return {'memory_usage': 0, 'cpu_usage': 0, 'disk_usage': 0}

    async def collect_metrics(self) -> Optional[PerformanceMetrics]:
        """Collect all performance metrics"""
        try:
            # Fetch blockchain data
            status_data = await self.fetch_node_status()
            net_data = await self.fetch_net_info()
            blockchain_data = await self.fetch_blockchain_info()
            system_metrics = self.get_system_metrics()

            if not status_data:
                return None

            # Extract metrics from status
            result = status_data.get('result', {})
            sync_info = result.get('sync_info', {})
            validator_info = result.get('validator_info', {})
            
            block_height = int(sync_info.get('latest_block_height', 0))
            latest_block_time = sync_info.get('latest_block_time', '')
            catching_up = sync_info.get('catching_up', True)
            
            # Calculate block time
            block_time = 0.0
            if len(self.metrics_history) > 0:
                prev_metrics = self.metrics_history[-1]
                time_diff = (datetime.now() - datetime.fromisoformat(prev_metrics.timestamp.replace('Z', '+00:00'))).total_seconds()
                height_diff = block_height - prev_metrics.block_height
                if height_diff > 0:
                    block_time = time_diff / height_diff

            # Get peer count
            peer_count = len(net_data.get('result', {}).get('peers', []))

            # Calculate TPS (approximate)
            tx_count = 0  # Would need to parse recent blocks for accurate count
            tps = 0.0
            if block_time > 0:
                tps = tx_count / block_time

            # Get validator voting power
            voting_power = int(validator_info.get('voting_power', 0))

            metrics = PerformanceMetrics(
                timestamp=datetime.now().isoformat(),
                block_height=block_height,
                block_time=block_time,
                tx_count=tx_count,
                tps=tps,
                memory_usage=system_metrics['memory_usage'],
                cpu_usage=system_metrics['cpu_usage'],
                disk_usage=system_metrics['disk_usage'],
                peer_count=peer_count,
                sync_status=not catching_up,
                gas_usage=0,  # Would need to analyze recent blocks
                validator_voting_power=voting_power
            )

            return metrics

        except Exception as e:
            self.logger.error(f"Error collecting metrics: {e}")
            return None

    async def send_alert(self, alert_type: str, message: str, severity: str = "warning"):
        """Send alert notification"""
        alert_key = f"{alert_type}_{severity}"
        
        # Prevent spam by only sending each alert type once per hour
        current_time = time.time()
        if alert_key in self.alerts_sent:
            last_sent = self.alerts_sent[alert_key]
            if current_time - last_sent < 3600:  # 1 hour
                return
        
        self.alerts_sent[alert_key] = current_time
        
        alert_data = {
            "timestamp": datetime.now().isoformat(),
            "alert_type": alert_type,
            "severity": severity,
            "message": message,
            "node_url": self.node_url
        }
        
        self.logger.warning(f"ALERT [{severity.upper()}] {alert_type}: {message}")
        
        if self.alert_webhook:
            try:
                async with aiohttp.ClientSession() as session:
                    async with session.post(self.alert_webhook, json=alert_data) as response:
                        if response.status == 200:
                            self.logger.info(f"Alert sent successfully: {alert_type}")
                        else:
                            self.logger.error(f"Failed to send alert: HTTP {response.status}")
            except Exception as e:
                self.logger.error(f"Error sending alert: {e}")

    async def check_alerts(self, metrics: PerformanceMetrics):
        """Check metrics against thresholds and send alerts"""
        
        # Block time alert
        if metrics.block_time > self.thresholds.max_block_time:
            await self.send_alert(
                "block_time", 
                f"Block time {metrics.block_time:.2f}s exceeds threshold {self.thresholds.max_block_time}s",
                "warning"
            )
        
        # TPS alert
        if metrics.tps < self.thresholds.min_tps and metrics.tps > 0:
            await self.send_alert(
                "low_tps",
                f"TPS {metrics.tps:.2f} below threshold {self.thresholds.min_tps}",
                "warning"
            )
        
        # Memory usage alert
        if metrics.memory_usage > self.thresholds.max_memory_usage:
            await self.send_alert(
                "high_memory",
                f"Memory usage {metrics.memory_usage:.1f}% exceeds threshold {self.thresholds.max_memory_usage}%",
                "critical" if metrics.memory_usage > 90 else "warning"
            )
        
        # CPU usage alert
        if metrics.cpu_usage > self.thresholds.max_cpu_usage:
            await self.send_alert(
                "high_cpu",
                f"CPU usage {metrics.cpu_usage:.1f}% exceeds threshold {self.thresholds.max_cpu_usage}%",
                "critical" if metrics.cpu_usage > 90 else "warning"
            )
        
        # Disk usage alert
        if metrics.disk_usage > self.thresholds.max_disk_usage:
            await self.send_alert(
                "high_disk",
                f"Disk usage {metrics.disk_usage:.1f}% exceeds threshold {self.thresholds.max_disk_usage}%",
                "critical"
            )
        
        # Peer count alert
        if metrics.peer_count < self.thresholds.min_peer_count:
            await self.send_alert(
                "low_peers",
                f"Peer count {metrics.peer_count} below threshold {self.thresholds.min_peer_count}",
                "warning"
            )
        
        # Sync status alert
        if not metrics.sync_status:
            await self.send_alert(
                "not_synced",
                "Node is not synced with network",
                "critical"
            )

    def print_metrics(self, metrics: PerformanceMetrics):
        """Print current metrics to console"""
        print(f"\n{'='*80}")
        print(f"DeshChain Performance Monitor - {metrics.timestamp}")
        print(f"{'='*80}")
        print(f"Block Height:      {metrics.block_height:,}")
        print(f"Block Time:        {metrics.block_time:.2f}s")
        print(f"TPS:               {metrics.tps:.2f}")
        print(f"Memory Usage:      {metrics.memory_usage:.1f}%")
        print(f"CPU Usage:         {metrics.cpu_usage:.1f}%")
        print(f"Disk Usage:        {metrics.disk_usage:.1f}%")
        print(f"Peer Count:        {metrics.peer_count}")
        print(f"Sync Status:       {'✅ Synced' if metrics.sync_status else '❌ Not Synced'}")
        print(f"Validator Power:   {metrics.validator_voting_power:,}")
        print(f"{'='*80}")

    def save_metrics(self, metrics: PerformanceMetrics, output_file: str):
        """Save metrics to file"""
        try:
            with open(output_file, 'a') as f:
                f.write(json.dumps(asdict(metrics)) + '\n')
        except Exception as e:
            self.logger.error(f"Error saving metrics: {e}")

    async def generate_report(self, output_file: str):
        """Generate performance report"""
        if not self.metrics_history:
            self.logger.warning("No metrics history available for report")
            return

        # Calculate statistics
        recent_metrics = self.metrics_history[-60:]  # Last 60 measurements
        
        avg_block_time = sum(m.block_time for m in recent_metrics if m.block_time > 0) / len([m for m in recent_metrics if m.block_time > 0])
        avg_tps = sum(m.tps for m in recent_metrics) / len(recent_metrics)
        avg_memory = sum(m.memory_usage for m in recent_metrics) / len(recent_metrics)
        avg_cpu = sum(m.cpu_usage for m in recent_metrics) / len(recent_metrics)
        
        max_memory = max(m.memory_usage for m in recent_metrics)
        max_cpu = max(m.cpu_usage for m in recent_metrics)
        
        report = {
            "report_timestamp": datetime.now().isoformat(),
            "monitoring_duration": len(self.metrics_history),
            "latest_metrics": asdict(self.metrics_history[-1]) if self.metrics_history else None,
            "statistics": {
                "average_block_time": avg_block_time,
                "average_tps": avg_tps,
                "average_memory_usage": avg_memory,
                "average_cpu_usage": avg_cpu,
                "peak_memory_usage": max_memory,
                "peak_cpu_usage": max_cpu
            },
            "thresholds": asdict(self.thresholds),
            "alerts_summary": {
                "total_alert_types": len(self.alerts_sent),
                "alert_types": list(set(alert.split('_')[0] for alert in self.alerts_sent.keys()))
            }
        }
        
        with open(output_file, 'w') as f:
            json.dump(report, f, indent=2)
        
        self.logger.info(f"Performance report saved to {output_file}")

    async def monitor(self, duration: int, interval: int, output_file: str = None):
        """Main monitoring loop"""
        self.logger.info(f"Starting DeshChain performance monitoring...")
        self.logger.info(f"Duration: {duration}s, Interval: {interval}s")
        self.logger.info(f"Node URL: {self.node_url}")
        
        start_time = time.time()
        
        try:
            while True:
                current_time = time.time()
                
                # Check if duration exceeded
                if duration > 0 and (current_time - start_time) >= duration:
                    break
                
                # Collect metrics
                metrics = await self.collect_metrics()
                
                if metrics:
                    self.metrics_history.append(metrics)
                    
                    # Keep only last 1000 metrics to prevent memory issues
                    if len(self.metrics_history) > 1000:
                        self.metrics_history = self.metrics_history[-1000:]
                    
                    # Print metrics
                    self.print_metrics(metrics)
                    
                    # Check for alerts
                    await self.check_alerts(metrics)
                    
                    # Save to file if specified
                    if output_file:
                        self.save_metrics(metrics, output_file)
                
                # Wait for next interval
                await asyncio.sleep(interval)
                
        except KeyboardInterrupt:
            self.logger.info("Monitoring stopped by user")
        except Exception as e:
            self.logger.error(f"Monitoring error: {e}")
        finally:
            # Generate final report
            report_file = f"performance-report-{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
            await self.generate_report(report_file)
            self.logger.info(f"Monitoring completed. Report saved to {report_file}")

def main():
    parser = argparse.ArgumentParser(description="DeshChain Performance Monitor")
    parser.add_argument("--node", default="http://localhost:26657", help="Node RPC URL")
    parser.add_argument("--duration", type=int, default=0, help="Monitoring duration in seconds (0 = infinite)")
    parser.add_argument("--interval", type=int, default=30, help="Monitoring interval in seconds")
    parser.add_argument("--output", help="Output file for metrics")
    parser.add_argument("--webhook", help="Alert webhook URL")
    parser.add_argument("--max-block-time", type=float, default=10.0, help="Max block time threshold")
    parser.add_argument("--min-tps", type=float, default=10.0, help="Min TPS threshold")
    parser.add_argument("--max-memory", type=float, default=80.0, help="Max memory usage threshold")
    parser.add_argument("--max-cpu", type=float, default=80.0, help="Max CPU usage threshold")
    
    args = parser.parse_args()
    
    # Create monitor
    monitor = DeshChainMonitor(args.node, args.webhook)
    
    # Update thresholds
    monitor.thresholds.max_block_time = args.max_block_time
    monitor.thresholds.min_tps = args.min_tps
    monitor.thresholds.max_memory_usage = args.max_memory
    monitor.thresholds.max_cpu_usage = args.max_cpu
    
    # Run monitoring
    asyncio.run(monitor.monitor(args.duration, args.interval, args.output))

if __name__ == "__main__":
    main()