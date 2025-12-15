package stealth

import (
	"math"
	"math/rand"
	"time"

	"linkedin-automation/config"
)

// TimingController implements randomized timing patterns (MANDATORY)
type TimingController struct {
	config config.TimingConfig
	rng    *rand.Rand
}

// NewTimingController creates a new timing controller
func NewTimingController(cfg config.TimingConfig) *TimingController {
	return &TimingController{
		config: cfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetActionDelay returns a randomized delay for general actions
// Uses normal distribution: N(μ=500ms, σ=200ms)
func (tc *TimingController) GetActionDelay() time.Duration {
	mu := 500.0 // mean in milliseconds
	sigma := 200.0

	// Box-Muller transform for normal distribution
	delay := tc.normalRandom(mu, sigma)

	// Ensure positive delay
	if delay < 100 {
		delay = 100 + tc.rng.Float64()*100
	}

	return time.Duration(delay) * time.Millisecond
}

// GetTypingDelay returns delay between keystrokes
// Uses normal distribution with configured min/max
func (tc *TimingController) GetTypingDelay() time.Duration {
	mu := float64(tc.config.TypingMinDelayMs+tc.config.TypingMaxDelayMs) / 2
	sigma := float64(tc.config.TypingMaxDelayMs-tc.config.TypingMinDelayMs) / 4

	delay := tc.normalRandom(mu, sigma)

	// Clamp to configured range
	if delay < float64(tc.config.TypingMinDelayMs) {
		delay = float64(tc.config.TypingMinDelayMs)
	}
	if delay > float64(tc.config.TypingMaxDelayMs) {
		delay = float64(tc.config.TypingMaxDelayMs)
	}

	return time.Duration(delay) * time.Millisecond
}

// GetThinkTime returns a "thinking" delay
func (tc *TimingController) GetThinkTime() time.Duration {
	minMs := tc.config.ThinkTimeMinMs
	maxMs := tc.config.ThinkTimeMaxMs

	if minMs == 0 {
		minMs = 2000
	}
	if maxMs == 0 {
		maxMs = 5000
	}

	delay := minMs + tc.rng.Intn(maxMs-minMs)
	return time.Duration(delay) * time.Millisecond
}

// GetPageLoadDelay returns delay after page navigation
func (tc *TimingController) GetPageLoadDelay() time.Duration {
	// Normal distribution N(2000ms, 500ms)
	delay := tc.normalRandom(2000, 500)

	if delay < 1000 {
		delay = 1000 + tc.rng.Float64()*500
	}

	return time.Duration(delay) * time.Millisecond
}

// GetCapitalLetterDelay returns additional delay before typing capital letters
func (tc *TimingController) GetCapitalLetterDelay() time.Duration {
	// Shift key hold simulation: 30-80ms extra
	return time.Duration(30+tc.rng.Intn(50)) * time.Millisecond
}

// ShouldIntroduceTypo determines if a typo should be introduced
func (tc *TimingController) ShouldIntroduceTypo() bool {
	return tc.rng.Float64() < tc.config.TypoProbability
}

// GetTypoCorrectionDelay returns delay for backspace correction
func (tc *TimingController) GetTypoCorrectionDelay() time.Duration {
	// Recognition delay + backspace delay
	return time.Duration(200+tc.rng.Intn(300)) * time.Millisecond
}

// GetBurstTypingCount returns number of characters to type in a "burst"
func (tc *TimingController) GetBurstTypingCount() int {
	// People sometimes type 3-8 characters quickly before pausing
	return 3 + tc.rng.Intn(6)
}

// GetBurstPauseDelay returns pause after a typing burst
func (tc *TimingController) GetBurstPauseDelay() time.Duration {
	// 100-300ms pause between bursts
	return time.Duration(100+tc.rng.Intn(200)) * time.Millisecond
}

// normalRandom generates a normally distributed random number
func (tc *TimingController) normalRandom(mu, sigma float64) float64 {
	// Box-Muller transform
	u1 := tc.rng.Float64()
	u2 := tc.rng.Float64()

	// Avoid log(0)
	if u1 == 0 {
		u1 = 0.0001
	}

	z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	return mu + sigma*z
}

// GetRandomizedDelay returns a delay within the given range with some randomization
func (tc *TimingController) GetRandomizedDelay(minMs, maxMs int) time.Duration {
	mu := float64(minMs+maxMs) / 2
	sigma := float64(maxMs-minMs) / 4

	delay := tc.normalRandom(mu, sigma)

	// Clamp
	if delay < float64(minMs) {
		delay = float64(minMs)
	}
	if delay > float64(maxMs) {
		delay = float64(maxMs)
	}

	return time.Duration(delay) * time.Millisecond
}

// SleepWithJitter sleeps for a duration with added jitter
func (tc *TimingController) SleepWithJitter(base time.Duration, jitterPercent float64) {
	jitter := float64(base) * jitterPercent * (tc.rng.Float64()*2 - 1)
	actualDelay := time.Duration(float64(base) + jitter)
	if actualDelay < 0 {
		actualDelay = base
	}
	time.Sleep(actualDelay)
}
