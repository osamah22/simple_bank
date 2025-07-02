package utils

import (
	"math/rand"
	"strings"
	"time"
)

var r *rand.Rand

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return r.Int63n(max-min+1) + min
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[r.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandonOwner() string {
	return RandomString(30)
}

func RandonAmount() int64 {
	return RandomInt(0, 2500)
}

func RandomCurrency() string {
	currencies := []string{"USD", "TRY", "EUR"}

	return currencies[r.Intn(len(currencies))]
}
