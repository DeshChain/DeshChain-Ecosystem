# DeshChain Identity SDK

A comprehensive TypeScript/JavaScript SDK for interacting with DeshChain's Identity Module, supporting W3C DIDs, Verifiable Credentials, and Zero-Knowledge Proofs.

## Installation

```bash
npm install @deshchain/identity-sdk
# or
yarn add @deshchain/identity-sdk
```

## Quick Start

```typescript
import { createIdentitySigningClient } from '@deshchain/identity-sdk';

// Initialize the SDK
const { sdk, address } = await createIdentitySigningClient({
  rpcEndpoint: 'https://rpc.deshchain.com',
  chainId: 'deshchain-1'
}, 'your mnemonic phrase here');

// Create an identity
const did = await sdk.createIdentity(address, publicKey);
console.log('Created DID:', did);
```

## Features

- 🆔 **W3C DID Management**: Create and manage decentralized identifiers
- 📜 **Verifiable Credentials**: Issue, verify, and present credentials
- 🔐 **Zero-Knowledge Proofs**: Privacy-preserving authentication
- 🇮🇳 **India Stack Integration**: Aadhaar, DigiLocker, UPI support
- 👆 **Biometric Authentication**: Multi-modal biometric support
- 🔒 **Selective Disclosure**: Share only required information

## Usage Examples

### Creating an Identity

```typescript
// Create a new identity with metadata
const did = await sdk.createIdentity(
  address,
  publicKey,
  {
    name: 'John Doe',
    type: 'individual',
    country: 'IN'
  }
);
```

### Issuing Credentials

```typescript
// Issue a KYC credential
const credentialId = await sdk.issueCredential(
  issuerAddress,
  holderDID,
  ['KYCCredential'],
  {
    kycLevel: 'enhanced',
    verifiedAt: new Date().toISOString(),
    riskRating: 'low'
  },
  new Date(Date.now() + 365 * 24 * 60 * 60 * 1000) // 1 year expiry
);

// Issue an education credential
const eduCredId = await sdk.issueCredential(
  universityAddress,
  studentDID,
  ['EducationCredential', 'DegreeCredential'],
  {
    degree: 'Bachelor of Technology',
    field: 'Computer Science',
    university: 'IIT Delhi',
    year: 2023,
    grade: 'A+'
  }
);
```

### Verifying Credentials

```typescript
// Verify a credential
const credential = await sdk.getCredentialById(credentialId);
const isValid = await sdk.verifyCredential(credential);

if (isValid) {
  console.log('Credential is valid');
} else {
  console.log('Credential verification failed');
}
```

### Selective Disclosure

```typescript
// Present specific attributes from credentials
const presentationId = await sdk.presentCredential(
  holderAddress,
  ['cred-123', 'cred-456'],
  verifierDID,
  {
    'cred-123': ['kycLevel', 'verifiedAt'],
    'cred-456': ['degree', 'university']
  }
);
```

### Zero-Knowledge Proofs

```typescript
// Prove age without revealing date of birth
const proofId = await sdk.createZKProof(
  proverAddress,
  {
    type: 'age-range',
    statement: 'age >= 18 AND age < 65',
    credentials: ['age-credential-id'],
    options: {
      anonymitySet: 10,
      expiryMinutes: 60
    }
  }
);

// Prove income range for loan eligibility
const incomeProofId = await sdk.createZKProof(
  proverAddress,
  {
    type: 'income-range',
    statement: 'annual_income >= 500000',
    credentials: ['income-credential-id']
  }
);
```

### Biometric Authentication

```typescript
// Register biometric
const biometricCredId = await sdk.registerBiometric(
  userAddress,
  {
    type: 'FINGERPRINT',
    templateHash: 'hash-of-biometric-template',
    deviceId: 'device-123'
  }
);

// Authenticate with biometric
const isAuthenticated = await sdk.authenticateBiometric(
  userAddress,
  {
    type: 'FINGERPRINT',
    templateHash: 'hash-to-verify',
    deviceId: 'device-123'
  }
);
```

### India Stack Integration

```typescript
// Link Aadhaar with consent
await sdk.linkAadhaar(
  userAddress,
  aadhaarHash, // Hashed Aadhaar number
  consentArtefact // Consent ID from DEPA
);

// Connect DigiLocker
await sdk.connectDigiLocker(
  userAddress,
  authToken,
  ['education', 'identity'] // Document types to access
);

// Link UPI ID
await sdk.linkUPI(
  userAddress,
  'user@paytm',
  verificationCode
);
```

### Privacy Settings

```typescript
// Update privacy settings
await sdk.updatePrivacySettings(
  userAddress,
  {
    disclosureLevel: 'minimal', // minimal, standard, full
    requireConsent: true,
    allowAnonymous: true
  }
);
```

## Advanced Usage

### Query Operations

```typescript
// Get identity by DID
const identity = await sdk.getIdentity('did:desh:abc123');

// Get all credentials for a subject
const credentials = await sdk.getCredentialsBySubject('did:desh:abc123');

// Get credentials by type
const kycCredentials = credentials.filter(
  cred => cred.type.includes('KYCCredential')
);
```

### Batch Operations

```typescript
// Issue multiple credentials
const credentialRequests = [
  { holder: did1, type: ['KYCCredential'], subject: {...} },
  { holder: did2, type: ['KYCCredential'], subject: {...} }
];

const credentialIds = await Promise.all(
  credentialRequests.map(req => 
    sdk.issueCredential(issuerAddress, req.holder, req.type, req.subject)
  )
);
```

### Error Handling

```typescript
try {
  const did = await sdk.createIdentity(address, publicKey);
} catch (error) {
  if (error.message.includes('already exists')) {
    console.log('Identity already exists for this address');
  } else {
    console.error('Failed to create identity:', error);
  }
}
```

## Integration with DeshChain Modules

### TradeFinance Integration

```typescript
// Verify KYC for trade finance
const hasValidKYC = await sdk.verifyKYCForTradeFinance(customerDID);

if (hasValidKYC) {
  // Proceed with trade finance operation
}
```

### MoneyOrder Integration

```typescript
// High-value transaction verification
const isVerified = await sdk.verifyBiometricForTransaction(
  userAddress,
  {
    amount: '1000000', // Amount in NAMO
    biometric: {
      type: 'FINGERPRINT',
      templateHash: 'hash',
      deviceId: 'device-123'
    }
  }
);
```

## Configuration

### Custom Registry

```typescript
import { Registry } from '@cosmjs/proto-signing';
import { identityTypes } from '@deshchain/identity-sdk';

const registry = new Registry();
registry.register('/deshchain.identity.MsgCreateIdentity', identityTypes.MsgCreateIdentity);
// Register other types...

const sdk = new DeshChainIdentitySDK({
  rpcEndpoint: 'https://rpc.deshchain.com',
  chainId: 'deshchain-1',
  registry // Use custom registry
});
```

### Offline Mode

```typescript
// Create credentials for offline verification
const offlineCredential = await sdk.createOfflineCredential(
  credential,
  { includeProof: true }
);

// Verify offline
const isValid = sdk.verifyOfflineCredential(offlineCredential);
```

## Best Practices

1. **Key Management**
   - Never expose private keys or mnemonics
   - Use hardware wallets for production
   - Implement key rotation policies

2. **Privacy**
   - Request only necessary attributes
   - Use selective disclosure
   - Implement consent management

3. **Performance**
   - Cache frequently accessed identities
   - Batch operations when possible
   - Use pagination for large queries

4. **Security**
   - Verify all credentials before use
   - Check credential expiry
   - Validate issuer trust

## API Reference

See [API Documentation](./docs/API.md) for complete reference.

## Examples

Check out the [examples](./examples) directory for more usage scenarios:
- Basic identity operations
- Credential workflows
- Zero-knowledge proofs
- Biometric authentication
- India Stack integration

## Contributing

We welcome contributions! Please see our [Contributing Guide](../../CONTRIBUTING.md).

## License

Apache 2.0 - see [LICENSE](../../LICENSE) for details.