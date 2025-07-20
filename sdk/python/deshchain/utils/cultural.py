"""
Cultural Utilities

Helper functions for cultural features.
"""

from typing import Dict, List


class CulturalUtils:
    """Cultural utility functions."""
    
    @staticmethod
    def get_state_from_pincode(pincode: str) -> str:
        """Get state from pincode."""
        first_digit = pincode[0] if pincode else ""
        
        state_map = {
            "1": "Delhi",
            "2": "Haryana/Himachal Pradesh", 
            "3": "Punjab/Jammu & Kashmir",
            "4": "Rajasthan",
            "5": "Uttar Pradesh/Uttarakhand",
            "6": "Bihar/Jharkhand",
            "7": "West Bengal/Sikkim",
            "8": "Odisha",
            "9": "Assam/Northeast States",
        }
        
        return state_map.get(first_digit, "Unknown")
    
    @staticmethod
    def get_regional_language(state: str) -> str:
        """Get regional language from state."""
        language_map = {
            "West Bengal": "Bengali",
            "Tamil Nadu": "Tamil",
            "Maharashtra": "Marathi",
            "Gujarat": "Gujarati",
            "Punjab": "Punjabi",
            "Karnataka": "Kannada",
            "Kerala": "Malayalam",
            "Andhra Pradesh": "Telugu",
            "Odisha": "Odia",
            "Assam": "Assamese",
        }
        
        return language_map.get(state, "Hindi")
    
    @staticmethod
    def get_major_festivals(state: str) -> List[str]:
        """Get major festivals by state."""
        festival_map = {
            "West Bengal": ["Durga Puja", "Kali Puja", "Poila Boishakh"],
            "Tamil Nadu": ["Pongal", "Diwali", "Navaratri"],
            "Kerala": ["Onam", "Vishu", "Thrissur Pooram"],
            "Maharashtra": ["Ganesh Chaturthi", "Gudi Padwa", "Navratri"],
            "Gujarat": ["Navratri", "Kite Festival", "Diwali"],
            "Punjab": ["Baisakhi", "Karva Chauth", "Lohri"],
        }
        
        return festival_map.get(state, ["Diwali", "Holi", "Dussehra"])
    
    @staticmethod
    def format_indian_currency(amount: float) -> str:
        """Format currency in Indian style."""
        if amount >= 10000000:  # 1 Crore
            return f"₹{amount/10000000:.2f} Cr"
        elif amount >= 100000:  # 1 Lakh
            return f"₹{amount/100000:.2f} L"
        elif amount >= 1000:  # 1 Thousand
            return f"₹{amount/1000:.2f} K"
        else:
            return f"₹{amount:.2f}"