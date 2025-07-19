/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VillagePool represents a community-managed liquidity pool
// Inspired by village self-help groups and cooperative societies
type VillagePool struct {
	// Pool identification
	PoolId            uint64         `json:"pool_id" yaml:"pool_id"`
	VillageName       string         `json:"village_name" yaml:"village_name"`
	PostalCode        string         `json:"postal_code" yaml:"postal_code"`
	StateCode         string         `json:"state_code" yaml:"state_code"`
	DistrictCode      string         `json:"district_code" yaml:"district_code"`
	
	// Community details
	PanchayatName     string         `json:"panchayat_name" yaml:"panchayat_name"`
	PanchayatHead     sdk.AccAddress `json:"panchayat_head" yaml:"panchayat_head"`
	CooperativeId     string         `json:"cooperative_id" yaml:"cooperative_id"`
	
	// Pool configuration
	SupportedTokens   []string       `json:"supported_tokens" yaml:"supported_tokens"`
	PrimaryToken      string         `json:"primary_token" yaml:"primary_token"` // Usually NAMO
	LocalCurrency     string         `json:"local_currency" yaml:"local_currency"` // INR
	
	// Liquidity management
	TotalLiquidity    sdk.Coins      `json:"total_liquidity" yaml:"total_liquidity"`
	AvailableLiquidity sdk.Coins     `json:"available_liquidity" yaml:"available_liquidity"`
	ReservedLiquidity sdk.Coins      `json:"reserved_liquidity" yaml:"reserved_liquidity"`
	MinimumLiquidity  sdk.Int        `json:"minimum_liquidity" yaml:"minimum_liquidity"`
	
	// Community validators (local trusted nodes)
	LocalValidators   []LocalValidator `json:"local_validators" yaml:"local_validators"`
	RequiredSignatures uint32         `json:"required_signatures" yaml:"required_signatures"`
	
	// Fee structure (community benefits)
	BaseTradingFee    sdk.Dec        `json:"base_trading_fee" yaml:"base_trading_fee"`
	CommunityFee      sdk.Dec        `json:"community_fee" yaml:"community_fee"` // Goes to village fund
	EducationFee      sdk.Dec        `json:"education_fee" yaml:"education_fee"` // For financial literacy
	InfrastructureFee sdk.Dec        `json:"infrastructure_fee" yaml:"infrastructure_fee"` // For local infra
	
	// Member management
	TotalMembers      uint64         `json:"total_members" yaml:"total_members"`
	ActiveTraders     uint64         `json:"active_traders" yaml:"active_traders"`
	MemberBenefits    MemberBenefits `json:"member_benefits" yaml:"member_benefits"`
	
	// Trading statistics
	TotalVolume       sdk.Int        `json:"total_volume" yaml:"total_volume"`
	DailyVolume       sdk.Int        `json:"daily_volume" yaml:"daily_volume"`
	MonthlyVolume     sdk.Int        `json:"monthly_volume" yaml:"monthly_volume"`
	TotalTransactions uint64         `json:"total_transactions" yaml:"total_transactions"`
	
	// Community funds
	CommunityFund     sdk.Coins      `json:"community_fund" yaml:"community_fund"`
	EducationFund     sdk.Coins      `json:"education_fund" yaml:"education_fund"`
	EmergencyFund     sdk.Coins      `json:"emergency_fund" yaml:"emergency_fund"`
	
	// Cultural features
	LocalFestivals    []string       `json:"local_festivals" yaml:"local_festivals"`
	LanguageSupport   []string       `json:"language_support" yaml:"language_support"`
	CulturalQuotes    []string       `json:"cultural_quotes" yaml:"cultural_quotes"`
	
	// Pool status
	Active            bool           `json:"active" yaml:"active"`
	Verified          bool           `json:"verified" yaml:"verified"` // Government verification
	EstablishedDate   time.Time      `json:"established_date" yaml:"established_date"`
	LastActivityDate  time.Time      `json:"last_activity_date" yaml:"last_activity_date"`
	
	// Achievements
	Achievements      []VillageAchievement `json:"achievements" yaml:"achievements"`
	TrustScore        uint32         `json:"trust_score" yaml:"trust_score"` // 0-100
}

// LocalValidator represents a trusted community validator
type LocalValidator struct {
	ValidatorAddress  sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	LocalName         string         `json:"local_name" yaml:"local_name"`
	Role              string         `json:"role" yaml:"role"` // "panchayat_member", "teacher", "postmaster"
	TrustLevel        uint32         `json:"trust_level" yaml:"trust_level"` // 1-10
	JoinedAt          time.Time      `json:"joined_at" yaml:"joined_at"`
}

// MemberBenefits defines benefits for village pool members
type MemberBenefits struct {
	FeeDiscount       sdk.Dec        `json:"fee_discount" yaml:"fee_discount"`
	PriorityExecution bool           `json:"priority_execution" yaml:"priority_execution"`
	VotingRights      bool           `json:"voting_rights" yaml:"voting_rights"`
	ProfitSharing     sdk.Dec        `json:"profit_sharing" yaml:"profit_sharing"`
	EducationAccess   bool           `json:"education_access" yaml:"education_access"`
	EmergencySupport  bool           `json:"emergency_support" yaml:"emergency_support"`
}

// VillageAchievement represents milestones achieved by the village
type VillageAchievement struct {
	AchievementId     string         `json:"achievement_id" yaml:"achievement_id"`
	Title             string         `json:"title" yaml:"title"`
	Description       string         `json:"description" yaml:"description"`
	Category          string         `json:"category" yaml:"category"` // "trading", "education", "community"
	AchievedAt        time.Time      `json:"achieved_at" yaml:"achieved_at"`
	RewardAmount      sdk.Coins      `json:"reward_amount" yaml:"reward_amount"`
}

// VillagePoolMember represents a member of the village pool
type VillagePoolMember struct {
	MemberAddress     sdk.AccAddress `json:"member_address" yaml:"member_address"`
	LocalName         string         `json:"local_name" yaml:"local_name"`
	MobileNumber      string         `json:"mobile_number" yaml:"mobile_number"`
	AadhaarHash       string         `json:"aadhaar_hash" yaml:"aadhaar_hash"` // Hashed for privacy
	
	// Membership details
	MembershipType    string         `json:"membership_type" yaml:"membership_type"` // "regular", "premium", "founder"
	JoinedAt          time.Time      `json:"joined_at" yaml:"joined_at"`
	Contribution      sdk.Coins      `json:"contribution" yaml:"contribution"`
	
	// Activity tracking
	TotalTrades       uint64         `json:"total_trades" yaml:"total_trades"`
	TotalVolume       sdk.Int        `json:"total_volume" yaml:"total_volume"`
	LastTradeAt       time.Time      `json:"last_trade_at" yaml:"last_trade_at"`
	
	// Benefits earned
	TotalEarnings     sdk.Coins      `json:"total_earnings" yaml:"total_earnings"`
	PendingRewards    sdk.Coins      `json:"pending_rewards" yaml:"pending_rewards"`
	EducationCredits  uint32         `json:"education_credits" yaml:"education_credits"`
}

// NewVillagePool creates a new village pool
func NewVillagePool(
	poolId uint64,
	villageName string,
	postalCode string,
	stateCode string,
	panchayatHead sdk.AccAddress,
) *VillagePool {
	now := time.Now()
	
	return &VillagePool{
		PoolId:          poolId,
		VillageName:     villageName,
		PostalCode:      postalCode,
		StateCode:       stateCode,
		PanchayatHead:   panchayatHead,
		
		// Default configuration
		PrimaryToken:    "unamo",
		LocalCurrency:   "uinr",
		SupportedTokens: []string{"unamo", "uinr"},
		
		// Initial liquidity settings
		MinimumLiquidity: MinVillagePoolLiquidity,
		
		// Default fee structure (favorable for villages)
		BaseTradingFee:    sdk.MustNewDecFromStr("0.002"), // 0.2% (lower than standard)
		CommunityFee:      sdk.MustNewDecFromStr("0.02"),  // 2% of fees to community
		EducationFee:      sdk.MustNewDecFromStr("0.01"),  // 1% for education
		InfrastructureFee: sdk.MustNewDecFromStr("0.01"),  // 1% for infrastructure
		
		// Default member benefits
		MemberBenefits: MemberBenefits{
			FeeDiscount:       sdk.MustNewDecFromStr(VillagePoolDiscount),
			PriorityExecution: true,
			VotingRights:      true,
			ProfitSharing:     sdk.MustNewDecFromStr("0.50"), // 50% profit sharing
			EducationAccess:   true,
			EmergencySupport:  true,
		},
		
		// Validator requirements
		RequiredSignatures: 2, // At least 2 local validators
		
		// Initial state
		Active:           true,
		Verified:         false, // Requires government verification
		EstablishedDate:  now,
		LastActivityDate: now,
		TrustScore:       50, // Start with medium trust
		
		// Initialize statistics
		TotalMembers:      0,
		ActiveTraders:     0,
		TotalVolume:       sdk.ZeroInt(),
		DailyVolume:       sdk.ZeroInt(),
		MonthlyVolume:     sdk.ZeroInt(),
		TotalTransactions: 0,
	}
}

// AddLiquidity adds liquidity to the village pool
func (vp *VillagePool) AddLiquidity(amount sdk.Coins) error {
	if !vp.Active {
		return fmt.Errorf("village pool is not active")
	}
	
	vp.TotalLiquidity = vp.TotalLiquidity.Add(amount...)
	vp.AvailableLiquidity = vp.AvailableLiquidity.Add(amount...)
	vp.LastActivityDate = time.Now()
	
	return nil
}

// CalculateMemberFee calculates trading fee for a member
func (vp *VillagePool) CalculateMemberFee(baseAmount sdk.Int, isMember bool) sdk.Int {
	fee := vp.BaseTradingFee.MulInt(baseAmount)
	
	if isMember {
		// Apply member discount
		discount := fee.Mul(vp.MemberBenefits.FeeDiscount)
		fee = fee.Sub(discount)
	}
	
	return fee.TruncateInt()
}

// DistributeFees distributes collected fees to various funds
func (vp *VillagePool) DistributeFees(totalFees sdk.Coin) map[string]sdk.Coin {
	distribution := make(map[string]sdk.Coin)
	
	// Community fund (for village development)
	communityAmount := vp.CommunityFee.MulInt(totalFees.Amount).TruncateInt()
	distribution["community"] = sdk.NewCoin(totalFees.Denom, communityAmount)
	
	// Education fund (for financial literacy)
	educationAmount := vp.EducationFee.MulInt(totalFees.Amount).TruncateInt()
	distribution["education"] = sdk.NewCoin(totalFees.Denom, educationAmount)
	
	// Infrastructure fund (for local development)
	infraAmount := vp.InfrastructureFee.MulInt(totalFees.Amount).TruncateInt()
	distribution["infrastructure"] = sdk.NewCoin(totalFees.Denom, infraAmount)
	
	// Remaining goes to liquidity providers
	totalDistributed := communityAmount.Add(educationAmount).Add(infraAmount)
	remainingAmount := totalFees.Amount.Sub(totalDistributed)
	distribution["liquidity_providers"] = sdk.NewCoin(totalFees.Denom, remainingAmount)
	
	return distribution
}

// UpdateTrustScore updates the village pool's trust score
func (vp *VillagePool) UpdateTrustScore(delta int32) {
	newScore := int32(vp.TrustScore) + delta
	
	// Keep score between 0 and 100
	if newScore < 0 {
		vp.TrustScore = 0
	} else if newScore > 100 {
		vp.TrustScore = 100
	} else {
		vp.TrustScore = uint32(newScore)
	}
}

// AddAchievement adds a new achievement to the village
func (vp *VillagePool) AddAchievement(achievement VillageAchievement) {
	vp.Achievements = append(vp.Achievements, achievement)
	
	// Update trust score based on achievement
	switch achievement.Category {
	case "trading":
		vp.UpdateTrustScore(5)
	case "education":
		vp.UpdateTrustScore(3)
	case "community":
		vp.UpdateTrustScore(4)
	}
}

// IsEligibleForGovernmentSupport checks if pool qualifies for government benefits
func (vp *VillagePool) IsEligibleForGovernmentSupport() bool {
	return vp.Verified && 
		vp.TrustScore >= 70 && 
		vp.TotalMembers >= 50 && 
		vp.Active
}

// ValidateVillagePool ensures pool parameters are valid
func (vp *VillagePool) ValidateVillagePool() error {
	if vp.VillageName == "" {
		return fmt.Errorf("village name cannot be empty")
	}
	
	if vp.PostalCode == "" || len(vp.PostalCode) != 6 {
		return fmt.Errorf("invalid postal code")
	}
	
	if vp.PanchayatHead.Empty() {
		return fmt.Errorf("panchayat head address cannot be empty")
	}
	
	if vp.RequiredSignatures == 0 {
		return fmt.Errorf("required signatures must be greater than 0")
	}
	
	if vp.BaseTradingFee.IsNegative() || vp.BaseTradingFee.GT(sdk.MustNewDecFromStr("0.01")) {
		return fmt.Errorf("invalid base trading fee")
	}
	
	return nil
}

// GetEffectiveVolume returns the effective trading volume for rewards
func (vp *VillagePool) GetEffectiveVolume() sdk.Int {
	// Apply trust score multiplier
	multiplier := sdk.NewDec(int64(vp.TrustScore)).Quo(sdk.NewDec(100))
	return multiplier.MulInt(vp.TotalVolume).TruncateInt()
}