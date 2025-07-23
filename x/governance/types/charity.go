package types

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event types for charity allocation
const (
	EventTypeCharityAllocationUpdate = "charity_allocation_update"
	AttributeKeyCharityPercentage    = "charity_percentage"
	AttributeKeyYear                 = "year"
	AttributeKeyTimestamp            = "timestamp"
)

// CharityAllocationInfo contains information about current charity allocation
type CharityAllocationInfo struct {
	CurrentPercentage  sdk.Dec   `json:"current_percentage"`
	YearsSinceGenesis  int32     `json:"years_since_genesis"`
	NextMilestone      string    `json:"next_milestone"`
	GenesisTime        time.Time `json:"genesis_time"`
	CurrentTime        time.Time `json:"current_time"`
}