"""
DeshChain SDK Exceptions

Custom exception classes for DeshChain SDK errors.
"""

from typing import Optional, Dict, Any


class DeshChainError(Exception):
    """Base exception for all DeshChain SDK errors."""
    
    def __init__(
        self, 
        message: str, 
        code: Optional[str] = None,
        details: Optional[Dict[str, Any]] = None
    ) -> None:
        super().__init__(message)
        self.message = message
        self.code = code
        self.details = details or {}
    
    def __str__(self) -> str:
        if self.code:
            return f"[{self.code}] {self.message}"
        return self.message
    
    def __repr__(self) -> str:
        return f"{self.__class__.__name__}(message='{self.message}', code='{self.code}')"


class NetworkError(DeshChainError):
    """Raised when network-related errors occur."""
    
    def __init__(
        self, 
        message: str, 
        status_code: Optional[int] = None,
        endpoint: Optional[str] = None
    ) -> None:
        super().__init__(message, "NETWORK_ERROR")
        self.status_code = status_code
        self.endpoint = endpoint


class TransactionError(DeshChainError):
    """Raised when transaction-related errors occur."""
    
    def __init__(
        self, 
        message: str, 
        tx_hash: Optional[str] = None,
        error_code: Optional[int] = None
    ) -> None:
        super().__init__(message, "TRANSACTION_ERROR")
        self.tx_hash = tx_hash
        self.error_code = error_code


class ValidationError(DeshChainError):
    """Raised when validation errors occur."""
    
    def __init__(
        self, 
        message: str, 
        field: Optional[str] = None,
        value: Optional[Any] = None
    ) -> None:
        super().__init__(message, "VALIDATION_ERROR")
        self.field = field
        self.value = value


class ConfigurationError(DeshChainError):
    """Raised when configuration errors occur."""
    
    def __init__(self, message: str, config_key: Optional[str] = None) -> None:
        super().__init__(message, "CONFIGURATION_ERROR")
        self.config_key = config_key


class AuthenticationError(DeshChainError):
    """Raised when authentication errors occur."""
    
    def __init__(self, message: str) -> None:
        super().__init__(message, "AUTHENTICATION_ERROR")


class InsufficientFundsError(DeshChainError):
    """Raised when account has insufficient funds."""
    
    def __init__(
        self, 
        message: str, 
        required: Optional[str] = None,
        available: Optional[str] = None
    ) -> None:
        super().__init__(message, "INSUFFICIENT_FUNDS")
        self.required = required
        self.available = available


class RateLimitError(DeshChainError):
    """Raised when rate limit is exceeded."""
    
    def __init__(
        self, 
        message: str, 
        retry_after: Optional[int] = None
    ) -> None:
        super().__init__(message, "RATE_LIMIT_EXCEEDED")
        self.retry_after = retry_after


class TimeoutError(DeshChainError):
    """Raised when operation times out."""
    
    def __init__(
        self, 
        message: str, 
        timeout: Optional[float] = None
    ) -> None:
        super().__init__(message, "TIMEOUT")
        self.timeout = timeout


class ContractError(DeshChainError):
    """Raised when smart contract errors occur."""
    
    def __init__(
        self, 
        message: str, 
        contract_address: Optional[str] = None,
        method: Optional[str] = None
    ) -> None:
        super().__init__(message, "CONTRACT_ERROR")
        self.contract_address = contract_address
        self.method = method


class GovernanceError(DeshChainError):
    """Raised when governance-related errors occur."""
    
    def __init__(
        self, 
        message: str, 
        proposal_id: Optional[str] = None
    ) -> None:
        super().__init__(message, "GOVERNANCE_ERROR")
        self.proposal_id = proposal_id


class CulturalError(DeshChainError):
    """Raised when cultural feature errors occur."""
    
    def __init__(
        self, 
        message: str, 
        festival_id: Optional[str] = None,
        region: Optional[str] = None
    ) -> None:
        super().__init__(message, "CULTURAL_ERROR")
        self.festival_id = festival_id
        self.region = region


class LendingError(DeshChainError):
    """Raised when lending module errors occur."""
    
    def __init__(
        self, 
        message: str, 
        loan_id: Optional[str] = None,
        module: Optional[str] = None
    ) -> None:
        super().__init__(message, "LENDING_ERROR")
        self.loan_id = loan_id
        self.module = module


class SikkebaazError(DeshChainError):
    """Raised when Sikkebaaz launchpad errors occur."""
    
    def __init__(
        self, 
        message: str, 
        token_symbol: Optional[str] = None,
        launch_id: Optional[str] = None
    ) -> None:
        super().__init__(message, "SIKKEBAAZ_ERROR")
        self.token_symbol = token_symbol
        self.launch_id = launch_id


class MoneyOrderError(DeshChainError):
    """Raised when Money Order DEX errors occur."""
    
    def __init__(
        self, 
        message: str, 
        order_id: Optional[str] = None,
        trading_pair: Optional[str] = None
    ) -> None:
        super().__init__(message, "MONEY_ORDER_ERROR")
        self.order_id = order_id
        self.trading_pair = trading_pair