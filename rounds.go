package crash

// RoundsPlayed returns the number of played rounds. It basically sums the number
// of times the function Game.Next has been called.
func (g *Game) RoundsPlayed() int {
	return (len(g.roundHashes) - g.roundIndex) / nOfBytesInHash
}

// FirstRoundHash returns the hash of the first round.
func (g *Game) FirstRoundHash() []byte {
	return g.roundHashes[len(g.roundHashes)-nOfBytesInHash:]
}
