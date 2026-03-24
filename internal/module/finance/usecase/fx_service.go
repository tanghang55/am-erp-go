package usecase

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	fallbackBaseCurrency      = "USD"
	fallbackExchangeRateScale = uint32(4)
)

var defaultBaseCurrencyResolver = func() string {
	return fallbackBaseCurrency
}

var defaultFXRateResolver = func(baseCurrency, originalCurrency string, occurredAt time.Time) (*FXRateSnapshot, error) {
	base := normalizeCurrency(baseCurrency)
	original := normalizeCurrency(originalCurrency)
	if base == "" || original == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	if base == original {
		return &FXRateSnapshot{
			Rate:        1,
			Source:      fxSourceIdentity,
			Version:     fxVersionIdentity,
			EffectiveAt: occurredAt,
		}, nil
	}
	return nil, fmt.Errorf("fx resolver is not configured for %s -> %s", original, base)
}

var defaultExchangeRateScaleResolver = func() uint32 {
	return fallbackExchangeRateScale
}

func normalizeCurrency(currency string) string {
	return strings.ToUpper(strings.TrimSpace(currency))
}

func SetDefaultBaseCurrencyResolver(resolver func() string) {
	if resolver == nil {
		defaultBaseCurrencyResolver = func() string {
			return fallbackBaseCurrency
		}
		return
	}
	defaultBaseCurrencyResolver = resolver
}

func SetFXRateResolver(resolver func(baseCurrency, originalCurrency string, occurredAt time.Time) (*FXRateSnapshot, error)) {
	if resolver == nil {
		defaultFXRateResolver = func(baseCurrency, originalCurrency string, occurredAt time.Time) (*FXRateSnapshot, error) {
			base := normalizeCurrency(baseCurrency)
			original := normalizeCurrency(originalCurrency)
			if base == "" || original == "" {
				return nil, fmt.Errorf("currency is required")
			}
			if occurredAt.IsZero() {
				occurredAt = time.Now()
			}
			if base == original {
				return &FXRateSnapshot{
					Rate:        1,
					Source:      fxSourceIdentity,
					Version:     fxVersionIdentity,
					EffectiveAt: occurredAt,
				}, nil
			}
			return nil, fmt.Errorf("fx resolver is not configured for %s -> %s", original, base)
		}
		return
	}
	defaultFXRateResolver = resolver
}

func SetExchangeRateScaleResolver(resolver func() uint32) {
	if resolver == nil {
		defaultExchangeRateScaleResolver = func() uint32 {
			return fallbackExchangeRateScale
		}
		return
	}
	defaultExchangeRateScaleResolver = resolver
}

func getDefaultBaseCurrency() string {
	value := normalizeCurrency(defaultBaseCurrencyResolver())
	if value == "" {
		return fallbackBaseCurrency
	}
	return value
}

func getExchangeRateScale() uint32 {
	value := defaultExchangeRateScaleResolver()
	if value > 8 {
		return fallbackExchangeRateScale
	}
	return value
}

func ensureBaseCurrency(value string) string {
	normalized := normalizeCurrency(value)
	if normalized == "" {
		return getDefaultBaseCurrency()
	}
	return normalized
}

// resolveStaticFXRate returns base/original rate.
func resolveFXRate(baseCurrency, originalCurrency string, occurredAt time.Time) (*FXRateSnapshot, error) {
	return defaultFXRateResolver(baseCurrency, originalCurrency, occurredAt)
}

func round6(v float64) float64 {
	return math.Round(v*1_000_000) / 1_000_000
}

func roundRate(v float64, scale uint32) float64 {
	factor := math.Pow10(int(scale))
	return math.Round(v*factor) / factor
}

func roundExchangeRate(v float64) float64 {
	return roundRate(v, getExchangeRateScale())
}
