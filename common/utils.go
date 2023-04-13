package common

import (
	"math/rand"
	"time"
)

func RandIntn(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomInt := r.Intn(n)
	return randomInt
}
