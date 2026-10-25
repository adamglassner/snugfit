package main

// Debug "screenshot tour": set SNUGFIT_TOUR=<dir> to walk through every screen
// automatically and save PNG screenshots. Used for automated visual testing.

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

type tourStep struct {
	wait   int
	action func(g *Game) bool // returns true when ready to move on
	shot   string
}

type tour struct {
	dir     string
	steps   []tourStep
	i       int
	frames  int
	pending string
	testPiece *Piece
}

func newTour(dir string) *tour {
	_ = os.MkdirAll(dir, 0o755)
	t := &tour{dir: dir}
	settings.Theme = "light"
	settings.MusicMuted = true
	settings.SFXMuted = true
	play := func() *PlayScene { ps, _ := game.scene.(*PlayScene); return ps }
	waitLoaded := func(g *Game) bool { ps := play(); return ps != nil && !ps.loading && g.fade == 0 }
	t.steps = []tourStep{
		{wait: 40, shot: "01_menu"},
		{action: func(g *Game) bool { g.SetScene(NewGalleryScene()); return true }},
		{wait: 40, shot: "02_gallery"},
		{action: func(g *Game) bool {
			if gs, ok := g.scene.(*GalleryScene); ok {
				gs.selected = puzzleByID["boston-terriers"]
			}
			return true
		}},
		{wait: 20, shot: "03_difficulty"},
		{action: func(g *Game) bool { g.SetScene(NewPlayScene(puzzleByID["pumpkin-patch"], 0, nil)); return true }},
		{wait: 10, action: waitLoaded},
		{wait: 30, shot: "04_play_start"},
		{action: func(g *Game) bool {
			// real mouse-drag test: drag one loose piece onto its spot on the board
			ps := play()
			var p *Piece
			for _, q := range ps.pieces {
				if !q.Locked && q.C == 2 && q.R == 2 {
					p = q
				}
			}
			sx := (p.X+ps.geo.PW/2-ps.camX)*ps.zoom/S
			sy := (p.Y+ps.geo.PH/2-ps.camY)*ps.zoom/S
			tx := (float64(p.C)*ps.geo.PW+ps.geo.PW/2+ps.geo.PW*0.1-ps.camX)*ps.zoom/S
			ty := (float64(p.R)*ps.geo.PH+ps.geo.PH/2-ps.camY)*ps.zoom/S
			fakeQueue = append(fakeQueue, Input{MX: sx, MY: sy, Pressed: true, Down: true})
			for k := 1; k <= 20; k++ {
				f := float64(k) / 20
				fakeQueue = append(fakeQueue, Input{MX: sx + (tx-sx)*f, MY: sy + (ty-sy)*f, Down: true})
			}
			fakeQueue = append(fakeQueue, Input{MX: tx, MY: ty, Released: true})
			t.testPiece = p
			return true
		}},
		{wait: 40, action: func(g *Game) bool {
			if t.testPiece != nil && t.testPiece.Locked {
				fmt.Println("DRAGTEST PASS: piece dragged with the mouse snapped into place")
			} else {
				fmt.Println("DRAGTEST FAIL")
			}
			return true
		}},
		{action: func(g *Game) bool {
			ps := play()
			// place ~45% of pieces and connect a loose cluster
			for i, p := range ps.pieces {
				if p.R < ps.geo.Rows/2 || p.C < 3 {
					tx, ty := float64(p.C)*ps.geo.PW, float64(p.R)*ps.geo.PH
					ps.moveGroup(p.Group, tx-p.X, ty-p.Y, false)
					_ = i
				}
			}
			for _, p := range ps.pieces {
				if p.R < ps.geo.Rows/2 || p.C < 3 {
					ps.drop(p.Group)
				}
			}
			ps.showPreview = true
			ps.hintT = 0
			return true
		}},
		{wait: 40, shot: "05_play_progress"},
		{action: func(g *Game) bool { ps := play(); ps.showPreview = false; ps.edgesOnly = true; ps.fit(false); return true }},
		{wait: 20, shot: "06_edges_only"},
		{action: func(g *Game) bool { ps := play(); ps.edgesOnly = false; ps.paused = true; ps.saveNow(false); return true }},
		{wait: 20, shot: "07_pause"},
		{action: func(g *Game) bool { g.SetScene(NewOptionsScene(nil)); return true }},
		{wait: 40, shot: "08_options"},
		{action: func(g *Game) bool { settings.Theme = "dark"; g.SetScene(NewMenuScene()); return true }},
		{wait: 40, shot: "09_menu_dark"},
		{action: func(g *Game) bool { g.SetScene(NewGalleryScene()); return true }},
		{wait: 40, shot: "10_gallery_dark"},
		{action: func(g *Game) bool { g.SetScene(NewSavesScene()); return true }},
		{wait: 40, shot: "11_saves_dark"},
		{action: func(g *Game) bool { g.SetScene(NewOptionsScene(nil)); return true }},
		{wait: 40, shot: "12_options_dark"},
		{action: func(g *Game) bool { settings.Theme = "light"; g.SetScene(NewChangelogScene()); return true }},
		{wait: 40, shot: "13_changelog"},
		{action: func(g *Game) bool { settings.TableColor = 1; g.SetScene(NewPlayScene(puzzleByID["ocean-reef"], 3, nil)); return true }},
		{wait: 10, action: waitLoaded},
		{wait: 30, shot: "14_advanced_1000"},
		{action: func(g *Game) bool {
			ps := play()
			ps.zoomAt(3.5, SW*0.18, SH*0.5)
			return true
		}},
		{wait: 20, shot: "15_advanced_zoom"},
		{action: func(g *Game) bool {
			ps := play()
			for _, p := range ps.pieces {
				tx, ty := float64(p.C)*ps.geo.PW, float64(p.R)*ps.geo.PH
				ps.moveGroup(p.Group, tx-p.X, ty-p.Y, false)
			}
			ps.drop(ps.pieces[0].Group)
			return true
		}},
		{wait: 110, shot: "16_complete"},
		{action: func(g *Game) bool { settings.TableColor = 0; quitRequested = true; return true }},
	}
	if st := os.Getenv("SNUGFIT_TOUR_START"); st != "" {
		n := 0
		fmt.Sscanf(st, "%d", &n)
		if n > 0 && n < len(t.steps) {
			t.i = n
		}
	}
	return t
}

func (t *tour) step(g *Game) {
	if t.i >= len(t.steps) || t.pending != "" {
		return
	}
	st := t.steps[t.i]
	t.frames++
	if t.frames < st.wait {
		return
	}
	if st.action != nil && !st.action(g) {
		return
	}
	if st.shot != "" {
		t.pending = st.shot
	}
	t.i++
	t.frames = 0
}

func (t *tour) capture(screen *ebiten.Image) {
	if t.pending == "" {
		return
	}
	b := screen.Bounds()
	img := image.NewRGBA(b)
	screen.ReadPixels(img.Pix)
	f, err := os.Create(filepath.Join(t.dir, t.pending+".png"))
	if err == nil {
		_ = png.Encode(f, img)
		f.Close()
	}
	fmt.Println("shot", t.pending)
	t.pending = ""
}
