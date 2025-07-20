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

// RegisterLegacyAminoCodec registers the necessary cultural interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Register messages
	// TODO: Register messages when implemented
	// cdc.RegisterConcrete(&MsgAddQuote{}, "cultural/AddQuote", nil)
	// cdc.RegisterConcrete(&MsgUpdateQuote{}, "cultural/UpdateQuote", nil)
	// cdc.RegisterConcrete(&MsgDeleteQuote{}, "cultural/DeleteQuote", nil)
}

// RegisterInterfaces registers the cultural interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Register messages
	// TODO: Register message implementations when created
	// registry.RegisterImplementations((*sdk.Msg)(nil),
	// 	&MsgAddQuote{},
	// 	&MsgUpdateQuote{},
	// 	&MsgDeleteQuote{},
	// )

	// TODO: Register when Msg service is generated from proto
	// msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}