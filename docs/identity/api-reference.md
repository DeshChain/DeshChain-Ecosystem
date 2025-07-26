# DeshChain Identity API Reference

## Overview

The DeshChain Identity API provides comprehensive access to decentralized identity features including DID management, verifiable credentials, privacy features, and India Stack integration. The API supports multiple protocols including REST, GraphQL, and gRPC.

## Base URLs

- **Mainnet**: `https://api.deshchain.network/identity/v1`
- **Testnet**: `https://testnet-api.deshchain.network/identity/v1`
- **Local**: `http://localhost:1317/identity/v1`

## Authentication

### API Key Authentication
```http
Authorization: Bearer <api-key>
```

### DID Authentication
```http
Authorization: DID <did>
Signature: <signed-challenge>
```

### Biometric Authentication
```http
Authorization: Biometric <biometric-token>
Device-ID: <device-identifier>
```

## Identity Management APIs

### Create Identity

Creates a new decentralized identity with DID.

**Endpoint**: `POST /identity`

**Request Body**:
```json
{
  "public_key": "base64-encoded-public-key",
  "recovery_methods": [
    {
      "type": "email",
      "value": "user@example.com"
    },
    {
      "type": "phone",
      "value": "+91-9876543210"
    }
  ],
  "privacy_level": "advanced",
  "metadata": {
    "name": "John Doe",
    "preferred_language": "en"
  }
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "did": "did:desh:1234567890abcdef",
    "address": "desh1abc123def456...",
    "did_document": {
      "@context": ["https://w3id.org/did/v1"],
      "id": "did:desh:1234567890abcdef",
      "controller": "did:desh:1234567890abcdef",
      "verification_method": [{
        "id": "did:desh:1234567890abcdef#key-1",
        "type": "Ed25519VerificationKey2018",
        "controller": "did:desh:1234567890abcdef",
        "public_key_base58": "H3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV"
      }],
      "authentication": ["did:desh:1234567890abcdef#key-1"],
      "service": []
    },
    "created_at": "2024-12-01T10:00:00Z"
  }
}
```

### Get Identity

Retrieves identity information by DID.

**Endpoint**: `GET /identity/{did}`

**Parameters**:
- `did` (path): The DID of the identity to retrieve

**Response**:
```json
{
  "success": true,
  "data": {
    "did": "did:desh:1234567890abcdef",
    "status": "active",
    "verification_level": "enhanced",
    "linked_credentials": 5,
    "last_activity": "2024-12-01T09:30:00Z",
    "privacy_settings": {
      "level": "advanced",
      "selective_disclosure": true,
      "anonymous_mode": false
    }
  }
}
```

### Update Identity

Updates identity information and settings.

**Endpoint**: `PUT /identity/{did}`

**Request Body**:
```json
{
  "metadata": {
    "name": "John Doe Updated",
    "preferred_language": "hi"
  },
  "privacy_level": "ultimate",
  "add_recovery_methods": [
    {
      "type": "guardian",
      "value": "did:desh:guardian123"
    }
  ]
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "did": "did:desh:1234567890abcdef",
    "updated_at": "2024-12-01T10:30:00Z",
    "changes": ["metadata", "privacy_level", "recovery_methods"]
  }
}
```

## DID Document APIs

### Get DID Document

Retrieves the W3C DID Document for a given DID.

**Endpoint**: `GET /did/{did}`

**Parameters**:
- `did` (path): The DID to resolve
- `version` (query, optional): Specific version of the document

**Response**:
```json
{
  "@context": [
    "https://w3id.org/did/v1",
    "https://deshchain.network/did/v1"
  ],
  "id": "did:desh:1234567890abcdef",
  "controller": "did:desh:1234567890abcdef",
  "verification_method": [
    {
      "id": "did:desh:1234567890abcdef#key-1",
      "type": "Ed25519VerificationKey2018",
      "controller": "did:desh:1234567890abcdef",
      "public_key_base58": "H3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV"
    }
  ],
  "authentication": ["did:desh:1234567890abcdef#key-1"],
  "assertion_method": ["did:desh:1234567890abcdef#key-1"],
  "key_agreement": ["did:desh:1234567890abcdef#key-1"],
  "capability_invocation": ["did:desh:1234567890abcdef#key-1"],
  "capability_delegation": ["did:desh:1234567890abcdef#key-1"],
  "service": [
    {
      "id": "did:desh:1234567890abcdef#deshchain-identity",
      "type": "DeshChainIdentityService",
      "service_endpoint": "https://identity.deshchain.network"
    }
  ],
  "created": "2024-12-01T10:00:00Z",
  "updated": "2024-12-01T10:30:00Z"
}
```

### Update DID Document

Updates the DID document with new verification methods or services.

**Endpoint**: `PUT /did/{did}`

**Request Body**:
```json
{
  "add_verification_methods": [
    {
      "id": "did:desh:1234567890abcdef#key-2",
      "type": "X25519KeyAgreementKey2019",
      "controller": "did:desh:1234567890abcdef",
      "public_key_base58": "LkZnDq9uVbJ8gJMvdpKnP2Rk3sF7xNm9qWe6tBv5cAz2"
    }
  ],
  "add_services": [
    {
      "id": "did:desh:1234567890abcdef#messaging",
      "type": "MessagingService",
      "service_endpoint": "https://messaging.example.com"
    }
  ]
}
```

## Verifiable Credentials APIs

### Issue Credential

Issues a new verifiable credential.

**Endpoint**: `POST /credentials/issue`

**Request Body**:
```json
{
  "issuer": "did:desh:issuer123",
  "subject": "did:desh:1234567890abcdef",
  "type": ["VerifiableCredential", "KYCCredential"],
  "credential_subject": {
    "id": "did:desh:1234567890abcdef",
    "kyc_level": "enhanced",
    "verification_date": "2024-12-01",
    "document_verified": true,
    "biometric_verified": true
  },
  "expiration_date": "2025-12-01T00:00:00Z",
  "evidence": [
    {
      "type": "DocumentVerification",
      "verifier": "did:desh:verifier123",
      "evidence_document": "aadhaar_verification",
      "verification_method": "biometric_match"
    }
  ]
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "@context": [
      "https://www.w3.org/2018/credentials/v1",
      "https://deshchain.network/credentials/v1"
    ],
    "id": "https://deshchain.network/credentials/kyc/12345",
    "type": ["VerifiableCredential", "KYCCredential"],
    "issuer": "did:desh:issuer123",
    "issuance_date": "2024-12-01T10:00:00Z",
    "expiration_date": "2025-12-01T00:00:00Z",
    "credential_subject": {
      "id": "did:desh:1234567890abcdef",
      "kyc_level": "enhanced",
      "verification_date": "2024-12-01",
      "document_verified": true,
      "biometric_verified": true
    },
    "proof": {
      "type": "Ed25519Signature2018",
      "created": "2024-12-01T10:00:00Z",
      "verification_method": "did:desh:issuer123#key-1",
      "proof_purpose": "assertion_method",
      "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..."
    }
  }
}
```

### Verify Credential

Verifies a verifiable credential's authenticity and validity.

**Endpoint**: `POST /credentials/verify`

**Request Body**:
```json
{
  "credential": {
    "@context": ["https://www.w3.org/2018/credentials/v1"],
    "id": "https://deshchain.network/credentials/kyc/12345",
    "type": ["VerifiableCredential", "KYCCredential"],
    "issuer": "did:desh:issuer123",
    "issuance_date": "2024-12-01T10:00:00Z",
    "credential_subject": {
      "id": "did:desh:1234567890abcdef",
      "kyc_level": "enhanced"
    },
    "proof": {
      "type": "Ed25519Signature2018",
      "created": "2024-12-01T10:00:00Z",
      "verification_method": "did:desh:issuer123#key-1",
      "proof_purpose": "assertion_method",
      "jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..."
    }
  },
  "options": {
    "check_revocation": true,
    "check_expiration": true,
    "require_fresh_proof": false
  }
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "verified": true,
    "checks": {
      "signature_valid": true,
      "issuer_trusted": true,
      "not_expired": true,
      "not_revoked": true,
      "schema_valid": true
    },
    "issuer_info": {
      "did": "did:desh:issuer123",
      "name": "DeshChain KYC Authority",
      "trust_level": "high"
    },
    "credential_info": {
      "issued": "2024-12-01T10:00:00Z",
      "expires": "2025-12-01T00:00:00Z",
      "type": "KYCCredential",
      "status": "valid"
    }
  }
}
```

### Get Credentials

Retrieves credentials for a specific subject.

**Endpoint**: `GET /credentials/subject/{did}`

**Parameters**:
- `did` (path): Subject DID
- `type` (query, optional): Filter by credential type
- `issuer` (query, optional): Filter by issuer DID
- `status` (query, optional): Filter by status (valid, expired, revoked)

**Response**:
```json
{
  "success": true,
  "data": {
    "credentials": [
      {
        "id": "https://deshchain.network/credentials/kyc/12345",
        "type": ["VerifiableCredential", "KYCCredential"],
        "issuer": "did:desh:issuer123",
        "issued": "2024-12-01T10:00:00Z",
        "expires": "2025-12-01T00:00:00Z",
        "status": "valid",
        "summary": {
          "kyc_level": "enhanced",
          "verification_date": "2024-12-01"
        }
      }
    ],
    "total_count": 1,
    "page": 1,
    "per_page": 10
  }
}
```

### Revoke Credential

Revokes a previously issued credential.

**Endpoint**: `POST /credentials/revoke`

**Request Body**:
```json
{
  "credential_id": "https://deshchain.network/credentials/kyc/12345",
  "revocation_reason": "credential_compromised",
  "revocation_date": "2024-12-01T15:00:00Z"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "credential_id": "https://deshchain.network/credentials/kyc/12345",
    "revoked": true,
    "revocation_date": "2024-12-01T15:00:00Z",
    "revocation_reason": "credential_compromised"
  }
}
```

## Privacy & Zero-Knowledge APIs

### Create ZK Proof

Generates a zero-knowledge proof for selective disclosure.

**Endpoint**: `POST /privacy/zk-proof`

**Request Body**:
```json
{
  "statement": "age >= 18",
  "credentials": [
    "https://deshchain.network/credentials/age/67890"
  ],
  "revealed_attributes": [],
  "proof_purpose": "age_verification"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "proof": {
      "type": "BBS+Signature2020",
      "created": "2024-12-01T10:00:00Z",
      "verification_method": "did:desh:1234567890abcdef#key-1",
      "proof_purpose": "assertion_method",
      "proof_value": "kTTk8ehCy3zs8VMrKX7U9..."
    },
    "disclosed_attributes": {},
    "verification_key": "zkVK_base64_encoded_key"
  }
}
```

### Verify ZK Proof

Verifies a zero-knowledge proof.

**Endpoint**: `POST /privacy/verify-zk-proof`

**Request Body**:
```json
{
  "proof": {
    "type": "BBS+Signature2020",
    "created": "2024-12-01T10:00:00Z",
    "verification_method": "did:desh:1234567890abcdef#key-1",
    "proof_purpose": "assertion_method",
    "proof_value": "kTTk8ehCy3zs8VMrKX7U9..."
  },
  "statement": "age >= 18",
  "verification_key": "zkVK_base64_encoded_key"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "verified": true,
    "statement_satisfied": true,
    "proof_valid": true,
    "verification_time": "2024-12-01T10:01:00Z"
  }
}
```

### Anonymous Credential

Creates an anonymous credential for privacy-preserving authentication.

**Endpoint**: `POST /privacy/anonymous-credential`

**Request Body**:
```json
{
  "base_credential": "https://deshchain.network/credentials/kyc/12345",
  "attributes_to_hide": ["name", "address"],
  "attributes_to_prove": ["age_over_18", "citizenship"],
  "validity_period": "24h"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "anonymous_credential": "anon_cred_base64_encoded",
    "presentation_token": "pres_token_12345",
    "expires_at": "2024-12-02T10:00:00Z",
    "proof_capabilities": ["age_over_18", "citizenship"]
  }
}
```

## India Stack Integration APIs

### Link Aadhaar

Links Aadhaar identity with DeshChain identity.

**Endpoint**: `POST /india-stack/aadhaar/link`

**Request Body**:
```json
{
  "did": "did:desh:1234567890abcdef",
  "aadhaar_hash": "sha256_hash_of_aadhaar",
  "consent_artifact": "consent_artifact_id",
  "otp": "123456",
  "biometric_data": "base64_encoded_biometric"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "linked": true,
    "aadhaar_token": "encrypted_aadhaar_reference",
    "verification_level": "enhanced",
    "kyc_credential_id": "https://deshchain.network/credentials/aadhaar-kyc/11111",
    "linked_at": "2024-12-01T10:00:00Z"
  }
}
```

### DigiLocker Connect

Connects to DigiLocker for document verification.

**Endpoint**: `POST /india-stack/digilocker/connect`

**Request Body**:
```json
{
  "did": "did:desh:1234567890abcdef",
  "authorization_code": "digilocker_auth_code",
  "consent_documents": ["aadhaar", "pan", "driving_license"]
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "connected": true,
    "access_token": "encrypted_access_token",
    "available_documents": [
      {
        "type": "aadhaar",
        "issuer": "UIDAI",
        "verified": true
      },
      {
        "type": "pan",
        "issuer": "Income Tax Department",
        "verified": true
      }
    ],
    "credentials_issued": [
      "https://deshchain.network/credentials/document/pan/22222"
    ]
  }
}
```

### UPI Identity Link

Links UPI identity for payment verification.

**Endpoint**: `POST /india-stack/upi/link`

**Request Body**:
```json
{
  "did": "did:desh:1234567890abcdef",
  "upi_id": "user@paytm",
  "verification_amount": 1.00,
  "consent_reference": "upi_consent_ref"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "linked": true,
    "upi_reference": "encrypted_upi_reference",
    "verification_status": "verified",
    "payment_credential_id": "https://deshchain.network/credentials/upi/33333"
  }
}
```

## Biometric Authentication APIs

### Register Biometric

Registers biometric data for authentication.

**Endpoint**: `POST /biometric/register`

**Request Body**:
```json
{
  "did": "did:desh:1234567890abcdef",
  "biometric_type": "fingerprint",
  "biometric_data": "base64_encoded_template",
  "device_info": {
    "device_id": "device_12345",
    "device_type": "fingerprint_scanner",
    "manufacturer": "SecuGen",
    "model": "Hamster Pro 20"
  },
  "liveness_proof": "base64_liveness_data"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "biometric_id": "bio_12345",
    "biometric_type": "fingerprint",
    "registered": true,
    "quality_score": 0.95,
    "device_bound": true,
    "expires_at": "2026-12-01T10:00:00Z"
  }
}
```

### Authenticate Biometric

Performs biometric authentication.

**Endpoint**: `POST /biometric/authenticate`

**Request Body**:
```json
{
  "did": "did:desh:1234567890abcdef",
  "biometric_type": "fingerprint",
  "biometric_sample": "base64_encoded_sample",
  "device_id": "device_12345",
  "challenge": "authentication_challenge",
  "liveness_proof": "base64_liveness_data"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "authenticated": true,
    "confidence_score": 0.98,
    "biometric_token": "bio_auth_token_67890",
    "valid_until": "2024-12-01T11:00:00Z",
    "authentication_level": "high"
  }
}
```

### Multi-Modal Authentication

Performs multi-modal biometric authentication.

**Endpoint**: `POST /biometric/multi-modal`

**Request Body**:
```json
{
  "did": "did:desh:1234567890abcdef",
  "modalities": [
    {
      "type": "face",
      "data": "base64_face_image",
      "liveness_proof": "base64_liveness_video"
    },
    {
      "type": "voice",
      "data": "base64_voice_sample",
      "passphrase": "my voice is my password"
    }
  ],
  "device_id": "device_12345",
  "challenge": "multi_modal_challenge"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "authenticated": true,
    "overall_confidence": 0.99,
    "modality_results": [
      {
        "type": "face",
        "confidence": 0.98,
        "liveness_verified": true
      },
      {
        "type": "voice",
        "confidence": 0.97,
        "passphrase_verified": true
      }
    ],
    "authentication_token": "multi_bio_token_99999",
    "valid_until": "2024-12-01T12:00:00Z"
  }
}
```

## Cross-Module APIs

### Request Module Access

Requests access to another module's identity data.

**Endpoint**: `POST /cross-module/request-access`

**Request Body**:
```json
{
  "requesting_module": "tradefinance",
  "target_did": "did:desh:1234567890abcdef",
  "requested_attributes": ["kyc_level", "verification_date"],
  "purpose": "trade_finance_compliance",
  "access_duration": "24h"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "request_id": "access_req_12345",
    "status": "approved",
    "access_token": "cross_module_token_67890",
    "expires_at": "2024-12-02T10:00:00Z",
    "allowed_attributes": ["kyc_level", "verification_date"]
  }
}
```

### Share Identity Data

Shares identity data with another module.

**Endpoint**: `POST /cross-module/share-data`

**Request Body**:
```json
{
  "access_token": "cross_module_token_67890",
  "requesting_module": "tradefinance",
  "data_minimization": true,
  "consent_verified": true
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "shared_data": {
      "kyc_level": "enhanced",
      "verification_date": "2024-12-01"
    },
    "access_logged": true,
    "usage_restrictions": ["compliance_only", "no_storage"]
  }
}
```

## Query APIs

### Search Identities

Searches for identities based on criteria (admin only).

**Endpoint**: `GET /query/identities`

**Parameters**:
- `kyc_level` (query): Filter by KYC level
- `verification_date_from` (query): Filter by verification date range
- `verification_date_to` (query): Filter by verification date range
- `status` (query): Filter by identity status
- `page` (query): Page number
- `per_page` (query): Results per page

**Response**:
```json
{
  "success": true,
  "data": {
    "identities": [
      {
        "did": "did:desh:1234567890abcdef",
        "kyc_level": "enhanced",
        "verification_date": "2024-12-01",
        "status": "active",
        "last_activity": "2024-12-01T09:30:00Z"
      }
    ],
    "total_count": 1,
    "page": 1,
    "per_page": 10
  }
}
```

### Identity Statistics

Gets identity system statistics.

**Endpoint**: `GET /query/statistics`

**Response**:
```json
{
  "success": true,
  "data": {
    "total_identities": 1000000,
    "verified_identities": 850000,
    "credentials_issued": 2500000,
    "active_sessions": 50000,
    "biometric_registrations": 600000,
    "india_stack_connections": 400000,
    "daily_authentications": 100000,
    "success_rate": 99.5
  }
}
```

## Error Responses

### Standard Error Format

```json
{
  "success": false,
  "error": {
    "code": "INVALID_DID",
    "message": "The provided DID is not valid",
    "details": {
      "field": "did",
      "value": "invalid_did_value",
      "expected_format": "did:desh:..."
    },
    "request_id": "req_12345",
    "timestamp": "2024-12-01T10:00:00Z"
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| INVALID_DID | 400 | Malformed or invalid DID |
| DID_NOT_FOUND | 404 | DID does not exist |
| CREDENTIAL_EXPIRED | 400 | Credential has expired |
| CREDENTIAL_REVOKED | 400 | Credential has been revoked |
| INSUFFICIENT_PERMISSIONS | 403 | Insufficient permissions for operation |
| BIOMETRIC_MISMATCH | 401 | Biometric authentication failed |
| CONSENT_REQUIRED | 403 | User consent required for operation |
| RATE_LIMIT_EXCEEDED | 429 | API rate limit exceeded |
| INDIA_STACK_ERROR | 502 | Error communicating with India Stack |
| ZK_PROOF_INVALID | 400 | Zero-knowledge proof verification failed |

## Rate Limits

### Default Limits

- **Identity Operations**: 1000 requests/hour
- **Credential Operations**: 5000 requests/hour
- **Verification Operations**: 10000 requests/hour
- **Query Operations**: 10000 requests/hour

### Rate Limit Headers

```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1701432000
```

## SDKs and Libraries

### JavaScript/TypeScript
```bash
npm install @deshchain/identity-sdk
```

### Python
```bash
pip install deshchain-identity
```

### Go
```bash
go get github.com/deshchain/identity-go
```

### Java
```xml
<dependency>
    <groupId>com.deshchain</groupId>
    <artifactId>identity-java-sdk</artifactId>
    <version>1.0.0</version>
</dependency>
```

---

**Last Updated**: December 2024  
**Version**: 1.0  
**API Version**: v1  
**Maintainers**: DeshChain Identity API Team