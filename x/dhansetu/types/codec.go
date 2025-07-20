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
	cdc.RegisterConcrete(&MsgRegisterDhanPataAddress{}, "dhansetu/RegisterDhanPataAddress", nil)
	cdc.RegisterConcrete(&MsgCreateKshetraCoin{}, "dhansetu/CreateKshetraCoin", nil)
	cdc.RegisterConcrete(&MsgRegisterEnhancedMitra{}, "dhansetu/RegisterEnhancedMitra", nil)
	cdc.RegisterConcrete(&MsgProcessMoneyOrderWithDhanPata{}, "dhansetu/ProcessMoneyOrderWithDhanPata", nil)
	cdc.RegisterConcrete(&MsgUpdateDhanPataMetadata{}, "dhansetu/UpdateDhanPataMetadata", nil)
}

// RegisterInterfaces registers the module interfaces
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterDhanPataAddress{},
		&MsgCreateKshetraCoin{},
		&MsgRegisterEnhancedMitra{},
		&MsgProcessMoneyOrderWithDhanPata{},
		&MsgUpdateDhanPataMetadata{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}