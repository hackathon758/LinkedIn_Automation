package stealth

import (
	"math"
	"math/rand"
	"time"

	"linkedin-automation/config"
)

// Point represents a 2D point
type Point struct {
	X, Y float64
}

// BezierMouse implements Bézier curve mouse movement (MANDATORY)
type BezierMouse struct {
	config config.BezierConfig
	rng    *rand.Rand
}

// NewBezierMouse creates a new Bézier mouse controller
func NewBezierMouse(cfg config.BezierConfig) *BezierMouse {
	return &BezierMouse{
		config: cfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GeneratePath generates a Bézier curve path from start to end point
// Uses cubic Bézier curves: P(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃
func (bm *BezierMouse) GeneratePath(startX, startY, endX, endY float64) []Point {
	if !bm.config.Enabled {
		// Return direct path if disabled
		return []Point{{startX, startY}, {endX, endY}}
	}

	// Calculate distance and determine number of steps (Fitts's Law)
	distance := math.Sqrt(math.Pow(endX-startX, 2) + math.Pow(endY-startY, 2))
	steps := bm.calculateSteps(distance)

	// Generate control points with randomization
	p0 := Point{startX, startY}
	p3 := Point{endX, endY}

	// Control points with random offsets
	p1 := bm.generateControlPoint(p0, p3, 0.3)
	p2 := bm.generateControlPoint(p0, p3, 0.7)

	// Generate path points
	path := make([]Point, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		point := bm.cubicBezier(p0, p1, p2, p3, t)
		path = append(path, point)
	}

	// Add overshoot with probability
	if bm.rng.Float64() < bm.config.OvershootProbability {
		path = bm.addOvershoot(path, p3)
	}

	return path
}

// cubicBezier calculates point on cubic Bézier curve at parameter t
func (bm *BezierMouse) cubicBezier(p0, p1, p2, p3 Point, t float64) Point {
	mt := 1 - t
	mt2 := mt * mt
	mt3 := mt2 * mt
	t2 := t * t
	t3 := t2 * t

	x := mt3*p0.X + 3*mt2*t*p1.X + 3*mt*t2*p2.X + t3*p3.X
	y := mt3*p0.Y + 3*mt2*t*p1.Y + 3*mt*t2*p2.Y + t3*p3.Y

	return Point{x, y}
}

// calculateSteps determines number of steps based on distance (Fitts's Law)
func (bm *BezierMouse) calculateSteps(distance float64) int {
	// Fitts's Law: MT = a + b * log2(A/W + 1)
	// We use this to determine granularity
	baseSteps := int(math.Log2(distance+1) * 10)

	// Add randomization
	steps := baseSteps + bm.rng.Intn(10) - 5

	// Clamp to configured range
	if steps < bm.config.MinSteps {
		steps = bm.config.MinSteps
	}
	if steps > bm.config.MaxSteps {
		steps = bm.config.MaxSteps
	}

	return steps
}

// generateControlPoint creates a control point with random offset
func (bm *BezierMouse) generateControlPoint(start, end Point, ratio float64) Point {
	// Base position along the line
	baseX := start.X + (end.X-start.X)*ratio
	baseY := start.Y + (end.Y-start.Y)*ratio

	// Calculate perpendicular offset
	dx := end.X - start.X
	dy := end.Y - start.Y
	length := math.Sqrt(dx*dx + dy*dy)

	// Perpendicular vector (normalized)
	perpX := -dy / length
	perpY := dx / length

	// Random offset magnitude (up to 30% of distance)
	offsetMagnitude := length * 0.3 * (bm.rng.Float64() - 0.5) * 2

	return Point{
		X: baseX + perpX*offsetMagnitude,
		Y: baseY + perpY*offsetMagnitude,
	}
}

// addOvershoot adds natural overshoot to the path
func (bm *BezierMouse) addOvershoot(path []Point, target Point) []Point {
	if len(path) == 0 {
		return path
	}

	lastPoint := path[len(path)-1]

	// Calculate overshoot direction
	dx := target.X - path[0].X
	dy := target.Y - path[0].Y
	length := math.Sqrt(dx*dx + dy*dy)

	// Overshoot by 2-8 pixels
	overshootDist := 2 + bm.rng.Float64()*6

	overshootPoint := Point{
		X: lastPoint.X + (dx/length)*overshootDist,
		Y: lastPoint.Y + (dy/length)*overshootDist,
	}

	// Add overshoot point
	path = append(path, overshootPoint)

	// Add correction back to target
	correctionSteps := 3 + bm.rng.Intn(3)
	for i := 1; i <= correctionSteps; i++ {
		t := float64(i) / float64(correctionSteps)
		path = append(path, Point{
			X: overshootPoint.X + (target.X-overshootPoint.X)*t,
			Y: overshootPoint.Y + (target.Y-overshootPoint.Y)*t,
		})
	}

	return path
}

// GetMovementDurations returns durations for each step (variable velocity)
func (bm *BezierMouse) GetMovementDurations(pathLength int, totalDuration time.Duration) []time.Duration {
	if pathLength <= 1 {
		return []time.Duration{totalDuration}
	}

	durations := make([]time.Duration, pathLength-1)
	baseInterval := totalDuration / time.Duration(pathLength-1)

	for i := 0; i < pathLength-1; i++ {
		// Add acceleration at start and deceleration at end
		progress := float64(i) / float64(pathLength-1)
		speedFactor := 1.0

		// Ease-in-out curve
		if progress < 0.2 {
			// Acceleration phase
			speedFactor = 0.5 + progress*2.5
		} else if progress > 0.8 {
			// Deceleration phase
			speedFactor = 0.5 + (1-progress)*2.5
		}

		// Add randomness (±20%)
		randomFactor := 0.8 + bm.rng.Float64()*0.4
		durations[i] = time.Duration(float64(baseInterval) / speedFactor * randomFactor)
	}

	return durations
}
