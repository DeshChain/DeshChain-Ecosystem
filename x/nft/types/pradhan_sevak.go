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
	"encoding/json"
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PradhanSevakNFT represents the unique NFT to be gifted to PMO India
type PradhanSevakNFT struct {
	// Core Identification
	TokenID      string    `json:"token_id"`
	Name         string    `json:"name"`
	Symbol       string    `json:"symbol"`
	Description  string    `json:"description"`
	
	// Tribute Details
	TributeTo    string    `json:"tribute_to"`
	Designation  string    `json:"designation"`
	TributeDate  time.Time `json:"tribute_date"`
	
	// Cultural Content
	Quotes       []CulturalQuote `json:"quotes"`
	Languages    []string        `json:"languages"`
	
	// Technical Metadata
	GenesisBlock    int64     `json:"genesis_block"`
	ChainID         string    `json:"chain_id"`
	CreationTime    time.Time `json:"creation_time"`
	CreatorAddress  string    `json:"creator_address"`
	
	// Royalty Configuration
	RoyaltyEnabled  bool      `json:"royalty_enabled"`
	RoyaltyPercent  sdk.Dec   `json:"royalty_percent"`
	RoyaltyRecipient string   `json:"royalty_recipient"`
	
	// Media and Display
	ImageURI        string    `json:"image_uri"`
	AnimationURI    string    `json:"animation_uri"`
	ExternalURI     string    `json:"external_uri"`
	
	// Special Properties
	Transferable    bool      `json:"transferable"`
	Burnable        bool      `json:"burnable"`
	Mutable         bool      `json:"mutable"`
}

// CulturalQuote represents a quote embedded in the NFT
type CulturalQuote struct {
	Text        string `json:"text"`
	Author      string `json:"author"`
	Language    string `json:"language"`
	Category    string `json:"category"`
	Translation string `json:"translation,omitempty"`
}

// CreatePradhanSevakNFT creates the unique tribute NFT
func CreatePradhanSevakNFT(chainID string, genesisHeight int64) *PradhanSevakNFT {
	return &PradhanSevakNFT{
		TokenID:      "PRADHAN-SEVAK-001",
		Name:         "Pradhan Sevak - Principal Servant of India",
		Symbol:       "SEVAK",
		Description:  "A perpetual tribute to Shri Narendra Modi Ji, Hon'ble Prime Minister of India, recognizing his transformative contributions to Digital India, financial inclusion, and national development. This unique NFT represents the gratitude of millions of Indians whose lives have been transformed through visionary leadership.",
		
		TributeTo:    "Shri Narendra Modi Ji",
		Designation:  "Hon'ble Prime Minister of India",
		TributeDate:  time.Now(),
		
		Quotes:       GetPradhanSevakQuotes(),
		Languages:    GetSupportedLanguages(),
		
		GenesisBlock:    genesisHeight,
		ChainID:         chainID,
		CreationTime:    time.Now(),
		CreatorAddress:  "deshchain1genesis", // Special genesis address
		
		RoyaltyEnabled:   true,
		RoyaltyPercent:   sdk.NewDecWithPrec(1, 4), // 0.01%
		RoyaltyRecipient: "PM CARES Fund",
		
		ImageURI:      "ipfs://QmPradhanSevakMainImage",
		AnimationURI:  "ipfs://QmPradhanSevakAnimation",
		ExternalURI:   "https://deshchain.com/nft/pradhan-sevak",
		
		Transferable:  false, // Cannot be transferred
		Burnable:      false, // Cannot be burned
		Mutable:       false, // Cannot be modified
	}
}

// GetPradhanSevakQuotes returns curated quotes for the NFT
func GetPradhanSevakQuotes() []CulturalQuote {
	return []CulturalQuote{
		{
			Text:     "सबका साथ, सबका विकास, सबका विश्वास, सबका प्रयास",
			Author:   "Narendra Modi",
			Language: "Hindi",
			Category: "Leadership",
			Translation: "Together with all, Development for all, Trust of all, Efforts of all",
		},
		{
			Text:     "Technology should be used as a tool to empower the poor",
			Author:   "Narendra Modi",
			Language: "English",
			Category: "Technology",
		},
		{
			Text:     "Digital India is not just about connectivity, it's about empowerment",
			Author:   "Narendra Modi",
			Language: "English",
			Category: "Digital India",
		},
		{
			Text:     "मैं देश का प्रधान सेवक हूं, प्रधान मंत्री नहीं",
			Author:   "Narendra Modi",
			Language: "Hindi",
			Category: "Service",
			Translation: "I am the principal servant of the nation, not just Prime Minister",
		},
		{
			Text:     "Every Indian must be financially included in the growth story",
			Author:   "Narendra Modi",
			Language: "English",
			Category: "Financial Inclusion",
		},
		// Add more quotes up to 1000
	}
}

// GetSupportedLanguages returns all 22 supported languages
func GetSupportedLanguages() []string {
	return []string{
		"Hindi", "English", "Bengali", "Telugu", "Marathi",
		"Tamil", "Gujarati", "Urdu", "Kannada", "Odia",
		"Malayalam", "Punjabi", "Assamese", "Maithili", "Sanskrit",
		"Konkani", "Nepali", "Manipuri", "Sindhi", "Dogri",
		"Kashmiri", "Bodo",
	}
}

// Validate ensures the NFT is properly configured
func (nft *PradhanSevakNFT) Validate() error {
	if nft.TokenID == "" {
		return fmt.Errorf("token ID cannot be empty")
	}
	if nft.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if nft.TributeTo == "" {
		return fmt.Errorf("tribute recipient cannot be empty")
	}
	if len(nft.Quotes) < 100 {
		return fmt.Errorf("NFT must contain at least 100 cultural quotes")
	}
	if nft.RoyaltyEnabled && nft.RoyaltyPercent.IsNegative() {
		return fmt.Errorf("royalty percentage cannot be negative")
	}
	if nft.Transferable || nft.Burnable || nft.Mutable {
		return fmt.Errorf("Pradhan Sevak NFT must be non-transferable, non-burnable, and immutable")
	}
	return nil
}

// GetMetadata returns the NFT metadata in standard format
func (nft *PradhanSevakNFT) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"name":        nft.Name,
		"description": nft.Description,
		"image":       nft.ImageURI,
		"animation_url": nft.AnimationURI,
		"external_url":  nft.ExternalURI,
		"attributes": []map[string]interface{}{
			{
				"trait_type": "Tribute To",
				"value":      nft.TributeTo,
			},
			{
				"trait_type": "Genesis Block",
				"value":      nft.GenesisBlock,
			},
			{
				"trait_type": "Total Quotes",
				"value":      len(nft.Quotes),
			},
			{
				"trait_type": "Languages Supported",
				"value":      len(nft.Languages),
			},
			{
				"trait_type": "Royalty Percentage",
				"value":      nft.RoyaltyPercent.String(),
			},
			{
				"trait_type": "Beneficiary",
				"value":      nft.RoyaltyRecipient,
			},
		},
		"properties": map[string]interface{}{
			"transferable": nft.Transferable,
			"burnable":     nft.Burnable,
			"mutable":      nft.Mutable,
			"category":     "Tribute",
			"creators": []map[string]interface{}{
				{
					"address": nft.CreatorAddress,
					"share":   100,
				},
			},
		},
	}
}

// ToJSON converts the NFT to JSON format
func (nft *PradhanSevakNFT) ToJSON() ([]byte, error) {
	return json.MarshalIndent(nft, "", "  ")
}

// GetRoyaltyInfo returns royalty information for the NFT
func (nft *PradhanSevakNFT) GetRoyaltyInfo(salePrice sdk.Coin) (recipient string, amount sdk.Coin) {
	if !nft.RoyaltyEnabled {
		return "", sdk.NewCoin(salePrice.Denom, sdk.ZeroInt())
	}
	
	royaltyAmount := salePrice.Amount.ToDec().Mul(nft.RoyaltyPercent).TruncateInt()
	return nft.RoyaltyRecipient, sdk.NewCoin(salePrice.Denom, royaltyAmount)
}

// IsSpecialNFT checks if this is the special Pradhan Sevak NFT
func IsSpecialNFT(tokenID string) bool {
	return tokenID == "PRADHAN-SEVAK-001"
}