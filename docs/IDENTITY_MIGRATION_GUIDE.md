# DeshChain Identity Module Migration Guide

## Overview

This guide provides step-by-step instructions for migrating existing KYC and biometric systems to the new DeshChain Identity Module while maintaining backward compatibility.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Migration Strategy](#migration-strategy)
3. [TradeFinance KYC Migration](#tradefinance-kyc-migration)
4. [MoneyOrder Biometric Migration](#moneyorder-biometric-migration)
5. [Data Migration](#data-migration)
6. [API Changes](#api-changes)
7. [Testing Strategy](#testing-strategy)
8. [Rollback Plan](#rollback-plan)

## Architecture Overview

The DeshChain Identity Module introduces:
- **W3C DID (Decentralized Identifiers)** support
- **Verifiable Credentials** infrastructure
- **Zero-Knowledge Proofs** for privacy
- **India Stack Integration** (Aadhaar, DigiLocker, UPI)
- **Unified Identity Management** across all modules

### Key Components

```
Identity Module
├── DID Management
│   ├── W3C DID Documents
│   ├── Verification Methods
│   └── Service Endpoints
├── Credential System
│   ├── Verifiable Credentials
│   ├── Selective Disclosure
│   └── Revocation Registry
├── Privacy Layer
│   ├── Zero-Knowledge Proofs
│   ├── Anonymous Credentials
│   └── Consent Management
└── Integration Adapters
    ├── TradeFinance Integration
    └── Biometric Integration
```

## Migration Strategy

### Phase 1: Parallel Operation (Recommended)
- Deploy identity module alongside existing systems
- Use integration adapters for dual-write operations
- Gradually migrate read operations to identity module
- Monitor both systems for consistency

### Phase 2: Data Migration
- Migrate historical KYC profiles to identity credentials
- Convert biometric registrations to verifiable credentials
- Maintain references between old and new identifiers

### Phase 3: Cutover
- Switch primary operations to identity module
- Keep legacy system in read-only mode
- Complete final data reconciliation

## TradeFinance KYC Migration

### 1. Update Module Dependencies

```go
// app/app.go
import (
    identitykeeper "github.com/namo/x/identity/keeper"
    identitytypes "github.com/namo/x/identity/types"
)

// Add identity keeper to app
type App struct {
    // ... existing keepers
    IdentityKeeper identitykeeper.Keeper
}
```

### 2. Initialize Identity Adapter

```go
// x/tradefinance/keeper/keeper.go
type Keeper struct {
    // ... existing fields
    identityAdapter *IdentityAdapter
}

// In NewKeeper
func NewKeeper(..., identityKeeper identitykeeper.Keeper) Keeper {
    k := Keeper{
        // ... existing initialization
    }
    k.identityAdapter = NewIdentityAdapter(&k, identityKeeper)
    return k
}
```

### 3. Update KYC Operations

Replace direct KYC calls with adapter methods:

```go
// Before
result, err := k.PerformCustomerKYC(ctx, customerID, submittedData)

// After
result, err := k.identityAdapter.PerformKYCWithIdentity(ctx, customerID, submittedData)
```

### 4. Implement Gradual Migration

```go
// Use feature flag for gradual rollout
if k.useIdentityModule {
    // New identity-based flow
    status, err := k.identityAdapter.VerifyCustomerIdentity(ctx, customerID)
} else {
    // Legacy flow
    profile, err := k.GetKYCProfile(ctx, customerID)
}
```

### 5. Migrate Existing Data

```go
// One-time migration script
func MigrateKYCProfiles(ctx sdk.Context, k Keeper) error {
    migrated, err := k.identityAdapter.MigrateKYCToIdentity(ctx)
    if err != nil {
        return err
    }
    k.Logger(ctx).Info("KYC migration completed", "profiles_migrated", migrated)
    return nil
}
```

## MoneyOrder Biometric Migration

### 1. Update Biometric Manager

```go
// x/moneyorder/keeper/keeper.go
type Keeper struct {
    // ... existing fields
    biometricAdapter *IdentityBiometricAdapter
}

// Initialize adapter
func NewKeeper(..., identityKeeper identitykeeper.Keeper) Keeper {
    k := Keeper{
        // ... existing initialization
    }
    k.biometricAdapter = NewIdentityBiometricAdapter(&k, identityKeeper)
    return k
}
```

### 2. Update Registration Flow

```go
// Before
err := k.biometricMgr.RegisterBiometric(ctx, userAddress, biometricType, templateHash, deviceID)

// After
err := k.biometricAdapter.RegisterBiometricWithIdentity(ctx, userAddress, biometricType, templateHash, deviceID)
```

### 3. Enhanced Authentication

```go
// Use enhanced authentication for high-value transactions
if amount.IsGTE(highValueThreshold) {
    verified, err := k.biometricAdapter.VerifyBiometricForHighValue(
        ctx, userAddress, amount, biometricData,
    )
    if !verified {
        return ErrBiometricRequired
    }
}
```

### 4. Status Checking

```go
// Get comprehensive biometric status
status, err := k.biometricAdapter.GetBiometricStatusWithIdentity(ctx, userAddress)
if err != nil {
    return err
}

// Check preferred system
if status.PreferredSystem == "identity" {
    // Use identity-based verification
} else {
    // Use traditional verification
}
```

## Data Migration

### 1. Export Existing Data

```bash
# Export KYC profiles
deshchaind query tradefinance export-kyc --output json > kyc_profiles.json

# Export biometric registrations  
deshchaind query moneyorder export-biometrics --output json > biometrics.json
```

### 2. Run Migration Tool

```bash
# Migrate KYC profiles
deshchaind tx identity migrate-kyc kyc_profiles.json --from admin

# Migrate biometrics
deshchaind tx identity migrate-biometrics biometrics.json --from admin
```

### 3. Verify Migration

```bash
# Check migration status
deshchaind query identity migration-status

# Verify specific customer
deshchaind query identity did did:desh:<customer_id>
```

## API Changes

### New Identity Endpoints

```bash
# Query identity
GET /deshchain/identity/v1/did/{did}

# Query credentials
GET /deshchain/identity/v1/credentials/{subject}

# Verify credential
POST /deshchain/identity/v1/verify
```

### Updated TradeFinance Endpoints

```bash
# Now returns both traditional and DID-based status
GET /deshchain/tradefinance/v1/kyc/{customer_id}
{
    "traditional_kyc": {...},
    "identity_did": "did:desh:...",
    "credentials": [...]
}
```

### Client SDK Updates

```typescript
// New identity-aware client
import { DeshChainClient } from '@deshchain/sdk';

const client = new DeshChainClient({
    useIdentity: true, // Enable identity module
    backwardCompatible: true // Support legacy operations
});

// Verify KYC with identity
const status = await client.tradefinance.verifyKYCWithIdentity(customerId);
```

## Testing Strategy

### 1. Unit Tests

```go
func TestKYCMigration(t *testing.T) {
    // Test dual-write operations
    result := adapter.PerformKYCWithIdentity(ctx, customerID, data)
    
    // Verify both systems have data
    assert.True(t, hasTraditionalKYC(customerID))
    assert.True(t, hasIdentityCredential(customerID))
}
```

### 2. Integration Tests

```bash
# Run integration test suite
make test-identity-integration

# Test specific scenarios
go test ./x/identity/integration/... -run TestKYCMigration
```

### 3. Load Testing

```bash
# Test performance with dual operations
deshchaind test identity-load \
  --concurrent-users 1000 \
  --operations-per-user 10 \
  --include-legacy
```

### 4. Data Validation

```sql
-- Verify data consistency
SELECT 
    COUNT(*) as total_customers,
    SUM(CASE WHEN has_traditional_kyc THEN 1 ELSE 0 END) as traditional_count,
    SUM(CASE WHEN has_identity_did THEN 1 ELSE 0 END) as identity_count
FROM customer_migration_status;
```

## Rollback Plan

### 1. Feature Flags

```go
// Disable identity integration
params.UseIdentityModule = false
params.UseIdentityBiometrics = false
```

### 2. Revert Adapters

```go
// Switch back to direct calls
if !params.UseIdentityModule {
    return k.legacyKYCEngine.PerformKYC(ctx, customerID, data)
}
```

### 3. Data Recovery

```bash
# Export identity data before rollback
deshchaind query identity export-all --height <pre_migration_height>

# Restore if needed
deshchaind tx identity restore-state backup.json --from admin
```

## Best Practices

1. **Gradual Rollout**
   - Start with read operations
   - Move to dual-write
   - Finally switch primary operations

2. **Monitoring**
   - Track both systems' performance
   - Monitor data consistency
   - Set up alerts for discrepancies

3. **Communication**
   - Notify users of new features
   - Provide migration timeline
   - Offer support during transition

4. **Backward Compatibility**
   - Maintain old APIs during transition
   - Support legacy data formats
   - Provide clear deprecation timeline

## Common Issues and Solutions

### Issue 1: Duplicate Identities
```go
// Check for existing DID before creating
if _, exists := k.GetIdentity(ctx, did); exists {
    return existingDID, nil
}
```

### Issue 2: Credential Expiry
```go
// Set appropriate expiry based on KYC level
expiryDays := map[string]uint32{
    "basic": 90,
    "standard": 180,
    "enhanced": 365,
}[kycLevel]
```

### Issue 3: Performance Impact
```go
// Use caching for frequently accessed identities
if cached, found := k.identityCache.Get(did); found {
    return cached.(*Identity), nil
}
```

## Migration Timeline

| Phase | Duration | Activities |
|-------|----------|------------|
| Preparation | 2 weeks | Deploy identity module, update adapters |
| Pilot | 1 week | Test with small user group |
| Gradual Rollout | 4 weeks | Increase usage percentage weekly |
| Full Migration | 2 weeks | Complete data migration |
| Cleanup | 1 week | Remove legacy code, finalize |

## Support and Resources

- **Documentation**: [Identity Module Docs](./modules/IDENTITY_MODULE.md)
- **API Reference**: [Identity API](./api/identity.md)
- **Examples**: [Integration Examples](./examples/identity/)
- **Support**: identity-support@deshchain.com

## Conclusion

The migration to DeshChain's Identity Module provides:
- Enhanced security through DID and Verifiable Credentials
- Better privacy with Zero-Knowledge Proofs
- Unified identity management across all modules
- Future-proof architecture for regulatory compliance

Follow this guide carefully and reach out to the support team for any assistance during migration.