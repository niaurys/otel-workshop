package random

import (
	rand "math/rand/v2"
)

func Item[T any](items []T) T {
	return items[rand.IntN(len(items))]
}

func Int(max int) int {
	return rand.IntN(max)
}

func Int64(max int64) int64 {
	return rand.Int64N(max)
}
