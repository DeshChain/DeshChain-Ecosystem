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

package types_test

import (
	"testing"
	"time"

	"github.com/deshchain/deshchain/x/validator/types"
	"github.com/stretchr/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"cosmossdk.io/math"
)

func TestGeographicValidatorCreation(t *testing.T) {
	validatorAddr := sdk.ValAddress("validator1")
	
	geoValidator := types.GeographicValidator{
		ValidatorAddress: validatorAddr,
		Country:         "India",
		State:           "Maharashtra",
		City:            "Mumbai",
		Tier:            types.TierOne,
		IsRural:         false,
		VerificationStatus: types.VerificationPending,
		SubmissionTime:  time.Now(),
		Documents: []types.VerificationDocument{
			{
				Type:     types.DocumentAadhaar,
				Hash:     "abc123",
				Verified: false,
			},
		},
	}

	require.Equal(t, validatorAddr, geoValidator.ValidatorAddress)
	require.Equal(t, "India", geoValidator.Country)
	require.Equal(t, types.TierOne, geoValidator.Tier)
	require.False(t, geoValidator.IsRural)
	require.Equal(t, types.VerificationPending, geoValidator.VerificationStatus)
}

func TestTierClassification(t *testing.T) {
	tests := []struct {
		name         string
		city         string
		expectedTier types.CityTier
	}{
		{
			name:         "Tier 1 City - Mumbai",
			city:         "Mumbai",
			expectedTier: types.TierOne,
		},
		{
			name:         "Tier 1 City - Delhi",
			city:         "Delhi",
			expectedTier: types.TierOne,
		},
		{
			name:         "Tier 1 City - Bangalore",
			city:         "Bangalore",
			expectedTier: types.TierOne,
		},
		{
			name:         "Tier 2 City - Jaipur",
			city:         "Jaipur",
			expectedTier: types.TierTwo,
		},
		{
			name:         "Tier 2 City - Lucknow",
			city:         "Lucknow",
			expectedTier: types.TierTwo,
		},
		{
			name:         "Tier 3 City - Unknown",
			city:         "UnknownCity",
			expectedTier: types.TierThree,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tier := types.ClassifyCity(tt.city)
			require.Equal(t, tt.expectedTier, tier, "City tier classification mismatch for %s", tt.city)
		})
	}
}

func TestBonusCalculation(t *testing.T) {
	tests := []struct {
		name         string
		tier         types.CityTier
		isRural      bool
		expectedRate math.LegacyDec
	}{
		{
			name:         "Tier 1 City Bonus",
			tier:         types.TierOne,
			isRural:      false,
			expectedRate: math.LegacyNewDecWithPrec(15, 3), // 1.5%
		},
		{
			name:         "Tier 2 City Bonus",
			tier:         types.TierTwo,
			isRural:      false,
			expectedRate: math.LegacyNewDecWithPrec(20, 3), // 2.0%
		},
		{
			name:         "Tier 3 City Bonus",
			tier:         types.TierThree,
			isRural:      false,
			expectedRate: math.LegacyNewDecWithPrec(25, 3), // 2.5%
		},
		{
			name:         "Rural Area Bonus",
			tier:         types.TierThree,
			isRural:      true,
			expectedRate: math.LegacyNewDecWithPrec(30, 3), // 3.0%
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate := types.CalculateGeographicBonus(tt.tier, tt.isRural)
			require.Equal(t, tt.expectedRate, rate, "Geographic bonus rate mismatch for %s", tt.name)
		})
	}
}

func TestKYCDocumentValidation(t *testing.T) {
	tests := []struct {
		name       string
		documents  []types.VerificationDocument
		shouldPass bool
	}{
		{
			name: "Valid KYC Documents",
			documents: []types.VerificationDocument{
				{Type: types.DocumentAadhaar, Hash: "aadhaar123", Verified: true},
				{Type: types.DocumentPAN, Hash: "pan123", Verified: true},
				{Type: types.DocumentAddress, Hash: "address123", Verified: true},
			},
			shouldPass: true,
		},
		{
			name: "Missing Aadhaar",
			documents: []types.VerificationDocument{
				{Type: types.DocumentPAN, Hash: "pan123", Verified: true},
				{Type: types.DocumentAddress, Hash: "address123", Verified: true},
			},
			shouldPass: false,
		},
		{
			name: "Unverified Documents",
			documents: []types.VerificationDocument{
				{Type: types.DocumentAadhaar, Hash: "aadhaar123", Verified: false},
				{Type: types.DocumentPAN, Hash: "pan123", Verified: true},
			},
			shouldPass: false,
		},
		{
			name:       "No Documents",
			documents:  []types.VerificationDocument{},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := types.ValidateKYCDocuments(tt.documents)
			require.Equal(t, tt.shouldPass, isValid, "KYC validation mismatch for %s", tt.name)
		})
	}
}

func TestVerificationStatusTransition(t *testing.T) {
	geoValidator := types.GeographicValidator{
		ValidatorAddress:   sdk.ValAddress("validator1"),
		VerificationStatus: types.VerificationPending,
	}

	// Test pending to approved
	err := geoValidator.UpdateVerificationStatus(types.VerificationApproved, "All documents verified")
	require.NoError(t, err)
	require.Equal(t, types.VerificationApproved, geoValidator.VerificationStatus)

	// Test approved to rejected (should fail)
	err = geoValidator.UpdateVerificationStatus(types.VerificationRejected, "Cannot downgrade")
	require.Error(t, err)
	require.Equal(t, types.VerificationApproved, geoValidator.VerificationStatus) // Should remain approved
}

func TestIndiaSpecificValidation(t *testing.T) {
	tests := []struct {
		name    string
		country string
		valid   bool
	}{
		{
			name:    "Valid India",
			country: "India",
			valid:   true,
		},
		{
			name:    "Valid Bharat",
			country: "Bharat",
			valid:   true,
		},
		{
			name:    "Invalid Country",
			country: "USA",
			valid:   false,
		},
		{
			name:    "Empty Country",
			country: "",
			valid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := types.IsValidIndianLocation(tt.country)
			require.Equal(t, tt.valid, isValid, "India validation mismatch for %s", tt.country)
		})
	}
}

func TestStateValidation(t *testing.T) {
	validStates := []string{
		"Maharashtra", "Karnataka", "Tamil Nadu", "Gujarat", "Rajasthan",
		"West Bengal", "Uttar Pradesh", "Madhya Pradesh", "Bihar", "Odisha",
		"Telangana", "Andhra Pradesh", "Kerala", "Haryana", "Punjab",
		"Chhattisgarh", "Jharkhand", "Uttarakhand", "Himachal Pradesh",
		"Tripura", "Meghalaya", "Manipur", "Nagaland", "Mizoram",
		"Arunachal Pradesh", "Sikkim", "Assam", "Goa",
		"Delhi", "Jammu and Kashmir", "Ladakh", "Chandigarh",
		"Dadra and Nagar Haveli and Daman and Diu", "Puducherry",
		"Andaman and Nicobar Islands", "Lakshadweep",
	}

	for _, state := range validStates {
		t.Run("Valid State: "+state, func(t *testing.T) {
			isValid := types.IsValidIndianState(state)
			require.True(t, isValid, "State should be valid: %s", state)
		})
	}

	// Test invalid states
	invalidStates := []string{"California", "Texas", "Ontario", ""}
	for _, state := range invalidStates {
		t.Run("Invalid State: "+state, func(t *testing.T) {
			isValid := types.IsValidIndianState(state)
			require.False(t, isValid, "State should be invalid: %s", state)
		})
	}
}

func TestGeographicBonusApplication(t *testing.T) {
	baseReward := math.NewInt(1000_000_000) // 1000 NAMO

	tests := []struct {
		name           string
		tier           types.CityTier
		isRural        bool
		expectedBonus  math.Int
	}{
		{
			name:          "Tier 1 Bonus",
			tier:          types.TierOne,
			isRural:       false,
			expectedBonus: math.NewInt(15_000_000), // 15 NAMO (1.5% of 1000)
		},
		{
			name:          "Tier 2 Bonus",
			tier:          types.TierTwo,
			isRural:       false,
			expectedBonus: math.NewInt(20_000_000), // 20 NAMO (2.0% of 1000)
		},
		{
			name:          "Tier 3 Bonus",
			tier:          types.TierThree,
			isRural:       false,
			expectedBonus: math.NewInt(25_000_000), // 25 NAMO (2.5% of 1000)
		},
		{
			name:          "Rural Bonus",
			tier:          types.TierThree,
			isRural:       true,
			expectedBonus: math.NewInt(30_000_000), // 30 NAMO (3.0% of 1000)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bonusRate := types.CalculateGeographicBonus(tt.tier, tt.isRural)
			actualBonus := bonusRate.MulInt(baseReward).TruncateInt()
			require.Equal(t, tt.expectedBonus, actualBonus, "Geographic bonus calculation mismatch for %s", tt.name)
		})
	}
}

func TestVerificationExpiry(t *testing.T) {
	now := time.Now()
	
	geoValidator := types.GeographicValidator{
		ValidatorAddress:   sdk.ValAddress("validator1"),
		VerificationStatus: types.VerificationApproved,
		VerificationTime:   now.Add(-13 * 30 * 24 * time.Hour), // 13 months ago
		ExpiryTime:         now.Add(-1 * 30 * 24 * time.Hour),  // 1 month ago
	}

	// Check if verification is expired
	isExpired := geoValidator.IsVerificationExpired(now)
	require.True(t, isExpired, "Verification should be expired")

	// Test non-expired verification
	geoValidator.ExpiryTime = now.Add(6 * 30 * 24 * time.Hour) // 6 months in future
	isExpired = geoValidator.IsVerificationExpired(now)
	require.False(t, isExpired, "Verification should not be expired")
}

func TestDocumentHashing(t *testing.T) {
	document := types.VerificationDocument{
		Type:     types.DocumentAadhaar,
		Hash:     "abc123def456",
		Verified: true,
	}

	// Test hash validation
	require.True(t, len(document.Hash) > 0, "Document hash should not be empty")
	require.True(t, types.IsValidDocumentHash(document.Hash), "Document hash should be valid")

	// Test invalid hashes
	invalidHashes := []string{"", "abc", "123"}
	for _, hash := range invalidHashes {
		require.False(t, types.IsValidDocumentHash(hash), "Hash should be invalid: %s", hash)
	}
}

func TestRuralClassification(t *testing.T) {
	tests := []struct {
		name       string
		population uint64
		isRural    bool
	}{
		{
			name:       "Rural Area",
			population: 25000,
			isRural:    true,
		},
		{
			name:       "Urban Area",
			population: 75000,
			isRural:    false,
		},
		{
			name:       "Boundary Case - Rural",
			population: 50000,
			isRural:    true,
		},
		{
			name:       "Boundary Case - Urban",
			population: 50001,
			isRural:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isRural := types.ClassifyRuralArea(tt.population)
			require.Equal(t, tt.isRural, isRural, "Rural classification mismatch for %s", tt.name)
		})
	}
}

func TestGeographicValidatorSerialization(t *testing.T) {
	geoValidator := types.GeographicValidator{
		ValidatorAddress:   sdk.ValAddress("validator1"),
		Country:           "India",
		State:             "Maharashtra",
		City:              "Mumbai",
		Tier:              types.TierOne,
		IsRural:           false,
		VerificationStatus: types.VerificationApproved,
		SubmissionTime:    time.Now(),
		VerificationTime:  time.Now(),
		ExpiryTime:        time.Now().Add(365 * 24 * time.Hour),
		Documents: []types.VerificationDocument{
			{
				Type:     types.DocumentAadhaar,
				Hash:     "aadhaar123",
				Verified: true,
			},
		},
		IPAddress:    "103.25.15.1",
		LastUpdated:  time.Now(),
	}

	// Test validation
	err := geoValidator.Validate()
	require.NoError(t, err, "Geographic validator should be valid")

	// Test required fields
	emptyValidator := types.GeographicValidator{}
	err = emptyValidator.Validate()
	require.Error(t, err, "Empty validator should be invalid")
}

func TestIPGeolocation(t *testing.T) {
	tests := []struct {
		name      string
		ip        string
		isIndian  bool
	}{
		{
			name:     "Indian IP Range 1",
			ip:       "103.25.15.1",
			isIndian: true,
		},
		{
			name:     "Indian IP Range 2", 
			ip:       "117.239.240.1",
			isIndian: true,
		},
		{
			name:     "US IP",
			ip:       "8.8.8.8",
			isIndian: false,
		},
		{
			name:     "Invalid IP",
			ip:       "invalid",
			isIndian: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isIndian := types.IsIndianIP(tt.ip)
			require.Equal(t, tt.isIndian, isIndian, "IP geolocation mismatch for %s", tt.ip)
		})
	}
}