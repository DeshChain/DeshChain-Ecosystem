"""Governance Client"""

from typing import List
from ..types import Proposal, Vote
from ..exceptions import NetworkError


class GovernanceClient:
    """Synchronous Governance client."""
    
    def __init__(self, client):
        self.client = client
    
    def get_proposals(self) -> List[Proposal]:
        """Get all proposals."""
        try:
            endpoint = "/deshchain/governance/v1/proposals"
            response = self.client._get(endpoint)
            
            proposals = []
            for proposal_data in response.get("proposals", []):
                proposals.append(Proposal(**proposal_data))
            
            return proposals
        except Exception as e:
            raise NetworkError(f"Failed to get proposals: {e}")


class AsyncGovernanceClient:
    """Asynchronous Governance client."""
    
    def __init__(self, client):
        self.client = client