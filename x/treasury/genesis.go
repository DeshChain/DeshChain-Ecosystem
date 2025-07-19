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

package treasury

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/treasury/keeper"
	"github.com/deshchain/deshchain/x/treasury/types"
)

// InitGenesis initializes the treasury module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Initialize community fund
	if genState.CommunityFund != nil {
		// Set initial balance
		k.SetCommunityFundBalance(ctx, types.CommunityFundBalance{
			TotalBalance:    sdk.NewCoin("namo", sdk.NewInt(types.CommunityFundAllocation)),
			AllocatedAmount: sdk.NewCoin("namo", sdk.ZeroInt()),
			AvailableAmount: sdk.NewCoin("namo", sdk.NewInt(types.CommunityFundAllocation)),
			ReservedAmount:  sdk.NewCoin("namo", sdk.ZeroInt()),
			PendingAmount:   sdk.NewCoin("namo", sdk.ZeroInt()),
			LastUpdateHeight: ctx.BlockHeight(),
			LastUpdateTime:  ctx.BlockTime(),
		})

		// Set proposals
		for _, proposal := range genState.CommunityFund.Proposals {
			k.SetCommunityFundProposal(ctx, proposal)
		}

		// Set transactions
		for _, tx := range genState.CommunityFund.Transactions {
			k.SetCommunityFundTransaction(ctx, tx)
		}

		// Set governance parameters
		if genState.CommunityFund.Governance != nil {
			k.SetCommunityGovernance(ctx, *genState.CommunityFund.Governance)
		}
	}

	// Initialize development fund
	if genState.DevelopmentFund != nil {
		// Set initial balance
		k.SetDevelopmentFundBalance(ctx, types.DevelopmentFundBalance{
			TotalBalance:     sdk.NewCoin("namo", sdk.NewInt(types.DevelopmentFundAllocation)),
			AllocatedAmount:  sdk.NewCoin("namo", sdk.ZeroInt()),
			AvailableAmount:  sdk.NewCoin("namo", sdk.NewInt(types.DevelopmentFundAllocation)),
			PendingAmount:    sdk.NewCoin("namo", sdk.ZeroInt()),
			EmergencyReserve: sdk.NewCoin("namo", sdk.NewInt(types.DevelopmentFundAllocation/10)), // 10% emergency reserve
			LastUpdateHeight: ctx.BlockHeight(),
			LastUpdateTime:   ctx.BlockTime(),
		})

		// Set proposals
		for _, proposal := range genState.DevelopmentFund.Proposals {
			k.SetDevelopmentFundProposal(ctx, proposal)
		}

		// Set transactions
		for _, tx := range genState.DevelopmentFund.Transactions {
			k.SetDevelopmentFundTransaction(ctx, tx)
		}
	}

	// Initialize multi-sig governance
	if len(genState.MultiSigGovernances) > 0 {
		for _, gov := range genState.MultiSigGovernances {
			k.SetMultiSigGovernance(ctx, gov)
		}
	} else {
		// Initialize default multi-sig governance for community fund
		communityGov := types.MultiSigGovernance{
			Id:          1,
			Name:        "Community Fund Governance",
			Description: "Multi-signature governance for community fund management",
			Type:        types.GovernanceTypeCommunity,
			Threshold:   5, // 5 out of 9
			Signers:     []types.Signer{}, // To be added through governance
			CreatedAt:   ctx.BlockTime(),
			LastUpdated: ctx.BlockTime(),
			Status:      "active",
			Rules: types.GovernanceRules{
				MinSigners:       9,
				MaxSigners:       15,
				ProposalDuration: 7 * 24 * 60 * 60, // 7 days in seconds
				ExecutionDelay:   24 * 60 * 60,      // 24 hours in seconds
				EmergencyThreshold: 7,               // 7 out of 9 for emergency
			},
		}
		k.SetMultiSigGovernance(ctx, communityGov)

		// Initialize default multi-sig governance for development fund
		developmentGov := types.MultiSigGovernance{
			Id:          2,
			Name:        "Development Fund Governance",
			Description: "Multi-signature governance for development fund management",
			Type:        types.GovernanceTypeDevelopment,
			Threshold:   6, // 6 out of 11
			Signers:     []types.Signer{}, // To be added through governance
			CreatedAt:   ctx.BlockTime(),
			LastUpdated: ctx.BlockTime(),
			Status:      "active",
			Rules: types.GovernanceRules{
				MinSigners:       11,
				MaxSigners:       17,
				ProposalDuration: 14 * 24 * 60 * 60, // 14 days in seconds
				ExecutionDelay:   48 * 60 * 60,      // 48 hours in seconds
				EmergencyThreshold: 8,               // 8 out of 11 for emergency
			},
		}
		k.SetMultiSigGovernance(ctx, developmentGov)
	}

	// Initialize community proposal system
	if genState.ProposalSystem != nil {
		k.SetCommunityProposalSystem(ctx, *genState.ProposalSystem)
		
		// Set initial governance phase
		k.SetGovernancePhase(ctx, genState.ProposalSystem.CurrentPhase)
	} else {
		// Initialize default proposal system
		system := types.CommunityProposalSystem{
			Id:             1,
			Name:           "DeshChain Community Governance",
			Description:    "Phased community governance system with gradual power transition",
			LaunchDate:     ctx.BlockTime(),
			ActivationDate: ctx.BlockTime(),
			CurrentPhase:   types.PhaseFounderDriven,
			Status:         "active",
			PhaseSchedule: types.PhaseSchedule{
				FounderDrivenEnd:      ctx.BlockTime().AddDate(3, 0, 0),  // 3 years
				TransitionalEnd:       ctx.BlockTime().AddDate(4, 0, 0),  // 4 years
				CommunityProposalEnd:  ctx.BlockTime().AddDate(7, 0, 0),  // 7 years
				FullGovernanceStart:   ctx.BlockTime().AddDate(7, 0, 0),  // 7+ years
			},
		}
		k.SetCommunityProposalSystem(ctx, system)
		k.SetGovernancePhase(ctx, types.PhaseFounderDriven)
	}

	// Initialize dashboard
	dashboard := types.RealTimeDashboard{
		LastUpdated:       ctx.BlockTime(),
		CommunityFund:     types.CommunityFundMetrics{},
		DevelopmentFund:   types.DevelopmentFundMetrics{},
		OverallMetrics:    types.OverallFundMetrics{},
		TransparencyScore: 10, // Maximum transparency
		ComplianceScore:   10, // Maximum compliance
		Status:            "active",
	}
	k.SetRealTimeDashboard(ctx, dashboard)

	// Initialize counters
	k.SetNextCommunityProposalID(ctx, 1)
	k.SetNextDevelopmentProposalID(ctx, 1)
	k.SetNextTransparencyReportID(ctx, 1)

	// Register module accounts
	k.RegisterModuleAccounts(ctx)
}

// ExportGenesis returns the treasury module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Export community fund state
	communityFund := &types.CommunityFundState{}
	
	// Get balance
	balance, _ := k.GetCommunityFundBalance(ctx)
	communityFund.Balance = &balance

	// Get all proposals
	proposals := []types.CommunityFundProposal{}
	iter, _ := k.CommunityFundProposals.Iterate(ctx, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		proposal, _ := iter.Value()
		proposals = append(proposals, proposal)
	}
	communityFund.Proposals = proposals

	// Get all transactions
	transactions := []types.CommunityFundTransaction{}
	txIter, _ := k.CommunityFundTransactions.Iterate(ctx, nil)
	defer txIter.Close()
	for ; txIter.Valid(); txIter.Next() {
		tx, _ := txIter.Value()
		transactions = append(transactions, tx)
	}
	communityFund.Transactions = transactions

	// Get governance
	governance, _ := k.GetCommunityGovernance(ctx)
	communityFund.Governance = &governance

	genesis.CommunityFund = communityFund

	// Export development fund state
	developmentFund := &types.DevelopmentFundState{}
	
	// Get balance
	devBalance, _ := k.GetDevelopmentFundBalance(ctx)
	developmentFund.Balance = &devBalance

	// Get all proposals
	devProposals := []types.DevelopmentFundProposal{}
	devIter, _ := k.DevelopmentFundProposals.Iterate(ctx, nil)
	defer devIter.Close()
	for ; devIter.Valid(); devIter.Next() {
		proposal, _ := devIter.Value()
		devProposals = append(devProposals, proposal)
	}
	developmentFund.Proposals = devProposals

	// Get all transactions
	devTransactions := []types.DevelopmentFundTransaction{}
	devTxIter, _ := k.DevelopmentFundTransactions.Iterate(ctx, nil)
	defer devTxIter.Close()
	for ; devTxIter.Valid(); devTxIter.Next() {
		tx, _ := devTxIter.Value()
		devTransactions = append(devTransactions, tx)
	}
	developmentFund.Transactions = devTransactions

	genesis.DevelopmentFund = developmentFund

	// Export multi-sig governances
	governances := []types.MultiSigGovernance{}
	govIter, _ := k.MultiSigGovernances.Iterate(ctx, nil)
	defer govIter.Close()
	for ; govIter.Valid(); govIter.Next() {
		gov, _ := govIter.Value()
		governances = append(governances, gov)
	}
	genesis.MultiSigGovernances = governances

	// Export proposal system
	system, _ := k.GetCommunityProposalSystem(ctx)
	system.CurrentPhase = k.GetCurrentGovernancePhase(ctx)
	genesis.ProposalSystem = &system

	// Export dashboard
	dashboard, _ := k.GetRealTimeDashboard(ctx)
	genesis.Dashboard = &dashboard

	// Export transparency reports
	reports := []types.TransparencyReport{}
	reportIter, _ := k.TransparencyReports.Iterate(ctx, nil)
	defer reportIter.Close()
	for ; reportIter.Valid(); reportIter.Next() {
		report, _ := reportIter.Value()
		reports = append(reports, report)
	}
	genesis.TransparencyReports = reports

	return genesis
}