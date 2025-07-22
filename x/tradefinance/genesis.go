package tradefinance

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/tradefinance/keeper"
	"github.com/deshchain/deshchain/x/tradefinance/types"
)

// InitGenesis initializes the trade finance module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module parameters
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	// Set trade parties
	for _, party := range genState.Parties {
		k.SetTradeParty(ctx, party)
		// Create address index
		k.SetPartyIDByAddress(ctx, party.DeshAddress, party.PartyId)
	}

	// Set letters of credit
	for _, lc := range genState.LettersOfCredit {
		k.SetLetterOfCredit(ctx, lc)
		// Create party indexes
		k.AddLcToPartyIndex(ctx, lc.ApplicantId, lc.LcId)
		k.AddLcToPartyIndex(ctx, lc.BeneficiaryId, lc.LcId)
		k.AddLcToPartyIndex(ctx, lc.IssuingBankId, lc.LcId)
	}

	// Set documents
	for _, doc := range genState.Documents {
		k.SetTradeDocument(ctx, doc)
		// Create LC index
		k.AddDocumentToLcIndex(ctx, doc.LcId, doc.DocumentId)
	}

	// Set insurance policies
	for _, policy := range genState.InsurancePolicies {
		k.SetInsurancePolicy(ctx, policy)
		k.AddPolicyToLcIndex(ctx, policy.LcId, policy.PolicyId)
	}

	// Set shipments
	for _, shipment := range genState.Shipments {
		k.SetShipmentTracking(ctx, shipment)
	}

	// Set payment instructions
	for _, payment := range genState.PaymentInstructions {
		k.SetPaymentInstruction(ctx, payment)
		k.AddPaymentToLcIndex(ctx, payment.LcId, payment.InstructionId)
	}

	// Set counters
	k.SetNextPartyID(ctx, genState.NextPartyId)
	k.SetNextLcID(ctx, genState.NextLcId)
	k.SetNextDocumentID(ctx, genState.NextDocumentId)
	k.SetNextPolicyID(ctx, genState.NextPolicyId)
	k.SetNextInstructionID(ctx, genState.NextInstructionId)

	// Set statistics
	k.SetTradeFinanceStats(ctx, genState.Stats)
}

// ExportGenesis returns the trade finance module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export all data
	genesis.Parties = k.GetAllTradeParties(ctx)
	genesis.LettersOfCredit = k.GetAllLettersOfCredit(ctx)
	genesis.Documents = k.GetAllTradeDocuments(ctx)
	genesis.InsurancePolicies = k.GetAllInsurancePolicies(ctx)
	genesis.Shipments = k.GetAllShipmentTrackings(ctx)
	genesis.PaymentInstructions = k.GetAllPaymentInstructions(ctx)

	// Export counters
	genesis.NextPartyId = k.GetNextPartyID(ctx)
	genesis.NextLcId = k.GetNextLcID(ctx)
	genesis.NextDocumentId = k.GetNextDocumentID(ctx)
	genesis.NextPolicyId = k.GetNextPolicyID(ctx)
	genesis.NextInstructionId = k.GetNextInstructionID(ctx)

	// Export statistics
	genesis.Stats = k.GetTradeFinanceStats(ctx)

	return genesis
}