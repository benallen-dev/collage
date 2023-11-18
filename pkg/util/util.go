package	util

import (
	"math/rand"
)

func GetRandomName() string {
	names := [7]string {
		"Dave",
		"Peter",
		"Jesus",
		"Mary",
		"Helen",
		"Alice",
		"Bob",
	}

	// make random number
	idx := rand.Intn(7)
	// return that element from the list
	return names[idx]
}

