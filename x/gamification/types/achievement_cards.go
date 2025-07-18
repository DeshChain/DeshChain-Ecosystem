package types

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// MoviePosterGenerator generates Bollywood movie poster style achievement cards
type MoviePosterGenerator struct {
	templates   []PosterTemplate
	colorThemes []ColorTheme
}

// PosterTemplate defines a movie poster template
type PosterTemplate struct {
	Name        string
	Layout      string
	Elements    []PosterElement
	AspectRatio string
}

// PosterElement defines an element on the poster
type PosterElement struct {
	Type     string // text, image, shape
	Content  string
	Position Position
	Style    ElementStyle
}

// Position defines x,y coordinates
type Position struct {
	X      int
	Y      int
	Width  int
	Height int
}

// ElementStyle defines styling
type ElementStyle struct {
	FontFamily   string
	FontSize     int
	Color        string
	Background   string
	Border       string
	Shadow       string
	Alignment    string
	Bold         bool
	Italic       bool
	Decorative   bool
}

// ColorTheme defines color schemes
type ColorTheme struct {
	Name       string
	Primary    string
	Secondary  string
	Accent     string
	Background string
	Text       string
	Gold       string
}

// NewMoviePosterGenerator creates a new generator
func NewMoviePosterGenerator() *MoviePosterGenerator {
	return &MoviePosterGenerator{
		templates:   getDefaultTemplates(),
		colorThemes: getDefaultColorThemes(),
	}
}

// GenerateAchievementPoster creates a movie poster for achievement
func (mpg *MoviePosterGenerator) GenerateAchievementPoster(
	card *AchievementCard,
	achievement *Achievement,
	profile *DeveloperProfile,
) map[string]interface{} {
	
	// Select random template and theme
	template := mpg.templates[rand.Intn(len(mpg.templates))]
	theme := mpg.colorThemes[rand.Intn(len(mpg.colorThemes))]
	
	// Create poster data
	posterData := map[string]interface{}{
		"template":     template.Name,
		"theme":        theme.Name,
		"aspectRatio":  template.AspectRatio,
		"elements":     []map[string]interface{}{},
		"metadata":     generatePosterMetadata(card, achievement, profile),
	}
	
	// Generate poster elements
	elements := mpg.generatePosterElements(card, achievement, profile, template, theme)
	posterData["elements"] = elements
	
	// Add special effects based on rarity
	posterData["effects"] = mpg.getSpecialEffects(achievement.RarityLevel)
	
	// Add background pattern
	posterData["background"] = mpg.generateBackground(achievement.Category, theme)
	
	return posterData
}

// generatePosterElements creates all elements for the poster
func (mpg *MoviePosterGenerator) generatePosterElements(
	card *AchievementCard,
	achievement *Achievement,
	profile *DeveloperProfile,
	template PosterTemplate,
	theme ColorTheme,
) []map[string]interface{} {
	
	elements := []map[string]interface{}{}
	
	// Main Title (Movie name style)
	elements = append(elements, map[string]interface{}{
		"type": "text",
		"content": card.MovieTitle,
		"position": map[string]int{
			"x": 50, "y": 30, "width": 700, "height": 100,
		},
		"style": map[string]interface{}{
			"fontFamily": "Bebas Neue",
			"fontSize": 72,
			"color": theme.Gold,
			"textShadow": "4px 4px 8px rgba(0,0,0,0.8)",
			"textAlign": "center",
			"fontWeight": "bold",
			"letterSpacing": "4px",
			"textTransform": "uppercase",
		},
	})
	
	// Tagline
	elements = append(elements, map[string]interface{}{
		"type": "text",
		"content": card.TagLine,
		"position": map[string]int{
			"x": 50, "y": 130, "width": 700, "height": 40,
		},
		"style": map[string]interface{}{
			"fontFamily": "Montserrat",
			"fontSize": 24,
			"color": theme.Secondary,
			"textAlign": "center",
			"fontStyle": "italic",
		},
	})
	
	// Avatar/Character Image
	avatar := GetAvatarByType(profile.ActiveAvatar)
	elements = append(elements, map[string]interface{}{
		"type": "image",
		"content": avatar.IconUrl,
		"position": map[string]int{
			"x": 250, "y": 200, "width": 300, "height": 400,
		},
		"style": map[string]interface{}{
			"filter": "drop-shadow(0 0 20px " + theme.Accent + ")",
			"borderRadius": "10px",
		},
	})
	
	// Character Name Plate
	elements = append(elements, map[string]interface{}{
		"type": "shape",
		"content": "rectangle",
		"position": map[string]int{
			"x": 200, "y": 550, "width": 400, "height": 60,
		},
		"style": map[string]interface{}{
			"background": "linear-gradient(90deg, " + theme.Primary + ", " + theme.Secondary + ")",
			"borderRadius": "30px",
			"opacity": 0.9,
		},
	})
	
	// Starring Text
	elements = append(elements, map[string]interface{}{
		"type": "text",
		"content": card.StarringText,
		"position": map[string]int{
			"x": 200, "y": 565, "width": 400, "height": 30,
		},
		"style": map[string]interface{}{
			"fontFamily": "Rajdhani",
			"fontSize": 22,
			"color": "#FFFFFF",
			"textAlign": "center",
			"fontWeight": "600",
		},
	})
	
	// Achievement Badge
	elements = append(elements, map[string]interface{}{
		"type": "badge",
		"content": getRarityBadge(achievement.RarityLevel),
		"position": map[string]int{
			"x": 600, "y": 220, "width": 120, "height": 120,
		},
		"style": map[string]interface{}{
			"animation": "pulse 2s infinite",
		},
	})
	
	// Box Office Collection
	elements = append(elements, map[string]interface{}{
		"type": "text",
		"content": card.BoxOfficeText,
		"position": map[string]int{
			"x": 50, "y": 650, "width": 700, "height": 40,
		},
		"style": map[string]interface{}{
			"fontFamily": "Oswald",
			"fontSize": 32,
			"color": theme.Gold,
			"textAlign": "center",
			"fontWeight": "bold",
			"textShadow": "2px 2px 4px rgba(0,0,0,0.6)",
		},
	})
	
	// Critic Reviews Section
	reviewY := 720
	for i, review := range card.CriticReviews[:min(3, len(card.CriticReviews))] {
		elements = append(elements, map[string]interface{}{
			"type": "text",
			"content": review,
			"position": map[string]int{
				"x": 50, "y": reviewY + (i * 30), "width": 700, "height": 25,
			},
			"style": map[string]interface{}{
				"fontFamily": "Roboto",
				"fontSize": 16,
				"color": theme.Text,
				"textAlign": "center",
			},
		})
	}
	
	// Director/Producer Credits
	elements = append(elements, map[string]interface{}{
		"type": "text",
		"content": card.DirectorText,
		"position": map[string]int{
			"x": 50, "y": 850, "width": 700, "height": 30,
		},
		"style": map[string]interface{}{
			"fontFamily": "Open Sans",
			"fontSize": 14,
			"color": theme.Secondary,
			"textAlign": "center",
		},
	})
	
	// Release Date
	elements = append(elements, map[string]interface{}{
		"type": "text",
		"content": fmt.Sprintf("Released: %s", card.ReleaseDate.Format("2 January 2006")),
		"position": map[string]int{
			"x": 50, "y": 880, "width": 700, "height": 25,
		},
		"style": map[string]interface{}{
			"fontFamily": "Roboto",
			"fontSize": 14,
			"color": theme.Text,
			"textAlign": "center",
		},
	})
	
	// DeshChain Logo
	elements = append(elements, map[string]interface{}{
		"type": "logo",
		"content": "/logos/deshchain.png",
		"position": map[string]int{
			"x": 50, "y": 920, "width": 150, "height": 50,
		},
	})
	
	// Social Share Icons
	elements = append(elements, map[string]interface{}{
		"type": "social_icons",
		"content": "share",
		"position": map[string]int{
			"x": 600, "y": 920, "width": 150, "height": 50,
		},
	})
	
	// Add decorative elements based on achievement type
	decorativeElements := mpg.getDecorativeElements(achievement.Category, theme)
	elements = append(elements, decorativeElements...)
	
	return elements
}

// getSpecialEffects returns effects based on rarity
func (mpg *MoviePosterGenerator) getSpecialEffects(rarity RarityLevel) map[string]interface{} {
	effects := map[string]interface{}{
		"particles": false,
		"glow": false,
		"animation": "none",
		"filter": "none",
	}
	
	switch rarity {
	case RarityLevel_RARITY_LEVEL_RARE:
		effects["glow"] = true
		effects["filter"] = "brightness(1.1)"
	case RarityLevel_RARITY_LEVEL_EPIC:
		effects["particles"] = true
		effects["glow"] = true
		effects["animation"] = "shimmer"
		effects["filter"] = "brightness(1.2) contrast(1.1)"
	case RarityLevel_RARITY_LEVEL_LEGENDARY:
		effects["particles"] = true
		effects["glow"] = true
		effects["animation"] = "rainbow"
		effects["filter"] = "brightness(1.3) contrast(1.2) saturate(1.2)"
	case RarityLevel_RARITY_LEVEL_MYTHIC:
		effects["particles"] = true
		effects["glow"] = true
		effects["animation"] = "mythic_aura"
		effects["filter"] = "brightness(1.4) contrast(1.3) saturate(1.3)"
		effects["special"] = "lightning_border"
	}
	
	return effects
}

// generateBackground creates background pattern
func (mpg *MoviePosterGenerator) generateBackground(category AchievementCategory, theme ColorTheme) map[string]interface{} {
	backgrounds := map[AchievementCategory]string{
		AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS:      "circuit_pattern",
		AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES:    "matrix_rain",
		AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES:     "star_burst",
		AchievementCategory_ACHIEVEMENT_CATEGORY_DOCUMENTATION: "book_pattern",
		AchievementCategory_ACHIEVEMENT_CATEGORY_PERFORMANCE:  "speed_lines",
		AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK:       "fire_pattern",
		AchievementCategory_ACHIEVEMENT_CATEGORY_SOCIAL:       "network_nodes",
		AchievementCategory_ACHIEVEMENT_CATEGORY_SPECIAL:      "mandala_pattern",
	}
	
	pattern := backgrounds[category]
	if pattern == "" {
		pattern = "gradient_mesh"
	}
	
	return map[string]interface{}{
		"pattern": pattern,
		"primaryColor": theme.Background,
		"secondaryColor": theme.Primary,
		"opacity": 0.1,
		"blend": "multiply",
	}
}

// getDecorativeElements adds category-specific decorations
func (mpg *MoviePosterGenerator) getDecorativeElements(category AchievementCategory, theme ColorTheme) []map[string]interface{} {
	elements := []map[string]interface{}{}
	
	switch category {
	case AchievementCategory_ACHIEVEMENT_CATEGORY_COMMITS:
		// Add commit graph decoration
		elements = append(elements, map[string]interface{}{
			"type": "decoration",
			"content": "commit_graph",
			"position": map[string]int{
				"x": 50, "y": 400, "width": 100, "height": 100,
			},
			"style": map[string]interface{}{
				"opacity": 0.3,
				"color": theme.Accent,
			},
		})
		
	case AchievementCategory_ACHIEVEMENT_CATEGORY_BUG_FIXES:
		// Add bug icons
		elements = append(elements, map[string]interface{}{
			"type": "decoration",
			"content": "bug_squash",
			"position": map[string]int{
				"x": 650, "y": 400, "width": 80, "height": 80,
			},
			"style": map[string]interface{}{
				"opacity": 0.4,
				"color": theme.Accent,
				"animation": "squash 3s infinite",
			},
		})
		
	case AchievementCategory_ACHIEVEMENT_CATEGORY_FEATURES:
		// Add star decorations
		for i := 0; i < 5; i++ {
			elements = append(elements, map[string]interface{}{
				"type": "decoration",
				"content": "star",
				"position": map[string]int{
					"x": 100 + (i * 120), "y": 350, "width": 30, "height": 30,
				},
				"style": map[string]interface{}{
					"opacity": 0.6,
					"color": theme.Gold,
					"animation": fmt.Sprintf("twinkle %ds infinite", 2+i),
				},
			})
		}
		
	case AchievementCategory_ACHIEVEMENT_CATEGORY_STREAK:
		// Add fire decoration
		elements = append(elements, map[string]interface{}{
			"type": "decoration",
			"content": "fire",
			"position": map[string]int{
				"x": 100, "y": 300, "width": 60, "height": 80,
			},
			"style": map[string]interface{}{
				"opacity": 0.7,
				"animation": "flicker 1s infinite",
			},
		})
	}
	
	// Add film reel borders
	elements = append(elements, map[string]interface{}{
		"type": "decoration",
		"content": "film_reel_top",
		"position": map[string]int{
			"x": 0, "y": 0, "width": 800, "height": 20,
		},
		"style": map[string]interface{}{
			"opacity": 0.8,
			"color": "#000000",
		},
	})
	
	elements = append(elements, map[string]interface{}{
		"type": "decoration",
		"content": "film_reel_bottom",
		"position": map[string]int{
			"x": 0, "y": 980, "width": 800, "height": 20,
		},
		"style": map[string]interface{}{
			"opacity": 0.8,
			"color": "#000000",
		},
	})
	
	return elements
}

// getRarityBadge returns badge design for rarity
func getRarityBadge(rarity RarityLevel) string {
	badges := map[RarityLevel]string{
		RarityLevel_RARITY_LEVEL_COMMON:    "bronze_badge",
		RarityLevel_RARITY_LEVEL_RARE:      "silver_badge",
		RarityLevel_RARITY_LEVEL_EPIC:      "gold_badge",
		RarityLevel_RARITY_LEVEL_LEGENDARY: "diamond_badge",
		RarityLevel_RARITY_LEVEL_MYTHIC:    "mythic_badge",
	}
	
	return badges[rarity]
}

// generatePosterMetadata creates metadata for the poster
func generatePosterMetadata(card *AchievementCard, achievement *Achievement, profile *DeveloperProfile) map[string]interface{} {
	return map[string]interface{}{
		"title": card.MovieTitle,
		"achievementName": achievement.Name,
		"developerName": profile.GithubUsername,
		"achievementId": achievement.AchievementId,
		"rarity": achievement.RarityLevel.String(),
		"category": achievement.Category.String(),
		"timestamp": time.Now().Unix(),
		"shareUrl": fmt.Sprintf("https://deshchain.io/achievement/%d", card.CardId),
		"nftMintable": !card.IsNft,
	}
}

// getDefaultTemplates returns predefined poster templates
func getDefaultTemplates() []PosterTemplate {
	return []PosterTemplate{
		{
			Name:        "Classic Bollywood",
			Layout:      "vertical",
			AspectRatio: "2:3",
		},
		{
			Name:        "Modern Multiplex",
			Layout:      "vertical",
			AspectRatio: "2:3",
		},
		{
			Name:        "Retro 90s",
			Layout:      "vertical",
			AspectRatio: "2:3",
		},
		{
			Name:        "South Blockbuster",
			Layout:      "vertical", 
			AspectRatio: "2:3",
		},
		{
			Name:        "Festival Release",
			Layout:      "vertical",
			AspectRatio: "2:3",
		},
	}
}

// getDefaultColorThemes returns predefined color themes
func getDefaultColorThemes() []ColorTheme {
	return []ColorTheme{
		{
			Name:       "Bollywood Gold",
			Primary:    "#FFD700",
			Secondary:  "#FF6B6B",
			Accent:     "#4ECDC4",
			Background: "#1A1A2E",
			Text:       "#FFFFFF",
			Gold:       "#FFD700",
		},
		{
			Name:       "Desi Masala",
			Primary:    "#FF6B35",
			Secondary:  "#F7931E",
			Accent:     "#FFCC00",
			Background: "#2C1810",
			Text:       "#FFF8DC",
			Gold:       "#FFD700",
		},
		{
			Name:       "Royal Blue",
			Primary:    "#0F3460",
			Secondary:  "#16537E",
			Accent:     "#E94560",
			Background: "#1A1A2E",
			Text:       "#FFFFFF",
			Gold:       "#FFD700",
		},
		{
			Name:       "Festival Vibes",
			Primary:    "#FF006E",
			Secondary:  "#FB5607",
			Accent:     "#FFBE0B",
			Background: "#3A0CA3",
			Text:       "#FFFFFF",
			Gold:       "#FFD700",
		},
		{
			Name:       "Tech Thriller",
			Primary:    "#00F5FF",
			Secondary:  "#00CED1",
			Accent:     "#FF1493",
			Background: "#0A0E27",
			Text:       "#FFFFFF",
			Gold:       "#00F5FF",
		},
	}
}

// GenerateAchievementCardHTML generates HTML for web display
func GenerateAchievementCardHTML(posterData map[string]interface{}) string {
	// Convert poster data to JSON
	jsonData, _ := json.Marshal(posterData)
	
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>DeshChain Achievement Card</title>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Bebas+Neue&family=Montserrat:ital@1&family=Oswald:wght@700&family=Rajdhani:wght@600&family=Roboto&display=swap');
        
        body {
            margin: 0;
            padding: 20px;
            background: #0a0a0a;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        
        .poster-container {
            width: 800px;
            height: 1000px;
            position: relative;
            box-shadow: 0 0 50px rgba(255, 215, 0, 0.5);
            overflow: hidden;
        }
        
        .poster-element {
            position: absolute;
        }
        
        @keyframes pulse {
            0%% { transform: scale(1); }
            50%% { transform: scale(1.1); }
            100%% { transform: scale(1); }
        }
        
        @keyframes shimmer {
            0%% { filter: brightness(1); }
            50%% { filter: brightness(1.5); }
            100%% { filter: brightness(1); }
        }
        
        @keyframes twinkle {
            0%%, 100%% { opacity: 0.3; }
            50%% { opacity: 1; }
        }
        
        @keyframes flicker {
            0%%, 100%% { opacity: 0.7; }
            50%% { opacity: 1; }
        }
        
        .share-button {
            margin-top: 20px;
            padding: 15px 30px;
            background: linear-gradient(45deg, #FFD700, #FF6B6B);
            border: none;
            border-radius: 30px;
            color: white;
            font-size: 18px;
            font-weight: bold;
            cursor: pointer;
            transition: transform 0.3s;
        }
        
        .share-button:hover {
            transform: scale(1.05);
        }
    </style>
</head>
<body>
    <div id="poster-container" class="poster-container"></div>
    <script>
        const posterData = %s;
        // Render poster elements dynamically
        renderPoster(posterData);
        
        function renderPoster(data) {
            const container = document.getElementById('poster-container');
            // Implementation would render each element based on posterData
        }
    </script>
</body>
</html>`, strings.ReplaceAll(string(jsonData), "%", "%%"))
	
	return html
}

// min returns minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}