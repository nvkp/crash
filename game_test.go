package crash_test

import (
	"encoding/hex"
	"testing"

	"github.com/nvkp/crash"
)

type crashTestCase struct {
	seed           []byte
	salt           []byte
	nOfRounds      int
	expectedRounds []round
	gameErr        error
}

type round struct {
	roundHash  string
	crashPoint float64
	roundErr   error
}

var crashTestCases = map[string]crashTestCase{
	"game-1": {
		seed:      []byte("77b271fe12fca03c618f63dfb79d4105726ba9d4a25bb3f1964e435ccf9cb208"),
		salt:      []byte("0000000000000000000fa3b65e43e4240d71762a5bf397d5304b2596d116859c"),
		nOfRounds: 10,
		expectedRounds: []round{
			{
				roundHash:  "90e96b6acaca0f7299d1acbc20a290b5af61148596c9ad3df2d004588f525419",
				crashPoint: 10.41,
			},
		},
	},
}

func TestCrash(t *testing.T) {
	for name, tc := range crashTestCases {
		t.Run(name, func(t *testing.T) {
			g, gameError := crash.New(
				tc.seed,
				crash.WithRounds(tc.nOfRounds),
				crash.WithSalt(tc.salt),
			)
			errorIs(t, gameError, tc.gameErr, "unexpected game error")
			for _, r := range tc.expectedRounds {
				next := g.Next()
				if !next {
					break
				}
				roundHash, crashPoint, roundErr := g.CrashPoint()
				roundHashString := hex.EncodeToString(roundHash)

				errorIs(t, roundErr, r.roundErr, "unexpected round error")
				equal(t, r.roundHash, roundHashString, "wrong round hash")
				equal(t, r.crashPoint, crashPoint, "wrong crash point")
			}
		})
	}
}

func TestVerifyGame(_ *testing.T) {
	g, _ := crash.New(
		[]byte("this is a seed"),
		crash.WithSalt([]byte("this is a salt")),
		crash.WithInstantCrashRate(30),
		crash.WithRounds(50),
	)

	type round struct {
		hash       string
		crashPoint float64
	}
	var rounds = []round{}

	for g.Next() {
		roundHash, crashPoint, err := g.CrashPoint()
		if err != nil {
			break
		}

		rounds = append(rounds, round{
			hash:       hex.EncodeToString(roundHash),
			crashPoint: crashPoint,
		})
	}

	_ = rounds

	roundHash, _ := hex.DecodeString(
		"94c1a1f23430dd9dbe78cc2ced06a2bb437c7a46ea378cfdfb2d051d5cf3f266",
	)

	_ = crash.CrashPoint(
		crash.Hash(roundHash),
		[]byte("this is a salt"),
		30,
	)
}

func TestRestoreGame(t *testing.T) {
	// create a new game
	g, _ := crash.New(
		[]byte("this is a seed"),
		crash.WithSalt([]byte("this is a salt")),
		crash.WithInstantCrashRate(30),
		crash.WithRounds(50),
	)

	// play some four rounds
	for i := 0; i < 4; i++ {
		_ = g.Next()
	}

	// store the hash chain and the round index
	hashChain, roundIndex := g.HashChain(), g.RoundIndex()

	// store the crash point and the hash of the next round
	_ = g.Next()
	crashPoint, hash, _ := g.CrashPoint()

	// create a new game from the stored hash chain and round index
	f, _ := crash.New(
		[]byte("this is a seed"),
		crash.WithSalt([]byte("this is a salt")),
		crash.WithInstantCrashRate(30),
		crash.WithRounds(50),
		crash.WithHashChain(hashChain),
		crash.WithRoundIndex(roundIndex),
	)

	// play one round and the crash point and the round hash of
	// the new game should equal the original game
	_ = f.Next()
	c, h, _ := f.CrashPoint()
	equal(t, crashPoint, c, "not matching crash point")
	equal(t, hash, h, "not matching round hash")
}
