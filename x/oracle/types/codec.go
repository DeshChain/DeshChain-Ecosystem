package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSubmitPrice{}, "oracle/MsgSubmitPrice", nil)
	cdc.RegisterConcrete(&MsgSubmitExchangeRate{}, "oracle/MsgSubmitExchangeRate", nil)
	cdc.RegisterConcrete(&MsgRegisterOracleValidator{}, "oracle/MsgRegisterOracleValidator", nil)
	cdc.RegisterConcrete(&MsgUpdateOracleValidator{}, "oracle/MsgUpdateOracleValidator", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "oracle/MsgUpdateParams", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitPrice{},
		&MsgSubmitExchangeRate{},
		&MsgRegisterOracleValidator{},
		&MsgUpdateOracleValidator{},
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