package types

// Event types for revenue module
const (
	EventTypeRevenueCollected   = "revenue_collected"
	EventTypeRevenueDistributed = "revenue_distributed"
	EventTypeRoyaltyClaimed     = "royalty_claimed"
	EventTypeBeneficiaryUpdated = "beneficiary_updated"
	EventTypeInheritanceTriggered = "inheritance_triggered"
	EventTypeRevenueStreamCreated = "revenue_stream_created"
	EventTypeRevenueStreamUpdated = "revenue_stream_updated"
	EventTypeMonthlyReport      = "monthly_revenue_report"
)

// Attribute keys for revenue events
const (
	AttributeKeyStreamID         = "stream_id"
	AttributeKeyStreamType       = "stream_type"
	AttributeKeyRevenueAmount    = "revenue_amount"
	AttributeKeyDistributionID   = "distribution_id"
	AttributeKeyDevelopmentAmount = "development_amount"
	AttributeKeyCommunityAmount  = "community_amount"
	AttributeKeyLiquidityAmount  = "liquidity_amount"
	AttributeKeyNGOAmount        = "ngo_amount"
	AttributeKeyEmergencyAmount  = "emergency_amount"
	AttributeKeyRoyaltyAmount    = "royalty_amount"
	AttributeKeyBeneficiary      = "beneficiary"
	AttributeKeyOldBeneficiary   = "old_beneficiary"
	AttributeKeyNewBeneficiary   = "new_beneficiary"
	AttributeKeyClaimAmount      = "claim_amount"
	AttributeKeyMonthlyTotal     = "monthly_total"
	AttributeKeyBlockHeight      = "block_height"
	AttributeKeyTimestamp        = "timestamp"
)

// EmitRevenueCollectedEvent creates attributes for revenue collection event
func EmitRevenueCollectedEvent(streamID, streamType string, amount string) []string {
	return []string{
		EventTypeRevenueCollected,
		AttributeKeyStreamID, streamID,
		AttributeKeyStreamType, streamType,
		AttributeKeyRevenueAmount, amount,
	}
}

// EmitRevenueDistributedEvent creates attributes for revenue distribution event
func EmitRevenueDistributedEvent(
	distributionID string,
	totalAmount string,
	devAmount string,
	communityAmount string,
	liquidityAmount string,
	ngoAmount string,
	emergencyAmount string,
	royaltyAmount string,
) []string {
	return []string{
		EventTypeRevenueDistributed,
		AttributeKeyDistributionID, distributionID,
		AttributeKeyRevenueAmount, totalAmount,
		AttributeKeyDevelopmentAmount, devAmount,
		AttributeKeyCommunityAmount, communityAmount,
		AttributeKeyLiquidityAmount, liquidityAmount,
		AttributeKeyNGOAmount, ngoAmount,
		AttributeKeyEmergencyAmount, emergencyAmount,
		AttributeKeyRoyaltyAmount, royaltyAmount,
	}
}

// EmitRoyaltyClaimedEvent creates attributes for royalty claim event
func EmitRoyaltyClaimedEvent(beneficiary, claimAmount string) []string {
	return []string{
		EventTypeRoyaltyClaimed,
		AttributeKeyBeneficiary, beneficiary,
		AttributeKeyClaimAmount, claimAmount,
	}
}

// EmitBeneficiaryUpdatedEvent creates attributes for beneficiary update event
func EmitBeneficiaryUpdatedEvent(oldBeneficiary, newBeneficiary string) []string {
	return []string{
		EventTypeBeneficiaryUpdated,
		AttributeKeyOldBeneficiary, oldBeneficiary,
		AttributeKeyNewBeneficiary, newBeneficiary,
	}
}

// EmitInheritanceTriggeredEvent creates attributes for inheritance trigger event
func EmitInheritanceTriggeredEvent(oldBeneficiary, newBeneficiary, amount string) []string {
	return []string{
		EventTypeInheritanceTriggered,
		AttributeKeyOldBeneficiary, oldBeneficiary,
		AttributeKeyNewBeneficiary, newBeneficiary,
		AttributeKeyRoyaltyAmount, amount,
	}
}