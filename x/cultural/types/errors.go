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
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cultural module sentinel errors
var (
	ErrInvalidQuote           = sdkerrors.Register(ModuleName, 1100, "invalid quote")
	ErrQuoteNotFound          = sdkerrors.Register(ModuleName, 1101, "quote not found")
	ErrNoQuotesAvailable      = sdkerrors.Register(ModuleName, 1102, "no quotes available")
	ErrInvalidAuthor          = sdkerrors.Register(ModuleName, 1103, "invalid author")
	ErrInvalidCategory        = sdkerrors.Register(ModuleName, 1104, "invalid category")
	ErrInvalidLanguage        = sdkerrors.Register(ModuleName, 1105, "invalid language")
	ErrQuoteTooLong           = sdkerrors.Register(ModuleName, 1106, "quote text too long")
	ErrQuoteTooShort          = sdkerrors.Register(ModuleName, 1107, "quote text too short")
	ErrDuplicateQuote         = sdkerrors.Register(ModuleName, 1108, "duplicate quote")
	ErrMaxQuotesExceeded      = sdkerrors.Register(ModuleName, 1109, "maximum quotes per user exceeded")
	
	ErrInvalidHistoricalEvent = sdkerrors.Register(ModuleName, 1110, "invalid historical event")
	ErrEventNotFound          = sdkerrors.Register(ModuleName, 1111, "historical event not found")
	ErrInvalidEventDate       = sdkerrors.Register(ModuleName, 1112, "invalid event date")
	ErrInvalidEventCategory   = sdkerrors.Register(ModuleName, 1113, "invalid event category")
	
	ErrInvalidCulturalWisdom  = sdkerrors.Register(ModuleName, 1120, "invalid cultural wisdom")
	ErrWisdomNotFound         = sdkerrors.Register(ModuleName, 1121, "cultural wisdom not found")
	ErrInvalidTradition       = sdkerrors.Register(ModuleName, 1122, "invalid tradition")
	ErrInvalidScripture       = sdkerrors.Register(ModuleName, 1123, "invalid scripture")
	
	ErrTransactionQuoteExists = sdkerrors.Register(ModuleName, 1130, "transaction quote already exists")
	ErrInvalidTransactionHash = sdkerrors.Register(ModuleName, 1131, "invalid transaction hash")
	
	ErrNFTCreationDisabled    = sdkerrors.Register(ModuleName, 1140, "NFT creation is disabled")
	ErrInsufficientFunds      = sdkerrors.Register(ModuleName, 1141, "insufficient funds for NFT minting")
	ErrInvalidNFTMetadata     = sdkerrors.Register(ModuleName, 1142, "invalid NFT metadata")
	
	ErrStatisticsDisabled     = sdkerrors.Register(ModuleName, 1150, "statistics collection is disabled")
	ErrInvalidDateRange       = sdkerrors.Register(ModuleName, 1151, "invalid date range for statistics")
	
	ErrQuotesDisabled         = sdkerrors.Register(ModuleName, 1160, "quotes feature is disabled")
	ErrInvalidSelectionAlgorithm = sdkerrors.Register(ModuleName, 1161, "invalid quote selection algorithm")
	ErrInvalidAmountRange     = sdkerrors.Register(ModuleName, 1162, "invalid amount range")
	
	ErrUnauthorized           = sdkerrors.Register(ModuleName, 1170, "unauthorized")
	ErrInvalidAddress         = sdkerrors.Register(ModuleName, 1171, "invalid address")
)