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

package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/revenue"
	"github.com/DeshChain/DeshChain-Ecosystem/x/revenue/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/revenue/types"
	taxtypes "github.com/DeshChain/DeshChain-Ecosystem/x/tax/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	keeper      keeper.Keeper
	bankKeeper  types.BankKeeper
	taxKeeper   types.TaxKeeper
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(types.StoreKey)
	tkey := sdk.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContext(key, tkey)
	suite.ctx = testCtx.Ctx

	encCfg := moduletestutil.MakeTestEncodingConfig(revenue.AppModuleBasic{})

	// Create mock keepers
	suite.bankKeeper = &MockBankKeeper{
		balances: make(map[string]sdk.Coins),
	}
	suite.taxKeeper = &MockTaxKeeper{}

	// Create keeper
	suite.keeper = keeper.NewKeeper(
		encCfg.Codec,
		sdk.NewKVStoreService(key),
		suite.bankKeeper,
		suite.taxKeeper,
		authtypes.NewModuleAddress(types.ModuleName).String(),
	)

	// Create message and query servers
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)

	queryHelper := baseapp.NewQueryServerTestHelper(testCtx.Ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, suite.keeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Test CollectRevenue
func (suite *KeeperTestSuite) TestCollectRevenue() {
	// Test source module
	sourceModule := "sikkebaaz"
	
	// Test fees
	platformFees := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000)))
	transactionVolume := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(100000)))

	// Collect revenue
	err := suite.keeper.CollectRevenue(suite.ctx, sourceModule, platformFees, transactionVolume)
	suite.Require().NoError(err)

	// Check that revenue was recorded
	moduleRevenue := suite.keeper.GetModuleRevenue(suite.ctx, sourceModule)
	suite.Equal(platformFees, moduleRevenue.PlatformFees)
	suite.Equal(transactionVolume, moduleRevenue.TransactionVolume)
}

// Test DistributeRevenue
func (suite *KeeperTestSuite) TestDistributeRevenue() {
	// Setup initial revenue in module account
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	initialRevenue := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000)))
	mockBank := suite.bankKeeper.(*MockBankKeeper)
	mockBank.balances[moduleAddr.String()] = initialRevenue

	// Enable revenue distribution
	params := types.DefaultParams()
	params.Enabled = true
	suite.keeper.SetParams(suite.ctx, params)

	// Distribute revenue
	err := suite.keeper.DistributeRevenue(suite.ctx)
	suite.Require().NoError(err)

	// Check distribution percentages
	// Liquidity: 30% = 300
	liquidityBalance := mockBank.GetBalance(suite.ctx, authtypes.NewModuleAddress(types.LiquidityPool), "namo")
	suite.Equal(sdk.NewInt(300), liquidityBalance.Amount)

	// Marketing: 20% = 200
	marketingBalance := mockBank.GetBalance(suite.ctx, authtypes.NewModuleAddress(types.MarketingPool), "namo")
	suite.Equal(sdk.NewInt(200), marketingBalance.Amount)

	// Operations: 15% = 150
	operationsBalance := mockBank.GetBalance(suite.ctx, authtypes.NewModuleAddress(types.OperationsPool), "namo")
	suite.Equal(sdk.NewInt(150), operationsBalance.Amount)
}

// Test ProcessGrossProfit
func (suite *KeeperTestSuite) TestProcessGrossProfit() {
	// Test gross profit
	grossProfit := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(10000)))

	// Process gross profit (30% should go to NGO)
	err := suite.keeper.ProcessGrossProfit(suite.ctx, grossProfit)
	suite.Require().NoError(err)

	// Check that 30% went to donation module
	mockBank := suite.bankKeeper.(*MockBankKeeper)
	donationBalance := mockBank.GetBalance(suite.ctx, authtypes.NewModuleAddress("donation"), "namo")
	expectedDonation := sdk.NewInt(3000) // 30% of 10000
	suite.Equal(expectedDonation, donationBalance.Amount)
}

// MockBankKeeper implements a mock bank keeper for testing
type MockBankKeeper struct {
	balances map[string]sdk.Coins
}

func (m *MockBankKeeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	senderAddr := authtypes.NewModuleAddress(senderModule)
	recipientAddr := authtypes.NewModuleAddress(recipientModule)
	
	// Deduct from sender
	senderBalance := m.balances[senderAddr.String()]
	newBalance, negative := senderBalance.SafeSub(amt...)
	if negative {
		return sdk.ErrInsufficientFunds
	}
	m.balances[senderAddr.String()] = newBalance
	
	// Add to recipient
	recipientBalance := m.balances[recipientAddr.String()]
	m.balances[recipientAddr.String()] = recipientBalance.Add(amt...)
	
	return nil
}

func (m *MockBankKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	recipientAddr := authtypes.NewModuleAddress(recipientModule)
	
	// Add to recipient
	recipientBalance := m.balances[recipientAddr.String()]
	m.balances[recipientAddr.String()] = recipientBalance.Add(amt...)
	
	return nil
}

func (m *MockBankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	balance := m.balances[addr.String()]
	return sdk.NewCoin(denom, balance.AmountOf(denom))
}

func (m *MockBankKeeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return m.balances[addr.String()]
}

// MockTaxKeeper implements a mock tax keeper for testing
type MockTaxKeeper struct{}

func (m *MockTaxKeeper) GetParams(ctx sdk.Context) taxtypes.Params {
	return taxtypes.DefaultParams()
}