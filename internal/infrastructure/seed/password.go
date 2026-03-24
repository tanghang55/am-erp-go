package seed

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const passwordCharset = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789!@#$%^&*"

func GeneratePassword(length int) (string, error) {
	if length < 8 {
		return "", fmt.Errorf("password length must be at least 8")
	}

	buf := make([]byte, length)
	for i := range buf {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordCharset))))
		if err != nil {
			return "", err
		}
		buf[i] = passwordCharset[n.Int64()]
	}
	return string(buf), nil
}

func SplitStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	statements := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		statements = append(statements, trimmed)
	}
	return statements
}
