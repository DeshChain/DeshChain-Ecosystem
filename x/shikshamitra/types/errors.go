package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/shikshamitra module sentinel errors
var (
	ErrInvalidDhanPata = sdkerrors.Register(ModuleName, 1100, "invalid DhanPata ID")
	ErrLoanNotFound = sdkerrors.Register(ModuleName, 1101, "loan not found")
	ErrNotEligible = sdkerrors.Register(ModuleName, 1102, "not eligible for loan")
	ErrUnauthorized = sdkerrors.Register(ModuleName, 1103, "unauthorized")
	ErrInvalidGPA = sdkerrors.Register(ModuleName, 1104, "invalid GPA")
	ErrInvalidRequest = sdkerrors.Register(ModuleName, 1105, "invalid request")
	ErrInvalidCoins = sdkerrors.Register(ModuleName, 1106, "invalid coins")
	ErrInvalidAddress = sdkerrors.Register(ModuleName, 1107, "invalid address")
	ErrStudentNotFound = sdkerrors.Register(ModuleName, 1108, "student not found")
	ErrInstitutionNotFound = sdkerrors.Register(ModuleName, 1109, "institution not found")
	ErrCourseNotFound = sdkerrors.Register(ModuleName, 1110, "course not found")
	ErrApplicationNotFound = sdkerrors.Register(ModuleName, 1111, "application not found")
	ErrScholarshipNotFound = sdkerrors.Register(ModuleName, 1112, "scholarship not found")
	ErrInvalidAcademicRecord = sdkerrors.Register(ModuleName, 1113, "invalid academic record")
	ErrInsufficientFunds = sdkerrors.Register(ModuleName, 1114, "insufficient funds")
	ErrExceedsLoanLimit = sdkerrors.Register(ModuleName, 1115, "exceeds loan limit")
	ErrInvalidRepayment = sdkerrors.Register(ModuleName, 1116, "invalid repayment")
	ErrLoanAlreadyActive = sdkerrors.Register(ModuleName, 1117, "loan already active")
	ErrInvalidEmploymentStatus = sdkerrors.Register(ModuleName, 1118, "invalid employment status")
	ErrMoratoriumActive = sdkerrors.Register(ModuleName, 1119, "moratorium period active")
	ErrRepaymentNotStarted = sdkerrors.Register(ModuleName, 1120, "repayment not yet started")
)