package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"deshchain/x/charitabletrust/types"
)

// Test QueryParams
func (suite *KeeperTestSuite) TestQueryParams() {
	ctx := suite.ctx

	// Set custom params
	params := types.Params{
		Enabled:                    true,
		MinAllocationAmount:        sdk.NewCoin("unamo", sdk.NewInt(200000000)),
		MaxMonthlyAllocationPerOrg: sdk.NewCoin("unamo", sdk.NewInt(200000000000)),
		ProposalVotingPeriod:       14 * 24 * 60 * 60,
		FraudInvestigationPeriod:   60,
		ImpactReportFrequency:      60,
		DistributionCategories:     []string{"education", "healthcare", "environment"},
	}
	suite.keeper.SetParams(ctx, params)

	// Query params
	res, err := suite.queryClient.Params(ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(params.Enabled, res.Params.Enabled)
	suite.Require().Equal(params.MinAllocationAmount, res.Params.MinAllocationAmount)
	suite.Require().Len(res.Params.DistributionCategories, 3)
}

// Test QueryTrustFundBalance
func (suite *KeeperTestSuite) TestQueryTrustFundBalance() {
	ctx := suite.ctx

	// Initially no balance
	res, err := suite.queryClient.TrustFundBalance(ctx, &types.QueryTrustFundBalanceRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NotNil(res.Balance)

	// Set balance
	balance := types.TrustFundBalance{
		TotalBalance:     sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		AllocatedAmount:  sdk.NewCoin("unamo", sdk.NewInt(300000000)),
		AvailableAmount:  sdk.NewCoin("unamo", sdk.NewInt(700000000)),
		TotalDistributed: sdk.NewCoin("unamo", sdk.NewInt(500000000)),
		LastUpdated:      ctx.BlockTime(),
	}
	suite.keeper.SetTrustFundBalance(ctx, balance)

	// Query balance
	res, err = suite.queryClient.TrustFundBalance(ctx, &types.QueryTrustFundBalanceRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(balance.TotalBalance, res.Balance.TotalBalance)
	suite.Require().Equal(balance.AllocatedAmount, res.Balance.AllocatedAmount)
	suite.Require().Equal(balance.AvailableAmount, res.Balance.AvailableAmount)
}

// Test QueryTrustGovernance
func (suite *KeeperTestSuite) TestQueryTrustGovernance() {
	ctx := suite.ctx

	// Set governance
	governance := types.TrustGovernance{
		Trustees: []types.Trustee{
			{
				Address:     "desh1trustee1test",
				Name:        "Test Trustee 1",
				Role:        "Chairman",
				Status:      "active",
				VotingPower: 100,
			},
			{
				Address:     "desh1trustee2test",
				Name:        "Test Trustee 2",
				Role:        "Secretary",
				Status:      "active",
				VotingPower: 100,
			},
		},
		Quorum:              2,
		ApprovalThreshold:   sdk.NewDecWithPrec(600, 3),
		TransparencyOfficer: "desh1transparency1test",
	}
	suite.keeper.SetTrustGovernance(ctx, governance)

	// Query governance
	res, err := suite.queryClient.TrustGovernance(ctx, &types.QueryTrustGovernanceRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Governance.Trustees, 2)
	suite.Require().Equal(governance.Quorum, res.Governance.Quorum)
	suite.Require().Equal(governance.ApprovalThreshold, res.Governance.ApprovalThreshold)
}

// Test QueryCharitableAllocation
func (suite *KeeperTestSuite) TestQueryCharitableAllocation() {
	ctx := suite.ctx

	// Create test allocation
	allocation := types.CharitableAllocation{
		Id:                    1,
		CharitableOrgWalletId: 1,
		OrganizationName:      "Test Charity",
		Amount:                sdk.NewCoin("unamo", sdk.NewInt(100000000)),
		Purpose:               "Education support",
		Category:              "education",
		Status:                "active",
	}
	suite.keeper.SetCharitableAllocation(ctx, allocation)

	// Query single allocation
	res, err := suite.queryClient.CharitableAllocation(ctx, &types.QueryCharitableAllocationRequest{
		AllocationId: 1,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(allocation.Id, res.Allocation.Id)
	suite.Require().Equal(allocation.OrganizationName, res.Allocation.OrganizationName)
	suite.Require().Equal(allocation.Amount, res.Allocation.Amount)

	// Query non-existent allocation
	res, err = suite.queryClient.CharitableAllocation(ctx, &types.QueryCharitableAllocationRequest{
		AllocationId: 999,
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not found")
}

// Test QueryCharitableAllocations
func (suite *KeeperTestSuite) TestQueryCharitableAllocations() {
	ctx := suite.ctx

	// Create multiple allocations
	for i := 1; i <= 5; i++ {
		allocation := types.CharitableAllocation{
			Id:                    uint64(i),
			CharitableOrgWalletId: uint64(i),
			OrganizationName:      fmt.Sprintf("Charity %d", i),
			Amount:                sdk.NewCoin("unamo", sdk.NewInt(int64(i*100000000))),
			Category:              "education",
			Status:                "active",
		}
		suite.keeper.SetCharitableAllocation(ctx, allocation)
	}

	// Query with pagination
	res, err := suite.queryClient.CharitableAllocations(ctx, &types.QueryCharitableAllocationsRequest{
		Pagination: &query.PageRequest{
			Limit: 3,
		},
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Allocations, 3)
	suite.Require().NotNil(res.Pagination.NextKey)

	// Query next page
	res, err = suite.queryClient.CharitableAllocations(ctx, &types.QueryCharitableAllocationsRequest{
		Pagination: &query.PageRequest{
			Key:   res.Pagination.NextKey,
			Limit: 3,
		},
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Allocations, 2)
	suite.Require().Nil(res.Pagination.NextKey)

	// Query with status filter
	res, err = suite.queryClient.CharitableAllocations(ctx, &types.QueryCharitableAllocationsRequest{
		Status: "active",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Allocations, 5)

	// Query with category filter
	res, err = suite.queryClient.CharitableAllocations(ctx, &types.QueryCharitableAllocationsRequest{
		Category: "education",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Allocations, 5)
}

// Test QueryAllocationProposal
func (suite *KeeperTestSuite) TestQueryAllocationProposal() {
	ctx := suite.ctx

	// Create test proposal
	proposal := types.AllocationProposal{
		Id:               1,
		Title:            "Q1 2025 Distribution",
		Description:      "Quarterly charity distribution",
		TotalAmount:      sdk.NewCoin("unamo", sdk.NewInt(1000000000)),
		Status:           "pending",
		ExpectedOutcomes: "Help 1000 beneficiaries",
	}
	suite.keeper.SetAllocationProposal(ctx, proposal)

	// Query proposal
	res, err := suite.queryClient.AllocationProposal(ctx, &types.QueryAllocationProposalRequest{
		ProposalId: 1,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(proposal.Id, res.Proposal.Id)
	suite.Require().Equal(proposal.Title, res.Proposal.Title)
	suite.Require().Equal(proposal.TotalAmount, res.Proposal.TotalAmount)
}

// Test QueryAllocationProposals
func (suite *KeeperTestSuite) TestQueryAllocationProposals() {
	ctx := suite.ctx

	// Create multiple proposals with different statuses
	statuses := []string{"pending", "approved", "rejected", "executed"}
	for i, status := range statuses {
		proposal := types.AllocationProposal{
			Id:          uint64(i + 1),
			Title:       fmt.Sprintf("Proposal %d", i+1),
			TotalAmount: sdk.NewCoin("unamo", sdk.NewInt(int64((i+1)*100000000))),
			Status:      status,
		}
		suite.keeper.SetAllocationProposal(ctx, proposal)
	}

	// Query all proposals
	res, err := suite.queryClient.AllocationProposals(ctx, &types.QueryAllocationProposalsRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Proposals, 4)

	// Query by status
	res, err = suite.queryClient.AllocationProposals(ctx, &types.QueryAllocationProposalsRequest{
		Status: "pending",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Proposals, 1)
	suite.Require().Equal("pending", res.Proposals[0].Status)
}

// Test QueryImpactReport
func (suite *KeeperTestSuite) TestQueryImpactReport() {
	ctx := suite.ctx

	// Create test impact report
	report := types.ImpactReport{
		Id:                   1,
		AllocationId:         1,
		OrganizationId:       1,
		ReportingPeriod:      "Q1 2025",
		BeneficiariesReached: 200,
		FundsUtilized:        sdk.NewCoin("unamo", sdk.NewInt(80000000)),
		ActivitiesConducted:  []string{"Activity 1", "Activity 2"},
		OutcomesAchieved:     []string{"Outcome 1", "Outcome 2"},
	}
	suite.keeper.SetImpactReport(ctx, report)

	// Query report
	res, err := suite.queryClient.ImpactReport(ctx, &types.QueryImpactReportRequest{
		ReportId: 1,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(report.Id, res.Report.Id)
	suite.Require().Equal(report.BeneficiariesReached, res.Report.BeneficiariesReached)
	suite.Require().Len(res.Report.ActivitiesConducted, 2)
}

// Test QueryImpactReports
func (suite *KeeperTestSuite) TestQueryImpactReports() {
	ctx := suite.ctx

	// Create multiple impact reports
	for i := 1; i <= 3; i++ {
		report := types.ImpactReport{
			Id:                   uint64(i),
			AllocationId:         1,
			OrganizationId:       uint64(i),
			BeneficiariesReached: int32(i * 100),
			FundsUtilized:        sdk.NewCoin("unamo", sdk.NewInt(int64(i*50000000))),
		}
		suite.keeper.SetImpactReport(ctx, report)
	}

	// Query by allocation
	res, err := suite.queryClient.ImpactReports(ctx, &types.QueryImpactReportsRequest{
		AllocationId: 1,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Reports, 3)

	// Query by organization
	res, err = suite.queryClient.ImpactReports(ctx, &types.QueryImpactReportsRequest{
		OrganizationId: 2,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Reports, 1)
	suite.Require().Equal(uint64(2), res.Reports[0].OrganizationId)

	// Query verified reports only
	// First verify one report
	verification := types.Verification{
		IsVerified: true,
		VerifiedBy: "desh1verifier1test",
		Score:      90,
	}
	err = suite.keeper.VerifyImpactReport(ctx, 1, verification)
	suite.Require().NoError(err)

	res, err = suite.queryClient.ImpactReports(ctx, &types.QueryImpactReportsRequest{
		VerifiedOnly: true,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Reports, 1)
	suite.Require().True(res.Reports[0].Verification.IsVerified)
}

// Test QueryFraudAlert
func (suite *KeeperTestSuite) TestQueryFraudAlert() {
	ctx := suite.ctx

	// Create test fraud alert
	alert := types.FraudAlert{
		Id:             1,
		AlertType:      "misuse_of_funds",
		Severity:       "high",
		AllocationId:   1,
		OrganizationId: 1,
		Description:    "Test fraud alert",
		Status:         "pending",
	}
	suite.keeper.SetFraudAlert(ctx, alert)

	// Query alert
	res, err := suite.queryClient.FraudAlert(ctx, &types.QueryFraudAlertRequest{
		AlertId: 1,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Equal(alert.Id, res.Alert.Id)
	suite.Require().Equal(alert.AlertType, res.Alert.AlertType)
	suite.Require().Equal(alert.Severity, res.Alert.Severity)
}

// Test QueryFraudAlerts
func (suite *KeeperTestSuite) TestQueryFraudAlerts() {
	ctx := suite.ctx

	// Create multiple fraud alerts with different severities and statuses
	severities := []string{"low", "medium", "high", "critical"}
	statuses := []string{"pending", "investigating", "investigated", "resolved"}

	for i := 0; i < 4; i++ {
		alert := types.FraudAlert{
			Id:             uint64(i + 1),
			AlertType:      "misuse_of_funds",
			Severity:       severities[i],
			Status:         statuses[i],
			AllocationId:   uint64(i + 1),
			OrganizationId: 1,
		}
		suite.keeper.SetFraudAlert(ctx, alert)
	}

	// Query all alerts
	res, err := suite.queryClient.FraudAlerts(ctx, &types.QueryFraudAlertsRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Alerts, 4)

	// Query by severity
	res, err = suite.queryClient.FraudAlerts(ctx, &types.QueryFraudAlertsRequest{
		Severity: "high",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Alerts, 1)
	suite.Require().Equal("high", res.Alerts[0].Severity)

	// Query by status
	res, err = suite.queryClient.FraudAlerts(ctx, &types.QueryFraudAlertsRequest{
		Status: "pending",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Alerts, 1)
	suite.Require().Equal("pending", res.Alerts[0].Status)

	// Query by organization
	res, err = suite.queryClient.FraudAlerts(ctx, &types.QueryFraudAlertsRequest{
		OrganizationId: 1,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().Len(res.Alerts, 4)
}

// Test QueryDonorStatistics
func (suite *KeeperTestSuite) TestQueryDonorStatistics() {
	ctx := suite.ctx

	// Create allocations and reports for statistics
	for i := 1; i <= 3; i++ {
		allocation := types.CharitableAllocation{
			Id:               uint64(i),
			OrganizationName: fmt.Sprintf("Charity %d", i),
			Amount:           sdk.NewCoin("unamo", sdk.NewInt(int64(i*100000000))),
			Category:         "education",
			Status:           "distributed",
		}
		suite.keeper.SetCharitableAllocation(ctx, allocation)

		report := types.ImpactReport{
			Id:                   uint64(i),
			AllocationId:         uint64(i),
			BeneficiariesReached: int32(i * 100),
			FundsUtilized:        sdk.NewCoin("unamo", sdk.NewInt(int64(i*80000000))),
		}
		suite.keeper.SetImpactReport(ctx, report)
	}

	// Query statistics - implementation would aggregate data
	res, err := suite.queryClient.DonorStatistics(ctx, &types.QueryDonorStatisticsRequest{
		DonorAddress: "desh1donor1test",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	// The actual implementation would calculate these statistics
	suite.Require().NotNil(res.Statistics)
}