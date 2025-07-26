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
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/dhansetu/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface for the provided Keeper
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// RegisterDhanPataAddress implements the DhanPata address registration
func (k msgServer) RegisterDhanPataAddress(goCtx context.Context, msg *types.MsgRegisterDhanPataAddress) (*types.MsgRegisterDhanPataAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate sender
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Get registration fee
	params := k.GetParams(ctx)
	registrationFee := sdk.NewCoin(types.DefaultDenom, params.DhanPataRegistrationFee)

	// Charge registration fee
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, senderAddr, types.ModuleName, sdk.NewCoins(registrationFee),
	); err != nil {
		return nil, err
	}

	// Create metadata from message
	metadata := &types.DhanPataMetadata{
		DisplayName:     msg.DisplayName,
		Description:     msg.Description,
		ProfileImageURL: msg.ProfileImageUrl,
		Tags:            msg.Tags,
		Verified:        false, // Initially unverified
	}

	// Register the address
	err = k.RegisterDhanPataAddress(ctx, msg.Name, msg.Sender, msg.AddressType, metadata)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterDhanPataAddressResponse{
		Success: true,
		DhanpataName: msg.Name,
	}, nil
}

// CreateKshetraCoin implements Kshetra coin creation
func (k msgServer) CreateKshetraCoin(goCtx context.Context, msg *types.MsgCreateKshetraCoin) (*types.MsgCreateKshetraCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate creator
	creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	// Get creation fee
	params := k.GetParams(ctx)
	creationFee := sdk.NewCoin(types.DefaultDenom, params.KshetraCoinCreationFee)

	// Charge creation fee
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, creatorAddr, types.ModuleName, sdk.NewCoins(creationFee),
	); err != nil {
		return nil, err
	}

	// Create coin symbol from pincode
	coinSymbol := fmt.Sprintf("%s%s", msg.CoinPrefix, msg.Pincode)

	// Create Kshetra coin
	coin := types.KshetraCoin{
		Pincode:           msg.Pincode,
		CoinName:          msg.CoinName,
		CoinSymbol:        coinSymbol,
		Creator:           msg.Creator,
		TotalSupply:       msg.TotalSupply,
		CirculatingSupply: sdk.ZeroInt(),
		MarketCap:         sdk.ZeroInt(),
		HolderCount:       0,
		CommunityFund:     sdk.ZeroInt(),
		NGOBeneficiary:    msg.NgoBeneficiary,
		Description:       msg.Description,
		LocalLandmarks:    msg.LocalLandmarks,
	}

	err = k.CreateKshetraCoin(ctx, coin)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateKshetraCoinResponse{
		Success:    true,
		CoinSymbol: coinSymbol,
		Pincode:    msg.Pincode,
	}, nil
}

// RegisterEnhancedMitra implements enhanced mitra registration
func (k msgServer) RegisterEnhancedMitra(goCtx context.Context, msg *types.MsgRegisterEnhancedMitra) (*types.MsgRegisterEnhancedMitraResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate sender
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Check minimum trust score requirement
	params := k.GetParams(ctx)
	if msg.TrustScore < params.MinTrustScoreForMitra {
		return nil, types.ErrInsufficientTrustScore
	}

	// Charge registration fee
	registrationFee := sdk.NewCoin(types.DefaultDenom, params.MitraRegistrationFee)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, senderAddr, types.ModuleName, sdk.NewCoins(registrationFee),
	); err != nil {
		return nil, err
	}

	// Create payment methods
	var paymentMethods []types.PaymentMethod
	for _, pm := range msg.PaymentMethods {
		paymentMethods = append(paymentMethods, types.PaymentMethod{
			Type:        pm.Type,
			Provider:    pm.Provider,
			Identifier:  pm.Identifier,
			IsPreferred: pm.IsPreferred,
			IsVerified:  false, // Initially unverified
		})
	}

	// Create enhanced mitra profile
	profile := types.EnhancedMitraProfile{
		MitraId:          msg.MitraId,
		DhanPataName:     msg.DhanpataName,
		MitraType:        msg.MitraType,
		TrustScore:       msg.TrustScore,
		DailyVolume:      sdk.ZeroInt(),
		MonthlyVolume:    sdk.ZeroInt(),
		TotalTrades:      0,
		SuccessfulTrades: 0,
		ActiveEscrows:    []string{},
		Specializations:  msg.Specializations,
		OperatingRegions: msg.OperatingRegions,
		PaymentMethods:   paymentMethods,
		KYCStatus:        "pending",
	}

	err = k.RegisterEnhancedMitra(ctx, profile)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterEnhancedMitraResponse{
		Success: true,
		MitraId: msg.MitraId,
	}, nil
}

// ProcessMoneyOrderWithDhanPata implements money order processing with DhanPata
func (k msgServer) ProcessMoneyOrderWithDhanPata(goCtx context.Context, msg *types.MsgProcessMoneyOrderWithDhanPata) (*types.MsgProcessMoneyOrderWithDhanPataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Process the money order with DhanPata integration
	err := k.ProcessMoneyOrderWithDhanPata(ctx, msg.Sender, msg.ReceiverDhanpata, msg.Amount, msg.Note)
	if err != nil {
		return nil, err
	}

	// Calculate fee
	params := k.GetParams(ctx)
	fee := params.DhanSetuFeeRate.MulInt(msg.Amount.Amount).TruncateInt()
	feeAmount := sdk.NewCoin(msg.Amount.Denom, fee)

	return &types.MsgProcessMoneyOrderWithDhanPataResponse{
		Success:          true,
		ReceiverDhanpata: msg.ReceiverDhanpata,
		ProcessedAmount:  msg.Amount,
		Fee:              feeAmount,
	}, nil
}

// UpdateDhanPataMetadata implements DhanPata metadata updates
func (k msgServer) UpdateDhanPataMetadata(goCtx context.Context, msg *types.MsgUpdateDhanPataMetadata) (*types.MsgUpdateDhanPataMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get existing DhanPata address
	dhanpataAddr, found := k.GetDhanPataAddress(ctx, msg.Name)
	if !found {
		return nil, types.ErrDhanPataNotFound
	}

	// Verify ownership
	if dhanpataAddr.Owner != msg.Owner {
		return nil, types.ErrUnauthorizedOperation
	}

	// Update metadata
	if dhanpataAddr.Metadata == nil {
		dhanpataAddr.Metadata = &types.DhanPataMetadata{}
	}

	if msg.DisplayName != "" {
		dhanpataAddr.Metadata.DisplayName = msg.DisplayName
	}
	if msg.Description != "" {
		dhanpataAddr.Metadata.Description = msg.Description
	}
	if msg.ProfileImageUrl != "" {
		dhanpataAddr.Metadata.ProfileImageURL = msg.ProfileImageUrl
	}
	if len(msg.Tags) > 0 {
		dhanpataAddr.Metadata.Tags = msg.Tags
	}

	dhanpataAddr.UpdatedAt = ctx.BlockTime()

	// Save updated address
	k.SetDhanPataAddress(ctx, dhanpataAddr)

	return &types.MsgUpdateDhanPataMetadataResponse{
		Success: true,
	}, nil
}

// GetParams returns the module parameters
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the module parameters
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}