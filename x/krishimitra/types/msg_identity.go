package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgMigrateFarmersToIdentity    = "migrate_farmers_to_identity"
	TypeMsgCreateFarmerCredential      = "create_farmer_credential"
	TypeMsgCreateLandRecordCredential  = "create_land_record_credential"
)

var (
	_ sdk.Msg = &MsgMigrateFarmersToIdentity{}
	_ sdk.Msg = &MsgCreateFarmerCredential{}
	_ sdk.Msg = &MsgCreateLandRecordCredential{}
)

// MsgMigrateFarmersToIdentity migrates farmers to identity system
type MsgMigrateFarmersToIdentity struct {
	Authority string `json:"authority"`
}

func NewMsgMigrateFarmersToIdentity(authority string) *MsgMigrateFarmersToIdentity {
	return &MsgMigrateFarmersToIdentity{
		Authority: authority,
	}
}

func (msg *MsgMigrateFarmersToIdentity) Route() string {
	return RouterKey
}

func (msg *MsgMigrateFarmersToIdentity) Type() string {
	return TypeMsgMigrateFarmersToIdentity
}

func (msg *MsgMigrateFarmersToIdentity) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgMigrateFarmersToIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMigrateFarmersToIdentity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

// MsgCreateFarmerCredential creates a credential for a farmer
type MsgCreateFarmerCredential struct {
	Authority     string `json:"authority"`
	FarmerAddress string `json:"farmer_address"`
	Pincode       string `json:"pincode"`
}

func NewMsgCreateFarmerCredential(authority, farmerAddress, pincode string) *MsgCreateFarmerCredential {
	return &MsgCreateFarmerCredential{
		Authority:     authority,
		FarmerAddress: farmerAddress,
		Pincode:       pincode,
	}
}

func (msg *MsgCreateFarmerCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateFarmerCredential) Type() string {
	return TypeMsgCreateFarmerCredential
}

func (msg *MsgCreateFarmerCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateFarmerCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateFarmerCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.FarmerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid farmer address (%s)", err)
	}
	
	if msg.Pincode == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pincode cannot be empty")
	}
	
	return nil
}

// MsgCreateLandRecordCredential creates a land record credential
type MsgCreateLandRecordCredential struct {
	Authority     string `json:"authority"`
	FarmerAddress string `json:"farmer_address"`
	KhataNumber   string `json:"khata_number"`
	TotalArea     string `json:"total_area"`
}

func NewMsgCreateLandRecordCredential(authority, farmerAddress, khataNumber, totalArea string) *MsgCreateLandRecordCredential {
	return &MsgCreateLandRecordCredential{
		Authority:     authority,
		FarmerAddress: farmerAddress,
		KhataNumber:   khataNumber,
		TotalArea:     totalArea,
	}
}

func (msg *MsgCreateLandRecordCredential) Route() string {
	return RouterKey
}

func (msg *MsgCreateLandRecordCredential) Type() string {
	return TypeMsgCreateLandRecordCredential
}

func (msg *MsgCreateLandRecordCredential) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateLandRecordCredential) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateLandRecordCredential) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	
	_, err = sdk.AccAddressFromBech32(msg.FarmerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid farmer address (%s)", err)
	}
	
	if msg.KhataNumber == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "khata number cannot be empty")
	}
	
	if msg.TotalArea == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "total area cannot be empty")
	}
	
	// Validate area is a valid decimal
	_, err = sdk.NewDecFromStr(msg.TotalArea)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid total area: %s", err)
	}
	
	return nil
}

// Query types for identity integration

// QueryFarmerIdentityRequest is the request for farmer identity
type QueryFarmerIdentityRequest struct {
	FarmerAddress string `json:"farmer_address"`
}

// QueryFarmerIdentityResponse is the response for farmer identity
type QueryFarmerIdentityResponse struct {
	Address            string   `json:"address"`
	HasIdentity        bool     `json:"has_identity"`
	DID                string   `json:"did,omitempty"`
	IsKYCVerified      bool     `json:"is_kyc_verified"`
	KYCLevel           string   `json:"kyc_level"`
	HasLandRecords     bool     `json:"has_land_records"`
	TotalLandArea      sdk.Dec  `json:"total_land_area,omitempty"`
	RegisteredCrops    []string `json:"registered_crops"`
	HasCropInsurance   bool     `json:"has_crop_insurance"`
	PINCode            string   `json:"pin_code"`
	VillageCode        string   `json:"village_code,omitempty"`
	ActiveLoans        []string `json:"active_loans"`
}

// QueryFarmerCredentialsRequest is the request for farmer credentials
type QueryFarmerCredentialsRequest struct {
	FarmerAddress string `json:"farmer_address"`
}

// QueryFarmerCredentialsResponse is the response for farmer credentials
type QueryFarmerCredentialsResponse struct {
	Address         string                  `json:"address"`
	DID             string                  `json:"did"`
	Credentials     []FarmerCredentialInfo  `json:"credentials"`
}

// FarmerCredentialInfo contains farmer credential information
type FarmerCredentialInfo struct {
	CredentialID string   `json:"credential_id"`
	Type         []string `json:"type"`
	LoanID       string   `json:"loan_id,omitempty"`
	Status       string   `json:"status"`
	CropType     string   `json:"crop_type,omitempty"`
	LandArea     string   `json:"land_area,omitempty"`
	IsValid      bool     `json:"is_valid"`
}

// QueryLandOwnershipRequest is the request for land ownership verification
type QueryLandOwnershipRequest struct {
	FarmerAddress string `json:"farmer_address"`
	RequiredArea  string `json:"required_area"`
}

// QueryLandOwnershipResponse is the response for land ownership verification
type QueryLandOwnershipResponse struct {
	Address           string  `json:"address"`
	TotalLandArea     sdk.Dec `json:"total_land_area"`
	RequiredArea      sdk.Dec `json:"required_area"`
	HasSufficientLand bool    `json:"has_sufficient_land"`
	LandCredentials   []string `json:"land_credentials"`
}