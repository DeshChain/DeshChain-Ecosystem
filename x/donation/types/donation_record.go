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

// DonationRecord represents a donation made to an NGO
type DonationRecord struct {
	Id                uint64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Donor             string    `protobuf:"bytes,2,opt,name=donor,proto3" json:"donor,omitempty"`
	NgoWalletId       uint64    `protobuf:"varint,3,opt,name=ngo_wallet_id,json=ngoWalletId,proto3" json:"ngo_wallet_id,omitempty"`
	Amount            sdk.Coins `protobuf:"bytes,4,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	Purpose           string    `protobuf:"bytes,5,opt,name=purpose,proto3" json:"purpose,omitempty"`
	Category          string    `protobuf:"bytes,6,opt,name=category,proto3" json:"category,omitempty"`
	IsAnonymous       bool      `protobuf:"varint,7,opt,name=is_anonymous,json=isAnonymous,proto3" json:"is_anonymous,omitempty"`
	DonatedAt         int64     `protobuf:"varint,8,opt,name=donated_at,json=donatedAt,proto3" json:"donated_at,omitempty"`
	TransactionHash   string    `protobuf:"bytes,9,opt,name=transaction_hash,json=transactionHash,proto3" json:"transaction_hash,omitempty"`
	ReceiptHash       string    `protobuf:"bytes,10,opt,name=receipt_hash,json=receiptHash,proto3" json:"receipt_hash,omitempty"`
	TaxBenefitAmount  sdk.Coins `protobuf:"bytes,11,rep,name=tax_benefit_amount,json=taxBenefitAmount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"tax_benefit_amount"`
	Message           string    `protobuf:"bytes,12,opt,name=message,proto3" json:"message,omitempty"`
	IsRecurring       bool      `protobuf:"varint,13,opt,name=is_recurring,json=isRecurring,proto3" json:"is_recurring,omitempty"`
	RecurringId       uint64    `protobuf:"varint,14,opt,name=recurring_id,json=recurringId,proto3" json:"recurring_id,omitempty"`
	CampaignId        uint64    `protobuf:"varint,15,opt,name=campaign_id,json=campaignId,proto3" json:"campaign_id,omitempty"`
	MatchingFundsFrom string    `protobuf:"bytes,16,opt,name=matching_funds_from,json=matchingFundsFrom,proto3" json:"matching_funds_from,omitempty"`
	MatchingAmount    sdk.Coins `protobuf:"bytes,17,rep,name=matching_amount,json=matchingAmount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"matching_amount"`
	CulturalQuoteId   uint64    `protobuf:"varint,18,opt,name=cultural_quote_id,json=culturalQuoteId,proto3" json:"cultural_quote_id,omitempty"`
	BlockHeight       int64     `protobuf:"varint,19,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

// ValidateDonationRecord validates a donation record
func ValidateDonationRecord(record DonationRecord) error {
	if len(record.Donor) == 0 {
		return fmt.Errorf("donor address cannot be empty")
	}

	if record.NgoWalletId == 0 {
		return fmt.Errorf("NGO wallet ID cannot be zero")
	}

	if record.Amount.IsZero() {
		return fmt.Errorf("donation amount cannot be zero")
	}

	if record.DonatedAt <= 0 {
		return fmt.Errorf("donation timestamp must be positive")
	}

	if record.BlockHeight < 0 {
		return fmt.Errorf("block height cannot be negative")
	}

	return nil
}

// DistributionRecord represents funds distributed by an NGO
type DistributionRecord struct {
	Id                 uint64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	NgoWalletId        uint64    `protobuf:"varint,2,opt,name=ngo_wallet_id,json=ngoWalletId,proto3" json:"ngo_wallet_id,omitempty"`
	Recipient          string    `protobuf:"bytes,3,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Amount             sdk.Coins `protobuf:"bytes,4,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	Purpose            string    `protobuf:"bytes,5,opt,name=purpose,proto3" json:"purpose,omitempty"`
	Category           string    `protobuf:"bytes,6,opt,name=category,proto3" json:"category,omitempty"`
	ProjectName        string    `protobuf:"bytes,7,opt,name=project_name,json=projectName,proto3" json:"project_name,omitempty"`
	BeneficiaryName    string    `protobuf:"bytes,8,opt,name=beneficiary_name,json=beneficiaryName,proto3" json:"beneficiary_name,omitempty"`
	BeneficiaryContact string    `protobuf:"bytes,9,opt,name=beneficiary_contact,json=beneficiaryContact,proto3" json:"beneficiary_contact,omitempty"`
	DocumentationHash  string    `protobuf:"bytes,10,opt,name=documentation_hash,json=documentationHash,proto3" json:"documentation_hash,omitempty"`
	PhotosHash         string    `protobuf:"bytes,11,opt,name=photos_hash,json=photosHash,proto3" json:"photos_hash,omitempty"`
	VideoHash          string    `protobuf:"bytes,12,opt,name=video_hash,json=videoHash,proto3" json:"video_hash,omitempty"`
	GpsCoordinates     string    `protobuf:"bytes,13,opt,name=gps_coordinates,json=gpsCoordinates,proto3" json:"gps_coordinates,omitempty"`
	Region             string    `protobuf:"bytes,14,opt,name=region,proto3" json:"region,omitempty"`
	DistributedAt      int64     `protobuf:"varint,15,opt,name=distributed_at,json=distributedAt,proto3" json:"distributed_at,omitempty"`
	TransactionHash    string    `protobuf:"bytes,16,opt,name=transaction_hash,json=transactionHash,proto3" json:"transaction_hash,omitempty"`
	ApprovedBy         []string  `protobuf:"bytes,17,rep,name=approved_by,json=approvedBy,proto3" json:"approved_by,omitempty"`
	ApprovalDate       int64     `protobuf:"varint,18,opt,name=approval_date,json=approvalDate,proto3" json:"approval_date,omitempty"`
	ExecutedBy         string    `protobuf:"bytes,19,opt,name=executed_by,json=executedBy,proto3" json:"executed_by,omitempty"`
	Status             string    `protobuf:"bytes,20,opt,name=status,proto3" json:"status,omitempty"`
	Notes              string    `protobuf:"bytes,21,opt,name=notes,proto3" json:"notes,omitempty"`
	ImpactMeasured     bool      `protobuf:"varint,22,opt,name=impact_measured,json=impactMeasured,proto3" json:"impact_measured,omitempty"`
	ImpactMetrics      []ImpactMetric `protobuf:"bytes,23,rep,name=impact_metrics,json=impactMetrics,proto3" json:"impact_metrics"`
	BlockHeight        int64     `protobuf:"varint,24,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

// ValidateDistributionRecord validates a distribution record
func ValidateDistributionRecord(record DistributionRecord) error {
	if record.NgoWalletId == 0 {
		return fmt.Errorf("NGO wallet ID cannot be zero")
	}

	if len(record.Recipient) == 0 {
		return fmt.Errorf("recipient cannot be empty")
	}

	if record.Amount.IsZero() {
		return fmt.Errorf("distribution amount cannot be zero")
	}

	if len(record.Purpose) == 0 {
		return fmt.Errorf("distribution purpose cannot be empty")
	}

	if record.DistributedAt <= 0 {
		return fmt.Errorf("distribution timestamp must be positive")
	}

	if record.BlockHeight < 0 {
		return fmt.Errorf("block height cannot be negative")
	}

	return nil
}

// AuditReport represents an audit report for an NGO
type AuditReport struct {
	Id                  uint64         `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	NgoWalletId         uint64         `protobuf:"varint,2,opt,name=ngo_wallet_id,json=ngoWalletId,proto3" json:"ngo_wallet_id,omitempty"`
	AuditType           string         `protobuf:"bytes,3,opt,name=audit_type,json=auditType,proto3" json:"audit_type,omitempty"`
	Auditor             string         `protobuf:"bytes,4,opt,name=auditor,proto3" json:"auditor,omitempty"`
	AuditorName         string         `protobuf:"bytes,5,opt,name=auditor_name,json=auditorName,proto3" json:"auditor_name,omitempty"`
	AuditorRegistration string         `protobuf:"bytes,6,opt,name=auditor_registration,json=auditorRegistration,proto3" json:"auditor_registration,omitempty"`
	PeriodStart         int64          `protobuf:"varint,7,opt,name=period_start,json=periodStart,proto3" json:"period_start,omitempty"`
	PeriodEnd           int64          `protobuf:"varint,8,opt,name=period_end,json=periodEnd,proto3" json:"period_end,omitempty"`
	Summary             string         `protobuf:"bytes,9,opt,name=summary,proto3" json:"summary,omitempty"`
	Findings            []AuditFinding `protobuf:"bytes,10,rep,name=findings,proto3" json:"findings"`
	Recommendations     []string       `protobuf:"bytes,11,rep,name=recommendations,proto3" json:"recommendations,omitempty"`
	OverallRating       int32          `protobuf:"varint,12,opt,name=overall_rating,json=overallRating,proto3" json:"overall_rating,omitempty"`
	ComplianceScore     int32          `protobuf:"varint,13,opt,name=compliance_score,json=complianceScore,proto3" json:"compliance_score,omitempty"`
	EfficiencyScore     int32          `protobuf:"varint,14,opt,name=efficiency_score,json=efficiencyScore,proto3" json:"efficiency_score,omitempty"`
	ImpactScore         int32          `protobuf:"varint,15,opt,name=impact_score,json=impactScore,proto3" json:"impact_score,omitempty"`
	DocumentHash        string         `protobuf:"bytes,16,opt,name=document_hash,json=documentHash,proto3" json:"document_hash,omitempty"`
	SubmittedAt         int64          `protobuf:"varint,17,opt,name=submitted_at,json=submittedAt,proto3" json:"submitted_at,omitempty"`
	VerifiedBy          string         `protobuf:"bytes,18,opt,name=verified_by,json=verifiedBy,proto3" json:"verified_by,omitempty"`
	VerifiedAt          int64          `protobuf:"varint,19,opt,name=verified_at,json=verifiedAt,proto3" json:"verified_at,omitempty"`
	IsPublic            bool           `protobuf:"varint,20,opt,name=is_public,json=isPublic,proto3" json:"is_public,omitempty"`
	BlockHeight         int64          `protobuf:"varint,21,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

// AuditFinding represents a single finding in an audit report
type AuditFinding struct {
	Title               string   `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description         string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Severity            string   `protobuf:"bytes,3,opt,name=severity,proto3" json:"severity,omitempty"`
	Category            string   `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	Status              string   `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	Recommendation      string   `protobuf:"bytes,6,opt,name=recommendation,proto3" json:"recommendation,omitempty"`
	ManagementResponse  string   `protobuf:"bytes,7,opt,name=management_response,json=managementResponse,proto3" json:"management_response,omitempty"`
	TargetResolutionDate int64   `protobuf:"varint,8,opt,name=target_resolution_date,json=targetResolutionDate,proto3" json:"target_resolution_date,omitempty"`
	ActualResolutionDate int64   `protobuf:"varint,9,opt,name=actual_resolution_date,json=actualResolutionDate,proto3" json:"actual_resolution_date,omitempty"`
	ResponsiblePerson   string   `protobuf:"bytes,10,opt,name=responsible_person,json=responsiblePerson,proto3" json:"responsible_person,omitempty"`
}

// ValidateAuditReport validates an audit report
func ValidateAuditReport(report AuditReport) error {
	if report.NgoWalletId == 0 {
		return fmt.Errorf("NGO wallet ID cannot be zero")
	}

	if len(report.AuditType) == 0 {
		return fmt.Errorf("audit type cannot be empty")
	}

	if len(report.Auditor) == 0 {
		return fmt.Errorf("auditor cannot be empty")
	}

	if report.PeriodStart <= 0 || report.PeriodEnd <= 0 {
		return fmt.Errorf("audit period timestamps must be positive")
	}

	if report.PeriodEnd <= report.PeriodStart {
		return fmt.Errorf("audit period end must be after start")
	}

	if report.OverallRating < 1 || report.OverallRating > 10 {
		return fmt.Errorf("overall rating must be between 1 and 10")
	}

	if report.SubmittedAt <= 0 {
		return fmt.Errorf("submission timestamp must be positive")
	}

	if report.BlockHeight < 0 {
		return fmt.Errorf("block height cannot be negative")
	}

	return nil
}

// BeneficiaryTestimonial represents feedback from a beneficiary
type BeneficiaryTestimonial struct {
	Id                  uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	NgoWalletId         uint64 `protobuf:"varint,2,opt,name=ngo_wallet_id,json=ngoWalletId,proto3" json:"ngo_wallet_id,omitempty"`
	DistributionId      uint64 `protobuf:"varint,3,opt,name=distribution_id,json=distributionId,proto3" json:"distribution_id,omitempty"`
	BeneficiaryName     string `protobuf:"bytes,4,opt,name=beneficiary_name,json=beneficiaryName,proto3" json:"beneficiary_name,omitempty"`
	BeneficiaryContact  string `protobuf:"bytes,5,opt,name=beneficiary_contact,json=beneficiaryContact,proto3" json:"beneficiary_contact,omitempty"`
	Testimonial         string `protobuf:"bytes,6,opt,name=testimonial,proto3" json:"testimonial,omitempty"`
	Rating              int32  `protobuf:"varint,7,opt,name=rating,proto3" json:"rating,omitempty"`
	PhotoHash           string `protobuf:"bytes,8,opt,name=photo_hash,json=photoHash,proto3" json:"photo_hash,omitempty"`
	VideoHash           string `protobuf:"bytes,9,opt,name=video_hash,json=videoHash,proto3" json:"video_hash,omitempty"`
	AudioHash           string `protobuf:"bytes,10,opt,name=audio_hash,json=audioHash,proto3" json:"audio_hash,omitempty"`
	Language            string `protobuf:"bytes,11,opt,name=language,proto3" json:"language,omitempty"`
	Translation         string `protobuf:"bytes,12,opt,name=translation,proto3" json:"translation,omitempty"`
	SubmittedAt         int64  `protobuf:"varint,13,opt,name=submitted_at,json=submittedAt,proto3" json:"submitted_at,omitempty"`
	VerifiedBy          string `protobuf:"bytes,14,opt,name=verified_by,json=verifiedBy,proto3" json:"verified_by,omitempty"`
	VerifiedAt          int64  `protobuf:"varint,15,opt,name=verified_at,json=verifiedAt,proto3" json:"verified_at,omitempty"`
	VerificationMethod  string `protobuf:"bytes,16,opt,name=verification_method,json=verificationMethod,proto3" json:"verification_method,omitempty"`
	IsPublic            bool   `protobuf:"varint,17,opt,name=is_public,json=isPublic,proto3" json:"is_public,omitempty"`
	ImpactDescription   string `protobuf:"bytes,18,opt,name=impact_description,json=impactDescription,proto3" json:"impact_description,omitempty"`
	LifeChangeMeasure   string `protobuf:"bytes,19,opt,name=life_change_measure,json=lifeChangeMeasure,proto3" json:"life_change_measure,omitempty"`
	WouldRecommend      bool   `protobuf:"varint,20,opt,name=would_recommend,json=wouldRecommend,proto3" json:"would_recommend,omitempty"`
	BlockHeight         int64  `protobuf:"varint,21,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

// ValidateBeneficiaryTestimonial validates a beneficiary testimonial
func ValidateBeneficiaryTestimonial(testimonial BeneficiaryTestimonial) error {
	if testimonial.NgoWalletId == 0 {
		return fmt.Errorf("NGO wallet ID cannot be zero")
	}

	if len(testimonial.BeneficiaryName) == 0 {
		return fmt.Errorf("beneficiary name cannot be empty")
	}

	if len(testimonial.Testimonial) == 0 {
		return fmt.Errorf("testimonial text cannot be empty")
	}

	if testimonial.Rating < 1 || testimonial.Rating > 10 {
		return fmt.Errorf("rating must be between 1 and 10")
	}

	if testimonial.SubmittedAt <= 0 {
		return fmt.Errorf("submission timestamp must be positive")
	}

	if testimonial.BlockHeight < 0 {
		return fmt.Errorf("block height cannot be negative")
	}

	return nil
}