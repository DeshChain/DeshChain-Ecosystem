package identity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/identity/types"
)

// InitGenesis initializes the identity module's state from a provided genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module params
	k.SetParams(ctx, genState.Params)
	
	// Initialize identities
	for _, identity := range genState.Identities {
		k.SetIdentity(ctx, identity)
		
		// Create identity -> DID index
		if identity.DID != "" {
			k.SetIdentityDIDIndex(ctx, identity.DID, identity.Address)
		}
	}
	
	// Initialize DID documents
	for _, did := range genState.DIDDocuments {
		k.SetDIDDocument(ctx, did)
	}
	
	// Initialize credentials
	for _, credential := range genState.Credentials {
		k.SetCredential(ctx, credential)
		
		// Create credential indexes
		if holder, err := credential.GetSubjectID(); err == nil {
			k.SetCredentialHolderIndex(ctx, holder, credential.ID)
		}
		k.SetCredentialIssuerIndex(ctx, credential.Issuer, credential.ID)
		
		// Index by type
		for _, credType := range credential.Type {
			k.SetCredentialTypeIndex(ctx, credType, credential.ID)
		}
	}
	
	// Initialize credential schemas
	for _, schema := range genState.CredentialSchemas {
		k.SetCredentialSchema(ctx, schema)
	}
	
	// Initialize revocation lists
	for _, revList := range genState.RevocationLists {
		k.SetRevocationList(ctx, revList)
	}
	
	// Initialize privacy settings
	for _, settings := range genState.PrivacySettings {
		k.SetPrivacySettings(ctx, settings)
	}
	
	// Initialize India Stack integrations
	for _, integration := range genState.IndiaStackIntegrations {
		k.SetIndiaStackIntegration(ctx, integration)
	}
	
	// Initialize service registry
	for _, service := range genState.ServiceRegistry {
		k.RegisterService(ctx, service)
	}
	
	// Initialize issuer registry
	for _, issuer := range genState.IssuerRegistry {
		k.RegisterIssuer(ctx, issuer)
	}
}

// ExportGenesis returns the identity module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)
	
	// Export identities
	k.IterateIdentities(ctx, func(identity types.Identity) bool {
		genesis.Identities = append(genesis.Identities, identity)
		return false
	})
	
	// Export DID documents
	k.IterateDIDDocuments(ctx, func(did types.DIDDocument) bool {
		genesis.DIDDocuments = append(genesis.DIDDocuments, did)
		return false
	})
	
	// Export credentials
	k.IterateCredentials(ctx, func(credential types.VerifiableCredential) bool {
		genesis.Credentials = append(genesis.Credentials, credential)
		return false
	})
	
	// Export credential schemas
	k.IterateCredentialSchemas(ctx, func(schema types.CredentialSchema) bool {
		genesis.CredentialSchemas = append(genesis.CredentialSchemas, schema)
		return false
	})
	
	// Export revocation lists
	k.IterateRevocationLists(ctx, func(revList types.RevocationList) bool {
		genesis.RevocationLists = append(genesis.RevocationLists, revList)
		return false
	})
	
	// Export privacy settings
	k.IteratePrivacySettings(ctx, func(settings types.PrivacySettings) bool {
		genesis.PrivacySettings = append(genesis.PrivacySettings, settings)
		return false
	})
	
	// Export India Stack integrations
	k.IterateIndiaStackIntegrations(ctx, func(integration types.IndiaStackIntegration) bool {
		genesis.IndiaStackIntegrations = append(genesis.IndiaStackIntegrations, integration)
		return false
	})
	
	// Export service registry
	k.IterateServiceRegistry(ctx, func(service types.ServiceRegistryEntry) bool {
		genesis.ServiceRegistry = append(genesis.ServiceRegistry, service)
		return false
	})
	
	// Export issuer registry
	k.IterateIssuerRegistry(ctx, func(issuer types.IssuerRegistryEntry) bool {
		genesis.IssuerRegistry = append(genesis.IssuerRegistry, issuer)
		return false
	})
	
	return genesis
}