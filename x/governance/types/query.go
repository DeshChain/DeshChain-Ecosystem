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
	"time"
)

// Query request types
type QueryGovernanceInfoRequest struct{}

type QueryProtectedParameterRequest struct {
	Name string `json:"name" yaml:"name"`
}

type QueryProtectedParametersRequest struct{}

type QueryProposalVetoStatusRequest struct {
	ProposalId uint64 `json:"proposal_id" yaml:"proposal_id"`
}

type QueryVetoedProposalsRequest struct{}

type QueryProposalFounderApprovalRequest struct {
	ProposalId uint64 `json:"proposal_id" yaml:"proposal_id"`
}

// Query response types
type QueryGovernanceInfoResponse struct {
	FounderAddress string           `json:"founder_address" yaml:"founder_address"`
	CurrentPhase   GovernancePhase  `json:"current_phase" yaml:"current_phase"`
	GenesisTime    time.Time        `json:"genesis_time" yaml:"genesis_time"`
}

type QueryProtectedParameterResponse struct {
	Parameter ProtectedParameter `json:"parameter" yaml:"parameter"`
}

type QueryProtectedParametersResponse struct {
	Parameters []ProtectedParameter `json:"parameters" yaml:"parameters"`
}

type QueryProposalVetoStatusResponse struct {
	IsVetoed bool `json:"is_vetoed" yaml:"is_vetoed"`
}

type QueryVetoedProposalsResponse struct {
	ProposalIds []uint64 `json:"proposal_ids" yaml:"proposal_ids"`
}

type QueryProposalFounderApprovalResponse struct {
	HasApproval bool `json:"has_approval" yaml:"has_approval"`
}