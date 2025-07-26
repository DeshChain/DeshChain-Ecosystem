# DeshChain Identity Offline Verification System

## Overview

The DeshChain Identity module includes a comprehensive offline verification system that enables identity verification without active network connectivity. This is crucial for rural areas, disaster scenarios, and environments with limited internet access.

## Features

### 🌐 Five Format Support
- **Self-Contained**: Complete offline package with all verification data
- **Compressed**: Space-efficient format for bandwidth-limited scenarios
- **QR Code**: Visual scanning compatible format for mobile apps
- **NFC**: Near Field Communication compatible for smart cards and IoT
- **Printable**: Human-readable format for paper-based verification

### 🛡️ Four Verification Modes
- **Full**: Complete verification with all security checks (Security Level 5)
- **Partial**: Cached verification with reduced security (Security Level 3)
- **Minimal**: Basic identity verification (Security Level 2)
- **Emergency**: Crisis mode with relaxed security (Security Level 1)

### 🔒 Security Features
- End-to-end encryption for all offline packages
- Cryptographic signatures for data integrity
- Biometric template encryption with secure storage
- Tamper detection and validation
- Expiration-based security with configurable timeouts

## Architecture

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                 Offline Verification System                │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────── │
│  │ Package Creator │  │ Format Manager  │  │ Device Manager │
│  └─────────────────┘  └─────────────────┘  └─────────────── │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────── │
│  │ Verification    │  │ Biometric       │  │ Localization  │
│  │ Engine          │  │ Processor       │  │ Manager       │
│  └─────────────────┘  └─────────────────┘  └─────────────── │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────│
│  │              Cryptographic Security Layer              │
│  └─────────────────────────────────────────────────────────│
└─────────────────────────────────────────────────────────────┘
```

### Data Structure

```go
type OfflineVerificationData struct {
    DID               string                         // Identity DID
    IdentityHash      string                         // Hash of identity data
    PublicKey         string                         // Public key for verification
    IssuerDID         string                         // Issuer's DID
    SecurityLevel     uint32                         // Security level (1-5)
    IssuedAt          time.Time                      // Issue timestamp
    ExpiresAt         time.Time                      // Expiration timestamp
    Format            OfflineCredentialFormat        // Package format
    
    // Cryptographic proofs
    IdentityProof     *CryptographicProof           // Identity verification proof
    KYCProof          *CryptographicProof           // KYC verification proof
    
    // Credentials and data
    Credentials       []*OfflineCredential          // Verifiable credentials
    BiometricTemplates map[string]*EncryptedBiometric // Encrypted biometric data
    
    // Revocation and security
    RevocationData    *RevocationData               // Revocation list data
    EmergencyContacts []*EmergencyContact           // Emergency contact information
    
    // Localization support
    LocalizedData     *OfflineLocalizationData      // Multi-language support
    
    // Integrity verification
    DataHash          string                        // Data integrity hash
    Signature         string                        // Digital signature
}
```

## Usage Guide

### 1. Preparing Offline Verification Data

#### CLI Command
```bash
# Basic offline package
deshchaind tx identity offline prepare-verification \
  did:desh:user123 \
  self_contained \
  24 \
  --from mykey

# Advanced package with biometrics
deshchaind tx identity offline prepare-verification \
  did:desh:user456 \
  compressed \
  72 \
  --include-biometric \
  --include-credentials="IdentityCredential,EducationCredential" \
  --required-level=3 \
  --device-id="mobile_device_001" \
  --from mykey
```

#### Go SDK Usage
```go
import "github.com/deshchain/deshchain/x/identity/types"

// Create offline verification data
offlineData := types.NewOfflineVerificationData(
    "did:desh:user123",
    "identity_hash_abc123",
    "public_key_def456",
    "did:desh:issuer",
    3, // Security level
)

// Add credentials
credential := types.NewOfflineCredential(
    "credential_001",
    []string{"VerifiableCredential", "IdentityCredential"},
    "did:desh:user123",
    credentialSubject,
)
offlineData.Credentials = []*types.OfflineCredential{credential}

// Set expiration (24 hours)
offlineData.ExpiresAt = time.Now().Add(24 * time.Hour)
```

### 2. Device Registration

#### Register Device for Offline Verification
```bash
deshchaind tx identity offline register-device \
  did:desh:user123 \
  "device_001" \
  "John's iPhone" \
  mobile \
  "pubkey_abc123" \
  --capabilities="identity_verification,biometric_capture" \
  --security-level=3 \
  --max-offline-hours=48 \
  --from mykey
```

#### Query Registered Devices
```bash
# List all devices for a DID
deshchaind query identity offline devices did:desh:user123

# Query specific device
deshchaind query identity offline devices did:desh:user123 --device-id="device_001"
```

### 3. Offline Verification Process

#### Verify Identity Offline
```bash
# Verify using offline data and verification request
deshchaind query identity offline verify \
  offline_data.json \
  verification_request.json \
  --output=json
```

#### Validate Offline Package
```bash
# Validate package integrity
deshchaind query identity offline validate offline_package.json

# Validate compressed package
deshchaind query identity offline validate \
  compressed_package.bin \
  --format=compressed
```

### 4. Offline Formats

#### Query Supported Formats
```bash
deshchaind query identity offline formats
```

Response includes:
- **Self-Contained**: 10MB max, complete verification package
- **Compressed**: 2MB max, space-efficient for mobile
- **QR Code**: 4KB max, visual scanning compatible
- **NFC**: 8KB max, contactless verification
- **Printable**: 100KB max, human-readable paper format

#### Query Verification Modes
```bash
deshchaind query identity offline modes
```

Available modes:
- **Full**: 24-hour validity, high security
- **Partial**: 7-day validity, cached verification
- **Minimal**: 30-day validity, basic identification
- **Emergency**: 1-hour validity, crisis situations

## Security Considerations

### 1. Cryptographic Security
- **Ed25519 signatures** for identity proofs
- **AES-256 encryption** for sensitive data
- **SHA-256 hashing** for data integrity
- **HMAC verification** for tamper detection

### 2. Biometric Security
- Templates encrypted with device-specific keys
- Never store raw biometric data offline
- Support for face, fingerprint, iris, voice, and palm
- Quality score validation for biometric matches

### 3. Expiration Management
- Configurable expiration periods per format
- Automatic cleanup of expired packages
- Grace period handling for emergency scenarios
- Real-time expiration checking

### 4. Revocation Handling
- Offline revocation list distribution
- Merkle tree-based efficient revocation checking
- Emergency revocation mechanisms
- Cached revocation data with periodic updates

## Configuration

### System Configuration
```bash
# Update offline verification settings
deshchaind tx identity offline update-config \
  --max-offline-hours=48 \
  --required-confidence=0.9 \
  --biometric-threshold=0.95 \
  --max-cache-mb=100 \
  --cache-expiration-hours=168 \
  --enable-compression=true \
  --compression-level=6 \
  --emergency-mode=true \
  --emergency-threshold=0.70 \
  --default-language=en \
  --supported-regions="india,global" \
  --from admin
```

### Query Configuration
```bash
deshchaind query identity offline config
```

## Multi-Language Support

### Supported Languages
The offline verification system supports 22 Indian languages:
- Hindi (hi), English (en), Bengali (bn)
- Tamil (ta), Telugu (te), Marathi (mr)
- Gujarati (gu), Kannada (kn), Malayalam (ml)
- Punjabi (pa), Odia (or), Assamese (as)
- Urdu (ur), Sanskrit (sa), Nepali (ne)
- Konkani (gom), Manipuri (mni), Bodo (brx)
- Santhali (sat), Maithili (mai), Dogri (doi)
- Kashmiri (ks)

### Localized Messages
```json
{
  "verification_success": {
    "en": "Verification successful",
    "hi": "सत्यापन सफल",
    "bn": "যাচাইকরণ সফল"
  },
  "verification_failed": {
    "en": "Verification failed",
    "hi": "सत्यापन असफल",
    "bn": "যাচাইকরণ ব্যর্থ"
  }
}
```

## Performance Specifications

### Throughput
- **10,000+ verifications/second** for cached data
- **1,000+ verifications/second** for full verification
- **100+ verifications/second** with biometric matching

### Latency
- **<100ms** for cached verification
- **<500ms** for full verification
- **<2s** for biometric verification

### Storage Requirements
- **Self-Contained**: 5-10MB per identity
- **Compressed**: 1-2MB per identity
- **QR Code**: 2-4KB per identity
- **NFC**: 4-8KB per identity
- **Printable**: 50-100KB per identity

## Use Cases

### 1. Rural Banking
- Village ATMs with limited connectivity
- Agricultural loan verification
- Pension distribution systems
- Microfinance identity checks

### 2. Emergency Services
- Disaster response identification
- Medical emergency access
- Evacuation center verification
- Emergency contact systems

### 3. Government Services
- Rural service center operations
- Election voter verification
- Ration distribution systems
- Healthcare access in remote areas

### 4. Educational Systems
- Exam center verification
- Remote learning access
- Scholarship verification
- Student identity management

## Best Practices

### 1. Package Management
- Generate packages just before offline periods
- Use shortest reasonable expiration times
- Regularly update revocation data
- Implement automatic cleanup procedures

### 2. Device Security
- Use device-specific encryption keys
- Implement secure element storage where available
- Regular device authentication checks
- Monitor for suspicious device activity

### 3. Biometric Handling
- Never store raw biometric data
- Use template encryption for all formats
- Implement quality score thresholds
- Support multiple biometric modalities

### 4. Network Optimization
- Compress packages for bandwidth-limited areas
- Use differential updates for revocation lists
- Implement intelligent caching strategies
- Support progressive package download

## Troubleshooting

### Common Issues

#### Package Validation Errors
```bash
# Check package integrity
deshchaind query identity offline validate package.json

# Common solutions:
# 1. Check expiration date
# 2. Verify signature
# 3. Validate format compliance
# 4. Check revocation status
```

#### Device Registration Issues
```bash
# Verify device capabilities
deshchaind query identity offline devices did:desh:user123

# Common solutions:
# 1. Check security level requirements
# 2. Verify public key format
# 3. Validate device capabilities
# 4. Check maximum offline duration
```

#### Verification Failures
```bash
# Debug verification process
deshchaind query identity offline verify data.json request.json --output=json

# Common solutions:
# 1. Check data expiration
# 2. Verify challenge-response
# 3. Validate credential requirements
# 4. Check biometric thresholds
```

## API Reference

### Query APIs
- `GET /cosmos/identity/v1/offline/config` - Get offline configuration
- `GET /cosmos/identity/v1/offline/devices/{did}` - List registered devices
- `GET /cosmos/identity/v1/offline/formats` - Get supported formats
- `GET /cosmos/identity/v1/offline/modes` - Get verification modes
- `POST /cosmos/identity/v1/offline/verify` - Perform offline verification
- `POST /cosmos/identity/v1/offline/validate` - Validate offline package

### Transaction APIs
- `POST /cosmos/identity/v1/tx/offline/prepare` - Prepare offline verification
- `POST /cosmos/identity/v1/tx/offline/backup` - Create offline backup
- `POST /cosmos/identity/v1/tx/offline/config` - Update configuration
- `POST /cosmos/identity/v1/tx/offline/register` - Register device
- `POST /cosmos/identity/v1/tx/offline/revoke` - Revoke offline access

## Integration Examples

### Mobile App Integration
```javascript
// React Native example
import { DeshChainIdentity } from '@deshchain/identity-sdk';

const identityClient = new DeshChainIdentity(config);

// Prepare offline package
const offlineData = await identityClient.prepareOfflineVerification({
  did: 'did:desh:user123',
  format: 'compressed',
  expirationHours: 24,
  includeBiometric: true
});

// Verify offline
const result = await identityClient.verifyOffline(offlineData, verificationRequest);
```

### IoT Device Integration
```c
// C SDK example for IoT devices
#include "deshchain_identity.h"

// Initialize offline verification
desh_offline_config_t config = {
    .format = DESH_FORMAT_NFC,
    .security_level = 2,
    .max_offline_hours = 48
};

// Verify identity
desh_verification_result_t result;
int status = desh_verify_offline(&config, offline_data, &result);
```

---

For more information, see the [DeshChain Identity Documentation](./README.md) and [API Reference](./api.md).