"""
DeshChain Python SDK

Official Python SDK for interacting with DeshChain blockchain.
Features cultural heritage integration, lending modules, festival celebrations.
"""

__version__ = "1.0.0"
__author__ = "DeshChain Development Team"
__email__ = "dev@deshchain.network"

# Core client
from .client import DeshChainClient, AsyncDeshChainClient

# Module clients
from .modules.cultural import CulturalClient, AsyncCulturalClient
from .modules.lending import LendingClient, AsyncLendingClient
from .modules.sikkebaaz import SikkebaazClient, AsyncSikkebaazClient
from .modules.moneyorder import MoneyOrderClient, AsyncMoneyOrderClient
from .modules.governance import GovernanceClient, AsyncGovernanceClient

# Types
from .types import *
from .types.cultural import *
from .types.lending import *
from .types.sikkebaaz import *
from .types.moneyorder import *
from .types.governance import *

# Utilities
from .utils import *
from .utils.cultural import CulturalUtils
from .utils.festival import FestivalUtils
from .utils.validation import ValidationUtils
from .utils.encoding import EncodingUtils

# Constants
from .constants import *

# Exceptions
from .exceptions import (
    DeshChainError,
    NetworkError,
    TransactionError,
    ValidationError,
    ConfigurationError,
)

__all__ = [
    # Version info
    "__version__",
    "__author__",
    "__email__",
    
    # Core clients
    "DeshChainClient",
    "AsyncDeshChainClient",
    
    # Module clients
    "CulturalClient",
    "AsyncCulturalClient",
    "LendingClient", 
    "AsyncLendingClient",
    "SikkebaazClient",
    "AsyncSikkebaazClient",
    "MoneyOrderClient",
    "AsyncMoneyOrderClient",
    "GovernanceClient",
    "AsyncGovernanceClient",
    
    # Utilities
    "CulturalUtils",
    "FestivalUtils",
    "ValidationUtils",
    "EncodingUtils",
    
    # Exceptions
    "DeshChainError",
    "NetworkError", 
    "TransactionError",
    "ValidationError",
    "ConfigurationError",
]