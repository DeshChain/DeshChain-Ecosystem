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
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary donation module interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgDonate{}, "donation/Donate")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "donation/UpdateParams")

	cdc.RegisterConcrete(&Params{}, "donation/Params", nil)
	cdc.RegisterConcrete(&NGOWallet{}, "donation/NGOWallet", nil)
	cdc.RegisterConcrete(&DonationRecord{}, "donation/DonationRecord", nil)
	cdc.RegisterConcrete(&DistributionRecord{}, "donation/DistributionRecord", nil)
	cdc.RegisterConcrete(&AuditReport{}, "donation/AuditReport", nil)
	cdc.RegisterConcrete(&BeneficiaryTestimonial{}, "donation/BeneficiaryTestimonial", nil)
	cdc.RegisterConcrete(&Campaign{}, "donation/Campaign", nil)
	cdc.RegisterConcrete(&RecurringDonation{}, "donation/RecurringDonation", nil)
	cdc.RegisterConcrete(&EmergencyPause{}, "donation/EmergencyPause", nil)
	cdc.RegisterConcrete(&Statistics{}, "donation/Statistics", nil)
	cdc.RegisterConcrete(&FundFlow{}, "donation/FundFlow", nil)
	cdc.RegisterConcrete(&TransparencyScore{}, "donation/TransparencyScore", nil)
	cdc.RegisterConcrete(&VerificationQueueItem{}, "donation/VerificationQueueItem", nil)
}

// RegisterInterfaces registers the donation module interfaces to the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDonate{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)