# DeshChain Money Order DEX - Complete User Guide

## Table of Contents
1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [P2P Money Orders](#p2p-money-orders)
4. [Seva Mitra System](#seva-mitra-system)
5. [Business Bulk Orders](#business-bulk-orders)
6. [Cross-Chain Transfers](#cross-chain-transfers)
7. [Biometric Authentication](#biometric-authentication)
8. [Analytics & Reporting](#analytics--reporting)
9. [Mobile App Guide](#mobile-app-guide)
10. [Troubleshooting](#troubleshooting)

---

## Introduction

The DeshChain Money Order DEX revolutionizes traditional money transfers by combining the familiarity of postal money orders with blockchain technology. Our platform enables secure, transparent, and culturally-rooted financial transactions with enhanced features for the digital age.

### Key Features
- **Traditional Money Order Experience**: Familiar interface inspired by Indian postal money orders
- **P2P Network**: Direct peer-to-peer matching without intermediaries
- **Seva Mitra System**: Community-driven cash-in/cash-out network
- **Bulk Processing**: Enterprise-grade bulk order capabilities
- **Cross-Chain Support**: Transfer across different blockchains via IBC
- **Biometric Security**: Multi-factor authentication using biometric data
- **Cultural Integration**: Support for 22 Indian languages and cultural themes

---

## Getting Started

### Account Registration

1. **Create Your Account**
   ```bash
   # Using DeshChain CLI
   deshchaind tx auth create-account --from mykey
   ```

2. **Complete KYC Verification**
   - Provide Aadhaar number
   - Upload government ID documents
   - Complete biometric enrollment (fingerprint, face, voice)
   - Village Panchayat verification (optional for enhanced trust)

3. **Fund Your Account**
   ```bash
   # Deposit NAMO tokens
   deshchaind tx bank send source_address your_address 1000namo
   ```

### Account Types

#### Individual Account
- Basic money order functionality
- Daily limit: ₹50,000
- Monthly limit: ₹10,00,000

#### Business Account
- Bulk order capabilities
- Higher transaction limits
- Advanced analytics and reporting
- Multi-user access controls

#### Seva Mitra Account
- Community service provider
- Cash-in/cash-out services
- Commission-based earnings
- Enhanced trust verification

---

## P2P Money Orders

### Creating a Money Order

1. **Access Money Order Interface**
   ```typescript
   // Navigate to Money Order section
   const orderData = {
     recipientAddress: "desh1abc123...",
     amount: "5000",
     memo: "Monthly salary",
     priority: "NORMAL"
   };
   ```

2. **Fill Order Details**
   - **Recipient Address**: DeshChain wallet address
   - **Amount**: Transaction amount in INR
   - **Memo**: Purpose of transfer (optional)
   - **Priority**: LOW, NORMAL, HIGH, URGENT

3. **Select Processing Method**
   - **Instant**: Direct wallet-to-wallet transfer
   - **P2P Matching**: Find local Seva Mitra for cash services
   - **Escrow**: Use smart contract for secure holding

### Order Status Tracking

```typescript
// Check order status
const status = await moneyOrderKeeper.GetOrder(orderID);

// Possible statuses:
// - CREATED: Order initialized
// - MATCHED: Paired with Seva Mitra
// - IN_TRANSIT: Being processed
// - COMPLETED: Successfully delivered
// - CANCELLED: Order cancelled
// - REFUNDED: Amount refunded
```

### Geographic Matching Algorithm

Our system matches orders based on postal code proximity:

```go
func (pmm *P2PMatchingManager) FindNearbySevaMap(
    postalCode string,
    radius int,
) ([]SevaMitra, error) {
    // Find Seva Mitras within specified radius
    // Priority matching based on:
    // 1. Distance (postal code proximity)
    // 2. Trust score
    // 3. Available capacity
    // 4. Commission rates
}
```

---

## Seva Mitra System

### Becoming a Seva Mitra

1. **Apply for Seva Mitra Status**
   ```bash
   deshchaind tx moneyorder register-seva-mitra \
     --name "Ram Kumar" \
     --location "110001" \
     --services "cash-in,cash-out" \
     --from mykey
   ```

2. **Complete Enhanced KYC**
   - Business registration (if applicable)
   - Police verification certificate
   - Bank account verification
   - Collateral deposit (₹50,000 minimum)

3. **Set Service Parameters**
   - Operating hours
   - Service radius (postal codes)
   - Commission rates
   - Maximum transaction limits

### Seva Mitra Dashboard

Access comprehensive dashboard with:

```typescript
interface SevaMitraDashboard {
  earnings: {
    today: number;
    thisWeek: number;
    thisMonth: number;
    total: number;
  };
  serviceRequests: ServiceRequest[];
  performanceMetrics: {
    successRate: number;
    averageResponseTime: number;
    customerRating: number;
    completedOrders: number;
  };
  trustScore: {
    current: number;
    factors: TrustFactor[];
    badges: Badge[];
  };
}
```

### Commission Structure

| Trust Level | Commission Rate | Benefits |
|-------------|----------------|----------|
| Bronze      | 0.5%           | Basic support |
| Silver      | 0.7%           | Priority matching |
| Gold        | 1.0%           | Enhanced visibility |
| Platinum    | 1.2%           | Premium features |
| Diamond     | 1.5%           | Exclusive partnerships |

---

## Business Bulk Orders

### Setting Up Business Account

1. **Business Registration**
   ```bash
   deshchaind tx moneyorder register-business \
     --business-name "Tech Solutions Pvt Ltd" \
     --registration-number "U72900DL2020PTC123456" \
     --gst-number "07AAACP1234A1ZN" \
     --from businesskey
   ```

2. **Configure Limits**
   - Daily transaction limit: ₹50,00,000
   - Monthly limit: ₹10,00,00,000
   - Maximum bulk order size: 10,000 orders

### Creating Bulk Orders

1. **Prepare CSV Template**
   ```csv
   recipient_address,amount,memo,priority,customer_ref
   desh1abc123...,1000,Salary payment,NORMAL,EMP001
   desh1def456...,2500,Vendor payment,HIGH,VND002
   desh1ghi789...,500,Refund,LOW,REF003
   ```

2. **Upload and Validate**
   ```typescript
   const bulkOrder = {
     orders: parsedOrders,
     metadata: {
       description: "Monthly salary payments",
       reference: "PAY-2024-001",
       department: "HR"
     },
     settings: {
       batchSize: 50,
       maxRetries: 3,
       validateRecipients: true
     }
   };
   ```

3. **Monitor Processing**
   - Real-time batch processing status
   - Individual order tracking
   - Error handling and retry logic
   - Completion notifications

### Bulk Order Analytics

```typescript
// Get bulk order analytics
const analytics = await analyticsManager.GetBusinessReport(
  businessAddress,
  startDate,
  endDate
);

// Includes:
// - Transaction volume and counts
// - Success/failure rates
// - Cost savings analysis
// - Processing time metrics
// - Compliance reporting
```

---

## Cross-Chain Transfers

### Supported Chains

Our IBC integration supports transfers to:
- Cosmos Hub
- Osmosis
- Juno Network
- Secret Network
- Akash Network

### Initiating Cross-Chain Transfer

```bash
# Send money order to another chain
deshchaind tx moneyorder send-cross-chain \
  --recipient-address "cosmos1abc123..." \
  --amount 1000namo \
  --recipient-chain "cosmoshub-4" \
  --memo "Cross-chain payment" \
  --timeout-minutes 60 \
  --from mykey
```

### Cross-Chain Status Tracking

```go
// Query cross-chain order status
type CrossChainStatus struct {
    OrderID          string
    Status           CrossChainStatus
    SenderChain      string
    RecipientChain   string
    ChannelID        string
    EstimatedTime    time.Duration
    CompletionRate   float64
}
```

### IBC Channel Management

```bash
# List supported chains
deshchaind query moneyorder supported-chains

# Check channel health
deshchaind query moneyorder channel-status --channel-id channel-0

# Get transfer fees
deshchaind query moneyorder cross-chain-fees --chain-id cosmoshub-4
```

---

## Biometric Authentication

### Supported Biometric Types

1. **Fingerprint Recognition**
   - Touch sensor or camera-based
   - Multiple finger enrollment
   - 85% match threshold

2. **Face Recognition**
   - 3D facial mapping
   - Liveness detection
   - Anti-spoofing measures

3. **Voice Recognition**
   - Voice pattern analysis
   - Passphrase verification
   - Background noise filtering

4. **Iris Scanning**
   - High-precision iris patterns
   - Infrared camera support
   - Medical condition considerations

5. **Palm Recognition**
   - Palm vein pattern analysis
   - Contact-free scanning
   - Hygiene-friendly option

### Enrollment Process

```typescript
// Enroll biometric template
const enrollmentResult = await biometricAuth.enrollTemplate({
  userId: userAddress,
  biometricType: "FINGERPRINT",
  templateData: capturedTemplate,
  deviceInfo: deviceIdentifier
});
```

### Authentication Flow

```go
func (bam *BiometricAuthManager) AuthenticateBiometric(
    ctx sdk.Context,
    userAddress string,
    biometricType BiometricType,
    templateData []byte,
    deviceID string,
) (bool, error) {
    // Verify device binding
    // Compare templates
    // Check attempt limits
    // Update security metrics
}
```

---

## Analytics & Reporting

### Dashboard Metrics

Access real-time insights including:

```typescript
interface DashboardMetrics {
  transactionStats: {
    totalCount: number;
    successRate: number;
    averageAmount: number;
  };
  volumeMetrics: {
    dailyVolume: number;
    weeklyGrowth: number;
    monthlyTrend: number;
  };
  performanceData: {
    averageProcessingTime: number;
    systemLoad: number;
    errorRate: number;
  };
}
```

### Report Generation

1. **System Reports**
   ```bash
   # Generate system analytics report
   deshchaind tx moneyorder generate-report \
     --report-type SYSTEM \
     --start-date 2024-01-01 \
     --end-date 2024-01-31 \
     --format PDF
   ```

2. **Business Reports**
   ```bash
   # Generate business analytics
   deshchaind tx moneyorder generate-business-report \
     --business-address desh1business... \
     --report-type TRANSACTION \
     --export-format EXCEL
   ```

### Anomaly Detection

Our AI-powered system monitors for:
- Unusual transaction patterns
- Potential fraud indicators
- System performance anomalies
- Compliance violations

```go
type DetectedAnomaly struct {
    Type        string
    Severity    AnomalySeverity  // LOW, MEDIUM, HIGH, CRITICAL
    Description string
    Confidence  float64
    Context     map[string]interface{}
}
```

---

## Mobile App Guide

### Installation

1. **Download Batua Wallet**
   - Android: Google Play Store
   - iOS: Apple App Store
   - Direct APK: [deshchain.org/batua](https://deshchain.org/batua)

2. **Initial Setup**
   ```dart
   // Initialize wallet
   final wallet = await BatuaWallet.initialize(
     networkConfig: DeshChainMainnet,
     biometricAuth: true,
     languageCode: 'hi', // Hindi
   );
   ```

### Key Mobile Features

#### Cultural Interface
- 22 Indian language support
- Festival-themed UI
- Cultural color schemes
- Regional customizations

#### Offline Capabilities
- Order creation without internet
- Queue management
- Sync when connected
- Emergency cash codes

#### Security Features
```dart
// Multi-layer security
final securityConfig = SecurityConfig(
  biometricAuth: true,
  deviceBinding: true,
  pinBackup: true,
  autoLock: Duration(minutes: 5),
);
```

### Mobile-Specific Functions

1. **QR Code Scanning**
   - Scan recipient addresses
   - Quick payment setup
   - Seva Mitra identification

2. **Voice Commands**
   - "Send money to Ram"
   - "Check my balance"
   - "Find nearby Seva Mitra"

3. **Offline Mode**
   - Create orders offline
   - Queue for processing
   - Emergency backup codes

---

## Troubleshooting

### Common Issues

#### Order Stuck in Pending
**Symptoms**: Order shows PENDING status for extended time
**Solutions**:
1. Check network connectivity
2. Verify sufficient balance
3. Contact assigned Seva Mitra
4. Cancel and recreate if necessary

```bash
# Check order details
deshchaind query moneyorder order --order-id ORDER123

# Cancel if needed
deshchaind tx moneyorder cancel-order --order-id ORDER123 --from mykey
```

#### Biometric Authentication Fails
**Symptoms**: Repeated authentication failures
**Solutions**:
1. Clean biometric sensor
2. Re-enroll templates
3. Use backup authentication
4. Check device compatibility

#### Cross-Chain Transfer Delays
**Symptoms**: Cross-chain order taking longer than expected
**Solutions**:
1. Check IBC channel status
2. Verify destination chain health
3. Monitor relayer activity
4. Contact support if timeout approaches

### Error Codes

| Code | Description | Solution |
|------|-------------|----------|
| E001 | Insufficient balance | Add funds to account |
| E002 | Invalid recipient | Verify address format |
| E003 | KYC incomplete | Complete verification |
| E004 | Daily limit exceeded | Wait for limit reset |
| E005 | Seva Mitra unavailable | Try different location |
| E006 | Biometric mismatch | Re-authenticate |
| E007 | Network congestion | Retry later |
| E008 | Cross-chain timeout | Check destination chain |

### Support Channels

1. **Documentation**: [docs.deshchain.org](https://docs.deshchain.org)
2. **Community Forum**: [forum.deshchain.org](https://forum.deshchain.org)
3. **Telegram**: @DeshChainSupport
4. **Email**: support@deshchain.org
5. **Emergency Hotline**: 1800-DESH-HELP

### Emergency Procedures

#### Lost Device
1. Immediately report to support
2. Use backup recovery phrase
3. Re-enroll biometrics
4. Update security settings

#### Compromised Account
1. Change passwords immediately
2. Revoke all active sessions
3. Review recent transactions
4. Contact fraud department

#### System Maintenance
- Check [status.deshchain.org](https://status.deshchain.org)
- Follow @DeshChainStatus for updates
- Use mobile app for urgent transactions
- Contact support for critical issues

---

## Advanced Features

### Smart Contract Integration

```go
// Custom money order logic
type SmartMoneyOrder struct {
    Conditions    []Condition
    AutoExecute   bool
    EscrowPeriod  time.Duration
    RefundPolicy  RefundPolicy
}
```

### API Integration

```javascript
// REST API example
const order = await fetch('/api/v1/moneyorder/create', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    recipientAddress: 'desh1...',
    amount: '5000',
    memo: 'API payment'
  })
});
```

### Webhook Integration

```json
{
  "event": "order.completed",
  "order_id": "ORDER123",
  "amount": "5000",
  "status": "COMPLETED",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## Security Best Practices

### Account Security
1. Use strong, unique passwords
2. Enable two-factor authentication
3. Regularly update biometric templates
4. Monitor account activity

### Transaction Security
1. Verify recipient addresses
2. Use memo fields appropriately
3. Set appropriate priority levels
4. Keep transaction records

### Device Security
1. Keep app updated
2. Use device lock screens
3. Avoid public Wi-Fi for transactions
4. Regular security scans

---

This comprehensive guide covers all aspects of the DeshChain Money Order DEX. For additional help or specific use cases, please refer to our detailed API documentation or contact our support team.

**Happy Banking with DeshChain! 🇮🇳**