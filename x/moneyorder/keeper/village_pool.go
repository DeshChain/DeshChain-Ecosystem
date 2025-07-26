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

package keeper

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/moneyorder/types"
)

// CreateVillagePool creates a new community-managed pool
func (k Keeper) CreateVillagePool(
	ctx sdk.Context,
	panchayatHead sdk.AccAddress,
	villageName string,
	postalCode string,
	stateCode string,
	districtCode string,
	initialLiquidity sdk.Coins,
	localValidators []sdk.ValAddress,
) (uint64, error) {
	params := k.GetParams(ctx)
	if !params.EnableVillagePools {
		return 0, types.ErrVillagePoolInactive
	}
	
	// Check if village pool already exists for this postal code
	if k.VillagePoolExists(ctx, postalCode) {
		return 0, fmt.Errorf("village pool already exists for postal code %s", postalCode)
	}
	
	// Validate initial liquidity meets minimum
	namoAmount := initialLiquidity.AmountOf("unamo")
	if namoAmount.LT(types.MinVillagePoolLiquidity) {
		return 0, fmt.Errorf("initial liquidity below minimum requirement")
	}
	
	// Generate new pool ID
	poolId := k.GetNextPoolId(ctx)
	
	// Create the village pool
	pool := types.NewVillagePool(
		poolId,
		villageName,
		postalCode,
		stateCode,
		panchayatHead,
	)
	
	// Set additional fields
	pool.DistrictCode = districtCode
	
	// Add local validators
	for i, val := range localValidators {
		if i >= 5 { // Maximum 5 local validators
			break
		}
		pool.LocalValidators = append(pool.LocalValidators, types.LocalValidator{
			ValidatorAddress: val,
			LocalName:        fmt.Sprintf("Validator-%d", i+1),
			Role:             "validator",
			TrustLevel:       5, // Start with medium trust
			JoinedAt:         ctx.BlockTime(),
		})
	}
	
	// Validate the pool
	if err := pool.ValidateVillagePool(); err != nil {
		return 0, err
	}
	
	// Transfer initial liquidity from panchayat head
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, panchayatHead, types.ModuleName, initialLiquidity,
	); err != nil {
		return 0, err
	}
	
	// Add liquidity to pool
	if err := pool.AddLiquidity(initialLiquidity); err != nil {
		return 0, err
	}
	
	// Store the pool
	k.SetVillagePool(ctx, pool)
	
	// Mark panchayat head as member
	member := &types.VillagePoolMember{
		MemberAddress:  panchayatHead,
		LocalName:      "Panchayat Head",
		MembershipType: "founder",
		JoinedAt:       ctx.BlockTime(),
		Contribution:   initialLiquidity,
	}
	k.SetVillagePoolMember(ctx, poolId, member)
	pool.TotalMembers = 1
	
	// Update pool
	k.SetVillagePool(ctx, pool)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVillagePoolCreated,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute(types.AttributeKeyVillageName, villageName),
			sdk.NewAttribute(types.AttributeKeyPostalCode, postalCode),
			sdk.NewAttribute(types.AttributeKeyPanchayatHead, panchayatHead.String()),
			sdk.NewAttribute(types.AttributeKeyLiquidity, initialLiquidity.String()),
		),
	)
	
	// Call hooks
	if k.hooks != nil {
		k.hooks.AfterVillagePoolCreated(ctx, poolId, villageName)
	}
	
	return poolId, nil
}

// SetVillagePool stores a village pool
func (k Keeper) SetVillagePool(ctx sdk.Context, pool *types.VillagePool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVillagePoolKey(pool.PostalCode)
	bz := k.cdc.MustMarshal(pool)
	store.Set(key, bz)
}

// GetVillagePool retrieves a village pool by postal code
func (k Keeper) GetVillagePool(ctx sdk.Context, postalCode string) (*types.VillagePool, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVillagePoolKey(postalCode)
	bz := store.Get(key)
	
	if bz == nil {
		return nil, false
	}
	
	var pool types.VillagePool
	k.cdc.MustUnmarshal(bz, &pool)
	return &pool, true
}

// GetVillagePoolById retrieves a village pool by ID
func (k Keeper) GetVillagePoolById(ctx sdk.Context, poolId uint64) (*types.VillagePool, bool) {
	// Iterate through all village pools to find by ID
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixVillagePool)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var pool types.VillagePool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)
		if pool.PoolId == poolId {
			return &pool, true
		}
	}
	
	return nil, false
}

// VillagePoolExists checks if a village pool exists for a postal code
func (k Keeper) VillagePoolExists(ctx sdk.Context, postalCode string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVillagePoolKey(postalCode)
	return store.Has(key)
}

// JoinVillagePool allows a user to join a village pool
func (k Keeper) JoinVillagePool(
	ctx sdk.Context,
	member sdk.AccAddress,
	poolId uint64,
	initialDeposit sdk.Coins,
	localName string,
	mobileNumber string,
) error {
	// Get the pool
	pool, found := k.GetVillagePoolById(ctx, poolId)
	if !found {
		return types.ErrVillagePoolNotFound
	}
	
	// Check if pool is active
	if !pool.Active {
		return types.ErrVillagePoolInactive
	}
	
	// Check if already a member
	if k.IsVillagePoolMemberById(ctx, poolId, member) {
		return fmt.Errorf("already a member of this village pool")
	}
	
	// Check member limit (max 1000 members per village)
	if pool.TotalMembers >= 1000 {
		return types.ErrMembershipLimitReached
	}
	
	// Transfer initial deposit
	if !initialDeposit.IsZero() {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(
			ctx, member, types.ModuleName, initialDeposit,
		); err != nil {
			return err
		}
		
		// Add to pool liquidity
		if err := pool.AddLiquidity(initialDeposit); err != nil {
			return err
		}
	}
	
	// Create member record
	memberRecord := &types.VillagePoolMember{
		MemberAddress:  member,
		LocalName:      localName,
		MobileNumber:   mobileNumber,
		MembershipType: "regular",
		JoinedAt:       ctx.BlockTime(),
		Contribution:   initialDeposit,
		TotalTrades:    0,
		TotalVolume:    sdk.ZeroInt(),
	}
	
	// Store member
	k.SetVillagePoolMember(ctx, poolId, memberRecord)
	
	// Update pool statistics
	pool.TotalMembers++
	pool.LastActivityDate = ctx.BlockTime()
	
	// Give trust score bonus for new member
	pool.UpdateTrustScore(1)
	
	// Save updated pool
	k.SetVillagePool(ctx, pool)
	
	// Mark user as village pool member (for discount calculation)
	k.MarkAsVillagePoolMember(ctx, member)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVillageMemberAdded,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute(types.AttributeKeyVillageName, pool.VillageName),
			sdk.NewAttribute("member", member.String()),
			sdk.NewAttribute("local_name", localName),
		),
	)
	
	// Call hooks
	if k.hooks != nil {
		k.hooks.AfterVillageMemberJoined(ctx, poolId, member)
	}
	
	return nil
}

// SetVillagePoolMember stores a village pool member
func (k Keeper) SetVillagePoolMember(ctx sdk.Context, poolId uint64, member *types.VillagePoolMember) {
	store := ctx.KVStore(k.storeKey)
	key := append(append([]byte("village_member:"), sdk.Uint64ToBigEndian(poolId)...), member.MemberAddress.Bytes()...)
	bz := k.cdc.MustMarshal(member)
	store.Set(key, bz)
}

// GetVillagePoolMember retrieves a village pool member
func (k Keeper) GetVillagePoolMember(ctx sdk.Context, poolId uint64, member sdk.AccAddress) (*types.VillagePoolMember, bool) {
	store := ctx.KVStore(k.storeKey)
	key := append(append([]byte("village_member:"), sdk.Uint64ToBigEndian(poolId)...), member.Bytes()...)
	bz := store.Get(key)
	
	if bz == nil {
		return nil, false
	}
	
	var memberRecord types.VillagePoolMember
	k.cdc.MustUnmarshal(bz, &memberRecord)
	return &memberRecord, true
}

// IsVillagePoolMemberById checks if user is a member of a specific village pool
func (k Keeper) IsVillagePoolMemberById(ctx sdk.Context, poolId uint64, member sdk.AccAddress) bool {
	_, found := k.GetVillagePoolMember(ctx, poolId, member)
	return found
}

// MarkAsVillagePoolMember marks a user as village pool member for discounts
func (k Keeper) MarkAsVillagePoolMember(ctx sdk.Context, member sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte("village_member:"), member.Bytes()...)
	store.Set(key, []byte{1})
}

// AddVillagePoolLiquidity adds liquidity to a village pool
func (k Keeper) AddVillagePoolLiquidity(
	ctx sdk.Context,
	poolId uint64,
	provider sdk.AccAddress,
	amount sdk.Coins,
) error {
	// Get the pool
	pool, found := k.GetVillagePoolById(ctx, poolId)
	if !found {
		return types.ErrVillagePoolNotFound
	}
	
	// Check if provider is a member
	if !k.IsVillagePoolMemberById(ctx, poolId, provider) {
		return types.ErrNotVillageMember
	}
	
	// Transfer liquidity
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, provider, types.ModuleName, amount,
	); err != nil {
		return err
	}
	
	// Add to pool
	if err := pool.AddLiquidity(amount); err != nil {
		return err
	}
	
	// Update member contribution
	member, found := k.GetVillagePoolMember(ctx, poolId, provider)
	if found {
		member.Contribution = member.Contribution.Add(amount...)
		k.SetVillagePoolMember(ctx, poolId, member)
	}
	
	// Update pool
	pool.LastActivityDate = ctx.BlockTime()
	k.SetVillagePool(ctx, pool)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddLiquidity,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute("provider", provider.String()),
			sdk.NewAttribute(types.AttributeKeyLiquidity, amount.String()),
		),
	)
	
	return nil
}

// ClaimVillagePoolRewards allows members to claim their rewards
func (k Keeper) ClaimVillagePoolRewards(
	ctx sdk.Context,
	poolId uint64,
	member sdk.AccAddress,
) (sdk.Coins, error) {
	// Get the pool
	pool, found := k.GetVillagePoolById(ctx, poolId)
	if !found {
		return nil, types.ErrVillagePoolNotFound
	}
	
	// Get member record
	memberRecord, found := k.GetVillagePoolMember(ctx, poolId, member)
	if !found {
		return nil, types.ErrNotVillageMember
	}
	
	// Check if member has pending rewards
	if memberRecord.PendingRewards.IsZero() {
		return nil, fmt.Errorf("no pending rewards")
	}
	
	// Transfer rewards from module to member
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, member, memberRecord.PendingRewards,
	); err != nil {
		return nil, err
	}
	
	// Update member record
	rewards := memberRecord.PendingRewards
	memberRecord.TotalEarnings = memberRecord.TotalEarnings.Add(rewards...)
	memberRecord.PendingRewards = sdk.NewCoins()
	k.SetVillagePoolMember(ctx, poolId, memberRecord)
	
	// Update pool
	pool.LastActivityDate = ctx.BlockTime()
	k.SetVillagePool(ctx, pool)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimRewards,
			sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute("claimer", member.String()),
			sdk.NewAttribute("rewards", rewards.String()),
		),
	)
	
	return rewards, nil
}

// DistributeVillagePoolFees distributes fees to village pool funds
func (k Keeper) DistributeVillagePoolFees(
	ctx sdk.Context,
	pool *types.VillagePool,
	fees sdk.Coin,
) error {
	// Distribute fees according to village pool configuration
	distribution := pool.DistributeFees(fees)
	
	for fundType, amount := range distribution {
		if amount.IsZero() {
			continue
		}
		
		switch fundType {
		case "community":
			pool.CommunityFund = pool.CommunityFund.Add(amount)
		case "education":
			pool.EducationFund = pool.EducationFund.Add(amount)
		case "infrastructure":
			// In production, would transfer to infrastructure account
		case "liquidity_providers":
			// Distribute to members based on contribution
			k.distributeToVillageMembers(ctx, pool.PoolId, amount)
		}
	}
	
	return nil
}

// distributeToVillageMembers distributes rewards to village pool members
func (k Keeper) distributeToVillageMembers(
	ctx sdk.Context,
	poolId uint64,
	amount sdk.Coin,
) {
	// Get all members and calculate total contribution
	totalContribution := sdk.ZeroInt()
	members := k.GetAllVillagePoolMembers(ctx, poolId)
	
	for _, member := range members {
		contribAmount := member.Contribution.AmountOf(amount.Denom)
		totalContribution = totalContribution.Add(contribAmount)
	}
	
	if totalContribution.IsZero() {
		return
	}
	
	// Distribute proportionally
	for _, member := range members {
		contribAmount := member.Contribution.AmountOf(amount.Denom)
		if contribAmount.IsZero() {
			continue
		}
		
		// Calculate member's share
		share := contribAmount.Mul(amount.Amount).Quo(totalContribution)
		if share.IsZero() {
			continue
		}
		
		// Add to pending rewards
		member.PendingRewards = member.PendingRewards.Add(sdk.NewCoin(amount.Denom, share))
		k.SetVillagePoolMember(ctx, poolId, member)
	}
}

// GetAllVillagePoolMembers returns all members of a village pool
func (k Keeper) GetAllVillagePoolMembers(ctx sdk.Context, poolId uint64) []*types.VillagePoolMember {
	store := ctx.KVStore(k.storeKey)
	prefix := append([]byte("village_member:"), sdk.Uint64ToBigEndian(poolId)...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	
	var members []*types.VillagePoolMember
	for ; iterator.Valid(); iterator.Next() {
		var member types.VillagePoolMember
		k.cdc.MustUnmarshal(iterator.Value(), &member)
		members = append(members, &member)
	}
	
	return members
}

// GetAllVillagePools returns all village pools
func (k Keeper) GetAllVillagePools(ctx sdk.Context) []*types.VillagePool {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixVillagePool)
	defer iterator.Close()
	
	var pools []*types.VillagePool
	for ; iterator.Valid(); iterator.Next() {
		var pool types.VillagePool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)
		pools = append(pools, &pool)
	}
	
	return pools
}

// UpdateVillagePoolVerification updates government verification status
func (k Keeper) UpdateVillagePoolVerification(
	ctx sdk.Context,
	postalCode string,
	verified bool,
) error {
	pool, found := k.GetVillagePool(ctx, postalCode)
	if !found {
		return types.ErrVillagePoolNotFound
	}
	
	pool.Verified = verified
	if verified {
		// Award achievement for verification
		achievement := types.VillageAchievement{
			AchievementId: fmt.Sprintf("VERIFIED-%d", pool.PoolId),
			Title:         "Government Verified",
			Description:   "Village pool has been verified by government authorities",
			Category:      "community",
			AchievedAt:    ctx.BlockTime(),
			RewardAmount:  sdk.NewCoins(sdk.NewCoin("unamo", sdk.NewInt(100_000_000))), // 100 NAMO reward
		}
		pool.AddAchievement(achievement)
		
		// Emit event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeVillageVerified,
				sdk.NewAttribute(types.AttributeKeyPoolId, fmt.Sprintf("%d", pool.PoolId)),
				sdk.NewAttribute(types.AttributeKeyVillageName, pool.VillageName),
			),
		)
	}
	
	k.SetVillagePool(ctx, pool)
	return nil
}