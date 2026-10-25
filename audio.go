package main

import (
	"bytes"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const sampleRate = 44100

var (
	audioCtx   *audio.Context
	sfxData    = map[string][]byte{}
	musicFiles []string
	musicIdx   int
	musicPl    *audio.Player
	musicFile  *os.File
	lastSFX    = map[string]int64{}
	musicFade  = 0.0
	musicPaused bool
)

func initAudio() {
	audioCtx = audio.NewContext(sampleRate)
	sfxDir := assetPath("sfx")
	entries, _ := os.ReadDir(sfxDir)
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".ogg") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(sfxDir, e.Name()))
		if err != nil {
			continue
		}
		s, err := vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(b))
		if err != nil {
			log.Println("sfx decode:", err)
			continue
		}
		pcm, err := io.ReadAll(s)
		if err != nil {
			continue
		}
		sfxData[strings.TrimSuffix(e.Name(), ".ogg")] = pcm
	}
	mus, _ := os.ReadDir(assetPath("music"))
	for _, e := range mus {
		if strings.HasSuffix(e.Name(), ".ogg") {
			musicFiles = append(musicFiles, filepath.Join(assetPath("music"), e.Name()))
		}
	}
	sort.Strings(musicFiles)
	if len(musicFiles) > 0 {
		musicIdx = rand.Intn(len(musicFiles))
	}
}

func playSFX(name string) {
	if audioCtx == nil || settings.SFXMuted || settings.SFXVolume <= 0 {
		return
	}
	b, ok := sfxData[name]
	if !ok {
		return
	}
	// throttle identical sounds within 3 ticks
	if t, ok := lastSFX[name]; ok && in.Tick-t < 3 {
		return
	}
	lastSFX[name] = in.Tick
	p := audioCtx.NewPlayerFromBytes(b)
	v := settings.SFXVolume
	p.SetVolume(v * v)
	p.Play()
}

func musicVolume() float64 {
	if settings.MusicMuted {
		return 0
	}
	v := settings.MusicVolume
	return v * v * 0.8
}

func startTrack(i int) {
	if len(musicFiles) == 0 {
		return
	}
	if musicPl != nil {
		musicPl.Close()
		musicPl = nil
	}
	if musicFile != nil {
		musicFile.Close()
		musicFile = nil
	}
	musicIdx = (i%len(musicFiles) + len(musicFiles)) % len(musicFiles)
	f, err := os.Open(musicFiles[musicIdx])
	if err != nil {
		return
	}
	s, err := vorbis.DecodeWithSampleRate(sampleRate, f)
	if err != nil {
		f.Close()
		return
	}
	p, err := audioCtx.NewPlayer(s)
	if err != nil {
		f.Close()
		return
	}
	musicFile = f
	musicPl = p
	musicFade = 0
	p.SetVolume(0)
	p.Play()
}

// updateMusic is called every tick: keeps the playlist going and applies volume.
func updateMusic() {
	if audioCtx == nil || len(musicFiles) == 0 {
		return
	}
	if musicVolume() <= 0 {
		if musicPl != nil && musicPl.IsPlaying() {
			musicPl.Pause()
			musicPaused = true
		}
		return
	}
	if musicPl == nil {
		startTrack(musicIdx)
		return
	}
	if !musicPl.IsPlaying() {
		if !musicPaused {
			// finished: move to next track
			startTrack(musicIdx + 1)
			return
		}
		musicPaused = false
		musicPl.Play()
	}
	musicFade = approach(musicFade, 1, 0.02)
	musicPl.SetVolume(musicVolume() * musicFade)
}

func musicTrackName() string {
	if len(musicFiles) == 0 {
		return ""
	}
	n := strings.TrimSuffix(filepath.Base(musicFiles[musicIdx]), ".ogg")
	n = strings.ReplaceAll(n, "-", " ")
	words := strings.Fields(n)
	for i, w := range words {
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return strings.Join(words, " ")
}
