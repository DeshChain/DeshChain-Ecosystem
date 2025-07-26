package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers the necessary types and interfaces for the identity module
func RegisterCodec(cdc *codec.LegacyAmino) {
	// Messages
	cdc.RegisterConcrete(&MsgCreateIdentity{}, "identity/CreateIdentity", nil)
	cdc.RegisterConcrete(&MsgUpdateIdentity{}, "identity/UpdateIdentity", nil)
	cdc.RegisterConcrete(&MsgRevokeIdentity{}, "identity/RevokeIdentity", nil)
	cdc.RegisterConcrete(&MsgRegisterDID{}, "identity/RegisterDID", nil)
	cdc.RegisterConcrete(&MsgUpdateDID{}, "identity/UpdateDID", nil)
	cdc.RegisterConcrete(&MsgDeactivateDID{}, "identity/DeactivateDID", nil)
	cdc.RegisterConcrete(&MsgIssueCredential{}, "identity/IssueCredential", nil)
	cdc.RegisterConcrete(&MsgRevokeCredential{}, "identity/RevokeCredential", nil)
	cdc.RegisterConcrete(&MsgPresentCredential{}, "identity/PresentCredential", nil)
	cdc.RegisterConcrete(&MsgCreateZKProof{}, "identity/CreateZKProof", nil)
	cdc.RegisterConcrete(&MsgVerifyZKProof{}, "identity/VerifyZKProof", nil)
	cdc.RegisterConcrete(&MsgLinkAadhaar{}, "identity/LinkAadhaar", nil)
	cdc.RegisterConcrete(&MsgConnectDigiLocker{}, "identity/ConnectDigiLocker", nil)
	cdc.RegisterConcrete(&MsgLinkUPI{}, "identity/LinkUPI", nil)
	cdc.RegisterConcrete(&MsgGiveConsent{}, "identity/GiveConsent", nil)
	cdc.RegisterConcrete(&MsgWithdrawConsent{}, "identity/WithdrawConsent", nil)
	cdc.RegisterConcrete(&MsgAddRecoveryMethod{}, "identity/AddRecoveryMethod", nil)
	cdc.RegisterConcrete(&MsgInitiateRecovery{}, "identity/InitiateRecovery", nil)
	cdc.RegisterConcrete(&MsgCompleteRecovery{}, "identity/CompleteRecovery", nil)
	cdc.RegisterConcrete(&MsgUpdatePrivacySettings{}, "identity/UpdatePrivacySettings", nil)
}

// RegisterInterfaces registers the identity module types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Messages
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateIdentity{},
		&MsgUpdateIdentity{},
		&MsgRevokeIdentity{},
		&MsgRegisterDID{},
		&MsgUpdateDID{},
		&MsgDeactivateDID{},
		&MsgIssueCredential{},
		&MsgRevokeCredential{},
		&MsgPresentCredential{},
		&MsgCreateZKProof{},
		&MsgVerifyZKProof{},
		&MsgLinkAadhaar{},
		&MsgConnectDigiLocker{},
		&MsgLinkUPI{},
		&MsgGiveConsent{},
		&MsgWithdrawConsent{},
		&MsgAddRecoveryMethod{},
		&MsgInitiateRecovery{},
		&MsgCompleteRecovery{},
		&MsgUpdatePrivacySettings{},
	)
	
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}