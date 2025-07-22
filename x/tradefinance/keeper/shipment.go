package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/tradefinance/types"
)

// UpdateShipment updates shipment tracking information
func (k Keeper) UpdateShipment(ctx sdk.Context, msg *types.MsgUpdateShipment) error {
	// Get LC
	lc, found := k.GetLetterOfCredit(ctx, msg.LcId)
	if !found {
		return types.ErrLCNotFound
	}

	// Validate updater is authorized (carrier or authorized party)
	updaterPartyID := k.GetPartyIDByAddress(ctx, msg.Updater)
	if updaterPartyID == "" {
		return types.ErrUnauthorized
	}

	// Get or create shipment tracking
	var tracking types.ShipmentTracking
	var exists bool
	
	if msg.TrackingId != "" {
		tracking, exists = k.GetShipmentTracking(ctx, msg.TrackingId)
		if !exists {
			return types.ErrInvalidShipmentStatus
		}
	} else {
		// Create new tracking
		tracking = types.ShipmentTracking{
			TrackingId:       lc.LcId + "_SHIP", // Simple ID for now
			LcId:             msg.LcId,
			Carrier:          updaterPartyID,
			VesselName:       "",
			ContainerNumber:  "",
			CurrentLocation:  msg.Location,
			Status:           msg.Status,
			Etd:              ctx.BlockTime(),
			Eta:              ctx.BlockTime().AddDate(0, 0, 30), // Default 30 days
			ActualDeparture:  ctx.BlockTime(),
			ActualArrival:    ctx.BlockTime(),
			Events:           []types.TrackingEvent{},
		}
	}

	// Add tracking event
	event := types.TrackingEvent{
		Timestamp:   ctx.BlockTime(),
		Location:    msg.Location,
		EventType:   msg.EventType,
		Description: msg.Description,
		ReportedBy:  updaterPartyID,
	}

	tracking.Events = append(tracking.Events, event)
	tracking.CurrentLocation = msg.Location
	tracking.Status = msg.Status

	// Save tracking
	k.SetShipmentTracking(ctx, tracking)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeShipmentUpdated,
			sdk.NewAttribute(types.AttributeKeyTrackingId, tracking.TrackingId),
			sdk.NewAttribute(types.AttributeKeyLcId, msg.LcId),
			sdk.NewAttribute(types.AttributeKeyStatus, msg.Status),
			sdk.NewAttribute(types.AttributeKeyLocation, msg.Location),
			sdk.NewAttribute(types.AttributeKeyEventType, msg.EventType),
		),
	)

	return nil
}

// GetShipmentTracking returns shipment tracking by ID
func (k Keeper) GetShipmentTracking(ctx sdk.Context, trackingID string) (types.ShipmentTracking, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ShipmentTrackingPrefix)
	
	bz := store.Get([]byte(trackingID))
	if bz == nil {
		return types.ShipmentTracking{}, false
	}

	var tracking types.ShipmentTracking
	k.cdc.MustUnmarshal(bz, &tracking)
	return tracking, true
}

// SetShipmentTracking saves shipment tracking
func (k Keeper) SetShipmentTracking(ctx sdk.Context, tracking types.ShipmentTracking) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ShipmentTrackingPrefix)
	bz := k.cdc.MustMarshal(&tracking)
	store.Set([]byte(tracking.TrackingId), bz)
}

// GetAllShipmentTrackings returns all shipment trackings
func (k Keeper) GetAllShipmentTrackings(ctx sdk.Context) []types.ShipmentTracking {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ShipmentTrackingPrefix)
	
	var trackings []types.ShipmentTracking
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	
	for ; iterator.Valid(); iterator.Next() {
		var tracking types.ShipmentTracking
		k.cdc.MustUnmarshal(iterator.Value(), &tracking)
		trackings = append(trackings, tracking)
	}
	
	return trackings
}

// GetShipmentByLc returns shipment tracking for an LC
func (k Keeper) GetShipmentByLc(ctx sdk.Context, lcID string) (types.ShipmentTracking, bool) {
	// For simplicity, using LC ID as part of tracking ID
	trackingID := lcID + "_SHIP"
	return k.GetShipmentTracking(ctx, trackingID)
}