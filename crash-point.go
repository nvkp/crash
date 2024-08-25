package crash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"math/big"
	"strconv"
)

// ErrNoMoreRounds signals that there are no more rounds to be played in this game.
var ErrNoMoreRounds = errors.New("no more rounds")

var e = math.Pow(2, 52)

const instantCrash float64 = 1.0

// Next shifts the game to a new round if there are some rounds not
// played yet and returns true. If there are no not yet played rounds,
// it returns false.
func (g *Game) Next() bool {
	if g.roundIndex < 0 {
		return false
	}

	g.roundIndex = g.roundIndex - nOfBytesInHash

	return true
}

// CrashPoint returns the hash and the crash point for the current round.
func (g *Game) CrashPoint() ([]byte, float64, error) {
	if g.roundIndex < 0 {
		return nil, 0, ErrNoMoreRounds
	}

	if g.roundIndex == len(g.roundHashes) {
		g.Next()
	}

	roundHash := g.roundHash()
	crashPoint := CrashPoint(roundHash, g.salt, g.instantCrashRate)

	return roundHash[:], crashPoint, nil
}

func (g *Game) roundHash() []byte {
	return g.roundHashes[g.roundIndex : g.roundIndex+nOfBytesInHash]
}

// CrashPoint return the crash round for the given round hash, salt,
// and the instant crash rate.
func CrashPoint(
	roundHash []byte,
	salt []byte,
	instantCrashRate int,
) float64 {
	mac := hmac.New(sha256.New, roundHash)
	mac.Write(salt)
	value := mac.Sum(nil)
	if isDivisible(value, instantCrashRate) {
		return instantCrash
	}

	x := hex.EncodeToString(value)
	i, _ := strconv.ParseInt(string(x[:13]), 16, 64)
	h := float64(i)

	return math.Floor((100*e-h)/(e-h)) / 100.0
}

func isDivisible(hash []byte, modulo int) bool {
	hashInt := new(big.Int).SetBytes(hash)
	mod := big.NewInt(int64(modulo))

	remainder := new(big.Int).Mod(hashInt, mod)

	return remainder.Cmp(big.NewInt(0)) == 0
}
