"""Lending Module Client - Krishi/Vyavasaya/Shiksha Mitra"""

from typing import List, Optional
from ..types import Loan, LoanStats
from ..exceptions import NetworkError


class LendingClient:
    """Synchronous lending client."""
    
    def __init__(self, client):
        self.client = client
    
    def get_krishi_mitra_stats(self) -> LoanStats:
        """Get Krishi Mitra statistics."""
        try:
            endpoint = "/deshchain/lending/v1/krishi/stats"
            response = self.client._get(endpoint)
            return LoanStats(**response["stats"])
        except Exception as e:
            raise NetworkError(f"Failed to get Krishi Mitra stats: {e}")
    
    def get_vyavasaya_mitra_stats(self) -> LoanStats:
        """Get Vyavasaya Mitra statistics."""
        try:
            endpoint = "/deshchain/lending/v1/vyavasaya/stats"
            response = self.client._get(endpoint)
            return LoanStats(**response["stats"])
        except Exception as e:
            raise NetworkError(f"Failed to get Vyavasaya Mitra stats: {e}")
    
    def get_shiksha_mitra_stats(self) -> LoanStats:
        """Get Shiksha Mitra statistics."""
        try:
            endpoint = "/deshchain/lending/v1/shiksha/stats"
            response = self.client._get(endpoint)
            return LoanStats(**response["stats"])
        except Exception as e:
            raise NetworkError(f"Failed to get Shiksha Mitra stats: {e}")
    
    def search_loans(self, query: str) -> List[Loan]:
        """Search loans across all modules."""
        try:
            endpoint = "/deshchain/lending/v1/loans/search"
            params = {"query": query}
            response = self.client._get(endpoint, params=params)
            
            loans = []
            for loan_data in response.get("loans", []):
                loans.append(Loan(**loan_data))
            
            return loans
        except Exception as e:
            raise NetworkError(f"Failed to search loans: {e}")


class AsyncLendingClient:
    """Asynchronous lending client."""
    
    def __init__(self, client):
        self.client = client
    
    async def get_krishi_mitra_stats(self) -> LoanStats:
        """Get Krishi Mitra statistics."""
        try:
            endpoint = "/deshchain/lending/v1/krishi/stats"
            response = await self.client._get(endpoint)
            return LoanStats(**response["stats"])
        except Exception as e:
            raise NetworkError(f"Failed to get Krishi Mitra stats: {e}")