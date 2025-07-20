# 🔧 DhanSetu Refactoring Instructions

## Quick Start: What to Do First

### Day 1-2: Immediate Actions
```bash
# 1. Clone and setup
git clone <repo>
cd deshchain

# 2. Create DhanSetu app structure
mkdir -p apps/dhansetu
cp -r batua/mobile apps/dhansetu/base  # Use Batua as base

# 3. Quick renames for compliance
find . -type f -name "*.go" -exec sed -i 's/Pension/Suraksha/g' {} +
find . -type f -name "*.go" -exec sed -i 's/pension/suraksha/g' {} +
```

## 📦 Task 1: Port Batua Wallet to DhanSetu

### Step 1.1: Analyze Existing Batua Code
```dart
// EXISTING: batua/mobile/lib/core/wallet/hd_wallet.dart
class HDWallet {
  // ✅ KEEP THIS - Just convert to TypeScript
  static Future<HDWallet> createWallet({String? mnemonic}) async {
    // BIP39 mnemonic generation - REUSE
    // BIP32 key derivation - REUSE
    // Multi-chain support - REUSE
  }
}

// EXISTING: batua/mobile/lib/core/blockchain/deshchain_client.dart
class DeshChainClient {
  // ✅ KEEP THIS - Convert to TypeScript
  Future<Balance> getBalance(String address) async {
    // RPC calls - REUSE
    // Transaction parsing - REUSE
  }
}
```

### Step 1.2: Convert to React Native
```typescript
// NEW: apps/dhansetu/src/core/wallet/HDWallet.ts
import * as bip39 from 'bip39';
import * as bip32 from 'bip32';

export class HDWallet {
  // Direct port from Dart
  static async createWallet(mnemonic?: string): Promise<HDWallet> {
    // Copy logic from hd_wallet.dart
    const seed = mnemonic || bip39.generateMnemonic();
    // ... rest of implementation
  }
}

// NEW: apps/dhansetu/src/core/blockchain/DeshChainClient.ts
export class DeshChainClient {
  // Port from deshchain_client.dart
  async getBalance(address: string): Promise<Balance> {
    // Copy RPC logic
  }
}
```

### Step 1.3: Reuse UI Components
```typescript
// EXISTING: batua/mobile/lib/ui/widgets/cultural_gradient_text.dart
// Has beautiful gradient text with Indian flag colors

// CONVERT TO: apps/dhansetu/src/components/CulturalGradientText.tsx
import { LinearGradient } from 'expo-linear-gradient';

export const CulturalGradientText: React.FC<Props> = ({ text, style }) => {
  // Port the gradient logic from Dart
  return (
    <LinearGradient colors={['#FF9933', '#FFFFFF', '#138808']}>
      <Text style={style}>{text}</Text>
    </LinearGradient>
  );
};
```

## 📚 Task 2: Enhance Cultural Module

### Step 2.1: Use Existing Quote Data
```go
// EXISTING: x/cultural/types/cultural_data.go
var gandhiQuotes = []CulturalQuote{
    // ✅ 10,000+ quotes already here!
    {
        Text: map[string]string{
            "en": "Be the change you wish to see",
            "hi": "वह बदलाव बनिए जो आप देखना चाहते हैं",
        },
        Author: "Mahatma Gandhi",
        Category: []string{"leadership", "change"},
    },
    // ... thousands more
}

// NEW: x/cultural/keeper/quotes.go
func (k Keeper) GetQuoteForTransaction(ctx sdk.Context, amount sdk.Int) (CulturalQuote, error) {
    // Add selection logic
    category := k.getCategoryByAmount(amount)
    quotes := k.getQuotesByCategory(category)
    
    // Add randomization
    index := rand.Intn(len(quotes))
    return quotes[index], nil
}
```

### Step 2.2: Add Keeper Methods
```go
// NEW: x/cultural/keeper/keeper.go
type Keeper struct {
    storeKey sdk.StoreKey
    cdc      codec.BinaryCodec
}

// NEW: x/cultural/keeper/msg_server.go
func (k msgServer) GetDailyWisdom(ctx context.Context, req *types.QueryDailyWisdomRequest) (*types.QueryDailyWisdomResponse, error) {
    // Use existing quotes
    dayOfYear := time.Now().YearDay()
    quote := gandhiQuotes[dayOfYear % len(gandhiQuotes)]
    
    return &types.QueryDailyWisdomResponse{
        Quote: &quote,
    }, nil
}
```

## 🏦 Task 3: Rename Pension to Suraksha

### Step 3.1: Module Rename
```bash
# Automated rename script
cd x/grampension

# 1. Rename files
find . -name "*pension*" -exec bash -c 'mv "$1" "${1//pension/suraksha}"' - {} \;

# 2. Update imports
find . -type f -name "*.go" -exec sed -i 's/grampension/gramsuraksha/g' {} +

# 3. Update types
sed -i 's/PensionScheme/SurakshaScheme/g' types/*.go
sed -i 's/GetPension/GetSuraksha/g' keeper/*.go
```

### Step 3.2: Update Scheme Types
```go
// BEFORE: x/grampension/types/scheme.go
type PensionScheme struct {
    MinInvestment sdk.Int
    MaturityYears uint32
    GuaranteedReturn sdk.Dec // 50%
}

// AFTER: x/gramsuraksha/types/scheme.go  
type SurakshaPool struct {
    MinDeposit sdk.Int      // DeFi terminology
    LockPeriod uint32       // Not "maturity"
    TargetAPY  sdk.Dec      // Not "guaranteed"
    // Add DeFi fields
    YieldStrategy string    // "validator-staking", "liquidity-provision"
    RiskLevel     string    // "low", "medium", "high"
}
```

## 🔄 Task 4: Integrate Money Order DEX

### Step 4.1: Connect to Frontend
```typescript
// EXISTING: x/moneyorder is complete!

// NEW: apps/dhansetu/src/services/MoneyOrderService.ts
import { MoneyOrderClient } from '@deshchain/sdk';

export class MoneyOrderService {
  async createP2POrder(params: OrderParams) {
    // Use existing proto definitions
    const msg = {
      typeUrl: '/deshchain.moneyorder.MsgCreateOrder',
      value: {
        creator: params.creator,
        orderType: 'p2p',
        amount: params.amount,
        rate: params.rate,
      },
    };
    
    return await this.client.signAndBroadcast(msg);
  }
}
```

### Step 4.2: Build UI Components
```typescript
// NEW: apps/dhansetu/src/screens/Exchange/MoneyOrderDEX.tsx
export const MoneyOrderDEX = () => {
  // Reuse components from frontend/apps/money-order-simple
  return (
    <View>
      <SevaMitraDiscovery />  {/* Exists */}
      <OrderBook />           {/* Exists */}
      <TradeForm />          {/* Exists */}
      <EscrowStatus />       {/* Build new */}
    </View>
  );
};
```

## 🆕 Task 5: Build DhanPata (New Module)

### Step 5.1: Create Module Structure
```bash
# Copy structure from money order (best implementation)
cp -r x/moneyorder x/dhanpata
cd x/dhanpata

# Update module name
find . -type f -name "*.go" -exec sed -i 's/moneyorder/dhanpata/g' {} +
```

### Step 5.2: Implement Core Logic
```go
// NEW: x/dhanpata/keeper/registration.go
func (k Keeper) RegisterName(ctx sdk.Context, name string, owner string, duration uint32) error {
    // Check availability
    if k.IsNameTaken(ctx, name) {
        return errors.New("name already taken")
    }
    
    // Calculate price
    price := k.CalculatePrice(ctx, name, duration)
    
    // Deduct payment
    err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, price)
    if err != nil {
        return err
    }
    
    // Store registration
    registration := types.DhanPataRegistration{
        Name:      name,
        Owner:     owner,
        ExpiresAt: ctx.BlockTime().Add(time.Duration(duration) * time.Second),
    }
    
    k.SetRegistration(ctx, registration)
    return nil
}
```

## 🎉 Task 6: Unify Festival System

### Step 6.1: Move Frontend Logic to Backend
```go
// EXISTING: frontend/packages/money-order-ui/src/constants/culturalData.ts
export const festivals = [
  { name: 'Diwali', discount: 0.5, duration: 5 },
  // ...
];

// MOVE TO: x/cultural/types/festivals.go
var Festivals = []Festival{
    {
        Name:     "Diwali",
        Type:     "religious",
        Duration: 5 * 24 * time.Hour,
        PreFestivalDays: 3,
        PostFestivalDays: 3,
        MaxBonus: sdk.NewDecWithPrec(50, 2), // 50%
        Greetings: map[string]string{
            "en": "Happy Diwali!",
            "hi": "दीपावली की शुभकामनाएं!",
        },
    },
    // Add all 500+ festivals
}
```

### Step 6.2: Create Festival Keeper
```go
// NEW: x/cultural/keeper/festivals.go
func (k Keeper) GetActiveFestival(ctx sdk.Context, pincode string) (*types.Festival, error) {
    currentTime := ctx.BlockTime()
    
    // Check national festivals
    for _, festival := range types.Festivals {
        if k.IsFestivalActive(festival, currentTime) {
            return &festival, nil
        }
    }
    
    // Check local festivals by pincode
    localFestivals := k.GetLocalFestivals(ctx, pincode)
    for _, festival := range localFestivals {
        if k.IsFestivalActive(festival, currentTime) {
            return &festival, nil
        }
    }
    
    return nil, nil
}
```

## 🚀 Task 7: Complete NAMO Token

### Step 7.1: Add Missing Implementation
```go
// EXISTING: x/namo/types/ has types but no logic

// NEW: x/namo/keeper/mint.go
func (k Keeper) MintTokens(ctx sdk.Context, amount sdk.Int, recipient string) error {
    // Only treasury can mint
    if !k.IsTreasury(ctx, ctx.MsgSender()) {
        return errors.New("unauthorized")
    }
    
    coins := sdk.NewCoins(sdk.NewCoin("namo", amount))
    err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
    if err != nil {
        return err
    }
    
    return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins)
}

// NEW: x/namo/keeper/burn.go
func (k Keeper) BurnTokens(ctx sdk.Context, amount sdk.Int) error {
    coins := sdk.NewCoins(sdk.NewCoin("namo", amount))
    err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, ctx.MsgSender(), types.ModuleName, coins)
    if err != nil {
        return err
    }
    
    return k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
}
```

## 📱 Task 8: Create Main DhanSetu App

### Step 8.1: Setup App Structure
```bash
# Use existing React Native setup
cd apps/dhansetu

# Install dependencies
npm install @react-navigation/native
npm install @reduxjs/toolkit react-redux
npm install react-native-svg react-native-linear-gradient
```

### Step 8.2: Create Navigation
```typescript
// NEW: apps/dhansetu/src/navigation/AppNavigator.tsx
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';

const Tab = createBottomTabNavigator();

export const AppNavigator = () => {
  return (
    <Tab.Navigator>
      <Tab.Screen name="Wallet" component={WalletScreen} />      {/* From Batua */}
      <Tab.Screen name="Exchange" component={ExchangeScreen} />   {/* Money Order */}
      <Tab.Screen name="LocalCoins" component={LocalCoinsScreen} />{/* New */}
      <Tab.Screen name="DeFi" component={DeFiScreen} />          {/* Suraksha */}
      <Tab.Screen name="Festival" component={FestivalScreen} />   {/* Cultural */}
    </Tab.Navigator>
  );
};
```

## 🏃 Quick Implementation Order

### Week 1: Foundation
1. **Day 1-2**: Setup project, rename pension → suraksha
2. **Day 3-4**: Port Batua wallet core to TypeScript
3. **Day 5**: Enhance cultural module with existing quotes

### Week 2: Integration
1. **Day 1-2**: Connect Money Order DEX to frontend
2. **Day 3-4**: Build main app navigation
3. **Day 5**: Test existing features

### Week 3-4: New Features
1. **Week 3**: Build DhanPata module
2. **Week 4**: Create Kshetra Coins factory

### Week 5-6: Completion
1. **Week 5**: Complete NAMO token, Sikkebaaz
2. **Week 6**: Testing and bug fixes

## 💡 Pro Tips

1. **Copy Module Pattern**: Always copy from `x/moneyorder/` - it's the best implemented
2. **Reuse Proto Files**: Don't recreate message types, extend existing ones
3. **Test Incrementally**: Test each refactored module before moving to next
4. **Keep Cultural Data**: The 10,000+ quotes and cultural data are gold - use them!
5. **Leverage Existing UI**: Money Order UI components are production-ready

---

**"Refactor with Purpose, Build with Precision"** 🎯