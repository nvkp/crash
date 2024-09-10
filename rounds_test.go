package crash_test

import (
	"encoding/hex"
	"testing"

	"github.com/nvkp/crash"
)

func TestRoundNumber(t *testing.T) {
	// create a game with 10 rounds
	game, _ := crash.New(
		[]byte("this is a seed"),
		crash.WithRounds(10),
	)

	// no round has been evoked, so the number of rounds played is 0
	equal(t, 0, game.RoundsPlayed(),
		"the current round number should be 0")

	// after 5 rounds are played, the current number of played rounds should be 5
	for i := 0; i < 5; i++ {
		_ = game.Next()
	}
	equal(t, 5, game.RoundsPlayed(),
		"the current round number should be 5")

	// after 5 more rounds are played, the number of played rounds is 10
	for i := 0; i < 5; i++ {
		_ = game.Next()
	}
	equal(t, 10, game.RoundsPlayed(),
		"the current round number should be 10")

	// the number of played rounds does not increases after 10 even when the Next
	// function is called again.
	_ = game.Next()
	equal(t, 10, game.RoundsPlayed(),
		"the current round number should be 10")
}

func TestFirstRoundHash(t *testing.T) {
	// create a game with 10 rounds
	game, _ := crash.New(
		[]byte("this is a seed"),
		crash.WithRounds(10),
	)

	// function should return the correct "terminating hash"
	expected := "ec00cd71f6aca6bd8744fc5a95dd3b121be25d17e42da65ab85989c62b07e57c"
	equal(t, expected, hex.EncodeToString(game.FirstRoundHash()),
		"function should return correct first round hash")

	// function returns the same hash even after few rounds are played
	for i := 0; i < 5; i++ {
		_ = game.Next()
	}
	equal(t, expected, hex.EncodeToString(game.FirstRoundHash()),
		"function should return correct first round hash")

	// function returns the same hash even after all rounds are played
	for i := 0; i < 5; i++ {
		_ = game.Next()
	}

	equal(t, expected, hex.EncodeToString(game.FirstRoundHash()),
		"function should return correct first round hash")
}

func TestRoundCount(t *testing.T) {
	// create a game with 10 rounds
	game, _ := crash.New(
		[]byte("this is a seed"),
		crash.WithRounds(10),
	)

	// go through all rounds and check the number of rounds
	var counter int
	for game.Next() {
		counter++
	}
	equal(t, 10, counter, "game should have the correct number of rounds")
}
