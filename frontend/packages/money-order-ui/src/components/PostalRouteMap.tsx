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

import React, { useEffect, useRef, useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Chip,
  IconButton,
  Tooltip,
  useTheme,
  alpha
} from '@mui/material';
import {
  ZoomIn as ZoomInIcon,
  ZoomOut as ZoomOutIcon,
  Fullscreen as FullscreenIcon,
  MyLocation as CenterIcon,
  Route as RouteIcon
} from '@mui/icons-material';

import { PostalRoute, PostalCodeInfo } from '../services/postalCodeService';

interface PostalRouteMapProps {
  route: PostalRoute;
  height?: number;
  interactive?: boolean;
  showControls?: boolean;
}

export const PostalRouteMap: React.FC<PostalRouteMapProps> = ({
  route,
  height = 400,
  interactive = true,
  showControls = true
}) => {
  const theme = useTheme();
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [zoom, setZoom] = useState(1);
  const [pan, setPan] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });

  // India map bounds (simplified)
  const mapBounds = {
    north: 37.0,
    south: 8.0,
    east: 97.0,
    west: 68.0
  };

  // Convert lat/lng to canvas coordinates
  const latLngToCanvas = (lat: number, lng: number, canvas: HTMLCanvasElement) => {
    const x = ((lng - mapBounds.west) / (mapBounds.east - mapBounds.west)) * canvas.width;
    const y = ((mapBounds.north - lat) / (mapBounds.north - mapBounds.south)) * canvas.height;
    
    return {
      x: x * zoom + pan.x,
      y: y * zoom + pan.y
    };
  };

  // Draw the route map
  const drawMap = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Draw background
    ctx.fillStyle = theme.palette.background.paper;
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // Draw India outline (simplified)
    ctx.strokeStyle = alpha(theme.palette.divider, 0.3);
    ctx.lineWidth = 1;
    ctx.beginPath();
    
    // Very simplified India outline
    const indiaOutline = [
      { lat: 35, lng: 75 }, // Kashmir
      { lat: 32, lng: 78 }, 
      { lat: 28, lng: 77 }, // Delhi
      { lat: 23, lng: 70 }, // Gujarat
      { lat: 20, lng: 73 }, // Mumbai
      { lat: 15, lng: 74 }, // Goa
      { lat: 8, lng: 77 },  // Kerala
      { lat: 10, lng: 79 }, // Tamil Nadu
      { lat: 13, lng: 80 }, // Chennai
      { lat: 20, lng: 85 }, // Odisha
      { lat: 22, lng: 88 }, // Kolkata
      { lat: 27, lng: 88 }, // Sikkim
      { lat: 28, lng: 94 }, // Arunachal
      { lat: 25, lng: 92 }, // Assam
      { lat: 24, lng: 88 }, // Bangladesh border
      { lat: 27, lng: 85 }, // Nepal border
      { lat: 30, lng: 79 }, // Uttarakhand
      { lat: 35, lng: 75 }  // Back to Kashmir
    ];

    indiaOutline.forEach((point, index) => {
      const { x, y } = latLngToCanvas(point.lat, point.lng, canvas);
      if (index === 0) {
        ctx.moveTo(x, y);
      } else {
        ctx.lineTo(x, y);
      }
    });
    ctx.closePath();
    ctx.stroke();

    // Draw route
    drawRoute(ctx, canvas);

    // Draw locations
    drawLocation(ctx, canvas, route.from, 'origin');
    if (route.hubs) {
      route.hubs.forEach(hub => drawLocation(ctx, canvas, hub, 'hub'));
    }
    drawLocation(ctx, canvas, route.to, 'destination');
  };

  // Draw route line
  const drawRoute = (ctx: CanvasRenderingContext2D, canvas: HTMLCanvasElement) => {
    ctx.save();

    // Route styling based on type
    if (route.routeType === 'direct') {
      ctx.strokeStyle = theme.palette.success.main;
      ctx.lineWidth = 3;
      ctx.setLineDash([]);
    } else {
      ctx.strokeStyle = theme.palette.primary.main;
      ctx.lineWidth = 2;
      ctx.setLineDash([10, 5]);
    }

    // Draw route segments
    const points: PostalCodeInfo[] = [route.from];
    if (route.hubs) points.push(...route.hubs);
    points.push(route.to);

    ctx.beginPath();
    points.forEach((point, index) => {
      const { x, y } = latLngToCanvas(
        point.latitude || 20 + index * 5,
        point.longitude || 75 + index * 5,
        canvas
      );

      if (index === 0) {
        ctx.moveTo(x, y);
      } else {
        // Draw curved line for better visualization
        const prevPoint = points[index - 1];
        const prevCoords = latLngToCanvas(
          prevPoint.latitude || 20 + (index - 1) * 5,
          prevPoint.longitude || 75 + (index - 1) * 5,
          canvas
        );
        
        const cp1x = prevCoords.x + (x - prevCoords.x) * 0.3;
        const cp1y = prevCoords.y - 50 * zoom;
        const cp2x = prevCoords.x + (x - prevCoords.x) * 0.7;
        const cp2y = y - 50 * zoom;
        
        ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y);
      }
    });
    ctx.stroke();
    ctx.restore();
  };

  // Draw location marker
  const drawLocation = (
    ctx: CanvasRenderingContext2D,
    canvas: HTMLCanvasElement,
    location: PostalCodeInfo,
    type: 'origin' | 'hub' | 'destination'
  ) => {
    const { x, y } = latLngToCanvas(
      location.latitude || 20,
      location.longitude || 75,
      canvas
    );

    ctx.save();

    // Marker styling
    let color = theme.palette.primary.main;
    let size = 8;
    
    if (type === 'origin') {
      color = theme.palette.success.main;
      size = 10;
    } else if (type === 'destination') {
      color = theme.palette.error.main;
      size = 10;
    }

    // Draw marker
    ctx.fillStyle = color;
    ctx.strokeStyle = theme.palette.background.paper;
    ctx.lineWidth = 2;
    
    ctx.beginPath();
    ctx.arc(x, y, size * zoom, 0, Math.PI * 2);
    ctx.fill();
    ctx.stroke();

    // Draw pin effect
    if (type !== 'hub') {
      ctx.beginPath();
      ctx.moveTo(x, y + size * zoom);
      ctx.lineTo(x - size * zoom * 0.5, y + size * zoom * 2);
      ctx.lineTo(x + size * zoom * 0.5, y + size * zoom * 2);
      ctx.closePath();
      ctx.fill();
    }

    // Draw label
    ctx.fillStyle = theme.palette.text.primary;
    ctx.font = `${12 * zoom}px ${theme.typography.fontFamily}`;
    ctx.textAlign = 'center';
    ctx.textBaseline = 'bottom';
    
    const label = location.pincode;
    const labelY = type === 'hub' ? y - size * zoom - 5 : y + size * zoom * 2 + 15;
    
    // Background for text
    const textWidth = ctx.measureText(label).width;
    ctx.fillStyle = alpha(theme.palette.background.paper, 0.9);
    ctx.fillRect(
      x - textWidth / 2 - 4,
      labelY - 14 * zoom,
      textWidth + 8,
      16 * zoom
    );
    
    // Text
    ctx.fillStyle = theme.palette.text.primary;
    ctx.fillText(label, x, labelY);

    ctx.restore();
  };

  // Handle mouse events
  const handleMouseDown = (e: React.MouseEvent) => {
    if (!interactive) return;
    setIsDragging(true);
    setDragStart({ x: e.clientX - pan.x, y: e.clientY - pan.y });
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!interactive || !isDragging) return;
    setPan({
      x: e.clientX - dragStart.x,
      y: e.clientY - dragStart.y
    });
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  const handleWheel = (e: React.WheelEvent) => {
    if (!interactive) return;
    e.preventDefault();
    const delta = e.deltaY > 0 ? 0.9 : 1.1;
    setZoom(prev => Math.max(0.5, Math.min(3, prev * delta)));
  };

  // Zoom controls
  const handleZoomIn = () => setZoom(prev => Math.min(3, prev * 1.2));
  const handleZoomOut = () => setZoom(prev => Math.max(0.5, prev / 1.2));
  const handleReset = () => {
    setZoom(1);
    setPan({ x: 0, y: 0 });
  };

  // Draw on canvas
  useEffect(() => {
    drawMap();
  }, [route, zoom, pan, theme]);

  // Resize canvas
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const updateSize = () => {
      const rect = canvas.getBoundingClientRect();
      canvas.width = rect.width;
      canvas.height = rect.height;
      drawMap();
    };

    updateSize();
    window.addEventListener('resize', updateSize);
    return () => window.removeEventListener('resize', updateSize);
  }, []);

  return (
    <Card>
      <CardContent sx={{ p: 0, position: 'relative' }}>
        {/* Canvas */}
        <Box
          sx={{
            position: 'relative',
            height,
            overflow: 'hidden',
            cursor: isDragging ? 'grabbing' : interactive ? 'grab' : 'default'
          }}
        >
          <canvas
            ref={canvasRef}
            style={{ width: '100%', height: '100%' }}
            onMouseDown={handleMouseDown}
            onMouseMove={handleMouseMove}
            onMouseUp={handleMouseUp}
            onMouseLeave={handleMouseUp}
            onWheel={handleWheel}
          />

          {/* Route Info Overlay */}
          <Box
            sx={{
              position: 'absolute',
              top: 16,
              left: 16,
              display: 'flex',
              gap: 1,
              flexWrap: 'wrap'
            }}
          >
            <Chip
              icon={<RouteIcon />}
              label={`${route.distance} km`}
              size="small"
              sx={{ bgcolor: alpha(theme.palette.background.paper, 0.9) }}
            />
            <Chip
              label={`${route.estimatedTime}h`}
              size="small"
              color={route.priority === 'express' ? 'error' : 'primary'}
              sx={{ bgcolor: alpha(theme.palette.background.paper, 0.9) }}
            />
          </Box>

          {/* Controls */}
          {showControls && (
            <Box
              sx={{
                position: 'absolute',
                bottom: 16,
                right: 16,
                display: 'flex',
                flexDirection: 'column',
                gap: 1
              }}
            >
              <Tooltip title="Zoom In">
                <IconButton
                  size="small"
                  onClick={handleZoomIn}
                  sx={{ bgcolor: alpha(theme.palette.background.paper, 0.9) }}
                >
                  <ZoomInIcon />
                </IconButton>
              </Tooltip>
              <Tooltip title="Zoom Out">
                <IconButton
                  size="small"
                  onClick={handleZoomOut}
                  sx={{ bgcolor: alpha(theme.palette.background.paper, 0.9) }}
                >
                  <ZoomOutIcon />
                </IconButton>
              </Tooltip>
              <Tooltip title="Reset View">
                <IconButton
                  size="small"
                  onClick={handleReset}
                  sx={{ bgcolor: alpha(theme.palette.background.paper, 0.9) }}
                >
                  <CenterIcon />
                </IconButton>
              </Tooltip>
            </Box>
          )}
        </Box>

        {/* Legend */}
        <Box sx={{ p: 2, borderTop: 1, borderColor: 'divider' }}>
          <Box display="flex" gap={3} justifyContent="center">
            <Box display="flex" alignItems="center" gap={0.5}>
              <Box
                sx={{
                  width: 12,
                  height: 12,
                  borderRadius: '50%',
                  bgcolor: 'success.main'
                }}
              />
              <Typography variant="caption">Origin</Typography>
            </Box>
            <Box display="flex" alignItems="center" gap={0.5}>
              <Box
                sx={{
                  width: 12,
                  height: 12,
                  borderRadius: '50%',
                  bgcolor: 'primary.main'
                }}
              />
              <Typography variant="caption">Hub</Typography>
            </Box>
            <Box display="flex" alignItems="center" gap={0.5}>
              <Box
                sx={{
                  width: 12,
                  height: 12,
                  borderRadius: '50%',
                  bgcolor: 'error.main'
                }}
              />
              <Typography variant="caption">Destination</Typography>
            </Box>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};