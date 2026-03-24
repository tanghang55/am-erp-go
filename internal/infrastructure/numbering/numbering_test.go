package numbering

import (
	"regexp"
	"testing"
	"time"
)

func TestGenerateUsesReadableNumericFormat(t *testing.T) {
	value := Generate("po", time.Date(2026, 3, 10, 13, 5, 45, 0, time.Local))
	matched, err := regexp.MatchString(`^PO202603101305\d{4}$`, value)
	if err != nil {
		t.Fatalf("match error: %v", err)
	}
	if !matched {
		t.Fatalf("unexpected generated number: %s", value)
	}
}
