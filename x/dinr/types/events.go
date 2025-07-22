package types

// DINR module event types
const (
	EventTypeMintDINR          = "mint_dinr"
	EventTypeBurnDINR          = "burn_dinr"
	EventTypeDepositCollateral = "deposit_collateral"
	EventTypeWithdrawCollateral = "withdraw_collateral"
	EventTypeLiquidate         = "liquidate"
	EventTypeUpdateParams      = "update_params"
	EventTypeStabilityUpdate   = "stability_update"
	EventTypeYieldDistribution = "yield_distribution"
	EventTypeInsuranceFundUpdate = "insurance_fund_update"

	AttributeKeyMinter             = "minter"
	AttributeKeyBurner             = "burner"
	AttributeKeyDepositor          = "depositor"
	AttributeKeyWithdrawer         = "withdrawer"
	AttributeKeyLiquidator         = "liquidator"
	AttributeKeyUser               = "user"
	AttributeKeyAuthority          = "authority"
	AttributeKeyCollateral         = "collateral"
	AttributeKeyDINRMinted         = "dinr_minted"
	AttributeKeyDINRBurned         = "dinr_burned"
	AttributeKeyCollateralReturned = "collateral_returned"
	AttributeKeyCollateralReceived = "collateral_received"
	AttributeKeyDINRCovered        = "dinr_covered"
	AttributeKeyHealthFactor       = "health_factor"
	AttributeKeyCurrentPrice       = "current_price"
	AttributeKeyTargetPrice        = "target_price"
	AttributeKeyPriceDeviation     = "price_deviation"
	AttributeKeyYieldAmount        = "yield_amount"
	AttributeKeyStrategy           = "strategy"
	AttributeKeyInsuranceBalance   = "insurance_balance"

	AttributeValueCategory = ModuleName
)