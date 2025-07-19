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

package types

import (
	"fmt"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MoneyOrderReceipt represents a digital receipt for Money Order transactions
// Designed to feel like traditional post office money order receipts
type MoneyOrderReceipt struct {
	// Basic Information (like traditional money order form)
	OrderId          string         `json:"order_id" yaml:"order_id"`
	ReferenceNumber  string         `json:"reference_number" yaml:"reference_number"` // MO-2024-DEL-001234
	TransactionType  string         `json:"transaction_type" yaml:"transaction_type"` // "instant", "normal", "scheduled"
	
	// Sender Details (भेजने वाला)
	SenderAddress    sdk.AccAddress `json:"sender_address" yaml:"sender_address"`
	SenderName       string         `json:"sender_name" yaml:"sender_name"`
	SenderMobile     string         `json:"sender_mobile" yaml:"sender_mobile"`
	SenderPostalCode string         `json:"sender_postal_code" yaml:"sender_postal_code"`
	
	// Receiver Details (प्राप्तकर्ता)
	ReceiverAddress    sdk.AccAddress `json:"receiver_address" yaml:"receiver_address"`
	ReceiverName       string         `json:"receiver_name" yaml:"receiver_name"`
	ReceiverMobile     string         `json:"receiver_mobile" yaml:"receiver_mobile"`
	ReceiverPostalCode string         `json:"receiver_postal_code" yaml:"receiver_postal_code"`
	
	// Transaction Details (लेनदेन विवरण)
	Amount           sdk.Coin       `json:"amount" yaml:"amount"`
	Fees             sdk.Coin       `json:"fees" yaml:"fees"`
	TotalAmount      sdk.Coin       `json:"total_amount" yaml:"total_amount"`
	ExchangeRate     sdk.Dec        `json:"exchange_rate" yaml:"exchange_rate"`
	
	// UPI-Style Simple Message
	Note             string         `json:"note" yaml:"note"` // Personal message like "Beta ke liye"
	Purpose          string         `json:"purpose" yaml:"purpose"` // "family", "business", "emergency"
	
	// Status Tracking
	Status           string         `json:"status" yaml:"status"`
	StatusMessage    string         `json:"status_message" yaml:"status_message"`
	
	// Cultural Touch
	CulturalQuote    string         `json:"cultural_quote" yaml:"cultural_quote"`
	FestivalGreeting string         `json:"festival_greeting" yaml:"festival_greeting"`
	Language         string         `json:"language" yaml:"language"`
	
	// Digital Features
	QRCode           string         `json:"qr_code" yaml:"qr_code"`
	TrackingURL      string         `json:"tracking_url" yaml:"tracking_url"`
	SMSNotification  bool           `json:"sms_notification" yaml:"sms_notification"`
	EmailNotification bool          `json:"email_notification" yaml:"email_notification"`
	
	// Timestamps
	CreatedAt        time.Time      `json:"created_at" yaml:"created_at"`
	ProcessedAt      time.Time      `json:"processed_at" yaml:"processed_at"`
	CompletedAt      time.Time      `json:"completed_at" yaml:"completed_at"`
	ExpiresAt        time.Time      `json:"expires_at" yaml:"expires_at"`
}

// UPIStyleTransfer represents a simplified UPI-like transfer request
type UPIStyleTransfer struct {
	// Simple Fields (like UPI apps)
	ReceiverUPI      string         `json:"receiver_upi" yaml:"receiver_upi"`       // name@deshchain or mobile@deshchain
	Amount           string         `json:"amount" yaml:"amount"`                   // Simple amount like "100"
	Note             string         `json:"note" yaml:"note"`                       // Optional message
	
	// Hidden complexity (handled internally)
	SenderAddress    sdk.AccAddress `json:"sender_address" yaml:"sender_address"`
	ReceiverAddress  sdk.AccAddress `json:"receiver_address" yaml:"receiver_address"`
	Fees             sdk.Coin       `json:"fees" yaml:"fees"`
}

// SimplifiedOrderStatus for UPI-style display
type SimplifiedOrderStatus struct {
	// Visual Status (with emojis for mobile)
	StatusIcon       string         `json:"status_icon" yaml:"status_icon"`         // ✓, ⏳, ❌
	StatusText       string         `json:"status_text" yaml:"status_text"`         // "Success", "Processing", "Failed"
	StatusColor      string         `json:"status_color" yaml:"status_color"`       // "green", "yellow", "red"
	
	// Simple Details
	Amount           string         `json:"amount" yaml:"amount"`                   // "₹1,000"
	ReceiverName     string         `json:"receiver_name" yaml:"receiver_name"`     // "Ramesh Kumar"
	Time             string         `json:"time" yaml:"time"`                       // "2 mins ago"
	
	// Action Buttons
	Actions          []OrderAction  `json:"actions" yaml:"actions"`
}

// OrderAction represents actions available for an order
type OrderAction struct {
	Label            string         `json:"label" yaml:"label"`                     // "Repeat", "Share", "Download"
	Icon             string         `json:"icon" yaml:"icon"`                       // Icon name
	ActionType       string         `json:"action_type" yaml:"action_type"`         // "repeat", "share", "download"
	Enabled          bool           `json:"enabled" yaml:"enabled"`
}

// QuickTransferTemplate for frequent transfers (like UPI)
type QuickTransferTemplate struct {
	TemplateId       string         `json:"template_id" yaml:"template_id"`
	TemplateName     string         `json:"template_name" yaml:"template_name"`     // "Send to Mom", "Pay Rent"
	IconEmoji        string         `json:"icon_emoji" yaml:"icon_emoji"`           // 👩‍👦, 🏠
	
	// Pre-filled Details
	ReceiverAddress  sdk.AccAddress `json:"receiver_address" yaml:"receiver_address"`
	ReceiverName     string         `json:"receiver_name" yaml:"receiver_name"`
	DefaultAmount    sdk.Int        `json:"default_amount" yaml:"default_amount"`
	DefaultNote      string         `json:"default_note" yaml:"default_note"`
	
	// Usage Stats
	UsageCount       uint64         `json:"usage_count" yaml:"usage_count"`
	LastUsed         time.Time      `json:"last_used" yaml:"last_used"`
}

// VirtualPaymentAddress for UPI-style addressing
type VirtualPaymentAddress struct {
	VPA              string         `json:"vpa" yaml:"vpa"`                         // ramesh@deshchain
	LinkedAddress    sdk.AccAddress `json:"linked_address" yaml:"linked_address"`
	DisplayName      string         `json:"display_name" yaml:"display_name"`
	ProfilePicture   string         `json:"profile_picture" yaml:"profile_picture"`
	VerifiedMerchant bool           `json:"verified_merchant" yaml:"verified_merchant"`
	CreatedAt        time.Time      `json:"created_at" yaml:"created_at"`
}

// NewMoneyOrderReceipt creates a new receipt with UPI-style simplicity
func NewMoneyOrderReceipt(
	orderId string,
	sender sdk.AccAddress,
	receiver sdk.AccAddress,
	amount sdk.Coin,
	note string,
) *MoneyOrderReceipt {
	now := time.Now()
	
	return &MoneyOrderReceipt{
		OrderId:         orderId,
		ReferenceNumber: GenerateReferenceNumber(now),
		TransactionType: "instant",
		
		SenderAddress:   sender,
		ReceiverAddress: receiver,
		Amount:          amount,
		Note:            note,
		Purpose:         "personal",
		
		Status:          OrderStatusPending,
		StatusMessage:   "Transaction initiated",
		Language:        LanguageEnglish,
		
		CreatedAt:       now,
		ExpiresAt:       now.Add(24 * time.Hour),
		
		SMSNotification: true,
		QRCode:          GenerateQRCode(orderId),
		TrackingURL:     fmt.Sprintf("https://moneyorder.deshchain.com/track/%s", orderId),
	}
}

// GenerateReferenceNumber creates a traditional money order style reference
func GenerateReferenceNumber(timestamp time.Time) string {
	// Format: MO-YYYY-STATE-NNNNNN
	return fmt.Sprintf("MO-%d-%s-%06d", 
		timestamp.Year(),
		"DL", // Default to Delhi, will be dynamic
		timestamp.UnixNano() % 1000000,
	)
}

// GenerateQRCode creates a QR code for the order
func GenerateQRCode(orderId string) string {
	// In real implementation, this would generate actual QR code
	return fmt.Sprintf("QR:%s", orderId)
}

// ToSimplifiedStatus converts receipt to UPI-style simple status
func (r *MoneyOrderReceipt) ToSimplifiedStatus() SimplifiedOrderStatus {
	status := SimplifiedOrderStatus{
		Amount:       FormatAmountSimple(r.Amount),
		ReceiverName: r.ReceiverName,
		Time:         FormatTimeAgo(r.CreatedAt),
	}
	
	switch r.Status {
	case OrderStatusCompleted:
		status.StatusIcon = "✓"
		status.StatusText = "Success"
		status.StatusColor = "green"
		status.Actions = []OrderAction{
			{Label: "Repeat", Icon: "repeat", ActionType: "repeat", Enabled: true},
			{Label: "Share", Icon: "share", ActionType: "share", Enabled: true},
		}
		
	case OrderStatusPending, OrderStatusProcessing:
		status.StatusIcon = "⏳"
		status.StatusText = "Processing"
		status.StatusColor = "yellow"
		status.Actions = []OrderAction{
			{Label: "Track", Icon: "track", ActionType: "track", Enabled: true},
		}
		
	case OrderStatusCancelled, OrderStatusRefunded:
		status.StatusIcon = "❌"
		status.StatusText = "Failed"
		status.StatusColor = "red"
		status.Actions = []OrderAction{
			{Label: "Retry", Icon: "retry", ActionType: "retry", Enabled: true},
			{Label: "Support", Icon: "help", ActionType: "support", Enabled: true},
		}
	}
	
	return status
}

// FormatAmountSimple formats amount in simple rupee format
func FormatAmountSimple(amount sdk.Coin) string {
	// Convert to rupee format with commas
	// This is a simplified version - real implementation would be more robust
	amountInRupees := amount.Amount.Quo(sdk.NewInt(1_000_000)).Int64()
	return fmt.Sprintf("₹%d", amountInRupees)
}

// FormatTimeAgo formats time in "X mins ago" style
func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)
	
	if duration < time.Minute {
		return "Just now"
	} else if duration < time.Hour {
		mins := int(duration.Minutes())
		return fmt.Sprintf("%d mins ago", mins)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
}

// ValidateReceipt ensures all required fields are present
func (r *MoneyOrderReceipt) ValidateReceipt() error {
	if r.OrderId == "" {
		return fmt.Errorf("order ID cannot be empty")
	}
	
	if r.SenderAddress.Empty() {
		return fmt.Errorf("sender address cannot be empty")
	}
	
	if r.ReceiverAddress.Empty() {
		return fmt.Errorf("receiver address cannot be empty")
	}
	
	if !r.Amount.IsValid() || r.Amount.IsZero() {
		return fmt.Errorf("invalid amount")
	}
	
	// Amount limits check
	if r.Amount.Amount.LT(MinMoneyOrderAmount) {
		return fmt.Errorf("amount below minimum: %s", MinMoneyOrderAmount)
	}
	
	if r.Amount.Amount.GT(MaxMoneyOrderAmount) {
		return fmt.Errorf("amount exceeds maximum: %s", MaxMoneyOrderAmount)
	}
	
	return nil
}