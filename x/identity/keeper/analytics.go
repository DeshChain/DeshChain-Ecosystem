package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/deshchain/x/identity/types"
)

// UsageMetrics represents usage metrics for analytics
type UsageMetrics struct {
	Date                    string `json:"date"`
	ActiveIdentities        int64  `json:"active_identities"`
	NewIdentities           int64  `json:"new_identities"`
	CredentialsIssued       int64  `json:"credentials_issued"`
	CredentialsRevoked      int64  `json:"credentials_revoked"`
	CredentialsPresented    int64  `json:"credentials_presented"`
	ZKProofsCreated         int64  `json:"zk_proofs_created"`
	ZKProofsVerified        int64  `json:"zk_proofs_verified"`
	ConsentsGiven           int64  `json:"consents_given"`
	ConsentsWithdrawn       int64  `json:"consents_withdrawn"`
	RecoveriesInitiated     int64  `json:"recoveries_initiated"`
	RecoveriesCompleted     int64  `json:"recoveries_completed"`
	AadhaarVerifications    int64  `json:"aadhaar_verifications"`
	DigiLockerConnections   int64  `json:"digilocker_connections"`
	UPILinkages             int64  `json:"upi_linkages"`
}

// SetUsageMetrics stores usage metrics
func (k Keeper) SetUsageMetrics(ctx sdk.Context, date string, metrics UsageMetrics) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UsageMetricsPrefix)
	b := k.cdc.MustMarshal(&metrics)
	store.Set([]byte(date), b)
}

// GetUsageMetrics retrieves usage metrics for a date
func (k Keeper) GetUsageMetrics(ctx sdk.Context, date string) (UsageMetrics, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UsageMetricsPrefix)
	b := store.Get([]byte(date))
	if b == nil {
		return UsageMetrics{Date: date}, false
	}
	
	var metrics UsageMetrics
	k.cdc.MustUnmarshal(b, &metrics)
	return metrics, true
}

// UpdateDailyActiveIdentities updates the daily active identities count
func (k Keeper) UpdateDailyActiveIdentities(ctx sdk.Context) {
	date := ctx.BlockTime().Format("2006-01-02")
	metrics, _ := k.GetUsageMetrics(ctx, date)
	
	// Count active identities (had activity today)
	activeCount := int64(0)
	k.IterateIdentities(ctx, func(identity types.Identity) bool {
		if identity.LastActivityAt.Format("2006-01-02") == date {
			activeCount++
		}
		return false
	})
	
	metrics.ActiveIdentities = activeCount
	k.SetUsageMetrics(ctx, date, metrics)
}

// UpdateCredentialMetrics updates credential-related metrics
func (k Keeper) UpdateCredentialMetrics(ctx sdk.Context) {
	date := ctx.BlockTime().Format("2006-01-02")
	metrics, _ := k.GetUsageMetrics(ctx, date)
	
	// These would be incremented during actual operations
	// For now, we'll just ensure the structure exists
	k.SetUsageMetrics(ctx, date, metrics)
}

// UpdateVerificationMetrics updates verification-related metrics
func (k Keeper) UpdateVerificationMetrics(ctx sdk.Context) {
	date := ctx.BlockTime().Format("2006-01-02")
	metrics, _ := k.GetUsageMetrics(ctx, date)
	
	// These would be incremented during actual operations
	k.SetUsageMetrics(ctx, date, metrics)
}

// IncrementMetric increments a specific metric
func (k Keeper) IncrementMetric(ctx sdk.Context, metricType string) {
	date := ctx.BlockTime().Format("2006-01-02")
	metrics, _ := k.GetUsageMetrics(ctx, date)
	
	switch metricType {
	case "new_identities":
		metrics.NewIdentities++
	case "credentials_issued":
		metrics.CredentialsIssued++
	case "credentials_revoked":
		metrics.CredentialsRevoked++
	case "credentials_presented":
		metrics.CredentialsPresented++
	case "zk_proofs_created":
		metrics.ZKProofsCreated++
	case "zk_proofs_verified":
		metrics.ZKProofsVerified++
	case "consents_given":
		metrics.ConsentsGiven++
	case "consents_withdrawn":
		metrics.ConsentsWithdrawn++
	case "recoveries_initiated":
		metrics.RecoveriesInitiated++
	case "recoveries_completed":
		metrics.RecoveriesCompleted++
	case "aadhaar_verifications":
		metrics.AadhaarVerifications++
	case "digilocker_connections":
		metrics.DigiLockerConnections++
	case "upi_linkages":
		metrics.UPILinkages++
	}
	
	k.SetUsageMetrics(ctx, date, metrics)
}

// GetMetricsForPeriod returns metrics for a date range
func (k Keeper) GetMetricsForPeriod(ctx sdk.Context, startDate, endDate string) []UsageMetrics {
	var metricsSlice []UsageMetrics
	
	// Parse dates
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	
	// Iterate through date range
	for date := start; !date.After(end); date = date.AddDate(0, 0, 1) {
		dateStr := date.Format("2006-01-02")
		if metrics, found := k.GetUsageMetrics(ctx, dateStr); found {
			metricsSlice = append(metricsSlice, metrics)
		}
	}
	
	return metricsSlice
}

// GetIdentityStats returns statistics about identities
func (k Keeper) GetIdentityStats(ctx sdk.Context) map[string]interface{} {
	stats := make(map[string]interface{})
	
	totalIdentities := int64(0)
	activeIdentities := int64(0)
	verifiedIdentities := int64(0)
	identitiesWithBiometrics := int64(0)
	identitiesWithRecovery := int64(0)
	
	k.IterateIdentities(ctx, func(identity types.Identity) bool {
		totalIdentities++
		
		if identity.IsActive() {
			activeIdentities++
		}
		
		if identity.KYCStatus.Status == types.VerificationStatus_VERIFIED {
			verifiedIdentities++
		}
		
		if identity.HasBiometrics() {
			identitiesWithBiometrics++
		}
		
		if len(identity.RecoveryMethods) > 0 {
			identitiesWithRecovery++
		}
		
		return false
	})
	
	stats["total_identities"] = totalIdentities
	stats["active_identities"] = activeIdentities
	stats["verified_identities"] = verifiedIdentities
	stats["identities_with_biometrics"] = identitiesWithBiometrics
	stats["identities_with_recovery"] = identitiesWithRecovery
	
	return stats
}

// GetCredentialStats returns statistics about credentials
func (k Keeper) GetCredentialStats(ctx sdk.Context) map[string]interface{} {
	stats := make(map[string]interface{})
	
	totalCredentials := int64(0)
	activeCredentials := int64(0)
	expiredCredentials := int64(0)
	revokedCredentials := int64(0)
	credentialsByType := make(map[string]int64)
	
	k.IterateCredentials(ctx, func(credential types.VerifiableCredential) bool {
		totalCredentials++
		
		if credential.IsExpired() {
			expiredCredentials++
		} else if k.IsCredentialRevoked(ctx, credential.ID) {
			revokedCredentials++
		} else {
			activeCredentials++
		}
		
		// Count by type
		for _, credType := range credential.Type {
			if credType != types.CredentialTypeVerifiable {
				credentialsByType[credType]++
			}
		}
		
		return false
	})
	
	stats["total_credentials"] = totalCredentials
	stats["active_credentials"] = activeCredentials
	stats["expired_credentials"] = expiredCredentials
	stats["revoked_credentials"] = revokedCredentials
	stats["credentials_by_type"] = credentialsByType
	
	return stats
}

// GetIndiaStackStats returns statistics about India Stack integrations
func (k Keeper) GetIndiaStackStats(ctx sdk.Context) map[string]interface{} {
	stats := make(map[string]interface{})
	
	aadhaarLinked := int64(0)
	digilockerConnected := int64(0)
	upiLinked := int64(0)
	panchayatKYCs := int64(0)
	
	k.IterateIndiaStackIntegrations(ctx, func(integration types.IndiaStackIntegration) bool {
		if integration.AadhaarLinked {
			aadhaarLinked++
		}
		if integration.DigiLockerLinked {
			digilockerConnected++
		}
		if integration.UPILinked {
			upiLinked++
		}
		if len(integration.VillagePanchayatKYCs) > 0 {
			panchayatKYCs++
		}
		return false
	})
	
	stats["aadhaar_linked"] = aadhaarLinked
	stats["digilocker_connected"] = digilockerConnected
	stats["upi_linked"] = upiLinked
	stats["panchayat_kycs"] = panchayatKYCs
	
	return stats
}

// GetPrivacyStats returns statistics about privacy features usage
func (k Keeper) GetPrivacyStats(ctx sdk.Context) map[string]interface{} {
	stats := make(map[string]interface{})
	
	zkProofsActive := int64(0)
	nullifiersUsed := int64(0)
	anonymousCredentials := int64(0)
	customPrivacySettings := int64(0)
	
	// Count active ZK proofs
	k.IterateZKProofs(ctx, func(proof types.ZKProof) bool {
		if proof.ExpiresAt == nil || proof.ExpiresAt.After(ctx.BlockTime()) {
			zkProofsActive++
		}
		return false
	})
	
	// Count nullifiers
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.NullifierPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		nullifiersUsed++
	}
	
	// Count anonymous credentials
	anonStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.AnonymousCredPrefix)
	anonIterator := sdk.KVStorePrefixIterator(anonStore, nil)
	defer anonIterator.Close()
	for ; anonIterator.Valid(); anonIterator.Next() {
		anonymousCredentials++
	}
	
	// Count custom privacy settings
	k.IteratePrivacySettings(ctx, func(settings types.PrivacySettings) bool {
		// Check if settings differ from defaults
		if settings.DefaultDisclosureLevel != types.DisclosureLevel_STANDARD ||
			settings.AllowAnonymousUsage != false ||
			settings.RequireExplicitConsent != true {
			customPrivacySettings++
		}
		return false
	})
	
	stats["zk_proofs_active"] = zkProofsActive
	stats["nullifiers_used"] = nullifiersUsed
	stats["anonymous_credentials"] = anonymousCredentials
	stats["custom_privacy_settings"] = customPrivacySettings
	
	return stats
}

// GenerateAnalyticsReport generates a comprehensive analytics report
func (k Keeper) GenerateAnalyticsReport(ctx sdk.Context) map[string]interface{} {
	report := make(map[string]interface{})
	
	// Add timestamp
	report["generated_at"] = ctx.BlockTime().Format(time.RFC3339)
	report["block_height"] = ctx.BlockHeight()
	
	// Add various stats
	report["identity_stats"] = k.GetIdentityStats(ctx)
	report["credential_stats"] = k.GetCredentialStats(ctx)
	report["india_stack_stats"] = k.GetIndiaStackStats(ctx)
	report["privacy_stats"] = k.GetPrivacyStats(ctx)
	
	// Add today's metrics
	todayMetrics, _ := k.GetUsageMetrics(ctx, ctx.BlockTime().Format("2006-01-02"))
	report["today_metrics"] = todayMetrics
	
	// Add module parameters
	params := k.GetParams(ctx)
	report["module_params"] = map[string]interface{}{
		"max_credentials_per_identity": params.MaxCredentialsPerIdentity,
		"credential_expiry_days":       params.CredentialExpiryDays,
		"kyc_expiry_days":              params.KYCExpiryDays,
		"enable_anonymous_credentials": params.EnableAnonymousCredentials,
		"enable_zk_proofs":             params.EnableZKProofs,
		"enable_india_stack":           params.EnableIndiaStack,
	}
	
	return report
}

// LogIdentityOperation logs an identity operation for analytics
func (k Keeper) LogIdentityOperation(ctx sdk.Context, operation string, address string, details map[string]string) {
	// Emit event for monitoring
	event := sdk.NewEvent(
		"identity_operation",
		sdk.NewAttribute("operation", operation),
		sdk.NewAttribute("address", address),
		sdk.NewAttribute("timestamp", ctx.BlockTime().Format(time.RFC3339)),
	)
	
	// Add details as attributes
	for key, value := range details {
		event = event.AppendAttributes(sdk.NewAttribute(key, value))
	}
	
	ctx.EventManager().EmitEvent(event)
	
	// Increment relevant metric
	switch operation {
	case "create_identity":
		k.IncrementMetric(ctx, "new_identities")
	case "issue_credential":
		k.IncrementMetric(ctx, "credentials_issued")
	case "revoke_credential":
		k.IncrementMetric(ctx, "credentials_revoked")
	case "present_credential":
		k.IncrementMetric(ctx, "credentials_presented")
	}
}