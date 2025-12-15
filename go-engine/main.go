package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

	"linkedin-automation/auth"
	"linkedin-automation/config"
	"linkedin-automation/database"
	"linkedin-automation/logger"
	"linkedin-automation/messaging"
	"linkedin-automation/search"
	"linkedin-automation/stealth"
)

// Automation represents the main automation controller
type Automation struct {
	config            *config.Config
	db                *database.DB
	logger            *logger.Logger
	browser           *rod.Browser
	page              *rod.Page
	authenticator     *auth.Authenticator
	searchModule      *search.Searcher
	connectionManager *messaging.ConnectionManager
	messageManager    *messaging.MessageManager
	stopChan          chan struct{}
	isRunning         bool
}

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	headless := flag.Bool("headless", true, "Run browser in headless mode")
	dryRun := flag.Bool("dry-run", false, "Run without actually sending requests")
	flag.Parse()

	fmt.Println("==================================================")
	fmt.Println("   LinkedIn Automation Tool v1.0.0")
	fmt.Println("   Go + Rod Browser Automation Engine")
	fmt.Println("==================================================")
	fmt.Println()

	// Load configuration
	fmt.Printf("Loading configuration from %s...\n", *configPath)
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.File)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	log.Info("LinkedIn Automation starting", "version", "1.0.0")

	// Initialize database
	fmt.Printf("Initializing database at %s...\n", cfg.Database.Path)
	db, err := database.New(cfg.Database.Path)
	if err != nil {
		log.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Initialize(); err != nil {
		log.Error("Failed to create database schema", "error", err)
		os.Exit(1)
	}

	// Create automation instance
	auto := &Automation{
		config:   cfg,
		db:       db,
		logger:   log,
		stopChan: make(chan struct{}),
	}

	// Initialize modules
	auto.authenticator = auth.NewAuthenticator(cfg.Credentials, db, log, cfg.Stealth)
	auto.searchModule = search.NewSearcher(cfg.Search, db, log, cfg.Stealth)
	auto.connectionManager = messaging.NewConnectionManager(cfg.Connection, db, log, cfg.Stealth)
	auto.messageManager = messaging.NewMessageManager(cfg.Messaging, db, log, cfg.Stealth)

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived shutdown signal, stopping gracefully...")
		auto.Stop()
	}()

	// Check for dry run mode
	if *dryRun {
		fmt.Println("\n[DRY RUN MODE] - No actual requests will be sent")
		auto.runDryRun()
		return
	}

	// Launch browser
	fmt.Println("\nLaunching browser...")
	err = auto.launchBrowser(*headless)
	if err != nil {
		log.Error("Failed to launch browser", "error", err)
		os.Exit(1)
	}
	defer auto.closeBrowser()

	// Run automation
	err = auto.Run()
	if err != nil {
		log.Error("Automation error", "error", err)
		os.Exit(1)
	}

	fmt.Println("\nAutomation completed successfully!")
}

// launchBrowser starts the Chromium browser
func (a *Automation) launchBrowser(headless bool) error {
	// Create launcher with stealth options
	fm := stealth.NewFingerprintMasker(a.config.Stealth.Fingerprint)

	l := launcher.New().
		Headless(headless).
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-dev-shm-usage").
		Set("no-first-run").
		Set("no-default-browser-check").
		Set("disable-infobars")

	// Set user agent
	userAgent := fm.GetRandomUserAgent()
	l.Set("user-agent", userAgent)

	url, err := l.Launch()
	if err != nil {
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	// Connect to browser
	a.browser = rod.New().ControlURL(url)
	err = a.browser.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}

	// Set viewport
	viewport := fm.GetRandomViewport()
	a.browser.DefaultDevice(devices.Device{
		Title:  "Desktop",
		Screen: devices.Screen{Width: viewport.Width, Height: viewport.Height},
	})

	a.logger.Info("Browser launched", "headless", headless, "userAgent", userAgent[:50]+"...")
	return nil
}

// closeBrowser closes the browser
func (a *Automation) closeBrowser() {
	if a.browser != nil {
		a.browser.Close()
	}
}

// Run executes the main automation workflow
func (a *Automation) Run() error {
	a.isRunning = true
	defer func() { a.isRunning = false }()

	a.logger.Info("Starting automation workflow")

	// Check business hours
	if !a.config.IsBusinessHours() {
		a.logger.Info("Outside business hours, waiting...")
		fmt.Println("\nOutside business hours. Automation will run during configured hours.")
		fmt.Printf("Business hours: %d:00 - %d:00\n",
			a.config.RateLimits.BusinessHoursStart,
			a.config.RateLimits.BusinessHoursEnd)
		return nil
	}

	// Step 1: Authenticate
	fmt.Println("\n[Step 1] Authenticating...")

	// Try session restore first
	page, restored, err := a.authenticator.TrySessionRestore(a.browser)
	if err != nil {
		a.logger.LogError("session restore", err, nil)
	}

	if restored {
		fmt.Println("✓ Session restored from saved cookies")
		a.page = page
	} else {
		// Perform fresh login
		page, result, err := a.authenticator.Login(a.browser)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		if result.SecurityChallenge {
			fmt.Printf("\n⚠ Security challenge detected: %s\n", result.ChallengeType)
			fmt.Println(result.ErrorMessage)
			fmt.Println("\nPlease complete the verification manually and restart the automation.")
			return nil
		}

		if !result.Success {
			return fmt.Errorf("login failed: %s", result.ErrorMessage)
		}

		fmt.Println("✓ Login successful")
		a.page = page
	}

	// Step 2: Search for profiles
	fmt.Println("\n[Step 2] Searching for profiles...")
	searchResult, err := a.searchModule.Search(a.page)
	if err != nil {
		a.logger.LogError("search", err, nil)
		fmt.Printf("⚠ Search error: %v\n", err)
	} else {
		fmt.Printf("✓ Found %d profiles (%d unique, %d duplicates)\n",
			searchResult.TotalFound,
			len(searchResult.Profiles),
			searchResult.Duplicates)
	}

	// Step 3: Send connection requests
	fmt.Println("\n[Step 3] Sending connection requests...")

	canSend, remaining, _ := a.connectionManager.CanSendMoreToday()
	if !canSend {
		fmt.Println("⚠ Daily connection limit reached")
	} else {
		fmt.Printf("Remaining connections today: %d\n", remaining)

		connectionsSent := 0
		for i, profile := range searchResult.Profiles {
			select {
			case <-a.stopChan:
				fmt.Println("\nStopping...")
				return nil
			default:
			}

			// Check if we can send more
			canSend, _, _ = a.connectionManager.CanSendMoreToday()
			if !canSend {
				fmt.Println("\n⚠ Daily limit reached, stopping connection requests")
				break
			}

			req := &messaging.ConnectionRequest{
				ProfileURL:  profile.ProfileURL,
				FirstName:   profile.FirstName,
				LastName:    profile.LastName,
				JobTitle:    profile.JobTitle,
				Company:     profile.Company,
				TemplateIdx: i,
			}

			result, err := a.connectionManager.SendConnectionRequest(a.page, req)
			if err != nil {
				a.logger.LogError("connection request", err, map[string]interface{}{"profile": profile.ProfileURL})
				continue
			}

			if result.Success {
				connectionsSent++
				fmt.Printf("  ✓ Sent to %s %s\n", profile.FirstName, profile.LastName)
			} else {
				fmt.Printf("  ⚠ Failed: %s\n", result.ErrorMessage)
			}

			// Rate limiting delay
			a.waitBetweenActions()
		}

		fmt.Printf("\n✓ Sent %d connection requests\n", connectionsSent)
	}

	// Step 4: Check for accepted connections and send follow-ups
	fmt.Println("\n[Step 4] Checking accepted connections...")
	accepted, err := a.messageManager.DetectAcceptedConnections(a.page)
	if err != nil {
		a.logger.LogError("detect accepted", err, nil)
	} else {
		fmt.Printf("✓ Found %d newly accepted connections\n", len(accepted))
	}

	// Send follow-up messages
	if len(a.config.Messaging.Templates) > 0 {
		needFollowUp, _ := a.messageManager.GetConnectionsNeedingFollowUp()
		if len(needFollowUp) > 0 {
			fmt.Printf("\n[Step 5] Sending follow-up messages to %d connections...\n", len(needFollowUp))

			messagesSent := 0
			for i, conn := range needFollowUp {
				select {
				case <-a.stopChan:
					return nil
				default:
				}

				canSend, _, _ := a.messageManager.CanSendMoreMessagesToday()
				if !canSend {
					fmt.Println("\n⚠ Daily message limit reached")
					break
				}

				req := &messaging.MessageRequest{
					ConnectionID: conn.ID,
					ProfileURL:   conn.ProfileURL,
					FirstName:    conn.FirstName,
					LastName:     conn.LastName,
					JobTitle:     conn.JobTitle,
					Company:      conn.Company,
					TemplateIdx:  i,
				}

				result, err := a.messageManager.SendMessage(a.page, req)
				if err != nil {
					a.logger.LogError("send message", err, nil)
					continue
				}

				if result.Success {
					messagesSent++
					fmt.Printf("  ✓ Message sent to %s %s\n", conn.FirstName, conn.LastName)
				}

				a.waitBetweenActions()
			}

			fmt.Printf("\n✓ Sent %d follow-up messages\n", messagesSent)
		}
	}

	// Print summary
	a.printSummary()

	return nil
}

// waitBetweenActions waits with randomized delay between actions
func (a *Automation) waitBetweenActions() {
	timing := stealth.NewTimingController(a.config.Stealth.Timing)
	delay := timing.GetRandomizedDelay(
		a.config.RateLimits.MinActionDelayMs,
		a.config.RateLimits.MaxActionDelayMs,
	)
	time.Sleep(delay)
}

// Stop signals the automation to stop
func (a *Automation) Stop() {
	if a.isRunning {
		close(a.stopChan)
	}
}

// runDryRun executes a dry run showing what would be done
func (a *Automation) runDryRun() {
	fmt.Println("\n--- Dry Run Configuration ---")
	fmt.Printf("Email: %s\n", maskEmail(a.config.Credentials.Email))
	fmt.Printf("Search Keywords: %v\n", a.config.Search.Keywords)
	fmt.Printf("Job Titles: %v\n", a.config.Search.JobTitles)
	fmt.Printf("Locations: %v\n", a.config.Search.Locations)
	fmt.Printf("Max Pages: %d\n", a.config.Search.MaxPages)
	fmt.Printf("Daily Connection Limit: %d\n", a.config.Connection.DailyLimit)
	fmt.Printf("Daily Message Limit: %d\n", a.config.Messaging.DailyLimit)
	fmt.Printf("Business Hours: %d:00 - %d:00\n",
		a.config.RateLimits.BusinessHoursStart,
		a.config.RateLimits.BusinessHoursEnd)

	fmt.Println("\n--- Stealth Configuration ---")
	fmt.Printf("Bézier Curves: %v\n", a.config.Stealth.Bezier.Enabled)
	fmt.Printf("Typing Delay: %d-%dms\n",
		a.config.Stealth.Timing.TypingMinDelayMs,
		a.config.Stealth.Timing.TypingMaxDelayMs)
	fmt.Printf("Typo Probability: %.2f\n", a.config.Stealth.Timing.TypoProbability)
	fmt.Printf("User Agent Rotation: %v\n", a.config.Stealth.Fingerprint.RotateUserAgent)
	fmt.Printf("Viewport Randomization: %v\n", a.config.Stealth.Fingerprint.RandomizeViewport)

	fmt.Println("\n--- Connection Templates ---")
	for i, t := range a.config.Connection.Templates {
		fmt.Printf("%d. %s\n", i+1, t)
	}

	fmt.Println("\n--- Message Templates ---")
	for i, t := range a.config.Messaging.Templates {
		fmt.Printf("%d. %s\n", i+1, t)
	}

	fmt.Println("\n✓ Configuration validated successfully")
	fmt.Println("Remove --dry-run flag to start actual automation")
}

// printSummary prints the automation summary
func (a *Automation) printSummary() {
	activity, _ := a.db.GetOrCreateDailyActivity()

	fmt.Println("\n==================================================")
	fmt.Println("   Automation Summary")
	fmt.Println("==================================================")
	fmt.Printf("Connections sent today: %d / %d\n", activity.ConnectionsSent, a.config.Connection.DailyLimit)
	fmt.Printf("Messages sent today: %d / %d\n", activity.MessagesSent, a.config.Messaging.DailyLimit)

	if a.config.IsBusinessHours() {
		fmt.Println("\nNext run: Whenever you start the automation again")
	} else {
		fmt.Printf("\nNote: Currently outside business hours (%d:00 - %d:00)\n",
			a.config.RateLimits.BusinessHoursStart,
			a.config.RateLimits.BusinessHoursEnd)
	}
}

// maskEmail masks email for logging
func maskEmail(email string) string {
	if len(email) < 5 {
		return "***"
	}
	return email[:3] + "***" + email[len(email)-4:]
}

// devices package stub for viewport
var devices = struct {
	Device func(title string, width, height int) Device
}{
	Device: func(title string, width, height int) Device {
		return Device{Title: title, Screen: Screen{Width: width, Height: height}}
	},
}

type Device struct {
	Title  string
	Screen Screen
}

type Screen struct {
	Width  int
	Height int
}
