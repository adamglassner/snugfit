package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	MusicVolume      float64          `json:"musicVolume"`
	SFXVolume        float64          `json:"sfxVolume"`
	MusicMuted       bool             `json:"musicMuted"`
	SFXMuted         bool             `json:"sfxMuted"`
	Theme            string           `json:"theme"` // "light" or "dark"
	Fullscreen       bool             `json:"fullscreen"`
	TableColor       int              `json:"tableColor"`
	GhostDefault     bool             `json:"ghostDefault"`
	AutoCheckUpdates bool             `json:"autoCheckUpdates"`
	LastSave         string           `json:"lastSave"`
	Completed        map[string]int64 `json:"completed"` // "puzzle|mode" -> best time in seconds
	SkippedVersion   string           `json:"skippedVersion"`
}

var settings = defaultSettings()

func defaultSettings() Settings {
	return Settings{
		MusicVolume:      0.55,
		SFXVolume:        0.8,
		Theme:            "light",
		TableColor:       0,
		GhostDefault:     false,
		AutoCheckUpdates: true,
		Completed:        map[string]int64{},
	}
}

func settingsPath() string { return filepath.Join(dataDir(), "settings.json") }

func loadSettings() {
	b, err := os.ReadFile(settingsPath())
	if err != nil {
		return
	}
	s := defaultSettings()
	if json.Unmarshal(b, &s) == nil {
		if s.Completed == nil {
			s.Completed = map[string]int64{}
		}
		if s.Theme != "dark" {
			s.Theme = "light"
		}
		settings = s
	}
}

func saveSettings() {
	b, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return
	}
	tmp := settingsPath() + ".tmp"
	if os.WriteFile(tmp, b, 0o644) == nil {
		_ = os.Rename(tmp, settingsPath())
	}
}
