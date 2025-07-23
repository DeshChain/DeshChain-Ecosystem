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
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Campaign represents a fundraising campaign
type Campaign struct {
	Id                uint64         `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	NgoWalletId       uint64         `protobuf:"varint,2,opt,name=ngo_wallet_id,json=ngoWalletId,proto3" json:"ngo_wallet_id,omitempty"`
	Name              string         `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description       string         `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	TargetAmount      sdk.Coins      `protobuf:"bytes,5,rep,name=target_amount,json=targetAmount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"target_amount"`
	RaisedAmount      sdk.Coins      `protobuf:"bytes,6,rep,name=raised_amount,json=raisedAmount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"raised_amount"`
	StartDate         int64          `protobuf:"varint,7,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate           int64          `protobuf:"varint,8,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	Category          string         `protobuf:"bytes,9,opt,name=category,proto3" json:"category,omitempty"`
	Status            string         `protobuf:"bytes,10,opt,name=status,proto3" json:"status,omitempty"`
	IsActive          bool           `protobuf:"varint,11,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	CreatedAt         int64          `protobuf:"varint,12,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt         int64          `protobuf:"varint,13,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	CompletedAt       int64          `protobuf:"varint,14,opt,name=completed_at,json=completedAt,proto3" json:"completed_at,omitempty"`
	DonorCount        uint64         `protobuf:"varint,15,opt,name=donor_count,json=donorCount,proto3" json:"donor_count,omitempty"`
	BeneficiaryTarget uint64         `protobuf:"varint,16,opt,name=beneficiary_target,json=beneficiaryTarget,proto3" json:"beneficiary_target,omitempty"`
	RegionsTargeted   []string       `protobuf:"bytes,17,rep,name=regions_targeted,json=regionsTargeted,proto3" json:"regions_targeted,omitempty"`
	DocumentsHash     []string       `protobuf:"bytes,18,rep,name=documents_hash,json=documentsHash,proto3" json:"documents_hash,omitempty"`
	ImagesHash        []string       `protobuf:"bytes,19,rep,name=images_hash,json=imagesHash,proto3" json:"images_hash,omitempty"`
	VideosHash        []string       `protobuf:"bytes,20,rep,name=videos_hash,json=videosHash,proto3" json:"videos_hash,omitempty"`
	Milestones        []Milestone    `protobuf:"bytes,21,rep,name=milestones,proto3" json:"milestones"`
	UpdateMessages    []UpdateMessage `protobuf:"bytes,22,rep,name=update_messages,json=updateMessages,proto3" json:"update_messages"`
	MatchingEnabled   bool           `protobuf:"varint,23,opt,name=matching_enabled,json=matchingEnabled,proto3" json:"matching_enabled,omitempty"`
	MatchingRatio     string         `protobuf:"bytes,24,opt,name=matching_ratio,json=matchingRatio,proto3" json:"matching_ratio,omitempty"`
	MatchingSponsor   string         `protobuf:"bytes,25,opt,name=matching_sponsor,json=matchingSponsor,proto3" json:"matching_sponsor,omitempty"`
	MatchingBudget    sdk.Coins      `protobuf:"bytes,26,rep,name=matching_budget,json=matchingBudget,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"matching_budget"`
	MatchingUsed      sdk.Coins      `protobuf:"bytes,27,rep,name=matching_used,json=matchingUsed,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"matching_used"`
}

// Milestone represents a campaign milestone
type Milestone struct {
	Name              string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description       string    `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	TargetAmount      sdk.Coins `protobuf:"bytes,3,rep,name=target_amount,json=targetAmount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"target_amount"`
	TargetDate        int64     `protobuf:"varint,4,opt,name=target_date,json=targetDate,proto3" json:"target_date,omitempty"`
	IsCompleted       bool      `protobuf:"varint,5,opt,name=is_completed,json=isCompleted,proto3" json:"is_completed,omitempty"`
	CompletedAt       int64     `protobuf:"varint,6,opt,name=completed_at,json=completedAt,proto3" json:"completed_at,omitempty"`
	ProofDocumentHash string    `protobuf:"bytes,7,opt,name=proof_document_hash,json=proofDocumentHash,proto3" json:"proof_document_hash,omitempty"`
}

// UpdateMessage represents a campaign update message
type UpdateMessage struct {
	Title        string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Content      string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	PostedAt     int64  `protobuf:"varint,3,opt,name=posted_at,json=postedAt,proto3" json:"posted_at,omitempty"`
	PostedBy     string `protobuf:"bytes,4,opt,name=posted_by,json=postedBy,proto3" json:"posted_by,omitempty"`
	AttachmentsHash []string `protobuf:"bytes,5,rep,name=attachments_hash,json=attachmentsHash,proto3" json:"attachments_hash,omitempty"`
}

// ValidateCampaign validates a campaign
func ValidateCampaign(campaign Campaign) error {
	if campaign.NgoWalletId == 0 {
		return fmt.Errorf("NGO wallet ID cannot be zero")
	}

	if len(campaign.Name) == 0 {
		return fmt.Errorf("campaign name cannot be empty")
	}

	if len(campaign.Description) == 0 {
		return fmt.Errorf("campaign description cannot be empty")
	}

	if campaign.TargetAmount.IsZero() {
		return fmt.Errorf("target amount cannot be zero")
	}

	if campaign.StartDate <= 0 || campaign.EndDate <= 0 {
		return fmt.Errorf("campaign dates must be positive")
	}

	if campaign.EndDate <= campaign.StartDate {
		return fmt.Errorf("campaign end date must be after start date")
	}

	if campaign.CreatedAt <= 0 {
		return fmt.Errorf("creation timestamp must be positive")
	}

	return nil
}

// RecurringDonation represents a recurring donation setup
type RecurringDonation struct {
	Id               uint64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Donor            string    `protobuf:"bytes,2,opt,name=donor,proto3" json:"donor,omitempty"`
	NgoWalletId      uint64    `protobuf:"varint,3,opt,name=ngo_wallet_id,json=ngoWalletId,proto3" json:"ngo_wallet_id,omitempty"`
	Amount           sdk.Coins `protobuf:"bytes,4,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	Frequency        string    `protobuf:"bytes,5,opt,name=frequency,proto3" json:"frequency,omitempty"`
	DayOfMonth       int32     `protobuf:"varint,6,opt,name=day_of_month,json=dayOfMonth,proto3" json:"day_of_month,omitempty"`
	DayOfWeek        int32     `protobuf:"varint,7,opt,name=day_of_week,json=dayOfWeek,proto3" json:"day_of_week,omitempty"`
	StartDate        int64     `protobuf:"varint,8,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate          int64     `protobuf:"varint,9,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	NextExecutionDate int64    `protobuf:"varint,10,opt,name=next_execution_date,json=nextExecutionDate,proto3" json:"next_execution_date,omitempty"`
	IsActive         bool      `protobuf:"varint,11,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	TotalDonated     sdk.Coins `protobuf:"bytes,12,rep,name=total_donated,json=totalDonated,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_donated"`
	ExecutionCount   uint64    `protobuf:"varint,13,opt,name=execution_count,json=executionCount,proto3" json:"execution_count,omitempty"`
	LastExecutionDate int64    `protobuf:"varint,14,opt,name=last_execution_date,json=lastExecutionDate,proto3" json:"last_execution_date,omitempty"`
	Purpose          string    `protobuf:"bytes,15,opt,name=purpose,proto3" json:"purpose,omitempty"`
	Category         string    `protobuf:"bytes,16,opt,name=category,proto3" json:"category,omitempty"`
	CreatedAt        int64     `protobuf:"varint,17,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt        int64     `protobuf:"varint,18,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	CancelledAt      int64     `protobuf:"varint,19,opt,name=cancelled_at,json=cancelledAt,proto3" json:"cancelled_at,omitempty"`
	CancelReason     string    `protobuf:"bytes,20,opt,name=cancel_reason,json=cancelReason,proto3" json:"cancel_reason,omitempty"`
}

// ValidateRecurringDonation validates a recurring donation
func ValidateRecurringDonation(recurring RecurringDonation) error {
	if len(recurring.Donor) == 0 {
		return fmt.Errorf("donor address cannot be empty")
	}

	if recurring.NgoWalletId == 0 {
		return fmt.Errorf("NGO wallet ID cannot be zero")
	}

	if recurring.Amount.IsZero() {
		return fmt.Errorf("recurring amount cannot be zero")
	}

	if len(recurring.Frequency) == 0 {
		return fmt.Errorf("frequency cannot be empty")
	}

	if recurring.StartDate <= 0 {
		return fmt.Errorf("start date must be positive")
	}

	if recurring.EndDate > 0 && recurring.EndDate <= recurring.StartDate {
		return fmt.Errorf("end date must be after start date")
	}

	if recurring.CreatedAt <= 0 {
		return fmt.Errorf("creation timestamp must be positive")
	}

	return nil
}

// EmergencyPause represents the emergency pause state
type EmergencyPause struct {
	IsPaused     bool   `protobuf:"varint,1,opt,name=is_paused,json=isPaused,proto3" json:"is_paused,omitempty"`
	PausedBy     string `protobuf:"bytes,2,opt,name=paused_by,json=pausedBy,proto3" json:"paused_by,omitempty"`
	PausedAt     int64  `protobuf:"varint,3,opt,name=paused_at,json=pausedAt,proto3" json:"paused_at,omitempty"`
	Reason       string `protobuf:"bytes,4,opt,name=reason,proto3" json:"reason,omitempty"`
	ExpectedDuration int64 `protobuf:"varint,5,opt,name=expected_duration,json=expectedDuration,proto3" json:"expected_duration,omitempty"`
	ResumedAt    int64  `protobuf:"varint,6,opt,name=resumed_at,json=resumedAt,proto3" json:"resumed_at,omitempty"`
	ResumedBy    string `protobuf:"bytes,7,opt,name=resumed_by,json=resumedBy,proto3" json:"resumed_by,omitempty"`
}

// ValidateEmergencyPause validates an emergency pause
func ValidateEmergencyPause(pause EmergencyPause) error {
	if pause.IsPaused {
		if len(pause.PausedBy) == 0 {
			return fmt.Errorf("paused by cannot be empty when paused")
		}
		if pause.PausedAt <= 0 {
			return fmt.Errorf("paused timestamp must be positive")
		}
		if len(pause.Reason) == 0 {
			return fmt.Errorf("pause reason cannot be empty")
		}
	}
	return nil
}

// Statistics represents donation module statistics
type Statistics struct {
	TotalDonations         sdk.Coins `protobuf:"bytes,1,rep,name=total_donations,json=totalDonations,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_donations"`
	TotalDistributed       sdk.Coins `protobuf:"bytes,2,rep,name=total_distributed,json=totalDistributed,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_distributed"`
	TotalDonors            uint64    `protobuf:"varint,3,opt,name=total_donors,json=totalDonors,proto3" json:"total_donors,omitempty"`
	TotalBeneficiaries     uint64    `protobuf:"varint,4,opt,name=total_beneficiaries,json=totalBeneficiaries,proto3" json:"total_beneficiaries,omitempty"`
	TotalNGOs              uint64    `protobuf:"varint,5,opt,name=total_ngos,json=totalNgos,proto3" json:"total_ngos,omitempty"`
	ActiveNGOs             uint64    `protobuf:"varint,6,opt,name=active_ngos,json=activeNgos,proto3" json:"active_ngos,omitempty"`
	TotalCampaigns         uint64    `protobuf:"varint,7,opt,name=total_campaigns,json=totalCampaigns,proto3" json:"total_campaigns,omitempty"`
	ActiveCampaigns        uint64    `protobuf:"varint,8,opt,name=active_campaigns,json=activeCampaigns,proto3" json:"active_campaigns,omitempty"`
	TotalRecurringDonations uint64   `protobuf:"varint,9,opt,name=total_recurring_donations,json=totalRecurringDonations,proto3" json:"total_recurring_donations,omitempty"`
	ActiveRecurringDonations uint64  `protobuf:"varint,10,opt,name=active_recurring_donations,json=activeRecurringDonations,proto3" json:"active_recurring_donations,omitempty"`
	AverageTransparencyScore float64 `protobuf:"fixed64,11,opt,name=average_transparency_score,json=averageTransparencyScore,proto3" json:"average_transparency_score,omitempty"`
	UtilizationRate        float64   `protobuf:"fixed64,12,opt,name=utilization_rate,json=utilizationRate,proto3" json:"utilization_rate,omitempty"`
	LastUpdated            int64     `protobuf:"varint,13,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
}

// ValidateStatistics validates statistics
func ValidateStatistics(stats Statistics) error {
	if stats.LastUpdated <= 0 {
		return fmt.Errorf("last updated timestamp must be positive")
	}
	return nil
}

// FundFlow represents a fund flow record
type FundFlow struct {
	Id              uint64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	FlowType        string    `protobuf:"bytes,2,opt,name=flow_type,json=flowType,proto3" json:"flow_type,omitempty"`
	FromAddress     string    `protobuf:"bytes,3,opt,name=from_address,json=fromAddress,proto3" json:"from_address,omitempty"`
	ToAddress       string    `protobuf:"bytes,4,opt,name=to_address,json=toAddress,proto3" json:"to_address,omitempty"`
	Amount          sdk.Coins `protobuf:"bytes,5,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	Purpose         string    `protobuf:"bytes,6,opt,name=purpose,proto3" json:"purpose,omitempty"`
	TransactionHash string    `protobuf:"bytes,7,opt,name=transaction_hash,json=transactionHash,proto3" json:"transaction_hash,omitempty"`
	Timestamp       int64     `protobuf:"varint,8,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	BlockHeight     int64     `protobuf:"varint,9,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	RelatedId       uint64    `protobuf:"varint,10,opt,name=related_id,json=relatedId,proto3" json:"related_id,omitempty"`
	RelatedType     string    `protobuf:"bytes,11,opt,name=related_type,json=relatedType,proto3" json:"related_type,omitempty"`
}

// ValidateFundFlow validates a fund flow
func ValidateFundFlow(flow FundFlow) error {
	if len(flow.FlowType) == 0 {
		return fmt.Errorf("flow type cannot be empty")
	}

	if len(flow.FromAddress) == 0 {
		return fmt.Errorf("from address cannot be empty")
	}

	if len(flow.ToAddress) == 0 {
		return fmt.Errorf("to address cannot be empty")
	}

	if flow.Amount.IsZero() {
		return fmt.Errorf("amount cannot be zero")
	}

	if flow.Timestamp <= 0 {
		return fmt.Errorf("timestamp must be positive")
	}

	if flow.BlockHeight < 0 {
		return fmt.Errorf("block height cannot be negative")
	}

	return nil
}

// TransparencyScore represents the transparency score for an NGO
type TransparencyScore struct {
	NgoWalletId             uint64  `protobuf:"varint,1,opt,name=ngo_wallet_id,json=ngoWalletId,proto3" json:"ngo_wallet_id,omitempty"`
	Score                   int32   `protobuf:"varint,2,opt,name=score,proto3" json:"score,omitempty"`
	AuditCompleteness       int32   `protobuf:"varint,3,opt,name=audit_completeness,json=auditCompleteness,proto3" json:"audit_completeness,omitempty"`
	ReportingFrequency      int32   `protobuf:"varint,4,opt,name=reporting_frequency,json=reportingFrequency,proto3" json:"reporting_frequency,omitempty"`
	DocumentationQuality    int32   `protobuf:"varint,5,opt,name=documentation_quality,json=documentationQuality,proto3" json:"documentation_quality,omitempty"`
	FundUtilization         int32   `protobuf:"varint,6,opt,name=fund_utilization,json=fundUtilization,proto3" json:"fund_utilization,omitempty"`
	BeneficiaryFeedback     int32   `protobuf:"varint,7,opt,name=beneficiary_feedback,json=beneficiaryFeedback,proto3" json:"beneficiary_feedback,omitempty"`
	ResponseTime            int32   `protobuf:"varint,8,opt,name=response_time,json=responseTime,proto3" json:"response_time,omitempty"`
	PublicAccessibility     int32   `protobuf:"varint,9,opt,name=public_accessibility,json=publicAccessibility,proto3" json:"public_accessibility,omitempty"`
	ComplianceAdherence     int32   `protobuf:"varint,10,opt,name=compliance_adherence,json=complianceAdherence,proto3" json:"compliance_adherence,omitempty"`
	LastCalculated          int64   `protobuf:"varint,11,opt,name=last_calculated,json=lastCalculated,proto3" json:"last_calculated,omitempty"`
	NextCalculationDue      int64   `protobuf:"varint,12,opt,name=next_calculation_due,json=nextCalculationDue,proto3" json:"next_calculation_due,omitempty"`
}

// ValidateTransparencyScore validates a transparency score
func ValidateTransparencyScore(score TransparencyScore) error {
	if score.NgoWalletId == 0 {
		return fmt.Errorf("NGO wallet ID cannot be zero")
	}

	if score.Score < 1 || score.Score > 10 {
		return fmt.Errorf("score must be between 1 and 10")
	}

	if score.LastCalculated <= 0 {
		return fmt.Errorf("last calculated timestamp must be positive")
	}

	return nil
}

// VerificationQueueItem represents an item in the verification queue
type VerificationQueueItem struct {
	NgoWalletId      uint64   `protobuf:"varint,1,opt,name=ngo_wallet_id,json=ngoWalletId,proto3" json:"ngo_wallet_id,omitempty"`
	VerificationType string   `protobuf:"bytes,2,opt,name=verification_type,json=verificationType,proto3" json:"verification_type,omitempty"`
	RequestedAt      int64    `protobuf:"varint,3,opt,name=requested_at,json=requestedAt,proto3" json:"requested_at,omitempty"`
	RequestedBy      string   `protobuf:"bytes,4,opt,name=requested_by,json=requestedBy,proto3" json:"requested_by,omitempty"`
	Priority         string   `protobuf:"bytes,5,opt,name=priority,proto3" json:"priority,omitempty"`
	Status           string   `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
	AssignedTo       string   `protobuf:"bytes,7,opt,name=assigned_to,json=assignedTo,proto3" json:"assigned_to,omitempty"`
	AssignedAt       int64    `protobuf:"varint,8,opt,name=assigned_at,json=assignedAt,proto3" json:"assigned_at,omitempty"`
	CompletedAt      int64    `protobuf:"varint,9,opt,name=completed_at,json=completedAt,proto3" json:"completed_at,omitempty"`
	Result           string   `protobuf:"bytes,10,opt,name=result,proto3" json:"result,omitempty"`
	Notes            string   `protobuf:"bytes,11,opt,name=notes,proto3" json:"notes,omitempty"`
	DocumentsHash    []string `protobuf:"bytes,12,rep,name=documents_hash,json=documentsHash,proto3" json:"documents_hash,omitempty"`
}

// ValidateVerificationQueueItem validates a verification queue item
func ValidateVerificationQueueItem(item VerificationQueueItem) error {
	if item.NgoWalletId == 0 {
		return fmt.Errorf("NGO wallet ID cannot be zero")
	}

	if len(item.VerificationType) == 0 {
		return fmt.Errorf("verification type cannot be empty")
	}

	if item.RequestedAt <= 0 {
		return fmt.Errorf("requested timestamp must be positive")
	}

	if len(item.RequestedBy) == 0 {
		return fmt.Errorf("requested by cannot be empty")
	}

	return nil
}