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
	"cosmossdk.io/errors"
)

// Error codes for DhanSetu module
var (
	ErrInvalidDhanPataName    = errors.Register(ModuleName, 1, "invalid DhanPata name format")
	ErrDhanPataAlreadyExists  = errors.Register(ModuleName, 2, "DhanPata name already exists")
	ErrDhanPataNotFound       = errors.Register(ModuleName, 3, "DhanPata address not found")
	ErrInvalidPincode         = errors.Register(ModuleName, 4, "invalid PIN code format")
	ErrKshetraCoinExists      = errors.Register(ModuleName, 5, "Kshetra coin already exists for this pincode")
	ErrInsufficientTrustScore = errors.Register(ModuleName, 6, "insufficient trust score for operation")
	ErrMitraLimitExceeded     = errors.Register(ModuleName, 7, "mitra daily/monthly limit exceeded")
	ErrInvalidMitraType       = errors.Register(ModuleName, 8, "invalid mitra type")
	ErrMitraNotFound          = errors.Register(ModuleName, 9, "mitra profile not found")
	ErrInvalidBridgeType      = errors.Register(ModuleName, 10, "invalid cross-module bridge type")
	ErrBridgeAlreadyExists    = errors.Register(ModuleName, 11, "bridge mapping already exists")
	ErrUnauthorizedOperation  = errors.Register(ModuleName, 12, "unauthorized operation")
	ErrInvalidPaymentMethod   = errors.Register(ModuleName, 13, "invalid payment method")
	ErrKYCRequired            = errors.Register(ModuleName, 14, "KYC verification required")
	ErrInvalidBusinessInfo    = errors.Register(ModuleName, 15, "invalid business information")
	ErrInsufficientFunds      = errors.Register(ModuleName, 16, "insufficient funds for operation")
	ErrCooldownPeriod         = errors.Register(ModuleName, 17, "operation in cooldown period")
	ErrRegionNotSupported     = errors.Register(ModuleName, 18, "region not supported by mitra")
)