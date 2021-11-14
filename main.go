package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"log"
	"os"
)

// GameProperties defines a game
type GameProperties struct {
	GameId          string `json:"game_id" csv:"game_id"`
	DisplayName     string `json:"display_name" csv:"display_name"`
	ScoreMultiplier int    `json:"score_multiplier" csv:"score_multiplier"`
}

// GameSession is one played game session
type GameSession struct {
	GameDate           string `json:"game_date"`
	GameId             string `json:"game_id"`
	ThresholdHeartRate int    `json:"threshold_hr"`
	HeartRates         []int  `json:"heartRates"`
}

// heartRateState represents the one of the two states that a user can be in
type heartRateState int
const (
	blue heartRateState = iota
	red
)

// algo1's state machine graph: blue --- +5 buffer ---> red --- +5 buffer, flush to total ---> blue
// modifier: score multiplier (provided via GameProperties)
func (g GameSession) algo1(gp GameProperties) int {
	buffer := 0
	total := 0
	var state heartRateState

	// init start state
	if g.HeartRates[0] > g.ThresholdHeartRate {
		state = red
	} else {
		state = blue
	}

	// loop through heart rates, process, and create score
	for _, hr := range g.HeartRates[1:] {
		// Set current heart rate state
		// go provides no ternary op, so we have if/else
		var currentHrState heartRateState
		if hr > g.ThresholdHeartRate {
			currentHrState = red
		} else {
			currentHrState = blue
		}

		// Compare current and state, if different add to buffer
		if currentHrState != state {
			buffer += 5
		}

		// Flush buffer to total. We only add score if the user returns to blue (ending on Red nets no points)
		if currentHrState == blue {
			total += buffer
			buffer = 0
		}

		state = currentHrState
	}

	// score multiplier
	return total * gp.ScoreMultiplier
}

// algo2's state machine graph: blue ------> red --- +20 ---> blue
// modifier: -10 pts if ending on red state
func (g GameSession) algo2(gp GameProperties) int {
	total := 0
	var state heartRateState

	// init start state
	if g.HeartRates[0] > g.ThresholdHeartRate {
		state = red
	} else {
		state = blue
	}

	// loop through heart rates, process, and create score
	for _, hr := range g.HeartRates[1:] {
		var currentHrState heartRateState
		if hr > g.ThresholdHeartRate {
			currentHrState = red
		} else {
			currentHrState = blue
		}

		if currentHrState != state && currentHrState == blue {
			total += 20
		}

		state = currentHrState
	}

	// deduct 10 points if ending in a red state
	if state == red {
		return total - 10
	}
	return total
}

func main() {
	csvPtr := flag.String("properties", "properties.csv", "the game properties csv file name (same directory)")
	jsonPtr := flag.String("session", "session.json", "the session json file (same directory)")

	flag.Parse()

	// Handle CSV (game properties)
	csv, err := os.ReadFile(*csvPtr)
	if err != nil {
		log.Fatalf("failed to load CSV: %s\n", err)
	}

	var games []*GameProperties
	err = gocsv.UnmarshalBytes(csv, &games)
	if err != nil {
		log.Fatalf("Failed to process CSV: %s\n", err)
	}

	// name -> properties map for quick lookup
	gm := make(map[string]GameProperties)
	for _, g := range games {
		gm[g.GameId] = *g
	}

	// Handle JSON (game sessions)
	gameSessionFile, err := os.ReadFile(*jsonPtr)
	if err != nil {
		log.Fatalf("Failed to load JSON: %s\n", err)
	}

	var session GameSession
	err = json.Unmarshal(gameSessionFile, &session)
	if err != nil {
		log.Fatalf("Failed to process JSON: %s\n", err)
	}

	p, ok := gm[session.GameId]
	if ok != true {
		log.Fatalf("Missing GameID in GameProperties\n")
	}

	// Process Game
	score1 := session.algo1(p)
	score2 := session.algo2(p)

	// Publish Results
	fmt.Println("Processed Data:")
	fmt.Println(session.GameDate)
	fmt.Println(p.DisplayName)
	fmt.Printf("Score1: %d\n", score1)
	fmt.Printf("Score2: %d\n", score2)
}
