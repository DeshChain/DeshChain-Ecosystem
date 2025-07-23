/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"context"
)

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// VetoProposal allows the founder to veto a proposal
	VetoProposal(context.Context, *MsgVetoProposal) (*MsgVetoProposalResponse, error)
	
	// ApproveFounderConsentProposal allows the founder to approve proposals requiring consent
	ApproveFounderConsentProposal(context.Context, *MsgApproveFounderConsentProposal) (*MsgApproveFounderConsentProposalResponse, error)
	
	// UpdateProtectedParameter allows updating protection levels of parameters
	UpdateProtectedParameter(context.Context, *MsgUpdateProtectedParameter) (*MsgUpdateProtectedParameterResponse, error)
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// GovernanceInfo queries the current governance information
	GovernanceInfo(context.Context, *QueryGovernanceInfoRequest) (*QueryGovernanceInfoResponse, error)
	
	// ProtectedParameter queries a specific protected parameter
	ProtectedParameter(context.Context, *QueryProtectedParameterRequest) (*QueryProtectedParameterResponse, error)
	
	// ProtectedParameters queries all protected parameters
	ProtectedParameters(context.Context, *QueryProtectedParametersRequest) (*QueryProtectedParametersResponse, error)
	
	// ProposalVetoStatus queries if a proposal has been vetoed
	ProposalVetoStatus(context.Context, *QueryProposalVetoStatusRequest) (*QueryProposalVetoStatusResponse, error)
	
	// VetoedProposals queries all vetoed proposal IDs
	VetoedProposals(context.Context, *QueryVetoedProposalsRequest) (*QueryVetoedProposalsResponse, error)
	
	// ProposalFounderApproval queries if a proposal has founder approval
	ProposalFounderApproval(context.Context, *QueryProposalFounderApprovalRequest) (*QueryProposalFounderApprovalResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct{}

func (*UnimplementedMsgServer) VetoProposal(context.Context, *MsgVetoProposal) (*MsgVetoProposalResponse, error) {
	return nil, nil
}

func (*UnimplementedMsgServer) ApproveFounderConsentProposal(context.Context, *MsgApproveFounderConsentProposal) (*MsgApproveFounderConsentProposalResponse, error) {
	return nil, nil
}

func (*UnimplementedMsgServer) UpdateProtectedParameter(context.Context, *MsgUpdateProtectedParameter) (*MsgUpdateProtectedParameterResponse, error) {
	return nil, nil
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct{}

func (*UnimplementedQueryServer) GovernanceInfo(context.Context, *QueryGovernanceInfoRequest) (*QueryGovernanceInfoResponse, error) {
	return nil, nil
}

func (*UnimplementedQueryServer) ProtectedParameter(context.Context, *QueryProtectedParameterRequest) (*QueryProtectedParameterResponse, error) {
	return nil, nil
}

func (*UnimplementedQueryServer) ProtectedParameters(context.Context, *QueryProtectedParametersRequest) (*QueryProtectedParametersResponse, error) {
	return nil, nil
}

func (*UnimplementedQueryServer) ProposalVetoStatus(context.Context, *QueryProposalVetoStatusRequest) (*QueryProposalVetoStatusResponse, error) {
	return nil, nil
}

func (*UnimplementedQueryServer) VetoedProposals(context.Context, *QueryVetoedProposalsRequest) (*QueryVetoedProposalsResponse, error) {
	return nil, nil
}

func (*UnimplementedQueryServer) ProposalFounderApproval(context.Context, *QueryProposalFounderApprovalRequest) (*QueryProposalFounderApprovalResponse, error) {
	return nil, nil
}

// Msg service descriptor for code generation
var _Msg_serviceDesc = struct {
	ServiceName string
	Methods     []struct {
		MethodName string
	}
}{
	ServiceName: "deshchain.governance.v1.Msg",
	Methods: []struct {
		MethodName string
	}{
		{MethodName: "VetoProposal"},
		{MethodName: "ApproveFounderConsentProposal"},
		{MethodName: "UpdateProtectedParameter"},
	},
}