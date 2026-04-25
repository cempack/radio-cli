package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/cempack/radio-cli/internal/domain"
)

type Store struct {
	dataDir   string
	configDir string
}

func New() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dataDir := filepath.Join(home, ".local", "share", "radiodrift")
	configDir := filepath.Join(home, ".config", "radiodrift")

	for _, dir := range []string{dataDir, configDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	return &Store{dataDir: dataDir, configDir: configDir}, nil
}

func loadJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func saveJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *Store) LoadFavorites() ([]domain.Station, error) {
	var out []domain.Station
	err := loadJSON(filepath.Join(s.dataDir, "favorites.json"), &out)
	if out == nil {
		out = []domain.Station{}
	}
	return out, err
}

func (s *Store) SaveFavorites(stations []domain.Station) error {
	return saveJSON(filepath.Join(s.dataDir, "favorites.json"), stations)
}

func (s *Store) LoadHistory() ([]domain.HistoryItem, error) {
	var out []domain.HistoryItem
	err := loadJSON(filepath.Join(s.dataDir, "history.json"), &out)
	if out == nil {
		out = []domain.HistoryItem{}
	}
	return out, err
}

func (s *Store) SaveHistory(items []domain.HistoryItem) error {
	return saveJSON(filepath.Join(s.dataDir, "history.json"), items)
}

func (s *Store) LoadPresets() ([]domain.Preset, error) {
	var out []domain.Preset
	err := loadJSON(filepath.Join(s.dataDir, "presets.json"), &out)
	if out == nil {
		out = []domain.Preset{}
	}
	return out, err
}

func (s *Store) SavePresets(presets []domain.Preset) error {
	return saveJSON(filepath.Join(s.dataDir, "presets.json"), presets)
}

func (s *Store) LoadCrates() ([]domain.Crate, error) {
	var out []domain.Crate
	err := loadJSON(filepath.Join(s.dataDir, "crates.json"), &out)
	if out == nil {
		out = []domain.Crate{}
	}
	return out, err
}

func (s *Store) SaveCrates(crates []domain.Crate) error {
	return saveJSON(filepath.Join(s.dataDir, "crates.json"), crates)
}

func (s *Store) LoadConfig() (domain.Config, error) {
	var out domain.Config
	err := loadJSON(filepath.Join(s.configDir, "config.json"), &out)
	return out, err
}

func (s *Store) SaveConfig(cfg domain.Config) error {
	return saveJSON(filepath.Join(s.configDir, "config.json"), cfg)
}
