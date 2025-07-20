/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package keeper

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// PostalRoute represents a delivery route between postal codes
type PostalRoute struct {
	FromPincode     string
	ToPincode       string
	Distance        float64
	EstimatedHours  int32
	RouteType       string // direct, hub, multi-hop
	TransportMode   string // air, rail, road, combined
	Hubs           []string
	Priority        string // express, standard, economy
	DeliveryFee     sdk.Coin
}

// RegionalHub represents a major postal hub
type RegionalHub struct {
	HubCode          string
	Name             string
	Pincode          string
	Region           string
	Capacity         int32
	ProcessingHours  int32
	ConnectedPincodes []string
}

// Major regional hubs in India
var regionalHubs = []RegionalHub{
	{
		HubCode:          "DEL-HUB",
		Name:             "Delhi NCR Hub",
		Pincode:          "110001",
		Region:           "North",
		Capacity:         100000,
		ProcessingHours:  2,
		ConnectedPincodes: []string{"110", "120", "201"},
	},
	{
		HubCode:          "MUM-HUB",
		Name:             "Mumbai Hub",
		Pincode:          "400001",
		Region:           "West",
		Capacity:         150000,
		ProcessingHours:  2,
		ConnectedPincodes: []string{"400", "410", "421"},
	},
	{
		HubCode:          "KOL-HUB",
		Name:             "Kolkata Hub",
		Pincode:          "700001",
		Region:           "East",
		Capacity:         80000,
		ProcessingHours:  3,
		ConnectedPincodes: []string{"700", "711", "712"},
	},
	{
		HubCode:          "CHN-HUB",
		Name:             "Chennai Hub",
		Pincode:          "600001",
		Region:           "South",
		Capacity:         90000,
		ProcessingHours:  2,
		ConnectedPincodes: []string{"600", "601", "603"},
	},
	{
		HubCode:          "BLR-HUB",
		Name:             "Bangalore Hub",
		Pincode:          "560001",
		Region:           "South",
		Capacity:         100000,
		ProcessingHours:  2,
		ConnectedPincodes: []string{"560", "561", "562"},
	},
}

// ValidateIndianPincode validates Indian postal code format
func (k Keeper) ValidateIndianPincode(pincode string) bool {
	// Indian postal codes are 6 digits starting with 1-9
	if len(pincode) != 6 {
		return false
	}
	
	// Check if all characters are digits
	if _, err := strconv.Atoi(pincode); err != nil {
		return false
	}
	
	// First digit should be 1-9
	firstDigit := pincode[0]
	if firstDigit < '1' || firstDigit > '9' {
		return false
	}
	
	return true
}

// GetPincodeRegion returns the region for a pincode
func (k Keeper) GetPincodeRegion(pincode string) string {
	if !k.ValidateIndianPincode(pincode) {
		return "Unknown"
	}
	
	firstDigit := pincode[0]
	switch firstDigit {
	case '1', '2':
		return "North"
	case '3', '4':
		return "West"
	case '5', '6':
		return "South"
	case '7':
		return "East"
	case '8', '9':
		return "North East"
	default:
		return "Unknown"
	}
}

// GetStateName returns the state name based on pincode prefix
func (k Keeper) GetStateName(pincode string) string {
	if !k.ValidateIndianPincode(pincode) {
		return "Unknown"
	}
	
	prefix := pincode[:2]
	stateMap := map[string]string{
		"11": "Delhi", "12": "Haryana", "13": "Punjab", "14": "Punjab",
		"20": "Uttar Pradesh", "21": "Uttar Pradesh", "22": "Uttar Pradesh",
		"30": "Rajasthan", "31": "Rajasthan", "32": "Rajasthan",
		"40": "Maharashtra", "41": "Maharashtra", "42": "Maharashtra",
		"50": "Madhya Pradesh", "51": "Chhattisgarh", "52": "Chhattisgarh",
		"56": "Karnataka", "57": "Karnataka", "58": "Karnataka",
		"60": "Tamil Nadu", "61": "Tamil Nadu", "62": "Tamil Nadu",
		"67": "Kerala", "68": "Kerala", "69": "Kerala",
		"70": "West Bengal", "71": "West Bengal", "72": "West Bengal",
		"75": "Odisha", "76": "Odisha", "77": "Odisha",
		"78": "Assam", "79": "North Eastern States",
		"80": "Bihar", "81": "Bihar", "82": "Jharkhand", "83": "Jharkhand",
	}
	
	if state, ok := stateMap[prefix]; ok {
		return state
	}
	
	return "Unknown"
}

// GetNearestHub returns the nearest hub for a pincode
func (k Keeper) GetNearestHub(pincode string) *RegionalHub {
	if !k.ValidateIndianPincode(pincode) {
		return nil
	}
	
	prefix := pincode[:3]
	
	// Check if pincode is directly connected to a hub
	for _, hub := range regionalHubs {
		for _, connectedPrefix := range hub.ConnectedPincodes {
			if strings.HasPrefix(pincode, connectedPrefix) {
				return &hub
			}
		}
	}
	
	// Default to nearest hub based on region
	region := k.GetPincodeRegion(pincode)
	for _, hub := range regionalHubs {
		if hub.Region == region {
			return &hub
		}
	}
	
	// Fallback to Delhi hub
	return &regionalHubs[0]
}

// CalculatePostalRoute calculates the optimal route between two pincodes
func (k Keeper) CalculatePostalRoute(ctx sdk.Context, fromPincode, toPincode string, priority string) (*PostalRoute, error) {
	// Validate pincodes
	if !k.ValidateIndianPincode(fromPincode) || !k.ValidateIndianPincode(toPincode) {
		return nil, types.ErrInvalidAddress
	}
	
	// Same pincode - local delivery
	if fromPincode == toPincode {
		return &PostalRoute{
			FromPincode:    fromPincode,
			ToPincode:      toPincode,
			Distance:       5,
			EstimatedHours: 2,
			RouteType:      "direct",
			TransportMode:  "road",
			Priority:       priority,
			DeliveryFee:    sdk.NewCoin(types.DefaultDenom, sdk.NewInt(10)),
		}, nil
	}
	
	// Calculate distance and route
	fromState := k.GetStateName(fromPincode)
	toState := k.GetStateName(toPincode)
	sameState := fromState == toState
	
	// Calculate approximate distance
	distance := k.calculatePincodeDistance(fromPincode, toPincode)
	
	var route PostalRoute
	route.FromPincode = fromPincode
	route.ToPincode = toPincode
	route.Distance = distance
	route.Priority = priority
	
	// Determine route type and transport mode
	if sameState && distance < 200 {
		// Intra-state, short distance
		route.RouteType = "direct"
		route.TransportMode = "road"
		route.EstimatedHours = 6
	} else if sameState && distance < 500 {
		// Intra-state, medium distance
		route.RouteType = "direct"
		route.TransportMode = "rail"
		route.EstimatedHours = 12
	} else if distance < 1000 {
		// Inter-state, medium distance
		route.RouteType = "hub"
		route.TransportMode = "combined"
		route.EstimatedHours = 24
		
		// Add hub routing
		fromHub := k.GetNearestHub(fromPincode)
		toHub := k.GetNearestHub(toPincode)
		if fromHub != nil && toHub != nil && fromHub.HubCode != toHub.HubCode {
			route.Hubs = []string{fromHub.Pincode, toHub.Pincode}
			route.EstimatedHours += fromHub.ProcessingHours + toHub.ProcessingHours
		}
	} else {
		// Long distance
		route.RouteType = "multi-hop"
		route.TransportMode = "air"
		route.EstimatedHours = 48
		
		// Multi-hop routing through hubs
		fromHub := k.GetNearestHub(fromPincode)
		toHub := k.GetNearestHub(toPincode)
		if fromHub != nil && toHub != nil {
			route.Hubs = []string{fromHub.Pincode}
			if fromHub.HubCode != toHub.HubCode {
				route.Hubs = append(route.Hubs, toHub.Pincode)
			}
			route.EstimatedHours += fromHub.ProcessingHours + toHub.ProcessingHours
		}
	}
	
	// Adjust for priority
	switch priority {
	case "express":
		route.EstimatedHours = int32(float64(route.EstimatedHours) * 0.6)
		route.TransportMode = "air" // Prefer air for express
		route.DeliveryFee = sdk.NewCoin(types.DefaultDenom, sdk.NewInt(int64(distance*2))) // 2x rate
	case "economy":
		route.EstimatedHours = int32(float64(route.EstimatedHours) * 1.5)
		route.TransportMode = "road" // Prefer road for economy
		route.DeliveryFee = sdk.NewCoin(types.DefaultDenom, sdk.NewInt(int64(distance*0.5))) // 0.5x rate
	default: // standard
		route.DeliveryFee = sdk.NewCoin(types.DefaultDenom, sdk.NewInt(int64(distance)))
	}
	
	// Minimum delivery fee
	if route.DeliveryFee.Amount.LT(sdk.NewInt(50)) {
		route.DeliveryFee = sdk.NewCoin(types.DefaultDenom, sdk.NewInt(50))
	}
	
	return &route, nil
}

// calculatePincodeDistance calculates approximate distance between pincodes
func (k Keeper) calculatePincodeDistance(from, to string) float64 {
	// Simple distance calculation based on pincode regions
	// In production, would use actual lat/lng data
	
	fromFirst := from[0]
	toFirst := to[0]
	
	// Base distance based on first digit difference
	digitDiff := math.Abs(float64(fromFirst - toFirst))
	baseDistance := digitDiff * 500 // Rough approximation
	
	// Add intra-state or inter-state distance
	fromState := k.GetStateName(from)
	toState := k.GetStateName(to)
	
	if fromState == toState {
		// Same state
		return baseDistance + 50 + (math.Mod(float64(time.Now().Unix()), 150))
	}
	
	// Different states
	return baseDistance + 200 + (math.Mod(float64(time.Now().Unix()), 300))
}

// UpdateMoneyOrderRoute updates the route information for a money order
func (k Keeper) UpdateMoneyOrderRoute(ctx sdk.Context, orderID string, route *PostalRoute) error {
	order, found := k.GetMoneyOrder(ctx, orderID)
	if !found {
		return types.ErrOrderNotFound
	}
	
	// Update delivery info
	order.DeliveryInfo = &types.DeliveryInfo{
		FromPincode:      route.FromPincode,
		ToPincode:        route.ToPincode,
		EstimatedDays:    int32(math.Ceil(float64(route.EstimatedHours) / 24)),
		TransportMode:    route.TransportMode,
		RouteType:        route.RouteType,
		DeliveryPriority: route.Priority,
		TrackingUpdates:  []types.TrackingUpdate{},
	}
	
	// Add initial tracking update
	order.DeliveryInfo.TrackingUpdates = append(order.DeliveryInfo.TrackingUpdates, types.TrackingUpdate{
		Timestamp:   ctx.BlockTime(),
		Location:    route.FromPincode,
		Status:      "Order Created",
		Description: fmt.Sprintf("Money order created at %s post office", k.GetStateName(route.FromPincode)),
	})
	
	// Save updated order
	k.SetMoneyOrder(ctx, order)
	
	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRouteCalculated,
			sdk.NewAttribute(types.AttributeKeyOrderID, orderID),
			sdk.NewAttribute("from_pincode", route.FromPincode),
			sdk.NewAttribute("to_pincode", route.ToPincode),
			sdk.NewAttribute("estimated_hours", fmt.Sprintf("%d", route.EstimatedHours)),
			sdk.NewAttribute("route_type", route.RouteType),
		),
	)
	
	return nil
}

// GetDeliveryEstimate returns delivery time estimate
func (k Keeper) GetDeliveryEstimate(ctx sdk.Context, fromPincode, toPincode string, priority string) (*types.DeliveryEstimate, error) {
	route, err := k.CalculatePostalRoute(ctx, fromPincode, toPincode, priority)
	if err != nil {
		return nil, err
	}
	
	minDays := int32(math.Ceil(float64(route.EstimatedHours) / 24))
	maxDays := minDays
	
	// Add buffer based on route type
	switch route.RouteType {
	case "direct":
		maxDays += 1
	case "hub":
		maxDays += 2
	case "multi-hop":
		maxDays += 3
	}
	
	// Calculate confidence based on route complexity
	confidence := float32(0.9)
	if route.RouteType == "multi-hop" {
		confidence = 0.75
	} else if route.RouteType == "hub" {
		confidence = 0.85
	}
	
	// Adjust confidence for priority
	if priority == "express" {
		confidence += 0.05
	} else if priority == "economy" {
		confidence -= 0.05
	}
	
	return &types.DeliveryEstimate{
		MinDays:    minDays,
		MaxDays:    maxDays,
		Confidence: confidence,
		Fee:        route.DeliveryFee,
	}, nil
}

// OptimizeBulkRoutes optimizes routes for multiple deliveries
func (k Keeper) OptimizeBulkRoutes(ctx sdk.Context, origin string, destinations []string, priority string) (map[string]*PostalRoute, error) {
	routes := make(map[string]*PostalRoute)
	
	// Group destinations by region for optimization
	regionGroups := make(map[string][]string)
	for _, dest := range destinations {
		if k.ValidateIndianPincode(dest) {
			region := k.GetPincodeRegion(dest)
			regionGroups[region] = append(regionGroups[region], dest)
		}
	}
	
	// Calculate routes for each destination
	for _, dest := range destinations {
		route, err := k.CalculatePostalRoute(ctx, origin, dest, priority)
		if err != nil {
			continue
		}
		routes[dest] = route
	}
	
	// Apply bulk discounts if applicable
	if len(routes) >= 10 {
		for _, route := range routes {
			// 10% discount for bulk orders
			discountedFee := route.DeliveryFee.Amount.MulRaw(9).QuoRaw(10)
			route.DeliveryFee = sdk.NewCoin(route.DeliveryFee.Denom, discountedFee)
		}
	}
	
	return routes, nil
}