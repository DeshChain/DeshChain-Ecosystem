// DeshChain Validators Page JavaScript

// Validator slot data with pricing information
const validatorSlots = {
    1: { contract: 100000, stake: 200000, total: 300000, tier: 1, lock: 6, vesting: 18, performanceBond: 0.20 },
    2: { contract: 120000, stake: 220000, total: 340000, tier: 1, lock: 6, vesting: 18, performanceBond: 0.20 },
    3: { contract: 140000, stake: 240000, total: 380000, tier: 1, lock: 6, vesting: 18, performanceBond: 0.20 },
    4: { contract: 160000, stake: 260000, total: 420000, tier: 1, lock: 6, vesting: 18, performanceBond: 0.20 },
    5: { contract: 180000, stake: 280000, total: 460000, tier: 1, lock: 6, vesting: 18, performanceBond: 0.20 },
    6: { contract: 200000, stake: 300000, total: 500000, tier: 1, lock: 6, vesting: 18, performanceBond: 0.20 },
    7: { contract: 220000, stake: 320000, total: 540000, tier: 1, lock: 6, vesting: 18, performanceBond: 0.20 },
    8: { contract: 240000, stake: 340000, total: 580000, tier: 1, lock: 6, vesting: 18, performanceBond: 0.20 },
    9: { contract: 260000, stake: 360000, total: 620000, tier: 1, lock: 6, vesting: 18, performanceBond: 0.20 },
    10: { contract: 280000, stake: 380000, total: 660000, tier: 1, lock: 6, vesting: 18, performanceBond: 0.20 },
    // Price doubles here
    11: { contract: 400000, stake: 800000, total: 1200000, tier: 2, lock: 9, vesting: 24, performanceBond: 0.25 },
    12: { contract: 420000, stake: 820000, total: 1240000, tier: 2, lock: 9, vesting: 24, performanceBond: 0.25 },
    13: { contract: 440000, stake: 840000, total: 1280000, tier: 2, lock: 9, vesting: 24, performanceBond: 0.25 },
    14: { contract: 460000, stake: 860000, total: 1320000, tier: 2, lock: 9, vesting: 24, performanceBond: 0.25 },
    15: { contract: 480000, stake: 880000, total: 1360000, tier: 2, lock: 9, vesting: 24, performanceBond: 0.25 },
    16: { contract: 500000, stake: 900000, total: 1400000, tier: 2, lock: 9, vesting: 24, performanceBond: 0.25 },
    17: { contract: 520000, stake: 920000, total: 1440000, tier: 2, lock: 9, vesting: 24, performanceBond: 0.25 },
    18: { contract: 540000, stake: 940000, total: 1480000, tier: 2, lock: 9, vesting: 24, performanceBond: 0.25 },
    19: { contract: 560000, stake: 960000, total: 1520000, tier: 2, lock: 9, vesting: 24, performanceBond: 0.25 },
    20: { contract: 580000, stake: 980000, total: 1560000, tier: 2, lock: 9, vesting: 24, performanceBond: 0.25 },
    // Premium finale
    21: { contract: 650000, stake: 1500000, total: 2150000, tier: 3, lock: 12, vesting: 36, performanceBond: 0.30, special: 'Param Rakshak NFT with 2x governance weight' }
};

// Total raise information
const totalRaise = {
    total: 20050000, // $20.05M
    contractFunds: 8250000, // $8.25M (Operations)
    stakingFunds: 11800000 // $11.8M (Locked)
};

// Revenue share percentages
const revenueShares = {
    transactionTax: { platformFee: 0.025, validatorShare: 0.5 }, // 2.5% fee, 50% to validators
    dexTrading: { platformFee: 0.0025, validatorShare: 0.25 }, // 0.25% fee, 25% to validators
    mevCapture: { validatorShare: 0.6 }, // 60% to validators
    launchpad: { platformFee: 0.05, validatorShare: 0.2 }, // 5% fee, 20% to validators
    blockRewards: { minAPY: 0.12, maxAPY: 0.18 } // 12-18% APY
};

// Calculate validator stake weight
function calculateStakeWeight(slotNumber) {
    const slot = validatorSlots[slotNumber];
    const totalStake = Object.values(validatorSlots).reduce((sum, v) => sum + v.stake, 0);
    return slot.stake / totalStake;
}

// Calculate ROI
function calculateROI() {
    const slotNumber = parseInt(document.getElementById('slot-selector').value);
    if (!slotNumber) {
        alert('Please select a validator slot');
        return;
    }

    const slot = validatorSlots[slotNumber];
    const dailyVolume = parseFloat(document.getElementById('daily-volume').value) || 0;
    const dexVolume = parseFloat(document.getElementById('dex-volume').value) || 0;
    const mevRevenue = parseFloat(document.getElementById('mev-revenue').value) || 0;
    const launchpadMonthly = parseFloat(document.getElementById('launchpad-monthly').value) || 0;
    const isIndiaBased = document.getElementById('india-based').checked;

    // Calculate stake weight
    const stakeWeight = calculateStakeWeight(slotNumber);
    const equalShare = 1 / 21; // Equal share among 21 validators

    // Calculate daily revenues
    const dailyTransactionFees = dailyVolume * revenueShares.transactionTax.platformFee * 
                                revenueShares.transactionTax.validatorShare * stakeWeight;
    
    const dailyDexFees = dexVolume * revenueShares.dexTrading.platformFee * 
                         revenueShares.dexTrading.validatorShare * stakeWeight;
    
    const dailyMEV = mevRevenue * revenueShares.mevCapture.validatorShare * stakeWeight;
    
    const dailyLaunchpadFees = (launchpadMonthly / 30) * revenueShares.launchpad.platformFee * 
                               revenueShares.launchpad.validatorShare * equalShare;

    // Calculate block rewards (average 15% APY)
    const avgBlockRewardAPY = (revenueShares.blockRewards.minAPY + revenueShares.blockRewards.maxAPY) / 2;
    const dailyBlockRewards = (slot.stake * avgBlockRewardAPY) / 365;

    // Apply India bonus if applicable
    const indiaBonus = isIndiaBased ? 0.2 : 0;
    
    // Calculate total daily revenue
    let totalDailyRevenue = dailyTransactionFees + dailyDexFees + dailyMEV + dailyLaunchpadFees + dailyBlockRewards;
    totalDailyRevenue *= (1 + indiaBonus);

    // Calculate monthly and annual revenues
    const monthlyRevenue = totalDailyRevenue * 30;
    const annualRevenue = totalDailyRevenue * 365;

    // Calculate ROI metrics
    const totalInvestment = slot.total;
    const contractAmount = slot.contract; // Non-refundable
    const annualROI = (annualRevenue / totalInvestment) * 100;
    const monthlyROI = (monthlyRevenue / totalInvestment) * 100;
    const paybackPeriod = totalInvestment / annualRevenue;

    // Display results
    displayResults({
        slotNumber,
        slot,
        dailyRevenue: totalDailyRevenue,
        monthlyRevenue,
        annualRevenue,
        annualROI,
        monthlyROI,
        paybackPeriod,
        contractAmount,
        breakdown: {
            transactionFees: dailyTransactionFees * 365,
            dexFees: dailyDexFees * 365,
            mevCapture: dailyMEV * 365,
            launchpadFees: dailyLaunchpadFees * 365,
            blockRewards: dailyBlockRewards * 365
        },
        indiaBonus,
        stakeWeight
    });
}

// Display calculation results
function displayResults(results) {
    const resultsContainer = document.getElementById('calculator-results');
    
    const html = `
        <h3>Projected Returns for Slot ${results.slotNumber}</h3>
        <div class="results-content">
            <!-- Risk Disclosure -->
            <div class="risk-disclosure">
                <h4>⚠️ Important Risk Disclosure</h4>
                <ul>
                    <li>Contract amount of $${results.contractAmount.toLocaleString()} is non-refundable</li>
                    <li>Returns are projections based on estimated platform metrics</li>
                    <li>Actual returns may vary significantly from projections</li>
                    <li>Cryptocurrency investments carry high risk of loss</li>
                </ul>
                <p>Past performance does not guarantee future results</p>
            </div>
            
            <!-- Investment Summary -->
            <div class="result-item">
                <span class="result-label">Total Investment:</span>
                <span class="result-value">$${results.slot.total.toLocaleString()}</span>
            </div>
            <div class="result-item">
                <span class="result-label">Contract (Non-refundable):</span>
                <span class="result-value">$${results.contractAmount.toLocaleString()}</span>
            </div>
            <div class="result-item">
                <span class="result-label">Staking Amount:</span>
                <span class="result-value">$${results.slot.stake.toLocaleString()}</span>
            </div>
            <div class="result-item">
                <span class="result-label">Your Stake Weight:</span>
                <span class="result-value">${(results.stakeWeight * 100).toFixed(2)}%</span>
            </div>
            
            <!-- Revenue Projections -->
            <h4 style="margin-top: 2rem;">Revenue Projections</h4>
            <div class="result-item">
                <span class="result-label">Daily Revenue:</span>
                <span class="result-value positive">$${results.dailyRevenue.toLocaleString('en-US', {maximumFractionDigits: 0})}</span>
            </div>
            <div class="result-item">
                <span class="result-label">Monthly Revenue:</span>
                <span class="result-value positive">$${results.monthlyRevenue.toLocaleString('en-US', {maximumFractionDigits: 0})}</span>
            </div>
            <div class="result-item">
                <span class="result-label">Annual Revenue:</span>
                <span class="result-value positive highlight">$${results.annualRevenue.toLocaleString('en-US', {maximumFractionDigits: 0})}</span>
            </div>
            
            <!-- Revenue Breakdown -->
            <h4 style="margin-top: 2rem;">Annual Revenue Breakdown</h4>
            <div class="result-item">
                <span class="result-label">Transaction Fees:</span>
                <span class="result-value">$${results.breakdown.transactionFees.toLocaleString('en-US', {maximumFractionDigits: 0})}</span>
            </div>
            <div class="result-item">
                <span class="result-label">DEX Trading Fees:</span>
                <span class="result-value">$${results.breakdown.dexFees.toLocaleString('en-US', {maximumFractionDigits: 0})}</span>
            </div>
            <div class="result-item">
                <span class="result-label">MEV Capture:</span>
                <span class="result-value">$${results.breakdown.mevCapture.toLocaleString('en-US', {maximumFractionDigits: 0})}</span>
            </div>
            <div class="result-item">
                <span class="result-label">Launchpad Fees:</span>
                <span class="result-value">$${results.breakdown.launchpadFees.toLocaleString('en-US', {maximumFractionDigits: 0})}</span>
            </div>
            <div class="result-item">
                <span class="result-label">Block Rewards:</span>
                <span class="result-value">$${results.breakdown.blockRewards.toLocaleString('en-US', {maximumFractionDigits: 0})}</span>
            </div>
            ${results.indiaBonus > 0 ? `
            <div class="result-item">
                <span class="result-label">India Location Bonus (20%):</span>
                <span class="result-value positive">Included</span>
            </div>
            ` : ''}
            
            <!-- ROI Summary -->
            <div class="roi-summary">
                <h4>Return on Investment</h4>
                <div class="roi-percentage">${results.annualROI.toFixed(1)}%</div>
                <div class="payback-period">Estimated payback period: ${results.paybackPeriod.toFixed(1)} years</div>
            </div>
            
            <!-- Additional Benefits -->
            <h4 style="margin-top: 2rem;">Additional Genesis Benefits</h4>
            <ul style="margin-left: 1.5rem; color: var(--text-secondary);">
                <li>Guaranteed 1% minimum pool share when >100 validators</li>
                <li>Tradeable Bharat Guardian NFT with 5% royalty</li>
                <li>Enhanced governance voting power (up to 3x)</li>
                <li>10-20% commission on referred validators</li>
                ${results.slotNumber === 21 ? '<li><strong>Exclusive Param Rakshak NFT with 2x governance weight</strong></li>' : ''}
            </ul>
            
            <!-- Lock and Vesting Info -->
            <h4 style="margin-top: 2rem;">Lock & Vesting Schedule</h4>
            <div class="result-item">
                <span class="result-label">Lock Period:</span>
                <span class="result-value">${results.slot.lock} months</span>
            </div>
            <div class="result-item">
                <span class="result-label">Vesting Period:</span>
                <span class="result-value">${results.slot.vesting} months</span>
            </div>
            <div class="result-item">
                <span class="result-label">Performance Bond:</span>
                <span class="result-value">$${(results.slot.stake * results.slot.performanceBond).toLocaleString()}</span>
            </div>
        </div>
    `;
    
    resultsContainer.innerHTML = html;
    
    // Scroll to results
    resultsContainer.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
}

// Add fund summary section
function addFundSummary() {
    // Check if summary already exists
    if (document.getElementById('fund-summary')) return;
    
    const auctionSection = document.querySelector('.auction-structure .container');
    const summaryHTML = `
        <div id="fund-summary" class="fund-summary" style="margin-top: 3rem; padding: 2rem; background: linear-gradient(135deg, rgba(255, 153, 51, 0.1), rgba(247, 37, 133, 0.1)); border-radius: var(--radius-xl);">
            <h3 style="text-align: center; margin-bottom: 2rem;">Total Validator Program Summary</h3>
            <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 2rem; text-align: center;">
                <div>
                    <span style="display: block; font-size: 2rem; font-weight: 700; color: var(--primary);">$${(totalRaise.total / 1000000).toFixed(2)}M</span>
                    <span style="color: var(--text-secondary);">Total Raise</span>
                </div>
                <div>
                    <span style="display: block; font-size: 2rem; font-weight: 700; color: var(--primary);">$${(totalRaise.contractFunds / 1000000).toFixed(2)}M</span>
                    <span style="color: var(--text-secondary);">Contract Funds (Operations)</span>
                </div>
                <div>
                    <span style="display: block; font-size: 2rem; font-weight: 700; color: var(--primary);">$${(totalRaise.stakingFunds / 1000000).toFixed(2)}M</span>
                    <span style="color: var(--text-secondary);">Staking Funds (Locked)</span>
                </div>
            </div>
            <div class="risk-disclosure" style="margin-top: 2rem;">
                <h4>⚠️ Important Notice</h4>
                <p>Contract amounts are non-refundable and will be used for DeshChain operations and development. Only staking funds are locked on-chain and subject to vesting schedules.</p>
            </div>
        </div>
    `;
    
    auctionSection.insertAdjacentHTML('beforeend', summaryHTML);
}

// Initialize page
document.addEventListener('DOMContentLoaded', function() {
    // Add fund summary
    addFundSummary();
    
    // Set up event listeners
    const calculateBtn = document.querySelector('.calculate-btn');
    if (calculateBtn) {
        calculateBtn.addEventListener('click', calculateROI);
    }
    
    // Add enter key support for inputs
    const inputs = document.querySelectorAll('input[type="number"]');
    inputs.forEach(input => {
        input.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                calculateROI();
            }
        });
    });
    
    // Add slot selector change handler
    const slotSelector = document.getElementById('slot-selector');
    if (slotSelector) {
        slotSelector.addEventListener('change', function() {
            // Clear results when slot changes
            const resultsContainer = document.getElementById('calculator-results');
            if (resultsContainer) {
                resultsContainer.innerHTML = `
                    <h3>Projected Returns</h3>
                    <div class="results-placeholder">
                        <p>Click calculate to see projected returns for Slot ${this.value}</p>
                    </div>
                `;
            }
        });
    }
});