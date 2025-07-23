package types

import (
	"cosmossdk.io/errors"
)

// Donation module errors
var (
	ErrInvalidType         = errors.Register(ModuleName, 1, "invalid type")
	ErrModuleDisabled      = errors.Register(ModuleName, 2, "module disabled")
	ErrNoActiveNGOs        = errors.Register(ModuleName, 3, "no active NGOs found")
	ErrNGONotFound         = errors.Register(ModuleName, 4, "NGO not found")
	ErrNGONotActive        = errors.Register(ModuleName, 5, "NGO is not active")
	ErrNGONotVerified      = errors.Register(ModuleName, 6, "NGO is not verified")
	ErrInvalidNGOAddress   = errors.Register(ModuleName, 7, "invalid NGO address")
	ErrInvalidDonation     = errors.Register(ModuleName, 8, "invalid donation amount")
	ErrDistributionFailed  = errors.Register(ModuleName, 9, "donation distribution failed")
)