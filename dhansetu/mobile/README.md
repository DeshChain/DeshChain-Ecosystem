# DhanSetu Mobile App

## Revolutionary DeFi Super App for Bharat 🇮🇳

DhanSetu is a comprehensive React Native mobile application that brings the power of DeshChain's revolutionary financial ecosystem to every Indian's fingertips. Built with cultural values at its core, DhanSetu combines traditional Indian financial wisdom with cutting-edge blockchain technology.

## 🌟 Key Features

### 💰 NAMO Token Integration
- Native NAMO token wallet with HD wallet support
- Multi-chain support (DeshChain, Ethereum, Bitcoin)
- Secure biometric authentication
- Cultural quotes on every transaction

### 💱 Money Order DEX
- Traditional money order system reimagined for blockchain
- PIN-based security for recipients
- Lowest fees in the industry (0.1%)
- Festival-themed UI during Indian celebrations

### 🚀 Sikkebaaz Memecoin Launchpad
- Launch your own Kshetra Coins (local memecoins)
- Anti-pump & dump protection
- Community veto mechanism
- Cultural integration with local festivals

### 🛡️ Gram Suraksha Pension
- Revolutionary minimum 8% guaranteed returns, up to 50% based on platform performance
- Community-pooled pension system
- Transparent on-chain management
- Early maturity options

### 💸 DhanPata Virtual Addresses
- Human-readable payment addresses (@username)
- No more complex wallet addresses
- Instant transfers with cultural messages
- QR code integration

### 🌾 Lending Suite (Mitra Services)
- **Krishi Mitra**: Agricultural loans at 6-9% interest
- **Vyavasaya Mitra**: Business loans for entrepreneurs
- **Shiksha Mitra**: Education loans with performance rewards
- Member-only exclusive rates

## 🚀 Getting Started

### Prerequisites
- Node.js >= 16
- React Native development environment
- Android Studio / Xcode
- Expo CLI (optional but recommended)

### Installation

```bash
# Clone the repository
git clone https://github.com/deshchain/dhansetu-mobile.git
cd dhansetu-mobile

# Install dependencies
npm install

# iOS specific
cd ios && pod install && cd ..

# Start Metro bundler
npm start

# Run on Android
npm run android

# Run on iOS
npm run ios
```

### Development Setup

1. **Environment Configuration**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Connect to DeshChain**
   - Default RPC: `https://rpc.deshchain.com`
   - Chain ID: `deshchain-1`
   - Explorer: `https://explorer.deshchain.com`

## 🏗️ Architecture

### Tech Stack
- **Framework**: React Native with Expo
- **State Management**: Redux Toolkit
- **Navigation**: React Navigation v6
- **Blockchain**: CosmJS, Ethers.js
- **UI Components**: React Native Paper + Custom Cultural Components
- **Security**: Biometric authentication, Secure storage
- **Styling**: Styled Components + Theme system

### Project Structure
```
src/
├── components/       # Reusable UI components
│   ├── common/      # Generic components
│   └── cultural/    # Cultural-themed components
├── screens/         # App screens
│   ├── onboarding/  # Wallet setup flows
│   ├── home/        # Dashboard
│   ├── wallet/      # Token management
│   ├── dex/         # Money Order DEX
│   ├── sikkebaaz/   # Memecoin launchpad
│   └── suraksha/    # Pension system
├── services/        # Core services
│   ├── blockchain/  # DeshChain client
│   ├── wallet/      # HD wallet implementation
│   └── cultural/    # Festival & quote services
├── store/           # Redux store
├── navigation/      # Navigation configuration
└── constants/       # App constants & theme
```

## 🎨 Cultural UI Features

### Dynamic Festival Themes
The app automatically adapts its UI during Indian festivals:
- **Diwali**: Golden gradients and diya icons
- **Holi**: Vibrant colors and playful animations
- **Independence Day**: Tricolor theme
- **Republic Day**: National emblems

### Multilingual Support
- Hindi (हिंदी)
- English
- 20+ Indian regional languages
- Sanskrit quotes for spiritual touch

### Cultural Elements
- 10,000+ curated Indian quotes
- Festival greetings on transactions
- Traditional motifs in UI design
- Patriotism score for users

## 🔐 Security Features

### Wallet Security
- BIP39/BIP32 HD wallet standard
- Secure encrypted storage
- Biometric authentication
- PIN protection with attempt limits
- Automatic session timeout

### Transaction Security
- Multi-signature support
- Hardware wallet integration (coming soon)
- Transaction simulation
- Phishing protection

## 🧪 Testing

```bash
# Run unit tests
npm test

# Run integration tests
npm run test:integration

# E2E tests with Detox
npm run test:e2e
```

## 📱 Build & Release

### Android
```bash
# Debug build
npm run android:debug

# Release build
npm run android:release

# Generate signed APK
cd android && ./gradlew assembleRelease
```

### iOS
```bash
# Debug build
npm run ios:debug

# Release build
npm run ios:release

# Archive for App Store
# Use Xcode or fastlane
```

## 🤝 Contributing

We welcome contributions from the community! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting PRs.

### Development Guidelines
1. Follow the existing code style
2. Add tests for new features
3. Update documentation
4. Ensure cultural sensitivity
5. Optimize for low-end devices

## 📄 License

Licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- Built with ❤️ for Bharat
- Inspired by traditional Indian financial systems
- Powered by DeshChain blockchain
- Cultural quotes curated by community

## 📞 Support

- **Discord**: [Join our community](https://discord.gg/deshchain)
- **Telegram**: [@DhanSetuSupport](https://t.me/dhansetu)
- **Email**: support@dhansetu.in
- **Website**: [dhansetu.in](https://dhansetu.in)

---

**DhanSetu** - *Your Bridge to Financial Freedom* 🌉

Jai Hind! 🇮🇳