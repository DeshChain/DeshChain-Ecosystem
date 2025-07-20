# Money Order UI Components Library

A comprehensive React component library for DeshChain's Money Order system with integrated Indian cultural features, multi-language support, and festival themes.

## Features

### 🎯 Core Components
- **MoneyOrderForm** - Complete form with validation and cultural integration
- **MoneyOrderReceipt** - Digital receipts with QR codes and cultural quotes
- **PoolSelector** - Pool selection with cultural themes and trust scores
- **TradingInterface** - Professional trading interface for advanced users

### 🎨 Cultural Integration
- **22 Indian Languages** - Full support for major Indian languages
- **Cultural Quotes** - 10,000+ curated quotes from Indian leaders and texts
- **Festival Themes** - Special UI themes during Indian festivals
- **Patriotism Scoring** - Gamified system rewarding community contribution

### 🎭 Festival Features
- **Diwali** - Golden theme with 15% transaction bonus
- **Holi** - Colorful theme with community bonuses
- **Independence Day** - Patriotic theme with special rewards
- **Regional Festivals** - State-specific celebrations and bonuses

## Installation

```bash
npm install @deshchain/money-order-ui
```

## Quick Start

```tsx
import React from 'react';
import {
  MoneyOrderProvider,
  CulturalProvider,
  MoneyOrderForm
} from '@deshchain/money-order-ui';

function App() {
  const config = {
    apiUrl: 'https://api.deshchain.org/v1/moneyorder',
    chainId: 'deshchain-1',
    enableCulturalFeatures: true,
    enableFestivalThemes: true
  };

  const handleMoneyOrderSubmit = async (data) => {
    console.log('Money order created:', data);
  };

  return (
    <MoneyOrderProvider config={config}>
      <CulturalProvider initialLanguage="hi">
        <MoneyOrderForm
          onSubmit={handleMoneyOrderSubmit}
          showCulturalFeatures={true}
          mode="simple"
        />
      </CulturalProvider>
    </MoneyOrderProvider>
  );
}
```

## Components

### MoneyOrderForm

A comprehensive form component for creating money orders with cultural integration.

```tsx
import { MoneyOrderForm } from '@deshchain/money-order-ui';

<MoneyOrderForm
  onSubmit={handleSubmit}
  mode="simple" // or "advanced"
  showCulturalFeatures={true}
  initialData={{
    culturalPreferences: {
      language: 'hi',
      theme: 'prosperity',
      includeQuote: true
    }
  }}
/>
```

#### Props
- `onSubmit: (data: MoneyOrderFormData) => Promise<void>` - Form submission handler
- `mode?: 'simple' | 'advanced'` - Form complexity level
- `showCulturalFeatures?: boolean` - Enable cultural features
- `initialData?: Partial<MoneyOrderFormData>` - Pre-fill form data
- `pools?: PoolInfo[]` - Available liquidity pools

### CulturalQuote

Display cultural quotes with rich formatting and interactions.

```tsx
import { CulturalQuote } from '@deshchain/money-order-ui';

<CulturalQuote
  quote={quoteData}
  variant="card" // or "banner", "inline"
  showActions={true}
  animation="fade"
  onFavorite={(quoteId) => console.log('Favorited:', quoteId)}
/>
```

### FestivalBanner

Celebrate Indian festivals with themed banners and bonuses.

```tsx
import { FestivalBanner } from '@deshchain/money-order-ui';

<FestivalBanner
  festival={festivalData}
  variant="full" // or "compact", "minimal"
  showBonusInfo={true}
  showGreeting={true}
/>
```

### PoolSelector

Select liquidity pools with cultural themes and performance metrics.

```tsx
import { PoolSelector } from '@deshchain/money-order-ui';

<PoolSelector
  pools={poolData}
  onSelect={(pool) => console.log('Selected:', pool)}
  showCulturalThemes={true}
  filterByRegion="north_india"
/>
```

## Hooks

### useMoneyOrder

Main hook for Money Order operations.

```tsx
import { useMoneyOrder } from '@deshchain/money-order-ui';

const {
  createMoneyOrder,
  getSwapQuote,
  executeSwap,
  addLiquidity,
  isLoading,
  error
} = useMoneyOrder();
```

### useCulturalContext

Hook for cultural features and localization.

```tsx
import { useCulturalContext } from '@deshchain/money-order-ui';

const {
  currentLanguage,
  currentQuote,
  currentFestival,
  changeLanguage,
  getContextualQuote,
  isCulturalFeaturesEnabled
} = useCulturalContext();
```

### usePoolData

Hook for pool information and analytics.

```tsx
import { usePoolData } from '@deshchain/money-order-ui';

const {
  pools,
  getPoolPerformance,
  getVillagePools,
  refreshPools,
  isLoading
} = usePoolData();
```

## Language Support

The library supports 22 Indian languages with proper scripts and cultural context:

- **Hindi** (हिन्दी) - Devanagari script
- **Bengali** (বাংলা) - Bengali script
- **Telugu** (తెలుగు) - Telugu script
- **Tamil** (தமிழ்) - Tamil script
- **Marathi** (मराठी) - Devanagari script
- **Gujarati** (ગુજરાતી) - Gujarati script
- **Kannada** (ಕನ್ನಡ) - Kannada script
- **Malayalam** (മലയാളം) - Malayalam script
- **Punjabi** (ਪੰਜਾਬੀ) - Gurmukhi script
- **Odia** (ଓଡ଼ିଆ) - Odia script
- **Assamese** (অসমীয়া) - Bengali-Assamese script
- **Urdu** (اردو) - Arabic script
- **Sanskrit** (संस्कृत) - Devanagari script
- And more regional languages...

## Cultural Themes

### Festivals with Special Themes
- **Diwali** - Golden prosperity theme with 15% bonus
- **Holi** - Colorful unity theme with 10% bonus
- **Independence Day** - Patriotic theme with 20% bonus
- **Eid** - Community celebration theme with 12% bonus
- **Durga Puja** - Divine power theme with 12% bonus

### Theme Categories
- **Prosperity** - Success and wealth themes
- **Unity** - Community and togetherness
- **Patriotism** - National pride and service
- **Wisdom** - Knowledge and enlightenment
- **Family** - Traditional values and relationships

## Configuration

```tsx
const config = {
  // API Configuration
  apiUrl: 'https://api.deshchain.org/v1/moneyorder',
  chainId: 'deshchain-1',
  
  // Cultural Features
  defaultLanguage: 'hi',
  enableCulturalFeatures: true,
  enableFestivalThemes: true,
  enablePatriotismRewards: true,
  
  // Trading Configuration
  maxSlippage: 0.05,
  defaultPoolType: 'fixed_rate',
  autoSelectPool: true,
  
  // Theme Configuration
  theme: {
    primary: '#FF6B35', // Saffron
    secondary: '#138808', // Green
    accent: '#000080', // Navy Blue
    culturalElements: {
      borderStyle: 'traditional',
      pattern: 'mandala',
      iconSet: 'indian_classical'
    }
  }
};
```

## Styling and Theming

The library uses Material-UI with custom cultural themes:

```tsx
import { ThemeProvider } from '@deshchain/money-order-ui';

<ThemeProvider theme={customTheme}>
  <YourApp />
</ThemeProvider>
```

### Festival-specific Styling
During festivals, the theme automatically adapts:

```tsx
// Diwali theme colors
const diwaliTheme = {
  primary: '#FFD700', // Gold
  secondary: '#FF4500', // Orange Red
  accent: '#8B0000' // Dark Red
};
```

## Validation

Built-in validation for Indian-specific formats:

```tsx
const validationRules = {
  // Address validation for DeshChain format
  addressPattern: /^desh1[a-z0-9]{38}$/,
  
  // Indian postal code validation
  postalCodePattern: /^[1-9][0-9]{5}$/,
  
  // Amount limits
  minAmount: '1',
  maxAmount: '10000000',
  
  // Cultural preferences
  supportedLanguages: ['hi', 'en', 'bn', 'te', ...],
  culturalThemes: ['prosperity', 'unity', 'patriotism', ...]
};
```

## Testing

```bash
# Run unit tests
npm test

# Run component tests with Storybook
npm run storybook

# Run integration tests
npm run test:integration
```

## Building

```bash
# Build the library
npm run build

# Build with watch mode for development
npm run build:watch

# Type checking
npm run type-check
```

## Contributing

1. **Cultural Accuracy** - Ensure all cultural references are respectful and accurate
2. **Language Support** - Test with multiple Indian languages and scripts
3. **Accessibility** - Follow WCAG guidelines for inclusive design
4. **Performance** - Optimize for mobile devices common in India

### Adding New Cultural Features

```tsx
// Add new festival
const newFestival: FestivalInfo = {
  festivalId: 'karva_chauth',
  name: 'Karva Chauth',
  description: 'Festival of married women',
  bonusRate: 0.08,
  culturalTheme: 'family',
  region: 'north_india'
};

// Add new cultural quote
const newQuote: CulturalQuoteData = {
  text: 'माता पिता गुरु देवा',
  author: 'Traditional Sanskrit',
  category: 'family',
  language: 'sanskrit',
  translation: 'Mother, Father, Teacher, God'
};
```

## Browser Support

- **Modern Browsers** - Chrome 70+, Firefox 65+, Safari 12+
- **Mobile Browsers** - iOS Safari 12+, Chrome Mobile 70+
- **Indian Regional Browsers** - UC Browser, Opera Mini support

## Performance

- **Bundle Size** - ~180KB gzipped
- **Tree Shaking** - Full support for unused component elimination
- **Lazy Loading** - Cultural data loaded on demand
- **Caching** - Intelligent caching of quotes and festival data

## License

Apache License 2.0 - See [LICENSE](./LICENSE) for details.

## Support

- **Documentation** - [docs.deshchain.org/ui-components](https://docs.deshchain.org/ui-components)
- **Issues** - [GitHub Issues](https://github.com/deshchain/deshchain/issues)
- **Community** - [Discord](https://discord.gg/deshchain)
- **Cultural Consultancy** - For adding new cultural features

---

**Made with ❤️ for the Indian blockchain community** 🇮🇳