#!/usr/bin/env python3
"""
DeshChain Realistic Monte Carlo Simulation
Fixed parameters based on actual documentation
"""

import numpy as np
import pandas as pd
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
    government_support: float  # 0.0 to 1.0 (NEW)

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
    cultural_adoption: float  # 0.0 to 1.0 (NEW)

class DeshChainRealisticSimulation:
    def __init__(self):
        # CORRECTED base case parameters from whitepaper
        self.base_monthly_users = 100000  # More realistic starting users
        self.base_transaction_volume = 23  # ₹23 Cr/month (₹275 Cr Year 1)
        self.base_tax_rate = 0.025  # 2.5%
        self.charity_percent_of_tax = 0.30  # 30% of tax goes to charity
        self.charity_percent_of_platform = 0.10  # 10% of platform revenue to charity
        self.fixed_costs_monthly = 18.3  # ₹18.3 Cr
        self.initial_investment = 500  # ₹500 Cr
        
        # Corrected revenue model
        self.tax_distribution = {
            'ngos': 0.30,  # 30% of tax
            'validators': 0.25,
            'community': 0.20,
            'operations': 0.25
        }
        
        self.platform_revenue_distribution = {
            'development': 0.30,
            'treasury': 0.25,
            'liquidity': 0.20,
            'ngos': 0.10,
            'emergency': 0.10,
            'founder': 0.05
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
            'failure_reason': None,
            'tax_revenue': [],
            'platform_revenue': []
        }
        
        # Initial conditions
        cash_balance = self.initial_investment
        current_users = self.base_monthly_users
        
        for month in range(months):
            # User growth with cultural adoption factor
            base_growth = 0.15 if month < 24 else 0.10  # Higher early growth
            
            user_growth_rate = (
                base_growth * 
                (0.5 + 0.5 * market.crypto_adoption_rate) *  # Less dependent on crypto
                (1 - market.competition_intensity * 0.3) *  # Less affected by competition
                (1 + market.economic_growth) *
                (0.5 + 0.5 * market.regulatory_favorability) *  # Less dependent on regulation
                (0.8 + 0.2 * revenue.cultural_adoption) *  # Cultural factor
                (0.7 + 0.3 * market.government_support)  # Government support
            )
            
            # Apply churn (reduced by cultural adoption)
            effective_churn = market.user_churn_rate * (1 - revenue.cultural_adoption * 0.5)
            current_users = current_users * (1 + user_growth_rate - effective_churn)
            
            # Calculate transaction volume (more realistic)
            user_activity_rate = 0.3 + 0.4 * revenue.cultural_adoption  # 30-70% active
            active_users = current_users * user_activity_rate
            
            # Average transaction per active user per month
            txns_per_user = 5 * (1 + market.economic_growth)
            avg_txn_size = revenue.average_transaction_size
            
            total_volume = active_users * txns_per_user * avg_txn_size / 1e7  # in Cr
            
            # TAX REVENUE (from all transactions)
            tax_rate = self.base_tax_rate * (1 - revenue.fee_compression * 0.5)  # Max 50% compression
            tax_revenue = total_volume * tax_rate
            
            # Tax goes to various parties, platform gets operational share
            platform_tax_share = tax_revenue * self.tax_distribution['operations']
            
            # PLATFORM REVENUE (from services)
            # 1. Lending revenue
            lending_users = current_users * 0.15 * (0.5 + 0.5 * market.crypto_adoption_rate)
            avg_loan = 50000
            interest_margin = 0.03 * (1 - revenue.lending_default_rate)
            lending_revenue = lending_users * avg_loan * interest_margin / 12 / 1e7
            
            # 2. DEX revenue (0.05% of volume)
            dex_volume = total_volume * 0.2  # 20% of transactions through DEX
            dex_revenue = dex_volume * 0.0005
            
            # 3. Pension scheme revenue
            pension_users = current_users * 0.10
            avg_pension_deposit = 10000
            pension_margin = 0.02  # 2% management fee
            pension_revenue = pension_users * avg_pension_deposit * pension_margin / 12 / 1e7
            
            # 4. Other services
            other_revenue = current_users * 100 * revenue.market_share / 1e7  # ₹100 per user
            
            # Total platform revenue
            platform_revenue = lending_revenue + dex_revenue + pension_revenue + other_revenue
            
            # TOTAL REVENUE
            total_revenue = platform_tax_share + platform_revenue
            
            # COSTS calculation
            # Variable costs scale with activity
            variable_costs = total_revenue * 0.25  # 25% of revenue
            
            # Fixed costs with efficiency improvements
            efficiency_factor = 1 - (revenue.operational_efficiency * min(month / 60, 1) * 0.3)
            fixed_costs = self.fixed_costs_monthly * (1 + market.inflation_rate / 12) * efficiency_factor
            
            # Fraud and tech failure costs
            fraud_losses = total_revenue * market.fraud_rate
            tech_failure_costs = total_revenue * market.technology_failure_rate * 2
            
            # Charity commitments
            charity_from_tax = tax_revenue * self.charity_percent_of_tax
            charity_from_platform = platform_revenue * self.charity_percent_of_platform
            total_charity = charity_from_tax + charity_from_platform
            
            total_costs = (
                variable_costs + fixed_costs + fraud_losses + 
                tech_failure_costs + total_charity
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
            results['tax_revenue'].append(tax_revenue)
            results['platform_revenue'].append(platform_revenue)
            
            # Check failure conditions (more realistic)
            if cash_balance < -200:  # ₹200 Cr debt limit (can raise more funds)
                results['failed'] = True
                results['failure_month'] = month
                results['failure_reason'] = 'Insolvency'
                break
                
            if current_users < 10000:  # Minimum viable user base
                results['failed'] = True
                results['failure_month'] = month
                results['failure_reason'] = 'User base collapse'
                break
                
            if total_revenue < fixed_costs * 0.3 and month > 36:  # After 3 years
                results['failed'] = True
                results['failure_month'] = month
                results['failure_reason'] = 'Unsustainable burn rate'
                break
                
            # Regulatory shutdown (rare with government support)
            shutdown_probability = 0.05 * (1 - market.government_support)
            if market.regulatory_favorability < 0.2 and np.random.random() < shutdown_probability:
                results['failed'] = True
                results['failure_month'] = month
                results['failure_reason'] = 'Regulatory shutdown'
                break
                
            # Catastrophic tech failure
            if market.technology_failure_rate > 0.03 and np.random.random() < 0.01:
                results['failed'] = True
                results['failure_month'] = month
                results['failure_reason'] = 'Catastrophic tech failure'
                break
        
        return results
    
    def run_monte_carlo(self, n_simulations: int = 10000) -> Dict:
        """Run Monte Carlo simulation with realistic parameters"""
        
        failure_scenarios = []
        success_count = 0
        failure_reasons = {}
        failure_months = []
        successful_outcomes = []
        
        for i in range(n_simulations):
            # Generate more realistic market conditions
            market = MarketConditions(
                crypto_adoption_rate=np.random.beta(3, 7),  # 30% average
                regulatory_favorability=np.random.beta(7, 3),  # 70% favorable
                competition_intensity=np.random.beta(3, 5),  # 37% average
                economic_growth=np.random.normal(0.06, 0.03),  # 6% growth
                inflation_rate=np.random.gamma(2, 0.025),  # 5% inflation
                user_churn_rate=np.random.beta(2, 18),  # 10% churn
                fraud_rate=np.random.beta(1, 99),  # 1% fraud
                technology_failure_rate=np.random.beta(1, 199),  # 0.5% failure
                government_support=np.random.beta(6, 4)  # 60% support
            )
            
            # Generate revenue parameters
            revenue = RevenueParameters(
                transaction_volume_growth=np.random.normal(0.25, 0.05),  # 25% growth
                average_transaction_size=np.random.gamma(5, 2000),  # ₹10,000 average
                fee_compression=np.random.beta(2, 8),  # 20% compression over time
                market_share=np.random.beta(3, 27),  # 10% market share
                lending_default_rate=np.random.beta(2, 48),  # 4% default
                pension_payout_ratio=np.random.uniform(0.9, 1.1),  # ±10% variation
                operational_efficiency=np.random.beta(7, 3),  # 70% efficiency gains
                cultural_adoption=np.random.beta(7, 3)  # 70% cultural adoption
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
                # Track successful outcome metrics
                final_users = result['user_count'][-1]
                final_revenue = result['monthly_revenue'][-1]
                total_profit = sum(result['monthly_profit'])
                successful_outcomes.append({
                    'final_users': final_users,
                    'final_monthly_revenue': final_revenue,
                    'total_profit': total_profit,
                    'market': market,
                    'revenue': revenue
                })
        
        # Analyze results
        failure_rate = len(failure_scenarios) / n_simulations
        
        # Analyze successful outcomes
        success_metrics = {}
        if successful_outcomes:
            success_metrics = {
                'avg_final_users': np.mean([s['final_users'] for s in successful_outcomes]),
                'avg_final_revenue': np.mean([s['final_monthly_revenue'] for s in successful_outcomes]),
                'avg_total_profit': np.mean([s['total_profit'] for s in successful_outcomes]),
                'min_total_profit': np.min([s['total_profit'] for s in successful_outcomes]),
                'max_total_profit': np.max([s['total_profit'] for s in successful_outcomes])
            }
        
        return {
            'failure_rate': failure_rate,
            'success_rate': 1 - failure_rate,
            'failure_reasons': failure_reasons,
            'avg_failure_month': np.mean(failure_months) if failure_months else None,
            'failure_scenarios': failure_scenarios,
            'critical_factors': self._identify_critical_factors(failure_scenarios, successful_outcomes),
            'success_metrics': success_metrics
        }
    
    def _identify_critical_factors(self, failure_scenarios: List[Dict], 
                                  successful_outcomes: List[Dict]) -> Dict:
        """Compare failure vs success scenarios"""
        
        factors = {
            'failure_thresholds': {},
            'success_thresholds': {},
            'critical_differences': {}
        }
        
        # Analyze failure scenarios
        if failure_scenarios:
            failure_params = {
                'crypto_adoption_rate': [s['market'].crypto_adoption_rate for s in failure_scenarios],
                'user_churn_rate': [s['market'].user_churn_rate for s in failure_scenarios],
                'cultural_adoption': [s['revenue'].cultural_adoption for s in failure_scenarios],
                'government_support': [s['market'].government_support for s in failure_scenarios]
            }
            
            for param, values in failure_params.items():
                factors['failure_thresholds'][param] = {
                    'mean': np.mean(values),
                    'std': np.std(values)
                }
        
        # Analyze successful scenarios
        if successful_outcomes:
            success_params = {
                'crypto_adoption_rate': [s['market'].crypto_adoption_rate for s in successful_outcomes],
                'user_churn_rate': [s['market'].user_churn_rate for s in successful_outcomes],
                'cultural_adoption': [s['revenue'].cultural_adoption for s in successful_outcomes],
                'government_support': [s['market'].government_support for s in successful_outcomes]
            }
            
            for param, values in success_params.items():
                factors['success_thresholds'][param] = {
                    'mean': np.mean(values),
                    'std': np.std(values)
                }
        
        # Calculate critical differences
        for param in ['crypto_adoption_rate', 'user_churn_rate', 'cultural_adoption', 'government_support']:
            if param in factors['failure_thresholds'] and param in factors['success_thresholds']:
                factors['critical_differences'][param] = (
                    factors['success_thresholds'][param]['mean'] - 
                    factors['failure_thresholds'][param]['mean']
                )
        
        return factors

# Run simulation
if __name__ == "__main__":
    print("Running DeshChain Realistic Monte Carlo Simulation...")
    print("=" * 60)
    
    sim = DeshChainRealisticSimulation()
    results = sim.run_monte_carlo(n_simulations=10000)
    
    print(f"\nSimulation Results (10,000 runs):")
    print(f"Failure Rate: {results['failure_rate']:.2%}")
    print(f"Success Rate: {results['success_rate']:.2%}")
    
    if results['failure_reasons']:
        print(f"\nFailure Reasons:")
        total_failures = sum(results['failure_reasons'].values())
        for reason, count in sorted(results['failure_reasons'].items(), 
                                   key=lambda x: x[1], reverse=True):
            print(f"  {reason}: {count} ({count/total_failures:.1%})")
    
    if results['avg_failure_month']:
        print(f"\nAverage Failure Month: {results['avg_failure_month']:.1f} ({results['avg_failure_month']/12:.1f} years)")
    
    print(f"\nCritical Success Factors:")
    if results['critical_factors']['critical_differences']:
        for factor, diff in results['critical_factors']['critical_differences'].items():
            print(f"  {factor}: {diff:+.3f} difference between success/failure")
    
    if results['success_metrics']:
        print(f"\nSuccessful Outcome Metrics:")
        print(f"  Average Final Users: {results['success_metrics']['avg_final_users']:,.0f}")
        print(f"  Average Monthly Revenue: ₹{results['success_metrics']['avg_final_revenue']:.1f} Cr")
        print(f"  Average Total Profit (10 years): ₹{results['success_metrics']['avg_total_profit']:,.0f} Cr")
        print(f"  Profit Range: ₹{results['success_metrics']['min_total_profit']:,.0f} - ₹{results['success_metrics']['max_total_profit']:,.0f} Cr")
    
    # Save detailed results
    with open('monte_carlo_realistic_results.json', 'w') as f:
        output = {
            'failure_rate': results['failure_rate'],
            'success_rate': results['success_rate'],
            'failure_reasons': results['failure_reasons'],
            'avg_failure_month': results['avg_failure_month'],
            'critical_factors': {
                'failure_thresholds': results['critical_factors']['failure_thresholds'],
                'success_thresholds': results['critical_factors']['success_thresholds'],
                'critical_differences': results['critical_factors']['critical_differences']
            },
            'success_metrics': results['success_metrics'],
            'total_simulations': 10000
        }
        json.dump(output, f, indent=2)
    
    print("\nDetailed results saved to monte_carlo_realistic_results.json")
    
    # Identify specific failure conditions
    print("\n" + "="*60)
    print("CRITICAL FAILURE CONDITIONS:")
    print("="*60)
    
    if results['critical_factors']['failure_thresholds']:
        print("\nWhen these conditions occur together, failure is likely:")
        for param, stats in results['critical_factors']['failure_thresholds'].items():
            success_mean = results['critical_factors']['success_thresholds'].get(param, {}).get('mean', 0)
            print(f"\n{param}:")
            print(f"  Failure avg: {stats['mean']:.3f}")
            print(f"  Success avg: {success_mean:.3f}")
            if stats['mean'] < success_mean:
                print(f"  CRITICAL: Keep above {stats['mean'] + stats['std']:.3f}")
            else:
                print(f"  CRITICAL: Keep below {stats['mean'] - stats['std']:.3f}")