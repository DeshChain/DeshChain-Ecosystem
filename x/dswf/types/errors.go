package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/dswf module sentinel errors
var (
	ErrInvalidProposal          = sdkerrors.Register(ModuleName, 2, "invalid proposal")
	ErrInsufficientFunds        = sdkerrors.Register(ModuleName, 3, "insufficient funds in DSWF")
	ErrAllocationNotFound       = sdkerrors.Register(ModuleName, 4, "allocation not found")
	ErrUnauthorized             = sdkerrors.Register(ModuleName, 5, "unauthorized")
	ErrInvalidAmount            = sdkerrors.Register(ModuleName, 6, "invalid amount")
	ErrAllocationLimitExceeded  = sdkerrors.Register(ModuleName, 7, "allocation limit exceeded")
	ErrInvalidCategory          = sdkerrors.Register(ModuleName, 8, "invalid allocation category")
	ErrDisbursementNotReady     = sdkerrors.Register(ModuleName, 9, "disbursement not ready")
	ErrMilestoneNotMet          = sdkerrors.Register(ModuleName, 10, "milestone not met")
	ErrInvalidStatus            = sdkerrors.Register(ModuleName, 11, "invalid status")
	ErrPortfolioRebalancing     = sdkerrors.Register(ModuleName, 12, "portfolio rebalancing in progress")
	ErrInvalidInvestmentStrategy = sdkerrors.Register(ModuleName, 13, "invalid investment strategy")
	ErrQuorumNotReached         = sdkerrors.Register(ModuleName, 14, "quorum not reached")
	ErrVotingPeriodExpired      = sdkerrors.Register(ModuleName, 15, "voting period expired")
	ErrDuplicateApproval        = sdkerrors.Register(ModuleName, 16, "duplicate approval")
	ErrReportNotFound           = sdkerrors.Register(ModuleName, 17, "report not found")
	ErrInvalidMetrics           = sdkerrors.Register(ModuleName, 18, "invalid performance metrics")
	ErrFundManagerNotFound      = sdkerrors.Register(ModuleName, 19, "fund manager not found")
	ErrMinimumBalanceRequired   = sdkerrors.Register(ModuleName, 20, "minimum fund balance required")
)