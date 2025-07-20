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

package urbanpension

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/urbansuraksha/keeper"
	"github.com/deshchain/deshchain/x/urbansuraksha/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all urban pension schemes
	for _, scheme := range genState.UrbanSurakshaSchemes {
		k.SetUrbanSurakshaScheme(ctx, scheme)
	}

	// Set all education loans
	for _, loan := range genState.EducationLoans {
		k.SetUrbanEducationLoan(ctx, loan)
	}

	// Set all insurance policies
	for _, policy := range genState.InsurancePolicies {
		k.SetUrbanInsurancePolicy(ctx, policy)
	}

	// Set all unified pools
	for _, pool := range genState.UnifiedPools {
		k.SetUrbanUnifiedPool(ctx, pool)
	}

	// Set all SME loans
	for _, loan := range genState.SMELoans {
		k.SetUdyamitraSMELoan(ctx, loan)
	}

	// Set all RBF investments
	for _, rbf := range genState.RBFInvestments {
		k.SetUdyamitraRBF(ctx, rbf)
	}

	// Set all referral rewards
	for _, reward := range genState.ReferralRewards {
		k.SetReferralReward(ctx, reward)
	}

	// Set module parameters
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Export all urban pension schemes
	genesis.UrbanSurakshaSchemes = k.GetAllUrbanSurakshaSchemes(ctx)

	// Export all education loans
	genesis.EducationLoans = k.GetAllUrbanEducationLoans(ctx)

	// Export all insurance policies
	genesis.InsurancePolicies = k.GetAllUrbanInsurancePolicies(ctx)

	// Export all unified pools
	genesis.UnifiedPools = k.GetAllUrbanUnifiedPools(ctx)

	// Export all SME loans
	genesis.SMELoans = k.GetAllUdyamitraSMELoans(ctx)

	// Export all RBF investments
	genesis.RBFInvestments = k.GetAllUdyamitraRBFs(ctx)

	// Export all referral rewards
	genesis.ReferralRewards = k.GetAllReferralRewards(ctx)

	return genesis
}