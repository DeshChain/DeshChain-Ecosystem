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

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/deshchain/deshchain/x/moneyorder/types"
)

func TestVillagePool_Validate(t *testing.T) {
	now := time.Now()
	validAddr := "desh1qperwt9wrnkg5k9e5gzfgjppzpqhyav5j24d66"
	validValAddr := "deshvaloper1qperwt9wrnkg5k9e5gzfgjppzpqhyav5j7z8kl"
	
	tests := []struct {
		name    string
		pool    types.VillagePool
		wantErr bool
	}{
		{
			name: "valid village pool",
			pool: types.VillagePool{
				PoolId:        1,
				VillageName:   "Rampur",
				PostalCode:    "110001",
				StateCode:     "DL",
				DistrictCode:  "ND",
				PanchayatName: "Rampur Gram Panchayat",
				PanchayatHead: validAddr,
				CooperativeId: "COOP001",
				SupportedTokens: []string{"unamo", "uinr"},
				PrimaryToken:   "unamo",
				LocalCurrency:  "INR",
				TotalLiquidity: sdk.NewCoins(
					sdk.NewCoin("unamo", sdk.NewInt(1000000)),
					sdk.NewCoin("uinr", sdk.NewInt(100000000)),
				),
				AvailableLiquidity: sdk.NewCoins(
					sdk.NewCoin("unamo", sdk.NewInt(900000)),
					sdk.NewCoin("uinr", sdk.NewInt(90000000)),
				),
				ReservedLiquidity: sdk.NewCoins(
					sdk.NewCoin("unamo", sdk.NewInt(100000)),
					sdk.NewCoin("uinr", sdk.NewInt(10000000)),
				),
				MinimumLiquidity: sdk.NewInt(100000),
				LocalValidators: []types.LocalValidator{
					{
						ValidatorAddress: validValAddr,
						LocalName:       "Shri Ram Kumar",
						Role:            "Village Elder",
						TrustLevel:      5,
						JoinedAt:        now,
					},
				},
				RequiredSignatures: 1,
				BaseTradingFee:     sdk.NewDecWithPrec(2, 3),
				CommunityFee:       sdk.NewDecWithPrec(5, 4),
				EducationFee:       sdk.NewDecWithPrec(3, 4),
				InfrastructureFee:  sdk.NewDecWithPrec(2, 4),
				TotalMembers:       100,
				ActiveTraders:      50,
				MemberBenefits: types.MemberBenefits{
					FeeDiscount:       sdk.NewDecWithPrec(30, 2),
					PriorityExecution: true,
					VotingRights:      true,
					ProfitSharing:     sdk.NewDecWithPrec(10, 2),
					EducationAccess:   true,
					EmergencySupport:  true,
				},
				TotalVolume:        sdk.NewInt(10000000),
				DailyVolume:        sdk.NewInt(100000),
				MonthlyVolume:      sdk.NewInt(1000000),
				TotalTransactions:  1000,
				CommunityFund:      sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(10000))),
				EducationFund:      sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(5000))),
				EmergencyFund:      sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(20000))),
				LocalFestivals:     []string{"Holi", "Diwali", "Dussehra"},
				LanguageSupport:    []string{"hi", "en"},
				CulturalQuotes:     []string{"वसुधैव कुटुम्बकम्"},
				Active:             true,
				Verified:           true,
				EstablishedDate:    now,
				LastActivityDate:   now,
				Achievements:       []types.VillageAchievement{},
				TrustScore:         80,
			},
			wantErr: false,
		},
		{
			name: "invalid pool id",
			pool: types.VillagePool{
				PoolId: 0,
			},
			wantErr: true,
		},
		{
			name: "empty village name",
			pool: types.VillagePool{
				PoolId:      1,
				VillageName: "",
			},
			wantErr: true,
		},
		{
			name: "invalid postal code",
			pool: types.VillagePool{
				PoolId:      1,
				VillageName: "Test Village",
				PostalCode:  "12345",
			},
			wantErr: true,
		},
		{
			name: "invalid panchayat head address",
			pool: types.VillagePool{
				PoolId:        1,
				VillageName:   "Test Village",
				PostalCode:    "110001",
				PanchayatName: "Test Panchayat",
				PanchayatHead: "invalid",
			},
			wantErr: true,
		},
		{
			name: "no supported tokens",
			pool: types.VillagePool{
				PoolId:        1,
				VillageName:   "Test Village",
				PostalCode:    "110001",
				PanchayatName: "Test Panchayat",
				PanchayatHead: validAddr,
				SupportedTokens: []string{},
			},
			wantErr: true,
		},
		{
			name: "negative balance",
			pool: types.VillagePool{
				PoolId:          1,
				VillageName:     "Test Village",
				PostalCode:      "110001",
				PanchayatName:   "Test Panchayat",
				PanchayatHead:   validAddr,
				SupportedTokens: []string{"unamo"},
				TotalLiquidity: sdk.Coins{
					sdk.Coin{Denom: "unamo", Amount: sdk.NewInt(-1000)},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid fee percentage",
			pool: types.VillagePool{
				PoolId:          1,
				VillageName:     "Test Village",
				PostalCode:      "110001",
				PanchayatName:   "Test Panchayat",
				PanchayatHead:   validAddr,
				SupportedTokens: []string{"unamo"},
				BaseTradingFee:  sdk.NewDec(2), // 200% fee
			},
			wantErr: true,
		},
		{
			name: "invalid trust score",
			pool: types.VillagePool{
				PoolId:          1,
				VillageName:     "Test Village",
				PostalCode:      "110001",
				PanchayatName:   "Test Panchayat",
				PanchayatHead:   validAddr,
				SupportedTokens: []string{"unamo"},
				TrustScore:      101, // Max is 100
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pool.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestVillagePoolMember_Validate(t *testing.T) {
	now := time.Now()
	validAddr := "desh1qperwt9wrnkg5k9e5gzfgjppzpqhyav5j24d66"
	
	tests := []struct {
		name    string
		member  types.VillagePoolMember
		wantErr bool
	}{
		{
			name: "valid member",
			member: types.VillagePoolMember{
				MemberAddress:  validAddr,
				LocalName:      "Ram Kumar",
				MobileNumber:   "+919876543210",
				AadhaarHash:    "0x1234567890abcdef",
				MembershipType: "standard",
				JoinedAt:       now,
				Contribution: sdk.NewCoins(
					sdk.NewCoin("unamo", sdk.NewInt(1000)),
				),
				TotalTrades:    10,
				TotalVolume:    sdk.NewInt(100000),
				LastTradeAt:    now,
				TotalEarnings:  sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(100))),
				PendingRewards: sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(10))),
				EducationCredits: 5,
			},
			wantErr: false,
		},
		{
			name: "invalid address",
			member: types.VillagePoolMember{
				MemberAddress: "invalid",
			},
			wantErr: true,
		},
		{
			name: "empty local name",
			member: types.VillagePoolMember{
				MemberAddress: validAddr,
				LocalName:     "",
			},
			wantErr: true,
		},
		{
			name: "invalid mobile number",
			member: types.VillagePoolMember{
				MemberAddress: validAddr,
				LocalName:     "Test",
				MobileNumber:  "123", // Too short
			},
			wantErr: true,
		},
		{
			name: "negative contribution",
			member: types.VillagePoolMember{
				MemberAddress: validAddr,
				LocalName:     "Test",
				MobileNumber:  "+919876543210",
				Contribution: sdk.Coins{
					sdk.Coin{Denom: "unamo", Amount: sdk.NewInt(-1000)},
				},
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.member.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLocalValidator_Validate(t *testing.T) {
	now := time.Now()
	validValAddr := "deshvaloper1qperwt9wrnkg5k9e5gzfgjppzpqhyav5j7z8kl"
	
	tests := []struct {
		name    string
		validator types.LocalValidator
		wantErr bool
	}{
		{
			name: "valid validator",
			validator: types.LocalValidator{
				ValidatorAddress: validValAddr,
				LocalName:       "Shri Ram Kumar",
				Role:            "Village Elder",
				TrustLevel:      5,
				JoinedAt:        now,
			},
			wantErr: false,
		},
		{
			name: "invalid address",
			validator: types.LocalValidator{
				ValidatorAddress: "invalid",
			},
			wantErr: true,
		},
		{
			name: "empty local name",
			validator: types.LocalValidator{
				ValidatorAddress: validValAddr,
				LocalName:        "",
			},
			wantErr: true,
		},
		{
			name: "invalid trust level",
			validator: types.LocalValidator{
				ValidatorAddress: validValAddr,
				LocalName:        "Test",
				Role:             "Elder",
				TrustLevel:       11, // Max is 10
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestVillageAchievement_Validate(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name    string
		achievement types.VillageAchievement
		wantErr bool
	}{
		{
			name: "valid achievement",
			achievement: types.VillageAchievement{
				AchievementId: "ACH001",
				Title:         "First 1000 Trades",
				Description:   "Village completed 1000 successful trades",
				Category:      "trading",
				AchievedAt:    now,
				RewardAmount:  sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(10000))),
			},
			wantErr: false,
		},
		{
			name: "empty achievement id",
			achievement: types.VillageAchievement{
				AchievementId: "",
			},
			wantErr: true,
		},
		{
			name: "empty title",
			achievement: types.VillageAchievement{
				AchievementId: "ACH001",
				Title:         "",
			},
			wantErr: true,
		},
		{
			name: "negative reward",
			achievement: types.VillageAchievement{
				AchievementId: "ACH001",
				Title:         "Test Achievement",
				Description:   "Test",
				RewardAmount: sdk.Coins{
					sdk.Coin{Denom: "unamo", Amount: sdk.NewInt(-1000)},
				},
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.achievement.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMemberBenefits_Validate(t *testing.T) {
	tests := []struct {
		name     string
		benefits types.MemberBenefits
		wantErr  bool
	}{
		{
			name: "valid benefits",
			benefits: types.MemberBenefits{
				FeeDiscount:       sdk.NewDecWithPrec(30, 2), // 30%
				PriorityExecution: true,
				VotingRights:      true,
				ProfitSharing:     sdk.NewDecWithPrec(10, 2), // 10%
				EducationAccess:   true,
				EmergencySupport:  true,
			},
			wantErr: false,
		},
		{
			name: "fee discount too high",
			benefits: types.MemberBenefits{
				FeeDiscount: sdk.NewDec(2), // 200%
			},
			wantErr: true,
		},
		{
			name: "negative fee discount",
			benefits: types.MemberBenefits{
				FeeDiscount: sdk.NewDec(-1),
			},
			wantErr: true,
		},
		{
			name: "profit sharing too high",
			benefits: types.MemberBenefits{
				FeeDiscount:   sdk.NewDecWithPrec(30, 2),
				ProfitSharing: sdk.NewDec(2), // 200%
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.benefits.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}