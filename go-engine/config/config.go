package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all configuration for the automation
type Config struct {
	Credentials CredentialsConfig `mapstructure:"credentials"`
	Search      SearchConfig      `mapstructure:"search"`
	Connection  ConnectionConfig  `mapstructure:"connection"`
	Messaging   MessagingConfig   `mapstructure:"messaging"`
	RateLimits  RateLimitsConfig  `mapstructure:"rate_limits"`
	Stealth     StealthConfig     `mapstructure:"stealth"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Logging     LoggingConfig     `mapstructure:"logging"`
	API         APIConfig         `mapstructure:"api"`
}

type CredentialsConfig struct {
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

type SearchConfig struct {
	JobTitles []string `mapstructure:"job_titles"`
	Companies []string `mapstructure:"companies"`
	Locations []string `mapstructure:"locations"`
	Keywords  []string `mapstructure:"keywords"`
	MaxPages  int      `mapstructure:"max_pages"`
}

type ConnectionConfig struct {
	DailyLimit    int      `mapstructure:"daily_limit"`
	Templates     []string `mapstructure:"templates"`
	MaxNoteLength int      `mapstructure:"max_note_length"`
}

type MessagingConfig struct {
	DailyLimit      int      `mapstructure:"daily_limit"`
	MinDelayMinutes int      `mapstructure:"min_delay_minutes"`
	MaxDelayMinutes int      `mapstructure:"max_delay_minutes"`
	Templates       []string `mapstructure:"templates"`
}

type RateLimitsConfig struct {
	MinActionDelayMs        int  `mapstructure:"min_action_delay_ms"`
	MaxActionDelayMs        int  `mapstructure:"max_action_delay_ms"`
	BusinessHoursStart      int  `mapstructure:"business_hours_start"`
	BusinessHoursEnd        int  `mapstructure:"business_hours_end"`
	SkipWeekends            bool `mapstructure:"skip_weekends"`
	CooldownAfterBulkSecs   int  `mapstructure:"cooldown_after_bulk_actions"`
}

type StealthConfig struct {
	Bezier      BezierConfig      `mapstructure:"bezier"`
	Timing      TimingConfig      `mapstructure:"timing"`
	Fingerprint FingerprintConfig `mapstructure:"fingerprint"`
	Scrolling   ScrollingConfig   `mapstructure:"scrolling"`
	Mouse       MouseConfig       `mapstructure:"mouse"`
	Scheduling  SchedulingConfig  `mapstructure:"scheduling"`
	Headers     HeadersConfig     `mapstructure:"headers"`
	Network     NetworkConfig     `mapstructure:"network"`
}

type BezierConfig struct {
	Enabled              bool    `mapstructure:"enabled"`
	OvershootProbability float64 `mapstructure:"overshoot_probability"`
	MinSteps             int     `mapstructure:"min_steps"`
	MaxSteps             int     `mapstructure:"max_steps"`
}

type TimingConfig struct {
	TypingMinDelayMs int     `mapstructure:"typing_min_delay_ms"`
	TypingMaxDelayMs int     `mapstructure:"typing_max_delay_ms"`
	TypoProbability  float64 `mapstructure:"typo_probability"`
	ThinkTimeMinMs   int     `mapstructure:"think_time_min_ms"`
	ThinkTimeMaxMs   int     `mapstructure:"think_time_max_ms"`
}

type FingerprintConfig struct {
	RotateUserAgent     bool `mapstructure:"rotate_user_agent"`
	RandomizeViewport   bool `mapstructure:"randomize_viewport"`
	DisableWebdriverFlag bool `mapstructure:"disable_webdriver_flag"`
	RandomizeTimezone   bool `mapstructure:"randomize_timezone"`
	ObfuscateCanvas     bool `mapstructure:"obfuscate_canvas"`
}

type ScrollingConfig struct {
	Enabled               bool    `mapstructure:"enabled"`
	MinSpeed              int     `mapstructure:"min_speed"`
	MaxSpeed              int     `mapstructure:"max_speed"`
	ScrollBackProbability float64 `mapstructure:"scroll_back_probability"`
}

type MouseConfig struct {
	HoverBeforeClick   bool `mapstructure:"hover_before_click"`
	RandomMovement     bool `mapstructure:"random_movement"`
	HoverDurationMinMs int  `mapstructure:"hover_duration_min_ms"`
	HoverDurationMaxMs int  `mapstructure:"hover_duration_max_ms"`
}

type SchedulingConfig struct {
	RespectBusinessHours bool `mapstructure:"respect_business_hours"`
	IncludeBreaks        bool `mapstructure:"include_breaks"`
	LunchBreakStart      int  `mapstructure:"lunch_break_start"`
	LunchBreakEnd        int  `mapstructure:"lunch_break_end"`
}

type HeadersConfig struct {
	Randomize          bool `mapstructure:"randomize"`
	VaryAcceptLanguage bool `mapstructure:"vary_accept_language"`
}

type NetworkConfig struct {
	SimulateLatency bool `mapstructure:"simulate_latency"`
	LatencyMinMs    int  `mapstructure:"latency_min_ms"`
	LatencyMaxMs    int  `mapstructure:"latency_max_ms"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	File   string `mapstructure:"file"`
	Format string `mapstructure:"format"`
}

type APIConfig struct {
	BackendURL  string `mapstructure:"backend_url"`
	SyncEnabled bool   `mapstructure:"sync_enabled"`
}

// Load loads configuration from file and environment
func Load(configPath string) (*Config, error) {
	// Load .env file if present
	godotenv.Load()

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set defaults
	v.SetDefault("search.max_pages", 5)
	v.SetDefault("connection.daily_limit", 50)
	v.SetDefault("connection.max_note_length", 300)
	v.SetDefault("messaging.daily_limit", 100)
	v.SetDefault("messaging.min_delay_minutes", 5)
	v.SetDefault("messaging.max_delay_minutes", 15)
	v.SetDefault("rate_limits.min_action_delay_ms", 5000)
	v.SetDefault("rate_limits.max_action_delay_ms", 15000)
	v.SetDefault("rate_limits.business_hours_start", 9)
	v.SetDefault("rate_limits.business_hours_end", 18)
	v.SetDefault("stealth.bezier.enabled", true)
	v.SetDefault("stealth.bezier.overshoot_probability", 0.15)
	v.SetDefault("stealth.bezier.min_steps", 20)
	v.SetDefault("stealth.bezier.max_steps", 50)
	v.SetDefault("stealth.timing.typing_min_delay_ms", 50)
	v.SetDefault("stealth.timing.typing_max_delay_ms", 150)
	v.SetDefault("stealth.timing.typo_probability", 0.05)
	v.SetDefault("database.path", "./linkedin_automation.db")
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Override with environment variables
	if email := os.Getenv("LINKEDIN_EMAIL"); email != "" {
		cfg.Credentials.Email = email
	}
	if password := os.Getenv("LINKEDIN_PASSWORD"); password != "" {
		cfg.Credentials.Password = password
	}

	return &cfg, nil
}

// IsBusinessHours checks if current time is within business hours
func (c *Config) IsBusinessHours() bool {
	if !c.Stealth.Scheduling.RespectBusinessHours {
		return true
	}

	now := time.Now()
	hour := now.Hour()
	weekday := now.Weekday()

	// Check weekends
	if c.RateLimits.SkipWeekends && (weekday == time.Saturday || weekday == time.Sunday) {
		return false
	}

	// Check business hours
	if hour < c.RateLimits.BusinessHoursStart || hour >= c.RateLimits.BusinessHoursEnd {
		return false
	}

	// Check lunch break
	if c.Stealth.Scheduling.IncludeBreaks {
		if hour >= c.Stealth.Scheduling.LunchBreakStart && hour < c.Stealth.Scheduling.LunchBreakEnd {
			return false
		}
	}

	return true
}
