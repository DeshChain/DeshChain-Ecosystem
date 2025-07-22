package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewSewaMitraAgent creates a new Sewa Mitra agent
func NewSewaMitraAgent(
	agentID, agentAddress, agentName, businessName string,
	country, state, city, postalCode string,
	addressLine1, addressLine2, phone, email string,
	supportedCurrencies []string,
	supportedMethods []SettlementMethod,
	liquidityLimit, dailyLimit sdk.Coin,
	baseCommissionRate, volumeBonusRate sdk.Dec,
	minimumCommission, maximumCommission sdk.Coin,
) SewaMitraAgent {
	return SewaMitraAgent{
		AgentId:                agentID,
		AgentAddress:           agentAddress,
		AgentName:              agentName,
		BusinessName:           businessName,
		Country:                country,
		State:                  state,
		City:                   city,
		PostalCode:             postalCode,
		AddressLine1:           addressLine1,
		AddressLine2:           addressLine2,
		Phone:                  phone,
		Email:                  email,
		Status:                 AGENT_STATUS_PENDING_VERIFICATION,
		SupportedCurrencies:    supportedCurrencies,
		SupportedMethods:       supportedMethods,
		LiquidityLimit:         liquidityLimit,
		DailyLimit:             dailyLimit,
		BaseCommissionRate:     baseCommissionRate,
		VolumeBonus:            volumeBonusRate,
		MinimumCommission:      minimumCommission,
		MaximumCommission:      maximumCommission,
		TotalTransactions:      0,
		TotalVolume:            sdk.NewCoin("usd", sdk.ZeroInt()),
		TotalCommissionsEarned: sdk.NewCoin("usd", sdk.ZeroInt()),
		SuccessRate:            sdk.ZeroDec(),
		AverageProcessingTime:  sdk.ZeroDec(),
		KycLevel:               KYC_LEVEL_BASIC,
		BackgroundVerified:     false,
		Certifications:         []string{},
		Metadata:               make(map[string]string),
	}
}

// Validate performs basic validation on a Sewa Mitra agent
func (a SewaMitraAgent) Validate() error {
	if a.AgentId == "" {
		return fmt.Errorf("agent ID cannot be empty")
	}

	if _, err := sdk.AccAddressFromBech32(a.AgentAddress); err != nil {
		return fmt.Errorf("invalid agent address: %w", err)
	}

	if a.AgentName == "" {
		return fmt.Errorf("agent name cannot be empty")
	}

	if a.Country == "" || a.City == "" {
		return fmt.Errorf("country and city are required")
	}

	if a.Phone == "" {
		return fmt.Errorf("phone number is required")
	}

	if a.BaseCommissionRate.IsNegative() || a.BaseCommissionRate.GT(sdk.NewDecWithPrec(10, 2)) {
		return fmt.Errorf("base commission rate must be between 0% and 10%%")
	}

	if a.VolumeBonus.IsNegative() || a.VolumeBonus.GT(sdk.NewDecWithPrec(5, 2)) {
		return fmt.Errorf("volume bonus rate must be between 0% and 5%%")
	}

	if !a.LiquidityLimit.IsValid() || !a.LiquidityLimit.IsPositive() {
		return fmt.Errorf("liquidity limit must be positive")
	}

	if !a.DailyLimit.IsValid() || !a.DailyLimit.IsPositive() {
		return fmt.Errorf("daily limit must be positive")
	}

	if len(a.SupportedCurrencies) == 0 {
		return fmt.Errorf("at least one supported currency is required")
	}

	if len(a.SupportedMethods) == 0 {
		return fmt.Errorf("at least one supported settlement method is required")
	}

	return nil
}

// IsActive returns whether the agent is currently active
func (a SewaMitraAgent) IsActive() bool {
	return a.Status == AGENT_STATUS_ACTIVE
}

// SupportsSettlement returns whether the agent supports a specific settlement method
func (a SewaMitraAgent) SupportsSettlement(method SettlementMethod) bool {
	for _, supported := range a.SupportedMethods {
		if supported == method {
			return true
		}
	}
	return false
}

// SupportsCurrency returns whether the agent supports a specific currency
func (a SewaMitraAgent) SupportsCurrency(currency string) bool {
	for _, supported := range a.SupportedCurrencies {
		if strings.EqualFold(supported, currency) {
			return true
		}
	}
	return false
}

// GetLocationKey returns a standardized location key for the agent
func (a SewaMitraAgent) GetLocationKey() string {
	return fmt.Sprintf("%s/%s", a.Country, a.City)
}

// NewSewaMitraCommission creates a new commission record
func NewSewaMitraCommission(
	commissionID, agentID, transferID string,
	baseCommission, volumeBonus, totalCommission sdk.Coin,
) SewaMitraCommission {
	return SewaMitraCommission{
		CommissionId:    commissionID,
		AgentId:         agentID,
		TransferId:      transferID,
		BaseCommission:  baseCommission,
		VolumeBonus:     volumeBonus,
		TotalCommission: totalCommission,
		Status:          COMMISSION_STATUS_EARNED,
	}
}

// Validate performs basic validation on a commission record
func (c SewaMitraCommission) Validate() error {
	if c.CommissionId == "" {
		return fmt.Errorf("commission ID cannot be empty")
	}

	if c.AgentId == "" {
		return fmt.Errorf("agent ID cannot be empty")
	}

	if c.TransferId == "" {
		return fmt.Errorf("transfer ID cannot be empty")
	}

	if !c.BaseCommission.IsValid() {
		return fmt.Errorf("base commission must be valid")
	}

	if !c.VolumeBonus.IsValid() {
		return fmt.Errorf("volume bonus must be valid")
	}

	if !c.TotalCommission.IsValid() {
		return fmt.Errorf("total commission must be valid")
	}

	expectedTotal := c.BaseCommission.Add(c.VolumeBonus)
	if !c.TotalCommission.Equal(expectedTotal) {
		return fmt.Errorf("total commission must equal base commission plus volume bonus")
	}

	return nil
}

// IsPaid returns whether the commission has been paid
func (c SewaMitraCommission) IsPaid() bool {
	return c.Status == COMMISSION_STATUS_PAID
}

// Enhanced RemittanceTransfer methods for Sewa Mitra integration

// UsesSewaMitra returns whether this transfer uses the Sewa Mitra network
func (t RemittanceTransfer) UsesSewaMitra() bool {
	return t.UsesSewaMitra && t.SewaMitraAgentId != ""
}

// GetTotalFees calculates the total fees including Sewa Mitra commission
func (t RemittanceTransfer) GetTotalFees() sdk.Coin {
	totalFee := sdk.NewCoin(t.Amount.Denom, sdk.ZeroInt())
	
	// Add regular fees
	for _, fee := range t.Fees {
		if fee.Amount.Denom == t.Amount.Denom {
			totalFee = totalFee.Add(fee.Amount)
		}
	}
	
	// Add Sewa Mitra commission
	if t.UsesSewaMitra() && !t.SewaMitraCommission.IsZero() {
		if t.SewaMitraCommission.Denom == t.Amount.Denom {
			totalFee = totalFee.Add(t.SewaMitraCommission)
		}
	}
	
	return totalFee
}

// GetEffectiveAmount returns the amount after deducting all fees
func (t RemittanceTransfer) GetEffectiveAmount() sdk.Coin {
	return t.Amount.Sub(t.GetTotalFees())
}

// Validation functions

// ValidateAgentID validates a Sewa Mitra agent ID format
func ValidateAgentID(agentID string) error {
	if agentID == "" {
		return fmt.Errorf("agent ID cannot be empty")
	}
	
	if len(agentID) < 3 || len(agentID) > 50 {
		return fmt.Errorf("agent ID must be between 3 and 50 characters")
	}
	
	// Agent ID should start with "SMA-" prefix
	if !strings.HasPrefix(agentID, "SMA-") {
		return fmt.Errorf("agent ID must start with 'SMA-' prefix")
	}
	
	return nil
}

// ValidateCommissionRate validates commission rate values
func ValidateCommissionRate(rate sdk.Dec, maxRate sdk.Dec) error {
	if rate.IsNegative() {
		return fmt.Errorf("commission rate cannot be negative")
	}
	
	if rate.GT(maxRate) {
		return fmt.Errorf("commission rate cannot exceed %s", maxRate.String())
	}
	
	return nil
}

// ValidateLocationInfo validates agent location information
func ValidateLocationInfo(country, state, city string) error {
	if country == "" {
		return fmt.Errorf("country is required")
	}
	
	if city == "" {
		return fmt.Errorf("city is required")
	}
	
	// Country should be ISO 3166-1 alpha-2 format
	if len(country) != 2 {
		return fmt.Errorf("country must be 2-character ISO code")
	}
	
	return nil
}

// Helper functions for Sewa Mitra operations

// GenerateAgentID generates a unique agent ID
func GenerateAgentID(sequence uint64) string {
	return fmt.Sprintf("SMA-%06d", sequence)
}

// GenerateCommissionID generates a unique commission ID
func GenerateCommissionID(sequence uint64) string {
	return fmt.Sprintf("COM-%08d", sequence)
}

// CalculateDistanceScore calculates a simple distance score between locations
// Returns 0-100 where 100 is same city, 80 is same state, 60 is same country
func CalculateDistanceScore(fromCountry, fromCity, toCountry, toCity string) int {
	if fromCountry == toCountry && fromCity == toCity {
		return 100 // Same city
	}
	
	if fromCountry == toCountry {
		return 60 // Same country, different city
	}
	
	return 0 // Different country
}

// FormatAgentSummary returns a formatted summary of an agent
func FormatAgentSummary(agent SewaMitraAgent) string {
	return fmt.Sprintf(
		"Agent %s (%s) in %s, %s - Status: %s, Success Rate: %s%%, Transactions: %d",
		agent.AgentId,
		agent.AgentName,
		agent.City,
		agent.Country,
		agent.Status.String(),
		agent.SuccessRate.Mul(sdk.NewDec(100)).String(),
		agent.TotalTransactions,
	)
}

// IsSettlementMethodSewaMitra checks if a settlement method uses Sewa Mitra
func IsSettlementMethodSewaMitra(method SettlementMethod) bool {
	return method == SETTLEMENT_METHOD_SEVA_MITRA_AGENT ||
		   method == SETTLEMENT_METHOD_CASH_PICKUP ||
		   method == SETTLEMENT_METHOD_HOME_DELIVERY
}

// GetSewaMitraCompatibleMethods returns settlement methods compatible with Sewa Mitra
func GetSewaMitraCompatibleMethods() []SettlementMethod {
	return []SettlementMethod{
		SETTLEMENT_METHOD_SEVA_MITRA_AGENT,
		SETTLEMENT_METHOD_CASH_PICKUP,
		SETTLEMENT_METHOD_HOME_DELIVERY,
		SETTLEMENT_METHOD_MOBILE_WALLET,
	}
}