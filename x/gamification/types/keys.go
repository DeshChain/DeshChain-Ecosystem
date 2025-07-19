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

import "encoding/binary"

const (
	// ModuleName defines the module name
	ModuleName = "gamification"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_gamification"

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KeyPrefix returns the store key prefix
func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	// ProfileKey is the key for developer profiles
	ProfileKey = "Profile/value/"

	// AchievementKey is the key for achievements
	AchievementKey = "Achievement/value/"

	// LeaderboardKey is the key for leaderboards
	LeaderboardKey = "Leaderboard/value/"

	// DailyChallengeKey is the key for daily challenges
	DailyChallengeKey = "DailyChallenge/value/"

	// TeamBattleKey is the key for team battles
	TeamBattleKey = "TeamBattle/value/"

	// SocialPostKey is the key for social media posts
	SocialPostKey = "SocialPost/value/"

	// QuoteKey is the key for humor quotes
	QuoteKey = "Quote/value/"

	// ParamsKey is the key for module parameters
	ParamsKey = "Params/value/"
)

// GetProfileKey returns the store key for a profile
func GetProfileKey(address string) []byte {
	return append(KeyPrefix(ProfileKey), []byte(address)...)
}

// GetAchievementKey returns the store key for an achievement
func GetAchievementKey(id uint64) []byte {
	return append(KeyPrefix(AchievementKey), GetUint64Bytes(id)...)
}

// GetLeaderboardKey returns the store key for a leaderboard
func GetLeaderboardKey(leaderboardType string) []byte {
	return append(KeyPrefix(LeaderboardKey), []byte(leaderboardType)...)
}

// GetDailyChallengeKey returns the store key for a daily challenge
func GetDailyChallengeKey(date string) []byte {
	return append(KeyPrefix(DailyChallengeKey), []byte(date)...)
}

// GetTeamBattleKey returns the store key for a team battle
func GetTeamBattleKey(id uint64) []byte {
	return append(KeyPrefix(TeamBattleKey), GetUint64Bytes(id)...)
}

// GetSocialPostKey returns the store key for a social post
func GetSocialPostKey(id uint64) []byte {
	return append(KeyPrefix(SocialPostKey), GetUint64Bytes(id)...)
}

// GetQuoteKey returns the store key for a quote
func GetQuoteKey(id uint64) []byte {
	return append(KeyPrefix(QuoteKey), GetUint64Bytes(id)...)
}

// GetUint64Bytes returns the byte representation of a uint64
func GetUint64Bytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetUint64FromBytes returns uint64 from byte representation
func GetUint64FromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}