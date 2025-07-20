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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/deshchain/deshchain/x/cultural/types"
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

// Quote returns a quote by ID
func (k Keeper) Quote(c context.Context, req *types.QueryQuoteRequest) (*types.QueryQuoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	quote, found := k.GetQuote(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "quote not found")
	}

	return &types.QueryQuoteResponse{Quote: quote}, nil
}

// Quotes returns all quotes with pagination
func (k Keeper) Quotes(c context.Context, req *types.QueryQuotesRequest) (*types.QueryQuotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	quotes := k.GetAllQuotes(ctx)

	// Apply pagination
	pageRes, err := query.Paginate(len(quotes), req.Pagination, func(offset, limit int) error {
		end := offset + limit
		if end > len(quotes) {
			end = len(quotes)
		}
		quotes = quotes[offset:end]
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryQuotesResponse{
		Quotes:     quotes,
		Pagination: pageRes,
	}, nil
}

// QuotesByCategory returns quotes filtered by category
func (k Keeper) QuotesByCategory(c context.Context, req *types.QueryQuotesByCategoryRequest) (*types.QueryQuotesByCategoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	quotes := k.GetQuotesByCategory(ctx, req.Category)

	// Apply pagination
	pageRes, err := query.Paginate(len(quotes), req.Pagination, func(offset, limit int) error {
		end := offset + limit
		if end > len(quotes) {
			end = len(quotes)
		}
		quotes = quotes[offset:end]
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryQuotesByCategoryResponse{
		Quotes:     quotes,
		Pagination: pageRes,
	}, nil
}

// QuotesByAuthor returns quotes filtered by author
func (k Keeper) QuotesByAuthor(c context.Context, req *types.QueryQuotesByAuthorRequest) (*types.QueryQuotesByAuthorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	quotes := k.GetQuotesByAuthor(ctx, req.Author)

	// Apply pagination
	pageRes, err := query.Paginate(len(quotes), req.Pagination, func(offset, limit int) error {
		end := offset + limit
		if end > len(quotes) {
			end = len(quotes)
		}
		quotes = quotes[offset:end]
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryQuotesByAuthorResponse{
		Quotes:     quotes,
		Pagination: pageRes,
	}, nil
}

// QuotesByRegion returns quotes filtered by region
func (k Keeper) QuotesByRegion(c context.Context, req *types.QueryQuotesByRegionRequest) (*types.QueryQuotesByRegionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	quotes := k.GetQuotesByRegion(ctx, req.Region)

	// Apply pagination
	pageRes, err := query.Paginate(len(quotes), req.Pagination, func(offset, limit int) error {
		end := offset + limit
		if end > len(quotes) {
			end = len(quotes)
		}
		quotes = quotes[offset:end]
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryQuotesByRegionResponse{
		Quotes:     quotes,
		Pagination: pageRes,
	}, nil
}

// DailyQuote returns the quote of the day
func (k Keeper) DailyQuote(c context.Context, req *types.QueryDailyQuoteRequest) (*types.QueryDailyQuoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	quote, err := k.GetDailyQuote(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDailyQuoteResponse{Quote: quote}, nil
}

// PopularQuotes returns the most popular quotes
func (k Keeper) PopularQuotes(c context.Context, req *types.QueryPopularQuotesRequest) (*types.QueryPopularQuotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	limit := int(req.Limit)
	if limit <= 0 || limit > 100 {
		limit = 10 // Default to 10 quotes
	}

	quotes := k.GetPopularQuotes(ctx, limit)

	return &types.QueryPopularQuotesResponse{Quotes: quotes}, nil
}

// TransactionQuote returns the quote associated with a transaction
func (k Keeper) TransactionQuote(c context.Context, req *types.QueryTransactionQuoteRequest) (*types.QueryTransactionQuoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	txQuote, found := k.GetTransactionQuote(ctx, req.TxHash)
	if !found {
		return nil, status.Error(codes.NotFound, "transaction quote not found")
	}

	quote, found := k.GetQuote(ctx, txQuote.QuoteId)
	if !found {
		return nil, status.Error(codes.NotFound, "associated quote not found")
	}

	return &types.QueryTransactionQuoteResponse{
		TransactionQuote: txQuote,
		Quote:            quote,
	}, nil
}

// HistoricalEvent returns a historical event by ID
func (k Keeper) HistoricalEvent(c context.Context, req *types.QueryHistoricalEventRequest) (*types.QueryHistoricalEventResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	event, found := k.GetHistoricalEvent(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "historical event not found")
	}

	return &types.QueryHistoricalEventResponse{Event: event}, nil
}

// HistoricalEvents returns all historical events with pagination
func (k Keeper) HistoricalEvents(c context.Context, req *types.QueryHistoricalEventsRequest) (*types.QueryHistoricalEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	events := k.GetAllHistoricalEvents(ctx)

	// Apply pagination
	pageRes, err := query.Paginate(len(events), req.Pagination, func(offset, limit int) error {
		end := offset + limit
		if end > len(events) {
			end = len(events)
		}
		events = events[offset:end]
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryHistoricalEventsResponse{
		Events:     events,
		Pagination: pageRes,
	}, nil
}

// CulturalWisdom returns cultural wisdom by ID
func (k Keeper) CulturalWisdom(c context.Context, req *types.QueryCulturalWisdomRequest) (*types.QueryCulturalWisdomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	wisdom, found := k.GetCulturalWisdom(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "cultural wisdom not found")
	}

	return &types.QueryCulturalWisdomResponse{Wisdom: wisdom}, nil
}

// CulturalWisdomList returns all cultural wisdom with pagination
func (k Keeper) CulturalWisdomList(c context.Context, req *types.QueryCulturalWisdomListRequest) (*types.QueryCulturalWisdomListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	wisdom := k.GetAllCulturalWisdom(ctx)

	// Apply pagination
	pageRes, err := query.Paginate(len(wisdom), req.Pagination, func(offset, limit int) error {
		end := offset + limit
		if end > len(wisdom) {
			end = len(wisdom)
		}
		wisdom = wisdom[offset:end]
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCulturalWisdomListResponse{
		Wisdom:     wisdom,
		Pagination: pageRes,
	}, nil
}

// ActiveFestivals queries currently active festivals
func (k Keeper) ActiveFestivals(goCtx context.Context, req *types.QueryActiveFestivalsRequest) (*types.QueryActiveFestivalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get active festivals from festival manager
	activeFestivals := k.GetActiveFestivals(ctx)
	
	// Convert to response format
	festivalInfos := make([]*types.FestivalInfo, len(activeFestivals))
	for i, festival := range activeFestivals {
		festivalInfos[i] = &types.FestivalInfo{
			Id:                  festival.ID,
			Name:                festival.Name,
			Description:         festival.Description,
			Date:                festival.Date.Format("2006-01-02T15:04:05Z07:00"),
			Type:                festival.Type,
			Region:              festival.Region,
			BonusRate:           festival.BonusRate,
			IsActive:            festival.IsActive,
			TraditionalGreeting: festival.TraditionalGreeting,
			CulturalTheme:       festival.CulturalTheme,
			DaysRemaining:       int32(festival.DaysRemaining),
			Significance:        festival.Significance,
		}
		
		if festival.EndDate != nil {
			endDate := festival.EndDate.Format("2006-01-02T15:04:05Z07:00")
			festivalInfos[i].EndDate = &endDate
		}
	}

	return &types.QueryActiveFestivalsResponse{
		Festivals: festivalInfos,
	}, nil
}

// UpcomingFestivals queries upcoming festivals
func (k Keeper) UpcomingFestivals(goCtx context.Context, req *types.QueryUpcomingFestivalsRequest) (*types.QueryUpcomingFestivalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get upcoming festivals from festival manager
	upcomingFestivals := k.GetUpcomingFestivals(ctx)
	
	// Convert to response format
	festivalInfos := make([]*types.FestivalInfo, len(upcomingFestivals))
	for i, festival := range upcomingFestivals {
		festivalInfos[i] = &types.FestivalInfo{
			Id:                  festival.ID,
			Name:                festival.Name,
			Description:         festival.Description,
			Date:                festival.Date.Format("2006-01-02T15:04:05Z07:00"),
			Type:                festival.Type,
			Region:              festival.Region,
			BonusRate:           festival.BonusRate,
			IsActive:            festival.IsActive,
			TraditionalGreeting: festival.TraditionalGreeting,
			CulturalTheme:       festival.CulturalTheme,
			DaysRemaining:       int32(festival.DaysRemaining),
			Significance:        festival.Significance,
		}
		
		if festival.EndDate != nil {
			endDate := festival.EndDate.Format("2006-01-02T15:04:05Z07:00")
			festivalInfos[i].EndDate = &endDate
		}
	}

	return &types.QueryUpcomingFestivalsResponse{
		Festivals: festivalInfos,
	}, nil
}

// FestivalStatus queries overall festival status
func (k Keeper) FestivalStatus(goCtx context.Context, req *types.QueryFestivalStatusRequest) (*types.QueryFestivalStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get active and upcoming festivals
	activeFestivals := k.GetActiveFestivals(ctx)
	upcomingFestivals := k.GetUpcomingFestivals(ctx)
	
	// Extract festival IDs
	activeIds := make([]string, len(activeFestivals))
	for i, festival := range activeFestivals {
		activeIds[i] = festival.ID
	}
	
	upcomingIds := make([]string, len(upcomingFestivals))
	for i, festival := range upcomingFestivals {
		upcomingIds[i] = festival.ID
	}
	
	// Get current bonus rate
	bonusRate := k.GetFestivalBonusRate(ctx)
	
	// Find next festival
	var nextFestival *types.FestivalInfo
	if len(upcomingFestivals) > 0 {
		festival := upcomingFestivals[0] // Assuming sorted by date
		nextFestival = &types.FestivalInfo{
			Id:                  festival.ID,
			Name:                festival.Name,
			Description:         festival.Description,
			Date:                festival.Date.Format("2006-01-02T15:04:05Z07:00"),
			Type:                festival.Type,
			Region:              festival.Region,
			BonusRate:           festival.BonusRate,
			IsActive:            festival.IsActive,
			TraditionalGreeting: festival.TraditionalGreeting,
			CulturalTheme:       festival.CulturalTheme,
			DaysRemaining:       int32(festival.DaysRemaining),
			Significance:        festival.Significance,
		}
		
		if festival.EndDate != nil {
			endDate := festival.EndDate.Format("2006-01-02T15:04:05Z07:00")
			nextFestival.EndDate = &endDate
		}
	}

	return &types.QueryFestivalStatusResponse{
		ActiveFestivals:   activeIds,
		UpcomingFestivals: upcomingIds,
		CurrentBonusRate:  bonusRate.String(),
		NextFestival:      nextFestival,
		LastUpdated:       ctx.BlockTime().Format("2006-01-02T15:04:05Z07:00"),
		Metadata:          make(map[string]string),
	}, nil
}

// FestivalById queries a specific festival by ID
func (k Keeper) FestivalById(goCtx context.Context, req *types.QueryFestivalByIdRequest) (*types.QueryFestivalByIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// Get festival by ID from festival manager
	festival, found := k.GetFestivalByID(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("festival %s not found", req.Id))
	}
	
	// Convert to response format
	festivalInfo := &types.FestivalInfo{
		Id:                  festival.ID,
		Name:                festival.Name,
		Description:         festival.Description,
		Date:                festival.Date.Format("2006-01-02T15:04:05Z07:00"),
		Type:                festival.Type,
		Region:              festival.Region,
		BonusRate:           festival.BonusRate,
		IsActive:            festival.IsActive,
		TraditionalGreeting: festival.TraditionalGreeting,
		CulturalTheme:       festival.CulturalTheme,
		DaysRemaining:       int32(festival.DaysRemaining),
		Significance:        festival.Significance,
	}
	
	if festival.EndDate != nil {
		endDate := festival.EndDate.Format("2006-01-02T15:04:05Z07:00")
		festivalInfo.EndDate = &endDate
	}

	return &types.QueryFestivalByIdResponse{
		Festival: festivalInfo,
	}, nil
}

// QuoteStatistics returns quote usage statistics
func (k Keeper) QuoteStatistics(c context.Context, req *types.QueryQuoteStatisticsRequest) (*types.QueryQuoteStatisticsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	stats := k.GetQuoteStatistics(ctx)

	// Convert map[string]interface{} to proper response format
	totalQuotes := uint64(0)
	totalUsage := uint64(0)
	categoryCount := make(map[string]uint64)
	authorCount := make(map[string]uint64)
	languageCount := make(map[string]uint64)

	if v, ok := stats["total_quotes"].(int); ok {
		totalQuotes = uint64(v)
	}
	if v, ok := stats["total_usage"].(uint64); ok {
		totalUsage = v
	}
	if v, ok := stats["categories"].(map[string]int); ok {
		for cat, count := range v {
			categoryCount[cat] = uint64(count)
		}
	}
	if v, ok := stats["authors"].(map[string]int); ok {
		for author, count := range v {
			authorCount[author] = uint64(count)
		}
	}

	// Count languages from quotes
	quotes := k.GetAllQuotes(ctx)
	for _, quote := range quotes {
		languageCount[quote.Language]++
	}

	return &types.QueryQuoteStatisticsResponse{
		TotalQuotes:   totalQuotes,
		TotalUsage:    totalUsage,
		CategoryCount: categoryCount,
		AuthorCount:   authorCount,
		LanguageCount: languageCount,
	}, nil
}

// paginateStore is a helper function for paginating store iterators
func paginateStore(
	ctx sdk.Context,
	prefixStore prefix.Store,
	pageReq *query.PageRequest,
	onResult func(key []byte, value []byte) error,
) (*query.PageResponse, error) {
	if pageReq == nil {
		pageReq = &query.PageRequest{
			Offset:     0,
			Limit:      100,
			CountTotal: true,
		}
	}

	limit := pageReq.Limit
	if limit == 0 {
		limit = 100 // Default limit
	}

	var count uint64
	var pageRes *query.PageResponse

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	// Count total if requested
	if pageReq.CountTotal {
		for ; iterator.Valid(); iterator.Next() {
			count++
		}
		iterator = prefixStore.Iterator(nil, nil)
		defer iterator.Close()
	}

	// Skip to offset
	skip := pageReq.Offset
	for i := uint64(0); i < skip && iterator.Valid(); i++ {
		iterator.Next()
	}

	// Collect results up to limit
	collected := uint64(0)
	for ; iterator.Valid() && collected < limit; iterator.Next() {
		if err := onResult(iterator.Key(), iterator.Value()); err != nil {
			return nil, err
		}
		collected++
	}

	pageRes = &query.PageResponse{
		Total: count,
	}

	return pageRes, nil
}