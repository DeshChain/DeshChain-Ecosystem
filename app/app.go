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

package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	upgradeclient "cosmossdk.io/x/upgrade/client"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	abci "github.com/tendermint/tendermint/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cometbft/cometbft/libs/log"
	tmos "github.com/cometbft/cometbft/libs/os"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/spf13/cast"

	// Money Order imports
	moneyorder "github.com/deshchain/deshchain/x/moneyorder"
	moneyorderkeeper "github.com/deshchain/deshchain/x/moneyorder/keeper"
	moneyordertypes "github.com/deshchain/deshchain/x/moneyorder/types"

	// Cultural imports
	cultural "github.com/deshchain/deshchain/x/cultural"
	culturalkeeper "github.com/deshchain/deshchain/x/cultural/keeper"
	culturaltypes "github.com/deshchain/deshchain/x/cultural/types"

	// NAMO imports
	namo "github.com/deshchain/deshchain/x/namo"
	namokeeper "github.com/deshchain/deshchain/x/namo/keeper"
	namotypes "github.com/deshchain/deshchain/x/namo/types"

	// DhanSetu imports
	dhansetu "github.com/deshchain/deshchain/x/dhansetu"
	dhansetukeeper "github.com/deshchain/deshchain/x/dhansetu/keeper"
	dhansettypes "github.com/deshchain/deshchain/x/dhansetu/types"

	// DINR imports
	dinr "github.com/deshchain/deshchain/x/dinr"
	dinrkeeper "github.com/deshchain/deshchain/x/dinr/keeper"
	dinrtypes "github.com/deshchain/deshchain/x/dinr/types"

	// Trade Finance imports
	tradefinance "github.com/deshchain/deshchain/x/tradefinance"
	tradefinancekeeper "github.com/deshchain/deshchain/x/tradefinance/keeper"
	tradefinancetypes "github.com/deshchain/deshchain/x/tradefinance/types"

	// Oracle imports
	oracle "github.com/deshchain/deshchain/x/oracle"
	oraclekeeper "github.com/deshchain/deshchain/x/oracle/keeper"
	oracletypes "github.com/deshchain/deshchain/x/oracle/types"

	// Sikkebaaz imports
	sikkebaaz "github.com/deshchain/namo/x/sikkebaaz"
	sikkebaazkeeper "github.com/deshchain/namo/x/sikkebaaz/keeper"
	sikkebaaztypes "github.com/deshchain/namo/x/sikkebaaz/types"

	// Lending Suite imports
	krishimitra "github.com/deshchain/namo/x/krishimitra"
	krishimitrakeeper "github.com/deshchain/namo/x/krishimitra/keeper"
	krishimitratypes "github.com/deshchain/namo/x/krishimitra/types"

	vyavasayamitra "github.com/deshchain/namo/x/vyavasayamitra"
	vyavasayamitrakeeper "github.com/deshchain/namo/x/vyavasayamitra/keeper"
	vyavasayamitratypes "github.com/deshchain/namo/x/vyavasayamitra/types"

	shikshamitra "github.com/deshchain/namo/x/shikshamitra"
	shikshamitrakeeper "github.com/deshchain/namo/x/shikshamitra/keeper"
	shikshamitratypes "github.com/deshchain/namo/x/shikshamitra/types"
)

const (
	AccountAddressPrefix = "desh"
	Name                 = "deshchain"
)

// These constants are derived from the above variables.
// These are the ones we will want to use in the code.
var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager that is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
				upgradeclient.LegacyProposalHandler,
				upgradeclient.LegacyCancelProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		vesting.AppModuleBasic{},
		consensus.AppModuleBasic{},
		moneyorder.AppModuleBasic{},
		cultural.AppModuleBasic{},
		namo.AppModuleBasic{},
		dhansetu.AppModuleBasic{},
		dinr.AppModuleBasic{},
		tradefinance.AppModuleBasic{},
		oracle.AppModuleBasic{},
		sikkebaaz.AppModuleBasic{},
		krishimitra.AppModuleBasic{},
		vyavasayamitra.AppModuleBasic{},
		shikshamitra.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		"army_ngo":                     nil, // Army NGO wallet
		"war_relief":                   nil, // War Relief wallet
		"development_fund":             nil, // Development fund wallet
		"operations_fund":              nil, // Operations fund wallet
		"burn_address":                 {authtypes.Burner}, // Token burn address
		moneyordertypes.ModuleName:     nil, // Money Order module
		culturaltypes.ModuleName:       nil, // Cultural module
		namotypes.ModuleName:           nil, // NAMO token module
		namotypes.VestingPoolName:      nil, // NAMO vesting pool
		namotypes.BurnPoolName:         {authtypes.Burner}, // NAMO burn pool
		dhansettypes.ModuleName:        nil, // DhanSetu integration module
		dinrtypes.ModuleName:           {authtypes.Minter, authtypes.Burner}, // DINR stablecoin module
		tradefinancetypes.ModuleName:   nil, // Trade Finance module
		oracletypes.ModuleName:         nil, // Oracle price feeds module
		sikkebaaztypes.ModuleName:      nil, // Sikkebaaz launchpad module
		// Sikkebaaz Module Accounts
		sikkebaaztypes.SikkebaazFeeCollector:    nil, // Fee collection
		sikkebaaztypes.LaunchEscrowAccount:      nil, // Launch escrow
		sikkebaaztypes.LiquidityLockAccount:     nil, // Liquidity locks
		sikkebaaztypes.CreatorRewardsPool:       nil, // Creator rewards
		sikkebaaztypes.CommunityIncentivePool:   nil, // Community incentives
		sikkebaaztypes.LocalNGOPool:             nil, // Local NGO donations
		sikkebaaztypes.FestivalBonusPool:        nil, // Festival bonuses
		sikkebaaztypes.SecurityAuditFund:        nil, // Security audits
		sikkebaaztypes.EmergencyFund:            nil, // Emergency controls
		// Platform Revenue Distribution Accounts
		"platform_development_fund":    nil, // 30% - Platform development
		"platform_community_treasury":  nil, // 25% - Community programs
		"platform_liquidity_pool":      nil, // 20% - Market liquidity
		"ngo_donation_pool":            nil, // 10% - Social impact
		"platform_emergency_reserve":   nil, // 10% - Risk management
		"founder_royalty_pool":         nil, // 5% - Founder compensation
		// Transaction Tax Distribution Accounts (for future use)
		"validator_pool":               nil, // Validator rewards
		"community_rewards_pool":       nil, // Community incentives
		"tech_innovation_pool":         nil, // R&D and acquisitions
		"operations_pool":              nil, // Platform maintenance
		"talent_acquisition_pool":      nil, // Global hiring
		"strategic_reserve_pool":       nil, // Emergency fund
		"co_founders_pool":             nil, // Co-founder compensation
		"angel_investors_pool":         nil, // Angel investor rewards
		// Lending Suite Module Accounts
		krishimitratypes.ModuleName:     nil, // Agricultural lending
		vyavasayamitratypes.ModuleName:   nil, // Business lending
		shikshamitratypes.ModuleName:     nil, // Education loans
		"krishi_loan_pool":              nil, // Agricultural loan disbursement
		"vyavasaya_loan_pool":           nil, // Business loan disbursement
		"shiksha_loan_pool":             nil, // Education loan disbursement
		"lending_insurance_pool":        nil, // Loan default insurance
		"lending_subsidy_pool":          nil, // Interest subsidy pool
	}
)

var (
	_ runtime.AppI            = (*DeshChainApp)(nil)
	_ servertypes.Application = (*DeshChainApp)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)
}

// DeshChainApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type DeshChainApp struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry
	txConfig          client.TxConfig

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	MoneyOrderKeeper      moneyorderkeeper.Keeper
	CulturalKeeper        culturalkeeper.Keeper
	NAMOKeeper            namokeeper.Keeper
	DhanSetuKeeper        dhansetukeeper.Keeper
	DINRKeeper            dinrkeeper.Keeper
	TradeFinanceKeeper    tradefinancekeeper.Keeper
	OracleKeeper          oraclekeeper.Keeper
	SikkebaazKeeper       sikkebaazkeeper.Keeper
	KrishiMitraKeeper     krishimitrakeeper.Keeper
	VyavasayaMitraKeeper  vyavasayamitrakeeper.Keeper
	ShikshaMitraKeeper    shikshamitrakeeper.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// module configurator
	configurator module.Configurator
}

// New returns a reference to an initialized DeshChainApp.
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *DeshChainApp {
	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry
	txConfig := encodingConfig.TxConfig

	bApp := baseapp.NewBaseApp(Name, logger, db, txConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey, crisistypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, upgradetypes.StoreKey, feegrant.StoreKey,
		evidencetypes.StoreKey, authzkeeper.StoreKey, consensusparamtypes.StoreKey,
		moneyordertypes.StoreKey, culturaltypes.StoreKey, namotypes.StoreKey, dhansettypes.StoreKey,
		dinrtypes.StoreKey, tradefinancetypes.StoreKey, oracletypes.StoreKey, sikkebaaztypes.StoreKey, krishimitratypes.StoreKey, vyavasayamitratypes.StoreKey, 
		shikshamitratypes.StoreKey,
	)
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := storetypes.NewMemoryStoreKeys(
		dhansettypes.MemStoreKey, tradefinancetypes.MemStoreKey, oracletypes.MemStoreKey, sikkebaaztypes.MemStoreKey,
		krishimitratypes.MemStoreKey, vyavasayamitratypes.MemStoreKey,
		shikshamitratypes.MemStoreKey,
	)

	app := &DeshChainApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		txConfig:          txConfig,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]), authtypes.NewModuleAddress(govtypes.ModuleName).String(), runtime.EventService{})
	bApp.SetParamStore(app.ConsensusParamsKeeper.ParamsStore)

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, runtime.NewKVStoreService(keys[authtypes.StoreKey]), authtypes.ProtoBaseAccount, maccPerms, AccountAddressPrefix, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, runtime.NewKVStoreService(keys[banktypes.StoreKey]), app.AccountKeeper, BlockedAddresses(), authtypes.NewModuleAddress(govtypes.ModuleName).String(), logger,
	)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(keys[stakingtypes.StoreKey]), app.AccountKeeper, app.BankKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(), authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()), authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(keys[minttypes.StoreKey]), app.StakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(keys[distrtypes.StoreKey]), app.AccountKeeper, app.BankKeeper,
		app.StakingKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, codec.NewLegacyAmino(), runtime.NewKVStoreService(keys[slashingtypes.StoreKey]), app.StakingKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(keys[crisistypes.StoreKey]), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String(), app.AccountKeeper.AddressCodec(),
	)

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(keys[feegrant.StoreKey]), app.AccountKeeper)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(runtime.NewKVStoreService(keys[authzkeeper.StoreKey]), appCodec, app.MsgServiceRouter(), app.AccountKeeper)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, runtime.NewKVStoreService(keys[upgradetypes.StoreKey]), appCodec, homePath, app.BaseApp, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper))

	govConfig := govtypes.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(keys[govtypes.StoreKey]), app.AccountKeeper, app.BankKeeper,
		app.StakingKeeper, app.DistrKeeper, app.MsgServiceRouter(), govConfig, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(keys[evidencetypes.StoreKey]), app.StakingKeeper, app.SlashingKeeper, app.AccountKeeper.AddressCodec(),
	)
	app.EvidenceKeeper = *evidenceKeeper

	// Initialize Money Order Keeper
	app.MoneyOrderKeeper = moneyorderkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[moneyordertypes.StoreKey]),
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(moneyordertypes.ModuleName),
	)

	// Initialize Cultural Keeper
	app.CulturalKeeper = culturalkeeper.NewKeeper(
		appCodec,
		keys[culturaltypes.StoreKey],
		app.GetSubspace(culturaltypes.ModuleName),
	)

	// Initialize NAMO Keeper
	app.NAMOKeeper = namokeeper.NewKeeper(
		appCodec,
		keys[namotypes.StoreKey],
		app.GetSubspace(namotypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
	)

	// Initialize DhanSetu Keeper
	app.DhanSetuKeeper = dhansetukeeper.NewKeeper(
		appCodec,
		keys[dhansettypes.StoreKey],
		memKeys[dhansettypes.MemStoreKey],
		app.GetSubspace(dhansettypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.MoneyOrderKeeper,
		app.CulturalKeeper,
		app.NAMOKeeper,
	)

	// Initialize Oracle Keeper (must be initialized before DINR)
	app.OracleKeeper = oraclekeeper.NewKeeper(
		appCodec,
		keys[oracletypes.StoreKey],
		memKeys[oracletypes.MemStoreKey],
		app.GetSubspace(oracletypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
	)

	// Initialize DINR Keeper (now with oracle keeper available)
	app.DINRKeeper = dinrkeeper.NewKeeper(
		appCodec,
		keys[dinrtypes.StoreKey],
		app.GetSubspace(dinrtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&app.OracleKeeper, // Oracle keeper is now available
		nil, // Revenue keeper will be added later
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Initialize Trade Finance Keeper
	app.TradeFinanceKeeper = tradefinancekeeper.NewKeeper(
		appCodec,
		keys[tradefinancetypes.StoreKey],
		memKeys[tradefinancetypes.MemStoreKey],
		app.GetSubspace(tradefinancetypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.DINRKeeper, // For DINR integration
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Initialize Sikkebaaz Keeper
	app.SikkebaazKeeper = sikkebaazkeeper.NewKeeper(
		appCodec,
		keys[sikkebaaztypes.StoreKey],
		memKeys[sikkebaaztypes.MemStoreKey],
		app.GetSubspace(sikkebaaztypes.ModuleName),
		app.BankKeeper,
		app.AccountKeeper,
		app.CulturalKeeper,
		app.NAMOKeeper, // Treasury keeper interface
	)

	// Initialize Lending Suite Keepers
	// Initialize Krishi Mitra (Agricultural Lending) Keeper
	app.KrishiMitraKeeper = krishimitrakeeper.NewKeeper(
		appCodec,
		keys[krishimitratypes.StoreKey],
		memKeys[krishimitratypes.MemStoreKey],
		app.GetSubspace(krishimitratypes.ModuleName),
		app.BankKeeper,
		app.AccountKeeper,
		app.DhanSetuKeeper, // For DhanPata verification
	)

	// Initialize Vyavasaya Mitra (Business Lending) Keeper
	app.VyavasayaMitraKeeper = vyavasayamitrakeeper.NewKeeper(
		appCodec,
		keys[vyavasayamitratypes.StoreKey],
		memKeys[vyavasayamitratypes.MemStoreKey],
		app.GetSubspace(vyavasayamitratypes.ModuleName),
		app.BankKeeper,
		app.AccountKeeper,
		app.DhanSetuKeeper, // For DhanPata verification
	)

	// Initialize Shiksha Mitra (Education Loans) Keeper
	app.ShikshaMitraKeeper = shikshamitrakeeper.NewKeeper(
		appCodec,
		keys[shikshamitratypes.StoreKey],
		memKeys[shikshamitratypes.MemStoreKey],
		app.GetSubspace(shikshamitratypes.ModuleName),
		app.BankKeeper,
		app.AccountKeeper,
		app.DhanSetuKeeper, // For DhanPata verification
	)

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, &app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(app.EvidenceKeeper),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),
		moneyorder.NewAppModule(appCodec, app.MoneyOrderKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(moneyordertypes.ModuleName)),
		cultural.NewAppModule(appCodec, app.CulturalKeeper, app.AccountKeeper, app.BankKeeper),
		namo.NewAppModule(appCodec, app.NAMOKeeper, app.AccountKeeper, app.BankKeeper),
		dhansetu.NewAppModule(appCodec, app.DhanSetuKeeper, app.AccountKeeper, app.BankKeeper),
		dinr.NewAppModule(appCodec, app.DINRKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(dinrtypes.ModuleName)),
		tradefinance.NewAppModule(appCodec, app.TradeFinanceKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(tradefinancetypes.ModuleName)),
		oracle.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
		sikkebaaz.NewAppModule(appCodec, app.SikkebaazKeeper, app.AccountKeeper, app.BankKeeper),
		krishimitra.NewAppModule(appCodec, app.KrishiMitraKeeper, app.AccountKeeper, app.BankKeeper),
		vyavasayamitra.NewAppModule(appCodec, app.VyavasayaMitraKeeper, app.AccountKeeper, app.BankKeeper),
		shikshamitra.NewAppModule(appCodec, app.ShikshaMitraKeeper, app.AccountKeeper, app.BankKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authz.ModuleName,
		genutiltypes.ModuleName,
		moneyordertypes.ModuleName,
		culturaltypes.ModuleName,
		namotypes.ModuleName,
		dhansettypes.ModuleName,
		oracletypes.ModuleName,
		dinrtypes.ModuleName,
		tradefinancetypes.ModuleName,
		sikkebaaztypes.ModuleName,
		krishimitratypes.ModuleName,
		vyavasayamitratypes.ModuleName,
		shikshamitratypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		feegrant.ModuleName,
		genutiltypes.ModuleName,
		moneyordertypes.ModuleName,
		culturaltypes.ModuleName,
		namotypes.ModuleName,
		dhansettypes.ModuleName,
		oracletypes.ModuleName,
		dinrtypes.ModuleName,
		tradefinancetypes.ModuleName,
		sikkebaaztypes.ModuleName,
		krishimitratypes.ModuleName,
		vyavasayamitratypes.ModuleName,
		shikshamitratypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	genesisModuleOrder := []string{
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		moneyordertypes.ModuleName,
		culturaltypes.ModuleName,
		namotypes.ModuleName,
		dhansettypes.ModuleName,
		oracletypes.ModuleName,
		dinrtypes.ModuleName,
		tradefinancetypes.ModuleName,
		sikkebaaztypes.ModuleName,
		krishimitratypes.ModuleName,
		vyavasayamitratypes.ModuleName,
		shikshamitratypes.ModuleName,
	}
	app.mm.SetOrderInitGenesis(genesisModuleOrder...)
	app.mm.SetOrderExportGenesis(genesisModuleOrder...)

	// Uncomment if you want to set a custom migration order here.
	// app.mm.SetOrderMigrations(custom order)

	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err := app.mm.RegisterServices(app.configurator)
	if err != nil {
		panic(err)
	}

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			SignModeHandler: txConfig.SignModeHandler(),
			FeegrantKeeper:  app.FeeGrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// Name returns the name of the App
func (app *DeshChainApp) Name() string { return app.BaseApp.Name() }

// GetBaseApp returns the base app of the application
func (app *DeshChainApp) GetBaseApp() *baseapp.BaseApp { return app.BaseApp }

// BeginBlocker application updates every begin block
func (app *DeshChainApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	// Run DINR begin blocker
	dinr.BeginBlocker(ctx, abci.RequestBeginBlock{}, app.DINRKeeper)
	// Run Trade Finance begin blocker
	tradefinance.BeginBlocker(ctx, abci.RequestBeginBlock{}, app.TradeFinanceKeeper)
	return app.mm.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *DeshChainApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *DeshChainApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *DeshChainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *DeshChainApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino returns DeshChainApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *DeshChainApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns DeshChainApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *DeshChainApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns DeshChainApp's InterfaceRegistry
func (app *DeshChainApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns DeshChainApp's TxConfig
func (app *DeshChainApp) TxConfig() client.TxConfig {
	return app.txConfig
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *DeshChainApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *DeshChainApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *DeshChainApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *DeshChainApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *DeshChainApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register app's OpenAPI routes.
	if apiConfig.Swagger {
		RegisterOpenAPIService(Name, apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *DeshChainApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *DeshChainApp) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

func (app *DeshChainApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// SimulationManager implements the SimulationApp interface
func (app *DeshChainApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// configure store loader that checks if version == upgradeHeight and applies store upgrades
func (app *DeshChainApp) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		if upgradeInfo.Name == upgrade.UpgradeName {
			storeUpgrades := upgrade.StoreUpgrades
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
		}
	}
}

func (app *DeshChainApp) setupUpgradeHandlers() {
	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(app.mm, app.configurator),
		)
	}
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	// Allow the following addresses to receive funds
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	return modAccAddrs
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(moneyordertypes.ModuleName)
	paramsKeeper.Subspace(culturaltypes.ModuleName)
	paramsKeeper.Subspace(namotypes.ModuleName)
	paramsKeeper.Subspace(dhansettypes.ModuleName)

	return paramsKeeper
}