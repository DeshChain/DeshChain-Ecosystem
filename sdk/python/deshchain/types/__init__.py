"""
DeshChain Python SDK Types

Type definitions for DeshChain SDK.
"""

from typing import Optional, List, Dict, Any, Union
from datetime import datetime
from dataclasses import dataclass
from pydantic import BaseModel

# Core Types
@dataclass
class ChainInfo:
    chain_id: str
    node_version: str
    block_height: int
    block_time: datetime
    validator_count: int
    catching_up: bool

@dataclass  
class NetworkStatus(ChainInfo):
    active_validators: int
    total_voting_power: int
    tps: float
    network_health: str
    cultural_events: List[Any]

@dataclass
class Account:
    address: str
    account_number: int
    sequence: int
    pub_key: Optional[Any] = None

@dataclass
class Balance:
    denom: str
    amount: str

@dataclass
class Transaction:
    txhash: str
    height: int
    timestamp: datetime
    fee: Dict[str, Any]
    memo: str
    messages: List[Any]
    events: List[Any]
    logs: List[Any]
    gas_used: str
    gas_wanted: str
    success: bool

@dataclass
class Block:
    height: int
    hash: str
    time: datetime
    proposer: str
    txs: List[str]
    num_txs: int
    last_commit: Any

# Festival Types
class Festival(BaseModel):
    festival_id: str
    name: str
    name_hindi: str
    description: str
    start_date: str
    end_date: str
    is_active: bool
    category: str
    regions: List[str]
    significance: str
    traditions: List[str]
    bonus_percentage: float
    cultural_impact: float

class CulturalQuote(BaseModel):
    quote_id: str
    text: str
    text_hindi: Optional[str] = None
    author: str
    author_hindi: Optional[str] = None
    category: str
    language: str
    source: str
    significance: str
    popularity: int
    tags: List[str]
    region: Optional[str] = None

# Lending Types  
class Loan(BaseModel):
    loan_id: str
    applicant_id: str
    amount: float
    interest_rate: float
    duration: int
    status: str
    disbursed_amount: Optional[float] = None
    repaid_amount: Optional[float] = None
    created_at: str
    due_date: str
    module: str

class LoanStats(BaseModel):
    total_loans: int
    total_disbursed: float
    average_rate: float
    active_loans: int
    default_rate: float

# Sikkebaaz Types
class LaunchpadToken(BaseModel):
    symbol: str
    name: str
    description: str
    creator: str
    total_supply: int
    current_supply: int
    launch_date: str
    status: str
    category: str
    logo_url: str
    cultural_theme: Optional[str] = None
    region: Optional[str] = None
    community_score: float

# Money Order Types
class MoneyOrder(BaseModel):
    order_id: str
    sender: str
    recipient: str
    amount: float
    currency: str
    status: str
    source_location: Dict[str, str]
    destination_location: Dict[str, str]
    estimated_delivery: str
    created_at: str
    tracking_id: str

# Governance Types
class Proposal(BaseModel):
    proposal_id: str
    title: str
    description: str
    proposer: str
    type: str
    status: str
    submit_time: str
    voting_start_time: str
    voting_end_time: str
    total_deposit: float

class Vote(BaseModel):
    proposal_id: str
    voter: str
    option: str
    weight: float
    timestamp: str