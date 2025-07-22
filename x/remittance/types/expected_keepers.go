package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI
	SetModuleAccount(context.Context, sdk.ModuleAccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	HasBalance(ctx context.Context, addr sdk.AccAddress, amt sdk.Coin) bool
}

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(
		ctx sdk.Context,
		chanCap *capabilitytypes.Capability,
		sourcePort string,
		sourceChannel string,
		timeoutHeight ibcexported.Height,
		timeoutTimestamp uint64,
		data []byte,
	) (uint64, error)
	ChanCloseInit(ctx sdk.Context, portID, channelID string, chanCap *capabilitytypes.Capability) error
}

// PortKeeper defines the expected IBC port keeper
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
	AuthenticatePort(ctx sdk.Context, cap *capabilitytypes.Capability, portID string) bool
}

// ScopedKeeper defines the expected IBC scoped keeper
type ScopedKeeper interface {
	GetCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, bool)
	AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool
	ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error
	NewCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, error)
}

// IBCKeeper defines the expected IBC keeper interface
type IBCKeeper interface {
	ClientKeeper() IBCClientKeeper
	ConnectionKeeper() IBCConnectionKeeper
	ChannelKeeper() IBCChannelKeeper
}

// IBCClientKeeper defines the expected IBC client keeper
type IBCClientKeeper interface {
	GetClientState(ctx sdk.Context, clientID string) (ibcexported.ClientState, bool)
	GetClientConsensusState(ctx sdk.Context, clientID string, height ibcexported.Height) (ibcexported.ConsensusState, bool)
}

// IBCConnectionKeeper defines the expected IBC connection keeper
type IBCConnectionKeeper interface {
	GetConnection(ctx sdk.Context, connectionID string) (connectiontypes.ConnectionEnd, bool)
}

// IBCChannelKeeper defines the expected IBC channel keeper
type IBCChannelKeeper interface {
	GetChannel(ctx sdk.Context, portID, channelID string) (channeltypes.Channel, bool)
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, ibcexported.ClientState, error)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	GetNextSequenceRecv(ctx sdk.Context, portID, channelID string) (uint64, bool)
	GetPacketCommitment(ctx sdk.Context, portID, channelID string, sequence uint64) []byte
	GetPacketReceipt(ctx sdk.Context, portID, channelID string, sequence uint64) (string, bool)
	GetPacketAcknowledgement(ctx sdk.Context, portID, channelID string, sequence uint64) ([]byte, bool)
	SendPacket(
		ctx sdk.Context,
		channelCap *capabilitytypes.Capability,
		sourcePort string,
		sourceChannel string,
		timeoutHeight ibcexported.Height,
		timeoutTimestamp uint64,
		data []byte,
	) (uint64, error)
	WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI, acknowledgement []byte) error
}

// OracleKeeper defines the expected oracle keeper for exchange rates
type OracleKeeper interface {
	GetPrice(ctx context.Context, baseCurrency, quoteCurrency string) (sdk.Dec, error)
	GetMultiplePrices(ctx context.Context, pairs []CurrencyPair) (map[string]sdk.Dec, error)
	IsSupported(ctx context.Context, baseCurrency, quoteCurrency string) bool
}

// CurrencyPair represents a currency pair for oracle queries
type CurrencyPair struct {
	Base  string
	Quote string
}

// ComplianceKeeper defines the expected compliance keeper for KYC/AML
type ComplianceKeeper interface {
	GetKYCLevel(ctx context.Context, address string) (KYCLevel, error)
	CheckSanctions(ctx context.Context, address string, country string) error
	CheckPEP(ctx context.Context, address string) error
	CalculateRiskScore(ctx context.Context, transfer RemittanceTransfer) (sdk.Dec, error)
	ReportTransaction(ctx context.Context, transfer RemittanceTransfer) error
}

// SettlementKeeper defines the expected settlement keeper
type SettlementKeeper interface {
	ProcessBankTransfer(ctx context.Context, transfer RemittanceTransfer, partner CorridorPartner) error
	ProcessMobileWallet(ctx context.Context, transfer RemittanceTransfer, partner CorridorPartner) error
	ProcessCashPickup(ctx context.Context, transfer RemittanceTransfer, partner CorridorPartner) error
	GetSettlementStatus(ctx context.Context, transferID string) (SettlementStatus, error)
	GetPartnerBalance(ctx context.Context, partnerID string) (sdk.Coins, error)
}

// SettlementStatus represents the status of a settlement
type SettlementStatus struct {
	Status        string    // pending, processing, completed, failed
	Reference     string    // External reference
	Timestamp     int64     // Unix timestamp
	ErrorMessage  string    // Error message if failed
	PartnerID     string    // Settlement partner ID
}

// LiquidityKeeper defines the expected liquidity keeper
type LiquidityKeeper interface {
	GetLiquidity(ctx context.Context, baseCurrency, quoteCurrency string) (sdk.Dec, error)
	CheckLiquidityAvailability(ctx context.Context, amount sdk.Coin, sourceCurrency, destCurrency string) bool
	ReserveLiquidity(ctx context.Context, poolID string, amount sdk.Coin) error
	ReleaseLiquidity(ctx context.Context, poolID string, amount sdk.Coin) error
	SwapCurrency(ctx context.Context, fromAmount sdk.Coin, toCurrency string) (sdk.Coin, error)
}

// GovernanceKeeper defines the expected governance keeper
type GovernanceKeeper interface {
	GetParams(ctx context.Context) (params GovernanceParams, err error)
	SetParams(ctx context.Context, params GovernanceParams) error
}

// GovernanceParams defines governance parameters
type GovernanceParams struct {
	VotingPeriod     int64 // Voting period for proposals
	DepositRequired  bool  // Whether deposit is required
	MinDeposit       sdk.Coins // Minimum deposit
}