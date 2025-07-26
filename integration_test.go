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

package app_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/cometbft/cometbft-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/DeshChain/DeshChain-Ecosystem/app"
	taxtypes "github.com/DeshChain/DeshChain-Ecosystem/x/tax/types"
	revenuetypes "github.com/DeshChain/DeshChain-Ecosystem/x/revenue/types"
	donationtypes "github.com/DeshChain/DeshChain-Ecosystem/x/donation/types"
	sikkebaaztypes "github.com/DeshChain/DeshChain-Ecosystem/x/sikkebaaz/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	app       *app.DeshChainApp
	ctx       sdk.Context
	val       sdk.ValAddress
	accAddrs  []sdk.AccAddress
}

func (suite *IntegrationTestSuite) SetupTest() {
	// Create app with in-memory database
	db := dbm.NewMemDB()
	suite.app = app.New(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		app.DefaultNodeHome,
		5,
		app.MakeEncodingConfig(),
		simtestutil.EmptyAppOptions{},
		baseapp.SetMinGasPrices("0namo"),
	)

	// Initialize the chain
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		Height: 1,
		Time:   time.Now().UTC(),
	})

	// Create test accounts
	suite.accAddrs = make([]sdk.AccAddress, 3)
	for i := range suite.accAddrs {
		pk := secp256k1.GenPrivKey().PubKey()
		suite.accAddrs[i] = sdk.AccAddress(pk.Address())
	}

	// Fund test accounts
	initCoins := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000000000)))
	for _, addr := range suite.accAddrs {
		acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr)
		suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
		suite.Require().NoError(
			simtestutil.FundAccount(suite.app.BankKeeper, suite.ctx, addr, initCoins),
		)
	}

	// Register test NGOs
	ngo1 := donationtypes.NGO{
		Id:            "ngo-001",
		Name:          "Education Foundation",
		Description:   "Supporting education initiatives",
		Website:       "https://education.org",
		WalletAddress: suite.accAddrs[1].String(),
		Category:      "education",
		IsActive:      true,
	}
	ngo2 := donationtypes.NGO{
		Id:            "ngo-002",
		Name:          "Healthcare Foundation",
		Description:   "Providing healthcare support",
		Website:       "https://healthcare.org",
		WalletAddress: suite.accAddrs[2].String(),
		Category:      "healthcare",
		IsActive:      true,
	}
	
	suite.app.DonationKeeper.SetNGO(suite.ctx, ngo1)
	suite.app.DonationKeeper.SetNGO(suite.ctx, ngo2)

	// Enable tax collection
	taxParams := taxtypes.DefaultParams()
	taxParams.Enabled = true
	suite.app.TaxKeeper.SetParams(suite.ctx, taxParams)

	// Enable revenue distribution
	revenueParams := revenuetypes.DefaultParams()
	revenueParams.Enabled = true
	suite.app.RevenueKeeper.SetParams(suite.ctx, revenueParams)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// TestCompleteRevenueFlow tests the complete flow from transaction to NGO donation
func (suite *IntegrationTestSuite) TestCompleteRevenueFlow() {
	// Step 1: Simulate a transaction that triggers tax collection
	sender := suite.accAddrs[0]
	recipient := sdk.AccAddress("recipient_address")
	amount := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(100000))) // 100,000 NAMO

	// Create bank send message
	msg := banktypes.NewMsgSend(sender, recipient, amount)
	
	// Tax should be collected automatically via ante handler
	// 2.5% of 100,000 = 2,500 NAMO tax

	// Step 2: Simulate platform revenue from Sikkebaaz
	platformFees := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(10000))) // 10,000 NAMO
	transactionVolume := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000000))) // 1M NAMO

	err := suite.app.RevenueKeeper.CollectRevenue(
		suite.ctx,
		sikkebaaztypes.ModuleName,
		platformFees,
		transactionVolume,
	)
	suite.Require().NoError(err)

	// Step 3: Distribute tax (happens in EndBlocker)
	err = suite.app.TaxKeeper.DistributeTax(suite.ctx)
	suite.Require().NoError(err)

	// Step 4: Distribute revenue
	err = suite.app.RevenueKeeper.DistributeRevenue(suite.ctx)
	suite.Require().NoError(err)

	// Step 5: Process gross profit (30% to NGO)
	grossProfit := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(5000))) // 5,000 NAMO profit
	err = suite.app.RevenueKeeper.ProcessGrossProfit(suite.ctx, grossProfit)
	suite.Require().NoError(err)

	// Step 6: Check NGO balances
	// Each NGO should have received donations from:
	// 1. Tax distribution (10% of 2,500 = 250 NAMO to NGO pool, split between 2 NGOs = 125 each)
	// 2. Gross profit (30% of 5,000 = 1,500 NAMO, split between 2 NGOs = 750 each)
	// Total per NGO = 125 + 750 = 875 NAMO

	ngo1Addr, _ := sdk.AccAddressFromBech32(suite.accAddrs[1].String())
	ngo1Balance := suite.app.BankKeeper.GetBalance(suite.ctx, ngo1Addr, "namo")
	
	// NGO1 started with 1,000,000,000 and should have received 875
	expectedNGO1Balance := sdk.NewInt(1000000875)
	suite.Equal(expectedNGO1Balance, ngo1Balance.Amount, "NGO1 balance mismatch")

	// Check donation records
	payouts := suite.app.DonationKeeper.GetAllPayouts(suite.ctx)
	suite.Greater(len(payouts), 0, "No payouts recorded")

	// Check revenue metrics
	totalRevenue := suite.app.RevenueKeeper.GetTotalRevenue(suite.ctx)
	suite.True(totalRevenue.IsAllPositive(), "Total revenue should be positive")

	moduleRevenue := suite.app.RevenueKeeper.GetModuleRevenue(suite.ctx, sikkebaaztypes.ModuleName)
	suite.Equal(platformFees, moduleRevenue.PlatformFees)
	suite.Equal(transactionVolume, moduleRevenue.TransactionVolume)
}

// TestSikkebaazTokenLaunchWithFees tests token launch with proper fee collection
func (suite *IntegrationTestSuite) TestSikkebaazTokenLaunchWithFees() {
	creator := suite.accAddrs[0]
	
	// Create token launch
	launch := &sikkebaaztypes.MsgCreateTokenLaunch{
		Creator:           creator.String(),
		TokenName:         "TestCoin",
		TokenSymbol:       "TEST",
		TokenSupply:       sdk.NewInt(1000000000), // 1B tokens
		InitialPrice:      sdk.NewCoin("namo", sdk.NewInt(100)), // 0.0001 NAMO per token
		LaunchDuration:    3600, // 1 hour
		VestingDuration:   86400, // 1 day
		MinPurchase:       sdk.NewCoin("namo", sdk.NewInt(1000)),
		MaxPurchase:       sdk.NewCoin("namo", sdk.NewInt(1000000)),
		AntiPumpConfig: &sikkebaaztypes.AntiPumpConfig{
			MaxPriceIncrease:    200, // 2x
			VestingPeriod:       86400,
			SellLimitPercent:    20,
			CooldownPeriod:      3600,
			MinLiquidityPercent: 50,
		},
	}

	// Execute token launch
	msgServer := sikkebaaztypes.NewMsgServerImpl(suite.app.SikkebaazKeeper)
	resp, err := msgServer.CreateTokenLaunch(suite.ctx, launch)
	suite.Require().NoError(err)
	suite.NotEmpty(resp.LaunchId)

	// Check that platform fee was collected
	moduleRevenue := suite.app.RevenueKeeper.GetModuleRevenue(suite.ctx, sikkebaaztypes.ModuleName)
	suite.True(moduleRevenue.PlatformFees.IsAllPositive(), "Platform fees should be collected")
}

// TestGovernancePhaseTransition tests governance phase transitions
func (suite *IntegrationTestSuite) TestGovernancePhaseTransition() {
	// Set genesis time to 2 years ago to trigger phase transition
	genesisTime := time.Now().Add(-2 * 365 * 24 * time.Hour)
	suite.app.GovernanceKeeper.SetGenesisTime(suite.ctx, genesisTime)
	
	// Set founder address
	founderAddr := suite.accAddrs[0].String()
	suite.app.GovernanceKeeper.SetFounderAddress(suite.ctx, founderAddr)

	// Process phase transitions
	suite.app.GovernanceKeeper.ProcessPhaseTransitions(suite.ctx)

	// Check that we're now in shared governance phase
	currentPhase := suite.app.GovernanceKeeper.GetGovernancePhase(suite.ctx)
	suite.Equal(governancetypes.GovernancePhase_SHARED_GOVERNANCE, currentPhase)
}