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

import axios from 'axios';

export interface NotificationPreferences {
  phoneNumber: string;
  enableSMS: boolean;
  enableWhatsApp: boolean;
  language: string;
  enableEmail?: boolean;
  email?: string;
}

export interface NotificationTemplate {
  id: string;
  name: string;
  type: 'sms' | 'whatsapp' | 'email';
  language: string;
  template: string;
  variables: string[];
}

export interface NotificationRequest {
  recipient: string;
  type: 'sms' | 'whatsapp' | 'email';
  template: string;
  data: Record<string, any>;
  language?: string;
  priority?: 'high' | 'normal' | 'low';
}

export interface NotificationStatus {
  id: string;
  status: 'pending' | 'sent' | 'delivered' | 'failed';
  timestamp: string;
  error?: string;
  deliveredAt?: string;
}

export class NotificationService {
  private baseUrl: string;
  private apiKey: string;
  private preferences: NotificationPreferences | null = null;

  constructor(baseUrl: string, apiKey: string) {
    this.baseUrl = baseUrl;
    this.apiKey = apiKey;
    this.loadPreferences();
  }

  // Load user preferences from localStorage
  private loadPreferences() {
    try {
      const saved = localStorage.getItem('deshchain-notification-preferences');
      if (saved) {
        this.preferences = JSON.parse(saved);
      }
    } catch (error) {
      console.error('Failed to load notification preferences:', error);
    }
  }

  // Save user preferences
  public savePreferences(preferences: NotificationPreferences) {
    this.preferences = preferences;
    localStorage.setItem('deshchain-notification-preferences', JSON.stringify(preferences));
  }

  // Get current preferences
  public getPreferences(): NotificationPreferences | null {
    return this.preferences;
  }

  // Send SMS notification
  public async sendSMS(
    phoneNumber: string,
    message: string,
    language: string = 'en'
  ): Promise<NotificationStatus> {
    try {
      const response = await axios.post(
        `${this.baseUrl}/notifications/sms`,
        {
          to: phoneNumber,
          message,
          language
        },
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`,
            'Content-Type': 'application/json'
          }
        }
      );

      return response.data;
    } catch (error) {
      console.error('Failed to send SMS:', error);
      throw error;
    }
  }

  // Send WhatsApp message
  public async sendWhatsApp(
    phoneNumber: string,
    template: string,
    data: Record<string, any>,
    mediaUrl?: string
  ): Promise<NotificationStatus> {
    try {
      const response = await axios.post(
        `${this.baseUrl}/notifications/whatsapp`,
        {
          to: phoneNumber,
          template,
          data,
          mediaUrl
        },
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`,
            'Content-Type': 'application/json'
          }
        }
      );

      return response.data;
    } catch (error) {
      console.error('Failed to send WhatsApp message:', error);
      throw error;
    }
  }

  // Send money order receipt notification
  public async sendMoneyOrderReceipt(
    receiptId: string,
    recipient: string,
    type: 'sms' | 'whatsapp' = 'whatsapp'
  ): Promise<NotificationStatus> {
    try {
      const response = await axios.post(
        `${this.baseUrl}/notifications/receipt`,
        {
          receiptId,
          recipient,
          type,
          language: this.preferences?.language || 'en'
        },
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`,
            'Content-Type': 'application/json'
          }
        }
      );

      return response.data;
    } catch (error) {
      console.error('Failed to send receipt notification:', error);
      throw error;
    }
  }

  // Send OTP for verification
  public async sendOTP(
    phoneNumber: string,
    purpose: 'transaction' | 'login' | 'registration'
  ): Promise<{ otp_id: string; expires_at: string }> {
    try {
      const response = await axios.post(
        `${this.baseUrl}/notifications/otp`,
        {
          phoneNumber,
          purpose,
          language: this.preferences?.language || 'en'
        },
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`,
            'Content-Type': 'application/json'
          }
        }
      );

      return response.data;
    } catch (error) {
      console.error('Failed to send OTP:', error);
      throw error;
    }
  }

  // Verify OTP
  public async verifyOTP(
    otpId: string,
    otp: string
  ): Promise<{ verified: boolean; error?: string }> {
    try {
      const response = await axios.post(
        `${this.baseUrl}/notifications/otp/verify`,
        {
          otp_id: otpId,
          otp
        },
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`,
            'Content-Type': 'application/json'
          }
        }
      );

      return response.data;
    } catch (error) {
      console.error('Failed to verify OTP:', error);
      return { verified: false, error: 'Verification failed' };
    }
  }

  // Get notification status
  public async getNotificationStatus(notificationId: string): Promise<NotificationStatus> {
    try {
      const response = await axios.get(
        `${this.baseUrl}/notifications/status/${notificationId}`,
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`
          }
        }
      );

      return response.data;
    } catch (error) {
      console.error('Failed to get notification status:', error);
      throw error;
    }
  }

  // Send bulk notifications
  public async sendBulkNotifications(
    notifications: NotificationRequest[]
  ): Promise<NotificationStatus[]> {
    try {
      const response = await axios.post(
        `${this.baseUrl}/notifications/bulk`,
        { notifications },
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`,
            'Content-Type': 'application/json'
          }
        }
      );

      return response.data;
    } catch (error) {
      console.error('Failed to send bulk notifications:', error);
      throw error;
    }
  }

  // Get available templates
  public async getTemplates(
    type?: 'sms' | 'whatsapp' | 'email',
    language?: string
  ): Promise<NotificationTemplate[]> {
    try {
      const params = new URLSearchParams();
      if (type) params.append('type', type);
      if (language) params.append('language', language);

      const response = await axios.get(
        `${this.baseUrl}/notifications/templates?${params}`,
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`
          }
        }
      );

      return response.data;
    } catch (error) {
      console.error('Failed to get templates:', error);
      throw error;
    }
  }

  // Format phone number for Indian format
  public static formatPhoneNumber(phoneNumber: string): string {
    // Remove all non-digits
    const digits = phoneNumber.replace(/\D/g, '');

    // Handle Indian phone numbers
    if (digits.length === 10) {
      return `+91${digits}`;
    } else if (digits.length === 12 && digits.startsWith('91')) {
      return `+${digits}`;
    } else if (digits.length === 13 && digits.startsWith('+91')) {
      return digits;
    }

    // Return as is if format is unknown
    return phoneNumber;
  }

  // Validate phone number
  public static isValidPhoneNumber(phoneNumber: string): boolean {
    const formatted = NotificationService.formatPhoneNumber(phoneNumber);
    const indianPhoneRegex = /^\+91[6-9]\d{9}$/;
    return indianPhoneRegex.test(formatted);
  }

  // Get notification history
  public async getNotificationHistory(
    limit: number = 50,
    offset: number = 0
  ): Promise<NotificationStatus[]> {
    try {
      const response = await axios.get(
        `${this.baseUrl}/notifications/history`,
        {
          params: { limit, offset },
          headers: {
            'Authorization': `Bearer ${this.apiKey}`
          }
        }
      );

      return response.data;
    } catch (error) {
      console.error('Failed to get notification history:', error);
      throw error;
    }
  }

  // Subscribe to push notifications
  public async subscribePushNotifications(
    subscription: PushSubscription
  ): Promise<{ success: boolean }> {
    try {
      const response = await axios.post(
        `${this.baseUrl}/notifications/push/subscribe`,
        { subscription },
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`,
            'Content-Type': 'application/json'
          }
        }
      );

      return response.data;
    } catch (error) {
      console.error('Failed to subscribe to push notifications:', error);
      throw error;
    }
  }

  // Generate message from template
  public static generateMessage(
    template: string,
    data: Record<string, any>,
    language: string = 'en'
  ): string {
    let message = template;

    // Replace template variables
    Object.keys(data).forEach(key => {
      const regex = new RegExp(`{{${key}}}`, 'g');
      message = message.replace(regex, data[key]);
    });

    // Add cultural elements based on language
    if (language === 'hi') {
      message += '\n\n🙏 धन्यवाद';
    } else if (language === 'bn') {
      message += '\n\n🙏 ধন্যবাদ';
    } else {
      message += '\n\n🙏 Thank you';
    }

    return message;
  }
}