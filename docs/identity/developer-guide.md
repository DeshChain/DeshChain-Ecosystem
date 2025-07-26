# DeshChain Identity Developer Guide

## Overview

This guide provides comprehensive examples and best practices for integrating with the DeshChain Identity system. Learn how to implement W3C DID-based identities, verifiable credentials, biometric authentication, and privacy-preserving features in your applications.

## Quick Start

### 1. Setup Identity SDK

#### JavaScript/TypeScript
```bash
npm install @deshchain/identity-sdk @deshchain/sdk
```

```typescript
import { IdentityClient, DIDManager, CredentialManager } from '@deshchain/identity-sdk';

const identityClient = new IdentityClient({
    rpcEndpoint: 'https://rpc.deshchain.com',
    chainId: 'deshchain-1',
    apiKey: 'your-api-key'
});
```

#### Python
```bash
pip install deshchain-identity deshchain-sdk
```

```python
from deshchain_identity import IdentityClient, DIDManager, CredentialManager

client = IdentityClient(
    rpc_endpoint="https://rpc.deshchain.com",
    chain_id="deshchain-1",
    api_key="your-api-key"
)
```

#### Go
```bash
go get github.com/deshchain/identity-go-sdk
```

```go
import (
    "github.com/deshchain/identity-go-sdk/client"
    "github.com/deshchain/identity-go-sdk/types"
)

client := client.NewIdentityClient(client.Config{
    RPCEndpoint: "https://rpc.deshchain.com",
    ChainID:     "deshchain-1",
    APIKey:      "your-api-key",
})
```

### 2. Create Your First Identity

#### TypeScript Example
```typescript
import { IdentityClient } from '@deshchain/identity-sdk';

async function createIdentity() {
    const client = new IdentityClient({
        rpcEndpoint: 'https://rpc.deshchain.com',
        chainId: 'deshchain-1'
    });

    // Create new identity with DID
    const identity = await client.createIdentity({
        recoveryMethods: [
            {
                type: 'email',
                value: 'user@example.com'
            },
            {
                type: 'phone',
                value: '+91-9876543210'
            }
        ],
        privacyLevel: 'advanced',
        metadata: {
            name: 'Rajesh Kumar',
            preferredLanguage: 'hi'
        }
    });

    console.log(`Created DID: ${identity.did}`);
    console.log(`Address: ${identity.address}`);
    
    return identity;
}
```

#### Python Example
```python
from deshchain_identity import IdentityClient
import asyncio

async def create_identity():
    client = IdentityClient(
        rpc_endpoint="https://rpc.deshchain.com",
        chain_id="deshchain-1"
    )
    
    # Create new identity
    identity = await client.create_identity(
        recovery_methods=[
            {
                "type": "email",
                "value": "user@example.com"
            },
            {
                "type": "guardian",
                "value": "did:desh:guardian123"
            }
        ],
        privacy_level="ultimate",
        metadata={
            "name": "Priya Sharma",
            "preferred_language": "ta"
        }
    )
    
    print(f"Created DID: {identity.did}")
    print(f"Address: {identity.address}")
    
    return identity

# Run the async function
identity = asyncio.run(create_identity())
```

#### Go Example
```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/deshchain/identity-go-sdk/client"
    "github.com/deshchain/identity-go-sdk/types"
)

func createIdentity() (*types.Identity, error) {
    client := client.NewIdentityClient(client.Config{
        RPCEndpoint: "https://rpc.deshchain.com",
        ChainID:     "deshchain-1",
    })

    identity, err := client.CreateIdentity(context.Background(), &types.CreateIdentityRequest{
        RecoveryMethods: []types.RecoveryMethod{
            {
                Type:  "email",
                Value: "user@example.com",
            },
            {
                Type:  "biometric",
                Value: "fingerprint_template_hash",
            },
        },
        PrivacyLevel: "advanced",
        Metadata: map[string]interface{}{
            "name":               "Amit Patel",
            "preferred_language": "gu",
        },
    })

    if err != nil {
        return nil, fmt.Errorf("failed to create identity: %w", err)
    }

    fmt.Printf("Created DID: %s\n", identity.DID)
    fmt.Printf("Address: %s\n", identity.Address)
    
    return identity, nil
}
```

## Identity Management

### Working with DIDs

#### Resolving DID Documents
```typescript
// Resolve a DID to get the DID Document
async function resolveDID(did: string) {
    const didDocument = await client.resolveDID(did);
    
    console.log('DID Document:', {
        id: didDocument.id,
        controller: didDocument.controller,
        verificationMethods: didDocument.verificationMethod.length,
        services: didDocument.service?.length || 0
    });
    
    return didDocument;
}

// Usage
const didDoc = await resolveDID('did:desh:1234567890abcdef');
```

#### Updating DID Documents
```typescript
async function updateDIDDocument(did: string) {
    const updateRequest = {
        did: did,
        addVerificationMethods: [
            {
                id: `${did}#key-2`,
                type: 'X25519KeyAgreementKey2019',
                controller: did,
                publicKeyBase58: 'LkZnDq9uVbJ8gJMvdpKnP2Rk3sF7xNm9qWe6tBv5cAz2'
            }
        ],
        addServices: [
            {
                id: `${did}#messaging`,
                type: 'MessagingService',
                serviceEndpoint: 'https://messaging.example.com'
            }
        ]
    };
    
    const result = await client.updateDIDDocument(updateRequest);
    console.log('DID Document updated:', result);
    
    return result;
}
```

### India Stack Integration

#### Linking Aadhaar Identity
```typescript
async function linkAadhaar(did: string) {
    // Step 1: Request OTP
    const otpRequest = await client.requestAadhaarOTP({
        did: did,
        aadhaarNumber: 'xxxx-xxxx-1234', // Last 4 digits only
        consentGiven: true
    });
    
    console.log(`OTP sent to: ${otpRequest.maskedMobile}`);
    
    // Step 2: Verify OTP and link
    const otp = await getUserInput('Enter OTP: '); // Get OTP from user
    
    const linkResult = await client.linkAadhaar({
        did: did,
        sessionId: otpRequest.sessionId,
        otp: otp,
        consentArtifact: otpRequest.consentArtifact
    });
    
    console.log('Aadhaar linked successfully:', {
        credentialId: linkResult.kycCredentialId,
        verificationLevel: linkResult.verificationLevel
    });
    
    return linkResult;
}
```

#### DigiLocker Integration
```typescript
async function connectDigiLocker(did: string) {
    // Step 1: Get authorization URL
    const authURL = await client.getDigiLockerAuthURL({
        did: did,
        requestedDocuments: ['aadhaar', 'pan', 'driving_license'],
        redirectURI: 'https://yourapp.com/callback'
    });
    
    console.log('Redirect user to:', authURL);
    
    // Step 2: After user authorization, exchange code
    const authCode = 'received_from_callback'; // Get from callback
    
    const connection = await client.connectDigiLocker({
        did: did,
        authorizationCode: authCode,
        consentDocuments: ['aadhaar', 'pan']
    });
    
    console.log('DigiLocker connected:', {
        documentsAvailable: connection.availableDocuments.length,
        credentialsIssued: connection.credentialsIssued.length
    });
    
    return connection;
}
```

#### UPI Identity Linking
```typescript
async function linkUPIIdentity(did: string, upiId: string) {
    const linkResult = await client.linkUPI({
        did: did,
        upiId: upiId,
        verificationAmount: 1.00, // ₹1 for verification
        consentReference: 'upi_consent_12345'
    });
    
    console.log('UPI linked:', {
        paymentCredentialId: linkResult.paymentCredentialId,
        verificationStatus: linkResult.verificationStatus
    });
    
    return linkResult;
}
```

## Verifiable Credentials

### Issuing Credentials

#### KYC Credential Example
```typescript
async function issueKYCCredential(issuerDID: string, subjectDID: string) {
    const credential = await client.issueCredential({
        issuer: issuerDID,
        subject: subjectDID,
        type: ['VerifiableCredential', 'KYCCredential'],
        credentialSubject: {
            id: subjectDID,
            kycLevel: 'enhanced',
            verificationDate: new Date().toISOString(),
            documentVerified: true,
            biometricVerified: true,
            riskLevel: 'low'
        },
        expirationDate: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000), // 1 year
        evidence: [
            {
                type: 'DocumentVerification',
                verifier: issuerDID,
                evidenceDocument: 'aadhaar_verification',
                verificationMethod: 'biometric_match'
            }
        ]
    });
    
    console.log('KYC Credential issued:', credential.id);
    return credential;
}
```

#### Education Credential Example
```python
async def issue_education_credential(issuer_did: str, student_did: str):
    client = IdentityClient()
    
    credential = await client.issue_credential(
        issuer=issuer_did,
        subject=student_did,
        type=["VerifiableCredential", "EducationCredential"],
        credential_subject={
            "id": student_did,
            "degree": "B.Tech Computer Science",
            "university": "IIT Bombay",
            "graduation_year": 2024,
            "cgpa": 8.5,
            "specialization": "Artificial Intelligence"
        },
        expiration_date="2034-12-31T23:59:59Z",
        evidence=[
            {
                "type": "AcademicTranscript",
                "verifier": issuer_did,
                "transcript_hash": "sha256_hash_of_transcript",
                "verification_method": "digital_signature"
            }
        ]
    )
    
    print(f"Education Credential issued: {credential.id}")
    return credential
```

### Verifying Credentials

#### Basic Verification
```typescript
async function verifyCredential(credential: any) {
    const verificationResult = await client.verifyCredential({
        credential: credential,
        options: {
            checkRevocation: true,
            checkExpiration: true,
            requireFreshProof: false
        }
    });
    
    if (verificationResult.verified) {
        console.log('Credential is valid ✅');
        console.log('Checks passed:', verificationResult.checks);
    } else {
        console.log('Credential verification failed ❌');
        console.log('Failed checks:', verificationResult.failedChecks);
    }
    
    return verificationResult;
}
```

#### Advanced Verification with Selective Disclosure
```typescript
async function verifyWithSelectiveDisclosure(credential: any, requiredClaims: string[]) {
    const presentation = await client.createPresentation({
        credentials: [credential],
        revealedClaims: {
            [credential.id]: requiredClaims
        },
        verifier: 'did:desh:verifier123',
        challenge: 'random_challenge_string'
    });
    
    const verificationResult = await client.verifyPresentation({
        presentation: presentation,
        challenge: 'random_challenge_string',
        verifier: 'did:desh:verifier123'
    });
    
    console.log('Selective disclosure verification:', verificationResult);
    return verificationResult;
}

// Usage: Only reveal degree and university, hide CGPA and personal details
await verifyWithSelectiveDisclosure(educationCredential, ['degree', 'university']);
```

## Biometric Authentication

### Registering Biometrics

#### Fingerprint Registration
```typescript
async function registerFingerprint(did: string, fingerprintTemplate: string) {
    const registration = await client.registerBiometric({
        did: did,
        biometricType: 'fingerprint',
        biometricData: fingerprintTemplate, // Base64 encoded template
        deviceInfo: {
            deviceId: 'device_12345',
            deviceType: 'fingerprint_scanner',
            manufacturer: 'SecuGen',
            model: 'Hamster Pro 20'
        },
        livenessProof: 'base64_liveness_data'
    });
    
    console.log('Fingerprint registered:', {
        biometricId: registration.biometricId,
        qualityScore: registration.qualityScore,
        expiresAt: registration.expiresAt
    });
    
    return registration;
}
```

#### Face Recognition Registration
```typescript
async function registerFaceRecognition(did: string, faceImage: string) {
    const registration = await client.registerBiometric({
        did: did,
        biometricType: 'face',
        biometricData: faceImage, // Base64 encoded image
        deviceInfo: {
            deviceId: 'camera_67890',
            deviceType: 'camera',
            manufacturer: 'Logitech',
            model: 'C920 HD Pro'
        },
        livenessProof: 'base64_liveness_video' // Video proof
    });
    
    console.log('Face recognition registered:', registration);
    return registration;
}
```

### Biometric Authentication

#### Single-Modal Authentication
```typescript
async function authenticateWithFingerprint(did: string, fingerprintSample: string) {
    const authResult = await client.authenticateBiometric({
        did: did,
        biometricType: 'fingerprint',
        biometricSample: fingerprintSample,
        deviceId: 'device_12345',
        challenge: 'auth_challenge_12345',
        livenessProof: 'base64_liveness_data'
    });
    
    if (authResult.authenticated) {
        console.log('Fingerprint authentication successful ✅');
        console.log(`Confidence: ${authResult.confidenceScore}`);
        console.log(`Token valid until: ${authResult.validUntil}`);
        
        // Use the biometric token for subsequent operations
        return authResult.biometricToken;
    } else {
        console.log('Fingerprint authentication failed ❌');
        return null;
    }
}
```

#### Multi-Modal Authentication
```typescript
async function multiModalAuthentication(did: string) {
    const authResult = await client.authenticateMultiModal({
        did: did,
        modalities: [
            {
                type: 'face',
                data: 'base64_face_image',
                livenessProof: 'base64_liveness_video'
            },
            {
                type: 'voice',
                data: 'base64_voice_sample',
                passphrase: 'my voice is my password'
            }
        ],
        deviceId: 'multi_device_12345',
        challenge: 'multi_modal_challenge'
    });
    
    console.log('Multi-modal authentication:', {
        authenticated: authResult.authenticated,
        overallConfidence: authResult.overallConfidence,
        modalityResults: authResult.modalityResults
    });
    
    return authResult;
}
```

## Privacy & Zero-Knowledge Proofs

### Creating Zero-Knowledge Proofs

#### Age Verification Without Revealing Birthdate
```typescript
async function proveAgeOver18(did: string, ageCredentialId: string) {
    const zkProof = await client.createZKProof({
        statement: 'age >= 18',
        credentials: [ageCredentialId],
        revealedAttributes: [], // Don't reveal any attributes
        proofPurpose: 'age_verification'
    });
    
    console.log('ZK Proof created for age verification');
    return zkProof;
}

// Verify the proof
async function verifyAgeProof(zkProof: any) {
    const verificationResult = await client.verifyZKProof({
        proof: zkProof.proof,
        statement: 'age >= 18',
        verificationKey: zkProof.verificationKey
    });
    
    console.log('Age proof verified:', verificationResult.verified);
    return verificationResult;
}
```

#### Income Range Proof
```python
async def prove_income_range(did: str, income_credential_id: str, min_income: int):
    client = IdentityClient()
    
    zk_proof = await client.create_zk_proof(
        statement=f"income >= {min_income}",
        credentials=[income_credential_id],
        revealed_attributes=[],
        proof_purpose="loan_eligibility"
    )
    
    print("Income range proof created (without revealing exact income)")
    return zk_proof

# Usage: Prove income >= ₹5,00,000 without revealing exact amount
income_proof = await prove_income_range(
    did="did:desh:user123", 
    income_credential_id="income_cred_456",
    min_income=500000
)
```

### Anonymous Credentials

```typescript
async function createAnonymousCredential(baseCredentialId: string) {
    const anonCredential = await client.createAnonymousCredential({
        baseCredential: baseCredentialId,
        attributesToHide: ['name', 'address', 'phone'],
        attributesToProve: ['age_over_18', 'indian_citizen', 'kyc_verified'],
        validityPeriod: '24h'
    });
    
    console.log('Anonymous credential created:', {
        presentationToken: anonCredential.presentationToken,
        expiresAt: anonCredential.expiresAt,
        proofCapabilities: anonCredential.proofCapabilities
    });
    
    return anonCredential;
}
```

## Cross-Module Integration

### Using Identity in Financial Modules

#### TradeFinance Integration
```typescript
import { TradeFinanceClient } from '@deshchain/sdk';
import { IdentityClient } from '@deshchain/identity-sdk';

async function createTradeFinanceLC(buyerDID: string, sellerDID: string) {
    const identityClient = new IdentityClient();
    const tradeClient = new TradeFinanceClient();
    
    // Verify KYC status for both parties
    const buyerKYC = await identityClient.getKYCStatus(buyerDID);
    const sellerKYC = await identityClient.getKYCStatus(sellerDID);
    
    if (buyerKYC.level !== 'enhanced' || sellerKYC.level !== 'enhanced') {
        throw new Error('Enhanced KYC required for trade finance');
    }
    
    // Create Letter of Credit with identity verification
    const lc = await tradeClient.createLC({
        buyer: buyerDID,
        seller: sellerDID,
        amount: 1000000, // $10,000
        currency: 'USD',
        identityVerificationRequired: true,
        kycCredentials: {
            buyer: buyerKYC.credentialId,
            seller: sellerKYC.credentialId
        }
    });
    
    console.log('LC created with identity verification:', lc.id);
    return lc;
}
```

#### Lending Module Integration
```typescript
async function applyForBusinessLoan(applicantDID: string) {
    const identityClient = new IdentityClient();
    const lendingClient = new VyavasayaMitraClient();
    
    // Get required credentials
    const kycCredential = await identityClient.getCredential(applicantDID, 'KYCCredential');
    const businessCredential = await identityClient.getCredential(applicantDID, 'BusinessCredential');
    const incomeCredential = await identityClient.getCredential(applicantDID, 'IncomeCredential');
    
    // Create ZK proof for income without revealing exact amount
    const incomeProof = await identityClient.createZKProof({
        statement: 'annual_income >= 1000000', // ₹10 lakh minimum
        credentials: [incomeCredential.id],
        revealedAttributes: [],
        proofPurpose: 'loan_eligibility'
    });
    
    // Apply for loan
    const application = await lendingClient.applyForLoan({
        applicant: applicantDID,
        loanType: 'working_capital',
        amount: 500000, // ₹5 lakh
        tenure: 24, // months
        credentials: {
            kyc: kycCredential.id,
            business: businessCredential.id
        },
        incomeProof: incomeProof,
        identityVerified: true
    });
    
    console.log('Loan application submitted:', application.id);
    return application;
}
```

### MoneyOrder DEX Integration

```typescript
async function createSecureMoneyOrder(senderDID: string, recipientAddress: string) {
    const identityClient = new IdentityClient();
    const moneyOrderClient = new MoneyOrderClient();
    
    // Perform biometric authentication
    const biometricToken = await identityClient.authenticateBiometric({
        did: senderDID,
        biometricType: 'fingerprint',
        biometricSample: 'user_fingerprint_sample',
        deviceId: 'user_device_123',
        challenge: 'money_order_challenge'
    });
    
    if (!biometricToken.authenticated) {
        throw new Error('Biometric authentication required for money orders');
    }
    
    // Create money order with biometric verification
    const moneyOrder = await moneyOrderClient.createOrder({
        sender: senderDID,
        recipient: recipientAddress,
        amount: 50000, // ₹50,000
        currency: 'NAMO',
        biometricVerification: biometricToken.biometricToken,
        requireRecipientKYC: true,
        privacyLevel: 'advanced' // Hide sender identity
    });
    
    console.log('Secure money order created:', moneyOrder.id);
    return moneyOrder;
}
```

## Best Practices

### 1. Security Best Practices

#### Secure Key Management
```typescript
class SecureKeyManager {
    private keystore: any;
    
    constructor() {
        // Use hardware security module when available
        this.keystore = new HardwareKeystore();
    }
    
    async generateKeys(did: string): Promise<KeyPair> {
        return await this.keystore.generateEd25519KeyPair({
            purpose: 'authentication',
            did: did,
            protected: true
        });
    }
    
    async signWithBiometric(data: string, did: string): Promise<string> {
        // Require biometric authentication for signing
        const biometricAuth = await this.authenticateBiometric(did);
        if (!biometricAuth.success) {
            throw new Error('Biometric authentication required');
        }
        
        return await this.keystore.sign(data, {
            keyId: biometricAuth.keyId,
            biometricToken: biometricAuth.token
        });
    }
}
```

#### Privacy Protection
```typescript
class PrivacyManager {
    async shareMinimalData(credentials: any[], requiredAttributes: string[]) {
        // Only share absolutely necessary attributes
        const presentation = await client.createSelectiveDisclosure({
            credentials: credentials,
            revealedAttributes: requiredAttributes,
            hideUnrequiredData: true,
            useZKProofs: true
        });
        
        return presentation;
    }
    
    async checkConsentBeforeSharing(userDID: string, requestingParty: string, attributes: string[]) {
        const consent = await client.getConsent({
            userDID: userDID,
            requestingParty: requestingParty,
            requestedAttributes: attributes
        });
        
        if (!consent.granted) {
            throw new Error('User consent required for data sharing');
        }
        
        return consent;
    }
}
```

### 2. Error Handling

```typescript
class IdentityErrorHandler {
    async handleIdentityOperation<T>(operation: () => Promise<T>): Promise<T> {
        try {
            return await operation();
        } catch (error) {
            if (error.code === 'BIOMETRIC_MISMATCH') {
                console.error('Biometric authentication failed. Please try again.');
                // Implement retry logic or alternative authentication
                return await this.handleBiometricFailure();
            } else if (error.code === 'CREDENTIAL_EXPIRED') {
                console.error('Credential has expired. Please renew.');
                // Guide user to credential renewal process
                return await this.handleCredentialRenewal();
            } else if (error.code === 'CONSENT_REQUIRED') {
                console.error('User consent required. Redirecting to consent page.');
                // Show consent UI to user
                return await this.handleConsentRequest();
            }
            
            throw error; // Re-throw unknown errors
        }
    }
}
```

### 3. Performance Optimization

```typescript
class IdentityCache {
    private cache = new Map();
    private readonly TTL = 5 * 60 * 1000; // 5 minutes
    
    async getDIDDocument(did: string) {
        const cacheKey = `did:${did}`;
        const cached = this.cache.get(cacheKey);
        
        if (cached && Date.now() - cached.timestamp < this.TTL) {
            return cached.data;
        }
        
        const didDocument = await client.resolveDID(did);
        this.cache.set(cacheKey, {
            data: didDocument,
            timestamp: Date.now()
        });
        
        return didDocument;
    }
    
    async getCredentialWithCache(credentialId: string) {
        // Implement credential caching with invalidation
        const cacheKey = `credential:${credentialId}`;
        // ... similar caching logic
    }
}
```

### 4. Testing Your Integration

#### Unit Tests
```typescript
describe('Identity Integration', () => {
    let client: IdentityClient;
    
    beforeEach(() => {
        client = new IdentityClient({
            rpcEndpoint: 'http://localhost:26657',
            chainId: 'deshchain-test'
        });
    });
    
    test('should create identity with valid DID', async () => {
        const identity = await client.createIdentity({
            recoveryMethods: [{ type: 'email', value: 'test@example.com' }],
            privacyLevel: 'basic'
        });
        
        expect(identity.did).toMatch(/^did:desh:[a-f0-9]+$/);
        expect(identity.address).toBeDefined();
    });
    
    test('should issue and verify credential', async () => {
        const credential = await client.issueCredential({
            issuer: 'did:desh:issuer123',
            subject: 'did:desh:subject456',
            type: ['VerifiableCredential', 'TestCredential'],
            credentialSubject: { test: true }
        });
        
        const verification = await client.verifyCredential({
            credential: credential
        });
        
        expect(verification.verified).toBe(true);
    });
});
```

#### Integration Tests
```bash
# Run identity integration tests
npm test -- --testPathPattern=identity

# Test with real network
npm run test:integration -- --network=testnet

# Performance tests
npm run test:performance -- --module=identity
```

## Advanced Examples

### 1. Multi-Signature Identity Operations

```typescript
async function createMultiSigIdentity(controllers: string[], threshold: number) {
    const identity = await client.createMultiSigIdentity({
        controllers: controllers,
        threshold: threshold,
        privacyLevel: 'ultimate'
    });
    
    console.log(`Multi-sig identity created with ${threshold}/${controllers.length} threshold`);
    return identity;
}

async function multiSigCredentialIssuance(multiSigDID: string, credentialData: any) {
    // Create credential proposal
    const proposal = await client.createCredentialProposal({
        issuer: multiSigDID,
        ...credentialData
    });
    
    // Collect signatures from required number of controllers
    const signatures = [];
    for (const controller of proposal.requiredSigners) {
        const signature = await getControllerSignature(controller, proposal);
        signatures.push(signature);
        
        if (signatures.length >= proposal.threshold) {
            break;
        }
    }
    
    // Issue credential with multi-sig approval
    const credential = await client.issueCredentialWithMultiSig({
        proposal: proposal,
        signatures: signatures
    });
    
    return credential;
}
```

### 2. Identity Recovery Scenarios

```typescript
async function performSocialRecovery(lostDID: string, guardianDIDs: string[]) {
    // Initiate recovery process
    const recoveryRequest = await client.initiateRecovery({
        lostIdentity: lostDID,
        recoveryMethod: 'social_recovery',
        guardians: guardianDIDs
    });
    
    console.log(`Recovery initiated. Need ${recoveryRequest.requiredGuardians} guardian approvals.`);
    
    // Collect guardian approvals
    const approvals = [];
    for (const guardianDID of guardianDIDs) {
        const approval = await client.approveRecovery({
            recoveryRequestId: recoveryRequest.id,
            guardianDID: guardianDID,
            approvalSignature: 'guardian_signature'
        });
        
        approvals.push(approval);
        
        if (approvals.length >= recoveryRequest.requiredGuardians) {
            break;
        }
    }
    
    // Complete recovery
    const newIdentity = await client.completeRecovery({
        recoveryRequestId: recoveryRequest.id,
        newKeys: 'new_key_pair',
        guardianApprovals: approvals
    });
    
    console.log('Identity recovered successfully:', newIdentity.did);
    return newIdentity;
}
```

### 3. Enterprise Identity Management

```typescript
class EnterpriseIdentityManager {
    private client: IdentityClient;
    
    constructor() {
        this.client = new IdentityClient({
            enterprise: true,
            batchOperations: true
        });
    }
    
    async bulkCreateEmployeeIdentities(employees: any[]) {
        const operations = employees.map(employee => ({
            type: 'createIdentity',
            data: {
                metadata: {
                    name: employee.name,
                    employeeId: employee.id,
                    department: employee.department
                },
                privacyLevel: 'advanced',
                enterpriseManaged: true
            }
        }));
        
        const results = await this.client.batchExecute(operations);
        console.log(`Created ${results.successful.length} employee identities`);
        
        return results;
    }
    
    async issueEmployeeCredentials(employeeIdentities: any[]) {
        const credentialOps = employeeIdentities.map(identity => ({
            type: 'issueCredential',
            data: {
                issuer: 'did:desh:company123',
                subject: identity.did,
                type: ['VerifiableCredential', 'EmployeeCredential'],
                credentialSubject: {
                    employeeId: identity.metadata.employeeId,
                    department: identity.metadata.department,
                    issueDate: new Date().toISOString()
                }
            }
        }));
        
        return await this.client.batchExecute(credentialOps);
    }
}
```

## Resources

### Documentation Links
- [Identity API Reference](./api-reference.md)
- [Architecture Overview](./architecture.md)
- [Security Guidelines](./privacy-security.md)
- [Compliance Guide](./compliance-guide.md)

### SDK Documentation
- [JavaScript SDK](https://www.npmjs.com/package/@deshchain/identity-sdk)
- [Python SDK](https://pypi.org/project/deshchain-identity/)
- [Go SDK](https://pkg.go.dev/github.com/deshchain/identity-go-sdk)

### Community Resources
- [GitHub Repository](https://github.com/deshchain/identity)
- [Discord Community](https://discord.gg/deshchain-identity)
- [Developer Forum](https://forum.deshchain.com/identity)
- [Stack Overflow](https://stackoverflow.com/questions/tagged/deshchain-identity)

### Support
- **Technical Support**: identity-support@deshchain.com
- **Enterprise Support**: enterprise@deshchain.com
- **Security Issues**: security@deshchain.com

---

**Last Updated**: December 2024  
**Version**: 1.0  
**Maintainers**: DeshChain Identity Development Team