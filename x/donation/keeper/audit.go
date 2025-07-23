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

package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"deshchain/x/donation/types"
)

// SetAuditReport sets an audit report in the store
func (k Keeper) SetAuditReport(ctx sdk.Context, report types.AuditReport) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditReportKey)
	bz := k.cdc.MustMarshal(&report)
	store.Set(sdk.Uint64ToBigEndian(report.Id), bz)
}

// GetAuditReport retrieves an audit report from the store
func (k Keeper) GetAuditReport(ctx sdk.Context, id uint64) (types.AuditReport, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditReportKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return types.AuditReport{}, false
	}
	var report types.AuditReport
	k.cdc.MustUnmarshal(bz, &report)
	return report, true
}

// GetAllAuditReports returns all audit reports
func (k Keeper) GetAllAuditReports(ctx sdk.Context) []types.AuditReport {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditReportKey)
	var reports []types.AuditReport
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var report types.AuditReport
		k.cdc.MustUnmarshal(iterator.Value(), &report)
		reports = append(reports, report)
	}
	
	return reports
}

// SetAuditReportCount sets the total audit report count
func (k Keeper) SetAuditReportCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditReportCountKey)
	store.Set([]byte{0}, sdk.Uint64ToBigEndian(count))
}

// GetAuditReportCount gets the total audit report count
func (k Keeper) GetAuditReportCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditReportCountKey)
	bz := store.Get([]byte{0})
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// AddAuditByNGO adds an audit ID to the NGO's audit list
func (k Keeper) AddAuditByNGO(ctx sdk.Context, ngoWalletId, auditId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditByNGOKey)
	key := append(sdk.Uint64ToBigEndian(ngoWalletId), sdk.Uint64ToBigEndian(auditId)...)
	store.Set(key, []byte{1})
}

// GetAuditsByNGO retrieves all audit IDs for an NGO
func (k Keeper) GetAuditsByNGO(ctx sdk.Context, ngoWalletId uint64) []uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AuditByNGOKey)
	var auditIds []uint64
	
	ngoKey := sdk.Uint64ToBigEndian(ngoWalletId)
	iterator := store.Iterator(ngoKey, sdk.PrefixEndBytes(ngoKey))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		// Extract audit ID from key (NGO ID + audit ID)
		if len(key) >= 16 {
			auditId := sdk.BigEndianToUint64(key[8:])
			auditIds = append(auditIds, auditId)
		}
	}
	
	return auditIds
}

// SetBeneficiaryTestimonial sets a beneficiary testimonial in the store
func (k Keeper) SetBeneficiaryTestimonial(ctx sdk.Context, testimonial types.BeneficiaryTestimonial) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BeneficiaryTestimonialKey)
	bz := k.cdc.MustMarshal(&testimonial)
	store.Set(sdk.Uint64ToBigEndian(testimonial.Id), bz)
}

// GetBeneficiaryTestimonial retrieves a beneficiary testimonial from the store
func (k Keeper) GetBeneficiaryTestimonial(ctx sdk.Context, id uint64) (types.BeneficiaryTestimonial, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BeneficiaryTestimonialKey)
	bz := store.Get(sdk.Uint64ToBigEndian(id))
	if bz == nil {
		return types.BeneficiaryTestimonial{}, false
	}
	var testimonial types.BeneficiaryTestimonial
	k.cdc.MustUnmarshal(bz, &testimonial)
	return testimonial, true
}

// GetAllBeneficiaryTestimonials returns all beneficiary testimonials
func (k Keeper) GetAllBeneficiaryTestimonials(ctx sdk.Context) []types.BeneficiaryTestimonial {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BeneficiaryTestimonialKey)
	var testimonials []types.BeneficiaryTestimonial
	
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var testimonial types.BeneficiaryTestimonial
		k.cdc.MustUnmarshal(iterator.Value(), &testimonial)
		testimonials = append(testimonials, testimonial)
	}
	
	return testimonials
}

// SetBeneficiaryTestimonialCount sets the total beneficiary testimonial count
func (k Keeper) SetBeneficiaryTestimonialCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BeneficiaryTestimonialCountKey)
	store.Set([]byte{0}, sdk.Uint64ToBigEndian(count))
}

// GetBeneficiaryTestimonialCount gets the total beneficiary testimonial count
func (k Keeper) GetBeneficiaryTestimonialCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BeneficiaryTestimonialCountKey)
	bz := store.Get([]byte{0})
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// AddTestimonialByNGO adds a testimonial ID to the NGO's testimonial list
func (k Keeper) AddTestimonialByNGO(ctx sdk.Context, ngoWalletId, testimonialId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TestimonialByNGOKey)
	key := append(sdk.Uint64ToBigEndian(ngoWalletId), sdk.Uint64ToBigEndian(testimonialId)...)
	store.Set(key, []byte{1})
}

// GetTestimonialsByNGO retrieves all testimonial IDs for an NGO
func (k Keeper) GetTestimonialsByNGO(ctx sdk.Context, ngoWalletId uint64) []uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.TestimonialByNGOKey)
	var testimonialIds []uint64
	
	ngoKey := sdk.Uint64ToBigEndian(ngoWalletId)
	iterator := store.Iterator(ngoKey, sdk.PrefixEndBytes(ngoKey))
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		// Extract testimonial ID from key (NGO ID + testimonial ID)
		if len(key) >= 16 {
			testimonialId := sdk.BigEndianToUint64(key[8:])
			testimonialIds = append(testimonialIds, testimonialId)
		}
	}
	
	return testimonialIds
}

// CalculateTransparencyScore calculates the transparency score for an NGO
func (k Keeper) CalculateTransparencyScore(ctx sdk.Context, ngoWalletId uint64) (int32, error) {
	ngo, found := k.GetNGOWallet(ctx, ngoWalletId)
	if !found {
		return 0, types.ErrNGONotFound
	}
	
	// Get audit reports
	auditIds := k.GetAuditsByNGO(ctx, ngoWalletId)
	auditCompleteness := int32(0)
	if len(auditIds) > 0 {
		// Check if audits are up to date
		lastAudit := auditIds[len(auditIds)-1]
		report, _ := k.GetAuditReport(ctx, lastAudit)
		if report.SubmittedAt > ctx.BlockTime().Unix()-86400*365 { // Within last year
			auditCompleteness = 10
		} else {
			auditCompleteness = 5
		}
	}
	
	// Get distribution records
	distributionIds := k.GetDistributionsByNGO(ctx, ngoWalletId)
	reportingFrequency := int32(0)
	if len(distributionIds) > 10 {
		reportingFrequency = 10
	} else if len(distributionIds) > 5 {
		reportingFrequency = 7
	} else if len(distributionIds) > 0 {
		reportingFrequency = 5
	}
	
	// Documentation quality (based on verification status)
	documentationQuality := int32(0)
	if ngo.IsVerified {
		documentationQuality = 10
	} else {
		documentationQuality = 5
	}
	
	// Fund utilization
	fundUtilization := int32(0)
	if !ngo.TotalReceived.IsZero() && !ngo.TotalDistributed.IsZero() {
		// Calculate utilization percentage
		for _, received := range ngo.TotalReceived {
			for _, distributed := range ngo.TotalDistributed {
				if received.Denom == distributed.Denom && !received.Amount.IsZero() {
					utilization := distributed.Amount.ToDec().Quo(received.Amount.ToDec()).MustFloat64()
					if utilization >= 0.8 {
						fundUtilization = 10
					} else if utilization >= 0.6 {
						fundUtilization = 7
					} else if utilization >= 0.4 {
						fundUtilization = 5
					} else {
						fundUtilization = 3
					}
					break
				}
			}
		}
	}
	
	// Beneficiary feedback
	testimonialIds := k.GetTestimonialsByNGO(ctx, ngoWalletId)
	beneficiaryFeedback := int32(0)
	if len(testimonialIds) > 0 {
		totalRating := int32(0)
		for _, id := range testimonialIds {
			testimonial, _ := k.GetBeneficiaryTestimonial(ctx, id)
			totalRating += testimonial.Rating
		}
		avgRating := totalRating / int32(len(testimonialIds))
		beneficiaryFeedback = avgRating
	}
	
	// Response time (based on distribution frequency)
	responseTime := reportingFrequency // Use same as reporting frequency for now
	
	// Public accessibility (always high for blockchain)
	publicAccessibility := int32(10)
	
	// Compliance adherence
	complianceAdherence := int32(0)
	if ngo.IsVerified && len(auditIds) > 0 {
		complianceAdherence = 10
	} else if ngo.IsVerified || len(auditIds) > 0 {
		complianceAdherence = 7
	} else {
		complianceAdherence = 5
	}
	
	// Calculate overall score
	components := []int32{
		auditCompleteness,
		reportingFrequency,
		documentationQuality,
		fundUtilization,
		beneficiaryFeedback,
		responseTime,
		publicAccessibility,
		complianceAdherence,
	}
	
	totalScore := int32(0)
	componentCount := int32(0)
	for _, score := range components {
		if score > 0 {
			totalScore += score
			componentCount++
		}
	}
	
	overallScore := int32(5) // Default score
	if componentCount > 0 {
		overallScore = totalScore / componentCount
	}
	
	// Store the transparency score
	transparencyScore := types.TransparencyScore{
		NgoWalletId:          ngoWalletId,
		Score:                overallScore,
		AuditCompleteness:    auditCompleteness,
		ReportingFrequency:   reportingFrequency,
		DocumentationQuality: documentationQuality,
		FundUtilization:      fundUtilization,
		BeneficiaryFeedback:  beneficiaryFeedback,
		ResponseTime:         responseTime,
		PublicAccessibility:  publicAccessibility,
		ComplianceAdherence:  complianceAdherence,
		LastCalculated:       ctx.BlockTime().Unix(),
		NextCalculationDue:   ctx.BlockTime().Unix() + 86400*30, // 30 days
	}
	
	k.SetTransparencyScore(ctx, transparencyScore)
	
	// Update NGO's transparency score
	if err := k.UpdateNGOTransparencyScore(ctx, ngoWalletId, overallScore); err != nil {
		return 0, err
	}
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTransparencyUpdate,
			sdk.NewAttribute(types.AttributeKeyNGOWalletID, sdk.FormatInvariant("%d", ngoWalletId)),
			sdk.NewAttribute(types.AttributeKeyTransparencyScore, sdk.FormatInvariant("%d", overallScore)),
		),
	)
	
	return overallScore, nil
}