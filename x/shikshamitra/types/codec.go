package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(amino)
	amino.Seal()
}

// RegisterCodec registers the necessary x/shikshamitra interfaces and concrete types
// on the provided Amino codec. These types are used for Amino JSON serialization.
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgApplyEducationLoan{}, "shikshamitra/ApplyEducationLoan", nil)
	cdc.RegisterConcrete(&MsgUpdateAcademicProgress{}, "shikshamitra/UpdateAcademicProgress", nil)
	cdc.RegisterConcrete(&MsgUpdateEmploymentStatus{}, "shikshamitra/UpdateEmploymentStatus", nil)
	cdc.RegisterConcrete(&MsgApplyScholarship{}, "shikshamitra/ApplyScholarship", nil)
	cdc.RegisterConcrete(&MsgApproveLoan{}, "shikshamitra/ApproveLoan", nil)
	cdc.RegisterConcrete(&MsgRejectLoan{}, "shikshamitra/RejectLoan", nil)
	cdc.RegisterConcrete(&MsgDisburseLoan{}, "shikshamitra/DisburseLoan", nil)
	cdc.RegisterConcrete(&MsgRepayLoan{}, "shikshamitra/RepayLoan", nil)
	cdc.RegisterConcrete(&MsgStartMoratorium{}, "shikshamitra/StartMoratorium", nil)
	cdc.RegisterConcrete(&MsgMigrateStudentsToIdentity{}, "shikshamitra/MigrateStudentsToIdentity", nil)
	cdc.RegisterConcrete(&MsgCreateStudentCredential{}, "shikshamitra/CreateStudentCredential", nil)
}

// RegisterInterfaces registers the x/shikshamitra interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgApplyEducationLoan{},
		&MsgUpdateAcademicProgress{},
		&MsgUpdateEmploymentStatus{},
		&MsgApplyScholarship{},
		&MsgApproveLoan{},
		&MsgRejectLoan{},
		&MsgDisburseLoan{},
		&MsgRepayLoan{},
		&MsgStartMoratorium{},
		&MsgMigrateStudentsToIdentity{},
		&MsgCreateStudentCredential{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}