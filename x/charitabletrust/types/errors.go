package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/charitabletrust module sentinel errors
var (
	ErrInvalidProposal            = sdkerrors.Register(ModuleName, 2, "invalid proposal")
	ErrInsufficientFunds          = sdkerrors.Register(ModuleName, 3, "insufficient funds in charitable trust")
	ErrAllocationNotFound         = sdkerrors.Register(ModuleName, 4, "allocation not found")
	ErrUnauthorized               = sdkerrors.Register(ModuleName, 5, "unauthorized")
	ErrInvalidAmount              = sdkerrors.Register(ModuleName, 6, "invalid amount")
	ErrOrganizationNotFound       = sdkerrors.Register(ModuleName, 7, "charitable organization not found")
	ErrOrganizationNotVerified    = sdkerrors.Register(ModuleName, 8, "charitable organization not verified")
	ErrOrganizationInactive       = sdkerrors.Register(ModuleName, 9, "charitable organization inactive")
	ErrReportNotFound             = sdkerrors.Register(ModuleName, 10, "impact report not found")
	ErrAlertNotFound              = sdkerrors.Register(ModuleName, 11, "fraud alert not found")
	ErrInvalidCategory            = sdkerrors.Register(ModuleName, 12, "invalid distribution category")
	ErrMonthlyLimitExceeded       = sdkerrors.Register(ModuleName, 13, "monthly allocation limit exceeded for organization")
	ErrQuorumNotReached           = sdkerrors.Register(ModuleName, 14, "quorum not reached")
	ErrVotingPeriodExpired        = sdkerrors.Register(ModuleName, 15, "voting period expired")
	ErrVotingPeriodActive         = sdkerrors.Register(ModuleName, 16, "voting period still active")
	ErrDuplicateVote              = sdkerrors.Register(ModuleName, 17, "duplicate vote")
	ErrNotTrustee                 = sdkerrors.Register(ModuleName, 18, "not a trustee")
	ErrInvestigationInProgress    = sdkerrors.Register(ModuleName, 19, "investigation in progress")
	ErrInvalidReportPeriod        = sdkerrors.Register(ModuleName, 20, "invalid report period")
	ErrMissingImpactMetrics       = sdkerrors.Register(ModuleName, 21, "missing required impact metrics")
	ErrVerificationFailed         = sdkerrors.Register(ModuleName, 22, "verification failed")
	ErrFraudDetected              = sdkerrors.Register(ModuleName, 23, "fraud detected")
	ErrMinAllocationAmount        = sdkerrors.Register(ModuleName, 24, "amount below minimum allocation")
)