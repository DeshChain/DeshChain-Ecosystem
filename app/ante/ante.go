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

package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	taxkeeper "deshchain/x/tax/keeper"
)

// NewAnteHandler creates a new ante handler for DeshChain
func NewAnteHandler(
	ak authkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	tk taxkeeper.Keeper,
	signModeHandler sdk.SignModeHandler,
	feegrantKeeper ante.FeegrantKeeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(),
		ante.NewExtensionOptionsDecorator(nil),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewDeductFeeDecorator(ak, bk, feegrantKeeper, nil),
		ante.NewSetPubKeyDecorator(ak),
		ante.NewValidateSigCountDecorator(ak),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
		// Add tax collection decorator
		NewTaxDecorator(tk, bk),
	)
}

// TaxDecorator applies transaction tax
type TaxDecorator struct {
	tk taxkeeper.Keeper
	bk bankkeeper.Keeper
}

// NewTaxDecorator creates a new TaxDecorator
func NewTaxDecorator(tk taxkeeper.Keeper, bk bankkeeper.Keeper) TaxDecorator {
	return TaxDecorator{
		tk: tk,
		bk: bk,
	}
}

// AnteHandle implements the AnteDecorator interface
func (td TaxDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Skip tax collection during simulation
	if simulate {
		return next(ctx, tx, simulate)
	}

	// Get fee payer
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, fmt.Errorf("tx must implement FeeTx interface")
	}

	feePayer := feeTx.FeePayer()
	if feePayer == nil {
		return ctx, fmt.Errorf("fee payer cannot be nil")
	}

	// Calculate transaction value (sum of all coin transfers in messages)
	totalValue := sdk.Coins{}
	for _, msg := range tx.GetMsgs() {
		// Check if message involves coin transfer
		switch m := msg.(type) {
		case interface{ GetAmount() sdk.Coins }:
			totalValue = totalValue.Add(m.GetAmount()...)
		}
	}

	// Apply tax if there's value being transferred
	if !totalValue.IsZero() {
		if err := td.tk.CollectTax(ctx, feePayer, totalValue); err != nil {
			// Log error but don't fail the transaction
			ctx.Logger().Error("failed to collect tax", "error", err)
		}
	}

	return next(ctx, tx, simulate)
}