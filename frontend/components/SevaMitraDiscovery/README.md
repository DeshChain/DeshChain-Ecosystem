# Seva Mitra Discovery UI - Component Documentation

## Overview

The Seva Mitra Discovery UI is a comprehensive React component system for finding and interacting with nearby Seva Mitras (community service friends) who provide cash-in, cash-out, remittance, and bill payment services. The system features an interactive map view, detailed list view, and advanced filtering capabilities.

## Features

### 🗺️ **Interactive Map View**
- **Google Maps Integration**: Real-time map with custom markers for each Seva Mitra
- **Trust Score Visualization**: Color-coded markers based on trust scores (Diamond/Platinum/Gold/Silver/Bronze)
- **Location Circles**: Visual radius showing maximum distance preferences
- **Rich Info Windows**: Detailed popup cards with comprehensive Seva Mitra information
- **Real-time Status**: Live open/closed status indicators
- **Direct Actions**: One-click call and directions functionality

### 📋 **Comprehensive List View**
- **Card-based Layout**: Beautiful cards with expandable details
- **Trust Badge System**: Visual trust score indicators with color coding
- **Real-time Availability**: Live status showing if currently open
- **Service Icons**: Visual representation of available services
- **Rating System**: Star ratings with review counts
- **Contact Integration**: Direct call and navigation buttons
- **Expandable Details**: Collapsible sections for additional information

### 🔍 **Advanced Filtering System**
- **Service Types**: Filter by Cash In, Cash Out, Remittance, Bill Payment
- **Geographic Filters**: Postal code, distance radius (1-50km)
- **Trust & Ratings**: Minimum trust score requirements with badge visualization
- **Availability**: Currently open only, KYC verified only
- **Language Support**: 22 Indian languages with chip selection
- **Payment Methods**: UPI, IMPS, NEFT, RTGS, Cash compatibility
- **Sorting Options**: Distance, trust score, rating, response time

### 📱 **Mobile-First Design**
- **Responsive Layout**: Optimized for mobile, tablet, and desktop
- **Touch-friendly**: Large tap targets and gesture support
- **Floating Action Button**: Quick filter access on mobile
- **Swipe Navigation**: Smooth tab switching
- **Progressive Enhancement**: Works across all device sizes

## Component Architecture

```
SevaMitraDiscovery/
├── index.tsx              # Main container component
├── MapView.tsx            # Google Maps integration
├── ListView.tsx           # List view with cards
├── FilterPanel.tsx        # Advanced filtering
└── README.md             # This documentation
```

## Data Structure

### SevaMitra Interface
```typescript
interface SevaMitra {
  mitraId: string;           // Unique identifier
  businessName: string;      // Display name
  address: string;           // Full address
  postalCode: string;        // 6-digit PIN code
  district: string;          // District name
  state: string;             // State name
  latitude: number;          // GPS coordinates
  longitude: number;         // GPS coordinates
  phone: string;             // Contact number
  email: string;             // Email address
  languages: string[];       // Supported languages
  services: string[];        // Available services
  trustScore: number;        // 0-100 trust rating
  averageRating: number;     // Star rating (1-5)
  totalRatings: number;      // Number of reviews
  isActive: boolean;         // Currently active
  isKYCVerified: boolean;    // KYC verification status
  operatingHours: OperatingHours[];
  distance?: number;         // Distance from user (km)
  responseTime?: string;     // Average response time
  commissionRate?: string;   // Commission percentage
}
```

### Operating Hours Structure
```typescript
interface OperatingHours {
  day: string;               // Day of week
  openTime: string;          // Opening time (24hr format)
  closeTime: string;         // Closing time (24hr format)
  isClosed: boolean;         // Closed on this day
}
```

## Key Features Explained

### 1. Trust Score System
- **Diamond (90+)**: Purple markers/badges - Highest trust level
- **Platinum (80-89)**: Blue-grey markers/badges - Very high trust
- **Gold (70-79)**: Orange markers/badges - High trust
- **Silver (60-69)**: Grey markers/badges - Good trust
- **Bronze (50-59)**: Brown markers/badges - Basic trust
- **New (<50)**: Red markers/badges - New users

### 2. Real-time Status Tracking
- **Live Open/Closed Status**: Based on current time vs operating hours
- **Status Indicators**: Green badges for open, red for closed
- **Response Time Tracking**: Average response time display
- **Auto-refresh Option**: 30-second intervals for live updates

### 3. Geographic Intelligence
- **Postal Code Integration**: 6-digit Indian PIN code support
- **Distance Calculation**: Accurate distance from user location
- **Radius Visualization**: Circle overlay showing search radius
- **Location Detection**: GPS-based current location with reverse geocoding

### 4. Cultural Localization
- **22 Indian Languages**: Complete regional language support
- **Cultural Design**: Indian color schemes and iconography
- **Local Payment Methods**: UPI, IMPS, NEFT preferred ordering
- **Regional Services**: India-specific financial services

## Usage Examples

### Basic Implementation
```tsx
import SevaMitraDiscovery from './components/SevaMitraDiscovery';

function App() {
  return (
    <div className="App">
      <SevaMitraDiscovery />
    </div>
  );
}
```

### With Custom Filters
```tsx
const customFilters = {
  services: ['CASH_IN', 'REMITTANCE'],
  maxDistance: 10,
  minTrustScore: 70,
  isKYCRequired: true,
  // ... other filter options
};

<SevaMitraDiscovery defaultFilters={customFilters} />
```

## Environment Setup

### Required Environment Variables
```bash
NEXT_PUBLIC_GOOGLE_MAPS_API_KEY=your_google_maps_api_key
```

### Dependencies
```json
{
  "@mui/material": "^5.x",
  "@mui/icons-material": "^5.x",
  "@emotion/react": "^11.x",
  "@emotion/styled": "^11.x",
  "@react-google-maps/api": "^2.x",
  "react": "^18.x",
  "typescript": "^5.x"
}
```

## API Integration

### Expected API Endpoints
```typescript
// Fetch Seva Mitras with filters
GET /api/seva-mitras?filters={filterObject}

// Get Seva Mitra details
GET /api/seva-mitras/{mitraId}

// Reverse geocoding for postal code detection
GET /api/geocoding/reverse?lat={lat}&lng={lng}
```

### Mock Data Structure
The component includes comprehensive mock data for development and testing, featuring realistic Indian business names, addresses, and operating patterns.

## Performance Optimizations

### 1. Map Performance
- **Marker Clustering**: Groups nearby markers at high zoom levels
- **Lazy Loading**: Loads map only when tab is active
- **Debounced Filtering**: Reduces API calls during filter changes
- **Viewport Optimization**: Only renders visible markers

### 2. List Performance
- **Virtual Scrolling**: Handles large lists efficiently
- **Image Lazy Loading**: Loads avatars on demand
- **Memoized Components**: Prevents unnecessary re-renders
- **Intersection Observer**: Tracks visible items

### 3. Filter Performance
- **Client-side Filtering**: Fast filtering for small datasets
- **Debounced Input**: Reduces search API calls
- **Cached Results**: Stores filter results temporarily
- **Progressive Enhancement**: Core functionality works without JavaScript

## Accessibility Features

### 1. Screen Reader Support
- **ARIA Labels**: Comprehensive screen reader descriptions
- **Semantic HTML**: Proper heading structure and landmarks
- **Focus Management**: Keyboard navigation support
- **Alternative Text**: Descriptive text for visual elements

### 2. Keyboard Navigation
- **Tab Order**: Logical keyboard navigation flow
- **Escape Key**: Closes dialogs and dropdowns
- **Enter/Space**: Activates buttons and links
- **Arrow Keys**: Navigate through lists and map

### 3. Visual Accessibility
- **High Contrast**: Meets WCAG AA standards
- **Color Independence**: Information not solely color-dependent
- **Scalable Text**: Responsive to browser zoom
- **Focus Indicators**: Clear visual focus states

## Browser Compatibility

### Supported Browsers
- **Chrome**: 90+ (Full support)
- **Firefox**: 88+ (Full support)
- **Safari**: 14+ (Full support)
- **Edge**: 90+ (Full support)
- **Mobile Browsers**: iOS Safari 14+, Chrome Mobile 90+

### Progressive Enhancement
- **Core Functionality**: Works without JavaScript
- **Map Fallback**: Static image map when Google Maps unavailable
- **Offline Support**: Cached data for offline viewing
- **Reduced Data Mode**: Lightweight version for slow connections

## Testing Strategy

### Unit Tests
- **Component Rendering**: All components render correctly
- **Filter Logic**: Filtering functions work as expected
- **Event Handlers**: User interactions trigger correct actions
- **Data Transformation**: API data properly formatted

### Integration Tests
- **Map Integration**: Google Maps loads and displays correctly
- **Filter Integration**: Filters update map and list views
- **Navigation**: Tab switching and mobile navigation
- **Location Services**: GPS and postal code detection

### E2E Tests
- **User Flows**: Complete search and selection workflow
- **Mobile Experience**: Touch interactions and responsive design
- **Performance**: Load times and interaction responsiveness
- **Accessibility**: Screen reader and keyboard navigation

## Future Enhancements

### 1. Advanced Features
- **Appointment Booking**: Schedule visits with Seva Mitras
- **Live Chat**: Real-time messaging system
- **Reviews & Ratings**: User feedback system
- **Favorites**: Save frequently used Seva Mitras

### 2. Performance Improvements
- **PWA Support**: Progressive Web App capabilities
- **Offline Mode**: Cached data and functionality
- **Push Notifications**: Real-time updates and alerts
- **Background Sync**: Update data in background

### 3. Enhanced Localization
- **RTL Support**: Right-to-left language support
- **Regional Customization**: State-specific features
- **Currency Localization**: Regional currency display
- **Cultural Themes**: Festival and regional themes

This comprehensive UI system provides a world-class user experience for discovering and connecting with Seva Mitras, combining modern web technologies with cultural sensitivity and accessibility best practices.