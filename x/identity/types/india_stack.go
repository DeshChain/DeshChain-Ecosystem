package types

import (
	"fmt"
	"time"
)

// AadhaarCredential represents Aadhaar-based identity credential
type AadhaarCredential struct {
	ID                   string            `json:"id"`
	AadhaarHash          string            `json:"aadhaar_hash"`         // Hashed Aadhaar number
	DemographicHash      string            `json:"demographic_hash"`     // Hash of demographic data
	BiometricHash        string            `json:"biometric_hash"`       // Hash of biometric data
	VerificationMethod   string            `json:"verification_method"`  // OTP, biometric, demographic
	VerificationScore    float64           `json:"verification_score"`
	ConsentArtefact      string            `json:"consent_artefact"`     // DEPA consent
	IssuedAt             time.Time         `json:"issued_at"`
	ExpiresAt            time.Time         `json:"expires_at"`
	IssuerDID            string            `json:"issuer_did"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

// DigiLockerDocument represents a document from DigiLocker
type DigiLockerDocument struct {
	ID               string            `json:"id"`
	DocumentType     string            `json:"document_type"`
	DocumentURI      string            `json:"document_uri"`
	DocumentHash     string            `json:"document_hash"`
	IssuerOrg        string            `json:"issuer_org"`
	IssuedOn         time.Time         `json:"issued_on"`
	ValidUntil       *time.Time        `json:"valid_until,omitempty"`
	VerificationHash string            `json:"verification_hash"`
	ConsentID        string            `json:"consent_id"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// UPIIdentity represents UPI-based payment identity
type UPIIdentity struct {
	ID                string            `json:"id"`
	VPAURI            string            `json:"vpa_uri"`              // Virtual Payment Address
	LinkedAccounts    []string          `json:"linked_accounts"`      // Linked bank accounts (encrypted)
	VerificationLevel string            `json:"verification_level"`
	PSPProvider       string            `json:"psp_provider"`
	IsActive          bool              `json:"is_active"`
	CreatedAt         time.Time         `json:"created_at"`
	LastUsedAt        time.Time         `json:"last_used_at"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// eKYCData represents data from Aadhaar e-KYC
type eKYCData struct {
	ReferenceID      string            `json:"reference_id"`
	Name             string            `json:"name"`
	DOB              string            `json:"dob"`
	Gender           string            `json:"gender"`
	Address          eKYCAddress       `json:"address"`
	PhotoHash        string            `json:"photo_hash"`
	MobileHash       string            `json:"mobile_hash"`
	EmailHash        string            `json:"email_hash"`
	ShareCode        string            `json:"share_code"`
	Timestamp        time.Time         `json:"timestamp"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// eKYCAddress represents address from e-KYC
type eKYCAddress struct {
	CareOf      string `json:"care_of"`
	House       string `json:"house"`
	Street      string `json:"street"`
	Landmark    string `json:"landmark"`
	Locality    string `json:"locality"`
	VTC         string `json:"vtc"`
	District    string `json:"district"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PinCode     string `json:"pin_code"`
}

// DEPAConsent represents Data Empowerment and Protection Architecture consent
type DEPAConsent struct {
	ID                   string            `json:"id"`
	ConsentArtefactID    string            `json:"consent_artefact_id"`
	DataPrincipal        string            `json:"data_principal"`       // User
	DataController       string            `json:"data_controller"`      // Service provider
	DataProcessors       []string          `json:"data_processors"`      // Third parties
	Purpose              string            `json:"purpose"`
	DataTypes            []string          `json:"data_types"`
	Frequency            string            `json:"frequency"`
	ConsentGivenAt       time.Time         `json:"consent_given_at"`
	ConsentExpiresAt     time.Time         `json:"consent_expires_at"`
	Status               string            `json:"status"`
	RevocationAllowed    bool              `json:"revocation_allowed"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

// AccountAggregatorConsent represents AA framework consent
type AccountAggregatorConsent struct {
	ID                string            `json:"id"`
	ConsentHandle     string            `json:"consent_handle"`
	ConsentID         string            `json:"consent_id"`
	CustomerID        string            `json:"customer_id"`
	FIUEntity         string            `json:"fiu_entity"`          // Financial Information User
	FIPEntities       []string          `json:"fip_entities"`        // Financial Information Providers
	ConsentTypes      []string          `json:"consent_types"`
	FITypes           []string          `json:"fi_types"`            // Financial information types
	ConsentStart      time.Time         `json:"consent_start"`
	ConsentExpiry     time.Time         `json:"consent_expiry"`
	ConsentMode       string            `json:"consent_mode"`
	FetchType         string            `json:"fetch_type"`
	Frequency         ConsentFrequency  `json:"frequency"`
	DataLife          DataLifeUnit      `json:"data_life"`
	Status            string            `json:"status"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// ConsentFrequency represents how often data can be fetched
type ConsentFrequency struct {
	Unit  string `json:"unit"`  // HOURLY, DAILY, MONTHLY, etc.
	Value int32  `json:"value"`
}

// DataLifeUnit represents how long data can be stored
type DataLifeUnit struct {
	Unit  string `json:"unit"`  // DAY, MONTH, YEAR
	Value int32  `json:"value"`
}

// VillagePanchayatKYC represents village-level KYC verification
type VillagePanchayatKYC struct {
	ID                  string            `json:"id"`
	PanchayatCode       string            `json:"panchayat_code"`
	PanchayatName       string            `json:"panchayat_name"`
	VerifierName        string            `json:"verifier_name"`
	VerifierDesignation string            `json:"verifier_designation"`
	VerifierID          string            `json:"verifier_id"`
	SubjectName         string            `json:"subject_name"`
	SubjectAddress      string            `json:"subject_address"`
	VerificationType    string            `json:"verification_type"`
	DocumentsVerified   []string          `json:"documents_verified"`
	PhotoHash           string            `json:"photo_hash"`
	Remarks             string            `json:"remarks"`
	VerifiedAt          time.Time         `json:"verified_at"`
	ValidUntil          time.Time         `json:"valid_until"`
	QRCode              string            `json:"qr_code"`
	Metadata            map[string]string `json:"metadata,omitempty"`
}

// JANDhanAccount represents Jan Dhan account details
type JANDhanAccount struct {
	ID                string            `json:"id"`
	AccountNumber     string            `json:"account_number"`      // Encrypted
	BankName          string            `json:"bank_name"`
	BranchCode        string            `json:"branch_code"`
	LinkedAadhaar     bool              `json:"linked_aadhaar"`
	RuPayCardIssued   bool              `json:"rupay_card_issued"`
	OverdraftEligible bool              `json:"overdraft_eligible"`
	InsuranceCovered  bool              `json:"insurance_covered"`
	OpenedDate        time.Time         `json:"opened_date"`
	Status            string            `json:"status"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// IndiaStackIntegration manages all India Stack integrations
type IndiaStackIntegration struct {
	UserAddress          string                    `json:"user_address"`
	AadhaarLinked        bool                      `json:"aadhaar_linked"`
	AadhaarCredentialID  string                    `json:"aadhaar_credential_id,omitempty"`
	DigiLockerLinked     bool                      `json:"digilocker_linked"`
	DigiLockerDocuments  []string                  `json:"digilocker_documents,omitempty"`
	UPILinked            bool                      `json:"upi_linked"`
	UPIIdentities        []string                  `json:"upi_identities,omitempty"`
	JANDhanLinked        bool                      `json:"jan_dhan_linked"`
	JANDhanAccounts      []string                  `json:"jan_dhan_accounts,omitempty"`
	DEPAConsents         []string                  `json:"depa_consents,omitempty"`
	AAConsents           []string                  `json:"aa_consents,omitempty"`
	VillagePanchayatKYCs []string                  `json:"village_panchayat_kycs,omitempty"`
	LastUpdated          time.Time                 `json:"last_updated"`
	Metadata             map[string]string         `json:"metadata,omitempty"`
}

// ValidateAadhaarCredential validates Aadhaar credential
func ValidateAadhaarCredential(cred *AadhaarCredential) error {
	if cred.AadhaarHash == "" {
		return fmt.Errorf("aadhaar hash is required")
	}
	
	if cred.ConsentArtefact == "" {
		return fmt.Errorf("consent artefact is required")
	}
	
	if cred.VerificationScore < 0 || cred.VerificationScore > 1 {
		return fmt.Errorf("verification score must be between 0 and 1")
	}
	
	if cred.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("credential has expired")
	}
	
	return nil
}

// ValidateDEPAConsent validates DEPA consent
func ValidateDEPAConsent(consent *DEPAConsent) error {
	if consent.ConsentArtefactID == "" {
		return fmt.Errorf("consent artefact ID is required")
	}
	
	if consent.DataPrincipal == "" || consent.DataController == "" {
		return fmt.Errorf("data principal and controller are required")
	}
	
	if len(consent.DataTypes) == 0 {
		return fmt.Errorf("at least one data type must be specified")
	}
	
	if consent.ConsentExpiresAt.Before(time.Now()) {
		return fmt.Errorf("consent has expired")
	}
	
	return nil
}

// India Stack specific constants
const (
	// Document types
	DocTypeAadhaar      = "AADHAAR"
	DocTypePAN          = "PAN"
	DocTypeDrivingLicense = "DRIVING_LICENSE"
	DocTypeVoterID      = "VOTER_ID"
	DocTypePassport     = "PASSPORT"
	DocTypeRationCard   = "RATION_CARD"
	
	// Verification methods
	VerificationMethodOTP       = "OTP"
	VerificationMethodBiometric = "BIOMETRIC"
	VerificationMethodDemographic = "DEMOGRAPHIC"
	VerificationMethodOffline   = "OFFLINE"
	
	// Consent status
	ConsentStatusActive   = "ACTIVE"
	ConsentStatusExpired  = "EXPIRED"
	ConsentStatusRevoked  = "REVOKED"
	ConsentStatusPaused   = "PAUSED"
	
	// FI Types for Account Aggregator
	FITypeDeposit     = "DEPOSIT"
	FITypeTermDeposit = "TERM_DEPOSIT"
	FITypeRecurringDeposit = "RECURRING_DEPOSIT"
	FITypeSIP         = "SIP"
	FITypeMutualFunds = "MUTUAL_FUNDS"
	FITypeInsurance   = "INSURANCE"
	FITypeGST         = "GST"
)

// India Stack events
const (
	EventTypeAadhaarLinked       = "aadhaar_linked"
	EventTypeAadhaarVerified     = "aadhaar_verified"
	EventTypeDigiLockerConnected = "digilocker_connected"
	EventTypeDocumentFetched     = "document_fetched"
	EventTypeUPILinked           = "upi_linked"
	EventTypeDEPAConsentGiven    = "depa_consent_given"
	EventTypeAAConsentCreated    = "aa_consent_created"
	EventTypePanchayatKYCDone    = "panchayat_kyc_done"
)

// India Stack attributes
const (
	AttributeKeyDocumentType    = "document_type"
	AttributeKeyVerificationMethod = "verification_method"
	AttributeKeyConsentArtefact = "consent_artefact"
	AttributeKeyFIUEntity       = "fiu_entity"
	AttributeKeyPanchayatCode   = "panchayat_code"
)