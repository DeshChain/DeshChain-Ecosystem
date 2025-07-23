package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers the necessary tax module interfaces and concrete types
// on the provided LegacyAmino codec.
func RegisterCodec(cdc *codec.LegacyAmino) {
	// Register messages here when we add them
}

// RegisterInterfaces registers the tax module interface types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Register message implementations here when we add them
	
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
	Amino.Seal()
}