package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
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

// StakingKeeper defines the expected staking keeper for NAMO staking
type StakingKeeper interface {
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, bool)
	GetAllDelegations(ctx sdk.Context) []stakingtypes.Delegation
	GetDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) []stakingtypes.Delegation
	BondDenom(ctx sdk.Context) string
}

// SurakshaKeeper defines the expected interface for Suraksha pool integration
type SurakshaKeeper interface {
	// Village Pool functions
	IsVillagePoolMember(ctx sdk.Context, member sdk.AccAddress) bool
	GetVillagePoolContribution(ctx sdk.Context, member sdk.AccAddress) sdk.Dec
	GetVillagePoolTotalValue(ctx sdk.Context) sdk.Dec
	
	// Urban Pool functions  
	IsUrbanPoolMember(ctx sdk.Context, member sdk.AccAddress) bool
	GetUrbanPoolContribution(ctx sdk.Context, member sdk.AccAddress) sdk.Dec
	GetUrbanPoolTotalValue(ctx sdk.Context) sdk.Dec
	
	// Pool management
	JoinVillagePool(ctx sdk.Context, member sdk.AccAddress, amount sdk.Dec) error
	JoinUrbanPool(ctx sdk.Context, member sdk.AccAddress, amount sdk.Dec) error
	LeavePool(ctx sdk.Context, member sdk.AccAddress, poolType string) error
}

// DEXKeeper defines the expected interface for Money Order DEX integration
type DEXKeeper interface {
	GetTotalLiquidity(ctx sdk.Context) sdk.Dec
	GetLiquidityPoolValue(ctx sdk.Context, pair string) sdk.Dec
	GetAllLiquidityPools(ctx sdk.Context) []LiquidityPool
}

// LiquidityPool represents a DEX liquidity pool
type LiquidityPool struct {
	Pair        string    `json:"pair"`
	Token0      string    `json:"token0"`
	Token1      string    `json:"token1"`
	Reserve0    sdk.Dec   `json:"reserve0"`
	Reserve1    sdk.Dec   `json:"reserve1"`
	TotalValue  sdk.Dec   `json:"total_value"`
	CreatedAt   time.Time `json:"created_at"`
}

// CollateralLoan represents a NAMO collateral loan
type CollateralLoan struct {
	Borrower         string    `json:"borrower"`
	LoanAmount       sdk.Dec   `json:"loan_amount"`
	CollateralAmount sdk.Dec   `json:"collateral_amount"`
	InterestRate     sdk.Dec   `json:"interest_rate"`
	Timestamp        time.Time `json:"timestamp"`
	DueDate          time.Time `json:"due_date"`
	IsActive         bool      `json:"is_active"`
	Module           string    `json:"module"`
}

// PoolMembership represents pool membership information
type PoolMembership struct {
	Member       string    `json:"member"`
	PoolType     string    `json:"pool_type"` // "village" or "urban"
	JoinedAt     time.Time `json:"joined_at"`
	Contribution sdk.Dec   `json:"contribution"`
	IsActive     bool      `json:"is_active"`
}

// NAMOStakeInfo represents NAMO staking information
type NAMOStakeInfo struct {
	Staker           string    `json:"staker"`
	StakedAmount     sdk.Dec   `json:"staked_amount"`
	LockedCollateral sdk.Dec   `json:"locked_collateral"`
	AvailableStake   sdk.Dec   `json:"available_stake"`
	StakedAt         time.Time `json:"staked_at"`
}

// LendingStats represents lending statistics
type LendingStats struct {
	TotalLoansIssued    sdk.Dec `json:"total_loans_issued"`
	TotalCollateralUsed sdk.Dec `json:"total_collateral_used"`
	ActiveLoans         int64   `json:"active_loans"`
	DefaultedLoans      int64   `json:"defaulted_loans"`
	TotalInterestEarned sdk.Dec `json:"total_interest_earned"`
	DailyVolume         sdk.Dec `json:"daily_volume"`
}