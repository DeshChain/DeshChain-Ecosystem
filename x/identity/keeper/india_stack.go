package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// SetIndiaStackIntegration stores India Stack integration data
func (k Keeper) SetIndiaStackIntegration(ctx sdk.Context, integration types.IndiaStackIntegration) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndiaStackIntegPrefix)
	b := k.cdc.MustMarshal(&integration)
	store.Set([]byte(integration.UserAddress), b)
}

// GetIndiaStackIntegration retrieves India Stack integration data
func (k Keeper) GetIndiaStackIntegration(ctx sdk.Context, address string) (types.IndiaStackIntegration, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndiaStackIntegPrefix)
	b := store.Get([]byte(address))
	if b == nil {
		// Return empty integration if not found
		return types.IndiaStackIntegration{
			UserAddress:  address,
			LastUpdated: ctx.BlockTime(),
		}, false
	}
	
	var integration types.IndiaStackIntegration
	k.cdc.MustUnmarshal(b, &integration)
	return integration, true
}

// IterateIndiaStackIntegrations iterates over all India Stack integrations
func (k Keeper) IterateIndiaStackIntegrations(ctx sdk.Context, cb func(integration types.IndiaStackIntegration) bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.IndiaStackIntegPrefix)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var integration types.IndiaStackIntegration
		k.cdc.MustUnmarshal(iterator.Value(), &integration)
		if cb(integration) {
			break
		}
	}
}

// LinkAadhaar links Aadhaar to user identity
func (k Keeper) LinkAadhaar(
	ctx sdk.Context,
	userAddress string,
	aadhaarHash string,
	demographicHash string,
	biometricHash string,
	verificationMethod string,
	consentArtefact string,
) error {
	// Check if India Stack is enabled
	if !k.EnableIndiaStack(ctx) {
		return types.ErrInvalidRequest
	}
	
	// Validate user has identity
	if !k.HasIdentity(ctx, userAddress) {
		return types.ErrIdentityNotFound
	}
	
	// Get or create India Stack integration
	integration, _ := k.GetIndiaStackIntegration(ctx, userAddress)
	
	// Check if already linked
	if integration.AadhaarLinked {
		return types.ErrInvalidRequest
	}
	
	// Create Aadhaar credential
	aadhaarCred := types.AadhaarCredential{
		ID:                 fmt.Sprintf("aadhaar:%s:%s", userAddress, sdk.NewRand().Str(16)),
		AadhaarHash:        aadhaarHash,
		DemographicHash:    demographicHash,
		BiometricHash:      biometricHash,
		VerificationMethod: verificationMethod,
		VerificationScore:  0.95, // Mock score, would come from actual verification
		ConsentArtefact:    consentArtefact,
		IssuedAt:           ctx.BlockTime(),
		ExpiresAt:          ctx.BlockTime().AddDate(0, 6, 0), // 6 months validity
		IssuerDID:          "did:desh:government-issuer",
	}
	
	// Validate Aadhaar credential
	if err := types.ValidateAadhaarCredential(&aadhaarCred); err != nil {
		return err
	}
	
	// Store Aadhaar credential
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AadhaarCredPrefix)
	b := k.cdc.MustMarshal(&aadhaarCred)
	store.Set([]byte(aadhaarCred.ID), b)
	
	// Update integration
	integration.AadhaarLinked = true
	integration.AadhaarCredentialID = aadhaarCred.ID
	integration.LastUpdated = ctx.BlockTime()
	k.SetIndiaStackIntegration(ctx, integration)
	
	// Issue verifiable credential
	claims := map[string]interface{}{
		"aadhaarVerified": true,
		"verificationMethod": verificationMethod,
		"verificationScore": aadhaarCred.VerificationScore,
	}
	
	_, err := k.IssueCredential(
		ctx,
		"did:desh:government-issuer",
		userAddress,
		types.DocTypeAadhaar,
		claims,
		180, // 6 months
	)
	if err != nil {
		return err
	}
	
	// Update identity KYC status
	identity, _ := k.GetIdentity(ctx, userAddress)
	identity.KYCStatus = types.KYCStatus{
		Level:        types.KYCLevel_ENHANCED,
		Status:       types.VerificationStatus_VERIFIED,
		VerifiedAt:   ctx.BlockTime(),
		ExpiresAt:    ctx.BlockTime().AddDate(0, 6, 0),
		Verifier:     "Aadhaar e-KYC",
		CredentialID: aadhaarCred.ID,
	}
	k.SetIdentity(ctx, identity)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAadhaarLinked,
			sdk.NewAttribute(types.AttributeKeyAddress, userAddress),
			sdk.NewAttribute(types.AttributeKeyVerificationMethod, verificationMethod),
		),
	)
	
	return nil
}

// ConnectDigiLocker connects DigiLocker account
func (k Keeper) ConnectDigiLocker(
	ctx sdk.Context,
	userAddress string,
	authToken string,
	consentID string,
	documentTypes []string,
) error {
	// Check if India Stack is enabled
	if !k.EnableIndiaStack(ctx) {
		return types.ErrInvalidRequest
	}
	
	// Validate user has identity
	if !k.HasIdentity(ctx, userAddress) {
		return types.ErrIdentityNotFound
	}
	
	// Get or create India Stack integration
	integration, _ := k.GetIndiaStackIntegration(ctx, userAddress)
	
	// Mock DigiLocker authentication (in production, would call actual service)
	userID := fmt.Sprintf("user_%s", sdk.NewRand().Str(16))
	
	// Fetch and store documents
	for _, docType := range documentTypes {
		doc := types.DigiLockerDocument{
			ID:               fmt.Sprintf("digilocker:%s:%s", userAddress, docType),
			DocumentType:     docType,
			DocumentURI:      fmt.Sprintf("https://digilocker.gov.in/api/v1/documents/%s", sdk.NewRand().Str(32)),
			DocumentHash:     sdk.NewRand().Str(64),
			IssuerOrg:        k.getDocumentIssuer(docType),
			IssuedOn:         ctx.BlockTime().AddDate(-1, 0, 0), // Mock: issued 1 year ago
			VerificationHash: sdk.NewRand().Str(64),
			ConsentID:        consentID,
		}
		
		// Store document reference
		store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DigiLockerDocPrefix)
		b := k.cdc.MustMarshal(&doc)
		store.Set([]byte(doc.ID), b)
		
		// Add to integration
		integration.DigiLockerDocuments = append(integration.DigiLockerDocuments, doc.ID)
		
		// Issue verifiable credential for document
		claims := map[string]interface{}{
			"documentType": docType,
			"documentHash": doc.DocumentHash,
			"issuerOrg":    doc.IssuerOrg,
		}
		
		_, err := k.IssueCredential(
			ctx,
			"did:desh:government-issuer",
			userAddress,
			docType,
			claims,
			365, // 1 year
		)
		if err != nil {
			return err
		}
	}
	
	// Update integration
	integration.DigiLockerLinked = true
	integration.LastUpdated = ctx.BlockTime()
	k.SetIndiaStackIntegration(ctx, integration)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDigiLockerConnected,
			sdk.NewAttribute(types.AttributeKeyAddress, userAddress),
			sdk.NewAttribute("document_count", fmt.Sprintf("%d", len(documentTypes))),
		),
	)
	
	return nil
}

// LinkUPI links UPI ID to identity
func (k Keeper) LinkUPI(
	ctx sdk.Context,
	userAddress string,
	vpaURI string,
	pspProvider string,
	authToken string,
) error {
	// Check if India Stack is enabled
	if !k.EnableIndiaStack(ctx) {
		return types.ErrInvalidRequest
	}
	
	// Validate user has identity
	if !k.HasIdentity(ctx, userAddress) {
		return types.ErrIdentityNotFound
	}
	
	// Get or create India Stack integration
	integration, _ := k.GetIndiaStackIntegration(ctx, userAddress)
	
	// Create UPI identity
	upiIdentity := types.UPIIdentity{
		ID:             fmt.Sprintf("upi:%s:%s", userAddress, sdk.NewRand().Str(16)),
		VPAURI:         vpaURI,
		PSPProvider:    pspProvider,
		IsActive:       true,
		CreatedAt:      ctx.BlockTime(),
		LastUsedAt:     ctx.BlockTime(),
	}
	
	// Store UPI identity
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UPIIdentityPrefix)
	b := k.cdc.MustMarshal(&upiIdentity)
	store.Set([]byte(upiIdentity.ID), b)
	
	// Update integration
	integration.UPILinked = true
	integration.UPIIdentities = append(integration.UPIIdentities, upiIdentity.ID)
	integration.LastUpdated = ctx.BlockTime()
	k.SetIndiaStackIntegration(ctx, integration)
	
	// Issue verifiable credential
	claims := map[string]interface{}{
		"vpa":         vpaURI,
		"pspProvider": pspProvider,
		"verified":    true,
	}
	
	_, err := k.IssueCredential(
		ctx,
		"did:desh:payment-issuer",
		userAddress,
		"PaymentIdentity",
		claims,
		365, // 1 year
	)
	if err != nil {
		return err
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUPILinked,
			sdk.NewAttribute(types.AttributeKeyAddress, userAddress),
			sdk.NewAttribute("vpa", vpaURI),
		),
	)
	
	return nil
}

// CreateDEPAConsent creates a DEPA consent artefact
func (k Keeper) CreateDEPAConsent(
	ctx sdk.Context,
	dataPrincipal string,
	dataController string,
	purpose string,
	dataTypes []string,
	expiryDays int32,
) (*types.DEPAConsent, error) {
	// Validate data principal has identity
	if !k.HasIdentity(ctx, dataPrincipal) {
		return nil, types.ErrIdentityNotFound
	}
	
	// Create consent
	consent := types.DEPAConsent{
		ID:                fmt.Sprintf("depa:%s:%s", dataPrincipal, sdk.NewRand().Str(16)),
		ConsentArtefactID: fmt.Sprintf("CA_%s", sdk.NewRand().Str(32)),
		DataPrincipal:     dataPrincipal,
		DataController:    dataController,
		Purpose:           purpose,
		DataTypes:         dataTypes,
		Frequency:         "ON_DEMAND",
		ConsentGivenAt:    ctx.BlockTime(),
		ConsentExpiresAt:  ctx.BlockTime().AddDate(0, 0, int(expiryDays)),
		Status:            types.ConsentStatusActive,
		RevocationAllowed: true,
	}
	
	// Validate consent
	if err := types.ValidateDEPAConsent(&consent); err != nil {
		return nil, err
	}
	
	// Store consent
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DEPAConsentPrefix)
	b := k.cdc.MustMarshal(&consent)
	store.Set([]byte(consent.ID), b)
	
	// Update integration
	integration, _ := k.GetIndiaStackIntegration(ctx, dataPrincipal)
	integration.DEPAConsents = append(integration.DEPAConsents, consent.ID)
	integration.LastUpdated = ctx.BlockTime()
	k.SetIndiaStackIntegration(ctx, integration)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDEPAConsentGiven,
			sdk.NewAttribute(types.AttributeKeyAddress, dataPrincipal),
			sdk.NewAttribute(types.AttributeKeyConsentArtefact, consent.ConsentArtefactID),
		),
	)
	
	return &consent, nil
}

// CreateVillagePanchayatKYC creates a village-level KYC verification
func (k Keeper) CreateVillagePanchayatKYC(
	ctx sdk.Context,
	subjectAddress string,
	panchayatCode string,
	verifierName string,
	verifierDesignation string,
	verifierID string,
	documentsVerified []string,
) (*types.VillagePanchayatKYC, error) {
	// Validate subject has identity
	if !k.HasIdentity(ctx, subjectAddress) {
		return nil, types.ErrIdentityNotFound
	}
	
	// Get subject identity for name
	identity, _ := k.GetIdentity(ctx, subjectAddress)
	
	// Create panchayat KYC
	panchayatKYC := types.VillagePanchayatKYC{
		ID:                  fmt.Sprintf("panchayat:%s:%s", subjectAddress, sdk.NewRand().Str(16)),
		PanchayatCode:       panchayatCode,
		PanchayatName:       k.getPanchayatName(panchayatCode), // Mock function
		VerifierName:        verifierName,
		VerifierDesignation: verifierDesignation,
		VerifierID:          verifierID,
		SubjectName:         "Subject Name", // Would come from identity
		SubjectAddress:      subjectAddress,
		VerificationType:    "IN_PERSON",
		DocumentsVerified:   documentsVerified,
		PhotoHash:           sdk.NewRand().Str(64),
		Remarks:             "Verified in person at Gram Panchayat office",
		VerifiedAt:          ctx.BlockTime(),
		ValidUntil:          ctx.BlockTime().AddDate(1, 0, 0), // 1 year
		QRCode:              fmt.Sprintf("QR_%s", sdk.NewRand().Str(32)),
	}
	
	// Store panchayat KYC
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.VillagePanchayatPrefix)
	b := k.cdc.MustMarshal(&panchayatKYC)
	store.Set([]byte(panchayatKYC.ID), b)
	
	// Update integration
	integration, _ := k.GetIndiaStackIntegration(ctx, subjectAddress)
	integration.VillagePanchayatKYCs = append(integration.VillagePanchayatKYCs, panchayatKYC.ID)
	integration.LastUpdated = ctx.BlockTime()
	k.SetIndiaStackIntegration(ctx, integration)
	
	// Issue verifiable credential
	claims := map[string]interface{}{
		"panchayatCode":    panchayatCode,
		"verificationType": "IN_PERSON",
		"verifierName":     verifierName,
		"documentsVerified": documentsVerified,
	}
	
	_, err := k.IssueCredential(
		ctx,
		"did:desh:panchayat-issuer",
		subjectAddress,
		"VillagePanchayatKYC",
		claims,
		365, // 1 year
	)
	if err != nil {
		return nil, err
	}
	
	// Update identity KYC status
	identity.KYCStatus = types.KYCStatus{
		Level:        types.KYCLevel_BASIC,
		Status:       types.VerificationStatus_VERIFIED,
		VerifiedAt:   ctx.BlockTime(),
		ExpiresAt:    ctx.BlockTime().AddDate(1, 0, 0),
		Verifier:     fmt.Sprintf("Gram Panchayat %s", panchayatCode),
		CredentialID: panchayatKYC.ID,
	}
	k.SetIdentity(ctx, identity)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePanchayatKYCDone,
			sdk.NewAttribute(types.AttributeKeyAddress, subjectAddress),
			sdk.NewAttribute(types.AttributeKeyPanchayatCode, panchayatCode),
		),
	)
	
	return &panchayatKYC, nil
}

// Helper functions

func (k Keeper) getDocumentIssuer(docType string) string {
	issuers := map[string]string{
		types.DocTypePAN:            "Income Tax Department",
		types.DocTypeDrivingLicense: "Regional Transport Office",
		types.DocTypeVoterID:        "Election Commission of India",
		types.DocTypePassport:       "Ministry of External Affairs",
		types.DocTypeRationCard:     "Food and Civil Supplies Department",
	}
	
	if issuer, ok := issuers[docType]; ok {
		return issuer
	}
	return "Government of India"
}

func (k Keeper) getPanchayatName(panchayatCode string) string {
	// In production, this would look up from a registry
	return fmt.Sprintf("Gram Panchayat %s", panchayatCode)
}