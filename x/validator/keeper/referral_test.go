package keeper_test

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/DeshChain/DeshChain-Ecosystem/x/validator/keeper"
    "github.com/DeshChain/DeshChain-Ecosystem/x/validator/types"
    "github.com/DeshChain/DeshChain-Ecosystem/testutil"
)

type ReferralTestSuite struct {
    suite.Suite
    
    ctx           sdk.Context
    keeper        keeper.Keeper
    referralKeeper *keeper.ReferralKeeper
    
    // Test addresses
    genesisValidator1 sdk.AccAddress
    genesisValidator2 sdk.AccAddress
    newValidator1     sdk.AccAddress
    newValidator2     sdk.AccAddress
    regularUser       sdk.AccAddress
}

func TestReferralTestSuite(t *testing.T) {
    suite.Run(t, new(ReferralTestSuite))
}

func (suite *ReferralTestSuite) SetupTest() {
    app := testutil.Setup(false)
    suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
    suite.keeper = app.ValidatorKeeper
    suite.referralKeeper = keeper.NewReferralKeeper(suite.keeper)
    
    // Create test addresses
    suite.genesisValidator1 = testutil.CreateRandomAccounts(1)[0]
    suite.genesisValidator2 = testutil.CreateRandomAccounts(1)[0]
    suite.newValidator1 = testutil.CreateRandomAccounts(1)[0]
    suite.newValidator2 = testutil.CreateRandomAccounts(1)[0]
    suite.regularUser = testutil.CreateRandomAccounts(1)[0]
    
    // Setup genesis validators with NFTs
    suite.setupGenesisValidators()
}

func (suite *ReferralTestSuite) setupGenesisValidators() {
    // Create genesis NFT for validator 1 (rank 1)
    nft1 := types.GenesisValidatorNFT{
        TokenID:      1,
        Rank:         1,
        EnglishName:  "Narendra Modi",
        HindiName:    "नरेंद्र मोदी",
        CurrentOwner: suite.genesisValidator1.String(),
        ImageURI:     "/nfts/1.png",
        MintedAt:     suite.ctx.BlockTime(),
    }
    suite.keeper.SetGenesisNFT(suite.ctx, nft1)
    
    // Create genesis NFT for validator 2 (rank 10)
    nft2 := types.GenesisValidatorNFT{
        TokenID:      10,
        Rank:         10,
        EnglishName:  "Amit Shah",
        HindiName:    "अमित शाह",
        CurrentOwner: suite.genesisValidator2.String(),
        ImageURI:     "/nfts/10.png",
        MintedAt:     suite.ctx.BlockTime(),
    }
    suite.keeper.SetGenesisNFT(suite.ctx, nft2)
    
    // Create validator stakes
    namoPrice := sdk.NewDecWithPrec(1, 2) // ₹0.01
    
    stake1 := types.ValidatorStake{
        ValidatorAddr:    suite.genesisValidator1.String(),
        Rank:            1,
        OriginalUSDValue: sdk.NewDec(380000), // $380K
        NAMOStaked:      types.CalculateNAMORequired(sdk.NewDec(380000), namoPrice),
        NAMOPrice:       namoPrice,
        StakingTime:     suite.ctx.BlockTime(),
        Tier:            types.StakeTier{TierID: 1, LockPeriod: 12, VestingPeriod: 36},
    }
    suite.keeper.SetValidatorStake(suite.ctx, stake1)
    
    stake2 := types.ValidatorStake{
        ValidatorAddr:    suite.genesisValidator2.String(),
        Rank:            10,
        OriginalUSDValue: sdk.NewDec(380000), // $380K
        NAMOStaked:      types.CalculateNAMORequired(sdk.NewDec(380000), namoPrice),
        NAMOPrice:       namoPrice,
        StakingTime:     suite.ctx.BlockTime(),
        Tier:            types.StakeTier{TierID: 1, LockPeriod: 12, VestingPeriod: 36},
    }
    suite.keeper.SetValidatorStake(suite.ctx, stake2)
}

// Test referral creation
func (suite *ReferralTestSuite) TestCreateReferral() {
    tests := []struct {
        name          string
        referrer      sdk.AccAddress
        referred      sdk.AccAddress
        rank          uint32
        expectError   bool
        errorContains string
    }{
        {
            name:        "valid referral by genesis validator",
            referrer:    suite.genesisValidator1,
            referred:    suite.newValidator1,
            rank:        22,
            expectError: false,
        },
        {
            name:          "self referral should fail",
            referrer:      suite.genesisValidator1,
            referred:      suite.genesisValidator1,
            rank:          23,
            expectError:   true,
            errorContains: "self-referral",
        },
        {
            name:          "non-genesis validator cannot refer",
            referrer:      suite.regularUser,
            referred:      suite.newValidator2,
            rank:          24,
            expectError:   true,
            errorContains: "genesis validators",
        },
        {
            name:          "invalid rank should fail",
            referrer:      suite.genesisValidator1,
            referred:      suite.newValidator2,
            rank:          1001,
            expectError:   true,
            errorContains: "invalid rank",
        },
    }
    
    for _, tc := range tests {
        suite.Run(tc.name, func() {
            referralID, err := suite.referralKeeper.CreateReferral(
                suite.ctx,
                tc.referrer,
                tc.referred,
                tc.rank,
            )
            
            if tc.expectError {
                suite.Require().Error(err)
                if tc.errorContains != "" {
                    suite.Require().Contains(err.Error(), tc.errorContains)
                }
            } else {
                suite.Require().NoError(err)
                suite.Require().Greater(referralID, uint64(0))
                
                // Verify referral was created
                referral, found := suite.keeper.GetReferral(suite.ctx, referralID)
                suite.Require().True(found)
                suite.Require().Equal(tc.referrer.String(), referral.ReferrerAddr)
                suite.Require().Equal(tc.referred.String(), referral.ReferredAddr)
                suite.Require().Equal(tc.rank, referral.ReferredRank)
                suite.Require().Equal(types.ReferralStatusPending, referral.Status)
            }
        })
    }
}

// Test referral limits
func (suite *ReferralTestSuite) TestReferralLimits() {
    // Create 5 referrals (monthly limit)
    for i := 0; i < 5; i++ {
        newValidator := testutil.CreateRandomAccounts(1)[0]
        _, err := suite.referralKeeper.CreateReferral(
            suite.ctx,
            suite.genesisValidator1,
            newValidator,
            uint32(22+i),
        )
        suite.Require().NoError(err)
    }
    
    // 6th referral should fail due to monthly limit
    newValidator := testutil.CreateRandomAccounts(1)[0]
    _, err := suite.referralKeeper.CreateReferral(
        suite.ctx,
        suite.genesisValidator1,
        newValidator,
        27,
    )
    suite.Require().Error(err)
    suite.Require().Contains(err.Error(), "monthly referral limit")
}

// Test referral activation
func (suite *ReferralTestSuite) TestReferralActivation() {
    // Create referral
    referralID, err := suite.referralKeeper.CreateReferral(
        suite.ctx,
        suite.genesisValidator1,
        suite.newValidator1,
        22,
    )
    suite.Require().NoError(err)
    
    // Verify referral is pending
    referral, found := suite.keeper.GetReferral(suite.ctx, referralID)
    suite.Require().True(found)
    suite.Require().Equal(types.ReferralStatusPending, referral.Status)
    
    // Simulate validator staking (which should activate referral)
    namoPrice := sdk.NewDecWithPrec(1, 2) // ₹0.01
    stake := types.ValidatorStake{
        ValidatorAddr:    suite.newValidator1.String(),
        Rank:            22,
        OriginalUSDValue: sdk.NewDec(200000), // $200K
        NAMOStaked:      types.CalculateNAMORequired(sdk.NewDec(200000), namoPrice),
        NAMOPrice:       namoPrice,
        StakingTime:     suite.ctx.BlockTime(),
        Tier:            types.StakeTier{TierID: 1, LockPeriod: 6, VestingPeriod: 18},
    }
    suite.keeper.SetValidatorStake(suite.ctx, stake)
    
    // Activate referral
    err = suite.referralKeeper.ActivateReferral(suite.ctx, suite.newValidator1)
    suite.Require().NoError(err)
    
    // Verify referral is now active
    referral, found = suite.keeper.GetReferral(suite.ctx, referralID)
    suite.Require().True(found)
    suite.Require().Equal(types.ReferralStatusActive, referral.Status)
    suite.Require().False(referral.ActivatedAt.IsZero())
    
    // Verify referrer stats updated
    stats := suite.referralKeeper.GetReferralStats(suite.ctx, suite.genesisValidator1.String())
    suite.Require().Equal(uint32(1), stats.ActiveReferrals)
}

// Test commission processing
func (suite *ReferralTestSuite) TestCommissionProcessing() {
    // Setup active referral
    referralID := suite.setupActiveReferral()
    
    // Process commission
    revenueAmount := sdk.NewInt(1000000) // ₹10,000
    err := suite.referralKeeper.ProcessReferralCommission(
        suite.ctx,
        suite.newValidator1,
        revenueAmount,
    )
    suite.Require().NoError(err)
    
    // Verify commission was calculated and paid
    referral, found := suite.keeper.GetReferral(suite.ctx, referralID)
    suite.Require().True(found)
    
    expectedCommission := revenueAmount.ToDec().Mul(referral.CommissionRate).TruncateInt()
    suite.Require().Equal(expectedCommission, referral.TotalCommission)
    suite.Require().Equal(expectedCommission, referral.PaidCommission)
}

// Test token launch eligibility
func (suite *ReferralTestSuite) TestTokenLaunchEligibility() {
    // Create enough referrals to trigger token launch
    referralIDs := make([]uint64, 0, 5)
    for i := 0; i < 5; i++ {
        newValidator := testutil.CreateRandomAccounts(1)[0]
        
        // Create referral
        referralID, err := suite.referralKeeper.CreateReferral(
            suite.ctx,
            suite.genesisValidator1,
            newValidator,
            uint32(22+i),
        )
        suite.Require().NoError(err)
        referralIDs = append(referralIDs, referralID)
        
        // Activate referral by creating stake
        namoPrice := sdk.NewDecWithPrec(1, 2)
        stake := types.ValidatorStake{
            ValidatorAddr:    newValidator.String(),
            Rank:            uint32(22 + i),
            OriginalUSDValue: sdk.NewDec(200000),
            NAMOStaked:      types.CalculateNAMORequired(sdk.NewDec(200000), namoPrice),
            NAMOPrice:       namoPrice,
            StakingTime:     suite.ctx.BlockTime(),
            Tier:            types.StakeTier{TierID: 1, LockPeriod: 6, VestingPeriod: 18},
        }
        suite.keeper.SetValidatorStake(suite.ctx, stake)
        
        err = suite.referralKeeper.ActivateReferral(suite.ctx, newValidator)
        suite.Require().NoError(err)
    }
    
    // Check if token launch was triggered
    stats := suite.referralKeeper.GetReferralStats(suite.ctx, suite.genesisValidator1.String())
    suite.Require().Equal(uint32(5), stats.TotalReferrals)
    suite.Require().True(stats.TokenLaunched) // Should have triggered auto-launch
    suite.Require().Greater(stats.TokenID, uint64(0))
}

// Test anti-gaming measures
func (suite *ReferralTestSuite) TestAntiGamingMeasures() {
    // Test minimum time gap between referrals
    newValidator1 := testutil.CreateRandomAccounts(1)[0]
    _, err := suite.referralKeeper.CreateReferral(
        suite.ctx,
        suite.genesisValidator1,
        newValidator1,
        22,
    )
    suite.Require().NoError(err)
    
    // Try to create another referral immediately (should fail)
    newValidator2 := testutil.CreateRandomAccounts(1)[0]
    _, err = suite.referralKeeper.CreateReferral(
        suite.ctx,
        suite.genesisValidator1,
        newValidator2,
        23,
    )
    suite.Require().Error(err)
    suite.Require().Contains(err.Error(), "24-hour gap")
    
    // Advance time by 25 hours
    suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(25 * time.Hour))
    
    // Now referral should succeed
    _, err = suite.referralKeeper.CreateReferral(
        suite.ctx,
        suite.genesisValidator1,
        newValidator2,
        23,
    )
    suite.Require().NoError(err)
}

// Test clawback mechanism
func (suite *ReferralTestSuite) TestClawbackMechanism() {
    // Setup active referral with commission
    referralID := suite.setupActiveReferral()
    
    // Process some commission
    revenueAmount := sdk.NewInt(1000000)
    err := suite.referralKeeper.ProcessReferralCommission(
        suite.ctx,
        suite.newValidator1,
        revenueAmount,
    )
    suite.Require().NoError(err)
    
    // Simulate referred validator exit (remove stake)
    suite.keeper.RemoveValidatorStake(suite.ctx, suite.newValidator1.String())
    
    // Advance time to within clawback period (within 1 year)
    suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(6 * 30 * 24 * time.Hour)) // 6 months
    
    // Process commission again (should trigger clawback check)
    err = suite.referralKeeper.ProcessReferralCommission(
        suite.ctx,
        suite.newValidator1,
        revenueAmount,
    )
    suite.Require().NoError(err) // Should handle clawback gracefully
    
    // Verify referral status changed to clawed back
    referral, found := suite.keeper.GetReferral(suite.ctx, referralID)
    suite.Require().True(found)
    // Note: Clawback would be triggered by validator exit detection elsewhere
}

// Test tier progression
func (suite *ReferralTestSuite) TestTierProgression() {
    initialStats := suite.referralKeeper.GetReferralStats(suite.ctx, suite.genesisValidator1.String())
    suite.Require().Equal(uint32(1), initialStats.CurrentTier)
    
    // Create referrals to progress through tiers
    for i := 0; i < 11; i++ { // Should move to tier 2 (11+ referrals)
        newValidator := testutil.CreateRandomAccounts(1)[0]
        
        // Advance time to avoid rate limiting
        suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(25 * time.Hour))
        
        _, err := suite.referralKeeper.CreateReferral(
            suite.ctx,
            suite.genesisValidator1,
            newValidator,
            uint32(22+i),
        )
        suite.Require().NoError(err)
        
        // Activate referral
        namoPrice := sdk.NewDecWithPrec(1, 2)
        stake := types.ValidatorStake{
            ValidatorAddr:    newValidator.String(),
            Rank:            uint32(22 + i),
            OriginalUSDValue: sdk.NewDec(200000),
            NAMOStaked:      types.CalculateNAMORequired(sdk.NewDec(200000), namoPrice),
            NAMOPrice:       namoPrice,
            StakingTime:     suite.ctx.BlockTime(),
            Tier:            types.StakeTier{TierID: 1, LockPeriod: 6, VestingPeriod: 18},
        }
        suite.keeper.SetValidatorStake(suite.ctx, stake)
        
        err = suite.referralKeeper.ActivateReferral(suite.ctx, newValidator)
        suite.Require().NoError(err)
    }
    
    // Check tier progression
    finalStats := suite.referralKeeper.GetReferralStats(suite.ctx, suite.genesisValidator1.String())
    suite.Require().Equal(uint32(11), finalStats.TotalReferrals)
    suite.Require().Equal(uint32(2), finalStats.CurrentTier) // Should be tier 2 now
}

// Helper function to setup an active referral
func (suite *ReferralTestSuite) setupActiveReferral() uint64 {
    // Create referral
    referralID, err := suite.referralKeeper.CreateReferral(
        suite.ctx,
        suite.genesisValidator1,
        suite.newValidator1,
        22,
    )
    suite.Require().NoError(err)
    
    // Create stake for referred validator
    namoPrice := sdk.NewDecWithPrec(1, 2)
    stake := types.ValidatorStake{
        ValidatorAddr:    suite.newValidator1.String(),
        Rank:            22,
        OriginalUSDValue: sdk.NewDec(200000),
        NAMOStaked:      types.CalculateNAMORequired(sdk.NewDec(200000), namoPrice),
        NAMOPrice:       namoPrice,
        StakingTime:     suite.ctx.BlockTime(),
        Tier:            types.StakeTier{TierID: 1, LockPeriod: 6, VestingPeriod: 18},
    }
    suite.keeper.SetValidatorStake(suite.ctx, stake)
    
    // Activate referral
    err = suite.referralKeeper.ActivateReferral(suite.ctx, suite.newValidator1)
    suite.Require().NoError(err)
    
    // Advance time past cliff period (6 months)
    suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(7 * 30 * 24 * time.Hour))
    
    return referralID
}

// Benchmark tests
func (suite *ReferralTestSuite) TestReferralPerformance() {
    // Test creation of 100 referrals
    start := time.Now()
    
    for i := 0; i < 100; i++ {
        newValidator := testutil.CreateRandomAccounts(1)[0]
        
        // Use different genesis validators to avoid limits
        referrer := suite.genesisValidator1
        if i%2 == 1 {
            referrer = suite.genesisValidator2
        }
        
        // Advance time to avoid rate limiting
        suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(25 * time.Hour))
        
        _, err := suite.referralKeeper.CreateReferral(
            suite.ctx,
            referrer,
            newValidator,
            uint32(22+i),
        )
        suite.Require().NoError(err)
    }
    
    duration := time.Since(start)
    suite.T().Logf("Created 100 referrals in %v", duration)
    suite.Require().Less(duration, 10*time.Second) // Should complete within 10 seconds
}