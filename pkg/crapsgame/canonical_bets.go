package crapsgame

import (
	"sort"
)

type BetCategory string

const (
	LineBets        BetCategory = "Line Bets"
	ComeBets        BetCategory = "Come Bets"
	OddsBets        BetCategory = "Odds Bets"
	FieldBets       BetCategory = "Field Bets"
	PlaceBets       BetCategory = "Place Bets"
	BuyBets         BetCategory = "Buy Bets"
	LayBets         BetCategory = "Lay Bets"
	PlaceToLoseBets BetCategory = "Place-to-Lose Bets"
	HardWayBets     BetCategory = "Hard Way Bets"
	PropositionBets BetCategory = "Proposition Bets"
	HornBets        BetCategory = "Horn Bets"
	HopBets         BetCategory = "Hop Bets"
	BigBets         BetCategory = "Big Bets"
	CombinationBets BetCategory = "Combination Bets"
)

// BetResolutionFunc defines the function signature for resolving a bet
// Returns (win, payout, remove):
//
//	win: true if the bet wins on this roll
//	payout: amount to pay (not including original bet)
//	remove: true if the bet should be removed after this roll
type BetResolutionFunc func(bet *Bet, roll *Roll, state GameState) (win bool, payout float64, remove bool)

// CanonicalBetDefinition represents a complete bet definition with all necessary information
// Now includes a Resolve function pointer for custom resolution logic
type CanonicalBetDefinition struct {
	Name              string
	Category          BetCategory
	Description       string
	Payout            string
	WorkingBehavior   string
	OneRoll           bool
	PayoutNumerator   int
	PayoutDenominator int
	ValidNumbers      []int
	RequiresPoint     bool
	RequiresComeOut   bool
	HouseEdge         float64
	Commission        float64 // Commission rate (0.05 for 5%)
}

// CanonicalBetDefinitions contains ALL bet definitions in a single source of truth
var CanonicalBetDefinitions = map[string]CanonicalBetDefinition{
	// Line Bets
	"PASS_LINE": {
		Name:              "Pass Line",
		Category:          LineBets,
		Description:       "Bet that shooter will win (7 or 11 on come out, then make point)",
		Payout:            "1:1",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   1,
		PayoutDenominator: 1,
		ValidNumbers:      []int{},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         1.41,
		Commission:        0.0,
	},
	"DONT_PASS": {
		Name:              "Don't Pass",
		Category:          LineBets,
		Description:       "Bet that shooter will lose (2, 3 on come out, then seven out)",
		Payout:            "1:1 (12 is push)",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   1,
		PayoutDenominator: 1,
		ValidNumbers:      []int{},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         1.36,
		Commission:        0.0,
	},

	// Come Bets
	"COME": {
		Name:              "Come",
		Category:          ComeBets,
		Description:       "Bet that next roll will be 7 or 11, or establish a point and make it",
		Payout:            "1:1",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   1,
		PayoutDenominator: 1,
		ValidNumbers:      []int{},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         1.41,
		Commission:        0.0,
	},
	"DONT_COME": {
		Name:              "Don't Come",
		Category:          ComeBets,
		Description:       "Bet that next roll will be 2, 3, or seven out before making point",
		Payout:            "1:1 (12 is push)",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   1,
		PayoutDenominator: 1,
		ValidNumbers:      []int{},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         1.36,
		Commission:        0.0,
	},

	// Odds Bets
	"PASS_ODDS": {
		Name:              "Pass Odds",
		Category:          OddsBets,
		Description:       "Odds behind pass line (only when point is established)",
		Payout:            "2:1 on 4/10, 3:2 on 5/9, 6:5 on 6/8",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   0, // Variable based on point
		PayoutDenominator: 0, // Variable based on point
		ValidNumbers:      []int{},
		RequiresPoint:     true,
		RequiresComeOut:   false,
		HouseEdge:         0.0,
		Commission:        0.0,
	},
	"DONT_PASS_ODDS": {
		Name:              "Don't Pass Odds",
		Category:          OddsBets,
		Description:       "Odds behind don't pass (only when point is established)",
		Payout:            "1:2 on 4/10, 2:3 on 5/9, 5:6 on 6/8",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   0, // Variable based on point
		PayoutDenominator: 0, // Variable based on point
		ValidNumbers:      []int{},
		RequiresPoint:     true,
		RequiresComeOut:   false,
		HouseEdge:         0.0,
		Commission:        0.0,
	},
	"COME_ODDS": {
		Name:              "Come Odds",
		Category:          OddsBets,
		Description:       "Odds behind come bets",
		Payout:            "Same as pass odds",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   0,
		PayoutDenominator: 0,
		ValidNumbers:      []int{},
		RequiresPoint:     true,
		RequiresComeOut:   false,
		HouseEdge:         0.0,
		Commission:        0.0,
	},
	"DONT_COME_ODDS": {
		Name:              "Don't Come Odds",
		Category:          OddsBets,
		Description:       "Odds behind don't come bets",
		Payout:            "Same as don't pass odds",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   0,
		PayoutDenominator: 0,
		ValidNumbers:      []int{},
		RequiresPoint:     true,
		RequiresComeOut:   false,
		HouseEdge:         0.0,
		Commission:        0.0,
	},

	// Field Bets
	"FIELD": {
		Name:              "Field",
		Category:          FieldBets,
		Description:       "Bet that next roll will be 2, 3, 4, 9, 10, 11, or 12",
		Payout:            "1:1 (2 pays 2:1, 12 pays 3:1)",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   1, // Default, but 2 and 12 are special
		PayoutDenominator: 1,
		ValidNumbers:      []int{}, // Field bets don't require specific numbers from player
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         2.78,
		Commission:        0.0,
	},

	// Place Bets
	"PLACE_4": {
		Name:              "Place 4",
		Category:          PlaceBets,
		Description:       "Bet that 4 will be rolled before 7",
		Payout:            "9:5",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   9,
		PayoutDenominator: 5,
		ValidNumbers:      []int{4},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         6.67,
		Commission:        0.0,
	},
	"PLACE_5": {
		Name:              "Place 5",
		Category:          PlaceBets,
		Description:       "Bet that 5 will be rolled before 7",
		Payout:            "7:5",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   7,
		PayoutDenominator: 5,
		ValidNumbers:      []int{5},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.0,
		Commission:        0.0,
	},
	"PLACE_6": {
		Name:              "Place 6",
		Category:          PlaceBets,
		Description:       "Bet that 6 will be rolled before 7",
		Payout:            "7:6",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   7,
		PayoutDenominator: 6,
		ValidNumbers:      []int{6},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         1.52,
		Commission:        0.0,
	},
	"PLACE_8": {
		Name:              "Place 8",
		Category:          PlaceBets,
		Description:       "Bet that 8 will be rolled before 7",
		Payout:            "7:6",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   7,
		PayoutDenominator: 6,
		ValidNumbers:      []int{8},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         1.52,
		Commission:        0.0,
	},
	"PLACE_9": {
		Name:              "Place 9",
		Category:          PlaceBets,
		Description:       "Bet that 9 will be rolled before 7",
		Payout:            "7:5",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   7,
		PayoutDenominator: 5,
		ValidNumbers:      []int{9},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.0,
		Commission:        0.0,
	},
	"PLACE_10": {
		Name:              "Place 10",
		Category:          PlaceBets,
		Description:       "Bet that 10 will be rolled before 7",
		Payout:            "9:5",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   9,
		PayoutDenominator: 5,
		ValidNumbers:      []int{10},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         6.67,
		Commission:        0.0,
	},

	// Buy Bets
	"BUY_4": {
		Name:              "Buy 4",
		Category:          BuyBets,
		Description:       "Buy bet on 4 (true odds with 5% commission)",
		Payout:            "2:1 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   2,
		PayoutDenominator: 1,
		ValidNumbers:      []int{4},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.76,
		Commission:        0.05,
	},
	"BUY_5": {
		Name:              "Buy 5",
		Category:          BuyBets,
		Description:       "Buy bet on 5 (true odds with 5% commission)",
		Payout:            "3:2 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   3,
		PayoutDenominator: 2,
		ValidNumbers:      []int{5},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.76,
		Commission:        0.05,
	},
	"BUY_6": {
		Name:              "Buy 6",
		Category:          BuyBets,
		Description:       "Buy bet on 6 (true odds with 5% commission)",
		Payout:            "6:5 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   6,
		PayoutDenominator: 5,
		ValidNumbers:      []int{6},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.76,
		Commission:        0.05,
	},
	"BUY_8": {
		Name:              "Buy 8",
		Category:          BuyBets,
		Description:       "Buy bet on 8 (true odds with 5% commission)",
		Payout:            "6:5 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   6,
		PayoutDenominator: 5,
		ValidNumbers:      []int{8},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.76,
		Commission:        0.05,
	},
	"BUY_9": {
		Name:              "Buy 9",
		Category:          BuyBets,
		Description:       "Buy bet on 9 (true odds with 5% commission)",
		Payout:            "3:2 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   3,
		PayoutDenominator: 2,
		ValidNumbers:      []int{9},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.76,
		Commission:        0.05,
	},
	"BUY_10": {
		Name:              "Buy 10",
		Category:          BuyBets,
		Description:       "Buy bet on 10 (true odds with 5% commission)",
		Payout:            "2:1 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   2,
		PayoutDenominator: 1,
		ValidNumbers:      []int{10},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.76,
		Commission:        0.05,
	},

	// Lay Bets
	"LAY_4": {
		Name:              "Lay 4",
		Category:          LayBets,
		Description:       "Lay bet against 4 (bet 7 comes before 4)",
		Payout:            "1:2 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   1,
		PayoutDenominator: 2,
		ValidNumbers:      []int{4},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         2.44,
		Commission:        0.05,
	},
	"LAY_5": {
		Name:              "Lay 5",
		Category:          LayBets,
		Description:       "Lay bet against 5 (bet 7 comes before 5)",
		Payout:            "2:3 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   2,
		PayoutDenominator: 3,
		ValidNumbers:      []int{5},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         3.23,
		Commission:        0.05,
	},
	"LAY_6": {
		Name:              "Lay 6",
		Category:          LayBets,
		Description:       "Lay bet against 6 (bet 7 comes before 6)",
		Payout:            "5:6 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   5,
		PayoutDenominator: 6,
		ValidNumbers:      []int{6},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.0,
		Commission:        0.05,
	},
	"LAY_8": {
		Name:              "Lay 8",
		Category:          LayBets,
		Description:       "Lay bet against 8 (bet 7 comes before 8)",
		Payout:            "5:6 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   5,
		PayoutDenominator: 6,
		ValidNumbers:      []int{8},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.0,
		Commission:        0.05,
	},
	"LAY_9": {
		Name:              "Lay 9",
		Category:          LayBets,
		Description:       "Lay bet against 9 (bet 7 comes before 9)",
		Payout:            "2:3 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   2,
		PayoutDenominator: 3,
		ValidNumbers:      []int{9},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         3.23,
		Commission:        0.05,
	},
	"LAY_10": {
		Name:              "Lay 10",
		Category:          LayBets,
		Description:       "Lay bet against 10 (bet 7 comes before 10)",
		Payout:            "1:2 (minus commission)",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   1,
		PayoutDenominator: 2,
		ValidNumbers:      []int{10},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         2.44,
		Commission:        0.05,
	},

	// Place-to-Lose Bets
	"PLACE_TO_LOSE_4": {
		Name:              "Place-to-Lose 4",
		Category:          PlaceToLoseBets,
		Description:       "Place-to-lose bet against 4 (bet 7 comes before 4)",
		Payout:            "1:2",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   1,
		PayoutDenominator: 2,
		ValidNumbers:      []int{4},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         2.44,
		Commission:        0.0,
	},
	"PLACE_TO_LOSE_5": {
		Name:              "Place-to-Lose 5",
		Category:          PlaceToLoseBets,
		Description:       "Place-to-lose bet against 5 (bet 7 comes before 5)",
		Payout:            "2:3",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   2,
		PayoutDenominator: 3,
		ValidNumbers:      []int{5},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         3.23,
		Commission:        0.0,
	},
	"PLACE_TO_LOSE_6": {
		Name:              "Place-to-Lose 6",
		Category:          PlaceToLoseBets,
		Description:       "Place-to-lose bet against 6 (bet 7 comes before 6)",
		Payout:            "5:6",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   5,
		PayoutDenominator: 6,
		ValidNumbers:      []int{6},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.0,
		Commission:        0.0,
	},
	"PLACE_TO_LOSE_8": {
		Name:              "Place-to-Lose 8",
		Category:          PlaceToLoseBets,
		Description:       "Place-to-lose bet against 8 (bet 7 comes before 8)",
		Payout:            "5:6",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   5,
		PayoutDenominator: 6,
		ValidNumbers:      []int{8},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         4.0,
		Commission:        0.0,
	},
	"PLACE_TO_LOSE_9": {
		Name:              "Place-to-Lose 9",
		Category:          PlaceToLoseBets,
		Description:       "Place-to-lose bet against 9 (bet 7 comes before 9)",
		Payout:            "2:3",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   2,
		PayoutDenominator: 3,
		ValidNumbers:      []int{9},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         3.23,
		Commission:        0.0,
	},
	"PLACE_TO_LOSE_10": {
		Name:              "Place-to-Lose 10",
		Category:          PlaceToLoseBets,
		Description:       "Place-to-lose bet against 10 (bet 7 comes before 10)",
		Payout:            "1:2",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   1,
		PayoutDenominator: 2,
		ValidNumbers:      []int{10},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         2.44,
		Commission:        0.0,
	},

	// Hard Way Bets
	"HARD_4": {
		Name:              "Hard 4",
		Category:          HardWayBets,
		Description:       "Bet that 4 will be rolled as 2-2",
		Payout:            "7:1",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   7,
		PayoutDenominator: 1,
		ValidNumbers:      []int{4},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HARD_6": {
		Name:              "Hard 6",
		Category:          HardWayBets,
		Description:       "Bet that 6 will be rolled as 3-3",
		Payout:            "9:1",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   9,
		PayoutDenominator: 1,
		ValidNumbers:      []int{6},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         9.09,
		Commission:        0.0,
	},
	"HARD_8": {
		Name:              "Hard 8",
		Category:          HardWayBets,
		Description:       "Bet that 8 will be rolled as 4-4",
		Payout:            "9:1",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   9,
		PayoutDenominator: 1,
		ValidNumbers:      []int{8},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         9.09,
		Commission:        0.0,
	},
	"HARD_10": {
		Name:              "Hard 10",
		Category:          HardWayBets,
		Description:       "Bet that 10 will be rolled as 5-5",
		Payout:            "7:1",
		WorkingBehavior:   "CONDITIONAL",
		OneRoll:           false,
		PayoutNumerator:   7,
		PayoutDenominator: 1,
		ValidNumbers:      []int{10},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},

	// Proposition Bets
	"ANY_SEVEN": {
		Name:              "Any Seven",
		Category:          PropositionBets,
		Description:       "Bet that next roll will be 7",
		Payout:            "4:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   4,
		PayoutDenominator: 1,
		ValidNumbers:      []int{7},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         16.67,
		Commission:        0.0,
	},
	"ANY_CRAPS": {
		Name:              "Any Craps",
		Category:          PropositionBets,
		Description:       "Bet that next roll will be 2, 3, or 12",
		Payout:            "7:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   7,
		PayoutDenominator: 1,
		ValidNumbers:      []int{2, 3, 12},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"ELEVEN": {
		Name:              "Eleven",
		Category:          PropositionBets,
		Description:       "Bet that next roll will be 11 (6-5)",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{11},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"ACE_DEUCE": {
		Name:              "Ace-Deuce",
		Category:          PropositionBets,
		Description:       "Bet that next roll will be 3 (1-2)",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{3},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"ACES": {
		Name:              "Aces",
		Category:          PropositionBets,
		Description:       "Bet that next roll will be 2 (1-1)",
		Payout:            "30:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   30,
		PayoutDenominator: 1,
		ValidNumbers:      []int{2},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         13.89,
		Commission:        0.0,
	},
	"BOXCARS": {
		Name:              "Boxcars",
		Category:          PropositionBets,
		Description:       "Bet that next roll will be 12 (6-6)",
		Payout:            "30:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   30,
		PayoutDenominator: 1,
		ValidNumbers:      []int{12},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         13.89,
		Commission:        0.0,
	},

	// Horn Bets
	"HORN_HIGH_2": {
		Name:              "Horn High 2",
		Category:          HornBets,
		Description:       "Horn bet with extra on 2",
		Payout:            "2 pays 27:4, 3 pays 3:1, 11 pays 3:1, 12 pays 3:1",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   27, // For 2
		PayoutDenominator: 4,
		ValidNumbers:      []int{2, 3, 11, 12},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         12.5,
		Commission:        0.0,
	},
	"HORN_HIGH_3": {
		Name:              "Horn High 3",
		Category:          HornBets,
		Description:       "Horn bet with extra on 3",
		Payout:            "2 pays 3:1, 3 pays 15:1, 11 pays 3:1, 12 pays 3:1",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   15, // For 3
		PayoutDenominator: 1,
		ValidNumbers:      []int{2, 3, 11, 12},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         12.5,
		Commission:        0.0,
	},
	"HORN_HIGH_11": {
		Name:              "Horn High 11",
		Category:          HornBets,
		Description:       "Horn bet with extra on 11",
		Payout:            "2 pays 3:1, 3 pays 3:1, 11 pays 15:1, 12 pays 3:1",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   15, // For 11
		PayoutDenominator: 1,
		ValidNumbers:      []int{2, 3, 11, 12},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         12.5,
		Commission:        0.0,
	},
	"HORN_HIGH_12": {
		Name:              "Horn High 12",
		Category:          HornBets,
		Description:       "Horn bet with extra on 12",
		Payout:            "2 pays 3:1, 3 pays 3:1, 11 pays 3:1, 12 pays 27:4",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   27, // For 12
		PayoutDenominator: 4,
		ValidNumbers:      []int{2, 3, 11, 12},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         12.5,
		Commission:        0.0,
	},

	// Hop Bets (Easy Hops)
	"HOP_1_2": {
		Name:              "Hop 1-2",
		Category:          HopBets,
		Description:       "Bet that next roll will be 1-2",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{3},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_1_3": {
		Name:              "Hop 1-3",
		Category:          HopBets,
		Description:       "Bet that next roll will be 1-3",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{4},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_1_4": {
		Name:              "Hop 1-4",
		Category:          HopBets,
		Description:       "Bet that next roll will be 1-4",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{5},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_1_5": {
		Name:              "Hop 1-5",
		Category:          HopBets,
		Description:       "Bet that next roll will be 1-5",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{6},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_1_6": {
		Name:              "Hop 1-6",
		Category:          HopBets,
		Description:       "Bet that next roll will be 1-6",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{7},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_2_3": {
		Name:              "Hop 2-3",
		Category:          HopBets,
		Description:       "Bet that next roll will be 2-3",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{5},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_2_4": {
		Name:              "Hop 2-4",
		Category:          HopBets,
		Description:       "Bet that next roll will be 2-4",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{6},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_2_5": {
		Name:              "Hop 2-5",
		Category:          HopBets,
		Description:       "Bet that next roll will be 2-5",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{7},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_2_6": {
		Name:              "Hop 2-6",
		Category:          HopBets,
		Description:       "Bet that next roll will be 2-6",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{8},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_3_4": {
		Name:              "Hop 3-4",
		Category:          HopBets,
		Description:       "Bet that next roll will be 3-4",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{7},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_3_5": {
		Name:              "Hop 3-5",
		Category:          HopBets,
		Description:       "Bet that next roll will be 3-5",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{8},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_3_6": {
		Name:              "Hop 3-6",
		Category:          HopBets,
		Description:       "Bet that next roll will be 3-6",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{9},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_4_5": {
		Name:              "Hop 4-5",
		Category:          HopBets,
		Description:       "Bet that next roll will be 4-5",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{9},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_4_6": {
		Name:              "Hop 4-6",
		Category:          HopBets,
		Description:       "Bet that next roll will be 4-6",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{10},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
	"HOP_5_6": {
		Name:              "Hop 5-6",
		Category:          HopBets,
		Description:       "Bet that next roll will be 5-6",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{11},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},

	// Big Bets
	"BIG_6": {
		Name:              "Big 6",
		Category:          BigBets,
		Description:       "Bet that 6 will be rolled before 7",
		Payout:            "1:1 (worse than place 6)",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   1,
		PayoutDenominator: 1,
		ValidNumbers:      []int{6},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         9.09,
		Commission:        0.0,
	},
	"BIG_8": {
		Name:              "Big 8",
		Category:          BigBets,
		Description:       "Bet that 8 will be rolled before 7",
		Payout:            "1:1 (worse than place 8)",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   1,
		PayoutDenominator: 1,
		ValidNumbers:      []int{8},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         9.09,
		Commission:        0.0,
	},

	// Place Numbers (composite bet type)
	"PLACE_NUMBERS": {
		Name:              "Place Numbers",
		Category:          PlaceBets,
		Description:       "Place bet on multiple numbers (requires specific numbers)",
		Payout:            "Variable based on numbers placed",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   0, // Variable
		PayoutDenominator: 0, // Variable
		ValidNumbers:      []int{4, 5, 6, 8, 9, 10},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         0.0, // Variable
		Commission:        0.0,
	},

	// Place Inside (covers 5, 6, 8, 9)
	"PLACE_INSIDE": {
		Name:              "Place Inside",
		Category:          PlaceBets,
		Description:       "Place bet covering 5, 6, 8, 9",
		Payout:            "7:5 on 5/9, 7:6 on 6/8",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   0, // Variable
		PayoutDenominator: 0, // Variable
		ValidNumbers:      []int{5, 6, 8, 9},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         0.0, // Variable
		Commission:        0.0,
	},

	// Place Outside (covers 4, 5, 9, 10)
	"PLACE_OUTSIDE": {
		Name:              "Place Outside",
		Category:          PlaceBets,
		Description:       "Place bet covering 4, 5, 9, 10",
		Payout:            "9:5 on 4/10, 7:5 on 5/9",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   0, // Variable
		PayoutDenominator: 0, // Variable
		ValidNumbers:      []int{4, 5, 9, 10},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         0.0, // Variable
		Commission:        0.0,
	},

	// All Hardways (covers all hard way bets)
	"ALL_HARDWAYS": {
		Name:              "All Hardways",
		Category:          HardWayBets,
		Description:       "Hard way bets on 4, 6, 8, 10",
		Payout:            "7:1 on 4/10, 9:1 on 6/8",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   0, // Variable
		PayoutDenominator: 0, // Variable
		ValidNumbers:      []int{4, 6, 8, 10},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         0.0, // Variable
		Commission:        0.0,
	},

	// Horn (covers 2, 3, 11, 12)
	"HORN": {
		Name:              "Horn",
		Category:          HornBets,
		Description:       "Horn bet covering 2, 3, 11, 12",
		Payout:            "2 pays 27:4, 3 pays 3:1, 11 pays 3:1, 12 pays 3:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           false,
		PayoutNumerator:   3,
		PayoutDenominator: 1,
		ValidNumbers:      []int{2, 3, 11, 12},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         12.5,
		Commission:        0.0,
	},

	// World (combination of any 7 and any craps)
	"WORLD": {
		Name:              "World",
		Category:          CombinationBets,
		Description:       "Combination of any 7 and any craps",
		Payout:            "4:1 on 7, 7:1 on any craps",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           false,
		PayoutNumerator:   4,
		PayoutDenominator: 1,
		ValidNumbers:      []int{2, 3, 7, 11, 12},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         13.33,
		Commission:        0.0,
	},

	// Hop Hard 6 (3-3)
	"HOP_HARD_6": {
		Name:              "Hop Hard 6",
		Category:          HopBets,
		Description:       "Hop bet on hard 6 (3-3)",
		Payout:            "30:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   30,
		PayoutDenominator: 1,
		ValidNumbers:      []int{6},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         13.89,
		Commission:        0.0,
	},

	// Hop Easy 8 (any combination except 4-4)
	"HOP_EASY_8": {
		Name:              "Hop Easy 8",
		Category:          HopBets,
		Description:       "Hop bet on easy 8 (any combination except 4-4)",
		Payout:            "15:1",
		WorkingBehavior:   "ONE_ROLL",
		OneRoll:           true,
		PayoutNumerator:   15,
		PayoutDenominator: 1,
		ValidNumbers:      []int{8},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},

	// Combination Bets
	"C_AND_E": {
		Name:              "C and E",
		Category:          CombinationBets,
		Description:       "Combination of craps and eleven",
		Payout:            "3:1 (any craps pays 7:1, eleven pays 15:1)",
		WorkingBehavior:   "ALWAYS",
		OneRoll:           false,
		PayoutNumerator:   3,
		PayoutDenominator: 1,
		ValidNumbers:      []int{2, 3, 11, 12},
		RequiresPoint:     false,
		RequiresComeOut:   false,
		HouseEdge:         11.11,
		Commission:        0.0,
	},
}

// GetBetDefinition returns the canonical bet definition for a given bet type string
func GetBetDefinition(betType string) (CanonicalBetDefinition, bool) {
	bet, ok := CanonicalBetDefinitions[betType]
	return bet, ok
}

// GetAllBetTypes returns a slice of all canonical bet type strings
func GetAllBetTypes() []string {
	betTypes := make([]string, 0, len(CanonicalBetDefinitions))
	for k := range CanonicalBetDefinitions {
		betTypes = append(betTypes, k)
	}
	return betTypes
}

// GetBetsByCategory returns all bets organized by category
func GetBetsByCategory() map[BetCategory][]string {
	categories := make(map[BetCategory][]string)
	for betType, info := range CanonicalBetDefinitions {
		categories[info.Category] = append(categories[info.Category], betType)
	}
	return categories
}

// GetOneRollBets returns all one-roll bets
func GetOneRollBets() []string {
	var oneRollBets []string
	for betType, info := range CanonicalBetDefinitions {
		if info.OneRoll {
			oneRollBets = append(oneRollBets, betType)
		}
	}
	return oneRollBets
}

// GetAlwaysWorkingBets returns all always-working bets
func GetAlwaysWorkingBets() []string {
	var alwaysWorkingBets []string
	for betType, info := range CanonicalBetDefinitions {
		if !info.OneRoll {
			alwaysWorkingBets = append(alwaysWorkingBets, betType)
		}
	}
	return alwaysWorkingBets
}

// GetBetsByHouseEdge returns bets sorted by house edge (lowest to highest)
func GetBetsByHouseEdge() []string {
	betTypes := make([]string, 0, len(CanonicalBetDefinitions))
	for betType := range CanonicalBetDefinitions {
		betTypes = append(betTypes, betType)
	}
	// Sort by house edge
	type betWithEdge struct {
		betType   string
		houseEdge float64
	}
	var betList []betWithEdge
	for _, betType := range betTypes {
		bet := CanonicalBetDefinitions[betType]
		betList = append(betList, betWithEdge{betType, bet.HouseEdge})
	}
	sort.Slice(betList, func(i, j int) bool {
		return betList[i].houseEdge < betList[j].houseEdge
	})
	result := make([]string, len(betList))
	for i, b := range betList {
		result[i] = b.betType
	}
	return result
}

// --- RESOLVER FUNCTIONS FOR ALL CANONICAL BET TYPES ---

// Generic resolver for Place bets (handles Place 4, 5, 6, 8, 9, 10)
func resolvePlaceBet(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	if state != StatePoint {
		return false, 0, false
	}
	if len(bet.Numbers) == 0 {
		return false, 0, false
	}
	num := bet.Numbers[0]
	if roll.Total == num {
		def, _ := CanonicalBetDefinitions[bet.Type]
		payout := bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator)
		return true, payout, true
	} else if roll.Total == 7 {
		return false, 0, true
	}
	return false, 0, false
}

// Generic resolver for Buy bets
func resolveBuyBet(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	if state != StatePoint {
		return false, 0, false
	}
	if len(bet.Numbers) == 0 {
		return false, 0, false
	}
	num := bet.Numbers[0]
	if roll.Total == num {
		def, _ := CanonicalBetDefinitions[bet.Type]
		gross := bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator)
		commission := bet.Amount * def.Commission
		return true, gross - commission, true
	} else if roll.Total == 7 {
		return false, 0, true
	}
	return false, 0, false
}

// Generic resolver for Lay bets
func resolveLayBet(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	if state != StatePoint {
		return false, 0, false
	}
	if len(bet.Numbers) == 0 {
		return false, 0, false
	}
	num := bet.Numbers[0]
	if roll.Total == 7 {
		def, _ := CanonicalBetDefinitions[bet.Type]
		gross := bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator)
		commission := bet.Amount * def.Commission
		return true, gross - commission, true
	} else if roll.Total == num {
		return false, 0, true
	}
	return false, 0, false
}

// Generic resolver for Place-to-Lose bets
func resolvePlaceToLoseBet(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	if state != StatePoint {
		return false, 0, false
	}
	if len(bet.Numbers) == 0 {
		return false, 0, false
	}
	num := bet.Numbers[0]
	if roll.Total == 7 {
		def, _ := CanonicalBetDefinitions[bet.Type]
		payout := bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator)
		return true, payout, true
	} else if roll.Total == num {
		return false, 0, true
	}
	return false, 0, false
}

// Generic resolver for Hardway bets
func resolveHardwayBet(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	if state != StatePoint {
		return false, 0, false
	}
	if len(bet.Numbers) == 0 {
		return false, 0, false
	}
	num := bet.Numbers[0]
	if roll.Total == num && roll.IsHard {
		def, _ := CanonicalBetDefinitions[bet.Type]
		payout := bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator)
		return true, payout, true
	} else if roll.Total == num && !roll.IsHard {
		return false, 0, true
	} else if roll.Total == 7 {
		return false, 0, true
	}
	return false, 0, false
}

// Pass Line resolver
func resolvePassLine(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	def, _ := CanonicalBetDefinitions[bet.Type]
	if state == StateComeOut {
		if roll.Total == 7 || roll.Total == 11 {
			return true, bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator), true
		} else if roll.Total == 2 || roll.Total == 3 || roll.Total == 12 {
			return false, 0, true
		}
		return false, 0, false
	} else if state == StatePoint {
		if len(bet.Numbers) == 0 {
			return false, 0, false
		}
		point := bet.Numbers[0]
		if roll.Total == point {
			return true, bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator), true
		} else if roll.Total == 7 {
			return false, 0, true
		}
	}
	return false, 0, false
}

// Don't Pass resolver
func resolveDontPass(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	def, _ := CanonicalBetDefinitions[bet.Type]
	if state == StateComeOut {
		if roll.Total == 2 || roll.Total == 3 {
			return true, bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator), true
		} else if roll.Total == 12 {
			return false, 0, true // push
		} else if roll.Total == 7 || roll.Total == 11 {
			return false, 0, true
		}
		return false, 0, false
	} else if state == StatePoint {
		if len(bet.Numbers) == 0 {
			return false, 0, false
		}
		point := bet.Numbers[0]
		if roll.Total == 7 {
			return true, bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator), true
		} else if roll.Total == point {
			return false, 0, true
		}
	}
	return false, 0, false
}

// Field bet resolver
func resolveFieldBet(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	if roll.Total == 2 {
		return true, bet.Amount * 2, true
	} else if roll.Total == 12 {
		return true, bet.Amount * 3, true
	} else if roll.Total == 3 || roll.Total == 4 || roll.Total == 9 || roll.Total == 10 || roll.Total == 11 {
		return true, bet.Amount, true
	}
	return false, 0, true
}

// Any Seven resolver
func resolveAnySeven(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	def, _ := CanonicalBetDefinitions[bet.Type]
	if roll.Total == 7 {
		return true, bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator), true
	}
	return false, 0, true
}

// Any Craps resolver
func resolveAnyCraps(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	def, _ := CanonicalBetDefinitions[bet.Type]
	if roll.Total == 2 || roll.Total == 3 || roll.Total == 12 {
		return true, bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator), true
	}
	return false, 0, true
}

// Eleven resolver
func resolveEleven(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	def, _ := CanonicalBetDefinitions[bet.Type]
	if roll.Total == 11 {
		return true, bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator), true
	}
	return false, 0, true
}

// Ace-Deuce resolver
func resolveAceDeuce(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	def, _ := CanonicalBetDefinitions[bet.Type]
	if roll.Total == 3 {
		return true, bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator), true
	}
	return false, 0, true
}

// Aces resolver
func resolveAces(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	def, _ := CanonicalBetDefinitions[bet.Type]
	if roll.Total == 2 {
		return true, bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator), true
	}
	return false, 0, true
}

// Boxcars resolver
func resolveBoxcars(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	def, _ := CanonicalBetDefinitions[bet.Type]
	if roll.Total == 12 {
		return true, bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator), true
	}
	return false, 0, true
}

// --- HORN BETS RESOLVER ---
func resolveHornBet(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	// Horn bets win on 2, 3, 11, or 12, with different payouts for "high" numbers
	// Use bet.Type to determine which is the "high" number
	win := false
	payout := 0.0
	switch bet.Type {
	case "HORN":
		if roll.Total == 2 || roll.Total == 3 || roll.Total == 11 || roll.Total == 12 {
			// Standard horn payout: 3:1 for 3, 11, 12; 27:4 for 2
			if roll.Total == 2 || roll.Total == 12 {
				payout = bet.Amount * 27.0 / 4.0
			} else {
				payout = bet.Amount * 3.0
			}
			win = true
		}
	case "HORN_HIGH_2":
		if roll.Total == 2 {
			payout = bet.Amount * 27.0 / 4.0
			win = true
		} else if roll.Total == 3 || roll.Total == 11 || roll.Total == 12 {
			payout = bet.Amount * 3.0
			win = true
		}
	case "HORN_HIGH_3":
		if roll.Total == 3 {
			payout = bet.Amount * 15.0
			win = true
		} else if roll.Total == 2 || roll.Total == 11 || roll.Total == 12 {
			payout = bet.Amount * 3.0
			win = true
		}
	case "HORN_HIGH_11":
		if roll.Total == 11 {
			payout = bet.Amount * 15.0
			win = true
		} else if roll.Total == 2 || roll.Total == 3 || roll.Total == 12 {
			payout = bet.Amount * 3.0
			win = true
		}
	case "HORN_HIGH_12":
		if roll.Total == 12 {
			payout = bet.Amount * 27.0 / 4.0
			win = true
		} else if roll.Total == 2 || roll.Total == 3 || roll.Total == 11 {
			payout = bet.Amount * 3.0
			win = true
		}
	case "HORN_HIGH_ACE_DEUCE":
		if roll.Total == 3 {
			payout = bet.Amount * 15.0
			win = true
		} else if roll.Total == 2 || roll.Total == 11 || roll.Total == 12 {
			payout = bet.Amount * 3.0
			win = true
		}
	}
	return win, payout, true
}

// --- HOP BETS RESOLVER ---
func resolveHopBet(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	// Hop bets are one-roll bets on a specific dice combination
	// Use ValidNumbers for the total, and bet.Type for hard/easy
	def, _ := CanonicalBetDefinitions[bet.Type]
	if len(def.ValidNumbers) == 0 {
		return false, 0, true
	}
	if roll.Total == def.ValidNumbers[0] {
		// For hard hops, check if IsHard is required
		if bet.Type == "HOP_HARD_6" && roll.Total == 6 && roll.IsHard {
			payout := bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator)
			return true, payout, true
		}
		if bet.Type == "HOP_EASY_8" && roll.Total == 8 && !roll.IsHard {
			payout := bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator)
			return true, payout, true
		}
		// For generic hops, just pay out
		payout := bet.Amount * float64(def.PayoutNumerator) / float64(def.PayoutDenominator)
		return true, payout, true
	}
	return false, 0, true
}

// --- COMBINATION BETS RESOLVER ---
func resolveCombinationBet(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	switch bet.Type {
	case "PLACE_NUMBERS":
		for _, num := range []int{4, 5, 6, 8, 9, 10} {
			tempBet := *bet
			tempBet.Numbers = []int{num}
			win, payout, remove := resolvePlaceBet(&tempBet, roll, state)
			if win {
				return win, payout, remove
			}
		}
		return false, 0, false
	case "PLACE_INSIDE":
		for _, num := range []int{5, 6, 8, 9} {
			tempBet := *bet
			tempBet.Numbers = []int{num}
			win, payout, remove := resolvePlaceBet(&tempBet, roll, state)
			if win {
				return win, payout, remove
			}
		}
		return false, 0, false
	case "PLACE_OUTSIDE":
		for _, num := range []int{4, 5, 9, 10} {
			tempBet := *bet
			tempBet.Numbers = []int{num}
			win, payout, remove := resolvePlaceBet(&tempBet, roll, state)
			if win {
				return win, payout, remove
			}
		}
		return false, 0, false
	case "ALL_HARDWAYS":
		for _, num := range []int{4, 6, 8, 10} {
			tempBet := *bet
			tempBet.Numbers = []int{num}
			win, payout, remove := resolveHardwayBet(&tempBet, roll, state)
			if win {
				return win, payout, remove
			}
		}
		return false, 0, false
	}
	return false, 0, false
}

// --- REGISTER HORN AND HOP BETS IN BetTypeResolvers ---
var BetTypeResolvers = map[string]BetResolutionFunc{
	// Place bets
	"PLACE_4":  resolvePlaceBet,
	"PLACE_5":  resolvePlaceBet,
	"PLACE_6":  resolvePlaceBet,
	"PLACE_8":  resolvePlaceBet,
	"PLACE_9":  resolvePlaceBet,
	"PLACE_10": resolvePlaceBet,
	// Buy bets
	"BUY_4":  resolveBuyBet,
	"BUY_5":  resolveBuyBet,
	"BUY_6":  resolveBuyBet,
	"BUY_8":  resolveBuyBet,
	"BUY_9":  resolveBuyBet,
	"BUY_10": resolveBuyBet,
	// Lay bets
	"LAY_4":  resolveLayBet,
	"LAY_5":  resolveLayBet,
	"LAY_6":  resolveLayBet,
	"LAY_8":  resolveLayBet,
	"LAY_9":  resolveLayBet,
	"LAY_10": resolveLayBet,
	// Place-to-lose bets
	"PLACE_TO_LOSE_4":  resolvePlaceToLoseBet,
	"PLACE_TO_LOSE_5":  resolvePlaceToLoseBet,
	"PLACE_TO_LOSE_6":  resolvePlaceToLoseBet,
	"PLACE_TO_LOSE_8":  resolvePlaceToLoseBet,
	"PLACE_TO_LOSE_9":  resolvePlaceToLoseBet,
	"PLACE_TO_LOSE_10": resolvePlaceToLoseBet,
	// Hardway bets
	"HARD_4":  resolveHardwayBet,
	"HARD_6":  resolveHardwayBet,
	"HARD_8":  resolveHardwayBet,
	"HARD_10": resolveHardwayBet,
	// Pass Line
	"PASS_LINE": resolvePassLine,
	// Don't Pass
	"DONT_PASS": resolveDontPass,
	// Field
	"FIELD": resolveFieldBet,
	// Any Seven
	"ANY_SEVEN": resolveAnySeven,
	// Any Craps
	"ANY_CRAPS": resolveAnyCraps,
	// Eleven
	"ELEVEN": resolveEleven,
	// Ace-Deuce
	"ACE_DEUCE": resolveAceDeuce,
	// Aces
	"ACES": resolveAces,
	// Boxcars
	"BOXCARS": resolveBoxcars,
	// Horn bets
	"HORN":                resolveHornBet,
	"HORN_HIGH_2":         resolveHornBet,
	"HORN_HIGH_3":         resolveHornBet,
	"HORN_HIGH_11":        resolveHornBet,
	"HORN_HIGH_12":        resolveHornBet,
	"HORN_HIGH_ACE_DEUCE": resolveHornBet,
	// Hop bets
	"HOP":        resolveHopBet,
	"HOP_HARD_6": resolveHopBet,
	"HOP_EASY_8": resolveHopBet,
	"HOP_1_2":    resolveHopBet,
	"HOP_1_3":    resolveHopBet,
	"HOP_1_4":    resolveHopBet,
	"HOP_1_5":    resolveHopBet,
	"HOP_1_6":    resolveHopBet,
	"HOP_2_3":    resolveHopBet,
	"HOP_2_4":    resolveHopBet,
	"HOP_2_5":    resolveHopBet,
	"HOP_2_6":    resolveHopBet,
	"HOP_3_4":    resolveHopBet,
	"HOP_3_5":    resolveHopBet,
	"HOP_3_6":    resolveHopBet,
	"HOP_4_5":    resolveHopBet,
	"HOP_4_6":    resolveHopBet,
	"HOP_5_6":    resolveHopBet,
	// Combination bets
	"PLACE_NUMBERS": resolveCombinationBet,
	"PLACE_INSIDE":  resolveCombinationBet,
	"PLACE_OUTSIDE": resolveCombinationBet,
	"ALL_HARDWAYS":  resolveCombinationBet,
}

// Central entry point for resolving a bet
func ResolveBet(bet *Bet, roll *Roll, state GameState) (bool, float64, bool) {
	resolver, ok := BetTypeResolvers[bet.Type]
	if !ok {
		// fallback or error: unknown bet type
		return false, 0, false
	}
	return resolver(bet, roll, state)
}
