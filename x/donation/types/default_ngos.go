package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultNGOWallets returns the default NGO wallets for the donation system
func DefaultNGOWallets() []NGOWallet {
	now := time.Now().Unix()
	
	return []NGOWallet{
		{
			Id:                   1,
			Name:                 "Army Welfare Fund",
			Address:              ArmyWelfareWalletAddress,
			Category:             CategoryArmyWelfare,
			Description:          "Dedicated to supporting the welfare of Indian Army personnel and their families, providing financial assistance, healthcare, education, and rehabilitation services.",
			RegistrationNumber:   "AWF/2023/001",
			Website:              "https://armywelfare.gov.in",
			ContactEmail:         "contact@armywelfare.gov.in",
			Signers:              DefaultGovernmentSigners,
			RequiredSignatures:   DefaultRequiredSignatures,
			IsVerified:           true,
			IsActive:             true,
			TotalReceived:        sdk.NewCoin("namo", sdk.ZeroInt()),
			TotalDistributed:     sdk.NewCoin("namo", sdk.ZeroInt()),
			CurrentBalance:       sdk.NewCoin("namo", sdk.ZeroInt()),
			VerificationDocuments: []string{
				"QmArmyWelfare1", "QmArmyWelfare2", "QmArmyWelfare3",
			},
			TaxExemptNumber:      "TAX/ARMY/2023/001",
			EightyGNumber:        "80G/ARMY/2023/001",
			CreatedAt:            now,
			UpdatedAt:            now,
			VerifiedBy:           DefaultVerificationAuthorities[0],
			AuditFrequency:       DefaultAuditFrequency,
			LastAuditDate:        now - 86400*30, // 30 days ago
			NextAuditDue:         now + 86400*335, // 335 days from now
			TransparencyScore:    10,
			ImpactMetrics: []ImpactMetric{
				{
					MetricName:        MetricTypeBeneficiaries,
					MetricValue:       "50000",
					MetricUnit:        "soldiers_and_families",
					TargetValue:       "100000",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
				{
					MetricName:        MetricTypeFundsUtilized,
					MetricValue:       "85",
					MetricUnit:        "percentage",
					TargetValue:       "90",
					MeasurementPeriod: "quarterly",
					LastUpdated:       now,
				},
				{
					MetricName:        MetricTypeProjectsCompleted,
					MetricValue:       "150",
					MetricUnit:        "projects",
					TargetValue:       "200",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
			},
			BeneficiaryCount: 50000,
			RegionsServed: []string{
				"Jammu & Kashmir", "Ladakh", "Himachal Pradesh", "Punjab", "Haryana",
				"Uttarakhand", "Uttar Pradesh", "Rajasthan", "Gujarat", "Madhya Pradesh",
				"Maharashtra", "Goa", "Karnataka", "Kerala", "Tamil Nadu", "Andhra Pradesh",
				"Telangana", "Odisha", "West Bengal", "Jharkhand", "Bihar", "Assam",
				"Arunachal Pradesh", "Manipur", "Meghalaya", "Mizoram", "Nagaland",
				"Sikkim", "Tripura", "Chhattisgarh", "Delhi", "Chandigarh", "Puducherry",
				"Andaman & Nicobar Islands", "Lakshadweep", "Dadra & Nagar Haveli",
				"Daman & Diu",
			},
			PriorityAreas: []string{
				"Soldier welfare", "Family support", "Medical care", "Education",
				"Housing assistance", "Disability support", "Pension support",
				"Emergency assistance", "Skill development", "Employment support",
			},
		},
		{
			Id:                   2,
			Name:                 "War Relief Fund",
			Address:              WarReliefWalletAddress,
			Category:             CategoryWarRelief,
			Description:          "Provides immediate relief and long-term support to war-affected areas and populations, including refugees, displaced persons, and conflict zones.",
			RegistrationNumber:   "WRF/2023/002",
			Website:              "https://warrelief.gov.in",
			ContactEmail:         "contact@warrelief.gov.in",
			Signers:              DefaultGovernmentSigners,
			RequiredSignatures:   DefaultRequiredSignatures,
			IsVerified:           true,
			IsActive:             true,
			TotalReceived:        sdk.NewCoin("namo", sdk.ZeroInt()),
			TotalDistributed:     sdk.NewCoin("namo", sdk.ZeroInt()),
			CurrentBalance:       sdk.NewCoin("namo", sdk.ZeroInt()),
			VerificationDocuments: []string{
				"QmWarRelief1", "QmWarRelief2", "QmWarRelief3",
			},
			TaxExemptNumber:      "TAX/WAR/2023/002",
			EightyGNumber:        "80G/WAR/2023/002",
			CreatedAt:            now,
			UpdatedAt:            now,
			VerifiedBy:           DefaultVerificationAuthorities[0],
			AuditFrequency:       DefaultAuditFrequency,
			LastAuditDate:        now - 86400*30,
			NextAuditDue:         now + 86400*335,
			TransparencyScore:    10,
			ImpactMetrics: []ImpactMetric{
				{
					MetricName:        MetricTypeBeneficiaries,
					MetricValue:       "25000",
					MetricUnit:        "affected_individuals",
					TargetValue:       "50000",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
				{
					MetricName:        MetricTypeFundsUtilized,
					MetricValue:       "92",
					MetricUnit:        "percentage",
					TargetValue:       "95",
					MeasurementPeriod: "quarterly",
					LastUpdated:       now,
				},
				{
					MetricName:        "emergency_response_time",
					MetricValue:       "24",
					MetricUnit:        "hours",
					TargetValue:       "12",
					MeasurementPeriod: "incident_based",
					LastUpdated:       now,
				},
			},
			BeneficiaryCount: 25000,
			RegionsServed: []string{
				"Border areas", "Conflict zones", "Refugee camps", "Displaced communities",
				"Emergency response areas", "International aid zones",
			},
			PriorityAreas: []string{
				"Emergency relief", "Medical aid", "Food distribution", "Shelter",
				"Rehabilitation", "Psychological support", "Community rebuilding",
				"Child protection", "Women's safety", "Elderly care",
			},
		},
		{
			Id:                   3,
			Name:                 "Disabled Soldiers Rehabilitation Fund",
			Address:              DisabledSoldiersWalletAddress,
			Category:             CategoryDisabledSoldiers,
			Description:          "Focused on rehabilitation, medical care, and empowerment of disabled soldiers and veterans, providing comprehensive support for their reintegration into society.",
			RegistrationNumber:   "DSRF/2023/003",
			Website:              "https://disabledsoldiers.gov.in",
			ContactEmail:         "contact@disabledsoldiers.gov.in",
			Signers:              DefaultGovernmentSigners,
			RequiredSignatures:   DefaultRequiredSignatures,
			IsVerified:           true,
			IsActive:             true,
			TotalReceived:        sdk.NewCoin("namo", sdk.ZeroInt()),
			TotalDistributed:     sdk.NewCoin("namo", sdk.ZeroInt()),
			CurrentBalance:       sdk.NewCoin("namo", sdk.ZeroInt()),
			VerificationDocuments: []string{
				"QmDisabledSoldiers1", "QmDisabledSoldiers2", "QmDisabledSoldiers3",
			},
			TaxExemptNumber:      "TAX/DISABLED/2023/003",
			EightyGNumber:        "80G/DISABLED/2023/003",
			CreatedAt:            now,
			UpdatedAt:            now,
			VerifiedBy:           DefaultVerificationAuthorities[0],
			AuditFrequency:       DefaultAuditFrequency,
			LastAuditDate:        now - 86400*30,
			NextAuditDue:         now + 86400*335,
			TransparencyScore:    10,
			ImpactMetrics: []ImpactMetric{
				{
					MetricName:        MetricTypeBeneficiaries,
					MetricValue:       "12000",
					MetricUnit:        "disabled_soldiers",
					TargetValue:       "15000",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
				{
					MetricName:        "prosthetics_provided",
					MetricValue:       "2500",
					MetricUnit:        "devices",
					TargetValue:       "3000",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
				{
					MetricName:        "employment_placements",
					MetricValue:       "1800",
					MetricUnit:        "jobs",
					TargetValue:       "2500",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
			},
			BeneficiaryCount: 12000,
			RegionsServed: []string{
				"All states and union territories", "Medical facilities",
				"Rehabilitation centers", "Training institutes",
			},
			PriorityAreas: []string{
				"Medical rehabilitation", "Prosthetic devices", "Skill training",
				"Employment support", "Psychological counseling", "Family support",
				"Accessibility improvements", "Technology assistance", "Sports activities",
				"Community integration",
			},
		},
		{
			Id:                   4,
			Name:                 "Border Area Schools Fund",
			Address:              BorderAreaSchoolsWalletAddress,
			Category:             CategoryBorderAreaSchools,
			Description:          "Dedicated to improving education infrastructure and quality in border areas, providing schools, teachers, and educational resources to remote communities.",
			RegistrationNumber:   "BASF/2023/004",
			Website:              "https://borderareaschools.gov.in",
			ContactEmail:         "contact@borderareaschools.gov.in",
			Signers:              DefaultGovernmentSigners,
			RequiredSignatures:   DefaultRequiredSignatures,
			IsVerified:           true,
			IsActive:             true,
			TotalReceived:        sdk.NewCoin("namo", sdk.ZeroInt()),
			TotalDistributed:     sdk.NewCoin("namo", sdk.ZeroInt()),
			CurrentBalance:       sdk.NewCoin("namo", sdk.ZeroInt()),
			VerificationDocuments: []string{
				"QmBorderSchools1", "QmBorderSchools2", "QmBorderSchools3",
			},
			TaxExemptNumber:      "TAX/SCHOOLS/2023/004",
			EightyGNumber:        "80G/SCHOOLS/2023/004",
			CreatedAt:            now,
			UpdatedAt:            now,
			VerifiedBy:           DefaultVerificationAuthorities[0],
			AuditFrequency:       DefaultAuditFrequency,
			LastAuditDate:        now - 86400*30,
			NextAuditDue:         now + 86400*335,
			TransparencyScore:    10,
			ImpactMetrics: []ImpactMetric{
				{
					MetricName:        "schools_built",
					MetricValue:       "250",
					MetricUnit:        "schools",
					TargetValue:       "500",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
				{
					MetricName:        "students_enrolled",
					MetricValue:       "35000",
					MetricUnit:        "students",
					TargetValue:       "50000",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
				{
					MetricName:        "teachers_trained",
					MetricValue:       "1500",
					MetricUnit:        "teachers",
					TargetValue:       "2000",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
			},
			BeneficiaryCount: 35000,
			RegionsServed: []string{
				"Jammu & Kashmir", "Ladakh", "Himachal Pradesh", "Uttarakhand",
				"Punjab", "Rajasthan", "Gujarat", "West Bengal", "Assam",
				"Arunachal Pradesh", "Manipur", "Meghalaya", "Mizoram",
				"Nagaland", "Tripura", "Sikkim",
			},
			PriorityAreas: []string{
				"School infrastructure", "Teacher training", "Educational materials",
				"Technology integration", "Scholarship programs", "Nutrition programs",
				"Transportation", "Hostel facilities", "Language support",
				"Cultural preservation",
			},
		},
		{
			Id:                   5,
			Name:                 "Martyrs' Children Education Fund",
			Address:              MartyrsChildrenWalletAddress,
			Category:             CategoryMartyrsChildren,
			Description:          "Provides comprehensive educational support to children of martyred soldiers, including scholarships, school fees, and career guidance.",
			RegistrationNumber:   "MCEF/2023/005",
			Website:              "https://martyrschildren.gov.in",
			ContactEmail:         "contact@martyrschildren.gov.in",
			Signers:              DefaultGovernmentSigners,
			RequiredSignatures:   DefaultRequiredSignatures,
			IsVerified:           true,
			IsActive:             true,
			TotalReceived:        sdk.NewCoin("namo", sdk.ZeroInt()),
			TotalDistributed:     sdk.NewCoin("namo", sdk.ZeroInt()),
			CurrentBalance:       sdk.NewCoin("namo", sdk.ZeroInt()),
			VerificationDocuments: []string{
				"QmMartyrsChildren1", "QmMartyrsChildren2", "QmMartyrsChildren3",
			},
			TaxExemptNumber:      "TAX/MARTYRS/2023/005",
			EightyGNumber:        "80G/MARTYRS/2023/005",
			CreatedAt:            now,
			UpdatedAt:            now,
			VerifiedBy:           DefaultVerificationAuthorities[0],
			AuditFrequency:       DefaultAuditFrequency,
			LastAuditDate:        now - 86400*30,
			NextAuditDue:         now + 86400*335,
			TransparencyScore:    10,
			ImpactMetrics: []ImpactMetric{
				{
					MetricName:        "scholarships_provided",
					MetricValue:       "8000",
					MetricUnit:        "scholarships",
					TargetValue:       "10000",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
				{
					MetricName:        "graduation_rate",
					MetricValue:       "95",
					MetricUnit:        "percentage",
					TargetValue:       "98",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
				{
					MetricName:        "career_placements",
					MetricValue:       "3500",
					MetricUnit:        "placements",
					TargetValue:       "5000",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
			},
			BeneficiaryCount: 8000,
			RegionsServed: []string{
				"All states and union territories", "Educational institutions",
				"Training centers", "Career guidance centers",
			},
			PriorityAreas: []string{
				"Educational scholarships", "Career guidance", "Skill development",
				"Higher education support", "Vocational training", "Counseling services",
				"Leadership development", "Mentorship programs", "Technology access",
				"Cultural activities",
			},
		},
		{
			Id:                   6,
			Name:                 "Disaster Relief Fund",
			Address:              DisasterReliefWalletAddress,
			Category:             CategoryDisasterRelief,
			Description:          "Provides immediate relief and long-term rehabilitation support during natural disasters, emergencies, and calamities across India.",
			RegistrationNumber:   "DRF/2023/006",
			Website:              "https://disasterrelief.gov.in",
			ContactEmail:         "contact@disasterrelief.gov.in",
			Signers:              DefaultGovernmentSigners,
			RequiredSignatures:   DefaultRequiredSignatures,
			IsVerified:           true,
			IsActive:             true,
			TotalReceived:        sdk.NewCoin("namo", sdk.ZeroInt()),
			TotalDistributed:     sdk.NewCoin("namo", sdk.ZeroInt()),
			CurrentBalance:       sdk.NewCoin("namo", sdk.ZeroInt()),
			VerificationDocuments: []string{
				"QmDisasterRelief1", "QmDisasterRelief2", "QmDisasterRelief3",
			},
			TaxExemptNumber:      "TAX/DISASTER/2023/006",
			EightyGNumber:        "80G/DISASTER/2023/006",
			CreatedAt:            now,
			UpdatedAt:            now,
			VerifiedBy:           DefaultVerificationAuthorities[0],
			AuditFrequency:       DefaultAuditFrequency,
			LastAuditDate:        now - 86400*30,
			NextAuditDue:         now + 86400*335,
			TransparencyScore:    10,
			ImpactMetrics: []ImpactMetric{
				{
					MetricName:        "disaster_responses",
					MetricValue:       "150",
					MetricUnit:        "responses",
					TargetValue:       "200",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
				{
					MetricName:        "people_assisted",
					MetricValue:       "200000",
					MetricUnit:        "individuals",
					TargetValue:       "250000",
					MeasurementPeriod: "annual",
					LastUpdated:       now,
				},
				{
					MetricName:        "response_time",
					MetricValue:       "6",
					MetricUnit:        "hours",
					TargetValue:       "4",
					MeasurementPeriod: "incident_based",
					LastUpdated:       now,
				},
			},
			BeneficiaryCount: 200000,
			RegionsServed: []string{
				"All states and union territories", "Disaster-prone areas",
				"Coastal regions", "Seismic zones", "Flood-prone areas",
				"Drought-affected areas", "Cyclone-prone areas",
			},
			PriorityAreas: []string{
				"Emergency response", "Relief distribution", "Medical aid",
				"Temporary shelter", "Food and water", "Rescue operations",
				"Rehabilitation", "Infrastructure repair", "Psychological support",
				"Community preparedness",
			},
		},
	}
}

// GetNGOWalletByCategory returns NGO wallets filtered by category
func GetNGOWalletByCategory(category string) []NGOWallet {
	wallets := DefaultNGOWallets()
	var filtered []NGOWallet
	
	for _, wallet := range wallets {
		if wallet.Category == category {
			filtered = append(filtered, wallet)
		}
	}
	
	return filtered
}

// GetNGOWalletByAddress returns an NGO wallet by address
func GetNGOWalletByAddress(address string) *NGOWallet {
	wallets := DefaultNGOWallets()
	
	for _, wallet := range wallets {
		if wallet.Address == address {
			return &wallet
		}
	}
	
	return nil
}

// GetVerifiedNGOWallets returns only verified NGO wallets
func GetVerifiedNGOWallets() []NGOWallet {
	wallets := DefaultNGOWallets()
	var verified []NGOWallet
	
	for _, wallet := range wallets {
		if wallet.IsVerified {
			verified = append(verified, wallet)
		}
	}
	
	return verified
}

// GetActiveNGOWallets returns only active NGO wallets
func GetActiveNGOWallets() []NGOWallet {
	wallets := DefaultNGOWallets()
	var active []NGOWallet
	
	for _, wallet := range wallets {
		if wallet.IsActive {
			active = append(active, wallet)
		}
	}
	
	return active
}

// GetNGOWalletsByTransparencyScore returns NGO wallets with transparency score >= threshold
func GetNGOWalletsByTransparencyScore(threshold int32) []NGOWallet {
	wallets := DefaultNGOWallets()
	var filtered []NGOWallet
	
	for _, wallet := range wallets {
		if wallet.TransparencyScore >= threshold {
			filtered = append(filtered, wallet)
		}
	}
	
	return filtered
}

// GetNGOWalletsByRegion returns NGO wallets serving a specific region
func GetNGOWalletsByRegion(region string) []NGOWallet {
	wallets := DefaultNGOWallets()
	var filtered []NGOWallet
	
	for _, wallet := range wallets {
		for _, served := range wallet.RegionsServed {
			if served == region {
				filtered = append(filtered, wallet)
				break
			}
		}
	}
	
	return filtered
}

// GetNGOWalletStatistics returns statistics about NGO wallets
func GetNGOWalletStatistics() map[string]interface{} {
	wallets := DefaultNGOWallets()
	stats := make(map[string]interface{})
	
	categoryCount := make(map[string]int)
	var totalBeneficiaries uint64
	var totalTransparencyScore int32
	verifiedCount := 0
	activeCount := 0
	
	for _, wallet := range wallets {
		categoryCount[wallet.Category]++
		totalBeneficiaries += wallet.BeneficiaryCount
		totalTransparencyScore += wallet.TransparencyScore
		
		if wallet.IsVerified {
			verifiedCount++
		}
		if wallet.IsActive {
			activeCount++
		}
	}
	
	stats["total_ngos"] = len(wallets)
	stats["verified_ngos"] = verifiedCount
	stats["active_ngos"] = activeCount
	stats["category_distribution"] = categoryCount
	stats["total_beneficiaries"] = totalBeneficiaries
	stats["average_transparency_score"] = float64(totalTransparencyScore) / float64(len(wallets))
	stats["verification_rate"] = float64(verifiedCount) / float64(len(wallets)) * 100
	stats["activation_rate"] = float64(activeCount) / float64(len(wallets)) * 100
	
	return stats
}

// ValidateNGOWallet validates an NGO wallet
func ValidateNGOWallet(wallet NGOWallet) error {
	if len(wallet.Name) == 0 {
		return fmt.Errorf("NGO name cannot be empty")
	}
	
	if len(wallet.Address) == 0 {
		return fmt.Errorf("NGO address cannot be empty")
	}
	
	if len(wallet.Category) == 0 {
		return fmt.Errorf("NGO category cannot be empty")
	}
	
	if len(wallet.Signers) < MinSigners {
		return fmt.Errorf("NGO must have at least %d signers", MinSigners)
	}
	
	if len(wallet.Signers) > MaxSigners {
		return fmt.Errorf("NGO cannot have more than %d signers", MaxSigners)
	}
	
	if wallet.RequiredSignatures < 1 {
		return fmt.Errorf("NGO must require at least 1 signature")
	}
	
	if wallet.RequiredSignatures > uint32(len(wallet.Signers)) {
		return fmt.Errorf("required signatures cannot exceed number of signers")
	}
	
	if wallet.TransparencyScore < 1 || wallet.TransparencyScore > 10 {
		return fmt.Errorf("transparency score must be between 1 and 10")
	}
	
	if wallet.AuditFrequency < 1 || wallet.AuditFrequency > 24 {
		return fmt.Errorf("audit frequency must be between 1 and 24 months")
	}
	
	return nil
}

// GetNGOWalletsByPriorityArea returns NGO wallets by priority area
func GetNGOWalletsByPriorityArea(priorityArea string) []NGOWallet {
	wallets := DefaultNGOWallets()
	var filtered []NGOWallet
	
	for _, wallet := range wallets {
		for _, area := range wallet.PriorityAreas {
			if area == priorityArea {
				filtered = append(filtered, wallet)
				break
			}
		}
	}
	
	return filtered
}

// GetNGOWalletImpactMetrics returns impact metrics for all NGO wallets
func GetNGOWalletImpactMetrics() map[uint64][]ImpactMetric {
	wallets := DefaultNGOWallets()
	metrics := make(map[uint64][]ImpactMetric)
	
	for _, wallet := range wallets {
		metrics[wallet.Id] = wallet.ImpactMetrics
	}
	
	return metrics
}

// GetTotalBeneficiaries returns total beneficiaries across all NGOs
func GetTotalBeneficiaries() uint64 {
	wallets := DefaultNGOWallets()
	var total uint64
	
	for _, wallet := range wallets {
		total += wallet.BeneficiaryCount
	}
	
	return total
}

// GetAverageTransparencyScore returns average transparency score
func GetAverageTransparencyScore() float64 {
	wallets := DefaultNGOWallets()
	var total int32
	
	for _, wallet := range wallets {
		total += wallet.TransparencyScore
	}
	
	return float64(total) / float64(len(wallets))
}

// GetRegionsCovered returns all regions covered by NGOs
func GetRegionsCovered() []string {
	wallets := DefaultNGOWallets()
	regionMap := make(map[string]bool)
	
	for _, wallet := range wallets {
		for _, region := range wallet.RegionsServed {
			regionMap[region] = true
		}
	}
	
	var regions []string
	for region := range regionMap {
		regions = append(regions, region)
	}
	
	return regions
}

// GetPriorityAreasCovered returns all priority areas covered by NGOs
func GetPriorityAreasCovered() []string {
	wallets := DefaultNGOWallets()
	areaMap := make(map[string]bool)
	
	for _, wallet := range wallets {
		for _, area := range wallet.PriorityAreas {
			areaMap[area] = true
		}
	}
	
	var areas []string
	for area := range areaMap {
		areas = append(areas, area)
	}
	
	return areas
}