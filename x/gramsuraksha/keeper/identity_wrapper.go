package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	identitykeeper "github.com/namo/x/identity/keeper"
	identitytypes "github.com/namo/x/identity/types"
)

// identityKeeperWrapper wraps the generic interface to provide type-safe access
type identityKeeperWrapper struct {
	keeper interface {
		GetIdentity(sdk.Context, string) (interface{}, bool)
		SetIdentity(sdk.Context, interface{})
		GetCredential(sdk.Context, string) (interface{}, bool)
		SetCredential(sdk.Context, interface{})
		GetCredentialsBySubject(sdk.Context, string) []string
		AddCredentialToSubject(sdk.Context, string, string)
		UpdateCredentialStatus(sdk.Context, string, interface{}) error
	}
}

// Implement the identitykeeper.Keeper interface methods

func (w *identityKeeperWrapper) GetIdentity(ctx sdk.Context, did string) (identitytypes.Identity, bool) {
	result, found := w.keeper.GetIdentity(ctx, did)
	if !found {
		return identitytypes.Identity{}, false
	}
	
	// Type assertion
	if identity, ok := result.(identitytypes.Identity); ok {
		return identity, true
	}
	
	return identitytypes.Identity{}, false
}

func (w *identityKeeperWrapper) SetIdentity(ctx sdk.Context, identity identitytypes.Identity) {
	w.keeper.SetIdentity(ctx, identity)
}

func (w *identityKeeperWrapper) GetCredential(ctx sdk.Context, id string) (identitytypes.VerifiableCredential, bool) {
	result, found := w.keeper.GetCredential(ctx, id)
	if !found {
		return identitytypes.VerifiableCredential{}, false
	}
	
	// Type assertion
	if credential, ok := result.(identitytypes.VerifiableCredential); ok {
		return credential, true
	}
	
	return identitytypes.VerifiableCredential{}, false
}

func (w *identityKeeperWrapper) SetCredential(ctx sdk.Context, credential identitytypes.VerifiableCredential) {
	w.keeper.SetCredential(ctx, credential)
}

func (w *identityKeeperWrapper) GetCredentialsBySubject(ctx sdk.Context, subject string) []string {
	return w.keeper.GetCredentialsBySubject(ctx, subject)
}

func (w *identityKeeperWrapper) AddCredentialToSubject(ctx sdk.Context, subject string, credentialID string) {
	w.keeper.AddCredentialToSubject(ctx, subject, credentialID)
}

func (w *identityKeeperWrapper) UpdateCredentialStatus(ctx sdk.Context, credentialID string, status *identitytypes.CredentialStatus) error {
	return w.keeper.UpdateCredentialStatus(ctx, credentialID, status)
}

// Additional methods that the identity keeper might have
func (w *identityKeeperWrapper) GetDIDDocument(ctx sdk.Context, did string) (identitytypes.DIDDocument, bool) {
	// This is a placeholder - implement based on actual identity keeper interface
	return identitytypes.DIDDocument{}, false
}

func (w *identityKeeperWrapper) SetDIDDocument(ctx sdk.Context, doc identitytypes.DIDDocument) error {
	// This is a placeholder - implement based on actual identity keeper interface
	return nil
}

func (w *identityKeeperWrapper) Logger(ctx sdk.Context) interface{} {
	// Return a logger interface
	return ctx.Logger()
}