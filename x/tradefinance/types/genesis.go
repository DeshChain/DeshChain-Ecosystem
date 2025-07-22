package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:               DefaultParams(),
		Parties:              []TradeParty{},
		LettersOfCredit:      []LetterOfCredit{},
		Documents:            []TradeDocument{},
		InsurancePolicies:    []InsurancePolicy{},
		Shipments:            []ShipmentTracking{},
		PaymentInstructions:  []PaymentInstruction{},
		NextPartyId:          1,
		NextLcId:             1,
		NextDocumentId:       1,
		NextPolicyId:         1,
		NextInstructionId:    1,
		Stats: TradeFinanceStats{
			TotalLcsIssued:         0,
			TotalTradeValue:        sdk.NewCoin("dinr", sdk.ZeroInt()),
			ActiveLcs:              0,
			CompletedTrades:        0,
			AverageProcessingHours: 4,
			TopTradeCorridor:       "",
			TotalFeesCollected:     sdk.NewCoin("dinr", sdk.ZeroInt()),
			DocumentsVerified:      0,
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	// Validate parties
	partyMap := make(map[string]bool)
	for _, party := range gs.Parties {
		if party.PartyId == "" {
			return fmt.Errorf("party ID cannot be empty")
		}
		if partyMap[party.PartyId] {
			return fmt.Errorf("duplicate party ID: %s", party.PartyId)
		}
		partyMap[party.PartyId] = true

		// Validate party fields
		if party.PartyType == "" {
			return fmt.Errorf("party type cannot be empty for party %s", party.PartyId)
		}
		if party.Name == "" {
			return fmt.Errorf("party name cannot be empty for party %s", party.PartyId)
		}
		if party.DeshAddress == "" {
			return fmt.Errorf("desh address cannot be empty for party %s", party.PartyId)
		}
	}

	// Validate LCs
	lcMap := make(map[string]bool)
	for _, lc := range gs.LettersOfCredit {
		if lc.LcId == "" {
			return fmt.Errorf("LC ID cannot be empty")
		}
		if lcMap[lc.LcId] {
			return fmt.Errorf("duplicate LC ID: %s", lc.LcId)
		}
		lcMap[lc.LcId] = true

		// Validate LC fields
		if lc.LcNumber == "" {
			return fmt.Errorf("LC number cannot be empty for LC %s", lc.LcId)
		}
		if !partyMap[lc.ApplicantId] {
			return fmt.Errorf("invalid applicant ID %s for LC %s", lc.ApplicantId, lc.LcId)
		}
		if !partyMap[lc.BeneficiaryId] {
			return fmt.Errorf("invalid beneficiary ID %s for LC %s", lc.BeneficiaryId, lc.LcId)
		}
		if !partyMap[lc.IssuingBankId] {
			return fmt.Errorf("invalid issuing bank ID %s for LC %s", lc.IssuingBankId, lc.LcId)
		}
		if lc.Amount.Amount.IsNil() || lc.Amount.Amount.IsNegative() {
			return fmt.Errorf("invalid LC amount for LC %s", lc.LcId)
		}
	}

	// Validate documents
	docMap := make(map[string]bool)
	for _, doc := range gs.Documents {
		if doc.DocumentId == "" {
			return fmt.Errorf("document ID cannot be empty")
		}
		if docMap[doc.DocumentId] {
			return fmt.Errorf("duplicate document ID: %s", doc.DocumentId)
		}
		docMap[doc.DocumentId] = true

		// Validate document refers to valid LC
		if !lcMap[doc.LcId] {
			return fmt.Errorf("invalid LC ID %s for document %s", doc.LcId, doc.DocumentId)
		}
	}

	// Validate insurance policies
	policyMap := make(map[string]bool)
	for _, policy := range gs.InsurancePolicies {
		if policy.PolicyId == "" {
			return fmt.Errorf("policy ID cannot be empty")
		}
		if policyMap[policy.PolicyId] {
			return fmt.Errorf("duplicate policy ID: %s", policy.PolicyId)
		}
		policyMap[policy.PolicyId] = true

		// Validate policy refers to valid LC
		if !lcMap[policy.LcId] {
			return fmt.Errorf("invalid LC ID %s for policy %s", policy.LcId, policy.PolicyId)
		}
	}

	// Validate payment instructions
	instructionMap := make(map[string]bool)
	for _, instruction := range gs.PaymentInstructions {
		if instruction.InstructionId == "" {
			return fmt.Errorf("instruction ID cannot be empty")
		}
		if instructionMap[instruction.InstructionId] {
			return fmt.Errorf("duplicate instruction ID: %s", instruction.InstructionId)
		}
		instructionMap[instruction.InstructionId] = true

		// Validate instruction refers to valid LC
		if !lcMap[instruction.LcId] {
			return fmt.Errorf("invalid LC ID %s for instruction %s", instruction.LcId, instruction.InstructionId)
		}
	}

	// Validate shipments
	for _, shipment := range gs.Shipments {
		if shipment.TrackingId == "" {
			return fmt.Errorf("tracking ID cannot be empty")
		}

		// Validate shipment refers to valid LC
		if !lcMap[shipment.LcId] {
			return fmt.Errorf("invalid LC ID %s for shipment %s", shipment.LcId, shipment.TrackingId)
		}
	}

	// Validate counters
	if gs.NextPartyId == 0 {
		return fmt.Errorf("next party ID must be greater than 0")
	}
	if gs.NextLcId == 0 {
		return fmt.Errorf("next LC ID must be greater than 0")
	}

	return nil
}