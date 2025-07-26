package keeper

import (
	"fmt"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tax/types"
)

// NAMOSwapRouter handles automatic swapping of any token to NAMO for fee payment
type NAMOSwapRouter struct {
	keeper *Keeper
}

// NewNAMOSwapRouter creates a new NAMO swap router
func NewNAMOSwapRouter(k *Keeper) *NAMOSwapRouter {
	return &NAMOSwapRouter{
		keeper: k,
	}
}

// SwapForNAMOFee swaps the required amount of user's token to NAMO for fee payment
func (nsr *NAMOSwapRouter) SwapForNAMOFee(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	feeAmount sdk.Coin,
	userToken sdk.Coin,
	inclusive bool,
) (sdk.Coin, error) {
	// If user already has NAMO, use it directly
	if userToken.Denom == "namo" {
		return feeAmount, nil
	}
	
	// Get swap rate from oracle or DEX
	swapRate := nsr.getSwapRate(ctx, userToken.Denom, "namo")
	if swapRate.IsZero() {
		return sdk.Coin{}, fmt.Errorf("no swap rate available for %s to NAMO", userToken.Denom)
	}
	
	// Calculate amount of user token needed
	requiredUserToken := sdk.NewDecFromInt(feeAmount.Amount).Quo(swapRate).TruncateInt()
	
	// Apply slippage protection (max 0.5%)
	slippageProtection := sdk.NewDecWithPrec(5, 3) // 0.005
	maxRequiredToken := sdk.NewDecFromInt(requiredUserToken).Mul(sdk.OneDec().Add(slippageProtection)).TruncateInt()
	
	// Check user balance
	userBalance := nsr.keeper.bankKeeper.GetBalance(ctx, userAddr, userToken.Denom)
	if userBalance.Amount.LT(maxRequiredToken) {
		return sdk.Coin{}, fmt.Errorf("insufficient balance for fee: need %s %s, have %s", 
			maxRequiredToken, userToken.Denom, userBalance.Amount)
	}
	
	// Execute swap through DEX or oracle-based conversion
	swappedAmount, err := nsr.executeSwap(ctx, userAddr, sdk.NewCoin(userToken.Denom, requiredUserToken), "namo")
	if err != nil {
		return sdk.Coin{}, err
	}
	
	// Verify we got enough NAMO
	if swappedAmount.Amount.LT(feeAmount.Amount) {
		// Try with slippage amount
		additionalNeeded := feeAmount.Amount.Sub(swappedAmount.Amount)
		additionalSwap, err := nsr.executeSwap(ctx, userAddr, 
			sdk.NewCoin(userToken.Denom, additionalNeeded), "namo")
		if err != nil {
			return sdk.Coin{}, fmt.Errorf("failed to swap sufficient NAMO for fees: %w", err)
		}
		swappedAmount = swappedAmount.Add(additionalSwap)
	}
	
	return swappedAmount, nil
}

// getSwapRate gets the current swap rate between two tokens
func (nsr *NAMOSwapRouter) getSwapRate(ctx sdk.Context, fromDenom, toDenom string) sdk.Dec {
	// Priority 1: Check internal DEX pools
	if nsr.keeper.dexKeeper != nil {
		rate := nsr.keeper.dexKeeper.GetSwapRate(ctx, fromDenom, toDenom)
		if !rate.IsZero() {
			return rate
		}
	}
	
	// Priority 2: Check oracle prices
	if nsr.keeper.oracleKeeper != nil {
		fromPrice := nsr.keeper.oracleKeeper.GetPrice(ctx, fromDenom)
		toPrice := nsr.keeper.oracleKeeper.GetPrice(ctx, toDenom)
		if !fromPrice.IsZero() && !toPrice.IsZero() {
			return fromPrice.Quo(toPrice)
		}
	}
	
	// Priority 3: Use hardcoded rates for common pairs (fallback)
	return nsr.getHardcodedRate(fromDenom, toDenom)
}

// executeSwap executes the actual token swap
func (nsr *NAMOSwapRouter) executeSwap(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	fromCoin sdk.Coin,
	toDenom string,
) (sdk.Coin, error) {
	// Try DEX swap first
	if nsr.keeper.dexKeeper != nil {
		swappedCoin, err := nsr.keeper.dexKeeper.Swap(ctx, userAddr, fromCoin, toDenom)
		if err == nil {
			return swappedCoin, nil
		}
	}
	
	// Fallback to oracle-based conversion
	rate := nsr.getSwapRate(ctx, fromCoin.Denom, toDenom)
	if rate.IsZero() {
		return sdk.Coin{}, fmt.Errorf("no swap route available")
	}
	
	// Calculate output amount
	outputAmount := sdk.NewDecFromInt(fromCoin.Amount).Mul(rate).TruncateInt()
	outputCoin := sdk.NewCoin(toDenom, outputAmount)
	
	// Burn input tokens and mint output tokens (simplified for testnet)
	err := nsr.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, userAddr, types.ModuleName, sdk.NewCoins(fromCoin))
	if err != nil {
		return sdk.Coin{}, err
	}
	
	err = nsr.keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(outputCoin))
	if err != nil {
		return sdk.Coin{}, err
	}
	
	err = nsr.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddr, sdk.NewCoins(outputCoin))
	if err != nil {
		return sdk.Coin{}, err
	}
	
	// Emit swap event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNAMOSwap,
			sdk.NewAttribute(types.AttributeKeyUser, userAddr.String()),
			sdk.NewAttribute(types.AttributeKeyFromToken, fromCoin.String()),
			sdk.NewAttribute(types.AttributeKeyToToken, outputCoin.String()),
			sdk.NewAttribute(types.AttributeKeySwapRate, rate.String()),
		),
	)
	
	return outputCoin, nil
}

// getHardcodedRate returns hardcoded swap rates for common pairs
func (nsr *NAMOSwapRouter) getHardcodedRate(fromDenom, toDenom string) sdk.Dec {
	// NAMO base rates (1 NAMO = 1 INR assumed)
	rates := map[string]map[string]sdk.Dec{
		"namo": {
			"dinr": sdk.OneDec(),                    // 1:1
			"dusd": sdk.NewDecWithPrec(12, 3),       // 1 NAMO = 0.012 DUSD (₹83 per USD)
			"usdt": sdk.NewDecWithPrec(12, 3),       // 1 NAMO = 0.012 USDT
			"usdc": sdk.NewDecWithPrec(12, 3),       // 1 NAMO = 0.012 USDC
		},
		"dinr": {
			"namo": sdk.OneDec(),                    // 1:1
			"dusd": sdk.NewDecWithPrec(12, 3),       // Same as NAMO
			"usdt": sdk.NewDecWithPrec(12, 3),
			"usdc": sdk.NewDecWithPrec(12, 3),
		},
		"dusd": {
			"namo": sdk.NewDec(83),                  // 1 DUSD = 83 NAMO
			"dinr": sdk.NewDec(83),                  // 1 DUSD = 83 DINR
			"usdt": sdk.OneDec(),                    // 1:1 with USDT
			"usdc": sdk.OneDec(),                    // 1:1 with USDC
		},
		"usdt": {
			"namo": sdk.NewDec(83),
			"dinr": sdk.NewDec(83),
			"dusd": sdk.OneDec(),
			"usdc": sdk.OneDec(),
		},
		"usdc": {
			"namo": sdk.NewDec(83),
			"dinr": sdk.NewDec(83),
			"dusd": sdk.OneDec(),
			"usdt": sdk.OneDec(),
		},
	}
	
	if fromRates, ok := rates[fromDenom]; ok {
		if rate, ok := fromRates[toDenom]; ok {
			return rate
		}
	}
	
	return sdk.ZeroDec()
}

// CalculateFeeWithSwap calculates the fee and handles auto-swap if needed
func (nsr *NAMOSwapRouter) CalculateFeeWithSwap(
	ctx sdk.Context,
	userAddr sdk.AccAddress,
	transactionAmount sdk.Coin,
	inclusive bool,
) (feeAmount sdk.Coin, finalAmount sdk.Coin, err error) {
	// Calculate tax in NAMO
	taxCalculator := types.NewTaxCalculator(nsr.keeper.GetTaxConfig(ctx))
	taxResult, err := taxCalculator.CalculateTax(transactionAmount, "transfer")
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}
	
	feeAmount = taxResult.TaxAmount
	
	// If user is paying in NAMO, no swap needed
	if transactionAmount.Denom == "namo" {
		if inclusive {
			// Fee deducted from transaction amount
			finalAmount = transactionAmount.Sub(feeAmount)
		} else {
			// Fee added on top
			finalAmount = transactionAmount
			// User needs to have transaction amount + fee
		}
		return feeAmount, finalAmount, nil
	}
	
	// User paying in different token - need to swap for NAMO fee
	swappedFee, err := nsr.SwapForNAMOFee(ctx, userAddr, feeAmount, transactionAmount, inclusive)
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, err
	}
	
	// For non-NAMO transactions, fee is always separate
	finalAmount = transactionAmount
	
	return swappedFee, finalAmount, nil
}

// BatchSwapForFees handles batch swapping for gas optimization
func (nsr *NAMOSwapRouter) BatchSwapForFees(
	ctx sdk.Context,
	swapRequests []types.SwapRequest,
) ([]types.SwapResult, error) {
	results := make([]types.SwapResult, len(swapRequests))
	
	// Group by token pairs for efficiency
	swapGroups := make(map[string][]int) // "from-to" -> indices
	for i, req := range swapRequests {
		key := fmt.Sprintf("%s-%s", req.FromCoin.Denom, "namo")
		swapGroups[key] = append(swapGroups[key], i)
	}
	
	// Process each group
	for _, indices := range swapGroups {
		// Aggregate amounts for batch processing
		var totalAmount sdk.Int
		fromDenom := swapRequests[indices[0]].FromCoin.Denom
		
		for _, idx := range indices {
			totalAmount = totalAmount.Add(swapRequests[idx].FromCoin.Amount)
		}
		
		// Get rate once for the group
		rate := nsr.getSwapRate(ctx, fromDenom, "namo")
		
		// Process individual swaps
		for _, idx := range indices {
			req := swapRequests[idx]
			outputAmount := sdk.NewDecFromInt(req.FromCoin.Amount).Mul(rate).TruncateInt()
			
			results[idx] = types.SwapResult{
				UserAddr:     req.UserAddr,
				InputCoin:    req.FromCoin,
				OutputCoin:   sdk.NewCoin("namo", outputAmount),
				SwapRate:     rate,
				Success:      true,
			}
		}
	}
	
	return results, nil
}