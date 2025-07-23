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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/deshchain/namo/x/governance/types"
)

var _ types.QueryServer = Keeper{}

// GovernanceInfo returns the current governance information
func (k Keeper) GovernanceInfo(goCtx context.Context, req *types.QueryGovernanceInfoRequest) (*types.QueryGovernanceInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryGovernanceInfoResponse{
		FounderAddress: k.GetFounderAddress(ctx),
		CurrentPhase:   k.GetGovernancePhase(ctx),
		GenesisTime:    k.GetGenesisTime(ctx),
	}, nil
}

// ProtectedParameter returns a specific protected parameter
func (k Keeper) ProtectedParameter(goCtx context.Context, req *types.QueryProtectedParameterRequest) (*types.QueryProtectedParameterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "parameter name cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	param, found := k.GetProtectedParameter(ctx, req.Name)
	if !found {
		return nil, status.Errorf(codes.NotFound, "parameter %s not found", req.Name)
	}

	return &types.QueryProtectedParameterResponse{
		Parameter: param,
	}, nil
}

// ProtectedParameters returns all protected parameters
func (k Keeper) ProtectedParameters(goCtx context.Context, req *types.QueryProtectedParametersRequest) (*types.QueryProtectedParametersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetAllProtectedParameters(ctx)

	return &types.QueryProtectedParametersResponse{
		Parameters: params,
	}, nil
}

// ProposalVetoStatus returns the veto status of a proposal
func (k Keeper) ProposalVetoStatus(goCtx context.Context, req *types.QueryProposalVetoStatusRequest) (*types.QueryProposalVetoStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	isVetoed := k.IsProposalVetoed(ctx, req.ProposalId)

	return &types.QueryProposalVetoStatusResponse{
		IsVetoed: isVetoed,
	}, nil
}

// VetoedProposals returns all vetoed proposal IDs
func (k Keeper) VetoedProposals(goCtx context.Context, req *types.QueryVetoedProposalsRequest) (*types.QueryVetoedProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	proposalIds := k.GetAllVetoedProposals(ctx)

	return &types.QueryVetoedProposalsResponse{
		ProposalIds: proposalIds,
	}, nil
}

// ProposalFounderApproval returns the founder approval status of a proposal
func (k Keeper) ProposalFounderApproval(goCtx context.Context, req *types.QueryProposalFounderApprovalRequest) (*types.QueryProposalFounderApprovalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	hasApproval := store.Has(types.GetFounderApprovalKey(req.ProposalId))

	return &types.QueryProposalFounderApprovalResponse{
		HasApproval: hasApproval,
	}, nil
}