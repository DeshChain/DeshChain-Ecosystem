"""
DeshChain SDK Constants

Constants used throughout the DeshChain SDK.
"""

from typing import Dict, List, Tuple

# Chain Information
CHAIN_IDS = {
    "MAINNET": "deshchain-1",
    "TESTNET": "deshchain-testnet-1", 
    "DEVNET": "deshchain-devnet-1",
}

# Network Endpoints
ENDPOINTS = {
    CHAIN_IDS["MAINNET"]: {
        "rpc": "https://rpc.deshchain.network",
        "rest": "https://api.deshchain.network",
        "explorer": "https://explorer.deshchain.network",
    },
    CHAIN_IDS["TESTNET"]: {
        "rpc": "https://testnet-rpc.deshchain.network",
        "rest": "https://testnet-api.deshchain.network", 
        "explorer": "https://testnet-explorer.deshchain.network",
    },
    CHAIN_IDS["DEVNET"]: {
        "rpc": "http://localhost:26657",
        "rest": "http://localhost:1317",
        "explorer": "http://localhost:3000",
    },
}

# Denominations
DENOMINATIONS = {
    "NAMO": "unamo",
    "MICRO_NAMO": "unamo",
}

# Gas Configuration
GAS = {
    "DEFAULT_GAS_PRICE": "0.025unamo",
    "DEFAULT_GAS_LIMIT": 200000,
    "MAX_GAS_LIMIT": 10000000,
    "MIN_GAS_PRICE": "0.001unamo",
    "MAX_GAS_PRICE": "1.0unamo",
}

# Module Names
MODULES = {
    "BANK": "bank",
    "STAKING": "staking", 
    "GOVERNANCE": "gov",
    "DISTRIBUTION": "distribution",
    "SLASHING": "slashing",
    "NAMO": "namo",
    "CULTURAL": "cultural",
    "LENDING": "lending", 
    "KRISHI_MITRA": "krishimitra",
    "VYAVASAYA_MITRA": "vyavasayamitra",
    "SHIKSHA_MITRA": "shikshamitra",
    "SIKKEBAAZ": "sikkebaaz",
    "MONEY_ORDER": "moneyorder",
}

# Transaction Types
TX_TYPES = {
    "SEND": "/cosmos.bank.v1beta1.MsgSend",
    "DELEGATE": "/cosmos.staking.v1beta1.MsgDelegate",
    "UNDELEGATE": "/cosmos.staking.v1beta1.MsgUndelegate",
    "REDELEGATE": "/cosmos.staking.v1beta1.MsgBeginRedelegate",
    "WITHDRAW_REWARDS": "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
    "VOTE": "/cosmos.gov.v1beta1.MsgVote",
    "SUBMIT_PROPOSAL": "/cosmos.gov.v1beta1.MsgSubmitProposal",
    "DEPOSIT": "/cosmos.gov.v1beta1.MsgDeposit",
    
    # DeshChain specific
    "BURN_NAMO": "/deshchain.namo.v1.MsgBurnNAMO",
    "VEST_NAMO": "/deshchain.namo.v1.MsgVestNAMO",
    "APPLY_LOAN": "/deshchain.krishimitra.v1.MsgApplyLoan",
    "LAUNCH_TOKEN": "/deshchain.sikkebaaz.v1.MsgLaunchToken",
    "CREATE_MONEY_ORDER": "/deshchain.moneyorder.v1.MsgCreateMoneyOrder",
}

# API Endpoints
API_ENDPOINTS = {
    "BLOCKS": "/cosmos/base/tendermint/v1beta1/blocks",
    "VALIDATORS": "/cosmos/staking/v1beta1/validators",
    "PROPOSALS": "/cosmos/gov/v1beta1/proposals",
    
    # DeshChain specific
    "CULTURAL_FESTIVALS": "/deshchain/cultural/v1/festivals",
    "CULTURAL_QUOTES": "/deshchain/cultural/v1/quotes",
    "LENDING_STATS": "/deshchain/lending/v1/stats",
    "SIKKEBAAZ_TOKENS": "/deshchain/sikkebaaz/v1/tokens",
    "MONEY_ORDER_STATS": "/deshchain/moneyorder/v1/stats",
}

# Error Codes
ERROR_CODES = {
    "NETWORK_ERROR": "NETWORK_ERROR",
    "TRANSACTION_ERROR": "TRANSACTION_ERROR",
    "VALIDATION_ERROR": "VALIDATION_ERROR",
    "INSUFFICIENT_FUNDS": "INSUFFICIENT_FUNDS",
    "INVALID_ADDRESS": "INVALID_ADDRESS", 
    "INVALID_AMOUNT": "INVALID_AMOUNT",
    "GAS_ESTIMATION_FAILED": "GAS_ESTIMATION_FAILED",
    "BROADCAST_FAILED": "BROADCAST_FAILED",
    "TIMEOUT": "TIMEOUT",
}

# Cultural Constants
CULTURAL = {
    "SUPPORTED_LANGUAGES": [
        "en", "hi", "bn", "te", "ta", "mr", "gu", "kn", "ml", "or",
        "pa", "as", "ur", "ne", "sd", "kok", "brx", "doi", "mni", "sat",
        "bho", "mai"
    ],
    
    "FESTIVAL_CATEGORIES": [
        "religious", "harvest", "seasonal", "national", "regional", "cultural", "historical"
    ],
    
    "QUOTE_CATEGORIES": [
        "wisdom", "patriotism", "philosophy", "spirituality", "motivation",
        "peace", "unity", "progress", "culture", "education"
    ],
    
    "INDIAN_STATES": [
        "Andhra Pradesh", "Arunachal Pradesh", "Assam", "Bihar", "Chhattisgarh",
        "Goa", "Gujarat", "Haryana", "Himachal Pradesh", "Jammu and Kashmir",
        "Jharkhand", "Karnataka", "Kerala", "Madhya Pradesh", "Maharashtra",
        "Manipur", "Meghalaya", "Mizoram", "Nagaland", "Odisha", "Punjab",
        "Rajasthan", "Sikkim", "Tamil Nadu", "Telangana", "Tripura",
        "Uttar Pradesh", "Uttarakhand", "West Bengal", "Delhi"
    ],
}

# Lending Constants
LENDING = {
    "INTEREST_RATES": {
        "KRISHI_MITRA": {"MIN": 6, "MAX": 9},  # Agriculture: 6-9%
        "VYAVASAYA_MITRA": {"MIN": 8, "MAX": 12},  # Business: 8-12%
        "SHIKSHA_MITRA": {"MIN": 4, "MAX": 7},  # Education: 4-7%
    },
    
    "LOAN_STATUSES": [
        "pending", "approved", "disbursed", "active", "completed", "defaulted", "rejected"
    ],
    
    "KYC_STATUSES": ["pending", "verified", "rejected", "expired"],
    
    "CREDIT_SCORE_RANGES": {
        "EXCELLENT": {"MIN": 750, "MAX": 900},
        "GOOD": {"MIN": 650, "MAX": 749},
        "FAIR": {"MIN": 550, "MAX": 649},
        "POOR": {"MIN": 300, "MAX": 549},
    },
}

# Sikkebaaz Constants
SIKKEBAAZ = {
    "TOKEN_CATEGORIES": [
        "cultural", "regional", "festival", "meme", "utility",
        "charity", "gaming", "art", "music", "sports"
    ],
    
    "TOKEN_STATUSES": [
        "pending", "review", "approved", "launched", "trading", "graduated", "failed", "vetoed"
    ],
    
    "ANTI_PUMP_LIMITS": {
        "MAX_WALLET_PERCENT": 5,  # 5% max wallet size
        "MAX_TRANSACTION_PERCENT": 1,  # 1% max transaction size 
        "TRADING_DELAY": 60,  # 60 seconds between trades
        "LIQUIDITY_LOCK_DURATION": 365,  # 365 days liquidity lock
    },
}

# Money Order Constants
MONEY_ORDER = {
    "STATUSES": [
        "created", "pending", "processing", "in_transit", "delivered", "failed", "cancelled", "refunded"
    ],
    
    "ORDER_TYPES": ["market", "limit", "stop"],
    
    "TRADING_SIDES": ["buy", "sell"],
    
    "FEE_STRUCTURE": {
        "BASE_FEE": 0.1,  # 0.1% base fee
        "DISTANCE_FEE_PER_KM": 0.001,  # 0.001% per km
        "URGENCY_MULTIPLIER": 1.5,  # 1.5x for urgent delivery
    },
}

# Governance Constants
GOVERNANCE = {
    "PROPOSAL_TYPES": [
        "text", "parameter_change", "software_upgrade", "community_pool_spend",
        "cancel_software_upgrade", "founder_protection", "emergency"
    ],
    
    "PROPOSAL_STATUSES": [
        "deposit_period", "voting_period", "passed", "rejected", "failed", "invalid"
    ],
    
    "VOTE_OPTIONS": ["yes", "no", "abstain", "no_with_veto"],
    
    "VETO_TYPES": [
        "founder_protection", "inheritance_protection", "revenue_protection", "emergency_veto"
    ],
    
    "VALIDATOR_STATUSES": ["bonded", "unbonded", "unbonding"],
}

# Time Constants
TIME = {
    "BLOCK_TIME": 6,  # 6 seconds average block time
    "UNBONDING_PERIOD": 21 * 24 * 60 * 60,  # 21 days in seconds
    "VOTING_PERIOD": 14 * 24 * 60 * 60,  # 14 days in seconds
    "DEPOSIT_PERIOD": 7 * 24 * 60 * 60,  # 7 days in seconds
}

# Precision Constants
PRECISION = {
    "NAMO_DECIMALS": 6,
    "PERCENTAGE_DECIMALS": 2,
    "PRICE_DECIMALS": 6,
    "AMOUNT_DECIMALS": 6,
}

# Limits
LIMITS = {
    "MAX_MEMO_LENGTH": 512,
    "MAX_VALIDATORS_PER_DELEGATOR": 100,
    "MAX_PROPOSAL_TITLE_LENGTH": 140,
    "MAX_PROPOSAL_DESCRIPTION_LENGTH": 10000,
    "MIN_DEPOSIT_AMOUNT": 1000000,  # 1 NAMO in micro units
    "MAX_DEPOSIT_AMOUNT": 1000000000000,  # 1M NAMO in micro units
}

# Revenue Distribution (Platform Model)
REVENUE_DISTRIBUTION = {
    "DEVELOPMENT": 0.30,  # 30% to development fund
    "COMMUNITY": 0.25,  # 25% to community rewards
    "LIQUIDITY": 0.20,  # 20% to liquidity provision
    "NGO": 0.10,  # 10% to NGO donations
    "EMERGENCY": 0.10,  # 10% to emergency fund
    "FOUNDER": 0.05,  # 5% to founder (reduced from 20%)
}

# Cultural Impact Scoring
CULTURAL_SCORING = {
    "WEIGHTS": {
        "FESTIVAL_IMPORTANCE": 0.3,
        "REGIONAL_PARTICIPATION": 0.3,
        "TRADITIONAL_VALUE": 0.2,
        "MODERN_RELEVANCE": 0.2,
    },
    
    "BONUS_MULTIPLIERS": {
        "NATIONAL_FESTIVAL": 1.5,
        "REGIONAL_FESTIVAL": 1.3,
        "LOCAL_CELEBRATION": 1.1,
        "CULTURAL_EVENT": 1.2,
    },
}

# HTTP Configuration
HTTP = {
    "DEFAULT_TIMEOUT": 30.0,
    "MAX_RETRIES": 3,
    "RETRY_BACKOFF": 2.0,
    "REQUEST_HEADERS": {
        "User-Agent": "DeshChain-Python-SDK/1.0.0",
        "Content-Type": "application/json",
        "Accept": "application/json",
    },
}

# WebSocket Configuration
WEBSOCKET = {
    "DEFAULT_TIMEOUT": 10.0,
    "PING_INTERVAL": 20.0,
    "PING_TIMEOUT": 10.0,
    "CLOSE_TIMEOUT": 10.0,
}

# Default Configuration
DEFAULT_CONFIG = {
    "chain_id": CHAIN_IDS["MAINNET"],
    "rpc_url": ENDPOINTS[CHAIN_IDS["MAINNET"]]["rpc"],
    "rest_url": ENDPOINTS[CHAIN_IDS["MAINNET"]]["rest"],
    "gas_price": GAS["DEFAULT_GAS_PRICE"],
    "gas_limit": GAS["DEFAULT_GAS_LIMIT"],
    "timeout": HTTP["DEFAULT_TIMEOUT"],
    "retries": HTTP["MAX_RETRIES"],
}