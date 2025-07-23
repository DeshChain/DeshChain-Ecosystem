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

// x/governance module sentinel errors
var (
	ErrNotFounder             = sdkerrors.Register(ModuleName, 2, "sender is not the founder")
	ErrInvalidPhase           = sdkerrors.Register(ModuleName, 3, "invalid governance phase for this action")
	ErrProposalNotFound       = sdkerrors.Register(ModuleName, 4, "proposal not found")
	ErrInvalidProposalStatus  = sdkerrors.Register(ModuleName, 5, "invalid proposal status")
	ErrNoConsentRequired      = sdkerrors.Register(ModuleName, 6, "proposal does not require founder consent")
	ErrInvalidAuthority       = sdkerrors.Register(ModuleName, 7, "invalid authority")
	ErrParameterNotFound      = sdkerrors.Register(ModuleName, 8, "protected parameter not found")
	ErrImmutableParameter     = sdkerrors.Register(ModuleName, 9, "parameter is immutable and cannot be changed")
	ErrProposalVetoed         = sdkerrors.Register(ModuleName, 10, "proposal has been vetoed")
	ErrSupermajorityRequired  = sdkerrors.Register(ModuleName, 11, "proposal requires supermajority")
	ErrFounderConsentRequired = sdkerrors.Register(ModuleName, 12, "proposal requires founder consent")
)