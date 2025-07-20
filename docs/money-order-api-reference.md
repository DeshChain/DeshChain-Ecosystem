# DeshChain Money Order DEX - API Reference

## Overview

The DeshChain Money Order API provides comprehensive access to all money order functionality, including P2P transfers, bulk processing, cross-chain operations, and analytics. This reference covers REST API endpoints, GraphQL schema, gRPC services, and SDK integration.

## Base Information

- **Base URL**: `https://api.deshchain.org/v1`
- **Protocol**: HTTPS only
- **Authentication**: API Key + JWT tokens
- **Rate Limiting**: 1000 requests/minute (higher limits available)
- **Response Format**: JSON
- **Error Format**: RFC 7807 Problem Details

---

## Authentication

### API Key Authentication

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  https://api.deshchain.org/v1/money-orders
```

### JWT Token Authentication

```javascript
// Get JWT token
const response = await fetch('/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    address: 'desh1...',
    signature: '0x...',
    message: 'Login request'
  })
});

const { token } = await response.json();

// Use token in subsequent requests
const ordersResponse = await fetch('/money-orders', {
  headers: { 'Authorization': `Bearer ${token}` }
});
```

---

## Core Money Order API

### Create Money Order

**Endpoint**: `POST /money-orders`

```bash
curl -X POST https://api.deshchain.org/v1/money-orders \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "recipientAddress": "desh1abc123...",
    "amount": "5000",
    "memo": "Payment for services",
    "priority": "NORMAL",
    "processingMethod": "P2P",
    "webhookUrl": "https://myapp.com/webhooks"
  }'
```

**Request Body**:
```typescript
interface CreateMoneyOrderRequest {
  recipientAddress: string;      // DeshChain address
  amount: string;                // Amount in smallest unit
  memo?: string;                 // Optional description
  priority?: "LOW" | "NORMAL" | "HIGH" | "URGENT";
  processingMethod?: "INSTANT" | "P2P" | "ESCROW";
  webhookUrl?: string;           // Webhook notification URL
  metadata?: Record<string, any>; // Additional metadata
}
```

**Response**:
```json
{
  "orderID": "ORDER_1234567890",
  "status": "CREATED",
  "senderAddress": "desh1sender...",
  "recipientAddress": "desh1recipient...",
  "amount": "5000",
  "fees": {
    "baseFee": "10",
    "processingFee": "5",
    "totalFee": "15"
  },
  "estimatedCompletion": "2024-01-15T10:45:00Z",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

### Get Money Order

**Endpoint**: `GET /money-orders/{orderID}`

```bash
curl https://api.deshchain.org/v1/money-orders/ORDER_1234567890 \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response**:
```json
{
  "orderID": "ORDER_1234567890",
  "status": "COMPLETED",
  "senderAddress": "desh1sender...",
  "recipientAddress": "desh1recipient...",
  "amount": "5000",
  "memo": "Payment for services",
  "priority": "NORMAL",
  "processingMethod": "P2P",
  "fees": {
    "baseFee": "10",
    "processingFee": "5",
    "totalFee": "15"
  },
  "sevaMitra": {
    "mitraId": "SM_123",
    "name": "Ram Kumar",
    "rating": 4.8,
    "location": "110001"
  },
  "timeline": [
    {
      "status": "CREATED",
      "timestamp": "2024-01-15T10:30:00Z"
    },
    {
      "status": "MATCHED",
      "timestamp": "2024-01-15T10:32:00Z"
    },
    {
      "status": "COMPLETED",
      "timestamp": "2024-01-15T10:45:00Z"
    }
  ],
  "createdAt": "2024-01-15T10:30:00Z",
  "completedAt": "2024-01-15T10:45:00Z"
}
```

### List Money Orders

**Endpoint**: `GET /money-orders`

```bash
curl "https://api.deshchain.org/v1/money-orders?status=COMPLETED&limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Query Parameters**:
- `status`: Filter by status (`CREATED`, `MATCHED`, `IN_TRANSIT`, `COMPLETED`, etc.)
- `senderAddress`: Filter by sender address
- `recipientAddress`: Filter by recipient address
- `startDate`: Filter by creation date (ISO 8601)
- `endDate`: Filter by creation date (ISO 8601)
- `minAmount`: Minimum amount filter
- `maxAmount`: Maximum amount filter
- `limit`: Number of results (max 100)
- `offset`: Pagination offset

### Cancel Money Order

**Endpoint**: `DELETE /money-orders/{orderID}`

```bash
curl -X DELETE https://api.deshchain.org/v1/money-orders/ORDER_1234567890 \
  -H "Authorization: Bearer YOUR_API_KEY"
```

---

## Bulk Orders API

### Create Bulk Order

**Endpoint**: `POST /bulk-orders`

```bash
curl -X POST https://api.deshchain.org/v1/bulk-orders \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "orders": [
      {
        "recipientAddress": "desh1emp001...",
        "amount": "45000",
        "memo": "January Salary",
        "customerRef": "EMP001"
      },
      {
        "recipientAddress": "desh1emp002...",
        "amount": "38000",
        "memo": "January Salary",
        "customerRef": "EMP002"
      }
    ],
    "metadata": {
      "description": "Monthly salary payments",
      "reference": "SAL-2024-01",
      "department": "HR"
    },
    "settings": {
      "batchSize": 50,
      "maxRetries": 3,
      "notifyOnCompletion": true
    }
  }'
```

**Request Body**:
```typescript
interface CreateBulkOrderRequest {
  orders: BulkOrderItem[];
  metadata: BulkOrderMetadata;
  settings: ProcessingSettings;
}

interface BulkOrderItem {
  recipientAddress: string;
  amount: string;
  memo?: string;
  priority?: string;
  customerRef?: string;
}

interface BulkOrderMetadata {
  description: string;
  reference: string;
  department?: string;
  projectCode?: string;
  notifyEmail?: string;
}

interface ProcessingSettings {
  batchSize: number;
  maxRetries: number;
  stopOnFirstFailure?: boolean;
  validateRecipients?: boolean;
  notifyOnCompletion?: boolean;
}
```

### Get Bulk Order Status

**Endpoint**: `GET /bulk-orders/{bulkOrderID}`

```bash
curl https://api.deshchain.org/v1/bulk-orders/BULK_1234567890 \
  -H "Authorization: Bearer YOUR_API_KEY"
```

### Upload CSV for Bulk Order

**Endpoint**: `POST /bulk-orders/upload`

```bash
curl -X POST https://api.deshchain.org/v1/bulk-orders/upload \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -F "file=@salary_payments.csv" \
  -F "metadata={\"description\":\"Salary payments\",\"reference\":\"SAL-2024-01\"}"
```

---

## Seva Mitra API

### Register as Seva Mitra

**Endpoint**: `POST /seva-mitra/register`

```bash
curl -X POST https://api.deshchain.org/v1/seva-mitra/register \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ram Kumar",
    "location": "110001",
    "services": ["CASH_IN", "CASH_OUT", "DOOR_DELIVERY"],
    "operatingHours": {
      "weekdays": {"start": "09:00", "end": "18:00"},
      "weekends": {"start": "10:00", "end": "16:00"}
    },
    "maxTransactionLimit": 25000,
    "commissionRates": {
      "cashIn": 0.8,
      "cashOut": 0.8,
      "delivery": 50
    }
  }'
```

### Get Seva Mitra Dashboard

**Endpoint**: `GET /seva-mitra/dashboard`

```bash
curl https://api.deshchain.org/v1/seva-mitra/dashboard \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response**:
```json
{
  "mitraInfo": {
    "mitraId": "SM_123",
    "name": "Ram Kumar",
    "rating": 4.8,
    "trustScore": 85,
    "badge": "GOLD"
  },
  "summary": {
    "todayEarnings": "450",
    "weeklyEarnings": "2800",
    "monthlyEarnings": "12500",
    "totalOrders": 156,
    "successRate": 98.7
  },
  "serviceRequests": [
    {
      "requestId": "REQ_001",
      "orderID": "ORDER_1234567890",
      "type": "CASH_OUT",
      "amount": "5000",
      "customerLocation": "110001",
      "urgency": "NORMAL",
      "estimatedCommission": "40"
    }
  ],
  "performanceStats": {
    "averageResponseTime": "3.2 minutes",
    "completionRate": 98.7,
    "customerRating": 4.8,
    "totalEarnings": "45600"
  }
}
```

### Accept Service Request

**Endpoint**: `POST /seva-mitra/requests/{requestID}/accept`

```bash
curl -X POST https://api.deshchain.org/v1/seva-mitra/requests/REQ_001/accept \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -d '{"estimatedArrival": "15 minutes"}'
```

---

## Cross-Chain API

### Get Supported Chains

**Endpoint**: `GET /cross-chain/supported-chains`

```bash
curl https://api.deshchain.org/v1/cross-chain/supported-chains \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response**:
```json
{
  "chains": [
    {
      "chainId": "cosmoshub-4",
      "chainName": "Cosmos Hub",
      "isActive": true,
      "channelId": "channel-0",
      "minAmount": "100",
      "maxAmount": "1000000",
      "fee": "0.002",
      "estimatedTime": "5-10 minutes"
    },
    {
      "chainId": "osmosis-1",
      "chainName": "Osmosis",
      "isActive": true,
      "channelId": "channel-1",
      "minAmount": "100",
      "maxAmount": "1000000",
      "fee": "0.002",
      "estimatedTime": "3-8 minutes"
    }
  ]
}
```

### Create Cross-Chain Money Order

**Endpoint**: `POST /cross-chain/money-orders`

```bash
curl -X POST https://api.deshchain.org/v1/cross-chain/money-orders \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "recipientAddress": "osmo1abc123...",
    "amount": "5000",
    "recipientChain": "osmosis-1",
    "memo": "Cross-chain payment",
    "timeoutMinutes": 60
  }'
```

### Get Cross-Chain Order Status

**Endpoint**: `GET /cross-chain/money-orders/{orderID}`

```bash
curl https://api.deshchain.org/v1/cross-chain/money-orders/XCHAIN_1234567890 \
  -H "Authorization: Bearer YOUR_API_KEY"
```

---

## Analytics API

### Get Dashboard Metrics

**Endpoint**: `GET /analytics/dashboard`

```bash
curl "https://api.deshchain.org/v1/analytics/dashboard?timeRange=7d" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response**:
```json
{
  "timeRange": "7d",
  "generatedAt": "2024-01-15T10:30:00Z",
  "transactionStats": {
    "totalCount": 1250,
    "successCount": 1238,
    "failedCount": 12,
    "successRate": 99.04
  },
  "volumeMetrics": {
    "totalVolume": "5675000",
    "averageAmount": "4540",
    "largestTransaction": "150000",
    "growth": 15.3
  },
  "performanceMetrics": {
    "averageProcessingTime": "4.2 minutes",
    "systemLoad": 45.6,
    "errorRate": 0.96
  }
}
```

### Generate Report

**Endpoint**: `POST /analytics/reports`

```bash
curl -X POST https://api.deshchain.org/v1/analytics/reports \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "reportType": "BUSINESS",
    "startDate": "2024-01-01",
    "endDate": "2024-01-31",
    "format": "PDF",
    "includeCharts": true
  }'
```

### Export Data

**Endpoint**: `POST /analytics/export`

```bash
curl -X POST https://api.deshchain.org/v1/analytics/export \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "query": {
      "startDate": "2024-01-01",
      "endDate": "2024-01-31",
      "filters": {
        "status": "COMPLETED",
        "minAmount": "1000"
      }
    },
    "format": "CSV",
    "compression": true
  }'
```

---

## Biometric Authentication API

### Enroll Biometric Template

**Endpoint**: `POST /biometric/enroll`

```bash
curl -X POST https://api.deshchain.org/v1/biometric/enroll \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "biometricType": "FINGERPRINT",
    "templateData": "base64_encoded_template",
    "deviceId": "device_unique_id"
  }'
```

### Authenticate with Biometric

**Endpoint**: `POST /biometric/authenticate`

```bash
curl -X POST https://api.deshchain.org/v1/biometric/authenticate \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "biometricType": "FINGERPRINT",
    "templateData": "base64_encoded_template",
    "deviceId": "device_unique_id"
  }'
```

---

## Webhooks

### Webhook Events

DeshChain sends webhooks for the following events:

- `money_order.created`
- `money_order.matched`
- `money_order.in_transit`
- `money_order.completed`
- `money_order.failed`
- `money_order.cancelled`
- `bulk_order.created`
- `bulk_order.batch_completed`
- `bulk_order.completed`
- `cross_chain.sent`
- `cross_chain.received`
- `cross_chain.completed`
- `seva_mitra.request_received`
- `biometric.authentication_failed`

### Webhook Payload

```json
{
  "id": "wh_1234567890",
  "event": "money_order.completed",
  "data": {
    "orderID": "ORDER_1234567890",
    "status": "COMPLETED",
    "senderAddress": "desh1sender...",
    "recipientAddress": "desh1recipient...",
    "amount": "5000",
    "completedAt": "2024-01-15T10:45:00Z"
  },
  "timestamp": "2024-01-15T10:45:05Z",
  "signature": "sha256=abcdef123456..."
}
```

### Webhook Verification

```javascript
const crypto = require('crypto');

function verifyWebhook(payload, signature, secret) {
  const expectedSignature = crypto
    .createHmac('sha256', secret)
    .update(payload)
    .digest('hex');
  
  return signature === `sha256=${expectedSignature}`;
}
```

---

## GraphQL API

### Schema Overview

```graphql
type Query {
  moneyOrder(id: ID!): MoneyOrder
  moneyOrders(filter: MoneyOrderFilter, pagination: Pagination): MoneyOrderConnection
  bulkOrder(id: ID!): BulkOrder
  sevaMitra(id: ID!): SevaMitra
  analytics(timeRange: String!): Analytics
}

type Mutation {
  createMoneyOrder(input: CreateMoneyOrderInput!): MoneyOrder
  cancelMoneyOrder(id: ID!): MoneyOrder
  createBulkOrder(input: CreateBulkOrderInput!): BulkOrder
  acceptServiceRequest(requestId: ID!): ServiceRequest
}

type Subscription {
  moneyOrderUpdates(orderId: ID!): MoneyOrder
  sevaMitraRequests(mitraId: ID!): ServiceRequest
}
```

### Example Queries

```graphql
# Get money order with timeline
query GetMoneyOrder($id: ID!) {
  moneyOrder(id: $id) {
    orderID
    status
    amount
    recipientAddress
    timeline {
      status
      timestamp
    }
    sevaMitra {
      name
      rating
    }
  }
}

# Create money order
mutation CreateMoneyOrder($input: CreateMoneyOrderInput!) {
  createMoneyOrder(input: $input) {
    orderID
    status
    estimatedCompletion
  }
}

# Subscribe to order updates
subscription MoneyOrderUpdates($orderId: ID!) {
  moneyOrderUpdates(orderId: $orderId) {
    orderID
    status
    progress
  }
}
```

---

## gRPC API

### Proto Definitions

```protobuf
service MoneyOrderService {
  rpc CreateMoneyOrder(CreateMoneyOrderRequest) returns (MoneyOrderResponse);
  rpc GetMoneyOrder(GetMoneyOrderRequest) returns (MoneyOrderResponse);
  rpc ListMoneyOrders(ListMoneyOrdersRequest) returns (ListMoneyOrdersResponse);
  rpc CancelMoneyOrder(CancelMoneyOrderRequest) returns (MoneyOrderResponse);
}

message CreateMoneyOrderRequest {
  string recipient_address = 1;
  string amount = 2;
  string memo = 3;
  Priority priority = 4;
  ProcessingMethod processing_method = 5;
}

message MoneyOrderResponse {
  string order_id = 1;
  Status status = 2;
  string sender_address = 3;
  string recipient_address = 4;
  string amount = 5;
  Fees fees = 6;
  google.protobuf.Timestamp created_at = 7;
}

enum Status {
  CREATED = 0;
  MATCHED = 1;
  IN_TRANSIT = 2;
  COMPLETED = 3;
  FAILED = 4;
  CANCELLED = 5;
}
```

### gRPC Client Example

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    pb "github.com/deshchain/api/proto/moneyorder"
)

func main() {
    conn, err := grpc.Dial("api.deshchain.org:443", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    client := pb.NewMoneyOrderServiceClient(conn)
    
    resp, err := client.CreateMoneyOrder(context.Background(), &pb.CreateMoneyOrderRequest{
        RecipientAddress: "desh1abc123...",
        Amount: "5000",
        Memo: "Payment for services",
        Priority: pb.Priority_NORMAL,
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Order created: %s", resp.OrderId)
}
```

---

## SDKs and Libraries

### JavaScript/TypeScript SDK

```bash
npm install @deshchain/sdk
```

```typescript
import { DeshChainSDK } from '@deshchain/sdk';

const sdk = new DeshChainSDK({
  apiKey: process.env.DESHCHAIN_API_KEY,
  network: 'mainnet'
});

// Create money order
const order = await sdk.moneyOrders.create({
  recipientAddress: 'desh1abc123...',
  amount: '5000',
  memo: 'Payment for services'
});

// Track order status
const status = await sdk.moneyOrders.getStatus(order.orderID);
```

### Python SDK

```bash
pip install deshchain-sdk
```

```python
from deshchain import DeshChainClient

client = DeshChainClient(
    api_key=os.environ['DESHCHAIN_API_KEY'],
    network='mainnet'
)

# Create money order
order = client.money_orders.create(
    recipient_address='desh1abc123...',
    amount='5000',
    memo='Payment for services'
)

# Get order status
status = client.money_orders.get_status(order.order_id)
```

### Go SDK

```bash
go get github.com/deshchain/go-sdk
```

```go
package main

import (
    "github.com/deshchain/go-sdk"
)

func main() {
    client := deshchain.NewClient(deshchain.Config{
        APIKey: os.Getenv("DESHCHAIN_API_KEY"),
        Network: "mainnet",
    })
    
    order, err := client.MoneyOrders.Create(&deshchain.CreateMoneyOrderRequest{
        RecipientAddress: "desh1abc123...",
        Amount: "5000",
        Memo: "Payment for services",
    })
}
```

---

## Error Handling

### Error Response Format

```json
{
  "type": "https://api.deshchain.org/errors/insufficient-balance",
  "title": "Insufficient Balance",
  "status": 400,
  "detail": "Account balance (1000 NAMO) is insufficient for transaction (5000 NAMO + 15 fees)",
  "instance": "/money-orders/ORDER_1234567890",
  "timestamp": "2024-01-15T10:30:00Z",
  "traceId": "abc123def456"
}
```

### Common Error Codes

| Code | Description | Resolution |
|------|-------------|------------|
| 400 | Bad Request | Check request parameters |
| 401 | Unauthorized | Verify API key |
| 403 | Forbidden | Check permissions |
| 404 | Not Found | Verify resource ID |
| 409 | Conflict | Resource already exists |
| 429 | Rate Limited | Reduce request frequency |
| 500 | Internal Error | Contact support |

### Error Handling Best Practices

```typescript
async function handleAPICall() {
  try {
    const response = await sdk.moneyOrders.create(orderData);
    return response;
  } catch (error) {
    if (error.status === 400) {
      // Handle validation errors
      console.error('Validation error:', error.detail);
    } else if (error.status === 429) {
      // Handle rate limiting
      await delay(error.retryAfter * 1000);
      return handleAPICall(); // Retry
    } else {
      // Handle other errors
      throw error;
    }
  }
}
```

---

## Rate Limiting

### Rate Limit Headers

```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1642261200
```

### Rate Limit Tiers

| Tier | Requests/Minute | Features |
|------|----------------|----------|
| Free | 100 | Basic API access |
| Developer | 1,000 | Enhanced limits |
| Business | 10,000 | Priority support |
| Enterprise | Unlimited | Custom limits |

---

## Testing

### Testnet Environment

- **Base URL**: `https://testnet-api.deshchain.org/v1`
- **Chain ID**: `deshchain-testnet-1`
- **Faucet**: `https://faucet.deshchain.org`

### Test Data

```json
{
  "testAccounts": {
    "sender": "deshtest1sender123...",
    "recipient": "deshtest1recipient456...",
    "sevaMitra": "deshtest1mitra789..."
  },
  "testTokens": "1000000namo"
}
```

### Integration Testing

```bash
# Run integration tests
npm test:integration

# Test specific endpoints
npm test:integration -- --grep "money-orders"
```

---

## Support and Resources

### Documentation
- API Reference: [api-docs.deshchain.org](https://api-docs.deshchain.org)
- Guides: [guides.deshchain.org](https://guides.deshchain.org)
- Examples: [examples.deshchain.org](https://examples.deshchain.org)

### Community
- Discord: [discord.gg/deshchain](https://discord.gg/deshchain)
- Telegram: [@DeshChainDev](https://t.me/DeshChainDev)
- Forum: [forum.deshchain.org](https://forum.deshchain.org)

### Support
- Email: api-support@deshchain.org
- Response Time: 24 hours (business days)
- Emergency: [status.deshchain.org](https://status.deshchain.org)

---

**API Version**: v1.0  
**Last Updated**: January 15, 2024  
**Status**: Production Ready 🚀