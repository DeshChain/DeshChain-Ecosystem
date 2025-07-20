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
	cdc.RegisterConcrete(&MsgCreateScheme{}, "grampension/CreateScheme", nil)
	cdc.RegisterConcrete(&MsgUpdateScheme{}, "grampension/UpdateScheme", nil)
	cdc.RegisterConcrete(&MsgEnrollParticipant{}, "grampension/EnrollParticipant", nil)
	cdc.RegisterConcrete(&MsgMakeContribution{}, "grampension/MakeContribution", nil)
	cdc.RegisterConcrete(&MsgProcessMaturity{}, "grampension/ProcessMaturity", nil)
	cdc.RegisterConcrete(&MsgRequestWithdrawal{}, "grampension/RequestWithdrawal", nil)
	cdc.RegisterConcrete(&MsgProcessWithdrawal{}, "grampension/ProcessWithdrawal", nil)
	cdc.RegisterConcrete(&MsgUpdateKYCStatus{}, "grampension/UpdateKYCStatus", nil)
	cdc.RegisterConcrete(&MsgClaimReferral{}, "grampension/ClaimReferral", nil)
}

// RegisterInterfaces registers the module's interface types
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateScheme{},
		&MsgUpdateScheme{},
		&MsgEnrollParticipant{},
		&MsgMakeContribution{},
		&MsgProcessMaturity{},
		&MsgRequestWithdrawal{},
		&MsgProcessWithdrawal{},
		&MsgUpdateKYCStatus{},
		&MsgClaimReferral{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}

// Placeholder for protobuf service descriptor
var _Msg_serviceDesc = msgservice.ServiceDesc{
	ServiceName: ModuleName,
}