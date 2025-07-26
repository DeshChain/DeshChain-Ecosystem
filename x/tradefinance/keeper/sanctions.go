package keeper

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/DeshChain/DeshChain-Ecosystem/x/tradefinance/types"
)

// SanctionsScreeningEngine implements comprehensive sanctions screening
type SanctionsScreeningEngine struct {
	keeper *Keeper
}

// NewSanctionsScreeningEngine creates a new sanctions screening engine
func NewSanctionsScreeningEngine(k *Keeper) *SanctionsScreeningEngine {
	return &SanctionsScreeningEngine{
		keeper: k,
	}
}

// ScreeningResult represents the result of sanctions screening
type ScreeningResult struct {
	EntityID           string                    `json:"entity_id"`
	EntityType         string                    `json:"entity_type"`        // individual, entity, vessel, address
	ScreeningID        string                    `json:"screening_id"`       // Unique screening reference
	IsMatch            bool                      `json:"is_match"`           // Any matches found
	RiskLevel          SanctionsRiskLevel        `json:"risk_level"`         // low, medium, high, blocked
	Matches            []SanctionsMatch          `json:"matches"`            // All matches found
	FalsePositives     []SanctionsMatch          `json:"false_positives"`    // Previously identified false positives
	RequiresReview     bool                      `json:"requires_review"`    // Manual review required
	AutoApproved       bool                      `json:"auto_approved"`      // Automatically cleared
	ScreeningTime      time.Duration             `json:"screening_time"`     // Time taken for screening
	ListsScreened      []string                  `json:"lists_screened"`     // Which lists were checked
	ComplianceScore    int                       `json:"compliance_score"`   // 0-100 compliance score
	RecommendedAction  string                    `json:"recommended_action"` // approve, review, block
	ScreenedAt         time.Time                 `json:"screened_at"`
	LastUpdated        time.Time                 `json:"last_updated"`
	ReviewerComments   string                    `json:"reviewer_comments"`
	ApprovalStatus     string                    `json:"approval_status"`    // pending, approved, denied
}

// SanctionsMatch represents a potential sanctions list match
type SanctionsMatch struct {
	MatchID           string                 `json:"match_id"`
	List              SanctionsList          `json:"list"`               // Which sanctions list
	EntryID           string                 `json:"entry_id"`           // ID on the sanctions list
	MatchedName       string                 `json:"matched_name"`       // Name that matched
	ListedName        string                 `json:"listed_name"`        // Name on the list
	MatchScore        float64                `json:"match_score"`        // 0.0-1.0 confidence score
	MatchType         SanctionsMatchType     `json:"match_type"`         // exact, fuzzy, alias, phonetic
	MatchedFields     []string               `json:"matched_fields"`     // Which fields matched
	AdditionalInfo    map[string]string      `json:"additional_info"`    // Additional data from list
	IsFalsePositive   bool                   `json:"is_false_positive"`  // Marked as false positive
	ReviewRequired    bool                   `json:"review_required"`    // Requires human review
	Severity          SanctionsMatchSeverity `json:"severity"`           // low, medium, high, critical
	CountryRisk       string                 `json:"country_risk"`       // Associated country risk
	ProgramType       string                 `json:"program_type"`       // Counter-terrorism, counter-narcotics, etc.
	EffectiveDate     time.Time              `json:"effective_date"`     // When sanctions became effective
	ExpiryDate        *time.Time             `json:"expiry_date,omitempty"` // When sanctions expire (if applicable)
}

// Enums for sanctions screening
type SanctionsRiskLevel int

const (
	RiskLevelLow SanctionsRiskLevel = iota
	RiskLevelMedium
	RiskLevelHigh
	RiskLevelBlocked
)

func (r SanctionsRiskLevel) String() string {
	switch r {
	case RiskLevelLow:
		return "low"
	case RiskLevelMedium:
		return "medium"
	case RiskLevelHigh:
		return "high"
	case RiskLevelBlocked:
		return "blocked"
	default:
		return "unknown"
	}
}

type SanctionsList int

const (
	SanctionsListOFAC_SDN SanctionsList = iota   // OFAC Specially Designated Nationals
	SanctionsListOFAC_Sectoral                   // OFAC Sectoral Sanctions
	SanctionsListOFAC_NonSDN                     // OFAC Non-SDN Sanctions
	SanctionsListUN                              // UN Security Council Sanctions
	SanctionsListEU                              // EU Consolidated List
	SanctionsListUK_HMT                          // UK HM Treasury
	SanctionsListCanada                          // Canada Sanctions
	SanctionsListAustralia                       // Australia DFAT
	SanctionsListBIS_DPL                         // Bureau of Industry & Security Denied Persons
	SanctionsListBIS_EL                          // Bureau of Industry & Security Entity List
	SanctionsListFBI_Most_Wanted                 // FBI Most Wanted
	SanctionsListInterpol                        // Interpol Red Notices
	SanctionsListPEP                             // Politically Exposed Persons
	SanctionsListCustom                          // Custom internal lists
)

func (s SanctionsList) String() string {
	switch s {
	case SanctionsListOFAC_SDN:
		return "OFAC_SDN"
	case SanctionsListOFAC_Sectoral:
		return "OFAC_Sectoral"
	case SanctionsListOFAC_NonSDN:
		return "OFAC_NonSDN"
	case SanctionsListUN:
		return "UN_Sanctions"
	case SanctionsListEU:
		return "EU_Consolidated"
	case SanctionsListUK_HMT:
		return "UK_HMT"
	case SanctionsListCanada:
		return "Canada_Sanctions"
	case SanctionsListAustralia:
		return "Australia_DFAT"
	case SanctionsListBIS_DPL:
		return "BIS_DPL"
	case SanctionsListBIS_EL:
		return "BIS_EL"
	case SanctionsListFBI_Most_Wanted:
		return "FBI_Most_Wanted"
	case SanctionsListInterpol:
		return "Interpol"
	case SanctionsListPEP:
		return "PEP"
	case SanctionsListCustom:
		return "Custom"
	default:
		return "Unknown"
	}
}

type SanctionsMatchType int

const (
	MatchTypeExact SanctionsMatchType = iota
	MatchTypeFuzzy
	MatchTypeAlias
	MatchTypePhonetic
	MatchTypePartial
	MatchTypeTransliteration
)

func (m SanctionsMatchType) String() string {
	switch m {
	case MatchTypeExact:
		return "exact"
	case MatchTypeFuzzy:
		return "fuzzy"
	case MatchTypeAlias:
		return "alias"
	case MatchTypePhonetic:
		return "phonetic"
	case MatchTypePartial:
		return "partial"
	case MatchTypeTransliteration:
		return "transliteration"
	default:
		return "unknown"
	}
}

type SanctionsMatchSeverity int

const (
	SeverityLow SanctionsMatchSeverity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s SanctionsMatchSeverity) String() string {
	switch s {
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	case SeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// SanctionsEntity represents an entity to be screened
type SanctionsEntity struct {
	ID             string            `json:"id"`
	Type           string            `json:"type"`              // individual, entity, vessel, address
	FullName       string            `json:"full_name"`
	FirstName      string            `json:"first_name"`
	LastName       string            `json:"last_name"`
	MiddleName     string            `json:"middle_name"`
	Aliases        []string          `json:"aliases"`
	DateOfBirth    *time.Time        `json:"date_of_birth,omitempty"`
	PlaceOfBirth   string            `json:"place_of_birth"`
	Nationality    []string          `json:"nationality"`
	Address        []string          `json:"address"`
	Country        string            `json:"country"`
	PhoneNumbers   []string          `json:"phone_numbers"`
	EmailAddresses []string          `json:"email_addresses"`
	Passport       []PassportInfo    `json:"passport"`
	NationalID     []NationalIDInfo  `json:"national_id"`
	CustomFields   map[string]string `json:"custom_fields"`
}

type PassportInfo struct {
	Number      string `json:"number"`
	Country     string `json:"country"`
	ExpiryDate  *time.Time `json:"expiry_date,omitempty"`
}

type NationalIDInfo struct {
	Number  string `json:"number"`
	Country string `json:"country"`
	Type    string `json:"type"` // ssn, national_id, tax_id, etc.
}

// SanctionsListEntry represents an entry on a sanctions list
type SanctionsListEntry struct {
	EntryID       string            `json:"entry_id"`
	List          SanctionsList     `json:"list"`
	FullName      string            `json:"full_name"`
	FirstName     string            `json:"first_name"`
	LastName      string            `json:"last_name"`
	Aliases       []string          `json:"aliases"`
	EntityType    string            `json:"entity_type"`
	DateOfBirth   *time.Time        `json:"date_of_birth,omitempty"`
	PlaceOfBirth  string            `json:"place_of_birth"`
	Nationality   []string          `json:"nationality"`
	Address       []string          `json:"address"`
	Program       string            `json:"program"`        // Sanctions program
	Remarks       string            `json:"remarks"`
	EffectiveDate time.Time         `json:"effective_date"`
	ExpiryDate    *time.Time        `json:"expiry_date,omitempty"`
	LastUpdated   time.Time         `json:"last_updated"`
	IsActive      bool              `json:"is_active"`
	Metadata      map[string]string `json:"metadata"`
}

// PerformSanctionsScreening performs comprehensive sanctions screening
func (sse *SanctionsScreeningEngine) PerformSanctionsScreening(
	ctx sdk.Context,
	entity SanctionsEntity,
	screeningOptions ScreeningOptions,
) (*ScreeningResult, error) {
	startTime := time.Now()
	
	// Generate screening ID
	screeningID := sse.generateScreeningID(ctx, entity.ID)
	
	result := &ScreeningResult{
		EntityID:        entity.ID,
		EntityType:      entity.Type,
		ScreeningID:     screeningID,
		IsMatch:         false,
		RiskLevel:       RiskLevelLow,
		Matches:         []SanctionsMatch{},
		FalsePositives:  []SanctionsMatch{},
		RequiresReview:  false,
		AutoApproved:    false,
		ListsScreened:   []string{},
		ComplianceScore: 100,
		ScreenedAt:      ctx.BlockTime(),
		LastUpdated:     ctx.BlockTime(),
		ApprovalStatus:  "pending",
	}

	// Screen against all configured lists
	listsToScreen := sse.determineListsToScreen(screeningOptions)
	
	for _, listType := range listsToScreen {
		result.ListsScreened = append(result.ListsScreened, listType.String())
		
		listMatches, err := sse.screenAgainstList(ctx, entity, listType, screeningOptions)
		if err != nil {
			sse.keeper.Logger(ctx).Error("Failed to screen against list", 
				"list", listType.String(), "error", err)
			continue
		}
		
		result.Matches = append(result.Matches, listMatches...)
	}

	// Check for previously identified false positives
	result.FalsePositives = sse.checkFalsePositives(ctx, entity.ID, result.Matches)
	
	// Filter out false positives from matches
	result.Matches = sse.filterFalsePositives(result.Matches, result.FalsePositives)
	
	// Calculate risk level and compliance score
	result.RiskLevel = sse.calculateRiskLevel(result.Matches)
	result.ComplianceScore = sse.calculateComplianceScore(result.Matches)
	result.IsMatch = len(result.Matches) > 0
	
	// Determine if review is required
	result.RequiresReview = sse.requiresManualReview(result.Matches, result.RiskLevel)
	
	// Auto-approve if possible
	if !result.RequiresReview && result.RiskLevel == RiskLevelLow {
		result.AutoApproved = true
		result.ApprovalStatus = "approved"
	}
	
	// Determine recommended action
	result.RecommendedAction = sse.determineRecommendedAction(result.RiskLevel, result.RequiresReview)
	
	result.ScreeningTime = time.Since(startTime)
	
	// Store screening result
	if err := sse.storeScreeningResult(ctx, result); err != nil {
		sse.keeper.Logger(ctx).Error("Failed to store screening result", "error", err)
	}
	
	// Emit screening event
	sse.emitScreeningEvent(ctx, result)
	
	return result, nil
}

// ScreeningOptions configures screening behavior
type ScreeningOptions struct {
	Lists              []SanctionsList `json:"lists"`                // Which lists to screen against
	MinMatchScore      float64         `json:"min_match_score"`      // Minimum score to consider a match
	EnableFuzzyMatch   bool            `json:"enable_fuzzy_match"`   // Enable fuzzy matching
	EnablePhonetic     bool            `json:"enable_phonetic"`      // Enable phonetic matching
	EnableTransliteration bool         `json:"enable_transliteration"` // Enable transliteration
	MaxMatches         int             `json:"max_matches"`          // Maximum matches to return
	IncludeFalsePositives bool         `json:"include_false_positives"` // Include known false positives
	BusinessRules      []string        `json:"business_rules"`       // Custom business rules to apply
	CountryRiskProfile string          `json:"country_risk_profile"` // low, medium, high
}

// screenAgainstList screens entity against a specific sanctions list
func (sse *SanctionsScreeningEngine) screenAgainstList(
	ctx sdk.Context,
	entity SanctionsEntity,
	listType SanctionsList,
	options ScreeningOptions,
) ([]SanctionsMatch, error) {
	// Get all entries from the specified list
	listEntries, err := sse.getListEntries(ctx, listType)
	if err != nil {
		return nil, err
	}

	var matches []SanctionsMatch

	for _, entry := range listEntries {
		// Skip inactive entries
		if !entry.IsActive {
			continue
		}

		// Check expiry date
		if entry.ExpiryDate != nil && entry.ExpiryDate.Before(ctx.BlockTime()) {
			continue
		}

		// Perform various types of matching
		entityMatches := sse.performMatching(entity, entry, options)
		matches = append(matches, entityMatches...)
	}

	// Sort matches by score (highest first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].MatchScore > matches[j].MatchScore
	})

	// Apply max matches limit
	if options.MaxMatches > 0 && len(matches) > options.MaxMatches {
		matches = matches[:options.MaxMatches]
	}

	return matches, nil
}

// performMatching performs various types of matching between entity and list entry
func (sse *SanctionsScreeningEngine) performMatching(
	entity SanctionsEntity,
	entry SanctionsListEntry,
	options ScreeningOptions,
) []SanctionsMatch {
	var matches []SanctionsMatch

	// Exact name matching
	if exactMatch := sse.exactNameMatch(entity, entry); exactMatch != nil {
		matches = append(matches, *exactMatch)
	}

	// Fuzzy name matching
	if options.EnableFuzzyMatch {
		if fuzzyMatches := sse.fuzzyNameMatch(entity, entry, options.MinMatchScore); len(fuzzyMatches) > 0 {
			matches = append(matches, fuzzyMatches...)
		}
	}

	// Alias matching
	if aliasMatches := sse.aliasMatch(entity, entry); len(aliasMatches) > 0 {
		matches = append(matches, aliasMatches...)
	}

	// Phonetic matching
	if options.EnablePhonetic {
		if phoneticMatches := sse.phoneticMatch(entity, entry); len(phoneticMatches) > 0 {
			matches = append(matches, phoneticMatches...)
		}
	}

	// Date of birth matching (for individuals)
	if entity.DateOfBirth != nil && entry.DateOfBirth != nil {
		for i := range matches {
			if sse.dateOfBirthMatch(*entity.DateOfBirth, *entry.DateOfBirth) {
				matches[i].MatchedFields = append(matches[i].MatchedFields, "date_of_birth")
				matches[i].MatchScore += 0.2 // Boost score for DOB match
			}
		}
	}

	// Address matching
	if len(entity.Address) > 0 && len(entry.Address) > 0 {
		for i := range matches {
			if sse.addressMatch(entity.Address, entry.Address) {
				matches[i].MatchedFields = append(matches[i].MatchedFields, "address")
				matches[i].MatchScore += 0.1 // Boost score for address match
			}
		}
	}

	// Nationality matching
	if len(entity.Nationality) > 0 && len(entry.Nationality) > 0 {
		for i := range matches {
			if sse.nationalityMatch(entity.Nationality, entry.Nationality) {
				matches[i].MatchedFields = append(matches[i].MatchedFields, "nationality")
				matches[i].MatchScore += 0.05 // Small boost for nationality match
			}
		}
	}

	// Filter by minimum match score
	var filteredMatches []SanctionsMatch
	for _, match := range matches {
		if match.MatchScore >= options.MinMatchScore {
			// Cap match score at 1.0
			if match.MatchScore > 1.0 {
				match.MatchScore = 1.0
			}
			filteredMatches = append(filteredMatches, match)
		}
	}

	return filteredMatches
}

// Matching algorithms

func (sse *SanctionsScreeningEngine) exactNameMatch(entity SanctionsEntity, entry SanctionsListEntry) *SanctionsMatch {
	entityName := sse.normalizeString(entity.FullName)
	entryName := sse.normalizeString(entry.FullName)
	
	if entityName == entryName {
		return &SanctionsMatch{
			MatchID:        sse.generateMatchID(),
			List:           entry.List,
			EntryID:        entry.EntryID,
			MatchedName:    entity.FullName,
			ListedName:     entry.FullName,
			MatchScore:     1.0,
			MatchType:      MatchTypeExact,
			MatchedFields:  []string{"full_name"},
			Severity:       sse.determineSeverity(entry.List, 1.0),
			ProgramType:    entry.Program,
			EffectiveDate:  entry.EffectiveDate,
			ExpiryDate:     entry.ExpiryDate,
			AdditionalInfo: map[string]string{
				"remarks": entry.Remarks,
				"entry_type": entry.EntityType,
			},
		}
	}
	
	return nil
}

func (sse *SanctionsScreeningEngine) fuzzyNameMatch(entity SanctionsEntity, entry SanctionsListEntry, minScore float64) []SanctionsMatch {
	var matches []SanctionsMatch
	
	// Compare full names
	if score := sse.calculateFuzzyScore(entity.FullName, entry.FullName); score >= minScore {
		matches = append(matches, SanctionsMatch{
			MatchID:        sse.generateMatchID(),
			List:           entry.List,
			EntryID:        entry.EntryID,
			MatchedName:    entity.FullName,
			ListedName:     entry.FullName,
			MatchScore:     score,
			MatchType:      MatchTypeFuzzy,
			MatchedFields:  []string{"full_name"},
			Severity:       sse.determineSeverity(entry.List, score),
			ProgramType:    entry.Program,
			EffectiveDate:  entry.EffectiveDate,
			ExpiryDate:     entry.ExpiryDate,
			AdditionalInfo: map[string]string{
				"remarks": entry.Remarks,
				"entry_type": entry.EntityType,
			},
		})
	}
	
	// Compare first/last names for individuals
	if entity.Type == "individual" && entry.EntityType == "individual" {
		if entity.FirstName != "" && entity.LastName != "" && entry.FirstName != "" && entry.LastName != "" {
			firstScore := sse.calculateFuzzyScore(entity.FirstName, entry.FirstName)
			lastScore := sse.calculateFuzzyScore(entity.LastName, entry.LastName)
			avgScore := (firstScore + lastScore) / 2
			
			if avgScore >= minScore {
				matches = append(matches, SanctionsMatch{
					MatchID:        sse.generateMatchID(),
					List:           entry.List,
					EntryID:        entry.EntryID,
					MatchedName:    fmt.Sprintf("%s %s", entity.FirstName, entity.LastName),
					ListedName:     fmt.Sprintf("%s %s", entry.FirstName, entry.LastName),
					MatchScore:     avgScore,
					MatchType:      MatchTypeFuzzy,
					MatchedFields:  []string{"first_name", "last_name"},
					Severity:       sse.determineSeverity(entry.List, avgScore),
					ProgramType:    entry.Program,
					EffectiveDate:  entry.EffectiveDate,
					ExpiryDate:     entry.ExpiryDate,
					AdditionalInfo: map[string]string{
						"first_name_score": fmt.Sprintf("%.2f", firstScore),
						"last_name_score": fmt.Sprintf("%.2f", lastScore),
						"remarks": entry.Remarks,
					},
				})
			}
		}
	}
	
	return matches
}

func (sse *SanctionsScreeningEngine) aliasMatch(entity SanctionsEntity, entry SanctionsListEntry) []SanctionsMatch {
	var matches []SanctionsMatch
	
	// Check entity aliases against entry name
	for _, alias := range entity.Aliases {
		if sse.normalizeString(alias) == sse.normalizeString(entry.FullName) {
			matches = append(matches, SanctionsMatch{
				MatchID:        sse.generateMatchID(),
				List:           entry.List,
				EntryID:        entry.EntryID,
				MatchedName:    alias,
				ListedName:     entry.FullName,
				MatchScore:     0.95, // High score for alias match
				MatchType:      MatchTypeAlias,
				MatchedFields:  []string{"alias"},
				Severity:       sse.determineSeverity(entry.List, 0.95),
				ProgramType:    entry.Program,
				EffectiveDate:  entry.EffectiveDate,
				ExpiryDate:     entry.ExpiryDate,
				AdditionalInfo: map[string]string{
					"matched_alias": alias,
					"remarks": entry.Remarks,
				},
			})
		}
	}
	
	// Check entity name against entry aliases
	for _, alias := range entry.Aliases {
		if sse.normalizeString(entity.FullName) == sse.normalizeString(alias) {
			matches = append(matches, SanctionsMatch{
				MatchID:        sse.generateMatchID(),
				List:           entry.List,
				EntryID:        entry.EntryID,
				MatchedName:    entity.FullName,
				ListedName:     alias,
				MatchScore:     0.95,
				MatchType:      MatchTypeAlias,
				MatchedFields:  []string{"name_vs_alias"},
				Severity:       sse.determineSeverity(entry.List, 0.95),
				ProgramType:    entry.Program,
				EffectiveDate:  entry.EffectiveDate,
				ExpiryDate:     entry.ExpiryDate,
				AdditionalInfo: map[string]string{
					"listed_alias": alias,
					"remarks": entry.Remarks,
				},
			})
		}
	}
	
	return matches
}

func (sse *SanctionsScreeningEngine) phoneticMatch(entity SanctionsEntity, entry SanctionsListEntry) []SanctionsMatch {
	var matches []SanctionsMatch
	
	// Simple phonetic matching using Soundex algorithm
	entitySoundex := sse.soundex(entity.FullName)
	entrySoundex := sse.soundex(entry.FullName)
	
	if entitySoundex != "" && entrySoundex != "" && entitySoundex == entrySoundex {
		// Phonetic match found, but calculate a more precise score
		score := sse.calculateFuzzyScore(entity.FullName, entry.FullName)
		if score < 0.7 { // Boost phonetic matches
			score = 0.7
		}
		
		matches = append(matches, SanctionsMatch{
			MatchID:        sse.generateMatchID(),
			List:           entry.List,
			EntryID:        entry.EntryID,
			MatchedName:    entity.FullName,
			ListedName:     entry.FullName,
			MatchScore:     score,
			MatchType:      MatchTypePhonetic,
			MatchedFields:  []string{"phonetic"},
			Severity:       sse.determineSeverity(entry.List, score),
			ProgramType:    entry.Program,
			EffectiveDate:  entry.EffectiveDate,
			ExpiryDate:     entry.ExpiryDate,
			AdditionalInfo: map[string]string{
				"entity_soundex": entitySoundex,
				"entry_soundex":  entrySoundex,
				"remarks":        entry.Remarks,
			},
		})
	}
	
	return matches
}

func (sse *SanctionsScreeningEngine) dateOfBirthMatch(entityDOB, entryDOB time.Time) bool {
	// Allow for small variations in date (±1 day for data entry errors)
	diff := entityDOB.Sub(entryDOB)
	return diff >= -24*time.Hour && diff <= 24*time.Hour
}

func (sse *SanctionsScreeningEngine) addressMatch(entityAddresses, entryAddresses []string) bool {
	for _, entityAddr := range entityAddresses {
		for _, entryAddr := range entryAddresses {
			if sse.calculateFuzzyScore(entityAddr, entryAddr) > 0.8 {
				return true
			}
		}
	}
	return false
}

func (sse *SanctionsScreeningEngine) nationalityMatch(entityNats, entryNats []string) bool {
	for _, entityNat := range entityNats {
		for _, entryNat := range entryNats {
			if strings.EqualFold(entityNat, entryNat) {
				return true
			}
		}
	}
	return false
}

// Utility functions

func (sse *SanctionsScreeningEngine) normalizeString(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	
	// Remove extra whitespace
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	
	// Remove common punctuation
	s = regexp.MustCompile(`[.,;:!?'"(){}[\]\-_]`).ReplaceAllString(s, "")
	
	// Remove diacritics (simplified)
	s = sse.removeDiacritics(s)
	
	return s
}

func (sse *SanctionsScreeningEngine) removeDiacritics(s string) string {
	// Simple diacritics removal (production would use more comprehensive mapping)
	replacements := map[rune]rune{
		'á': 'a', 'à': 'a', 'ä': 'a', 'â': 'a', 'ā': 'a', 'ã': 'a',
		'é': 'e', 'è': 'e', 'ë': 'e', 'ê': 'e', 'ē': 'e',
		'í': 'i', 'ì': 'i', 'ï': 'i', 'î': 'i', 'ī': 'i',
		'ó': 'o', 'ò': 'o', 'ö': 'o', 'ô': 'o', 'ō': 'o', 'õ': 'o',
		'ú': 'u', 'ù': 'u', 'ü': 'u', 'û': 'u', 'ū': 'u',
		'ñ': 'n', 'ç': 'c',
	}
	
	var result strings.Builder
	for _, r := range s {
		if replacement, found := replacements[unicode.ToLower(r)]; found {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	
	return result.String()
}

func (sse *SanctionsScreeningEngine) calculateFuzzyScore(s1, s2 string) float64 {
	// Levenshtein distance based similarity
	s1 = sse.normalizeString(s1)
	s2 = sse.normalizeString(s2)
	
	if s1 == s2 {
		return 1.0
	}
	
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}
	
	distance := sse.levenshteinDistance(s1, s2)
	maxLen := math.Max(float64(len(s1)), float64(len(s2)))
	
	return 1.0 - (float64(distance) / maxLen)
}

func (sse *SanctionsScreeningEngine) levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}
	
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}
	
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			
			matrix[i][j] = sse.min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}
	
	return matrix[len(s1)][len(s2)]
}

func (sse *SanctionsScreeningEngine) min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

// Simplified Soundex algorithm
func (sse *SanctionsScreeningEngine) soundex(s string) string {
	if len(s) == 0 {
		return ""
	}
	
	s = strings.ToUpper(sse.normalizeString(s))
	if len(s) == 0 {
		return ""
	}
	
	// Keep first letter
	result := string(s[0])
	
	// Mapping for consonants
	mapping := map[rune]rune{
		'B': '1', 'F': '1', 'P': '1', 'V': '1',
		'C': '2', 'G': '2', 'J': '2', 'K': '2', 'Q': '2', 'S': '2', 'X': '2', 'Z': '2',
		'D': '3', 'T': '3',
		'L': '4',
		'M': '5', 'N': '5',
		'R': '6',
	}
	
	var prevCode rune
	for _, r := range s[1:] {
		if code, exists := mapping[r]; exists {
			if code != prevCode {
				result += string(code)
				prevCode = code
			}
		} else {
			prevCode = 0 // Reset for vowels and 'H', 'W', 'Y'
		}
		
		if len(result) >= 4 {
			break
		}
	}
	
	// Pad with zeros if needed
	for len(result) < 4 {
		result += "0"
	}
	
	return result[:4]
}

// Risk assessment functions

func (sse *SanctionsScreeningEngine) calculateRiskLevel(matches []SanctionsMatch) SanctionsRiskLevel {
	if len(matches) == 0 {
		return RiskLevelLow
	}
	
	highestSeverity := SeverityLow
	highestScore := 0.0
	
	for _, match := range matches {
		if match.Severity > highestSeverity {
			highestSeverity = match.Severity
		}
		if match.MatchScore > highestScore {
			highestScore = match.MatchScore
		}
	}
	
	// Determine risk level based on severity and score
	if highestSeverity == SeverityCritical || highestScore >= 0.95 {
		return RiskLevelBlocked
	}
	if highestSeverity == SeverityHigh || highestScore >= 0.85 {
		return RiskLevelHigh
	}
	if highestSeverity == SeverityMedium || highestScore >= 0.70 {
		return RiskLevelMedium
	}
	
	return RiskLevelLow
}

func (sse *SanctionsScreeningEngine) calculateComplianceScore(matches []SanctionsMatch) int {
	if len(matches) == 0 {
		return 100
	}
	
	// Start with perfect score and deduct for matches
	score := 100
	
	for _, match := range matches {
		deduction := 0
		switch match.Severity {
		case SeverityCritical:
			deduction = 50
		case SeverityHigh:
			deduction = 30
		case SeverityMedium:
			deduction = 15
		case SeverityLow:
			deduction = 5
		}
		
		// Adjust based on match score
		deduction = int(float64(deduction) * match.MatchScore)
		score -= deduction
	}
	
	if score < 0 {
		score = 0
	}
	
	return score
}

func (sse *SanctionsScreeningEngine) determineSeverity(list SanctionsList, score float64) SanctionsMatchSeverity {
	// Base severity on the list type
	baseSeverity := SeverityMedium
	
	switch list {
	case SanctionsListOFAC_SDN, SanctionsListUN:
		baseSeverity = SeverityHigh
	case SanctionsListBIS_EL, SanctionsListFBI_Most_Wanted:
		baseSeverity = SeverityCritical
	case SanctionsListEU, SanctionsListUK_HMT:
		baseSeverity = SeverityHigh
	case SanctionsListPEP:
		baseSeverity = SeverityLow
	}
	
	// Adjust based on match score
	if score >= 0.95 {
		if baseSeverity < SeverityHigh {
			baseSeverity = SeverityHigh
		}
	} else if score >= 0.90 {
		if baseSeverity < SeverityMedium {
			baseSeverity = SeverityMedium
		}
	}
	
	return baseSeverity
}

func (sse *SanctionsScreeningEngine) requiresManualReview(matches []SanctionsMatch, riskLevel SanctionsRiskLevel) bool {
	if riskLevel >= RiskLevelHigh {
		return true
	}
	
	// Check for specific conditions that require review
	for _, match := range matches {
		if match.MatchType == MatchTypeFuzzy && match.MatchScore < 0.90 {
			return true
		}
		if match.MatchType == MatchTypePhonetic {
			return true
		}
		if match.List == SanctionsListPEP && match.MatchScore < 0.95 {
			return true
		}
	}
	
	return false
}

func (sse *SanctionsScreeningEngine) determineRecommendedAction(riskLevel SanctionsRiskLevel, requiresReview bool) string {
	if riskLevel == RiskLevelBlocked {
		return "block"
	}
	if requiresReview || riskLevel >= RiskLevelMedium {
		return "review"
	}
	return "approve"
}

// Storage and retrieval functions

func (sse *SanctionsScreeningEngine) getListEntries(ctx sdk.Context, listType SanctionsList) ([]SanctionsListEntry, error) {
	// In production, this would retrieve from database/external sources
	// For now, return mock data for demonstration
	return sse.getMockListEntries(listType), nil
}

func (sse *SanctionsScreeningEngine) getMockListEntries(listType SanctionsList) []SanctionsListEntry {
	// Mock data for testing - in production, this would come from official sources
	mockEntries := []SanctionsListEntry{
		{
			EntryID:       "OFAC-001",
			List:          SanctionsListOFAC_SDN,
			FullName:      "JOHN DOE",
			FirstName:     "JOHN",
			LastName:      "DOE",
			EntityType:    "individual",
			Nationality:   []string{"XX"},
			Program:       "TERRORISM",
			EffectiveDate: time.Now().AddDate(-1, 0, 0),
			IsActive:      true,
		},
		{
			EntryID:       "OFAC-002",
			List:          SanctionsListOFAC_SDN,
			FullName:      "BAD COMPANY LLC",
			EntityType:    "entity",
			Program:       "SANCTIONS",
			EffectiveDate: time.Now().AddDate(-2, 0, 0),
			IsActive:      true,
		},
	}
	
	// Filter by list type
	var filteredEntries []SanctionsListEntry
	for _, entry := range mockEntries {
		if entry.List == listType {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	
	return filteredEntries
}

func (sse *SanctionsScreeningEngine) checkFalsePositives(ctx sdk.Context, entityID string, matches []SanctionsMatch) []SanctionsMatch {
	// Check stored false positives for this entity
	// In production, this would query a database
	return []SanctionsMatch{} // Return empty for now
}

func (sse *SanctionsScreeningEngine) filterFalsePositives(matches, falsePositives []SanctionsMatch) []SanctionsMatch {
	fpMap := make(map[string]bool)
	for _, fp := range falsePositives {
		fpMap[fp.EntryID] = true
	}
	
	var filtered []SanctionsMatch
	for _, match := range matches {
		if !fpMap[match.EntryID] {
			filtered = append(filtered, match)
		}
	}
	
	return filtered
}

func (sse *SanctionsScreeningEngine) storeScreeningResult(ctx sdk.Context, result *ScreeningResult) error {
	store := ctx.KVStore(sse.keeper.storeKey)
	key := []byte("sanctions_screening:" + result.ScreeningID)
	bz := sse.keeper.cdc.MustMarshal(result)
	store.Set(key, bz)
	return nil
}

func (sse *SanctionsScreeningEngine) determineListsToScreen(options ScreeningOptions) []SanctionsList {
	if len(options.Lists) > 0 {
		return options.Lists
	}
	
	// Default lists to screen
	return []SanctionsList{
		SanctionsListOFAC_SDN,
		SanctionsListOFAC_Sectoral,
		SanctionsListUN,
		SanctionsListEU,
		SanctionsListBIS_EL,
	}
}

// Utility ID generation functions

func (sse *SanctionsScreeningEngine) generateScreeningID(ctx sdk.Context, entityID string) string {
	timestamp := ctx.BlockTime().Unix()
	return fmt.Sprintf("SCR-%s-%d", entityID, timestamp)
}

func (sse *SanctionsScreeningEngine) generateMatchID() string {
	// Simple UUID-like generation
	return fmt.Sprintf("MAT-%d", time.Now().UnixNano())
}

func (sse *SanctionsScreeningEngine) emitScreeningEvent(ctx sdk.Context, result *ScreeningResult) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"sanctions_screening_completed",
			sdk.NewAttribute("screening_id", result.ScreeningID),
			sdk.NewAttribute("entity_id", result.EntityID),
			sdk.NewAttribute("risk_level", result.RiskLevel.String()),
			sdk.NewAttribute("matches_found", strconv.Itoa(len(result.Matches))),
			sdk.NewAttribute("compliance_score", strconv.Itoa(result.ComplianceScore)),
			sdk.NewAttribute("requires_review", strconv.FormatBool(result.RequiresReview)),
			sdk.NewAttribute("recommended_action", result.RecommendedAction),
		),
	)
}

// Integration functions for Trade Finance and Remittance modules

// ScreenTradeFinanceParties screens all parties involved in a trade finance transaction
func (k Keeper) ScreenTradeFinanceParties(ctx sdk.Context, lcID string) (*ScreeningResult, error) {
	engine := NewSanctionsScreeningEngine(&k)
	
	lc, found := k.GetLetterOfCredit(ctx, lcID)
	if !found {
		return nil, types.ErrLCNotFound
	}
	
	// Screen applicant
	applicant, found := k.GetTradeParty(ctx, lc.ApplicantId)
	if !found {
		return nil, types.ErrPartyNotFound
	}
	
	entity := SanctionsEntity{
		ID:       applicant.PartyId,
		Type:     "entity",
		FullName: applicant.Name,
		Country:  applicant.Country,
		Address:  []string{applicant.Address},
	}
	
	options := ScreeningOptions{
		MinMatchScore:    0.75,
		EnableFuzzyMatch: true,
		EnablePhonetic:   true,
		MaxMatches:       10,
	}
	
	return engine.PerformSanctionsScreening(ctx, entity, options)
}

// ScreenRemittanceParties screens parties in a remittance transfer
func (k Keeper) ScreenRemittanceParties(ctx sdk.Context, transferID string) (*ScreeningResult, error) {
	engine := NewSanctionsScreeningEngine(&k)
	
	// This would integrate with the remittance module to get transfer details
	// For now, return a mock result
	entity := SanctionsEntity{
		ID:       transferID,
		Type:     "individual",
		FullName: "Test User",
	}
	
	options := ScreeningOptions{
		MinMatchScore:    0.80,
		EnableFuzzyMatch: true,
		MaxMatches:       5,
	}
	
	return engine.PerformSanctionsScreening(ctx, entity, options)
}