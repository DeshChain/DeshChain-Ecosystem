package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IBC constants and identifiers
const (
	// MoneyOrderPortID is the port ID for the money order module
	MoneyOrderPortID = "moneyorder"
	
	// Version defines the current version of the money order IBC application
	Version = "moneyorder-1"
)

// PacketType represents different types of IBC packets
type PacketType string

const (
	PacketType_MONEY_ORDER_TRANSFER     PacketType = "MONEY_ORDER_TRANSFER"
	PacketType_MONEY_ORDER_CONFIRMATION PacketType = "MONEY_ORDER_CONFIRMATION"
	PacketType_MONEY_ORDER_REFUND       PacketType = "MONEY_ORDER_REFUND"
	PacketType_MONEY_ORDER_QUERY        PacketType = "MONEY_ORDER_QUERY"
)

// ChannelState represents the state of an IBC channel
type ChannelState int32

const (
	ChannelState_UNINITIALIZED ChannelState = 0
	ChannelState_INIT          ChannelState = 1
	ChannelState_TRYOPEN       ChannelState = 2
	ChannelState_OPEN          ChannelState = 3
	ChannelState_ACTIVE        ChannelState = 4
	ChannelState_CLOSING       ChannelState = 5
	ChannelState_CLOSED        ChannelState = 6
)

// CrossChainStatus represents the status of a cross-chain money order
type CrossChainStatus int32

const (
	CrossChainStatus_PENDING   CrossChainStatus = 0
	CrossChainStatus_SENT      CrossChainStatus = 1
	CrossChainStatus_RECEIVED  CrossChainStatus = 2
	CrossChainStatus_CONFIRMED CrossChainStatus = 3
	CrossChainStatus_COMPLETED CrossChainStatus = 4
	CrossChainStatus_FAILED    CrossChainStatus = 5
	CrossChainStatus_REFUNDED  CrossChainStatus = 6
	CrossChainStatus_TIMEOUT   CrossChainStatus = 7
)

// String returns the string representation of CrossChainStatus
func (ccs CrossChainStatus) String() string {
	switch ccs {
	case CrossChainStatus_PENDING:
		return "PENDING"
	case CrossChainStatus_SENT:
		return "SENT"
	case CrossChainStatus_RECEIVED:
		return "RECEIVED"
	case CrossChainStatus_CONFIRMED:
		return "CONFIRMED"
	case CrossChainStatus_COMPLETED:
		return "COMPLETED"
	case CrossChainStatus_FAILED:
		return "FAILED"
	case CrossChainStatus_REFUNDED:
		return "REFUNDED"
	case CrossChainStatus_TIMEOUT:
		return "TIMEOUT"
	default:
		return "UNKNOWN"
	}
}

// IBCMoneyOrderPacketData represents the packet data for IBC money order transfers
type IBCMoneyOrderPacketData struct {
	Type             PacketType `json:"type"`
	OrderID          string     `json:"order_id"`
	SenderAddress    string     `json:"sender_address"`
	RecipientAddress string     `json:"recipient_address"`
	Amount           string     `json:"amount"`
	SenderChain      string     `json:"sender_chain"`
	RecipientChain   string     `json:"recipient_chain"`
	Memo             string     `json:"memo,omitempty"`
	Timestamp        int64      `json:"timestamp"`
	Sequence         uint64     `json:"sequence"`
	Status           string     `json:"status,omitempty"`
	ErrorMessage     string     `json:"error_message,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// Validate validates the IBC packet data
func (data IBCMoneyOrderPacketData) Validate() error {
	if data.Type == "" {
		return sdkerrors.Wrap(ErrInvalidInput, "packet type cannot be empty")
	}
	
	if data.OrderID == "" {
		return sdkerrors.Wrap(ErrInvalidInput, "order ID cannot be empty")
	}
	
	if data.SenderAddress == "" {
		return sdkerrors.Wrap(ErrInvalidInput, "sender address cannot be empty")
	}
	
	if data.RecipientAddress == "" {
		return sdkerrors.Wrap(ErrInvalidInput, "recipient address cannot be empty")
	}
	
	if data.Amount == "" {
		return sdkerrors.Wrap(ErrInvalidAmount, "amount cannot be empty")
	}
	
	// Validate amount is positive
	amount, ok := sdk.NewIntFromString(data.Amount)
	if !ok {
		return sdkerrors.Wrap(ErrInvalidAmount, "invalid amount format")
	}
	
	if amount.IsZero() || amount.IsNegative() {
		return sdkerrors.Wrap(ErrInvalidAmount, "amount must be positive")
	}
	
	if data.SenderChain == "" {
		return sdkerrors.Wrap(ErrInvalidInput, "sender chain cannot be empty")
	}
	
	if data.RecipientChain == "" {
		return sdkerrors.Wrap(ErrInvalidInput, "recipient chain cannot be empty")
	}
	
	if data.Timestamp <= 0 {
		return sdkerrors.Wrap(ErrInvalidInput, "timestamp must be positive")
	}
	
	return nil
}

// GetBytes returns the marshaled packet data
func (data IBCMoneyOrderPacketData) GetBytes() []byte {
	bz, err := ModuleCdc.MarshalJSON(&data)
	if err != nil {
		panic(err)
	}
	return bz
}

// IBCMoneyOrderMetadata represents metadata for IBC money order channels
type IBCMoneyOrderMetadata struct {
	ChainID                string        `json:"chain_id"`
	ChainName              string        `json:"chain_name"`
	Version                string        `json:"version"`
	Capabilities           []string      `json:"capabilities"`
	MinTransferAmount      sdk.Int       `json:"min_transfer_amount"`
	MaxTransferAmount      sdk.Int       `json:"max_transfer_amount"`
	TransferFee            sdk.Dec       `json:"transfer_fee"`
	EstimatedTransferTime  time.Duration `json:"estimated_transfer_time"`
	SupportedDenoms        []string      `json:"supported_denoms"`
	RequiresKYC            bool          `json:"requires_kyc"`
	ComplianceLevel        string        `json:"compliance_level"`
	OperatingHours         string        `json:"operating_hours,omitempty"`
	MaintenanceWindows     []string      `json:"maintenance_windows,omitempty"`
	ContactInfo            ContactInfo   `json:"contact_info,omitempty"`
}

// ContactInfo represents contact information for the connected chain
type ContactInfo struct {
	Organization string `json:"organization,omitempty"`
	Email        string `json:"email,omitempty"`
	Website      string `json:"website,omitempty"`
	Support      string `json:"support,omitempty"`
}

// Validate validates the IBC metadata
func (metadata IBCMoneyOrderMetadata) Validate() error {
	if metadata.ChainID == "" {
		return sdkerrors.Wrap(ErrInvalidInput, "chain ID cannot be empty")
	}
	
	if metadata.Version != Version {
		return sdkerrors.Wrapf(ErrInvalidInput, "unsupported version: expected %s, got %s", Version, metadata.Version)
	}
	
	if metadata.MinTransferAmount.IsNegative() {
		return sdkerrors.Wrap(ErrInvalidAmount, "minimum transfer amount cannot be negative")
	}
	
	if metadata.MaxTransferAmount.IsNegative() {
		return sdkerrors.Wrap(ErrInvalidAmount, "maximum transfer amount cannot be negative")
	}
	
	if metadata.MinTransferAmount.GT(metadata.MaxTransferAmount) {
		return sdkerrors.Wrap(ErrInvalidAmount, "minimum transfer amount cannot be greater than maximum")
	}
	
	if metadata.TransferFee.IsNegative() {
		return sdkerrors.Wrap(ErrInvalidInput, "transfer fee cannot be negative")
	}
	
	return nil
}

// IBCChannelInfo represents information about an IBC channel
type IBCChannelInfo struct {
	ChannelID             string                `json:"channel_id"`
	CounterpartyPortID    string                `json:"counterparty_port_id"`
	CounterpartyChannelID string                `json:"counterparty_channel_id"`
	ConnectionID          string                `json:"connection_id"`
	State                 ChannelState          `json:"state"`
	Metadata              IBCMoneyOrderMetadata `json:"metadata"`
	CreatedAt             time.Time             `json:"created_at"`
	ClosedAt              time.Time             `json:"closed_at,omitempty"`
	LastActivity          time.Time             `json:"last_activity,omitempty"`
	TotalPacketsSent      uint64                `json:"total_packets_sent"`
	TotalPacketsReceived  uint64                `json:"total_packets_received"`
	SuccessfulTransfers   uint64                `json:"successful_transfers"`
	FailedTransfers       uint64                `json:"failed_transfers"`
	TotalVolume           sdk.Int               `json:"total_volume"`
}

// CrossChainMoneyOrder represents a cross-chain money order
type CrossChainMoneyOrder struct {
	OrderID          string           `json:"order_id"`
	SenderAddress    string           `json:"sender_address"`
	RecipientAddress string           `json:"recipient_address"`
	Amount           sdk.Int          `json:"amount"`
	SenderChain      string           `json:"sender_chain"`
	RecipientChain   string           `json:"recipient_chain"`
	ChannelID        string           `json:"channel_id"`
	Status           CrossChainStatus `json:"status"`
	CreatedAt        time.Time        `json:"created_at"`
	SentAt           time.Time        `json:"sent_at,omitempty"`
	ReceivedAt       time.Time        `json:"received_at,omitempty"`
	ConfirmedAt      time.Time        `json:"confirmed_at,omitempty"`
	CompletedAt      time.Time        `json:"completed_at,omitempty"`
	FailedAt         time.Time        `json:"failed_at,omitempty"`
	RefundedAt       time.Time        `json:"refunded_at,omitempty"`
	TimeoutAt        time.Time        `json:"timeout_at,omitempty"`
	Memo             string           `json:"memo,omitempty"`
	TimeoutHeight    uint64           `json:"timeout_height,omitempty"`
	TimeoutTimestamp uint64           `json:"timeout_timestamp,omitempty"`
	ErrorMessage     string           `json:"error_message,omitempty"`
	PacketSequence   uint64           `json:"packet_sequence,omitempty"`
	RelayerAddress   string           `json:"relayer_address,omitempty"`
	Fees             CrossChainFees   `json:"fees"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// CrossChainFees represents fees associated with cross-chain transfers
type CrossChainFees struct {
	BaseFee     sdk.Int `json:"base_fee"`
	RelayerFee  sdk.Int `json:"relayer_fee"`
	ProtocolFee sdk.Int `json:"protocol_fee"`
	TotalFee    sdk.Int `json:"total_fee"`
}

// SupportedChain represents a chain that supports cross-chain money orders
type SupportedChain struct {
	ChainID           string        `json:"chain_id"`
	ChainName         string        `json:"chain_name"`
	ChannelID         string        `json:"channel_id"`
	PortID            string        `json:"port_id"`
	IsActive          bool          `json:"is_active"`
	Capabilities      []string      `json:"capabilities"`
	MinAmount         sdk.Int       `json:"min_amount"`
	MaxAmount         sdk.Int       `json:"max_amount"`
	Fee               sdk.Dec       `json:"fee"`
	EstimatedTime     time.Duration `json:"estimated_time"`
	SupportedAssets   []string      `json:"supported_assets"`
	RequiresKYC       bool          `json:"requires_kyc"`
	ComplianceLevel   string        `json:"compliance_level"`
	NetworkHealth     string        `json:"network_health"` // "healthy", "degraded", "down"
	LastHealthCheck   time.Time     `json:"last_health_check"`
}

// CrossChainStatusResponse represents the response for cross-chain status queries
type CrossChainStatusResponse struct {
	Order       CrossChainMoneyOrder `json:"order"`
	ChannelInfo IBCChannelInfo       `json:"channel_info"`
	Timeline    []TimelineEvent      `json:"timeline"`
	Fees        CrossChainFees       `json:"fees"`
	Estimated   EstimatedInfo        `json:"estimated"`
}

// TimelineEvent represents an event in the cross-chain transfer timeline
type TimelineEvent struct {
	EventType   string    `json:"event_type"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
	BlockHeight uint64    `json:"block_height,omitempty"`
	TxHash      string    `json:"tx_hash,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// EstimatedInfo provides estimated information about the transfer
type EstimatedInfo struct {
	EstimatedCompletion time.Time     `json:"estimated_completion"`
	EstimatedDuration   time.Duration `json:"estimated_duration"`
	Confidence          float64       `json:"confidence"`
	NextStep            string        `json:"next_step"`
	RequiresAction      bool          `json:"requires_action"`
	ActionDescription   string        `json:"action_description,omitempty"`
}

// IBCMoneyOrderAcknowledgement represents the acknowledgement for IBC money order packets
type IBCMoneyOrderAcknowledgement struct {
	Success bool   `json:"success"`
	Result  string `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
	OrderID string `json:"order_id,omitempty"`
}

// CrossChainMoneyOrderRequest represents a request to create a cross-chain money order
type CrossChainMoneyOrderRequest struct {
	SenderAddress    string `json:"sender_address"`
	RecipientAddress string `json:"recipient_address"`
	Amount           string `json:"amount"`
	RecipientChain   string `json:"recipient_chain"`
	Memo             string `json:"memo,omitempty"`
	TimeoutMinutes   uint64 `json:"timeout_minutes,omitempty"`
	Priority         string `json:"priority,omitempty"`
	NotifyEmail      string `json:"notify_email,omitempty"`
}

// Validate validates the cross-chain money order request
func (req CrossChainMoneyOrderRequest) Validate() error {
	if req.SenderAddress == "" {
		return sdkerrors.Wrap(ErrInvalidSender, "sender address cannot be empty")
	}
	
	if req.RecipientAddress == "" {
		return sdkerrors.Wrap(ErrInvalidRecipient, "recipient address cannot be empty")
	}
	
	if req.Amount == "" {
		return sdkerrors.Wrap(ErrInvalidAmount, "amount cannot be empty")
	}
	
	amount, ok := sdk.NewIntFromString(req.Amount)
	if !ok {
		return sdkerrors.Wrap(ErrInvalidAmount, "invalid amount format")
	}
	
	if amount.IsZero() || amount.IsNegative() {
		return sdkerrors.Wrap(ErrInvalidAmount, "amount must be positive")
	}
	
	if req.RecipientChain == "" {
		return sdkerrors.Wrap(ErrInvalidInput, "recipient chain cannot be empty")
	}
	
	if len(req.Memo) > 256 {
		return sdkerrors.Wrap(ErrInvalidInput, "memo too long (max 256 characters)")
	}
	
	return nil
}

// IBCMoneyOrderStats represents statistics for IBC money order operations
type IBCMoneyOrderStats struct {
	TotalChannels        int64                        `json:"total_channels"`
	ActiveChannels       int64                        `json:"active_channels"`
	TotalTransfers       int64                        `json:"total_transfers"`
	SuccessfulTransfers  int64                        `json:"successful_transfers"`
	FailedTransfers      int64                        `json:"failed_transfers"`
	TotalVolume          sdk.Int                      `json:"total_volume"`
	AverageTransferTime  time.Duration                `json:"average_transfer_time"`
	SuccessRate          float64                      `json:"success_rate"`
	TopChains            []ChainStats                 `json:"top_chains"`
	RecentActivity       []RecentTransfer             `json:"recent_activity"`
	ChannelHealthStatus  map[string]ChannelHealth     `json:"channel_health_status"`
}

// ChainStats represents statistics for a specific chain
type ChainStats struct {
	ChainID     string  `json:"chain_id"`
	ChainName   string  `json:"chain_name"`
	Transfers   int64   `json:"transfers"`
	Volume      sdk.Int `json:"volume"`
	SuccessRate float64 `json:"success_rate"`
	AvgTime     time.Duration `json:"avg_time"`
}

// RecentTransfer represents a recent cross-chain transfer
type RecentTransfer struct {
	OrderID        string           `json:"order_id"`
	SenderChain    string           `json:"sender_chain"`
	RecipientChain string           `json:"recipient_chain"`
	Amount         sdk.Int          `json:"amount"`
	Status         CrossChainStatus `json:"status"`
	Timestamp      time.Time        `json:"timestamp"`
}

// ChannelHealth represents the health status of an IBC channel
type ChannelHealth struct {
	ChannelID       string        `json:"channel_id"`
	IsHealthy       bool          `json:"is_healthy"`
	LastPacketTime  time.Time     `json:"last_packet_time"`
	PacketSuccess   float64       `json:"packet_success_rate"`
	AverageLatency  time.Duration `json:"average_latency"`
	Issues          []string      `json:"issues,omitempty"`
	LastHealthCheck time.Time     `json:"last_health_check"`
}

// IBCRelayerInfo represents information about IBC relayers
type IBCRelayerInfo struct {
	RelayerAddress   string    `json:"relayer_address"`
	PacketsRelayed   int64     `json:"packets_relayed"`
	SuccessRate      float64   `json:"success_rate"`
	AverageLatency   time.Duration `json:"average_latency"`
	LastActive       time.Time `json:"last_active"`
	IsActive         bool      `json:"is_active"`
	ReputationScore  float64   `json:"reputation_score"`
	TotalFees        sdk.Int   `json:"total_fees"`
}

// CrossChainConfiguration represents configuration for cross-chain operations
type CrossChainConfiguration struct {
	EnabledChains        []string      `json:"enabled_chains"`
	MinTransferAmount    sdk.Int       `json:"min_transfer_amount"`
	MaxTransferAmount    sdk.Int       `json:"max_transfer_amount"`
	DefaultTimeoutMinutes uint64       `json:"default_timeout_minutes"`
	BaseFee              sdk.Dec       `json:"base_fee"`
	RelayerFeePercentage sdk.Dec       `json:"relayer_fee_percentage"`
	RequireKYC           bool          `json:"require_kyc"`
	EnableAntiMoney      bool          `json:"enable_anti_money_laundering"`
	MaxDailyVolume       sdk.Int       `json:"max_daily_volume"`
	AlertThresholds      AlertThresholds `json:"alert_thresholds"`
}

// AlertThresholds represents thresholds for generating alerts
type AlertThresholds struct {
	LargeTransferAmount  sdk.Int `json:"large_transfer_amount"`
	HighFrequencyCount   int64   `json:"high_frequency_count"`
	HighFrequencyWindow  time.Duration `json:"high_frequency_window"`
	SuspiciousPatterns   []string `json:"suspicious_patterns"`
	FailureRateThreshold float64 `json:"failure_rate_threshold"`
}

// Key prefixes for IBC storage
var (
	IBCChannelPrefix        = []byte{0x50}
	CrossChainOrderPrefix   = []byte{0x51}
	IBCStatsPrefix          = []byte{0x52}
	IBCConfigPrefix         = []byte{0x53}
	IBCRelayerPrefix        = []byte{0x54}
	IBCEscrowPrefix         = []byte{0x55}
	IBCTimelinePrefix       = []byte{0x56}
)

// Event types for IBC operations
const (
	EventTypeIBCChannelOpen             = "ibc_channel_open"
	EventTypeIBCChannelClose            = "ibc_channel_close"
	EventTypeCrossChainMoneyOrderSent   = "cross_chain_money_order_sent"
	EventTypeCrossChainMoneyOrderReceived = "cross_chain_money_order_received"
	EventTypeCrossChainMoneyOrderConfirmed = "cross_chain_money_order_confirmed"
	EventTypeCrossChainMoneyOrderCompleted = "cross_chain_money_order_completed"
	EventTypeCrossChainMoneyOrderFailed   = "cross_chain_money_order_failed"
	EventTypeCrossChainMoneyOrderRefunded = "cross_chain_money_order_refunded"
	EventTypeCrossChainMoneyOrderTimeout  = "cross_chain_money_order_timeout"
)

// Attribute keys for IBC events
const (
	AttributeKeyChannelID             = "channel_id"
	AttributeKeyCounterpartyChannelID = "counterparty_channel_id"
	AttributeKeyPortID                = "port_id"
	AttributeKeyPacketSequence        = "packet_sequence"
	AttributeKeyPacketTimeoutHeight   = "packet_timeout_height"
	AttributeKeyPacketTimeoutTimestamp = "packet_timeout_timestamp"
	AttributeKeyRecipientChain        = "recipient_chain"
	AttributeKeySenderChain           = "sender_chain"
	AttributeKeyErrorMessage          = "error_message"
	AttributeKeyRelayerAddress        = "relayer_address"
)

// Helper functions

// NewIBCMoneyOrderPacketData creates a new IBC money order packet data
func NewIBCMoneyOrderPacketData(
	packetType PacketType,
	orderID string,
	senderAddress string,
	recipientAddress string,
	amount string,
	senderChain string,
	recipientChain string,
	memo string,
	sequence uint64,
) IBCMoneyOrderPacketData {
	return IBCMoneyOrderPacketData{
		Type:             packetType,
		OrderID:          orderID,
		SenderAddress:    senderAddress,
		RecipientAddress: recipientAddress,
		Amount:           amount,
		SenderChain:      senderChain,
		RecipientChain:   recipientChain,
		Memo:             memo,
		Timestamp:        time.Now().Unix(),
		Sequence:         sequence,
	}
}

// NewIBCMoneyOrderAcknowledgement creates a new acknowledgement
func NewIBCMoneyOrderAcknowledgement(success bool, result string, error string, orderID string) IBCMoneyOrderAcknowledgement {
	return IBCMoneyOrderAcknowledgement{
		Success: success,
		Result:  result,
		Error:   error,
		OrderID: orderID,
	}
}

// CalculateTimeout calculates timeout based on minutes
func CalculateTimeout(timeoutMinutes uint64) (uint64, uint64) {
	if timeoutMinutes == 0 {
		timeoutMinutes = 60 // Default 1 hour
	}
	
	timeoutTimestamp := uint64(time.Now().Add(time.Duration(timeoutMinutes) * time.Minute).UnixNano())
	timeoutHeight := uint64(0) // Use timestamp-based timeout
	
	return timeoutHeight, timeoutTimestamp
}

// IsValidChainID validates if a chain ID is in the correct format
func IsValidChainID(chainID string) bool {
	return len(chainID) > 0 && len(chainID) <= 64
}

// FormatCrossChainOrderID formats a cross-chain order ID
func FormatCrossChainOrderID(prefix string, chainID string, sequence uint64) string {
	return fmt.Sprintf("%s_%s_%d_%d", prefix, chainID, time.Now().Unix(), sequence)
}