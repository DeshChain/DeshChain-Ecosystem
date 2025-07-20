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

export interface PostalCodeInfo {
  pincode: string;
  officeName: string;
  officeType: 'HO' | 'SO' | 'BO'; // Head Office, Sub Office, Branch Office
  deliveryStatus: 'Delivery' | 'Non-Delivery';
  divisionName: string;
  regionName: string;
  circleName: string;
  taluk: string;
  districtName: string;
  stateName: string;
  country: string;
  latitude?: number;
  longitude?: number;
}

export interface PostalRoute {
  from: PostalCodeInfo;
  to: PostalCodeInfo;
  distance: number;
  estimatedTime: number; // in hours
  routeType: 'direct' | 'hub' | 'multi-hop';
  hubs?: PostalCodeInfo[];
  transportMode: 'air' | 'rail' | 'road' | 'combined';
  priority: 'express' | 'standard' | 'economy';
}

export interface RegionalHub {
  hubCode: string;
  name: string;
  pincode: string;
  region: string;
  capacity: number;
  processingTime: number; // in hours
  connectedPincodes: string[];
  transportModes: string[];
}

// Major postal circles in India
export const POSTAL_CIRCLES = {
  'AP': 'Andhra Pradesh',
  'AS': 'Assam',
  'BR': 'Bihar',
  'CG': 'Chhattisgarh',
  'DL': 'Delhi',
  'GA': 'Goa',
  'GJ': 'Gujarat',
  'HP': 'Himachal Pradesh',
  'HR': 'Haryana',
  'JH': 'Jharkhand',
  'JK': 'Jammu & Kashmir',
  'KA': 'Karnataka',
  'KL': 'Kerala',
  'MH': 'Maharashtra',
  'MP': 'Madhya Pradesh',
  'NE': 'North Eastern',
  'OD': 'Odisha',
  'PB': 'Punjab',
  'RJ': 'Rajasthan',
  'TN': 'Tamil Nadu',
  'TG': 'Telangana',
  'UK': 'Uttarakhand',
  'UP': 'Uttar Pradesh',
  'WB': 'West Bengal'
};

// Regional hubs for efficient routing
export const REGIONAL_HUBS: RegionalHub[] = [
  {
    hubCode: 'DEL-HUB',
    name: 'Delhi NCR Hub',
    pincode: '110001',
    region: 'North',
    capacity: 100000,
    processingTime: 2,
    connectedPincodes: ['110xxx', '120xxx', '201xxx'],
    transportModes: ['air', 'rail', 'road']
  },
  {
    hubCode: 'MUM-HUB',
    name: 'Mumbai Hub',
    pincode: '400001',
    region: 'West',
    capacity: 150000,
    processingTime: 2,
    connectedPincodes: ['400xxx', '410xxx', '421xxx'],
    transportModes: ['air', 'rail', 'road', 'sea']
  },
  {
    hubCode: 'KOL-HUB',
    name: 'Kolkata Hub',
    pincode: '700001',
    region: 'East',
    capacity: 80000,
    processingTime: 3,
    connectedPincodes: ['700xxx', '711xxx', '712xxx'],
    transportModes: ['air', 'rail', 'road']
  },
  {
    hubCode: 'CHN-HUB',
    name: 'Chennai Hub',
    pincode: '600001',
    region: 'South',
    capacity: 90000,
    processingTime: 2,
    connectedPincodes: ['600xxx', '601xxx', '603xxx'],
    transportModes: ['air', 'rail', 'road', 'sea']
  },
  {
    hubCode: 'BLR-HUB',
    name: 'Bangalore Hub',
    pincode: '560001',
    region: 'South',
    capacity: 100000,
    processingTime: 2,
    connectedPincodes: ['560xxx', '561xxx', '562xxx'],
    transportModes: ['air', 'rail', 'road']
  }
];

export class PostalCodeService {
  private static pincodeCache = new Map<string, PostalCodeInfo>();
  private static routeCache = new Map<string, PostalRoute>();

  // Validate Indian postal code
  static isValidPincode(pincode: string): boolean {
    // Indian postal codes are 6 digits
    const pincodeRegex = /^[1-9][0-9]{5}$/;
    return pincodeRegex.test(pincode);
  }

  // Get postal code info (mock implementation - would connect to real API)
  static async getPostalCodeInfo(pincode: string): Promise<PostalCodeInfo | null> {
    if (!this.isValidPincode(pincode)) {
      return null;
    }

    // Check cache
    if (this.pincodeCache.has(pincode)) {
      return this.pincodeCache.get(pincode)!;
    }

    // Mock data based on pincode patterns
    const info = this.generateMockPostalInfo(pincode);
    
    // Cache the result
    this.pincodeCache.set(pincode, info);
    
    return info;
  }

  // Calculate optimal route between two postal codes
  static async calculateRoute(
    fromPincode: string,
    toPincode: string,
    priority: 'express' | 'standard' | 'economy' = 'standard'
  ): Promise<PostalRoute | null> {
    const cacheKey = `${fromPincode}-${toPincode}-${priority}`;
    
    // Check cache
    if (this.routeCache.has(cacheKey)) {
      return this.routeCache.get(cacheKey)!;
    }

    const fromInfo = await this.getPostalCodeInfo(fromPincode);
    const toInfo = await this.getPostalCodeInfo(toPincode);

    if (!fromInfo || !toInfo) {
      return null;
    }

    // Calculate route
    const route = this.calculateOptimalRoute(fromInfo, toInfo, priority);
    
    // Cache the result
    this.routeCache.set(cacheKey, route);
    
    return route;
  }

  // Get nearest hub for a pincode
  static getNearestHub(pincode: string): RegionalHub | null {
    const firstThree = pincode.substring(0, 3);
    
    // Find hub that serves this pincode
    for (const hub of REGIONAL_HUBS) {
      const matches = hub.connectedPincodes.some(pattern => {
        const prefix = pattern.replace(/x/g, '');
        return pincode.startsWith(prefix);
      });
      
      if (matches) {
        return hub;
      }
    }
    
    // Default to nearest major hub based on first digit
    const firstDigit = parseInt(pincode[0]);
    if (firstDigit <= 2) return REGIONAL_HUBS[0]; // Delhi
    if (firstDigit <= 4) return REGIONAL_HUBS[1]; // Mumbai
    if (firstDigit <= 5) return REGIONAL_HUBS[4]; // Bangalore
    if (firstDigit <= 6) return REGIONAL_HUBS[3]; // Chennai
    return REGIONAL_HUBS[2]; // Kolkata
  }

  // Calculate distance between two pincodes (simplified)
  private static calculateDistance(from: PostalCodeInfo, to: PostalCodeInfo): number {
    // In real implementation, would use actual coordinates
    const fromFirst = parseInt(from.pincode[0]);
    const toFirst = parseInt(to.pincode[0]);
    
    // Rough distance calculation based on pincode regions
    const regionDistance = Math.abs(fromFirst - toFirst) * 500; // km
    
    // Add intra-state distance
    if (from.stateName === to.stateName) {
      return regionDistance + 50 + Math.random() * 150;
    }
    
    return regionDistance + 200 + Math.random() * 300;
  }

  // Calculate optimal route
  private static calculateOptimalRoute(
    from: PostalCodeInfo,
    to: PostalCodeInfo,
    priority: 'express' | 'standard' | 'economy'
  ): PostalRoute {
    const distance = this.calculateDistance(from, to);
    const sameState = from.stateName === to.stateName;
    const sameCity = from.districtName === to.districtName;
    
    let routeType: 'direct' | 'hub' | 'multi-hop' = 'direct';
    let transportMode: 'air' | 'rail' | 'road' | 'combined' = 'road';
    let estimatedTime = 24; // hours
    let hubs: PostalCodeInfo[] = [];

    // Determine route type and transport mode
    if (sameCity) {
      routeType = 'direct';
      transportMode = 'road';
      estimatedTime = 4 + Math.random() * 4;
    } else if (sameState) {
      routeType = 'direct';
      transportMode = distance > 300 ? 'rail' : 'road';
      estimatedTime = 8 + Math.random() * 8;
    } else if (distance > 1000) {
      routeType = 'hub';
      transportMode = priority === 'express' ? 'air' : 'combined';
      
      // Add hub routing
      const fromHub = this.getNearestHub(from.pincode);
      const toHub = this.getNearestHub(to.pincode);
      
      if (fromHub && toHub && fromHub.hubCode !== toHub.hubCode) {
        hubs = [
          this.generateMockPostalInfo(fromHub.pincode),
          this.generateMockPostalInfo(toHub.pincode)
        ];
        routeType = 'multi-hop';
      }
      
      estimatedTime = priority === 'express' ? 24 : 48;
    } else {
      routeType = 'hub';
      transportMode = 'combined';
      estimatedTime = 24 + Math.random() * 24;
    }

    // Adjust time based on priority
    if (priority === 'express') {
      estimatedTime *= 0.6;
    } else if (priority === 'economy') {
      estimatedTime *= 1.5;
    }

    return {
      from,
      to,
      distance,
      estimatedTime: Math.round(estimatedTime),
      routeType,
      hubs,
      transportMode,
      priority
    };
  }

  // Generate mock postal info
  private static generateMockPostalInfo(pincode: string): PostalCodeInfo {
    const stateCode = this.getStateFromPincode(pincode);
    const stateName = this.getStateName(stateCode);
    
    return {
      pincode,
      officeName: `Post Office ${pincode}`,
      officeType: pincode.endsWith('001') ? 'HO' : pincode.endsWith('00') ? 'SO' : 'BO',
      deliveryStatus: 'Delivery',
      divisionName: `Division ${pincode.substring(0, 3)}`,
      regionName: this.getRegionFromPincode(pincode),
      circleName: stateName,
      taluk: `Taluk ${pincode.substring(0, 4)}`,
      districtName: `District ${pincode.substring(0, 3)}`,
      stateName,
      country: 'India',
      latitude: 20 + Math.random() * 15,
      longitude: 70 + Math.random() * 20
    };
  }

  // Get state from pincode (simplified mapping)
  private static getStateFromPincode(pincode: string): string {
    const firstTwo = pincode.substring(0, 2);
    const stateMap: Record<string, string> = {
      '11': 'DL', '12': 'HR', '13': 'PB', '14': 'PB',
      '20': 'UP', '21': 'UP', '22': 'UP', '23': 'UP',
      '30': 'RJ', '31': 'RJ', '32': 'RJ', '33': 'RJ',
      '40': 'MH', '41': 'MH', '42': 'MH', '43': 'MH',
      '50': 'MP', '51': 'CG', '52': 'CG',
      '56': 'KA', '57': 'KA', '58': 'KA', '59': 'KA',
      '60': 'TN', '61': 'TN', '62': 'TN', '63': 'TN',
      '67': 'KL', '68': 'KL', '69': 'KL',
      '70': 'WB', '71': 'WB', '72': 'WB', '73': 'WB',
      '75': 'OD', '76': 'OD', '77': 'OD',
      '78': 'AS', '79': 'NE',
      '80': 'BR', '81': 'BR', '82': 'JH', '83': 'JH'
    };
    
    return stateMap[firstTwo] || 'DL';
  }

  // Get state name from code
  private static getStateName(stateCode: string): string {
    return POSTAL_CIRCLES[stateCode] || 'Unknown';
  }

  // Get region from pincode
  private static getRegionFromPincode(pincode: string): string {
    const firstDigit = parseInt(pincode[0]);
    
    if (firstDigit <= 2) return 'North';
    if (firstDigit <= 4) return 'West';
    if (firstDigit <= 6) return 'South';
    if (firstDigit <= 7) return 'East';
    return 'North East';
  }

  // Bulk route optimization for multiple deliveries
  static async optimizeBulkRoutes(
    origin: string,
    destinations: string[],
    priority: 'express' | 'standard' | 'economy' = 'standard'
  ): Promise<Map<string, PostalRoute>> {
    const routes = new Map<string, PostalRoute>();
    
    // Group destinations by region for efficiency
    const groupedDestinations = new Map<string, string[]>();
    
    for (const dest of destinations) {
      const region = this.getRegionFromPincode(dest);
      if (!groupedDestinations.has(region)) {
        groupedDestinations.set(region, []);
      }
      groupedDestinations.get(region)!.push(dest);
    }
    
    // Calculate routes for each destination
    for (const dest of destinations) {
      const route = await this.calculateRoute(origin, dest, priority);
      if (route) {
        routes.set(dest, route);
      }
    }
    
    return routes;
  }

  // Get delivery estimate
  static getDeliveryEstimate(route: PostalRoute): {
    minDays: number;
    maxDays: number;
    confidence: number;
  } {
    const baseDays = Math.ceil(route.estimatedTime / 24);
    
    // Add buffer based on route type
    let buffer = 0;
    let confidence = 0.9;
    
    switch (route.routeType) {
      case 'direct':
        buffer = 1;
        confidence = 0.95;
        break;
      case 'hub':
        buffer = 2;
        confidence = 0.85;
        break;
      case 'multi-hop':
        buffer = 3;
        confidence = 0.75;
        break;
    }
    
    // Adjust for priority
    if (route.priority === 'express') {
      confidence += 0.05;
    } else if (route.priority === 'economy') {
      buffer += 1;
      confidence -= 0.05;
    }
    
    return {
      minDays: baseDays,
      maxDays: baseDays + buffer,
      confidence: Math.min(confidence, 0.99)
    };
  }
}