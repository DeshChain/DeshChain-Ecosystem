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

package nft

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/nft/keeper"
	"github.com/DeshChain/DeshChain-Ecosystem/x/nft/types"
)

// MintPradhanSevakNFT mints the special tribute NFT in the genesis block
func MintPradhanSevakNFT(ctx sdk.Context, k keeper.Keeper) error {
	// Check if we're in genesis block
	if ctx.BlockHeight() != 0 {
		return fmt.Errorf("Pradhan Sevak NFT can only be minted in genesis block")
	}
	
	// Check if NFT already exists
	if k.HasNFT(ctx, "tribute-collection", "PRADHAN-SEVAK-001") {
		return fmt.Errorf("Pradhan Sevak NFT already exists")
	}
	
	// Create the special NFT
	pradhanSevakNFT := types.CreatePradhanSevakNFT(ctx.ChainID(), ctx.BlockHeight())
	
	// Validate the NFT
	if err := pradhanSevakNFT.Validate(); err != nil {
		return fmt.Errorf("invalid Pradhan Sevak NFT: %w", err)
	}
	
	// Create special genesis address as owner (will be transferred to PMO later)
	genesisAddr, err := sdk.AccAddressFromBech32("deshchain1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq93t2hq")
	if err != nil {
		return fmt.Errorf("invalid genesis address: %w", err)
	}
	
	// Create the NFT collection if it doesn't exist
	if !k.HasCollection(ctx, "tribute-collection") {
		collection := types.Collection{
			Denom: types.Denom{
				Id:               "tribute-collection",
				Name:             "DeshChain Tribute Collection",
				Schema:           "",
				Creator:          genesisAddr.String(),
				Symbol:           "TRIBUTE",
				MintRestricted:   true,
				UpdateRestricted: true,
				Description:      "Special collection for tribute NFTs honoring leaders who have transformed India",
				Uri:              "https://deshchain.com/collections/tribute",
				UriHash:          "",
				Data:             "",
			},
		}
		
		if err := k.SaveCollection(ctx, collection); err != nil {
			return fmt.Errorf("failed to create tribute collection: %w", err)
		}
	}
	
	// Create the NFT
	nft := types.BaseNFT{
		Id:      pradhanSevakNFT.TokenID,
		Name:    pradhanSevakNFT.Name,
		URI:     pradhanSevakNFT.ExternalURI,
		Data:    string(mustMarshalJSON(pradhanSevakNFT)),
		Owner:   genesisAddr.String(),
		UriHash: "",
	}
	
	// Mint the NFT
	if err := k.MintNFT(ctx, "tribute-collection", nft.Id, nft.Name, nft.URI, nft.UriHash, nft.Data, genesisAddr); err != nil {
		return fmt.Errorf("failed to mint Pradhan Sevak NFT: %w", err)
	}
	
	// Set special properties to make it non-transferable
	k.SetNFTProperty(ctx, "tribute-collection", nft.Id, "transferable", "false")
	k.SetNFTProperty(ctx, "tribute-collection", nft.Id, "burnable", "false")
	k.SetNFTProperty(ctx, "tribute-collection", nft.Id, "mutable", "false")
	
	// Set royalty configuration
	k.SetNFTRoyalty(ctx, "tribute-collection", nft.Id, types.Royalty{
		Percentage: pradhanSevakNFT.RoyaltyPercent,
		Recipient:  pradhanSevakNFT.RoyaltyRecipient,
	})
	
	// Emit special event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"pradhan_sevak_nft_minted",
			sdk.NewAttribute("token_id", pradhanSevakNFT.TokenID),
			sdk.NewAttribute("tribute_to", pradhanSevakNFT.TributeTo),
			sdk.NewAttribute("genesis_block", fmt.Sprintf("%d", ctx.BlockHeight())),
			sdk.NewAttribute("chain_id", ctx.ChainID()),
			sdk.NewAttribute("royalty_recipient", pradhanSevakNFT.RoyaltyRecipient),
			sdk.NewAttribute("royalty_percent", pradhanSevakNFT.RoyaltyPercent.String()),
		),
	)
	
	// Log the historic moment
	ctx.Logger().Info("Historic moment: Pradhan Sevak NFT minted in genesis block",
		"token_id", pradhanSevakNFT.TokenID,
		"tribute_to", pradhanSevakNFT.TributeTo,
		"total_quotes", len(pradhanSevakNFT.Quotes),
		"languages", len(pradhanSevakNFT.Languages),
	)
	
	return nil
}

// TransferPradhanSevakToPMO handles the formal transfer to PMO
// This will be called after mainnet launch during the presentation ceremony
func TransferPradhanSevakToPMO(ctx sdk.Context, k keeper.Keeper, pmoAddress sdk.AccAddress, transferAuthority sdk.AccAddress) error {
	// Verify the NFT exists
	if !k.HasNFT(ctx, "tribute-collection", "PRADHAN-SEVAK-001") {
		return fmt.Errorf("Pradhan Sevak NFT not found")
	}
	
	// Get current owner
	owner, err := k.GetOwner(ctx, "tribute-collection", "PRADHAN-SEVAK-001") 
	if err != nil {
		return fmt.Errorf("failed to get NFT owner: %w", err)
	}
	
	// Special override for genesis NFT transfer (normally non-transferable)
	// This requires special governance proposal or transfer authority
	if !isTransferAuthorized(ctx, transferAuthority) {
		return fmt.Errorf("unauthorized transfer attempt")
	}
	
	// Perform the transfer
	if err := k.TransferOwner(ctx, "tribute-collection", "PRADHAN-SEVAK-001", owner, pmoAddress); err != nil {
		return fmt.Errorf("failed to transfer NFT: %w", err)
	}
	
	// Emit historic event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"pradhan_sevak_nft_presented",
			sdk.NewAttribute("token_id", "PRADHAN-SEVAK-001"),
			sdk.NewAttribute("from", owner.String()),
			sdk.NewAttribute("to_pmo", pmoAddress.String()),
			sdk.NewAttribute("ceremony_time", ctx.BlockTime().String()),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
		),
	)
	
	ctx.Logger().Info("Historic presentation: Pradhan Sevak NFT transferred to PMO India",
		"recipient", pmoAddress.String(),
		"ceremony_time", ctx.BlockTime().String(),
	)
	
	return nil
}

// Helper functions

func mustMarshalJSON(v interface{}) []byte {
	bz, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bz
}

func isTransferAuthorized(ctx sdk.Context, authority sdk.AccAddress) bool {
	// In production, this would check:
	// 1. Governance proposal approval
	// 2. Multi-sig authorization
	// 3. Special transfer window
	// For now, simplified check
	return ctx.BlockHeight() > 1000 // After mainnet stabilization
}