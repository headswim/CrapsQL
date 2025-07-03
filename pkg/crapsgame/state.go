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
}

// NewTable creates a new craps table
func NewTable(minBet, maxBet float64, maxOdds int) *Table {
	table := &Table{
		State:     StateComeOut,
		Point:     PointOff,
		Players:   make(map[string]*Player),
		MinBet:    minBet,
		MaxBet:    maxBet,
		MaxOdds:   maxOdds,
		CreatedAt: time.Now(),
	}
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

	// Remove all player's bets and return money for active bets
	for _, bet := range player.Bets {
		if bet.Working {
			// Return bet amount to player's bankroll
			player.Bankroll += bet.Amount
		}
	}

	delete(t.Players, id)

	// If this was the shooter, assign new shooter
	if t.Shooter == id {
		t.assignNewShooter()
	}

	return nil
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

	// Note: State updates are handled by the caller (ExecuteGameTurn)
	// This prevents double state updates when ROLL DICE is called

	return roll
}

// UpdateGameState updates the game state based on the current roll
func (t *Table) UpdateGameState(roll *Roll) {
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

// UpdateGameStateOnly updates only the game state based on the roll, without bet resolution
func (t *Table) UpdateGameStateOnly(roll *Roll) {
	switch t.State {
	case StateComeOut:
		switch roll.Total {
		case 7, 11:
			// Natural - stay in come out
			fmt.Printf("Natural: %d - staying in come out\n", roll.Total)
		case 2, 3, 12:
			// Craps - stay in come out
			fmt.Printf("Craps: %d - staying in come out\n", roll.Total)
		default:
			// Point established
			point, err := rollTotalToPoint(roll.Total)
			if err != nil {
				fmt.Printf("Error converting roll total to point: %v\n", err)
				return
			}
			t.State = StatePoint
			t.Point = point
			fmt.Printf("Point established: %d\n", roll.Total)
		}
	case StatePoint:
		if roll.Total == 7 {
			// Seven out - back to come out
			t.State = StateComeOut
			t.Point = PointOff
			t.assignNewShooter()
			fmt.Printf("Seven out! New shooter: %s\n", t.Shooter)
		} else {
			pointNumber, err := PointToNumber(t.Point)
			if err != nil {
				fmt.Printf("Error getting point number: %v\n", err)
				return
			}
			if roll.Total == pointNumber {
				// Point made - back to come out
				t.State = StateComeOut
				t.Point = PointOff
				fmt.Printf("Point resolved: %d\n", roll.Total)
			}
			// Other numbers don't change the point
		}
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

	// Convert roll total to proper Point enum value
	point, err := rollTotalToPoint(roll.Total)
	if err != nil {
		fmt.Printf("Error converting roll total to point: %v\n", err)
		return
	}

	// Update state
	t.State = StatePoint
	t.Point = point

	// Log state transition
	t.LogStateTransition(fromState, t.State, roll, "point establishment")
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

	// Change state
	t.State = StateComeOut
	t.Point = PointOff

	// Log state transition
	t.LogStateTransition(fromState, t.State, roll, "point resolution")
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

	// Change state
	t.State = StateSevenOut

	// Log state transition
	t.LogStateTransition(fromState, t.State, roll, "seven out")

	// Assign new shooter
	t.assignNewShooter()

	// Reset to come out phase
	t.State = StateComeOut
	t.Point = PointOff

	// Log final state transition
	t.LogStateTransition(StateSevenOut, t.State, roll, "come out after seven out")
	fmt.Printf("Seven out! New shooter: %s\n", t.Shooter)
}

// natural handles natural wins (7 or 11) during come out
func (t *Table) natural(roll *Roll) {
	// Validate state transition
	if err := t.validateStateTransition(StateComeOut, StateComeOut, roll); err != nil {
		fmt.Printf("Error natural: %v\n", err)
		return
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

	// Comprehensive validation using validation functions from crapsql package
	// Import the validation functions to ensure consistent validation across the codebase

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
				// Remove bet from slice (no money returned - this is for losing bets)
				player.Bets = append(player.Bets[:i], player.Bets[i+1:]...)
				return
			}
		}
	}
}

// removeBetWithRefund removes a bet and returns the money to the player (for voluntary removal)
func (t *Table) removeBetWithRefund(betID string) {
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

// ResolveAllBets resolves all bets using the unified ResolveBet function
func (t *Table) ResolveAllBets(roll *Roll) []string {
	var results []string

	// Process all player bets
	for _, player := range t.Players {
		var betsToRemove []*Bet

		for _, bet := range player.Bets {
			if !bet.Working {
				continue
			}

			// Use the unified ResolveBet function from canonical_bets.go
			win, payout, remove := ResolveBet(bet, roll, t.State)

			if win {
				// Bet wins - add payout to bankroll
				player.Bankroll += bet.Amount + payout
				results = append(results, fmt.Sprintf("🎉 %s wins $%.2f (payout: $%.2f)", bet.Type, bet.Amount+payout, payout))
			} else if remove {
				// Bet loses and should be removed
				results = append(results, fmt.Sprintf("💥 %s loses $%.2f", bet.Type, bet.Amount))
			}
			// If neither win nor remove, bet continues (no action needed)

			if remove {
				betsToRemove = append(betsToRemove, bet)
			}
		}

		// Remove resolved bets
		for _, bet := range betsToRemove {
			t.removeBet(bet.ID)
		}
	}

	return results
}

// RollDiceAndResolve follows the simplified game flow: roll dice, resolve bets, update state
func (t *Table) RollDiceAndResolve() (*Roll, []string) {
	// Validate shooter before roll
	if err := t.validateShooter(t.Shooter); err != nil {
		fmt.Printf("Warning: Invalid shooter before roll: %v\n", err)
		t.assignNewShooter()
	}

	// Step 1: Roll the dice
	roll := &Roll{
		Die1: rollDieSecure(),
		Die2: rollDieSecure(),
		Time: time.Now(),
	}
	roll.Total = roll.Die1 + roll.Die2
	roll.IsHard = roll.Die1 == roll.Die2

	t.CurrentRoll = roll
	t.LastRoll = roll.Time

	fmt.Printf("Rolled: %d-%d = %d\n", roll.Die1, roll.Die2, roll.Total)

	// Step 2: Resolve all bets using unified ResolveBet function
	betResults := t.ResolveAllBets(roll)

	// Step 3: Update game state (after bet resolution)
	t.UpdateGameStateOnly(roll)

	return roll, betResults
}

func (t *Table) UpdateBetWorkingStatus() {
	for _, player := range t.Players {
		for _, bet := range player.Bets {
			bet.Working = t.shouldBetBeWorking(bet, t.State)
		}
	}
}

func (t *Table) shouldBetBeWorking(bet *Bet, state GameState) bool {
	return true
}

// PlayGame implements the simplified game flow:
// 1. Place Bets (deduct bankroll + store bet)
// 2. Roll Dice
// 3. Pay/collect every bet using unified ResolveBet
// 4. Update game state based on dice
// 5. Repeat
func (t *Table) PlayGame() {
	fmt.Println("=== CRAPS GAME STARTED ===")
	fmt.Printf("Current state: %s\n", t.State.String())
	if t.State == StatePoint {
		fmt.Printf("Point: %s\n", t.Point.String())
	}
	fmt.Printf("Shooter: %s\n", t.Shooter)
	fmt.Println("Ready for bets...")
}

// ExecuteGameTurn executes one complete turn of the game
// This is the main game loop that follows your desired pattern
func (t *Table) ExecuteGameTurn() (*Roll, []string) {
	// Step 1: Roll the dice
	roll := t.RollDice()

	// Step 2: Pay/collect every bet using unified ResolveBet
	betResults := t.ResolveAllBets(roll)

	// Step 3: Update game state based on dice
	t.UpdateGameStateOnly(roll)

	return roll, betResults
}

// RemoveBet removes a specific bet type for a player
func (t *Table) RemoveBet(playerID, betType string) error {
	player, err := t.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("player %s not found", playerID)
	}

	var remainingBets []*Bet
	removedCount := 0

	for _, bet := range player.Bets {
		if bet.Type == betType && bet.Working {
			// Return bet amount to player's bankroll
			player.Bankroll += bet.Amount
			removedCount++
		} else {
			remainingBets = append(remainingBets, bet)
		}
	}

	player.Bets = remainingBets

	if removedCount == 0 {
		return fmt.Errorf("no active %s bets to remove", betType)
	}

	return nil
}

// PressBet increases the amount of a specific bet type for a player
func (t *Table) PressBet(playerID, betType string, amount float64) error {
	player, err := t.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("player %s not found", playerID)
	}

	if amount <= 0 {
		return fmt.Errorf("press amount must be positive")
	}

	if player.Bankroll < amount {
		return fmt.Errorf("insufficient bankroll for press")
	}

	pressedCount := 0
	for _, bet := range player.Bets {
		if bet.Type == betType && bet.Working {
			bet.Amount += amount
			player.Bankroll -= amount
			pressedCount++
		}
	}

	if pressedCount == 0 {
		return fmt.Errorf("no active %s bets to press", betType)
	}

	return nil
}

// TurnBet turns a specific bet type on or off for a player
func (t *Table) TurnBet(playerID, betType string, working bool) error {
	player, err := t.GetPlayer(playerID)
	if err != nil {
		return fmt.Errorf("player %s not found", playerID)
	}

	turnedCount := 0
	for _, bet := range player.Bets {
		if bet.Type == betType {
			bet.Working = working
			turnedCount++
		}
	}

	if turnedCount == 0 {
		return fmt.Errorf("no %s bets to turn", betType)
	}

	return nil
}

// rollDieSecure generates a secure random die roll (1-6)
func rollDieSecure() int {
	n, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		// Fallback to timestamp-based random if crypto/rand fails
		return int(time.Now().UnixNano()%6) + 1
	}
	return int(n.Int64()) + 1
}
