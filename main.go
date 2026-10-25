package main

import (
	_ "embed"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed CHANGELOG.md
var changelogMD string

// SW, SH are the logical screen size.
var SW, SH = 1280.0, 800.0

type Scene interface {
	Update()
	Draw(dst *ebiten.Image)
}

type Game struct {
	scene     Scene
	next      Scene
	fade      float64
	fadingOut bool
	bgPieces  []floatPiece
	tour      *tour
}

var game *Game

func (g *Game) SetScene(s Scene) {
	if g.scene == nil {
		g.scene = s
		return
	}
	g.next = s
	g.fadingOut = true
}

func (g *Game) Update() error {
	if quitRequested {
		if ps, ok := g.scene.(*PlayScene); ok {
			ps.saveNow(false)
		}
		saveSettings()
		return ebiten.Termination
	}
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	in.poll()
	if g.tour != nil {
		g.tour.step(g)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) || (inpututil.IsKeyJustPressed(ebiten.KeyEnter) && ebiten.IsKeyPressed(ebiten.KeyAlt)) {
		settings.Fullscreen = !settings.Fullscreen
		ebiten.SetFullscreen(settings.Fullscreen)
		saveSettings()
	}
	updateMusic()
	updateChecker.poll()
	if g.fadingOut {
		g.fade = math.Min(1, g.fade+0.16)
		if g.fade >= 1 {
			g.scene = g.next
			g.next = nil
			g.fadingOut = false
		}
		return nil
	}
	if g.fade > 0 {
		g.fade = math.Max(0, g.fade-0.12)
	}
	if g.scene != nil {
		g.scene.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.scene != nil {
		g.scene.Draw(screen)
	}
	if g.fade > 0 {
		FillRect(screen, Rect{0, 0, SW, SH}, alpha(P().Bg, g.fade))
	}
	if g.tour != nil {
		g.tour.capture(screen)
	}
}

func (g *Game) Layout(ow, oh int) (int, int) {
	panic("use LayoutF")
}

func (g *Game) LayoutF(ow, oh float64) (float64, float64) {
	s := ebiten.Monitor().DeviceScaleFactor()
	if s <= 0 {
		s = 1
	}
	S = s
	SW, SH = ow, oh
	return math.Ceil(ow * S), math.Ceil(oh * S)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	assetDir = findAssetDir()
	loadSettings()
	ebiten.SetWindowTitle(AppName)
	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowSizeLimits(1024, 640, -1, -1)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(settings.Fullscreen)
	ebiten.SetTPS(60)
	ebiten.SetVsyncEnabled(true)
	loadAssets()
	initAudio()
	game = &Game{}
	initFloatPieces(game)
	if dir := os.Getenv("SNUGFIT_TOUR"); dir != "" {
		game.tour = newTour(dir)
	}
	if settings.AutoCheckUpdates && game.tour == nil {
		updateChecker.check(false)
	}
	game.SetScene(NewMenuScene())
	op := &ebiten.RunGameOptions{}
	if err := ebiten.RunGameWithOptions(game, op); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
	saveSettings()
}
