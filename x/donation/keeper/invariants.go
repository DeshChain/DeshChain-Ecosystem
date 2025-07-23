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

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"deshchain/x/donation/types"
)

// RegisterInvariants registers all donation module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "ngo-balances", NGOBalancesInvariant(k))
	ir.RegisterRoute(types.ModuleName, "donation-consistency", DonationConsistencyInvariant(k))
	ir.RegisterRoute(types.ModuleName, "transparency-scores", TransparencyScoresInvariant(k))
	ir.RegisterRoute(types.ModuleName, "campaign-funds", CampaignFundsInvariant(k))
}

// NGOBalancesInvariant checks that NGO balances are consistent
func NGOBalancesInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			broken bool
			msg    string
		)

		ngos := k.GetAllNGOWallets(ctx)
		for _, ngo := range ngos {
			// Check that current balance = total received - total distributed
			expectedBalance := ngo.TotalReceived.Sub(ngo.TotalDistributed...)
			if !expectedBalance.IsEqual(ngo.CurrentBalance) {
				broken = true
				msg += fmt.Sprintf("NGO %d balance mismatch: expected %s, got %s\n", 
					ngo.Id, expectedBalance.String(), ngo.CurrentBalance.String())
			}

			// Check that balances are non-negative
			if ngo.CurrentBalance.IsAnyNegative() {
				broken = true
				msg += fmt.Sprintf("NGO %d has negative balance: %s\n", ngo.Id, ngo.CurrentBalance.String())
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "ngo-balances", 
			fmt.Sprintf("NGO balance consistency check\n%s", msg)), broken
	}
}

// DonationConsistencyInvariant checks that donation records are consistent
func DonationConsistencyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			broken bool
			msg    string
		)

		// Calculate total donations from records
		totalFromRecords := sdk.Coins{}
		donations := k.GetAllDonationRecords(ctx)
		for _, donation := range donations {
			totalFromRecords = totalFromRecords.Add(donation.Amount...)
		}

		// Get total from statistics
		stats, found := k.GetStatistics(ctx)
		if found {
			if !totalFromRecords.IsEqual(stats.TotalDonations) {
				broken = true
				msg += fmt.Sprintf("Total donations mismatch: records=%s, stats=%s\n",
					totalFromRecords.String(), stats.TotalDonations.String())
			}
		}

		// Verify donor count
		donorMap := make(map[string]bool)
		for _, donation := range donations {
			if !donation.IsAnonymous {
				donorMap[donation.Donor] = true
			}
		}
		actualDonorCount := uint64(len(donorMap))
		if found && actualDonorCount != stats.TotalDonors {
			broken = true
			msg += fmt.Sprintf("Donor count mismatch: actual=%d, stats=%d\n",
				actualDonorCount, stats.TotalDonors)
		}

		return sdk.FormatInvariant(types.ModuleName, "donation-consistency",
			fmt.Sprintf("Donation consistency check\n%s", msg)), broken
	}
}

// TransparencyScoresInvariant checks that transparency scores are valid
func TransparencyScoresInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			broken bool
			msg    string
		)

		ngos := k.GetAllNGOWallets(ctx)
		for _, ngo := range ngos {
			// Check that transparency score is within valid range
			if ngo.TransparencyScore < 1 || ngo.TransparencyScore > 10 {
				broken = true
				msg += fmt.Sprintf("NGO %d has invalid transparency score: %d\n",
					ngo.Id, ngo.TransparencyScore)
			}

			// Check that verified NGOs have minimum score
			if ngo.IsVerified && ngo.TransparencyScore < 5 {
				broken = true
				msg += fmt.Sprintf("Verified NGO %d has low transparency score: %d\n",
					ngo.Id, ngo.TransparencyScore)
			}
		}

		// Check transparency score records
		scores := k.GetAllTransparencyScores(ctx)
		for _, score := range scores {
			if score.Score < 1 || score.Score > 10 {
				broken = true
				msg += fmt.Sprintf("Invalid transparency score record for NGO %d: %d\n",
					score.NgoWalletId, score.Score)
			}

			// Check component scores
			components := []struct {
				name  string
				value int32
			}{
				{"AuditCompleteness", score.AuditCompleteness},
				{"ReportingFrequency", score.ReportingFrequency},
				{"DocumentationQuality", score.DocumentationQuality},
				{"FundUtilization", score.FundUtilization},
				{"BeneficiaryFeedback", score.BeneficiaryFeedback},
				{"ResponseTime", score.ResponseTime},
				{"PublicAccessibility", score.PublicAccessibility},
				{"ComplianceAdherence", score.ComplianceAdherence},
			}

			for _, comp := range components {
				if comp.value < 0 || comp.value > 10 {
					broken = true
					msg += fmt.Sprintf("Invalid %s score for NGO %d: %d\n",
						comp.name, score.NgoWalletId, comp.value)
				}
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "transparency-scores",
			fmt.Sprintf("Transparency scores check\n%s", msg)), broken
	}
}

// CampaignFundsInvariant checks that campaign funds are properly tracked
func CampaignFundsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			broken bool
			msg    string
		)

		campaigns := k.GetAllCampaigns(ctx)
		for _, campaign := range campaigns {
			// Check that raised amount doesn't exceed target for non-matching campaigns
			if !campaign.MatchingEnabled {
				for i, raised := range campaign.RaisedAmount {
					if i < len(campaign.TargetAmount) {
						target := campaign.TargetAmount[i]
						if raised.Amount.GT(target.Amount.Mul(sdk.NewInt(2))) {
							// Allow some buffer for multiple donations processed in same block
							broken = true
							msg += fmt.Sprintf("Campaign %d raised excessive amount: %s > %s\n",
								campaign.Id, raised.String(), target.String())
						}
					}
				}
			}

			// Check matching funds consistency
			if campaign.MatchingEnabled {
				if !campaign.MatchingUsed.IsAllLTE(campaign.MatchingBudget) {
					broken = true
					msg += fmt.Sprintf("Campaign %d matching funds exceeded budget: used=%s, budget=%s\n",
						campaign.Id, campaign.MatchingUsed.String(), campaign.MatchingBudget.String())
				}
			}

			// Check campaign dates
			if campaign.EndDate <= campaign.StartDate {
				broken = true
				msg += fmt.Sprintf("Campaign %d has invalid date range: start=%d, end=%d\n",
					campaign.Id, campaign.StartDate, campaign.EndDate)
			}

			// Check completion status
			if campaign.CompletedAt > 0 && campaign.IsActive {
				broken = true
				msg += fmt.Sprintf("Campaign %d is marked completed but still active\n", campaign.Id)
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "campaign-funds",
			fmt.Sprintf("Campaign funds consistency check\n%s", msg)), broken
	}
}

// GetAllInvariants returns all the donation module invariants
func GetAllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := NGOBalancesInvariant(k)(ctx)
		if stop {
			return res, stop
		}

		res, stop = DonationConsistencyInvariant(k)(ctx)
		if stop {
			return res, stop
		}

		res, stop = TransparencyScoresInvariant(k)(ctx)
		if stop {
			return res, stop
		}

		return CampaignFundsInvariant(k)(ctx)
	}
}