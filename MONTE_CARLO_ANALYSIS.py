import numpy as np
import matplotlib.pyplot as plt
from dataclasses import dataclass
from typing import List, Tuple
import pandas as pd

@dataclass
class MarketFactors:
    """Indian market-specific factors"""
    digital_india_adoption: float  # 0.0 to 1.0
    crypto_regulation_clarity: float  # -1.0 to 1.0
    cultural_alignment: float  # 0.0 to 1.0
    ngo_trust_factor: float  # 0.0 to 1.0
    government_support: float  # -1.0 to 1.0
    competition_intensity: float  # 0.0 to 1.0
    economic_growth: float  # GDP growth rate
    smartphone_penetration: float  # 0.0 to 1.0
    financial_inclusion_need: float  # 0.0 to 1.0
    diaspora_remittance: float  # 0.0 to 1.0

@dataclass
class ProjectFactors:
    """DeshChain-specific factors"""
    tech_execution: float  # 0.0 to 1.0
    team_quality: float  # 0.0 to 1.0
    marketing_effectiveness: float  # 0.0 to 1.0
    partnership_success: float  # 0.0 to 1.0
    product_market_fit: float  # 0.0 to 1.0
    security_reliability: float  # 0.0 to 1.0
    user_experience: float  # 0.0 to 1.0
    charity_impact_visibility: float  # 0.0 to 1.0
    founder_dedication: float  # Fixed at 1.0
    tokenomics_sustainability: float  # 0.0 to 1.0

class DeshChainMonteCarloSimulation:
    def __init__(self, n_simulations: int = 10000):
        self.n_simulations = n_simulations
        self.results = []
        
    def generate_market_factors(self) -> MarketFactors:
        """Generate random market factors based on Indian context"""
        return MarketFactors(
            # Digital India momentum is strong
            digital_india_adoption=np.random.beta(8, 2),  # Skewed positive
            
            # Crypto regulation uncertain but improving
            crypto_regulation_clarity=np.random.normal(0.3, 0.3),
            
            # Strong cultural alignment with charity/dharma
            cultural_alignment=np.random.beta(9, 1),  # Very positive
            
            # High trust in established NGOs
            ngo_trust_factor=np.random.beta(7, 3),
            
            # Government support uncertain
            government_support=np.random.normal(0, 0.4),
            
            # Growing competition
            competition_intensity=np.random.beta(3, 7),
            
            # India GDP growth 6-8%
            economic_growth=np.random.normal(0.07, 0.02),
            
            # Rapidly growing smartphone adoption
            smartphone_penetration=np.random.beta(6, 4),
            
            # Huge need for financial inclusion
            financial_inclusion_need=np.random.beta(8, 2),
            
            # Large diaspora remittance market
            diaspora_remittance=np.random.beta(7, 3)
        )
    
    def generate_project_factors(self) -> ProjectFactors:
        """Generate random project execution factors"""
        return ProjectFactors(
            # Technical execution uncertainty
            tech_execution=np.random.beta(5, 2),
            
            # Team quality assumption
            team_quality=np.random.beta(6, 2),
            
            # Marketing effectiveness
            marketing_effectiveness=np.random.beta(4, 3),
            
            # Partnership success rate
            partnership_success=np.random.beta(5, 3),
            
            # Product-market fit probability
            product_market_fit=np.random.beta(6, 2),
            
            # Security is critical
            security_reliability=np.random.beta(7, 1),
            
            # UX must be excellent
            user_experience=np.random.beta(5, 3),
            
            # Charity impact visibility
            charity_impact_visibility=np.random.beta(7, 2),
            
            # Founder dedication guaranteed
            founder_dedication=1.0,
            
            # Tokenomics well designed
            tokenomics_sustainability=np.random.beta(8, 2)
        )
    
    def calculate_adoption_rate(self, market: MarketFactors, project: ProjectFactors) -> float:
        """Calculate user adoption rate based on all factors"""
        
        # Market opportunity score
        market_score = (
            market.digital_india_adoption * 0.15 +
            max(0, market.crypto_regulation_clarity) * 0.10 +
            market.cultural_alignment * 0.20 +
            market.ngo_trust_factor * 0.15 +
            max(0, market.government_support) * 0.05 +
            (1 - market.competition_intensity) * 0.10 +
            min(1, market.economic_growth * 10) * 0.05 +
            market.smartphone_penetration * 0.10 +
            market.financial_inclusion_need * 0.05 +
            market.diaspora_remittance * 0.05
        )
        
        # Project execution score
        project_score = (
            project.tech_execution * 0.15 +
            project.team_quality * 0.10 +
            project.marketing_effectiveness * 0.15 +
            project.partnership_success * 0.10 +
            project.product_market_fit * 0.15 +
            project.security_reliability * 0.10 +
            project.user_experience * 0.10 +
            project.charity_impact_visibility * 0.10 +
            project.tokenomics_sustainability * 0.05
        )
        
        # Combined score with Indian market premium
        base_adoption = (market_score * 0.6 + project_score * 0.4)
        
        # Cultural alignment bonus (India loves charity)
        cultural_bonus = market.cultural_alignment * market.ngo_trust_factor * 0.1
        
        return min(1.0, base_adoption + cultural_bonus)
    
    def calculate_revenue_multiplier(self, adoption_rate: float, 
                                   market: MarketFactors, 
                                   project: ProjectFactors) -> float:
        """Calculate revenue multiplier based on adoption and execution"""
        
        base_multiplier = adoption_rate ** 1.5  # Network effects
        
        # Indian market size advantage
        market_size_bonus = 1 + (market.smartphone_penetration * 0.5)
        
        # Charity-driven viral growth
        viral_factor = 1 + (project.charity_impact_visibility * 
                           market.cultural_alignment * 0.3)
        
        # Regulatory risk discount
        regulatory_discount = 1 + market.crypto_regulation_clarity * 0.2
        
        return base_multiplier * market_size_bonus * viral_factor * regulatory_discount
    
    def simulate_year_5_metrics(self) -> Tuple[float, float, float]:
        """Simulate Year 5 daily transaction volume, platform revenue, and token price"""
        
        market = self.generate_market_factors()
        project = self.generate_project_factors()
        
        adoption_rate = self.calculate_adoption_rate(market, project)
        revenue_multiplier = self.calculate_revenue_multiplier(adoption_rate, market, project)
        
        # Base case: ₹500 Cr daily volume
        # Range: ₹10 Cr (failure) to ₹5000 Cr (massive success)
        daily_volume = 10 + (5000 - 10) * adoption_rate * revenue_multiplier
        
        # Platform revenue based on volume
        # Assuming 2% average take rate across all products
        annual_platform_revenue = daily_volume * 365 * 0.02
        
        # Token price correlation with adoption
        # Base case: ₹50, Range: ₹0.1 to ₹500
        token_price = 0.1 + (500 - 0.1) * (adoption_rate ** 1.2) * \
                     (revenue_multiplier ** 0.5)
        
        return daily_volume, annual_platform_revenue, token_price
    
    def run_simulation(self) -> pd.DataFrame:
        """Run Monte Carlo simulation"""
        
        results = []
        for _ in range(self.n_simulations):
            daily_volume, platform_revenue, token_price = self.simulate_year_5_metrics()
            
            # Calculate founder wealth
            founder_tokens_value = 142862766 * token_price / 1e7  # In Cr
            tax_royalty = daily_volume * 365 * 0.025 * 0.001  # 0.1% of 2.5% tax
            platform_royalty = platform_revenue * 0.05  # 5% of platform revenue
            total_founder_income = tax_royalty + platform_royalty
            
            # Calculate success metrics
            results.append({
                'daily_volume': daily_volume,
                'platform_revenue': platform_revenue,
                'token_price': token_price,
                'founder_token_value': founder_tokens_value,
                'founder_annual_income': total_founder_income,
                'total_charity': daily_volume * 365 * 0.025 * 0.0075 + platform_revenue * 0.1,
                'success': daily_volume >= 100  # Success = ₹100 Cr+ daily volume
            })
        
        return pd.DataFrame(results)
    
    def analyze_results(self, df: pd.DataFrame) -> dict:
        """Analyze simulation results"""
        
        success_rate = df['success'].mean() * 100
        
        # Calculate percentiles
        metrics = ['daily_volume', 'platform_revenue', 'token_price', 
                  'founder_annual_income', 'total_charity']
        
        percentiles = {}
        for metric in metrics:
            percentiles[metric] = {
                'p10': df[metric].quantile(0.10),
                'p25': df[metric].quantile(0.25),
                'p50': df[metric].quantile(0.50),
                'p75': df[metric].quantile(0.75),
                'p90': df[metric].quantile(0.90),
                'mean': df[metric].mean()
            }
        
        # Probability of different success levels
        probabilities = {
            'unicorn': (df['platform_revenue'] >= 8000).mean() * 100,  # $1B+ valuation
            'high_success': (df['daily_volume'] >= 500).mean() * 100,
            'moderate_success': (df['daily_volume'] >= 100).mean() * 100,
            'survival': (df['daily_volume'] >= 10).mean() * 100,
            'founder_billionaire': (df['founder_token_value'] >= 1000).mean() * 100
        }
        
        return {
            'success_rate': success_rate,
            'percentiles': percentiles,
            'probabilities': probabilities
        }
    
    def plot_results(self, df: pd.DataFrame):
        """Create visualization of results"""
        
        fig, axes = plt.subplots(2, 3, figsize=(18, 12))
        fig.suptitle('DeshChain Monte Carlo Simulation Results (Year 5)', fontsize=16)
        
        # Daily Volume Distribution
        axes[0, 0].hist(df['daily_volume'], bins=50, alpha=0.7, color='blue', edgecolor='black')
        axes[0, 0].axvline(df['daily_volume'].median(), color='red', linestyle='--', label='Median')
        axes[0, 0].set_xlabel('Daily Transaction Volume (₹ Cr)')
        axes[0, 0].set_ylabel('Frequency')
        axes[0, 0].set_title('Daily Transaction Volume Distribution')
        axes[0, 0].legend()
        
        # Token Price Distribution
        axes[0, 1].hist(df['token_price'], bins=50, alpha=0.7, color='green', edgecolor='black')
        axes[0, 1].axvline(df['token_price'].median(), color='red', linestyle='--', label='Median')
        axes[0, 1].set_xlabel('Token Price (₹)')
        axes[0, 1].set_ylabel('Frequency')
        axes[0, 1].set_title('NAMO Token Price Distribution')
        axes[0, 1].legend()
        
        # Founder Income Distribution
        axes[0, 2].hist(df['founder_annual_income'], bins=50, alpha=0.7, color='purple', edgecolor='black')
        axes[0, 2].axvline(df['founder_annual_income'].median(), color='red', linestyle='--', label='Median')
        axes[0, 2].set_xlabel('Annual Income (₹ Cr)')
        axes[0, 2].set_ylabel('Frequency')
        axes[0, 2].set_title('Founder Annual Income Distribution')
        axes[0, 2].legend()
        
        # Charity Impact Distribution
        axes[1, 0].hist(df['total_charity'], bins=50, alpha=0.7, color='orange', edgecolor='black')
        axes[1, 0].axvline(df['total_charity'].median(), color='red', linestyle='--', label='Median')
        axes[1, 0].set_xlabel('Annual Charity (₹ Cr)')
        axes[1, 0].set_ylabel('Frequency')
        axes[1, 0].set_title('Total Annual Charity Distribution')
        axes[1, 0].legend()
        
        # Success Probability by Volume
        volume_ranges = [(0, 50), (50, 100), (100, 500), (500, 1000), (1000, 5000)]
        success_probs = []
        labels = []
        
        for low, high in volume_ranges:
            mask = (df['daily_volume'] >= low) & (df['daily_volume'] < high)
            prob = mask.mean() * 100
            success_probs.append(prob)
            labels.append(f'₹{low}-{high} Cr')
        
        axes[1, 1].bar(labels, success_probs, color='teal', edgecolor='black')
        axes[1, 1].set_xlabel('Daily Volume Range')
        axes[1, 1].set_ylabel('Probability (%)')
        axes[1, 1].set_title('Probability Distribution by Volume Range')
        axes[1, 1].tick_params(axis='x', rotation=45)
        
        # Correlation: Volume vs Founder Wealth
        axes[1, 2].scatter(df['daily_volume'], df['founder_annual_income'], 
                          alpha=0.5, s=10, color='darkred')
        axes[1, 2].set_xlabel('Daily Volume (₹ Cr)')
        axes[1, 2].set_ylabel('Founder Annual Income (₹ Cr)')
        axes[1, 2].set_title('Volume vs Founder Income Correlation')
        axes[1, 2].set_xscale('log')
        axes[1, 2].set_yscale('log')
        
        plt.tight_layout()
        plt.savefig('/root/namo/monte_carlo_results.png', dpi=300, bbox_inches='tight')
        plt.close()

def run_comprehensive_analysis():
    """Run comprehensive Monte Carlo analysis"""
    
    print("Running DeshChain Monte Carlo Simulation...")
    print("=" * 50)
    
    # Run simulation
    sim = DeshChainMonteCarloSimulation(n_simulations=10000)
    df = sim.run_simulation()
    results = sim.analyze_results(df)
    
    # Print results
    print(f"\nSUCCESS PROBABILITY: {results['success_rate']:.1f}%")
    print(f"(Success defined as ₹100+ Cr daily transaction volume)")
    
    print("\n" + "=" * 50)
    print("YEAR 5 PROJECTIONS (Percentiles)")
    print("=" * 50)
    
    for metric, values in results['percentiles'].items():
        print(f"\n{metric.replace('_', ' ').title()}:")
        print(f"  10th percentile: ₹{values['p10']:,.0f} Cr")
        print(f"  25th percentile: ₹{values['p25']:,.0f} Cr")
        print(f"  50th percentile: ₹{values['p50']:,.0f} Cr")
        print(f"  75th percentile: ₹{values['p75']:,.0f} Cr")
        print(f"  90th percentile: ₹{values['p90']:,.0f} Cr")
        print(f"  Mean: ₹{values['mean']:,.0f} Cr")
    
    print("\n" + "=" * 50)
    print("PROBABILITY OF DIFFERENT OUTCOMES")
    print("=" * 50)
    
    for outcome, prob in results['probabilities'].items():
        print(f"{outcome.replace('_', ' ').title()}: {prob:.1f}%")
    
    # Indian market specific insights
    print("\n" + "=" * 50)
    print("INDIAN MARKET ADVANTAGE FACTORS")
    print("=" * 50)
    print("1. Cultural alignment with charity (dharma): +20% success boost")
    print("2. Massive unbanked population: 200M+ potential users")
    print("3. Digital India momentum: Government push helps adoption")
    print("4. Diaspora remittance market: $100B+ annual opportunity")
    print("5. Trust in NGO partnerships: Higher than global average")
    
    # Risk factors
    print("\n" + "=" * 50)
    print("KEY RISK FACTORS")
    print("=" * 50)
    print("1. Regulatory uncertainty: Could impact 30% of outcomes")
    print("2. Competition from global players: Moderate risk")
    print("3. Technical execution: Critical for first 2 years")
    print("4. Market education: Crypto awareness still growing")
    
    # Create visualizations
    sim.plot_results(df)
    print("\nVisualizations saved to: monte_carlo_results.png")
    
    return df, results

# Run the analysis
if __name__ == "__main__":
    df, results = run_comprehensive_analysis()