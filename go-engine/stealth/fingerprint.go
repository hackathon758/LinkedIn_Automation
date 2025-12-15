package stealth

import (
	"math/rand"
	"time"

	"linkedin-automation/config"
)

// FingerprintMasker implements browser fingerprint masking (MANDATORY)
type FingerprintMasker struct {
	config config.FingerprintConfig
	rng    *rand.Rand
}

// NewFingerprintMasker creates a new fingerprint masker
func NewFingerprintMasker(cfg config.FingerprintConfig) *FingerprintMasker {
	return &FingerprintMasker{
		config: cfg,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Viewport represents browser viewport dimensions
type Viewport struct {
	Width  int
	Height int
}

// Common user agents for rotation
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
}

// Common viewport sizes (realistic desktop resolutions)
var viewports = []Viewport{
	{1920, 1080},
	{1366, 768},
	{1536, 864},
	{1440, 900},
	{1280, 720},
	{1600, 900},
	{2560, 1440},
	{1680, 1050},
}

// Common timezones
var timezones = []string{
	"America/New_York",
	"America/Chicago",
	"America/Denver",
	"America/Los_Angeles",
	"America/Phoenix",
	"America/Detroit",
	"America/Indianapolis",
}

// Common Accept-Language headers
var acceptLanguages = []string{
	"en-US,en;q=0.9",
	"en-US,en;q=0.8",
	"en-GB,en;q=0.9,en-US;q=0.8",
	"en-US,en;q=0.9,es;q=0.8",
}

// GetRandomUserAgent returns a random user agent string
func (fm *FingerprintMasker) GetRandomUserAgent() string {
	if !fm.config.RotateUserAgent {
		return userAgents[0]
	}
	return userAgents[fm.rng.Intn(len(userAgents))]
}

// GetRandomViewport returns a random viewport size
func (fm *FingerprintMasker) GetRandomViewport() Viewport {
	if !fm.config.RandomizeViewport {
		return Viewport{1920, 1080}
	}

	base := viewports[fm.rng.Intn(len(viewports))]

	// Add small random variation (±5 pixels)
	return Viewport{
		Width:  base.Width + fm.rng.Intn(11) - 5,
		Height: base.Height + fm.rng.Intn(11) - 5,
	}
}

// GetRandomTimezone returns a random timezone
func (fm *FingerprintMasker) GetRandomTimezone() string {
	if !fm.config.RandomizeTimezone {
		return "America/New_York"
	}
	return timezones[fm.rng.Intn(len(timezones))]
}

// GetRandomAcceptLanguage returns a random Accept-Language header
func (fm *FingerprintMasker) GetRandomAcceptLanguage() string {
	return acceptLanguages[fm.rng.Intn(len(acceptLanguages))]
}

// GetWebdriverDisableScript returns JavaScript to disable webdriver detection
func (fm *FingerprintMasker) GetWebdriverDisableScript() string {
	if !fm.config.DisableWebdriverFlag {
		return ""
	}

	return `
		// Remove webdriver property
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined,
		});

		// Override the navigator.plugins to appear non-automated
		Object.defineProperty(navigator, 'plugins', {
			get: () => [
				{ name: 'Chrome PDF Plugin', filename: 'internal-pdf-viewer' },
				{ name: 'Chrome PDF Viewer', filename: 'mhjfbmdgcfjbbpaeojofohoefgiehjai' },
				{ name: 'Native Client', filename: 'internal-nacl-plugin' }
			],
		});

		// Override navigator.languages
		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en'],
		});

		// Remove automation indicators from Chrome
		if (window.chrome) {
			window.chrome.runtime = undefined;
		}

		// Override permissions query
		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({ state: Notification.permission }) :
				originalQuery(parameters)
		);
	`
}

// GetCanvasObfuscationScript returns JavaScript to obfuscate canvas fingerprinting
func (fm *FingerprintMasker) GetCanvasObfuscationScript() string {
	if !fm.config.ObfuscateCanvas {
		return ""
	}

	return `
		// Add noise to canvas fingerprinting
		const originalToDataURL = HTMLCanvasElement.prototype.toDataURL;
		HTMLCanvasElement.prototype.toDataURL = function(type) {
			if (type === 'image/png' || type === undefined) {
				const context = this.getContext('2d');
				if (context) {
					const imageData = context.getImageData(0, 0, this.width, this.height);
					for (let i = 0; i < imageData.data.length; i += 4) {
						// Add minimal noise that won't be visible but changes fingerprint
						imageData.data[i] ^= (Math.random() * 2) | 0;
					}
					context.putImageData(imageData, 0, 0);
				}
			}
			return originalToDataURL.apply(this, arguments);
		};

		// Override getImageData
		const originalGetImageData = CanvasRenderingContext2D.prototype.getImageData;
		CanvasRenderingContext2D.prototype.getImageData = function() {
			const imageData = originalGetImageData.apply(this, arguments);
			for (let i = 0; i < imageData.data.length; i += 4) {
				imageData.data[i] ^= (Math.random() * 2) | 0;
			}
			return imageData;
		};
	`
}

// GetAllMaskingScripts returns all JavaScript for fingerprint masking
func (fm *FingerprintMasker) GetAllMaskingScripts() string {
	scripts := fm.GetWebdriverDisableScript()
	scripts += fm.GetCanvasObfuscationScript()
	return scripts
}

// GetBrowserArgs returns Chrome/Chromium arguments for stealth
func (fm *FingerprintMasker) GetBrowserArgs() []string {
	args := []string{
		"--disable-blink-features=AutomationControlled",
		"--disable-dev-shm-usage",
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-infobars",
		"--disable-extensions",
		"--disable-popup-blocking",
	}

	// Add random window position
	if fm.config.RandomizeViewport {
		x := fm.rng.Intn(100)
		y := fm.rng.Intn(100)
		args = append(args, "--window-position="+string(rune(x))+","+string(rune(y)))
	}

	return args
}
