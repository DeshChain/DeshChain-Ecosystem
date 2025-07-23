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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "deshgovernance"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for governance
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_deshgovernance"
)

// Store key prefixes
var (
	// Key for governance phase
	KeyGovernancePhase = []byte{0x01}
	
	// Key for founder address
	KeyFounderAddress = []byte{0x02}
	
	// Key for phase transition time
	KeyPhaseTransitionTime = []byte{0x03}
	
	// Key for protected parameters
	KeyProtectedParams = []byte{0x04}
	
	// Key for proposal veto
	KeyProposalVeto = []byte{0x05}
	
	// Key for founder approval
	KeyFounderApproval = []byte{0x06}
	
	// Key for supermajority proposals
	KeySupermajorityProposal = []byte{0x07}
	
	// Key for genesis time
	KeyGenesisTime = []byte{0x08}
	
	// Key for charity allocation params
	KeyCharityAllocationParams = []byte{0x09}
	
	// Key for last charity update year
	KeyLastCharityUpdateYear = []byte{0x0A}
)

// Governance phases
type GovernancePhase int32

const (
	// Phase 1: Founder Control (0-2 years)
	GovernancePhase_FOUNDER_CONTROL GovernancePhase = 0
	
	// Phase 2: Shared Governance (2-3 years)
	GovernancePhase_SHARED_GOVERNANCE GovernancePhase = 1
	
	// Phase 3: Community Governance (3+ years)
	GovernancePhase_COMMUNITY_GOVERNANCE GovernancePhase = 2
)

// Protection types for parameters
type ProtectionType int32

const (
	// No special protection
	ProtectionType_NONE ProtectionType = 0
	
	// Immutable - can never be changed
	ProtectionType_IMMUTABLE ProtectionType = 1
	
	// Requires founder consent
	ProtectionType_FOUNDER_CONSENT ProtectionType = 2
	
	// Requires supermajority (80%)
	ProtectionType_SUPERMAJORITY ProtectionType = 3
)

// GetProposalVetoKey returns the key for a proposal veto
func GetProposalVetoKey(proposalID uint64) []byte {
	return append(KeyProposalVeto, sdk.Uint64ToBigEndian(proposalID)...)
}

// GetFounderApprovalKey returns the key for founder approval
func GetFounderApprovalKey(proposalID uint64) []byte {
	return append(KeyFounderApproval, sdk.Uint64ToBigEndian(proposalID)...)
}