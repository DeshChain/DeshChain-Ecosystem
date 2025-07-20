/**
 * DeshChain JavaScript/TypeScript SDK
 * 
 * Official SDK for interacting with DeshChain blockchain
 * Features: Cultural heritage integration, lending modules, festival celebrations
 */

// Core client
export { DeshChainClient } from './client/DeshChainClient'
export { SigningDeshChainClient } from './client/SigningDeshChainClient'

// Module clients
export { LendingClient } from './modules/lending/LendingClient'
export { CulturalClient } from './modules/cultural/CulturalClient'
export { SikkebaazClient } from './modules/sikkebaaz/SikkebaazClient'
export { MoneyOrderClient } from './modules/moneyorder/MoneyOrderClient'
export { GovernanceClient } from './modules/governance/GovernanceClient'

// Types
export * from './types'
export * from './types/lending'
export * from './types/cultural'
export * from './types/sikkebaaz'
export * from './types/moneyorder'
export * from './types/governance'

// Utilities
export * from './utils/encoding'
export * from './utils/validation'
export * from './utils/cultural'
export * from './utils/festival'

// Constants
export * from './constants'

// Version
export const VERSION = '1.0.0'

// Default exports for convenience
export { DeshChainClient as default } from './client/DeshChainClient'