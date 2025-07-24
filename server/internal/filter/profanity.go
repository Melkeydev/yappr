package filter

import (
	"regexp"
	"strings"
)

// ProfanityFilter contains the filtering logic for inappropriate content
type ProfanityFilter struct {
	bannedWords []string
	patterns    []*regexp.Regexp
}

// NewProfanityFilter creates a new profanity filter with a comprehensive list of banned words
func NewProfanityFilter() *ProfanityFilter {
	filter := &ProfanityFilter{
		bannedWords: getBannedWords(),
	}
	
	// Compile regex patterns for banned words (case-insensitive, with word boundaries)
	for _, word := range filter.bannedWords {
		// Create pattern that matches the word with optional special characters
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(word) + `\b`)
		filter.patterns = append(filter.patterns, pattern)
		
		// Also check for l33t speak and common substitutions
		substituted := applyCommonSubstitutions(word)
		if substituted != word {
			substPattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(substituted) + `\b`)
			filter.patterns = append(filter.patterns, substPattern)
		}
	}
	
	return filter
}

// ContainsProfanity checks if the given text contains profanity
func (pf *ProfanityFilter) ContainsProfanity(text string) bool {
	// Normalize the text - remove extra spaces, convert to lowercase for checking
	normalized := strings.ToLower(strings.TrimSpace(text))
	
	// Check against all patterns
	for _, pattern := range pf.patterns {
		if pattern.MatchString(normalized) {
			return true
		}
	}
	
	return false
}

// applyCommonSubstitutions applies common character substitutions used to bypass filters
func applyCommonSubstitutions(word string) string {
	substitutions := map[string]string{
		"a": "@",
		"e": "3",
		"i": "1",
		"o": "0",
		"s": "$",
		"t": "7",
	}
	
	result := word
	for original, substitute := range substitutions {
		result = strings.ReplaceAll(result, original, substitute)
	}
	
	return result
}

// getBannedWords returns a comprehensive list of inappropriate words
// This includes profanity, slurs, and racist terms
func getBannedWords() []string {
	return []string{
		// Common profanity
		"fuck", "shit", "damn", "bitch", "asshole", "bastard", "crap",
		"piss", "dickhead", "jackass", "dumbass", "bullshit",
		
		// Stronger profanity
		"motherfucker", "cocksucker", "son of a bitch", "piece of shit",
		
		// Sexual content
		"porn", "sex", "naked", "nude", "xxx", "adult", "escort",
		
		// Racist slurs and terms (censored for safety)
		"nigger", "nigga", "negro", "spic", "wetback", "chink", "gook",
		"kike", "hymie", "raghead", "towelhead", "sand nigger",
		"cracker", "honky", "whitey", "gringo", "beaner", "border hopper",
		
		// Homophobic slurs
		"faggot", "fag", "dyke", "homo", "queer", "tranny",
		
		// Religious slurs
		"kike", "hymie", "christ killer", "jihad", "terrorist",
		
		// Disability slurs
		"retard", "retarded", "spastic", "cripple", "invalid",
		
		// Common offensive terms
		"nazi", "hitler", "genocide", "kill yourself", "kys",
		"suicide", "die", "death", "murder", "rape",
		
		// Hate speech indicators
		"hate", "supremacy", "master race", "inferior race",
		"pure blood", "ethnic cleansing",
		
		// Drug references
		"cocaine", "heroin", "meth", "crack", "weed", "marijuana",
		"drugs", "dealer", "pusher",
		
		// Violence
		"violence", "beating", "assault", "abuse", "torture",
		"bomb", "explosion", "terrorist", "attack",
		
		// Spam/inappropriate
		"admin", "moderator", "official", "staff", "bot",
		"advertisement", "promotion", "scam", "phishing",
		
		// Common variations and misspellings
		"fuk", "shyt", "btch", "azz", "phuck", "biatch",
		"n1gger", "n1gga", "f4ggot", "f4g", "sh1t", "fck",
	}
}