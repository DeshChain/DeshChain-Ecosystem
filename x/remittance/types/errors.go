package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// DONTCOVER

// x/remittance module sentinel errors
var (
	// General errors
	ErrInvalidSigner       = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidTransferID   = sdkerrors.Register(ModuleName, 1101, "invalid transfer ID")
	ErrTransferNotFound    = sdkerrors.Register(ModuleName, 1102, "remittance transfer not found")
	ErrTransferExists      = sdkerrors.Register(ModuleName, 1103, "remittance transfer already exists")
	ErrInvalidAmount       = sdkerrors.Register(ModuleName, 1104, "invalid transfer amount")
	ErrInsufficientFunds   = sdkerrors.Register(ModuleName, 1105, "insufficient funds for transfer")
	
	// Transfer status errors
	ErrInvalidTransferStatus = sdkerrors.Register(ModuleName, 1110, "invalid transfer status")
	ErrTransferAlreadyProcessed = sdkerrors.Register(ModuleName, 1111, "transfer already processed")
	ErrTransferExpired       = sdkerrors.Register(ModuleName, 1112, "transfer has expired")
	ErrTransferCancelled     = sdkerrors.Register(ModuleName, 1113, "transfer has been cancelled")
	ErrTransferNotPending    = sdkerrors.Register(ModuleName, 1114, "transfer is not in pending status")
	ErrTransferNotProcessing = sdkerrors.Register(ModuleName, 1115, "transfer is not in processing status")
	ErrCannotCancelTransfer  = sdkerrors.Register(ModuleName, 1116, "cannot cancel transfer in current status")
	
	// Currency and exchange rate errors
	ErrInvalidCurrency       = sdkerrors.Register(ModuleName, 1120, "invalid currency code")
	ErrCurrencyNotSupported  = sdkerrors.Register(ModuleName, 1121, "currency not supported")
	ErrExchangeRateNotFound  = sdkerrors.Register(ModuleName, 1122, "exchange rate not found")
	ErrInvalidExchangeRate   = sdkerrors.Register(ModuleName, 1123, "invalid exchange rate")
	ErrExchangeRateExpired   = sdkerrors.Register(ModuleName, 1124, "exchange rate has expired")
	ErrExchangeRateStale     = sdkerrors.Register(ModuleName, 1125, "exchange rate is stale")
	
	// Country and corridor errors
	ErrInvalidCountryCode    = sdkerrors.Register(ModuleName, 1130, "invalid country code")
	ErrCountryNotSupported   = sdkerrors.Register(ModuleName, 1131, "country not supported")
	ErrCorridorNotFound      = sdkerrors.Register(ModuleName, 1132, "remittance corridor not found")
	ErrCorridorInactive      = sdkerrors.Register(ModuleName, 1133, "remittance corridor is inactive")
	ErrCorridorExists        = sdkerrors.Register(ModuleName, 1134, "remittance corridor already exists")
	ErrInvalidCorridorID     = sdkerrors.Register(ModuleName, 1135, "invalid corridor ID")
	
	// Liquidity pool errors
	ErrPoolNotFound          = sdkerrors.Register(ModuleName, 1140, "liquidity pool not found")
	ErrPoolExists            = sdkerrors.Register(ModuleName, 1141, "liquidity pool already exists")
	ErrInvalidPoolID         = sdkerrors.Register(ModuleName, 1142, "invalid pool ID")
	ErrInsufficientLiquidity = sdkerrors.Register(ModuleName, 1143, "insufficient liquidity in pool")
	ErrInvalidLiquidityAmount = sdkerrors.Register(ModuleName, 1144, "invalid liquidity amount")
	ErrPoolInactive          = sdkerrors.Register(ModuleName, 1145, "liquidity pool is inactive")
	ErrMinLiquidityNotMet    = sdkerrors.Register(ModuleName, 1146, "minimum liquidity requirement not met")
	ErrLiquidityProviderNotFound = sdkerrors.Register(ModuleName, 1147, "liquidity provider not found")
	ErrInvalidLPTokens       = sdkerrors.Register(ModuleName, 1148, "invalid LP token amount")
	
	// Settlement and partner errors
	ErrPartnerNotFound       = sdkerrors.Register(ModuleName, 1150, "settlement partner not found")
	ErrPartnerInactive       = sdkerrors.Register(ModuleName, 1151, "settlement partner is inactive")
	ErrInvalidPartnerID      = sdkerrors.Register(ModuleName, 1152, "invalid partner ID")
	ErrPartnerExists         = sdkerrors.Register(ModuleName, 1153, "settlement partner already exists")
	ErrInvalidSettlementMethod = sdkerrors.Register(ModuleName, 1154, "invalid settlement method")
	ErrSettlementNotFound    = sdkerrors.Register(ModuleName, 1155, "settlement record not found")
	ErrSettlementFailed      = sdkerrors.Register(ModuleName, 1156, "settlement processing failed")
	ErrInvalidSettlementDetails = sdkerrors.Register(ModuleName, 1157, "invalid settlement details")
	
	// KYC and compliance errors
	ErrKYCRequired           = sdkerrors.Register(ModuleName, 1160, "KYC verification required")
	ErrInsufficientKYCLevel  = sdkerrors.Register(ModuleName, 1161, "insufficient KYC level for this transaction")
	ErrComplianceCheckFailed = sdkerrors.Register(ModuleName, 1162, "compliance check failed")
	ErrSanctionsScreeningFailed = sdkerrors.Register(ModuleName, 1163, "sanctions screening failed")
	ErrPEPScreeningFailed    = sdkerrors.Register(ModuleName, 1164, "PEP screening failed")
	ErrHighRiskTransaction   = sdkerrors.Register(ModuleName, 1165, "transaction flagged as high risk")
	ErrManualReviewRequired  = sdkerrors.Register(ModuleName, 1166, "manual review required for this transaction")
	ErrComplianceHold        = sdkerrors.Register(ModuleName, 1167, "transaction on compliance hold")
	
	// Transaction limit errors
	ErrAmountBelowMinimum    = sdkerrors.Register(ModuleName, 1170, "amount below minimum transfer limit")
	ErrAmountAboveMaximum    = sdkerrors.Register(ModuleName, 1171, "amount above maximum transfer limit")
	ErrDailyLimitExceeded    = sdkerrors.Register(ModuleName, 1172, "daily transfer limit exceeded")
	ErrMonthlyLimitExceeded  = sdkerrors.Register(ModuleName, 1173, "monthly transfer limit exceeded")
	ErrTransactionLimitExceeded = sdkerrors.Register(ModuleName, 1174, "transaction count limit exceeded")
	
	// Recipient and address errors
	ErrInvalidRecipient      = sdkerrors.Register(ModuleName, 1180, "invalid recipient information")
	ErrInvalidAddress        = sdkerrors.Register(ModuleName, 1181, "invalid address")
	ErrInvalidPhoneNumber    = sdkerrors.Register(ModuleName, 1182, "invalid phone number")
	ErrInvalidEmail          = sdkerrors.Register(ModuleName, 1183, "invalid email address")
	ErrInvalidBankDetails    = sdkerrors.Register(ModuleName, 1184, "invalid bank account details")
	ErrInvalidMobileWallet   = sdkerrors.Register(ModuleName, 1185, "invalid mobile wallet information")
	ErrInvalidPickupLocation = sdkerrors.Register(ModuleName, 1186, "invalid cash pickup location")
	
	// Fee and pricing errors
	ErrInvalidFee            = sdkerrors.Register(ModuleName, 1190, "invalid fee structure")
	ErrFeeCalculationFailed  = sdkerrors.Register(ModuleName, 1191, "fee calculation failed")
	ErrInvalidPricing        = sdkerrors.Register(ModuleName, 1192, "invalid pricing configuration")
	ErrSlippageExceeded      = sdkerrors.Register(ModuleName, 1193, "slippage tolerance exceeded")
	
	// Authorization and permission errors
	ErrUnauthorized          = sdkerrors.Register(ModuleName, 1200, "unauthorized operation")
	ErrNotOwner              = sdkerrors.Register(ModuleName, 1201, "not the owner of this transfer")
	ErrNotRecipient          = sdkerrors.Register(ModuleName, 1202, "not the recipient of this transfer")
	ErrInvalidAuthority      = sdkerrors.Register(ModuleName, 1203, "invalid authority address")
	ErrPermissionDenied      = sdkerrors.Register(ModuleName, 1204, "permission denied")
	
	// IBC and cross-chain errors
	ErrIBCChannelNotFound    = sdkerrors.Register(ModuleName, 1210, "IBC channel not found")
	ErrIBCChannelInactive    = sdkerrors.Register(ModuleName, 1211, "IBC channel is inactive")
	ErrInvalidIBCPacket      = sdkerrors.Register(ModuleName, 1212, "invalid IBC packet")
	ErrCrossChainTransferFailed = sdkerrors.Register(ModuleName, 1213, "cross-chain transfer failed")
	ErrBridgeNotAvailable    = sdkerrors.Register(ModuleName, 1214, "cross-chain bridge not available")
	ErrRelayerNotFound       = sdkerrors.Register(ModuleName, 1215, "IBC relayer not found")
	
	// Configuration and parameter errors
	ErrInvalidParams         = sdkerrors.Register(ModuleName, 1220, "invalid module parameters")
	ErrConfigurationError    = sdkerrors.Register(ModuleName, 1221, "configuration error")
	ErrServiceUnavailable    = sdkerrors.Register(ModuleName, 1222, "remittance service unavailable")
	ErrMaintenanceMode       = sdkerrors.Register(ModuleName, 1223, "system in maintenance mode")
	
	// Data validation errors
	ErrInvalidInput          = sdkerrors.Register(ModuleName, 1230, "invalid input data")
	ErrDataCorruption        = sdkerrors.Register(ModuleName, 1231, "data corruption detected")
	ErrInvalidTimestamp      = sdkerrors.Register(ModuleName, 1232, "invalid timestamp")
	ErrInvalidSignature      = sdkerrors.Register(ModuleName, 1233, "invalid signature")
	ErrInvalidProof          = sdkerrors.Register(ModuleName, 1234, "invalid cryptographic proof")
	
	// Rate limiting and throttling errors
	ErrRateLimitExceeded     = sdkerrors.Register(ModuleName, 1240, "rate limit exceeded")
	ErrTooManyRequests       = sdkerrors.Register(ModuleName, 1241, "too many requests")
	ErrTemporarilyUnavailable = sdkerrors.Register(ModuleName, 1242, "service temporarily unavailable")
	
	// Business logic errors
	ErrBusinessHoursOnly     = sdkerrors.Register(ModuleName, 1250, "transfers only allowed during business hours")
	ErrHolidayRestriction    = sdkerrors.Register(ModuleName, 1251, "transfers restricted on holidays")
	ErrWeekendRestriction    = sdkerrors.Register(ModuleName, 1252, "transfers restricted on weekends")
	ErrCountryRestriction    = sdkerrors.Register(ModuleName, 1253, "transfers restricted for this country pair")
	ErrPurposeRestriction    = sdkerrors.Register(ModuleName, 1254, "invalid purpose code for this corridor")
	
	// External service errors
	ErrExternalServiceFailure = sdkerrors.Register(ModuleName, 1260, "external service failure")
	ErrBankAPIFailure        = sdkerrors.Register(ModuleName, 1261, "bank API failure")
	ErrMobileWalletAPIFailure = sdkerrors.Register(ModuleName, 1262, "mobile wallet API failure")
	ErrOracleServiceFailure  = sdkerrors.Register(ModuleName, 1263, "oracle service failure")
	ErrComplianceAPIFailure  = sdkerrors.Register(ModuleName, 1264, "compliance API failure")
	
	// Timeout and expiration errors
	ErrTransactionTimeout    = sdkerrors.Register(ModuleName, 1270, "transaction timeout")
	ErrSettlementTimeout     = sdkerrors.Register(ModuleName, 1271, "settlement timeout")
	ErrConfirmationTimeout   = sdkerrors.Register(ModuleName, 1272, "confirmation timeout")
	ErrIBCPacketTimeout      = sdkerrors.Register(ModuleName, 1273, "IBC packet timeout")
	
	// Reconciliation errors
	ErrReconciliationFailed  = sdkerrors.Register(ModuleName, 1280, "reconciliation failed")
	ErrBalanceMismatch       = sdkerrors.Register(ModuleName, 1281, "balance mismatch detected")
	ErrDuplicateTransaction  = sdkerrors.Register(ModuleName, 1282, "duplicate transaction detected")
	ErrInconsistentState     = sdkerrors.Register(ModuleName, 1283, "inconsistent state detected")
	
	// Sewa Mitra specific errors
	ErrInvalidAgentID        = sdkerrors.Register(ModuleName, 1300, "invalid Sewa Mitra agent ID")
	ErrAgentNotFound         = sdkerrors.Register(ModuleName, 1301, "Sewa Mitra agent not found")
	ErrAgentAlreadyExists    = sdkerrors.Register(ModuleName, 1302, "Sewa Mitra agent already exists")
	ErrInvalidAgentName      = sdkerrors.Register(ModuleName, 1303, "invalid agent name")
	ErrInvalidLocation       = sdkerrors.Register(ModuleName, 1304, "invalid agent location")
	ErrInvalidPhone          = sdkerrors.Register(ModuleName, 1305, "invalid phone number")
	ErrInvalidCommissionRate = sdkerrors.Register(ModuleName, 1306, "invalid commission rate")
	ErrInvalidBonusRate      = sdkerrors.Register(ModuleName, 1307, "invalid volume bonus rate")
	ErrInvalidLiquidity      = sdkerrors.Register(ModuleName, 1308, "invalid liquidity limit")
	ErrInvalidDailyLimit     = sdkerrors.Register(ModuleName, 1309, "invalid daily limit")
	ErrNoSupportedCurrencies = sdkerrors.Register(ModuleName, 1310, "no supported currencies specified")
	ErrNoSupportedMethods    = sdkerrors.Register(ModuleName, 1311, "no supported settlement methods specified")
	ErrInvalidAgentStatus    = sdkerrors.Register(ModuleName, 1312, "invalid agent status")
	ErrNoAvailableAgent      = sdkerrors.Register(ModuleName, 1313, "no available Sewa Mitra agent found")
	ErrAgentSuspended        = sdkerrors.Register(ModuleName, 1314, "Sewa Mitra agent is suspended")
	ErrAgentDeactivated      = sdkerrors.Register(ModuleName, 1315, "Sewa Mitra agent is deactivated")
	ErrAgentLiquidityExceeded = sdkerrors.Register(ModuleName, 1316, "agent liquidity limit exceeded")
	ErrAgentDailyLimitExceeded = sdkerrors.Register(ModuleName, 1317, "agent daily limit exceeded")
	
	// Sewa Mitra commission errors
	ErrCommissionNotFound    = sdkerrors.Register(ModuleName, 1320, "commission record not found")
	ErrCommissionAlreadyPaid = sdkerrors.Register(ModuleName, 1321, "commission already paid")
	ErrInvalidCommissionID   = sdkerrors.Register(ModuleName, 1322, "invalid commission ID")
	ErrCommissionCalculationFailed = sdkerrors.Register(ModuleName, 1323, "commission calculation failed")
	ErrCommissionPaymentFailed = sdkerrors.Register(ModuleName, 1324, "commission payment failed")
	
	// Sewa Mitra transfer specific errors
	ErrInvalidTransferType   = sdkerrors.Register(ModuleName, 1330, "invalid transfer type")
	ErrNotSewaMitraTransfer  = sdkerrors.Register(ModuleName, 1331, "transfer does not use Sewa Mitra")
	ErrAgentMismatch         = sdkerrors.Register(ModuleName, 1332, "agent mismatch for this transfer")
	ErrPickupLocationMismatch = sdkerrors.Register(ModuleName, 1333, "pickup location mismatch")
	ErrInvalidPickupProof    = sdkerrors.Register(ModuleName, 1334, "invalid pickup proof")
)