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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/deshchain/deshchain/x/moneyorder/types"
)

// NotificationService handles SMS and WhatsApp notifications
type NotificationService struct {
	keeper          Keeper
	smsProvider     SMSProvider
	whatsappProvider WhatsAppProvider
	templateEngine  *template.Template
}

// SMSProvider interface for SMS service integration
type SMSProvider interface {
	SendSMS(ctx context.Context, to string, message string, language string) error
	SendOTP(ctx context.Context, to string, otp string) error
	GetBalance(ctx context.Context) (float64, error)
}

// WhatsAppProvider interface for WhatsApp Business API
type WhatsAppProvider interface {
	SendMessage(ctx context.Context, to string, message WhatsAppMessage) error
	SendTemplate(ctx context.Context, to string, templateName string, params map[string]interface{}) error
	SendReceipt(ctx context.Context, to string, receipt types.MoneyOrderReceipt) error
}

// WhatsAppMessage represents a WhatsApp message
type WhatsAppMessage struct {
	Type       string                 `json:"type"`
	Text       string                 `json:"text,omitempty"`
	MediaURL   string                 `json:"media_url,omitempty"`
	Template   string                 `json:"template,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Language   string                 `json:"language"`
}

// TwilioSMSProvider implements SMS using Twilio
type TwilioSMSProvider struct {
	AccountSID string
	AuthToken  string
	FromNumber string
	BaseURL    string
}

// SendSMS sends an SMS using Twilio
func (t *TwilioSMSProvider) SendSMS(ctx context.Context, to string, message string, language string) error {
	// Translate message if needed
	translatedMessage := translateMessage(message, language)
	
	// Prepare request
	data := map[string]string{
		"To":   to,
		"From": t.FromNumber,
		"Body": translatedMessage,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create request
	url := fmt.Sprintf("%s/Accounts/%s/Messages.json", t.BaseURL, t.AccountSID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.SetBasicAuth(t.AccountSID, t.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send SMS: status %d", resp.StatusCode)
	}

	return nil
}

// SendOTP sends an OTP SMS
func (t *TwilioSMSProvider) SendOTP(ctx context.Context, to string, otp string) error {
	message := fmt.Sprintf("Your DeshChain verification code is: %s. Valid for 10 minutes.", otp)
	return t.SendSMS(ctx, to, message, "en")
}

// GetBalance gets SMS credits balance
func (t *TwilioSMSProvider) GetBalance(ctx context.Context) (float64, error) {
	// Implementation for getting Twilio balance
	return 100.0, nil // Mock balance
}

// WhatsAppBusinessProvider implements WhatsApp using Business API
type WhatsAppBusinessProvider struct {
	APIKey      string
	PhoneNumber string
	BaseURL     string
	Templates   map[string]string
}

// SendMessage sends a WhatsApp message
func (w *WhatsAppBusinessProvider) SendMessage(ctx context.Context, to string, message WhatsAppMessage) error {
	// Prepare WhatsApp API request
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              message.Type,
	}

	switch message.Type {
	case "text":
		data["text"] = map[string]string{"body": message.Text}
	case "template":
		data["template"] = map[string]interface{}{
			"name":     message.Template,
			"language": map[string]string{"code": message.Language},
			"components": []map[string]interface{}{
				{
					"type": "body",
					"parameters": formatTemplateParams(message.Parameters),
				},
			},
		}
	case "image":
		data["image"] = map[string]string{
			"link":    message.MediaURL,
			"caption": message.Text,
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Send request
	req, err := http.NewRequestWithContext(ctx, "POST", w.BaseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+w.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send WhatsApp message: status %d", resp.StatusCode)
	}

	return nil
}

// SendTemplate sends a WhatsApp template message
func (w *WhatsAppBusinessProvider) SendTemplate(ctx context.Context, to string, templateName string, params map[string]interface{}) error {
	message := WhatsAppMessage{
		Type:       "template",
		Template:   templateName,
		Parameters: params,
		Language:   "en", // Default to English
	}
	return w.SendMessage(ctx, to, message)
}

// SendReceipt sends a money order receipt via WhatsApp
func (w *WhatsAppBusinessProvider) SendReceipt(ctx context.Context, to string, receipt types.MoneyOrderReceipt) error {
	// Generate QR code for receipt
	qrCodeURL := generateReceiptQRCode(receipt)

	// Send receipt template
	params := map[string]interface{}{
		"receipt_id":   receipt.ReceiptId,
		"amount":       formatAmount(receipt.Amount),
		"sender":       formatAddress(receipt.Sender),
		"receiver":     formatAddress(receipt.Receiver),
		"status":       receipt.Status,
		"timestamp":    receipt.Timestamp.Format("02 Jan 2006, 03:04 PM"),
		"verification": receipt.VerificationCode,
		"qr_code_url":  qrCodeURL,
	}

	return w.SendTemplate(ctx, to, "money_order_receipt", params)
}

// NewNotificationService creates a new notification service
func NewNotificationService(k Keeper) *NotificationService {
	// Initialize SMS provider (Twilio)
	smsProvider := &TwilioSMSProvider{
		AccountSID: getEnvVar("TWILIO_ACCOUNT_SID"),
		AuthToken:  getEnvVar("TWILIO_AUTH_TOKEN"),
		FromNumber: getEnvVar("TWILIO_FROM_NUMBER"),
		BaseURL:    "https://api.twilio.com/2010-04-01",
	}

	// Initialize WhatsApp provider
	whatsappProvider := &WhatsAppBusinessProvider{
		APIKey:      getEnvVar("WHATSAPP_API_KEY"),
		PhoneNumber: getEnvVar("WHATSAPP_PHONE_NUMBER"),
		BaseURL:     "https://graph.facebook.com/v17.0",
		Templates:   loadWhatsAppTemplates(),
	}

	// Load notification templates
	tmpl := loadNotificationTemplates()

	return &NotificationService{
		keeper:           k,
		smsProvider:      smsProvider,
		whatsappProvider: whatsappProvider,
		templateEngine:   tmpl,
	}
}

// SendMoneyOrderNotification sends notification for money order
func (ns *NotificationService) SendMoneyOrderNotification(
	ctx sdk.Context,
	order types.MoneyOrder,
	receipt types.MoneyOrderReceipt,
	notificationType string,
) error {
	// Get user preferences
	senderPrefs := ns.getUserPreferences(ctx, order.Sender)
	receiverPrefs := ns.getUserPreferences(ctx, order.Receiver)

	// Send notifications based on preferences
	switch notificationType {
	case "created":
		// Notify sender
		if senderPrefs.EnableSMS {
			ns.sendSMSNotification(ctx, senderPrefs.PhoneNumber, "money_order_sent", receipt, senderPrefs.Language)
		}
		if senderPrefs.EnableWhatsApp {
			ns.sendWhatsAppNotification(ctx, senderPrefs.PhoneNumber, receipt)
		}

	case "completed":
		// Notify both sender and receiver
		if receiverPrefs.EnableSMS {
			ns.sendSMSNotification(ctx, receiverPrefs.PhoneNumber, "money_order_received", receipt, receiverPrefs.Language)
		}
		if receiverPrefs.EnableWhatsApp {
			ns.sendWhatsAppNotification(ctx, receiverPrefs.PhoneNumber, receipt)
		}

	case "failed":
		// Notify sender about failure
		if senderPrefs.EnableSMS {
			ns.sendSMSNotification(ctx, senderPrefs.PhoneNumber, "money_order_failed", receipt, senderPrefs.Language)
		}
	}

	return nil
}

// sendSMSNotification sends SMS notification
func (ns *NotificationService) sendSMSNotification(
	ctx sdk.Context,
	phoneNumber string,
	templateName string,
	receipt types.MoneyOrderReceipt,
	language string,
) error {
	// Generate message from template
	message, err := ns.generateMessage(templateName, receipt, language)
	if err != nil {
		return err
	}

	// Send SMS
	return ns.smsProvider.SendSMS(context.Background(), phoneNumber, message, language)
}

// sendWhatsAppNotification sends WhatsApp notification
func (ns *NotificationService) sendWhatsAppNotification(
	ctx sdk.Context,
	phoneNumber string,
	receipt types.MoneyOrderReceipt,
) error {
	return ns.whatsappProvider.SendReceipt(context.Background(), phoneNumber, receipt)
}

// generateMessage generates message from template
func (ns *NotificationService) generateMessage(
	templateName string,
	receipt types.MoneyOrderReceipt,
	language string,
) (string, error) {
	// Template data
	data := map[string]interface{}{
		"ReceiptID":   receipt.ReceiptId,
		"Amount":      formatAmount(receipt.Amount),
		"Currency":    receipt.Amount.Denom,
		"Sender":      formatAddress(receipt.Sender),
		"Receiver":    formatAddress(receipt.Receiver),
		"Status":      receipt.Status,
		"Timestamp":   receipt.Timestamp.Format("02 Jan 2006, 03:04 PM"),
		"Verification": receipt.VerificationCode,
	}

	// Add cultural elements
	if receipt.CulturalQuote != "" {
		data["Quote"] = receipt.CulturalQuote
	}
	if receipt.PatriotismBonus != nil && receipt.PatriotismBonus.IsPositive() {
		data["Bonus"] = formatAmount(*receipt.PatriotismBonus)
	}

	// Execute template
	var buf bytes.Buffer
	tmplName := fmt.Sprintf("%s_%s", templateName, language)
	if err := ns.templateEngine.ExecuteTemplate(&buf, tmplName, data); err != nil {
		// Fallback to English
		tmplName = fmt.Sprintf("%s_en", templateName)
		if err := ns.templateEngine.ExecuteTemplate(&buf, tmplName, data); err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

// getUserPreferences gets user notification preferences
func (ns *NotificationService) getUserPreferences(ctx sdk.Context, address sdk.AccAddress) UserPreferences {
	store := ctx.KVStore(ns.keeper.storeKey)
	key := types.GetUserPreferencesKey(address)

	var prefs UserPreferences
	bz := store.Get(key)
	if bz == nil {
		// Return default preferences
		return UserPreferences{
			EnableSMS:      false,
			EnableWhatsApp: false,
			Language:       "en",
		}
	}

	ns.keeper.cdc.MustUnmarshal(bz, &prefs)
	return prefs
}

// UserPreferences represents user notification preferences
type UserPreferences struct {
	PhoneNumber    string `json:"phone_number"`
	EnableSMS      bool   `json:"enable_sms"`
	EnableWhatsApp bool   `json:"enable_whatsapp"`
	Language       string `json:"language"`
}

// Helper functions

func translateMessage(message, language string) string {
	// Simple translation logic - in production, use proper i18n
	translations := map[string]map[string]string{
		"hi": {
			"Your DeshChain verification code is": "आपका DeshChain सत्यापन कोड है",
			"Valid for 10 minutes":                "10 मिनट के लिए वैध",
		},
		"bn": {
			"Your DeshChain verification code is": "আপনার DeshChain যাচাইকরণ কোড হল",
			"Valid for 10 minutes":                "10 মিনিটের জন্য বৈধ",
		},
	}

	if trans, ok := translations[language]; ok {
		for en, local := range trans {
			message = bytes.ReplaceAll([]byte(message), []byte(en), []byte(local))
		}
	}

	return message
}

func formatTemplateParams(params map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, v := range params {
		result = append(result, map[string]interface{}{
			"type": "text",
			"text": fmt.Sprintf("%v", v),
		})
	}
	return result
}

func formatAmount(coin sdk.Coin) string {
	// Convert to readable format
	amount := coin.Amount.ToDec().QuoInt64(1000000)
	return fmt.Sprintf("%.2f %s", amount.MustFloat64(), coin.Denom)
}

func formatAddress(addr sdk.AccAddress) string {
	address := addr.String()
	if len(address) > 20 {
		return address[:8] + "..." + address[len(address)-8:]
	}
	return address
}

func generateReceiptQRCode(receipt types.MoneyOrderReceipt) string {
	// Generate QR code URL - in production, actually generate QR code
	return fmt.Sprintf("https://api.deshchain.org/qr/receipt/%s", receipt.ReceiptId)
}

func getEnvVar(key string) string {
	// In production, properly handle environment variables
	defaults := map[string]string{
		"TWILIO_ACCOUNT_SID":    "AC_test",
		"TWILIO_AUTH_TOKEN":     "auth_test",
		"TWILIO_FROM_NUMBER":    "+1234567890",
		"WHATSAPP_API_KEY":      "whatsapp_test",
		"WHATSAPP_PHONE_NUMBER": "+1234567890",
	}
	if val, ok := defaults[key]; ok {
		return val
	}
	return ""
}

func loadWhatsAppTemplates() map[string]string {
	return map[string]string{
		"money_order_receipt": "money_order_receipt_template",
		"money_order_status":  "money_order_status_template",
		"otp_verification":    "otp_verification_template",
	}
}

func loadNotificationTemplates() *template.Template {
	// Load all notification templates
	tmpl := template.New("notifications")

	// English templates
	tmpl.New("money_order_sent_en").Parse(`DeshChain: Money Order #{{.ReceiptID}} sent! Amount: {{.Amount}} to {{.Receiver}}. Track: deshchain.org/track/{{.ReceiptID}}`)
	tmpl.New("money_order_received_en").Parse(`DeshChain: Received {{.Amount}} from {{.Sender}}. Receipt #{{.ReceiptID}}. Verify: {{.Verification}}`)
	tmpl.New("money_order_failed_en").Parse(`DeshChain: Money Order #{{.ReceiptID}} failed. Amount {{.Amount}} will be refunded.`)

	// Hindi templates
	tmpl.New("money_order_sent_hi").Parse(`DeshChain: मनी ऑर्डर #{{.ReceiptID}} भेजा गया! राशि: {{.Amount}} को {{.Receiver}}। ट्रैक करें: deshchain.org/track/{{.ReceiptID}}`)
	tmpl.New("money_order_received_hi").Parse(`DeshChain: {{.Sender}} से {{.Amount}} प्राप्त हुआ। रसीद #{{.ReceiptID}}। सत्यापित करें: {{.Verification}}`)

	// Bengali templates
	tmpl.New("money_order_sent_bn").Parse(`DeshChain: মানি অর্ডার #{{.ReceiptID}} পাঠানো হয়েছে! পরিমাণ: {{.Amount}} প্রাপক {{.Receiver}}। ট্র্যাক করুন: deshchain.org/track/{{.ReceiptID}}`)

	return tmpl
}