package domain

import "time"

type Station struct {
	ID              string
	Name            string
	StreamURL       string
	HomepageURL     string
	FaviconURL      string
	Country         string
	CountryCode     string
	State           string
	Language        string
	Tags            []string
	Codec           string
	Bitrate         int
	Votes           int
	ClickCount      int
	LastCheckOK     bool
	LastCheckTime   time.Time
	LastCheckStatus string
	ResolvedURL     string
}

type Crate struct {
	ID          string
	Name        string
	Description string
	Tags        []string
	StationIDs  []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Preset struct {
	Key       string
	Name      string
	StationID string
	CreatedAt time.Time
}

type HistoryItem struct {
	StationID string
	StartedAt time.Time
	EndedAt   time.Time
	Duration  time.Duration
	Context   string
}

type Config struct {
	PlayerBackend     string
	VolumeStep        int
	Theme             string
	CompactMode       bool
	StartView         string
	AutoReconnect     bool
	ReconnectAttempts int
	SaveHistory       bool
	SaveSession       bool
	CountryBias       []string
	TagBias           []string
}
