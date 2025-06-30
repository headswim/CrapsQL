package crapsql

import (
	"fmt"

	"github.com/headswim/CrapsQL/pkg/crapsgame"
)

// Table represents a craps table with game state and players
type Table = crapsgame.Table

// Player represents a player at the table
type Player = crapsgame.Player

// Bet represents a bet on the table
type Bet = crapsgame.Bet

// Roll represents a dice roll
type Roll = crapsgame.Roll

// NewTable creates a new craps table with specified limits
func NewTable(minBet, maxBet float64, maxOdds int) *Table {
	return crapsgame.NewTable(minBet, maxBet, maxOdds)
}

// ExecuteString is a convenience function to execute CrapsQL commands
func ExecuteString(input string, table *Table) ([]string, error) {
	interpreter := NewInterpreter(table)
	return interpreter.ExecuteString(input)
}

// RollDice rolls the dice and resolves all bets
func RollDice(table *Table) *Roll {
	return table.RollDice()
}

// AddPlayer adds a player to the table
func AddPlayer(table *Table, id, name string, bankroll float64) error {
	return table.AddPlayer(id, name, bankroll)
}

// GetPlayer returns a player by ID
func GetPlayer(table *Table, id string) (*Player, error) {
	return table.GetPlayer(id)
}

// PlaceBet places a bet on the table
func PlaceBet(table *Table, playerID, betType string, amount float64) (*Bet, error) {
	return table.PlaceBet(playerID, betType, amount, []int{})
}

// GetState returns the current game state
func GetState(table *Table) crapsgame.GameState {
	return table.GetState()
}

// GetPoint returns the current point
func GetPoint(table *Table) crapsgame.Point {
	return table.GetPoint()
}

// IsComeOut returns true if we're in come out phase
func IsComeOut(table *Table) bool {
	return table.IsComeOut()
}

// IsPoint returns true if we have a point established
func IsPoint(table *Table) bool {
	return table.IsPoint()
}

// GetShooter returns the current shooter
func GetShooter(table *Table) string {
	return table.GetShooter()
}

// RemovePlayer removes a player from the table
func RemovePlayer(table *Table, id string) error {
	return table.RemovePlayer(id)
}

// Example usage and documentation
func Example() {
	// Create a new craps table
	table := NewTable(5.0, 100.0, 3) // $5 min, $100 max, 3x odds

	// Add a player
	err := AddPlayer(table, "player1", "John", 1000.0)
	if err != nil {
		fmt.Printf("Error adding player: %v\n", err)
		return
	}

	// Create interpreter
	interpreter := NewInterpreter(table)

	// Execute CrapsQL commands
	results, err := interpreter.ExecuteString("PLACE $25 ON PASS_LINE;")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, result := range results {
		fmt.Println(result)
	}

	// Roll the dice
	roll := RollDice(table)
	fmt.Printf("Rolled: %d-%d = %d\n", roll.Die1, roll.Die2, roll.Total)

	// Check game state
	if IsComeOut(table) {
		fmt.Println("Come out roll")
	} else {
		fmt.Printf("Point is: %d\n", GetPoint(table))
	}
}
