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

package grampension

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/gramsuraksha/types"
)

// NewHandler creates an sdk.Handler for all the gram pension type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateScheme:
			return handleMsgCreateScheme(ctx, k, msg)
		case *types.MsgUpdateScheme:
			return handleMsgUpdateScheme(ctx, k, msg)
		case *types.MsgEnrollParticipant:
			return handleMsgEnrollParticipant(ctx, k, msg)
		case *types.MsgMakeContribution:
			return handleMsgMakeContribution(ctx, k, msg)
		case *types.MsgProcessMaturity:
			return handleMsgProcessMaturity(ctx, k, msg)
		case *types.MsgRequestWithdrawal:
			return handleMsgRequestWithdrawal(ctx, k, msg)
		case *types.MsgProcessWithdrawal:
			return handleMsgProcessWithdrawal(ctx, k, msg)
		case *types.MsgUpdateKYCStatus:
			return handleMsgUpdateKYCStatus(ctx, k, msg)
		case *types.MsgClaimReferral:
			return handleMsgClaimReferral(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateScheme(ctx sdk.Context, k keeper.Keeper, msg *types.MsgCreateScheme) (*sdk.Result, error) {
	// Generate unique scheme ID
	schemeID := k.GenerateSchemeID(ctx, msg.SchemeName)

	// Create scheme object
	scheme := types.SurakshaScheme{
		SchemeID:               schemeID,
		SchemeName:             msg.SchemeName,
		Description:            msg.Description,
		MonthlyContribution:    msg.MonthlyContribution,
		ContributionPeriod:     msg.ContributionPeriod,
		MaturityBonus:          msg.MaturityBonus,
		MinAge:                 msg.MinAge,
		MaxAge:                 msg.MaxAge,
		MaxParticipants:        msg.MaxParticipants,
		GracePeriodDays:        msg.GracePeriodDays,
		EarlyWithdrawalPenalty: msg.EarlyWithdrawalPenalty,
		LatePaymentPenalty:     msg.LatePaymentPenalty,
		ReferralRewardPercent:  msg.ReferralRewardPercent,
		OnTimeBonusPercent:     msg.OnTimeBonusPercent,
		KYCRequired:            msg.KYCRequired,
		LiquidityProvision:     msg.LiquidityProvision,
		LiquidityUtilization:   msg.LiquidityUtilization,
		CulturalIntegration:    true, // Always enable cultural integration
	}

	// Create the scheme
	if err := k.CreateScheme(ctx, scheme); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleMsgUpdateScheme(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateScheme) (*sdk.Result, error) {
	// Only authority can update schemes
	if msg.Authority.String() != k.GetAuthority() {
		return nil, types.ErrUnauthorized
	}

	if err := k.UpdateScheme(ctx, msg.SchemeID, *msg); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleMsgEnrollParticipant(ctx sdk.Context, k keeper.Keeper, msg *types.MsgEnrollParticipant) (*sdk.Result, error) {
	// Generate participant ID
	participantID := k.GenerateParticipantID(ctx, msg.Participant, msg.SchemeID)

	// Create participant object
	participant := types.SurakshaParticipant{
		ParticipantID:     participantID,
		SchemeID:          msg.SchemeID,
		Address:           msg.Participant,
		Name:              msg.Name,
		Age:               msg.Age,
		VillagePostalCode: msg.VillagePostalCode,
		ReferrerAddress:   msg.ReferrerAddress,
	}

	// Enroll the participant
	if err := k.EnrollParticipant(ctx, participant); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleMsgMakeContribution(ctx sdk.Context, k keeper.Keeper, msg *types.MsgMakeContribution) (*sdk.Result, error) {
	if err := k.MakeContribution(ctx, msg.ParticipantID, msg.Amount, msg.Month); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleMsgProcessMaturity(ctx sdk.Context, k keeper.Keeper, msg *types.MsgProcessMaturity) (*sdk.Result, error) {
	// Only authority can process maturities manually
	if msg.Authority.String() != k.GetAuthority() {
		return nil, types.ErrUnauthorized
	}

	if err := k.ProcessMaturity(ctx, msg.ParticipantID); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}

func handleMsgRequestWithdrawal(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRequestWithdrawal) (*sdk.Result, error) {
	// Request withdrawal logic would be implemented here
	// For now, return not implemented
	return nil, sdkerrors.Wrap(sdkerrors.ErrNotSupported, "withdrawal requests not yet implemented")
}

func handleMsgProcessWithdrawal(ctx sdk.Context, k keeper.Keeper, msg *types.MsgProcessWithdrawal) (*sdk.Result, error) {
	// Process withdrawal logic would be implemented here
	// For now, return not implemented
	return nil, sdkerrors.Wrap(sdkerrors.ErrNotSupported, "withdrawal processing not yet implemented")
}

func handleMsgUpdateKYCStatus(ctx sdk.Context, k keeper.Keeper, msg *types.MsgUpdateKYCStatus) (*sdk.Result, error) {
	// Only authority can update KYC status
	if msg.Authority.String() != k.GetAuthority() {
		return nil, types.ErrUnauthorized
	}

	// KYC update logic would be implemented here
	// For now, return not implemented
	return nil, sdkerrors.Wrap(sdkerrors.ErrNotSupported, "KYC updates not yet implemented")
}

func handleMsgClaimReferral(ctx sdk.Context, k keeper.Keeper, msg *types.MsgClaimReferral) (*sdk.Result, error) {
	// Referral claim logic would be implemented here
	// For now, return not implemented
	return nil, sdkerrors.Wrap(sdkerrors.ErrNotSupported, "referral claims not yet implemented")
}