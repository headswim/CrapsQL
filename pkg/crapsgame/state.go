package crapsgame

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// GameState represents the current state of the craps table
type GameState int

const (
	StateComeOut GameState = iota
	StatePoint
	StateSevenOut
)

// Point represents the current point number
type Point int

const (
	PointOff Point = iota
	Point4
	Point5
	Point6
	Point8
	Point9
	Point10
)

// Roll represents a dice roll
type Roll struct {
	Die1   int
	Die2   int
	Total  int
	IsHard bool // true if both dice show the same number
	Time   time.Time
}

// Bet represents a single bet on the table
type Bet struct {
	ID       string
	Type     string
	Amount   float64
	Player   string
	PlacedAt time.Time
	Working  bool    // true if bet is active for current roll
	Odds     float64 // for odds bets
	Numbers  []int   // for bets on specific numbers (e.g., place numbers)
}

// Player represents a player at the table
type Player struct {
	ID           string
	Name         string
	Bankroll     float64
	Bets         []*Bet
	MaxBet       float64
	MinBet       float64
	WinGoal      float64
	LossLimit    float64
	SessionStart time.Time
}

// ComeBet represents a come bet with its lifecycle
type ComeBet struct {
	ID        string
	PlayerID  string
	Amount    float64
	ComePoint Point   // Point number when established, PointOff when not established
	Working   bool    // true if bet is active for current roll
	Odds      float64 // odds bet amount (0 if no odds)
	PlacedAt  time.Time
}

// OddsBet represents an odds bet with its lifecycle
type OddsBet struct {
	ID          string
	PlayerID    string
	Amount      float64
	BaseBetType string // Type of base bet (PASS_LINE, DONT_PASS, COME, DONT_COME)
	Point       Point  // Point number when odds were placed
	Working     bool   // true if bet is active for current roll
	PlacedAt    time.Time
}

// Table represents the craps table
type Table struct {
	State       GameState
	Point       Point
	CurrentRoll *Roll
	Players     map[string]*Player
	Shooter     string // current shooter's ID
	MinBet      float64
	MaxBet      float64
	MaxOdds     int // maximum odds allowed (e.g., 3x, 5x)
	CreatedAt   time.Time
	LastRoll    time.Time
	BetResolver *BetResolution      // New bet resolution handler
	ComeBets    map[string]*ComeBet // Come bet state tracking
	OddsBets    map[string]*OddsBet // Odds bet state tracking
}

// NewTable creates a new craps table
func NewTable(minBet, maxBet float64, maxOdds int) *Table {
	table := &Table{
		State:     StateComeOut,
		Point:     PointOff,
		Players:   make(map[string]*Player),
		ComeBets:  make(map[string]*ComeBet),
		OddsBets:  make(map[string]*OddsBet),
		MinBet:    minBet,
		MaxBet:    maxBet,
		MaxOdds:   maxOdds,
		CreatedAt: time.Now(),
	}
	table.BetResolver = NewBetResolution(table)
	return table
}

// AddPlayer adds a player to the table
func (t *Table) AddPlayer(id, name string, bankroll float64) error {
	if _, exists := t.Players[id]; exists {
		return fmt.Errorf("player %s already exists", id)
	}

	t.Players[id] = &Player{
		ID:           id,
		Name:         name,
		Bankroll:     bankroll,
		Bets:         []*Bet{},
		MaxBet:       t.MaxBet,
		MinBet:       t.MinBet,
		SessionStart: time.Now(),
	}

	// Set first player as shooter if no shooter exists
	if t.Shooter == "" {
		t.Shooter = id
	}

	return nil
}

// RemovePlayer removes a player from the table
func (t *Table) RemovePlayer(id string) error {
	player, exists := t.Players[id]
	if !exists {
		return fmt.Errorf("player %s not found", id)
	}

	// Remove all player's bets
	for _, bet := range player.Bets {
		t.removeBet(bet.ID)
	}

	// Remove all player's come bets
	t.removePlayerComeBets(id)

	// Remove all player's odds bets
	t.removePlayerOddsBets(id)

	delete(t.Players, id)

	// If this was the shooter, assign new shooter
	if t.Shooter == id {
		t.assignNewShooter()
	}

	return nil
}

// removePlayerComeBets removes all come bets for a specific player
func (t *Table) removePlayerComeBets(playerID string) {
	var comeBetIDsToRemove []string

	// Find all come bets for this player
	for comeBetID, comeBet := range t.ComeBets {
		if comeBet.PlayerID == playerID {
			comeBetIDsToRemove = append(comeBetIDsToRemove, comeBetID)
		}
	}

	// Remove the come bets
	for _, comeBetID := range comeBetIDsToRemove {
		t.removeComeBet(comeBetID)
	}
}

// removePlayerOddsBets removes all odds bets for a specific player
func (t *Table) removePlayerOddsBets(playerID string) {
	var oddsBetIDsToRemove []string

	// Find all odds bets for this player
	for oddsBetID, oddsBet := range t.OddsBets {
		if oddsBet.PlayerID == playerID {
			oddsBetIDsToRemove = append(oddsBetIDsToRemove, oddsBetID)
		}
	}

	// Remove the odds bets
	for _, oddsBetID := range oddsBetIDsToRemove {
		t.removeOddsBet(oddsBetID)
	}
}

// assignNewShooter assigns a new shooter from available players
func (t *Table) assignNewShooter() {
	for id := range t.Players {
		t.Shooter = id
		return
	}
	t.Shooter = "" // no players left
}

// RollDice simulates a dice roll using secure RNG
func (t *Table) RollDice() *Roll {
	// Validate shooter before roll
	if err := t.validateShooter(t.Shooter); err != nil {
		fmt.Printf("Warning: Invalid shooter before roll: %v\n", err)
		// Assign new shooter if current one is invalid
		t.assignNewShooter()
	}

	// Validate table state before roll
	if err := t.validateTableState(); err != nil {
		fmt.Printf("Warning: Invalid table state before roll: %v\n", err)
	}

	roll := &Roll{
		Die1: rollDieSecure(),
		Die2: rollDieSecure(),
		Time: time.Now(),
	}
	roll.Total = roll.Die1 + roll.Die2
	roll.IsHard = roll.Die1 == roll.Die2

	t.CurrentRoll = roll
	t.LastRoll = roll.Time

	// Update game state based on roll
	t.updateGameState(roll)

	// Validate table state after roll
	if err := t.validateTableState(); err != nil {
		fmt.Printf("Warning: Invalid table state after roll: %v\n", err)
	}

	// Note: Bet resolution is handled by the caller (interpreter)
	// This prevents double resolution when ROLL DICE is called

	return roll
}

// updateGameState updates the game state based on the current roll
func (t *Table) updateGameState(roll *Roll) {
	switch t.State {
	case StateComeOut:
		switch roll.Total {
		case 7, 11:
			// Natural - pass line wins, don't pass loses
			t.natural(roll)
		case 2, 3, 12:
			// Craps - pass line loses, don't pass wins (except 12)
			t.craps(roll)
		default:
			// Point established
			t.establishPoint(roll)
		}
	case StatePoint:
		if roll.Total == 7 {
			// Seven out - pass line loses, don't pass wins
			t.sevenOut(roll)
		} else {
			pointNumber, err := PointToNumber(t.Point)
			if err != nil {
				fmt.Printf("Error getting point number: %v\n", err)
				return
			}
			if roll.Total == pointNumber {
				// Point made - pass line wins, don't pass loses
				t.resolvePoint(roll)
			}
		}
		// Other numbers don't change the point
	}
}

// establishPoint establishes a point when a point number is rolled during come out
func (t *Table) establishPoint(roll *Roll) {
	// Validate state transition
	if err := t.validateStateTransition(StateComeOut, StatePoint, roll); err != nil {
		fmt.Printf("Error establishing point: %v\n", err)
		return
	}

	fromState := t.State
	t.State = StatePoint

	// Convert roll total to proper Point enum value
	point, err := rollTotalToPoint(roll.Total)
	if err != nil {
		fmt.Printf("Error converting roll total to point: %v\n", err)
		return
	}
	t.Point = point

	// Log state transition
	t.LogStateTransition(fromState, t.State, roll, "point establishment")

	// Trigger bet resolution for come out bets
	t.BetResolver.ResolveBets(roll)

	// Resolve come bets
	comeBetResults := t.resolveAllComeBets(roll)
	for _, result := range comeBetResults {
		fmt.Printf("Come bet result: %s\n", result)
	}

	// Resolve odds bets
	oddsBetResults := t.resolveAllOddsBets(roll)
	for _, result := range oddsBetResults {
		fmt.Printf("Odds bet result: %s\n", result)
	}

	// Log state transition
	fmt.Printf("Point established: %d\n", roll.Total)
}

// resolvePoint resolves the point when the point number is rolled again
func (t *Table) resolvePoint(roll *Roll) {
	// Validate state transition
	if err := t.validateStateTransition(StatePoint, StateComeOut, roll); err != nil {
		fmt.Printf("Error resolving point: %v\n", err)
		return
	}

	fromState := t.State
	t.State = StateComeOut
	t.Point = PointOff

	// Log state transition
	t.LogStateTransition(fromState, t.State, roll, "point resolution")

	// Trigger bet resolution for point resolution
	t.BetResolver.ResolveBets(roll)

	// Resolve come bets
	comeBetResults := t.resolveAllComeBets(roll)
	for _, result := range comeBetResults {
		fmt.Printf("Come bet result: %s\n", result)
	}

	// Resolve odds bets
	oddsBetResults := t.resolveAllOddsBets(roll)
	for _, result := range oddsBetResults {
		fmt.Printf("Odds bet result: %s\n", result)
	}

	// Log state transition
	fmt.Printf("Point resolved: %d\n", roll.Total)
}

// sevenOut handles seven-out when a 7 is rolled during point phase
func (t *Table) sevenOut(roll *Roll) {
	// Validate state transition
	if err := t.validateStateTransition(StatePoint, StateSevenOut, roll); err != nil {
		fmt.Printf("Error seven out: %v\n", err)
		return
	}

	fromState := t.State
	t.State = StateSevenOut

	// Log state transition
	t.LogStateTransition(fromState, t.State, roll, "seven out")

	// Trigger bet resolution for seven-out
	t.BetResolver.ResolveBets(roll)

	// Resolve come bets
	comeBetResults := t.resolveAllComeBets(roll)
	for _, result := range comeBetResults {
		fmt.Printf("Come bet result: %s\n", result)
	}

	// Resolve odds bets
	oddsBetResults := t.resolveAllOddsBets(roll)
	for _, result := range oddsBetResults {
		fmt.Printf("Odds bet result: %s\n", result)
	}

	// Assign new shooter
	t.assignNewShooter()

	// Reset to come out phase
	t.State = StateComeOut
	t.Point = PointOff

	// Log final state transition
	t.LogStateTransition(StateSevenOut, t.State, roll, "come out after seven out")

	// Log state transition
	fmt.Printf("Seven out! New shooter: %s\n", t.Shooter)
}

// natural handles natural wins (7 or 11) during come out
func (t *Table) natural(roll *Roll) {
	// Validate state transition
	if err := t.validateStateTransition(StateComeOut, StateComeOut, roll); err != nil {
		fmt.Printf("Error natural: %v\n", err)
		return
	}

	// Trigger bet resolution for natural
	t.BetResolver.ResolveBets(roll)

	// Resolve come bets
	comeBetResults := t.resolveAllComeBets(roll)
	for _, result := range comeBetResults {
		fmt.Printf("Come bet result: %s\n", result)
	}

	// Resolve odds bets
	oddsBetResults := t.resolveAllOddsBets(roll)
	for _, result := range oddsBetResults {
		fmt.Printf("Odds bet result: %s\n", result)
	}

	// Log state transition
	fmt.Printf("Natural: %d\n", roll.Total)
}

// craps handles craps (2, 3, or 12) during come out
func (t *Table) craps(roll *Roll) {
	// Validate state transition
	if err := t.validateStateTransition(StateComeOut, StateComeOut, roll); err != nil {
		fmt.Printf("Error craps: %v\n", err)
		return
	}

	// Trigger bet resolution for craps
	t.BetResolver.ResolveBets(roll)

	// Resolve come bets
	comeBetResults := t.resolveAllComeBets(roll)
	for _, result := range comeBetResults {
		fmt.Printf("Come bet result: %s\n", result)
	}

	// Resolve odds bets
	oddsBetResults := t.resolveAllOddsBets(roll)
	for _, result := range oddsBetResults {
		fmt.Printf("Odds bet result: %s\n", result)
	}

	// Log state transition
	fmt.Printf("Craps: %d\n", roll.Total)
}

// validateStateTransition validates if a state transition is valid
func (t *Table) validateStateTransition(fromState GameState, toState GameState, roll *Roll) error {
	// Validate that we're not in an invalid state
	if fromState < StateComeOut || fromState > StateSevenOut {
		return fmt.Errorf("invalid current state: %d", fromState)
	}

	if toState < StateComeOut || toState > StateSevenOut {
		return fmt.Errorf("invalid target state: %d", toState)
	}

	// Validate specific transitions
	switch fromState {
	case StateComeOut:
		switch toState {
		case StatePoint:
			// Valid: point establishment
			if roll.Total < 4 || roll.Total > 10 || roll.Total == 7 {
				return fmt.Errorf("invalid point number: %d", roll.Total)
			}
		case StateComeOut:
			// Valid: natural or craps
			if roll.Total != 2 && roll.Total != 3 && roll.Total != 7 && roll.Total != 11 && roll.Total != 12 {
				return fmt.Errorf("invalid come out roll: %d", roll.Total)
			}
		default:
			return fmt.Errorf("invalid transition from come out to state %d", toState)
		}
	case StatePoint:
		switch toState {
		case StateComeOut:
			// Valid: point resolution
			pointNumber, err := PointToNumber(t.Point)
			if err != nil {
				return fmt.Errorf("invalid point state: %v", err)
			}
			if roll.Total != pointNumber && roll.Total != 7 {
				return fmt.Errorf("invalid point phase roll: %d", roll.Total)
			}
		case StateSevenOut:
			// Valid: seven out
			if roll.Total != 7 {
				return fmt.Errorf("invalid seven out roll: %d", roll.Total)
			}
		default:
			return fmt.Errorf("invalid transition from point to state %d", toState)
		}
	case StateSevenOut:
		// Seven out should only transition to come out
		if toState != StateComeOut {
			return fmt.Errorf("invalid transition from seven out to state %d", toState)
		}
	}

	return nil
}

// validatePoint validates that a point number is valid
func (t *Table) validatePoint(point Point) error {
	switch point {
	case PointOff, Point4, Point5, Point6, Point8, Point9, Point10:
		return nil
	default:
		return fmt.Errorf("invalid point: %d", point)
	}
}

// validateShooter validates that the shooter exists and is valid
func (t *Table) validateShooter(shooterID string) error {
	if shooterID == "" {
		return fmt.Errorf("no shooter assigned")
	}

	_, exists := t.Players[shooterID]
	if !exists {
		return fmt.Errorf("shooter %s not found", shooterID)
	}

	return nil
}

// validateTableState validates the overall table state
func (t *Table) validateTableState() error {
	// Validate game state
	if t.State < StateComeOut || t.State > StateSevenOut {
		return fmt.Errorf("invalid game state: %d", t.State)
	}

	// Validate point
	if err := t.validatePoint(t.Point); err != nil {
		return err
	}

	// Validate shooter if we have players
	if len(t.Players) > 0 {
		if err := t.validateShooter(t.Shooter); err != nil {
			return err
		}
	}

	// Validate that point is only set during point phase
	if t.State == StateComeOut && t.Point != PointOff {
		return fmt.Errorf("point should be off during come out phase")
	}

	if t.State == StatePoint && t.Point == PointOff {
		return fmt.Errorf("point should be set during point phase")
	}

	return nil
}

// PlaceBet places a bet on the table
func (t *Table) PlaceBet(playerID, betType string, amount float64, numbers []int) (*Bet, error) {
	player, exists := t.Players[playerID]
	if !exists {
		return nil, fmt.Errorf("player %s not found", playerID)
	}

	// Create bet object for comprehensive validation
	bet := &Bet{
		ID:       generateBetID(),
		Type:     betType,
		Amount:   amount,
		Player:   playerID,
		PlacedAt: time.Now(),
		Working:  true,
		Numbers:  numbers,
	}

	// Comprehensive validation using validation functions
	// Note: Since we can't directly import the validation functions from crapsql package,
	// we'll implement the same validation logic here to ensure comprehensive validation

	// Validate bet amount
	if err := t.validateBetAmount(amount); err != nil {
		return nil, fmt.Errorf("bet amount validation failed: %v", err)
	}

	// Validate bankroll
	if err := t.validateBankroll(player, amount); err != nil {
		return nil, fmt.Errorf("bankroll validation failed: %v", err)
	}

	// Validate bet type
	if err := t.validateBetType(betType); err != nil {
		return nil, fmt.Errorf("bet type validation failed: %v", err)
	}

	// Validate game state for this bet type
	if err := t.validateGameState(betType, t.State); err != nil {
		return nil, fmt.Errorf("game state validation failed: %v", err)
	}

	// Validate bet placement (comprehensive validation)
	if err := t.validateBetPlacement(bet, player); err != nil {
		return nil, fmt.Errorf("bet placement validation failed: %v", err)
	}

	// Deduct from bankroll
	player.Bankroll -= amount
	player.Bets = append(player.Bets, bet)

	return bet, nil
}

// removeBet removes a bet from the table
func (t *Table) removeBet(betID string) {
	for _, player := range t.Players {
		for i, bet := range player.Bets {
			if bet.ID == betID {
				// Return bet amount to bankroll if bet is still working
				if bet.Working {
					player.Bankroll += bet.Amount
				}
				// Remove bet from slice
				player.Bets = append(player.Bets[:i], player.Bets[i+1:]...)
				return
			}
		}
	}
}

// GetPlayer returns a player by ID
func (t *Table) GetPlayer(id string) (*Player, error) {
	player, exists := t.Players[id]
	if !exists {
		return nil, fmt.Errorf("player %s not found", id)
	}
	return player, nil
}

// GetState returns the current game state
func (t *Table) GetState() GameState {
	return t.State
}

// GetPoint returns the current point
func (t *Table) GetPoint() Point {
	return t.Point
}

// GetShooter returns the current shooter
func (t *Table) GetShooter() string {
	return t.Shooter
}

// IsComeOut returns true if we're in come out phase
func (t *Table) IsComeOut() bool {
	return t.State == StateComeOut
}

// IsPoint returns true if we have a point established
func (t *Table) IsPoint() bool {
	return t.State == StatePoint
}

func generateBetID() string {
	// Generate a random 8-character alphanumeric ID
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to timestamp if crypto/rand fails
			return fmt.Sprintf("bet_%d", time.Now().UnixNano())
		}
		result[i] = charset[n.Int64()]
	}

	return fmt.Sprintf("bet_%s", string(result))
}

// String returns the string representation of a GameState
func (s GameState) String() string {
	switch s {
	case StateComeOut:
		return "COME_OUT"
	case StatePoint:
		return "POINT"
	case StateSevenOut:
		return "SEVEN_OUT"
	default:
		return "UNKNOWN"
	}
}

// String returns the string representation of a Point
func (p Point) String() string {
	switch p {
	case PointOff:
		return "OFF"
	case Point4:
		return "4"
	case Point5:
		return "5"
	case Point6:
		return "6"
	case Point8:
		return "8"
	case Point9:
		return "9"
	case Point10:
		return "10"
	default:
		return "UNKNOWN"
	}
}

// GetStateString returns the current game state as a string
func (t *Table) GetStateString() string {
	return t.State.String()
}

// GetPointString returns the current point as a string
func (t *Table) GetPointString() string {
	return t.Point.String()
}

// IsPointEstablished returns true if a point is currently established
func (t *Table) IsPointEstablished() bool {
	return t.State == StatePoint && t.Point != PointOff
}

// GetPointNumber returns the current point number as an integer
func (t *Table) GetPointNumber() int {
	if t.State == StatePoint && t.Point != PointOff {
		pointNumber, err := PointToNumber(t.Point)
		if err != nil {
			return 0
		}
		return pointNumber
	}
	return 0
}

// LogStateTransition logs a state transition for debugging
func (t *Table) LogStateTransition(fromState GameState, toState GameState, roll *Roll, reason string) {
	fmt.Printf("State transition: %s -> %s (roll: %d, reason: %s)\n",
		fromState.String(), toState.String(), roll.Total, reason)
}

// establishComePoint establishes a come point when a point number is rolled
func (t *Table) establishComePoint(comeBet *ComeBet, roll *Roll) error {
	// Validate come bet exists
	if comeBet == nil {
		return fmt.Errorf("come bet is nil")
	}

	// Validate come bet is not already established
	if comeBet.ComePoint != PointOff {
		return fmt.Errorf("come bet %s is already established on point %s", comeBet.ID, comeBet.ComePoint.String())
	}

	// Convert roll total to proper Point enum value
	point, err := rollTotalToPoint(roll.Total)
	if err != nil {
		return fmt.Errorf("invalid point number for come bet: %d", roll.Total)
	}

	// Establish the come point
	comeBet.ComePoint = point

	// Log come point establishment
	fmt.Printf("Come bet %s established on point %d\n", comeBet.ID, roll.Total)

	return nil
}

// resolveComeBet resolves a come bet based on the current roll
func (t *Table) resolveComeBet(comeBet *ComeBet, roll *Roll, player *Player) string {
	if comeBet == nil || player == nil {
		return "Error: Invalid come bet or player"
	}

	// Check if come bet is established
	if comeBet.ComePoint == PointOff {
		// Come bet not established yet - check for immediate win/loss
		switch roll.Total {
		case 7, 11:
			// Natural - come bet wins
			winnings := comeBet.Amount * 2 // 1:1 payout
			player.Bankroll += winnings
			t.removeComeBet(comeBet.ID)
			return fmt.Sprintf("Come bet %s wins on natural %d: $%.2f", comeBet.ID, roll.Total, winnings)
		case 2, 3, 12:
			// Craps - come bet loses
			t.removeComeBet(comeBet.ID)
			return fmt.Sprintf("Come bet %s loses on craps %d", comeBet.ID, roll.Total)
		default:
			// Point number - establish come point
			if err := t.establishComePoint(comeBet, roll); err != nil {
				return fmt.Sprintf("Error establishing come point: %v", err)
			}
			return fmt.Sprintf("Come bet %s established on point %d", comeBet.ID, roll.Total)
		}
	} else {
		// Come bet is established - check for resolution
		comePointNumber, err := PointToNumber(comeBet.ComePoint)
		if err != nil {
			return fmt.Sprintf("Error getting come point number: %v", err)
		}

		if roll.Total == comePointNumber {
			// Point made - come bet wins
			winnings := comeBet.Amount * 2 // 1:1 payout
			if comeBet.Odds > 0 {
				// Calculate odds payout based on point
				oddsPayout := t.calculateComeOddsPayout(comeBet.Odds, comePointNumber)
				winnings += oddsPayout
			}
			player.Bankroll += winnings
			t.removeComeBet(comeBet.ID)
			return fmt.Sprintf("Come bet %s wins on point %d: $%.2f", comeBet.ID, roll.Total, winnings)
		} else if roll.Total == 7 {
			// Seven out - come bet loses
			t.removeComeBet(comeBet.ID)
			return fmt.Sprintf("Come bet %s loses on seven out", comeBet.ID)
		}
		// Other numbers don't affect established come bet
		return fmt.Sprintf("Come bet %s continues on point %d", comeBet.ID, comeBet.ComePoint)
	}
}

// calculateComeOddsPayout calculates the payout for come odds based on point
func (t *Table) calculateComeOddsPayout(oddsAmount float64, point int) float64 {
	switch point {
	case 4, 10:
		return oddsAmount * 2 // 2:1 odds
	case 5, 9:
		return oddsAmount * 1.5 // 3:2 odds
	case 6, 8:
		return oddsAmount * 1.2 // 6:5 odds
	default:
		return 0
	}
}

// removeComeBet removes a come bet from tracking
func (t *Table) removeComeBet(comeBetID string) {
	delete(t.ComeBets, comeBetID)
}

// resolveAllComeBets resolves all come bets for a given roll
func (t *Table) resolveAllComeBets(roll *Roll) []string {
	var results []string

	for _, comeBet := range t.ComeBets {
		player := t.Players[comeBet.PlayerID]
		if player != nil {
			result := t.resolveComeBet(comeBet, roll, player)
			results = append(results, result)
		}
	}

	return results
}

// calculateOddsPayout calculates the payout for odds based on point
func (t *Table) calculateOddsPayout(oddsBet *OddsBet, point int) float64 {
	if oddsBet == nil {
		return 0
	}

	// Calculate true odds based on point number
	switch point {
	case 4, 10:
		return oddsBet.Amount * 2 // 2:1 odds
	case 5, 9:
		return oddsBet.Amount * 1.5 // 3:2 odds
	case 6, 8:
		return oddsBet.Amount * 1.2 // 6:5 odds
	default:
		return 0
	}
}

// resolveOddsBet resolves an odds bet based on the current roll
func (t *Table) resolveOddsBet(oddsBet *OddsBet, roll *Roll, player *Player) string {
	if oddsBet == nil || player == nil {
		return "Error: Invalid odds bet or player"
	}

	// Odds bets are resolved based on the base bet type
	switch oddsBet.BaseBetType {
	case "PASS_LINE":
		return t.resolvePassLineOdds(oddsBet, roll, player)
	case "DONT_PASS":
		return t.resolveDontPassOdds(oddsBet, roll, player)
	case "COME":
		return t.resolveComeOdds(oddsBet, roll, player)
	case "DONT_COME":
		return t.resolveDontComeOdds(oddsBet, roll, player)
	default:
		return fmt.Sprintf("Unknown odds bet type: %s", oddsBet.BaseBetType)
	}
}

// resolvePassLineOdds resolves pass line odds
func (t *Table) resolvePassLineOdds(oddsBet *OddsBet, roll *Roll, player *Player) string {
	// Pass line odds win when point is made, lose on seven out
	pointNumber, err := PointToNumber(oddsBet.Point)
	if err != nil {
		return fmt.Sprintf("Error getting point number: %v", err)
	}

	if roll.Total == pointNumber {
		// Point made - odds win
		winnings := t.calculateOddsPayout(oddsBet, pointNumber)
		player.Bankroll += winnings
		t.removeOddsBet(oddsBet.ID)
		return fmt.Sprintf("Pass line odds %s wins on point %d: $%.2f", oddsBet.ID, roll.Total, winnings)
	} else if roll.Total == 7 {
		// Seven out - odds lose
		t.removeOddsBet(oddsBet.ID)
		return fmt.Sprintf("Pass line odds %s loses on seven out", oddsBet.ID)
	}
	// Other numbers don't affect pass line odds
	return fmt.Sprintf("Pass line odds %s continues on point %d", oddsBet.ID, oddsBet.Point)
}

// resolveDontPassOdds resolves don't pass odds
func (t *Table) resolveDontPassOdds(oddsBet *OddsBet, roll *Roll, player *Player) string {
	// Don't pass odds win on seven out, lose when point is made
	pointNumber, err := PointToNumber(oddsBet.Point)
	if err != nil {
		return fmt.Sprintf("Error getting point number: %v", err)
	}

	if roll.Total == 7 {
		// Seven out - odds win
		winnings := t.calculateOddsPayout(oddsBet, pointNumber)
		player.Bankroll += winnings
		t.removeOddsBet(oddsBet.ID)
		return fmt.Sprintf("Don't pass odds %s wins on seven out: $%.2f", oddsBet.ID, winnings)
	} else if roll.Total == pointNumber {
		// Point made - odds lose
		t.removeOddsBet(oddsBet.ID)
		return fmt.Sprintf("Don't pass odds %s loses on point %d", oddsBet.ID, roll.Total)
	}
	// Other numbers don't affect don't pass odds
	return fmt.Sprintf("Don't pass odds %s continues on point %d", oddsBet.ID, oddsBet.Point)
}

// resolveComeOdds resolves come odds
func (t *Table) resolveComeOdds(oddsBet *OddsBet, roll *Roll, player *Player) string {
	// Come odds win when come point is made, lose on seven out
	pointNumber, err := PointToNumber(oddsBet.Point)
	if err != nil {
		return fmt.Sprintf("Error getting point number: %v", err)
	}

	if roll.Total == pointNumber {
		// Come point made - odds win
		winnings := t.calculateOddsPayout(oddsBet, pointNumber)
		player.Bankroll += winnings
		t.removeOddsBet(oddsBet.ID)
		return fmt.Sprintf("Come odds %s wins on point %d: $%.2f", oddsBet.ID, roll.Total, winnings)
	} else if roll.Total == 7 {
		// Seven out - odds lose
		t.removeOddsBet(oddsBet.ID)
		return fmt.Sprintf("Come odds %s loses on seven out", oddsBet.ID)
	}
	// Other numbers don't affect come odds
	return fmt.Sprintf("Come odds %s continues on point %d", oddsBet.ID, oddsBet.Point)
}

// resolveDontComeOdds resolves don't come odds
func (t *Table) resolveDontComeOdds(oddsBet *OddsBet, roll *Roll, player *Player) string {
	// Don't come odds win on seven out, lose when come point is made
	pointNumber, err := PointToNumber(oddsBet.Point)
	if err != nil {
		return fmt.Sprintf("Error getting point number: %v", err)
	}

	if roll.Total == 7 {
		// Seven out - odds win
		winnings := t.calculateOddsPayout(oddsBet, pointNumber)
		player.Bankroll += winnings
		t.removeOddsBet(oddsBet.ID)
		return fmt.Sprintf("Don't come odds %s wins on seven out: $%.2f", oddsBet.ID, winnings)
	} else if roll.Total == pointNumber {
		// Come point made - odds lose
		t.removeOddsBet(oddsBet.ID)
		return fmt.Sprintf("Don't come odds %s loses on point %d", oddsBet.ID, roll.Total)
	}
	// Other numbers don't affect don't come odds
	return fmt.Sprintf("Don't come odds %s continues on point %d", oddsBet.ID, oddsBet.Point)
}

// removeOddsBet removes an odds bet from tracking
func (t *Table) removeOddsBet(oddsBetID string) {
	delete(t.OddsBets, oddsBetID)
}

// resolveAllOddsBets resolves all odds bets for a given roll
func (t *Table) resolveAllOddsBets(roll *Roll) []string {
	var results []string

	for _, oddsBet := range t.OddsBets {
		player := t.Players[oddsBet.PlayerID]
		if player != nil {
			result := t.resolveOddsBet(oddsBet, roll, player)
			results = append(results, result)
		}
	}

	return results
}

// validateBetAmount validates that the bet amount is within table limits
func (t *Table) validateBetAmount(amount float64) error {
	if amount < t.MinBet {
		return fmt.Errorf("bet amount $%.2f is below minimum $%.2f", amount, t.MinBet)
	}
	if amount > t.MaxBet {
		return fmt.Errorf("bet amount $%.2f exceeds maximum $%.2f", amount, t.MaxBet)
	}
	if amount <= 0 {
		return fmt.Errorf("bet amount must be positive, got $%.2f", amount)
	}
	return nil
}

// validateBankroll validates that the player has sufficient bankroll
func (t *Table) validateBankroll(player *Player, amount float64) error {
	if amount > player.Bankroll {
		return fmt.Errorf("insufficient bankroll: $%.2f available, $%.2f required", player.Bankroll, amount)
	}
	return nil
}

// validateBetType validates that the bet type is valid
func (t *Table) validateBetType(betType string) error {
	// Check if bet type exists in canonical definitions
	_, exists := CanonicalBetDefinitions[betType]
	if !exists {
		return fmt.Errorf("unknown bet type: %s", betType)
	}
	return nil
}

// validateGameState validates that the bet type is valid for current game state
func (t *Table) validateGameState(betType string, state GameState) error {
	// Get bet definition
	betDef, exists := CanonicalBetDefinitions[betType]
	if !exists {
		return fmt.Errorf("unknown bet type: %s", betType)
	}

	// Check if bet requires come-out phase
	if betDef.RequiresComeOut && state != StateComeOut {
		return fmt.Errorf("bet type %s can only be placed during come-out phase", betType)
	}

	// Check if bet requires point phase
	if betDef.RequiresPoint && state != StatePoint {
		return fmt.Errorf("bet type %s can only be placed during point phase", betType)
	}

	return nil
}

// validateBetPlacement performs comprehensive validation of bet placement
func (t *Table) validateBetPlacement(bet *Bet, player *Player) error {
	// Validate bet object
	if bet == nil {
		return fmt.Errorf("bet object is nil")
	}

	// Validate player
	if player == nil {
		return fmt.Errorf("player object is nil")
	}

	// Validate bet type
	if err := t.validateBetType(bet.Type); err != nil {
		return err
	}

	// Validate bet amount
	if err := t.validateBetAmount(bet.Amount); err != nil {
		return err
	}

	// Validate bankroll
	if err := t.validateBankroll(player, bet.Amount); err != nil {
		return err
	}

	// Validate game state
	if err := t.validateGameState(bet.Type, t.State); err != nil {
		return err
	}

	// Validate numbers for bets that require specific numbers
	if len(bet.Numbers) > 0 {
		for _, num := range bet.Numbers {
			if num < 1 || num > 12 {
				return fmt.Errorf("invalid number %d for bet type %s", num, bet.Type)
			}
		}
	}

	return nil
}

// rollTotalToPoint converts a roll total to the corresponding Point enum value
func rollTotalToPoint(total int) (Point, error) {
	switch total {
	case 4:
		return Point4, nil
	case 5:
		return Point5, nil
	case 6:
		return Point6, nil
	case 8:
		return Point8, nil
	case 9:
		return Point9, nil
	case 10:
		return Point10, nil
	default:
		return PointOff, fmt.Errorf("invalid point number: %d", total)
	}
}

// PointToNumber converts a Point enum value to its corresponding number
func PointToNumber(point Point) (int, error) {
	switch point {
	case Point4:
		return 4, nil
	case Point5:
		return 5, nil
	case Point6:
		return 6, nil
	case Point8:
		return 8, nil
	case Point9:
		return 9, nil
	case Point10:
		return 10, nil
	case PointOff:
		return 0, nil
	default:
		return 0, fmt.Errorf("invalid point enum value: %d", point)
	}
}
