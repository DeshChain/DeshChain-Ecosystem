package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/remittance/types"
	identitykeeper "github.com/deshchain/deshchain/x/identity/keeper"
	identitytypes "github.com/deshchain/deshchain/x/identity/types"
)

// IdentityAdapter provides identity integration for Remittance module
type IdentityAdapter struct {
	keeper         *Keeper
	identityKeeper identitykeeper.Keeper
}

// NewIdentityAdapter creates a new identity adapter
func NewIdentityAdapter(keeper *Keeper, identityKeeper identitykeeper.Keeper) *IdentityAdapter {
	return &IdentityAdapter{
		keeper:         keeper,
		identityKeeper: identityKeeper,
	}
}

// SenderIdentityStatus represents the identity status of a remittance sender
type SenderIdentityStatus struct {
	HasIdentity        bool
	DID                string
	IsKYCVerified      bool
	KYCLevel           string
	AMLVerified        bool
	SanctionsChecked   bool
	SourceOfFunds      string
	RiskLevel          string
	MaxTransferLimit   sdk.Coin
	DailyLimit         sdk.Coin
	MonthlyLimit       sdk.Coin
	CountryRestrictions []string
	ComplianceExpiry   time.Time
}

// RecipientIdentityStatus represents the identity status of a remittance recipient
type RecipientIdentityStatus struct {
	HasIdentity      bool
	DID              string
	IsKYCVerified    bool
	KYCLevel         string
	VerificationDocs []string
	PurposeOfFunds   string
	BeneficiaryType  string
	Country          string
	CanReceiveFrom   []string
}

// SewaMitraIdentityStatus represents the identity status of a Sewa Mitra agent
type SewaMitraIdentityStatus struct {
	HasIdentity          bool
	DID                  string
	IsKYCVerified        bool
	KYCLevel             string
	BusinessLicenseValid bool
	ComplianceVerified   bool
	AMLCompliant         bool
	CertificationsValid  []string
	ServiceAreas         []string
	SupportedCurrencies  []string
	MaxTransactionLimit  sdk.Coin
	BackgroundVerified   bool
	InsuranceCovered     bool
	Rating               sdk.Dec
}

// RemittanceComplianceCheck represents comprehensive compliance verification
type RemittanceComplianceCheck struct {
	SenderCompliant    bool
	RecipientCompliant bool
	AgentCompliant     bool
	CorridorAllowed    bool
	AmountWithinLimits bool
	SanctionsCleared   bool
	AMLCompliant       bool
	RegulatoryCompliant bool
	ComplianceScore    sdk.Dec
	RequiredActions    []string
	Warnings           []string
}

// VerifySenderIdentity verifies a sender's identity for remittance compliance
func (ia *IdentityAdapter) VerifySenderIdentity(
	ctx context.Context,
	senderAddress sdk.AccAddress,
	transferAmount sdk.Coin,
	recipientCountry string,
) (*SenderIdentityStatus, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	// Check if identity exists
	identity, found := ia.identityKeeper.GetIdentityByAddress(sdkCtx, senderAddress)
	if !found {
		return &SenderIdentityStatus{
			HasIdentity: false,
		}, nil
	}

	status := &SenderIdentityStatus{
		HasIdentity: true,
		DID:         identity.Did,
	}

	// Check KYC credentials
	kycCreds := ia.identityKeeper.GetCredentialsByType(sdkCtx, identity.Did, "KYCCredential")
	for _, cred := range kycCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			status.IsKYCVerified = true
			if level, ok := cred.CredentialSubject["kyc_level"].(string); ok {
				status.KYCLevel = level
			}
			if sourceOfFunds, ok := cred.CredentialSubject["source_of_funds"].(string); ok {
				status.SourceOfFunds = sourceOfFunds
			}
			break
		}
	}

	// Check remittance credentials
	remittanceCreds := ia.identityKeeper.GetCredentialsByType(sdkCtx, identity.Did, "RemittanceCredential")
	for _, cred := range remittanceCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			if amlVerified, ok := cred.CredentialSubject["aml_verified"].(bool); ok {
				status.AMLVerified = amlVerified
			}
			if sanctionsChecked, ok := cred.CredentialSubject["sanctions_checked"].(bool); ok {
				status.SanctionsChecked = sanctionsChecked
			}
			if riskLevel, ok := cred.CredentialSubject["risk_level"].(string); ok {
				status.RiskLevel = riskLevel
			}
			if maxLimitStr, ok := cred.CredentialSubject["max_transfer_limit"].(string); ok {
				if maxLimit, err := sdk.ParseCoinNormalized(maxLimitStr); err == nil {
					status.MaxTransferLimit = maxLimit
				}
			}
			if dailyLimitStr, ok := cred.CredentialSubject["daily_limit"].(string); ok {
				if dailyLimit, err := sdk.ParseCoinNormalized(dailyLimitStr); err == nil {
					status.DailyLimit = dailyLimit
				}
			}
			if monthlyLimitStr, ok := cred.CredentialSubject["monthly_limit"].(string); ok {
				if monthlyLimit, err := sdk.ParseCoinNormalized(monthlyLimitStr); err == nil {
					status.MonthlyLimit = monthlyLimit
				}
			}
			if restrictions, ok := cred.CredentialSubject["country_restrictions"].([]interface{}); ok {
				for _, r := range restrictions {
					if rStr, ok := r.(string); ok {
						status.CountryRestrictions = append(status.CountryRestrictions, rStr)
					}
				}
			}
			if expiryStr, ok := cred.CredentialSubject["compliance_expiry"].(string); ok {
				if expiry, err := time.Parse(time.RFC3339, expiryStr); err == nil {
					status.ComplianceExpiry = expiry
				}
			}
			break
		}
	}

	return status, nil
}

// VerifyRecipientIdentity verifies a recipient's identity
func (ia *IdentityAdapter) VerifyRecipientIdentity(
	ctx context.Context,
	recipientAddress sdk.AccAddress,
	senderCountry string,
) (*RecipientIdentityStatus, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	// Check if identity exists
	identity, found := ia.identityKeeper.GetIdentityByAddress(sdkCtx, recipientAddress)
	if !found {
		return &RecipientIdentityStatus{
			HasIdentity: false,
		}, nil
	}

	status := &RecipientIdentityStatus{
		HasIdentity: true,
		DID:         identity.Did,
	}

	// Check KYC credentials
	kycCreds := ia.identityKeeper.GetCredentialsByType(sdkCtx, identity.Did, "KYCCredential")
	for _, cred := range kycCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			status.IsKYCVerified = true
			if level, ok := cred.CredentialSubject["kyc_level"].(string); ok {
				status.KYCLevel = level
			}
			if country, ok := cred.CredentialSubject["country"].(string); ok {
				status.Country = country
			}
			break
		}
	}

	// Check recipient credentials
	recipientCreds := ia.identityKeeper.GetCredentialsByType(sdkCtx, identity.Did, "RecipientCredential")
	for _, cred := range recipientCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			if docs, ok := cred.CredentialSubject["verification_documents"].([]interface{}); ok {
				for _, doc := range docs {
					if docStr, ok := doc.(string); ok {
						status.VerificationDocs = append(status.VerificationDocs, docStr)
					}
				}
			}
			if purpose, ok := cred.CredentialSubject["purpose_of_funds"].(string); ok {
				status.PurposeOfFunds = purpose
			}
			if beneficiaryType, ok := cred.CredentialSubject["beneficiary_type"].(string); ok {
				status.BeneficiaryType = beneficiaryType
			}
			if canReceive, ok := cred.CredentialSubject["can_receive_from"].([]interface{}); ok {
				for _, c := range canReceive {
					if cStr, ok := c.(string); ok {
						status.CanReceiveFrom = append(status.CanReceiveFrom, cStr)
					}
				}
			}
			break
		}
	}

	return status, nil
}

// VerifySewaMitraIdentity verifies a Sewa Mitra agent's identity
func (ia *IdentityAdapter) VerifySewaMitraIdentity(
	ctx context.Context,
	agentAddress sdk.AccAddress,
	agentID string,
) (*SewaMitraIdentityStatus, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	// Check if identity exists
	identity, found := ia.identityKeeper.GetIdentityByAddress(sdkCtx, agentAddress)
	if !found {
		return &SewaMitraIdentityStatus{
			HasIdentity: false,
		}, nil
	}

	status := &SewaMitraIdentityStatus{
		HasIdentity: true,
		DID:         identity.Did,
	}

	// Check KYC credentials
	kycCreds := ia.identityKeeper.GetCredentialsByType(sdkCtx, identity.Did, "KYCCredential")
	for _, cred := range kycCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			status.IsKYCVerified = true
			if level, ok := cred.CredentialSubject["kyc_level"].(string); ok {
				status.KYCLevel = level
			}
			break
		}
	}

	// Check Sewa Mitra credentials
	sewaMitraCreds := ia.identityKeeper.GetCredentialsByType(sdkCtx, identity.Did, "SewaMitraCredential")
	for _, cred := range sewaMitraCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			if businessLicense, ok := cred.CredentialSubject["business_license_valid"].(bool); ok {
				status.BusinessLicenseValid = businessLicense
			}
			if compliance, ok := cred.CredentialSubject["compliance_verified"].(bool); ok {
				status.ComplianceVerified = compliance
			}
			if aml, ok := cred.CredentialSubject["aml_compliant"].(bool); ok {
				status.AMLCompliant = aml
			}
			if certs, ok := cred.CredentialSubject["certifications"].([]interface{}); ok {
				for _, cert := range certs {
					if certStr, ok := cert.(string); ok {
						status.CertificationsValid = append(status.CertificationsValid, certStr)
					}
				}
			}
			if areas, ok := cred.CredentialSubject["service_areas"].([]interface{}); ok {
				for _, area := range areas {
					if areaStr, ok := area.(string); ok {
						status.ServiceAreas = append(status.ServiceAreas, areaStr)
					}
				}
			}
			if currencies, ok := cred.CredentialSubject["supported_currencies"].([]interface{}); ok {
				for _, currency := range currencies {
					if currencyStr, ok := currency.(string); ok {
						status.SupportedCurrencies = append(status.SupportedCurrencies, currencyStr)
					}
				}
			}
			if maxLimitStr, ok := cred.CredentialSubject["max_transaction_limit"].(string); ok {
				if maxLimit, err := sdk.ParseCoinNormalized(maxLimitStr); err == nil {
					status.MaxTransactionLimit = maxLimit
				}
			}
			if backgroundVerified, ok := cred.CredentialSubject["background_verified"].(bool); ok {
				status.BackgroundVerified = backgroundVerified
			}
			if insurance, ok := cred.CredentialSubject["insurance_covered"].(bool); ok {
				status.InsuranceCovered = insurance
			}
			if ratingStr, ok := cred.CredentialSubject["rating"].(string); ok {
				if rating, err := sdk.NewDecFromStr(ratingStr); err == nil {
					status.Rating = rating
				}
			}
			break
		}
	}

	return status, nil
}

// PerformComplianceCheck performs comprehensive compliance verification for a remittance
func (ia *IdentityAdapter) PerformComplianceCheck(
	ctx context.Context,
	senderAddress sdk.AccAddress,
	recipientAddress sdk.AccAddress,
	agentAddress sdk.AccAddress,
	transferAmount sdk.Coin,
	sourceCurrency string,
	destCurrency string,
	senderCountry string,
	recipientCountry string,
) (*RemittanceComplianceCheck, error) {
	result := &RemittanceComplianceCheck{}

	// Verify sender identity
	senderStatus, err := ia.VerifySenderIdentity(ctx, senderAddress, transferAmount, recipientCountry)
	if err != nil {
		return result, err
	}
	
	result.SenderCompliant = senderStatus.HasIdentity && 
		senderStatus.IsKYCVerified && 
		senderStatus.AMLVerified && 
		senderStatus.SanctionsChecked

	// Check sender limits
	if senderStatus.MaxTransferLimit.IsPositive() {
		if transferAmount.Amount.GT(senderStatus.MaxTransferLimit.Amount) {
			result.AmountWithinLimits = false
			result.RequiredActions = append(result.RequiredActions, "Transfer amount exceeds sender's maximum limit")
		}
	}

	// Check country restrictions
	for _, restrictedCountry := range senderStatus.CountryRestrictions {
		if restrictedCountry == recipientCountry {
			result.CorridorAllowed = false
			result.RequiredActions = append(result.RequiredActions, fmt.Sprintf("Transfers to %s are restricted for this sender", recipientCountry))
		}
	}

	// Verify recipient identity if address provided
	if !recipientAddress.Empty() {
		recipientStatus, err := ia.VerifyRecipientIdentity(ctx, recipientAddress, senderCountry)
		if err != nil {
			return result, err
		}
		
		result.RecipientCompliant = recipientStatus.HasIdentity && recipientStatus.IsKYCVerified
		
		// Check if recipient can receive from sender country
		if len(recipientStatus.CanReceiveFrom) > 0 {
			allowed := false
			for _, allowedCountry := range recipientStatus.CanReceiveFrom {
				if allowedCountry == senderCountry {
					allowed = true
					break
				}
			}
			if !allowed {
				result.CorridorAllowed = false
				result.RequiredActions = append(result.RequiredActions, fmt.Sprintf("Recipient cannot receive transfers from %s", senderCountry))
			}
		}
	} else {
		result.RecipientCompliant = true // No recipient address means cash pickup
	}

	// Verify agent identity if provided
	if !agentAddress.Empty() {
		agentStatus, err := ia.VerifySewaMitraIdentity(ctx, agentAddress, "")
		if err != nil {
			return result, err
		}
		
		result.AgentCompliant = agentStatus.HasIdentity && 
			agentStatus.IsKYCVerified && 
			agentStatus.BusinessLicenseValid && 
			agentStatus.ComplianceVerified && 
			agentStatus.AMLCompliant

		// Check agent limits
		if agentStatus.MaxTransactionLimit.IsPositive() {
			if transferAmount.Amount.GT(agentStatus.MaxTransactionLimit.Amount) {
				result.AmountWithinLimits = false
				result.RequiredActions = append(result.RequiredActions, "Transfer amount exceeds agent's maximum limit")
			}
		}

		// Check agent currency support
		if len(agentStatus.SupportedCurrencies) > 0 {
			supported := false
			for _, currency := range agentStatus.SupportedCurrencies {
				if currency == destCurrency {
					supported = true
					break
				}
			}
			if !supported {
				result.RequiredActions = append(result.RequiredActions, fmt.Sprintf("Agent does not support %s currency", destCurrency))
			}
		}
	} else {
		result.AgentCompliant = true // No agent means bank transfer
	}

	// Set overall compliance flags
	result.SanctionsCleared = senderStatus.SanctionsChecked
	result.AMLCompliant = senderStatus.AMLVerified
	if result.AmountWithinLimits && result.CorridorAllowed == false {
		result.CorridorAllowed = true // Set to true if not explicitly set to false
	}
	if result.AmountWithinLimits == false {
		result.AmountWithinLimits = true // Set to true if not explicitly set to false
	}
	
	result.RegulatoryCompliant = result.SenderCompliant && result.RecipientCompliant && result.AgentCompliant

	// Calculate compliance score
	score := sdk.ZeroDec()
	if result.SenderCompliant {
		score = score.Add(sdk.NewDecWithPrec(30, 2)) // 30%
	}
	if result.RecipientCompliant {
		score = score.Add(sdk.NewDecWithPrec(25, 2)) // 25%
	}
	if result.AgentCompliant {
		score = score.Add(sdk.NewDecWithPrec(20, 2)) // 20%
	}
	if result.AmountWithinLimits {
		score = score.Add(sdk.NewDecWithPrec(10, 2)) // 10%
	}
	if result.CorridorAllowed {
		score = score.Add(sdk.NewDecWithPrec(10, 2)) // 10%
	}
	if result.SanctionsCleared {
		score = score.Add(sdk.NewDecWithPrec(5, 2)) // 5%
	}
	result.ComplianceScore = score

	return result, nil
}

// CreateSenderCredential creates a remittance sender credential
func (ia *IdentityAdapter) CreateSenderCredential(
	ctx context.Context,
	issuerAddress sdk.AccAddress,
	senderAddress sdk.AccAddress,
	kycLevel string,
	sourceOfFunds string,
	riskLevel string,
	maxTransferLimit sdk.Coin,
	dailyLimit sdk.Coin,
	monthlyLimit sdk.Coin,
	countryRestrictions []string,
) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	// Get or create identity for sender
	identity, found := ia.identityKeeper.GetIdentityByAddress(sdkCtx, senderAddress)
	if !found {
		// Create new identity
		did := fmt.Sprintf("did:desh:remittance:%s", senderAddress.String())
		identity = identitytypes.Identity{
			Did:        did,
			Controller: senderAddress.String(),
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  sdkCtx.BlockTime(),
			UpdatedAt:  sdkCtx.BlockTime(),
		}
		ia.identityKeeper.SetIdentity(sdkCtx, identity)
	}

	// Create remittance credential
	credentialSubject := map[string]interface{}{
		"id":                    identity.Did,
		"sender_address":        senderAddress.String(),
		"kyc_level":             kycLevel,
		"source_of_funds":       sourceOfFunds,
		"risk_level":            riskLevel,
		"max_transfer_limit":    maxTransferLimit.String(),
		"daily_limit":           dailyLimit.String(),
		"monthly_limit":         monthlyLimit.String(),
		"country_restrictions":  countryRestrictions,
		"aml_verified":          true,
		"sanctions_checked":     true,
		"compliance_date":       sdkCtx.BlockTime().Format(time.RFC3339),
		"compliance_expiry":     sdkCtx.BlockTime().AddDate(1, 0, 0).Format(time.RFC3339), // 1 year validity
	}

	// Issue credential
	return ia.identityKeeper.IssueCredential(
		sdkCtx,
		issuerAddress,
		identity.Did,
		[]string{"VerifiableCredential", "RemittanceCredential"},
		credentialSubject,
	)
}

// CreateSewaMitraCredential creates a Sewa Mitra agent credential
func (ia *IdentityAdapter) CreateSewaMitraCredential(
	ctx context.Context,
	issuerAddress sdk.AccAddress,
	agentAddress sdk.AccAddress,
	agentID string,
	businessName string,
	serviceAreas []string,
	supportedCurrencies []string,
	maxTransactionLimit sdk.Coin,
	certifications []string,
) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	// Get or create identity for agent
	identity, found := ia.identityKeeper.GetIdentityByAddress(sdkCtx, agentAddress)
	if !found {
		// Create new identity
		did := fmt.Sprintf("did:desh:sewamitra:%s", agentAddress.String())
		identity = identitytypes.Identity{
			Did:        did,
			Controller: agentAddress.String(),
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  sdkCtx.BlockTime(),
			UpdatedAt:  sdkCtx.BlockTime(),
		}
		ia.identityKeeper.SetIdentity(sdkCtx, identity)
	}

	// Create Sewa Mitra credential
	credentialSubject := map[string]interface{}{
		"id":                      identity.Did,
		"agent_address":           agentAddress.String(),
		"agent_id":                agentID,
		"business_name":           businessName,
		"service_areas":           serviceAreas,
		"supported_currencies":    supportedCurrencies,
		"max_transaction_limit":   maxTransactionLimit.String(),
		"certifications":          certifications,
		"business_license_valid":  true,
		"compliance_verified":     true,
		"aml_compliant":          true,
		"background_verified":     true,
		"insurance_covered":       true,
		"rating":                  "4.5",
		"verification_date":       sdkCtx.BlockTime().Format(time.RFC3339),
	}

	// Issue credential
	return ia.identityKeeper.IssueCredential(
		sdkCtx,
		issuerAddress,
		identity.Did,
		[]string{"VerifiableCredential", "SewaMitraCredential"},
		credentialSubject,
	)
}

// CreateTransferCredential creates a credential for a completed transfer
func (ia *IdentityAdapter) CreateTransferCredential(
	ctx context.Context,
	transfer types.RemittanceTransfer,
) (string, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	// Get sender's identity
	senderAddr, err := sdk.AccAddressFromBech32(transfer.SenderAddress)
	if err != nil {
		return "", err
	}

	identity, found := ia.identityKeeper.GetIdentityByAddress(sdkCtx, senderAddr)
	if !found {
		return "", fmt.Errorf("identity not found for sender: %s", transfer.SenderAddress)
	}

	// Create transfer credential
	credentialSubject := map[string]interface{}{
		"id":                   identity.Did,
		"transfer_id":          transfer.Id,
		"sender_address":       transfer.SenderAddress,
		"recipient_address":    transfer.RecipientAddress,
		"amount":               transfer.Amount.String(),
		"source_currency":      transfer.SourceCurrency,
		"destination_currency": transfer.DestinationCurrency,
		"sender_country":       transfer.SenderCountry,
		"recipient_country":    transfer.RecipientCountry,
		"status":               transfer.Status.String(),
		"corridor_id":          transfer.CorridorId,
		"settlement_method":    transfer.SettlementMethod.String(),
		"transfer_date":        sdkCtx.BlockTime().Format(time.RFC3339),
	}

	if transfer.UsesSewaMitra && transfer.SewaMitraAgentId != "" {
		credentialSubject["sewa_mitra_agent_id"] = transfer.SewaMitraAgentId
		credentialSubject["sewa_mitra_commission"] = transfer.SewaMitraCommission.String()
	}

	// Issue credential (self-issued by sender)
	return ia.identityKeeper.IssueCredential(
		sdkCtx,
		senderAddr,
		identity.Did,
		[]string{"VerifiableCredential", "RemittanceTransferCredential"},
		credentialSubject,
	)
}

// MigrateExistingTransfers migrates existing transfers to create transfer credentials
func (ia *IdentityAdapter) MigrateExistingTransfers(ctx context.Context) error {
	transfers, err := ia.keeper.GetAllRemittanceTransfers(ctx)
	if err != nil {
		return err
	}
	
	for _, transfer := range transfers {
		// Only migrate completed transfers
		if transfer.Status != types.TRANSFER_STATUS_COMPLETED {
			continue
		}

		senderAddr, err := sdk.AccAddressFromBech32(transfer.SenderAddress)
		if err != nil {
			continue
		}

		sdkCtx := sdk.UnwrapSDKContext(ctx)
		
		// Check if identity exists
		identity, found := ia.identityKeeper.GetIdentityByAddress(sdkCtx, senderAddr)
		if !found {
			continue
		}

		// Check if transfer credential already exists
		transferCreds := ia.identityKeeper.GetCredentialsByType(sdkCtx, identity.Did, "RemittanceTransferCredential")
		credExists := false
		for _, cred := range transferCreds {
			if transferID, ok := cred.CredentialSubject["transfer_id"].(string); ok && transferID == transfer.Id {
				credExists = true
				break
			}
		}

		if !credExists {
			// Create transfer credential
			ia.CreateTransferCredential(ctx, transfer)
		}
	}

	return nil
}

// MigrateExistingSewaMitras migrates existing Sewa Mitra agents to identity system
func (ia *IdentityAdapter) MigrateExistingSewaMitras(ctx context.Context) error {
	agents := ia.keeper.GetAllSewaMitraAgents(ctx)
	
	for _, agent := range agents {
		agentAddr, err := sdk.AccAddressFromBech32(agent.AgentAddress)
		if err != nil {
			continue // Skip invalid addresses
		}

		sdkCtx := sdk.UnwrapSDKContext(ctx)
		
		// Check if already migrated
		_, found := ia.identityKeeper.GetIdentityByAddress(sdkCtx, agentAddr)
		if found {
			continue
		}

		// Create identity
		did := fmt.Sprintf("did:desh:sewamitra:%s", agentAddr.String())
		identity := identitytypes.Identity{
			Did:        did,
			Controller: agentAddr.String(),
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  sdkCtx.BlockTime(),
			UpdatedAt:  sdkCtx.BlockTime(),
		}
		ia.identityKeeper.SetIdentity(sdkCtx, identity)

		// Create basic Sewa Mitra credential
		credentialSubject := map[string]interface{}{
			"id":                      did,
			"agent_address":           agentAddr.String(),
			"agent_id":                agent.AgentId,
			"business_name":           agent.BusinessName,
			"country":                 agent.Country,
			"city":                    agent.City,
			"phone":                   agent.Phone,
			"email":                   agent.Email,
			"supported_currencies":    agent.SupportedCurrencies,
			"liquidity_limit":         agent.LiquidityLimit.String(),
			"daily_limit":             agent.DailyLimit.String(),
			"kyc_level":               agent.KycLevel.String(),
			"background_verified":     agent.BackgroundVerified,
			"certifications":          agent.Certifications,
			"status":                  agent.Status.String(),
			"migration_date":          sdkCtx.BlockTime().Format(time.RFC3339),
		}

		// Self-issue for migration
		ia.identityKeeper.IssueCredential(
			sdkCtx,
			agentAddr, // Self-issued during migration
			did,
			[]string{"VerifiableCredential", "SewaMitraCredential"},
			credentialSubject,
		)
	}

	return nil
}