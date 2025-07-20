"""
Encoding Utilities

Helper functions for encoding/decoding data.
"""

import base64
import json
import hashlib
from typing import Any, Dict


class EncodingUtils:
    """Encoding utility functions."""
    
    @staticmethod
    def to_base64(data: str) -> str:
        """Convert string to base64."""
        return base64.b64encode(data.encode()).decode()
    
    @staticmethod
    def from_base64(data: str) -> str:
        """Convert base64 to string."""
        return base64.b64decode(data).decode()
    
    @staticmethod
    def to_hex(data: bytes) -> str:
        """Convert bytes to hex string."""
        return data.hex()
    
    @staticmethod
    def from_hex(data: str) -> bytes:
        """Convert hex string to bytes."""
        return bytes.fromhex(data)
    
    @staticmethod
    def json_to_base64(data: Any) -> str:
        """Encode JSON to base64."""
        json_str = json.dumps(data)
        return EncodingUtils.to_base64(json_str)
    
    @staticmethod
    def base64_to_json(data: str) -> Any:
        """Decode JSON from base64."""
        json_str = EncodingUtils.from_base64(data)
        return json.loads(json_str)
    
    @staticmethod
    def hash_message(message: str) -> str:
        """Create message hash."""
        return hashlib.sha256(message.encode()).hexdigest()
    
    @staticmethod
    def is_valid_address(address: str, prefix: str = "deshchain") -> bool:
        """Validate address format."""
        if not address.startswith(prefix):
            return False
        if len(address) != len(prefix) + 39:
            return False
        
        address_part = address[len(prefix):]
        return address_part.isalnum() and address_part.islower()
    
    @staticmethod
    def is_valid_tx_hash(hash_str: str) -> bool:
        """Validate transaction hash."""
        return len(hash_str) == 64 and all(c in "0123456789ABCDEFabcdef" for c in hash_str)
    
    @staticmethod
    def format_amount(amount: int, decimals: int = 6) -> str:
        """Format amount with decimals."""
        return f"{amount / (10 ** decimals):.{decimals}f}"
    
    @staticmethod
    def parse_amount(amount: float, decimals: int = 6) -> int:
        """Parse amount to micro units."""
        return int(amount * (10 ** decimals))