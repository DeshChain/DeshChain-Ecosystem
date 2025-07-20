# Money Order Simple Interface

A user-friendly, culturally-integrated money transfer application built on DeshChain. This simple interface makes blockchain-based money orders accessible to everyone with support for 22 Indian languages and festival celebrations.

## Features

### 🎯 Simple Money Transfers
- **One-Page Interface** - Everything you need on a single, intuitive screen
- **Address Book** - Save and manage frequently used addresses
- **Amount Presets** - Quick selection for common amounts (₹100 - ₹10,000)
- **Real-time Exchange Rates** - Live NAMO to INR conversion
- **Transaction Preview** - Review before sending

### 🎨 Cultural Integration
- **Festival Bonuses** - Special rewards during Indian festivals (up to 20%)
- **Cultural Quotes** - Inspirational quotes from Indian heritage
- **Multi-language Support** - Interface in 22 Indian languages
- **Patriotism Scoring** - Rewards for community contribution

### 💰 Smart Features
- **Auto Pool Selection** - Automatically finds the best liquidity pool
- **Fee Transparency** - Clear breakdown of all fees
- **Priority Options** - Standard, Fast, or Instant processing
- **Receipt Generation** - Digital receipts with QR codes

## Getting Started

### Prerequisites
- Node.js 16+ and npm
- DeshChain wallet (Keplr compatible)
- Some NAMO tokens for transactions

### Installation

```bash
# Clone the repository
git clone https://github.com/deshchain/money-order-simple.git
cd money-order-simple

# Install dependencies
npm install

# Run development server
npm run dev

# Open http://localhost:3001
```

### Configuration

Create a `.env.local` file:

```env
NEXT_PUBLIC_API_URL=https://api.deshchain.org/v1
NEXT_PUBLIC_CHAIN_ID=deshchain-1
NEXT_PUBLIC_RPC_URL=https://rpc.deshchain.org
NEXT_PUBLIC_ANALYTICS_ENABLED=true
```

## Usage Guide

### 1. Connect Your Wallet
Click the wallet icon to connect your DeshChain-compatible wallet (Keplr).

### 2. Enter Transaction Details
- **From**: Your wallet address (auto-filled when connected)
- **To**: Recipient's DeshChain address (use address book for saved contacts)
- **Amount**: Enter NAMO amount or use presets

### 3. Review and Send
- Check the transaction preview
- Verify fees and exchange rate
- Click "Send Money Order"
- Receive digital receipt

## Features in Detail

### Address Book
Save frequently used addresses with custom names and tags:
- Family members
- Business partners
- Community pools
- Regular merchants

### Festival Bonuses
Automatic bonuses during Indian festivals:
- **Diwali**: 15% bonus on transactions
- **Independence Day**: 20% patriotic bonus
- **Holi**: 10% community bonus
- **Eid**: 12% celebration bonus

### Cultural Quotes
Each transaction includes an inspirational quote from:
- Ancient Indian texts
- Freedom fighters
- Cultural leaders
- Traditional wisdom

### Multi-Language Interface
Full support for:
- Hindi (हिन्दी)
- Bengali (বাংলা)
- Telugu (తెలుగు)
- Tamil (தமிழ்)
- And 18 more Indian languages

## Components

### SimpleMoneyOrderForm
Main form component with:
- Input validation
- Real-time quote updates
- Cultural integration
- Error handling

### AddressBook
Contact management with:
- Search functionality
- Favorite contacts
- Tag system
- Edit/Delete options

### FeeDisplay
Transparent fee breakdown:
- Base network fee
- Priority multiplier
- Festival discounts
- Total cost in NAMO and INR

### TransactionPreview
Visual confirmation showing:
- Sender/Receiver details
- Amount and conversion
- Fees and total
- Memo (if provided)

## Keyboard Shortcuts

- `Ctrl/Cmd + Enter` - Send transaction
- `Ctrl/Cmd + K` - Open address book
- `Ctrl/Cmd + L` - Change language
- `Esc` - Close dialogs

## Security Features

- **Address Validation** - Prevents sending to invalid addresses
- **Amount Limits** - Min: 1 NAMO, Max: 10M NAMO
- **Slippage Protection** - Maximum 5% price impact
- **Secure Storage** - Address book encrypted locally

## Performance

- **Fast Loading** - < 2s initial load
- **Optimized Bundle** - < 200KB gzipped
- **PWA Support** - Works offline for viewing
- **Mobile First** - Responsive design

## Accessibility

- **WCAG 2.1 AA** compliant
- **Keyboard Navigation** - Full keyboard support
- **Screen Reader** - Optimized for NVDA/JAWS
- **High Contrast** - Support for visual impairments

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request

### Development Guidelines

- Follow TypeScript best practices
- Test on multiple screen sizes
- Verify cultural accuracy
- Check accessibility compliance

## Troubleshooting

### Common Issues

**Wallet not connecting**
- Ensure Keplr is installed
- Check if DeshChain network is added
- Try refreshing the page

**Transaction failing**
- Verify sufficient NAMO balance
- Check network connectivity
- Ensure valid recipient address

**Language not displaying correctly**
- Install required fonts
- Clear browser cache
- Check language support

## Support

- **Documentation**: [docs.deshchain.org](https://docs.deshchain.org)
- **Discord**: [discord.gg/deshchain](https://discord.gg/deshchain)
- **Email**: support@deshchain.org
- **GitHub Issues**: [Report bugs](https://github.com/deshchain/money-order-simple/issues)

## License

Apache License 2.0 - See [LICENSE](./LICENSE) for details.

---

**Built with ❤️ for India's digital future** 🇮🇳

*Simple, Cultural, Secure - Money Orders for Everyone*