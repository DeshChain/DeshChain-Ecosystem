package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"deshchain/x/charitabletrust/keeper"
	"deshchain/x/charitabletrust/types"
)

// Test InitGenesis and ExportGenesis
func (suite *KeeperTestSuite) TestGenesisInitAndExport() {
	ctx := suite.ctx
	k := suite.keeper

	// Create test genesis state
	genState := &types.GenesisState{
		Params: types.Params{
			Enabled:                    true,
			MinAllocationAmount:        sdk.NewCoin("unamo", sdk.NewInt(100000000)),
			MaxMonthlyAllocationPerOrg: sdk.NewCoin("unamo", sdk.NewInt(100000000000)),
			ProposalVotingPeriod:       7 * 24 * 60 * 60, // 7 days
			FraudInvestigationPeriod:   30,
			ImpactReportFrequency:      30,
			DistributionCategories: []string{
				"education", "healthcare", "rural_development",
				"women_empowerment", "emergency_relief",
			},
		},
		TrustGovernance: &types.TrustGovernance{
			Trustees: []types.Trustee{
				{
					Address:      "desh1trustee1genesis",
					Name:         "Genesis Trustee 1",
					Role:         "Chairman",
					Expertise:    "Social Impact",
					AppointedAt:  ctx.BlockTime(),
					TermEndDate:  ctx.BlockTime().Add(2 * 365 * 24 * time.Hour),
					Status:       "active",
					VotingPower:  100,
				},
				{
					Address:      "desh1trustee2genesis",
					Name:         "Genesis Trustee 2",
					Role:         "Secretary",
					Expertise:    "Finance",
					AppointedAt:  ctx.BlockTime(),
					TermEndDate:  ctx.BlockTime().Add(2 * 365 * 24 * time.Hour),
					Status:       "active",
					VotingPower:  100,
				},
				{
					Address:      "desh1trustee3genesis",
					Name:         "Genesis Trustee 3",
					Role:         "Member",
					Expertise:    "Healthcare",
					AppointedAt:  ctx.BlockTime(),
					TermEndDate:  ctx.BlockTime().Add(2 * 365 * 24 * time.Hour),
					Status:       "active",
					VotingPower:  100,
				},
			},
			Quorum:              2,
			ApprovalThreshold:   sdk.NewDecWithPrec(667, 3), // 66.7%
			AdvisoryCommittee:   []string{},
			TransparencyOfficer: "desh1transparencygenesis",
			NextElection:        ctx.BlockTime().Add(2 * 365 * 24 * time.Hour),
		},
		TrustFundBalance: &types.TrustFundBalance{
			TotalBalance:     sdk.NewCoin("unamo", sdk.NewInt(5000000000)),
			AllocatedAmount:  sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
			AvailableAmount:  sdk.NewCoin("unamo", sdk.NewInt(4000000000)),
			TotalDistributed: sdk.NewCoin("unamo", sdk.NewInt(2000000000)),
			LastUpdated:      ctx.BlockTime(),
		},
		Allocations: []types.CharitableAllocation{
			{
				Id:                    1,
				CharitableOrgWalletId: 1,
				OrganizationName:      "Genesis Education Foundation",
				Amount:                sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				Purpose:               "Rural education development",
				Category:              "education",
				ProposalId:            1,
				ApprovedBy:            []string{"desh1trustee1genesis", "desh1trustee2genesis"},
				AllocatedAt:           ctx.BlockTime(),
				ExpectedImpact:        "Educate 1000 rural children",
				Monitoring: &types.MonitoringRequirements{
					ReportingFrequency:     30,
					RequiredReports:        []string{"impact", "financial"},
					Kpis:                   []string{"students_enrolled", "literacy_rate"},
					MonitoringDuration:     365,
					SiteVisitsRequired:     true,
					FinancialAuditRequired: true,
				},
				Status: "active",
				Distribution: &types.Distribution{
					TxHash:        "",
					DistributedAt: time.Time{},
					DistributedBy: "",
				},
			},
			{
				Id:                    2,
				CharitableOrgWalletId: 2,
				OrganizationName:      "Genesis Health Trust",
				Amount:                sdk.NewCoin("unamo", sdk.NewInt(500000000)),
				Purpose:               "Healthcare accessibility",
				Category:              "healthcare",
				ProposalId:            1,
				ApprovedBy:            []string{"desh1trustee1genesis", "desh1trustee3genesis"},
				AllocatedAt:           ctx.BlockTime(),
				ExpectedImpact:        "Serve 5000 patients",
				Monitoring: &types.MonitoringRequirements{
					ReportingFrequency:     30,
					RequiredReports:        []string{"impact", "financial"},
					Kpis:                   []string{"patients_served", "recovery_rate"},
					MonitoringDuration:     365,
					SiteVisitsRequired:     true,
					FinancialAuditRequired: true,
				},
				Status: "distributed",
				Distribution: &types.Distribution{
					TxHash:        "0x1234567890abcdef",
					DistributedAt: ctx.BlockTime().Add(-30 * 24 * time.Hour),
					DistributedBy: "desh1treasurygenesis",
				},
			},
		},
		Proposals: []types.AllocationProposal{
			{
				Id:               1,
				Title:            "Genesis Distribution Q1 2025",
				Description:      "Initial distribution to verified charities",
				TotalAmount:      sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
				Allocations:      []types.ProposedAllocation{},
				ProposedBy:       "desh1trustee1genesis",
				ProposedAt:       ctx.BlockTime().Add(-60 * 24 * time.Hour),
				VotingStartTime:  ctx.BlockTime().Add(-60 * 24 * time.Hour),
				VotingEndTime:    ctx.BlockTime().Add(-53 * 24 * time.Hour),
				Status:           "executed",
				Votes: []types.Vote{
					{
						Voter:     "desh1trustee1genesis",
						VoteType:  "yes",
						VotedAt:   ctx.BlockTime().Add(-59 * 24 * time.Hour),
						Comments:  "Approved for launch",
						VoteValue: 100,
					},
					{
						Voter:     "desh1trustee2genesis",
						VoteType:  "yes",
						VotedAt:   ctx.BlockTime().Add(-58 * 24 * time.Hour),
						Comments:  "Good distribution plan",
						VoteValue: 100,
					},
				},
				Documents:        []string{"https://ipfs.io/ipfs/Qm...genesisproposal"},
				ExpectedOutcomes: "Launch charitable distribution system",
			},
		},
		ImpactReports: []types.ImpactReport{
			{
				Id:                   1,
				AllocationId:         2,
				OrganizationId:       2,
				ReportingPeriod:      "Genesis Month 1",
				SubmittedAt:          ctx.BlockTime(),
				BeneficiariesReached: 500,
				FundsUtilized:        sdk.NewCoin("unamo", sdk.NewInt(400000000)),
				ActivitiesConducted: []string{
					"Established 3 mobile health units",
					"Conducted 50 health camps",
					"Distributed medicines to 500 families",
				},
				OutcomesAchieved: []string{
					"500 patients treated",
					"95% treatment success rate",
					"3 villages covered",
				},
				Challenges: []string{
					"Remote area accessibility",
					"Language barriers",
				},
				SupportingDocuments: []string{
					"https://ipfs.io/ipfs/Qm...healthreport1",
					"https://ipfs.io/ipfs/Qm...photos1",
				},
				Verification: &types.Verification{
					VerifiedBy:       "desh1verifiergenesis",
					VerifiedAt:       ctx.BlockTime().Add(-7 * 24 * time.Hour),
					IsVerified:       true,
					VerificationNote: "Comprehensive verification completed",
					SiteVisitReport:  "All activities verified on ground",
					Score:            92,
				},
				ImpactMetrics: map[string]string{
					"patients_served":     "500",
					"treatment_success":   "95%",
					"villages_covered":    "3",
					"medicine_distributed": "500 families",
				},
			},
		},
		FraudAlerts: []types.FraudAlert{
			{
				Id:               1,
				AlertType:        "suspicious_activity",
				Severity:         "medium",
				AllocationId:     1,
				OrganizationId:   1,
				ReportedBy:       "desh1whistleblowergenesis",
				ReportedAt:       ctx.BlockTime().Add(-15 * 24 * time.Hour),
				Description:      "Unusual fund transfer patterns detected",
				Evidence:         []string{"Bank statement analysis", "Transaction logs"},
				Status:           "investigated",
				Investigation: &types.Investigation{
					InvestigatorId:   "desh1investigatorgenesis",
					StartedAt:        ctx.BlockTime().Add(-14 * 24 * time.Hour),
					CompletedAt:      ctx.BlockTime().Add(-7 * 24 * time.Hour),
					Findings:         []string{"No fraud detected", "Legitimate bulk purchases"},
					Recommendation:   "no_action_required",
					ActionsTaken:     []string{"Verified with organization", "Reviewed documentation"},
					EvidenceGathered: []string{"Invoice copies", "Supplier contracts"},
				},
				Resolution: &types.Resolution{
					ResolvedAt:    ctx.BlockTime().Add(-7 * 24 * time.Hour),
					ResolvedBy:    "desh1investigatorgenesis",
					Resolution:    "Alert cleared - no fraud found",
					ActionsTaken:  []string{"Case closed", "Documentation updated"},
					RecoveredFunds: sdk.NewCoin("unamo", sdk.NewInt(0)),
				},
				MonetaryImpact:   sdk.NewCoin("unamo", sdk.NewInt(0)),
				AffectedEntities: []string{"Genesis Education Foundation"},
			},
		},
		AllocationCount: 2,
		ProposalCount:   1,
		ReportCount:     1,
		AlertCount:      1,
	}

	// Initialize genesis
	keeper.InitGenesis(ctx, k, genState)

	// Verify params were set
	params := k.GetParams(ctx)
	suite.Require().Equal(genState.Params.Enabled, params.Enabled)
	suite.Require().Equal(genState.Params.MinAllocationAmount, params.MinAllocationAmount)
	suite.Require().Len(params.DistributionCategories, 5)

	// Verify governance was set
	governance, found := k.GetTrustGovernance(ctx)
	suite.Require().True(found)
	suite.Require().Len(governance.Trustees, 3)
	suite.Require().Equal("desh1trustee1genesis", governance.Trustees[0].Address)
	suite.Require().Equal(uint32(2), governance.Quorum)

	// Verify trust fund balance was set
	balance, found := k.GetTrustFundBalance(ctx)
	suite.Require().True(found)
	suite.Require().Equal(genState.TrustFundBalance.TotalBalance, balance.TotalBalance)
	suite.Require().Equal(genState.TrustFundBalance.AllocatedAmount, balance.AllocatedAmount)

	// Verify allocations were set
	allocations := k.GetAllCharitableAllocations(ctx)
	suite.Require().Len(allocations, 2)
	suite.Require().Equal("Genesis Education Foundation", allocations[0].OrganizationName)
	suite.Require().Equal("Genesis Health Trust", allocations[1].OrganizationName)
	suite.Require().Equal("active", allocations[0].Status)
	suite.Require().Equal("distributed", allocations[1].Status)

	// Verify proposals were set
	proposals := k.GetAllAllocationProposals(ctx)
	suite.Require().Len(proposals, 1)
	suite.Require().Equal("Genesis Distribution Q1 2025", proposals[0].Title)
	suite.Require().Equal("executed", proposals[0].Status)
	suite.Require().Len(proposals[0].Votes, 2)

	// Verify impact reports were set
	report, found := k.GetImpactReport(ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(int32(500), report.BeneficiariesReached)
	suite.Require().NotNil(report.Verification)
	suite.Require().True(report.Verification.IsVerified)
	suite.Require().Equal(int32(92), report.Verification.Score)

	// Verify fraud alerts were set
	alert, found := k.GetFraudAlert(ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal("suspicious_activity", alert.AlertType)
	suite.Require().Equal("investigated", alert.Status)
	suite.Require().NotNil(alert.Investigation)
	suite.Require().NotNil(alert.Resolution)

	// Verify counters were set
	nextAllocationID := k.IncrementAllocationCount(ctx)
	suite.Require().Equal(uint64(3), nextAllocationID)
	nextProposalID := k.IncrementProposalCount(ctx)
	suite.Require().Equal(uint64(2), nextProposalID)
	nextReportID := k.IncrementReportCount(ctx)
	suite.Require().Equal(uint64(2), nextReportID)
	nextAlertID := k.IncrementAlertCount(ctx)
	suite.Require().Equal(uint64(2), nextAlertID)

	// Export genesis and verify it matches
	exportedGenesis := keeper.ExportGenesis(ctx, k)

	// Verify exported params
	suite.Require().Equal(genState.Params.Enabled, exportedGenesis.Params.Enabled)
	suite.Require().Equal(genState.Params.MinAllocationAmount, exportedGenesis.Params.MinAllocationAmount)

	// Verify exported governance
	suite.Require().NotNil(exportedGenesis.TrustGovernance)
	suite.Require().Len(exportedGenesis.TrustGovernance.Trustees, 3)
	suite.Require().Equal(genState.TrustGovernance.Quorum, exportedGenesis.TrustGovernance.Quorum)

	// Verify exported balance
	suite.Require().NotNil(exportedGenesis.TrustFundBalance)
	suite.Require().Equal(genState.TrustFundBalance.TotalBalance, exportedGenesis.TrustFundBalance.TotalBalance)

	// Verify exported allocations
	suite.Require().Len(exportedGenesis.Allocations, 2)
	suite.Require().Equal(genState.Allocations[0].OrganizationName, exportedGenesis.Allocations[0].OrganizationName)

	// Verify exported proposals
	suite.Require().Len(exportedGenesis.Proposals, 1)
	suite.Require().Equal(genState.Proposals[0].Title, exportedGenesis.Proposals[0].Title)

	// Verify exported reports
	suite.Require().Len(exportedGenesis.ImpactReports, 1)
	suite.Require().Equal(genState.ImpactReports[0].BeneficiariesReached, exportedGenesis.ImpactReports[0].BeneficiariesReached)

	// Verify exported alerts
	suite.Require().Len(exportedGenesis.FraudAlerts, 1)
	suite.Require().Equal(genState.FraudAlerts[0].AlertType, exportedGenesis.FraudAlerts[0].AlertType)

	// Verify exported counters (should be incremented by the tests above)
	suite.Require().Equal(uint64(2), exportedGenesis.AllocationCount) // Original value from genesis
	suite.Require().Equal(uint64(1), exportedGenesis.ProposalCount)
	suite.Require().Equal(uint64(1), exportedGenesis.ReportCount)
	suite.Require().Equal(uint64(1), exportedGenesis.AlertCount)
}

// Test empty genesis initialization
func (suite *KeeperTestSuite) TestEmptyGenesisInit() {
	ctx := suite.ctx
	k := suite.keeper

	// Create minimal genesis state
	genState := &types.GenesisState{
		Params: types.DefaultParams(),
	}

	// Initialize genesis
	keeper.InitGenesis(ctx, k, genState)

	// Verify default params were set
	params := k.GetParams(ctx)
	suite.Require().True(params.Enabled)
	suite.Require().Equal("unamo", params.MinAllocationAmount.Denom)

	// Verify empty collections
	allocations := k.GetAllCharitableAllocations(ctx)
	suite.Require().Len(allocations, 0)

	proposals := k.GetAllAllocationProposals(ctx)
	suite.Require().Len(proposals, 0)

	// Verify governance is not set
	_, found := k.GetTrustGovernance(ctx)
	suite.Require().False(found)

	// Verify counters start at 0
	nextAllocationID := k.IncrementAllocationCount(ctx)
	suite.Require().Equal(uint64(1), nextAllocationID)
}

// Test invalid genesis data handling
func (suite *KeeperTestSuite) TestInvalidGenesisData() {
	ctx := suite.ctx
	k := suite.keeper

	// Test with invalid trustee data
	genState := &types.GenesisState{
		Params: types.DefaultParams(),
		TrustGovernance: &types.TrustGovernance{
			Trustees: []types.Trustee{
				{
					Address:     "", // Invalid empty address
					Name:        "Invalid Trustee",
					Status:      "active",
					VotingPower: 100,
				},
			},
			Quorum: 1,
		},
	}

	// This should handle the invalid data gracefully
	// In a real implementation, you might want to validate and panic/return error
	keeper.InitGenesis(ctx, k, genState)

	// For now, we just verify it doesn't crash
	governance, found := k.GetTrustGovernance(ctx)
	suite.Require().True(found)
	suite.Require().Len(governance.Trustees, 1)
}

// Test genesis with large dataset
func (suite *KeeperTestSuite) TestLargeGenesisDataset() {
	ctx := suite.ctx
	k := suite.keeper

	// Create genesis with many allocations
	allocations := make([]types.CharitableAllocation, 100)
	for i := 0; i < 100; i++ {
		allocations[i] = types.CharitableAllocation{
			Id:                    uint64(i + 1),
			CharitableOrgWalletId: uint64(i + 1),
			OrganizationName:      fmt.Sprintf("Charity %d", i+1),
			Amount:                sdk.NewCoin("unamo", sdk.NewInt(int64((i+1)*1000000))),
			Category:              "education",
			Status:                "active",
		}
	}

	genState := &types.GenesisState{
		Params:          types.DefaultParams(),
		Allocations:     allocations,
		AllocationCount: 100,
	}

	// Initialize genesis
	keeper.InitGenesis(ctx, k, genState)

	// Verify all allocations were set
	storedAllocations := k.GetAllCharitableAllocations(ctx)
	suite.Require().Len(storedAllocations, 100)

	// Verify counter was set correctly
	nextID := k.IncrementAllocationCount(ctx)
	suite.Require().Equal(uint64(101), nextID)

	// Export and verify performance
	exportedGenesis := keeper.ExportGenesis(ctx, k)
	suite.Require().Len(exportedGenesis.Allocations, 100)
}