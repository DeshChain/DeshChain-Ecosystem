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

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc authtypes.ModuleAccountI)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
}

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	BlockedAddr(addr sdk.AccAddress) bool
}

// DistributionKeeper defines the expected distribution keeper
type DistributionKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
	GetFeePool(ctx sdk.Context) (feePool DistributionFeePool)
	SetFeePool(ctx sdk.Context, feePool DistributionFeePool)
}

// DistributionFeePool defines the expected fee pool interface
type DistributionFeePool interface {
	GetCommunityPool() sdk.DecCoins
}

// StakingKeeper defines the expected staking keeper (for validator queries)
type StakingKeeper interface {
	Validator(sdk.Context, sdk.ValAddress) (StakingValidatorI, error)
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) (StakingValidatorI, error)
	GetAllValidators(ctx sdk.Context) []StakingValidatorI
}

// StakingValidatorI defines the expected validator interface
type StakingValidatorI interface {
	GetOperator() sdk.ValAddress
	GetConsAddr() (sdk.ConsAddress, error)
	GetMoniker() string
	GetStatus() BondStatus
	IsBonded() bool
	IsUnbonded() bool
	IsUnbonding() bool
	GetTokens() sdk.Int
	GetBondedTokens() sdk.Int
	GetDelegatorShares() sdk.Dec
}

// BondStatus is the validator bond status
type BondStatus int32

const (
	Unbonded BondStatus = 1
	Unbonding BondStatus = 2
	Bonded BondStatus = 3
)

// ICS20Keeper defines the expected IBC transfer keeper
type ICS20Keeper interface {
	SendTransfer(
		ctx sdk.Context,
		sourcePort,
		sourceChannel string,
		token sdk.Coin,
		sender sdk.AccAddress,
		receiver string,
		timeoutHeight uint64,
		timeoutTimestamp uint64,
	) error
}

// RevenueKeeper defines the expected revenue keeper interface
type RevenueKeeper interface {
	CollectTradingFee(ctx sdk.Context, moduleName string, trader sdk.AccAddress, fee sdk.Coins, pair string) error
	CollectServiceFee(ctx sdk.Context, moduleName string, user sdk.AccAddress, fee sdk.Coins, service string) error
	RecordRevenue(ctx sdk.Context, source, category string, amount sdk.Coins, description string) error
}

// MoneyOrderHooks defines the hooks for money order module
type MoneyOrderHooks interface {
	// Pool lifecycle hooks
	AfterPoolCreated(ctx sdk.Context, poolId uint64)
	BeforePoolCreation(ctx sdk.Context, creator sdk.AccAddress, poolType string) error
	AfterPoolCreation(ctx sdk.Context, poolId uint64, poolType string, creator sdk.AccAddress) error
	
	// Trading hooks
	AfterSwap(ctx sdk.Context, poolId uint64, tokenIn, tokenOut sdk.Coin)
	AfterLiquidityAdded(ctx sdk.Context, poolId uint64, provider sdk.AccAddress, shares sdk.Int)
	AfterLiquidityRemoved(ctx sdk.Context, poolId uint64, provider sdk.AccAddress, shares sdk.Int)
	AfterTradingFeeCollection(ctx sdk.Context, poolId uint64, tradingFees sdk.Coins) error
	
	// Money Order hooks
	AfterMoneyOrderCreated(ctx sdk.Context, orderId string, sender, receiver sdk.AccAddress, amount sdk.Coin)
	AfterMoneyOrderCompleted(ctx sdk.Context, orderId string)
	
	// Village pool hooks
	AfterVillagePoolCreated(ctx sdk.Context, poolId uint64, villageName string)
	AfterVillageMemberJoined(ctx sdk.Context, poolId uint64, member sdk.AccAddress)
	
	// Pension integration hooks
	AfterSurakshaContribution(ctx sdk.Context, pensionAccountId string, contributor sdk.AccAddress, contribution sdk.Coin, villagePostalCode string) error
	AfterSurakshaMaturity(ctx sdk.Context, pensionAccountId string, beneficiary sdk.AccAddress, maturityAmount sdk.Coin) error
	
	// Agricultural lending integration hooks
	BeforeAgriLoanApproval(ctx sdk.Context, borrower sdk.AccAddress, loanAmount sdk.Coin, villagePostalCode string) error
	AfterAgriLoanDisbursement(ctx sdk.Context, loanId string, borrower sdk.AccAddress, loanAmount sdk.Coin, loanType string, duration uint32, villagePostalCode string) error
	AfterAgriLoanRepayment(ctx sdk.Context, loanId string, repaymentAmount sdk.Coin) error
	
	// Monthly distribution hook
	MonthlyRevenueDistribution(ctx sdk.Context) error
}

// ParamKeyTable returns the parameter key table
func ParamKeyTable() ParamKeyTable {
	return NewParamKeyTable().RegisterParamSet(&MoneyOrderParams{})
}

// ParamKeyTable is an interface for parameter key table
type ParamKeyTable interface {
	RegisterParamSet(ParamSet) ParamKeyTable
}

// NewParamKeyTable creates a new parameter key table
func NewParamKeyTable() ParamKeyTable {
	// This will be properly implemented when integrating with Cosmos SDK
	// For now, returning a placeholder
	return &paramKeyTable{}
}

type paramKeyTable struct{}

func (p *paramKeyTable) RegisterParamSet(ps ParamSet) ParamKeyTable {
	return p
}

// ParamSet defines an interface for parameter sets
type ParamSet interface {
	ParamSetPairs() []ParamSetPair
}

// ParamSetPair defines a parameter set pair
type ParamSetPair struct {
	Key   []byte
	Value interface{}
	ValidatorFn func(interface{}) error
}