package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// AccountKeeper defines the expected interface needed to retrieve account balances.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	HasDhanPataID(ctx sdk.Context, addr sdk.AccAddress, dhanPataID string) bool
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	IterateAllBalances(ctx sdk.Context, cb func(address sdk.AccAddress, coin sdk.Coin) bool)
	IterateAccountBalances(ctx sdk.Context, addr sdk.AccAddress, cb func(coin sdk.Coin) bool)
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
}

// DhanPataKeeper defines the expected interface for DhanPata virtual address system
type DhanPataKeeper interface {
	GetVirtualAddress(ctx sdk.Context, realAddress sdk.AccAddress) (string, bool)
	SetVirtualAddress(ctx sdk.Context, realAddress sdk.AccAddress, virtualAddress string) error
	ResolveVirtualAddress(ctx sdk.Context, virtualAddress string) (sdk.AccAddress, bool)
	IsValidVirtualAddress(virtualAddress string) bool
	GenerateVirtualAddress(ctx sdk.Context, realAddress sdk.AccAddress) (string, error)
}

// LiquidityManagerKeeper defines the expected interface for liquidity management and member verification
type LiquidityManagerKeeper interface {
	// Liquidity status functions
	GetLiquidityInfo(ctx sdk.Context) LiquidityInfo
	IsLendingAvailable(ctx sdk.Context) bool
	CanProcessLoan(ctx sdk.Context, amount sdk.Dec, module string, borrower sdk.AccAddress) (bool, string)
	
	// Member verification functions
	IsPoolMember(ctx sdk.Context, user sdk.AccAddress) bool
	
	// NAMO collateral functions
	CanProcessCollateralLoan(ctx sdk.Context, loanAmount sdk.Dec, collateralAmount sdk.Dec, borrower sdk.AccAddress) (bool, string)
	GetStakedNAMO(ctx sdk.Context, user sdk.AccAddress) sdk.Dec
	LockCollateral(ctx sdk.Context, user sdk.AccAddress, amount sdk.Dec) error
	UnlockCollateral(ctx sdk.Context, user sdk.AccAddress, amount sdk.Dec) error
	
	// Pool membership management
	SetPoolMembership(ctx sdk.Context, user sdk.AccAddress, poolType string, active bool)
	SetStakedNAMO(ctx sdk.Context, user sdk.AccAddress, amount sdk.Dec)
	
	// Loan processing
	ProcessLoan(ctx sdk.Context, borrower sdk.AccAddress, amount sdk.Dec, module string) error
	ProcessCollateralLoan(ctx sdk.Context, borrower sdk.AccAddress, loanAmount, collateralAmount sdk.Dec) error
	RepayCollateralLoan(ctx sdk.Context, borrower sdk.AccAddress, loanAmount, collateralAmount sdk.Dec) error
}

// LiquidityInfo represents comprehensive liquidity information
type LiquidityInfo struct {
	TotalPoolValue      sdk.Dec  `json:"total_pool_value"`
	AvailableForLending sdk.Dec  `json:"available_for_lending"`
	ReserveAmount       sdk.Dec  `json:"reserve_amount"`
	EmergencyReserve    sdk.Dec  `json:"emergency_reserve"`
	Status              string   `json:"status"`
	MaxLoanAmount       sdk.Dec  `json:"max_loan_amount"`
	DailyLendingLimit   sdk.Dec  `json:"daily_lending_limit"`
	AvailableModules    []string `json:"available_modules"`
	NextThreshold       sdk.Dec  `json:"next_threshold"`
	ProgressToNext      sdk.Dec  `json:"progress_to_next"`
	EstimatedDaysToNext int64    `json:"estimated_days_to_next"`
}