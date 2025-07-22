package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterParty{}, "tradefinance/MsgRegisterParty", nil)
	cdc.RegisterConcrete(&MsgIssueLc{}, "tradefinance/MsgIssueLc", nil)
	cdc.RegisterConcrete(&MsgAcceptLc{}, "tradefinance/MsgAcceptLc", nil)
	cdc.RegisterConcrete(&MsgSubmitDocuments{}, "tradefinance/MsgSubmitDocuments", nil)
	cdc.RegisterConcrete(&MsgVerifyDocument{}, "tradefinance/MsgVerifyDocument", nil)
	cdc.RegisterConcrete(&MsgRequestPayment{}, "tradefinance/MsgRequestPayment", nil)
	cdc.RegisterConcrete(&MsgMakePayment{}, "tradefinance/MsgMakePayment", nil)
	cdc.RegisterConcrete(&MsgAmendLc{}, "tradefinance/MsgAmendLc", nil)
	cdc.RegisterConcrete(&MsgCancelLc{}, "tradefinance/MsgCancelLc", nil)
	cdc.RegisterConcrete(&MsgCreateInsurancePolicy{}, "tradefinance/MsgCreateInsurancePolicy", nil)
	cdc.RegisterConcrete(&MsgUpdateShipment{}, "tradefinance/MsgUpdateShipment", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "tradefinance/MsgUpdateParams", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterParty{},
		&MsgIssueLc{},
		&MsgAcceptLc{},
		&MsgSubmitDocuments{},
		&MsgVerifyDocument{},
		&MsgRequestPayment{},
		&MsgMakePayment{},
		&MsgAmendLc{},
		&MsgCancelLc{},
		&MsgCreateInsurancePolicy{},
		&MsgUpdateShipment{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	RegisterInterfaces(ModuleCdc.InterfaceRegistry())
	Amino.Seal()
}