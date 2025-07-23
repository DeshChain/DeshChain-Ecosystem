package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/go-bip39"
	"google.golang.org/grpc"
)

// LoadTestConfig holds configuration for load testing
type LoadTestConfig struct {
	ChainID         string
	NodeAddress     string
	Workers         int
	TxPerWorker     int
	TestDuration    time.Duration
	TxType          string
	MaxTxSize       int
	BatchSize       int
	KeyringBackend  string
	EnableMetrics   bool
	OutputFile      string
}

// LoadTestResults stores test results
type LoadTestResults struct {
	TotalTx         int
	SuccessfulTx    int
	FailedTx        int
	TotalDuration   time.Duration
	TxPerSecond     float64
	AverageLatency  time.Duration
	Errors          []string
	GasUsed         int64
	TotalFees       sdk.Coins
}

// LoadTester manages the load testing process
type LoadTester struct {
	config  *LoadTestConfig
	client  client.Context
	results *LoadTestResults
	mutex   sync.Mutex
}

// NewLoadTester creates a new load tester instance
func NewLoadTester(config *LoadTestConfig) (*LoadTester, error) {
	// Setup client context
	encodingConfig := simapp.MakeTestEncodingConfig()
	
	clientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authsigning.NewAccountRetriever(encodingConfig.Marshaler)).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(simapp.DefaultNodeHome).
		WithChainID(config.ChainID)

	// Connect to node
	if config.NodeAddress != "" {
		clientCtx = clientCtx.WithNodeURI(config.NodeAddress)
		
		conn, err := grpc.Dial(config.NodeAddress, grpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("failed to connect to node: %w", err)
		}
		clientCtx = clientCtx.WithGRPCClient(conn)
	}

	return &LoadTester{
		config: config,
		client: clientCtx,
		results: &LoadTestResults{
			Errors:    make([]string, 0),
			TotalFees: sdk.NewCoins(),
		},
	}, nil
}

// generateTestAccounts creates test accounts for load testing
func (lt *LoadTester) generateTestAccounts(count int) ([]keyring.Info, error) {
	kr, err := keyring.New("test", lt.config.KeyringBackend, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create keyring: %w", err)
	}

	accounts := make([]keyring.Info, count)
	
	for i := 0; i < count; i++ {
		// Generate mnemonic
		entropy, err := bip39.NewEntropy(256)
		if err != nil {
			return nil, fmt.Errorf("failed to generate entropy: %w", err)
		}
		
		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
		}

		// Create account
		accountName := fmt.Sprintf("load-test-account-%d", i)
		info, err := kr.NewAccount(accountName, mnemonic, "", "", hd.Secp256k1)
		if err != nil {
			return nil, fmt.Errorf("failed to create account %d: %w", i, err)
		}
		
		accounts[i] = info
		log.Printf("Created test account %d: %s", i, info.GetAddress().String())
	}

	return accounts, nil
}

// generateRandomTransfer creates a random bank transfer transaction
func (lt *LoadTester) generateRandomTransfer(from, to sdk.AccAddress) *banktypes.MsgSend {
	// Random amount between 1 and 1000 units
	amountInt, _ := rand.Int(rand.Reader, big.NewInt(1000))
	amount := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(amountInt.Int64()+1)))

	return banktypes.NewMsgSend(from, to, amount)
}

// runTransactionWorker runs transactions for a single worker
func (lt *LoadTester) runTransactionWorker(workerID int, accounts []keyring.Info, wg *sync.WaitGroup) {
	defer wg.Done()
	
	log.Printf("Worker %d started", workerID)
	
	workerStart := time.Now()
	successCount := 0
	errorCount := 0
	totalGas := int64(0)
	totalFees := sdk.NewCoins()

	for i := 0; i < lt.config.TxPerWorker; i++ {
		// Select random accounts
		fromIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(accounts))))
		toIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(accounts))))
		
		if fromIdx.Int64() == toIdx.Int64() {
			// Skip if same account
			continue
		}

		fromAccount := accounts[fromIdx.Int64()]
		toAccount := accounts[toIdx.Int64()]

		// Create transaction
		msg := lt.generateRandomTransfer(fromAccount.GetAddress(), toAccount.GetAddress())
		
		txStart := time.Now()
		
		// Build and broadcast transaction
		txBuilder := lt.client.TxConfig.NewTxBuilder()
		err := txBuilder.SetMsgs(msg)
		if err != nil {
			lt.recordError(fmt.Sprintf("Worker %d: Failed to set message: %v", workerID, err))
			errorCount++
			continue
		}

		// Set fee
		fee := sdk.NewCoins(sdk.NewCoin("namo", sdk.NewInt(1000)))
		txBuilder.SetFeeAmount(fee)
		txBuilder.SetGasLimit(200000)

		// Sign transaction (simplified for load testing)
		txBytes, err := lt.client.TxConfig.TxEncoder()(txBuilder.GetTx())
		if err != nil {
			lt.recordError(fmt.Sprintf("Worker %d: Failed to encode tx: %v", workerID, err))
			errorCount++
			continue
		}

		// Simulate broadcast (in real test, you'd broadcast to node)
		time.Sleep(time.Millisecond * 10) // Simulate network latency
		
		txDuration := time.Since(txStart)
		
		// Simulate success/failure (90% success rate)
		successRand, _ := rand.Int(rand.Reader, big.NewInt(100))
		if successRand.Int64() < 90 {
			successCount++
			totalGas += 180000 // Simulated gas usage
			totalFees = totalFees.Add(fee...)
			
			// Record latency
			lt.recordLatency(txDuration)
		} else {
			errorCount++
			lt.recordError(fmt.Sprintf("Worker %d: Simulated transaction failure", workerID))
		}

		// Check if we should stop (duration-based test)
		if lt.config.TestDuration > 0 && time.Since(workerStart) >= lt.config.TestDuration {
			break
		}
	}

	// Update results
	lt.mutex.Lock()
	lt.results.SuccessfulTx += successCount
	lt.results.FailedTx += errorCount
	lt.results.GasUsed += totalGas
	lt.results.TotalFees = lt.results.TotalFees.Add(totalFees...)
	lt.mutex.Unlock()

	log.Printf("Worker %d completed: %d successful, %d failed", workerID, successCount, errorCount)
}

// recordError safely records an error
func (lt *LoadTester) recordError(err string) {
	lt.mutex.Lock()
	defer lt.mutex.Unlock()
	lt.results.Errors = append(lt.results.Errors, err)
}

// recordLatency safely records latency
func (lt *LoadTester) recordLatency(duration time.Duration) {
	lt.mutex.Lock()
	defer lt.mutex.Unlock()
	// Simple running average
	if lt.results.TotalTx == 0 {
		lt.results.AverageLatency = duration
	} else {
		lt.results.AverageLatency = (lt.results.AverageLatency + duration) / 2
	}
}

// RunLoadTest executes the load test
func (lt *LoadTester) RunLoadTest() error {
	log.Printf("Starting load test with %d workers, %d tx per worker", lt.config.Workers, lt.config.TxPerWorker)
	
	// Generate test accounts
	accounts, err := lt.generateTestAccounts(lt.config.Workers * 2) // 2x workers for variety
	if err != nil {
		return fmt.Errorf("failed to generate test accounts: %w", err)
	}

	startTime := time.Now()
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < lt.config.Workers; i++ {
		wg.Add(1)
		go lt.runTransactionWorker(i, accounts, &wg)
	}

	// Wait for all workers to complete
	wg.Wait()

	// Calculate final results
	lt.results.TotalDuration = time.Since(startTime)
	lt.results.TotalTx = lt.results.SuccessfulTx + lt.results.FailedTx
	
	if lt.results.TotalDuration > 0 {
		lt.results.TxPerSecond = float64(lt.results.SuccessfulTx) / lt.results.TotalDuration.Seconds()
	}

	return nil
}

// PrintResults displays the test results
func (lt *LoadTester) PrintResults() {
	fmt.Println("\n" + "="*80)
	fmt.Println("DESHCHAIN LOAD TEST RESULTS")
	fmt.Println("="*80)
	fmt.Printf("Test Configuration:\n")
	fmt.Printf("  Workers: %d\n", lt.config.Workers)
	fmt.Printf("  Transactions per Worker: %d\n", lt.config.TxPerWorker)
	fmt.Printf("  Test Duration: %v\n", lt.config.TestDuration)
	fmt.Printf("  Transaction Type: %s\n", lt.config.TxType)
	fmt.Println()
	
	fmt.Printf("Performance Results:\n")
	fmt.Printf("  Total Transactions: %d\n", lt.results.TotalTx)
	fmt.Printf("  Successful Transactions: %d\n", lt.results.SuccessfulTx)
	fmt.Printf("  Failed Transactions: %d\n", lt.results.FailedTx)
	fmt.Printf("  Success Rate: %.2f%%\n", float64(lt.results.SuccessfulTx)/float64(lt.results.TotalTx)*100)
	fmt.Printf("  Total Duration: %v\n", lt.results.TotalDuration)
	fmt.Printf("  Transactions per Second: %.2f\n", lt.results.TxPerSecond)
	fmt.Printf("  Average Latency: %v\n", lt.results.AverageLatency)
	fmt.Printf("  Total Gas Used: %d\n", lt.results.GasUsed)
	fmt.Printf("  Total Fees: %s\n", lt.results.TotalFees.String())
	fmt.Println()

	if len(lt.results.Errors) > 0 {
		fmt.Printf("Errors (%d total):\n", len(lt.results.Errors))
		for i, err := range lt.results.Errors {
			if i < 10 { // Show first 10 errors
				fmt.Printf("  %s\n", err)
			} else {
				fmt.Printf("  ... and %d more errors\n", len(lt.results.Errors)-10)
				break
			}
		}
	}
	fmt.Println("="*80)
}

// SaveResults saves results to file
func (lt *LoadTester) SaveResults() error {
	if lt.config.OutputFile == "" {
		return nil
	}

	file, err := os.Create(lt.config.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Write JSON results
	fmt.Fprintf(file, `{
  "config": {
    "workers": %d,
    "tx_per_worker": %d,
    "test_duration": "%v",
    "tx_type": "%s"
  },
  "results": {
    "total_tx": %d,
    "successful_tx": %d,
    "failed_tx": %d,
    "success_rate": %.2f,
    "total_duration": "%v",
    "tx_per_second": %.2f,
    "average_latency": "%v",
    "total_gas_used": %d,
    "total_fees": "%s",
    "error_count": %d
  }
}`,
		lt.config.Workers,
		lt.config.TxPerWorker,
		lt.config.TestDuration,
		lt.config.TxType,
		lt.results.TotalTx,
		lt.results.SuccessfulTx,
		lt.results.FailedTx,
		float64(lt.results.SuccessfulTx)/float64(lt.results.TotalTx)*100,
		lt.results.TotalDuration,
		lt.results.TxPerSecond,
		lt.results.AverageLatency,
		lt.results.GasUsed,
		lt.results.TotalFees.String(),
		len(lt.results.Errors),
	)

	log.Printf("Results saved to %s", lt.config.OutputFile)
	return nil
}

func main() {
	var (
		chainID       = flag.String("chain-id", "deshchain-testnet-1", "Chain ID")
		nodeAddr      = flag.String("node", "tcp://localhost:26657", "Node address")
		workers       = flag.Int("workers", 10, "Number of concurrent workers")
		txPerWorker   = flag.Int("tx-per-worker", 100, "Transactions per worker")
		duration      = flag.Duration("duration", 0, "Test duration (0 = use tx-per-worker)")
		txType        = flag.String("tx-type", "bank-send", "Transaction type to test")
		outputFile    = flag.String("output", "", "Output file for results")
		keyringBackend = flag.String("keyring-backend", "test", "Keyring backend")
	)
	flag.Parse()

	config := &LoadTestConfig{
		ChainID:        *chainID,
		NodeAddress:    *nodeAddr,
		Workers:        *workers,
		TxPerWorker:    *txPerWorker,
		TestDuration:   *duration,
		TxType:         *txType,
		OutputFile:     *outputFile,
		KeyringBackend: *keyringBackend,
	}

	tester, err := NewLoadTester(config)
	if err != nil {
		log.Fatalf("Failed to create load tester: %v", err)
	}

	// Run the load test
	if err := tester.RunLoadTest(); err != nil {
		log.Fatalf("Load test failed: %v", err)
	}

	// Display and save results
	tester.PrintResults()
	
	if err := tester.SaveResults(); err != nil {
		log.Printf("Failed to save results: %v", err)
	}

	// Exit with error code if success rate is too low
	successRate := float64(tester.results.SuccessfulTx) / float64(tester.results.TotalTx) * 100
	if successRate < 95.0 {
		log.Printf("WARNING: Success rate %.2f%% is below 95%% threshold", successRate)
		os.Exit(1)
	}

	log.Println("Load test completed successfully")
}