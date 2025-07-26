package keeper

import (
    "fmt"
    "time"
    "strings"
    "encoding/csv"
    "bytes"
    
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/DeshChain/DeshChain-Ecosystem/x/validator/types"
)

// AirdropManager handles token airdrops and distribution campaigns
type AirdropManager struct {
    keeper        Keeper
    tokenManager  *TokenManager
}

// NewAirdropManager creates a new airdrop manager
func NewAirdropManager(k Keeper) *AirdropManager {
    return &AirdropManager{
        keeper:       k,
        tokenManager: NewTokenManager(k),
    }
}

// CreateAirdropCampaign creates a new airdrop campaign
func (am *AirdropManager) CreateAirdropCampaign(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    tokenID uint64,
    campaign types.AirdropCampaign,
) (uint64, error) {
    // Verify token ownership
    token, found := am.keeper.GetValidatorToken(ctx, tokenID)
    if !found {
        return 0, fmt.Errorf("token not found: %d", tokenID)
    }
    
    if token.ValidatorAddr != validatorAddr.String() {
        return 0, fmt.Errorf("validator does not own this token")
    }
    
    // Validate campaign parameters
    if err := am.validateCampaign(campaign); err != nil {
        return 0, fmt.Errorf("invalid campaign: %w", err)
    }
    
    // Check available allocation
    remainingAllocation := am.calculateRemainingAirdropAllocation(ctx, token)
    if campaign.TotalAmount.GT(remainingAllocation) {
        return 0, fmt.Errorf("insufficient allocation: requested %s, available %s",
            campaign.TotalAmount.String(), remainingAllocation.String())
    }
    
    // Create campaign
    campaignID := am.keeper.GetNextAirdropCampaignID(ctx)
    campaign.CampaignID = campaignID
    campaign.TokenID = tokenID
    campaign.ValidatorAddr = validatorAddr.String()
    campaign.CreatedAt = ctx.BlockTime()
    campaign.Status = types.AirdropStatusPending
    
    am.keeper.SetAirdropCampaign(ctx, campaign)
    am.keeper.SetNextAirdropCampaignID(ctx, campaignID+1)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "airdrop_campaign_created",
            sdk.NewAttribute("campaign_id", fmt.Sprintf("%d", campaignID)),
            sdk.NewAttribute("token_id", fmt.Sprintf("%d", tokenID)),
            sdk.NewAttribute("validator", validatorAddr.String()),
            sdk.NewAttribute("total_amount", campaign.TotalAmount.String()),
            sdk.NewAttribute("recipient_count", fmt.Sprintf("%d", len(campaign.Recipients))),
        ),
    )
    
    return campaignID, nil
}

// ExecuteAirdropCampaign executes a pending airdrop campaign
func (am *AirdropManager) ExecuteAirdropCampaign(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    campaignID uint64,
) error {
    campaign, found := am.keeper.GetAirdropCampaign(ctx, campaignID)
    if !found {
        return fmt.Errorf("campaign not found: %d", campaignID)
    }
    
    if campaign.ValidatorAddr != validatorAddr.String() {
        return fmt.Errorf("validator does not own this campaign")
    }
    
    if campaign.Status != types.AirdropStatusPending {
        return fmt.Errorf("campaign is not pending: %s", campaign.Status)
    }
    
    // Check if campaign has passed start time
    if ctx.BlockTime().Before(campaign.StartTime) {
        return fmt.Errorf("campaign has not started yet")
    }
    
    // Execute the airdrop
    successCount := 0
    failedRecipients := make([]types.AirdropRecipient, 0)
    
    for _, recipient := range campaign.Recipients {
        if err := am.executeIndividualAirdrop(ctx, campaign, recipient); err != nil {
            failedRecipients = append(failedRecipients, recipient)
            continue
        }
        successCount++
    }
    
    // Update campaign status
    campaign.Status = types.AirdropStatusCompleted
    campaign.ExecutedAt = ctx.BlockTime()
    campaign.SuccessfulDrops = uint32(successCount)
    campaign.FailedRecipients = failedRecipients
    
    am.keeper.SetAirdropCampaign(ctx, campaign)
    
    // Create execution record
    executionRecord := types.AirdropExecution{
        CampaignID:       campaignID,
        ExecutedAt:       ctx.BlockTime(),
        ExecutedBy:       validatorAddr.String(),
        SuccessfulDrops:  uint32(successCount),
        FailedDrops:      uint32(len(failedRecipients)),
        TotalDistributed: campaign.TotalAmount.Sub(am.calculateFailedAmount(failedRecipients)),
        BlockHeight:      ctx.BlockHeight(),
    }
    
    am.keeper.SetAirdropExecution(ctx, executionRecord)
    
    // Emit event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "airdrop_campaign_executed",
            sdk.NewAttribute("campaign_id", fmt.Sprintf("%d", campaignID)),
            sdk.NewAttribute("successful_drops", fmt.Sprintf("%d", successCount)),
            sdk.NewAttribute("failed_drops", fmt.Sprintf("%d", len(failedRecipients))),
            sdk.NewAttribute("total_distributed", executionRecord.TotalDistributed.String()),
        ),
    )
    
    return nil
}

// CreateBulkAirdrop creates an airdrop from CSV data
func (am *AirdropManager) CreateBulkAirdrop(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    tokenID uint64,
    csvData string,
    campaignName string,
    description string,
) (uint64, error) {
    // Parse CSV data
    recipients, err := am.parseCSVRecipients(csvData)
    if err != nil {
        return 0, fmt.Errorf("failed to parse CSV: %w", err)
    }
    
    if len(recipients) == 0 {
        return 0, fmt.Errorf("no valid recipients found in CSV")
    }
    
    if len(recipients) > 10000 {
        return 0, fmt.Errorf("too many recipients (max 10,000): %d", len(recipients))
    }
    
    // Calculate total amount
    totalAmount := sdk.ZeroInt()
    for _, recipient := range recipients {
        totalAmount = totalAmount.Add(recipient.Amount)
    }
    
    // Create campaign
    campaign := types.AirdropCampaign{
        TokenID:       tokenID,
        ValidatorAddr: validatorAddr.String(),
        CampaignName:  campaignName,
        Description:   description,
        Recipients:    recipients,
        TotalAmount:   totalAmount,
        StartTime:     ctx.BlockTime(),
        CampaignType:  types.AirdropTypeBulk,
    }
    
    return am.CreateAirdropCampaign(ctx, validatorAddr, tokenID, campaign)
}

// CreateTimedAirdrop creates an airdrop with a specific start time
func (am *AirdropManager) CreateTimedAirdrop(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    tokenID uint64,
    recipients []types.AirdropRecipient,
    startTime time.Time,
    campaignName string,
    description string,
) (uint64, error) {
    if startTime.Before(ctx.BlockTime()) {
        return 0, fmt.Errorf("start time cannot be in the past")
    }
    
    // Calculate total amount
    totalAmount := sdk.ZeroInt()
    for _, recipient := range recipients {
        totalAmount = totalAmount.Add(recipient.Amount)
    }
    
    // Create campaign
    campaign := types.AirdropCampaign{
        TokenID:       tokenID,
        ValidatorAddr: validatorAddr.String(),
        CampaignName:  campaignName,
        Description:   description,
        Recipients:    recipients,
        TotalAmount:   totalAmount,
        StartTime:     startTime,
        CampaignType:  types.AirdropTypeTimed,
    }
    
    return am.CreateAirdropCampaign(ctx, validatorAddr, tokenID, campaign)
}

// CreateVestingAirdrop creates an airdrop with vesting schedule
func (am *AirdropManager) CreateVestingAirdrop(
    ctx sdk.Context,
    validatorAddr sdk.AccAddress,
    tokenID uint64,
    recipients []types.AirdropRecipient,
    vestingSchedule types.VestingSchedule,
    campaignName string,
    description string,
) (uint64, error) {
    // Validate vesting schedule
    if err := am.validateVestingSchedule(vestingSchedule); err != nil {
        return 0, fmt.Errorf("invalid vesting schedule: %w", err)
    }
    
    // Calculate total amount
    totalAmount := sdk.ZeroInt()
    for _, recipient := range recipients {
        totalAmount = totalAmount.Add(recipient.Amount)
    }
    
    // Create campaign
    campaign := types.AirdropCampaign{
        TokenID:         tokenID,
        ValidatorAddr:   validatorAddr.String(),
        CampaignName:    campaignName,
        Description:     description,
        Recipients:      recipients,
        TotalAmount:     totalAmount,
        StartTime:       ctx.BlockTime(),
        CampaignType:    types.AirdropTypeVesting,
        VestingSchedule: &vestingSchedule,
    }
    
    return am.CreateAirdropCampaign(ctx, validatorAddr, tokenID, campaign)
}

// GetAirdropCampaigns returns all campaigns for a validator
func (am *AirdropManager) GetAirdropCampaigns(
    ctx sdk.Context,
    validatorAddr string,
) ([]types.AirdropCampaign, error) {
    return am.keeper.GetAirdropCampaignsByValidator(ctx, validatorAddr), nil
}

// GetAirdropAnalytics returns analytics for validator's airdrop campaigns
func (am *AirdropManager) GetAirdropAnalytics(
    ctx sdk.Context,
    validatorAddr string,
) (types.AirdropAnalytics, error) {
    campaigns := am.keeper.GetAirdropCampaignsByValidator(ctx, validatorAddr)
    
    analytics := types.AirdropAnalytics{
        ValidatorAddr: validatorAddr,
    }
    
    for _, campaign := range campaigns {
        analytics.TotalCampaigns++
        analytics.TotalRecipients += uint32(len(campaign.Recipients))
        analytics.TotalDistributed = analytics.TotalDistributed.Add(campaign.TotalAmount)
        
        switch campaign.Status {
        case types.AirdropStatusPending:
            analytics.PendingCampaigns++
        case types.AirdropStatusCompleted:
            analytics.CompletedCampaigns++
            analytics.SuccessfulDrops += campaign.SuccessfulDrops
        case types.AirdropStatusCancelled:
            analytics.CancelledCampaigns++
        }
    }
    
    // Calculate success rate
    if analytics.TotalRecipients > 0 {
        analytics.SuccessRate = sdk.NewDec(int64(analytics.SuccessfulDrops)).
            Quo(sdk.NewDec(int64(analytics.TotalRecipients))).
            Mul(sdk.NewDec(100))
    }
    
    return analytics, nil
}

// Helper functions

func (am *AirdropManager) validateCampaign(campaign types.AirdropCampaign) error {
    if len(campaign.CampaignName) == 0 {
        return fmt.Errorf("campaign name cannot be empty")
    }
    
    if len(campaign.Recipients) == 0 {
        return fmt.Errorf("no recipients specified")
    }
    
    if len(campaign.Recipients) > 10000 {
        return fmt.Errorf("too many recipients (max 10,000)")
    }
    
    if campaign.TotalAmount.IsZero() || campaign.TotalAmount.IsNegative() {
        return fmt.Errorf("invalid total amount: %s", campaign.TotalAmount)
    }
    
    // Validate recipients
    addressMap := make(map[string]bool)
    for i, recipient := range campaign.Recipients {
        if _, err := sdk.AccAddressFromBech32(recipient.Address); err != nil {
            return fmt.Errorf("invalid recipient address at index %d: %v", i, err)
        }
        
        if recipient.Amount.IsZero() || recipient.Amount.IsNegative() {
            return fmt.Errorf("invalid amount for recipient %d: %s", i, recipient.Amount)
        }
        
        // Check for duplicate addresses
        if addressMap[recipient.Address] {
            return fmt.Errorf("duplicate recipient address: %s", recipient.Address)
        }
        addressMap[recipient.Address] = true
    }
    
    return nil
}

func (am *AirdropManager) calculateRemainingAirdropAllocation(
    ctx sdk.Context,
    token types.ValidatorToken,
) sdk.Int {
    // Get all executed airdrops for this token
    campaigns := am.keeper.GetAirdropCampaignsByToken(ctx, token.TokenID)
    
    totalAirdropped := sdk.ZeroInt()
    for _, campaign := range campaigns {
        if campaign.Status == types.AirdropStatusCompleted {
            totalAirdropped = totalAirdropped.Add(campaign.TotalAmount)
        }
    }
    
    originalAllocation := token.TotalSupply.MulRaw(15).QuoRaw(100) // 15%
    remaining := originalAllocation.Sub(totalAirdropped)
    
    if remaining.IsNegative() {
        return sdk.ZeroInt()
    }
    
    return remaining
}

func (am *AirdropManager) executeIndividualAirdrop(
    ctx sdk.Context,
    campaign types.AirdropCampaign,
    recipient types.AirdropRecipient,
) error {
    // For vesting airdrops, create vesting account
    if campaign.CampaignType == types.AirdropTypeVesting && campaign.VestingSchedule != nil {
        return am.createVestingAirdrop(ctx, campaign, recipient)
    }
    
    // Regular airdrop - transfer tokens immediately
    token, _ := am.keeper.GetValidatorToken(ctx, campaign.TokenID)
    return am.tokenManager.transferTokens(ctx, token, recipient.Address, recipient.Amount)
}

func (am *AirdropManager) createVestingAirdrop(
    ctx sdk.Context,
    campaign types.AirdropCampaign,
    recipient types.AirdropRecipient,
) error {
    // Create vesting schedule for this recipient
    vestingAccount := types.VestingAirdrop{
        CampaignID:       campaign.CampaignID,
        RecipientAddr:    recipient.Address,
        TotalAmount:      recipient.Amount,
        VestingSchedule:  *campaign.VestingSchedule,
        CreatedAt:        ctx.BlockTime(),
        UnlockedAmount:   sdk.ZeroInt(),
        LastUnlockTime:   ctx.BlockTime(),
    }
    
    am.keeper.SetVestingAirdrop(ctx, vestingAccount)
    
    // Emit vesting created event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "vesting_airdrop_created",
            sdk.NewAttribute("campaign_id", fmt.Sprintf("%d", campaign.CampaignID)),
            sdk.NewAttribute("recipient", recipient.Address),
            sdk.NewAttribute("total_amount", recipient.Amount.String()),
            sdk.NewAttribute("vesting_duration", fmt.Sprintf("%d", campaign.VestingSchedule.DurationMonths)),
        ),
    )
    
    return nil
}

func (am *AirdropManager) calculateFailedAmount(failedRecipients []types.AirdropRecipient) sdk.Int {
    total := sdk.ZeroInt()
    for _, recipient := range failedRecipients {
        total = total.Add(recipient.Amount)
    }
    return total
}

func (am *AirdropManager) parseCSVRecipients(csvData string) ([]types.AirdropRecipient, error) {
    reader := csv.NewReader(strings.NewReader(csvData))
    records, err := reader.ReadAll()
    if err != nil {
        return nil, fmt.Errorf("failed to parse CSV: %w", err)
    }
    
    if len(records) == 0 {
        return nil, fmt.Errorf("CSV is empty")
    }
    
    // Skip header if present
    startIndex := 0
    if len(records) > 0 && (strings.ToLower(records[0][0]) == "address" || strings.ToLower(records[0][0]) == "wallet") {
        startIndex = 1
    }
    
    var recipients []types.AirdropRecipient
    for i := startIndex; i < len(records); i++ {
        record := records[i]
        if len(record) < 2 {
            continue // Skip invalid rows
        }
        
        address := strings.TrimSpace(record[0])
        amountStr := strings.TrimSpace(record[1])
        
        // Validate address
        if _, err := sdk.AccAddressFromBech32(address); err != nil {
            continue // Skip invalid addresses
        }
        
        // Parse amount
        amount, ok := sdk.NewIntFromString(amountStr)
        if !ok || amount.IsZero() {
            continue // Skip invalid amounts
        }
        
        recipients = append(recipients, types.AirdropRecipient{
            Address: address,
            Amount:  amount,
        })
    }
    
    return recipients, nil
}

func (am *AirdropManager) validateVestingSchedule(schedule types.VestingSchedule) error {
    if schedule.DurationMonths == 0 {
        return fmt.Errorf("vesting duration cannot be zero")
    }
    
    if schedule.DurationMonths > 60 { // Max 5 years
        return fmt.Errorf("vesting duration too long (max 60 months)")
    }
    
    if schedule.CliffMonths > schedule.DurationMonths {
        return fmt.Errorf("cliff period cannot exceed total duration")
    }
    
    if schedule.UnlockPercentage.IsZero() || schedule.UnlockPercentage.GT(sdk.OneDec()) {
        return fmt.Errorf("unlock percentage must be between 0 and 1")
    }
    
    return nil
}

// ProcessVestingUnlocks processes pending vesting unlocks
func (am *AirdropManager) ProcessVestingUnlocks(ctx sdk.Context) error {
    vestingAirdrops := am.keeper.GetAllVestingAirdrops(ctx)
    
    for _, vesting := range vestingAirdrops {
        if am.shouldUnlockVesting(ctx, vesting) {
            if err := am.processVestingUnlock(ctx, vesting); err != nil {
                // Log error but continue processing others
                continue
            }
        }
    }
    
    return nil
}

func (am *AirdropManager) shouldUnlockVesting(ctx sdk.Context, vesting types.VestingAirdrop) bool {
    // Check if cliff period has passed
    cliffTime := vesting.CreatedAt.Add(time.Duration(vesting.VestingSchedule.CliffMonths) * 30 * 24 * time.Hour)
    if ctx.BlockTime().Before(cliffTime) {
        return false
    }
    
    // Check if it's time for next unlock
    monthsSinceCreation := int(ctx.BlockTime().Sub(vesting.CreatedAt).Hours() / (24 * 30))
    expectedUnlocks := monthsSinceCreation * int(vesting.VestingSchedule.UnlockPercentage.MulInt64(100).TruncateInt64())
    currentUnlocks := vesting.UnlockedAmount.ToDec().Quo(vesting.TotalAmount.ToDec()).MulInt64(100).TruncateInt64()
    
    return int64(expectedUnlocks) > currentUnlocks
}

func (am *AirdropManager) processVestingUnlock(ctx sdk.Context, vesting types.VestingAirdrop) error {
    // Calculate unlock amount
    monthsSinceCreation := int(ctx.BlockTime().Sub(vesting.CreatedAt).Hours() / (24 * 30))
    totalUnlockPercentage := sdk.NewDec(int64(monthsSinceCreation)).Mul(vesting.VestingSchedule.UnlockPercentage)
    
    if totalUnlockPercentage.GT(sdk.OneDec()) {
        totalUnlockPercentage = sdk.OneDec()
    }
    
    totalUnlockAmount := vesting.TotalAmount.ToDec().Mul(totalUnlockPercentage).TruncateInt()
    unlockAmount := totalUnlockAmount.Sub(vesting.UnlockedAmount)
    
    if unlockAmount.IsZero() {
        return nil
    }
    
    // Transfer unlocked tokens
    campaign, _ := am.keeper.GetAirdropCampaign(ctx, vesting.CampaignID)
    token, _ := am.keeper.GetValidatorToken(ctx, campaign.TokenID)
    
    if err := am.tokenManager.transferTokens(ctx, token, vesting.RecipientAddr, unlockAmount); err != nil {
        return err
    }
    
    // Update vesting record
    vesting.UnlockedAmount = vesting.UnlockedAmount.Add(unlockAmount)
    vesting.LastUnlockTime = ctx.BlockTime()
    am.keeper.SetVestingAirdrop(ctx, vesting)
    
    // Emit unlock event
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(
            "vesting_tokens_unlocked",
            sdk.NewAttribute("campaign_id", fmt.Sprintf("%d", vesting.CampaignID)),
            sdk.NewAttribute("recipient", vesting.RecipientAddr),
            sdk.NewAttribute("unlocked_amount", unlockAmount.String()),
            sdk.NewAttribute("total_unlocked", vesting.UnlockedAmount.String()),
        ),
    )
    
    return nil
}