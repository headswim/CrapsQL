package crapsql

import (
	"fmt"
)

// BetType is assumed to be defined in types.go
// type BetType int

var (
	stringToBetType = make(map[string]BetType)
	betTypeToString = make(map[BetType]string)
)

// StringToBetType converts a string to a BetType enum value
func StringToBetType(betString string) (BetType, error) {
	bt, ok := stringToBetType[betString]
	if !ok {
		return 0, fmt.Errorf("unknown bet type: %s", betString)
	}
	return bt, nil
}

// BetTypeToString converts a BetType enum value to its string representation
func BetTypeToString(betType BetType) (string, error) {
	str, ok := betTypeToString[betType]
	if !ok {
		return "", fmt.Errorf("unknown BetType: %v", betType)
	}
	return str, nil
}

// ValidateBetType checks if a string is a valid bet type
func ValidateBetType(betType string) error {
	if _, ok := stringToBetType[betType]; !ok {
		return fmt.Errorf("invalid bet type: %s", betType)
	}
	return nil
}

// GetAllRegisteredBetTypes returns all registered bet type strings
func GetAllRegisteredBetTypes() []string {
	result := make([]string, 0, len(stringToBetType))
	for k := range stringToBetType {
		result = append(result, k)
	}
	return result
}

// IsValidBetType returns true if the string is a valid bet type
func IsValidBetType(betType string) bool {
	_, ok := stringToBetType[betType]
	return ok
}

// TODO: Add a function to initialize the registry from canonical bet types

// initRegistry initializes the bet type registry with all canonical bet types
func initRegistry() {
	// Line Bets
	stringToBetType["PASS_LINE"] = BetPassLine
	stringToBetType["DONT_PASS"] = BetDontPass
	stringToBetType["COME"] = BetCome
	stringToBetType["DONT_COME"] = BetDontCome

	// Field and one-roll bets
	stringToBetType["FIELD"] = BetField
	stringToBetType["ANY_SEVEN"] = BetAnySeven
	stringToBetType["ANY_CRAPS"] = BetAnyCraps
	stringToBetType["ELEVEN"] = BetEleven
	stringToBetType["ACE_DEUCE"] = BetAceDeuce
	stringToBetType["ACES"] = BetAces
	stringToBetType["BOXCARS"] = BetBoxcars

	// Place bets
	stringToBetType["PLACE_4"] = BetPlace4
	stringToBetType["PLACE_5"] = BetPlace5
	stringToBetType["PLACE_6"] = BetPlace6
	stringToBetType["PLACE_8"] = BetPlace8
	stringToBetType["PLACE_9"] = BetPlace9
	stringToBetType["PLACE_10"] = BetPlace10
	stringToBetType["PLACE_NUMBERS"] = BetPlaceNumbers
	stringToBetType["PLACE_INSIDE"] = BetPlaceInside
	stringToBetType["PLACE_OUTSIDE"] = BetPlaceOutside

	// Hard ways
	stringToBetType["HARD_4"] = BetHard4
	stringToBetType["HARD_6"] = BetHard6
	stringToBetType["HARD_8"] = BetHard8
	stringToBetType["HARD_10"] = BetHard10
	stringToBetType["ALL_HARDWAYS"] = BetAllHardways

	// Odds bets
	stringToBetType["PASS_ODDS"] = BetPassOdds
	stringToBetType["DONT_PASS_ODDS"] = BetDontPassOdds
	stringToBetType["COME_ODDS"] = BetComeOdds
	stringToBetType["DONT_COME_ODDS"] = BetDontComeOdds

	// Buy/Lay
	stringToBetType["BUY_4"] = BetBuy4
	stringToBetType["BUY_10"] = BetBuy10
	stringToBetType["LAY_4"] = BetLay4
	stringToBetType["LAY_10"] = BetLay10

	// Big 6/8
	stringToBetType["BIG_6"] = BetBig6
	stringToBetType["BIG_8"] = BetBig8

	// Hop bets
	stringToBetType["HOP"] = BetHop
	stringToBetType["HOP_HARD_6"] = BetHopHard6
	stringToBetType["HOP_EASY_8"] = BetHopEasy8

	// Horn bets
	stringToBetType["HORN"] = BetHorn
	stringToBetType["HORN_HIGH_11"] = BetHornHigh11
	stringToBetType["HORN_HIGH_ACE_DEUCE"] = BetHornHighAceDeuce

	// Proposition bets
	stringToBetType["WORLD"] = BetWorld
	stringToBetType["C_AND_E"] = BetCAndE
	stringToBetType["HORN_HIGH_2"] = BetHornHigh2
	stringToBetType["HORN_HIGH_3"] = BetHornHigh3
	stringToBetType["HORN_HIGH_12"] = BetHornHigh12

	// Missing bet types from canonical definitions
	// Buy bets
	stringToBetType["BUY_5"] = BetBuy5
	stringToBetType["BUY_6"] = BetBuy6
	stringToBetType["BUY_8"] = BetBuy8
	stringToBetType["BUY_9"] = BetBuy9

	// Lay bets
	stringToBetType["LAY_5"] = BetLay5
	stringToBetType["LAY_6"] = BetLay6
	stringToBetType["LAY_8"] = BetLay8
	stringToBetType["LAY_9"] = BetLay9

	// Place-to-lose bets
	stringToBetType["PLACE_TO_LOSE_4"] = BetPlaceToLose4
	stringToBetType["PLACE_TO_LOSE_5"] = BetPlaceToLose5
	stringToBetType["PLACE_TO_LOSE_6"] = BetPlaceToLose6
	stringToBetType["PLACE_TO_LOSE_8"] = BetPlaceToLose8
	stringToBetType["PLACE_TO_LOSE_9"] = BetPlaceToLose9
	stringToBetType["PLACE_TO_LOSE_10"] = BetPlaceToLose10

	// Horn high bets
	stringToBetType["HORN_HIGH_2"] = BetHornHigh2
	stringToBetType["HORN_HIGH_3"] = BetHornHigh3
	stringToBetType["HORN_HIGH_11"] = BetHornHigh11
	stringToBetType["HORN_HIGH_12"] = BetHornHigh12

	// Hop bets (all combinations)
	stringToBetType["HOP_1_2"] = BetHop12
	stringToBetType["HOP_1_3"] = BetHop13
	stringToBetType["HOP_1_4"] = BetHop14
	stringToBetType["HOP_1_5"] = BetHop15
	stringToBetType["HOP_1_6"] = BetHop16
	stringToBetType["HOP_2_3"] = BetHop23
	stringToBetType["HOP_2_4"] = BetHop24
	stringToBetType["HOP_2_5"] = BetHop25
	stringToBetType["HOP_2_6"] = BetHop26
	stringToBetType["HOP_3_4"] = BetHop34
	stringToBetType["HOP_3_5"] = BetHop35
	stringToBetType["HOP_3_6"] = BetHop36
	stringToBetType["HOP_4_5"] = BetHop45
	stringToBetType["HOP_4_6"] = BetHop46
	stringToBetType["HOP_5_6"] = BetHop56

	// Create reverse mapping
	for str, bt := range stringToBetType {
		betTypeToString[bt] = str
	}
}

// init initializes the registry when the package is imported
func init() {
	initRegistry()
}
