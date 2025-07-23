#!/usr/bin/env python3
"""
DeshChain Monte Carlo Simulation
Stress testing financial model to identify failure points
"""

import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from dataclasses import dataclass
from typing import Dict, List, Tuple
import json

@dataclass
class MarketConditions:
    """Market parameters for simulation"""
    crypto_adoption_rate: float  # 0.0 to 1.0
    regulatory_favorability: float  # 0.0 to 1.0
    competition_intensity: float  # 0.0 to 1.0
    economic_growth: float  # -0.2 to 0.2
    inflation_rate: float  # 0.0 to 0.15
    user_churn_rate: float  # 0.0 to 0.5
    fraud_rate: float  # 0.0 to 0.1
    technology_failure_rate: float  # 0.0 to 0.05

@dataclass
class RevenueParameters:
    """Revenue model parameters"""
    transaction_volume_growth: float
    average_transaction_size: float
    fee_compression: float  # How much fees reduce over time
    market_share: float
    lending_default_rate: float
    pension_payout_ratio: float  # Actual vs promised returns
    operational_efficiency: float  # Cost reduction over time

class DeshChainSimulation:
    def __init__(self):
        # Base case parameters from whitepaper
        self.base_monthly_users = 10000  # Starting users
        self.base_transaction_volume = 275  # ₹275 Cr Year 1
        self.base_tax_rate = 0.025  # 2.5%
        self.charity_commitment = 0.40  # 40% to charity
        self.fixed_costs_monthly = 18.3  # ₹18.3 Cr
        self.initial_investment = 500  # ₹500 Cr
        
        # Revenue stream weights
        self.revenue_weights = {
            'lending': 0.30,
            'transaction_fees': 0.25,
            'defi_services': 0.20,
            'investment_returns': 0.15,
            'other': 0.10
        }
        
    def simulate_scenario(self, 
                         market: MarketConditions, 
                         revenue: RevenueParameters,
                         months: int = 120) -> Dict:
        """Run a single simulation scenario"""
        
        results = {
            'monthly_revenue': [],
            'monthly_costs': [],
            'monthly_profit': [],
            'cumulative_cash': [],
            'user_count': [],
            'failed': False,
            'failure_month': None,
            'failure_reason': None
        }
        
        # Initial conditions
        cash_balance = self.initial_investment
        current_users = self.base_monthly_users
        monthly_revenue = self.base_transaction_volume / 12
        
        for month in range(months):
            # User growth affected by market conditions
            user_growth_rate = (
                0.15 * market.crypto_adoption_rate * 
                (1 - market.competition_intensity) * 
                (1 + market.economic_growth) *
                market.regulatory_favorability
            )
            
            # Apply churn
            current_users = current_users * (1 + user_growth_rate - market.user_churn_rate)
            
            # Revenue calculations
            transaction_revenue = self._calculate_transaction_revenue(
                current_users, revenue, market, month
            )
            
            lending_revenue = self._calculate_lending_revenue(
                current_users, revenue, market
            )
            
            defi_revenue = self._calculate_defi_revenue(
                current_users, revenue, market
            )
            
            investment_revenue = self._calculate_investment_revenue(
                cash_balance, market
            )
            
            other_revenue = self._calculate_other_revenue(
                current_users, revenue, market
            )
            
            # Total revenue
            total_revenue = (
                transaction_revenue + lending_revenue + 
                defi_revenue + investment_revenue + other_revenue
            )
            
            # Costs calculation
            variable_costs = total_revenue * 0.15  # 15% of revenue
            fixed_costs = self.fixed_costs_monthly * (1 + market.inflation_rate / 12)
            fraud_losses = total_revenue * market.fraud_rate
            tech_failure_costs = total_revenue * market.technology_failure_rate * 2
            
            # Charity commitment (immutable)
            charity_payment = total_revenue * self.charity_commitment
            
            total_costs = (
                variable_costs + fixed_costs + fraud_losses + 
                tech_failure_costs + charity_payment
            )
            
            # Profit/Loss
            monthly_profit = total_revenue - total_costs
            cash_balance += monthly_profit
            
            # Store results
            results['monthly_revenue'].append(total_revenue)
            results['monthly_costs'].append(total_costs)
            results['monthly_profit'].append(monthly_profit)
            results['cumulative_cash'].append(cash_balance)
            results['user_count'].append(current_users)
            
            # Check failure conditions
            if cash_balance < -100:  # ₹100 Cr debt limit
                results['failed'] = True
                results['failure_month'] = month
                results['failure_reason'] = 'Insolvency'
                break
                
            if current_users < 1000:  # Minimum viable user base
                results['failed'] = True
                results['failure_month'] = month
                results['failure_reason'] = 'User base collapse'
                break
                
            if monthly_revenue < fixed_costs * 0.5:  # Revenue below 50% of fixed costs
                results['failed'] = True
                results['failure_month'] = month
                results['failure_reason'] = 'Unsustainable burn rate'
                break
                
            # Regulatory shutdown
            if market.regulatory_favorability < 0.1 and np.random.random() < 0.3:
                results['failed'] = True
                results['failure_month'] = month
                results['failure_reason'] = 'Regulatory shutdown'
                break
        
        return results
    
    def _calculate_transaction_revenue(self, users, revenue, market, month):
        """Calculate transaction fee revenue"""
        # Base transaction volume per user
        volume_per_user = 10000 * (1 + revenue.transaction_volume_growth) ** (month / 12)
        
        # Adjust for market conditions
        volume_per_user *= (1 - market.competition_intensity * 0.5)
        volume_per_user *= (1 + market.economic_growth)
        
        # Total volume
        total_volume = users * volume_per_user
        
        # Fee rate with compression
        fee_rate = self.base_tax_rate * (1 - revenue.fee_compression) ** (month / 12)
        
        # Platform share (25% of transaction fees)
        return total_volume * fee_rate * 0.25 / 1e7  # Convert to Cr
    
    def _calculate_lending_revenue(self, users, revenue, market):
        """Calculate lending revenue"""
        # Lending adoption rate
        lending_users = users * 0.2 * market.crypto_adoption_rate
        
        # Average loan size
        avg_loan = 50000 * (1 + market.economic_growth)
        
        # Interest margin after defaults
        interest_margin = 0.05 * (1 - revenue.lending_default_rate)
        
        # Monthly revenue
        return lending_users * avg_loan * interest_margin / 12 / 1e7
    
    def _calculate_defi_revenue(self, users, revenue, market):
        """Calculate DeFi services revenue"""
        # DeFi adoption
        defi_users = users * 0.15 * market.crypto_adoption_rate
        
        # Average DeFi volume
        avg_volume = 100000 * revenue.market_share
        
        # Fee rate
        fee_rate = 0.003 * (1 - revenue.fee_compression)
        
        return defi_users * avg_volume * fee_rate / 1e7
    
    def _calculate_investment_revenue(self, cash_balance, market):
        """Calculate treasury investment returns"""
        if cash_balance <= 0:
            return 0
        
        # Risk-adjusted returns
        return_rate = 0.08 * (1 - market.technology_failure_rate * 10)
        
        return cash_balance * return_rate / 12
    
    def _calculate_other_revenue(self, users, revenue, market):
        """Calculate other revenue streams"""
        # NFT, launchpad, etc.
        other_users = users * 0.05
        revenue_per_user = 5000 * market.crypto_adoption_rate
        
        return other_users * revenue_per_user / 1e7
    
    def run_monte_carlo(self, n_simulations: int = 10000) -> Dict:
        """Run Monte Carlo simulation"""
        
        failure_scenarios = []
        success_count = 0
        failure_reasons = {}
        failure_months = []
        
        for i in range(n_simulations):
            # Generate random market conditions
            market = MarketConditions(
                crypto_adoption_rate=np.random.beta(2, 5),  # Skewed low
                regulatory_favorability=np.random.beta(5, 2),  # Skewed high
                competition_intensity=np.random.beta(3, 3),  # Normal
                economic_growth=np.random.normal(0.05, 0.05),
                inflation_rate=np.random.gamma(2, 0.03),
                user_churn_rate=np.random.beta(2, 8),  # Low churn expected
                fraud_rate=np.random.beta(1, 20),  # Very low fraud
                technology_failure_rate=np.random.beta(1, 50)  # Very low tech failure
            )
            
            # Generate revenue parameters
            revenue = RevenueParameters(
                transaction_volume_growth=np.random.normal(0.3, 0.1),
                average_transaction_size=np.random.gamma(5, 2000),
                fee_compression=np.random.beta(2, 8),  # Slow compression
                market_share=np.random.beta(2, 18),  # 10% average share
                lending_default_rate=np.random.beta(2, 48),  # 4% default
                pension_payout_ratio=np.random.uniform(0.8, 1.2),
                operational_efficiency=np.random.beta(5, 2)  # Improving over time
            )
            
            # Run simulation
            result = self.simulate_scenario(market, revenue)
            
            if result['failed']:
                failure_scenarios.append({
                    'market': market,
                    'revenue': revenue,
                    'month': result['failure_month'],
                    'reason': result['failure_reason']
                })
                
                failure_reasons[result['failure_reason']] = \
                    failure_reasons.get(result['failure_reason'], 0) + 1
                
                failure_months.append(result['failure_month'])
            else:
                success_count += 1
        
        # Analyze results
        failure_rate = len(failure_scenarios) / n_simulations
        
        return {
            'failure_rate': failure_rate,
            'success_rate': 1 - failure_rate,
            'failure_reasons': failure_reasons,
            'avg_failure_month': np.mean(failure_months) if failure_months else None,
            'failure_scenarios': failure_scenarios,
            'critical_factors': self._identify_critical_factors(failure_scenarios)
        }
    
    def _identify_critical_factors(self, failure_scenarios: List[Dict]) -> Dict:
        """Identify which factors most contribute to failure"""
        
        if not failure_scenarios:
            return {}
        
        # Extract all parameters from failure scenarios
        factors = {
            'crypto_adoption_rate': [],
            'regulatory_favorability': [],
            'competition_intensity': [],
            'user_churn_rate': [],
            'lending_default_rate': [],
            'fee_compression': []
        }
        
        for scenario in failure_scenarios:
            factors['crypto_adoption_rate'].append(scenario['market'].crypto_adoption_rate)
            factors['regulatory_favorability'].append(scenario['market'].regulatory_favorability)
            factors['competition_intensity'].append(scenario['market'].competition_intensity)
            factors['user_churn_rate'].append(scenario['market'].user_churn_rate)
            factors['lending_default_rate'].append(scenario['revenue'].lending_default_rate)
            factors['fee_compression'].append(scenario['revenue'].fee_compression)
        
        # Calculate average values for failure scenarios
        critical_thresholds = {}
        for factor, values in factors.items():
            if values:
                critical_thresholds[factor] = {
                    'mean': np.mean(values),
                    'std': np.std(values),
                    'min': np.min(values),
                    'max': np.max(values)
                }
        
        return critical_thresholds

# Run simulation
if __name__ == "__main__":
    print("Running DeshChain Monte Carlo Simulation...")
    print("=" * 50)
    
    sim = DeshChainSimulation()
    results = sim.run_monte_carlo(n_simulations=10000)
    
    print(f"\nSimulation Results (10,000 runs):")
    print(f"Failure Rate: {results['failure_rate']:.2%}")
    print(f"Success Rate: {results['success_rate']:.2%}")
    
    if results['failure_reasons']:
        print(f"\nFailure Reasons:")
        for reason, count in sorted(results['failure_reasons'].items(), 
                                   key=lambda x: x[1], reverse=True):
            print(f"  {reason}: {count} ({count/sum(results['failure_reasons'].values()):.1%})")
    
    if results['avg_failure_month']:
        print(f"\nAverage Failure Month: {results['avg_failure_month']:.1f} ({results['avg_failure_month']/12:.1f} years)")
    
    print(f"\nCritical Failure Thresholds:")
    for factor, stats in results['critical_factors'].items():
        print(f"\n  {factor}:")
        print(f"    Mean in failures: {stats['mean']:.3f}")
        print(f"    Std deviation: {stats['std']:.3f}")
        print(f"    Range: {stats['min']:.3f} - {stats['max']:.3f}")
    
    # Save detailed results
    with open('monte_carlo_results.json', 'w') as f:
        # Convert to serializable format
        output = {
            'failure_rate': results['failure_rate'],
            'success_rate': results['success_rate'],
            'failure_reasons': results['failure_reasons'],
            'avg_failure_month': results['avg_failure_month'],
            'critical_factors': results['critical_factors'],
            'total_simulations': 10000
        }
        json.dump(output, f, indent=2)
    
    print("\nDetailed results saved to monte_carlo_results.json")