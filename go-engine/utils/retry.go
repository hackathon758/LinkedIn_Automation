package utils

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxRetries     int
	InitialDelay   time.Duration
	MaxDelay       time.Duration
	JitterPercent  float64
}

// DefaultRetryConfig returns sensible defaults
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    5,
		InitialDelay:  1 * time.Second,
		MaxDelay:      30 * time.Second,
		JitterPercent: 0.5,
	}
}

// RetryWithBackoff executes a function with exponential backoff
func RetryWithBackoff(cfg RetryConfig, fn func() error) error {
	var lastErr error
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err

		if attempt == cfg.MaxRetries {
			break
		}

		// Calculate exponential backoff with jitter
		delay := cfg.InitialDelay * time.Duration(math.Pow(2, float64(attempt)))
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}

		// Add jitter
		jitter := float64(delay) * cfg.JitterPercent * (rng.Float64()*2 - 1)
		delay = time.Duration(float64(delay) + jitter)

		time.Sleep(delay)
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// IsTransientError checks if an error is transient and retryable
func IsTransientError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Common transient errors
	transientPatterns := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"temporary",
		"network",
		"502",
		"503",
		"504",
	}
	for _, pattern := range transientPatterns {
		if containsIgnoreCase(errStr, pattern) {
			return true
		}
	}
	return false
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsIgnoreCase(s[1:], substr))
}
