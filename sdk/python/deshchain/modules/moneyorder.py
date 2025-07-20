"""Money Order DEX Client"""

from typing import List
from ..types import MoneyOrder
from ..exceptions import NetworkError


class MoneyOrderClient:
    """Synchronous Money Order client."""
    
    def __init__(self, client):
        self.client = client
    
    def get_money_order(self, order_id: str) -> MoneyOrder:
        """Get money order by ID."""
        try:
            endpoint = f"/deshchain/moneyorder/v1/orders/{order_id}"
            response = self.client._get(endpoint)
            return MoneyOrder(**response["order"])
        except Exception as e:
            raise NetworkError(f"Failed to get money order: {e}")


class AsyncMoneyOrderClient:
    """Asynchronous Money Order client."""
    
    def __init__(self, client):
        self.client = client