"""
DeshChain Client

Main client for interacting with DeshChain blockchain.
"""

import asyncio
from typing import Optional, Dict, Any, List, Union
from datetime import datetime

import httpx
import requests
from cosmpy.clients.tendermint import TendermintClient

from .constants import ENDPOINTS, CHAIN_IDS, DEFAULT_CONFIG
from .exceptions import NetworkError, ConfigurationError
from .types import ChainInfo, NetworkStatus, Account, Balance, Transaction, Block
from .modules.cultural import CulturalClient, AsyncCulturalClient
from .modules.lending import LendingClient, AsyncLendingClient  
from .modules.sikkebaaz import SikkebaazClient, AsyncSikkebaazClient
from .modules.moneyorder import MoneyOrderClient, AsyncMoneyOrderClient
from .modules.governance import GovernanceClient, AsyncGovernanceClient


class DeshChainClient:
    """
    Synchronous DeshChain client for read-only operations.
    
    Provides access to all DeshChain modules and blockchain data.
    """
    
    def __init__(
        self,
        rpc_url: str,
        rest_url: Optional[str] = None,
        chain_id: Optional[str] = None,
        timeout: float = 30.0,
        retries: int = 3,
        **kwargs
    ):
        """
        Initialize DeshChain client.
        
        Args:
            rpc_url: Tendermint RPC endpoint
            rest_url: REST API endpoint (optional)
            chain_id: Chain ID (optional, will be detected)
            timeout: Request timeout in seconds
            retries: Number of retry attempts
        """
        self.rpc_url = rpc_url
        self.rest_url = rest_url or rpc_url.replace(":26657", ":1317")
        self._chain_id = chain_id
        self.timeout = timeout
        self.retries = retries
        
        # Initialize HTTP session
        self.session = requests.Session()
        self.session.timeout = timeout
        
        # Initialize Tendermint client
        self.tm_client = TendermintClient(rpc_url)
        
        # Initialize module clients
        self.cultural = CulturalClient(self)
        self.lending = LendingClient(self)
        self.sikkebaaz = SikkebaazClient(self)
        self.money_order = MoneyOrderClient(self)
        self.governance = GovernanceClient(self)
    
    @classmethod
    def connect(
        cls,
        network: str = "mainnet",
        **kwargs
    ) -> "DeshChainClient":
        """
        Connect to a DeshChain network.
        
        Args:
            network: Network name ('mainnet', 'testnet', 'devnet')
            **kwargs: Additional client options
            
        Returns:
            DeshChainClient instance
        """
        network_map = {
            "mainnet": CHAIN_IDS["MAINNET"],
            "testnet": CHAIN_IDS["TESTNET"], 
            "devnet": CHAIN_IDS["DEVNET"],
        }
        
        if network not in network_map:
            raise ConfigurationError(f"Unknown network: {network}")
        
        chain_id = network_map[network]
        endpoints = ENDPOINTS[chain_id]
        
        return cls(
            rpc_url=endpoints["rpc"],
            rest_url=endpoints["rest"],
            chain_id=chain_id,
            **kwargs
        )
    
    def get_chain_info(self) -> ChainInfo:
        """Get chain information."""
        try:
            status = self.tm_client.status()
            validators = self.tm_client.validators()
            
            return ChainInfo(
                chain_id=status.node_info.network,
                node_version=status.node_info.version,
                block_height=status.sync_info.latest_block_height,
                block_time=status.sync_info.latest_block_time,
                validator_count=len(validators.validators),
                catching_up=status.sync_info.catching_up,
            )
        except Exception as e:
            raise NetworkError(f"Failed to get chain info: {e}")
    
    def get_network_status(self) -> NetworkStatus:
        """Get network status with DeshChain-specific metrics."""
        try:
            chain_info = self.get_chain_info()
            validators = self.tm_client.validators()
            
            # Calculate network health metrics
            active_validators = len([v for v in validators.validators if int(v.voting_power) > 0])
            total_voting_power = sum(int(v.voting_power) for v in validators.validators)
            
            # Get recent blocks for TPS calculation  
            recent_blocks = []
            for i in range(3):
                height = chain_info.block_height - i
                block = self.tm_client.block(height)
                recent_blocks.append(block)
            
            total_txs = sum(len(block.block.data.txs) for block in recent_blocks)
            if len(recent_blocks) >= 2:
                time_span = (recent_blocks[0].block.header.time - recent_blocks[-1].block.header.time).total_seconds()
                tps = total_txs / max(time_span, 1) if time_span > 0 else 0
            else:
                tps = 0
            
            # Get cultural events
            cultural_events = self.cultural.get_active_festivals()
            
            return NetworkStatus(
                chain_id=chain_info.chain_id,
                node_version=chain_info.node_version,
                block_height=chain_info.block_height,
                block_time=chain_info.block_time,
                validator_count=chain_info.validator_count,
                catching_up=chain_info.catching_up,
                active_validators=active_validators,
                total_voting_power=total_voting_power,
                tps=round(tps, 2),
                network_health="syncing" if chain_info.catching_up else "healthy",
                cultural_events=cultural_events,
            )
        except Exception as e:
            raise NetworkError(f"Failed to get network status: {e}")
    
    def get_account(self, address: str) -> Optional[Account]:
        """Get account information."""
        try:
            endpoint = f"/cosmos/auth/v1beta1/accounts/{address}"
            response = self._get(endpoint)
            
            if "account" not in response:
                return None
                
            account_data = response["account"]
            return Account(
                address=account_data["address"],
                account_number=int(account_data["account_number"]),
                sequence=int(account_data["sequence"]),
                pub_key=account_data.get("pub_key"),
            )
        except Exception as e:
            raise NetworkError(f"Failed to get account: {e}")
    
    def get_all_balances(self, address: str) -> List[Balance]:
        """Get all balances for an address."""
        try:
            endpoint = f"/cosmos/bank/v1beta1/balances/{address}"
            response = self._get(endpoint)
            
            balances = []
            for balance_data in response.get("balances", []):
                balances.append(Balance(
                    denom=balance_data["denom"],
                    amount=balance_data["amount"],
                ))
            
            return balances
        except Exception as e:
            raise NetworkError(f"Failed to get balances: {e}")
    
    def get_balance(self, address: str, denom: str) -> Optional[Balance]:
        """Get balance for specific denomination."""
        try:
            endpoint = f"/cosmos/bank/v1beta1/balances/{address}/by_denom"
            params = {"denom": denom}
            response = self._get(endpoint, params=params)
            
            if "balance" not in response:
                return None
                
            balance_data = response["balance"]
            return Balance(
                denom=balance_data["denom"],
                amount=balance_data["amount"],
            )
        except Exception as e:
            raise NetworkError(f"Failed to get balance: {e}")
    
    def get_transaction(self, tx_hash: str) -> Optional[Transaction]:
        """Get transaction by hash."""
        try:
            tx = self.tm_client.tx(tx_hash)
            if not tx:
                return None
                
            return Transaction(
                txhash=tx_hash,
                height=tx.height,
                timestamp=tx.timestamp,
                fee=tx.fee,
                memo=tx.memo,
                messages=tx.messages,
                events=tx.events,
                logs=tx.logs,
                gas_used=tx.gas_used,
                gas_wanted=tx.gas_wanted,
                success=tx.code == 0,
            )
        except Exception as e:
            raise NetworkError(f"Failed to get transaction: {e}")
    
    def get_block(self, height: Optional[int] = None) -> Block:
        """Get block by height."""
        try:
            block = self.tm_client.block(height)
            
            return Block(
                height=block.block.header.height,
                hash=block.block_id.hash,
                time=block.block.header.time,
                proposer=block.block.header.proposer_address,
                txs=block.block.data.txs,
                num_txs=len(block.block.data.txs),
                last_commit=block.block.last_commit,
            )
        except Exception as e:
            raise NetworkError(f"Failed to get block: {e}")
    
    def search_transactions(
        self, 
        query: str,
        page: int = 1,
        limit: int = 30
    ) -> List[Transaction]:
        """Search transactions."""
        try:
            results = self.tm_client.tx_search(query, page=page, per_page=limit)
            
            transactions = []
            for tx_result in results.txs:
                transactions.append(Transaction(
                    txhash=tx_result.hash,
                    height=tx_result.height,
                    timestamp=tx_result.timestamp,
                    fee=tx_result.fee,
                    memo=tx_result.memo,
                    messages=tx_result.messages,
                    events=tx_result.events,
                    logs=tx_result.logs,
                    gas_used=tx_result.gas_used,
                    gas_wanted=tx_result.gas_wanted,
                    success=tx_result.code == 0,
                ))
            
            return transactions
        except Exception as e:
            raise NetworkError(f"Failed to search transactions: {e}")
    
    def get_current_festival(self):
        """Get currently active festival."""
        return self.cultural.get_current_festival()
    
    def get_lending_stats(self):
        """Get lending statistics across all modules."""
        krishi_stats = self.lending.get_krishi_mitra_stats()
        vyavasaya_stats = self.lending.get_vyavasaya_mitra_stats()
        shiksha_stats = self.lending.get_shiksha_mitra_stats()
        
        return {
            "krishi_mitra": krishi_stats,
            "vyavasaya_mitra": vyavasaya_stats,
            "shiksha_mitra": shiksha_stats,
            "combined": {
                "total_disbursed": (
                    krishi_stats.total_disbursed + 
                    vyavasaya_stats.total_disbursed + 
                    shiksha_stats.total_disbursed
                ),
                "avg_interest_rate": (
                    krishi_stats.average_rate + 
                    vyavasaya_stats.average_rate + 
                    shiksha_stats.average_rate
                ) / 3,
                "total_borrowers": (
                    krishi_stats.active_loans + 
                    vyavasaya_stats.active_loans + 
                    shiksha_stats.active_loans
                ),
            }
        }
    
    def search(self, query: str) -> Dict[str, Any]:
        """Search across all modules."""
        results = {
            "transactions": [],
            "blocks": [],
            "addresses": [],
            "loans": [],
            "tokens": [],
        }
        
        try:
            # Search transactions
            if len(query) == 64:  # Transaction hash
                tx = self.get_transaction(query.upper())
                if tx:
                    results["transactions"] = [tx]
            
            # Search blocks
            if query.isdigit():
                try:
                    block = self.get_block(int(query))
                    results["blocks"] = [block]
                except:
                    pass
            
            # Search addresses
            if query.startswith('deshchain1') and len(query) == 45:
                account = self.get_account(query)
                if account:
                    results["addresses"] = [{"address": query, "account": account}]
            
            # Search loans
            try:
                loans = self.lending.search_loans(query)
                results["loans"] = loans
            except:
                pass
            
            # Search tokens
            try:
                tokens = self.sikkebaaz.search_tokens(query)
                results["tokens"] = tokens
            except:
                pass
                
        except Exception:
            pass  # Ignore search errors for individual modules
        
        return results
    
    def is_deshchain(self) -> bool:
        """Check if connected to DeshChain network."""
        try:
            chain_info = self.get_chain_info()
            return chain_info.chain_id.startswith("deshchain")
        except:
            return False
    
    def close(self):
        """Close client connections."""
        self.session.close()
        if hasattr(self.tm_client, 'close'):
            self.tm_client.close()
    
    def _get(self, endpoint: str, params: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Make GET request to REST API."""
        url = f"{self.rest_url}{endpoint}"
        
        try:
            response = self.session.get(url, params=params, timeout=self.timeout)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            raise NetworkError(f"Request failed: {e}", endpoint=endpoint)
    
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()


class AsyncDeshChainClient:
    """
    Asynchronous DeshChain client for read-only operations.
    
    Provides async access to all DeshChain modules and blockchain data.
    """
    
    def __init__(
        self,
        rpc_url: str,
        rest_url: Optional[str] = None,
        chain_id: Optional[str] = None,
        timeout: float = 30.0,
        retries: int = 3,
        **kwargs
    ):
        """Initialize async DeshChain client."""
        self.rpc_url = rpc_url
        self.rest_url = rest_url or rpc_url.replace(":26657", ":1317")
        self._chain_id = chain_id
        self.timeout = timeout
        self.retries = retries
        
        # Will be initialized in async context
        self._client: Optional[httpx.AsyncClient] = None
        
        # Initialize module clients
        self.cultural = AsyncCulturalClient(self)
        self.lending = AsyncLendingClient(self)
        self.sikkebaaz = AsyncSikkebaazClient(self)
        self.money_order = AsyncMoneyOrderClient(self)
        self.governance = AsyncGovernanceClient(self)
    
    @classmethod
    def connect(
        cls,
        network: str = "mainnet",
        **kwargs
    ) -> "AsyncDeshChainClient":
        """Connect to a DeshChain network."""
        network_map = {
            "mainnet": CHAIN_IDS["MAINNET"],
            "testnet": CHAIN_IDS["TESTNET"],
            "devnet": CHAIN_IDS["DEVNET"],
        }
        
        if network not in network_map:
            raise ConfigurationError(f"Unknown network: {network}")
        
        chain_id = network_map[network]
        endpoints = ENDPOINTS[chain_id]
        
        return cls(
            rpc_url=endpoints["rpc"],
            rest_url=endpoints["rest"],
            chain_id=chain_id,
            **kwargs
        )
    
    async def __aenter__(self):
        """Async context manager entry."""
        self._client = httpx.AsyncClient(timeout=self.timeout)
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit."""
        if self._client:
            await self._client.aclose()
    
    @property
    def client(self) -> httpx.AsyncClient:
        """Get HTTP client, creating if necessary."""
        if self._client is None:
            self._client = httpx.AsyncClient(timeout=self.timeout)
        return self._client
    
    async def get_chain_info(self) -> ChainInfo:
        """Get chain information."""
        try:
            # For async implementation, we'd need an async Tendermint client
            # This is a simplified version
            endpoint = "/status"
            response = await self._get(endpoint)
            
            result = response["result"]
            node_info = result["node_info"]
            sync_info = result["sync_info"]
            
            return ChainInfo(
                chain_id=node_info["network"],
                node_version=node_info["version"],
                block_height=int(sync_info["latest_block_height"]),
                block_time=datetime.fromisoformat(sync_info["latest_block_time"].replace("Z", "+00:00")),
                validator_count=0,  # Would need separate call
                catching_up=sync_info["catching_up"],
            )
        except Exception as e:
            raise NetworkError(f"Failed to get chain info: {e}")
    
    async def get_account(self, address: str) -> Optional[Account]:
        """Get account information."""
        try:
            endpoint = f"/cosmos/auth/v1beta1/accounts/{address}"
            response = await self._get(endpoint)
            
            if "account" not in response:
                return None
                
            account_data = response["account"]
            return Account(
                address=account_data["address"],
                account_number=int(account_data["account_number"]),
                sequence=int(account_data["sequence"]),
                pub_key=account_data.get("pub_key"),
            )
        except Exception as e:
            raise NetworkError(f"Failed to get account: {e}")
    
    async def get_all_balances(self, address: str) -> List[Balance]:
        """Get all balances for an address."""
        try:
            endpoint = f"/cosmos/bank/v1beta1/balances/{address}"
            response = await self._get(endpoint)
            
            balances = []
            for balance_data in response.get("balances", []):
                balances.append(Balance(
                    denom=balance_data["denom"],
                    amount=balance_data["amount"],
                ))
            
            return balances
        except Exception as e:
            raise NetworkError(f"Failed to get balances: {e}")
    
    async def close(self):
        """Close client connections."""
        if self._client:
            await self._client.aclose()
    
    async def _get(self, endpoint: str, params: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Make async GET request."""
        if endpoint.startswith("/"):
            url = f"{self.rest_url}{endpoint}"
        else:
            url = f"{self.rpc_url}/{endpoint}"
        
        try:
            response = await self.client.get(url, params=params)
            response.raise_for_status()
            return response.json()
        except httpx.RequestError as e:
            raise NetworkError(f"Request failed: {e}", endpoint=endpoint)