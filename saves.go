package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type SavedPiece struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Group  int     `json:"g"`
	Locked bool    `json:"l,omitempty"`
}

type SaveGame struct {
	Version  string       `json:"version"`
	PuzzleID string       `json:"puzzle"`
	Mode     int          `json:"mode"`
	Cols     int          `json:"cols"`
	Rows     int          `json:"rows"`
	Seed     int64        `json:"seed"`
	Elapsed  float64      `json:"elapsed"`
	Placed   int          `json:"placed"`
	Total    int          `json:"total"`
	Updated  time.Time    `json:"updated"`
	CamX     float64      `json:"camX"`
	CamY     float64      `json:"camY"`
	Zoom     float64      `json:"zoom"`
	Ghost    bool         `json:"ghost"`
	Pieces   []SavedPiece `json:"pieces"`
}

func saveKey(puzzleID string, mode int) string {
	return puzzleID + "_" + modes[mode].Key
}

func savePath(key string) string { return filepath.Join(savesDir(), key+".json") }

func writeSave(s *SaveGame) error {
	s.Version = Version
	s.Updated = time.Now()
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	p := savePath(saveKey(s.PuzzleID, s.Mode))
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

func readSave(key string) (*SaveGame, error) {
	b, err := os.ReadFile(savePath(key))
	if err != nil {
		return nil, err
	}
	var s SaveGame
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	if s.Mode < 0 || s.Mode >= len(modes) {
		s.Mode = 0
	}
	return &s, nil
}

func deleteSave(key string) {
	_ = os.Remove(savePath(key))
	if settings.LastSave == key {
		settings.LastSave = ""
		saveSettings()
	}
}

func hasSave(puzzleID string, mode int) bool {
	_, err := os.Stat(savePath(saveKey(puzzleID, mode)))
	return err == nil
}

type SaveInfo struct {
	Key     string
	Game    *SaveGame
	Puzzle  *PuzzleDef
	Percent float64
}

// listSaves returns saves sorted by most recently played. Pieces are dropped to save memory.
func listSaves() []SaveInfo {
	entries, _ := os.ReadDir(savesDir())
	var out []SaveInfo
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		key := strings.TrimSuffix(e.Name(), ".json")
		s, err := readSave(key)
		if err != nil {
			continue
		}
		p := puzzleByID[s.PuzzleID]
		if p == nil {
			continue
		}
		pct := 0.0
		if s.Total > 0 {
			pct = float64(s.Placed) / float64(s.Total)
		}
		s.Pieces = nil
		out = append(out, SaveInfo{Key: key, Game: s, Puzzle: p, Percent: pct})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Game.Updated.After(out[j].Game.Updated) })
	return out
}

func formatDuration(sec float64) string {
	s := int64(sec)
	h := s / 3600
	m := (s % 3600) / 60
	ss := s % 60
	if h > 0 {
		return itoa(h) + ":" + pad2(m) + ":" + pad2(ss)
	}
	return itoa(m) + ":" + pad2(ss)
}

func itoa(v int64) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var b []byte
	for v > 0 {
		b = append([]byte{byte('0' + v%10)}, b...)
		v /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}

func pad2(v int64) string {
	if v < 10 {
		return "0" + itoa(v)
	}
	return itoa(v)
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int64(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return itoa(m) + " minutes ago"
	case d < 24*time.Hour:
		h := int64(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return itoa(h) + " hours ago"
	case d < 48*time.Hour:
		return "yesterday"
	default:
		return t.Format("Jan 2")
	}
}
