package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Params returns the module parameters
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// Party returns a specific trade party
func (k Keeper) Party(c context.Context, req *types.QueryPartyRequest) (*types.QueryPartyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	party, found := k.GetTradeParty(ctx, req.PartyId)
	if !found {
		return nil, status.Error(codes.NotFound, "party not found")
	}

	return &types.QueryPartyResponse{Party: &party}, nil
}

// Parties returns all trade parties with pagination
func (k Keeper) Parties(c context.Context, req *types.QueryPartiesRequest) (*types.QueryPartiesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var parties []types.TradeParty
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	partyStore := prefix.NewStore(store, types.TradePartyPrefix)

	pageRes, err := query.Paginate(partyStore, req.Pagination, func(key []byte, value []byte) error {
		var party types.TradeParty
		if err := k.cdc.Unmarshal(value, &party); err != nil {
			return err
		}

		// Filter by party type if specified
		if req.PartyType != "" && party.PartyType != req.PartyType {
			return nil
		}

		// Filter by country if specified
		if req.Country != "" && party.Country != req.Country {
			return nil
		}

		parties = append(parties, party)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPartiesResponse{
		Parties:    parties,
		Pagination: pageRes,
	}, nil
}

// LetterOfCredit returns a specific LC
func (k Keeper) LetterOfCredit(c context.Context, req *types.QueryLetterOfCreditRequest) (*types.QueryLetterOfCreditResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	lc, found := k.GetLetterOfCredit(ctx, req.LcId)
	if !found {
		return nil, status.Error(codes.NotFound, "letter of credit not found")
	}

	return &types.QueryLetterOfCreditResponse{Lc: &lc}, nil
}

// LettersOfCredit returns all LCs with pagination
func (k Keeper) LettersOfCredit(c context.Context, req *types.QueryLettersOfCreditRequest) (*types.QueryLettersOfCreditResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var lcs []types.LetterOfCredit
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	lcStore := prefix.NewStore(store, types.LetterOfCreditPrefix)

	pageRes, err := query.Paginate(lcStore, req.Pagination, func(key []byte, value []byte) error {
		var lc types.LetterOfCredit
		if err := k.cdc.Unmarshal(value, &lc); err != nil {
			return err
		}

		// Filter by status if specified
		if req.Status != "" && lc.Status != req.Status {
			return nil
		}

		// Filter by issuing bank if specified
		if req.IssuingBankId != "" && lc.IssuingBankId != req.IssuingBankId {
			return nil
		}

		lcs = append(lcs, lc)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLettersOfCreditResponse{
		Lcs:        lcs,
		Pagination: pageRes,
	}, nil
}

// LcsByParty returns LCs for a specific party
func (k Keeper) LcsByParty(c context.Context, req *types.QueryLcsByPartyRequest) (*types.QueryLcsByPartyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	lcs := k.GetLcsByParty(ctx, req.PartyId)

	// Filter by role if specified
	if req.Role != "" {
		var filteredLcs []types.LetterOfCredit
		for _, lc := range lcs {
			switch req.Role {
			case "applicant":
				if lc.ApplicantId == req.PartyId {
					filteredLcs = append(filteredLcs, lc)
				}
			case "beneficiary":
				if lc.BeneficiaryId == req.PartyId {
					filteredLcs = append(filteredLcs, lc)
				}
			case "bank":
				if lc.IssuingBankId == req.PartyId || lc.AdvisingBankId == req.PartyId {
					filteredLcs = append(filteredLcs, lc)
				}
			}
		}
		lcs = filteredLcs
	}

	return &types.QueryLcsByPartyResponse{
		Lcs: lcs,
	}, nil
}

// Documents returns documents for an LC
func (k Keeper) Documents(c context.Context, req *types.QueryDocumentsRequest) (*types.QueryDocumentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	documents := k.GetDocumentsByLc(ctx, req.LcId)

	// Filter by status if specified
	if req.Status != "" {
		var filteredDocs []types.TradeDocument
		for _, doc := range documents {
			if doc.Status == req.Status {
				filteredDocs = append(filteredDocs, doc)
			}
		}
		documents = filteredDocs
	}

	return &types.QueryDocumentsResponse{
		Documents: documents,
	}, nil
}

// InsurancePolicies returns insurance policies for an LC
func (k Keeper) InsurancePolicies(c context.Context, req *types.QueryInsurancePoliciesRequest) (*types.QueryInsurancePoliciesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	policies := k.GetInsurancePoliciesByLc(ctx, req.LcId)

	return &types.QueryInsurancePoliciesResponse{
		Policies: policies,
	}, nil
}

// ShipmentTracking returns shipment tracking for an LC
func (k Keeper) ShipmentTracking(c context.Context, req *types.QueryShipmentTrackingRequest) (*types.QueryShipmentTrackingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	tracking, found := k.GetShipmentByLc(ctx, req.LcId)
	if !found {
		return &types.QueryShipmentTrackingResponse{}, nil
	}

	return &types.QueryShipmentTrackingResponse{
		Tracking: &tracking,
	}, nil
}

// PaymentInstructions returns payment instructions for an LC
func (k Keeper) PaymentInstructions(c context.Context, req *types.QueryPaymentInstructionsRequest) (*types.QueryPaymentInstructionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	instructions := k.GetPaymentsByLc(ctx, req.LcId)

	// Filter by status if specified
	if req.Status != "" {
		var filteredInstructions []types.PaymentInstruction
		for _, instruction := range instructions {
			if instruction.Status == req.Status {
				filteredInstructions = append(filteredInstructions, instruction)
			}
		}
		instructions = filteredInstructions
	}

	return &types.QueryPaymentInstructionsResponse{
		Instructions: instructions,
	}, nil
}

// TradeFinanceStats returns module statistics
func (k Keeper) TradeFinanceStats(c context.Context, req *types.QueryTradeFinanceStatsRequest) (*types.QueryTradeFinanceStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	stats := k.GetTradeFinanceStats(ctx)

	return &types.QueryTradeFinanceStatsResponse{
		Stats: &stats,
	}, nil
}

// EstimateFees estimates fees for LC issuance
func (k Keeper) EstimateFees(c context.Context, req *types.QueryEstimateFeesRequest) (*types.QueryEstimateFeesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// Parse amount
	lcAmount, err := sdk.ParseCoinNormalized(req.LcAmount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid LC amount")
	}

	issuanceFee, documentFees, insuranceFee, totalFees, processingTimeHours := k.EstimateFees(
		ctx,
		lcAmount,
		req.PaymentTerms,
		req.InsuranceRequired,
	)

	return &types.QueryEstimateFeesResponse{
		IssuanceFee:          issuanceFee.String(),
		DocumentFees:         documentFees.String(),
		InsuranceFee:         insuranceFee.String(),
		TotalFees:           totalFees.String(),
		ProcessingTimeHours: sdk.NewInt(int64(processingTimeHours)).String(),
	}, nil
}