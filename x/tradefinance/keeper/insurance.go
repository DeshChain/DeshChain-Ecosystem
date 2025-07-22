package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/tradefinance/types"
)

// CreateInsurancePolicy creates a new insurance policy for an LC
func (k Keeper) CreateInsurancePolicy(ctx sdk.Context, msg *types.MsgCreateInsurancePolicy) (string, error) {
	// Get LC
	lc, found := k.GetLetterOfCredit(ctx, msg.LcId)
	if !found {
		return "", types.ErrLCNotFound
	}

	// Validate insurer
	insurerPartyID := k.GetPartyIDByAddress(ctx, msg.Insurer)
	insurer, found := k.GetTradeParty(ctx, insurerPartyID)
	if !found || insurer.PartyType != "insurance" {
		return "", types.ErrPartyNotFound
	}

	// Validate coverage doesn't exceed LC amount
	if msg.CoverageAmount.Amount.GT(lc.Amount.Amount) {
		return "", types.ErrInsufficientInsurance
	}

	// Transfer premium from LC applicant to insurer
	applicant, found := k.GetTradeParty(ctx, lc.ApplicantId)
	if !found {
		return "", types.ErrPartyNotFound
	}

	applicantAddr, err := sdk.AccAddressFromBech32(applicant.DeshAddress)
	if err != nil {
		return "", err
	}

	insurerAddr, err := sdk.AccAddressFromBech32(msg.Insurer)
	if err != nil {
		return "", err
	}

	err = k.bankKeeper.SendCoins(ctx, applicantAddr, insurerAddr, sdk.NewCoins(msg.Premium))
	if err != nil {
		return "", err
	}

	// Generate policy ID
	policyID := k.GetNextPolicyID(ctx)
	policyIDStr := fmt.Sprintf("POL%06d", policyID)

	// Create policy
	policy := types.InsurancePolicy{
		PolicyId:       policyIDStr,
		LcId:           msg.LcId,
		InsurerId:      insurerPartyID,
		PolicyType:     msg.PolicyType,
		CoverageAmount: msg.CoverageAmount,
		Premium:        msg.Premium,
		StartDate:      msg.StartDate,
		EndDate:        msg.EndDate,
		TermsIpfsHash:  msg.TermsIpfsHash,
		IsActive:       true,
		CoveredRisks:   msg.CoveredRisks,
	}

	// Save policy
	k.SetInsurancePolicy(ctx, policy)
	k.AddPolicyToLcIndex(ctx, msg.LcId, policyIDStr)
	k.SetNextPolicyID(ctx, policyID+1)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeInsuranceCreated,
			sdk.NewAttribute(types.AttributeKeyPolicyId, policyIDStr),
			sdk.NewAttribute(types.AttributeKeyLcId, msg.LcId),
			sdk.NewAttribute(types.AttributeKeyInsurer, insurerPartyID),
			sdk.NewAttribute(types.AttributeKeyCoverageAmount, msg.CoverageAmount.String()),
			sdk.NewAttribute(types.AttributeKeyPremium, msg.Premium.String()),
		),
	)

	return policyIDStr, nil
}

// GetInsurancePolicy returns an insurance policy by ID
func (k Keeper) GetInsurancePolicy(ctx sdk.Context, policyID string) (types.InsurancePolicy, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InsurancePolicyPrefix)
	
	bz := store.Get([]byte(policyID))
	if bz == nil {
		return types.InsurancePolicy{}, false
	}

	var policy types.InsurancePolicy
	k.cdc.MustUnmarshal(bz, &policy)
	return policy, true
}

// SetInsurancePolicy saves an insurance policy
func (k Keeper) SetInsurancePolicy(ctx sdk.Context, policy types.InsurancePolicy) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InsurancePolicyPrefix)
	bz := k.cdc.MustMarshal(&policy)
	store.Set([]byte(policy.PolicyId), bz)
}

// GetAllInsurancePolicies returns all insurance policies
func (k Keeper) GetAllInsurancePolicies(ctx sdk.Context) []types.InsurancePolicy {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.InsurancePolicyPrefix)
	
	var policies []types.InsurancePolicy
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var policy types.InsurancePolicy
		k.cdc.MustUnmarshal(iterator.Value(), &policy)
		policies = append(policies, policy)
	}
	
	return policies
}

// GetInsurancePoliciesByLc returns all policies for an LC
func (k Keeper) GetInsurancePoliciesByLc(ctx sdk.Context, lcID string) []types.InsurancePolicy {
	policyIDs := k.GetPolicyIDsByLc(ctx, lcID)
	
	var policies []types.InsurancePolicy
	for _, policyID := range policyIDs {
		policy, found := k.GetInsurancePolicy(ctx, policyID)
		if found {
			policies = append(policies, policy)
		}
	}
	
	return policies
}

// AddPolicyToLcIndex adds a policy to LC's index
func (k Keeper) AddPolicyToLcIndex(ctx sdk.Context, lcID, policyID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PolicyByLcPrefix)
	key := append([]byte(lcID), []byte(policyID)...)
	store.Set(key, []byte{1})
}

// GetPolicyIDsByLc returns policy IDs for an LC
func (k Keeper) GetPolicyIDsByLc(ctx sdk.Context, lcID string) []string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PolicyByLcPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte(lcID))
	defer iterator.Close()
	
	var policyIDs []string
	for ; iterator.Valid(); iterator.Next() {
		// Extract policy ID from key
		key := iterator.Key()
		policyID := string(key[len(lcID):])
		policyIDs = append(policyIDs, policyID)
	}
	
	return policyIDs
}

// GetNextPolicyID returns the next policy ID
func (k Keeper) GetNextPolicyID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextPolicyIDKey)
	
	if bz == nil {
		return 1
	}
	
	return sdk.BigEndianToUint64(bz)
}

// SetNextPolicyID sets the next policy ID
func (k Keeper) SetNextPolicyID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.NextPolicyIDKey, sdk.Uint64ToBigEndian(id))
}