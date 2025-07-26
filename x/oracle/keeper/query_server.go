package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DeshChain/DeshChain-Ecosystem/x/oracle/types"
)

type queryServer struct {
	Keeper
}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}

func (k queryServer) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k queryServer) Price(goCtx context.Context, req *types.QueryPriceRequest) (*types.QueryPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Symbol == "" {
		return nil, status.Error(codes.InvalidArgument, "symbol cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	priceData, found := k.GetPriceData(ctx, req.Symbol)
	if !found {
		return nil, status.Error(codes.NotFound, "price not found")
	}

	return &types.QueryPriceResponse{PriceData: priceData}, nil
}

func (k queryServer) Prices(goCtx context.Context, req *types.QueryPricesRequest) (*types.QueryPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	priceStore := prefix.NewStore(store, types.PriceDataKeyPrefix)

	var prices []types.PriceData
	pageRes, err := query.Paginate(priceStore, req.Pagination, func(key []byte, value []byte) error {
		var priceData types.PriceData
		if err := k.cdc.Unmarshal(value, &priceData); err != nil {
			return err
		}
		prices = append(prices, priceData)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPricesResponse{
		Prices:     prices,
		Pagination: pageRes,
	}, nil
}

func (k queryServer) ExchangeRate(goCtx context.Context, req *types.QueryExchangeRateRequest) (*types.QueryExchangeRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Base == "" || req.Target == "" {
		return nil, status.Error(codes.InvalidArgument, "base and target currencies cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	exchangeRate, found := k.GetExchangeRate(ctx, req.Base, req.Target)
	if !found {
		return nil, status.Error(codes.NotFound, "exchange rate not found")
	}

	return &types.QueryExchangeRateResponse{ExchangeRate: exchangeRate}, nil
}

func (k queryServer) ExchangeRates(goCtx context.Context, req *types.QueryExchangeRatesRequest) (*types.QueryExchangeRatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	exchangeRateStore := prefix.NewStore(store, types.ExchangeRateKeyPrefix)

	var exchangeRates []types.ExchangeRate
	pageRes, err := query.Paginate(exchangeRateStore, req.Pagination, func(key []byte, value []byte) error {
		var exchangeRate types.ExchangeRate
		if err := k.cdc.Unmarshal(value, &exchangeRate); err != nil {
			return err
		}
		exchangeRates = append(exchangeRates, exchangeRate)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryExchangeRatesResponse{
		ExchangeRates: exchangeRates,
		Pagination:    pageRes,
	}, nil
}

func (k queryServer) OracleValidator(goCtx context.Context, req *types.QueryOracleValidatorRequest) (*types.QueryOracleValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Validator == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	oracleValidator, found := k.GetOracleValidator(ctx, req.Validator)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle validator not found")
	}

	return &types.QueryOracleValidatorResponse{OracleValidator: oracleValidator}, nil
}

func (k queryServer) OracleValidators(goCtx context.Context, req *types.QueryOracleValidatorsRequest) (*types.QueryOracleValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	validatorStore := prefix.NewStore(store, types.OracleValidatorKeyPrefix)

	var validators []types.OracleValidator
	pageRes, err := query.Paginate(validatorStore, req.Pagination, func(key []byte, value []byte) error {
		var validator types.OracleValidator
		if err := k.cdc.Unmarshal(value, &validator); err != nil {
			return err
		}
		validators = append(validators, validator)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOracleValidatorsResponse{
		OracleValidators: validators,
		Pagination:       pageRes,
	}, nil
}

func (k queryServer) PriceHistory(goCtx context.Context, req *types.QueryPriceHistoryRequest) (*types.QueryPriceHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Symbol == "" {
		return nil, status.Error(codes.InvalidArgument, "symbol cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	
	historicalPrices := k.GetPriceHistory(ctx, req.Symbol, req.Limit)

	return &types.QueryPriceHistoryResponse{
		Symbol:            req.Symbol,
		HistoricalPrices: historicalPrices,
	}, nil
}