package numbering

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

const (
	timeLayout       = "200601021504"
	randomDigitCount = 4
)

func Generate(prefix string, t time.Time) string {
	return strings.ToUpper(strings.TrimSpace(prefix)) + t.Format(timeLayout) + randomDigits(randomDigitCount)
}

func randomDigits(length int) string {
	if length <= 0 {
		return ""
	}
	max := big.NewInt(1)
	for i := 0; i < length; i++ {
		max.Mul(max, big.NewInt(10))
	}
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return strings.Repeat("0", length)
	}
	return fmt.Sprintf("%0*d", length, n.Int64())
}
