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

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/deshchain/deshchain/x/treasury/types"
)

// Keeper defines the treasury keeper
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	logger       log.Logger

	// Authority address for governance proposals
	authority string

	// Collections for storing treasury data
	Schema collections.Schema
	
	// Community Fund storage
	CommunityFundProposals    collections.Map[uint64, types.CommunityFundProposal]
	CommunityFundBalance      collections.Item[types.CommunityFundBalance]
	CommunityFundTransactions collections.Map[string, types.CommunityFundTransaction]
	CommunityGovernance       collections.Item[types.CommunityGovernance]
	
	// Development Fund storage
	DevelopmentFundProposals    collections.Map[uint64, types.DevelopmentFundProposal]
	DevelopmentFundBalance      collections.Item[types.DevelopmentFundBalance]
	DevelopmentFundTransactions collections.Map[string, types.DevelopmentFundTransaction]
	
	// Multi-signature governance storage
	MultiSigGovernance collections.Map[uint64, types.MultiSigGovernance]
	
	// Community Proposal System storage
	CommunityProposalSystem collections.Item[types.CommunityProposalSystem]
	PhaseTransitions        collections.Map[uint64, types.PhaseTransition]
	
	// Real-time dashboard storage
	RealTimeDashboard collections.Item[types.RealTimeDashboard]
	
	// External keeper references
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
}

// NewKeeper creates a new treasury keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	
	k := Keeper{
		cdc:           cdc,
		storeService:  storeService,
		logger:        logger,
		authority:     authority,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		
		// Initialize collections
		CommunityFundProposals: collections.NewMap(
			sb, 
			types.CommunityFundProposalKey,
			"community_fund_proposals",
			collections.Uint64Key,
			codec.CollValue[types.CommunityFundProposal](cdc),
		),
		CommunityFundBalance: collections.NewItem(
			sb,
			types.CommunityFundBalanceKey,
			"community_fund_balance",
			codec.CollValue[types.CommunityFundBalance](cdc),
		),
		CommunityFundTransactions: collections.NewMap(
			sb,
			types.CommunityFundTransactionKey,
			"community_fund_transactions",
			collections.StringKey,
			codec.CollValue[types.CommunityFundTransaction](cdc),
		),
		CommunityGovernance: collections.NewItem(
			sb,
			types.CommunityFundGovernanceKey,
			"community_governance",
			codec.CollValue[types.CommunityGovernance](cdc),
		),
		DevelopmentFundProposals: collections.NewMap(
			sb,
			types.DevelopmentFundProposalKey,
			"development_fund_proposals",
			collections.Uint64Key,
			codec.CollValue[types.DevelopmentFundProposal](cdc),
		),
		DevelopmentFundBalance: collections.NewItem(
			sb,
			types.DevelopmentFundBalanceKey,
			"development_fund_balance",
			codec.CollValue[types.DevelopmentFundBalance](cdc),
		),
		DevelopmentFundTransactions: collections.NewMap(
			sb,
			types.DevelopmentFundTransactionKey,
			"development_fund_transactions",
			collections.StringKey,
			codec.CollValue[types.DevelopmentFundTransaction](cdc),
		),
		MultiSigGovernance: collections.NewMap(
			sb,
			types.MultiSigGovernanceKey,
			"multisig_governance",
			collections.Uint64Key,
			codec.CollValue[types.MultiSigGovernance](cdc),
		),
		CommunityProposalSystem: collections.NewItem(
			sb,
			types.CommunityProposalSystemKey,
			"community_proposal_system",
			codec.CollValue[types.CommunityProposalSystem](cdc),
		),
		PhaseTransitions: collections.NewMap(
			sb,
			types.PhaseTransitionKey,
			"phase_transitions",
			collections.Uint64Key,
			codec.CollValue[types.PhaseTransition](cdc),
		),
		RealTimeDashboard: collections.NewItem(
			sb,
			collections.NewPrefix(500), // Custom prefix for dashboard
			"realtime_dashboard",
			codec.CollValue[types.RealTimeDashboard](cdc),
		),
	}
	
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	
	return k
}

// GetAuthority returns the module's authority address
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns the module logger
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", "x/"+types.ModuleName)
}

// InitializeModule initializes the treasury module with default settings
func (k Keeper) InitializeModule(ctx sdk.Context) error {
	// Initialize Community Fund Balance
	communityBalance := types.CommunityFundBalance{
		TotalBalance:     sdk.NewCoin("namo", sdk.NewInt(types.CommunityFundAllocation)),
		AllocatedAmount:  sdk.NewCoin("namo", sdk.ZeroInt()),
		AvailableAmount:  sdk.NewCoin("namo", sdk.NewInt(types.CommunityFundAllocation)),
		ReservedAmount:   sdk.NewCoin("namo", sdk.ZeroInt()),
		PendingAmount:    sdk.NewCoin("namo", sdk.ZeroInt()),
		LastUpdateHeight: ctx.BlockHeight(),
		LastUpdateTime:   ctx.BlockTime(),
	}
	if err := k.CommunityFundBalance.Set(ctx, communityBalance); err != nil {
		return err
	}

	// Initialize Development Fund Balance
	developmentBalance := types.DevelopmentFundBalance{
		TotalBalance:     sdk.NewCoin("namo", sdk.NewInt(types.DevelopmentFundAllocation)),
		AllocatedAmount:  sdk.NewCoin("namo", sdk.ZeroInt()),
		AvailableAmount:  sdk.NewCoin("namo", sdk.NewInt(types.DevelopmentFundAllocation)),
		PendingAmount:    sdk.NewCoin("namo", sdk.ZeroInt()),
		EmergencyReserve: sdk.NewCoin("namo", sdk.ZeroInt()),
		LastUpdateHeight: ctx.BlockHeight(),
		LastUpdateTime:   ctx.BlockTime(),
	}
	if err := k.DevelopmentFundBalance.Set(ctx, developmentBalance); err != nil {
		return err
	}

	// Initialize Community Governance with defaults
	if err := k.CommunityGovernance.Set(ctx, types.DefaultCommunityGovernance); err != nil {
		return err
	}

	// Initialize MultiSig Governance
	if err := k.MultiSigGovernance.Set(ctx, 1, types.DefaultCommunityFundMultiSig); err != nil {
		return err
	}
	if err := k.MultiSigGovernance.Set(ctx, 2, types.DefaultDevelopmentFundMultiSig); err != nil {
		return err
	}

	// Initialize Community Proposal System
	launchDate := ctx.BlockTime()
	activationDate := launchDate.AddDate(3, 0, 0) // 3 years after launch
	
	proposalSystem := types.CommunityProposalSystem{
		ID:             1,
		Name:           "DeshChain Community Proposal System",
		Description:    "Phased governance system transitioning from founder-led to community-driven",
		LaunchDate:     launchDate,
		ActivationDate: activationDate,
		CurrentPhase:   types.PhaseFounderDriven,
		GovernanceParameters: types.PhaseGovernanceParams{
			CurrentPhase:  types.PhaseFounderDriven,
			FounderDriven: types.DefaultFounderDrivenParams,
		},
		Status: types.SystemStatusActive,
	}
	if err := k.CommunityProposalSystem.Set(ctx, proposalSystem); err != nil {
		return err
	}

	// Initialize Real-Time Dashboard
	dashboard := types.RealTimeDashboard{
		LastUpdated:       ctx.BlockTime(),
		UpdateFrequency:   types.DashboardUpdateInterval,
		TransparencyScore: 10,
		ComplianceScore:   10,
		Status:            types.DashboardStatusOnline,
	}
	if err := k.RealTimeDashboard.Set(ctx, dashboard); err != nil {
		return err
	}

	k.Logger().Info("Treasury module initialized successfully")
	return nil
}

// UpdateDashboard updates the real-time dashboard with latest metrics
func (k Keeper) UpdateDashboard(ctx sdk.Context) error {
	dashboard, err := k.RealTimeDashboard.Get(ctx)
	if err != nil {
		return err
	}

	// Update community fund metrics
	communityBalance, err := k.CommunityFundBalance.Get(ctx)
	if err != nil {
		return err
	}
	
	dashboard.CommunityFund = types.CommunityFundMetrics{
		TotalAllocation:  sdk.NewCoin("namo", sdk.NewInt(types.CommunityFundAllocation)),
		CurrentBalance:   communityBalance.TotalBalance,
		AllocatedAmount:  communityBalance.AllocatedAmount,
		SpentAmount:      communityBalance.TotalBalance.Sub(communityBalance.AvailableAmount),
		RemainingAmount:  communityBalance.AvailableAmount,
		PendingAmount:    communityBalance.PendingAmount,
		ReservedAmount:   communityBalance.ReservedAmount,
		LastUpdated:      ctx.BlockTime(),
	}

	// Update development fund metrics
	developmentBalance, err := k.DevelopmentFundBalance.Get(ctx)
	if err != nil {
		return err
	}
	
	dashboard.DevelopmentFund = types.DevelopmentFundMetrics{
		TotalAllocation:  sdk.NewCoin("namo", sdk.NewInt(types.DevelopmentFundAllocation)),
		CurrentBalance:   developmentBalance.TotalBalance,
		AllocatedAmount:  developmentBalance.AllocatedAmount,
		SpentAmount:      developmentBalance.TotalBalance.Sub(developmentBalance.AvailableAmount),
		RemainingAmount:  developmentBalance.AvailableAmount,
		PendingAmount:    developmentBalance.PendingAmount,
		EmergencyReserve: developmentBalance.EmergencyReserve,
		LastUpdated:      ctx.BlockTime(),
	}

	// Update overall metrics
	totalFunds := communityBalance.TotalBalance.Add(developmentBalance.TotalBalance)
	totalAllocated := communityBalance.AllocatedAmount.Add(developmentBalance.AllocatedAmount)
	totalSpent := totalFunds.Sub(communityBalance.AvailableAmount.Add(developmentBalance.AvailableAmount))
	totalRemaining := communityBalance.AvailableAmount.Add(developmentBalance.AvailableAmount)

	dashboard.OverallMetrics = types.OverallFundMetrics{
		TotalFunds:          totalFunds,
		TotalAllocated:      totalAllocated,
		TotalSpent:          totalSpent,
		TotalRemaining:      totalRemaining,
		OverallTransparency: 10,
		OverallCompliance:   10,
		LastUpdated:         ctx.BlockTime(),
	}

	dashboard.LastUpdated = ctx.BlockTime()
	
	return k.RealTimeDashboard.Set(ctx, dashboard)
}

// CheckPhaseTransition checks if it's time to transition to the next governance phase
func (k Keeper) CheckPhaseTransition(ctx sdk.Context) error {
	system, err := k.CommunityProposalSystem.Get(ctx)
	if err != nil {
		return err
	}

	currentTime := ctx.BlockTime()
	yearsSinceLaunch := currentTime.Sub(system.LaunchDate).Hours() / (24 * 365)

	var newPhase types.ProposalSystemPhase
	switch {
	case yearsSinceLaunch < 3:
		newPhase = types.PhaseFounderDriven
	case yearsSinceLaunch < 4:
		newPhase = types.PhaseTransitional
	case yearsSinceLaunch < 7:
		newPhase = types.PhaseCommunityProposal
	default:
		newPhase = types.PhaseFullGovernance
	}

	// If phase needs to change
	if newPhase != system.CurrentPhase {
		// Record the transition
		transition := types.PhaseTransition{
			From:           system.CurrentPhase,
			To:             newPhase,
			TransitionDate: currentTime,
			Reason:         "Automatic phase transition based on timeline",
			ApprovedBy:     sdk.AccAddress(k.authority),
			TransitionMetrics: types.TransitionMetrics{
				CommunityReadiness: 8,
				SystemHealth:       9,
				TransitionRisk:     2,
			},
		}

		// Get next transition ID
		transitionID := k.getNextTransitionID(ctx)
		if err := k.PhaseTransitions.Set(ctx, transitionID, transition); err != nil {
			return err
		}

		// Update system phase
		system.CurrentPhase = newPhase
		system.PhaseHistory = append(system.PhaseHistory, transition)

		// Update governance parameters based on new phase
		switch newPhase {
		case types.PhaseTransitional:
			system.GovernanceParameters.CurrentPhase = types.PhaseTransitional
			system.GovernanceParameters.Transitional = types.DefaultTransitionalParams
		case types.PhaseCommunityProposal:
			system.GovernanceParameters.CurrentPhase = types.PhaseCommunityProposal
			system.GovernanceParameters.CommunityProposal = types.DefaultCommunityProposalParams
		case types.PhaseFullGovernance:
			system.GovernanceParameters.CurrentPhase = types.PhaseFullGovernance
			system.GovernanceParameters.FullGovernance = types.DefaultFullGovernanceParams
		}

		if err := k.CommunityProposalSystem.Set(ctx, system); err != nil {
			return err
		}

		// Emit event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePhaseTransition,
				sdk.NewAttribute("from", string(transition.From)),
				sdk.NewAttribute("to", string(transition.To)),
				sdk.NewAttribute("date", transition.TransitionDate.String()),
			),
		)

		k.Logger().Info("Governance phase transitioned",
			"from", transition.From,
			"to", transition.To,
			"years_since_launch", yearsSinceLaunch)
	}

	return nil
}

// getNextTransitionID returns the next available transition ID
func (k Keeper) getNextTransitionID(ctx sdk.Context) uint64 {
	var maxID uint64
	iterator, err := k.PhaseTransitions.Iterate(ctx, nil)
	if err != nil {
		return 1
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key, err := iterator.Key()
		if err != nil {
			continue
		}
		if key > maxID {
			maxID = key
		}
	}
	return maxID + 1
}

// GetCommunityFundBalance returns the current community fund balance
func (k Keeper) GetCommunityFundBalance(ctx sdk.Context) (types.CommunityFundBalance, error) {
	return k.CommunityFundBalance.Get(ctx)
}

// GetDevelopmentFundBalance returns the current development fund balance
func (k Keeper) GetDevelopmentFundBalance(ctx sdk.Context) (types.DevelopmentFundBalance, error) {
	return k.DevelopmentFundBalance.Get(ctx)
}

// GetCurrentGovernancePhase returns the current governance phase
func (k Keeper) GetCurrentGovernancePhase(ctx sdk.Context) (types.ProposalSystemPhase, error) {
	system, err := k.CommunityProposalSystem.Get(ctx)
	if err != nil {
		return "", err
	}
	return system.CurrentPhase, nil
}

// ProcessCommunityProposal processes a community fund proposal
func (k Keeper) ProcessCommunityProposal(ctx sdk.Context, proposal types.CommunityFundProposal) error {
	// Validate proposal against current phase rules
	phase, err := k.GetCurrentGovernancePhase(ctx)
	if err != nil {
		return err
	}

	// Check if community proposals are allowed in current phase
	if phase == types.PhaseFounderDriven {
		return fmt.Errorf("community proposals not allowed in founder-driven phase")
	}

	// Validate proposal amount
	balance, err := k.GetCommunityFundBalance(ctx)
	if err != nil {
		return err
	}

	if proposal.RequestedAmount.IsGTE(balance.AvailableAmount) {
		return fmt.Errorf("requested amount exceeds available balance")
	}

	// Check category limits
	categoryLimit, exists := types.CategoryLimits[proposal.Category]
	if !exists {
		return fmt.Errorf("invalid proposal category: %s", proposal.Category)
	}

	// Calculate category allocation
	maxCategoryAmount := sdk.NewDecFromInt(balance.TotalBalance.Amount).
		Mul(categoryLimit).
		TruncateInt()
	
	if proposal.RequestedAmount.Amount.GT(maxCategoryAmount) {
		return fmt.Errorf("requested amount exceeds category limit")
	}

	// Set proposal fields
	proposal.SubmissionTime = ctx.BlockTime()
	proposal.VotingEndTime = ctx.BlockTime().Add(types.DefaultCommunityGovernance.VotingPeriod)
	proposal.Status = types.StatusPending
	proposal.TransparencyScore = 10

	// Store proposal
	proposalID := k.getNextCommunityProposalID(ctx)
	proposal.ProposalID = proposalID
	
	if err := k.CommunityFundProposals.Set(ctx, proposalID, proposal); err != nil {
		return err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalSubmitted,
			sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute("category", string(proposal.Category)),
			sdk.NewAttribute("amount", proposal.RequestedAmount.String()),
		),
	)

	return nil
}

// getNextCommunityProposalID returns the next available community proposal ID
func (k Keeper) getNextCommunityProposalID(ctx sdk.Context) uint64 {
	var maxID uint64
	iterator, err := k.CommunityFundProposals.Iterate(ctx, nil)
	if err != nil {
		return 1
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key, err := iterator.Key()
		if err != nil {
			continue
		}
		if key > maxID {
			maxID = key
		}
	}
	return maxID + 1
}

// AllocateCommunityFunds allocates funds from the community fund
func (k Keeper) AllocateCommunityFunds(ctx sdk.Context, proposalID uint64) error {
	proposal, err := k.CommunityFundProposals.Get(ctx, proposalID)
	if err != nil {
		return err
	}

	if proposal.Status != types.StatusPassed {
		return fmt.Errorf("proposal not passed")
	}

	balance, err := k.GetCommunityFundBalance(ctx)
	if err != nil {
		return err
	}

	// Update balance
	balance.AllocatedAmount = balance.AllocatedAmount.Add(proposal.RequestedAmount)
	balance.AvailableAmount = balance.AvailableAmount.Sub(proposal.RequestedAmount)
	balance.LastUpdateHeight = ctx.BlockHeight()
	balance.LastUpdateTime = ctx.BlockTime()

	if err := k.CommunityFundBalance.Set(ctx, balance); err != nil {
		return err
	}

	// Create transaction record
	tx := types.CommunityFundTransaction{
		TxID:        fmt.Sprintf("cf-%d-%d", proposalID, ctx.BlockHeight()),
		ProposalID:  proposalID,
		From:        sdk.AccAddress(types.CommunityFundModuleName),
		Amount:      proposal.RequestedAmount,
		Type:        types.TxTypeAllocation,
		Category:    proposal.Category,
		Description: proposal.Description,
		Timestamp:   ctx.BlockTime(),
		BlockHeight: ctx.BlockHeight(),
		Status:      types.TxStatusConfirmed,
		Verified:    true,
	}

	if err := k.CommunityFundTransactions.Set(ctx, tx.TxID, tx); err != nil {
		return err
	}

	// Update proposal status
	proposal.Status = types.StatusExecuting
	proposal.ExecutionTime = &ctx.BlockTime()
	
	return k.CommunityFundProposals.Set(ctx, proposalID, proposal)
}

// AllocateDevelopmentFunds allocates funds from the development fund
func (k Keeper) AllocateDevelopmentFunds(ctx sdk.Context, proposalID uint64) error {
	proposal, err := k.DevelopmentFundProposals.Get(ctx, proposalID)
	if err != nil {
		return err
	}

	if proposal.Status != types.DevStatusApproved {
		return fmt.Errorf("proposal not approved")
	}

	balance, err := k.GetDevelopmentFundBalance(ctx)
	if err != nil {
		return err
	}

	// Update balance
	balance.AllocatedAmount = balance.AllocatedAmount.Add(proposal.RequestedAmount)
	balance.AvailableAmount = balance.AvailableAmount.Sub(proposal.RequestedAmount)
	balance.LastUpdateHeight = ctx.BlockHeight()
	balance.LastUpdateTime = ctx.BlockTime()

	if err := k.DevelopmentFundBalance.Set(ctx, balance); err != nil {
		return err
	}

	// Create transaction record
	tx := types.DevelopmentFundTransaction{
		TxID:        fmt.Sprintf("df-%d-%d", proposalID, ctx.BlockHeight()),
		ProposalID:  proposalID,
		From:        sdk.AccAddress(types.DevelopmentFundModuleName),
		Amount:      proposal.RequestedAmount,
		Type:        types.DevTxTypeAllocation,
		Category:    proposal.Category,
		Description: proposal.Description,
		Timestamp:   ctx.BlockTime(),
		BlockHeight: ctx.BlockHeight(),
		Status:      types.DevTxStatusApproved,
		Reviewed:    true,
		Approved:    true,
	}

	if err := k.DevelopmentFundTransactions.Set(ctx, tx.TxID, tx); err != nil {
		return err
	}

	// Update proposal status
	proposal.Status = types.DevStatusInProgress
	now := ctx.BlockTime()
	proposal.ExecutionStartTime = &now
	
	return k.DevelopmentFundProposals.Set(ctx, proposalID, proposal)
}