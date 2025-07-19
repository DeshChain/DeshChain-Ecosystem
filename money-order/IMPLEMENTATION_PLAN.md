# Money Order DEX - Clean Implementation Plan

## 📁 Folder Structure

```
x/moneyorder/                    # Single clean module (no x/dex)
├── types/
│   ├── keys.go                  # Module constants and keys
│   ├── params.go                # Money Order parameters
│   ├── fixed_rate_pool.go       # Fixed rate exchange pools
│   ├── money_order_receipt.go   # Receipt generation system
│   ├── village_pool.go          # Community liquidity pools
│   ├── trading_pair.go          # Trading pair definitions
│   ├── msgs.go                  # Transaction messages
│   ├── errors.go                # Error definitions
│   ├── events.go                # Event definitions
│   ├── codec.go                 # Protobuf registration
│   ├── genesis.go               # Genesis state types
│   └── expected_keepers.go      # Interface definitions
├── keeper/
│   ├── keeper.go                # Main keeper logic
│   ├── params.go                # Parameter management
│   ├── fixed_rate.go            # Fixed rate exchange logic
│   ├── amm.go                   # AMM pool logic (from Osmosis)
│   ├── concentrated_liquidity.go # CL pools (from Osmosis)
│   ├── receipt.go               # Receipt generation/tracking
│   ├── village_pool.go          # Village pool management
│   ├── cultural_integration.go  # Festival bonuses, quotes
│   ├── postal_routing.go        # Postal code based routing
│   ├── kyc_integration.go       # KYC for large orders
│   ├── msg_server.go            # Message handler
│   ├── grpc_query.go            # Query server
│   ├── hooks.go                 # Module hooks
│   ├── invariants.go            # Invariant checks
│   └── migrations.go            # State migrations
├── client/
│   ├── cli/
│   │   ├── tx.go                # CLI transactions
│   │   ├── query.go             # CLI queries
│   │   └── flags.go             # CLI flags
│   └── rest/
│       ├── rest.go              # REST endpoints
│       ├── tx.go                # REST transactions
│       └── query.go             # REST queries
├── simulation/
│   ├── operations.go            # Simulation operations
│   ├── params.go                # Parameter simulation
│   └── genesis.go               # Genesis simulation
├── spec/                        # Module specifications
│   ├── 01_concepts.md           # Core concepts
│   ├── 02_state.md              # State management
│   ├── 03_messages.md           # Message types
│   ├── 04_events.md             # Event types
│   ├── 05_params.md             # Parameters
│   └── README.md                # Module overview
├── module.go                    # Module interface
├── genesis.go                   # Genesis import/export
├── handler.go                   # Message handler
├── abci.go                      # ABCI hooks
└── README.md                    # Module documentation
```

## 🔄 Migration Strategy from x/dex

### Step 1: Copy Essential Logic
```bash
# Create new module structure
mkdir -p x/moneyorder/{types,keeper,client/cli,client/rest,simulation,spec}

# Copy and adapt params.go
cp x/dex/types/params.go x/moneyorder/types/params.go
# Modify to include Money Order specific parameters

# Create new types based on existing
# - fixed_rate_pool.go (new concept)
# - money_order_receipt.go (new concept)
# - village_pool.go (new concept)
```

### Step 2: Integrate Osmosis Logic
```go
// keeper/amm.go - Adapted from Osmosis
package keeper

import (
    "github.com/deshchain/deshchain/x/moneyorder/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateAMMPool creates a new AMM pool with cultural features
func (k Keeper) CreateAMMPool(ctx sdk.Context, msg *types.MsgCreateAMMPool) error {
    // Osmosis logic + cultural adaptations
    // Add festival bonus multipliers
    // Add village pool connections
    // Add receipt generation
}
```

### Step 3: Remove Old DEX Module
```bash
# After successful migration and testing
rm -rf x/dex/
# Update app.go to remove dex module references
# Update genesis to use moneyorder module
```

## 📋 Implementation Tasks

### Phase 1: Core Infrastructure (Weeks 1-4)

#### Week 1: Foundation
- [ ] Create x/moneyorder folder structure
- [ ] Implement basic types (keys.go, params.go, errors.go)
- [ ] Define proto messages for Money Order types
- [ ] Set up keeper structure

#### Week 2: Fixed Rate Pools
- [ ] Implement fixed_rate_pool.go types
- [ ] Create keeper logic for fixed rate exchanges
- [ ] Add Money Order receipt generation
- [ ] Implement postal code routing

#### Week 3: AMM Integration
- [ ] Adapt Osmosis AMM logic
- [ ] Add concentrated liquidity support
- [ ] Implement village pool system
- [ ] Add cultural trading pairs

#### Week 4: Module Integration
- [ ] Complete module.go implementation
- [ ] Add genesis import/export
- [ ] Implement ABCI hooks
- [ ] Create migration from x/dex

### Phase 2: Features & Testing (Weeks 5-8)

#### Week 5: Advanced Features
- [ ] KYC integration for large orders
- [ ] Festival bonus system
- [ ] Bulk order support
- [ ] Agent dashboard backend

#### Week 6: Client Implementation
- [ ] CLI commands for Money Orders
- [ ] REST API endpoints
- [ ] Query implementations
- [ ] Event streaming

#### Week 7: Testing
- [ ] Unit tests for all components
- [ ] Integration tests
- [ ] Simulation tests
- [ ] Security audit preparation

#### Week 8: Documentation
- [ ] API documentation
- [ ] User guides
- [ ] Migration guide
- [ ] Deployment instructions

## 🏗️ Key Components

### 1. Fixed Rate Pools
```go
// types/fixed_rate_pool.go
type FixedRatePool struct {
    PoolId          uint64
    TokenPairDenom  string   // e.g., "namo:inr"
    FixedRate       sdk.Dec  // e.g., 75.00
    MaxOrderSize    sdk.Int  // e.g., 50,000 NAMO
    DailyLimit      sdk.Int  // e.g., 1,000,000 NAMO
    ActivePeriod    Duration // e.g., 24 hours
    RequiresKYC     bool
    PostalCodes     []string // Supported regions
}
```

### 2. Money Order Receipt
```go
// types/money_order_receipt.go
type MoneyOrderReceipt struct {
    OrderId         string
    ReferenceNumber string    // MO-2024-001234
    Sender          AccAddress
    Receiver        AccAddress
    Amount          sdk.Coin
    ExchangeRate    sdk.Dec
    Fees            sdk.Coin
    PostalCodeFrom  string
    PostalCodeTo    string
    Status          OrderStatus
    CulturalQuote   string
    QRCode          []byte
    CreatedAt       time.Time
    CompletedAt     time.Time
}
```

### 3. Village Pool
```go
// types/village_pool.go
type VillagePool struct {
    PoolId          uint64
    VillageName     string
    PostalCode      string
    StateCode       string
    Liquidity       sdk.Coins
    LocalValidators []ValAddress
    CommunityFund   sdk.Dec    // 2% of fees
    ActiveTraders   uint64
    TotalVolume     sdk.Int
    Established     time.Time
}
```

## 🎯 Clean Architecture Benefits

1. **Single Module**: All Money Order logic in x/moneyorder
2. **No Confusion**: Remove x/dex completely
3. **Clear Ownership**: Money Order specific implementations
4. **Easy Testing**: Isolated module testing
5. **Better Maintenance**: All code in one place

## 📝 Migration Checklist

- [ ] Create x/moneyorder module structure
- [ ] Migrate fee distribution from x/dex/types/params.go
- [ ] Adapt trading logic for Money Orders
- [ ] Implement new Money Order specific features
- [ ] Update app/app.go to use moneyorder module
- [ ] Update genesis configuration
- [ ] Test migration path
- [ ] Remove x/dex folder
- [ ] Update all imports
- [ ] Update documentation

## 🚀 Next Steps

1. **Approve** this clean structure approach
2. **Start** with Phase 1 implementation
3. **Test** each component thoroughly
4. **Migrate** existing DEX logic carefully
5. **Remove** old x/dex module completely

This approach ensures a clean, maintainable, and culturally-aligned Money Order DEX implementation!