"""
Cultural Heritage Client

Client for interacting with DeshChain cultural features.
"""

from typing import List, Optional, Dict, Any
from ..types import Festival, CulturalQuote
from ..exceptions import NetworkError


class CulturalClient:
    """Synchronous cultural heritage client."""
    
    def __init__(self, client):
        self.client = client
    
    def get_current_festival(self) -> Optional[Festival]:
        """Get currently active festival."""
        try:
            endpoint = "/deshchain/cultural/v1/festivals/current"
            response = self.client._get(endpoint)
            
            if "festival" not in response:
                return None
                
            return Festival(**response["festival"])
        except Exception as e:
            raise NetworkError(f"Failed to get current festival: {e}")
    
    def get_active_festivals(self) -> List[Festival]:
        """Get all active festivals."""
        try:
            endpoint = "/deshchain/cultural/v1/festivals/active"
            response = self.client._get(endpoint)
            
            festivals = []
            for festival_data in response.get("festivals", []):
                festivals.append(Festival(**festival_data))
            
            return festivals
        except Exception as e:
            raise NetworkError(f"Failed to get active festivals: {e}")
    
    def get_upcoming_festivals(self, days: int = 30) -> List[Festival]:
        """Get upcoming festivals."""
        try:
            endpoint = "/deshchain/cultural/v1/festivals/upcoming"
            params = {"days": days}
            response = self.client._get(endpoint, params=params)
            
            festivals = []
            for festival_data in response.get("festivals", []):
                festivals.append(Festival(**festival_data))
            
            return festivals
        except Exception as e:
            raise NetworkError(f"Failed to get upcoming festivals: {e}")
    
    def get_festival(self, festival_id: str) -> Optional[Festival]:
        """Get festival by ID."""
        try:
            endpoint = f"/deshchain/cultural/v1/festivals/{festival_id}"
            response = self.client._get(endpoint)
            
            if "festival" not in response:
                return None
                
            return Festival(**response["festival"])
        except Exception as e:
            raise NetworkError(f"Failed to get festival: {e}")
    
    def get_daily_quote(self) -> CulturalQuote:
        """Get daily cultural quote."""
        try:
            endpoint = "/deshchain/cultural/v1/quotes/daily"
            response = self.client._get(endpoint)
            
            return CulturalQuote(**response["quote"])
        except Exception as e:
            raise NetworkError(f"Failed to get daily quote: {e}")
    
    def get_quote_by_category(
        self, 
        category: str, 
        language: str = "en"
    ) -> CulturalQuote:
        """Get random quote by category."""
        try:
            endpoint = "/deshchain/cultural/v1/quotes/category"
            params = {"category": category, "language": language}
            response = self.client._get(endpoint, params=params)
            
            return CulturalQuote(**response["quote"])
        except Exception as e:
            raise NetworkError(f"Failed to get quote by category: {e}")
    
    def search_quotes(
        self, 
        query: str, 
        language: str = "en"
    ) -> List[CulturalQuote]:
        """Search quotes."""
        try:
            endpoint = "/deshchain/cultural/v1/quotes/search"
            params = {"query": query, "language": language}
            response = self.client._get(endpoint, params=params)
            
            quotes = []
            for quote_data in response.get("quotes", []):
                quotes.append(CulturalQuote(**quote_data))
            
            return quotes
        except Exception as e:
            raise NetworkError(f"Failed to search quotes: {e}")
    
    def get_festival_bonuses(self) -> List[Dict[str, Any]]:
        """Get current festival bonuses."""
        try:
            endpoint = "/deshchain/cultural/v1/bonuses/festival"
            response = self.client._get(endpoint)
            
            return response.get("bonuses", [])
        except Exception as e:
            raise NetworkError(f"Failed to get festival bonuses: {e}")
    
    def calculate_festival_bonus(
        self, 
        amount: float, 
        transaction_type: str
    ) -> Dict[str, Any]:
        """Calculate festival bonus for transaction."""
        try:
            endpoint = "/deshchain/cultural/v1/bonuses/calculate"
            params = {
                "amount": amount,
                "transaction_type": transaction_type
            }
            response = self.client._get(endpoint, params=params)
            
            return {
                "bonus": response["bonus"],
                "percentage": response["percentage"],
                "festival": response["festival"]
            }
        except Exception as e:
            raise NetworkError(f"Failed to calculate festival bonus: {e}")
    
    def get_regional_celebrations(self, state: str) -> List[Dict[str, Any]]:
        """Get regional celebrations."""
        try:
            endpoint = "/deshchain/cultural/v1/celebrations/regional"
            params = {"state": state}
            response = self.client._get(endpoint, params=params)
            
            return response.get("celebrations", [])
        except Exception as e:
            raise NetworkError(f"Failed to get regional celebrations: {e}")
    
    def get_cultural_preferences(self, pincode: str) -> Dict[str, Any]:
        """Get cultural preferences by pincode."""
        try:
            endpoint = "/deshchain/cultural/v1/preferences"
            params = {"pincode": pincode}
            response = self.client._get(endpoint, params=params)
            
            return response.get("preferences", {})
        except Exception as e:
            raise NetworkError(f"Failed to get cultural preferences: {e}")


class AsyncCulturalClient:
    """Asynchronous cultural heritage client."""
    
    def __init__(self, client):
        self.client = client
    
    async def get_current_festival(self) -> Optional[Festival]:
        """Get currently active festival."""
        try:
            endpoint = "/deshchain/cultural/v1/festivals/current"
            response = await self.client._get(endpoint)
            
            if "festival" not in response:
                return None
                
            return Festival(**response["festival"])
        except Exception as e:
            raise NetworkError(f"Failed to get current festival: {e}")
    
    async def get_active_festivals(self) -> List[Festival]:
        """Get all active festivals."""
        try:
            endpoint = "/deshchain/cultural/v1/festivals/active"
            response = await self.client._get(endpoint)
            
            festivals = []
            for festival_data in response.get("festivals", []):
                festivals.append(Festival(**festival_data))
            
            return festivals
        except Exception as e:
            raise NetworkError(f"Failed to get active festivals: {e}")
    
    async def get_daily_quote(self) -> CulturalQuote:
        """Get daily cultural quote."""
        try:
            endpoint = "/deshchain/cultural/v1/quotes/daily"
            response = await self.client._get(endpoint)
            
            return CulturalQuote(**response["quote"])
        except Exception as e:
            raise NetworkError(f"Failed to get daily quote: {e}")
    
    async def calculate_festival_bonus(
        self, 
        amount: float, 
        transaction_type: str
    ) -> Dict[str, Any]:
        """Calculate festival bonus for transaction."""
        try:
            endpoint = "/deshchain/cultural/v1/bonuses/calculate"
            params = {
                "amount": amount,
                "transaction_type": transaction_type
            }
            response = await self.client._get(endpoint, params=params)
            
            return {
                "bonus": response["bonus"],
                "percentage": response["percentage"],
                "festival": response["festival"]
            }
        except Exception as e:
            raise NetworkError(f"Failed to calculate festival bonus: {e}")