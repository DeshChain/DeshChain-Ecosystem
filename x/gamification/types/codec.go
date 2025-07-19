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

// RegisterLegacyAminoCodec registers the necessary x/gamification interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateProfile{}, "gamification/CreateProfile", nil)
	cdc.RegisterConcrete(&MsgUpdateProfile{}, "gamification/UpdateProfile", nil)
	cdc.RegisterConcrete(&MsgSelectAvatar{}, "gamification/SelectAvatar", nil)
	cdc.RegisterConcrete(&MsgClaimAchievement{}, "gamification/ClaimAchievement", nil)
	cdc.RegisterConcrete(&MsgRecordAction{}, "gamification/RecordAction", nil)
	cdc.RegisterConcrete(&MsgShareAchievement{}, "gamification/ShareAchievement", nil)
	cdc.RegisterConcrete(&MsgJoinTeamBattle{}, "gamification/JoinTeamBattle", nil)
	cdc.RegisterConcrete(&MsgCreateTeamBattle{}, "gamification/CreateTeamBattle", nil)
	cdc.RegisterConcrete(&MsgCompleteDailyChallenge{}, "gamification/CompleteDailyChallenge", nil)
}

// RegisterInterfaces registers the x/gamification interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateProfile{},
		&MsgUpdateProfile{},
		&MsgSelectAvatar{},
		&MsgClaimAchievement{},
		&MsgRecordAction{},
		&MsgShareAchievement{},
		&MsgJoinTeamBattle{},
		&MsgCreateTeamBattle{},
		&MsgCompleteDailyChallenge{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
	Amino.Seal()
}