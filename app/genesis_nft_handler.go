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

package app

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/nft"
)

// InitializeGenesisNFTs mints special NFTs during genesis
func (app *App) InitializeGenesisNFTs(ctx sdk.Context) error {
	// Only execute in genesis block
	if ctx.BlockHeight() != 0 {
		return nil
	}
	
	app.Logger().Info("Initializing Genesis NFTs...")
	
	// Mint the Pradhan Sevak NFT
	if err := nft.MintPradhanSevakNFT(ctx, app.NFTKeeper); err != nil {
		return fmt.Errorf("failed to mint Pradhan Sevak NFT: %w", err)
	}
	
	app.Logger().Info("Successfully minted Pradhan Sevak NFT in Genesis Block",
		"block_height", ctx.BlockHeight(),
		"chain_id", ctx.ChainID(),
		"tribute_to", "Shri Narendra Modi Ji",
	)
	
	// Future: Add more genesis NFTs here if needed
	
	return nil
}

// SetupGenesisHooks configures hooks for genesis initialization
func (app *App) SetupGenesisHooks() {
	// This will be called during InitChain
	app.SetInitChainer(func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		// Run default init chain
		res := app.BaseApp.InitChain(req)
		
		// Initialize genesis NFTs
		if err := app.InitializeGenesisNFTs(ctx); err != nil {
			panic(fmt.Sprintf("failed to initialize genesis NFTs: %v", err))
		}
		
		return res
	})
}

// PMOTransferProposal represents a governance proposal to transfer Pradhan Sevak NFT to PMO
type PMOTransferProposal struct {
	Title              string `json:"title"`
	Description        string `json:"description"`
	PMOAddress         string `json:"pmo_address"`
	CeremonyDate       string `json:"ceremony_date"`
	TransferAuthority  string `json:"transfer_authority"`
}

// GetTitle returns the title of the proposal
func (p PMOTransferProposal) GetTitle() string { return p.Title }

// GetDescription returns the description of the proposal
func (p PMOTransferProposal) GetDescription() string { return p.Description }

// ProposalRoute returns the routing key of the proposal
func (p PMOTransferProposal) ProposalRoute() string { return "nft" }

// ProposalType returns the type of the proposal
func (p PMOTransferProposal) ProposalType() string { return "PMOTransfer" }

// ValidateBasic performs basic validation of the proposal
func (p PMOTransferProposal) ValidateBasic() error {
	if p.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if p.Description == "" {
		return fmt.Errorf("description cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(p.PMOAddress); err != nil {
		return fmt.Errorf("invalid PMO address: %w", err)
	}
	if _, err := sdk.AccAddressFromBech32(p.TransferAuthority); err != nil {
		return fmt.Errorf("invalid transfer authority address: %w", err)
	}
	return nil
}

// String returns the string representation of the proposal
func (p PMOTransferProposal) String() string {
	return fmt.Sprintf(`PMO Transfer Proposal:
  Title: %s
  Description: %s
  PMO Address: %s
  Ceremony Date: %s
  Transfer Authority: %s`,
		p.Title, p.Description, p.PMOAddress, p.CeremonyDate, p.TransferAuthority)
}

// HandlePMOTransferProposal handles the execution of approved PMO transfer proposal
func HandlePMOTransferProposal(ctx sdk.Context, k nft.Keeper, p PMOTransferProposal) error {
	pmoAddr, err := sdk.AccAddressFromBech32(p.PMOAddress)
	if err != nil {
		return err
	}
	
	authAddr, err := sdk.AccAddressFromBech32(p.TransferAuthority)
	if err != nil {
		return err
	}
	
	// Execute the transfer
	if err := nft.TransferPradhanSevakToPMO(ctx, k, pmoAddr, authAddr); err != nil {
		return fmt.Errorf("failed to transfer Pradhan Sevak NFT to PMO: %w", err)
	}
	
	// Log the historic event
	ctx.Logger().Info("PMO Transfer Proposal Executed Successfully",
		"proposal_title", p.Title,
		"pmo_address", p.PMOAddress,
		"ceremony_date", p.CeremonyDate,
		"block_height", ctx.BlockHeight(),
	)
	
	return nil
}