# DeshChain Explorer Frontend

A comprehensive blockchain explorer for DeshChain that combines cultural heritage with advanced blockchain analytics.

## 🌟 Features

### Core Explorer Functionality
- **Real-time Blockchain Data**: Live updates of blocks, transactions, and network status
- **Advanced Search**: Search transactions, blocks, addresses, and lending data
- **Multi-module Support**: Explorer for lending modules (Krishi, Vyavasaya, Shiksha Mitra)
- **Validator Monitoring**: Track validator performance and network health
- **Transaction Analytics**: Detailed transaction analysis and history

### Cultural Integration
- **Festival Celebrations**: Dynamic theming during Indian festivals
- **Multilingual Support**: 10+ Indian languages including Hindi, Bengali, Tamil
- **Cultural Quotes**: Daily inspirational quotes from Indian heritage
- **Festival Bonuses**: Real-time tracking of festival-based transaction bonuses

### DeFi Analytics
- **Lending Module Stats**: Comprehensive analytics for all three lending modules
- **Interest Rate Tracking**: Real-time rate monitoring across modules
- **Default Rate Analysis**: Risk assessment and loan performance metrics
- **Cultural Impact Metrics**: How festivals and events affect lending patterns

### User Experience
- **Dark/Light Themes**: Automatic theme switching with festival themes
- **Responsive Design**: Optimized for mobile, tablet, and desktop
- **Real-time Updates**: WebSocket connections for live data
- **Performance Optimized**: Fast loading with efficient caching

## 🏗️ Architecture

### Technology Stack
- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS with custom cultural themes
- **Animations**: Framer Motion
- **Charts**: Recharts for data visualization
- **State Management**: React Query + Context API
- **Blockchain Integration**: CosmJS for DeshChain connectivity

### Project Structure
```
├── app/                    # Next.js App Router
│   ├── layout.tsx         # Root layout with providers
│   ├── page.tsx           # Home dashboard
│   ├── providers.tsx      # Global providers setup
│   └── globals.css        # Global styles and themes
├── components/            # Reusable components
│   ├── dashboard/         # Dashboard components
│   ├── layout/           # Layout components
│   ├── search/           # Search functionality
│   ├── cultural/         # Cultural heritage components
│   └── network/          # Network status components
├── hooks/                # Custom React hooks
├── providers/            # Context providers
├── utils/               # Utility functions
└── types/               # TypeScript type definitions
```

## 🚀 Getting Started

### Prerequisites
- Node.js 18+ 
- npm or yarn
- DeshChain node running locally or access to testnet

### Installation

1. **Clone and install dependencies**:
   ```bash
   cd frontend/apps/deshchain-explorer
   npm install
   ```

2. **Environment Setup**:
   ```bash
   cp .env.example .env.local
   ```
   
   Configure your environment variables:
   ```env
   NEXT_PUBLIC_CHAIN_ID=deshchain-1
   NEXT_PUBLIC_RPC_ENDPOINT=http://localhost:26657
   NEXT_PUBLIC_REST_ENDPOINT=http://localhost:1317
   NEXT_PUBLIC_EXPLORER_NAME=DeshChain Explorer
   ```

3. **Development Server**:
   ```bash
   npm run dev
   ```
   
   Open [http://localhost:3001](http://localhost:3001) in your browser.

4. **Production Build**:
   ```bash
   npm run build
   npm start
   ```

## 🎨 Cultural Features

### Festival Integration
The explorer automatically detects and celebrates Indian festivals:

- **Dynamic Theming**: Color schemes change based on festivals
- **Bonus Tracking**: Real-time festival bonus calculations
- **Cultural Quotes**: Daily rotating quotes from Indian heritage
- **Event Calendar**: Upcoming cultural events and celebrations

### Multilingual Support
Supports 10+ Indian languages:
- English, Hindi (हिन्दी), Bengali (বাংলা)
- Tamil (தமிழ்), Telugu (తెలుగు), Gujarati (ગુજરાતી)
- Marathi (मराठी), Kannada (ಕನ್ನಡ), Malayalam (മലയാളം)
- Punjabi (ਪੰਜਾਬੀ)

### Cultural UI Elements
- Traditional Indian color palettes (Saffron, White, Green)
- Cultural patterns and motifs
- Festival-specific animations and effects
- Region-specific customizations

## 📊 Analytics Dashboard

### Real-time Metrics
- **Block Height**: Current blockchain height with trend indicators
- **Transaction Volume**: Real-time transaction monitoring
- **Network Health**: Validator status and network performance
- **NAMO Supply**: Token supply tracking and distribution

### Lending Module Analytics
- **Krishi Mitra**: Agricultural lending statistics (6-9% rates)
- **Vyavasaya Mitra**: Business lending metrics (8-12% rates)  
- **Shiksha Mitra**: Education loan analytics (4-7% rates)
- **Cross-module Insights**: Comparative analysis and trends

### Performance Metrics
- **Default Rates**: Risk analysis across modules
- **Interest Rate Trends**: Historical rate movements
- **Geographic Distribution**: State-wise lending patterns
- **Cultural Impact**: Festival effects on lending activity

## 🔍 Search Functionality

### Advanced Search
- **Transaction Hashes**: Direct transaction lookup
- **Block Numbers**: Block details and transaction lists
- **Addresses**: Account balances and transaction history
- **Lending Data**: Loan applications and status tracking

### Smart Suggestions
- Auto-complete for addresses and transaction hashes
- Recent search history
- Quick links to popular sections
- Context-aware suggestions

### Filter Options
- Date range filtering
- Transaction type filtering
- Amount range filtering
- Module-specific filtering

## 🎯 Components Guide

### Core Components

#### StatCard
```tsx
<StatCard
  title="Total Transactions"
  value={1234567}
  icon={<TrendingUp />}
  trend={12.5}
  color="green"
  format="number"
/>
```

#### LendingModuleStats
```tsx
<LendingModuleStats 
  stats={lendingData}
  loading={false}
/>
```

#### FestivalBanner
```tsx
<FestivalBanner 
  festival={currentFestival}
/>
```

#### QuickSearch
```tsx
<QuickSearch />
```

### Cultural Components

#### CulturalProvider
Manages festival data, quotes, and cultural events:
```tsx
const { currentFestival, dailyQuote, setLanguage } = useCultural()
```

#### ExplorerProvider
Handles blockchain connectivity and data:
```tsx
const { chainInfo, isConnected, getTransaction } = useExplorer()
```

## 🛠️ Development

### Code Standards
- TypeScript strict mode
- ESLint + Prettier configuration
- Component documentation required
- Cultural sensitivity guidelines

### Testing
```bash
# Unit tests
npm run test

# Component tests
npm run test:components

# E2E tests
npm run test:e2e

# Performance tests
npm run test:performance
```

### Build Optimization
- Automatic code splitting
- Image optimization
- Bundle analysis
- Performance monitoring

## 🌐 Deployment

### Vercel (Recommended)
```bash
npm run build
vercel deploy
```

### Docker
```bash
docker build -t deshchain-explorer .
docker run -p 3001:3001 deshchain-explorer
```

### Static Export
```bash
npm run export
```

## 📱 Mobile Support

### Progressive Web App
- Offline support for core features
- App-like experience on mobile
- Push notifications for important events
- Home screen installation

### Responsive Design
- Mobile-first approach
- Touch-optimized interactions
- Adaptive navigation
- Performance optimizations

## 🔒 Security

### Data Protection
- No private key handling
- Read-only blockchain access
- Secure API endpoints
- HTTPS enforcement

### Privacy
- No user tracking
- Local preference storage
- Minimal data collection
- GDPR compliance

## 🤝 Contributing

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Follow coding standards
4. Add tests for new features
5. Submit pull request

### Cultural Sensitivity
- Respect for all Indian cultures
- Accurate cultural representations
- Inclusive design principles
- Community feedback integration

### Feature Requests
- Cultural feature suggestions welcome
- Accessibility improvements priority
- Performance optimizations valued
- Multi-language support expansion

## 📋 Roadmap

### Upcoming Features
- **Mobile App**: React Native version
- **Advanced Analytics**: ML-powered insights
- **Social Features**: Community discussions
- **API Integration**: REST and GraphQL APIs

### Cultural Enhancements
- **Regional Themes**: State-specific customizations
- **Cultural Calendar**: Comprehensive event tracking
- **Heritage Stories**: Historical context integration
- **Community Features**: User-generated cultural content

## 🐛 Troubleshooting

### Common Issues
1. **Connection Failed**: Check RPC endpoint configuration
2. **Slow Loading**: Verify network connection and endpoint status
3. **Theme Issues**: Clear browser cache and refresh
4. **Search Problems**: Ensure proper input formatting

### Debug Mode
```bash
NEXT_PUBLIC_DEBUG=true npm run dev
```

### Performance Monitoring
Built-in performance metrics and error tracking available in development mode.

## 📄 License

MIT License - see LICENSE file for details.

## 🙏 Acknowledgments

- DeshChain development team
- Indian cultural heritage consultants
- Open source community contributors
- Festival and cultural celebration experts

---

**Built with ❤️ for preserving Indian cultural heritage while embracing blockchain innovation.**