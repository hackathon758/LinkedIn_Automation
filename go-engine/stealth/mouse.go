package stealth

import (
	"math/rand"
	"time"

	"linkedin-automation/config"
)

// MouseHoverController implements mouse hovering and movement
type MouseHoverController struct {
	config config.MouseConfig
	rng    *rand.Rand
}

// NewMouseHoverController creates a new mouse hover controller
func NewMouseHoverController(cfg config.MouseConfig) *MouseHoverController {
	return &MouseHoverController{
		config: cfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// HoverAction represents a hover/movement action
type HoverAction struct {
	X        float64
	Y        float64
	Duration time.Duration
	IsHover  bool
}

// GeneratePreClickSequence generates mouse movements before clicking
func (mh *MouseHoverController) GeneratePreClickSequence(targetX, targetY float64, viewportWidth, viewportHeight int) []HoverAction {
	var actions []HoverAction

	if !mh.config.HoverBeforeClick {
		return actions
	}

	// Random starting position
	currentX := float64(mh.rng.Intn(viewportWidth))
	currentY := float64(mh.rng.Intn(viewportHeight))

	// Add some random wandering if enabled
	if mh.config.RandomMovement {
		numWanders := 1 + mh.rng.Intn(3)
		for i := 0; i < numWanders; i++ {
			// Random intermediate point
			wanderX := currentX + float64(mh.rng.Intn(200)-100)
			wanderY := currentY + float64(mh.rng.Intn(200)-100)

			// Clamp to viewport
			if wanderX < 0 {
				wanderX = 0
			}
			if wanderX > float64(viewportWidth) {
				wanderX = float64(viewportWidth)
			}
			if wanderY < 0 {
				wanderY = 0
			}
			if wanderY > float64(viewportHeight) {
				wanderY = float64(viewportHeight)
			}

			actions = append(actions, HoverAction{
				X:        wanderX,
				Y:        wanderY,
				Duration: time.Duration(100+mh.rng.Intn(200)) * time.Millisecond,
				IsHover:  false,
			})

			currentX = wanderX
			currentY = wanderY
		}
	}

	// Move towards target with hover
	actions = append(actions, HoverAction{
		X:        targetX,
		Y:        targetY,
		Duration: mh.getHoverDuration(),
		IsHover:  true,
	})

	return actions
}

// getHoverDuration returns a random hover duration
func (mh *MouseHoverController) getHoverDuration() time.Duration {
	min := mh.config.HoverDurationMinMs
	max := mh.config.HoverDurationMaxMs

	if min == 0 {
		min = 100
	}
	if max == 0 {
		max = 500
	}

	return time.Duration(min+mh.rng.Intn(max-min+1)) * time.Millisecond
}

// ShouldPerformRandomMovement determines if random cursor movement should occur
func (mh *MouseHoverController) ShouldPerformRandomMovement() bool {
	if !mh.config.RandomMovement {
		return false
	}
	return mh.rng.Float64() < 0.2 // 20% chance
}

// GenerateRandomMovement generates a random cursor movement
func (mh *MouseHoverController) GenerateRandomMovement(currentX, currentY float64, viewportWidth, viewportHeight int) HoverAction {
	// Small random offset
	offsetX := float64(mh.rng.Intn(100) - 50)
	offsetY := float64(mh.rng.Intn(100) - 50)

	newX := currentX + offsetX
	newY := currentY + offsetY

	// Clamp to viewport
	if newX < 0 {
		newX = 0
	}
	if newX > float64(viewportWidth) {
		newX = float64(viewportWidth)
	}
	if newY < 0 {
		newY = 0
	}
	if newY > float64(viewportHeight) {
		newY = float64(viewportHeight)
	}

	return HoverAction{
		X:        newX,
		Y:        newY,
		Duration: time.Duration(50+mh.rng.Intn(150)) * time.Millisecond,
		IsHover:  false,
	}
}

// GeneratePostClickMovement generates natural movement after a click
func (mh *MouseHoverController) GeneratePostClickMovement(clickX, clickY float64, viewportWidth, viewportHeight int) []HoverAction {
	var actions []HoverAction

	// Small drift after click
	driftX := clickX + float64(mh.rng.Intn(20)-10)
	driftY := clickY + float64(mh.rng.Intn(20)-10)

	actions = append(actions, HoverAction{
		X:        driftX,
		Y:        driftY,
		Duration: time.Duration(50+mh.rng.Intn(100)) * time.Millisecond,
		IsHover:  false,
	})

	return actions
}
