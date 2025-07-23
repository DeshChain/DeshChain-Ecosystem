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

package keeper

import (
	"context"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"deshchain/x/donation/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the donation MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// UpdateParams updates the module parameters
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Check if module is paused
	if k.IsModulePaused(ctx) {
		return nil, types.ErrModuleDisabled
	}
	
	// Validate authority
	if msg.Authority != k.authority {
		return nil, sdk.ErrUnauthorized
	}
	
	// Validate and set params
	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}
	
	k.SetParams(ctx, msg.Params)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateParams,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.Authority),
		),
	)
	
	return &types.MsgUpdateParamsResponse{}, nil
}

// Donate handles donation messages
func (k msgServer) Donate(goCtx context.Context, msg *types.MsgDonate) (*types.MsgDonateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Check if module is paused
	if k.IsModulePaused(ctx) {
		return nil, types.ErrModuleDisabled
	}
	
	// Check if module is enabled
	params := k.GetParams(ctx)
	if !params.Enabled {
		return nil, types.ErrModuleDisabled
	}
	
	donorAddr, err := sdk.AccAddressFromBech32(msg.Donor)
	if err != nil {
		return nil, err
	}
	
	// Transfer funds from donor to module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, donorAddr, types.ModuleName, msg.Amount); err != nil {
		return nil, err
	}
	
	// Create donation record
	donationId, err := k.CreateDonationRecord(ctx, msg.Donor, msg.NgoWalletId, msg.Amount, msg.Purpose, msg.IsAnonymous)
	if err != nil {
		// Refund the amount back to donor
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, donorAddr, msg.Amount)
		return nil, err
	}
	
	// Handle campaign donation if specified
	if msg.CampaignId > 0 {
		if err := k.UpdateCampaignProgress(ctx, msg.CampaignId, msg.Amount); err != nil {
			// Don't fail the donation, just log the error
			ctx.Logger().Error("failed to update campaign progress", "error", err)
		}
	}
	
	// Transfer funds to NGO
	ngo, _ := k.GetNGOWallet(ctx, msg.NgoWalletId)
	ngoAddr, err := sdk.AccAddressFromBech32(ngo.Address)
	if err != nil {
		return nil, err
	}
	
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ngoAddr, msg.Amount); err != nil {
		return nil, err
	}
	
	// Record fund flow
	k.RecordFundFlow(ctx, "donation", msg.Donor, ngo.Address, msg.Amount, msg.Purpose, donationId, "donation")
	
	// Generate receipt hash
	receiptHash := sdk.FormatInvariant("receipt_%d_%s", donationId, ctx.BlockHeader().Time.String())
	
	// Update donation record with receipt
	donation, _ := k.GetDonationRecord(ctx, donationId)
	donation.ReceiptHash = receiptHash
	k.SetDonationRecord(ctx, donation)
	
	// Emit receipt generated event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeReceiptGenerated,
			sdk.NewAttribute(types.AttributeKeyDonationID, sdk.FormatInvariant("%d", donationId)),
			sdk.NewAttribute(types.AttributeKeyReceiptHash, receiptHash),
		),
	)
	
	return &types.MsgDonateResponse{
		DonationId:  donationId,
		ReceiptHash: receiptHash,
	}, nil
}