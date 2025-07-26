# Identity Module Documentation

## Overview

The DeshChain Identity Module is a comprehensive, production-ready decentralized identity solution that implements W3C DID standards, Verifiable Credentials, and privacy-preserving authentication mechanisms. It provides seamless integration with India Stack services while enabling self-sovereign identity for all DeshChain users.

## Key Features

### 1. Decentralized Identity (DID)
- **W3C Compliant**: Full implementation of W3C DID specification
- **Multiple DID Methods**: Support for `did:desh` and `did:web`
- **DID Documents**: Comprehensive identity documents with verification methods
- **Service Endpoints**: Discoverable services associated with DIDs

### 2. Verifiable Credentials
- **Issuance**: Create tamper-proof digital credentials
- **Verification**: Cryptographically verify credential authenticity
- **Selective Disclosure**: Share only required information
- **Revocation**: Manage credential lifecycle with revocation registry

### 3. Privacy Features
- **Zero-Knowledge Proofs**: Prove claims without revealing data
- **Anonymous Credentials**: Use credentials without revealing identity
- **Consent Management**: GDPR/DPDP Act compliant consent framework
- **Privacy Levels**: Basic, Advanced, and Ultimate privacy options

### 4. India Stack Integration
- **Aadhaar Integration**: Secure Aadhaar-based KYC with consent
- **DigiLocker Support**: Access and verify government documents
- **UPI Identity**: Link payment identities
- **Village Panchayat KYC**: Grassroots-level identity verification

### 5. Biometric Authentication
- **Multi-Modal Support**: Fingerprint, Face, Iris, Voice, Palm
- **Device Binding**: Secure device registration
- **Liveness Detection**: Anti-spoofing measures
- **Template Protection**: Secure biometric template storage

## Technical Architecture

### Module Structure
```
x/identity/
├── types/          # Core types and messages
│   ├── did.go      # DID document structures
│   ├── vc.go       # Verifiable Credential types
│   ├── identity.go # Identity management types
│   ├── privacy.go  # Privacy and ZKP types
│   └── india_stack.go # India Stack integration
├── keeper/         # Business logic
│   ├── did_keeper.go     # DID operations
│   ├── vc_keeper.go      # Credential management
│   ├── privacy_keeper.go # Privacy features
│   └── integration/      # Module integrations
└── client/         # CLI and REST interfaces
```

### State Management

#### Identity Storage
```go
// Primary identity record
type Identity struct {
    Did              string
    Controller       string  
    Status           IdentityStatus
    RecoveryMethods  []RecoveryMethod
    LinkedIdentities []string
    Metadata         map[string]string
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

#### DID Document
```go
type DIDDocument struct {
    Context            []string
    ID                 string
    Controller         string
    VerificationMethod []VerificationMethod
    Authentication     []interface{}
    Service            []Service
    Created            time.Time
    Updated            time.Time
}
```

#### Verifiable Credential
```go
type VerifiableCredential struct {
    Context           []string
    ID                string
    Type              []string
    Issuer            string
    IssuanceDate      time.Time
    ExpirationDate    *time.Time
    CredentialSubject interface{}
    CredentialStatus  *CredentialStatus
    Proof             *Proof
}
```

## Integration with Other Modules

### TradeFinance Integration
The Identity module provides KYC verification for trade finance operations:
- Creates DID-based identities for customers
- Issues verifiable KYC credentials
- Supports selective disclosure for compliance
- Maintains backward compatibility with existing KYC

### MoneyOrder Integration
Biometric authentication for secure money transfers:
- Registers biometric credentials
- Performs multi-factor authentication
- Supports high-value transaction verification
- Integrates with existing biometric systems

### Validator Integration
Identity verification for validators:
- Validator identity credentials
- Reputation tracking
- Multi-signature support

## Usage Examples

### Creating an Identity
```bash
# Create a new identity with DID
deshchaind tx identity create-identity \
  --public-key="<base64-encoded-public-key>" \
  --from=<account>

# Link Aadhaar for KYC (with consent)
deshchaind tx identity link-aadhaar \
  --aadhaar-hash="<hash>" \
  --consent-artefact="<consent-id>" \
  --from=<account>
```

### Issuing Credentials
```bash
# Issue a KYC credential
deshchaind tx identity issue-credential \
  --holder=<holder-address> \
  --type="KYCCredential" \
  --claims='{"level":"enhanced","verified":true}' \
  --from=<issuer>

# Issue an education credential
deshchaind tx identity issue-credential \
  --holder=<holder-address> \
  --type="EducationCredential" \
  --claims='{"degree":"B.Tech","university":"IIT"}' \
  --from=<issuer>
```

### Selective Disclosure
```bash
# Present specific attributes from credentials
deshchaind tx identity present-credential \
  --credentials="cred1,cred2" \
  --verifier=<verifier-address> \
  --revealed-claims='{"cred1":["degree"]}' \
  --from=<holder>
```

### Zero-Knowledge Proofs
```bash
# Create age proof without revealing date of birth
deshchaind tx identity create-zk-proof \
  --type="age-range" \
  --statement="age >= 18" \
  --credentials="age-credential-id" \
  --from=<account>
```

## Privacy & Security

### Data Protection
- **Encryption**: All sensitive data encrypted at rest
- **Minimal Storage**: Only essential data stored on-chain
- **IPFS Integration**: Large data stored off-chain with encryption
- **Key Management**: Hierarchical deterministic keys

### Access Control
- **Capability-Based**: Fine-grained access permissions
- **Time-Bound**: Temporary access with expiry
- **Delegation**: Controlled delegation chains
- **Audit Trail**: Immutable access logs

### Compliance
- **GDPR/DPDP Act**: Full compliance with data protection laws
- **RBI Guidelines**: Meets KYC/AML requirements
- **UIDAI Regulations**: Compliant Aadhaar integration
- **W3C Standards**: Full DID and VC specification compliance

## API Reference

### gRPC Services
```protobuf
service Msg {
  // Identity Management
  rpc CreateIdentity(MsgCreateIdentity) returns (MsgCreateIdentityResponse);
  rpc UpdateIdentity(MsgUpdateIdentity) returns (MsgUpdateIdentityResponse);
  rpc RevokeIdentity(MsgRevokeIdentity) returns (MsgRevokeIdentityResponse);
  
  // DID Operations
  rpc RegisterDID(MsgRegisterDID) returns (MsgRegisterDIDResponse);
  rpc UpdateDID(MsgUpdateDID) returns (MsgUpdateDIDResponse);
  rpc DeactivateDID(MsgDeactivateDID) returns (MsgDeactivateDIDResponse);
  
  // Credential Management
  rpc IssueCredential(MsgIssueCredential) returns (MsgIssueCredentialResponse);
  rpc RevokeCredential(MsgRevokeCredential) returns (MsgRevokeCredentialResponse);
  rpc PresentCredential(MsgPresentCredential) returns (MsgPresentCredentialResponse);
  
  // Privacy Features
  rpc CreateZKProof(MsgCreateZKProof) returns (MsgCreateZKProofResponse);
  rpc VerifyZKProof(MsgVerifyZKProof) returns (MsgVerifyZKProofResponse);
  
  // India Stack Integration
  rpc LinkAadhaar(MsgLinkAadhaar) returns (MsgLinkAadhaarResponse);
  rpc ConnectDigiLocker(MsgConnectDigiLocker) returns (MsgConnectDigiLockerResponse);
  rpc LinkUPI(MsgLinkUPI) returns (MsgLinkUPIResponse);
}
```

### REST Endpoints
```bash
# Query identity
GET /deshchain/identity/v1/identity/{did}

# Query DID document
GET /deshchain/identity/v1/did/{did}

# Query credentials by subject
GET /deshchain/identity/v1/credentials/{subject}

# Verify credential
POST /deshchain/identity/v1/verify
{
  "credential": {...},
  "options": {...}
}

# Query India Stack status
GET /deshchain/identity/v1/india-stack/{did}
```

## Module Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| MaxDIDDocumentSize | 64KB | Maximum size of DID document |
| MaxCredentialSize | 32KB | Maximum size of credential |
| CredentialExpiryDays | 365 | Default credential validity |
| KYCExpiryDays | 180 | KYC credential validity |
| BiometricExpiryDays | 730 | Biometric credential validity |
| EnableZKProofs | true | Enable zero-knowledge proofs |
| EnableIndiaStack | true | Enable India Stack integration |
| RequireKYCForHighValue | true | Require KYC for high-value transactions |
| HighValueThreshold | 100,000 NAMO | Threshold for high-value transactions |

## Performance Metrics

- **DID Resolution**: < 50ms
- **Credential Verification**: < 100ms
- **ZK Proof Generation**: < 500ms
- **Biometric Matching**: < 200ms
- **Throughput**: 10,000+ identity operations/second

## Migration Guide

For existing modules integrating with Identity:
1. See [Identity Migration Guide](../IDENTITY_MIGRATION_GUIDE.md)
2. Use integration adapters for backward compatibility
3. Gradually migrate to DID-based identities

## Best Practices

### For Users
1. **Backup Recovery Methods**: Always set multiple recovery options
2. **Credential Management**: Regularly review and revoke unused credentials
3. **Privacy Settings**: Configure appropriate privacy levels
4. **Biometric Security**: Register biometrics only on trusted devices

### For Developers
1. **Use DIDs**: Prefer DIDs over addresses for identity
2. **Verify Credentials**: Always verify credential validity and expiry
3. **Respect Privacy**: Request only necessary attributes
4. **Handle Errors**: Gracefully handle identity verification failures

## Governance

Identity module parameters can be updated through governance:
```bash
# Propose parameter change
deshchaind tx gov submit-proposal param-change proposal.json --from proposer

# Example proposal.json
{
  "title": "Update Identity Module Parameters",
  "description": "Increase KYC expiry to 1 year",
  "changes": [{
    "subspace": "identity",
    "key": "KYCExpiryDays",
    "value": "365"
  }]
}
```

## Future Enhancements

1. **Cross-Chain Identity Bridge**: Share identities across blockchains
2. **Quantum-Resistant Cryptography**: Future-proof security
3. **AI-Powered Fraud Detection**: Enhanced security
4. **Decentralized Biometric Matching**: Privacy-preserving biometrics
5. **Self-Sovereign Organization Identity**: Corporate DIDs

## Support

- **Technical Documentation**: [GitHub Wiki](https://github.com/deshchain/wiki)
- **API Support**: identity-support@deshchain.com
- **Bug Reports**: [GitHub Issues](https://github.com/deshchain/issues)
- **Community**: [Discord](https://discord.gg/deshchain)