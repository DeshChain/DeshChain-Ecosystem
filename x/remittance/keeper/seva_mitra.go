package keeper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/core/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/deshchain/deshchain/x/remittance/types"
)

// ========================= Sewa Mitra Agent Management =========================

// RegisterSewaMitraAgent registers a new Sewa Mitra agent
func (k Keeper) RegisterSewaMitraAgent(ctx context.Context, agent types.SewaMitraAgent) error {
	// Validate agent data
	if err := k.validateSewaMitraAgent(ctx, &agent); err != nil {
		return err
	}

	// Check if agent already exists
	if k.HasSewaMitraAgent(ctx, agent.AgentId) {
		return types.ErrAgentAlreadyExists
	}

	// Set initial status and timestamps
	agent.Status = types.AGENT_STATUS_PENDING_VERIFICATION
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()
	agent.TotalTransactions = 0
	agent.TotalVolume = sdk.NewCoin("usd", sdk.ZeroInt())
	agent.TotalCommissionsEarned = sdk.NewCoin("usd", sdk.ZeroInt())
	agent.SuccessRate = sdk.ZeroDec()
	agent.AverageProcessingTime = sdk.ZeroDec()

	// Store the agent
	if err := k.SetSewaMitraAgent(ctx, agent); err != nil {
		return err
	}

	// Update global counters
	counters := k.GetCounters(ctx)
	counters.TotalSevaMitraAgents++
	k.SetCounters(ctx, counters)

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"seva_mitra_agent_registered",
			sdk.NewAttribute("agent_id", agent.AgentId),
			sdk.NewAttribute("agent_name", agent.AgentName),
			sdk.NewAttribute("country", agent.Country),
			sdk.NewAttribute("city", agent.City),
		),
	)

	return nil
}

// SetSewaMitraAgent stores a Sewa Mitra agent
func (k Keeper) SetSewaMitraAgent(ctx context.Context, agent types.SewaMitraAgent) error {
	store := k.GetStore(ctx)
	key := types.SewaMitraAgentKey(agent.AgentId)
	bz := k.cdc.MustMarshal(&agent)
	store.Set(key, bz)

	// Set indexes
	k.setSewaMitraAgentIndexes(ctx, agent)

	return nil
}

// GetSewaMitraAgent retrieves a Sewa Mitra agent by ID
func (k Keeper) GetSewaMitraAgent(ctx context.Context, agentID string) (types.SewaMitraAgent, error) {
	store := k.GetStore(ctx)
	key := types.SewaMitraAgentKey(agentID)
	bz := store.Get(key)
	if bz == nil {
		return types.SewaMitraAgent{}, types.ErrAgentNotFound
	}

	var agent types.SewaMitraAgent
	k.cdc.MustUnmarshal(bz, &agent)
	return agent, nil
}

// HasSewaMitraAgent checks if a Sewa Mitra agent exists
func (k Keeper) HasSewaMitraAgent(ctx context.Context, agentID string) bool {
	store := k.GetStore(ctx)
	key := types.SewaMitraAgentKey(agentID)
	return store.Has(key)
}

// DeleteSewaMitraAgent removes a Sewa Mitra agent
func (k Keeper) DeleteSewaMitraAgent(ctx context.Context, agentID string) error {
	agent, err := k.GetSewaMitraAgent(ctx, agentID)
	if err != nil {
		return err
	}

	store := k.GetStore(ctx)
	key := types.SewaMitraAgentKey(agentID)
	store.Delete(key)

	// Remove indexes
	k.removeSewaMitraAgentIndexes(ctx, agent)

	// Update global counters
	counters := k.GetCounters(ctx)
	if counters.TotalSevaMitraAgents > 0 {
		counters.TotalSevaMitraAgents--
	}
	k.SetCounters(ctx, counters)

	return nil
}

// setSewaMitraAgentIndexes sets the various indexes for a Sewa Mitra agent
func (k Keeper) setSewaMitraAgentIndexes(ctx context.Context, agent types.SewaMitraAgent) {
	store := k.GetStore(ctx)

	// Index by country
	countryKey := types.AgentByCountryKey(agent.Country, agent.AgentId)
	store.Set(countryKey, []byte(agent.AgentId))

	// Index by city
	cityKey := types.AgentByCityKey(agent.Country, agent.City, agent.AgentId)
	store.Set(cityKey, []byte(agent.AgentId))

	// Index by status
	statusKey := types.AgentByStatusKey(agent.Status, agent.AgentId)
	store.Set(statusKey, []byte(agent.AgentId))

	// Index by currency support
	for _, currency := range agent.SupportedCurrencies {
		currencyKey := types.AgentByCurrencyKey(currency, agent.AgentId)
		store.Set(currencyKey, []byte(agent.AgentId))
	}
}

// removeSewaMitraAgentIndexes removes the various indexes for a Sewa Mitra agent
func (k Keeper) removeSewaMitraAgentIndexes(ctx context.Context, agent types.SewaMitraAgent) {
	store := k.GetStore(ctx)

	// Remove country index
	countryKey := types.AgentByCountryKey(agent.Country, agent.AgentId)
	store.Delete(countryKey)

	// Remove city index
	cityKey := types.AgentByCityKey(agent.Country, agent.City, agent.AgentId)
	store.Delete(cityKey)

	// Remove status index
	statusKey := types.AgentByStatusKey(agent.Status, agent.AgentId)
	store.Delete(statusKey)

	// Remove currency indexes
	for _, currency := range agent.SupportedCurrencies {
		currencyKey := types.AgentByCurrencyKey(currency, agent.AgentId)
		store.Delete(currencyKey)
	}
}

// GetAllSewaMitraAgents returns all Sewa Mitra agents
func (k Keeper) GetAllSewaMitraAgents(ctx context.Context) ([]types.SewaMitraAgent, error) {
	store := k.GetStore(ctx)
	iterator := store.Iterator(types.SewaMitraAgentKeyPrefix, nil)
	defer iterator.Close()

	var agents []types.SewaMitraAgent
	for ; iterator.Valid(); iterator.Next() {
		var agent types.SewaMitraAgent
		k.cdc.MustUnmarshal(iterator.Value(), &agent)
		agents = append(agents, agent)
	}

	return agents, nil
}

// GetSewaMitraAgentsByCountry returns agents filtered by country
func (k Keeper) GetSewaMitraAgentsByCountry(ctx context.Context, country string) ([]types.SewaMitraAgent, error) {
	store := k.GetStore(ctx)
	countryPrefix := append(types.AgentByCountryKeyPrefix, []byte(country)...)
	iterator := store.Iterator(countryPrefix, nil)
	defer iterator.Close()

	var agents []types.SewaMitraAgent
	for ; iterator.Valid(); iterator.Next() {
		_, agentID := types.ParseAgentByCountryKey(iterator.Key())
		if agentID != "" {
			agent, err := k.GetSewaMitraAgent(ctx, agentID)
			if err == nil {
				agents = append(agents, agent)
			}
		}
	}

	return agents, nil
}

// GetSewaMitraAgentsByCity returns agents filtered by city
func (k Keeper) GetSewaMitraAgentsByCity(ctx context.Context, country, city string) ([]types.SewaMitraAgent, error) {
	store := k.GetStore(ctx)
	cityPrefix := append(types.AgentByCityKeyPrefix, []byte(country+"/"+city)...)
	iterator := store.Iterator(cityPrefix, nil)
	defer iterator.Close()

	var agents []types.SewaMitraAgent
	for ; iterator.Valid(); iterator.Next() {
		_, _, agentID := types.ParseAgentByCityKey(iterator.Key())
		if agentID != "" {
			agent, err := k.GetSewaMitraAgent(ctx, agentID)
			if err == nil {
				agents = append(agents, agent)
			}
		}
	}

	return agents, nil
}

// GetSewaMitraAgentsByStatus returns agents filtered by status
func (k Keeper) GetSewaMitraAgentsByStatus(ctx context.Context, status types.AgentStatus) ([]types.SewaMitraAgent, error) {
	store := k.GetStore(ctx)
	statusPrefix := append(types.AgentByStatusKeyPrefix, sdk.Uint64ToBigEndian(uint64(status))...)
	iterator := store.Iterator(statusPrefix, nil)
	defer iterator.Close()

	var agents []types.SewaMitraAgent
	for ; iterator.Valid(); iterator.Next() {
		_, agentID := types.ParseAgentByStatusKey(iterator.Key())
		if agentID != "" {
			agent, err := k.GetSewaMitraAgent(ctx, agentID)
			if err == nil {
				agents = append(agents, agent)
			}
		}
	}

	return agents, nil
}

// validateSewaMitraAgent performs basic validation on a Sewa Mitra agent
func (k Keeper) validateSewaMitraAgent(ctx context.Context, agent *types.SewaMitraAgent) error {
	// Validate agent ID
	if agent.AgentId == "" {
		return types.ErrInvalidAgentID
	}

	// Validate agent address
	if _, err := sdk.AccAddressFromBech32(agent.AgentAddress); err != nil {
		return types.ErrInvalidAddress
	}

	// Validate required fields
	if agent.AgentName == "" {
		return types.ErrInvalidAgentName
	}

	if agent.Country == "" || agent.City == "" {
		return types.ErrInvalidLocation
	}

	if agent.Phone == "" {
		return types.ErrInvalidPhone
	}

	// Validate commission rates
	if agent.BaseCommissionRate.IsNegative() || agent.BaseCommissionRate.GT(sdk.NewDec(10)) {
		return types.ErrInvalidCommissionRate
	}

	if agent.VolumeBonus.IsNegative() || agent.VolumeBonus.GT(sdk.NewDec(5)) {
		return types.ErrInvalidBonusRate
	}

	// Validate limits
	if !agent.LiquidityLimit.IsValid() || !agent.LiquidityLimit.IsPositive() {
		return types.ErrInvalidLiquidity
	}

	if !agent.DailyLimit.IsValid() || !agent.DailyLimit.IsPositive() {
		return types.ErrInvalidDailyLimit
	}

	// Validate supported currencies
	if len(agent.SupportedCurrencies) == 0 {
		return types.ErrNoSupportedCurrencies
	}

	// Validate supported methods
	if len(agent.SupportedMethods) == 0 {
		return types.ErrNoSupportedMethods
	}

	return nil
}

// ActivateSewaMitraAgent activates a pending agent after verification
func (k Keeper) ActivateSewaMitraAgent(ctx context.Context, agentID string) error {
	agent, err := k.GetSewaMitraAgent(ctx, agentID)
	if err != nil {
		return err
	}

	if agent.Status != types.AGENT_STATUS_PENDING_VERIFICATION {
		return types.ErrInvalidAgentStatus
	}

	// Remove old status index
	k.removeSewaMitraAgentIndexes(ctx, agent)

	// Update status
	agent.Status = types.AGENT_STATUS_ACTIVE
	agent.UpdatedAt = time.Now()

	// Store updated agent with new indexes
	if err := k.SetSewaMitraAgent(ctx, agent); err != nil {
		return err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"seva_mitra_agent_activated",
			sdk.NewAttribute("agent_id", agent.AgentId),
			sdk.NewAttribute("agent_name", agent.AgentName),
		),
	)

	return nil
}

// SuspendSewaMitraAgent suspends an active agent
func (k Keeper) SuspendSewaMitraAgent(ctx context.Context, agentID string, suspendedUntil time.Time, reason string) error {
	agent, err := k.GetSewaMitraAgent(ctx, agentID)
	if err != nil {
		return err
	}

	if agent.Status != types.AGENT_STATUS_ACTIVE {
		return types.ErrInvalidAgentStatus
	}

	// Remove old status index
	k.removeSewaMitraAgentIndexes(ctx, agent)

	// Update status
	agent.Status = types.AGENT_STATUS_SUSPENDED
	agent.SuspendedUntil = &suspendedUntil
	agent.UpdatedAt = time.Now()

	// Add reason to metadata
	if agent.Metadata == nil {
		agent.Metadata = make(map[string]string)
	}
	agent.Metadata["suspension_reason"] = reason

	// Store updated agent with new indexes
	if err := k.SetSewaMitraAgent(ctx, agent); err != nil {
		return err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"seva_mitra_agent_suspended",
			sdk.NewAttribute("agent_id", agent.AgentId),
			sdk.NewAttribute("reason", reason),
			sdk.NewAttribute("suspended_until", suspendedUntil.Format(time.RFC3339)),
		),
	)

	return nil
}

// ========================= Sewa Mitra Commission Management =========================

// RecordSewaMitraCommission records commission earned by an agent
func (k Keeper) RecordSewaMitraCommission(ctx context.Context, transferID, agentID string, baseCommission, volumeBonus sdk.Coin) error {
	// Generate commission ID
	counters := k.GetCounters(ctx)
	commissionID := fmt.Sprintf("COM-%d", counters.NextCommissionId)
	counters.NextCommissionId++
	k.SetCounters(ctx, counters)

	// Calculate total commission
	totalCommission := baseCommission.Add(volumeBonus)

	// Create commission record
	commission := types.SewaMitraCommission{
		CommissionId:   commissionID,
		AgentId:        agentID,
		TransferId:     transferID,
		BaseCommission: baseCommission,
		VolumeBonus:    volumeBonus,
		TotalCommission: totalCommission,
		Status:         types.COMMISSION_STATUS_EARNED,
		EarnedAt:       time.Now(),
	}

	// Store commission
	if err := k.SetSewaMitraCommission(ctx, commission); err != nil {
		return err
	}

	// Update agent statistics
	agent, err := k.GetSewaMitraAgent(ctx, agentID)
	if err != nil {
		return err
	}

	agent.TotalCommissionsEarned = agent.TotalCommissionsEarned.Add(totalCommission)
	agent.UpdatedAt = time.Now()

	if err := k.SetSewaMitraAgent(ctx, agent); err != nil {
		return err
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"seva_mitra_commission_recorded",
			sdk.NewAttribute("commission_id", commissionID),
			sdk.NewAttribute("agent_id", agentID),
			sdk.NewAttribute("transfer_id", transferID),
			sdk.NewAttribute("total_commission", totalCommission.String()),
		),
	)

	return nil
}

// SetSewaMitraCommission stores a commission record
func (k Keeper) SetSewaMitraCommission(ctx context.Context, commission types.SewaMitraCommission) error {
	store := k.GetStore(ctx)
	key := types.SewaMitraCommissionKey(commission.CommissionId)
	bz := k.cdc.MustMarshal(&commission)
	store.Set(key, bz)

	// Set indexes
	k.setSewaMitraCommissionIndexes(ctx, commission)

	return nil
}

// GetSewaMitraCommission retrieves a commission record by ID
func (k Keeper) GetSewaMitraCommission(ctx context.Context, commissionID string) (types.SewaMitraCommission, error) {
	store := k.GetStore(ctx)
	key := types.SewaMitraCommissionKey(commissionID)
	bz := store.Get(key)
	if bz == nil {
		return types.SewaMitraCommission{}, types.ErrCommissionNotFound
	}

	var commission types.SewaMitraCommission
	k.cdc.MustUnmarshal(bz, &commission)
	return commission, nil
}

// setSewaMitraCommissionIndexes sets indexes for a commission record
func (k Keeper) setSewaMitraCommissionIndexes(ctx context.Context, commission types.SewaMitraCommission) {
	store := k.GetStore(ctx)

	// Index by agent
	agentKey := types.CommissionByAgentKey(commission.AgentId, commission.CommissionId)
	store.Set(agentKey, []byte(commission.CommissionId))

	// Index by transfer
	transferKey := types.CommissionByTransferKey(commission.TransferId, commission.CommissionId)
	store.Set(transferKey, []byte(commission.CommissionId))

	// Index by status
	statusKey := types.CommissionByStatusKey(commission.Status, commission.CommissionId)
	store.Set(statusKey, []byte(commission.CommissionId))
}

// GetSewaMitraCommissionsByAgent returns all commissions for an agent
func (k Keeper) GetSewaMitraCommissionsByAgent(ctx context.Context, agentID string) ([]types.SewaMitraCommission, error) {
	store := k.GetStore(ctx)
	agentPrefix := append(types.CommissionByAgentKeyPrefix, []byte(agentID)...)
	iterator := store.Iterator(agentPrefix, nil)
	defer iterator.Close()

	var commissions []types.SewaMitraCommission
	for ; iterator.Valid(); iterator.Next() {
		_, commissionID := types.ParseCommissionByAgentKey(iterator.Key())
		if commissionID != "" {
			commission, err := k.GetSewaMitraCommission(ctx, commissionID)
			if err == nil {
				commissions = append(commissions, commission)
			}
		}
	}

	return commissions, nil
}

// ========================= Sewa Mitra Transfer Routing =========================

// FindNearestSewaMitraAgent finds the nearest active agent for a given location and currency
func (k Keeper) FindNearestSewaMitraAgent(ctx context.Context, country, city, currency string) (types.SewaMitraAgent, error) {
	// First try to find agents in the same city
	agents, err := k.GetSewaMitraAgentsByCity(ctx, country, city)
	if err != nil {
		return types.SewaMitraAgent{}, err
	}

	// Filter by active status and currency support
	var suitableAgents []types.SewaMitraAgent
	for _, agent := range agents {
		if agent.Status == types.AGENT_STATUS_ACTIVE && k.agentSupportsCurrency(agent, currency) {
			suitableAgents = append(suitableAgents, agent)
		}
	}

	if len(suitableAgents) > 0 {
		// Return the agent with highest success rate
		bestAgent := suitableAgents[0]
		for _, agent := range suitableAgents[1:] {
			if agent.SuccessRate.GT(bestAgent.SuccessRate) {
				bestAgent = agent
			}
		}
		return bestAgent, nil
	}

	// If no agents in city, try country level
	countryAgents, err := k.GetSewaMitraAgentsByCountry(ctx, country)
	if err != nil {
		return types.SewaMitraAgent{}, err
	}

	for _, agent := range countryAgents {
		if agent.Status == types.AGENT_STATUS_ACTIVE && k.agentSupportsCurrency(agent, currency) {
			suitableAgents = append(suitableAgents, agent)
		}
	}

	if len(suitableAgents) == 0 {
		return types.SewaMitraAgent{}, types.ErrNoAvailableAgent
	}

	// Return the agent with highest success rate
	bestAgent := suitableAgents[0]
	for _, agent := range suitableAgents[1:] {
		if agent.SuccessRate.GT(bestAgent.SuccessRate) {
			bestAgent = agent
		}
	}

	return bestAgent, nil
}

// agentSupportsCurrency checks if an agent supports a specific currency
func (k Keeper) agentSupportsCurrency(agent types.SewaMitraAgent, currency string) bool {
	for _, supportedCurrency := range agent.SupportedCurrencies {
		if strings.EqualFold(supportedCurrency, currency) {
			return true
		}
	}
	return false
}

// CalculateSewaMitraCommission calculates commission for a Sewa Mitra agent
func (k Keeper) CalculateSewaMitraCommission(ctx context.Context, agentID string, transferAmount sdk.Coin) (baseCommission, volumeBonus sdk.Coin, err error) {
	agent, err := k.GetSewaMitraAgent(ctx, agentID)
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}

	// Calculate base commission
	baseCommissionAmount := agent.BaseCommissionRate.MulInt(transferAmount.Amount)
	baseCommission = sdk.NewCoin(transferAmount.Denom, baseCommissionAmount.TruncateInt())

	// Apply minimum/maximum limits
	if baseCommission.IsLT(agent.MinimumCommission) {
		baseCommission = agent.MinimumCommission
	}
	if baseCommission.IsGT(agent.MaximumCommission) {
		baseCommission = agent.MaximumCommission
	}

	// Calculate volume bonus if applicable
	volumeBonus = sdk.NewCoin(transferAmount.Denom, sdk.ZeroInt())
	if agent.TotalVolume.Amount.GT(sdk.NewInt(100000)) { // $100k+ volume threshold
		volumeBonusAmount := agent.VolumeBonus.MulInt(transferAmount.Amount)
		volumeBonus = sdk.NewCoin(transferAmount.Denom, volumeBonusAmount.TruncateInt())
	}

	return baseCommission, volumeBonus, nil
}

// UpdateSewaMitraAgentStats updates agent performance statistics after a transfer
func (k Keeper) UpdateSewaMitraAgentStats(ctx context.Context, agentID string, transferAmount sdk.Coin, processingTimeMinutes int64, success bool) error {
	agent, err := k.GetSewaMitraAgent(ctx, agentID)
	if err != nil {
		return err
	}

	// Update transaction count and volume
	agent.TotalTransactions++
	agent.TotalVolume = agent.TotalVolume.Add(transferAmount)

	// Update success rate
	if success {
		successfulTransactions := agent.SuccessRate.MulInt64(agent.TotalTransactions - 1).Add(sdk.OneDec())
		agent.SuccessRate = successfulTransactions.QuoInt64(agent.TotalTransactions)
	} else {
		successfulTransactions := agent.SuccessRate.MulInt64(agent.TotalTransactions - 1)
		agent.SuccessRate = successfulTransactions.QuoInt64(agent.TotalTransactions)
	}

	// Update average processing time
	currentAverage := agent.AverageProcessingTime.MulInt64(agent.TotalTransactions - 1)
	newAverage := currentAverage.Add(sdk.NewDec(processingTimeMinutes)).QuoInt64(agent.TotalTransactions)
	agent.AverageProcessingTime = newAverage

	agent.UpdatedAt = time.Now()

	return k.SetSewaMitraAgent(ctx, agent)
}