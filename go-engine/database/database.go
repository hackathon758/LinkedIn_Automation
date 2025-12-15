package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB represents the database connection
type DB struct {
	*sql.DB
}

// Connection represents a LinkedIn connection
type Connection struct {
	ID                string
	ProfileURL        string
	FirstName         string
	LastName          string
	JobTitle          string
	Company           string
	Location          string
	NoteSent          string
	Status            string // pending, accepted, declined, failed
	SearchCriteriaID  string
	CreatedAt         time.Time
	AcceptedAt        *time.Time
}

// Message represents a sent message
type Message struct {
	ID           string
	ConnectionID string
	Content      string
	TemplateID   string
	Status       string // sent, delivered, read, failed
	SentAt       time.Time
}

// DailyActivity tracks daily activity for rate limiting
type DailyActivity struct {
	ID                string
	Date              string
	ConnectionsSent   int
	MessagesSent      int
	LastConnectionAt  *time.Time
	LastMessageAt     *time.Time
}

// SessionCookie stores LinkedIn session cookies
type SessionCookie struct {
	ID        string
	Name      string
	Value     string
	Domain    string
	Path      string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrent access
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	return &DB{db}, nil
}

// Initialize creates all required tables
func (db *DB) Initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS connections (
		id TEXT PRIMARY KEY,
		profile_url TEXT NOT NULL UNIQUE,
		first_name TEXT,
		last_name TEXT,
		job_title TEXT,
		company TEXT,
		location TEXT,
		note_sent TEXT,
		status TEXT CHECK(status IN ('pending', 'accepted', 'declined', 'failed')) DEFAULT 'pending',
		search_criteria_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		accepted_at DATETIME
	);

	CREATE INDEX IF NOT EXISTS idx_connections_status ON connections(status);
	CREATE INDEX IF NOT EXISTS idx_connections_created ON connections(created_at);

	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		connection_id TEXT NOT NULL,
		content TEXT NOT NULL,
		template_id TEXT,
		status TEXT CHECK(status IN ('sent', 'delivered', 'read', 'failed')) DEFAULT 'sent',
		sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (connection_id) REFERENCES connections(id)
	);

	CREATE INDEX IF NOT EXISTS idx_messages_connection ON messages(connection_id);
	CREATE INDEX IF NOT EXISTS idx_messages_sent ON messages(sent_at);

	CREATE TABLE IF NOT EXISTS daily_activity (
		id TEXT PRIMARY KEY,
		date TEXT NOT NULL UNIQUE,
		connections_sent INTEGER DEFAULT 0,
		messages_sent INTEGER DEFAULT 0,
		last_connection_at DATETIME,
		last_message_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS session_cookies (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		value TEXT NOT NULL,
		domain TEXT,
		path TEXT,
		expires_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS processed_profiles (
		profile_url TEXT PRIMARY KEY,
		processed_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// ============== Connection Methods ==============

// SaveConnection saves a new connection to the database
func (db *DB) SaveConnection(conn *Connection) error {
	query := `
	INSERT INTO connections (id, profile_url, first_name, last_name, job_title, company, location, note_sent, status, search_criteria_id, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(profile_url) DO UPDATE SET
		note_sent = excluded.note_sent,
		status = excluded.status
	`
	_, err := db.Exec(query, conn.ID, conn.ProfileURL, conn.FirstName, conn.LastName, 
		conn.JobTitle, conn.Company, conn.Location, conn.NoteSent, conn.Status, 
		conn.SearchCriteriaID, conn.CreatedAt)
	return err
}

// UpdateConnectionStatus updates the status of a connection
func (db *DB) UpdateConnectionStatus(profileURL, status string) error {
	query := `UPDATE connections SET status = ? WHERE profile_url = ?`
	if status == "accepted" {
		query = `UPDATE connections SET status = ?, accepted_at = CURRENT_TIMESTAMP WHERE profile_url = ?`
	}
	_, err := db.Exec(query, status, profileURL)
	return err
}

// GetPendingConnections returns all pending connections
func (db *DB) GetPendingConnections() ([]Connection, error) {
	query := `SELECT id, profile_url, first_name, last_name, job_title, company, location, note_sent, status, search_criteria_id, created_at, accepted_at FROM connections WHERE status = 'pending'`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []Connection
	for rows.Next() {
		var c Connection
		err := rows.Scan(&c.ID, &c.ProfileURL, &c.FirstName, &c.LastName, &c.JobTitle, 
			&c.Company, &c.Location, &c.NoteSent, &c.Status, &c.SearchCriteriaID, 
			&c.CreatedAt, &c.AcceptedAt)
		if err != nil {
			return nil, err
		}
		connections = append(connections, c)
	}
	return connections, nil
}

// GetAcceptedConnections returns all accepted connections
func (db *DB) GetAcceptedConnections() ([]Connection, error) {
	query := `SELECT id, profile_url, first_name, last_name, job_title, company, location, note_sent, status, search_criteria_id, created_at, accepted_at FROM connections WHERE status = 'accepted'`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []Connection
	for rows.Next() {
		var c Connection
		err := rows.Scan(&c.ID, &c.ProfileURL, &c.FirstName, &c.LastName, &c.JobTitle, 
			&c.Company, &c.Location, &c.NoteSent, &c.Status, &c.SearchCriteriaID, 
			&c.CreatedAt, &c.AcceptedAt)
		if err != nil {
			return nil, err
		}
		connections = append(connections, c)
	}
	return connections, nil
}

// IsProfileProcessed checks if a profile URL has been processed
func (db *DB) IsProfileProcessed(profileURL string) (bool, error) {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM processed_profiles WHERE profile_url = ?)`, profileURL).Scan(&exists)
	return exists, err
}

// MarkProfileProcessed marks a profile URL as processed
func (db *DB) MarkProfileProcessed(profileURL string) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO processed_profiles (profile_url) VALUES (?)`, profileURL)
	return err
}

// ============== Message Methods ==============

// SaveMessage saves a new message to the database
func (db *DB) SaveMessage(msg *Message) error {
	query := `INSERT INTO messages (id, connection_id, content, template_id, status, sent_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, msg.ID, msg.ConnectionID, msg.Content, msg.TemplateID, msg.Status, msg.SentAt)
	return err
}

// GetMessagesForConnection returns all messages for a specific connection
func (db *DB) GetMessagesForConnection(connectionID string) ([]Message, error) {
	query := `SELECT id, connection_id, content, template_id, status, sent_at FROM messages WHERE connection_id = ? ORDER BY sent_at DESC`
	rows, err := db.Query(query, connectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		err := rows.Scan(&m.ID, &m.ConnectionID, &m.Content, &m.TemplateID, &m.Status, &m.SentAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// HasSentFollowUp checks if a follow-up message has been sent to a connection
func (db *DB) HasSentFollowUp(connectionID string) (bool, error) {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM messages WHERE connection_id = ?)`, connectionID).Scan(&exists)
	return exists, err
}

// ============== Daily Activity Methods ==============

// GetOrCreateDailyActivity gets or creates today's activity record
func (db *DB) GetOrCreateDailyActivity() (*DailyActivity, error) {
	today := time.Now().Format("2006-01-02")
	
	var activity DailyActivity
	err := db.QueryRow(`SELECT id, date, connections_sent, messages_sent, last_connection_at, last_message_at FROM daily_activity WHERE date = ?`, today).Scan(
		&activity.ID, &activity.Date, &activity.ConnectionsSent, &activity.MessagesSent, 
		&activity.LastConnectionAt, &activity.LastMessageAt)
	
	if err == sql.ErrNoRows {
		// Create new record for today
		id := fmt.Sprintf("activity_%s", today)
		_, err = db.Exec(`INSERT INTO daily_activity (id, date, connections_sent, messages_sent) VALUES (?, ?, 0, 0)`, id, today)
		if err != nil {
			return nil, err
		}
		return &DailyActivity{ID: id, Date: today}, nil
	}
	
	return &activity, err
}

// IncrementConnectionCount increments today's connection count
func (db *DB) IncrementConnectionCount() error {
	today := time.Now().Format("2006-01-02")
	_, err := db.Exec(`UPDATE daily_activity SET connections_sent = connections_sent + 1, last_connection_at = CURRENT_TIMESTAMP WHERE date = ?`, today)
	return err
}

// IncrementMessageCount increments today's message count
func (db *DB) IncrementMessageCount() error {
	today := time.Now().Format("2006-01-02")
	_, err := db.Exec(`UPDATE daily_activity SET messages_sent = messages_sent + 1, last_message_at = CURRENT_TIMESTAMP WHERE date = ?`, today)
	return err
}

// ============== Session Cookie Methods ==============

// SaveCookies saves session cookies
func (db *DB) SaveCookies(cookies []SessionCookie) error {
	// Clear existing cookies
	_, err := db.Exec(`DELETE FROM session_cookies`)
	if err != nil {
		return err
	}

	// Insert new cookies
	for _, cookie := range cookies {
		_, err = db.Exec(`INSERT INTO session_cookies (id, name, value, domain, path, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			cookie.ID, cookie.Name, cookie.Value, cookie.Domain, cookie.Path, cookie.ExpiresAt, cookie.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCookies retrieves stored session cookies
func (db *DB) GetCookies() ([]SessionCookie, error) {
	rows, err := db.Query(`SELECT id, name, value, domain, path, expires_at, created_at FROM session_cookies`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cookies []SessionCookie
	for rows.Next() {
		var c SessionCookie
		err := rows.Scan(&c.ID, &c.Name, &c.Value, &c.Domain, &c.Path, &c.ExpiresAt, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		// Skip expired cookies
		if c.ExpiresAt.After(time.Now()) {
			cookies = append(cookies, c)
		}
	}
	return cookies, nil
}

// ClearCookies removes all stored cookies
func (db *DB) ClearCookies() error {
	_, err := db.Exec(`DELETE FROM session_cookies`)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
