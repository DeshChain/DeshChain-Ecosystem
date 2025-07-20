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

import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Stepper,
  Step,
  StepLabel,
  StepContent,
  TextField,
  Button,
  Alert,
  AlertTitle,
  FormControl,
  FormLabel,
  FormGroup,
  FormControlLabel,
  Checkbox,
  RadioGroup,
  Radio,
  InputAdornment,
  Select,
  MenuItem,
  Grid,
  Chip,
  LinearProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  useTheme,
  alpha
} from '@mui/material';
import {
  Store as StoreIcon,
  Business as BusinessIcon,
  LocationOn as LocationIcon,
  Schedule as ScheduleIcon,
  Payment as PaymentIcon,
  Security as SecurityIcon,
  Check as CheckIcon,
  Info as InfoIcon,
  AccountBalance as BankIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

import { PostalCodeInput } from './PostalCodeInput';
import { useLanguage } from '../hooks/useLanguage';

interface AgentRegistrationProps {
  onComplete: (agentData: any) => void;
}

interface OperatingHours {
  day: string;
  open: string;
  close: string;
  isClosed: boolean;
}

export const AgentRegistration: React.FC<AgentRegistrationProps> = ({ onComplete }) => {
  const theme = useTheme();
  const { t, formatCurrency, currentLanguage } = useLanguage();
  const [activeStep, setActiveStep] = useState(0);
  const [showInfoDialog, setShowInfoDialog] = useState(false);
  
  // Form data
  const [businessData, setBusinessData] = useState({
    businessName: '',
    registrationNumber: '',
    gstNumber: '',
    businessType: 'retail'
  });
  
  const [locationData, setLocationData] = useState({
    postalCode: '',
    fullAddress: '',
    landmark: '',
    latitude: '',
    longitude: ''
  });
  
  const [servicesData, setServicesData] = useState({
    cashIn: true,
    cashOut: true,
    remittance: false,
    billPayment: false
  });
  
  const [operatingHours, setOperatingHours] = useState<OperatingHours[]>([
    { day: 'Monday', open: '09:00', close: '18:00', isClosed: false },
    { day: 'Tuesday', open: '09:00', close: '18:00', isClosed: false },
    { day: 'Wednesday', open: '09:00', close: '18:00', isClosed: false },
    { day: 'Thursday', open: '09:00', close: '18:00', isClosed: false },
    { day: 'Friday', open: '09:00', close: '18:00', isClosed: false },
    { day: 'Saturday', open: '10:00', close: '16:00', isClosed: false },
    { day: 'Sunday', open: '', close: '', isClosed: true }
  ]);
  
  const [limitsData, setLimitsData] = useState({
    dailyLimit: '1000000',
    perTransactionLimit: '100000',
    securityDeposit: '100000'
  });
  
  const [bankData, setBankData] = useState({
    accountHolderName: '',
    accountNumber: '',
    confirmAccountNumber: '',
    ifscCode: '',
    bankName: '',
    branchName: ''
  });
  
  const [contactData, setContactData] = useState({
    phone: '',
    email: '',
    languages: ['hi', 'en']
  });
  
  const steps = [
    'Business Information',
    'Location Details',
    'Services Offered',
    'Operating Hours',
    'Transaction Limits',
    'Bank Details',
    'Contact & Languages',
    'Review & Submit'
  ];
  
  const handleNext = () => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };
  
  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };
  
  const handleOperatingHoursChange = (index: number, field: string, value: any) => {
    const newHours = [...operatingHours];
    newHours[index] = {
      ...newHours[index],
      [field]: value
    };
    setOperatingHours(newHours);
  };
  
  const validateStep = (step: number): boolean => {
    switch (step) {
      case 0:
        return businessData.businessName.length > 3 && 
               businessData.registrationNumber.length > 5;
      case 1:
        return locationData.postalCode.length === 6 && 
               locationData.fullAddress.length > 10;
      case 2:
        return Object.values(servicesData).some(v => v);
      case 3:
        return operatingHours.some(h => !h.isClosed);
      case 4:
        return parseInt(limitsData.dailyLimit) > 0 &&
               parseInt(limitsData.perTransactionLimit) > 0 &&
               parseInt(limitsData.securityDeposit) >= 100000;
      case 5:
        return bankData.accountNumber === bankData.confirmAccountNumber &&
               bankData.ifscCode.length === 11;
      case 6:
        return contactData.phone.length === 10 &&
               contactData.email.includes('@') &&
               contactData.languages.length > 0;
      default:
        return true;
    }
  };
  
  const handleSubmit = () => {
    const agentData = {
      ...businessData,
      ...locationData,
      services: Object.entries(servicesData)
        .filter(([_, enabled]) => enabled)
        .map(([service, _]) => service.toUpperCase()),
      operatingHours,
      ...limitsData,
      bankDetails: {
        ...bankData,
        accountNumber: undefined, // Don't include in final data
        confirmAccountNumber: undefined
      },
      ...contactData
    };
    
    onComplete(agentData);
  };
  
  const getStepContent = (step: number) => {
    switch (step) {
      case 0:
        return (
          <Box>
            <TextField
              fullWidth
              label="Business Name"
              value={businessData.businessName}
              onChange={(e) => setBusinessData({ ...businessData, businessName: e.target.value })}
              margin="normal"
              required
              helperText="Legal name of your business"
            />
            
            <TextField
              fullWidth
              label="Registration Number"
              value={businessData.registrationNumber}
              onChange={(e) => setBusinessData({ ...businessData, registrationNumber: e.target.value })}
              margin="normal"
              required
              helperText="Business registration or license number"
            />
            
            <TextField
              fullWidth
              label="GST Number (Optional)"
              value={businessData.gstNumber}
              onChange={(e) => setBusinessData({ ...businessData, gstNumber: e.target.value })}
              margin="normal"
              helperText="Required for registered businesses"
            />
            
            <FormControl fullWidth margin="normal">
              <FormLabel>Business Type</FormLabel>
              <RadioGroup
                value={businessData.businessType}
                onChange={(e) => setBusinessData({ ...businessData, businessType: e.target.value })}
              >
                <FormControlLabel value="retail" control={<Radio />} label="Retail Store" />
                <FormControlLabel value="kiosk" control={<Radio />} label="Kiosk/Counter" />
                <FormControlLabel value="mobile" control={<Radio />} label="Mobile Agent" />
                <FormControlLabel value="other" control={<Radio />} label="Other" />
              </RadioGroup>
            </FormControl>
          </Box>
        );
        
      case 1:
        return (
          <Box>
            <PostalCodeInput
              value={locationData.postalCode}
              onChange={(value) => setLocationData({ ...locationData, postalCode: value })}
              label="Business Postal Code"
              required
              autoDetect={false}
            />
            
            <TextField
              fullWidth
              multiline
              rows={3}
              label="Full Address"
              value={locationData.fullAddress}
              onChange={(e) => setLocationData({ ...locationData, fullAddress: e.target.value })}
              margin="normal"
              required
              helperText="Complete address including street, area"
            />
            
            <TextField
              fullWidth
              label="Landmark"
              value={locationData.landmark}
              onChange={(e) => setLocationData({ ...locationData, landmark: e.target.value })}
              margin="normal"
              helperText="Nearby landmark for easy location"
            />
            
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <TextField
                  fullWidth
                  label="Latitude (Optional)"
                  value={locationData.latitude}
                  onChange={(e) => setLocationData({ ...locationData, latitude: e.target.value })}
                  margin="normal"
                  type="number"
                />
              </Grid>
              <Grid item xs={6}>
                <TextField
                  fullWidth
                  label="Longitude (Optional)"
                  value={locationData.longitude}
                  onChange={(e) => setLocationData({ ...locationData, longitude: e.target.value })}
                  margin="normal"
                  type="number"
                />
              </Grid>
            </Grid>
            
            <Button
              variant="outlined"
              startIcon={<LocationIcon />}
              sx={{ mt: 2 }}
            >
              Get Current Location
            </Button>
          </Box>
        );
        
      case 2:
        return (
          <Box>
            <FormControl component="fieldset">
              <FormLabel component="legend">Select Services You'll Offer</FormLabel>
              <FormGroup>
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={servicesData.cashIn}
                      onChange={(e) => setServicesData({ ...servicesData, cashIn: e.target.checked })}
                    />
                  }
                  label={
                    <Box>
                      <Typography>Cash to NAMO (Cash In)</Typography>
                      <Typography variant="caption" color="text.secondary">
                        Accept cash and credit NAMO to customer wallet
                      </Typography>
                    </Box>
                  }
                />
                
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={servicesData.cashOut}
                      onChange={(e) => setServicesData({ ...servicesData, cashOut: e.target.checked })}
                    />
                  }
                  label={
                    <Box>
                      <Typography>NAMO to Cash (Cash Out)</Typography>
                      <Typography variant="caption" color="text.secondary">
                        Dispense cash for NAMO withdrawals
                      </Typography>
                    </Box>
                  }
                />
                
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={servicesData.remittance}
                      onChange={(e) => setServicesData({ ...servicesData, remittance: e.target.checked })}
                    />
                  }
                  label={
                    <Box>
                      <Typography>Money Transfer/Remittance</Typography>
                      <Typography variant="caption" color="text.secondary">
                        Send money to other locations
                      </Typography>
                    </Box>
                  }
                />
                
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={servicesData.billPayment}
                      onChange={(e) => setServicesData({ ...servicesData, billPayment: e.target.checked })}
                    />
                  }
                  label={
                    <Box>
                      <Typography>Bill Payments</Typography>
                      <Typography variant="caption" color="text.secondary">
                        Utility bills, recharges, etc.
                      </Typography>
                    </Box>
                  }
                />
              </FormGroup>
            </FormControl>
            
            <Alert severity="info" sx={{ mt: 2 }}>
              <AlertTitle>Commission Rates</AlertTitle>
              <Typography variant="body2">
                • Cash In/Out: 2-3% commission<br />
                • Remittance: 1-2% commission<br />
                • Bill Payment: Fixed ₹5-20 per transaction
              </Typography>
            </Alert>
          </Box>
        );
        
      case 3:
        return (
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              Set Your Operating Hours
            </Typography>
            
            {operatingHours.map((hours, index) => (
              <Box key={hours.day} sx={{ mb: 2 }}>
                <Grid container spacing={2} alignItems="center">
                  <Grid item xs={3}>
                    <Typography>{hours.day}</Typography>
                  </Grid>
                  
                  <Grid item xs={2}>
                    <FormControlLabel
                      control={
                        <Checkbox
                          checked={hours.isClosed}
                          onChange={(e) => handleOperatingHoursChange(index, 'isClosed', e.target.checked)}
                        />
                      }
                      label="Closed"
                    />
                  </Grid>
                  
                  {!hours.isClosed && (
                    <>
                      <Grid item xs={3}>
                        <TextField
                          type="time"
                          value={hours.open}
                          onChange={(e) => handleOperatingHoursChange(index, 'open', e.target.value)}
                          InputLabelProps={{ shrink: true }}
                          label="Open"
                          fullWidth
                        />
                      </Grid>
                      
                      <Grid item xs={3}>
                        <TextField
                          type="time"
                          value={hours.close}
                          onChange={(e) => handleOperatingHoursChange(index, 'close', e.target.value)}
                          InputLabelProps={{ shrink: true }}
                          label="Close"
                          fullWidth
                        />
                      </Grid>
                    </>
                  )}
                </Grid>
              </Box>
            ))}
            
            <Button
              variant="outlined"
              size="small"
              sx={{ mt: 1 }}
              onClick={() => {
                setOperatingHours(operatingHours.map(h => ({
                  ...h,
                  open: '09:00',
                  close: '18:00',
                  isClosed: h.day === 'Sunday'
                })));
              }}
            >
              Set Standard Hours
            </Button>
          </Box>
        );
        
      case 4:
        return (
          <Box>
            <TextField
              fullWidth
              label="Daily Transaction Limit"
              value={limitsData.dailyLimit}
              onChange={(e) => setLimitsData({ ...limitsData, dailyLimit: e.target.value })}
              margin="normal"
              type="number"
              required
              InputProps={{
                startAdornment: <InputAdornment position="start">₹</InputAdornment>
              }}
              helperText="Maximum daily transaction volume"
            />
            
            <TextField
              fullWidth
              label="Per Transaction Limit"
              value={limitsData.perTransactionLimit}
              onChange={(e) => setLimitsData({ ...limitsData, perTransactionLimit: e.target.value })}
              margin="normal"
              type="number"
              required
              InputProps={{
                startAdornment: <InputAdornment position="start">₹</InputAdornment>
              }}
              helperText="Maximum single transaction amount"
            />
            
            <TextField
              fullWidth
              label="Security Deposit"
              value={limitsData.securityDeposit}
              onChange={(e) => setLimitsData({ ...limitsData, securityDeposit: e.target.value })}
              margin="normal"
              type="number"
              required
              InputProps={{
                startAdornment: <InputAdornment position="start">₹</InputAdornment>
              }}
              helperText="Minimum ₹100,000 required (in NAMO tokens)"
            />
            
            <Alert severity="warning" sx={{ mt: 2 }}>
              <AlertTitle>Security Deposit Required</AlertTitle>
              <Typography variant="body2">
                A refundable security deposit of ₹{formatCurrency(parseInt(limitsData.securityDeposit))} 
                worth of NAMO tokens will be locked for the duration of your agent status.
              </Typography>
            </Alert>
          </Box>
        );
        
      case 5:
        return (
          <Box>
            <TextField
              fullWidth
              label="Account Holder Name"
              value={bankData.accountHolderName}
              onChange={(e) => setBankData({ ...bankData, accountHolderName: e.target.value })}
              margin="normal"
              required
              helperText="As per bank records"
            />
            
            <TextField
              fullWidth
              label="Account Number"
              value={bankData.accountNumber}
              onChange={(e) => setBankData({ ...bankData, accountNumber: e.target.value })}
              margin="normal"
              type="password"
              required
            />
            
            <TextField
              fullWidth
              label="Confirm Account Number"
              value={bankData.confirmAccountNumber}
              onChange={(e) => setBankData({ ...bankData, confirmAccountNumber: e.target.value })}
              margin="normal"
              required
              error={bankData.confirmAccountNumber !== '' && bankData.accountNumber !== bankData.confirmAccountNumber}
              helperText={
                bankData.confirmAccountNumber !== '' && bankData.accountNumber !== bankData.confirmAccountNumber
                  ? "Account numbers don't match"
                  : ""
              }
            />
            
            <TextField
              fullWidth
              label="IFSC Code"
              value={bankData.ifscCode}
              onChange={(e) => setBankData({ ...bankData, ifscCode: e.target.value.toUpperCase() })}
              margin="normal"
              required
              inputProps={{ maxLength: 11 }}
              helperText="11-character bank IFSC code"
            />
            
            <TextField
              fullWidth
              label="Bank Name"
              value={bankData.bankName}
              onChange={(e) => setBankData({ ...bankData, bankName: e.target.value })}
              margin="normal"
              required
            />
            
            <TextField
              fullWidth
              label="Branch Name"
              value={bankData.branchName}
              onChange={(e) => setBankData({ ...bankData, branchName: e.target.value })}
              margin="normal"
              required
            />
            
            <Alert severity="info" sx={{ mt: 2 }}>
              Daily settlements will be made to this account after deducting commission.
            </Alert>
          </Box>
        );
        
      case 6:
        return (
          <Box>
            <TextField
              fullWidth
              label="Mobile Number"
              value={contactData.phone}
              onChange={(e) => setContactData({ ...contactData, phone: e.target.value })}
              margin="normal"
              required
              InputProps={{
                startAdornment: <InputAdornment position="start">+91</InputAdornment>
              }}
              inputProps={{ maxLength: 10 }}
              helperText="For customer support and notifications"
            />
            
            <TextField
              fullWidth
              label="Email Address"
              value={contactData.email}
              onChange={(e) => setContactData({ ...contactData, email: e.target.value })}
              margin="normal"
              type="email"
              required
              helperText="For important updates and reports"
            />
            
            <FormControl fullWidth margin="normal">
              <FormLabel>Languages Spoken</FormLabel>
              <Select
                multiple
                value={contactData.languages}
                onChange={(e) => setContactData({ ...contactData, languages: e.target.value as string[] })}
                renderValue={(selected) => (
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                    {selected.map((value) => (
                      <Chip key={value} label={value.toUpperCase()} size="small" />
                    ))}
                  </Box>
                )}
              >
                <MenuItem value="hi">Hindi</MenuItem>
                <MenuItem value="en">English</MenuItem>
                <MenuItem value="bn">Bengali</MenuItem>
                <MenuItem value="te">Telugu</MenuItem>
                <MenuItem value="mr">Marathi</MenuItem>
                <MenuItem value="ta">Tamil</MenuItem>
                <MenuItem value="gu">Gujarati</MenuItem>
                <MenuItem value="ur">Urdu</MenuItem>
                <MenuItem value="kn">Kannada</MenuItem>
                <MenuItem value="or">Odia</MenuItem>
                <MenuItem value="ml">Malayalam</MenuItem>
                <MenuItem value="pa">Punjabi</MenuItem>
              </Select>
              <Typography variant="caption" color="text.secondary" sx={{ mt: 1 }}>
                Select all languages you can provide service in
              </Typography>
            </FormControl>
          </Box>
        );
        
      case 7:
        return (
          <Box>
            <Alert severity="success" sx={{ mb: 3 }}>
              <AlertTitle>Review Your Information</AlertTitle>
              Please review all details before submitting your application.
            </Alert>
            
            <Box sx={{ mb: 3 }}>
              <Typography variant="subtitle2" color="primary" gutterBottom>
                Business Information
              </Typography>
              <Typography variant="body2">
                {businessData.businessName}<br />
                Registration: {businessData.registrationNumber}<br />
                {businessData.gstNumber && `GST: ${businessData.gstNumber}`}
              </Typography>
            </Box>
            
            <Box sx={{ mb: 3 }}>
              <Typography variant="subtitle2" color="primary" gutterBottom>
                Location
              </Typography>
              <Typography variant="body2">
                {locationData.fullAddress}<br />
                Postal Code: {locationData.postalCode}
              </Typography>
            </Box>
            
            <Box sx={{ mb: 3 }}>
              <Typography variant="subtitle2" color="primary" gutterBottom>
                Services & Limits
              </Typography>
              <Typography variant="body2">
                Services: {Object.entries(servicesData)
                  .filter(([_, enabled]) => enabled)
                  .map(([service, _]) => service)
                  .join(', ')}<br />
                Daily Limit: ₹{formatCurrency(parseInt(limitsData.dailyLimit))}<br />
                Security Deposit: ₹{formatCurrency(parseInt(limitsData.securityDeposit))}
              </Typography>
            </Box>
            
            <Box sx={{ mb: 3 }}>
              <Typography variant="subtitle2" color="primary" gutterBottom>
                Contact
              </Typography>
              <Typography variant="body2">
                Phone: +91 {contactData.phone}<br />
                Email: {contactData.email}<br />
                Languages: {contactData.languages.join(', ')}
              </Typography>
            </Box>
            
            <FormControlLabel
              control={<Checkbox required />}
              label="I agree to the terms and conditions and confirm all information is accurate"
            />
          </Box>
        );
        
      default:
        return 'Unknown step';
    }
  };
  
  return (
    <Card>
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Typography variant="h5">
            Become a DeshChain Agent
          </Typography>
          <Button
            startIcon={<InfoIcon />}
            onClick={() => setShowInfoDialog(true)}
          >
            Benefits
          </Button>
        </Box>
        
        <Stepper activeStep={activeStep} orientation="vertical">
          {steps.map((label, index) => (
            <Step key={label}>
              <StepLabel>{label}</StepLabel>
              <StepContent>
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.3 }}
                >
                  {getStepContent(index)}
                </motion.div>
                
                <Box sx={{ mt: 3 }}>
                  <Button
                    variant="contained"
                    onClick={index === steps.length - 1 ? handleSubmit : handleNext}
                    disabled={!validateStep(index)}
                    sx={{ mr: 1 }}
                  >
                    {index === steps.length - 1 ? 'Submit Application' : 'Continue'}
                  </Button>
                  <Button
                    disabled={index === 0}
                    onClick={handleBack}
                  >
                    Back
                  </Button>
                </Box>
              </StepContent>
            </Step>
          ))}
        </Stepper>
        
        {activeStep === steps.length && (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <CheckIcon sx={{ fontSize: 60, color: 'success.main', mb: 2 }} />
            <Typography variant="h6" gutterBottom>
              Application Submitted Successfully!
            </Typography>
            <Typography color="text.secondary">
              Your agent registration is under review. You'll receive KYC verification instructions within 24 hours.
            </Typography>
          </Box>
        )}
      </CardContent>
      
      {/* Benefits Dialog */}
      <Dialog
        open={showInfoDialog}
        onClose={() => setShowInfoDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Agent Benefits</DialogTitle>
        <DialogContent>
          <Box display="flex" flexDirection="column" gap={2}>
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                💰 Earn Commission
              </Typography>
              <Typography variant="body2" color="text.secondary">
                2-3% on cash transactions, 1-2% on remittances
              </Typography>
            </Box>
            
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                📈 Build Your Business
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Attract more customers with digital financial services
              </Typography>
            </Box>
            
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                🛡️ Trust & Security
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Blockchain-backed transactions with full transparency
              </Typography>
            </Box>
            
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                🎯 Marketing Support
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Branded materials and listing on DeshChain app
              </Typography>
            </Box>
            
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                🏆 Performance Rewards
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Bonuses for high volume and customer satisfaction
              </Typography>
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowInfoDialog(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </Card>
  );
};