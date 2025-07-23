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

package donation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"deshchain/x/donation/keeper"
	"deshchain/x/donation/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module params
	k.SetParams(ctx, genState.Params)

	// Initialize NGO wallets
	for _, ngo := range genState.NgoWallets {
		k.SetNGOWallet(ctx, ngo)
		// Index by address
		k.SetNGOByAddress(ctx, ngo.Address, ngo.Id)
	}

	// Set NGO wallet counter
	k.SetNGOWalletCount(ctx, genState.NgoWalletCount)

	// Initialize donation records
	for _, donation := range genState.DonationRecords {
		k.SetDonationRecord(ctx, donation)
		// Index by donor
		k.AddDonationByDonor(ctx, donation.Donor, donation.Id)
		// Index by NGO
		k.AddDonationByNGO(ctx, donation.NgoWalletId, donation.Id)
	}

	// Set donation record counter
	k.SetDonationRecordCount(ctx, genState.DonationRecordCount)

	// Initialize distribution records
	for _, distribution := range genState.DistributionRecords {
		k.SetDistributionRecord(ctx, distribution)
		// Index by NGO
		k.AddDistributionByNGO(ctx, distribution.NgoWalletId, distribution.Id)
	}

	// Set distribution record counter
	k.SetDistributionRecordCount(ctx, genState.DistributionRecordCount)

	// Initialize audit reports
	for _, audit := range genState.AuditReports {
		k.SetAuditReport(ctx, audit)
		// Index by NGO
		k.AddAuditByNGO(ctx, audit.NgoWalletId, audit.Id)
	}

	// Set audit report counter
	k.SetAuditReportCount(ctx, genState.AuditReportCount)

	// Initialize beneficiary testimonials
	for _, testimonial := range genState.BeneficiaryTestimonials {
		k.SetBeneficiaryTestimonial(ctx, testimonial)
		// Index by NGO
		k.AddTestimonialByNGO(ctx, testimonial.NgoWalletId, testimonial.Id)
	}

	// Set beneficiary testimonial counter
	k.SetBeneficiaryTestimonialCount(ctx, genState.BeneficiaryTestimonialCount)

	// Initialize campaigns
	for _, campaign := range genState.Campaigns {
		k.SetCampaign(ctx, campaign)
	}

	// Set campaign counter
	k.SetCampaignCount(ctx, genState.CampaignCount)

	// Initialize recurring donations
	for _, recurring := range genState.RecurringDonations {
		k.SetRecurringDonation(ctx, recurring)
	}

	// Set recurring donation counter
	k.SetRecurringDonationCount(ctx, genState.RecurringDonationCount)

	// Initialize emergency pause status
	if genState.EmergencyPause != nil {
		k.SetEmergencyPause(ctx, *genState.EmergencyPause)
	}

	// Initialize statistics
	if genState.Statistics != nil {
		k.SetStatistics(ctx, *genState.Statistics)
	}

	// Initialize fund flows
	for _, flow := range genState.FundFlows {
		k.SetFundFlow(ctx, flow)
	}

	// Initialize transparency scores
	for _, score := range genState.TransparencyScores {
		k.SetTransparencyScore(ctx, score)
	}

	// Initialize verification queue
	for _, item := range genState.VerificationQueue {
		k.AddToVerificationQueue(ctx, item)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export NGO wallets
	genesis.NgoWallets = k.GetAllNGOWallets(ctx)
	genesis.NgoWalletCount = k.GetNGOWalletCount(ctx)

	// Export donation records
	genesis.DonationRecords = k.GetAllDonationRecords(ctx)
	genesis.DonationRecordCount = k.GetDonationRecordCount(ctx)

	// Export distribution records
	genesis.DistributionRecords = k.GetAllDistributionRecords(ctx)
	genesis.DistributionRecordCount = k.GetDistributionRecordCount(ctx)

	// Export audit reports
	genesis.AuditReports = k.GetAllAuditReports(ctx)
	genesis.AuditReportCount = k.GetAuditReportCount(ctx)

	// Export beneficiary testimonials
	genesis.BeneficiaryTestimonials = k.GetAllBeneficiaryTestimonials(ctx)
	genesis.BeneficiaryTestimonialCount = k.GetBeneficiaryTestimonialCount(ctx)

	// Export campaigns
	genesis.Campaigns = k.GetAllCampaigns(ctx)
	genesis.CampaignCount = k.GetCampaignCount(ctx)

	// Export recurring donations
	genesis.RecurringDonations = k.GetAllRecurringDonations(ctx)
	genesis.RecurringDonationCount = k.GetRecurringDonationCount(ctx)

	// Export emergency pause
	pause, found := k.GetEmergencyPause(ctx)
	if found {
		genesis.EmergencyPause = &pause
	}

	// Export statistics
	stats, found := k.GetStatistics(ctx)
	if found {
		genesis.Statistics = &stats
	}

	// Export fund flows
	genesis.FundFlows = k.GetAllFundFlows(ctx)

	// Export transparency scores
	genesis.TransparencyScores = k.GetAllTransparencyScores(ctx)

	// Export verification queue
	genesis.VerificationQueue = k.GetAllVerificationQueueItems(ctx)

	return genesis
}