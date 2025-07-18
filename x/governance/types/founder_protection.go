package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Immutable founder protections - these constants can NEVER be changed
const (
	// FounderTokenAllocationPercent is immutably set to 10%
	FounderTokenAllocationPercent = 10
	
	// FounderTaxRoyaltyPercent is immutably set to 0.10%
	FounderTaxRoyaltyPercent = 0.10
	
	// FounderPlatformRoyaltyPercent is immutably set to 5%
	FounderPlatformRoyaltyPercent = 5
	
	// FounderMinimumVotingPowerPercent is guaranteed 15%
	FounderMinimumVotingPowerPercent = 15
	
	// FounderVetoDurationYears is 3 years from genesis
	FounderVetoDurationYears = 3
	
	// SupermajorityThresholdPercent to override founder is 80%
	SupermajorityThresholdPercent = 80
	
	// GovernanceChangeNoticeDays is 90 days
	GovernanceChangeNoticeDays = 90
)

// Protected parameter keys that require special governance
var (
	// Immutable parameters - can NEVER be changed
	ImmutableParameters = []string{
		"founder_token_allocation",
		"founder_tax_royalty",
		"founder_platform_royalty",
		"founder_inheritance_mechanism",
		"founder_minimum_voting_power",
		"founder_protection_removal", // This protection itself cannot be removed
	}
	
	// Founder consent required parameters
	FounderConsentParameters = []string{
		"chain_upgrade_handler",
		"crisis_module_permissions",
		"slashing_parameters",
		"consensus_parameters",
		"ibc_transfer_enabled",
		"wasm_permissions",
	}
	
	// Supermajority required parameters (80%)
	SupermajorityParameters = []string{
		"governance_voting_period",
		"governance_deposit_amount",
		"distribution_community_tax",
		"mint_inflation_rate",
		"staking_unbonding_time",
	}
)

// ValidateFounderProtection validates that founder protections are not violated
func ValidateFounderProtection(paramKey string, currentValue, proposedValue interface{}) error {
	// Check if parameter is immutable
	for _, immutable := range ImmutableParameters {
		if paramKey == immutable {
			return fmt.Errorf("parameter %s is immutable and cannot be changed", paramKey)
		}
	}
	
	// Validate specific protections
	switch paramKey {
	case "tax_distribution":
		// Ensure founder royalty is maintained at 0.10%
		if !validateTaxDistribution(proposedValue) {
			return fmt.Errorf("founder tax royalty must remain at 0.10%%")
		}
	case "revenue_distribution":
		// Ensure founder platform royalty is maintained at 5%
		if !validateRevenueDistribution(proposedValue) {
			return fmt.Errorf("founder platform royalty must remain at 5%%")
		}
	case "token_distribution":
		// Token distribution cannot be changed after genesis
		return fmt.Errorf("token distribution is immutable after genesis")
	}
	
	return nil
}

// validateTaxDistribution ensures founder gets 0.10% from tax
func validateTaxDistribution(value interface{}) bool {
	// This would check the actual tax distribution structure
	// For now, returning true as placeholder
	// In real implementation, would parse the distribution and verify founder gets 0.10%
	return true
}

// validateRevenueDistribution ensures founder gets 5% from platform revenue
func validateRevenueDistribution(value interface{}) bool {
	// This would check the actual revenue distribution structure
	// For now, returning true as placeholder
	// In real implementation, would parse the distribution and verify founder gets 5%
	return true
}

// CalculateFounderVotingPower ensures founder always has minimum 15% voting power
func CalculateFounderVotingPower(totalVotingPower sdk.Int, founderTokens sdk.Int) sdk.Int {
	// Calculate actual voting power from tokens
	actualVotingPower := founderTokens
	
	// Calculate minimum guaranteed voting power (15% of total)
	minVotingPower := totalVotingPower.MulRaw(FounderMinimumVotingPowerPercent).QuoRaw(100)
	
	// Return the higher of actual or minimum
	if actualVotingPower.GT(minVotingPower) {
		return actualVotingPower
	}
	return minVotingPower
}

// CanFounderVeto checks if founder can veto a proposal
func CanFounderVeto(proposalType ProposalType, submitTime time.Time, genesisTime time.Time) bool {
	// Calculate veto expiry (3 years from genesis)
	vetoExpiry := genesisTime.Add(time.Duration(FounderVetoDurationYears) * 365 * 24 * time.Hour)
	
	// Check if we're still within veto period
	if submitTime.After(vetoExpiry) {
		return false
	}
	
	// Founder can veto these proposal types within the veto period
	vetoableTypes := []ProposalType{
		ProposalType_PROPOSAL_TYPE_PARAMETER_CHANGE,
		ProposalType_PROPOSAL_TYPE_SOFTWARE_UPGRADE,
		ProposalType_PROPOSAL_TYPE_REVENUE_DISTRIBUTION,
		ProposalType_PROPOSAL_TYPE_TAX_ADJUSTMENT,
		ProposalType_PROPOSAL_TYPE_FOUNDER_RELATED,
	}
	
	for _, vType := range vetoableTypes {
		if proposalType == vType {
			return true
		}
	}
	
	return false
}

// RequiresFounderConsent checks if a proposal requires founder approval
func RequiresFounderConsent(proposalType ProposalType, affectedParams []string) bool {
	// Any proposal affecting founder-related parameters requires consent
	if proposalType == ProposalType_PROPOSAL_TYPE_FOUNDER_RELATED {
		return true
	}
	
	// Check if any affected parameters require founder consent
	for _, param := range affectedParams {
		for _, protected := range FounderConsentParameters {
			if param == protected {
				return true
			}
		}
	}
	
	return false
}

// RequiresSupermajority checks if a proposal requires 80% supermajority
func RequiresSupermajority(proposalType ProposalType, affectedParams []string) bool {
	// Overriding founder decisions always requires supermajority
	if proposalType == ProposalType_PROPOSAL_TYPE_FOUNDER_RELATED {
		return true
	}
	
	// Check if any affected parameters require supermajority
	for _, param := range affectedParams {
		for _, param80 := range SupermajorityParameters {
			if param == param80 {
				return true
			}
		}
	}
	
	return false
}

// ValidateEmergencyAction validates if founder can take an emergency action
func ValidateEmergencyAction(actionType EmergencyActionType, founderAddr string, executorAddr string) error {
	// Only founder can execute emergency actions
	if founderAddr != executorAddr {
		return fmt.Errorf("only founder can execute emergency actions")
	}
	
	// Validate action type
	validActions := []EmergencyActionType{
		EmergencyActionType_EMERGENCY_ACTION_TYPE_HALT_CHAIN,
		EmergencyActionType_EMERGENCY_ACTION_TYPE_FREEZE_MODULE,
		EmergencyActionType_EMERGENCY_ACTION_TYPE_ROLLBACK_UPGRADE,
		EmergencyActionType_EMERGENCY_ACTION_TYPE_PATCH_VULNERABILITY,
	}
	
	valid := false
	for _, action := range validActions {
		if actionType == action {
			valid = true
			break
		}
	}
	
	if !valid {
		return fmt.Errorf("invalid emergency action type")
	}
	
	return nil
}

// GetProtectedParameters returns all parameters with their protection levels
func GetProtectedParameters() []ProtectedParameter {
	params := []ProtectedParameter{}
	
	// Add immutable parameters
	for _, param := range ImmutableParameters {
		params = append(params, ProtectedParameter{
			ParameterKey:    param,
			ProtectionType:  ProtectionType_PROTECTION_TYPE_IMMUTABLE,
			Description:     fmt.Sprintf("%s is permanently protected and cannot be changed", param),
		})
	}
	
	// Add founder consent parameters
	for _, param := range FounderConsentParameters {
		params = append(params, ProtectedParameter{
			ParameterKey:    param,
			ProtectionType:  ProtectionType_PROTECTION_TYPE_FOUNDER_CONSENT,
			Description:     fmt.Sprintf("%s requires founder approval to change", param),
		})
	}
	
	// Add supermajority parameters
	for _, param := range SupermajorityParameters {
		params = append(params, ProtectedParameter{
			ParameterKey:    param,
			ProtectionType:  ProtectionType_PROTECTION_TYPE_SUPERMAJORITY,
			Description:     fmt.Sprintf("%s requires 80%% supermajority to change", param),
		})
	}
	
	return params
}

// FounderProtectionConfig is the immutable configuration for founder protections
type FounderProtectionConfig struct {
	// These fields are set at genesis and can NEVER be changed
	FounderAddress              string
	TokenAllocation            sdk.Int
	TaxRoyaltyRate             sdk.Dec
	PlatformRoyaltyRate        sdk.Dec
	MinimumVotingPower         sdk.Dec
	VetoExpiryTime             time.Time
	SupermajorityThreshold     sdk.Dec
	GovernanceChangeNoticeDays int
}

// NewFounderProtectionConfig creates the immutable founder protection configuration
func NewFounderProtectionConfig(founderAddr string, genesisTime time.Time) FounderProtectionConfig {
	return FounderProtectionConfig{
		FounderAddress:              founderAddr,
		TokenAllocation:            sdk.NewInt(142862766), // 10% of total supply
		TaxRoyaltyRate:             sdk.NewDecWithPrec(10, 4), // 0.10%
		PlatformRoyaltyRate:        sdk.NewDecWithPrec(5, 2), // 5%
		MinimumVotingPower:         sdk.NewDecWithPrec(15, 2), // 15%
		VetoExpiryTime:             genesisTime.Add(time.Duration(FounderVetoDurationYears) * 365 * 24 * time.Hour),
		SupermajorityThreshold:     sdk.NewDecWithPrec(80, 2), // 80%
		GovernanceChangeNoticeDays: GovernanceChangeNoticeDays,
	}
}

// Validate ensures the protection config is valid
func (f FounderProtectionConfig) Validate() error {
	if f.FounderAddress == "" {
		return fmt.Errorf("founder address cannot be empty")
	}
	
	if !f.TokenAllocation.Equal(sdk.NewInt(142862766)) {
		return fmt.Errorf("founder token allocation must be exactly 142,862,766 NAMO")
	}
	
	if !f.TaxRoyaltyRate.Equal(sdk.NewDecWithPrec(10, 4)) {
		return fmt.Errorf("founder tax royalty must be exactly 0.10%%")
	}
	
	if !f.PlatformRoyaltyRate.Equal(sdk.NewDecWithPrec(5, 2)) {
		return fmt.Errorf("founder platform royalty must be exactly 5%%")
	}
	
	if !f.MinimumVotingPower.Equal(sdk.NewDecWithPrec(15, 2)) {
		return fmt.Errorf("founder minimum voting power must be exactly 15%%")
	}
	
	if !f.SupermajorityThreshold.Equal(sdk.NewDecWithPrec(80, 2)) {
		return fmt.Errorf("supermajority threshold must be exactly 80%%")
	}
	
	return nil
}

// IsProtected checks if a parameter is protected
func IsProtected(paramKey string) bool {
	// Check all protection lists
	allProtected := append(ImmutableParameters, FounderConsentParameters...)
	allProtected = append(allProtected, SupermajorityParameters...)
	
	for _, protected := range allProtected {
		if paramKey == protected {
			return true
		}
	}
	
	return false
}

// GetProtectionLevel returns the protection level for a parameter
func GetProtectionLevel(paramKey string) ProtectionType {
	// Check immutable
	for _, param := range ImmutableParameters {
		if paramKey == param {
			return ProtectionType_PROTECTION_TYPE_IMMUTABLE
		}
	}
	
	// Check founder consent
	for _, param := range FounderConsentParameters {
		if paramKey == param {
			return ProtectionType_PROTECTION_TYPE_FOUNDER_CONSENT
		}
	}
	
	// Check supermajority
	for _, param := range SupermajorityParameters {
		if paramKey == param {
			return ProtectionType_PROTECTION_TYPE_SUPERMAJORITY
		}
	}
	
	return ProtectionType_PROTECTION_TYPE_UNSPECIFIED
}