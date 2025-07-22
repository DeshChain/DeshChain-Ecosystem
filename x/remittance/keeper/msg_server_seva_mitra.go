package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/deshchain/deshchain/x/remittance/types"
)

// MsgServer implements the remittance MsgServer interface for Sewa Mitra operations
type sewaMitraMsgServer struct {
	Keeper
}

// NewSewaMitraMsgServer returns an implementation of the remittance MsgServer interface for Sewa Mitra operations
func NewSewaMitraMsgServer(keeper Keeper) types.SewaMitraMsgServer {
	return &sewaMitraMsgServer{Keeper: keeper}
}

// RegisterSewaMitraAgent implements the MsgRegisterSewaMitraAgent message handler
func (k sewaMitraMsgServer) RegisterSewaMitraAgent(goCtx context.Context, msg *types.MsgRegisterSewaMitraAgent) (*types.MsgRegisterSewaMitraAgentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the authority (only authorized addresses can register agents)
	if msg.Authority != k.GetAuthority() {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	// Validate agent address
	agentAddr, err := sdk.AccAddressFromBech32(msg.AgentAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid agent address")
	}

	// Check if agent already exists
	if k.HasSewaMitraAgent(ctx, msg.AgentId) {
		return nil, sdkerrors.Wrap(types.ErrAgentAlreadyExists, msg.AgentId)
	}

	// Parse commission rates
	baseCommissionRate, err := sdk.NewDecFromStr(msg.BaseCommissionRate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid base commission rate")
	}

	volumeBonusRate, err := sdk.NewDecFromStr(msg.VolumeBonusRate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid volume bonus rate")
	}

	// Create new agent
	agent := types.NewSewaMitraAgent(
		msg.AgentId,
		msg.AgentAddress,
		msg.AgentName,
		msg.BusinessName,
		msg.Country,
		msg.State,
		msg.City,
		msg.PostalCode,
		msg.AddressLine1,
		msg.AddressLine2,
		msg.Phone,
		msg.Email,
		msg.SupportedCurrencies,
		msg.SupportedMethods,
		msg.LiquidityLimit,
		msg.DailyLimit,
		baseCommissionRate,
		volumeBonusRate,
		msg.MinimumCommission,
		msg.MaximumCommission,
	)

	// Set creation timestamps
	agent.CreatedAt = ctx.BlockTime()
	agent.UpdatedAt = ctx.BlockTime()

	// Register the agent
	if err := k.RegisterSewaMitraAgent(ctx, agent); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to register sewa mitra agent")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"sewa_mitra_agent_registered",
			sdk.NewAttribute(types.AttributeKeySewaMitraAgentID, msg.AgentId),
			sdk.NewAttribute(types.AttributeKeyAgentName, msg.AgentName),
			sdk.NewAttribute(types.AttributeKeyAgentLocation, fmt.Sprintf("%s, %s", msg.City, msg.Country)),
			sdk.NewAttribute("business_name", msg.BusinessName),
		),
	)

	return &types.MsgRegisterSewaMitraAgentResponse{
		AgentId: msg.AgentId,
		Status:  "pending_verification",
	}, nil
}

// UpdateSewaMitraAgent implements the MsgUpdateSewaMitraAgent message handler
func (k sewaMitraMsgServer) UpdateSewaMitraAgent(goCtx context.Context, msg *types.MsgUpdateSewaMitraAgent) (*types.MsgUpdateSewaMitraAgentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get existing agent
	agent, err := k.GetSewaMitraAgent(ctx, msg.AgentId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "agent not found")
	}

	// Verify the sender is the agent or authority
	if msg.Sender != agent.AgentAddress && msg.Sender != k.GetAuthority() {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "sender must be agent or authority")
	}

	// Update agent information
	if msg.Phone != "" {
		agent.Phone = msg.Phone
	}
	if msg.Email != "" {
		agent.Email = msg.Email
	}
	if len(msg.SupportedCurrencies) > 0 {
		agent.SupportedCurrencies = msg.SupportedCurrencies
	}
	if len(msg.SupportedMethods) > 0 {
		agent.SupportedMethods = msg.SupportedMethods
	}
	if msg.LiquidityLimit.IsPositive() {
		agent.LiquidityLimit = msg.LiquidityLimit
	}
	if msg.DailyLimit.IsPositive() {
		agent.DailyLimit = msg.DailyLimit
	}

	agent.UpdatedAt = ctx.BlockTime()

	// Save updated agent
	if err := k.SetSewaMitraAgent(ctx, agent); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to update agent")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"sewa_mitra_agent_updated",
			sdk.NewAttribute(types.AttributeKeySewaMitraAgentID, msg.AgentId),
			sdk.NewAttribute("updated_by", msg.Sender),
		),
	)

	return &types.MsgUpdateSewaMitraAgentResponse{
		Success: true,
	}, nil
}

// ActivateSewaMitraAgent implements the MsgActivateSewaMitraAgent message handler
func (k sewaMitraMsgServer) ActivateSewaMitraAgent(goCtx context.Context, msg *types.MsgActivateSewaMitraAgent) (*types.MsgActivateSewaMitraAgentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate authority (only authorized addresses can activate agents)
	if msg.Authority != k.GetAuthority() {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	// Activate the agent
	if err := k.ActivateSewaMitraAgent(ctx, msg.AgentId); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to activate agent")
	}

	return &types.MsgActivateSewaMitraAgentResponse{
		Success: true,
	}, nil
}

// SuspendSewaMitraAgent implements the MsgSuspendSewaMitraAgent message handler
func (k sewaMitraMsgServer) SuspendSewaMitraAgent(goCtx context.Context, msg *types.MsgSuspendSewaMitraAgent) (*types.MsgSuspendSewaMitraAgentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate authority (only authorized addresses can suspend agents)
	if msg.Authority != k.GetAuthority() {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	// Parse suspension duration
	suspendedUntil := ctx.BlockTime().Add(time.Duration(msg.SuspensionDurationDays) * 24 * time.Hour)

	// Suspend the agent
	if err := k.SuspendSewaMitraAgent(ctx, msg.AgentId, suspendedUntil, msg.Reason); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to suspend agent")
	}

	return &types.MsgSuspendSewaMitraAgentResponse{
		Success:        true,
		SuspendedUntil: suspendedUntil,
	}, nil
}

// InitiateRemittanceWithSewaMitra implements the MsgInitiateRemittanceWithSewaMitra message handler
func (k sewaMitraMsgServer) InitiateRemittanceWithSewaMitra(goCtx context.Context, msg *types.MsgInitiateRemittanceWithSewaMitra) (*types.MsgInitiateRemittanceWithSewaMitraResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate sender address
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	// Validate recipient address if provided
	if msg.RecipientAddress != "" {
		if _, err := sdk.AccAddressFromBech32(msg.RecipientAddress); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid recipient address")
		}
	}

	// Parse expiry time
	expiresAt := ctx.BlockTime().Add(24 * time.Hour) // Default 24 hour expiry
	if msg.ExpiresInHours > 0 {
		expiresAt = ctx.BlockTime().Add(time.Duration(msg.ExpiresInHours) * time.Hour)
	}

	// Create recipient info
	recipientInfo := types.RecipientInfo{
		Name:           msg.RecipientName,
		Phone:          msg.RecipientPhone,
		DocumentType:   msg.RecipientDocumentType,
		DocumentNumber: msg.RecipientDocumentNumber,
		Address: types.RecipientAddress{
			Line1:      msg.RecipientAddress,
			City:       msg.RecipientCity,
			State:      msg.RecipientState,
			Country:    msg.RecipientCountry,
			PostalCode: msg.RecipientPostalCode,
		},
	}

	// Create settlement details for Sewa Mitra
	settlementDetails := types.SettlementDetails{
		Method:      msg.SettlementMethod,
		Details:     msg.SettlementDetails,
		BankDetails: nil, // Not used for Sewa Mitra cash pickup
		WalletDetails: &types.WalletDetails{
			WalletType: "seva_mitra",
			WalletId:   msg.PreferredSewaMitraLocation,
		},
	}

	// Initiate transfer with Sewa Mitra preference
	transferID, err := k.InitiateTransfer(
		ctx,
		msg.Sender,
		msg.RecipientAddress,
		msg.SenderCountry,
		msg.RecipientCountry,
		msg.Amount,
		msg.SourceCurrency,
		msg.DestinationCurrency,
		msg.SettlementMethod,
		msg.PurposeOfTransfer,
		msg.Memo,
		expiresAt,
		recipientInfo,
		settlementDetails,
		true, // prefer Sewa Mitra
		msg.PreferredSewaMitraLocation,
	)

	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to initiate remittance transfer")
	}

	// Get the created transfer to return agent info
	transfer, err := k.GetRemittanceTransfer(ctx, transferID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to retrieve created transfer")
	}

	response := &types.MsgInitiateRemittanceWithSewaMitraResponse{
		TransferId: transferID,
		Status:     transfer.Status.String(),
		ExpiresAt:  expiresAt,
	}

	// Include Sewa Mitra agent info if one was assigned
	if transfer.UsesSewaMitra {
		agent, err := k.GetSewaMitraAgent(ctx, transfer.SewaMitraAgentId)
		if err == nil {
			response.AssignedSewaMitraAgent = &types.AssignedSewaMitraAgent{
				AgentId:       agent.AgentId,
				AgentName:     agent.AgentName,
				BusinessName:  agent.BusinessName,
				Phone:         agent.Phone,
				Address:       fmt.Sprintf("%s, %s, %s", agent.AddressLine1, agent.City, agent.Country),
				Commission:    transfer.SewaMitraCommission,
			}
		}
	}

	return response, nil
}

// ConfirmSewaMitraPickup implements the MsgConfirmSewaMitraPickup message handler
func (k sewaMitraMsgServer) ConfirmSewaMitraPickup(goCtx context.Context, msg *types.MsgConfirmSewaMitraPickup) (*types.MsgConfirmSewaMitraPickupResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the transfer
	transfer, err := k.GetRemittanceTransfer(ctx, msg.TransferId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "transfer not found")
	}

	// Verify this is a Sewa Mitra transfer
	if !transfer.UsesSewaMitra {
		return nil, sdkerrors.Wrap(types.ErrInvalidTransferType, "transfer does not use Sewa Mitra")
	}

	// Verify the sender is the assigned agent
	if msg.AgentAddress != "" {
		agent, err := k.GetSewaMitraAgent(ctx, transfer.SewaMitraAgentId)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "assigned agent not found")
		}
		if msg.AgentAddress != agent.AgentAddress {
			return nil, sdkerrors.Wrap(types.ErrUnauthorized, "sender is not the assigned agent")
		}
	}

	// Confirm the transfer
	err = k.ConfirmTransfer(
		ctx,
		transfer.RecipientAddress,
		msg.TransferId,
		msg.ConfirmationCode,
		msg.PickupProof,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to confirm transfer")
	}

	// Get updated transfer
	updatedTransfer, err := k.GetRemittanceTransfer(ctx, msg.TransferId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to retrieve updated transfer")
	}

	return &types.MsgConfirmSewaMitraPickupResponse{
		Success:     true,
		Status:      updatedTransfer.Status.String(),
		CompletedAt: updatedTransfer.CompletedAt,
	}, nil
}

// PaySewaMitraCommission implements the MsgPaySewaMitraCommission message handler
func (k sewaMitraMsgServer) PaySewaMitraCommission(goCtx context.Context, msg *types.MsgPaySewaMitraCommission) (*types.MsgPaySewaMitraCommissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate authority (only authorized addresses can pay commissions)
	if msg.Authority != k.GetAuthority() {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	// Get the commission record
	commission, err := k.GetSewaMitraCommission(ctx, msg.CommissionId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "commission not found")
	}

	// Check if already paid
	if commission.IsPaid() {
		return nil, sdkerrors.Wrap(types.ErrCommissionAlreadyPaid, msg.CommissionId)
	}

	// Get the agent
	agent, err := k.GetSewaMitraAgent(ctx, commission.AgentId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "agent not found")
	}

	// Parse agent address
	agentAddr, err := sdk.AccAddressFromBech32(agent.AgentAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid agent address")
	}

	// Transfer commission from module to agent
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		agentAddr,
		sdk.NewCoins(commission.TotalCommission),
	); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to transfer commission")
	}

	// Update commission status
	commission.Status = types.COMMISSION_STATUS_PAID
	commission.PaidAt = &ctx.BlockTime()

	if err := k.SetSewaMitraCommission(ctx, commission); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to update commission")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"sewa_mitra_commission_paid",
			sdk.NewAttribute("commission_id", msg.CommissionId),
			sdk.NewAttribute(types.AttributeKeySewaMitraAgentID, commission.AgentId),
			sdk.NewAttribute(types.AttributeKeyCommissionAmount, commission.TotalCommission.String()),
			sdk.NewAttribute("paid_to", agent.AgentAddress),
		),
	)

	return &types.MsgPaySewaMitraCommissionResponse{
		Success:   true,
		PaidAt:    ctx.BlockTime(),
		Amount:    commission.TotalCommission,
	}, nil
}

// DeactivateSewaMitraAgent implements the MsgDeactivateSewaMitraAgent message handler
func (k sewaMitraMsgServer) DeactivateSewaMitraAgent(goCtx context.Context, msg *types.MsgDeactivateSewaMitraAgent) (*types.MsgDeactivateSewaMitraAgentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the agent
	agent, err := k.GetSewaMitraAgent(ctx, msg.AgentId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "agent not found")
	}

	// Verify the sender is the agent or authority
	if msg.Sender != agent.AgentAddress && msg.Sender != k.GetAuthority() {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "sender must be agent or authority")
	}

	// Remove old status index
	k.removeSewaMitraAgentIndexes(ctx, agent)

	// Update agent status
	agent.Status = types.AGENT_STATUS_DEACTIVATED
	agent.UpdatedAt = ctx.BlockTime()

	// Add reason to metadata
	if agent.Metadata == nil {
		agent.Metadata = make(map[string]string)
	}
	agent.Metadata["deactivation_reason"] = msg.Reason
	agent.Metadata["deactivated_by"] = msg.Sender

	// Store updated agent with new indexes
	if err := k.SetSewaMitraAgent(ctx, agent); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to deactivate agent")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"sewa_mitra_agent_deactivated",
			sdk.NewAttribute(types.AttributeKeySewaMitraAgentID, msg.AgentId),
			sdk.NewAttribute("reason", msg.Reason),
			sdk.NewAttribute("deactivated_by", msg.Sender),
		),
	)

	return &types.MsgDeactivateSewaMitraAgentResponse{
		Success: true,
	}, nil
}