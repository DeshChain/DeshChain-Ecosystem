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

package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers the necessary types and interfaces for the module
func RegisterCodec(cdc *codec.LegacyAmino) {
	// Register messages
	cdc.RegisterConcrete(&MsgCreateMoneyOrder{}, "moneyorder/CreateMoneyOrder", nil)
	cdc.RegisterConcrete(&MsgCreateFixedRatePool{}, "moneyorder/CreateFixedRatePool", nil)
	cdc.RegisterConcrete(&MsgCreateVillagePool{}, "moneyorder/CreateVillagePool", nil)
	cdc.RegisterConcrete(&MsgAddLiquidity{}, "moneyorder/AddLiquidity", nil)
	cdc.RegisterConcrete(&MsgRemoveLiquidity{}, "moneyorder/RemoveLiquidity", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountIn{}, "moneyorder/SwapExactAmountIn", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountOut{}, "moneyorder/SwapExactAmountOut", nil)
	cdc.RegisterConcrete(&MsgJoinVillagePool{}, "moneyorder/JoinVillagePool", nil)
	cdc.RegisterConcrete(&MsgClaimRewards{}, "moneyorder/ClaimRewards", nil)
	cdc.RegisterConcrete(&MsgUpdatePoolParams{}, "moneyorder/UpdatePoolParams", nil)
}

// RegisterInterfaces registers the interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Register messages
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateMoneyOrder{},
		&MsgCreateFixedRatePool{},
		&MsgCreateVillagePool{},
		&MsgAddLiquidity{},
		&MsgRemoveLiquidity{},
		&MsgSwapExactAmountIn{},
		&MsgSwapExactAmountOut{},
		&MsgJoinVillagePool{},
		&MsgClaimRewards{},
		&MsgUpdatePoolParams{},
	)
	
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// ModuleCdc references the global moneyorder module codec
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	
	// AminoCodec is the legacy amino codec
	AminoCodec = codec.NewLegacyAmino()
)

func init() {
	RegisterCodec(AminoCodec)
	RegisterInterfaces(cdctypes.NewInterfaceRegistry())
	
	// Seal the amino codec to prevent further modifications
	AminoCodec.Seal()
}

// Note: _Msg_serviceDesc will be generated from proto files
// This is a placeholder for the generated code
var _Msg_serviceDesc = msgservice.ServiceDesc{
	ServiceName: "deshchain.moneyorder.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods:     []msgservice.MethodDesc{},
}

// MsgServer is the server API for Msg service (interface placeholder)
type MsgServer interface {
	CreateMoneyOrder(sdk.Context, *MsgCreateMoneyOrder) (*MsgCreateMoneyOrderResponse, error)
	CreateFixedRatePool(sdk.Context, *MsgCreateFixedRatePool) (*MsgCreateFixedRatePoolResponse, error)
	CreateVillagePool(sdk.Context, *MsgCreateVillagePool) (*MsgCreateVillagePoolResponse, error)
	AddLiquidity(sdk.Context, *MsgAddLiquidity) (*MsgAddLiquidityResponse, error)
	RemoveLiquidity(sdk.Context, *MsgRemoveLiquidity) (*MsgRemoveLiquidityResponse, error)
	SwapExactAmountIn(sdk.Context, *MsgSwapExactAmountIn) (*MsgSwapExactAmountInResponse, error)
	SwapExactAmountOut(sdk.Context, *MsgSwapExactAmountOut) (*MsgSwapExactAmountOutResponse, error)
	JoinVillagePool(sdk.Context, *MsgJoinVillagePool) (*MsgJoinVillagePoolResponse, error)
	ClaimRewards(sdk.Context, *MsgClaimRewards) (*MsgClaimRewardsResponse, error)
	UpdatePoolParams(sdk.Context, *MsgUpdatePoolParams) (*MsgUpdatePoolParamsResponse, error)
}

// Response types (these will be properly defined in proto files)
type (
	MsgCreateMoneyOrderResponse struct {
		OrderId         string `json:"order_id" yaml:"order_id"`
		ReferenceNumber string `json:"reference_number" yaml:"reference_number"`
	}
	
	MsgCreateFixedRatePoolResponse struct {
		PoolId uint64 `json:"pool_id" yaml:"pool_id"`
	}
	
	MsgCreateVillagePoolResponse struct {
		PoolId uint64 `json:"pool_id" yaml:"pool_id"`
	}
	
	MsgAddLiquidityResponse struct {
		SharesOut sdk.Int `json:"shares_out" yaml:"shares_out"`
	}
	
	MsgRemoveLiquidityResponse struct {
		TokensOut sdk.Coins `json:"tokens_out" yaml:"tokens_out"`
	}
	
	MsgSwapExactAmountInResponse struct {
		TokenOut sdk.Coin `json:"token_out" yaml:"token_out"`
	}
	
	MsgSwapExactAmountOutResponse struct {
		TokenIn sdk.Coin `json:"token_in" yaml:"token_in"`
	}
	
	MsgJoinVillagePoolResponse struct {
		MemberId string `json:"member_id" yaml:"member_id"`
	}
	
	MsgClaimRewardsResponse struct {
		RewardsClaimed sdk.Coins `json:"rewards_claimed" yaml:"rewards_claimed"`
	}
	
	MsgUpdatePoolParamsResponse struct {
		Success bool `json:"success" yaml:"success"`
	}
)