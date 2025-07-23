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

	"github.com/deshchain/namo/x/donation"
	"github.com/deshchain/namo/x/donation/keeper"
	"github.com/deshchain/namo/x/donation/types"
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

	encCfg := moduletestutil.MakeTestEncodingConfig(donation.AppModuleBasic{})

	// Create mock bank keeper
	suite.bankKeeper = &MockBankKeeper{
		balances: make(map[string]sdk.Coins),
	}

	// Create keeper
	suite.keeper = keeper.NewKeeper(
		encCfg.Codec,
		sdk.NewKVStoreService(key),
		suite.bankKeeper,
		nil, // AccountKeeper not needed for these tests
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

// Test RegisterNGO
func (suite *KeeperTestSuite) TestRegisterNGO() {
	// Create test NGO
	ngo := types.NGO{
		Id:          "ngo-001",
		Name:        "Test Foundation",
		Description: "A test NGO for helping communities",
		Website:     "https://testfoundation.org",
		WalletAddress: "desh1qypqxpq9qcrsszgse4wwrm07l3eyqq4mnz44tx",
		Category:    "education",
		IsActive:    true,
	}

	// Register NGO
	suite.keeper.SetNGO(suite.ctx, ngo)

	// Retrieve NGO
	retrievedNGO, found := suite.keeper.GetNGO(suite.ctx, ngo.Id)
	suite.Require().True(found)
	suite.Equal(ngo, retrievedNGO)

	// Test query
	resp, err := suite.queryClient.NGO(suite.ctx, &types.QueryNGORequest{Id: ngo.Id})
	suite.Require().NoError(err)
	suite.Equal(ngo, resp.Ngo)
}

// Test CreatePayout
func (suite *KeeperTestSuite) TestCreatePayout() {
	// Register test NGOs first
	ngo1 := types.NGO{
		Id:            "ngo-001",
		Name:          "Education Foundation",
		WalletAddress: "desh1qypqxpq9qcrsszgse4wwrm07l3eyqq4mnz44tx",
		IsActive:      true,
	}
	ngo2 := types.NGO{
		Id:            "ngo-002",
		Name:          "Healthcare Foundation",
		WalletAddress: "desh1qypqxpq9qcrsszgse4wwrm07l3eyqq4mnz44ty",
		IsActive:      true,
	}
	
	suite.keeper.SetNGO(suite.ctx, ngo1)
	suite.keeper.SetNGO(suite.ctx, ngo2)

	// Create payout
	payoutId := suite.keeper.CreatePayout(
		suite.ctx,
		sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000))),
		"Monthly NGO payout",
	)

	// Retrieve payout
	payout, found := suite.keeper.GetPayout(suite.ctx, payoutId)
	suite.Require().True(found)
	suite.Equal(uint64(1), payout.Id)
	suite.Equal(sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000))), payout.TotalAmount)
	suite.Equal("Monthly NGO payout", payout.Description)
	suite.Equal(types.PayoutStatus_PENDING, payout.Status)
}

// Test ProcessPayout
func (suite *KeeperTestSuite) TestProcessPayout() {
	// Setup module account with funds
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	mockBank := suite.bankKeeper.(*MockBankKeeper)
	mockBank.balances[moduleAddr.String()] = sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000)))

	// Register two active NGOs
	ngo1 := types.NGO{
		Id:            "ngo-001",
		Name:          "Education Foundation",
		WalletAddress: "desh1qypqxpq9qcrsszgse4wwrm07l3eyqq4mnz44tx",
		IsActive:      true,
	}
	ngo2 := types.NGO{
		Id:            "ngo-002",
		Name:          "Healthcare Foundation",
		WalletAddress: "desh1qypqxpq9qcrsszgse4wwrm07l3eyqq4mnz44ty",
		IsActive:      true,
	}
	
	suite.keeper.SetNGO(suite.ctx, ngo1)
	suite.keeper.SetNGO(suite.ctx, ngo2)

	// Process payout
	err := suite.keeper.ProcessPayout(suite.ctx, sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000))))
	suite.Require().NoError(err)

	// Check that each NGO received equal share (500 each)
	ngo1Addr, _ := sdk.AccAddressFromBech32(ngo1.WalletAddress)
	ngo1Balance := mockBank.GetBalance(suite.ctx, ngo1Addr, "namo")
	suite.Equal(sdk.NewInt(500), ngo1Balance.Amount)

	ngo2Addr, _ := sdk.AccAddressFromBech32(ngo2.WalletAddress)
	ngo2Balance := mockBank.GetBalance(suite.ctx, ngo2Addr, "namo")
	suite.Equal(sdk.NewInt(500), ngo2Balance.Amount)

	// Verify payout was recorded
	payouts := suite.keeper.GetAllPayouts(suite.ctx)
	suite.Require().Len(payouts, 1)
	suite.Equal(types.PayoutStatus_COMPLETED, payouts[0].Status)
}

// MockBankKeeper implements a mock bank keeper for testing
type MockBankKeeper struct {
	balances map[string]sdk.Coins
}

func (m *MockBankKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	senderAddr := authtypes.NewModuleAddress(senderModule)
	
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

func (m *MockBankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	balance := m.balances[addr.String()]
	return sdk.NewCoin(denom, balance.AmountOf(denom))
}

func (m *MockBankKeeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return m.balances[addr.String()]
}