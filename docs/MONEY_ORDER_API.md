# Money Order API Documentation

## Overview

The Money Order API provides comprehensive access to DeshChain's culturally-integrated financial exchange system. This RESTful API enables developers to build applications that leverage traditional money order functionality with modern DeFi features and Indian cultural heritage integration.

## Base URL

```
https://api.deshchain.org/v1/moneyorder
```

## Authentication

### API Key Authentication
```http
Authorization: Bearer YOUR_API_KEY
```

### Wallet Signature Authentication
```http
X-Wallet-Address: desh1...
X-Signature: 0x...
X-Timestamp: 1640995200
```

## Rate Limiting

- **Standard**: 1000 requests per hour
- **Premium**: 10,000 requests per hour
- **Enterprise**: Unlimited

Rate limit headers:
```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

## Common Headers

```http
Content-Type: application/json
Accept: application/json
X-Chain-ID: deshchain-1
X-Culture-Language: hi-IN
X-Festival-Context: diwali
```

## Error Handling

### Error Response Format
```json
{
  "error": {
    "code": "INSUFFICIENT_LIQUIDITY",
    "message": "Insufficient liquidity in the selected pool",
    "details": {
      "pool_id": "123",
      "required_amount": "1000000",
      "available_amount": "500000"
    },
    "cultural_message": "धैर्य रखें, अधिक तरलता आ रही है - Have patience, more liquidity is coming"
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "req_abc123def456"
}
```

### Error Codes

| Code | Description | HTTP Status |
|------|-------------|-------------|
| `INVALID_POOL` | Pool not found or invalid | 404 |
| `INSUFFICIENT_LIQUIDITY` | Not enough liquidity for transaction | 400 |
| `EXCESSIVE_SLIPPAGE` | Price impact exceeds maximum allowed | 400 |
| `RATE_LIMIT_EXCEEDED` | Too many requests | 429 |
| `INVALID_SIGNATURE` | Invalid wallet signature | 401 |
| `UNSUPPORTED_TOKEN` | Token not supported in pool | 400 |

---

## Pool Management

### Get All Pools

```http
GET /pools
```

**Query Parameters:**
- `type` (string, optional): Filter by pool type (`amm`, `fixed_rate`, `village`, `concentrated`)
- `tokens` (string, optional): Filter by token pair (e.g., `unamo,inr`)
- `village_code` (string, optional): Filter by postal code for village pools
- `cultural_theme` (string, optional): Filter by cultural theme
- `min_liquidity` (string, optional): Minimum liquidity amount
- `page` (integer, optional): Page number (default: 1)
- `limit` (integer, optional): Items per page (default: 20, max: 100)

**Response:**
```json
{
  "pools": [
    {
      "pool_id": "1",
      "type": "fixed_rate",
      "token_a": "unamo",
      "token_b": "inr",
      "exchange_rate": "0.075",
      "reserve_a": {
        "denom": "unamo",
        "amount": "1000000000"
      },
      "reserve_b": {
        "denom": "inr",
        "amount": "75000000"
      },
      "cultural_theme": "independence",
      "trust_score": "0.98",
      "price_stability": "0.99",
      "patriotism_quote": "स्वतंत्रता हमारा जन्मसिद्ध अधिकार है - Freedom is our birthright",
      "created_at": "2024-01-01T00:00:00Z",
      "status": "active"
    }
  ],
  "pagination": {
    "total": 150,
    "page": 1,
    "limit": 20,
    "total_pages": 8
  },
  "cultural_context": {
    "current_festival": "diwali",
    "festival_bonus": "0.15",
    "patriotism_day": false
  }
}
```

### Get Pool Details

```http
GET /pools/{pool_id}
```

**Path Parameters:**
- `pool_id` (string, required): Pool identifier

**Response:**
```json
{
  "pool": {
    "pool_id": "1",
    "type": "village",
    "village_name": "Gram Panchayat Prosperity",
    "postal_code": "110001",
    "token_a": "unamo",
    "token_b": "inr",
    "reserve_a": {
      "denom": "unamo",
      "amount": "5000000000"
    },
    "reserve_b": {
      "denom": "inr",
      "amount": "375000000"
    },
    "total_shares": "2236067977",
    "members": [
      "desh1abc...",
      "desh1def..."
    ],
    "coordinator": "desh1xyz...",
    "swap_fee": "0.0015",
    "cultural_theme": "community",
    "trust_score": "0.96",
    "community_impact": "0.88",
    "local_economy_boost": "0.92",
    "monthly_transactions": 1250,
    "patriotism_quote": "वसुधैव कुटुम्बकम् - The world is one family",
    "verification_status": "verified",
    "created_at": "2024-01-01T00:00:00Z",
    "last_activity": "2024-01-15T10:25:00Z"
  },
  "performance_metrics": {
    "24h_volume": "2500000",
    "7d_volume": "18750000",
    "30d_volume": "78125000",
    "apy": "0.18",
    "price_impact_24h": "0.002"
  },
  "cultural_features": {
    "festival_bonuses_active": true,
    "current_bonus": "0.10",
    "patriotism_rewards_enabled": true,
    "community_support_level": "high"
  }
}
```

### Create Pool

```http
POST /pools
```

**Request Body:**
```json
{
  "type": "village",
  "token_a": "unamo",
  "token_b": "inr",
  "initial_reserve_a": {
    "denom": "unamo",
    "amount": "1000000000"
  },
  "initial_reserve_b": {
    "denom": "inr",
    "amount": "75000000"
  },
  "village_name": "New Village Pool",
  "postal_code": "110002",
  "cultural_theme": "prosperity",
  "swap_fee": "0.0015"
}
```

**Response:**
```json
{
  "pool_id": "151",
  "transaction_hash": "0xabc123...",
  "cultural_quote": {
    "text": "सर्वे भवन्तु सुखिनः - May all beings be happy",
    "author": "Vedic Tradition",
    "category": "prosperity"
  },
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

## Money Order Operations

### Create Money Order

```http
POST /money-orders
```

**Request Body:**
```json
{
  "sender": "desh1abc...",
  "receiver": "desh1def...",
  "amount": {
    "denom": "unamo",
    "amount": "1000000"
  },
  "pool_id": "1",
  "memo": "Monthly allowance",
  "cultural_preferences": {
    "language": "hi-IN",
    "theme": "family",
    "include_quote": true
  },
  "priority": "standard"
}
```

**Response:**
```json
{
  "order_id": "MO-20240115-001",
  "receipt": {
    "receipt_id": "RCP-20240115-001",
    "sender": "desh1abc...",
    "receiver": "desh1def...",
    "amount": {
      "denom": "unamo",
      "amount": "1000000"
    },
    "fee": {
      "denom": "unamo",
      "amount": "1000"
    },
    "exchange_rate": "0.075",
    "cultural_quote": {
      "text": "माता पिता गुरु देवा - Mother, Father, Teacher, God",
      "author": "Traditional Sanskrit",
      "category": "family",
      "language": "sanskrit"
    },
    "patriotism_bonus": {
      "denom": "unamo",
      "amount": "500"
    },
    "verification_code": "VERIFY-ABC123",
    "qr_code": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
    "timestamp": "2024-01-15T10:30:00Z",
    "status": "completed"
  },
  "transaction_hash": "0xdef456...",
  "estimated_delivery": "2024-01-15T10:35:00Z"
}
```

### Get Money Order Status

```http
GET /money-orders/{order_id}
```

**Response:**
```json
{
  "order": {
    "order_id": "MO-20240115-001",
    "status": "completed",
    "sender": "desh1abc...",
    "receiver": "desh1def...",
    "amount": {
      "denom": "unamo",
      "amount": "1000000"
    },
    "created_at": "2024-01-15T10:30:00Z",
    "completed_at": "2024-01-15T10:30:05Z",
    "receipt_id": "RCP-20240115-001"
  },
  "tracking": {
    "status_history": [
      {
        "status": "created",
        "timestamp": "2024-01-15T10:30:00Z"
      },
      {
        "status": "processing",
        "timestamp": "2024-01-15T10:30:02Z"
      },
      {
        "status": "completed",
        "timestamp": "2024-01-15T10:30:05Z"
      }
    ]
  }
}
```

---

## Trading Operations

### Swap Tokens

```http
POST /swaps
```

**Request Body:**
```json
{
  "pool_id": "1",
  "token_in": {
    "denom": "unamo",
    "amount": "1000000"
  },
  "token_out_denom": "inr",
  "min_amount_out": "70000",
  "max_slippage": "0.05",
  "cultural_preferences": {
    "language": "hi-IN",
    "include_quote": true
  }
}
```

**Response:**
```json
{
  "swap_id": "SWAP-20240115-001",
  "token_in": {
    "denom": "unamo",
    "amount": "1000000"
  },
  "token_out": {
    "denom": "inr",
    "amount": "74250"
  },
  "exchange_rate": "0.07425",
  "price_impact": "0.01",
  "fee": {
    "denom": "unamo",
    "amount": "3000"
  },
  "cultural_bonus": {
    "denom": "inr",
    "amount": "750"
  },
  "cultural_quote": {
    "text": "यत्र नार्यस्तु पूज्यन्ते रमन्ते तत्र देवताः",
    "author": "Manusmriti",
    "category": "prosperity",
    "translation": "Where women are honored, divinity blossoms there"
  },
  "transaction_hash": "0x789abc...",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Get Swap Quote

```http
GET /swaps/quote
```

**Query Parameters:**
- `pool_id` (string, required): Pool identifier
- `token_in_denom` (string, required): Input token denomination
- `token_in_amount` (string, required): Input token amount
- `token_out_denom` (string, required): Output token denomination

**Response:**
```json
{
  "quote": {
    "token_in": {
      "denom": "unamo",
      "amount": "1000000"
    },
    "token_out": {
      "denom": "inr",
      "amount": "74625"
    },
    "exchange_rate": "0.074625",
    "price_impact": "0.005",
    "fee": {
      "denom": "unamo",
      "amount": "3000"
    },
    "minimum_received": "74250",
    "route": [
      {
        "pool_id": "1",
        "token_in": "unamo",
        "token_out": "inr"
      }
    ]
  },
  "cultural_context": {
    "festival_bonus_available": true,
    "bonus_percentage": "0.10",
    "cultural_theme": "prosperity"
  },
  "valid_until": "2024-01-15T10:35:00Z"
}
```

---

## Liquidity Management

### Add Liquidity

```http
POST /liquidity/add
```

**Request Body:**
```json
{
  "pool_id": "1",
  "token_a_amount": {
    "denom": "unamo",
    "amount": "1000000"
  },
  "token_b_amount": {
    "denom": "inr",
    "amount": "75000"
  },
  "min_shares": "86602540",
  "cultural_preferences": {
    "language": "hi-IN",
    "theme": "community"
  }
}
```

**Response:**
```json
{
  "position_id": "LP-20240115-001",
  "shares_minted": "86602540",
  "token_a_deposited": {
    "denom": "unamo",
    "amount": "1000000"
  },
  "token_b_deposited": {
    "denom": "inr",
    "amount": "75000"
  },
  "cultural_quote": {
    "text": "सहयोग से समृद्धि - Prosperity through cooperation",
    "author": "Traditional Wisdom",
    "category": "community"
  },
  "transaction_hash": "0x123def...",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Remove Liquidity

```http
POST /liquidity/remove
```

**Request Body:**
```json
{
  "position_id": "LP-20240115-001",
  "shares_to_remove": "43301270",
  "min_token_a": "450000",
  "min_token_b": "33750"
}
```

**Response:**
```json
{
  "token_a_received": {
    "denom": "unamo",
    "amount": "500000"
  },
  "token_b_received": {
    "denom": "inr",
    "amount": "37500"
  },
  "rewards_claimed": [
    {
      "denom": "unamo",
      "amount": "5000"
    }
  ],
  "cultural_bonus": {
    "denom": "unamo",
    "amount": "250"
  },
  "transaction_hash": "0x456ghi...",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Get Liquidity Positions

```http
GET /liquidity/positions
```

**Query Parameters:**
- `owner` (string, required): Owner address
- `pool_id` (string, optional): Filter by pool
- `status` (string, optional): Filter by status (`active`, `closed`)

**Response:**
```json
{
  "positions": [
    {
      "position_id": "LP-20240115-001",
      "pool_id": "1",
      "owner": "desh1abc...",
      "shares": "43301270",
      "token_a_amount": {
        "denom": "unamo",
        "amount": "500000"
      },
      "token_b_amount": {
        "denom": "inr",
        "amount": "37500"
      },
      "rewards_accumulated": [
        {
          "denom": "unamo",
          "amount": "5000"
        }
      ],
      "cultural_bonus": {
        "denom": "unamo",
        "amount": "250"
      },
      "patriotism_reward": {
        "denom": "unamo",
        "amount": "100"
      },
      "community_contribution": "0.75",
      "created_at": "2024-01-15T10:30:00Z",
      "last_claim_at": "2024-01-15T10:25:00Z",
      "status": "active"
    }
  ]
}
```

---

## Cultural Features

### Get Cultural Quotes

```http
GET /cultural/quotes
```

**Query Parameters:**
- `category` (string, optional): Filter by category (`wisdom`, `motivation`, `patriotism`, `prosperity`)
- `language` (string, optional): Filter by language (`hi`, `en`, `sa`, `bn`, etc.)
- `occasion` (string, optional): Filter by occasion (`general`, `festival`, `independence_day`)
- `author` (string, optional): Filter by author

**Response:**
```json
{
  "quotes": [
    {
      "quote_id": "1",
      "text": "सत्यमेव जयते - Truth alone triumphs",
      "author": "Mundaka Upanishad",
      "category": "truth",
      "language": "sanskrit",
      "occasion": "general",
      "translation": "Truth alone triumphs",
      "context": "National motto of India",
      "weight": 10,
      "active": true
    }
  ],
  "total": 10000,
  "cultural_context": {
    "current_festival": "diwali",
    "recommended_categories": ["prosperity", "celebration", "community"]
  }
}
```

### Get Festival Information

```http
GET /cultural/festivals
```

**Query Parameters:**
- `current` (boolean, optional): Show only currently active festivals
- `region` (string, optional): Filter by region (`north_india`, `south_india`, `pan_india`)
- `upcoming` (boolean, optional): Show upcoming festivals

**Response:**
```json
{
  "festivals": [
    {
      "festival_id": "1",
      "name": "Diwali",
      "description": "Festival of Lights",
      "start_date": "2024-11-01",
      "end_date": "2024-11-05",
      "bonus_rate": "0.15",
      "cultural_theme": "prosperity",
      "region": "pan_india",
      "significance": "Victory of light over darkness",
      "traditional_greeting": "दीपावली की शुभकामनाएं",
      "active": true,
      "days_remaining": 290
    }
  ],
  "current_active": [
    {
      "name": "Republic Day",
      "bonus_rate": "0.20",
      "theme": "patriotism"
    }
  ]
}
```

### Update Cultural Preferences

```http
PUT /cultural/preferences
```

**Request Body:**
```json
{
  "language": "hi-IN",
  "preferred_themes": ["family", "prosperity", "community"],
  "quote_frequency": "every_transaction",
  "festival_notifications": true,
  "patriotism_features": true,
  "regional_focus": "north_india"
}
```

**Response:**
```json
{
  "preferences_updated": true,
  "effective_from": "2024-01-15T10:30:00Z",
  "cultural_quote": {
    "text": "यदा यदा हि धर्मस्य ग्लानिर्भवति भारत।",
    "author": "Bhagavad Gita",
    "category": "wisdom",
    "message": "Preferences updated successfully"
  }
}
```

---

## Analytics & Reporting

### Pool Performance

```http
GET /analytics/pools/{pool_id}/performance
```

**Query Parameters:**
- `period` (string, optional): Time period (`24h`, `7d`, `30d`, `1y`)
- `metrics` (string, optional): Comma-separated list of metrics

**Response:**
```json
{
  "pool_id": "1",
  "period": "7d",
  "performance": {
    "volume": "125000000",
    "transactions": 2500,
    "unique_users": 850,
    "average_transaction": "50000",
    "apy": "0.18",
    "price_stability": "0.99",
    "cultural_engagement": "0.85"
  },
  "cultural_metrics": {
    "quotes_displayed": 2500,
    "festival_bonuses_distributed": "375000",
    "patriotism_rewards": "125000",
    "community_impact_score": "0.88"
  },
  "trend_data": [
    {
      "date": "2024-01-08",
      "volume": "15000000",
      "transactions": 300
    }
  ]
}
```

### Community Statistics

```http
GET /analytics/community
```

**Response:**
```json
{
  "community_stats": {
    "total_pools": 150,
    "village_pools": 45,
    "total_users": 25000,
    "active_users_24h": 1250,
    "total_volume": "5000000000",
    "cultural_engagement_rate": "0.78"
  },
  "cultural_impact": {
    "quotes_served_today": 3500,
    "festivals_celebrated": 12,
    "patriotism_score_average": "0.82",
    "community_pools_growth": "0.15"
  },
  "regional_breakdown": {
    "north_india": {
      "pools": 65,
      "volume": "2000000000"
    },
    "south_india": {
      "pools": 45,
      "volume": "1500000000"
    },
    "west_india": {
      "pools": 25,
      "volume": "1000000000"
    },
    "east_india": {
      "pools": 15,
      "volume": "500000000"
    }
  }
}
```

---

## Receipts & Documentation

### Get Receipt

```http
GET /receipts/{receipt_id}
```

**Response:**
```json
{
  "receipt": {
    "receipt_id": "RCP-20240115-001",
    "order_id": "MO-20240115-001",
    "sender": "desh1abc...",
    "receiver": "desh1def...",
    "amount": {
      "denom": "unamo",
      "amount": "1000000"
    },
    "fee": {
      "denom": "unamo",
      "amount": "1000"
    },
    "exchange_rate": "0.075",
    "cultural_quote": {
      "text": "माता पिता गुरु देवा - Mother, Father, Teacher, God",
      "author": "Traditional Sanskrit",
      "category": "family"
    },
    "patriotism_bonus": {
      "denom": "unamo",
      "amount": "500"
    },
    "verification_code": "VERIFY-ABC123",
    "qr_code": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
    "digital_signature": "0x789def...",
    "timestamp": "2024-01-15T10:30:00Z",
    "status": "completed",
    "blockchain_confirmation": {
      "transaction_hash": "0xdef456...",
      "block_height": 12345678,
      "confirmations": 15
    }
  },
  "print_ready_url": "https://api.deshchain.org/v1/receipts/RCP-20240115-001/print",
  "share_url": "https://deshchain.org/receipt/RCP-20240115-001"
}
```

### Generate Receipt PDF

```http
GET /receipts/{receipt_id}/pdf
```

**Query Parameters:**
- `language` (string, optional): Language for PDF (`hi`, `en`)
- `template` (string, optional): Template style (`traditional`, `modern`, `cultural`)

**Response:**
```
Content-Type: application/pdf
Content-Disposition: attachment; filename="receipt-RCP-20240115-001.pdf"

[PDF Binary Data]
```

---

## Webhook Integration

### Register Webhook

```http
POST /webhooks
```

**Request Body:**
```json
{
  "url": "https://yourdomain.com/webhook",
  "events": [
    "money_order.created",
    "money_order.completed",
    "swap.executed",
    "liquidity.added",
    "cultural.festival_started"
  ],
  "secret": "your_webhook_secret"
}
```

### Webhook Event Format

```json
{
  "event": "money_order.completed",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "order_id": "MO-20240115-001",
    "receipt_id": "RCP-20240115-001",
    "amount": {
      "denom": "unamo",
      "amount": "1000000"
    },
    "cultural_context": {
      "quote_displayed": true,
      "festival_bonus_applied": false,
      "patriotism_score": "0.85"
    }
  },
  "signature": "sha256=abc123def456..."
}
```

---

## SDK Examples

### JavaScript/TypeScript

```typescript
import { MoneyOrderClient } from '@deshchain/sdk';

const client = new MoneyOrderClient({
  apiKey: 'your_api_key',
  baseUrl: 'https://api.deshchain.org/v1/moneyorder',
  culture: {
    language: 'hi-IN',
    includeQuotes: true
  }
});

// Create money order
const moneyOrder = await client.createMoneyOrder({
  sender: 'desh1abc...',
  receiver: 'desh1def...',
  amount: { denom: 'unamo', amount: '1000000' },
  poolId: '1'
});

console.log(moneyOrder.receipt.cultural_quote);
```

### Python

```python
from deshchain_sdk import MoneyOrderClient

client = MoneyOrderClient(
    api_key='your_api_key',
    base_url='https://api.deshchain.org/v1/moneyorder',
    culture={'language': 'hi-IN', 'include_quotes': True}
)

# Get pool information
pool = client.get_pool('1')
print(f"Trust Score: {pool['trust_score']}")
print(f"Cultural Quote: {pool['patriotism_quote']}")
```

### Go

```go
package main

import (
    "github.com/deshchain/sdk-go/moneyorder"
)

func main() {
    client := moneyorder.NewClient(&moneyorder.Config{
        APIKey:  "your_api_key",
        BaseURL: "https://api.deshchain.org/v1/moneyorder",
        Culture: &moneyorder.CultureConfig{
            Language:      "hi-IN",
            IncludeQuotes: true,
        },
    })

    quote, err := client.GetRandomQuote("prosperity")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Quote: %s - %s\n", quote.Text, quote.Author)
}
```

---

## Testing & Development

### Testnet Endpoints

```
Base URL: https://testnet-api.deshchain.org/v1/moneyorder
Chain ID: deshchain-testnet-1
```

### Development Tools

- **API Explorer**: https://docs.deshchain.org/api-explorer
- **SDK Documentation**: https://sdk.deshchain.org
- **Postman Collection**: Available in docs repository
- **OpenAPI Specification**: https://api.deshchain.org/v1/openapi.json

### Cultural Testing Data

Test cultural quotes and festivals are available in testnet for development purposes:

```json
{
  "test_quotes": [
    {
      "id": "test_1",
      "text": "Test cultural quote in Hindi - परीक्षा उद्धरण",
      "category": "testing"
    }
  ],
  "test_festivals": [
    {
      "name": "Dev Festival",
      "bonus_rate": "0.50",
      "always_active": true
    }
  ]
}
```

## Support

- **Documentation**: https://docs.deshchain.org
- **API Support**: api-support@deshchain.org
- **Community**: https://discord.gg/deshchain
- **GitHub**: https://github.com/deshchain/deshchain