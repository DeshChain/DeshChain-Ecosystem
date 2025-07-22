package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Message type constants for Sewa Mitra operations
const (
	TypeMsgRegisterSewaMitraAgent          = "register_seva_mitra_agent"
	TypeMsgUpdateSewaMitraAgent           = "update_seva_mitra_agent"
	TypeMsgActivateSewaMitraAgent         = "activate_seva_mitra_agent"
	TypeMsgSuspendSewaMitraAgent          = "suspend_seva_mitra_agent"
	TypeMsgDeactivateSewaMitraAgent       = "deactivate_seva_mitra_agent"
	TypeMsgInitiateRemittanceWithSewaMitra = "initiate_remittance_with_seva_mitra"
	TypeMsgConfirmSewaMitraPickup         = "confirm_seva_mitra_pickup"
	TypeMsgPaySewaMitraCommission         = "pay_seva_mitra_commission"
)

// ========== MsgRegisterSewaMitraAgent ==========

// MsgRegisterSewaMitraAgent represents a message to register a new Sewa Mitra agent
type MsgRegisterSewaMitraAgent struct {
	Authority            string               `json:"authority" yaml:"authority"`
	AgentId              string               `json:"agent_id" yaml:"agent_id"`
	AgentAddress         string               `json:"agent_address" yaml:"agent_address"`
	AgentName            string               `json:"agent_name" yaml:"agent_name"`
	BusinessName         string               `json:"business_name" yaml:"business_name"`
	Country              string               `json:"country" yaml:"country"`
	State                string               `json:"state" yaml:"state"`
	City                 string               `json:"city" yaml:"city"`
	PostalCode           string               `json:"postal_code" yaml:"postal_code"`
	AddressLine1         string               `json:"address_line1" yaml:"address_line1"`
	AddressLine2         string               `json:"address_line2" yaml:"address_line2"`
	Phone                string               `json:"phone" yaml:"phone"`
	Email                string               `json:"email" yaml:"email"`
	SupportedCurrencies  []string             `json:"supported_currencies" yaml:"supported_currencies"`
	SupportedMethods     []SettlementMethod   `json:"supported_methods" yaml:"supported_methods"`
	LiquidityLimit       sdk.Coin             `json:"liquidity_limit" yaml:"liquidity_limit"`
	DailyLimit           sdk.Coin             `json:"daily_limit" yaml:"daily_limit"`
	BaseCommissionRate   string               `json:"base_commission_rate" yaml:"base_commission_rate"`
	VolumeBonusRate      string               `json:"volume_bonus_rate" yaml:"volume_bonus_rate"`
	MinimumCommission    sdk.Coin             `json:"minimum_commission" yaml:"minimum_commission"`
	MaximumCommission    sdk.Coin             `json:"maximum_commission" yaml:"maximum_commission"`
}

// NewMsgRegisterSewaMitraAgent creates a new MsgRegisterSewaMitraAgent instance
func NewMsgRegisterSewaMitraAgent(
	authority, agentId, agentAddress, agentName, businessName string,
	country, state, city, postalCode, addressLine1, addressLine2, phone, email string,
	supportedCurrencies []string, supportedMethods []SettlementMethod,
	liquidityLimit, dailyLimit, minimumCommission, maximumCommission sdk.Coin,
	baseCommissionRate, volumeBonusRate string,
) *MsgRegisterSewaMitraAgent {
	return &MsgRegisterSewaMitraAgent{
		Authority:           authority,
		AgentId:             agentId,
		AgentAddress:        agentAddress,
		AgentName:           agentName,
		BusinessName:        businessName,
		Country:             country,
		State:               state,
		City:                city,
		PostalCode:          postalCode,
		AddressLine1:        addressLine1,
		AddressLine2:        addressLine2,
		Phone:               phone,
		Email:               email,
		SupportedCurrencies: supportedCurrencies,
		SupportedMethods:    supportedMethods,
		LiquidityLimit:      liquidityLimit,
		DailyLimit:          dailyLimit,
		BaseCommissionRate:  baseCommissionRate,
		VolumeBonusRate:     volumeBonusRate,
		MinimumCommission:   minimumCommission,
		MaximumCommission:   maximumCommission,
	}
}

// Route returns the message route for routing
func (msg MsgRegisterSewaMitraAgent) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg MsgRegisterSewaMitraAgent) Type() string {
	return TypeMsgRegisterSewaMitraAgent
}

// GetSigners returns the expected signers for the message
func (msg MsgRegisterSewaMitraAgent) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgRegisterSewaMitraAgent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation of the message
func (msg MsgRegisterSewaMitraAgent) ValidateBasic() error {
	// Validate authority
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	// Validate agent ID
	if err := ValidateAgentID(msg.AgentId); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Validate agent address
	if _, err := sdk.AccAddressFromBech32(msg.AgentAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid agent address (%s)", err)
	}

	// Validate required fields
	if msg.AgentName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "agent name is required")
	}

	if err := ValidateLocationInfo(msg.Country, msg.State, msg.City); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Validate limits
	if !msg.LiquidityLimit.IsValid() || !msg.LiquidityLimit.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "liquidity limit must be positive")
	}

	if !msg.DailyLimit.IsValid() || !msg.DailyLimit.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "daily limit must be positive")
	}

	// Validate commission rates
	baseRate, err := sdk.NewDecFromStr(msg.BaseCommissionRate)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid base commission rate")
	}
	if err := ValidateCommissionRate(baseRate, sdk.NewDecWithPrec(10, 2)); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	volumeRate, err := sdk.NewDecFromStr(msg.VolumeBonusRate)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid volume bonus rate")
	}
	if err := ValidateCommissionRate(volumeRate, sdk.NewDecWithPrec(5, 2)); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

// MsgRegisterSewaMitraAgentResponse is the response for MsgRegisterSewaMitraAgent
type MsgRegisterSewaMitraAgentResponse struct {
	AgentId string `json:"agent_id" yaml:"agent_id"`
	Status  string `json:"status" yaml:"status"`
}

// ========== MsgUpdateSewaMitraAgent ==========

// MsgUpdateSewaMitraAgent represents a message to update a Sewa Mitra agent
type MsgUpdateSewaMitraAgent struct {
	Sender              string             `json:"sender" yaml:"sender"`
	AgentId             string             `json:"agent_id" yaml:"agent_id"`
	Phone               string             `json:"phone,omitempty" yaml:"phone,omitempty"`
	Email               string             `json:"email,omitempty" yaml:"email,omitempty"`
	SupportedCurrencies []string           `json:"supported_currencies,omitempty" yaml:"supported_currencies,omitempty"`
	SupportedMethods    []SettlementMethod `json:"supported_methods,omitempty" yaml:"supported_methods,omitempty"`
	LiquidityLimit      sdk.Coin           `json:"liquidity_limit,omitempty" yaml:"liquidity_limit,omitempty"`
	DailyLimit          sdk.Coin           `json:"daily_limit,omitempty" yaml:"daily_limit,omitempty"`
}

// NewMsgUpdateSewaMitraAgent creates a new MsgUpdateSewaMitraAgent instance
func NewMsgUpdateSewaMitraAgent(sender, agentId string) *MsgUpdateSewaMitraAgent {
	return &MsgUpdateSewaMitraAgent{
		Sender:  sender,
		AgentId: agentId,
	}
}

// Route returns the message route for routing
func (msg MsgUpdateSewaMitraAgent) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg MsgUpdateSewaMitraAgent) Type() string {
	return TypeMsgUpdateSewaMitraAgent
}

// GetSigners returns the expected signers for the message
func (msg MsgUpdateSewaMitraAgent) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgUpdateSewaMitraAgent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation of the message
func (msg MsgUpdateSewaMitraAgent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if err := ValidateAgentID(msg.AgentId); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

// MsgUpdateSewaMitraAgentResponse is the response for MsgUpdateSewaMitraAgent
type MsgUpdateSewaMitraAgentResponse struct {
	Success bool `json:"success" yaml:"success"`
}

// ========== MsgActivateSewaMitraAgent ==========

// MsgActivateSewaMitraAgent represents a message to activate a Sewa Mitra agent
type MsgActivateSewaMitraAgent struct {
	Authority string `json:"authority" yaml:"authority"`
	AgentId   string `json:"agent_id" yaml:"agent_id"`
}

// NewMsgActivateSewaMitraAgent creates a new MsgActivateSewaMitraAgent instance
func NewMsgActivateSewaMitraAgent(authority, agentId string) *MsgActivateSewaMitraAgent {
	return &MsgActivateSewaMitraAgent{
		Authority: authority,
		AgentId:   agentId,
	}
}

// Route returns the message route for routing
func (msg MsgActivateSewaMitraAgent) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg MsgActivateSewaMitraAgent) Type() string {
	return TypeMsgActivateSewaMitraAgent
}

// GetSigners returns the expected signers for the message
func (msg MsgActivateSewaMitraAgent) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgActivateSewaMitraAgent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation of the message
func (msg MsgActivateSewaMitraAgent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if err := ValidateAgentID(msg.AgentId); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

// MsgActivateSewaMitraAgentResponse is the response for MsgActivateSewaMitraAgent
type MsgActivateSewaMitraAgentResponse struct {
	Success bool `json:"success" yaml:"success"`
}

// ========== MsgSuspendSewaMitraAgent ==========

// MsgSuspendSewaMitraAgent represents a message to suspend a Sewa Mitra agent
type MsgSuspendSewaMitraAgent struct {
	Authority             string `json:"authority" yaml:"authority"`
	AgentId               string `json:"agent_id" yaml:"agent_id"`
	Reason                string `json:"reason" yaml:"reason"`
	SuspensionDurationDays int32  `json:"suspension_duration_days" yaml:"suspension_duration_days"`
}

// NewMsgSuspendSewaMitraAgent creates a new MsgSuspendSewaMitraAgent instance
func NewMsgSuspendSewaMitraAgent(authority, agentId, reason string, durationDays int32) *MsgSuspendSewaMitraAgent {
	return &MsgSuspendSewaMitraAgent{
		Authority:             authority,
		AgentId:               agentId,
		Reason:                reason,
		SuspensionDurationDays: durationDays,
	}
}

// Route returns the message route for routing
func (msg MsgSuspendSewaMitraAgent) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg MsgSuspendSewaMitraAgent) Type() string {
	return TypeMsgSuspendSewaMitraAgent
}

// GetSigners returns the expected signers for the message
func (msg MsgSuspendSewaMitraAgent) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgSuspendSewaMitraAgent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation of the message
func (msg MsgSuspendSewaMitraAgent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if err := ValidateAgentID(msg.AgentId); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if msg.Reason == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "reason is required")
	}

	if msg.SuspensionDurationDays <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "suspension duration must be positive")
	}

	return nil
}

// MsgSuspendSewaMitraAgentResponse is the response for MsgSuspendSewaMitraAgent
type MsgSuspendSewaMitraAgentResponse struct {
	Success        bool      `json:"success" yaml:"success"`
	SuspendedUntil time.Time `json:"suspended_until" yaml:"suspended_until"`
}

// ========== MsgInitiateRemittanceWithSewaMitra ==========

// MsgInitiateRemittanceWithSewaMitra represents a message to initiate remittance with Sewa Mitra preference
type MsgInitiateRemittanceWithSewaMitra struct {
	Sender                       string `json:"sender" yaml:"sender"`
	RecipientAddress             string `json:"recipient_address,omitempty" yaml:"recipient_address,omitempty"`
	RecipientName                string `json:"recipient_name" yaml:"recipient_name"`
	RecipientPhone               string `json:"recipient_phone" yaml:"recipient_phone"`
	RecipientDocumentType        string `json:"recipient_document_type" yaml:"recipient_document_type"`
	RecipientDocumentNumber      string `json:"recipient_document_number" yaml:"recipient_document_number"`
	RecipientCity                string `json:"recipient_city" yaml:"recipient_city"`
	RecipientState               string `json:"recipient_state" yaml:"recipient_state"`
	RecipientCountry             string `json:"recipient_country" yaml:"recipient_country"`
	RecipientPostalCode          string `json:"recipient_postal_code" yaml:"recipient_postal_code"`
	SenderCountry                string `json:"sender_country" yaml:"sender_country"`
	Amount                       sdk.Coin `json:"amount" yaml:"amount"`
	SourceCurrency               string `json:"source_currency" yaml:"source_currency"`
	DestinationCurrency          string `json:"destination_currency" yaml:"destination_currency"`
	SettlementMethod             string `json:"settlement_method" yaml:"settlement_method"`
	SettlementDetails            string `json:"settlement_details" yaml:"settlement_details"`
	PurposeOfTransfer            string `json:"purpose_of_transfer" yaml:"purpose_of_transfer"`
	Memo                         string `json:"memo,omitempty" yaml:"memo,omitempty"`
	ExpiresInHours               int32  `json:"expires_in_hours,omitempty" yaml:"expires_in_hours,omitempty"`
	PreferredSewaMitraLocation   string `json:"preferred_seva_mitra_location" yaml:"preferred_seva_mitra_location"`
}

// NewMsgInitiateRemittanceWithSewaMitra creates a new MsgInitiateRemittanceWithSewaMitra instance
func NewMsgInitiateRemittanceWithSewaMitra(
	sender, recipientName, recipientPhone, recipientCountry, senderCountry string,
	amount sdk.Coin, sourceCurrency, destinationCurrency, settlementMethod, purposeOfTransfer string,
	preferredSewaMitraLocation string,
) *MsgInitiateRemittanceWithSewaMitra {
	return &MsgInitiateRemittanceWithSewaMitra{
		Sender:                     sender,
		RecipientName:              recipientName,
		RecipientPhone:             recipientPhone,
		RecipientCountry:           recipientCountry,
		SenderCountry:              senderCountry,
		Amount:                     amount,
		SourceCurrency:             sourceCurrency,
		DestinationCurrency:        destinationCurrency,
		SettlementMethod:           settlementMethod,
		PurposeOfTransfer:          purposeOfTransfer,
		PreferredSewaMitraLocation: preferredSewaMitraLocation,
	}
}

// Route returns the message route for routing
func (msg MsgInitiateRemittanceWithSewaMitra) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg MsgInitiateRemittanceWithSewaMitra) Type() string {
	return TypeMsgInitiateRemittanceWithSewaMitra
}

// GetSigners returns the expected signers for the message
func (msg MsgInitiateRemittanceWithSewaMitra) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgInitiateRemittanceWithSewaMitra) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation of the message
func (msg MsgInitiateRemittanceWithSewaMitra) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.RecipientName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "recipient name is required")
	}

	if msg.RecipientCountry == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "recipient country is required")
	}

	if msg.SenderCountry == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "sender country is required")
	}

	if !msg.Amount.IsValid() || !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}

	if msg.SourceCurrency == "" || msg.DestinationCurrency == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "currencies are required")
	}

	if msg.PurposeOfTransfer == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "purpose of transfer is required")
	}

	return nil
}

// AssignedSewaMitraAgent represents details of the assigned Sewa Mitra agent
type AssignedSewaMitraAgent struct {
	AgentId       string   `json:"agent_id" yaml:"agent_id"`
	AgentName     string   `json:"agent_name" yaml:"agent_name"`
	BusinessName  string   `json:"business_name" yaml:"business_name"`
	Phone         string   `json:"phone" yaml:"phone"`
	Address       string   `json:"address" yaml:"address"`
	Commission    sdk.Coin `json:"commission" yaml:"commission"`
}

// MsgInitiateRemittanceWithSewaMitraResponse is the response for MsgInitiateRemittanceWithSewaMitra
type MsgInitiateRemittanceWithSewaMitraResponse struct {
	TransferId               string                  `json:"transfer_id" yaml:"transfer_id"`
	Status                   string                  `json:"status" yaml:"status"`
	ExpiresAt                time.Time               `json:"expires_at" yaml:"expires_at"`
	AssignedSewaMitraAgent   *AssignedSewaMitraAgent `json:"assigned_seva_mitra_agent,omitempty" yaml:"assigned_seva_mitra_agent,omitempty"`
}

// ========== MsgConfirmSewaMitraPickup ==========

// MsgConfirmSewaMitraPickup represents a message to confirm Sewa Mitra pickup
type MsgConfirmSewaMitraPickup struct {
	TransferId       string `json:"transfer_id" yaml:"transfer_id"`
	AgentAddress     string `json:"agent_address,omitempty" yaml:"agent_address,omitempty"`
	ConfirmationCode string `json:"confirmation_code" yaml:"confirmation_code"`
	PickupProof      string `json:"pickup_proof" yaml:"pickup_proof"`
}

// NewMsgConfirmSewaMitraPickup creates a new MsgConfirmSewaMitraPickup instance
func NewMsgConfirmSewaMitraPickup(transferId, agentAddress, confirmationCode, pickupProof string) *MsgConfirmSewaMitraPickup {
	return &MsgConfirmSewaMitraPickup{
		TransferId:       transferId,
		AgentAddress:     agentAddress,
		ConfirmationCode: confirmationCode,
		PickupProof:      pickupProof,
	}
}

// Route returns the message route for routing
func (msg MsgConfirmSewaMitraPickup) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg MsgConfirmSewaMitraPickup) Type() string {
	return TypeMsgConfirmSewaMitraPickup
}

// GetSigners returns the expected signers for the message
func (msg MsgConfirmSewaMitraPickup) GetSigners() []sdk.AccAddress {
	if msg.AgentAddress != "" {
		agent, err := sdk.AccAddressFromBech32(msg.AgentAddress)
		if err != nil {
			panic(err)
		}
		return []sdk.AccAddress{agent}
	}
	// If no agent address, anyone can confirm (recipient or agent)
	return []sdk.AccAddress{}
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgConfirmSewaMitraPickup) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation of the message
func (msg MsgConfirmSewaMitraPickup) ValidateBasic() error {
	if msg.TransferId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "transfer ID is required")
	}

	if msg.AgentAddress != "" {
		if _, err := sdk.AccAddressFromBech32(msg.AgentAddress); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid agent address (%s)", err)
		}
	}

	if msg.ConfirmationCode == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "confirmation code is required")
	}

	return nil
}

// MsgConfirmSewaMitraPickupResponse is the response for MsgConfirmSewaMitraPickup
type MsgConfirmSewaMitraPickupResponse struct {
	Success     bool       `json:"success" yaml:"success"`
	Status      string     `json:"status" yaml:"status"`
	CompletedAt *time.Time `json:"completed_at,omitempty" yaml:"completed_at,omitempty"`
}

// ========== MsgPaySewaMitraCommission ==========

// MsgPaySewaMitraCommission represents a message to pay Sewa Mitra commission
type MsgPaySewaMitraCommission struct {
	Authority    string `json:"authority" yaml:"authority"`
	CommissionId string `json:"commission_id" yaml:"commission_id"`
}

// NewMsgPaySewaMitraCommission creates a new MsgPaySewaMitraCommission instance
func NewMsgPaySewaMitraCommission(authority, commissionId string) *MsgPaySewaMitraCommission {
	return &MsgPaySewaMitraCommission{
		Authority:    authority,
		CommissionId: commissionId,
	}
}

// Route returns the message route for routing
func (msg MsgPaySewaMitraCommission) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg MsgPaySewaMitraCommission) Type() string {
	return TypeMsgPaySewaMitraCommission
}

// GetSigners returns the expected signers for the message
func (msg MsgPaySewaMitraCommission) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgPaySewaMitraCommission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation of the message
func (msg MsgPaySewaMitraCommission) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.CommissionId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "commission ID is required")
	}

	return nil
}

// MsgPaySewaMitraCommissionResponse is the response for MsgPaySewaMitraCommission
type MsgPaySewaMitraCommissionResponse struct {
	Success bool      `json:"success" yaml:"success"`
	PaidAt  time.Time `json:"paid_at" yaml:"paid_at"`
	Amount  sdk.Coin  `json:"amount" yaml:"amount"`
}

// ========== MsgDeactivateSewaMitraAgent ==========

// MsgDeactivateSewaMitraAgent represents a message to deactivate a Sewa Mitra agent
type MsgDeactivateSewaMitraAgent struct {
	Sender  string `json:"sender" yaml:"sender"`
	AgentId string `json:"agent_id" yaml:"agent_id"`
	Reason  string `json:"reason" yaml:"reason"`
}

// NewMsgDeactivateSewaMitraAgent creates a new MsgDeactivateSewaMitraAgent instance
func NewMsgDeactivateSewaMitraAgent(sender, agentId, reason string) *MsgDeactivateSewaMitraAgent {
	return &MsgDeactivateSewaMitraAgent{
		Sender:  sender,
		AgentId: agentId,
		Reason:  reason,
	}
}

// Route returns the message route for routing
func (msg MsgDeactivateSewaMitraAgent) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg MsgDeactivateSewaMitraAgent) Type() string {
	return TypeMsgDeactivateSewaMitraAgent
}

// GetSigners returns the expected signers for the message
func (msg MsgDeactivateSewaMitraAgent) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// GetSignBytes returns the message bytes to sign over
func (msg MsgDeactivateSewaMitraAgent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation of the message
func (msg MsgDeactivateSewaMitraAgent) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if err := ValidateAgentID(msg.AgentId); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if msg.Reason == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "reason is required")
	}

	return nil
}

// MsgDeactivateSewaMitraAgentResponse is the response for MsgDeactivateSewaMitraAgent
type MsgDeactivateSewaMitraAgentResponse struct {
	Success bool `json:"success" yaml:"success"`
}