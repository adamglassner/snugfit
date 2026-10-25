package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var quitRequested bool

// ---- floating background pieces

type floatPiece struct {
	img            *ebiten.Image
	x, y, vx, vy   float64
	rot, vr, scale float64
	col            int
}

var pieceSilhouettes []*ebiten.Image

func initFloatPieces(g *Game) {
	white := image.NewRGBA(image.Rect(0, 0, 400, 400))
	for i := range white.Pix {
		white.Pix[i] = 255
	}
	for k := 0; k < 6; k++ {
		geo := newGeometry(300, 300, 3, 3, int64(1000+k*17))
		pr := rasterPiece(geo, white, 1, 1)
		pieceSilhouettes = append(pieceSilhouettes, ebiten.NewImageFromImage(pr.RGBA))
	}
	for i := 0; i < 14; i++ {
		g.bgPieces = append(g.bgPieces, floatPiece{
			img: pieceSilhouettes[i%len(pieceSilhouettes)],
			x:   rand.Float64() * 1600, y: rand.Float64() * 1000,
			vx: (rand.Float64() - 0.5) * 0.25, vy: (rand.Float64() - 0.5) * 0.25,
			rot: rand.Float64() * 6.28, vr: (rand.Float64() - 0.5) * 0.002,
			scale: 0.35 + rand.Float64()*0.55, col: i % 4,
		})
	}
}

func drawFloatPieces(dst *ebiten.Image) {
	p := P()
	cols := []color.RGBA{p.Accent, p.Teal, p.Plum, p.Gold}
	a := float32(0.09)
	if settings.Theme == "dark" {
		a = 0.07
	}
	for i := range game.bgPieces {
		f := &game.bgPieces[i]
		f.x += f.vx
		f.y += f.vy
		f.rot += f.vr
		if f.x < -150 {
			f.x = SW + 150
		} else if f.x > SW+150 {
			f.x = -150
		}
		if f.y < -150 {
			f.y = SH + 150
		} else if f.y > SH+150 {
			f.y = -150
		}
		b := f.img.Bounds()
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
		op.GeoM.Rotate(f.rot)
		op.GeoM.Scale(f.scale*S, f.scale*S)
		op.GeoM.Translate(f.x*S, f.y*S)
		c := cols[f.col]
		op.ColorScale.Scale(float32(c.R)/255*a, float32(c.G)/255*a, float32(c.B)/255*a, a)
		dst.DrawImage(f.img, op)
	}
}

// ---- main menu

type MenuScene struct {
	widgets []Widget
	saves   []SaveInfo
	last    *SaveInfo
	t       float64
}

func NewMenuScene() *MenuScene {
	m := &MenuScene{}
	m.saves = listSaves()
	for i := range m.saves {
		if m.saves[i].Key == settings.LastSave {
			m.last = &m.saves[i]
		}
	}
	if m.last == nil && len(m.saves) > 0 {
		m.last = &m.saves[0]
	}
	return m
}

func (m *MenuScene) continueLast() {
	if m.last == nil {
		return
	}
	sv, err := readSave(m.last.Key)
	if err != nil {
		return
	}
	game.SetScene(NewPlayScene(m.last.Puzzle, sv.Mode, sv))
}

func (m *MenuScene) layout() {
	m.widgets = m.widgets[:0]
	lx := math.Max(56, SW*0.06)
	bw := 320.0
	y := math.Max(250, SH*0.36)
	add := func(id, label string, kind ButtonKind, badge bool, fn func()) {
		m.widgets = append(m.widgets, &Button{ID: id, R: Rect{lx, y, bw, 54}, Label: label, Kind: kind, Size: 16, Align: AlignLeft, Badge: badge, OnClick: fn})
		y += 64
	}
	if m.last != nil {
		add("m_cont", "Continue puzzle", BtnPrimary, false, m.continueLast)
		add("m_new", "New puzzle", BtnSecondary, false, func() { game.SetScene(NewGalleryScene()) })
	} else {
		add("m_new", "New puzzle", BtnPrimary, false, func() { game.SetScene(NewGalleryScene()) })
	}
	add("m_saves", fmt.Sprintf("Saved puzzles (%d)", len(m.saves)), BtnSecondary, false, func() { game.SetScene(NewSavesScene()) })
	add("m_opts", "Options", BtnSecondary, updateChecker.available(), func() { game.SetScene(NewOptionsScene(nil)) })
	add("m_log", "What's new", BtnSecondary, false, func() { game.SetScene(NewChangelogScene()) })
	add("m_quit", "Quit", BtnGhost, false, func() { quitRequested = true })

	hero := m.heroRect()
	if m.last != nil {
		m.widgets = append(m.widgets, &Button{ID: "h_cont", R: Rect{hero.Right() - 190, hero.Bottom() - 70, 166, 46}, Label: "Continue", Kind: BtnPrimary, OnClick: m.continueLast})
	} else {
		m.widgets = append(m.widgets, &Button{ID: "h_browse", R: Rect{hero.Right() - 190, hero.Bottom() - 70, 166, 46}, Label: "Browse puzzles", Kind: BtnPrimary, OnClick: func() { game.SetScene(NewGalleryScene()) }})
	}
	if updateChecker.available() {
		w := 300.0
		m.widgets = append(m.widgets, &Button{ID: "upd", R: Rect{SW - w - 24, 20, w, 40}, Label: "Update available: v" + updateChecker.latest.Version, Kind: BtnTeal, Size: 14, OnClick: func() { game.SetScene(NewOptionsScene(nil)) }})
	}
}

func (m *MenuScene) heroRect() Rect {
	x := math.Max(440, SW*0.42)
	w := SW - x - math.Max(48, SW*0.05)
	h := math.Min(w*0.62, SH-330)
	return Rect{x, math.Max(100, SH*0.14), w, h}
}

func (m *MenuScene) Update() {
	m.t += 1.0 / 60
	m.layout()
	for _, w := range m.widgets {
		w.Update()
	}
}

func (m *MenuScene) Draw(dst *ebiten.Image) {
	p := P()
	dst.Fill(p.Bg)
	drawFloatPieces(dst)
	lx := math.Max(56, SW*0.06)
	if lg := logo(); lg != nil {
		DrawImageFit(dst, lg, Rect{lx - 8, math.Max(70, SH*0.1), 380, 114}, 1)
	}
	Text(dst, "Relaxing jigsaw puzzles for every season", Regular, 15, lx, math.Max(70, SH*0.1)+136, p.TextDim, AlignLeft)

	hero := m.heroRect()
	Panel(dst, hero, 24)
	img := Rect{hero.X + 16, hero.Y + 16, hero.W - 32, hero.H - 110}
	if m.last != nil {
		DrawImageRR(dst, m.last.Puzzle.Thumb(), img, 16, 1)
		pill := "In progress"
		pw := TextWidth(pill, Bold, 12) + 24
		FillRR(dst, Rect{img.X + 14, img.Y + 14, pw, 28}, 14, alpha(p.Overlay, 0.65))
		Text(dst, pill, Bold, 12, img.X+14+pw/2, img.Y+28, color.RGBA{255, 255, 255, 255}, AlignCenter)
		Text(dst, m.last.Puzzle.Title, Bold, 20, hero.X+28, hero.Bottom()-68, p.Text, AlignLeft)
		g := m.last.Game
		Text(dst, fmt.Sprintf("%s  -  %d of %d pieces  -  %s", modes[g.Mode].Name, g.Placed, g.Total, formatDuration(g.Elapsed)), Regular, 13.5, hero.X+28, hero.Bottom()-42, p.TextDim, AlignLeft)
		bar := Rect{hero.X + 28, hero.Bottom() - 24, math.Min(320, hero.W-260), 7}
		FillRR(dst, bar, 3.5, p.Track)
		FillRR(dst, Rect{bar.X, bar.Y, math.Max(7, bar.W*m.last.Percent), bar.H}, 3.5, p.Teal)
	} else if len(manifest.Collections) > 0 {
		col := manifest.Collections[0]
		n := len(col.Puzzles)
		if n > 0 {
			// mosaic: one large + others small
			big := Rect{img.X, img.Y, img.W * 0.6, img.H}
			DrawImageRR(dst, col.Puzzles[0].Thumb(), big, 14, 1)
			rest := col.Puzzles[1:]
			if len(rest) > 4 {
				rest = rest[:4]
			}
			gx := big.Right() + 10
			gw := (img.Right() - gx - 10) / 2
			gh := (img.H - 10) / 2
			for i, pz := range rest {
				r := Rect{gx + float64(i%2)*(gw+10), img.Y + float64(i/2)*(gh+10), gw, gh}
				DrawImageRR(dst, pz.Thumb(), r, 12, 1)
			}
		}
		Text(dst, "New this month: "+col.Title, Bold, 20, hero.X+28, hero.Bottom()-64, p.Text, AlignLeft)
		Text(dst, fmt.Sprintf("%d new puzzles  -  %s", n, col.Subtitle), Regular, 13.5, hero.X+28, hero.Bottom()-38, p.TextDim, AlignLeft)
	}

	// coming soon card
	cs := Rect{hero.X, hero.Bottom() + 20, hero.W, 86}
	if cs.Bottom() < SH-40 {
		FillRR(dst, cs, 18, alpha(p.Accent, 0.12))
		StrokeRR(dst, cs, 18, 1.5, alpha(p.Accent, 0.5))
		FillCircle(dst, cs.X+44, cs.CY(), 22, p.Accent)
		Text(dst, "+5", Bold, 16, cs.X+44, cs.CY(), p.AccentText, AlignCenter)
		Text(dst, "More puzzles coming soon!", Bold, 17, cs.X+82, cs.CY()-12, p.Text, AlignLeft)
		Text(dst, "Five new puzzles arrive every month. Updates install right from Options.", Regular, 13.5, cs.X+82, cs.CY()+13, p.TextDim, AlignLeft)
	}

	for _, w := range m.widgets {
		w.Draw(dst)
	}
	Text(dst, "Version "+Version, Regular, 12.5, lx, SH-28, p.TextDim, AlignLeft)
	if n := musicTrackName(); n != "" && musicVolume() > 0 {
		Text(dst, "Now playing: "+n, Regular, 12.5, SW-28, SH-28, p.TextDim, AlignRight)
	}
	if assetsError != "" {
		r := Rect{SW/2 - 300, SH - 90, 600, 44}
		FillRR(dst, r, 12, p.Danger)
		Text(dst, assetsError, Medium, 13, r.CX(), r.CY(), color.RGBA{255, 255, 255, 255}, AlignCenter)
	}
}
