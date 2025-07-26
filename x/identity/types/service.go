package types

import (
	"context"
)

// This file is a placeholder for the protobuf-generated service descriptor
// In a real implementation, this would be generated from .proto files

// _Msg_serviceDesc is a placeholder for the gRPC service descriptor
var _Msg_serviceDesc = struct{}{}

// MsgServer is the server API for Msg service
type MsgServer interface {
	CreateIdentity(context.Context, *MsgCreateIdentity) (*MsgCreateIdentityResponse, error)
	UpdateIdentity(context.Context, *MsgUpdateIdentity) (*MsgUpdateIdentityResponse, error)
	RevokeIdentity(context.Context, *MsgRevokeIdentity) (*MsgRevokeIdentityResponse, error)
	RegisterDID(context.Context, *MsgRegisterDID) (*MsgRegisterDIDResponse, error)
	UpdateDID(context.Context, *MsgUpdateDID) (*MsgUpdateDIDResponse, error)
	DeactivateDID(context.Context, *MsgDeactivateDID) (*MsgDeactivateDIDResponse, error)
	IssueCredential(context.Context, *MsgIssueCredential) (*MsgIssueCredentialResponse, error)
	RevokeCredential(context.Context, *MsgRevokeCredential) (*MsgRevokeCredentialResponse, error)
	PresentCredential(context.Context, *MsgPresentCredential) (*MsgPresentCredentialResponse, error)
	CreateZKProof(context.Context, *MsgCreateZKProof) (*MsgCreateZKProofResponse, error)
	VerifyZKProof(context.Context, *MsgVerifyZKProof) (*MsgVerifyZKProofResponse, error)
	LinkAadhaar(context.Context, *MsgLinkAadhaar) (*MsgLinkAadhaarResponse, error)
	ConnectDigiLocker(context.Context, *MsgConnectDigiLocker) (*MsgConnectDigiLockerResponse, error)
	LinkUPI(context.Context, *MsgLinkUPI) (*MsgLinkUPIResponse, error)
	GiveConsent(context.Context, *MsgGiveConsent) (*MsgGiveConsentResponse, error)
	WithdrawConsent(context.Context, *MsgWithdrawConsent) (*MsgWithdrawConsentResponse, error)
	AddRecoveryMethod(context.Context, *MsgAddRecoveryMethod) (*MsgAddRecoveryMethodResponse, error)
	InitiateRecovery(context.Context, *MsgInitiateRecovery) (*MsgInitiateRecoveryResponse, error)
	CompleteRecovery(context.Context, *MsgCompleteRecovery) (*MsgCompleteRecoveryResponse, error)
	UpdatePrivacySettings(context.Context, *MsgUpdatePrivacySettings) (*MsgUpdatePrivacySettingsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations
type UnimplementedMsgServer struct{}

func (*UnimplementedMsgServer) CreateIdentity(context.Context, *MsgCreateIdentity) (*MsgCreateIdentityResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) UpdateIdentity(context.Context, *MsgUpdateIdentity) (*MsgUpdateIdentityResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) RevokeIdentity(context.Context, *MsgRevokeIdentity) (*MsgRevokeIdentityResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) RegisterDID(context.Context, *MsgRegisterDID) (*MsgRegisterDIDResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) UpdateDID(context.Context, *MsgUpdateDID) (*MsgUpdateDIDResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) DeactivateDID(context.Context, *MsgDeactivateDID) (*MsgDeactivateDIDResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) IssueCredential(context.Context, *MsgIssueCredential) (*MsgIssueCredentialResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) RevokeCredential(context.Context, *MsgRevokeCredential) (*MsgRevokeCredentialResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) PresentCredential(context.Context, *MsgPresentCredential) (*MsgPresentCredentialResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) CreateZKProof(context.Context, *MsgCreateZKProof) (*MsgCreateZKProofResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) VerifyZKProof(context.Context, *MsgVerifyZKProof) (*MsgVerifyZKProofResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) LinkAadhaar(context.Context, *MsgLinkAadhaar) (*MsgLinkAadhaarResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) ConnectDigiLocker(context.Context, *MsgConnectDigiLocker) (*MsgConnectDigiLockerResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) LinkUPI(context.Context, *MsgLinkUPI) (*MsgLinkUPIResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) GiveConsent(context.Context, *MsgGiveConsent) (*MsgGiveConsentResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) WithdrawConsent(context.Context, *MsgWithdrawConsent) (*MsgWithdrawConsentResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) AddRecoveryMethod(context.Context, *MsgAddRecoveryMethod) (*MsgAddRecoveryMethodResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) InitiateRecovery(context.Context, *MsgInitiateRecovery) (*MsgInitiateRecoveryResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) CompleteRecovery(context.Context, *MsgCompleteRecovery) (*MsgCompleteRecoveryResponse, error) {
	return nil, ErrInvalidRequest
}

func (*UnimplementedMsgServer) UpdatePrivacySettings(context.Context, *MsgUpdatePrivacySettings) (*MsgUpdatePrivacySettingsResponse, error) {
	return nil, ErrInvalidRequest
}