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
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                       DefaultParams(),
		NgoWallets:                   DefaultNGOWallets(),
		NgoWalletCount:               uint64(len(DefaultNGOWallets())),
		DonationRecords:              []DonationRecord{},
		DonationRecordCount:          0,
		DistributionRecords:          []DistributionRecord{},
		DistributionRecordCount:      0,
		AuditReports:                 []AuditReport{},
		AuditReportCount:             0,
		BeneficiaryTestimonials:      []BeneficiaryTestimonial{},
		BeneficiaryTestimonialCount:  0,
		Campaigns:                    []Campaign{},
		CampaignCount:                0,
		RecurringDonations:           []RecurringDonation{},
		RecurringDonationCount:       0,
		EmergencyPause:               nil,
		Statistics:                   nil,
		FundFlows:                    []FundFlow{},
		TransparencyScores:           []TransparencyScore{},
		VerificationQueue:            []VerificationQueueItem{},
	}
}

// ValidateGenesis validates the donation genesis state
func ValidateGenesis(data GenesisState) error {
	// Validate params
	if err := data.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Validate NGO wallets
	ngoMap := make(map[uint64]bool)
	addressMap := make(map[string]bool)
	for _, ngo := range data.NgoWallets {
		// Check for duplicate IDs
		if ngoMap[ngo.Id] {
			return fmt.Errorf("duplicate NGO wallet ID: %d", ngo.Id)
		}
		ngoMap[ngo.Id] = true

		// Check for duplicate addresses
		if addressMap[ngo.Address] {
			return fmt.Errorf("duplicate NGO wallet address: %s", ngo.Address)
		}
		addressMap[ngo.Address] = true

		// Validate NGO wallet
		if err := ValidateNGOWallet(ngo); err != nil {
			return fmt.Errorf("invalid NGO wallet %d: %w", ngo.Id, err)
		}
	}

	// Validate donation records
	donationMap := make(map[uint64]bool)
	for _, donation := range data.DonationRecords {
		// Check for duplicate IDs
		if donationMap[donation.Id] {
			return fmt.Errorf("duplicate donation record ID: %d", donation.Id)
		}
		donationMap[donation.Id] = true

		// Check NGO exists
		if !ngoMap[donation.NgoWalletId] {
			return fmt.Errorf("donation record %d references non-existent NGO wallet %d", donation.Id, donation.NgoWalletId)
		}

		// Validate donation record
		if err := ValidateDonationRecord(donation); err != nil {
			return fmt.Errorf("invalid donation record %d: %w", donation.Id, err)
		}
	}

	// Validate distribution records
	distributionMap := make(map[uint64]bool)
	for _, distribution := range data.DistributionRecords {
		// Check for duplicate IDs
		if distributionMap[distribution.Id] {
			return fmt.Errorf("duplicate distribution record ID: %d", distribution.Id)
		}
		distributionMap[distribution.Id] = true

		// Check NGO exists
		if !ngoMap[distribution.NgoWalletId] {
			return fmt.Errorf("distribution record %d references non-existent NGO wallet %d", distribution.Id, distribution.NgoWalletId)
		}

		// Validate distribution record
		if err := ValidateDistributionRecord(distribution); err != nil {
			return fmt.Errorf("invalid distribution record %d: %w", distribution.Id, err)
		}
	}

	// Validate audit reports
	auditMap := make(map[uint64]bool)
	for _, audit := range data.AuditReports {
		// Check for duplicate IDs
		if auditMap[audit.Id] {
			return fmt.Errorf("duplicate audit report ID: %d", audit.Id)
		}
		auditMap[audit.Id] = true

		// Check NGO exists
		if !ngoMap[audit.NgoWalletId] {
			return fmt.Errorf("audit report %d references non-existent NGO wallet %d", audit.Id, audit.NgoWalletId)
		}

		// Validate audit report
		if err := ValidateAuditReport(audit); err != nil {
			return fmt.Errorf("invalid audit report %d: %w", audit.Id, err)
		}
	}

	// Validate beneficiary testimonials
	testimonialMap := make(map[uint64]bool)
	for _, testimonial := range data.BeneficiaryTestimonials {
		// Check for duplicate IDs
		if testimonialMap[testimonial.Id] {
			return fmt.Errorf("duplicate beneficiary testimonial ID: %d", testimonial.Id)
		}
		testimonialMap[testimonial.Id] = true

		// Check NGO exists
		if !ngoMap[testimonial.NgoWalletId] {
			return fmt.Errorf("beneficiary testimonial %d references non-existent NGO wallet %d", testimonial.Id, testimonial.NgoWalletId)
		}

		// Validate beneficiary testimonial
		if err := ValidateBeneficiaryTestimonial(testimonial); err != nil {
			return fmt.Errorf("invalid beneficiary testimonial %d: %w", testimonial.Id, err)
		}
	}

	// Validate campaigns
	campaignMap := make(map[uint64]bool)
	for _, campaign := range data.Campaigns {
		// Check for duplicate IDs
		if campaignMap[campaign.Id] {
			return fmt.Errorf("duplicate campaign ID: %d", campaign.Id)
		}
		campaignMap[campaign.Id] = true

		// Check NGO exists
		if !ngoMap[campaign.NgoWalletId] {
			return fmt.Errorf("campaign %d references non-existent NGO wallet %d", campaign.Id, campaign.NgoWalletId)
		}

		// Validate campaign
		if err := ValidateCampaign(campaign); err != nil {
			return fmt.Errorf("invalid campaign %d: %w", campaign.Id, err)
		}
	}

	// Validate recurring donations
	recurringMap := make(map[uint64]bool)
	for _, recurring := range data.RecurringDonations {
		// Check for duplicate IDs
		if recurringMap[recurring.Id] {
			return fmt.Errorf("duplicate recurring donation ID: %d", recurring.Id)
		}
		recurringMap[recurring.Id] = true

		// Check NGO exists
		if !ngoMap[recurring.NgoWalletId] {
			return fmt.Errorf("recurring donation %d references non-existent NGO wallet %d", recurring.Id, recurring.NgoWalletId)
		}

		// Validate recurring donation
		if err := ValidateRecurringDonation(recurring); err != nil {
			return fmt.Errorf("invalid recurring donation %d: %w", recurring.Id, err)
		}
	}

	// Validate emergency pause if present
	if data.EmergencyPause != nil {
		if err := ValidateEmergencyPause(*data.EmergencyPause); err != nil {
			return fmt.Errorf("invalid emergency pause: %w", err)
		}
	}

	// Validate statistics if present
	if data.Statistics != nil {
		if err := ValidateStatistics(*data.Statistics); err != nil {
			return fmt.Errorf("invalid statistics: %w", err)
		}
	}

	// Validate fund flows
	for i, flow := range data.FundFlows {
		if err := ValidateFundFlow(flow); err != nil {
			return fmt.Errorf("invalid fund flow %d: %w", i, err)
		}
	}

	// Validate transparency scores
	for _, score := range data.TransparencyScores {
		// Check NGO exists
		if !ngoMap[score.NgoWalletId] {
			return fmt.Errorf("transparency score references non-existent NGO wallet %d", score.NgoWalletId)
		}

		if err := ValidateTransparencyScore(score); err != nil {
			return fmt.Errorf("invalid transparency score for NGO %d: %w", score.NgoWalletId, err)
		}
	}

	// Validate verification queue
	for i, item := range data.VerificationQueue {
		// Check NGO exists
		if !ngoMap[item.NgoWalletId] {
			return fmt.Errorf("verification queue item %d references non-existent NGO wallet %d", i, item.NgoWalletId)
		}

		if err := ValidateVerificationQueueItem(item); err != nil {
			return fmt.Errorf("invalid verification queue item %d: %w", i, err)
		}
	}

	return nil
}

// GenesisState defines the donation module's genesis state.
type GenesisState struct {
	Params                       Params                    `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	NgoWallets                   []NGOWallet               `protobuf:"bytes,2,rep,name=ngo_wallets,json=ngoWallets,proto3" json:"ngo_wallets"`
	NgoWalletCount               uint64                    `protobuf:"varint,3,opt,name=ngo_wallet_count,json=ngoWalletCount,proto3" json:"ngo_wallet_count,omitempty"`
	DonationRecords              []DonationRecord          `protobuf:"bytes,4,rep,name=donation_records,json=donationRecords,proto3" json:"donation_records"`
	DonationRecordCount          uint64                    `protobuf:"varint,5,opt,name=donation_record_count,json=donationRecordCount,proto3" json:"donation_record_count,omitempty"`
	DistributionRecords          []DistributionRecord      `protobuf:"bytes,6,rep,name=distribution_records,json=distributionRecords,proto3" json:"distribution_records"`
	DistributionRecordCount      uint64                    `protobuf:"varint,7,opt,name=distribution_record_count,json=distributionRecordCount,proto3" json:"distribution_record_count,omitempty"`
	AuditReports                 []AuditReport             `protobuf:"bytes,8,rep,name=audit_reports,json=auditReports,proto3" json:"audit_reports"`
	AuditReportCount             uint64                    `protobuf:"varint,9,opt,name=audit_report_count,json=auditReportCount,proto3" json:"audit_report_count,omitempty"`
	BeneficiaryTestimonials      []BeneficiaryTestimonial  `protobuf:"bytes,10,rep,name=beneficiary_testimonials,json=beneficiaryTestimonials,proto3" json:"beneficiary_testimonials"`
	BeneficiaryTestimonialCount  uint64                    `protobuf:"varint,11,opt,name=beneficiary_testimonial_count,json=beneficiaryTestimonialCount,proto3" json:"beneficiary_testimonial_count,omitempty"`
	Campaigns                    []Campaign                `protobuf:"bytes,12,rep,name=campaigns,proto3" json:"campaigns"`
	CampaignCount                uint64                    `protobuf:"varint,13,opt,name=campaign_count,json=campaignCount,proto3" json:"campaign_count,omitempty"`
	RecurringDonations           []RecurringDonation       `protobuf:"bytes,14,rep,name=recurring_donations,json=recurringDonations,proto3" json:"recurring_donations"`
	RecurringDonationCount       uint64                    `protobuf:"varint,15,opt,name=recurring_donation_count,json=recurringDonationCount,proto3" json:"recurring_donation_count,omitempty"`
	EmergencyPause               *EmergencyPause           `protobuf:"bytes,16,opt,name=emergency_pause,json=emergencyPause,proto3" json:"emergency_pause,omitempty"`
	Statistics                   *Statistics               `protobuf:"bytes,17,opt,name=statistics,proto3" json:"statistics,omitempty"`
	FundFlows                    []FundFlow                `protobuf:"bytes,18,rep,name=fund_flows,json=fundFlows,proto3" json:"fund_flows"`
	TransparencyScores           []TransparencyScore       `protobuf:"bytes,19,rep,name=transparency_scores,json=transparencyScores,proto3" json:"transparency_scores"`
	VerificationQueue            []VerificationQueueItem   `protobuf:"bytes,20,rep,name=verification_queue,json=verificationQueue,proto3" json:"verification_queue"`
}