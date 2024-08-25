package crash_test

import (
	"testing"

	"github.com/nvkp/crash"
)

const seed = "this is a seed"

func BenchmarkGame(b *testing.B) {
	for n := 0; n < b.N; n++ {
		g, _ := crash.New(
			[]byte(seed),
			crash.WithRounds(1000),
		)

		for g.Next() {
			_, _, err := g.CrashPoint()
			if err != nil {
				break
			}
		}
	}
}
