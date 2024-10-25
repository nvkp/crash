package crash

import (
	"errors"
)

// ErrTooFewRounds signals that the set number of rounds is too small.
var ErrTooFewRounds = errors.New("not enough rounds")

const (
	defaultNOfRounds        = 100
	defaultInstantCrashRate = 20
	nOfBytesInHash          = 32
	minNOfRounds            = 2
)

type Game struct {
	roundHashes      []byte
	salt             []byte
	nOfRounds        int
	roundIndex       int
	instantCrashRate int
}

// New creates a new crash game.
func New(
	seed []byte,
	options ...func(*Game),
) (*Game, error) {
	g := Game{
		instantCrashRate: defaultInstantCrashRate,
		nOfRounds:        defaultNOfRounds,
	}

	for _, o := range options {
		o(&g)
	}

	// return game as is if chain already set
	if len(g.roundHashes) > 0 && g.roundIndex > 0 {
		return &g, nil
	}

	// proceed to generate the hash chain otherwise
	chain, err := hashChain(seed, g.nOfRounds)
	if err != nil {
		return nil, err
	}

	g.roundHashes = chain
	g.roundIndex = len(chain)

	return &g, nil
}

// WithSalt sets a salt for a game.
func WithSalt(salt []byte) func(*Game) {
	return func(g *Game) {
		g.salt = salt
	}
}

// WithRounds sets an arbitrary number of rounds for the game.
func WithRounds(n int) func(*Game) {
	return func(g *Game) {
		g.nOfRounds = n
	}
}

// WithInstantCrashRate sets the custom instant crash rate.
func WithInstantCrashRate(rate int) func(*Game) {
	return func(g *Game) {
		g.instantCrashRate = rate
	}
}

// WithHashChain sets custom hash chain for the game. This
// can be used for loading a game from some persisted state.
func WithHashChain(chain []byte) func(*Game) {
	return func(g *Game) {
		g.roundHashes = chain
		if g.roundIndex > 0 {
			return
		}
		g.roundIndex = len(chain)
	}
}

// WithRoundIndex sets custom cursor to the hash chain. This
// can be used for loading a game from some persisted state.
func WithRoundIndex(roundIndex int) func(*Game) {
	return func(g *Game) {
		g.roundIndex = roundIndex
	}
}
