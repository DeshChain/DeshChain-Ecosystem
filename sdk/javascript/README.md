# DeshChain JavaScript/TypeScript SDK

Official SDK for interacting with the DeshChain blockchain - a cultural heritage-focused blockchain platform built on Cosmos SDK.

## Features

- 🌟 **Cultural Integration**: Festival bonuses, regional quotes, and cultural heritage preservation
- 💰 **DeFi Lending**: Krishi Mitra (Agriculture), Vyavasaya Mitra (Business), Shiksha Mitra (Education)
- 🚀 **Sikkebaaz Launchpad**: Anti-pump memecoin platform with community governance
- 💸 **Money Order DEX**: Traditional money transfers reimagined for blockchain
- 🏛️ **Governance**: Founder protection with community participation
- 🎭 **22 Languages**: Multi-language support for Indian regional languages

## Installation

```bash
npm install @deshchain/sdk
# or
yarn add @deshchain/sdk
```

## Quick Start

```typescript
import { DeshChainClient } from '@deshchain/sdk'

// Connect to DeshChain network
const client = await DeshChainClient.connect('https://rpc.deshchain.network')

// Get chain info
const chainInfo = await client.getChainInfo()
console.log('Connected to:', chainInfo.chainId)

// Check current festival
const festival = await client.getCurrentFestival()
if (festival) {
  console.log('Active festival:', festival.name)
}

// Get lending statistics
const lendingStats = await client.getLendingStats()
console.log('Total loans disbursed:', lendingStats.combined.totalDisbursed)
```

## Module Clients

### Cultural Heritage

```typescript
// Get daily cultural quote
const quote = await client.cultural.getDailyQuote()
console.log(quote.text, '-', quote.author)

// Get active festivals
const festivals = await client.cultural.getActiveFestivals()

// Calculate festival bonus
const bonus = await client.cultural.calculateFestivalBonus(1000, 'send')
```

### Lending (Krishi/Vyavasaya/Shiksha Mitra)

```typescript
// Get agricultural loan rates
const rateQuote = await client.lending.getAgriculturalRateQuote(
  50000, // amount
  12,    // duration in months
  'wheat', // crop type
  2.5    // land size in acres
)

// Check farmer eligibility
const eligibility = await client.lending.checkFarmerEligibility(application)

// Get education scholarships
const scholarships = await client.lending.getScholarships()
```

### Sikkebaaz Launchpad

```typescript
// Get featured tokens
const featuredTokens = await client.sikkebaaz.getFeaturedTokens()

// Check anti-pump restrictions
const canTrade = await client.sikkebaaz.canTrade(address, 'MYTOKEN', 1000)

// Get community votes
const votes = await client.sikkebaaz.getCommunityVotes('MYTOKEN')
```

### Money Order DEX

```typescript
// Get trading pairs
const pairs = await client.moneyOrder.getTradingPairs()

// Calculate trade
const tradeInfo = await client.moneyOrder.calculateTrade(
  'NAMO', 'USDC', 'buy', 1000, 'market'
)

// Get money order fees
const fees = await client.moneyOrder.calculateMoneyOrderFees(
  5000, '110001', '400001' // amount, source pin, destination pin
)
```

### Governance

```typescript
// Get active proposals
const proposals = await client.governance.getActiveProposals()

// Get voting power
const votingPower = await client.governance.getVotingPower(address)

// Check founder protection
const protection = await client.governance.checkFounderProtection(proposalId)
```

## Network Configuration

```typescript
// Testnet
const testnetClient = await DeshChainClient.connect(
  'https://testnet-rpc.deshchain.network',
  { chainId: 'deshchain-testnet-1' }
)

// Custom configuration
const client = await DeshChainClient.connect(endpoint, {
  chainId: 'deshchain-1',
  prefix: 'deshchain',
  gasPrice: '0.025unamo'
})
```

## Cultural Features

### Festival Integration

```typescript
import { FestivalUtils } from '@deshchain/sdk'

// Check current festivals
const { isFestival, festivals } = FestivalUtils.isFestivalDay()

// Get festival bonus
const bonus = FestivalUtils.getFestivalBonus()

// Get upcoming festivals
const upcoming = FestivalUtils.getUpcomingFestivals(30) // next 30 days
```

### Regional Customization

```typescript
import { CulturalUtils } from '@deshchain/sdk'

// Get state from pincode
const state = CulturalUtils.getStateFromPincode('110001')

// Get regional language
const language = CulturalUtils.getRegionalLanguage('Maharashtra')

// Get major festivals by state
const festivals = CulturalUtils.getMajorFestivals('West Bengal')

// Format currency in Indian style
const formatted = CulturalUtils.formatIndianCurrency(150000) // ₹1,50,000
```

## Error Handling

```typescript
import { DeshChainError, NetworkError, TransactionError } from '@deshchain/sdk'

try {
  const result = await client.someOperation()
} catch (error) {
  if (error instanceof NetworkError) {
    console.error('Network issue:', error.message)
  } else if (error instanceof TransactionError) {
    console.error('Transaction failed:', error.message, error.txHash)
  } else if (error instanceof DeshChainError) {
    console.error('DeshChain error:', error.code, error.message)
  }
}
```

## Validation Utilities

```typescript
import { ValidationUtils } from '@deshchain/sdk'

// Validate addresses
const { valid, error } = ValidationUtils.validateAddress('deshchain1...')

// Validate Indian-specific data
ValidationUtils.validatePincode('110001')
ValidationUtils.validateAadhar('1234 5678 9012')
ValidationUtils.validatePAN('ABCDE1234F')
ValidationUtils.validateGST('07AAGFF2194N1Z1')
```

## TypeScript Support

The SDK is written in TypeScript and provides full type definitions:

```typescript
import type { 
  Festival, 
  Loan, 
  LaunchpadToken, 
  MoneyOrder, 
  Proposal 
} from '@deshchain/sdk'

// All API responses are fully typed
const festivals: Festival[] = await client.cultural.getActiveFestivals()
const loans: Loan[] = await client.lending.getFarmerLoans(farmerId)
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Documentation

- [API Documentation](https://docs.deshchain.network/sdk)
- [Cultural Integration Guide](https://docs.deshchain.network/cultural)
- [Lending Modules](https://docs.deshchain.network/lending)
- [Sikkebaaz Launchpad](https://docs.deshchain.network/sikkebaaz)

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- GitHub Issues: [deshchain/deshchain-sdk-js/issues](https://github.com/deshchain/deshchain-sdk-js/issues)
- Discord: [DeshChain Community](https://discord.gg/deshchain)
- Documentation: [docs.deshchain.network](https://docs.deshchain.network)

---

**Jai Hind! 🇮🇳** - Building the future of Indian blockchain technology with cultural values at its core.