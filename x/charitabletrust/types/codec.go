package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/charitabletrust interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateAllocationProposal{}, "charitabletrust/CreateAllocationProposal", nil)
	cdc.RegisterConcrete(&MsgVoteOnProposal{}, "charitabletrust/VoteOnProposal", nil)
	cdc.RegisterConcrete(&MsgExecuteAllocation{}, "charitabletrust/ExecuteAllocation", nil)
	cdc.RegisterConcrete(&MsgSubmitImpactReport{}, "charitabletrust/SubmitImpactReport", nil)
	cdc.RegisterConcrete(&MsgVerifyImpactReport{}, "charitabletrust/VerifyImpactReport", nil)
	cdc.RegisterConcrete(&MsgReportFraud{}, "charitabletrust/ReportFraud", nil)
	cdc.RegisterConcrete(&MsgInvestigateFraud{}, "charitabletrust/InvestigateFraud", nil)
	cdc.RegisterConcrete(&MsgUpdateTrustees{}, "charitabletrust/UpdateTrustees", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "charitabletrust/UpdateParams", nil)
}

// RegisterInterfaces registers the x/charitabletrust interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAllocationProposal{},
		&MsgVoteOnProposal{},
		&MsgExecuteAllocation{},
		&MsgSubmitImpactReport{},
		&MsgVerifyImpactReport{},
		&MsgReportFraud{},
		&MsgInvestigateFraud{},
		&MsgUpdateTrustees{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)