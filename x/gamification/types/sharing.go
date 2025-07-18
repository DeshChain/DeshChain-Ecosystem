package types

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// SharingManager handles social media sharing
type SharingManager struct {
	baseURL string
}

// NewSharingManager creates a new sharing manager
func NewSharingManager(baseURL string) *SharingManager {
	return &SharingManager{
		baseURL: baseURL,
	}
}

// GenerateShareLinks creates platform-specific share links
func (sm *SharingManager) GenerateShareLinks(
	post *SocialMediaPost,
	achievement *Achievement,
	profile *DeveloperProfile,
) map[string]string {
	
	links := make(map[string]string)
	
	// Generate unique share URL
	shareURL := fmt.Sprintf("%s/achievement/%d?user=%s", 
		sm.baseURL, 
		achievement.AchievementId,
		profile.GithubUsername,
	)
	
	// Twitter/X
	twitterText := url.QueryEscape(fmt.Sprintf(
		"🎊 Just unlocked '%s' on @DeshChain! %s",
		achievement.Name,
		achievement.UnlockQuote,
	))
	twitterHashtags := url.QueryEscape(strings.Join(post.Hashtags, ","))
	links["twitter"] = fmt.Sprintf(
		"https://twitter.com/intent/tweet?text=%s&hashtags=%s&url=%s",
		twitterText, twitterHashtags, url.QueryEscape(shareURL),
	)
	
	// LinkedIn
	linkedinTitle := url.QueryEscape(fmt.Sprintf("Achievement Unlocked: %s", achievement.Name))
	linkedinSummary := url.QueryEscape(post.Content)
	links["linkedin"] = fmt.Sprintf(
		"https://www.linkedin.com/sharing/share-offsite/?url=%s&title=%s&summary=%s",
		url.QueryEscape(shareURL), linkedinTitle, linkedinSummary,
	)
	
	// WhatsApp
	whatsappText := url.QueryEscape(fmt.Sprintf(
		"*🎊 Achievement Unlocked!*\n\n*%s*\n%s\n\n_%s_\n\nCheck it out: %s",
		achievement.Name,
		achievement.Description,
		achievement.UnlockQuote,
		shareURL,
	))
	links["whatsapp"] = fmt.Sprintf("https://wa.me/?text=%s", whatsappText)
	
	// Telegram
	telegramText := url.QueryEscape(post.Content)
	links["telegram"] = fmt.Sprintf(
		"https://t.me/share/url?url=%s&text=%s",
		url.QueryEscape(shareURL), telegramText,
	)
	
	// Facebook
	links["facebook"] = fmt.Sprintf(
		"https://www.facebook.com/sharer/sharer.php?u=%s",
		url.QueryEscape(shareURL),
	)
	
	// Reddit
	redditTitle := url.QueryEscape(fmt.Sprintf("[DeshChain] Unlocked: %s", achievement.Name))
	links["reddit"] = fmt.Sprintf(
		"https://reddit.com/submit?url=%s&title=%s",
		url.QueryEscape(shareURL), redditTitle,
	)
	
	// Email
	emailSubject := url.QueryEscape(fmt.Sprintf("Check out my DeshChain achievement: %s", achievement.Name))
	emailBody := url.QueryEscape(fmt.Sprintf(
		"Hi!\n\nI just unlocked an amazing achievement on DeshChain:\n\n%s\n%s\n\n%s\n\nCheck it out here: %s\n\nJoin me on India's first cultural blockchain!",
		achievement.Name,
		achievement.Description,
		achievement.UnlockQuote,
		shareURL,
	))
	links["email"] = fmt.Sprintf("mailto:?subject=%s&body=%s", emailSubject, emailBody)
	
	return links
}

// GenerateShareableImage creates shareable image data
func (sm *SharingManager) GenerateShareableImage(
	card *AchievementCard,
	posterData map[string]interface{},
) ShareableImage {
	
	return ShareableImage{
		Type:         "achievement_poster",
		Format:       "PNG",
		Width:        800,
		Height:       1000,
		Quality:      90,
		PosterData:   posterData,
		Watermark:    true,
		WatermarkText: "DeshChain.io",
		Metadata: ImageMetadata{
			Title:       card.MovieTitle,
			Description: card.TagLine,
			Author:      card.StarringText,
			Created:     card.ReleaseDate.Format(time.RFC3339),
		},
	}
}

// ShareableImage represents an image ready for sharing
type ShareableImage struct {
	Type          string
	Format        string
	Width         int
	Height        int
	Quality       int
	PosterData    map[string]interface{}
	Watermark     bool
	WatermarkText string
	Metadata      ImageMetadata
}

// ImageMetadata contains image metadata
type ImageMetadata struct {
	Title       string
	Description string
	Author      string
	Created     string
}

// GenerateOGTags creates Open Graph tags for sharing
func GenerateOGTags(achievement *Achievement, profile *DeveloperProfile, imageURL string) map[string]string {
	return map[string]string{
		"og:title":       fmt.Sprintf("%s - DeshChain Achievement", achievement.Name),
		"og:description": achievement.Description,
		"og:image":       imageURL,
		"og:url":         fmt.Sprintf("https://deshchain.io/achievement/%d", achievement.AchievementId),
		"og:type":        "article",
		"og:site_name":   "DeshChain",
		"twitter:card":   "summary_large_image",
		"twitter:site":   "@DeshChain",
		"twitter:creator": fmt.Sprintf("@%s", profile.GithubUsername),
		"twitter:title":  achievement.Name,
		"twitter:description": achievement.UnlockQuote,
		"twitter:image":  imageURL,
	}
}

// GenerateShareCode creates a unique share code
func GenerateShareCode(achievementId uint64, userId string) string {
	data := fmt.Sprintf("%d:%s:%d", achievementId, userId, time.Now().Unix())
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	// Make it URL safe and shorter
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.TrimRight(encoded, "=")
	
	if len(encoded) > 12 {
		encoded = encoded[:12]
	}
	
	return encoded
}

// DecodeShareCode decodes a share code
func DecodeShareCode(code string) (achievementId uint64, userId string, timestamp int64, err error) {
	// Add padding if needed
	padding := (4 - len(code)%4) % 4
	code += strings.Repeat("=", padding)
	
	// Restore URL encoding
	code = strings.ReplaceAll(code, "-", "+")
	code = strings.ReplaceAll(code, "_", "/")
	
	decoded, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return 0, "", 0, err
	}
	
	parts := strings.Split(string(decoded), ":")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid share code format")
	}
	
	fmt.Sscanf(parts[0], "%d", &achievementId)
	userId = parts[1]
	fmt.Sscanf(parts[2], "%d", &timestamp)
	
	return
}

// GenerateQRCode generates QR code data for sharing
func GenerateQRCode(shareURL string) map[string]interface{} {
	return map[string]interface{}{
		"url":        shareURL,
		"size":       256,
		"errorLevel": "M", // Medium error correction
		"margin":     4,
		"darkColor":  "#000000",
		"lightColor": "#FFFFFF",
		"logo":       "/logos/deshchain-small.png",
		"logoSize":   48,
	}
}

// TrackShareEvent tracks when content is shared
func TrackShareEvent(post *SocialMediaPost, platform string) ShareEvent {
	return ShareEvent{
		PostId:    post.PostId,
		Platform:  platform,
		Timestamp: time.Now(),
		UserAgent: "DeshChain/1.0",
		Success:   true,
	}
}

// ShareEvent represents a sharing event
type ShareEvent struct {
	PostId    uint64
	Platform  string
	Timestamp time.Time
	UserAgent string
	Success   bool
	Error     string
}

// GenerateShareStats creates sharing statistics
func GenerateShareStats(events []ShareEvent) map[string]interface{} {
	stats := make(map[string]interface{})
	platformCounts := make(map[string]int)
	successCount := 0
	
	for _, event := range events {
		platformCounts[event.Platform]++
		if event.Success {
			successCount++
		}
	}
	
	stats["total_shares"] = len(events)
	stats["successful_shares"] = successCount
	stats["platform_breakdown"] = platformCounts
	stats["success_rate"] = float64(successCount) / float64(len(events)) * 100
	
	// Find most popular platform
	maxShares := 0
	popularPlatform := ""
	for platform, count := range platformCounts {
		if count > maxShares {
			maxShares = count
			popularPlatform = platform
		}
	}
	stats["most_popular_platform"] = popularPlatform
	
	return stats
}

// CreateShareBundle creates a complete sharing package
func CreateShareBundle(
	achievement *Achievement,
	profile *DeveloperProfile,
	card *AchievementCard,
	posterData map[string]interface{},
) ShareBundle {
	
	sm := NewSharingManager("https://deshchain.io")
	
	// Generate all share links
	post := &SocialMediaPost{
		Content:  fmt.Sprintf("Unlocked: %s", achievement.Name),
		Hashtags: GetHashtags(achievement.Category.String()),
	}
	
	shareLinks := sm.GenerateShareLinks(post, achievement, profile)
	shareCode := GenerateShareCode(achievement.AchievementId, profile.Address)
	shareURL := fmt.Sprintf("https://deshchain.io/share/%s", shareCode)
	
	return ShareBundle{
		ShareCode:   shareCode,
		ShareURL:    shareURL,
		ShareLinks:  shareLinks,
		QRCode:      GenerateQRCode(shareURL),
		OGTags:      GenerateOGTags(achievement, profile, card.PosterImageUrl),
		Image:       sm.GenerateShareableImage(card, posterData),
		Achievement: achievement,
		Profile:     profile,
		Card:        card,
	}
}

// ShareBundle contains all sharing data
type ShareBundle struct {
	ShareCode   string
	ShareURL    string
	ShareLinks  map[string]string
	QRCode      map[string]interface{}
	OGTags      map[string]string
	Image       ShareableImage
	Achievement *Achievement
	Profile     *DeveloperProfile
	Card        *AchievementCard
}