package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/dswf interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgProposeAllocation{}, "dswf/ProposeAllocation", nil)
	cdc.RegisterConcrete(&MsgApproveAllocation{}, "dswf/ApproveAllocation", nil)
	cdc.RegisterConcrete(&MsgExecuteDisbursement{}, "dswf/ExecuteDisbursement", nil)
	cdc.RegisterConcrete(&MsgUpdatePortfolio{}, "dswf/UpdatePortfolio", nil)
	cdc.RegisterConcrete(&MsgSubmitMonthlyReport{}, "dswf/SubmitMonthlyReport", nil)
	cdc.RegisterConcrete(&MsgUpdateGovernance{}, "dswf/UpdateGovernance", nil)
	cdc.RegisterConcrete(&MsgRecordReturns{}, "dswf/RecordReturns", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "dswf/UpdateParams", nil)
}

// RegisterInterfaces registers the x/dswf interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgProposeAllocation{},
		&MsgApproveAllocation{},
		&MsgExecuteDisbursement{},
		&MsgUpdatePortfolio{},
		&MsgSubmitMonthlyReport{},
		&MsgUpdateGovernance{},
		&MsgRecordReturns{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)