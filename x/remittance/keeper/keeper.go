package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/DeshChain/DeshChain-Ecosystem/x/remittance/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec
	paramstore   paramtypes.Subspace
	logger       log.Logger

	// External keepers
	accountKeeper   types.AccountKeeper
	bankKeeper      types.BankKeeper
	ibcKeeper       types.IBCKeeper
	channelKeeper   types.ChannelKeeper
	portKeeper      types.PortKeeper
	scopedKeeper    types.ScopedKeeper

	// Authority is the address capable of executing governance proposals
	authority string
}

// NewKeeper creates a new remittance keeper
func NewKeeper(
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	logger log.Logger,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	ibcKeeper types.IBCKeeper,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	scopedKeeper types.ScopedKeeper,
	authority string,
) Keeper {
	// Ensure that authority is a valid AccAddress
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	// Set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeService:    storeService,
		cdc:             cdc,
		paramstore:      ps,
		logger:          logger,
		accountKeeper:   accountKeeper,
		bankKeeper:      bankKeeper,
		ibcKeeper:       ibcKeeper,
		channelKeeper:   channelKeeper,
		portKeeper:      portKeeper,
		scopedKeeper:    scopedKeeper,
		authority:       authority,
	}
}

// GetAuthority returns the module's authority address
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName), "height", sdkCtx.BlockHeight())
}

// GetStore returns the module store
func (k Keeper) GetStore(ctx context.Context) store.KVStore {
	return k.storeService.OpenKVStore(ctx)
}

// SetParams sets the module parameters
func (k Keeper) SetParams(ctx context.Context, params types.RemittanceParams) error {
	store := k.GetStore(ctx)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetParams returns the module parameters
func (k Keeper) GetParams(ctx context.Context) (params types.RemittanceParams, err error) {
	store := k.GetStore(ctx)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params, fmt.Errorf("parameters not found")
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params, nil
}

// GetCounters gets the module counters
func (k Keeper) GetCounters(ctx context.Context) types.Counters {
	store := k.GetStore(ctx)
	bz := store.Get(types.CountersKey)
	if bz == nil {
		return types.Counters{
			NextTransferId:    1,
			NextPoolId:        1,
			NextCorridorId:    1,
			NextPartnerId:     1,
			NextSettlementId:  1,
			TotalTransfers:    0,
			TotalVolumeUsd:    0,
		}
	}
	var counters types.Counters
	k.cdc.MustUnmarshal(bz, &counters)
	return counters
}

// SetCounters sets the module counters
func (k Keeper) SetCounters(ctx context.Context, counters types.Counters) {
	store := k.GetStore(ctx)
	bz := k.cdc.MustMarshal(&counters)
	store.Set(types.CountersKey, bz)
}

// GetNextTransferID returns the next available transfer ID
func (k Keeper) GetNextTransferID(ctx context.Context) string {
	counters := k.GetCounters(ctx)
	transferID := fmt.Sprintf("RMT-%d", counters.NextTransferId)
	counters.NextTransferId++
	k.SetCounters(ctx, counters)
	return transferID
}

// GetNextPoolID returns the next available pool ID
func (k Keeper) GetNextPoolID(ctx context.Context) string {
	counters := k.GetCounters(ctx)
	poolID := fmt.Sprintf("POOL-%d", counters.NextPoolId)
	counters.NextPoolId++
	k.SetCounters(ctx, counters)
	return poolID
}

// GetNextCorridorID returns the next available corridor ID
func (k Keeper) GetNextCorridorID(ctx context.Context) string {
	counters := k.GetCounters(ctx)
	corridorID := fmt.Sprintf("COR-%d", counters.NextCorridorId)
	counters.NextCorridorId++
	k.SetCounters(ctx, counters)
	return corridorID
}

// GetNextPartnerID returns the next available partner ID
func (k Keeper) GetNextPartnerID(ctx context.Context) string {
	counters := k.GetCounters(ctx)
	partnerID := fmt.Sprintf("PTR-%d", counters.NextPartnerId)
	counters.NextPartnerId++
	k.SetCounters(ctx, counters)
	return partnerID
}

// GetNextSettlementID returns the next available settlement ID
func (k Keeper) GetNextSettlementID(ctx context.Context) string {
	counters := k.GetCounters(ctx)
	settlementID := fmt.Sprintf("STL-%d", counters.NextSettlementId)
	counters.NextSettlementId++
	k.SetCounters(ctx, counters)
	return settlementID
}

// ========================= Remittance Transfer Methods =========================

// SetRemittanceTransfer stores a remittance transfer
func (k Keeper) SetRemittanceTransfer(ctx context.Context, transfer types.RemittanceTransfer) error {
	store := k.GetStore(ctx)
	key := types.RemittanceTransferKey(transfer.Id)
	bz := k.cdc.MustMarshal(&transfer)
	store.Set(key, bz)

	// Set indexes
	k.setTransferIndexes(ctx, transfer)

	return nil
}

// GetRemittanceTransfer retrieves a remittance transfer by ID
func (k Keeper) GetRemittanceTransfer(ctx context.Context, transferID string) (types.RemittanceTransfer, error) {
	store := k.GetStore(ctx)
	key := types.RemittanceTransferKey(transferID)
	bz := store.Get(key)
	if bz == nil {
		return types.RemittanceTransfer{}, types.ErrTransferNotFound
	}

	var transfer types.RemittanceTransfer
	k.cdc.MustUnmarshal(bz, &transfer)
	return transfer, nil
}

// HasRemittanceTransfer checks if a transfer exists
func (k Keeper) HasRemittanceTransfer(ctx context.Context, transferID string) bool {
	store := k.GetStore(ctx)
	key := types.RemittanceTransferKey(transferID)
	return store.Has(key)
}

// DeleteRemittanceTransfer removes a remittance transfer
func (k Keeper) DeleteRemittanceTransfer(ctx context.Context, transferID string) error {
	transfer, err := k.GetRemittanceTransfer(ctx, transferID)
	if err != nil {
		return err
	}

	store := k.GetStore(ctx)
	key := types.RemittanceTransferKey(transferID)
	store.Delete(key)

	// Remove indexes
	k.removeTransferIndexes(ctx, transfer)

	return nil
}

// setTransferIndexes sets the various indexes for a transfer
func (k Keeper) setTransferIndexes(ctx context.Context, transfer types.RemittanceTransfer) {
	store := k.GetStore(ctx)

	// Index by sender
	senderKey := types.TransferBySenderKey(transfer.SenderAddress, transfer.Id)
	store.Set(senderKey, []byte(transfer.Id))

	// Index by recipient
	recipientKey := types.TransferByRecipientKey(transfer.RecipientAddress, transfer.Id)
	store.Set(recipientKey, []byte(transfer.Id))

	// Index by status
	statusKey := types.TransferByStatusKey(transfer.Status, transfer.Id)
	store.Set(statusKey, []byte(transfer.Id))

	// Index by corridor
	if transfer.CorridorId != "" {
		corridorKey := types.TransferByCorridorKey(transfer.CorridorId, transfer.Id)
		store.Set(corridorKey, []byte(transfer.Id))
	}
}

// removeTransferIndexes removes the various indexes for a transfer
func (k Keeper) removeTransferIndexes(ctx context.Context, transfer types.RemittanceTransfer) {
	store := k.GetStore(ctx)

	// Remove sender index
	senderKey := types.TransferBySenderKey(transfer.SenderAddress, transfer.Id)
	store.Delete(senderKey)

	// Remove recipient index
	recipientKey := types.TransferByRecipientKey(transfer.RecipientAddress, transfer.Id)
	store.Delete(recipientKey)

	// Remove status index
	statusKey := types.TransferByStatusKey(transfer.Status, transfer.Id)
	store.Delete(statusKey)

	// Remove corridor index
	if transfer.CorridorId != "" {
		corridorKey := types.TransferByCorridorKey(transfer.CorridorId, transfer.Id)
		store.Delete(corridorKey)
	}
}

// GetAllRemittanceTransfers returns all remittance transfers
func (k Keeper) GetAllRemittanceTransfers(ctx context.Context) ([]types.RemittanceTransfer, error) {
	store := k.GetStore(ctx)
	iterator := store.Iterator(types.RemittanceTransferKeyPrefix, nil)
	defer iterator.Close()

	var transfers []types.RemittanceTransfer
	for ; iterator.Valid(); iterator.Next() {
		var transfer types.RemittanceTransfer
		k.cdc.MustUnmarshal(iterator.Value(), &transfer)
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

// GetTransfersByStatus returns transfers filtered by status
func (k Keeper) GetTransfersByStatus(ctx context.Context, status types.TransferStatus) ([]types.RemittanceTransfer, error) {
	store := k.GetStore(ctx)
	statusPrefix := append(types.TransferByStatusKeyPrefix, sdk.Uint64ToBigEndian(uint64(status))...)
	iterator := store.Iterator(statusPrefix, nil)
	defer iterator.Close()

	var transfers []types.RemittanceTransfer
	for ; iterator.Valid(); iterator.Next() {
		_, transferID := types.ParseTransferByStatusKey(iterator.Key())
		if transferID != "" {
			transfer, err := k.GetRemittanceTransfer(ctx, transferID)
			if err == nil {
				transfers = append(transfers, transfer)
			}
		}
	}

	return transfers, nil
}

// GetTransfersBySender returns transfers sent by a specific address
func (k Keeper) GetTransfersBySender(ctx context.Context, sender string) ([]types.RemittanceTransfer, error) {
	store := k.GetStore(ctx)
	senderPrefix := append(types.TransferBySenderKeyPrefix, []byte(sender)...)
	iterator := store.Iterator(senderPrefix, nil)
	defer iterator.Close()

	var transfers []types.RemittanceTransfer
	for ; iterator.Valid(); iterator.Next() {
		_, transferID := types.ParseTransferBySenderKey(iterator.Key())
		if transferID != "" {
			transfer, err := k.GetRemittanceTransfer(ctx, transferID)
			if err == nil {
				transfers = append(transfers, transfer)
			}
		}
	}

	return transfers, nil
}

// GetTransfersByRecipient returns transfers received by a specific address
func (k Keeper) GetTransfersByRecipient(ctx context.Context, recipient string) ([]types.RemittanceTransfer, error) {
	store := k.GetStore(ctx)
	recipientPrefix := append(types.TransferByRecipientKeyPrefix, []byte(recipient)...)
	iterator := store.Iterator(recipientPrefix, nil)
	defer iterator.Close()

	var transfers []types.RemittanceTransfer
	for ; iterator.Valid(); iterator.Next() {
		_, transferID := types.ParseTransferByRecipientKey(iterator.Key())
		if transferID != "" {
			transfer, err := k.GetRemittanceTransfer(ctx, transferID)
			if err == nil {
				transfers = append(transfers, transfer)
			}
		}
	}

	return transfers, nil
}

// ValidateBasic performs basic validation on a transfer
func (k Keeper) ValidateBasic(ctx context.Context, transfer *types.RemittanceTransfer) error {
	// Validate transfer ID
	if transfer.Id == "" {
		return types.ErrInvalidTransferID
	}

	// Validate addresses
	if _, err := sdk.AccAddressFromBech32(transfer.SenderAddress); err != nil {
		return types.ErrInvalidAddress
	}

	if transfer.RecipientAddress != "" {
		if _, err := sdk.AccAddressFromBech32(transfer.RecipientAddress); err != nil {
			return types.ErrInvalidAddress
		}
	}

	// Validate countries
	if transfer.SenderCountry == "" || transfer.RecipientCountry == "" {
		return types.ErrInvalidCountryCode
	}

	// Validate amounts
	if !transfer.Amount.IsValid() || !transfer.Amount.IsPositive() {
		return types.ErrInvalidAmount
	}

	// Validate currencies
	if transfer.SourceCurrency == "" || transfer.DestinationCurrency == "" {
		return types.ErrInvalidCurrency
	}

	return nil
}