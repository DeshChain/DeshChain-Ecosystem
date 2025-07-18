package types

import (
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RoyaltyType represents the type of royalty
type RoyaltyType string

const (
	// RoyaltyTypeTransaction for transaction tax royalties (0.10%)
	RoyaltyTypeTransaction RoyaltyType = "transaction"
	
	// RoyaltyTypePlatform for platform revenue royalties (5%)
	RoyaltyTypePlatform RoyaltyType = "platform"
)

// RoyaltyConfig represents the main royalty configuration
type RoyaltyConfig struct {
	// Current beneficiary who receives royalties
	Beneficiary string `json:"beneficiary"`
	
	// Backup beneficiaries for inheritance (in order of priority)
	BackupBeneficiaries []string `json:"backup_beneficiaries"`
	
	// Transaction royalty rate (0.10% of transaction tax)
	TransactionRoyaltyRate sdk.Dec `json:"transaction_royalty_rate"`
	
	// Platform revenue royalty rate (5% of platform revenues)
	PlatformRoyaltyRate sdk.Dec `json:"platform_royalty_rate"`
	
	// Whether inheritance is enabled
	InheritanceEnabled bool `json:"inheritance_enabled"`
	
	// Lock period before inheritance can be triggered (in days)
	InheritanceLockDays int64 `json:"inheritance_lock_days"`
	
	// Last activity timestamp of current beneficiary
	LastActivityTime time.Time `json:"last_activity_time"`
	
	// Total royalties earned lifetime
	TotalEarned RoyaltyBalance `json:"total_earned"`
	
	// Current unclaimed balance
	UnclaimedBalance RoyaltyBalance `json:"unclaimed_balance"`
	
	// Configuration active status
	Active bool `json:"active"`
	
	// Creation timestamp
	CreatedAt time.Time `json:"created_at"`
	
	// Last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
}

// RoyaltyBalance tracks royalty amounts by type
type RoyaltyBalance struct {
	// Transaction-based royalties
	TransactionRoyalties sdk.Coin `json:"transaction_royalties"`
	
	// Platform revenue royalties
	PlatformRoyalties sdk.Coin `json:"platform_royalties"`
	
	// Total combined royalties
	Total sdk.Coin `json:"total"`
}

// RoyaltyClaim represents a royalty claim by beneficiary
type RoyaltyClaim struct {
	// Unique claim ID
	ID uint64 `json:"id"`
	
	// Beneficiary who made the claim
	Beneficiary string `json:"beneficiary"`
	
	// Amount claimed
	ClaimAmount RoyaltyBalance `json:"claim_amount"`
	
	// Claim timestamp
	ClaimTime time.Time `json:"claim_time"`
	
	// Block height of claim
	BlockHeight int64 `json:"block_height"`
	
	// Transaction hash
	TxHash string `json:"tx_hash"`
	
	// Status of claim
	Status ClaimStatus `json:"status"`
}

// ClaimStatus represents the status of a royalty claim
type ClaimStatus string

const (
	ClaimStatusPending   ClaimStatus = "pending"
	ClaimStatusCompleted ClaimStatus = "completed"
	ClaimStatusFailed    ClaimStatus = "failed"
)

// InheritanceRecord tracks inheritance events
type InheritanceRecord struct {
	// Record ID
	ID uint64 `json:"id"`
	
	// Previous beneficiary
	PreviousBeneficiary string `json:"previous_beneficiary"`
	
	// New beneficiary
	NewBeneficiary string `json:"new_beneficiary"`
	
	// Reason for inheritance trigger
	TriggerReason InheritanceTrigger `json:"trigger_reason"`
	
	// Balance transferred
	TransferredBalance RoyaltyBalance `json:"transferred_balance"`
	
	// Timestamp of inheritance
	InheritanceTime time.Time `json:"inheritance_time"`
	
	// Block height
	BlockHeight int64 `json:"block_height"`
	
	// Additional notes
	Notes string `json:"notes"`
}

// InheritanceTrigger represents reasons for inheritance
type InheritanceTrigger string

const (
	// InheritanceTriggerInactivity when beneficiary is inactive beyond lock period
	InheritanceTriggerInactivity InheritanceTrigger = "inactivity"
	
	// InheritanceTriggerVoluntary when beneficiary voluntarily transfers
	InheritanceTriggerVoluntary InheritanceTrigger = "voluntary"
	
	// InheritanceTriggerGovernance when triggered by governance
	InheritanceTriggerGovernance InheritanceTrigger = "governance"
	
	// InheritanceTriggerEmergency for emergency situations
	InheritanceTriggerEmergency InheritanceTrigger = "emergency"
)

// BeneficiaryHistory tracks historical beneficiaries
type BeneficiaryHistory struct {
	// Beneficiary address
	Address string `json:"address"`
	
	// Start time as beneficiary
	StartTime time.Time `json:"start_time"`
	
	// End time as beneficiary (empty if current)
	EndTime *time.Time `json:"end_time,omitempty"`
	
	// Total earned during tenure
	TotalEarned RoyaltyBalance `json:"total_earned"`
	
	// Total claimed during tenure
	TotalClaimed RoyaltyBalance `json:"total_claimed"`
	
	// Reason for change
	ChangeReason string `json:"change_reason"`
}

// RoyaltyAccumulator tracks accumulated royalties by source
type RoyaltyAccumulator struct {
	// Source identifier (e.g., "dex_trading", "nft_marketplace")
	Source string `json:"source"`
	
	// Type of royalty
	Type RoyaltyType `json:"type"`
	
	// Accumulated amount
	Amount sdk.Coin `json:"amount"`
	
	// Last accumulation timestamp
	LastAccumulation time.Time `json:"last_accumulation"`
	
	// Number of accumulations
	AccumulationCount uint64 `json:"accumulation_count"`
}

// NewRoyaltyConfig creates a new royalty configuration
func NewRoyaltyConfig(beneficiary string) RoyaltyConfig {
	transactionRate, _ := sdk.NewDecFromStr("0.001") // 0.10%
	platformRate, _ := sdk.NewDecFromStr("0.05")     // 5%
	
	return RoyaltyConfig{
		Beneficiary:            beneficiary,
		BackupBeneficiaries:    []string{},
		TransactionRoyaltyRate: transactionRate,
		PlatformRoyaltyRate:    platformRate,
		InheritanceEnabled:     true,
		InheritanceLockDays:    90, // 90 days inactivity before inheritance
		LastActivityTime:       time.Now(),
		TotalEarned:            NewRoyaltyBalance(),
		UnclaimedBalance:       NewRoyaltyBalance(),
		Active:                 true,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
}

// NewRoyaltyBalance creates a new royalty balance
func NewRoyaltyBalance() RoyaltyBalance {
	return RoyaltyBalance{
		TransactionRoyalties: sdk.NewCoin("namo", sdk.ZeroInt()),
		PlatformRoyalties:    sdk.NewCoin("namo", sdk.ZeroInt()),
		Total:                sdk.NewCoin("namo", sdk.ZeroInt()),
	}
}

// Add adds royalty amounts to the balance
func (rb *RoyaltyBalance) Add(royaltyType RoyaltyType, amount sdk.Coin) error {
	if amount.IsNegative() {
		return fmt.Errorf("royalty amount cannot be negative")
	}
	
	switch royaltyType {
	case RoyaltyTypeTransaction:
		rb.TransactionRoyalties = rb.TransactionRoyalties.Add(amount)
	case RoyaltyTypePlatform:
		rb.PlatformRoyalties = rb.PlatformRoyalties.Add(amount)
	default:
		return fmt.Errorf("unknown royalty type: %s", royaltyType)
	}
	
	// Update total
	rb.Total = rb.TransactionRoyalties.Add(rb.PlatformRoyalties)
	
	return nil
}

// CanTriggerInheritance checks if inheritance can be triggered
func (rc RoyaltyConfig) CanTriggerInheritance(currentTime time.Time) (bool, InheritanceTrigger) {
	if !rc.InheritanceEnabled {
		return false, ""
	}
	
	// Check for inactivity
	inactiveDays := int64(currentTime.Sub(rc.LastActivityTime).Hours() / 24)
	if inactiveDays >= rc.InheritanceLockDays {
		return true, InheritanceTriggerInactivity
	}
	
	return false, ""
}

// GetNextBeneficiary returns the next beneficiary in line for inheritance
func (rc RoyaltyConfig) GetNextBeneficiary() (string, error) {
	if len(rc.BackupBeneficiaries) == 0 {
		return "", fmt.Errorf("no backup beneficiaries configured")
	}
	
	return rc.BackupBeneficiaries[0], nil
}

// UpdateActivity updates the last activity time
func (rc *RoyaltyConfig) UpdateActivity(timestamp time.Time) {
	rc.LastActivityTime = timestamp
	rc.UpdatedAt = timestamp
}

// ValidateConfig validates the royalty configuration
func (rc RoyaltyConfig) ValidateConfig() error {
	// Validate beneficiary address
	if _, err := sdk.AccAddressFromBech32(rc.Beneficiary); err != nil {
		return fmt.Errorf("invalid beneficiary address: %w", err)
	}
	
	// Validate backup beneficiaries
	for i, backup := range rc.BackupBeneficiaries {
		if _, err := sdk.AccAddressFromBech32(backup); err != nil {
			return fmt.Errorf("invalid backup beneficiary at index %d: %w", i, err)
		}
		
		// Check for duplicates
		if backup == rc.Beneficiary {
			return fmt.Errorf("backup beneficiary cannot be same as primary beneficiary")
		}
	}
	
	// Validate royalty rates
	if rc.TransactionRoyaltyRate.IsNegative() || rc.TransactionRoyaltyRate.GT(sdk.OneDec()) {
		return fmt.Errorf("transaction royalty rate must be between 0 and 1")
	}
	
	if rc.PlatformRoyaltyRate.IsNegative() || rc.PlatformRoyaltyRate.GT(sdk.OneDec()) {
		return fmt.Errorf("platform royalty rate must be between 0 and 1")
	}
	
	// Validate inheritance lock days
	if rc.InheritanceLockDays < 30 || rc.InheritanceLockDays > 365 {
		return fmt.Errorf("inheritance lock days must be between 30 and 365")
	}
	
	return nil
}

// CalculateRoyalty calculates royalty amount based on type
func (rc RoyaltyConfig) CalculateRoyalty(royaltyType RoyaltyType, baseAmount sdk.Coin) (sdk.Coin, error) {
	var rate sdk.Dec
	
	switch royaltyType {
	case RoyaltyTypeTransaction:
		rate = rc.TransactionRoyaltyRate
	case RoyaltyTypePlatform:
		rate = rc.PlatformRoyaltyRate
	default:
		return sdk.Coin{}, fmt.Errorf("unknown royalty type: %s", royaltyType)
	}
	
	royaltyAmount := rate.MulInt(baseAmount.Amount).TruncateInt()
	return sdk.NewCoin(baseAmount.Denom, royaltyAmount), nil
}

// GetClaimableBalance returns the total claimable balance
func (rc RoyaltyConfig) GetClaimableBalance() sdk.Coin {
	return rc.UnclaimedBalance.Total
}

// CanClaim checks if beneficiary can claim royalties
func (rc RoyaltyConfig) CanClaim(minClaimAmount sdk.Coin) bool {
	return rc.Active && rc.UnclaimedBalance.Total.IsGTE(minClaimAmount)
}

// ProcessClaim processes a royalty claim
func (rc *RoyaltyConfig) ProcessClaim(claimAmount RoyaltyBalance) error {
	// Verify sufficient balance
	if rc.UnclaimedBalance.Total.IsLT(claimAmount.Total) {
		return fmt.Errorf("insufficient balance: available %s, requested %s", 
			rc.UnclaimedBalance.Total, claimAmount.Total)
	}
	
	// Deduct from unclaimed balance
	rc.UnclaimedBalance.TransactionRoyalties = rc.UnclaimedBalance.TransactionRoyalties.Sub(claimAmount.TransactionRoyalties)
	rc.UnclaimedBalance.PlatformRoyalties = rc.UnclaimedBalance.PlatformRoyalties.Sub(claimAmount.PlatformRoyalties)
	rc.UnclaimedBalance.Total = rc.UnclaimedBalance.Total.Sub(claimAmount.Total)
	
	// Update activity
	rc.UpdateActivity(time.Now())
	
	return nil
}

// TriggerInheritance transfers beneficiary rights to next in line
func (rc *RoyaltyConfig) TriggerInheritance(newBeneficiary string, trigger InheritanceTrigger) (*InheritanceRecord, error) {
	if !rc.InheritanceEnabled {
		return nil, fmt.Errorf("inheritance is not enabled")
	}
	
	// Create inheritance record
	record := &InheritanceRecord{
		PreviousBeneficiary: rc.Beneficiary,
		NewBeneficiary:      newBeneficiary,
		TriggerReason:       trigger,
		TransferredBalance:  rc.UnclaimedBalance,
		InheritanceTime:     time.Now(),
		BlockHeight:         0, // To be set by keeper
	}
	
	// Update beneficiary
	rc.Beneficiary = newBeneficiary
	
	// Remove new beneficiary from backup list if present
	newBackups := []string{}
	for _, backup := range rc.BackupBeneficiaries {
		if backup != newBeneficiary {
			newBackups = append(newBackups, backup)
		}
	}
	rc.BackupBeneficiaries = newBackups
	
	// Reset activity time
	rc.UpdateActivity(time.Now())
	
	return record, nil
}