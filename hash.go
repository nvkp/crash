package crash

import "crypto/sha256"

// HashChain returns a copy of the game's internal hash chain.
func (g *Game) HashChain() []byte {
	out := make([]byte, len(g.roundHashes))
	copy(out, g.roundHashes)
	return out
}

// Hash returns the hash of the previous round.
func Hash(in []byte) []byte {
	h := hash(in)
	return h[:]
}

func hashChain(seed []byte, nOfRounds int) ([]byte, error) {
	if nOfRounds < minNOfRounds {
		return nil, ErrTooFewRounds
	}

	chain := make([]byte, 0, nOfRounds*nOfBytesInHash)
	hashedSeed := hash(seed)

	for i := 0; i < nOfRounds; i++ {
		chain = append(chain, hashedSeed[:]...)
		hashedSeed = hash(hashedSeed[:])
	}

	return chain, nil
}

func hash(in []byte) [32]byte {
	return sha256.Sum256(in)
}
