package keeper

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DigitalSignatureSystem manages digital signatures for documents
type DigitalSignatureSystem struct {
	keeper               Keeper
	certificateAuthority *CertificateAuthority
	signatureValidator   *SignatureValidator
	timestampAuthority   *TimestampAuthority
	revocationManager    *RevocationManager
	auditLogger          *SignatureAuditLogger
	mu                   sync.RWMutex
}

// CertificateAuthority manages digital certificates
type CertificateAuthority struct {
	rootCertificate     *x509.Certificate
	rootPrivateKey      *rsa.PrivateKey
	issuedCertificates  map[string]*IssuedCertificate
	certificateStore    *CertificateStore
	validationRules     []CertificateValidationRule
	certificateProfiles map[string]*CertificateProfile
}

// IssuedCertificate represents a certificate issued by the CA
type IssuedCertificate struct {
	CertificateID    string
	Certificate      *x509.Certificate
	PrivateKey       *rsa.PrivateKey
	IssuedTo         string
	IssuedAt         time.Time
	ExpiresAt        time.Time
	Status           CertificateStatus
	Purpose          []string
	RevocationReason string
	Metadata         map[string]string
}

// SignatureValidator validates digital signatures
type SignatureValidator struct {
	trustedRoots        []*x509.Certificate
	validationPolicies  map[string]*ValidationPolicy
	certificateChecker  *CertificateChecker
	timestampVerifier   *TimestampVerifier
	revocationChecker   *RevocationChecker
}

// TimestampAuthority provides trusted timestamps
type TimestampAuthority struct {
	tsaCertificate    *x509.Certificate
	tsaPrivateKey     *rsa.PrivateKey
	timestampRecords  map[string]*TimestampRecord
	accuracyTolerance time.Duration
	hashAlgorithm     crypto.Hash
}

// TimestampRecord represents a trusted timestamp
type TimestampRecord struct {
	TimestampID    string
	DocumentHash   string
	Timestamp      time.Time
	Accuracy       time.Duration
	SerialNumber   string
	TSACertificate string
	Signature      []byte
}

// RevocationManager handles certificate revocation
type RevocationManager struct {
	revocationList      map[string]*RevokedCertificate
	crlNumber           uint64
	lastCRLUpdate       time.Time
	nextCRLUpdate       time.Time
	ocspResponder       *OCSPResponder
	distributionPoints  []string
}

// RevokedCertificate represents a revoked certificate
type RevokedCertificate struct {
	SerialNumber     string
	RevocationTime   time.Time
	RevocationReason RevocationReason
	InvalidityDate   *time.Time
	RevokedBy        string
}

// SignatureAuditLogger logs all signature operations
type SignatureAuditLogger struct {
	auditRecords    []SignatureAuditRecord
	retentionPeriod time.Duration
	encryptionKey   []byte
}

// SignatureAuditRecord represents an audit log entry
type SignatureAuditRecord struct {
	RecordID        string
	Timestamp       time.Time
	Operation       SignatureOperation
	DocumentID      string
	SignerID        string
	CertificateID   string
	Result          OperationResult
	IPAddress       string
	UserAgent       string
	AdditionalInfo  map[string]string
}

// Enums and constants
type CertificateStatus int
type RevocationReason int
type SignatureOperation int
type OperationResult int
type SignatureFormat int
type SignatureLevel int

const (
	// Certificate Status
	CertificateActive CertificateStatus = iota
	CertificateRevoked
	CertificateExpired
	CertificateSuspended
	
	// Revocation Reasons
	RevocationUnspecified RevocationReason = iota
	RevocationKeyCompromise
	RevocationCACompromise
	RevocationAffiliationChanged
	RevocationSuperseded
	RevocationCessationOfOperation
	RevocationCertificateHold
	
	// Signature Operations
	SignOperation SignatureOperation = iota
	VerifyOperation
	TimestampOperation
	RevocationCheckOperation
	
	// Signature Formats
	PKCS7Format SignatureFormat = iota
	XMLDSigFormat
	PDFSignatureFormat
	JSONSignatureFormat
	
	// Signature Levels
	BasicSignature SignatureLevel = iota
	AdvancedSignature
	QualifiedSignature
)

// Core digital signature methods

// SignDocument creates a digital signature for a document
func (k Keeper) SignDocument(ctx context.Context, documentID string, signerID string, certificateID string) (*DigitalSignature, error) {
	dss := k.getDigitalSignatureSystem()
	
	// Get document
	document, err := k.getDocument(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}
	
	// Get signer's certificate
	cert, err := dss.certificateAuthority.getCertificate(certificateID)
	if err != nil {
		return nil, fmt.Errorf("certificate not found: %w", err)
	}
	
	// Validate certificate
	if err := dss.signatureValidator.validateCertificate(cert); err != nil {
		return nil, fmt.Errorf("certificate validation failed: %w", err)
	}
	
	// Create signature
	signature := &DigitalSignature{
		SignatureID:   generateID("sig"),
		DocumentID:    documentID,
		SignerID:      signerID,
		CertificateID: certificateID,
		SignatureTime: time.Now(),
		SignatureFormat: PKCS7Format,
		SignatureLevel: AdvancedSignature,
	}
	
	// Calculate document hash
	documentHash := sha256.Sum256(document.Content)
	signature.DocumentHash = base64.StdEncoding.EncodeToString(documentHash[:])
	
	// Sign the hash
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, cert.PrivateKey, crypto.SHA256, documentHash[:])
	if err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}
	signature.SignatureValue = base64.StdEncoding.EncodeToString(signatureBytes)
	
	// Get timestamp
	timestamp, err := dss.timestampAuthority.createTimestamp(documentHash[:])
	if err != nil {
		return nil, fmt.Errorf("timestamp creation failed: %w", err)
	}
	signature.TimestampToken = timestamp
	
	// Create signature attributes
	signature.SignedAttributes = map[string]string{
		"signingTime":     signature.SignatureTime.Format(time.RFC3339),
		"messageDigest":   signature.DocumentHash,
		"signingCertificate": cert.Certificate.SerialNumber.String(),
		"contentType":     document.DocumentType,
	}
	
	// Store signature
	if err := k.storeSignature(ctx, signature); err != nil {
		return nil, fmt.Errorf("failed to store signature: %w", err)
	}
	
	// Audit log
	dss.auditLogger.logOperation(SignOperation, documentID, signerID, certificateID, OperationSuccess)
	
	return signature, nil
}

// VerifySignature verifies a digital signature
func (k Keeper) VerifySignature(ctx context.Context, signatureID string) (*VerificationResult, error) {
	dss := k.getDigitalSignatureSystem()
	
	// Get signature
	signature, err := k.getSignature(ctx, signatureID)
	if err != nil {
		return nil, fmt.Errorf("signature not found: %w", err)
	}
	
	// Get document
	document, err := k.getDocument(ctx, signature.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}
	
	result := &VerificationResult{
		SignatureID:      signatureID,
		VerificationTime: time.Now(),
		SignatureValid:   true,
	}
	
	// Verify document integrity
	documentHash := sha256.Sum256(document.Content)
	expectedHash, _ := base64.StdEncoding.DecodeString(signature.DocumentHash)
	if !bytesEqual(documentHash[:], expectedHash) {
		result.SignatureValid = false
		result.Errors = append(result.Errors, "Document has been modified")
	}
	
	// Get certificate
	cert, err := dss.certificateAuthority.getCertificate(signature.CertificateID)
	if err != nil {
		result.SignatureValid = false
		result.Errors = append(result.Errors, "Certificate not found")
		return result, nil
	}
	
	// Verify certificate validity
	certValidation := dss.signatureValidator.validateCertificateWithDetails(cert)
	result.CertificateValid = certValidation.IsValid
	result.CertificateStatus = certValidation.Status
	if !certValidation.IsValid {
		result.SignatureValid = false
		result.Errors = append(result.Errors, certValidation.Errors...)
	}
	
	// Check revocation status
	revocationStatus, err := dss.revocationManager.checkRevocation(cert.Certificate.SerialNumber.String())
	if err == nil && revocationStatus.IsRevoked {
		result.SignatureValid = false
		result.CertificateStatus = CertificateRevoked
		result.Errors = append(result.Errors, fmt.Sprintf("Certificate revoked: %s", revocationStatus.Reason))
	}
	
	// Verify signature
	signatureBytes, _ := base64.StdEncoding.DecodeString(signature.SignatureValue)
	err = rsa.VerifyPKCS1v15(cert.Certificate.PublicKey.(*rsa.PublicKey), crypto.SHA256, documentHash[:], signatureBytes)
	if err != nil {
		result.SignatureValid = false
		result.Errors = append(result.Errors, "Signature verification failed")
	}
	
	// Verify timestamp
	if signature.TimestampToken != nil {
		timestampValid := dss.timestampAuthority.verifyTimestamp(signature.TimestampToken, documentHash[:])
		result.TimestampValid = timestampValid
		if !timestampValid {
			result.Warnings = append(result.Warnings, "Timestamp verification failed")
		}
	}
	
	// Audit log
	dss.auditLogger.logOperation(VerifyOperation, signature.DocumentID, "", signature.CertificateID, OperationSuccess)
	
	return result, nil
}

// Certificate Authority methods

// IssueCertificate issues a new digital certificate
func (ca *CertificateAuthority) issueCertificate(request *CertificateRequest) (*IssuedCertificate, error) {
	// Validate request
	if err := ca.validateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid certificate request: %w", err)
	}
	
	// Generate key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("key generation failed: %w", err)
	}
	
	// Create certificate template
	template := &x509.Certificate{
		SerialNumber: generateSerialNumber(),
		Subject: x509.Name{
			Organization:  []string{request.Organization},
			Country:       []string{request.Country},
			Province:      []string{request.Province},
			Locality:      []string{request.Locality},
			StreetAddress: []string{request.StreetAddress},
			PostalCode:    []string{request.PostalCode},
			CommonName:    request.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(request.Validity),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		BasicConstraintsValid: true,
	}
	
	// Apply certificate profile
	if profile, ok := ca.certificateProfiles[request.ProfileType]; ok {
		applyProfile(template, profile)
	}
	
	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, template, ca.rootCertificate, &privateKey.PublicKey, ca.rootPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("certificate creation failed: %w", err)
	}
	
	// Parse certificate
	certificate, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("certificate parsing failed: %w", err)
	}
	
	// Create issued certificate record
	issued := &IssuedCertificate{
		CertificateID: generateID("cert"),
		Certificate:   certificate,
		PrivateKey:    privateKey,
		IssuedTo:      request.CommonName,
		IssuedAt:      time.Now(),
		ExpiresAt:     certificate.NotAfter,
		Status:        CertificateActive,
		Purpose:       request.Purpose,
		Metadata:      request.Metadata,
	}
	
	// Store certificate
	ca.issuedCertificates[issued.CertificateID] = issued
	ca.certificateStore.store(issued)
	
	return issued, nil
}

// Revoke certificate
func (ca *CertificateAuthority) revokeCertificate(certificateID string, reason RevocationReason) error {
	cert, ok := ca.issuedCertificates[certificateID]
	if !ok {
		return fmt.Errorf("certificate not found")
	}
	
	if cert.Status == CertificateRevoked {
		return fmt.Errorf("certificate already revoked")
	}
	
	// Update status
	cert.Status = CertificateRevoked
	cert.RevocationReason = reason.String()
	
	// Add to revocation list
	ca.keeper.getDigitalSignatureSystem().revocationManager.addRevocation(&RevokedCertificate{
		SerialNumber:     cert.Certificate.SerialNumber.String(),
		RevocationTime:   time.Now(),
		RevocationReason: reason,
	})
	
	return nil
}

// Timestamp Authority methods

func (tsa *TimestampAuthority) createTimestamp(documentHash []byte) (*TimestampToken, error) {
	timestamp := &TimestampRecord{
		TimestampID:  generateID("ts"),
		DocumentHash: base64.StdEncoding.EncodeToString(documentHash),
		Timestamp:    time.Now(),
		Accuracy:     tsa.accuracyTolerance,
		SerialNumber: generateSerialNumber().String(),
	}
	
	// Create timestamp token
	token := &TimestampToken{
		Version:       1,
		Policy:        "1.2.3.4.5.6", // TSA policy OID
		MessageImprint: MessageImprint{
			HashAlgorithm: "SHA256",
			HashedMessage: timestamp.DocumentHash,
		},
		SerialNumber: timestamp.SerialNumber,
		GenTime:      timestamp.Timestamp,
		Accuracy:     timestamp.Accuracy,
	}
	
	// Sign timestamp token
	tokenBytes, _ := json.Marshal(token)
	tokenHash := sha256.Sum256(tokenBytes)
	signature, err := rsa.SignPKCS1v15(rand.Reader, tsa.tsaPrivateKey, crypto.SHA256, tokenHash[:])
	if err != nil {
		return nil, fmt.Errorf("timestamp signing failed: %w", err)
	}
	
	token.Signature = base64.StdEncoding.EncodeToString(signature)
	timestamp.Signature = signature
	
	// Store timestamp record
	tsa.timestampRecords[timestamp.TimestampID] = timestamp
	
	return token, nil
}

// Multi-signature support

// CreateMultiSignature creates a document requiring multiple signatures
func (k Keeper) CreateMultiSignature(ctx context.Context, documentID string, signers []SignerRequirement) (*MultiSignatureRequest, error) {
	request := &MultiSignatureRequest{
		RequestID:        generateID("msig"),
		DocumentID:       documentID,
		RequiredSigners:  signers,
		CreatedAt:        time.Now(),
		Status:           MultiSigPending,
		CollectedSignatures: make(map[string]*DigitalSignature),
	}
	
	// Validate signers
	for _, signer := range signers {
		if err := k.validateSigner(ctx, signer); err != nil {
			return nil, fmt.Errorf("invalid signer %s: %w", signer.SignerID, err)
		}
	}
	
	// Store request
	if err := k.storeMultiSignatureRequest(ctx, request); err != nil {
		return nil, err
	}
	
	return request, nil
}

// AddSignatureToMultiSig adds a signature to a multi-signature request
func (k Keeper) AddSignatureToMultiSig(ctx context.Context, requestID string, signature *DigitalSignature) error {
	request, err := k.getMultiSignatureRequest(ctx, requestID)
	if err != nil {
		return err
	}
	
	if request.Status != MultiSigPending {
		return fmt.Errorf("multi-signature request is not pending")
	}
	
	// Verify signer is authorized
	authorized := false
	for _, req := range request.RequiredSigners {
		if req.SignerID == signature.SignerID {
			authorized = true
			break
		}
	}
	if !authorized {
		return fmt.Errorf("signer not authorized for this document")
	}
	
	// Add signature
	request.CollectedSignatures[signature.SignerID] = signature
	
	// Check if all signatures collected
	if len(request.CollectedSignatures) == len(request.RequiredSigners) {
		request.Status = MultiSigComplete
		request.CompletedAt = timePtr(time.Now())
	}
	
	// Update request
	return k.updateMultiSignatureRequest(ctx, request)
}

// Helper types

type DigitalSignature struct {
	SignatureID      string
	DocumentID       string
	SignerID         string
	CertificateID    string
	SignatureTime    time.Time
	SignatureValue   string
	DocumentHash     string
	SignatureFormat  SignatureFormat
	SignatureLevel   SignatureLevel
	TimestampToken   *TimestampToken
	SignedAttributes map[string]string
	UnsignedAttributes map[string]string
}

type VerificationResult struct {
	SignatureID      string
	VerificationTime time.Time
	SignatureValid   bool
	CertificateValid bool
	TimestampValid   bool
	CertificateStatus CertificateStatus
	Errors           []string
	Warnings         []string
}

type CertificateRequest struct {
	CommonName     string
	Organization   string
	Country        string
	Province       string
	Locality       string
	StreetAddress  string
	PostalCode     string
	EmailAddress   string
	Validity       time.Duration
	ProfileType    string
	Purpose        []string
	Metadata       map[string]string
}

type TimestampToken struct {
	Version        int
	Policy         string
	MessageImprint MessageImprint
	SerialNumber   string
	GenTime        time.Time
	Accuracy       time.Duration
	Ordering       bool
	Nonce          string
	TSA            string
	Extensions     []Extension
	Signature      string
}

type MessageImprint struct {
	HashAlgorithm string
	HashedMessage string
}

type MultiSignatureRequest struct {
	RequestID           string
	DocumentID          string
	RequiredSigners     []SignerRequirement
	CollectedSignatures map[string]*DigitalSignature
	CreatedAt           time.Time
	CompletedAt         *time.Time
	Status              MultiSignatureStatus
	SigningOrder        SigningOrder
}

type SignerRequirement struct {
	SignerID         string
	SignerName       string
	SignerRole       string
	MandatoryOrder   int
	SigningDeadline  *time.Time
}

type MultiSignatureStatus int
type SigningOrder int

const (
	MultiSigPending MultiSignatureStatus = iota
	MultiSigPartial
	MultiSigComplete
	MultiSigExpired
	
	ParallelSigning SigningOrder = iota
	SequentialSigning
)

// Utility functions

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func generateSerialNumber() *big.Int {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(128), nil).Sub(max, big.NewInt(1))
	n, _ := rand.Int(rand.Reader, max)
	return n
}

func applyProfile(cert *x509.Certificate, profile *CertificateProfile) {
	cert.KeyUsage = profile.KeyUsage
	cert.ExtKeyUsage = profile.ExtKeyUsage
	if profile.IsCA {
		cert.IsCA = true
		cert.MaxPathLen = profile.MaxPathLen
	}
}