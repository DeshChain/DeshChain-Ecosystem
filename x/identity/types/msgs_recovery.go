package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Identity Backup and Recovery Messages

const (
	TypeMsgCreateIdentityBackup     = "create_identity_backup"
	TypeMsgInitiateRecovery         = "initiate_recovery"
	TypeMsgSubmitRecoveryProof      = "submit_recovery_proof"
	TypeMsgExecuteRecovery          = "execute_recovery"
	TypeMsgAddSocialRecoveryGuardian = "add_social_recovery_guardian"
	TypeMsgSubmitGuardianVote       = "submit_guardian_vote"
	TypeMsgVerifyBackupIntegrity    = "verify_backup_integrity"
	TypeMsgCreateDisasterRecoveryConfig = "create_disaster_recovery_config"
)

// MsgCreateIdentityBackup creates a complete backup of an identity
type MsgCreateIdentityBackup struct {
	Authority        string          `json:"authority"`
	HolderDID        string          `json:"holder_did"`
	RecoveryMethods  []RecoveryMethod `json:"recovery_methods"`
	EncryptionKey    string          `json:"encryption_key"` // Base64 encoded
	RetentionDays    int64           `json:"retention_days"`
}

// NewMsgCreateIdentityBackup creates a new MsgCreateIdentityBackup
func NewMsgCreateIdentityBackup(
	authority string,
	holderDID string,
	recoveryMethods []RecoveryMethod,
	encryptionKey string,
	retentionDays int64,
) *MsgCreateIdentityBackup {
	return &MsgCreateIdentityBackup{
		Authority:       authority,
		HolderDID:       holderDID,
		RecoveryMethods: recoveryMethods,
		EncryptionKey:   encryptionKey,
		RetentionDays:   retentionDays,
	}
}

// Route returns the message route
func (msg *MsgCreateIdentityBackup) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgCreateIdentityBackup) Type() string { return TypeMsgCreateIdentityBackup }

// GetSigners returns the signers
func (msg *MsgCreateIdentityBackup) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgCreateIdentityBackup) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgCreateIdentityBackup) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	if msg.HolderDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "holder DID cannot be empty")
	}
	if len(msg.RecoveryMethods) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one recovery method is required")
	}
	if msg.EncryptionKey == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "encryption key cannot be empty")
	}
	if msg.RetentionDays <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "retention days must be positive")
	}

	// Validate individual recovery methods
	for i, method := range msg.RecoveryMethods {
		if err := method.ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "invalid recovery method at index %d", i)
		}
	}

	return nil
}

// MsgInitiateRecovery initiates an identity recovery process
type MsgInitiateRecovery struct {
	Authority string `json:"authority"`
	HolderDID string `json:"holder_did"`
	BackupID  string `json:"backup_id"`
	Reason    string `json:"reason"`
}

// NewMsgInitiateRecovery creates a new MsgInitiateRecovery
func NewMsgInitiateRecovery(authority, holderDID, backupID, reason string) *MsgInitiateRecovery {
	return &MsgInitiateRecovery{
		Authority: authority,
		HolderDID: holderDID,
		BackupID:  backupID,
		Reason:    reason,
	}
}

// Route returns the message route
func (msg *MsgInitiateRecovery) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgInitiateRecovery) Type() string { return TypeMsgInitiateRecovery }

// GetSigners returns the signers
func (msg *MsgInitiateRecovery) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgInitiateRecovery) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgInitiateRecovery) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	if msg.HolderDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "holder DID cannot be empty")
	}
	if msg.BackupID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "backup ID cannot be empty")
	}
	if msg.Reason == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "reason cannot be empty")
	}

	return nil
}

// MsgSubmitRecoveryProof submits proof for a recovery method
type MsgSubmitRecoveryProof struct {
	Authority        string                 `json:"authority"`
	RequestID        string                 `json:"request_id"`
	MethodID         string                 `json:"method_id"`
	ProofData        string                 `json:"proof_data"` // Base64 encoded
	VerificationData map[string]interface{} `json:"verification_data,omitempty"`
}

// NewMsgSubmitRecoveryProof creates a new MsgSubmitRecoveryProof
func NewMsgSubmitRecoveryProof(
	authority string,
	requestID string,
	methodID string,
	proofData string,
	verificationData map[string]interface{},
) *MsgSubmitRecoveryProof {
	return &MsgSubmitRecoveryProof{
		Authority:        authority,
		RequestID:        requestID,
		MethodID:         methodID,
		ProofData:        proofData,
		VerificationData: verificationData,
	}
}

// Route returns the message route
func (msg *MsgSubmitRecoveryProof) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgSubmitRecoveryProof) Type() string { return TypeMsgSubmitRecoveryProof }

// GetSigners returns the signers
func (msg *MsgSubmitRecoveryProof) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgSubmitRecoveryProof) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgSubmitRecoveryProof) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	if msg.RequestID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "request ID cannot be empty")
	}
	if msg.MethodID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "method ID cannot be empty")
	}
	if msg.ProofData == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "proof data cannot be empty")
	}

	return nil
}

// MsgExecuteRecovery executes the recovery process
type MsgExecuteRecovery struct {
	Authority           string `json:"authority"`
	RequestID           string `json:"request_id"`
	NewControllerAddress string `json:"new_controller_address"`
	DecryptionKey       string `json:"decryption_key"` // Base64 encoded
}

// NewMsgExecuteRecovery creates a new MsgExecuteRecovery
func NewMsgExecuteRecovery(
	authority string,
	requestID string,
	newControllerAddress string,
	decryptionKey string,
) *MsgExecuteRecovery {
	return &MsgExecuteRecovery{
		Authority:            authority,
		RequestID:            requestID,
		NewControllerAddress: newControllerAddress,
		DecryptionKey:        decryptionKey,
	}
}

// Route returns the message route
func (msg *MsgExecuteRecovery) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgExecuteRecovery) Type() string { return TypeMsgExecuteRecovery }

// GetSigners returns the signers
func (msg *MsgExecuteRecovery) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgExecuteRecovery) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgExecuteRecovery) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.NewControllerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new controller address: %v", err)
	}

	if msg.RequestID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "request ID cannot be empty")
	}
	if msg.DecryptionKey == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "decryption key cannot be empty")
	}

	return nil
}

// MsgAddSocialRecoveryGuardian adds a guardian for social recovery
type MsgAddSocialRecoveryGuardian struct {
	Authority       string `json:"authority"`
	HolderDID       string `json:"holder_did"`
	GuardianDID     string `json:"guardian_did"`
	GuardianAddress string `json:"guardian_address"`
	GuardianName    string `json:"guardian_name"`
	Weight          int    `json:"weight"`
	ContactInfo     string `json:"contact_info"`
	PublicKey       string `json:"public_key"`
}

// NewMsgAddSocialRecoveryGuardian creates a new MsgAddSocialRecoveryGuardian
func NewMsgAddSocialRecoveryGuardian(
	authority string,
	holderDID string,
	guardianDID string,
	guardianAddress string,
	guardianName string,
	weight int,
	contactInfo string,
	publicKey string,
) *MsgAddSocialRecoveryGuardian {
	return &MsgAddSocialRecoveryGuardian{
		Authority:       authority,
		HolderDID:       holderDID,
		GuardianDID:     guardianDID,
		GuardianAddress: guardianAddress,
		GuardianName:    guardianName,
		Weight:          weight,
		ContactInfo:     contactInfo,
		PublicKey:       publicKey,
	}
}

// Route returns the message route
func (msg *MsgAddSocialRecoveryGuardian) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgAddSocialRecoveryGuardian) Type() string { return TypeMsgAddSocialRecoveryGuardian }

// GetSigners returns the signers
func (msg *MsgAddSocialRecoveryGuardian) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgAddSocialRecoveryGuardian) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgAddSocialRecoveryGuardian) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.GuardianAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid guardian address: %v", err)
	}

	if msg.HolderDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "holder DID cannot be empty")
	}
	if msg.GuardianDID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "guardian DID cannot be empty")
	}
	if msg.GuardianName == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "guardian name cannot be empty")
	}
	if msg.Weight <= 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "guardian weight must be positive")
	}
	if msg.PublicKey == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "guardian public key cannot be empty")
	}

	return nil
}

// MsgSubmitGuardianVote submits a guardian vote for recovery
type MsgSubmitGuardianVote struct {
	Authority string   `json:"authority"`
	RequestID string   `json:"request_id"`
	Vote      VoteType `json:"vote"`
	Reason    string   `json:"reason,omitempty"`
	Signature string   `json:"signature"`
}

// NewMsgSubmitGuardianVote creates a new MsgSubmitGuardianVote
func NewMsgSubmitGuardianVote(
	authority string,
	requestID string,
	vote VoteType,
	reason string,
	signature string,
) *MsgSubmitGuardianVote {
	return &MsgSubmitGuardianVote{
		Authority: authority,
		RequestID: requestID,
		Vote:      vote,
		Reason:    reason,
		Signature: signature,
	}
}

// Route returns the message route
func (msg *MsgSubmitGuardianVote) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgSubmitGuardianVote) Type() string { return TypeMsgSubmitGuardianVote }

// GetSigners returns the signers
func (msg *MsgSubmitGuardianVote) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgSubmitGuardianVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgSubmitGuardianVote) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	if msg.RequestID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "request ID cannot be empty")
	}
	if msg.Signature == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "signature cannot be empty")
	}

	return nil
}

// MsgVerifyBackupIntegrity verifies the integrity of a backup
type MsgVerifyBackupIntegrity struct {
	Authority       string `json:"authority"`
	BackupID        string `json:"backup_id"`
	VerificationKey string `json:"verification_key,omitempty"` // Base64 encoded
}

// NewMsgVerifyBackupIntegrity creates a new MsgVerifyBackupIntegrity
func NewMsgVerifyBackupIntegrity(
	authority string,
	backupID string,
	verificationKey string,
) *MsgVerifyBackupIntegrity {
	return &MsgVerifyBackupIntegrity{
		Authority:       authority,
		BackupID:        backupID,
		VerificationKey: verificationKey,
	}
}

// Route returns the message route
func (msg *MsgVerifyBackupIntegrity) Route() string { return RouterKey }

// Type returns the message type
func (msg *MsgVerifyBackupIntegrity) Type() string { return TypeMsgVerifyBackupIntegrity }

// GetSigners returns the signers
func (msg *MsgVerifyBackupIntegrity) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes returns the bytes to sign
func (msg *MsgVerifyBackupIntegrity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic performs basic validation
func (msg *MsgVerifyBackupIntegrity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address: %v", err)
	}

	if msg.BackupID == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "backup ID cannot be empty")
	}

	return nil
}

// Response types

// MsgCreateIdentityBackupResponse is the response for MsgCreateIdentityBackup
type MsgCreateIdentityBackupResponse struct {
	BackupID      string    `json:"backup_id"`
	BackupVersion int64     `json:"backup_version"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// MsgInitiateRecoveryResponse is the response for MsgInitiateRecovery
type MsgInitiateRecoveryResponse struct {
	RequestID     string    `json:"request_id"`
	RequiredScore int       `json:"required_score"`
	ExpiresAt     time.Time `json:"expires_at"`
	MaxAttempts   int       `json:"max_attempts"`
}

// MsgSubmitRecoveryProofResponse is the response for MsgSubmitRecoveryProof
type MsgSubmitRecoveryProofResponse struct {
	AttemptID       string        `json:"attempt_id"`
	Status          AttemptStatus `json:"status"`
	Confidence      int           `json:"confidence"`
	CurrentScore    int           `json:"current_score"`
	RequiredScore   int           `json:"required_score"`
	CanAttemptMore  bool          `json:"can_attempt_more"`
}

// MsgExecuteRecoveryResponse is the response for MsgExecuteRecovery
type MsgExecuteRecoveryResponse struct {
	RecoveredAt         time.Time `json:"recovered_at"`
	NewControllerAddress string    `json:"new_controller_address"`
	DataRestored        bool      `json:"data_restored"`
}

// MsgAddSocialRecoveryGuardianResponse is the response for MsgAddSocialRecoveryGuardian
type MsgAddSocialRecoveryGuardianResponse struct {
	GuardianID string    `json:"guardian_id"`
	AddedAt    time.Time `json:"added_at"`
}

// MsgSubmitGuardianVoteResponse is the response for MsgSubmitGuardianVote
type MsgSubmitGuardianVoteResponse struct {
	VoteID        string    `json:"vote_id"`
	VotedAt       time.Time `json:"voted_at"`
	Weight        int       `json:"weight"`
	CurrentVotes  int       `json:"current_votes"`
	RequiredVotes int       `json:"required_votes"`
}

// MsgVerifyBackupIntegrityResponse is the response for MsgVerifyBackupIntegrity
type MsgVerifyBackupIntegrityResponse struct {
	VerificationID      string   `json:"verification_id"`
	IntegrityValid      bool     `json:"integrity_valid"`
	DecryptionValid     bool     `json:"decryption_valid"`
	DataCompleteness    float64  `json:"data_completeness"`
	RecoverabilityScore int      `json:"recoverability_score"`
	IssuesFound         []string `json:"issues_found,omitempty"`
	Recommendations     []string `json:"recommendations,omitempty"`
}
