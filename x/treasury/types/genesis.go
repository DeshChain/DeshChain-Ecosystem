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
	"fmt"
	"time"
)

// DefaultGenesis returns the default treasury genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		CommunityFund:        DefaultCommunityFundState(),
		DevelopmentFund:      DefaultDevelopmentFundState(),
		MultiSigGovernances:  []MultiSigGovernance{},
		ProposalSystem:       DefaultProposalSystem(),
		Dashboard:            DefaultDashboard(),
		TransparencyReports:  []TransparencyReport{},
	}
}

// DefaultCommunityFundState returns default community fund state
func DefaultCommunityFundState() *CommunityFundState {
	return &CommunityFundState{
		Balance:      nil, // Will be initialized in InitGenesis
		Proposals:    []CommunityFundProposal{},
		Transactions: []CommunityFundTransaction{},
		Governance:   DefaultCommunityGovernance(),
	}
}

// DefaultDevelopmentFundState returns default development fund state
func DefaultDevelopmentFundState() *DevelopmentFundState {
	return &DevelopmentFundState{
		Balance:      nil, // Will be initialized in InitGenesis
		Proposals:    []DevelopmentFundProposal{},
		Transactions: []DevelopmentFundTransaction{},
	}
}

// DefaultCommunityGovernance returns default community governance parameters
func DefaultCommunityGovernance() *CommunityGovernance {
	return &CommunityGovernance{
		QuorumPercentage:      fmt.Sprintf("%d", DefaultQuorumPercentage),
		PassingPercentage:     fmt.Sprintf("%d", DefaultPassingPercentage),
		VotingPeriod:          time.Duration(DefaultVotingPeriodDays * 24 * time.Hour),
		MinDeposit:            NewCoin("namo", 1000000), // 1 NAMO minimum
		MaxProposalSize:       NewCoin("namo", 10000000000), // 10,000 NAMO max per proposal
		RequiredStake:         NewCoin("namo", 100000), // 0.1 NAMO to vote
		AuditThreshold:        NewCoin("namo", 5000000000), // 5,000 NAMO requires audit
		TransparencyRequired:  true,
		CommunityApprovalReq:  true,
		MultiSigRequired:      true,
		MinReputation:         50,
	}
}

// DefaultProposalSystem returns default proposal system
func DefaultProposalSystem() *CommunityProposalSystem {
	now := time.Now()
	return &CommunityProposalSystem{
		Id:             1,
		Name:           "DeshChain Community Governance",
		Description:    "Phased community governance system with gradual power transition",
		LaunchDate:     now,
		ActivationDate: now,
		CurrentPhase:   PhaseFounderDriven,
		Status:         "active",
		PhaseSchedule: PhaseSchedule{
			FounderDrivenEnd:     now.AddDate(3, 0, 0),  // 3 years
			TransitionalEnd:      now.AddDate(4, 0, 0),  // 4 years
			CommunityProposalEnd: now.AddDate(7, 0, 0),  // 7 years
			FullGovernanceStart:  now.AddDate(7, 0, 0),  // 7+ years
		},
		PhaseConfigs: []PhaseConfig{
			{
				Phase:                  PhaseFounderDriven,
				Description:            "Founder-driven development with community input",
				FounderAllocationPower: 100,
				CommunityProposalPower: 0,
				CommunityVotingEnabled: false,
				FounderVetoEnabled:     true,
				AllowedProposalTypes:   []string{"feedback", "suggestion"},
			},
			{
				Phase:                  PhaseTransitional,
				Description:            "Transitional phase with shared governance",
				FounderAllocationPower: 70,
				CommunityProposalPower: 30,
				CommunityVotingEnabled: true,
				FounderVetoEnabled:     true,
				AllowedProposalTypes:   []string{"feedback", "suggestion", "minor_allocation", "community_project"},
			},
			{
				Phase:                  PhaseCommunityProposal,
				Description:            "Community-led proposals with founder oversight",
				FounderAllocationPower: 30,
				CommunityProposalPower: 70,
				CommunityVotingEnabled: true,
				FounderVetoEnabled:     true,
				AllowedProposalTypes:   []string{"all"},
			},
			{
				Phase:                  PhaseFullGovernance,
				Description:            "Full community governance with minimal founder control",
				FounderAllocationPower: 10,
				CommunityProposalPower: 90,
				CommunityVotingEnabled: true,
				FounderVetoEnabled:     false,
				AllowedProposalTypes:   []string{"all"},
			},
		},
		EmergencyProtocol: EmergencyProtocol{
			Enabled:           true,
			MinSigners:        7,
			Threshold:         5,
			MaxAmount:         NewCoin("namo", 100000000000), // 100,000 NAMO emergency limit
			CooldownPeriod:    time.Duration(7 * 24 * time.Hour),
			AllowedCategories: []string{"security", "critical_bug", "exploit_prevention"},
		},
	}
}

// DefaultDashboard returns default dashboard
func DefaultDashboard() *RealTimeDashboard {
	return &RealTimeDashboard{
		LastUpdated:       time.Now(),
		TransparencyScore: 10,
		ComplianceScore:   10,
		Status:            "active",
	}
}

// ValidateGenesis performs basic genesis state validation
func ValidateGenesis(data GenesisState) error {
	// Validate community fund
	if data.CommunityFund != nil {
		if err := validateCommunityFund(data.CommunityFund); err != nil {
			return err
		}
	}

	// Validate development fund
	if data.DevelopmentFund != nil {
		if err := validateDevelopmentFund(data.DevelopmentFund); err != nil {
			return err
		}
	}

	// Validate multi-sig governances
	for _, gov := range data.MultiSigGovernances {
		if err := validateMultiSigGovernance(gov); err != nil {
			return err
		}
	}

	// Validate proposal system
	if data.ProposalSystem != nil {
		if err := validateProposalSystem(data.ProposalSystem); err != nil {
			return err
		}
	}

	// Validate transparency reports
	for _, report := range data.TransparencyReports {
		if err := validateTransparencyReport(report); err != nil {
			return err
		}
	}

	return nil
}

// validateCommunityFund validates community fund state
func validateCommunityFund(fund *CommunityFundState) error {
	// Validate proposals
	proposalIDs := make(map[uint64]bool)
	for _, proposal := range fund.Proposals {
		if proposalIDs[proposal.ProposalId] {
			return fmt.Errorf("duplicate community proposal ID: %d", proposal.ProposalId)
		}
		proposalIDs[proposal.ProposalId] = true

		if err := validateCommunityProposal(proposal); err != nil {
			return err
		}
	}

	// Validate transactions
	for _, tx := range fund.Transactions {
		if err := validateTransaction(tx); err != nil {
			return err
		}
	}

	// Validate governance
	if fund.Governance != nil {
		if err := validateGovernance(fund.Governance); err != nil {
			return err
		}
	}

	return nil
}

// validateDevelopmentFund validates development fund state
func validateDevelopmentFund(fund *DevelopmentFundState) error {
	// Validate proposals
	proposalIDs := make(map[uint64]bool)
	for _, proposal := range fund.Proposals {
		if proposalIDs[proposal.ProposalId] {
			return fmt.Errorf("duplicate development proposal ID: %d", proposal.ProposalId)
		}
		proposalIDs[proposal.ProposalId] = true

		if err := validateDevelopmentProposal(proposal); err != nil {
			return err
		}
	}

	// Validate transactions
	for _, tx := range fund.Transactions {
		if err := validateDevelopmentTransaction(tx); err != nil {
			return err
		}
	}

	return nil
}

// validateMultiSigGovernance validates multi-sig governance
func validateMultiSigGovernance(gov MultiSigGovernance) error {
	if gov.Id == 0 {
		return fmt.Errorf("invalid governance ID: must be greater than 0")
	}

	if gov.Threshold == 0 || gov.Threshold > uint8(len(gov.Signers)) {
		return fmt.Errorf("invalid threshold: must be between 1 and number of signers")
	}

	if gov.Name == "" {
		return fmt.Errorf("governance name cannot be empty")
	}

	// Validate signers
	signerAddresses := make(map[string]bool)
	for _, signer := range gov.Signers {
		if signerAddresses[signer.Address] {
			return fmt.Errorf("duplicate signer address: %s", signer.Address)
		}
		signerAddresses[signer.Address] = true

		if err := validateSigner(signer); err != nil {
			return err
		}
	}

	return nil
}

// validateProposalSystem validates proposal system
func validateProposalSystem(system *CommunityProposalSystem) error {
	if system.Id == 0 {
		return fmt.Errorf("invalid system ID: must be greater than 0")
	}

	if system.Name == "" {
		return fmt.Errorf("system name cannot be empty")
	}

	// Validate phase schedule
	if system.PhaseSchedule.FounderDrivenEnd.Before(system.LaunchDate) {
		return fmt.Errorf("founder driven end date cannot be before launch date")
	}

	if system.PhaseSchedule.TransitionalEnd.Before(system.PhaseSchedule.FounderDrivenEnd) {
		return fmt.Errorf("transitional end date must be after founder driven end date")
	}

	if system.PhaseSchedule.CommunityProposalEnd.Before(system.PhaseSchedule.TransitionalEnd) {
		return fmt.Errorf("community proposal end date must be after transitional end date")
	}

	if system.PhaseSchedule.FullGovernanceStart.Before(system.PhaseSchedule.CommunityProposalEnd) {
		return fmt.Errorf("full governance start date must be after community proposal end date")
	}

	// Validate phase configs
	phases := make(map[string]bool)
	for _, config := range system.PhaseConfigs {
		if phases[config.Phase] {
			return fmt.Errorf("duplicate phase config: %s", config.Phase)
		}
		phases[config.Phase] = true

		if config.FounderAllocationPower+config.CommunityProposalPower > 100 {
			return fmt.Errorf("total allocation power cannot exceed 100 for phase: %s", config.Phase)
		}
	}

	return nil
}

// validateTransparencyReport validates transparency report
func validateTransparencyReport(report TransparencyReport) error {
	if report.ReportId == 0 {
		return fmt.Errorf("invalid report ID: must be greater than 0")
	}

	if report.EndDate.Before(report.StartDate) {
		return fmt.Errorf("report end date must be after start date")
	}

	if !report.TotalFunds.IsValid() {
		return fmt.Errorf("invalid total funds amount")
	}

	if !report.AllocatedFunds.IsValid() {
		return fmt.Errorf("invalid allocated funds amount")
	}

	if !report.SpentFunds.IsValid() {
		return fmt.Errorf("invalid spent funds amount")
	}

	if !report.RemainingFunds.IsValid() {
		return fmt.Errorf("invalid remaining funds amount")
	}

	// Verify funds consistency
	calculatedRemaining := report.TotalFunds.Amount.Sub(report.SpentFunds.Amount)
	if !calculatedRemaining.Equal(report.RemainingFunds.Amount) {
		return fmt.Errorf("remaining funds mismatch: expected %s, got %s", 
			calculatedRemaining.String(), report.RemainingFunds.Amount.String())
	}

	return nil
}

// Additional validation helper functions
func validateCommunityProposal(proposal CommunityFundProposal) error {
	if proposal.ProposalId == 0 {
		return fmt.Errorf("invalid proposal ID: must be greater than 0")
	}

	if proposal.Title == "" {
		return fmt.Errorf("proposal title cannot be empty")
	}

	if !proposal.RequestedAmount.IsValid() || proposal.RequestedAmount.IsZero() {
		return fmt.Errorf("invalid requested amount for proposal %d", proposal.ProposalId)
	}

	return nil
}

func validateDevelopmentProposal(proposal DevelopmentFundProposal) error {
	if proposal.ProposalId == 0 {
		return fmt.Errorf("invalid proposal ID: must be greater than 0")
	}

	if proposal.Title == "" {
		return fmt.Errorf("proposal title cannot be empty")
	}

	if !proposal.RequestedAmount.IsValid() || proposal.RequestedAmount.IsZero() {
		return fmt.Errorf("invalid requested amount for proposal %d", proposal.ProposalId)
	}

	return nil
}

func validateTransaction(tx CommunityFundTransaction) error {
	if tx.TxId == "" {
		return fmt.Errorf("transaction ID cannot be empty")
	}

	if !tx.Amount.IsValid() || tx.Amount.IsZero() {
		return fmt.Errorf("invalid transaction amount for tx %s", tx.TxId)
	}

	return nil
}

func validateDevelopmentTransaction(tx DevelopmentFundTransaction) error {
	if tx.TxId == "" {
		return fmt.Errorf("transaction ID cannot be empty")
	}

	if !tx.Amount.IsValid() || tx.Amount.IsZero() {
		return fmt.Errorf("invalid transaction amount for tx %s", tx.TxId)
	}

	return nil
}

func validateGovernance(gov *CommunityGovernance) error {
	if gov.VotingPeriod <= 0 {
		return fmt.Errorf("voting period must be positive")
	}

	if !gov.MinDeposit.IsValid() {
		return fmt.Errorf("invalid minimum deposit")
	}

	if !gov.MaxProposalSize.IsValid() {
		return fmt.Errorf("invalid maximum proposal size")
	}

	return nil
}

func validateSigner(signer Signer) error {
	if signer.Address == "" {
		return fmt.Errorf("signer address cannot be empty")
	}

	if signer.Weight == 0 {
		return fmt.Errorf("signer weight must be greater than 0")
	}

	return nil
}

// Helper function
func NewCoin(denom string, amount int64) Coin {
	return Coin{
		Denom:  denom,
		Amount: NewInt(amount),
	}
}

// Mock Int type for amount (replace with actual SDK Int)
type Int struct {
	value int64
}

func NewInt(value int64) Int {
	return Int{value: value}
}

func (i Int) Sub(other Int) Int {
	return Int{value: i.value - other.value}
}

func (i Int) Equal(other Int) bool {
	return i.value == other.value
}

func (i Int) String() string {
	return fmt.Sprintf("%d", i.value)
}