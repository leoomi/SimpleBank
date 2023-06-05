package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + random.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[random.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(12)
}

func RandomMoney() int64 {
	return RandomInt(0, 10000)
}

func RandomCurrency() string {
	currencies := []string{
		"USD",
		"EUR",
	}

	return currencies[random.Intn(len(currencies))]
}
