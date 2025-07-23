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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"deshchain/x/donation/types"
)

var _ types.QueryServer = Keeper{}

// Params queries the module parameters
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// NGOWallet queries an NGO wallet by ID
func (k Keeper) NGOWallet(c context.Context, req *types.QueryNGOWalletRequest) (*types.QueryNGOWalletResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	ngo, found := k.GetNGOWallet(ctx, req.NgoWalletId)
	if !found {
		return nil, status.Error(codes.NotFound, "NGO wallet not found")
	}

	return &types.QueryNGOWalletResponse{NgoWallet: ngo}, nil
}

// NGOWallets queries all NGO wallets
func (k Keeper) NGOWallets(c context.Context, req *types.QueryNGOWalletsRequest) (*types.QueryNGOWalletsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var ngos []types.NGOWallet
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	ngoStore := prefix.NewStore(store, types.NGOWalletKey)

	pageRes, err := query.Paginate(ngoStore, req.Pagination, func(key []byte, value []byte) error {
		var ngo types.NGOWallet
		if err := k.cdc.Unmarshal(value, &ngo); err != nil {
			return err
		}

		ngos = append(ngos, ngo)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryNGOWalletsResponse{NgoWallets: ngos, Pagination: pageRes}, nil
}

// ActiveNGOs queries all active and verified NGOs
func (k Keeper) ActiveNGOs(c context.Context, req *types.QueryActiveNGOsRequest) (*types.QueryActiveNGOsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	activeNGOs := k.GetActiveNGOs(ctx)

	return &types.QueryActiveNGOsResponse{NgoWallets: activeNGOs}, nil
}

// DonationRecord queries a donation record by ID
func (k Keeper) DonationRecord(c context.Context, req *types.QueryDonationRecordRequest) (*types.QueryDonationRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	record, found := k.GetDonationRecord(ctx, req.DonationId)
	if !found {
		return nil, status.Error(codes.NotFound, "donation record not found")
	}

	return &types.QueryDonationRecordResponse{DonationRecord: record}, nil
}

// DonationsByDonor queries all donations by a donor
func (k Keeper) DonationsByDonor(c context.Context, req *types.QueryDonationsByDonorRequest) (*types.QueryDonationsByDonorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	donationIds := k.GetDonationsByDonor(ctx, req.Donor)
	var donations []types.DonationRecord

	// Paginate through donation IDs
	start, end := 0, len(donationIds)
	if req.Pagination != nil {
		limit := req.Pagination.Limit
		if limit > 0 {
			start = int(req.Pagination.Offset)
			end = start + int(limit)
			if end > len(donationIds) {
				end = len(donationIds)
			}
		}
	}

	for i := start; i < end; i++ {
		record, found := k.GetDonationRecord(ctx, donationIds[i])
		if found {
			donations = append(donations, record)
		}
	}

	pageRes := &query.PageResponse{
		Total: uint64(len(donationIds)),
	}

	return &types.QueryDonationsByDonorResponse{
		DonationRecords: donations,
		Pagination:      pageRes,
	}, nil
}

// DonationsByNGO queries all donations received by an NGO
func (k Keeper) DonationsByNGO(c context.Context, req *types.QueryDonationsByNGORequest) (*types.QueryDonationsByNGOResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	donationIds := k.GetDonationsByNGO(ctx, req.NgoWalletId)
	var donations []types.DonationRecord

	// Paginate through donation IDs
	start, end := 0, len(donationIds)
	if req.Pagination != nil {
		limit := req.Pagination.Limit
		if limit > 0 {
			start = int(req.Pagination.Offset)
			end = start + int(limit)
			if end > len(donationIds) {
				end = len(donationIds)
			}
		}
	}

	for i := start; i < end; i++ {
		record, found := k.GetDonationRecord(ctx, donationIds[i])
		if found {
			donations = append(donations, record)
		}
	}

	pageRes := &query.PageResponse{
		Total: uint64(len(donationIds)),
	}

	return &types.QueryDonationsByNGOResponse{
		DonationRecords: donations,
		Pagination:      pageRes,
	}, nil
}

// Statistics queries the donation module statistics
func (k Keeper) Statistics(c context.Context, req *types.QueryStatisticsRequest) (*types.QueryStatisticsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	stats := k.GetModuleStatistics(ctx)

	return &types.QueryStatisticsResponse{Statistics: stats}, nil
}

// Campaign queries a campaign by ID
func (k Keeper) Campaign(c context.Context, req *types.QueryCampaignRequest) (*types.QueryCampaignResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	campaign, found := k.GetCampaign(ctx, req.CampaignId)
	if !found {
		return nil, status.Error(codes.NotFound, "campaign not found")
	}

	return &types.QueryCampaignResponse{Campaign: campaign}, nil
}

// ActiveCampaigns queries all active campaigns
func (k Keeper) ActiveCampaigns(c context.Context, req *types.QueryActiveCampaignsRequest) (*types.QueryActiveCampaignsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	activeCampaigns := k.GetActiveCampaigns(ctx)

	return &types.QueryActiveCampaignsResponse{Campaigns: activeCampaigns}, nil
}

// TransparencyScore queries the transparency score for an NGO
func (k Keeper) TransparencyScore(c context.Context, req *types.QueryTransparencyScoreRequest) (*types.QueryTransparencyScoreResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	score, found := k.GetTransparencyScore(ctx, req.NgoWalletId)
	if !found {
		// Calculate score if not found
		calculatedScore, err := k.CalculateTransparencyScore(ctx, req.NgoWalletId)
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		score, _ = k.GetTransparencyScore(ctx, req.NgoWalletId)
	}

	return &types.QueryTransparencyScoreResponse{TransparencyScore: score}, nil
}