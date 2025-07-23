#!/usr/bin/env python3
"""
DeshChain Failure Analysis - Identify exact failure points
"""

import numpy as np
import json

class FailureAnalysis:
    def __init__(self):
        # From whitepaper projections
        self.year1_revenue = 275  # ₹275 Cr
        self.year1_monthly = 275 / 12  # ₹22.9 Cr/month
        self.fixed_costs = 18.3  # ₹18.3 Cr/month
        self.initial_investment = 500  # ₹500 Cr
        
    def analyze_baseline_scenario(self):
        """Analyze the baseline scenario from whitepaper"""
        print("BASELINE SCENARIO ANALYSIS (from Whitepaper)")
        print("=" * 60)
        
        months = 60  # 5 years
        cash = self.initial_investment
        
        for year in range(1, 6):
            # Whitepaper projections
            revenues = {
                1: 275,
                2: 1046,
                3: 3736,
                4: 7632,
                5: 12506
            }
            
            annual_revenue = revenues[year]
            monthly_revenue = annual_revenue / 12
            
            # Cost structure from whitepaper
            if year == 1:
                margin = 0.20  # 20% margin
            elif year == 2:
                margin = 0.40
            elif year == 3:
                margin = 0.50
            elif year == 4:
                margin = 0.55
            else:
                margin = 0.60
                
            annual_costs = annual_revenue * (1 - margin)
            monthly_costs = annual_costs / 12
            monthly_profit = monthly_revenue - monthly_costs
            
            print(f"\nYear {year}:")
            print(f"  Monthly Revenue: ₹{monthly_revenue:.1f} Cr")
            print(f"  Monthly Costs: ₹{monthly_costs:.1f} Cr")
            print(f"  Monthly Profit: ₹{monthly_profit:.1f} Cr")
            print(f"  Net Margin: {margin:.0%}")
            
            # Update cash for the year
            cash += annual_revenue - annual_costs
            print(f"  Year-end Cash: ₹{cash:.1f} Cr")
        
        return cash
    
    def identify_failure_points(self):
        """Identify specific failure conditions"""
        print("\n\nCRITICAL FAILURE POINT ANALYSIS")
        print("=" * 60)
        
        # 1. User acquisition failure
        print("\n1. USER ACQUISITION FAILURE:")
        print("   Required monthly active users for breakeven:")
        
        # From the model: need enough transaction volume
        # Assuming ₹10,000 avg transaction, 5 txns/user/month
        # Tax rate: 2.5%, Platform gets 25% of tax = 0.625%
        # Plus platform services revenue
        
        for month in [6, 12, 24, 36]:
            # Costs grow over time
            monthly_costs = self.fixed_costs * (1.05 ** (month/12))  # 5% annual inflation
            
            # To break even from transaction tax alone:
            # Revenue = Users * Txns * Amount * Tax * Platform_Share
            # monthly_costs = Users * 5 * 10000 * 0.025 * 0.25 / 1e7
            # monthly_costs = Users * 0.03125
            
            users_needed_tax_only = monthly_costs * 1e7 / (5 * 10000 * 0.025 * 0.25)
            
            # With platform revenue (est. ₹200/user/month)
            platform_rev_per_user = 200 / 1e7  # in Cr
            total_rev_per_user = 0.03125 + platform_rev_per_user
            users_needed_total = monthly_costs / total_rev_per_user
            
            print(f"   Month {month}: {users_needed_total:,.0f} users (costs: ₹{monthly_costs:.1f} Cr)")
        
        # 2. Fee compression failure
        print("\n2. FEE COMPRESSION FAILURE:")
        print("   Maximum sustainable fee reduction:")
        
        base_fee = 0.025
        for compression in [0.2, 0.4, 0.6, 0.8]:
            new_fee = base_fee * (1 - compression)
            revenue_impact = (1 - compression)
            print(f"   {compression:.0%} compression: {new_fee:.3%} fee (revenue down {1-revenue_impact:.0%})")
        
        # 3. Competition impact
        print("\n3. COMPETITION IMPACT:")
        print("   Market share required at different competition levels:")
        
        total_market = 50_000_000  # 50 Cr potential users
        for competition in [0.3, 0.5, 0.7, 0.9]:
            # Higher competition = need higher market share
            min_users = 500_000 * (1 + competition)
            market_share = min_users / total_market
            print(f"   {competition:.0%} competition: {market_share:.1%} market share ({min_users:,.0f} users)")
        
        # 4. Charity commitment impact
        print("\n4. CHARITY COMMITMENT ANALYSIS:")
        print("   Effective margin after 40% charity commitment:")
        
        for gross_margin in [0.3, 0.4, 0.5, 0.6, 0.7]:
            # 30% of tax + 10% of platform revenue to charity
            # Assuming 50/50 tax/platform revenue split
            charity_impact = 0.3 * 0.5 + 0.1 * 0.5  # 20% of total revenue
            net_margin = gross_margin - charity_impact
            print(f"   {gross_margin:.0%} gross → {net_margin:.0%} net margin")
            if net_margin < 0:
                print(f"     WARNING: Negative margin!")
        
        # 5. Burn rate analysis
        print("\n5. BURN RATE SUSTAINABILITY:")
        print("   Months until cash depletion at different burn rates:")
        
        initial_cash = 500  # ₹500 Cr
        for monthly_burn in [10, 20, 30, 40, 50]:
            months_to_zero = initial_cash / monthly_burn
            print(f"   ₹{monthly_burn} Cr/month burn: {months_to_zero:.1f} months ({months_to_zero/12:.1f} years)")
        
        # 6. Critical thresholds
        print("\n6. CRITICAL SURVIVAL THRESHOLDS:")
        print("   Minimum requirements for sustainability:")
        
        print(f"\n   USERS:")
        print(f"   - Year 1: 500,000+ active users")
        print(f"   - Year 2: 2,000,000+ active users")
        print(f"   - Year 3: 5,000,000+ active users")
        
        print(f"\n   REVENUE:")
        print(f"   - Break-even: ₹{self.fixed_costs:.1f} Cr/month")
        print(f"   - Sustainable: ₹{self.fixed_costs * 2:.1f} Cr/month (2x costs)")
        print(f"   - Growth mode: ₹{self.fixed_costs * 3:.1f} Cr/month (3x costs)")
        
        print(f"\n   MARKET CONDITIONS:")
        print(f"   - Crypto adoption: >20% of population")
        print(f"   - Regulatory support: >60% favorable")
        print(f"   - User churn: <15% monthly")
        print(f"   - Competition: <50% market saturation")
        
    def calculate_realistic_scenario(self):
        """Calculate a realistic survival scenario"""
        print("\n\nREALISTIC SURVIVAL SCENARIO")
        print("=" * 60)
        
        # More conservative assumptions
        print("Assumptions:")
        print("- Start with 50,000 users (not 100,000)")
        print("- 20% monthly growth (not 30%)")
        print("- ₹5,000 avg transaction (not ₹10,000)")
        print("- 3 transactions/user/month (not 5)")
        print("- 40% charity commitment maintained")
        print("- Operating costs: ₹15 Cr/month initially")
        
        months = 36
        users = 50_000
        cash = 500  # ₹500 Cr
        
        results = []
        
        for month in range(1, months + 1):
            # User growth
            users = users * 1.20 if month < 12 else users * 1.10
            
            # Revenue calculation
            active_users = users * 0.5  # 50% active
            txn_volume = active_users * 3 * 5000 / 1e7  # in Cr
            
            # Tax revenue (2.5%)
            tax_revenue = txn_volume * 0.025
            platform_tax_share = tax_revenue * 0.25  # Platform gets 25% of tax
            
            # Platform services revenue
            platform_revenue = users * 150 / 1e7  # ₹150/user
            
            total_revenue = platform_tax_share + platform_revenue
            
            # Costs
            operating_costs = 15 * (1.02 ** (month/12))  # 2% inflation
            
            # Charity (simplified: 20% of total revenue)
            charity = total_revenue * 0.20
            
            total_costs = operating_costs + charity
            
            # Profit/Loss
            profit = total_revenue - total_costs
            cash += profit
            
            if month % 6 == 0:
                print(f"\nMonth {month}:")
                print(f"  Users: {users:,.0f}")
                print(f"  Revenue: ₹{total_revenue:.1f} Cr")
                print(f"  Costs: ₹{total_costs:.1f} Cr")
                print(f"  Profit: ₹{profit:+.1f} Cr")
                print(f"  Cash: ₹{cash:.1f} Cr")
                
                if cash < 0:
                    print(f"  STATUS: FAILED - Out of cash!")
                    break
                elif profit > 0:
                    print(f"  STATUS: PROFITABLE!")
                else:
                    print(f"  STATUS: Burning ₹{-profit:.1f} Cr/month")
        
        return cash > 0

# Run analysis
if __name__ == "__main__":
    analyzer = FailureAnalysis()
    
    # Analyze baseline
    final_cash = analyzer.analyze_baseline_scenario()
    
    # Identify failure points
    analyzer.identify_failure_points()
    
    # Calculate realistic scenario
    survives = analyzer.calculate_realistic_scenario()
    
    print("\n" + "=" * 60)
    print("CONCLUSION:")
    print("=" * 60)
    
    if survives:
        print("✓ Model CAN survive with realistic assumptions")
    else:
        print("✗ Model WILL FAIL without significant adjustments")
    
    print("\nKEY FAILURE RISKS:")
    print("1. User acquisition below 500K in Year 1")
    print("2. Fee compression beyond 40%")
    print("3. Competition capturing >70% market")
    print("4. Operating costs exceeding ₹30 Cr/month")
    print("5. Charity commitment making margins negative")
    print("6. Regulatory shutdown in first 2 years")
    
    print("\nCRITICAL SUCCESS FACTORS:")
    print("1. Government partnerships for user acquisition")
    print("2. Cultural adoption driving retention")
    print("3. Operational efficiency improvements")
    print("4. Multiple revenue streams activation")
    print("5. Strong early traction (>20% monthly growth)")
    print("6. Maintain gross margins >50%")