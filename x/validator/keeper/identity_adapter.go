package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/namo/x/validator/types"
	identitykeeper "github.com/deshchain/deshchain/x/identity/keeper"
	identitytypes "github.com/deshchain/deshchain/x/identity/types"
)

// IdentityAdapter provides identity integration for Validator module
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

// ValidatorIdentityStatus represents the identity status of a validator
type ValidatorIdentityStatus struct {
	HasIdentity        bool
	DID                string
	IsKYCVerified      bool
	KYCLevel           string
	ValidatorRank      uint32
	StakeVerified      bool
	NFTBound           bool
	ReferralCredential bool
	TokenLaunched      bool
	ComplianceStatus   string
	JurisdictionCode   string
	GeographicRegion   string
}

// ValidatorComplianceStatus represents compliance verification
type ValidatorComplianceStatus struct {
	IsCompliant        bool
	ComplianceLevel    string
	JurisdictionCodes  []string
	RequiredDocuments  []string
	VerifiedDocuments  []string
	ComplianceExpiry   time.Time
	AMLVerified        bool
	SanctionsChecked   bool
	KYBCompleted       bool
}

// VerifyValidatorIdentity verifies a validator's identity and compliance
func (ia *IdentityAdapter) VerifyValidatorIdentity(
	ctx sdk.Context,
	validatorAddress sdk.AccAddress,
	operatorAddress string,
) (*ValidatorIdentityStatus, error) {
	// Check if identity exists
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, validatorAddress)
	if !found {
		return &ValidatorIdentityStatus{
			HasIdentity: false,
		}, nil
	}

	status := &ValidatorIdentityStatus{
		HasIdentity: true,
		DID:         identity.Did,
	}

	// Check KYC credentials
	kycCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "KYCCredential")
	for _, cred := range kycCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			status.IsKYCVerified = true
			if level, ok := cred.CredentialSubject["kyc_level"].(string); ok {
				status.KYCLevel = level
			}
			if jurisdiction, ok := cred.CredentialSubject["jurisdiction"].(string); ok {
				status.JurisdictionCode = jurisdiction
			}
			break
		}
	}

	// Check validator credentials
	validatorCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "ValidatorCredential")
	for _, cred := range validatorCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			if rank, ok := cred.CredentialSubject["validator_rank"].(float64); ok {
				status.ValidatorRank = uint32(rank)
			}
			if stakeVerified, ok := cred.CredentialSubject["stake_verified"].(bool); ok {
				status.StakeVerified = stakeVerified
			}
			if region, ok := cred.CredentialSubject["geographic_region"].(string); ok {
				status.GeographicRegion = region
			}
			break
		}
	}

	// Check NFT binding credentials
	nftCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "NFTBindingCredential")
	for _, cred := range nftCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			status.NFTBound = true
			break
		}
	}

	// Check referral credentials
	referralCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "ReferralCredential")
	status.ReferralCredential = len(referralCreds) > 0

	// Check token launch credentials
	tokenCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "TokenLaunchCredential")
	status.TokenLaunched = len(tokenCreds) > 0

	// Check compliance credentials
	complianceCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "ComplianceCredential")
	for _, cred := range complianceCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			if complianceStatus, ok := cred.CredentialSubject["compliance_status"].(string); ok {
				status.ComplianceStatus = complianceStatus
			}
			break
		}
	}

	return status, nil
}

// VerifyValidatorCompliance verifies comprehensive compliance for a validator
func (ia *IdentityAdapter) VerifyValidatorCompliance(
	ctx sdk.Context,
	validatorAddress sdk.AccAddress,
	requiredJurisdictions []string,
) (*ValidatorComplianceStatus, error) {
	// Get validator's identity
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, validatorAddress)
	if !found {
		return &ValidatorComplianceStatus{
			IsCompliant: false,
		}, nil
	}

	status := &ValidatorComplianceStatus{
		JurisdictionCodes: requiredJurisdictions,
	}

	// Check compliance credentials
	complianceCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "ComplianceCredential")
	for _, cred := range complianceCreds {
		if cred.Status == identitytypes.CredentialStatus_ACTIVE {
			if level, ok := cred.CredentialSubject["compliance_level"].(string); ok {
				status.ComplianceLevel = level
			}
			
			if amlVerified, ok := cred.CredentialSubject["aml_verified"].(bool); ok {
				status.AMLVerified = amlVerified
			}
			
			if sanctionsChecked, ok := cred.CredentialSubject["sanctions_checked"].(bool); ok {
				status.SanctionsChecked = sanctionsChecked
			}
			
			if kybCompleted, ok := cred.CredentialSubject["kyb_completed"].(bool); ok {
				status.KYBCompleted = kybCompleted
			}
			
			if expiryStr, ok := cred.CredentialSubject["compliance_expiry"].(string); ok {
				if expiry, err := time.Parse(time.RFC3339, expiryStr); err == nil {
					status.ComplianceExpiry = expiry
				}
			}
			
			if requiredDocs, ok := cred.CredentialSubject["required_documents"].([]interface{}); ok {
				for _, doc := range requiredDocs {
					if docStr, ok := doc.(string); ok {
						status.RequiredDocuments = append(status.RequiredDocuments, docStr)
					}
				}
			}
			
			if verifiedDocs, ok := cred.CredentialSubject["verified_documents"].([]interface{}); ok {
				for _, doc := range verifiedDocs {
					if docStr, ok := doc.(string); ok {
						status.VerifiedDocuments = append(status.VerifiedDocuments, docStr)
					}
				}
			}
			
			break
		}
	}

	// Check if all requirements are met
	status.IsCompliant = status.AMLVerified && 
		status.SanctionsChecked && 
		status.KYBCompleted &&
		(status.ComplianceExpiry.IsZero() || status.ComplianceExpiry.After(ctx.BlockTime())) &&
		len(status.VerifiedDocuments) >= len(status.RequiredDocuments)

	return status, nil
}

// CreateValidatorCredential creates a validator credential
func (ia *IdentityAdapter) CreateValidatorCredential(
	ctx sdk.Context,
	issuerAddress sdk.AccAddress,
	validatorAddress sdk.AccAddress,
	operatorAddress string,
	validatorRank uint32,
	stakeAmount sdk.Int,
) (string, error) {
	// Get or create identity for validator
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, validatorAddress)
	if !found {
		// Create new identity
		did := fmt.Sprintf("did:desh:validator:%s", validatorAddress.String())
		identity = identitytypes.Identity{
			Did:        did,
			Controller: validatorAddress.String(),
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		}
		ia.identityKeeper.SetIdentity(ctx, identity)
	}

	// Get validator stake info
	stake, found := ia.keeper.GetValidatorStake(ctx, operatorAddress)
	var stakeInfo map[string]interface{}
	if found {
		stakeInfo = map[string]interface{}{
			"namo_tokens_staked":   stake.NAMOTokensStaked.String(),
			"original_usd_value":   stake.OriginalUSDValue.String(),
			"stake_timestamp":      stake.StakeTimestamp.Format(time.RFC3339),
			"lock_end_time":        stake.LockEndTime.Format(time.RFC3339),
			"vesting_end_time":     stake.VestingEndTime.Format(time.RFC3339),
			"performance_bond":     stake.PerformanceBond.String(),
			"tier":                 stake.Tier,
		}
	}

	// Create validator credential
	credentialSubject := map[string]interface{}{
		"id":                 identity.Did,
		"validator_address":  validatorAddress.String(),
		"operator_address":   operatorAddress,
		"validator_rank":     validatorRank,
		"stake_amount":       stakeAmount.String(),
		"verification_date":  ctx.BlockTime().Format(time.RFC3339),
		"validator_type":     "genesis",
		"stake_verified":     found,
	}

	if stakeInfo != nil {
		credentialSubject["stake_info"] = stakeInfo
	}

	// Add geographic information if available
	if validatorRank > 0 && validatorRank <= 21 {
		tier, found := types.GetTierForRank(validatorRank)
		if found {
			credentialSubject["tier_id"] = tier.TierID
			credentialSubject["lock_period_months"] = tier.LockPeriodMonths
			credentialSubject["vesting_months"] = tier.VestingMonths
		}
	}

	// Issue credential
	return ia.identityKeeper.IssueCredential(
		ctx,
		issuerAddress,
		identity.Did,
		[]string{"VerifiableCredential", "ValidatorCredential"},
		credentialSubject,
	)
}

// CreateNFTBindingCredential creates an NFT binding credential
func (ia *IdentityAdapter) CreateNFTBindingCredential(
	ctx sdk.Context,
	issuerAddress sdk.AccAddress,
	validatorAddress sdk.AccAddress,
	nftTokenID uint64,
) (string, error) {
	// Get identity
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, validatorAddress)
	if !found {
		return "", fmt.Errorf("identity not found for validator: %s", validatorAddress.String())
	}

	// Get NFT info
	nft, found := ia.keeper.GetGenesisNFT(ctx, nftTokenID)
	var nftInfo map[string]interface{}
	if found {
		nftInfo = map[string]interface{}{
			"token_id":         nft.TokenID,
			"validator_rank":   nft.ValidatorRank,
			"initial_value":    nft.InitialValue.String(),
			"geographic_type":  nft.GeographicType,
			"cultural_theme":   nft.CulturalTheme,
			"special_rights":   nft.SpecialRights,
		}
	}

	// Create NFT binding credential
	credentialSubject := map[string]interface{}{
		"id":               identity.Did,
		"validator_address": validatorAddress.String(),
		"nft_token_id":     nftTokenID,
		"binding_date":     ctx.BlockTime().Format(time.RFC3339),
		"binding_active":   true,
	}

	if nftInfo != nil {
		credentialSubject["nft_info"] = nftInfo
	}

	// Issue credential
	return ia.identityKeeper.IssueCredential(
		ctx,
		issuerAddress,
		identity.Did,
		[]string{"VerifiableCredential", "NFTBindingCredential"},
		credentialSubject,
	)
}

// CreateReferralCredential creates a referral credential
func (ia *IdentityAdapter) CreateReferralCredential(
	ctx sdk.Context,
	issuerAddress sdk.AccAddress,
	referrerAddress sdk.AccAddress,
	referredAddress sdk.AccAddress,
	referralID uint64,
) (string, error) {
	// Get referrer's identity
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, referrerAddress)
	if !found {
		return "", fmt.Errorf("identity not found for referrer: %s", referrerAddress.String())
	}

	// Get referral info
	referral, found := ia.keeper.GetReferral(ctx, referralID)
	if !found {
		return "", fmt.Errorf("referral not found: %d", referralID)
	}

	// Create referral credential
	credentialSubject := map[string]interface{}{
		"id":                identity.Did,
		"referrer_address":  referrerAddress.String(),
		"referred_address":  referredAddress.String(),
		"referral_id":       referralID,
		"referred_rank":     referral.ReferredRank,
		"commission_rate":   referral.CommissionRate.String(),
		"total_commission":  referral.TotalCommission.String(),
		"referral_date":     referral.CreatedAt.Format(time.RFC3339),
		"status":            string(referral.Status),
	}

	// Issue credential
	return ia.identityKeeper.IssueCredential(
		ctx,
		issuerAddress,
		identity.Did,
		[]string{"VerifiableCredential", "ReferralCredential"},
		credentialSubject,
	)
}

// CreateTokenLaunchCredential creates a token launch credential
func (ia *IdentityAdapter) CreateTokenLaunchCredential(
	ctx sdk.Context,
	issuerAddress sdk.AccAddress,
	validatorAddress sdk.AccAddress,
	tokenID uint64,
) (string, error) {
	// Get identity
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, validatorAddress)
	if !found {
		return "", fmt.Errorf("identity not found for validator: %s", validatorAddress.String())
	}

	// Get token info
	token, found := ia.keeper.GetValidatorToken(ctx, tokenID)
	if !found {
		return "", fmt.Errorf("validator token not found: %d", tokenID)
	}

	// Create token launch credential
	credentialSubject := map[string]interface{}{
		"id":                   identity.Did,
		"validator_address":    validatorAddress.String(),
		"token_id":             tokenID,
		"token_name":           token.TokenName,
		"token_symbol":         token.TokenSymbol,
		"total_supply":         token.TotalSupply.String(),
		"launch_date":          token.LaunchedAt.Format(time.RFC3339),
		"launch_trigger":       token.LaunchTrigger,
		"validator_allocation": token.ValidatorAllocation.String(),
		"liquidity_allocation": token.LiquidityAllocation.String(),
		"current_price":        token.CurrentPrice.String(),
		"market_cap":           token.MarketCap.String(),
	}

	// Issue credential
	return ia.identityKeeper.IssueCredential(
		ctx,
		issuerAddress,
		identity.Did,
		[]string{"VerifiableCredential", "TokenLaunchCredential"},
		credentialSubject,
	)
}

// CreateComplianceCredential creates a compliance credential
func (ia *IdentityAdapter) CreateComplianceCredential(
	ctx sdk.Context,
	issuerAddress sdk.AccAddress,
	validatorAddress sdk.AccAddress,
	complianceLevel string,
	jurisdictions []string,
	requiredDocuments []string,
	verifiedDocuments []string,
) (string, error) {
	// Get identity
	identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, validatorAddress)
	if !found {
		return "", fmt.Errorf("identity not found for validator: %s", validatorAddress.String())
	}

	// Create compliance credential
	credentialSubject := map[string]interface{}{
		"id":                   identity.Did,
		"validator_address":    validatorAddress.String(),
		"compliance_level":     complianceLevel,
		"jurisdictions":        jurisdictions,
		"required_documents":   requiredDocuments,
		"verified_documents":   verifiedDocuments,
		"aml_verified":         len(verifiedDocuments) > 0,
		"sanctions_checked":    true,
		"kyb_completed":        true,
		"compliance_date":      ctx.BlockTime().Format(time.RFC3339),
		"compliance_expiry":    ctx.BlockTime().AddDate(1, 0, 0).Format(time.RFC3339), // 1 year validity
	}

	// Issue credential
	return ia.identityKeeper.IssueCredential(
		ctx,
		issuerAddress,
		identity.Did,
		[]string{"VerifiableCredential", "ComplianceCredential"},
		credentialSubject,
	)
}

// MigrateExistingValidators migrates existing validators to identity system
func (ia *IdentityAdapter) MigrateExistingValidators(ctx sdk.Context) error {
	validators := ia.keeper.GetAllActiveValidators(ctx)
	
	for _, validator := range validators {
		validatorAddr, err := sdk.AccAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			continue // Skip invalid addresses
		}

		// Check if already migrated
		_, found := ia.identityKeeper.GetIdentityByAddress(ctx, validatorAddr)
		if found {
			continue
		}

		// Create identity
		did := fmt.Sprintf("did:desh:validator:%s", validatorAddr.String())
		identity := identitytypes.Identity{
			Did:        did,
			Controller: validatorAddr.String(),
			Status:     identitytypes.IdentityStatus_ACTIVE,
			CreatedAt:  ctx.BlockTime(),
			UpdatedAt:  ctx.BlockTime(),
		}
		ia.identityKeeper.SetIdentity(ctx, identity)

		// Create basic validator credential
		credentialSubject := map[string]interface{}{
			"id":                did,
			"validator_address": validatorAddr.String(),
			"operator_address":  validator.OperatorAddress,
			"tokens":            validator.Tokens.String(),
			"status":            validator.Status.String(),
			"join_order":        validator.JoinOrder,
			"migration_date":    ctx.BlockTime().Format(time.RFC3339),
		}

		// Self-issue for migration
		ia.identityKeeper.IssueCredential(
			ctx,
			validatorAddr, // Self-issued during migration
			did,
			[]string{"VerifiableCredential", "ValidatorCredential"},
			credentialSubject,
		)
	}

	return nil
}

// MigrateExistingReferrals migrates existing referrals to create referral credentials
func (ia *IdentityAdapter) MigrateExistingReferrals(ctx sdk.Context) error {
	referrals := ia.keeper.GetAllReferrals(ctx)
	
	for _, referral := range referrals {
		if referral.Status != types.ReferralStatusActive {
			continue
		}

		referrerAddr, err := sdk.AccAddressFromBech32(referral.ReferrerAddr)
		if err != nil {
			continue
		}

		referredAddr, err := sdk.AccAddressFromBech32(referral.ReferredAddr)
		if err != nil {
			continue
		}

		// Check if identity exists for referrer
		identity, found := ia.identityKeeper.GetIdentityByAddress(ctx, referrerAddr)
		if !found {
			continue
		}

		// Check if referral credential already exists
		referralCreds := ia.identityKeeper.GetCredentialsByType(ctx, identity.Did, "ReferralCredential")
		credExists := false
		for _, cred := range referralCreds {
			if referralIDFloat, ok := cred.CredentialSubject["referral_id"].(float64); ok {
				if uint64(referralIDFloat) == referral.ReferralID {
					credExists = true
					break
				}
			}
		}

		if !credExists {
			// Create referral credential (self-issued during migration)
			ia.CreateReferralCredential(ctx, referrerAddr, referrerAddr, referredAddr, referral.ReferralID)
		}
	}

	return nil
}