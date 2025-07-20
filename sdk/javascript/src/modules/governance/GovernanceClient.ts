import { StargateClient } from '@cosmjs/stargate'
import { Tendermint34Client } from '@cosmjs/tendermint-rpc'
import { 
  Proposal, 
  Vote, 
  GovernanceStats,
  VotingPower,
  ProposalStatus,
  TallyResult,
  GovernanceParams,
  VetoVote,
  CommunityPool,
  DelegationInfo
} from '../../types/governance'

/**
 * Client for interacting with DeshChain governance system
 * Includes founder protection and community governance features
 */
export class GovernanceClient {
  constructor(
    private readonly client: StargateClient,
    private readonly tmClient: Tendermint34Client
  ) {}

  /**
   * Proposal Management
   */

  /**
   * Get all proposals
   */
  async getProposals(
    status?: ProposalStatus,
    limit: number = 50,
    offset: number = 0
  ): Promise<Proposal[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_proposals: { status, limit, offset } }
    )
    return response.proposals || []
  }

  /**
   * Get proposal by ID
   */
  async getProposal(proposalId: string): Promise<Proposal | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_proposal: { proposal_id: proposalId } }
      )
      return response.proposal
    } catch {
      return null
    }
  }

  /**
   * Get active proposals
   */
  async getActiveProposals(): Promise<Proposal[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_active_proposals: {} }
    )
    return response.proposals || []
  }

  /**
   * Get proposals by proposer
   */
  async getProposalsByProposer(proposerAddress: string): Promise<Proposal[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_proposals_by_proposer: { proposer: proposerAddress } }
    )
    return response.proposals || []
  }

  /**
   * Get proposal tally
   */
  async getProposalTally(proposalId: string): Promise<TallyResult> {
    const response = await this.client.queryContractSmart(
      '',
      { get_proposal_tally: { proposal_id: proposalId } }
    )
    return {
      yes: response.yes,
      no: response.no,
      abstain: response.abstain,
      noWithVeto: response.no_with_veto,
      totalVotingPower: response.total_voting_power,
      turnout: response.turnout,
    }
  }

  /**
   * Voting System
   */

  /**
   * Get vote by voter and proposal
   */
  async getVote(proposalId: string, voterAddress: string): Promise<Vote | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_vote: { proposal_id: proposalId, voter: voterAddress } }
      )
      return response.vote
    } catch {
      return null
    }
  }

  /**
   * Get all votes for a proposal
   */
  async getVotes(proposalId: string, limit: number = 100): Promise<Vote[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_votes: { proposal_id: proposalId, limit } }
    )
    return response.votes || []
  }

  /**
   * Get votes by voter
   */
  async getVotesByVoter(voterAddress: string): Promise<Vote[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_votes_by_voter: { voter: voterAddress } }
    )
    return response.votes || []
  }

  /**
   * Get voting power for address
   */
  async getVotingPower(address: string): Promise<VotingPower> {
    const response = await this.client.queryContractSmart(
      '',
      { get_voting_power: { address } }
    )
    return {
      totalPower: response.total_power,
      delegatedPower: response.delegated_power,
      ownPower: response.own_power,
      bondedTokens: response.bonded_tokens,
      delegations: response.delegations || [],
    }
  }

  /**
   * Check if address can vote on proposal
   */
  async canVote(proposalId: string, voterAddress: string): Promise<{
    canVote: boolean
    reason?: string
    votingPower: number
  }> {
    const response = await this.client.queryContractSmart(
      '',
      { can_vote: { proposal_id: proposalId, voter: voterAddress } }
    )
    return response
  }

  /**
   * Founder Protection System
   */

  /**
   * Check if proposal affects founder protection
   */
  async checkFounderProtection(proposalId: string): Promise<{
    affectsFounder: boolean
    protectionType: string[]
    requiresVeto: boolean
    vetoDeadline?: string
  }> {
    const response = await this.client.queryContractSmart(
      '',
      { check_founder_protection: { proposal_id: proposalId } }
    )
    return response
  }

  /**
   * Get founder veto votes
   */
  async getFounderVetoVotes(proposalId: string): Promise<VetoVote[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_founder_veto_votes: { proposal_id: proposalId } }
    )
    return response.veto_votes || []
  }

  /**
   * Check founder veto power
   */
  async getFounderVetoPower(address: string): Promise<{
    hasVetoPower: boolean
    vetoType: string
    expiryDate?: string
    inherited: boolean
  }> {
    const response = await this.client.queryContractSmart(
      '',
      { get_founder_veto_power: { address } }
    )
    return response
  }

  /**
   * Get protected parameters
   */
  async getProtectedParameters(): Promise<string[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_protected_parameters: {} }
    )
    return response.parameters || []
  }

  /**
   * Delegation System
   */

  /**
   * Get delegations by delegator
   */
  async getDelegations(delegatorAddress: string): Promise<DelegationInfo[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_delegations: { delegator: delegatorAddress } }
    )
    return response.delegations || []
  }

  /**
   * Get delegations to validator
   */
  async getDelegationsToValidator(validatorAddress: string): Promise<DelegationInfo[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_delegations_to_validator: { validator: validatorAddress } }
    )
    return response.delegations || []
  }

  /**
   * Get delegation info
   */
  async getDelegation(
    delegatorAddress: string, 
    validatorAddress: string
  ): Promise<DelegationInfo | null> {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_delegation: { delegator: delegatorAddress, validator: validatorAddress } }
      )
      return response.delegation
    } catch {
      return null
    }
  }

  /**
   * Get unbonding delegations
   */
  async getUnbondingDelegations(delegatorAddress: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_unbonding_delegations: { delegator: delegatorAddress } }
    )
    return response.unbonding_delegations || []
  }

  /**
   * Get redelegations
   */
  async getRedelegations(delegatorAddress: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_redelegations: { delegator: delegatorAddress } }
    )
    return response.redelegations || []
  }

  /**
   * Validator Information
   */

  /**
   * Get all validators
   */
  async getValidators(status: 'bonded' | 'unbonded' | 'unbonding' | 'all' = 'all') {
    const response = await this.client.queryContractSmart(
      '',
      { get_validators: { status } }
    )
    return response.validators || []
  }

  /**
   * Get validator by address
   */
  async getValidator(validatorAddress: string) {
    try {
      const response = await this.client.queryContractSmart(
        '',
        { get_validator: { validator: validatorAddress } }
      )
      return response.validator
    } catch {
      return null
    }
  }

  /**
   * Get validator delegations
   */
  async getValidatorDelegations(validatorAddress: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_validator_delegations: { validator: validatorAddress } }
    )
    return response.delegations || []
  }

  /**
   * Get staking pool info
   */
  async getStakingPool() {
    const response = await this.client.queryContractSmart(
      '',
      { get_staking_pool: {} }
    )
    return response.pool
  }

  /**
   * Community Pool
   */

  /**
   * Get community pool balance
   */
  async getCommunityPool(): Promise<CommunityPool> {
    const response = await this.client.queryContractSmart(
      '',
      { get_community_pool: {} }
    )
    return {
      balance: response.balance || [],
      totalValue: response.total_value,
      availableForSpending: response.available_for_spending,
    }
  }

  /**
   * Get community pool spend proposals
   */
  async getCommunityPoolSpendProposals(): Promise<Proposal[]> {
    const response = await this.client.queryContractSmart(
      '',
      { get_community_pool_spend_proposals: {} }
    )
    return response.proposals || []
  }

  /**
   * Parameters and Configuration
   */

  /**
   * Get governance parameters
   */
  async getGovernanceParams(): Promise<GovernanceParams> {
    const response = await this.client.queryContractSmart(
      '',
      { get_governance_params: {} }
    )
    return {
      votingPeriod: response.voting_period,
      depositPeriod: response.deposit_period,
      minDeposit: response.min_deposit,
      quorum: response.quorum,
      threshold: response.threshold,
      vetoThreshold: response.veto_threshold,
      maxDepositPeriod: response.max_deposit_period,
    }
  }

  /**
   * Get staking parameters
   */
  async getStakingParams() {
    const response = await this.client.queryContractSmart(
      '',
      { get_staking_params: {} }
    )
    return response.params
  }

  /**
   * Get distribution parameters
   */
  async getDistributionParams() {
    const response = await this.client.queryContractSmart(
      '',
      { get_distribution_params: {} }
    )
    return response.params
  }

  /**
   * Statistics and Analytics
   */

  /**
   * Get governance statistics
   */
  async getGovernanceStats(): Promise<GovernanceStats> {
    const response = await this.client.queryContractSmart(
      '',
      { get_governance_stats: {} }
    )
    return {
      totalProposals: response.total_proposals,
      activeProposals: response.active_proposals,
      passedProposals: response.passed_proposals,
      rejectedProposals: response.rejected_proposals,
      totalDeposits: response.total_deposits,
      averageVotingPower: response.average_voting_power,
      participationRate: response.participation_rate,
      proposalSuccessRate: response.proposal_success_rate,
    }
  }

  /**
   * Get voting statistics for proposal
   */
  async getVotingStats(proposalId: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_voting_stats: { proposal_id: proposalId } }
    )
    return response.stats
  }

  /**
   * Get voter participation history
   */
  async getVoterParticipation(voterAddress: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_voter_participation: { voter: voterAddress } }
    )
    return response.participation
  }

  /**
   * Get top voters by participation
   */
  async getTopVoters(limit: number = 10) {
    const response = await this.client.queryContractSmart(
      '',
      { get_top_voters: { limit } }
    )
    return response.voters || []
  }

  /**
   * Cultural and Regional Features
   */

  /**
   * Get regional voting patterns
   */
  async getRegionalVotingPatterns(state?: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_regional_voting_patterns: { state } }
    )
    return response.patterns
  }

  /**
   * Get cultural governance incentives
   */
  async getCulturalGovernanceIncentives() {
    const response = await this.client.queryContractSmart(
      '',
      { get_cultural_governance_incentives: {} }
    )
    return response.incentives || []
  }

  /**
   * Get festival governance bonuses
   */
  async getFestivalGovernanceBonuses() {
    const response = await this.client.queryContractSmart(
      '',
      { get_festival_governance_bonuses: {} }
    )
    return response.bonuses || []
  }

  /**
   * Utility Functions
   */

  /**
   * Search proposals
   */
  async searchProposals(query: string): Promise<Proposal[]> {
    const response = await this.client.queryContractSmart(
      '',
      { search_proposals: { query } }
    )
    return response.proposals || []
  }

  /**
   * Get proposal deposit info
   */
  async getProposalDeposits(proposalId: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_proposal_deposits: { proposal_id: proposalId } }
    )
    return response.deposits || []
  }

  /**
   * Check if address is founder
   */
  async isFounder(address: string): Promise<boolean> {
    const response = await this.client.queryContractSmart(
      '',
      { is_founder: { address } }
    )
    return response.is_founder
  }

  /**
   * Get governance timeline
   */
  async getGovernanceTimeline(proposalId: string) {
    const response = await this.client.queryContractSmart(
      '',
      { get_governance_timeline: { proposal_id: proposalId } }
    )
    return response.timeline || []
  }

  /**
   * Health check
   */
  async healthCheck() {
    const response = await this.client.queryContractSmart(
      '',
      { health_check: {} }
    )
    return {
      healthy: response.healthy,
      governanceStatus: response.governance_status,
      stakingStatus: response.staking_status,
      lastUpdated: response.last_updated,
    }
  }
}