"""
Validation Utilities

Helper functions for data validation.
"""

import re
from typing import Dict, Any, Optional


class ValidationUtils:
    """Validation utility functions."""
    
    @staticmethod
    def validate_address(address: str) -> Dict[str, Any]:
        """Validate DeshChain address."""
        if not address:
            return {"valid": False, "error": "Address is required"}
        
        if not address.startswith("deshchain"):
            return {"valid": False, "error": "Address must start with 'deshchain'"}
        
        if len(address) != 45:
            return {"valid": False, "error": "Address must be 45 characters long"}
        
        address_part = address[9:]  # Remove 'deshchain' prefix
        if not re.match("^[a-z0-9]+$", address_part):
            return {"valid": False, "error": "Address contains invalid characters"}
        
        return {"valid": True}
    
    @staticmethod
    def validate_amount(amount: float) -> Dict[str, Any]:
        """Validate amount."""
        if not isinstance(amount, (int, float)):
            return {"valid": False, "error": "Amount must be a number"}
        
        if amount <= 0:
            return {"valid": False, "error": "Amount must be greater than 0"}
        
        if amount > 1e15:  # Very large number check
            return {"valid": False, "error": "Amount is too large"}
        
        return {"valid": True}
    
    @staticmethod
    def validate_pincode(pincode: str) -> Dict[str, Any]:
        """Validate Indian pincode."""
        if not re.match("^[1-9][0-9]{5}$", pincode):
            return {"valid": False, "error": "Invalid pincode format"}
        
        return {"valid": True}
    
    @staticmethod
    def validate_phone_number(phone: str) -> Dict[str, Any]:
        """Validate Indian phone number."""
        if not re.match("^[6-9][0-9]{9}$", phone):
            return {"valid": False, "error": "Invalid phone number format"}
        
        return {"valid": True}
    
    @staticmethod
    def validate_aadhar(aadhar: str) -> Dict[str, Any]:
        """Validate Aadhar number."""
        # Remove spaces and hyphens
        clean_aadhar = re.sub(r"[\s-]", "", aadhar)
        
        if not re.match("^[0-9]{12}$", clean_aadhar):
            return {"valid": False, "error": "Aadhar must be 12 digits"}
        
        # Simple validation - in real implementation would use proper algorithm
        return {"valid": True}
    
    @staticmethod
    def validate_pan(pan: str) -> Dict[str, Any]:
        """Validate PAN number."""
        if not re.match("^[A-Z]{5}[0-9]{4}[A-Z]{1}$", pan.upper()):
            return {"valid": False, "error": "Invalid PAN format"}
        
        return {"valid": True}
    
    @staticmethod
    def validate_email(email: str) -> Dict[str, Any]:
        """Validate email address."""
        pattern = r"^[^\s@]+@[^\s@]+\.[^\s@]+$"
        if not re.match(pattern, email):
            return {"valid": False, "error": "Invalid email format"}
        
        return {"valid": True}
    
    @staticmethod
    def validate_token_symbol(symbol: str) -> Dict[str, Any]:
        """Validate token symbol."""
        if not symbol:
            return {"valid": False, "error": "Token symbol is required"}
        
        if len(symbol) < 2 or len(symbol) > 8:
            return {"valid": False, "error": "Token symbol must be 2-8 characters"}
        
        if not re.match("^[A-Z][A-Z0-9]*$", symbol):
            return {"valid": False, "error": "Token symbol must start with letter and contain only uppercase letters and numbers"}
        
        return {"valid": True}