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
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/tax"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tax/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tax/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	keeper      keeper.Keeper
	bankKeeper  types.BankKeeper
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(types.StoreKey)
	tkey := sdk.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContext(key, tkey)
	suite.ctx = testCtx.Ctx

	encCfg := moduletestutil.MakeTestEncodingConfig(tax.AppModuleBasic{})

	// Create mock bank keeper
	suite.bankKeeper = MockBankKeeper{
		balances: make(map[string]sdk.Coins),
	}

	// Create keeper
	suite.keeper = keeper.NewKeeper(
		encCfg.Codec,
		sdk.NewKVStoreService(key),
		suite.bankKeeper,
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

// Test CollectTax
func (suite *KeeperTestSuite) TestCollectTax() {
	// Enable tax collection
	params := types.DefaultParams()
	params.Enabled = true
	suite.keeper.SetParams(suite.ctx, params)

	// Test address
	addr := sdk.AccAddress("test_address")
	
	// Test amount
	amount := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000)))

	// Collect tax
	err := suite.keeper.CollectTax(suite.ctx, addr, amount)
	suite.Require().NoError(err)

	// Check that tax was collected (2.5% of 1000 = 25)
	expectedTax := sdk.NewInt(25)
	
	// Verify tax was sent to module account
	mockBank := suite.bankKeeper.(MockBankKeeper)
	moduleBalance := mockBank.GetBalance(suite.ctx, authtypes.NewModuleAddress(types.ModuleName), "namo")
	suite.Equal(expectedTax, moduleBalance.Amount)
}

// Test DistributeTax
func (suite *KeeperTestSuite) TestDistributeTax() {
	// Setup initial tax in module account
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	initialTax := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000)))
	mockBank := suite.bankKeeper.(MockBankKeeper)
	mockBank.balances[moduleAddr.String()] = initialTax

	// Enable tax collection
	params := types.DefaultParams()
	params.Enabled = true
	suite.keeper.SetParams(suite.ctx, params)

	// Distribute tax
	err := suite.keeper.DistributeTax(suite.ctx)
	suite.Require().NoError(err)

	// Check distribution
	// Each pool should get 10% = 100 NAMO
	expectedPerPool := sdk.NewInt(100)

	// Verify NGO pool got 10%
	ngoBalance := mockBank.GetBalance(suite.ctx, authtypes.NewModuleAddress(types.NGOPool), "namo")
	suite.Equal(expectedPerPool, ngoBalance.Amount)

	// Verify Founder pool got 10%
	founderBalance := mockBank.GetBalance(suite.ctx, authtypes.NewModuleAddress(types.FounderPool), "namo")
	suite.Equal(expectedPerPool, founderBalance.Amount)
}

// MockBankKeeper implements a mock bank keeper for testing
type MockBankKeeper struct {
	balances map[string]sdk.Coins
}

func (m MockBankKeeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
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

func (m MockBankKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	recipientAddr := authtypes.NewModuleAddress(recipientModule)
	
	// Add to recipient (we don't track sender balance in this mock)
	recipientBalance := m.balances[recipientAddr.String()]
	m.balances[recipientAddr.String()] = recipientBalance.Add(amt...)
	
	return nil
}

func (m MockBankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	balance := m.balances[addr.String()]
	return sdk.NewCoin(denom, balance.AmountOf(denom))
}

func (m MockBankKeeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return m.balances[addr.String()]
}