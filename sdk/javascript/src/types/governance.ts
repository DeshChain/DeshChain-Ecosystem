/**
 * Governance types for DeshChain
 */

export interface Proposal {
  proposalId: string
  title: string
  description: string
  proposer: string
  type: ProposalType
  status: ProposalStatus
  submitTime: string
  depositEndTime: string
  votingStartTime: string
  votingEndTime: string
  totalDeposit: number
  votes: VoteCounts
  tally: TallyResult
  metadata?: ProposalMetadata
}

export type ProposalType = 
  | 'text'
  | 'parameter_change'
  | 'software_upgrade'
  | 'community_pool_spend'
  | 'cancel_software_upgrade'
  | 'founder_protection'
  | 'emergency'

export type ProposalStatus = 
  | 'deposit_period'
  | 'voting_period'
  | 'passed'
  | 'rejected'
  | 'failed'
  | 'invalid'

export interface ProposalMetadata {
  category: string
  impact: 'low' | 'medium' | 'high' | 'critical'
  affectsFounder: boolean
  requiresVeto: boolean
  urgency: 'low' | 'medium' | 'high' | 'emergency'
  culturalRelevance?: string
}

export interface Vote {
  proposalId: string
  voter: string
  option: VoteOption
  weight: number
  timestamp: string
  reason?: string
}

export type VoteOption = 'yes' | 'no' | 'abstain' | 'no_with_veto'

export interface VoteCounts {
  yes: number
  no: number
  abstain: number
  noWithVeto: number
  total: number
}

export interface TallyResult extends VoteCounts {
  totalVotingPower: number
  turnout: number
  quorum: number
  threshold: number
  vetoThreshold: number
  passed: boolean
}

export interface VotingPower {
  totalPower: number
  delegatedPower: number
  ownPower: number
  bondedTokens: number
  delegations: DelegationSummary[]
}

export interface DelegationSummary {
  validatorAddress: string
  delegatedAmount: number
  votingPower: number
  validatorMoniker: string
}

export interface GovernanceStats {
  totalProposals: number
  activeProposals: number
  passedProposals: number
  rejectedProposals: number
  totalDeposits: number
  averageVotingPower: number
  participationRate: number
  proposalSuccessRate: number
}

export interface VetoVote {
  proposalId: string
  vetoer: string
  vetoType: VetoType
  reason: string
  timestamp: string
  validUntil: string
}

export type VetoType = 
  | 'founder_protection'
  | 'inheritance_protection'
  | 'revenue_protection'
  | 'emergency_veto'

export interface CommunityPool {
  balance: Array<{
    denom: string
    amount: number
  }>
  totalValue: number
  availableForSpending: number
}

export interface DelegationInfo {
  delegatorAddress: string
  validatorAddress: string
  shares: number
  balance: number
  validatorInfo: ValidatorInfo
}

export interface ValidatorInfo {
  operatorAddress: string
  consensusPubkey: string
  jailed: boolean
  status: ValidatorStatus
  tokens: number
  delegatorShares: number
  description: ValidatorDescription
  unbondingHeight: number
  unbondingTime: string
  commission: Commission
  minSelfDelegation: number
}

export type ValidatorStatus = 'bonded' | 'unbonded' | 'unbonding'

export interface ValidatorDescription {
  moniker: string
  identity: string
  website: string
  securityContact: string
  details: string
}

export interface Commission {
  commissionRates: {
    rate: number
    maxRate: number
    maxChangeRate: number
  }
  updateTime: string
}

export interface GovernanceParams {
  votingPeriod: number
  depositPeriod: number
  minDeposit: Array<{
    denom: string
    amount: number
  }>
  quorum: number
  threshold: number
  vetoThreshold: number
  maxDepositPeriod: number
}

export interface ProposalDeposit {
  proposalId: string
  depositor: string
  amount: number
  timestamp: string
}

export interface VoterParticipation {
  voter: string
  totalProposals: number
  votedProposals: number
  participationRate: number
  votingPower: number
  averageWeight: number
  streak: number
  lastVote: string
}

export interface GovernanceTimeline {
  proposalId: string
  events: TimelineEvent[]
}

export interface TimelineEvent {
  type: TimelineEventType
  timestamp: string
  description: string
  actor?: string
  data?: any
}

export type TimelineEventType = 
  | 'submitted'
  | 'deposit_received'
  | 'voting_started'
  | 'vote_cast'
  | 'veto_cast'
  | 'voting_ended'
  | 'proposal_passed'
  | 'proposal_rejected'
  | 'proposal_executed'

export interface FounderProtection {
  founderAddress: string
  protectionType: FounderProtectionType[]
  expiryDate?: string
  inherited: boolean
  inheritanceConditions?: string[]
  protectedParameters: string[]
  vetoRights: VetoRight[]
}

export type FounderProtectionType = 
  | 'token_allocation'
  | 'revenue_stream'
  | 'governance_veto'
  | 'parameter_protection'
  | 'inheritance_rights'

export interface VetoRight {
  rightType: string
  scope: string[]
  conditions: string[]
  expiryDate?: string
}

export interface RegionalGovernance {
  state: string
  district?: string
  activeProposals: number
  voterTurnout: number
  popularIssues: string[]
  culturalInfluence: number
  delegates: number
  proposalTypes: Record<ProposalType, number>
}

export interface CulturalGovernanceBonus {
  bonusType: 'participation_reward' | 'cultural_proposal' | 'festival_voting'
  percentage: number
  maxAmount: number
  conditions: string[]
  festival?: string
  region?: string
  validUntil: string
}

export interface StakingPool {
  notBondedTokens: number
  bondedTokens: number
  totalSupply: number
  bondedRatio: number
  inflation: number
  annualProvisions: number
}

export interface UnbondingDelegation {
  delegatorAddress: string
  validatorAddress: string
  entries: UnbondingEntry[]
}

export interface UnbondingEntry {
  creationHeight: number
  completionTime: string
  initialBalance: number
  balance: number
}

export interface Redelegation {
  delegatorAddress: string
  validatorSrcAddress: string
  validatorDstAddress: string
  entries: RedelegationEntry[]
}

export interface RedelegationEntry {
  creationHeight: number
  completionTime: string
  initialBalance: number
  sharesDst: number
}

export interface DistributionInfo {
  delegatorAddress: string
  validatorAddress: string
  reward: number
  commission?: number
  withdrawalAddress?: string
}

export interface SlashingInfo {
  validatorConsAddress: string
  infraction: InfractionType
  power: number
  slashFraction: number
  jailUntil: string
  tombstoned: boolean
}

export type InfractionType = 'double_sign' | 'downtime'