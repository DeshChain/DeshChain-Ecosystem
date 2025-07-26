package keeper

import (
	"context"
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/oracle/types"
)

// OracleClient interface for external price data sources
type OracleClient interface {
	GetPrice(ctx context.Context, symbol string) (*big.Int, error)
	IsHealthy() bool
	GetName() string
	GetWeight() sdk.Dec
}

// ChainlinkClient implements Oracle interface for Chainlink price feeds
type ChainlinkClient struct {
	endpoint    string
	apiKey      string
	weight      sdk.Dec
	timeout     time.Duration
	lastUpdate time.Time
	healthy     bool
}

// NewChainlinkClient creates a new Chainlink oracle client
func NewChainlinkClient(endpoint, apiKey string, weight sdk.Dec) *ChainlinkClient {
	return &ChainlinkClient{
		endpoint: endpoint,
		apiKey:   apiKey,
		weight:   weight,
		timeout:  30 * time.Second,
		healthy:  true,
	}
}

// GetPrice fetches price from Chainlink
func (c *ChainlinkClient) GetPrice(ctx context.Context, symbol string) (*big.Int, error) {
	// TODO: Implement actual Chainlink API integration
	// For now, return mock data with proper error handling
	
	// Simulate Chainlink price feed response
	mockPrices := map[string]*big.Int{
		"DINR": big.NewInt(100000000), // $1.00 in 8 decimals
		"BTC":  big.NewInt(4500000000000), // $45,000 in 8 decimals
		"ETH":  big.NewInt(300000000000), // $3,000 in 8 decimals
		"NAMO": big.NewInt(5000000), // $0.05 in 8 decimals
	}
	
	price, exists := mockPrices[symbol]
	if !exists {
		c.healthy = false
		return nil, fmt.Errorf("price not available for symbol %s", symbol)
	}
	
	c.lastUpdate = time.Now()
	c.healthy = true
	return price, nil
}

// IsHealthy returns the health status of the oracle
func (c *ChainlinkClient) IsHealthy() bool {
	// Consider unhealthy if no update in last 10 minutes
	if time.Since(c.lastUpdate) > 10*time.Minute {
		c.healthy = false
	}
	return c.healthy
}

// GetName returns the oracle name
func (c *ChainlinkClient) GetName() string {
	return "Chainlink"
}

// GetWeight returns the oracle weight in aggregation
func (c *ChainlinkClient) GetWeight() sdk.Dec {
	return c.weight
}

// BandProtocolClient implements Oracle interface for Band Protocol
type BandProtocolClient struct {
	endpoint    string
	weight      sdk.Dec
	timeout     time.Duration
	lastUpdate time.Time
	healthy     bool
}

// NewBandProtocolClient creates a new Band Protocol oracle client
func NewBandProtocolClient(endpoint string, weight sdk.Dec) *BandProtocolClient {
	return &BandProtocolClient{
		endpoint: endpoint,
		weight:   weight,
		timeout:  30 * time.Second,
		healthy:  true,
	}
}

// GetPrice fetches price from Band Protocol
func (b *BandProtocolClient) GetPrice(ctx context.Context, symbol string) (*big.Int, error) {
	// TODO: Implement actual Band Protocol integration
	// For now, return mock data with slight variation for testing aggregation
	
	mockPrices := map[string]*big.Int{
		"DINR": big.NewInt(100100000), // $1.001 in 8 decimals (slight variance)
		"BTC":  big.NewInt(4505000000000), // $45,050 in 8 decimals
		"ETH":  big.NewInt(300200000000), // $3,002 in 8 decimals
		"NAMO": big.NewInt(5010000), // $0.0501 in 8 decimals
	}
	
	price, exists := mockPrices[symbol]
	if !exists {
		b.healthy = false
		return nil, fmt.Errorf("price not available for symbol %s", symbol)
	}
	
	b.lastUpdate = time.Now()
	b.healthy = true
	return price, nil
}

// IsHealthy returns the health status of the oracle
func (b *BandProtocolClient) IsHealthy() bool {
	if time.Since(b.lastUpdate) > 10*time.Minute {
		b.healthy = false
	}
	return b.healthy
}

// GetName returns the oracle name
func (b *BandProtocolClient) GetName() string {
	return "BandProtocol"
}

// GetWeight returns the oracle weight in aggregation
func (b *BandProtocolClient) GetWeight() sdk.Dec {
	return b.weight
}

// PythNetworkClient implements Oracle interface for Pyth Network
type PythNetworkClient struct {
	endpoint    string
	weight      sdk.Dec
	timeout     time.Duration
	lastUpdate time.Time
	healthy     bool
}

// NewPythNetworkClient creates a new Pyth Network oracle client
func NewPythNetworkClient(endpoint string, weight sdk.Dec) *PythNetworkClient {
	return &PythNetworkClient{
		endpoint: endpoint,
		weight:   weight,
		timeout:  30 * time.Second,
		healthy:  true,
	}
}

// GetPrice fetches price from Pyth Network
func (p *PythNetworkClient) GetPrice(ctx context.Context, symbol string) (*big.Int, error) {
	// TODO: Implement actual Pyth Network integration
	// For now, return mock data with slight variation
	
	mockPrices := map[string]*big.Int{
		"DINR": big.NewInt(99950000), // $0.9995 in 8 decimals (slight variance)
		"BTC":  big.NewInt(4498000000000), // $44,980 in 8 decimals
		"ETH":  big.NewInt(299800000000), // $2,998 in 8 decimals
		"NAMO": big.NewInt(4995000), // $0.04995 in 8 decimals
	}
	
	price, exists := mockPrices[symbol]
	if !exists {
		p.healthy = false
		return nil, fmt.Errorf("price not available for symbol %s", symbol)
	}
	
	p.lastUpdate = time.Now()
	p.healthy = true
	return price, nil
}

// IsHealthy returns the health status of the oracle
func (p *PythNetworkClient) IsHealthy() bool {
	if time.Since(p.lastUpdate) > 10*time.Minute {
		p.healthy = false
	}
	return p.healthy
}

// GetName returns the oracle name
func (p *PythNetworkClient) GetName() string {
	return "PythNetwork"
}

// GetWeight returns the oracle weight in aggregation
func (p *PythNetworkClient) GetWeight() sdk.Dec {
	return p.weight
}

// PriceAggregator aggregates prices from multiple oracle sources
type PriceAggregator struct {
	sources     []OracleClient
	minSources  int
	maxDeviation sdk.Dec // Maximum allowed deviation between sources
}

// NewPriceAggregator creates a new price aggregator
func NewPriceAggregator(minSources int, maxDeviation sdk.Dec) *PriceAggregator {
	return &PriceAggregator{
		sources:      make([]OracleClient, 0),
		minSources:   minSources,
		maxDeviation: maxDeviation,
	}
}

// AddSource adds an oracle source to the aggregator
func (pa *PriceAggregator) AddSource(source OracleClient) {
	pa.sources = append(pa.sources, source)
}

// GetAggregatedPrice calculates weighted average price from all healthy sources
func (pa *PriceAggregator) GetAggregatedPrice(ctx context.Context, symbol string) (sdk.Dec, error) {
	if len(pa.sources) == 0 {
		return sdk.ZeroDec(), fmt.Errorf("no oracle sources configured")
	}
	
	var validPrices []priceData
	totalWeight := sdk.ZeroDec()
	
	// Collect prices from all healthy sources
	for _, source := range pa.sources {
		if !source.IsHealthy() {
			continue
		}
		
		price, err := source.GetPrice(ctx, symbol)
		if err != nil {
			continue
		}
		
		priceDecimal := sdk.NewDecFromBigInt(price).QuoInt64(100000000) // Convert from 8 decimals
		weight := source.GetWeight()
		
		validPrices = append(validPrices, priceData{
			price:  priceDecimal,
			weight: weight,
			source: source.GetName(),
		})
		totalWeight = totalWeight.Add(weight)
	}
	
	// Check if we have minimum required sources
	if len(validPrices) < pa.minSources {
		return sdk.ZeroDec(), fmt.Errorf("insufficient healthy oracle sources: %d < %d", len(validPrices), pa.minSources)
	}
	
	// Validate price deviation
	if err := pa.validatePriceDeviation(validPrices); err != nil {
		return sdk.ZeroDec(), err
	}
	
	// Calculate weighted average
	weightedSum := sdk.ZeroDec()
	for _, data := range validPrices {
		weightedSum = weightedSum.Add(data.price.Mul(data.weight))
	}
	
	if totalWeight.IsZero() {
		return sdk.ZeroDec(), fmt.Errorf("total weight is zero")
	}
	
	aggregatedPrice := weightedSum.Quo(totalWeight)
	return aggregatedPrice, nil
}

type priceData struct {
	price  sdk.Dec
	weight sdk.Dec
	source string
}

// validatePriceDeviation checks if prices from different sources are within acceptable range
func (pa *PriceAggregator) validatePriceDeviation(prices []priceData) error {
	if len(prices) < 2 {
		return nil // No validation needed for single source
	}
	
	// Find min and max prices
	minPrice := prices[0].price
	maxPrice := prices[0].price
	
	for _, data := range prices {
		if data.price.LT(minPrice) {
			minPrice = data.price
		}
		if data.price.GT(maxPrice) {
			maxPrice = data.price
		}
	}
	
	// Calculate deviation as percentage
	if minPrice.IsZero() {
		return fmt.Errorf("invalid price data: zero price detected")
	}
	
	deviation := maxPrice.Sub(minPrice).Quo(minPrice)
	
	if deviation.GT(pa.maxDeviation) {
		return fmt.Errorf("price deviation too high: %s > %s", deviation.String(), pa.maxDeviation.String())
	}
	
	return nil
}

// GetHealthySourceCount returns the number of healthy oracle sources
func (pa *PriceAggregator) GetHealthySourceCount() int {
	count := 0
	for _, source := range pa.sources {
		if source.IsHealthy() {
			count++
		}
	}
	return count
}

// GetSourceStatus returns the status of all configured sources
func (pa *PriceAggregator) GetSourceStatus() []types.OracleSourceStatus {
	status := make([]types.OracleSourceStatus, len(pa.sources))
	
	for i, source := range pa.sources {
		status[i] = types.OracleSourceStatus{
			Name:    source.GetName(),
			Healthy: source.IsHealthy(),
			Weight:  source.GetWeight(),
		}
	}
	
	return status
}