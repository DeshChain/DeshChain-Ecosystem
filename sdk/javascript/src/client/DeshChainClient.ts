import { StargateClient, QueryClient } from '@cosmjs/stargate'
import { Tendermint34Client } from '@cosmjs/tendermint-rpc'
import { LendingClient } from '../modules/lending/LendingClient'
import { CulturalClient } from '../modules/cultural/CulturalClient'
import { SikkebaazClient } from '../modules/sikkebaaz/SikkebaazClient'
import { MoneyOrderClient } from '../modules/moneyorder/MoneyOrderClient'
import { GovernanceClient } from '../modules/governance/GovernanceClient'
import { DeshChainClientOptions, ChainInfo, NetworkStatus } from '../types'

/**
 * Main DeshChain client for read-only operations
 * 
 * Provides access to all DeshChain modules and blockchain data
 */
export class DeshChainClient {
  protected readonly client: StargateClient
  protected readonly tmClient: Tendermint34Client
  protected readonly options: DeshChainClientOptions

  // Module clients
  public readonly lending: LendingClient
  public readonly cultural: CulturalClient
  public readonly sikkebaaz: SikkebaazClient
  public readonly moneyOrder: MoneyOrderClient
  public readonly governance: GovernanceClient

  protected constructor(
    client: StargateClient,
    tmClient: Tendermint34Client,
    options: DeshChainClientOptions
  ) {
    this.client = client
    this.tmClient = tmClient
    this.options = options

    // Initialize module clients
    this.lending = new LendingClient(client, tmClient)
    this.cultural = new CulturalClient(client, tmClient)
    this.sikkebaaz = new SikkebaazClient(client, tmClient)
    this.moneyOrder = new MoneyOrderClient(client, tmClient)
    this.governance = new GovernanceClient(client, tmClient)
  }

  /**
   * Connect to DeshChain network
   */
  public static async connect(
    endpoint: string,
    options: Partial<DeshChainClientOptions> = {}
  ): Promise<DeshChainClient> {
    const fullOptions: DeshChainClientOptions = {
      chainId: 'deshchain-1',
      prefix: 'deshchain',
      gasPrice: '0.025unamo',
      ...options,
    }

    const tmClient = await Tendermint34Client.connect(endpoint)
    const client = await StargateClient.create(tmClient)

    return new DeshChainClient(client, tmClient, fullOptions)
  }

  /**
   * Get chain information
   */
  public async getChainInfo(): Promise<ChainInfo> {
    const status = await this.tmClient.status()
    const validators = await this.tmClient.validatorsAll()

    return {
      chainId: status.nodeInfo.network,
      nodeVersion: status.nodeInfo.version,
      blockHeight: status.syncInfo.latestBlockHeight,
      blockTime: status.syncInfo.latestBlockTime,
      validatorCount: validators.validators.length,
      catchingUp: status.syncInfo.catchingUp,
    }
  }

  /**
   * Get network status with DeshChain-specific metrics
   */
  public async getNetworkStatus(): Promise<NetworkStatus> {
    const chainInfo = await this.getChainInfo()
    const validators = await this.tmClient.validatorsAll()
    
    // Calculate network health metrics
    const activeValidators = validators.validators.filter(v => v.votingPower.toNumber() > 0).length
    const totalVotingPower = validators.validators.reduce(
      (sum, v) => sum + v.votingPower.toNumber(), 
      0
    )

    // Get recent blocks for TPS calculation
    const recentBlocks = await Promise.all([
      this.tmClient.block(chainInfo.blockHeight - 2),
      this.tmClient.block(chainInfo.blockHeight - 1),
      this.tmClient.block(chainInfo.blockHeight),
    ])

    const totalTxs = recentBlocks.reduce(
      (sum, block) => sum + block.block.data.txs.length, 
      0
    )
    const timeSpan = recentBlocks[2].block.header.time.getTime() - 
                   recentBlocks[0].block.header.time.getTime()
    const tps = Math.round((totalTxs / (timeSpan / 1000)) * 100) / 100

    return {
      ...chainInfo,
      activeValidators,
      totalVotingPower,
      tps,
      networkHealth: chainInfo.catchingUp ? 'syncing' : 'healthy',
      culturalEvents: await this.cultural.getActiveFestivals(),
    }
  }

  /**
   * Get account information
   */
  public async getAccount(address: string) {
    return this.client.getAccount(address)
  }

  /**
   * Get all balances for an address
   */
  public async getAllBalances(address: string) {
    return this.client.getAllBalances(address)
  }

  /**
   * Get balance for specific denomination
   */
  public async getBalance(address: string, denom: string) {
    return this.client.getBalance(address, denom)
  }

  /**
   * Get transaction by hash
   */
  public async getTx(hash: string) {
    return this.client.getTx(hash)
  }

  /**
   * Search transactions
   */
  public async searchTx(query: string) {
    return this.client.searchTx(query)
  }

  /**
   * Get block by height
   */
  public async getBlock(height?: number) {
    return this.tmClient.block(height)
  }

  /**
   * Get validators
   */
  public async getValidators(height?: number) {
    return this.tmClient.validatorsAll(height)
  }

  /**
   * Disconnect from the network
   */
  public disconnect() {
    this.client.disconnect()
    this.tmClient.disconnect()
  }

  /**
   * Check if connected to DeshChain network
   */
  public async isDeshChain(): Promise<boolean> {
    try {
      const chainInfo = await this.getChainInfo()
      return chainInfo.chainId.startsWith('deshchain')
    } catch {
      return false
    }
  }

  /**
   * Get current festival information
   */
  public async getCurrentFestival() {
    return this.cultural.getCurrentFestival()
  }

  /**
   * Get lending statistics across all modules
   */
  public async getLendingStats() {
    const [krishiStats, vyavasayaStats, shikshaMitraStats] = await Promise.all([
      this.lending.getKrishiMitraStats(),
      this.lending.getVyavasayaMitraStats(),
      this.lending.getShikshaMitraStats(),
    ])

    return {
      krishiMitra: krishiStats,
      vyavasayaMitra: vyavasayaStats,
      shikshamitra: shikshaMitraStats,
      combined: {
        totalDisbursed: [krishiStats, vyavasayaStats, shikshaMitraStats]
          .reduce((sum, stats) => sum + parseFloat(stats.totalDisbursed.replace(/[₹,\s]/g, '')), 0),
        avgInterestRate: [krishiStats, vyavasayaStats, shikshaMitraStats]
          .reduce((sum, stats) => sum + stats.averageRate, 0) / 3,
        totalBorrowers: [krishiStats, vyavasayaStats, shikshaMitraStats]
          .reduce((sum, stats) => sum + stats.activeLoans, 0),
      }
    }
  }

  /**
   * Search across all modules
   */
  public async search(query: string) {
    const results = await Promise.allSettled([
      this.searchTransactions(query),
      this.searchBlocks(query),
      this.searchAddresses(query),
      this.lending.searchLoans(query),
      this.sikkebaaz.searchTokens(query),
    ])

    return {
      transactions: results[0].status === 'fulfilled' ? results[0].value : [],
      blocks: results[1].status === 'fulfilled' ? results[1].value : [],
      addresses: results[2].status === 'fulfilled' ? results[2].value : [],
      loans: results[3].status === 'fulfilled' ? results[3].value : [],
      tokens: results[4].status === 'fulfilled' ? results[4].value : [],
    }
  }

  private async searchTransactions(query: string) {
    if (query.match(/^[A-F0-9]{64}$/i)) {
      const tx = await this.getTx(query.toUpperCase())
      return tx ? [tx] : []
    }
    return []
  }

  private async searchBlocks(query: string) {
    if (query.match(/^\d+$/)) {
      const block = await this.getBlock(parseInt(query))
      return block ? [block] : []
    }
    return []
  }

  private async searchAddresses(query: string) {
    if (query.startsWith('deshchain1') && query.length === 45) {
      const account = await this.getAccount(query)
      return account ? [{ address: query, account }] : []
    }
    return []
  }

  /**
   * Get node info
   */
  public async getNodeInfo() {
    const status = await this.tmClient.status()
    return status.nodeInfo
  }

  /**
   * Get genesis information
   */
  public async getGenesis() {
    return this.tmClient.genesis()
  }

  /**
   * Get consensus state
   */
  public async getConsensusState() {
    return this.tmClient.consensusState()
  }

  /**
   * Get ABCl info
   */
  public async getABCIInfo() {
    return this.tmClient.abciInfo()
  }

  /**
   * Health check
   */
  public async health() {
    return this.tmClient.health()
  }
}