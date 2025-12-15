# LinkedIn Automation Tool

A sophisticated Go-based LinkedIn automation tool with a React dashboard for autonomous connection requests, personalized messaging, and advanced anti-detection capabilities.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Technology Stack](#technology-stack)
4. [Features](#features)
5. [Anti-Bot Detection System](#anti-bot-detection-system)
6. [Algorithms](#algorithms)
7. [API Reference](#api-reference)
8. [Configuration](#configuration)
9. [Database Schema](#database-schema)
10. [Installation & Setup](#installation--setup)
11. [Usage](#usage)
12. [Security Considerations](#security-considerations)

---

## Overview

This LinkedIn Automation Tool is a comprehensive technical proof-of-concept demonstrating advanced browser automation, human-like behavior simulation, and sophisticated anti-bot detection techniques. The system enables:

- **Autonomous Connection Requests**: Automatically send personalized connection requests
- **Follow-up Messaging**: Send messages to accepted connections
- **Profile Discovery**: Search and extract LinkedIn profiles based on criteria
- **Anti-Detection**: 10+ stealth techniques to evade bot detection
- **State Persistence**: SQLite-based tracking and resumption

### Target Use Cases
- Lead generation and networking automation
- Technical demonstration of browser automation patterns
- Educational resource for anti-detection strategies
- Recruitment outreach automation

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              LINKEDIN AUTOMATION                              │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────┐         ┌─────────────────────────────────────────┐ │
│  │    React Frontend   │◄───────►│         FastAPI Backend (Python)       │ │
│  │   (Dashboard UI)    │   HTTP  │        - REST API Endpoints            │ │
│  │                     │         │        - Configuration Management      │ │
│  │  • Dashboard        │         │        - State Synchronization         │ │
│  │  • Credentials      │         │        - Activity Logging              │ │
│  │  • Search Criteria  │         └─────────────────┬───────────────────────┘ │
│  │  • Templates        │                           │                         │
│  │  • Connections      │                           │ MongoDB                 │
│  │  • Messages         │                           │                         │
│  │  • Stealth Config   │                           ▼                         │
│  │  • Rate Limits      │         ┌─────────────────────────────────────────┐ │
│  │  • Activity Logs    │         │            MongoDB Database             │ │
│  └────────────────────┘         │  • connections    • messages            │ │
│                                 │  • credentials    • templates           │ │
│                                 │  • search_criteria • stealth_config    │ │
│                                 │  • rate_limits    • activity_logs      │ │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                     GO AUTOMATION ENGINE                               │ │
│  │  ┌──────────────────────────────────────────────────────────────────┐ │ │
│  │  │                        MAIN APPLICATION                           │ │ │
│  │  │  • Workflow orchestration    • Signal handling                   │ │ │
│  │  │  • Browser lifecycle         • State management                  │ │ │
│  │  └──────────────────────────────────────────────────────────────────┘ │ │
│  │                                                                        │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐ │ │
│  │  │   AUTH   │ │  SEARCH  │ │MESSAGING │ │ STEALTH  │ │   DATABASE   │ │ │
│  │  │          │ │          │ │          │ │          │ │              │ │ │
│  │  │• Login   │ │• Query   │ │• Connect │ │• Bézier  │ │• Connections │ │ │
│  │  │• Session │ │• Extract │ │• Request │ │• Timing  │ │• Messages    │ │ │
│  │  │• Cookies │ │• Paginate│ │• Message │ │• Print   │ │• Daily Stats │ │ │
│  │  │• 2FA Det │ │• Dedup   │ │• Template│ │• Scroll  │ │• Cookies     │ │ │
│  │  └──────────┘ └──────────┘ └──────────┘ │• Typing  │ └──────────────┘ │ │
│  │                                          │• Mouse   │                  │ │
│  │                                          │• Headers │ ┌──────────────┐ │ │
│  │                                          │• Network │ │    UTILS     │ │ │
│  │                                          └──────────┘ │• Retry       │ │ │
│  │                                                       │• Validation  │ │ │
│  │                    Rod Browser (CDP)                  └──────────────┘ │ │
│  │                         │                                              │ │
│  │                         ▼                                              │ │
│  │              ┌─────────────────────┐        ┌───────────────────────┐ │ │
│  │              │ Chromium (Headless) │◄──────►│   SQLite Database     │ │ │
│  │              └─────────────────────┘        │  (State Persistence)  │ │ │
│  │                         │                   └───────────────────────┘ │ │
│  │                         ▼                                              │ │
│  │              ┌─────────────────────┐                                  │ │
│  │              │     LinkedIn.com    │                                  │ │
│  │              └─────────────────────┘                                  │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## Technology Stack

### Go Automation Engine
| Component | Technology | Purpose |
|-----------|------------|---------|
| Language | Go 1.21+ | Compiled performance, concurrency support |
| Browser Automation | Rod (Chrome DevTools Protocol) | High-level browser control |
| Browser | Chromium (Headless) | Full web platform support |
| Database | SQLite3 | Lightweight state persistence |
| Configuration | Viper + godotenv | YAML/JSON config, env variables |
| Logging | Go's slog package | Structured, leveled logging |

### Web Dashboard
| Component | Technology | Purpose |
|-----------|------------|---------|
| Frontend | React 18 + Tailwind CSS | Modern responsive UI |
| Backend | FastAPI (Python) | REST API endpoints |
| Database | MongoDB | Configuration & log storage |
| Icons | Lucide React | UI iconography |
| Notifications | Sonner | Toast notifications |

### Key Dependencies
```go
// Go Engine
github.com/go-rod/rod      // Browser automation
github.com/mattn/go-sqlite3 // Database operations
github.com/spf13/viper     // Configuration
github.com/joho/godotenv   // Environment variables
```

```javascript
// Frontend
axios          // HTTP client
sonner         // Notifications
lucide-react   // Icons
tailwindcss    // Styling
```

---

## Features

### 1. Authentication System

**Functionality:**
- Automated LinkedIn login with credential management
- Session cookie persistence for seamless resumption
- Security checkpoint detection (2FA, CAPTCHA)
- Encrypted credential storage

**Implementation Details:**
- Credentials loaded from environment variables (`LINKEDIN_EMAIL`, `LINKEDIN_PASSWORD`)
- Realistic typing simulation for credential entry
- Bézier curve mouse movement for button clicks
- Automatic cookie extraction and injection

```go
// Login flow
1. Navigate to linkedin.com/login
2. Wait for page load with randomized delay
3. Enter email with realistic typing
4. Enter password with typo simulation
5. Click login with Bézier mouse movement
6. Detect security challenges
7. Verify successful login
8. Save session cookies
```

### 2. Search & Targeting System

**Functionality:**
- Search users by job title, company, location, keywords
- Profile URL extraction with CSS selectors
- Pagination handling for large result sets
- Duplicate detection across sessions

**Search URL Construction:**
```
https://www.linkedin.com/search/results/people/?keywords={terms}&geoUrn={location}
```

**Profile Extraction:**
- Selector: `a[href*="/in/"]`
- Extracts: Profile URL, Name, Job Title, Company, Location
- Normalizes URLs by removing tracking parameters

### 3. Connection Request System

**Functionality:**
- Navigate to profiles with error handling
- Locate and click Connect button
- Send personalized notes with template variables
- Track sent requests with daily limits

**Template Variables:**
| Variable | Description |
|----------|-------------|
| `{{firstName}}` | Target's first name |
| `{{lastName}}` | Target's last name |
| `{{jobTitle}}` | Target's current position |
| `{{company}}` | Target's current company |
| `{{location}}` | Target's location |

**Example Templates:**
```
"Hi {{firstName}}, I noticed your work at {{company}} and would love to connect!"
"Hello {{firstName}}, I'm impressed by your experience as a {{jobTitle}}. Let's connect!"
```

### 4. Messaging System

**Functionality:**
- Detect newly accepted connections
- Send follow-up messages automatically
- Support template personalization
- Track message delivery status

**Message Flow:**
```
1. Query pending connections from database
2. Navigate to My Network to verify accepted
3. Open messaging interface
4. Type message with realistic behavior
5. Send and record in database
```

### 5. Dashboard Features

**Dashboard Sections:**
| Section | Purpose |
|---------|---------|
| Dashboard | Real-time stats, automation controls |
| Credentials | LinkedIn login configuration |
| Search Criteria | Define targeting filters |
| Templates | Connection & follow-up messages |
| Connections | View all connection requests |
| Messages | View sent follow-up messages |
| Anti-Bot Settings | Configure stealth techniques |
| Rate Limits | Daily limits and timing |
| Activity Logs | View automation history |

---

## Anti-Bot Detection System

The system implements **10+ stealth techniques** to evade LinkedIn's bot detection:

### Mandatory Techniques (3)

#### 1. Bézier Curve Mouse Movement
Simulates natural human mouse paths using cubic Bézier curves.

**Mathematical Formula:**
```
P(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃

Where:
- P₀ = Start point
- P₁, P₂ = Control points (randomized)
- P₃ = End point
- t = Parameter [0, 1]
```

**Features:**
- Variable velocity with acceleration/deceleration
- Natural overshoot and correction
- Randomized control points
- 50-100 discrete movement steps
- Fitts's Law for movement duration

#### 2. Randomized Timing Patterns
Human "think time" simulation using statistical distributions.

**Delay Types:**
| Action | Distribution | Parameters |
|--------|--------------|------------|
| Action-level | Normal | μ=500ms, σ=200ms |
| Keystroke | Normal | μ=75ms, σ=30ms |
| Page load | Normal | μ=2000ms, σ=500ms |
| Think time | Uniform | 2000-5000ms |

**Implementation:**
```go
// Box-Muller transform for normal distribution
z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
return mu + sigma*z
```

#### 3. Browser Fingerprint Masking
Counters fingerprinting-based detection.

**Masked Properties:**
- User Agent rotation (8+ browser variants)
- Viewport randomization (8 common resolutions)
- WebDriver flag removal
- Navigator plugins emulation
- Canvas fingerprint obfuscation
- Timezone randomization

**JavaScript Injection:**
```javascript
Object.defineProperty(navigator, 'webdriver', { get: () => undefined });
Object.defineProperty(navigator, 'plugins', { get: () => [...] });
window.chrome.runtime = undefined;
```

### Additional Techniques (7)

#### 4. Random Scrolling Behavior
- Variable scroll speeds (50-300 px/operation)
- Natural acceleration/deceleration
- Occasional scroll-back movements
- Random pauses during scrolling

#### 5. Realistic Typing Simulation
- Variable keystroke intervals (50-150ms)
- Occasional typos with corrections (5% probability)
- Burst typing patterns (3-8 characters)
- Extra delay before capital letters

**Adjacent Key Map:**
```go
adjacentKeys = map[rune][]rune{
    'a': {'s', 'q', 'w', 'z'},
    'b': {'v', 'n', 'g', 'h'},
    // ... full QWERTY layout
}
```

#### 6. Mouse Hovering & Movement
- Pre-click hover events
- Natural cursor wandering
- Post-click drift movement
- Hover-to-click sequences

#### 7. Activity Scheduling
- Business hours operation (default 9 AM - 6 PM)
- Lunch break simulation
- Weekend skipping option
- Daily start time variation (±30 minutes)

#### 8. Rate Limiting & Throttling
- Token Bucket algorithm for rate control
- Daily connection limits (default 50)
- Message spacing (5-15 minute intervals)
- Cooldown periods after bulk activity

#### 9. Request Header Randomization
- Accept-Language variation
- Accept-Encoding combinations
- Referer header modification
- DNT header presence variation

#### 10. Network Behavior Simulation
- Variable latency simulation (50-200ms)
- Request retry patterns
- Realistic caching behavior

---

## Algorithms

### 1. Token Bucket Algorithm (Rate Limiting)

Controls the frequency of actions to prevent detection.

```go
type TokenBucket struct {
    capacity    int           // Maximum tokens (daily limit)
    tokens      int           // Current available tokens
    refillRate  time.Duration // Time between refills
    lastRefill  time.Time
}

func (tb *TokenBucket) CanPerformAction() bool {
    tb.refill()
    if tb.tokens > 0 {
        tb.tokens--
        return true
    }
    return false
}
```

**Parameters:**
| Setting | Default | Description |
|---------|---------|-------------|
| Capacity | 50 | Daily connection limit |
| Refill Rate | 24 hours | Full refill period |
| Action Cost | 1 | Tokens per action |

### 2. Exponential Backoff (Error Recovery)

Handles transient failures with increasing retry delays.

```go
func RetryWithBackoff(cfg RetryConfig, fn func() error) error {
    for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        // Calculate delay: base * 2^attempt + jitter
        delay := cfg.InitialDelay * time.Duration(math.Pow(2, float64(attempt)))
        jitter := float64(delay) * cfg.JitterPercent * (rand.Float64()*2 - 1)
        time.Sleep(time.Duration(float64(delay) + jitter))
    }
    return fmt.Errorf("max retries exceeded")
}
```

**Parameters:**
| Setting | Default | Description |
|---------|---------|-------------|
| MaxRetries | 5 | Maximum retry attempts |
| InitialDelay | 1s | First retry delay |
| MaxDelay | 30s | Maximum delay cap |
| JitterPercent | 50% | Random variation |

### 3. Fitts's Law (Mouse Movement Duration)

Determines how long mouse movement should take based on distance and target size.

```
MT = a + b × log₂(A/W + 1)

Where:
- MT = Movement Time
- A = Amplitude (distance to target)
- W = Width (target size)
- a, b = Constants (empirically derived)
```

**Implementation:**
```go
func (bm *BezierMouse) calculateSteps(distance float64) int {
    baseSteps := int(math.Log2(distance+1) * 10)
    steps := baseSteps + bm.rng.Intn(10) - 5
    return clamp(steps, bm.config.MinSteps, bm.config.MaxSteps)
}
```

### 4. Box-Muller Transform (Normal Distribution)

Generates normally distributed random numbers for realistic timing.

```go
func normalRandom(mu, sigma float64) float64 {
    u1 := rand.Float64()
    u2 := rand.Float64()
    
    z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
    return mu + sigma*z
}
```

### 5. Pagination Algorithm

Handles multi-page search result extraction.

```go
func (s *Searcher) Search(page *rod.Page) (*SearchResult, error) {
    for pageNum := 1; pageNum <= s.config.MaxPages; pageNum++ {
        // Extract profiles
        profiles := s.extractProfiles(page)
        
        // Filter duplicates
        for _, profile := range profiles {
            if !s.db.IsProfileProcessed(profile.ProfileURL) {
                result.Profiles = append(result.Profiles, profile)
            }
        }
        
        // Navigate to next page
        if !s.goToNextPage(page) {
            break
        }
        
        // Random delay between pages
        time.Sleep(s.timing.GetThinkTime())
    }
    return result, nil
}
```

### 6. Duplicate Detection Algorithm

Prevents processing the same profile multiple times.

```go
func (s *Searcher) isDuplicate(profileURL string) bool {
    // Check in-memory set (current session)
    if s.processedSet[profileURL] {
        return true
    }
    
    // Check SQLite database (historical)
    processed, _ := s.db.IsProfileProcessed(profileURL)
    return processed
}
```

---

## API Reference

### Dashboard API Endpoints

#### Credentials
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/credentials` | Save LinkedIn credentials |
| GET | `/api/credentials` | Get credentials status |
| DELETE | `/api/credentials` | Delete stored credentials |

#### Search Criteria
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/search-criteria` | Create search criteria |
| GET | `/api/search-criteria` | List all criteria |
| PUT | `/api/search-criteria/{id}` | Update criteria |
| DELETE | `/api/search-criteria/{id}` | Delete criteria |

#### Templates
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/templates` | Create message template |
| GET | `/api/templates` | List templates |
| PUT | `/api/templates/{id}` | Update template |
| DELETE | `/api/templates/{id}` | Delete template |

#### Connections & Messages
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/connections` | List connections |
| GET | `/api/connections/{id}` | Get connection details |
| GET | `/api/messages` | List sent messages |

#### Automation Control
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/automation/start` | Start automation |
| POST | `/api/automation/stop` | Stop automation |
| POST | `/api/automation/pause` | Pause automation |
| POST | `/api/automation/resume` | Resume automation |

#### Configuration
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/rate-limits` | Get rate limit config |
| PUT | `/api/rate-limits` | Update rate limits |
| GET | `/api/stealth-config` | Get stealth settings |
| PUT | `/api/stealth-config` | Update stealth settings |
| GET | `/api/dashboard/stats` | Get dashboard statistics |
| GET | `/api/activity-logs` | Get activity history |
| GET | `/api/export-config` | Export config for Go engine |

---

## Configuration

### Go Engine Configuration (config.yaml)

```yaml
# LinkedIn Credentials
credentials:
  email: ""           # Set via LINKEDIN_EMAIL
  password: ""        # Set via LINKEDIN_PASSWORD

# Search Targeting
search:
  job_titles:
    - "Software Engineer"
    - "Developer"
  companies: []
  locations:
    - "San Francisco Bay Area"
  keywords:
    - "hiring"
  max_pages: 5

# Connection Settings
connection:
  daily_limit: 50
  templates:
    - "Hi {{firstName}}, I noticed your work at {{company}}..."
  max_note_length: 300

# Messaging Settings
messaging:
  daily_limit: 100
  min_delay_minutes: 5
  max_delay_minutes: 15
  templates:
    - "Thanks for connecting, {{firstName}}!"

# Rate Limits
rate_limits:
  min_action_delay_ms: 5000
  max_action_delay_ms: 15000
  business_hours_start: 9
  business_hours_end: 18
  skip_weekends: true
  cooldown_after_bulk_actions: 300

# Stealth Configuration
stealth:
  bezier:
    enabled: true
    overshoot_probability: 0.15
    min_steps: 20
    max_steps: 50
    
  timing:
    typing_min_delay_ms: 50
    typing_max_delay_ms: 150
    typo_probability: 0.05
    think_time_min_ms: 2000
    think_time_max_ms: 5000
    
  fingerprint:
    rotate_user_agent: true
    randomize_viewport: true
    disable_webdriver_flag: true
    randomize_timezone: true
    obfuscate_canvas: true
    
  scrolling:
    enabled: true
    min_speed: 50
    max_speed: 300
    scroll_back_probability: 0.1
    
  mouse:
    hover_before_click: true
    random_movement: true
    hover_duration_min_ms: 100
    hover_duration_max_ms: 500
    
  scheduling:
    respect_business_hours: true
    include_breaks: true
    lunch_break_start: 12
    lunch_break_end: 13
    
  headers:
    randomize: true
    vary_accept_language: true
    
  network:
    simulate_latency: true
    latency_min_ms: 50
    latency_max_ms: 200

# Database
database:
  path: "./linkedin_automation.db"

# Logging
logging:
  level: "info"      # debug, info, warn, error
  file: "./automation.log"
  format: "json"     # json, text
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `LINKEDIN_EMAIL` | LinkedIn login email | Yes |
| `LINKEDIN_PASSWORD` | LinkedIn password | Yes |
| `MONGO_URL` | MongoDB connection string | Yes (Dashboard) |
| `REACT_APP_BACKEND_URL` | Backend API URL | Yes (Frontend) |

---

## Database Schema

### SQLite (Go Engine)

```sql
-- Connection tracking
CREATE TABLE connections (
    id TEXT PRIMARY KEY,
    profile_url TEXT NOT NULL UNIQUE,
    first_name TEXT,
    last_name TEXT,
    job_title TEXT,
    company TEXT,
    location TEXT,
    note_sent TEXT,
    status TEXT CHECK(status IN ('pending', 'accepted', 'declined', 'failed')),
    search_criteria_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    accepted_at DATETIME
);

-- Message history
CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    connection_id TEXT NOT NULL,
    content TEXT NOT NULL,
    template_id TEXT,
    status TEXT CHECK(status IN ('sent', 'delivered', 'read', 'failed')),
    sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (connection_id) REFERENCES connections(id)
);

-- Daily activity tracking
CREATE TABLE daily_activity (
    id TEXT PRIMARY KEY,
    date TEXT NOT NULL UNIQUE,
    connections_sent INTEGER DEFAULT 0,
    messages_sent INTEGER DEFAULT 0,
    last_connection_at DATETIME,
    last_message_at DATETIME
);

-- Session cookies
CREATE TABLE session_cookies (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    domain TEXT,
    path TEXT,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Processed profiles (deduplication)
CREATE TABLE processed_profiles (
    profile_url TEXT PRIMARY KEY,
    processed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### MongoDB (Dashboard)

| Collection | Purpose |
|------------|---------|
| `credentials` | LinkedIn login credentials |
| `search_criteria` | Targeting filters |
| `templates` | Message templates |
| `connections` | Connection request records |
| `messages` | Sent message records |
| `rate_limits` | Rate limiting configuration |
| `stealth_config` | Anti-bot settings |
| `activity_logs` | Automation activity history |
| `automation_state` | Current automation status |

---

## Installation & Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- Python 3.9+
- MongoDB 4.4+
- Chromium browser

### Backend Setup

```bash
# Install Python dependencies
cd backend
pip install -r requirements.txt

# Set environment variables
export MONGO_URL="mongodb://localhost:27017"
export DB_NAME="linkedin_automation"

# Start the server
uvicorn server:app --host 0.0.0.0 --port 8001
```

### Frontend Setup

```bash
# Install dependencies
cd frontend
yarn install

# Set environment variables
echo "REACT_APP_BACKEND_URL=http://localhost:8001" > .env

# Start development server
yarn start
```

### Go Engine Setup

```bash
# Navigate to Go engine
cd go-engine

# Install dependencies
go mod download

# Set credentials
export LINKEDIN_EMAIL="your-email@example.com"
export LINKEDIN_PASSWORD="your-password"

# Run automation
go run main.go --config config.yaml

# Or with dry run (no actual requests)
go run main.go --config config.yaml --dry-run
```

### Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | config.yaml | Configuration file path |
| `--headless` | true | Run browser in headless mode |
| `--dry-run` | false | Validate config without sending requests |

---

## Usage

### Dashboard Workflow

1. **Configure Credentials**
   - Navigate to Credentials section
   - Enter LinkedIn email and password
   - Save credentials

2. **Define Search Criteria**
   - Go to Search Criteria section
   - Add job titles, companies, locations, keywords
   - Save targeting filters

3. **Create Message Templates**
   - Go to Templates section
   - Create connection request templates
   - Create follow-up message templates
   - Use variables: `{{firstName}}`, `{{company}}`, etc.

4. **Configure Settings**
   - Adjust rate limits (daily limits, delays)
   - Enable/disable stealth features
   - Set business hours

5. **Start Automation**
   - Go to Dashboard
   - Click "Start" to begin automation
   - Monitor progress in real-time

### Go Engine Workflow

```bash
# Standard run
go run main.go

# Output:
# ==================================================
#    LinkedIn Automation Tool v1.0.0
#    Go + Rod Browser Automation Engine
# ==================================================
#
# Loading configuration from config.yaml...
# Initializing database at ./linkedin_automation.db...
# Launching browser...
#
# [Step 1] Authenticating...
# ✓ Session restored from saved cookies
#
# [Step 2] Searching for profiles...
# ✓ Found 47 profiles (45 unique, 2 duplicates)
#
# [Step 3] Sending connection requests...
# Remaining connections today: 50
#   ✓ Sent to John Smith
#   ✓ Sent to Jane Doe
# ...
#
# [Step 4] Checking accepted connections...
# ✓ Found 3 newly accepted connections
#
# [Step 5] Sending follow-up messages...
#   ✓ Message sent to Alice Johnson
# ...
#
# ==================================================
#    Automation Summary
# ==================================================
# Connections sent today: 12 / 50
# Messages sent today: 3 / 100
```

---

## Security Considerations

### Credential Protection
- Never log credentials in plain text
- Store in environment variables only
- Encrypt session cookies (AES-256 recommended)
- Clear sensitive data from memory

### Data Privacy
- SQLite database encrypted at rest
- No transmission to external services
- Minimal data retention
- Support data deletion on request

### Ethical Usage
- Respect LinkedIn's Terms of Service
- Use reasonable rate limits
- Avoid aggressive automation
- Consider using dedicated accounts

### Best Practices
1. Use a dedicated LinkedIn account for automation
2. Start with conservative rate limits
3. Monitor for security checkpoints
4. Implement proper error handling
5. Regular session rotation
6. Keep logs for troubleshooting

---

## Project Structure

```
/app/
├── backend/                    # FastAPI Backend
│   ├── server.py              # Main API server
│   ├── requirements.txt       # Python dependencies
│   └── .env                   # Environment variables
│
├── frontend/                   # React Dashboard
│   ├── src/
│   │   ├── App.js            # Main application
│   │   ├── App.css           # Styles
│   │   └── components/       # UI components
│   ├── package.json          # Node dependencies
│   └── .env                  # Environment variables
│
├── go-engine/                  # Go Automation Engine
│   ├── main.go               # Entry point
│   ├── config.yaml           # Configuration
│   ├── go.mod                # Go module
│   │
│   ├── auth/                 # Authentication module
│   │   ├── auth.go          # Login automation
│   │   └── session.go       # Session management
│   │
│   ├── search/               # Search module
│   │   ├── search.go        # Query execution
│   │   ├── parser.go        # Profile extraction
│   │   └── pagination.go    # Pagination handling
│   │
│   ├── messaging/            # Messaging module
│   │   ├── connection.go    # Connection requests
│   │   ├── message.go       # Follow-up messages
│   │   └── templates.go     # Template rendering
│   │
│   ├── stealth/              # Anti-detection module
│   │   ├── bezier.go        # Mouse movement
│   │   ├── timing.go        # Delay patterns
│   │   ├── fingerprint.go   # Browser masking
│   │   ├── scrolling.go     # Scroll behavior
│   │   ├── typing.go        # Typing simulation
│   │   └── mouse.go         # Hover/movement
│   │
│   ├── database/             # Data persistence
│   │   └── database.go      # SQLite operations
│   │
│   ├── config/               # Configuration
│   │   └── config.go        # YAML loading
│   │
│   ├── logger/               # Logging
│   │   └── logger.go        # Structured logging
│   │
│   └── utils/                # Utilities
│       ├── retry.go         # Backoff logic
│       └── validation.go    # Input validation
│
├── tests/                      # Test files
├── README.md                   # This documentation
└── test_result.md             # Testing protocol
```

---

## License

This project is for educational and demonstration purposes. Use responsibly and in compliance with LinkedIn's Terms of Service.

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | July 2025 | Initial release with full feature set |

---

## Support

For issues and feature requests, please create an issue in the repository.
