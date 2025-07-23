package types

// DUSD module event types
const (
	EventTypeDUSDMinted          = "dusd_minted"
	EventTypeDUSDBurned          = "dusd_burned"
	EventTypeDUSDFeeCalculated   = "dusd_fee_calculated"
	EventTypeDUSDStabilityAction = "dusd_stability_action"
	EventTypeDUSDPriceUpdate     = "dusd_price_update"
	EventTypeDUSDPositionCreated = "dusd_position_created"
	EventTypeDUSDPositionClosed  = "dusd_position_closed"
	EventTypeDUSDLiquidation     = "dusd_liquidation"
	
	AttributeKeyMinter          = "minter"
	AttributeKeyAmount          = "amount"
	AttributeKeyFee             = "fee"
	AttributeKeyFeeRate         = "fee_rate"
	AttributeKeyMonthlyVolume   = "monthly_volume"
	AttributeKeyPositionID      = "position_id"
	AttributeKeyHealthFactor    = "health_factor"
	AttributeKeyCollateral      = "collateral"
	AttributeKeyPriceSource     = "price_source"
	AttributeKeyPrice           = "price"
	AttributeKeyDeviation       = "deviation"
	AttributeKeyActionType      = "action_type"
	AttributeKeyActionAmount    = "action_amount"
)