package stealth

import (
	"math/rand"
	"strings"
	"time"
	"unicode"

	"linkedin-automation/config"
)

// TypingSimulator implements realistic typing simulation
type TypingSimulator struct {
	config config.TimingConfig
	rng    *rand.Rand
}

// NewTypingSimulator creates a new typing simulator
func NewTypingSimulator(cfg config.TimingConfig) *TypingSimulator {
	return &TypingSimulator{
		config: cfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// TypedChar represents a character to be typed with timing info
type TypedChar struct {
	Char         rune
	Delay        time.Duration
	IsBackspace  bool
	IsShiftHeld  bool
	IsBurstPause bool
}

// Common typos map (adjacent keys on QWERTY keyboard)
var adjacentKeys = map[rune][]rune{
	'a': {'s', 'q', 'w', 'z'},
	'b': {'v', 'n', 'g', 'h'},
	'c': {'x', 'v', 'd', 'f'},
	'd': {'s', 'f', 'e', 'r', 'c', 'x'},
	'e': {'w', 'r', 'd', 's'},
	'f': {'d', 'g', 'r', 't', 'v', 'c'},
	'g': {'f', 'h', 't', 'y', 'b', 'v'},
	'h': {'g', 'j', 'y', 'u', 'n', 'b'},
	'i': {'u', 'o', 'k', 'j'},
	'j': {'h', 'k', 'u', 'i', 'm', 'n'},
	'k': {'j', 'l', 'i', 'o', 'm'},
	'l': {'k', 'o', 'p'},
	'm': {'n', 'j', 'k'},
	'n': {'b', 'm', 'h', 'j'},
	'o': {'i', 'p', 'k', 'l'},
	'p': {'o', 'l'},
	'q': {'w', 'a'},
	'r': {'e', 't', 'd', 'f'},
	's': {'a', 'd', 'w', 'e', 'x', 'z'},
	't': {'r', 'y', 'f', 'g'},
	'u': {'y', 'i', 'h', 'j'},
	'v': {'c', 'b', 'f', 'g'},
	'w': {'q', 'e', 'a', 's'},
	'x': {'z', 'c', 's', 'd'},
	'y': {'t', 'u', 'g', 'h'},
	'z': {'a', 's', 'x'},
}

// GenerateTypingSequence generates a realistic typing sequence for the given text
func (ts *TypingSimulator) GenerateTypingSequence(text string) []TypedChar {
	var sequence []TypedChar
	runes := []rune(text)

	burstCounter := 0
	burstTarget := ts.getBurstTarget()

	for i, char := range runes {
		// Check if we should introduce a typo
		if ts.rng.Float64() < ts.config.TypoProbability && len(adjacentKeys[unicode.ToLower(char)]) > 0 {
			// Add typo character
			typoChar := ts.getTypoChar(char)
			sequence = append(sequence, TypedChar{
				Char:  typoChar,
				Delay: ts.getTypingDelay(),
			})

			// Add recognition pause
			sequence = append(sequence, TypedChar{
				Delay:       time.Duration(200+ts.rng.Intn(300)) * time.Millisecond,
				IsBurstPause: true,
			})

			// Add backspace
			sequence = append(sequence, TypedChar{
				IsBackspace: true,
				Delay:       time.Duration(50+ts.rng.Intn(50)) * time.Millisecond,
			})
		}

		// Determine if shift is needed
		isShift := unicode.IsUpper(char)

		// Get base delay
		delay := ts.getTypingDelay()

		// Add extra delay for capital letters (shift key)
		if isShift {
			delay += time.Duration(30+ts.rng.Intn(50)) * time.Millisecond
		}

		// Add extra delay after punctuation or space
		if i > 0 {
			prevChar := runes[i-1]
			if prevChar == '.' || prevChar == '!' || prevChar == '?' {
				delay += time.Duration(100+ts.rng.Intn(200)) * time.Millisecond
			} else if prevChar == ' ' && ts.rng.Float64() < 0.3 {
				delay += time.Duration(50+ts.rng.Intn(100)) * time.Millisecond
			}
		}

		sequence = append(sequence, TypedChar{
			Char:        char,
			Delay:       delay,
			IsShiftHeld: isShift,
		})

		// Handle burst typing rhythm
		burstCounter++
		if burstCounter >= burstTarget {
			// Add burst pause
			sequence = append(sequence, TypedChar{
				Delay:       time.Duration(100+ts.rng.Intn(200)) * time.Millisecond,
				IsBurstPause: true,
			})
			burstCounter = 0
			burstTarget = ts.getBurstTarget()
		}
	}

	return sequence
}

// getTypingDelay returns a randomized typing delay
func (ts *TypingSimulator) getTypingDelay() time.Duration {
	min := ts.config.TypingMinDelayMs
	max := ts.config.TypingMaxDelayMs

	if min == 0 {
		min = 50
	}
	if max == 0 {
		max = 150
	}

	delay := min + ts.rng.Intn(max-min+1)
	return time.Duration(delay) * time.Millisecond
}

// getTypoChar returns an adjacent character for a typo
func (ts *TypingSimulator) getTypoChar(original rune) rune {
	lower := unicode.ToLower(original)
	adjacent, ok := adjacentKeys[lower]
	if !ok || len(adjacent) == 0 {
		return original
	}

	typo := adjacent[ts.rng.Intn(len(adjacent))]

	// Preserve case
	if unicode.IsUpper(original) {
		return unicode.ToUpper(typo)
	}
	return typo
}

// getBurstTarget returns the number of characters to type before a burst pause
func (ts *TypingSimulator) getBurstTarget() int {
	return 3 + ts.rng.Intn(6) // 3-8 characters
}

// GetTotalTypingDuration estimates total time to type a string
func (ts *TypingSimulator) GetTotalTypingDuration(text string) time.Duration {
	avgDelay := time.Duration((ts.config.TypingMinDelayMs+ts.config.TypingMaxDelayMs)/2) * time.Millisecond
	return avgDelay * time.Duration(len(text))
}

// ShouldDoubleCharacter determines if a character should be typed twice (common typo)
func (ts *TypingSimulator) ShouldDoubleCharacter() bool {
	return ts.rng.Float64() < 0.02 // 2% chance
}

// GetRandomNeighborKey gets a random adjacent key for realistic typos
func (ts *TypingSimulator) GetRandomNeighborKey(char rune) rune {
	return ts.getTypoChar(char)
}

// SimulateTextSelection simulates selecting text (e.g., Ctrl+A behavior)
func (ts *TypingSimulator) SimulateTextSelection() time.Duration {
	// Selection + processing time
	return time.Duration(100+ts.rng.Intn(200)) * time.Millisecond
}

// SubstituteTemplate replaces template variables with actual values
func SubstituteTemplate(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
	}
	return result
}
