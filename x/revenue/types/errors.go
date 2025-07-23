package types

import (
	"cosmossdk.io/errors"
)

// Revenue module errors
var (
	ErrInvalidType         = errors.Register(ModuleName, 1, "invalid type")
	ErrInvalidRevenue      = errors.Register(ModuleName, 2, "invalid revenue amount")
	ErrInvalidDistribution = errors.Register(ModuleName, 3, "invalid distribution")
	ErrInvalidRecipient    = errors.Register(ModuleName, 4, "invalid recipient address")
	ErrInvalidSource       = errors.Register(ModuleName, 5, "invalid revenue source")
	ErrDistributionFailed  = errors.Register(ModuleName, 6, "revenue distribution failed")
	ErrModuleDisabled      = errors.Register(ModuleName, 7, "module disabled")
	ErrInsufficientFunds   = errors.Register(ModuleName, 8, "insufficient funds for distribution")
)