package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NGOWallet represents an NGO receiving donations
type NGOWallet struct {
	Id                    uint64         `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                  string         `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Address               string         `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	Category              string         `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	Description           string         `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	RegistrationNumber    string         `protobuf:"bytes,6,opt,name=registration_number,json=registrationNumber,proto3" json:"registration_number,omitempty"`
	Website               string         `protobuf:"bytes,7,opt,name=website,proto3" json:"website,omitempty"`
	ContactEmail          string         `protobuf:"bytes,8,opt,name=contact_email,json=contactEmail,proto3" json:"contact_email,omitempty"`
	Signers               []string       `protobuf:"bytes,9,rep,name=signers,proto3" json:"signers,omitempty"`
	RequiredSignatures    uint32         `protobuf:"varint,10,opt,name=required_signatures,json=requiredSignatures,proto3" json:"required_signatures,omitempty"`
	IsVerified            bool           `protobuf:"varint,11,opt,name=is_verified,json=isVerified,proto3" json:"is_verified,omitempty"`
	IsActive              bool           `protobuf:"varint,12,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	TotalReceived         sdk.Coins      `protobuf:"bytes,13,rep,name=total_received,json=totalReceived,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_received"`
	TotalDistributed      sdk.Coins      `protobuf:"bytes,14,rep,name=total_distributed,json=totalDistributed,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_distributed"`
	CurrentBalance        sdk.Coins      `protobuf:"bytes,15,rep,name=current_balance,json=currentBalance,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"current_balance"`
	VerificationDocuments []string       `protobuf:"bytes,16,rep,name=verification_documents,json=verificationDocuments,proto3" json:"verification_documents,omitempty"`
	TaxExemptNumber       string         `protobuf:"bytes,17,opt,name=tax_exempt_number,json=taxExemptNumber,proto3" json:"tax_exempt_number,omitempty"`
	EightyGNumber         string         `protobuf:"bytes,18,opt,name=eighty_g_number,json=eightyGNumber,proto3" json:"eighty_g_number,omitempty"`
	CreatedAt             int64          `protobuf:"varint,19,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt             int64          `protobuf:"varint,20,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	VerifiedBy            string         `protobuf:"bytes,21,opt,name=verified_by,json=verifiedBy,proto3" json:"verified_by,omitempty"`
	AuditFrequency        uint32         `protobuf:"varint,22,opt,name=audit_frequency,json=auditFrequency,proto3" json:"audit_frequency,omitempty"`
	LastAuditDate         int64          `protobuf:"varint,23,opt,name=last_audit_date,json=lastAuditDate,proto3" json:"last_audit_date,omitempty"`
	NextAuditDue          int64          `protobuf:"varint,24,opt,name=next_audit_due,json=nextAuditDue,proto3" json:"next_audit_due,omitempty"`
	TransparencyScore     int32          `protobuf:"varint,25,opt,name=transparency_score,json=transparencyScore,proto3" json:"transparency_score,omitempty"`
	ImpactMetrics         []ImpactMetric `protobuf:"bytes,26,rep,name=impact_metrics,json=impactMetrics,proto3" json:"impact_metrics"`
	BeneficiaryCount      uint64         `protobuf:"varint,27,opt,name=beneficiary_count,json=beneficiaryCount,proto3" json:"beneficiary_count,omitempty"`
	RegionsServed         []string       `protobuf:"bytes,28,rep,name=regions_served,json=regionsServed,proto3" json:"regions_served,omitempty"`
	PriorityAreas         []string       `protobuf:"bytes,29,rep,name=priority_areas,json=priorityAreas,proto3" json:"priority_areas,omitempty"`
}

// ImpactMetric represents a measurable impact metric for an NGO
type ImpactMetric struct {
	MetricName        string `protobuf:"bytes,1,opt,name=metric_name,json=metricName,proto3" json:"metric_name,omitempty"`
	MetricValue       string `protobuf:"bytes,2,opt,name=metric_value,json=metricValue,proto3" json:"metric_value,omitempty"`
	MetricUnit        string `protobuf:"bytes,3,opt,name=metric_unit,json=metricUnit,proto3" json:"metric_unit,omitempty"`
	TargetValue       string `protobuf:"bytes,4,opt,name=target_value,json=targetValue,proto3" json:"target_value,omitempty"`
	MeasurementPeriod string `protobuf:"bytes,5,opt,name=measurement_period,json=measurementPeriod,proto3" json:"measurement_period,omitempty"`
	LastUpdated       int64  `protobuf:"varint,6,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty"`
}

// Constants for NGO wallets
const (
	// NGO Categories
	CategoryArmyWelfare       = "army_welfare"
	CategoryWarRelief         = "war_relief"
	CategoryDisabledSoldiers  = "disabled_soldiers"
	CategoryBorderAreaSchools = "border_area_schools"
	CategoryMartyrsChildren   = "martyrs_children"
	CategoryDisasterRelief    = "disaster_relief"
	
	// Default addresses (these would be replaced with actual addresses in production)
	ArmyWelfareWalletAddress       = "desh1armywelfare0000000000000000000000000000"
	WarReliefWalletAddress         = "desh1warrelief000000000000000000000000000000"
	DisabledSoldiersWalletAddress  = "desh1disabledsoldiers000000000000000000000000"
	BorderAreaSchoolsWalletAddress = "desh1borderareaschools00000000000000000000000"
	MartyrsChildrenWalletAddress   = "desh1martyrschildren0000000000000000000000000"
	DisasterReliefWalletAddress    = "desh1disasterrelief00000000000000000000000000"
	
	// Default configuration
	DefaultRequiredSignatures = 3
	DefaultAuditFrequency     = 12 // months
	MinSigners                = 3
	MaxSigners                = 10
	
	// Impact metric types
	MetricTypeBeneficiaries    = "beneficiaries_served"
	MetricTypeFundsUtilized    = "funds_utilization_rate"
	MetricTypeProjectsCompleted = "projects_completed"
)

// Default signers for government-controlled NGOs
var DefaultGovernmentSigners = []string{
	"desh1gov1000000000000000000000000000000000000",
	"desh1gov2000000000000000000000000000000000000",
	"desh1gov3000000000000000000000000000000000000",
	"desh1gov4000000000000000000000000000000000000",
	"desh1gov5000000000000000000000000000000000000",
}

// Default verification authorities
var DefaultVerificationAuthorities = []string{
	"Ministry of Home Affairs",
	"Ministry of Defence",
	"Ministry of Finance",
	"Controller General of Accounts",
	"Comptroller and Auditor General of India",
}