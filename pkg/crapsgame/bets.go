package crapsgame

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

const (
	BuyLayCommissionRate = 0.05 // 5% commission for buy and lay bets
)

type BetResolution struct { // BetResolution handles all bet payout calculations
	table *Table
}

func NewBetResolution(table *Table) *BetResolution { // NewBetResolution creates a new bet resolution handler
	return &BetResolution{table: table}
}

func (br *BetResolution) calculateCommission(winnings float64, rate float64) float64 {
	if rate <= 0 || rate >= 1 {
		return 0 // Invalid rate, no commission
	}
	return winnings * rate
}

func (br *BetResolution) calculateWinningsWithCommission(betAmount float64, odds float64, commissionRate float64) float64 {
	// Validate commission rate
	if err := br.validateCommissionRate(commissionRate); err != nil {
		// If commission rate is invalid, return winnings without commission
		return betAmount * odds
	}

	grossWinnings := betAmount * odds
	commission := br.calculateCommission(grossWinnings, commissionRate)
	return grossWinnings - commission
}

func (br *BetResolution) getCommissionAmount(betAmount float64, odds float64, commissionRate float64) float64 {
	if err := br.validateCommissionRate(commissionRate); err != nil {
		return 0
	}
	grossWinnings := betAmount * odds
	return br.calculateCommission(grossWinnings, commissionRate)
}

func (br *BetResolution) formatCommissionMessage(betAmount float64, odds float64, commissionRate float64) string {
	if commissionRate <= 0 {
		return ""
	}
	commission := br.getCommissionAmount(betAmount, odds, commissionRate)
	return fmt.Sprintf(" (commission: $%.2f)", commission)
}

func (br *BetResolution) validateCommissionRate(rate float64) error {
	if rate < 0 {
		return fmt.Errorf("commission rate cannot be negative: %.2f", rate)
	}
	if rate >= 1 {
		return fmt.Errorf("commission rate cannot be 100%% or greater: %.2f", rate)
	}
	return nil
}

func (br *BetResolution) ResolveBets(roll *Roll) []string {
	var results []string

	br.updateBetWorkingStatus()

	for _, player := range br.table.Players {
		for _, bet := range player.Bets {
			// Validate bet state before processing
			if err := br.validateBetState(bet); err != nil {
				results = append(results, fmt.Sprintf("Invalid bet state for player %s: %v", player.ID, err))
				continue
			}

			if !bet.Working {
				continue
			}

			result := br.resolveBet(bet, roll, player)
			if result != "" {
				results = append(results, result)
			}
		}
	}

	br.cleanupResolvedBets()

	return results
}

func (br *BetResolution) resolveBet(bet *Bet, roll *Roll, player *Player) string {
	if err := br.validateBetState(bet); err != nil {
		return fmt.Sprintf("Invalid bet state: %v", err)
	}

	betDef, exists := CanonicalBetDefinitions[bet.Type]
	if !exists {
		return fmt.Sprintf("Unknown bet type: %s", bet.Type)
	}

	if !br.shouldBetBeWorking(bet, betDef) {
		br.setBetWorking(bet, false)
		return ""
	}

	if betDef.OneRoll {
		return br.resolveOneRollBet(bet, betDef, roll, player)
	}

	return br.resolveAlwaysWorkingBet(bet, betDef, roll, player)
}

func (br *BetResolution) shouldBetBeWorking(bet *Bet, betDef CanonicalBetDefinition) bool {
	switch betDef.WorkingBehavior {
	case "ONE_ROLL":
		return bet.Working
	case "ALWAYS":
		return bet.Working
	case "CONDITIONAL":
		return br.checkConditionalBetState(bet, br.table.State)
	default:
		return bet.Working
	}
}

func (br *BetResolution) resolveOneRollBet(bet *Bet, betDef CanonicalBetDefinition, roll *Roll, player *Player) string {
	if err := br.validateBetState(bet); err != nil {
		return fmt.Sprintf("Invalid bet state: %v", err)
	}

	for _, validNum := range betDef.ValidNumbers {
		if roll.Total == validNum {
			winnings := (bet.Amount * float64(betDef.PayoutNumerator)) / float64(betDef.PayoutDenominator)
			player.Bankroll += winnings

			br.setBetWorking(bet, false)
			br.removeOneRollBet(bet)

			return fmt.Sprintf("🎉 %s wins $%.2f (%d)", betDef.Name, winnings, roll.Total)
		}
	}

	br.setBetWorking(bet, false)
	br.removeOneRollBet(bet)

	return fmt.Sprintf("💥 %s loses $%.2f (%d)", betDef.Name, bet.Amount, roll.Total)
}

func (br *BetResolution) resolveAlwaysWorkingBet(bet *Bet, betDef CanonicalBetDefinition, roll *Roll, player *Player) string {
	if err := br.validateBetState(bet); err != nil {
		return fmt.Sprintf("Invalid bet state: %v", err)
	}

	switch betDef.Category {
	case "Line Bets":
		result := br.resolveLineBet(bet, betDef, roll, player)
		// Only persist if bet didn't win or lose (result is empty)
		if result == "" {
			br.persistAlwaysWorkingBet(bet)
		}
		return result
	case "Come Bets":
		result := br.resolveComeBet(bet, betDef, roll, player)
		// Only persist if bet didn't win or lose (result is empty)
		if result == "" {
			br.persistAlwaysWorkingBet(bet)
		}
		return result
	case "Odds Bets":
		result := br.resolveOddsBet(bet, betDef, roll, player)
		// Only persist if bet didn't win or lose (result is empty)
		if result == "" {
			br.persistAlwaysWorkingBet(bet)
		}
		return result
	case "Place Bets", "Buy Bets", "Lay Bets", "Place-to-Lose Bets":
		result := br.resolvePlaceStyleBet(bet, betDef, roll, player)
		// Only persist if bet didn't win or lose (result is empty)
		if result == "" {
			br.persistAlwaysWorkingBet(bet)
		}
		return result
	case "Hard Way Bets":
		result := br.resolveHardWayBet(bet, betDef, roll, player)
		// Hard way bets are removed by the resolveHardWay function when they win or lose
		// Only persist if no result (bet continues working)
		if result == "" {
			br.persistAlwaysWorkingBet(bet)
		}
		return result
	case "Big Bets":
		result := br.resolveBigBetWithDef(bet, betDef, roll, player)
		// Only persist if bet didn't win or lose (result is empty)
		if result == "" {
			br.persistAlwaysWorkingBet(bet)
		}
		return result
	case "Horn Bets":
		switch bet.Type {
		case "HORN":
			result := br.resolveHornBet(bet, roll, player)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HORN_HIGH_2":
			result := br.resolveHornHighBet(bet, roll, player, 2)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HORN_HIGH_3":
			result := br.resolveHornHighBet(bet, roll, player, 3)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HORN_HIGH_11":
			result := br.resolveHornHighBet(bet, roll, player, 11)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HORN_HIGH_12":
			result := br.resolveHornHighBet(bet, roll, player, 12)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		default:
			return fmt.Sprintf("Unknown horn bet type: %s", bet.Type)
		}
	case "Combination Bets":
		switch bet.Type {
		case "WORLD":
			result := br.resolveWorldBet(bet, roll, player)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "C_AND_E":
			result := br.resolveCAndE(bet, roll, player)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		default:
			return fmt.Sprintf("Unknown combination bet type: %s", bet.Type)
		}
	case "Horn High Bets":
		switch bet.Type {
		case "HORN_HIGH_2":
			result := br.resolveHornHighBet(bet, roll, player, 2)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HORN_HIGH_3":
			result := br.resolveHornHighBet(bet, roll, player, 3)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HORN_HIGH_11":
			result := br.resolveHornHighBet(bet, roll, player, 11)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HORN_HIGH_12":
			result := br.resolveHornHighBet(bet, roll, player, 12)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		default:
			return fmt.Sprintf("Unknown horn high bet type: %s", bet.Type)
		}
	case "Hop Bets":
		switch betDef.Name {
		case "HOP_1_2":
			result := br.resolveHopBet(bet, roll, player, 1, 2)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_1_3":
			result := br.resolveHopBet(bet, roll, player, 1, 3)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_1_4":
			result := br.resolveHopBet(bet, roll, player, 1, 4)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_1_5":
			result := br.resolveHopBet(bet, roll, player, 1, 5)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_1_6":
			result := br.resolveHopBet(bet, roll, player, 1, 6)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_2_3":
			result := br.resolveHopBet(bet, roll, player, 2, 3)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_2_4":
			result := br.resolveHopBet(bet, roll, player, 2, 4)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_2_5":
			result := br.resolveHopBet(bet, roll, player, 2, 5)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_2_6":
			result := br.resolveHopBet(bet, roll, player, 2, 6)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_3_4":
			result := br.resolveHopBet(bet, roll, player, 3, 4)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_3_5":
			result := br.resolveHopBet(bet, roll, player, 3, 5)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_3_6":
			result := br.resolveHopBet(bet, roll, player, 3, 6)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_4_5":
			result := br.resolveHopBet(bet, roll, player, 4, 5)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_4_6":
			result := br.resolveHopBet(bet, roll, player, 4, 6)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		case "HOP_5_6":
			result := br.resolveHopBet(bet, roll, player, 5, 6)
			// Only persist if bet didn't win or lose (result is empty)
			if result == "" {
				br.persistAlwaysWorkingBet(bet)
			}
			return result
		default:
			return fmt.Sprintf("Unknown hop bet type: %s", betDef.Name)
		}
	default:
		// For any other bet types, use the generic resolver
		result := br.resolveGenericBet(bet, betDef, roll, player)
		// Only persist if bet didn't win or lose (result is empty)
		if result == "" {
			br.persistAlwaysWorkingBet(bet)
		}
		return result
	}
}

func (br *BetResolution) resolveLineBet(bet *Bet, betDef CanonicalBetDefinition, roll *Roll, player *Player) string {
	switch betDef.Name {
	case "PASS_LINE":
		return br.resolvePassLine(bet, roll, player)
	case "DONT_PASS":
		return br.resolveDontPass(bet, roll, player)
	default:
		return fmt.Sprintf("Unknown line bet type: %s", betDef.Name)
	}
}

func (br *BetResolution) resolveComeBet(bet *Bet, betDef CanonicalBetDefinition, roll *Roll, player *Player) string {
	switch betDef.Name {
	case "COME":
		return br.resolveCome(bet, roll, player)
	case "DONT_COME":
		return br.resolveDontCome(bet, roll, player)
	default:
		return fmt.Sprintf("Unknown come bet type: %s", betDef.Name)
	}
}

func (br *BetResolution) resolveOddsBet(bet *Bet, betDef CanonicalBetDefinition, roll *Roll, player *Player) string {
	switch betDef.Name {
	case "PASS_ODDS":
		return br.resolvePassOdds(bet, roll, player)
	case "DONT_PASS_ODDS":
		return br.resolveDontPassOdds(bet, roll, player)
	case "COME_ODDS":
		return br.resolveComeOdds(bet, roll, player)
	case "DONT_COME_ODDS":
		return br.resolveDontComeOdds(bet, roll, player)
	default:
		return fmt.Sprintf("Unknown odds bet type: %s", betDef.Name)
	}
}

func (br *BetResolution) resolvePlaceStyleBet(bet *Bet, betDef CanonicalBetDefinition, roll *Roll, player *Player) string {
	// Use the bet type key instead of the Name field for routing
	switch bet.Type {
	case "PLACE_NUMBERS":
		return br.resolvePlaceNumbers(bet, roll, player)
	case "PLACE_INSIDE":
		return br.resolvePlaceInside(bet, roll, player)
	case "PLACE_OUTSIDE":
		return br.resolvePlaceOutside(bet, roll, player)
	case "PLACE_4":
		return br.resolvePlaceBet(bet, roll, player, 4, 9, 5)
	case "PLACE_5":
		return br.resolvePlaceBet(bet, roll, player, 5, 7, 5)
	case "PLACE_6":
		return br.resolvePlaceBet(bet, roll, player, 6, 7, 6)
	case "PLACE_8":
		return br.resolvePlaceBet(bet, roll, player, 8, 7, 6)
	case "PLACE_9":
		return br.resolvePlaceBet(bet, roll, player, 9, 7, 5)
	case "PLACE_10":
		return br.resolvePlaceBet(bet, roll, player, 10, 9, 5)
	case "BUY_4":
		return br.resolveBuyBet(bet, roll, player, 4)
	case "BUY_5":
		return br.resolveBuyBet(bet, roll, player, 5)
	case "BUY_6":
		return br.resolveBuyBet(bet, roll, player, 6)
	case "BUY_8":
		return br.resolveBuyBet(bet, roll, player, 8)
	case "BUY_9":
		return br.resolveBuyBet(bet, roll, player, 9)
	case "BUY_10":
		return br.resolveBuyBet(bet, roll, player, 10)
	case "LAY_4":
		return br.resolveLayBet(bet, roll, player, 4)
	case "LAY_5":
		return br.resolveLayBet(bet, roll, player, 5)
	case "LAY_6":
		return br.resolveLayBet(bet, roll, player, 6)
	case "LAY_8":
		return br.resolveLayBet(bet, roll, player, 8)
	case "LAY_9":
		return br.resolveLayBet(bet, roll, player, 9)
	case "LAY_10":
		return br.resolveLayBet(bet, roll, player, 10)
	case "PLACE_TO_LOSE_4":
		return br.resolvePlaceToLoseBet(bet, roll, player, 4)
	case "PLACE_TO_LOSE_5":
		return br.resolvePlaceToLoseBet(bet, roll, player, 5)
	case "PLACE_TO_LOSE_6":
		return br.resolvePlaceToLoseBet(bet, roll, player, 6)
	case "PLACE_TO_LOSE_8":
		return br.resolvePlaceToLoseBet(bet, roll, player, 8)
	case "PLACE_TO_LOSE_9":
		return br.resolvePlaceToLoseBet(bet, roll, player, 9)
	case "PLACE_TO_LOSE_10":
		return br.resolvePlaceToLoseBet(bet, roll, player, 10)
	default:
		return fmt.Sprintf("Unknown place style bet type: %s", bet.Type)
	}
}

func (br *BetResolution) resolveHardWayBet(bet *Bet, betDef CanonicalBetDefinition, roll *Roll, player *Player) string {
	switch bet.Type {
	case "HARD_4":
		return br.resolveHardWay(bet, roll, player, 4, 7, 1)
	case "HARD_6":
		return br.resolveHardWay(bet, roll, player, 6, 9, 1)
	case "HARD_8":
		return br.resolveHardWay(bet, roll, player, 8, 9, 1)
	case "HARD_10":
		return br.resolveHardWay(bet, roll, player, 10, 7, 1)
	case "ALL_HARDWAYS":
		return br.resolveAllHardways(bet, roll, player)
	default:
		return fmt.Sprintf("Unknown hard way bet type: %s", bet.Type)
	}
}

func (br *BetResolution) resolveBigBetWithDef(bet *Bet, betDef CanonicalBetDefinition, roll *Roll, player *Player) string {
	switch betDef.Name {
	case "Big 6":
		return br.resolveBigBet(bet, roll, player, 6)
	case "Big 8":
		return br.resolveBigBet(bet, roll, player, 8)
	default:
		return fmt.Sprintf("Unknown big bet type: %s", betDef.Name)
	}
}

func (br *BetResolution) resolveBigBet(bet *Bet, roll *Roll, player *Player, number int) string {
	switch roll.Total {
	case number:
		winnings := bet.Amount
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Big %d wins $%.2f", number, winnings)
	case 7:
		bet.Working = false
		return fmt.Sprintf("💥 Big %d loses $%.2f", number, bet.Amount)
	default:
		return ""
	}
}

func (br *BetResolution) resolveGenericBet(bet *Bet, betDef CanonicalBetDefinition, roll *Roll, player *Player) string {
	switch betDef.Name {
	case "FIELD":
		return br.resolveField(bet, roll, player)
	case "ANY_SEVEN":
		return br.resolveAnySeven(bet, roll, player)
	case "ANY_CRAPS":
		return br.resolveAnyCraps(bet, roll, player)
	case "ELEVEN":
		return br.resolveEleven(bet, roll, player)
	case "ACE_DEUCE":
		return br.resolveAceDeuce(bet, roll, player)
	case "ACES":
		return br.resolveAces(bet, roll, player)
	case "BOXCARS":
		return br.resolveBoxcars(bet, roll, player)
	case "PLACE_4":
		return br.resolvePlaceBet(bet, roll, player, 4, 9, 5)
	case "PLACE_5":
		return br.resolvePlaceBet(bet, roll, player, 5, 7, 5)
	case "PLACE_6":
		return br.resolvePlaceBet(bet, roll, player, 6, 7, 6)
	case "PLACE_8":
		return br.resolvePlaceBet(bet, roll, player, 8, 7, 6)
	case "PLACE_9":
		return br.resolvePlaceBet(bet, roll, player, 9, 7, 5)
	case "PLACE_10":
		return br.resolvePlaceBet(bet, roll, player, 10, 9, 5)
	case "BUY_4":
		return br.resolveBuyBet(bet, roll, player, 4)
	case "BUY_5":
		return br.resolveBuyBet(bet, roll, player, 5)
	case "BUY_6":
		return br.resolveBuyBet(bet, roll, player, 6)
	case "BUY_8":
		return br.resolveBuyBet(bet, roll, player, 8)
	case "BUY_9":
		return br.resolveBuyBet(bet, roll, player, 9)
	case "BUY_10":
		return br.resolveBuyBet(bet, roll, player, 10)
	case "LAY_4":
		return br.resolveLayBet(bet, roll, player, 4)
	case "LAY_5":
		return br.resolveLayBet(bet, roll, player, 5)
	case "LAY_6":
		return br.resolveLayBet(bet, roll, player, 6)
	case "LAY_8":
		return br.resolveLayBet(bet, roll, player, 8)
	case "LAY_9":
		return br.resolveLayBet(bet, roll, player, 9)
	case "LAY_10":
		return br.resolveLayBet(bet, roll, player, 10)
	case "PLACE_TO_LOSE_4":
		return br.resolvePlaceToLoseBet(bet, roll, player, 4)
	case "PLACE_TO_LOSE_5":
		return br.resolvePlaceToLoseBet(bet, roll, player, 5)
	case "PLACE_TO_LOSE_6":
		return br.resolvePlaceToLoseBet(bet, roll, player, 6)
	case "PLACE_TO_LOSE_8":
		return br.resolvePlaceToLoseBet(bet, roll, player, 8)
	case "PLACE_TO_LOSE_9":
		return br.resolvePlaceToLoseBet(bet, roll, player, 9)
	case "PLACE_TO_LOSE_10":
		return br.resolvePlaceToLoseBet(bet, roll, player, 10)
	case "HORN_HIGH_2":
		return br.resolveHornHighBet(bet, roll, player, 2)
	case "HORN_HIGH_3":
		return br.resolveHornHighBet(bet, roll, player, 3)
	case "HORN_HIGH_11":
		return br.resolveHornHighBet(bet, roll, player, 11)
	case "HORN_HIGH_12":
		return br.resolveHornHighBet(bet, roll, player, 12)
	case "HOP_1_2":
		return br.resolveHopBet(bet, roll, player, 1, 2)
	case "HOP_1_3":
		return br.resolveHopBet(bet, roll, player, 1, 3)
	case "HOP_1_4":
		return br.resolveHopBet(bet, roll, player, 1, 4)
	case "HOP_1_5":
		return br.resolveHopBet(bet, roll, player, 1, 5)
	case "HOP_1_6":
		return br.resolveHopBet(bet, roll, player, 1, 6)
	case "HOP_2_3":
		return br.resolveHopBet(bet, roll, player, 2, 3)
	case "HOP_2_4":
		return br.resolveHopBet(bet, roll, player, 2, 4)
	case "HOP_2_5":
		return br.resolveHopBet(bet, roll, player, 2, 5)
	case "HOP_2_6":
		return br.resolveHopBet(bet, roll, player, 2, 6)
	case "HOP_3_4":
		return br.resolveHopBet(bet, roll, player, 3, 4)
	case "HOP_3_5":
		return br.resolveHopBet(bet, roll, player, 3, 5)
	case "HOP_3_6":
		return br.resolveHopBet(bet, roll, player, 3, 6)
	case "HOP_4_5":
		return br.resolveHopBet(bet, roll, player, 4, 5)
	case "HOP_4_6":
		return br.resolveHopBet(bet, roll, player, 4, 6)
	case "HOP_5_6":
		return br.resolveHopBet(bet, roll, player, 5, 6)
	case "WORLD":
		return br.resolveWorldBet(bet, roll, player)
	case "C_AND_E":
		return br.resolveCAndE(bet, roll, player)
	case "HORN":
		return br.resolveHornBet(bet, roll, player)
	case "HOP_HARD_6":
		return br.resolveHopHard6(bet, roll, player)
	case "HOP_EASY_8":
		return br.resolveHopEasy8(bet, roll, player)
	case "HOP":
		return br.resolveHopBet(bet, roll, player, 1, 1) // Default to 1-1 for generic HOP
	default:
		return fmt.Sprintf("Unknown generic bet type: %s", betDef.Name)
	}
}

func (br *BetResolution) resolvePassLine(bet *Bet, roll *Roll, player *Player) string {
	switch br.table.State {
	case StateComeOut:
		switch roll.Total {
		case 7, 11:
			winnings := bet.Amount
			player.Bankroll += winnings
			bet.Working = false
			return fmt.Sprintf("🎉 Pass line wins $%.2f (Natural)", winnings)
		case 2, 3, 12:
			bet.Working = false
			return fmt.Sprintf("💥 Pass line loses $%.2f (Craps)", bet.Amount)
		}
		return ""
	case StatePoint:
		switch roll.Total {
		case 7:
			bet.Working = false
			return fmt.Sprintf("💥 Pass line loses $%.2f (Seven out)", bet.Amount)
		default:
			pointNumber, err := PointToNumber(br.table.Point)
			if err != nil {
				return fmt.Sprintf("Error getting point number: %v", err)
			}
			if roll.Total == pointNumber {
				winnings := bet.Amount
				player.Bankroll += winnings
				bet.Working = false
				return fmt.Sprintf("🎉 Pass line wins $%.2f (Point made)", winnings)
			}
		}
		return ""
	}
	return ""
}

func (br *BetResolution) resolveDontPass(bet *Bet, roll *Roll, player *Player) string {
	switch br.table.State {
	case StateComeOut:
		switch roll.Total {
		case 7, 11:
			bet.Working = false
			return fmt.Sprintf("💥 Don't pass loses $%.2f (Natural)", bet.Amount)
		case 2, 3:
			winnings := bet.Amount
			player.Bankroll += winnings
			bet.Working = false
			return fmt.Sprintf("🎉 Don't pass wins $%.2f (Craps)", winnings)
		case 12:
			bet.Working = false
			player.Bankroll += bet.Amount
			return fmt.Sprintf("🤝 Don't pass push $%.2f (12)", bet.Amount)
		}
		return ""
	case StatePoint:
		switch roll.Total {
		case 7:
			winnings := bet.Amount
			player.Bankroll += winnings
			bet.Working = false
			return fmt.Sprintf("🎉 Don't pass wins $%.2f (Seven out)", winnings)
		default:
			pointNumber, err := PointToNumber(br.table.Point)
			if err != nil {
				return fmt.Sprintf("Error getting point number: %v", err)
			}
			if roll.Total == pointNumber {
				bet.Working = false
				return fmt.Sprintf("💥 Don't pass loses $%.2f (Point made)", bet.Amount)
			}
		}
		return ""
	}
	return ""
}

func (br *BetResolution) resolveField(bet *Bet, roll *Roll, player *Player) string {
	winningNumbers := map[int]float64{
		2:  2.0, // 2:1
		3:  1.0, // 1:1
		4:  1.0, // 1:1
		9:  1.0, // 1:1
		10: 1.0, // 1:1
		11: 1.0, // 1:1
		12: 2.0, // 2:1
	}

	if multiplier, wins := winningNumbers[roll.Total]; wins {
		winnings := bet.Amount * multiplier
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Field wins $%.2f (%d)", winnings, roll.Total)
	} else {
		bet.Working = false
		return fmt.Sprintf("💥 Field loses $%.2f (%d)", bet.Amount, roll.Total)
	}
}

func (br *BetResolution) resolveAnySeven(bet *Bet, roll *Roll, player *Player) string {
	switch roll.Total {
	case 7:
		winnings := bet.Amount * 4 // 4:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Any seven wins $%.2f", winnings)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 Any seven loses $%.2f", bet.Amount)
	}
}

func (br *BetResolution) resolveAnyCraps(bet *Bet, roll *Roll, player *Player) string {
	switch roll.Total {
	case 2, 3, 12:
		winnings := bet.Amount * 7 // 7:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Any craps wins $%.2f (%d)", winnings, roll.Total)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 Any craps loses $%.2f", bet.Amount)
	}
}

func (br *BetResolution) resolveEleven(bet *Bet, roll *Roll, player *Player) string {
	switch roll.Total {
	case 11:
		winnings := bet.Amount * 15 // 15:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Eleven wins $%.2f", winnings)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 Eleven loses $%.2f", bet.Amount)
	}
}

func (br *BetResolution) resolveAceDeuce(bet *Bet, roll *Roll, player *Player) string {
	switch roll.Total {
	case 3:
		winnings := bet.Amount * 15 // 15:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Ace-deuce wins $%.2f", winnings)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 Ace-deuce loses $%.2f", bet.Amount)
	}
}

func (br *BetResolution) resolveAces(bet *Bet, roll *Roll, player *Player) string {
	switch roll.Total {
	case 2:
		winnings := bet.Amount * 30 // 30:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Aces wins $%.2f", winnings)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 Aces loses $%.2f", bet.Amount)
	}
}

func (br *BetResolution) resolveBoxcars(bet *Bet, roll *Roll, player *Player) string {
	switch roll.Total {
	case 12:
		winnings := bet.Amount * 30 // 30:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Boxcars wins $%.2f", winnings)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 Boxcars loses $%.2f", bet.Amount)
	}
}

func (br *BetResolution) resolvePlaceBet(bet *Bet, roll *Roll, player *Player, number, numerator, denominator int) string {
	switch roll.Total {
	case number:
		winnings := (bet.Amount * float64(numerator)) / float64(denominator)
		player.Bankroll += winnings
		// Place bets continue working after they win - don't set Working = false
		return fmt.Sprintf("🎉 Place %d wins $%.2f", number, winnings)
	case 7:
		bet.Working = false
		return fmt.Sprintf("💥 Place %d loses $%.2f", number, bet.Amount)
	default:
		return ""
	}
}

func (br *BetResolution) resolveHardWay(bet *Bet, roll *Roll, player *Player, number, numerator, denominator int) string {
	switch {
	case roll.Total == number && roll.IsHard:
		winnings := (bet.Amount * float64(numerator)) / float64(denominator)
		player.Bankroll += winnings
		bet.Working = false
		// Don't call removeOneRollBet - hard way bets are not one-roll bets
		// They will be removed by cleanupResolvedBets() when Working = false
		return fmt.Sprintf("🎉 Hard %d wins $%.2f", number, winnings)
	case roll.Total == number && !roll.IsHard:
		bet.Working = false
		// Don't call removeOneRollBet - hard way bets are not one-roll bets
		// They will be removed by cleanupResolvedBets() when Working = false
		return fmt.Sprintf("💥 Hard %d loses $%.2f (Easy way)", number, bet.Amount)
	case roll.Total == 7:
		bet.Working = false
		// Don't call removeOneRollBet - hard way bets are not one-roll bets
		// They will be removed by cleanupResolvedBets() when Working = false
		return fmt.Sprintf("💥 Hard %d loses $%.2f (Seven out)", number, bet.Amount)
	default:
		return ""
	}
}

func (br *BetResolution) resolveWorldBet(bet *Bet, roll *Roll, player *Player) string {
	switch roll.Total {
	case 7:
		// World bet is split: 16.0 on any 7 (4:1), 4.0 on any craps (7:1)
		// 4:1 payout for any 7 portion
		winnings := 16.0 * 4.0 // 64.0
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 World wins $%.2f", winnings)
	case 2, 3, 12:
		// World bet is split: 16.0 on any 7 (4:1), 4.0 on any craps (7:1)
		// 7:1 payout for any craps portion
		winnings := 4.0 * 7.0 // 28.0
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 World wins $%.2f", winnings)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 World bet loses $%.2f", bet.Amount)
	}
}

func (br *BetResolution) resolveCAndE(bet *Bet, roll *Roll, player *Player) string {
	switch roll.Total {
	case 11:
		// C_AND_E bet is split: 10.0 on eleven (15:1), 10.0 on any craps (7:1)
		// 15:1 payout for eleven portion
		winnings := 10.0 * 15.0 // 150.0
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 C and E wins $%.2f", winnings)
	case 2, 3, 12:
		// C_AND_E bet is split: 10.0 on eleven (15:1), 10.0 on any craps (7:1)
		// 7:1 payout for any craps portion
		winnings := 10.0 * 7.0 // 70.0
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 C and E wins $%.2f", winnings)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 C & E loses $%.2f", bet.Amount)
	}
}

func (br *BetResolution) resolveHornBet(bet *Bet, roll *Roll, player *Player) string {
	betPerNumber := bet.Amount / 4.0

	switch roll.Total {
	case 2:
		winnings := betPerNumber * 30 // 30:1 for aces
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Horn wins $%.2f (Aces)", winnings)
	case 3:
		winnings := betPerNumber * 15 // 15:1 for ace-deuce
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Horn wins $%.2f (Ace-deuce)", winnings)
	case 11:
		winnings := betPerNumber * 15 // 15:1 for eleven
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Horn wins $%.2f (Eleven)", winnings)
	case 12:
		winnings := betPerNumber * 30 // 30:1 for boxcars
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Horn wins $%.2f (Boxcars)", winnings)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 Horn loses $%.2f", bet.Amount)
	}
}

func (br *BetResolution) resolveCome(bet *Bet, roll *Roll, player *Player) string {
	switch br.table.State {
	case StateComeOut:
		switch roll.Total {
		case 7, 11:
			winnings := bet.Amount
			player.Bankroll += winnings
			bet.Working = false
			return fmt.Sprintf("🎉 Come bet wins $%.2f (Natural)", winnings)
		case 2, 3, 12:
			bet.Working = false
			return fmt.Sprintf("💥 Come bet loses $%.2f (Craps)", bet.Amount)
		}
		bet.Working = true
		return ""
	case StatePoint:
		switch roll.Total {
		case 7:
			bet.Working = false
			return fmt.Sprintf("💥 Come bet loses $%.2f (Seven out)", bet.Amount)
		default:
			pointNumber, err := PointToNumber(br.table.Point)
			if err != nil {
				return fmt.Sprintf("Error getting point number: %v", err)
			}
			if roll.Total == pointNumber {
				winnings := bet.Amount
				player.Bankroll += winnings
				bet.Working = false
				return fmt.Sprintf("🎉 Come bet wins $%.2f (Point made)", winnings)
			}
		}
		return ""
	}
	return ""
}

func (br *BetResolution) resolveDontCome(bet *Bet, roll *Roll, player *Player) string {
	switch br.table.State {
	case StateComeOut:
		switch roll.Total {
		case 7, 11:
			bet.Working = false
			return fmt.Sprintf("💥 Don't come loses $%.2f (Natural)", bet.Amount)
		case 2, 3:
			winnings := bet.Amount
			player.Bankroll += winnings
			bet.Working = false
			return fmt.Sprintf("🎉 Don't come wins $%.2f (Craps)", winnings)
		case 12:
			bet.Working = false
			player.Bankroll += bet.Amount
			return fmt.Sprintf("🤝 Don't come push $%.2f (12)", bet.Amount)
		}
		bet.Working = true
		return ""
	case StatePoint:
		switch roll.Total {
		case 7:
			winnings := bet.Amount
			player.Bankroll += winnings
			bet.Working = false
			return fmt.Sprintf("🎉 Don't come wins $%.2f (Seven out)", winnings)
		default:
			pointNumber, err := PointToNumber(br.table.Point)
			if err != nil {
				return fmt.Sprintf("Error getting point number: %v", err)
			}
			if roll.Total == pointNumber {
				bet.Working = false
				return fmt.Sprintf("💥 Don't come loses $%.2f (Point made)", bet.Amount)
			}
		}
		return ""
	}
	return ""
}

func (br *BetResolution) resolvePlaceNumbers(bet *Bet, roll *Roll, player *Player) string {
	if len(bet.Numbers) == 0 {
		return "Invalid place numbers bet: no numbers specified"
	}

	if roll.Total == 7 {
		bet.Working = false
		return fmt.Sprintf("💥 Place numbers lose $%.2f (Seven out)", bet.Amount)
	}

	odds := map[int][2]float64{
		4:  {9, 5},
		5:  {7, 5},
		6:  {7, 6},
		8:  {7, 6},
		9:  {7, 5},
		10: {9, 5},
	}

	betPerNumber := bet.Amount / float64(len(bet.Numbers))
	for _, num := range bet.Numbers {
		if roll.Total == num {
			payout, ok := odds[num]
			if !ok {
				continue
			}
			winnings := betPerNumber * payout[0] / payout[1]
			player.Bankroll += winnings
			bet.Working = false
			return fmt.Sprintf("🎉 Place numbers win $%.2f (%d)", winnings, num)
		}
	}

	return ""
}

func (br *BetResolution) resolvePlaceInside(bet *Bet, roll *Roll, player *Player) string {
	betPerNumber := bet.Amount / 4.0

	switch roll.Total {
	case 5, 9:
		winnings := betPerNumber * 7.0 / 5.0 // 7:5 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Place inside wins $%.2f (%d)", winnings, roll.Total)
	case 6, 8:
		winnings := betPerNumber * 7.0 / 6.0 // 7:6 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Place inside wins $%.2f (%d)", winnings, roll.Total)
	case 7:
		bet.Working = false
		return fmt.Sprintf("💥 Place inside loses $%.2f (Seven out)", bet.Amount)
	}
	return ""
}

func (br *BetResolution) resolvePlaceOutside(bet *Bet, roll *Roll, player *Player) string {
	betPerNumber := bet.Amount / 4.0

	switch roll.Total {
	case 4, 10:
		winnings := betPerNumber * 9.0 / 5.0 // 9:5 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Place outside wins $%.2f (%d)", winnings, roll.Total)
	case 5, 9:
		winnings := betPerNumber * 7.0 / 5.0 // 7:5 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Place outside wins $%.2f (%d)", winnings, roll.Total)
	case 7:
		bet.Working = false
		return fmt.Sprintf("💥 Place outside loses $%.2f (Seven out)", bet.Amount)
	}
	return ""
}

func (br *BetResolution) resolveAllHardways(bet *Bet, roll *Roll, player *Player) string {
	betPerNumber := bet.Amount / 4.0

	switch {
	case roll.Total == 4 && roll.IsHard:
		winnings := betPerNumber * 7.0 // 7:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 All hardways win $%.2f (Hard 4)", winnings)
	case roll.Total == 6 && roll.IsHard:
		winnings := betPerNumber * 9.0 // 9:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 All hardways win $%.2f (Hard 6)", winnings)
	case roll.Total == 8 && roll.IsHard:
		winnings := betPerNumber * 9.0 // 9:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 All hardways win $%.2f (Hard 8)", winnings)
	case roll.Total == 10 && roll.IsHard:
		winnings := betPerNumber * 7.0 // 7:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 All hardways win $%.2f (Hard 10)", winnings)
	case roll.Total == 7:
		bet.Working = false
		return fmt.Sprintf("💥 All hardways lose $%.2f (Seven out)", bet.Amount)
	}
	return ""
}

// resolveHopHard6 resolves hop hard 6 bet (3-3)
func (br *BetResolution) resolveHopHard6(bet *Bet, roll *Roll, player *Player) string {
	switch {
	case roll.Total == 6 && roll.IsHard:
		winnings := bet.Amount * 30.0 // 30:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Hop hard 6 wins $%.2f", winnings)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 Hop hard 6 loses $%.2f", bet.Amount)
	}
}

// resolveHopEasy8 resolves hop easy 8 bet (any 8 except 4-4)
func (br *BetResolution) resolveHopEasy8(bet *Bet, roll *Roll, player *Player) string {
	switch {
	case roll.Total == 8 && !roll.IsHard:
		winnings := bet.Amount * 15.0 // 15:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Hop easy 8 wins $%.2f", winnings)
	default:
		bet.Working = false
		return fmt.Sprintf("💥 Hop easy 8 loses $%.2f", bet.Amount)
	}
}

// resolveHopBet resolves hop bets (one-roll)
func (br *BetResolution) resolveHopBet(bet *Bet, roll *Roll, player *Player, die1, die2 int) string {
	// Check if roll matches the exact hop combination
	if (roll.Die1 == die1 && roll.Die2 == die2) || (roll.Die1 == die2 && roll.Die2 == die1) {
		// Exact combination wins
		winnings := bet.Amount * 15.0 // 15:1 payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Hop %d-%d wins $%.2f", die1, die2, winnings)
	} else {
		// Any other roll loses
		bet.Working = false
		return fmt.Sprintf("💥 Hop %d-%d loses $%.2f", die1, die2, bet.Amount)
	}
}

// resolvePassOdds resolves pass odds bets
func (br *BetResolution) resolvePassOdds(bet *Bet, roll *Roll, player *Player) string {
	if br.table.State != StatePoint {
		return ("Pass odds not valid in current state")
	}

	pointNumber, err := PointToNumber(br.table.Point)
	if err != nil {
		return fmt.Sprintf("Error getting point number: %v", err)
	}

	switch roll.Total {
	case pointNumber:
		// Point made - pass odds win
		var odds float64
		switch pointNumber {
		case 4, 10:
			odds = 2.0 // 2:1
		case 5, 9:
			odds = 1.5 // 3:2
		case 6, 8:
			odds = 1.2 // 6:5
		default:
			return fmt.Sprintf("Invalid point for pass odds: %d", pointNumber)
		}

		winnings := bet.Amount * odds
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Pass odds win $%.2f (point made)", winnings)
	case 7:
		// Seven out - pass odds lose
		bet.Working = false
		return fmt.Sprintf("💥 Pass odds lose $%.2f (seven out)", bet.Amount)
	default:
		// Other numbers - bet stays working
		return ""
	}
}

// resolveDontPassOdds resolves don't pass odds bets
func (br *BetResolution) resolveDontPassOdds(bet *Bet, roll *Roll, player *Player) string {
	if br.table.State != StatePoint {
		return ("Don't pass odds not valid in current state")
	}

	pointNumber, err := PointToNumber(br.table.Point)
	if err != nil {
		return fmt.Sprintf("Error getting point number: %v", err)
	}

	switch roll.Total {
	case 7:
		// Seven out - don't pass odds win
		var odds float64
		switch pointNumber {
		case 4, 10:
			odds = 0.5 // 1:2
		case 5, 9:
			odds = 0.667 // 2:3
		case 6, 8:
			odds = 0.833 // 5:6
		default:
			return fmt.Sprintf("Invalid point for don't pass odds: %d", pointNumber)
		}

		winnings := bet.Amount * odds
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Don't pass odds win $%.2f (seven out)", winnings)
	case pointNumber:
		// Point made - don't pass odds lose
		bet.Working = false
		return fmt.Sprintf("💥 Don't pass odds lose $%.2f (point made)", bet.Amount)
	default:
		// Other numbers - bet stays working
		return ""
	}
}

// resolveComeOdds resolves come odds bets
func (br *BetResolution) resolveComeOdds(bet *Bet, roll *Roll, player *Player) string {
	// Come odds work the same as pass odds but for come bets
	// This is a simplified implementation - in a full system, come odds
	// would be tied to specific come bet points
	return br.resolvePassOdds(bet, roll, player)
}

// resolveDontComeOdds resolves don't come odds bets
func (br *BetResolution) resolveDontComeOdds(bet *Bet, roll *Roll, player *Player) string {
	// Don't come odds work the same as don't pass odds but for don't come bets
	// This is a simplified implementation - in a full system, don't come odds
	// would be tied to specific don't come bet points
	return br.resolveDontPassOdds(bet, roll, player)
}

// resolveBuyBet resolves buy bets with commission
func (br *BetResolution) resolveBuyBet(bet *Bet, roll *Roll, player *Player, number int) string {
	switch roll.Total {
	case number:
		// Buy bet wins - calculate true odds payout
		var odds float64
		switch number {
		case 4, 10:
			odds = 2.0 // 2:1
		case 5, 9:
			odds = 1.5 // 3:2
		case 6, 8:
			odds = 1.2 // 6:5
		default:
			return fmt.Sprintf("Invalid buy bet number: %d", number)
		}

		// Calculate winnings with commission using the new function
		winnings := br.calculateWinningsWithCommission(bet.Amount, odds, BuyLayCommissionRate)
		commissionMsg := br.formatCommissionMessage(bet.Amount, odds, BuyLayCommissionRate)
		player.Bankroll += winnings
		// Buy bets continue working after they win - don't set Working = false
		return fmt.Sprintf("🎉 Buy %d wins $%.2f%s", number, winnings, commissionMsg)
	case 7:
		// Seven out - buy bet loses
		bet.Working = false
		return fmt.Sprintf("💥 Buy %d loses $%.2f", number, bet.Amount)
	default:
		// Other numbers - bet stays working
		return ""
	}
}

// resolveLayBet resolves lay bets with commission
func (br *BetResolution) resolveLayBet(bet *Bet, roll *Roll, player *Player, number int) string {
	switch roll.Total {
	case 7:
		// Seven out - lay bet wins
		var odds float64
		switch number {
		case 4, 10:
			odds = 0.5 // 1:2
		case 5, 9:
			odds = 0.667 // 2:3
		case 6, 8:
			odds = 0.833 // 5:6
		default:
			return fmt.Sprintf("Invalid lay bet number: %d", number)
		}

		// Calculate winnings with commission using the new function
		winnings := br.calculateWinningsWithCommission(bet.Amount, odds, BuyLayCommissionRate)
		commissionMsg := br.formatCommissionMessage(bet.Amount, odds, BuyLayCommissionRate)
		player.Bankroll += winnings
		// Lay bets continue working after they win - don't set Working = false
		return fmt.Sprintf("🎉 Lay %d wins $%.2f%s", number, winnings, commissionMsg)
	case number:
		// Number rolled - lay bet loses
		bet.Working = false
		return fmt.Sprintf("💥 Lay %d loses $%.2f", number, bet.Amount)
	default:
		// Other numbers - bet stays working
		return ""
	}
}

// resolvePlaceToLoseBet resolves place-to-lose bets (no commission)
func (br *BetResolution) resolvePlaceToLoseBet(bet *Bet, roll *Roll, player *Player, number int) string {
	switch roll.Total {
	case 7:
		// Seven out - place-to-lose bet wins
		var odds float64
		switch number {
		case 4, 10:
			odds = 0.5 // 1:2
		case 5, 9:
			odds = 0.667 // 2:3
		case 6, 8:
			odds = 0.833 // 5:6
		default:
			return fmt.Sprintf("Invalid place-to-lose bet number: %d", number)
		}

		// Calculate winnings without commission
		winnings := bet.Amount * odds
		player.Bankroll += winnings
		// Place-to-lose bets continue working after they win - don't set Working = false
		return fmt.Sprintf("🎉 Place-to-lose %d wins $%.2f", number, winnings)
	case number:
		// Number rolled - place-to-lose bet loses
		bet.Working = false
		return fmt.Sprintf("💥 Place-to-lose %d loses $%.2f", number, bet.Amount)
	default:
		// Other numbers - bet stays working
		return ""
	}
}

// resolveHornHighBet resolves horn high bets
func (br *BetResolution) resolveHornHighBet(bet *Bet, roll *Roll, player *Player, highNumber int) string {
	switch roll.Total {
	case highNumber:
		// High number wins with enhanced payout
		var payout float64
		switch highNumber {
		case 2, 12:
			payout = 6.75 // 27:4
		case 3, 11:
			payout = 15.0 // 15:1
		default:
			return fmt.Sprintf("Invalid horn high number: %d", highNumber)
		}

		winnings := bet.Amount * payout
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Horn High %d wins $%.2f", highNumber, winnings)
	case 2, 3, 11, 12:
		// Other horn numbers win with standard payout
		winnings := bet.Amount * 3.0 // 3:1
		player.Bankroll += winnings
		bet.Working = false
		return fmt.Sprintf("🎉 Horn High %d wins $%.2f (horn number)", highNumber, winnings)
	default:
		// Other numbers - bet loses
		bet.Working = false
		return fmt.Sprintf("💥 Horn High %d loses $%.2f", highNumber, bet.Amount)
	}
}

// Proper RNG for dice rolls
func rollDieSecure() int {
	n, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		// Fallback to time-based if crypto/rand fails
		return 1 + (int(time.Now().UnixNano()) % 6)
	}
	return int(n.Int64()) + 1
}

// ===== BET LIFECYCLE MANAGEMENT FUNCTIONS =====

// setBetWorking sets the working status of a bet
func (br *BetResolution) setBetWorking(bet *Bet, working bool) {
	if bet == nil {
		return
	}
	bet.Working = working
}

// removeOneRollBet removes a one-roll bet after resolution
func (br *BetResolution) removeOneRollBet(bet *Bet) {
	if bet == nil {
		return
	}

	// Find the player who owns this bet
	player, exists := br.table.Players[bet.Player]
	if !exists {
		return
	}

	// Remove bet from player's bet list
	for i, playerBet := range player.Bets {
		if playerBet.ID == bet.ID {
			// Remove bet from slice
			player.Bets = append(player.Bets[:i], player.Bets[i+1:]...)
			break
		}
	}
}

// persistAlwaysWorkingBet keeps an always-working bet active after resolution
func (br *BetResolution) persistAlwaysWorkingBet(bet *Bet) {
	if bet == nil {
		return
	}

	// Always-working bets stay active unless explicitly removed
	// This function ensures the bet remains in the player's bet list
	// and maintains its working status for future rolls
	bet.Working = true
}

// checkConditionalBetState determines if a conditional bet should be working based on game state
func (br *BetResolution) checkConditionalBetState(bet *Bet, gameState GameState) bool {
	if bet == nil {
		return false
	}

	// Look up bet definition to check conditional requirements
	betDef, exists := CanonicalBetDefinitions[bet.Type]
	if !exists {
		return false
	}

	// Special handling for hard way bets
	if betDef.Category == HardWayBets {
		// Hard way bets should be working unless they've been resolved (won or lost)
		// The Working field is set to false when they win or lose
		return bet.Working
	}

	// Special handling for place bets, buy bets, lay bets, and place-to-lose bets
	// These bets are off during come-out rolls and on during point phases
	if betDef.Category == PlaceBets || betDef.Category == BuyBets || betDef.Category == LayBets || betDef.Category == PlaceToLoseBets {
		// These bets are only working during point phase
		return gameState == StatePoint
	}

	// Check if bet requires specific game state
	if betDef.RequiresPoint && gameState != StatePoint {
		return false
	}

	if betDef.RequiresComeOut && gameState != StateComeOut {
		return false
	}

	// If no specific requirements, bet is working
	return true
}

// validateBetState validates the current state of a bet
func (br *BetResolution) validateBetState(bet *Bet) error {
	if bet == nil {
		return fmt.Errorf("bet is nil")
	}

	if bet.Amount <= 0 {
		return fmt.Errorf("bet amount must be positive: $%.2f", bet.Amount)
	}

	if bet.Player == "" {
		return fmt.Errorf("bet must have a player ID")
	}

	if bet.Type == "" {
		return fmt.Errorf("bet must have a type")
	}

	// Check if bet type exists in canonical definitions
	if _, exists := CanonicalBetDefinitions[bet.Type]; !exists {
		return fmt.Errorf("unknown bet type: %s", bet.Type)
	}

	// Check if player exists
	if _, exists := br.table.Players[bet.Player]; !exists {
		return fmt.Errorf("player %s not found for bet", bet.Player)
	}

	return nil
}

// cleanupResolvedBets removes resolved bets and updates working status
func (br *BetResolution) cleanupResolvedBets() {
	for _, player := range br.table.Players {
		var activeBets []*Bet

		for _, bet := range player.Bets {
			// Validate bet state
			if err := br.validateBetState(bet); err != nil {
				// Remove invalid bets
				continue
			}

			// Check if bet should remain active
			if bet.Working {
				activeBets = append(activeBets, bet)
			}
		}

		// Update player's bet list to only include active bets
		player.Bets = activeBets
	}
}

// updateBetWorkingStatus updates working status for all bets based on current game state
func (br *BetResolution) updateBetWorkingStatus() {
	for _, player := range br.table.Players {
		for _, bet := range player.Bets {
			// Get bet definition
			betDef, exists := CanonicalBetDefinitions[bet.Type]
			if !exists {
				continue
			}

			// Don't override bets that have been set to Working = false by resolution
			// This prevents resolved bets from being reactivated
			if !bet.Working {
				continue
			}

			// Update working status based on bet behavior and game state
			switch betDef.WorkingBehavior {
			case "ONE_ROLL":
				// One-roll bets are only working for one roll
				// They should be set to false after resolution
				// This is handled in resolveOneRollBet
			case "ALWAYS":
				// Always working bets stay active
				bet.Working = true
			case "CONDITIONAL":
				// Conditional bets depend on game state
				bet.Working = br.checkConditionalBetState(bet, br.table.State)
			default:
				// Default to working if behavior is not specified
				bet.Working = true
			}
		}
	}
}
