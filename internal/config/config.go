package config

import "github.com/cempack/radio-cli/internal/domain"

func DefaultConfig() domain.Config {
	return domain.Config{
		PlayerBackend:     "mpv",
		VolumeStep:        5,
		Theme:             "dark",
		CompactMode:       false,
		StartView:         "home",
		AutoReconnect:     true,
		ReconnectAttempts: 3,
		SaveHistory:       true,
		SaveSession:       true,
	}
}
