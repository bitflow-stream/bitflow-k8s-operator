package scheduler

import (
	"math/rand"
	"time"
)

// Default seed: current time
var rng = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

// Mainly for tests in other packages
func Seed(seed int) {
	rng.Seed(int64(seed))
}

func selectRandomNode(nodes []string) string {
	if len(nodes) == 0 {
		return ""
	}
	index := rng.Int31n(int32(len(nodes)))
	return nodes[index]
}
