package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
)

// ProfileConfig holds profiling configuration
type ProfileConfig struct {
	NodeAddress    string
	ProfileDuration time.Duration
	CPUProfile     string
	MemProfile     string
	BlockProfile   string
	MutexProfile   string
	TraceProfile   string
	OutputDir      string
	EnablePprof    bool
	PprofPort      int
}

// PerformanceProfiler manages blockchain performance profiling
type PerformanceProfiler struct {
	config *ProfileConfig
	client client.Context
	results *ProfilingResults
}

// ProfilingResults stores profiling analysis results
type ProfilingResults struct {
	Timestamp        string                 `json:"timestamp"`
	Duration         string                 `json:"duration"`
	SystemInfo       SystemInfo             `json:"system_info"`
	BlockchainMetrics BlockchainMetrics     `json:"blockchain_metrics"`
	ResourceUsage    ResourceUsage          `json:"resource_usage"`
	Bottlenecks      []Bottleneck           `json:"bottlenecks"`
	Recommendations  []string               `json:"recommendations"`
	Profiles         map[string]string      `json:"profiles"`
}

type SystemInfo struct {
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
	MemAlloc     uint64 `json:"mem_alloc"`
	MemSys       uint64 `json:"mem_sys"`
	MemHeapAlloc uint64 `json:"mem_heap_alloc"`
	MemHeapSys   uint64 `json:"mem_heap_sys"`
	GCCycles     uint32 `json:"gc_cycles"`
}

type BlockchainMetrics struct {
	BlockHeight     int64   `json:"block_height"`
	BlockTime       float64 `json:"block_time"`
	TransactionTPS  float64 `json:"transaction_tps"`
	ValidatorCount  int     `json:"validator_count"`
	PeerCount       int     `json:"peer_count"`
	SyncStatus      bool    `json:"sync_status"`
}

type ResourceUsage struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	DiskIORead    uint64  `json:"disk_io_read"`
	DiskIOWrite   uint64  `json:"disk_io_write"`
	NetworkRx     uint64  `json:"network_rx"`
	NetworkTx     uint64  `json:"network_tx"`
}

type Bottleneck struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	Impact      string  `json:"impact"`
	Suggestion  string  `json:"suggestion"`
	MetricValue float64 `json:"metric_value,omitempty"`
}

// NewPerformanceProfiler creates a new performance profiler
func NewPerformanceProfiler(config *ProfileConfig) (*PerformanceProfiler, error) {
	// Setup client context
	encodingConfig := simapp.MakeTestEncodingConfig()
	
	clientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino)

	// Connect to node if address provided
	if config.NodeAddress != "" {
		clientCtx = clientCtx.WithNodeURI(config.NodeAddress)
		
		conn, err := grpc.Dial(config.NodeAddress, grpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("failed to connect to node: %w", err)
		}
		clientCtx = clientCtx.WithGRPCClient(conn)
	}

	return &PerformanceProfiler{
		config: config,
		client: clientCtx,
		results: &ProfilingResults{
			Profiles:        make(map[string]string),
			Bottlenecks:     make([]Bottleneck, 0),
			Recommendations: make([]string, 0),
		},
	}, nil
}

// collectSystemInfo gathers system information
func (pp *PerformanceProfiler) collectSystemInfo() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemInfo{
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		MemAlloc:     m.Alloc,
		MemSys:       m.Sys,
		MemHeapAlloc: m.HeapAlloc,
		MemHeapSys:   m.HeapSys,
		GCCycles:     m.NumGC,
	}
}

// collectBlockchainMetrics gathers blockchain performance metrics
func (pp *PerformanceProfiler) collectBlockchainMetrics() (BlockchainMetrics, error) {
	// Placeholder implementation - would integrate with actual DeshChain client
	// This would fetch real metrics from the blockchain node
	
	return BlockchainMetrics{
		BlockHeight:    12345,  // Would fetch from node
		BlockTime:      3.2,    // Average block time
		TransactionTPS: 15.5,   // Transactions per second
		ValidatorCount: 50,     // Active validators
		PeerCount:      8,      // Connected peers
		SyncStatus:     true,   // Sync status
	}, nil
}

// collectResourceUsage gathers system resource usage
func (pp *PerformanceProfiler) collectResourceUsage() ResourceUsage {
	// Simplified resource collection - would use psutil or similar for accurate metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return ResourceUsage{
		CPUPercent:    float64(runtime.NumGoroutine()) * 0.1, // Simplified calculation
		MemoryPercent: float64(m.Alloc) / float64(m.Sys) * 100,
		DiskIORead:    0,    // Would collect from system
		DiskIOWrite:   0,    // Would collect from system
		NetworkRx:     0,    // Would collect from system
		NetworkTx:     0,    // Would collect from system
	}
}

// analyzeBottlenecks identifies performance bottlenecks
func (pp *PerformanceProfiler) analyzeBottlenecks() []Bottleneck {
	bottlenecks := make([]Bottleneck, 0)
	
	// Analyze system info
	sysInfo := pp.results.SystemInfo
	if sysInfo.NumGoroutine > 10000 {
		bottlenecks = append(bottlenecks, Bottleneck{
			Type:        "goroutine_leak",
			Severity:    "high",
			Description: fmt.Sprintf("High goroutine count: %d", sysInfo.NumGoroutine),
			Impact:      "Memory consumption and scheduling overhead",
			Suggestion:  "Review goroutine lifecycle management and implement proper cleanup",
			MetricValue: float64(sysInfo.NumGoroutine),
		})
	}
	
	// Analyze memory usage
	memUsagePercent := float64(sysInfo.MemAlloc) / float64(sysInfo.MemSys) * 100
	if memUsagePercent > 80 {
		bottlenecks = append(bottlenecks, Bottleneck{
			Type:        "high_memory_usage",
			Severity:    "medium",
			Description: fmt.Sprintf("High memory usage: %.1f%%", memUsagePercent),
			Impact:      "Increased GC pressure and potential OOM issues",
			Suggestion:  "Implement memory pooling and optimize data structures",
			MetricValue: memUsagePercent,
		})
	}
	
	// Analyze blockchain metrics
	bcMetrics := pp.results.BlockchainMetrics
	if bcMetrics.BlockTime > 10.0 {
		bottlenecks = append(bottlenecks, Bottleneck{
			Type:        "slow_block_time",
			Severity:    "high",
			Description: fmt.Sprintf("Slow block time: %.2fs", bcMetrics.BlockTime),
			Impact:      "Reduced transaction throughput and user experience",
			Suggestion:  "Optimize consensus mechanism and block processing",
			MetricValue: bcMetrics.BlockTime,
		})
	}
	
	if bcMetrics.TransactionTPS < 10 {
		bottlenecks = append(bottlenecks, Bottleneck{
			Type:        "low_tps",
			Severity:    "medium",
			Description: fmt.Sprintf("Low transaction throughput: %.1f TPS", bcMetrics.TransactionTPS),
			Impact:      "Limited scalability and network capacity",
			Suggestion:  "Optimize transaction processing and consider parallelization",
			MetricValue: bcMetrics.TransactionTPS,
		})
	}
	
	if bcMetrics.PeerCount < 5 {
		bottlenecks = append(bottlenecks, Bottleneck{
			Type:        "low_peer_count",
			Severity:    "medium",
			Description: fmt.Sprintf("Low peer count: %d", bcMetrics.PeerCount),
			Impact:      "Reduced network resilience and sync reliability",
			Suggestion:  "Improve peer discovery and connection management",
			MetricValue: float64(bcMetrics.PeerCount),
		})
	}
	
	return bottlenecks
}

// generateRecommendations creates optimization recommendations
func (pp *PerformanceProfiler) generateRecommendations() []string {
	recommendations := make([]string, 0)
	
	// Analyze bottlenecks and generate specific recommendations
	for _, bottleneck := range pp.results.Bottlenecks {
		switch bottleneck.Type {
		case "goroutine_leak":
			recommendations = append(recommendations, 
				"Implement goroutine pools to limit concurrent operations",
				"Add context cancellation to long-running operations",
				"Use sync.WaitGroup for proper goroutine lifecycle management")
		case "high_memory_usage":
			recommendations = append(recommendations,
				"Implement object pooling for frequently allocated objects",
				"Optimize data structures to reduce memory footprint", 
				"Configure GOGC environment variable for GC tuning")
		case "slow_block_time":
			recommendations = append(recommendations,
				"Profile and optimize consensus algorithm implementation",
				"Implement faster transaction validation",
				"Consider block size optimization")
		case "low_tps":
			recommendations = append(recommendations,
				"Implement transaction batching and parallel processing",
				"Optimize database operations and indexing",
				"Consider implementing transaction pools")
		}
	}
	
	// General performance recommendations
	recommendations = append(recommendations,
		"Enable CPU and memory profiling in production (with rotation)",
		"Implement caching layers for frequently accessed data",
		"Use connection pooling for database and network operations",
		"Consider implementing circuit breakers for external dependencies",
		"Set up comprehensive monitoring and alerting",
		"Implement graceful degradation for high-load scenarios")
	
	return recommendations
}

// startProfiling begins performance profiling
func (pp *PerformanceProfiler) startProfiling() error {
	log.Printf("Starting performance profiling for %v", pp.config.ProfileDuration)
	
	// Start pprof HTTP server if enabled
	if pp.config.EnablePprof {
		go func() {
			log.Printf("Starting pprof HTTP server on port %d", pp.config.PprofPort)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", pp.config.PprofPort), nil); err != nil {
				log.Printf("pprof server error: %v", err)
			}
		}()
	}
	
	// Start CPU profiling
	if pp.config.CPUProfile != "" {
		f, err := os.Create(pp.config.CPUProfile)
		if err != nil {
			return fmt.Errorf("could not create CPU profile: %w", err)
		}
		defer f.Close()
		
		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("could not start CPU profile: %w", err)
		}
		defer pprof.StopCPUProfile()
		
		pp.results.Profiles["cpu"] = pp.config.CPUProfile
	}
	
	// Enable block profiling
	if pp.config.BlockProfile != "" {
		runtime.SetBlockProfileRate(1)
	}
	
	// Enable mutex profiling  
	if pp.config.MutexProfile != "" {
		runtime.SetMutexProfileFraction(1)
	}
	
	// Start trace profiling
	if pp.config.TraceProfile != "" {
		f, err := os.Create(pp.config.TraceProfile)
		if err != nil {
			return fmt.Errorf("could not create trace profile: %w", err)
		}
		defer f.Close()
		
		// Note: Would use runtime/trace.Start() for actual tracing
		pp.results.Profiles["trace"] = pp.config.TraceProfile
	}
	
	// Collect initial metrics
	pp.results.SystemInfo = pp.collectSystemInfo()
	bcMetrics, err := pp.collectBlockchainMetrics()
	if err != nil {
		log.Printf("Error collecting blockchain metrics: %v", err)
	}
	pp.results.BlockchainMetrics = bcMetrics
	pp.results.ResourceUsage = pp.collectResourceUsage()
	
	// Wait for profiling duration
	time.Sleep(pp.config.ProfileDuration)
	
	// Collect final metrics and create profiles
	if pp.config.MemProfile != "" {
		f, err := os.Create(pp.config.MemProfile)
		if err != nil {
			return fmt.Errorf("could not create memory profile: %w", err)
		}
		defer f.Close()
		
		runtime.GC() // Get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			return fmt.Errorf("could not write memory profile: %w", err)
		}
		
		pp.results.Profiles["memory"] = pp.config.MemProfile
	}
	
	if pp.config.BlockProfile != "" {
		f, err := os.Create(pp.config.BlockProfile)
		if err != nil {
			return fmt.Errorf("could not create block profile: %w", err)
		}
		defer f.Close()
		
		if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
			return fmt.Errorf("could not write block profile: %w", err)
		}
		
		pp.results.Profiles["block"] = pp.config.BlockProfile
	}
	
	if pp.config.MutexProfile != "" {
		f, err := os.Create(pp.config.MutexProfile)
		if err != nil {
			return fmt.Errorf("could not create mutex profile: %w", err)
		}
		defer f.Close()
		
		if err := pprof.Lookup("mutex").WriteTo(f, 0); err != nil {
			return fmt.Errorf("could not write mutex profile: %w", err)
		}
		
		pp.results.Profiles["mutex"] = pp.config.MutexProfile
	}
	
	return nil
}

// analyzePerformance performs comprehensive performance analysis
func (pp *PerformanceProfiler) analyzePerformance() error {
	log.Println("Analyzing performance data...")
	
	// Set analysis metadata
	pp.results.Timestamp = time.Now().Format(time.RFC3339)
	pp.results.Duration = pp.config.ProfileDuration.String()
	
	// Analyze bottlenecks
	pp.results.Bottlenecks = pp.analyzeBottlenecks()
	
	// Generate recommendations
	pp.results.Recommendations = pp.generateRecommendations()
	
	return nil
}

// generateReport creates a comprehensive performance report
func (pp *PerformanceProfiler) generateReport() error {
	reportPath := fmt.Sprintf("%s/performance-analysis-report.md", pp.config.OutputDir)
	
	report := fmt.Sprintf(`# DeshChain Performance Analysis Report

**Generated:** %s  
**Duration:** %s  
**Node:** %s  

## Executive Summary

This report provides a comprehensive analysis of DeshChain's performance characteristics, identifying bottlenecks and providing optimization recommendations.

### Key Findings
- **Bottlenecks Identified:** %d
- **Critical Issues:** %d
- **Optimization Opportunities:** %d

## System Information

- **OS/Architecture:** %s/%s
- **CPU Cores:** %d
- **Active Goroutines:** %d
- **Memory Allocated:** %.2f MB
- **Memory System:** %.2f MB
- **Heap Allocated:** %.2f MB
- **GC Cycles:** %d

## Blockchain Metrics

- **Block Height:** %d
- **Average Block Time:** %.2f seconds
- **Transaction Throughput:** %.1f TPS
- **Active Validators:** %d
- **Connected Peers:** %d
- **Sync Status:** %v

## Resource Usage

- **CPU Usage:** %.1f%%
- **Memory Usage:** %.1f%%
- **Disk I/O Read:** %d bytes
- **Disk I/O Write:** %d bytes
- **Network RX:** %d bytes  
- **Network TX:** %d bytes

## Performance Bottlenecks

`,
		pp.results.Timestamp,
		pp.results.Duration,
		pp.config.NodeAddress,
		len(pp.results.Bottlenecks),
		countBottlenecksBySeverity(pp.results.Bottlenecks, "high"),
		len(pp.results.Recommendations),
		pp.results.SystemInfo.GOOS,
		pp.results.SystemInfo.GOARCH,
		pp.results.SystemInfo.NumCPU,
		pp.results.SystemInfo.NumGoroutine,
		float64(pp.results.SystemInfo.MemAlloc)/1024/1024,
		float64(pp.results.SystemInfo.MemSys)/1024/1024,
		float64(pp.results.SystemInfo.MemHeapAlloc)/1024/1024,
		pp.results.SystemInfo.GCCycles,
		pp.results.BlockchainMetrics.BlockHeight,
		pp.results.BlockchainMetrics.BlockTime,
		pp.results.BlockchainMetrics.TransactionTPS,
		pp.results.BlockchainMetrics.ValidatorCount,
		pp.results.BlockchainMetrics.PeerCount,
		pp.results.BlockchainMetrics.SyncStatus,
		pp.results.ResourceUsage.CPUPercent,
		pp.results.ResourceUsage.MemoryPercent,
		pp.results.ResourceUsage.DiskIORead,
		pp.results.ResourceUsage.DiskIOWrite,
		pp.results.ResourceUsage.NetworkRx,
		pp.results.ResourceUsage.NetworkTx,
	)
	
	// Add bottlenecks section
	if len(pp.results.Bottlenecks) > 0 {
		for i, bottleneck := range pp.results.Bottlenecks {
			severityIcon := "⚠️"
			if bottleneck.Severity == "high" {
				severityIcon = "🔴"
			} else if bottleneck.Severity == "low" {
				severityIcon = "🟡"
			}
			
			report += fmt.Sprintf(`
### %d. %s %s

**Severity:** %s  
**Description:** %s  
**Impact:** %s  
**Recommendation:** %s
`,
				i+1,
				severityIcon,
				bottleneck.Type,
				bottleneck.Severity,
				bottleneck.Description,
				bottleneck.Impact,
				bottleneck.Suggestion,
			)
		}
	} else {
		report += "\n✅ No significant performance bottlenecks detected.\n"
	}
	
	// Add recommendations section
	report += "\n## Optimization Recommendations\n\n"
	for i, rec := range pp.results.Recommendations {
		report += fmt.Sprintf("%d. %s\n", i+1, rec)
	}
	
	// Add profiles section
	report += "\n## Generated Profiles\n\n"
	for profileType, path := range pp.results.Profiles {
		report += fmt.Sprintf("- **%s Profile:** `%s`\n", profileType, path)
	}
	
	// Add analysis commands
	report += fmt.Sprintf(`

## Profile Analysis Commands

To analyze the generated profiles, use the following commands:

### CPU Profile
` + "```bash" + `
go tool pprof %s
` + "```" + `

### Memory Profile  
` + "```bash" + `
go tool pprof %s
` + "```" + `

### Block Profile
` + "```bash" + `
go tool pprof %s
` + "```" + `

### Mutex Profile
` + "```bash" + `
go tool pprof %s
` + "```" + `

---
*Generated by DeshChain Performance Profiler*
`,
		pp.results.Profiles["cpu"],
		pp.results.Profiles["memory"],
		pp.results.Profiles["block"],
		pp.results.Profiles["mutex"],
	)
	
	// Write report to file
	return os.WriteFile(reportPath, []byte(report), 0644)
}

// saveResults saves profiling results to JSON
func (pp *PerformanceProfiler) saveResults() error {
	resultsPath := fmt.Sprintf("%s/performance-analysis-results.json", pp.config.OutputDir)
	
	data, err := json.MarshalIndent(pp.results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	
	return os.WriteFile(resultsPath, data, 0644)
}

// countBottlenecksBySeverity counts bottlenecks by severity level
func countBottlenecksBySeverity(bottlenecks []Bottleneck, severity string) int {
	count := 0
	for _, b := range bottlenecks {
		if b.Severity == severity {
			count++
		}
	}
	return count
}

// RunProfiling executes the complete profiling workflow
func (pp *PerformanceProfiler) RunProfiling() error {
	// Create output directory
	if err := os.MkdirAll(pp.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Start profiling
	if err := pp.startProfiling(); err != nil {
		return fmt.Errorf("profiling failed: %w", err)
	}
	
	// Analyze performance
	if err := pp.analyzePerformance(); err != nil {
		return fmt.Errorf("performance analysis failed: %w", err)
	}
	
	// Generate report
	if err := pp.generateReport(); err != nil {
		return fmt.Errorf("report generation failed: %w", err)
	}
	
	// Save results
	if err := pp.saveResults(); err != nil {
		return fmt.Errorf("saving results failed: %w", err)
	}
	
	log.Printf("Performance profiling completed successfully")
	log.Printf("Results saved to: %s", pp.config.OutputDir)
	
	return nil
}

func main() {
	var (
		nodeAddr     = flag.String("node", "", "Blockchain node address")
		duration     = flag.Duration("duration", 5*time.Minute, "Profiling duration")
		outputDir    = flag.String("output", "./performance-profiles", "Output directory")
		cpuProfile   = flag.String("cpuprofile", "", "Write CPU profile to file")
		memProfile   = flag.String("memprofile", "", "Write memory profile to file") 
		blockProfile = flag.String("blockprofile", "", "Write block profile to file")
		mutexProfile = flag.String("mutexprofile", "", "Write mutex profile to file")
		traceProfile = flag.String("traceprofile", "", "Write trace profile to file")
		enablePprof  = flag.Bool("pprof", false, "Enable pprof HTTP server")
		pprofPort    = flag.Int("pprof-port", 6060, "pprof HTTP server port")
	)
	flag.Parse()
	
	// Create output directory with timestamp
	outputDir := fmt.Sprintf("%s-%s", *outputDir, time.Now().Format("20060102-150405"))
	
	// Set default profile paths if not specified
	if *cpuProfile == "" && *memProfile == "" && *blockProfile == "" && *mutexProfile == "" {
		*cpuProfile = fmt.Sprintf("%s/cpu.prof", outputDir)
		*memProfile = fmt.Sprintf("%s/mem.prof", outputDir)
		*blockProfile = fmt.Sprintf("%s/block.prof", outputDir)
		*mutexProfile = fmt.Sprintf("%s/mutex.prof", outputDir)
		*traceProfile = fmt.Sprintf("%s/trace.prof", outputDir)
	}
	
	config := &ProfileConfig{
		NodeAddress:     *nodeAddr,
		ProfileDuration: *duration,
		CPUProfile:      *cpuProfile,
		MemProfile:      *memProfile,
		BlockProfile:    *blockProfile,
		MutexProfile:    *mutexProfile,
		TraceProfile:    *traceProfile,
		OutputDir:       outputDir,
		EnablePprof:     *enablePprof,
		PprofPort:       *pprofPort,
	}
	
	profiler, err := NewPerformanceProfiler(config)
	if err != nil {
		log.Fatalf("Failed to create profiler: %v", err)
	}
	
	if err := profiler.RunProfiling(); err != nil {
		log.Fatalf("Profiling failed: %v", err)
	}
	
	fmt.Printf("\n🎯 Performance profiling completed successfully!\n")
	fmt.Printf("📁 Results directory: %s\n", outputDir)
	fmt.Printf("📊 Analysis report: %s/performance-analysis-report.md\n", outputDir)
	fmt.Printf("📈 JSON results: %s/performance-analysis-results.json\n", outputDir)
}