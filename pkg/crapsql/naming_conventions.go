package crapsql

// NamingConventions defines the standard naming conventions for CrapsQL
// This file establishes a single source of truth for bet type naming

// BetTypeNamingConventions maps bet types to their canonical string representations
var BetTypeNamingConventions = map[string]string{
	// Line bets
	"PASS_LINE": "PASS_LINE",
	"DONT_PASS": "DONT_PASS",
	"COME":      "COME",
	"DONT_COME": "DONT_COME",

	// Field and one-roll bets
	"FIELD":     "FIELD",
	"ANY_SEVEN": "ANY_SEVEN",
	"ANY_CRAPS": "ANY_CRAPS",
	"ELEVEN":    "ELEVEN",
	"ACE_DEUCE": "ACE_DEUCE",
	"ACES":      "ACES",
	"BOXCARS":   "BOXCARS",

	// Place bets
	"PLACE_4":       "PLACE_4",
	"PLACE_5":       "PLACE_5",
	"PLACE_6":       "PLACE_6",
	"PLACE_8":       "PLACE_8",
	"PLACE_9":       "PLACE_9",
	"PLACE_10":      "PLACE_10",
	"PLACE_NUMBERS": "PLACE_NUMBERS",
	"PLACE_INSIDE":  "PLACE_INSIDE",
	"PLACE_OUTSIDE": "PLACE_OUTSIDE",

	// Place-to-lose bets
	"PLACE_TO_LOSE_4":  "PLACE_TO_LOSE_4",
	"PLACE_TO_LOSE_5":  "PLACE_TO_LOSE_5",
	"PLACE_TO_LOSE_6":  "PLACE_TO_LOSE_6",
	"PLACE_TO_LOSE_8":  "PLACE_TO_LOSE_8",
	"PLACE_TO_LOSE_9":  "PLACE_TO_LOSE_9",
	"PLACE_TO_LOSE_10": "PLACE_TO_LOSE_10",

	// Hard ways
	"HARD_4":       "HARD_4",
	"HARD_6":       "HARD_6",
	"HARD_8":       "HARD_8",
	"HARD_10":      "HARD_10",
	"ALL_HARDWAYS": "ALL_HARDWAYS",

	// Odds bets (standardized naming)
	"PASS_ODDS":      "PASS_ODDS",      // was "ODDS"
	"DONT_PASS_ODDS": "DONT_PASS_ODDS", // was "DONT_ODDS"
	"COME_ODDS":      "COME_ODDS",
	"DONT_COME_ODDS": "DONT_COME_ODDS",

	// Buy bets
	"BUY_4":  "BUY_4",
	"BUY_5":  "BUY_5",
	"BUY_6":  "BUY_6",
	"BUY_8":  "BUY_8",
	"BUY_9":  "BUY_9",
	"BUY_10": "BUY_10",

	// Lay bets
	"LAY_4":  "LAY_4",
	"LAY_5":  "LAY_5",
	"LAY_6":  "LAY_6",
	"LAY_8":  "LAY_8",
	"LAY_9":  "LAY_9",
	"LAY_10": "LAY_10",

	// Big 6/8
	"BIG_6": "BIG_6",
	"BIG_8": "BIG_8",

	// Hop bets (standardized naming)
	"HOP_1_2": "HOP_1_2",
	"HOP_1_3": "HOP_1_3",
	"HOP_1_4": "HOP_1_4",
	"HOP_1_5": "HOP_1_5",
	"HOP_1_6": "HOP_1_6",
	"HOP_2_3": "HOP_2_3",
	"HOP_2_4": "HOP_2_4",
	"HOP_2_5": "HOP_2_5",
	"HOP_2_6": "HOP_2_6",
	"HOP_3_4": "HOP_3_4",
	"HOP_3_5": "HOP_3_5",
	"HOP_3_6": "HOP_3_6",
	"HOP_4_5": "HOP_4_5",
	"HOP_4_6": "HOP_4_6",
	"HOP_5_6": "HOP_5_6",

	// Special hop bets
	"HOP":        "HOP",
	"HOP_HARD_6": "HOP_HARD_6",
	"HOP_EASY_8": "HOP_EASY_8",

	// Proposition bets
	"WORLD":   "WORLD",
	"C_AND_E": "C_AND_E",
	"HORN":    "HORN",

	// Horn high bets (standardized naming)
	"HORN_HIGH_2":  "HORN_HIGH_2",
	"HORN_HIGH_3":  "HORN_HIGH_3",
	"HORN_HIGH_11": "HORN_HIGH_11", // was "HORN_HIGH_YO"
	"HORN_HIGH_12": "HORN_HIGH_12", // was "HORN_HIGH_ACE_DEUCE"
}


// GetCanonicalBetTypeName returns the canonical name for a bet type
func GetCanonicalBetTypeName(betType string) string {
	if canonical, exists := BetTypeNamingConventions[betType]; exists {
		return canonical
	}
	return betType // Return original if not found in conventions
}

// IsValidBetTypeName checks if a bet type name follows the established conventions
func IsValidBetTypeName(betType string) bool {
	_, exists := BetTypeNamingConventions[betType]
	return exists
}
