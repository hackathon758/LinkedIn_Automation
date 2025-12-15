package stealth

import (
	"math/rand"
	"time"

	"linkedin-automation/config"
)

// ScrollController implements random scrolling behavior
type ScrollController struct {
	config config.ScrollingConfig
	rng    *rand.Rand
}

// NewScrollController creates a new scroll controller
func NewScrollController(cfg config.ScrollingConfig) *ScrollController {
	return &ScrollController{
		config: cfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ScrollStep represents a single scroll action
type ScrollStep struct {
	DeltaY   int           // Pixels to scroll (negative = up, positive = down)
	Duration time.Duration // Time for this scroll step
}

// GenerateScrollSequence generates a natural scrolling sequence
func (sc *ScrollController) GenerateScrollSequence(targetScrollY int, currentScrollY int) []ScrollStep {
	if !sc.config.Enabled {
		return []ScrollStep{{
			DeltaY:   targetScrollY - currentScrollY,
			Duration: 100 * time.Millisecond,
		}}
	}

	distance := targetScrollY - currentScrollY
	if distance == 0 {
		return nil
	}

	var steps []ScrollStep
	remaining := distance

	for remaining != 0 {
		// Random scroll amount between min and max speed
		speed := sc.config.MinSpeed + sc.rng.Intn(sc.config.MaxSpeed-sc.config.MinSpeed+1)

		// Calculate step with acceleration/deceleration
		progress := float64(distance-remaining) / float64(distance)
		speedMultiplier := sc.getSpeedMultiplier(progress)
		stepSize := int(float64(speed) * speedMultiplier)

		if stepSize == 0 {
			stepSize = 1
		}

		// Ensure we don't overshoot
		if remaining > 0 {
			if stepSize > remaining {
				stepSize = remaining
			}
			remaining -= stepSize
		} else {
			stepSize = -stepSize
			if stepSize < remaining {
				stepSize = remaining
			}
			remaining -= stepSize
		}

		// Random duration for this step
		duration := time.Duration(20+sc.rng.Intn(30)) * time.Millisecond

		steps = append(steps, ScrollStep{
			DeltaY:   stepSize,
			Duration: duration,
		})

		// Occasionally add a small pause
		if sc.rng.Float64() < 0.1 {
			steps = append(steps, ScrollStep{
				DeltaY:   0,
				Duration: time.Duration(200+sc.rng.Intn(300)) * time.Millisecond,
			})
		}
	}

	// Add scroll-back with probability
	if sc.rng.Float64() < sc.config.ScrollBackProbability {
		scrollBackAmount := sc.config.MinSpeed + sc.rng.Intn(sc.config.MaxSpeed-sc.config.MinSpeed)
		steps = append(steps, ScrollStep{
			DeltaY:   -scrollBackAmount,
			Duration: time.Duration(50+sc.rng.Intn(50)) * time.Millisecond,
		})
		// Scroll forward again
		steps = append(steps, ScrollStep{
			DeltaY:   scrollBackAmount,
			Duration: time.Duration(50+sc.rng.Intn(50)) * time.Millisecond,
		})
	}

	return steps
}

// getSpeedMultiplier returns a speed multiplier for natural acceleration/deceleration
func (sc *ScrollController) getSpeedMultiplier(progress float64) float64 {
	// Ease-in-out curve
	if progress < 0.2 {
		// Acceleration phase
		return 0.3 + (progress * 3.5)
	} else if progress > 0.8 {
		// Deceleration phase
		return 0.3 + ((1 - progress) * 3.5)
	}
	// Constant speed in the middle
	return 1.0
}

// GetRandomScrollPause returns a random pause during scrolling
func (sc *ScrollController) GetRandomScrollPause() time.Duration {
	// 200-800ms pause
	return time.Duration(200+sc.rng.Intn(600)) * time.Millisecond
}

// ShouldPauseWhileScrolling determines if we should pause during scrolling
func (sc *ScrollController) ShouldPauseWhileScrolling() bool {
	return sc.rng.Float64() < 0.15 // 15% chance to pause
}
