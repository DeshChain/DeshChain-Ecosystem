# DeshChain Python SDK

Official Python SDK for interacting with the DeshChain blockchain - a cultural heritage-focused blockchain platform built on Cosmos SDK.

## Features

- 🌟 **Cultural Integration**: Festival bonuses, regional quotes, and cultural heritage preservation
- 💰 **DeFi Lending**: Krishi Mitra (Agriculture), Vyavasaya Mitra (Business), Shiksha Mitra (Education)
- 🚀 **Sikkebaaz Launchpad**: Anti-pump memecoin platform with community governance
- 💸 **Money Order DEX**: Traditional money transfers reimagined for blockchain
- 🏛️ **Governance**: Founder protection with community participation
- 🎭 **22 Languages**: Multi-language support for Indian regional languages
- ⚡ **Async Support**: Both sync and async client implementations

## Installation

```bash
pip install deshchain-sdk
```

## Quick Start

### Synchronous Client

```python
from deshchain import DeshChainClient

# Connect to DeshChain network
client = DeshChainClient.connect("mainnet")

# Get chain info
chain_info = client.get_chain_info()
print(f"Connected to: {chain_info.chain_id}")

# Check current festival
festival = client.get_current_festival()
if festival:
    print(f"Active festival: {festival.name}")

# Get lending statistics
lending_stats = client.get_lending_stats()
print(f"Total loans disbursed: {lending_stats['combined']['total_disbursed']}")

# Close client
client.close()
```

### Asynchronous Client

```python
import asyncio
from deshchain import AsyncDeshChainClient

async def main():
    async with AsyncDeshChainClient.connect("mainnet") as client:
        # Get chain info
        chain_info = await client.get_chain_info()
        print(f"Connected to: {chain_info.chain_id}")
        
        # Get current festival
        festival = await client.cultural.get_current_festival()
        if festival:
            print(f"Active festival: {festival.name}")

# Run async client
asyncio.run(main())
```

## Module Clients

### Cultural Heritage

```python
# Get daily cultural quote
quote = client.cultural.get_daily_quote()
print(f"{quote.text} - {quote.author}")

# Get active festivals
festivals = client.cultural.get_active_festivals()

# Calculate festival bonus
bonus = client.cultural.calculate_festival_bonus(1000, "send")
print(f"Festival bonus: {bonus['percentage']}%")

# Get regional celebrations
celebrations = client.cultural.get_regional_celebrations("West Bengal")
```

### Lending (Krishi/Vyavasaya/Shiksha Mitra)

```python
# Get agricultural loan stats
krishi_stats = client.lending.get_krishi_mitra_stats()
print(f"Krishi Mitra - Total loans: {krishi_stats.total_loans}")

# Get business loan stats
vyavasaya_stats = client.lending.get_vyavasaya_mitra_stats()
print(f"Vyavasaya Mitra - Average rate: {vyavasaya_stats.average_rate}%")

# Get education loan stats
shiksha_stats = client.lending.get_shiksha_mitra_stats()
print(f"Shiksha Mitra - Default rate: {shiksha_stats.default_rate}%")

# Search loans
loans = client.lending.search_loans("farmer_123")
```

### Sikkebaaz Launchpad

```python
# Get featured tokens
featured_tokens = client.sikkebaaz.get_featured_tokens()
for token in featured_tokens:
    print(f"{token.symbol}: {token.name}")

# Search tokens
results = client.sikkebaaz.search_tokens("cultural")
print(f"Found {len(results)} cultural tokens")
```

### Money Order DEX

```python
# Get money order
order = client.money_order.get_money_order("order_123")
print(f"Order status: {order.status}")
```

### Governance

```python
# Get all proposals
proposals = client.governance.get_proposals()
print(f"Total proposals: {len(proposals)}")

# Get active proposals only
active_proposals = [p for p in proposals if p.status == "voting_period"]
```

## Network Configuration

```python
# Testnet
testnet_client = DeshChainClient.connect("testnet")

# Custom RPC endpoint
custom_client = DeshChainClient(
    rpc_url="https://custom-rpc.example.com",
    rest_url="https://custom-api.example.com",
    chain_id="deshchain-1"
)

# With custom timeout and retries
client = DeshChainClient.connect(
    "mainnet",
    timeout=60.0,
    retries=5
)
```

## Cultural Features

### Festival Integration

```python
from deshchain.utils import FestivalUtils

# Check current festivals
festival_info = FestivalUtils.is_festival_day()
if festival_info["is_festival"]:
    print("Today is a festival day!")
    for festival in festival_info["festivals"]:
        print(f"- {festival['name']}")

# Get festival bonus
bonus = FestivalUtils.get_festival_bonus()
print(f"Current festival bonus: {bonus}x")

# Get festival greeting
greeting = FestivalUtils.get_festival_greeting("Diwali")
print(greeting["english"])
print(greeting["hindi"])
```

### Regional Customization

```python
from deshchain.utils import CulturalUtils

# Get state from pincode
state = CulturalUtils.get_state_from_pincode("110001")
print(f"State: {state}")

# Get regional language
language = CulturalUtils.get_regional_language("Maharashtra")
print(f"Language: {language}")

# Get major festivals by state
festivals = CulturalUtils.get_major_festivals("West Bengal")
print(f"Major festivals: {festivals}")

# Format currency in Indian style
formatted = CulturalUtils.format_indian_currency(1500000)
print(formatted)  # ₹15.00 L
```

## Error Handling

```python
from deshchain.exceptions import (
    DeshChainError, 
    NetworkError, 
    TransactionError,
    ValidationError
)

try:
    result = client.get_account("invalid_address")
except NetworkError as e:
    print(f"Network issue: {e.message}")
    print(f"Status code: {e.status_code}")
except ValidationError as e:
    print(f"Validation error: {e.message}")
    print(f"Field: {e.field}")
except DeshChainError as e:
    print(f"DeshChain error [{e.code}]: {e.message}")
```

## Validation Utilities

```python
from deshchain.utils import ValidationUtils

# Validate addresses
result = ValidationUtils.validate_address("deshchain1...")
if not result["valid"]:
    print(f"Invalid address: {result['error']}")

# Validate Indian-specific data
ValidationUtils.validate_pincode("110001")
ValidationUtils.validate_aadhar("1234 5678 9012")
ValidationUtils.validate_pan("ABCDE1234F")
ValidationUtils.validate_phone_number("9876543210")
```

## Type Hints

The SDK provides full type hints for better development experience:

```python
from deshchain.types import Festival, Loan, LaunchpadToken

# All responses are properly typed
festivals: List[Festival] = client.cultural.get_active_festivals()
loans: List[Loan] = client.lending.search_loans("query")
tokens: List[LaunchpadToken] = client.sikkebaaz.get_featured_tokens()
```

## Context Managers

Use context managers for automatic resource cleanup:

```python
# Synchronous
with DeshChainClient.connect("mainnet") as client:
    # Your code here
    pass  # Client is automatically closed

# Asynchronous
async with AsyncDeshChainClient.connect("mainnet") as client:
    # Your async code here
    pass  # Client is automatically closed
```

## Development

### Installing for Development

```bash
git clone https://github.com/deshchain/deshchain-sdk-python
cd deshchain-sdk-python
pip install -e ".[dev]"
```

### Running Tests

```bash
pytest
pytest --cov=deshchain  # With coverage
```

### Code Formatting

```bash
black deshchain/
isort deshchain/
flake8 deshchain/
mypy deshchain/
```

## API Reference

Full API documentation is available at [docs.deshchain.network/sdk/python](https://docs.deshchain.network/sdk/python)

## Examples

More examples available in the [examples/](examples/) directory:

- [Basic Usage](examples/basic_usage.py)
- [Cultural Features](examples/cultural_features.py)
- [Lending Integration](examples/lending_example.py)
- [Async Usage](examples/async_example.py)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`pytest`)
4. Format code (`black . && isort .`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- GitHub Issues: [deshchain/deshchain-sdk-python/issues](https://github.com/deshchain/deshchain-sdk-python/issues)
- Discord: [DeshChain Community](https://discord.gg/deshchain)
- Documentation: [docs.deshchain.network](https://docs.deshchain.network)

---

**Jai Hind! 🇮🇳** - Building the future of Indian blockchain technology with cultural values at its core.