package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Trade Finance module errors
var (
	ErrInvalidPartyType          = sdkerrors.Register(ModuleName, 1201, "invalid party type")
	ErrPartyNotFound             = sdkerrors.Register(ModuleName, 1202, "party not found")
	ErrPartyAlreadyExists        = sdkerrors.Register(ModuleName, 1203, "party already exists")
	ErrLCNotFound                = sdkerrors.Register(ModuleName, 1204, "letter of credit not found")
	ErrInvalidLCStatus           = sdkerrors.Register(ModuleName, 1205, "invalid LC status for operation")
	ErrInsufficientCollateral    = sdkerrors.Register(ModuleName, 1206, "insufficient collateral")
	ErrDocumentNotFound          = sdkerrors.Register(ModuleName, 1207, "document not found")
	ErrInvalidDocumentType       = sdkerrors.Register(ModuleName, 1208, "invalid document type")
	ErrDocumentAlreadyVerified   = sdkerrors.Register(ModuleName, 1209, "document already verified")
	ErrMissingRequiredDocuments  = sdkerrors.Register(ModuleName, 1210, "missing required documents")
	ErrPaymentInstructionNotFound = sdkerrors.Register(ModuleName, 1211, "payment instruction not found")
	ErrPaymentAlreadyCompleted   = sdkerrors.Register(ModuleName, 1212, "payment already completed")
	ErrInvalidPaymentAmount      = sdkerrors.Register(ModuleName, 1213, "invalid payment amount")
	ErrUnauthorized              = sdkerrors.Register(ModuleName, 1214, "unauthorized operation")
	ErrLCExpired                 = sdkerrors.Register(ModuleName, 1215, "letter of credit has expired")
	ErrShipmentDeadlinePassed    = sdkerrors.Register(ModuleName, 1216, "shipment deadline has passed")
	ErrInvalidInsurancePolicy    = sdkerrors.Register(ModuleName, 1217, "invalid insurance policy")
	ErrInsufficientInsurance     = sdkerrors.Register(ModuleName, 1218, "insufficient insurance coverage")
	ErrInvalidShipmentStatus     = sdkerrors.Register(ModuleName, 1219, "invalid shipment status")
	ErrAmountExceedsLC           = sdkerrors.Register(ModuleName, 1220, "amount exceeds LC limit")
	ErrInvalidCurrency           = sdkerrors.Register(ModuleName, 1221, "invalid or unsupported currency")
	ErrMinimumLCAmount           = sdkerrors.Register(ModuleName, 1222, "LC amount below minimum")
	ErrMaximumLCDuration         = sdkerrors.Register(ModuleName, 1223, "LC duration exceeds maximum")
	ErrInvalidIncoterms          = sdkerrors.Register(ModuleName, 1224, "invalid incoterms")
	ErrDisputePeriodActive       = sdkerrors.Register(ModuleName, 1225, "dispute resolution period active")
	ErrModuleDisabled            = sdkerrors.Register(ModuleName, 1226, "trade finance module is disabled")
)