# crash

[![Go Reference](https://pkg.go.dev/badge/github.com/nvkp/crash.svg)](https://pkg.go.dev/github.com/nvkp/crash)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

This package provides an implementation of the Crash game often offered by online casinos. The host of the game chooses any *seed*. This seed is hashed and the hashed seed is hashed again etc. until the number of generated hashes equals the chosen number of rounds. The host of the game can commit to a not (in the moment of choosing the seed) known number or any piece of data, and later use it as salt to prove that he had control over what exact crash points would be generated for the game. 

The last generated hash serves as the hash of the first game. Together with the optional salt it is passed to a function that has the crash point as the output. The host can also specify an instant crash rate *X* causing every *xth* round to instantly crash and generate a crash point of *1* The inverted value of the instant crash rate more or less equals the resulting margin of the game.

With this package, one can create a new game with an arbitrary number of rounds, seed to generate the round hashes and a salt to alter the calculated crash points. The package exposes functions to prove that previous rounds were *fair*, meaning their outcome could not be altered once the first round of the game started to being played.

## Usage

To add this package as a dependency to you Golang module, run:

```shell
go get github.com/nvkp/crash
```

A new game can be created by saving the result of the `crash.New` function with a seed provided to generate the hash for each round. This generates a game of *100* rounds with the instant crash rate of *20*.

```golang
g, err := crash.New([]byte("this is a seed"))
```

The game can be tweaked by using a custom salt to alter the generated crash points, a custom instant crash rate, and a custom number of rounds. The minimum number of rounds is *2*.

```golang
g, err := crash.New(
    []byte("this is a seed"),
    crash.WithSalt([]byte("this is a salt")),
    crash.WithInstantCrashRate(30),
    crash.WithRounds(50),
)
```

Once a game is created, the crash points of individual rounds can be retrieved by iteration with the `crash.Game.Next` function. The function `crash.Game.CrashPoint` returns the hash of the given round and the calculated crash point.

```golang
for g.Next() {
    roundHash, crashPoint, err := g.CrashPoint()
    if err != nil {
        break
    }
}
```

A crash point for any number of rounds can be verified. By calling `crash.Hash` with the given round hash we get the hash of the previous game and with the function `crash.CrashPoint`, given that we know the used salt and the instant crash rate, we can verify that the previously generated crash points are valid.

```golang
roundHash, err := hex.DecodeString(
    "94c1a1f23430dd9dbe78cc2ced06a2bb437c7a46ea378cfdfb2d051d5cf3f266",
)

prevGameCrashPoint := crash.CrashPoint(
    crash.Hash(roundHash),
    []byte("this is a salt"),
    30,
)
```