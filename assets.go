package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type PuzzleDef struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Image      string `json:"image"`
	Collection string `json:"-"`
	thumb      *ebiten.Image
	aspect     float64
}

// Aspect returns the full image's width/height ratio.
func (p *PuzzleDef) Aspect() float64 {
	if p.aspect == 0 {
		p.aspect = 16.0 / 9
		if f, err := os.Open(assetPath("puzzles", p.Image)); err == nil {
			if cfg, _, err := image.DecodeConfig(f); err == nil && cfg.Height > 0 {
				p.aspect = float64(cfg.Width) / float64(cfg.Height)
			}
			f.Close()
		}
	}
	return p.aspect
}

type Collection struct {
	ID       string       `json:"id"`
	Title    string       `json:"title"`
	Subtitle string       `json:"subtitle"`
	Puzzles  []*PuzzleDef `json:"puzzles"`
}

type Manifest struct {
	Collections []*Collection `json:"collections"`
}

var (
	manifest    Manifest
	puzzleByID  = map[string]*PuzzleDef{}
	logoLight   *ebiten.Image
	logoDark    *ebiten.Image
	iconImg     *ebiten.Image
	assetsError string
)

func loadImageFile(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func loadEbitenImage(path string) *ebiten.Image {
	img, err := loadImageFile(path)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}

func loadAssets() {
	for i, n := range []string{"Poppins-Regular.ttf", "Poppins-Medium.ttf", "Poppins-Bold.ttf"} {
		f, err := os.Open(assetPath("fonts", n))
		if err != nil {
			assetsError = "Missing font " + n
			continue
		}
		src, err := text.NewGoTextFaceSource(f)
		f.Close()
		if err == nil {
			fontSources[i] = src
		}
	}
	// fall back so we never crash on a missing weight
	for i := range fontSources {
		if fontSources[i] == nil {
			for j := range fontSources {
				if fontSources[j] != nil {
					fontSources[i] = fontSources[j]
				}
			}
		}
	}
	b, err := os.ReadFile(assetPath("puzzles", "puzzles.json"))
	if err != nil {
		assetsError = fmt.Sprintf("Could not find the puzzles folder (%s).", assetPath("puzzles"))
	} else if err := json.Unmarshal(b, &manifest); err != nil {
		assetsError = "puzzles.json is invalid: " + err.Error()
	}
	for _, c := range manifest.Collections {
		for _, p := range c.Puzzles {
			p.Collection = c.ID
			puzzleByID[p.ID] = p
		}
	}
	logoLight = loadEbitenImage(assetPath("branding", "logo_light.png"))
	logoDark = loadEbitenImage(assetPath("branding", "logo_dark.png"))
	iconImg = loadEbitenImage(assetPath("branding", "icon_256.png"))
	var icons []image.Image
	for _, n := range []string{"icon_16.png", "icon_32.png", "icon_64.png", "icon_256.png"} {
		if im, err := loadImageFile(assetPath("branding", n)); err == nil {
			icons = append(icons, im)
		}
	}
	if len(icons) > 0 {
		ebiten.SetWindowIcon(icons)
	}
}

func (p *PuzzleDef) Thumb() *ebiten.Image {
	if p.thumb == nil {
		p.thumb = loadEbitenImage(assetPath("puzzles", "thumbs", p.Image))
		if p.thumb == nil {
			p.thumb = loadEbitenImage(assetPath("puzzles", p.Image))
		}
		if p.thumb == nil {
			p.thumb = ebiten.NewImage(16, 9)
			p.thumb.Fill(P().Surface2)
		}
	}
	return p.thumb
}

func logo() *ebiten.Image {
	if settings.Theme == "dark" && logoDark != nil {
		return logoDark
	}
	return logoLight
}

func totalPuzzles() int {
	n := 0
	for _, c := range manifest.Collections {
		n += len(c.Puzzles)
	}
	return n
}
