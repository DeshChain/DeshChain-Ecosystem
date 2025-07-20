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

import React from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Switch,
  FormControlLabel,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Chip,
  Grid,
  IconButton,
  Divider,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  ListItemIcon,
  Avatar,
  useTheme
} from '@mui/material';
import {
  Celebration as CelebrationIcon,
  LocationOn as LocationIcon,
  Favorite as FavoriteIcon,
  FavoriteBorder as FavoriteBorderIcon,
  Palette as PaletteIcon,
  NotificationsActive as NotificationIcon
} from '@mui/icons-material';

import { useFestival } from '../../hooks/useFestival';
import { useLanguage } from '../../hooks/useLanguage';
import { FESTIVALS } from '../../themes/festivals';

export const FestivalSettings: React.FC = () => {
  const theme = useTheme();
  const { currentLanguage } = useLanguage();
  const {
    festivalThemeEnabled,
    toggleFestivalTheme,
    userRegion,
    setUserRegion,
    favoritesFestivals,
    toggleFavoriteFestival,
    currentFestival
  } = useFestival();

  const regions = [
    { value: 'all', label: 'All India' },
    { value: 'North India', label: 'North India' },
    { value: 'South India', label: 'South India' },
    { value: 'East India', label: 'East India' },
    { value: 'West India', label: 'West India' },
    { value: 'North East', label: 'North East' },
    { value: 'Central India', label: 'Central India' }
  ];

  return (
    <Grid container spacing={3}>
      {/* Theme Settings */}
      <Grid item xs={12} md={6}>
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center" gap={1} mb={2}>
              <PaletteIcon color="primary" />
              <Typography variant="h6">Festival Theme</Typography>
            </Box>
            
            <FormControlLabel
              control={
                <Switch
                  checked={festivalThemeEnabled}
                  onChange={(e) => toggleFestivalTheme(e.target.checked)}
                  color="primary"
                />
              }
              label="Enable festival themes"
            />
            
            <Typography variant="body2" color="textSecondary" sx={{ mt: 1 }}>
              Automatically apply festive themes and decorations during Indian festivals
            </Typography>
            
            {currentFestival && festivalThemeEnabled && (
              <Box mt={2} p={2} bgcolor="primary.light" borderRadius={1}>
                <Typography variant="body2" color="primary.contrastText">
                  🎉 {currentFestival.name} theme is active!
                </Typography>
              </Box>
            )}
          </CardContent>
        </Card>
      </Grid>

      {/* Region Settings */}
      <Grid item xs={12} md={6}>
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center" gap={1} mb={2}>
              <LocationIcon color="primary" />
              <Typography variant="h6">Region</Typography>
            </Box>
            
            <FormControl fullWidth>
              <InputLabel>Your Region</InputLabel>
              <Select
                value={userRegion}
                onChange={(e) => setUserRegion(e.target.value)}
                label="Your Region"
              >
                {regions.map(region => (
                  <MenuItem key={region.value} value={region.value}>
                    {region.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            
            <Typography variant="body2" color="textSecondary" sx={{ mt: 1 }}>
              Get notifications for festivals celebrated in your region
            </Typography>
          </CardContent>
        </Card>
      </Grid>

      {/* Favorite Festivals */}
      <Grid item xs={12}>
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center" gap={1} mb={2}>
              <CelebrationIcon color="primary" />
              <Typography variant="h6">Festival Preferences</Typography>
            </Box>
            
            <Typography variant="body2" color="textSecondary" gutterBottom>
              Select your favorite festivals to get special notifications and offers
            </Typography>
            
            <Divider sx={{ my: 2 }} />
            
            <List>
              {FESTIVALS.map(festival => {
                const isFavorite = favoritesFestivals.includes(festival.id);
                const localName = festival.localNames[currentLanguage] || festival.name;
                
                return (
                  <ListItem key={festival.id}>
                    <ListItemIcon>
                      <Avatar
                        sx={{
                          bgcolor: isFavorite ? 'primary.main' : 'grey.300',
                          width: 40,
                          height: 40
                        }}
                      >
                        <CelebrationIcon />
                      </Avatar>
                    </ListItemIcon>
                    
                    <ListItemText
                      primary={localName}
                      secondary={
                        <Box>
                          <Typography variant="caption" component="div">
                            {festival.name} • {new Date(festival.startDate).toLocaleDateString()}
                          </Typography>
                          <Box display="flex" gap={0.5} mt={0.5}>
                            <Chip
                              size="small"
                              label={festival.type}
                              variant="outlined"
                            />
                            {festival.regions.length > 0 && festival.regions[0] !== 'all' && (
                              <Chip
                                size="small"
                                label={festival.regions[0]}
                                variant="outlined"
                              />
                            )}
                          </Box>
                        </Box>
                      }
                    />
                    
                    <ListItemSecondaryAction>
                      <IconButton
                        edge="end"
                        onClick={() => toggleFavoriteFestival(festival.id)}
                        color={isFavorite ? 'primary' : 'default'}
                      >
                        {isFavorite ? <FavoriteIcon /> : <FavoriteBorderIcon />}
                      </IconButton>
                    </ListItemSecondaryAction>
                  </ListItem>
                );
              })}
            </List>
          </CardContent>
        </Card>
      </Grid>

      {/* Notification Settings */}
      <Grid item xs={12}>
        <Card>
          <CardContent>
            <Box display="flex" alignItems="center" gap={1} mb={2}>
              <NotificationIcon color="primary" />
              <Typography variant="h6">Festival Notifications</Typography>
            </Box>
            
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <FormControlLabel
                  control={<Switch defaultChecked />}
                  label="Festival reminders"
                />
                <Typography variant="caption" display="block" color="textSecondary">
                  Get notified 1 day before festivals
                </Typography>
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <FormControlLabel
                  control={<Switch defaultChecked />}
                  label="Festival offers"
                />
                <Typography variant="caption" display="block" color="textSecondary">
                  Receive special festival discount notifications
                </Typography>
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <FormControlLabel
                  control={<Switch defaultChecked />}
                  label="Festival greetings"
                />
                <Typography variant="caption" display="block" color="textSecondary">
                  Receive festival wishes in your language
                </Typography>
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <FormControlLabel
                  control={<Switch />}
                  label="Festival animations"
                />
                <Typography variant="caption" display="block" color="textSecondary">
                  Show particle effects and animations
                </Typography>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  );
};