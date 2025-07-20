"""
Festival Utilities

Helper functions for festival features.
"""

from typing import List, Dict, Any, Optional
from datetime import datetime, date


class FestivalUtils:
    """Festival utility functions."""
    
    @staticmethod
    def get_current_festivals(current_date: Optional[date] = None) -> List[Dict[str, Any]]:
        """Get currently active festivals."""
        if current_date is None:
            current_date = date.today()
        
        # This would normally fetch from API, but here's a simple implementation
        festivals = FestivalUtils._get_all_festivals()
        
        active = []
        for festival in festivals:
            festival_date = datetime.strptime(festival["date"], "%Y-%m-%d").date()
            end_date = festival_date  # Simplified - festivals are 1 day
            
            if festival_date <= current_date <= end_date:
                active.append(festival)
        
        return active
    
    @staticmethod
    def is_festival_day(check_date: Optional[date] = None) -> Dict[str, Any]:
        """Check if date is during a festival."""
        if check_date is None:
            check_date = date.today()
            
        active_festivals = FestivalUtils.get_current_festivals(check_date)
        
        return {
            "is_festival": len(active_festivals) > 0,
            "festivals": active_festivals
        }
    
    @staticmethod
    def get_festival_bonus(check_date: Optional[date] = None) -> float:
        """Get festival bonus multiplier for date."""
        if check_date is None:
            check_date = date.today()
            
        active_festivals = FestivalUtils.get_current_festivals(check_date)
        
        if not active_festivals:
            return 1.0
        
        # Return the highest bonus multiplier
        max_bonus = max(f.get("bonus_multiplier", 1.0) for f in active_festivals)
        return max_bonus
    
    @staticmethod
    def get_festival_greeting(festival_name: str) -> Dict[str, str]:
        """Generate festival greeting."""
        greetings = {
            "Diwali": {
                "english": "Happy Diwali! May the festival of lights bring joy and prosperity.",
                "hindi": "दीपावली की हार्दिक शुभकामनाएं!"
            },
            "Holi": {
                "english": "Happy Holi! May your life be filled with colors of joy.",
                "hindi": "होली की शुभकामनाएं!"
            },
            "Independence Day": {
                "english": "Happy Independence Day! Jai Hind!",
                "hindi": "स्वतंत्रता दिवस की शुभकामनाएं! जय हिन्द!"
            }
        }
        
        return greetings.get(festival_name, {
            "english": f"Happy {festival_name}!",
            "hindi": f"{festival_name} की शुभकामनाएं!"
        })
    
    @staticmethod
    def _get_all_festivals() -> List[Dict[str, Any]]:
        """Get all festivals (static data for demo)."""
        current_year = datetime.now().year
        
        return [
            {
                "name": "Republic Day",
                "name_hindi": "गणतंत्र दिवस",
                "date": f"{current_year}-01-26",
                "category": "National",
                "bonus_multiplier": 1.5,
                "regions": ["All India"]
            },
            {
                "name": "Independence Day", 
                "name_hindi": "स्वतंत्रता दिवस",
                "date": f"{current_year}-08-15",
                "category": "National",
                "bonus_multiplier": 1.5,
                "regions": ["All India"]
            },
            {
                "name": "Diwali",
                "name_hindi": "दीपावली", 
                "date": f"{current_year}-11-12",  # Approximate
                "category": "Religious",
                "bonus_multiplier": 1.5,
                "regions": ["All India"]
            },
            {
                "name": "Holi",
                "name_hindi": "होली",
                "date": f"{current_year}-03-13",  # Approximate
                "category": "Religious", 
                "bonus_multiplier": 1.3,
                "regions": ["All India"]
            }
        ]