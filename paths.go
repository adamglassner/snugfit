package main

import (
	"os"
	"path/filepath"
)

var assetDir string

func findAssetDir() string {
	var candidates []string
	if exe, err := os.Executable(); err == nil {
		if real, err := filepath.EvalSymlinks(exe); err == nil {
			exe = real
		}
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "assets"))
	}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, "assets"))
	}
	for _, c := range candidates {
		if st, err := os.Stat(filepath.Join(c, "puzzles", "puzzles.json")); err == nil && !st.IsDir() {
			return c
		}
	}
	if len(candidates) > 0 {
		return candidates[0]
	}
	return "assets"
}

func assetPath(parts ...string) string {
	return filepath.Join(append([]string{assetDir}, parts...)...)
}

// dataDir is where settings and saves live: %APPDATA%\Snugfit on Windows.
func dataDir() string {
	base, err := os.UserConfigDir()
	if err != nil {
		base = "."
	}
	d := filepath.Join(base, AppName)
	_ = os.MkdirAll(d, 0o755)
	return d
}

func savesDir() string {
	d := filepath.Join(dataDir(), "saves")
	_ = os.MkdirAll(d, 0o755)
	return d
}
