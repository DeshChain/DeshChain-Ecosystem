"""Sikkebaaz Launchpad Client"""

from typing import List
from ..types import LaunchpadToken
from ..exceptions import NetworkError


class SikkebaazClient:
    """Synchronous Sikkebaaz client."""
    
    def __init__(self, client):
        self.client = client
    
    def get_featured_tokens(self) -> List[LaunchpadToken]:
        """Get featured tokens."""
        try:
            endpoint = "/deshchain/sikkebaaz/v1/tokens/featured"
            response = self.client._get(endpoint)
            
            tokens = []
            for token_data in response.get("tokens", []):
                tokens.append(LaunchpadToken(**token_data))
            
            return tokens
        except Exception as e:
            raise NetworkError(f"Failed to get featured tokens: {e}")
    
    def search_tokens(self, query: str) -> List[LaunchpadToken]:
        """Search tokens."""
        try:
            endpoint = "/deshchain/sikkebaaz/v1/tokens/search"
            params = {"query": query}
            response = self.client._get(endpoint, params=params)
            
            tokens = []
            for token_data in response.get("tokens", []):
                tokens.append(LaunchpadToken(**token_data))
            
            return tokens
        except Exception as e:
            raise NetworkError(f"Failed to search tokens: {e}")


class AsyncSikkebaazClient:
    """Asynchronous Sikkebaaz client."""
    
    def __init__(self, client):
        self.client = client