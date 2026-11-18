package main

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GalleryScene struct {
	widgets  []Widget
	scroll   ScrollState
	selected *PuzzleDef
	confirm  int // mode index awaiting "start over" confirmation, -1 none
	cards    []cardInfo
	sectionY []float64
	saved    map[string]int
	cols     []*Collection // collections shown on this screen
	title    string
	sub      string
	back     func()
	isMain   bool
}

func (g *GalleryScene) has(id string, mode int) bool {
	_, ok := g.saved[saveKey(id, mode)]
	return ok
}

type cardInfo struct {
	p    *PuzzleDef
	r    Rect
	soon bool
	col  *Collection
	cat  *Collection // a category card that opens its subcategories
}

func newGallery(cols []*Collection, title, sub string, back func(), isMain bool) *GalleryScene {
	g := &GalleryScene{confirm: -1, saved: map[string]int{}, cols: cols, title: title, sub: sub, back: back, isMain: isMain}
	for _, si := range listSaves() {
		g.saved[si.Key] = int(si.Percent * 100)
	}
	return g
}

func NewGalleryScene() *GalleryScene {
	return newGallery(manifest.Collections, "Choose a puzzle",
		fmt.Sprintf("%d puzzles  -  4 difficulty levels each", totalPuzzles()),
		func() { game.SetScene(NewMenuScene()) }, true)
}

// NewSubcategoryGallery shows the puzzles of one subcategory (e.g. one state).
func NewSubcategoryGallery(parent, c *Collection) *GalleryScene {
	return newGallery([]*Collection{c}, c.Title, fmt.Sprintf("%s  -  %d puzzles", parent.Title, len(c.Puzzles)),
		func() { game.SetScene(NewCategoryScene(parent)) }, false)
}

func header(widgets *[]Widget, title string, back func()) {
	*widgets = append(*widgets, &Button{ID: "back", R: Rect{24, 22, 104, 42}, Label: "< Back", Kind: BtnGhost, OnClick: back})
}

func drawHeader(dst *ebiten.Image, title, sub string) {
	p := P()
	Text(dst, title, Bold, 28, 148, 40, p.Text, AlignLeft)
	if sub != "" {
		Text(dst, sub, Regular, 14, 150, 68, p.TextDim, AlignLeft)
	}
}

func (g *GalleryScene) gridMetrics() (float64, float64, float64, int) {
	margin := math.Max(40, SW*0.05)
	gap := 22.0
	cols := int((SW - 2*margin + gap) / (250 + gap))
	if cols > 5 {
		cols = 5
	}
	if cols < 2 {
		cols = 2
	}
	cw := (SW - 2*margin - gap*float64(cols-1)) / float64(cols)
	return margin, gap, cw, cols
}

func (g *GalleryScene) layoutCards() float64 {
	g.cards = g.cards[:0]
	g.sectionY = g.sectionY[:0]
	margin, gap, cw, cols := g.gridMetrics()
	ch := cw*0.5625 + 62
	y := 110.0 - g.scroll.Y
	for ci, c := range g.cols {
		y += 50 // section title
		g.sectionY = append(g.sectionY, y-30)
		if len(c.Subcategories) > 0 {
			// a category of categories: one card that opens it
			r := Rect{margin, y, cw, ch}
			g.cards = append(g.cards, cardInfo{r: r, cat: c})
			y += ch + gap + 16
			continue
		}
		items := len(c.Puzzles)
		extra := 0
		if (ci == 0 && g.isMain) || c.ComingSoon {
			extra = 1
		}
		for i := 0; i < items+extra; i++ {
			col, row := i%cols, i/cols
			r := Rect{margin + float64(col)*(cw+gap), y + float64(row)*(ch+gap), cw, ch}
			if i < items {
				g.cards = append(g.cards, cardInfo{p: c.Puzzles[i], r: r})
			} else {
				if c.ComingSoon && items == 0 && cols >= 3 {
					r.W = cw*2 + gap // an empty upcoming category gets a wider card
				}
				g.cards = append(g.cards, cardInfo{r: r, soon: true, col: c})
			}
		}
		rows := (items + extra + cols - 1) / cols
		y += float64(rows)*(ch+gap) + 16
	}
	return y + g.scroll.Y
}

func (g *GalleryScene) Update() {
	g.widgets = g.widgets[:0]
	if g.selected != nil {
		g.updateModal()
		return
	}
	header(&g.widgets, "", g.back)
	contentH := g.layoutCards()
	g.scroll.Max = contentH - SH + 30
	area := Rect{0, 90, SW, SH - 90}
	g.scroll.Update(area)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.back()
		return
	}
	for _, w := range g.widgets {
		w.Update()
	}
	for i := range g.cards {
		c := &g.cards[i]
		if c.soon {
			continue
		}
		if c.cat != nil {
			h := in.Hover(c.r) && in.MY > 90
			animHover("cat_"+c.cat.ID, h)
			if h {
				ebiten.SetCursorShape(ebiten.CursorShapePointer)
			}
			if h && in.Click(c.r) {
				playSFX("click")
				game.SetScene(NewCategoryScene(c.cat))
			}
			continue
		}
		h := in.Hover(c.r) && in.MY > 90
		animHover("card_"+c.p.ID, h)
		if h {
			ebiten.SetCursorShape(ebiten.CursorShapePointer)
		}
		if h && in.Click(c.r) {
			playSFX("click")
			g.selected = c.p
			g.confirm = -1
		}
	}
}

func (g *GalleryScene) Draw(dst *ebiten.Image) {
	p := P()
	dst.Fill(p.Bg)
	drawFloatPieces(dst)
	if g.selected == nil || true {
		g.layoutCards()
		y := 110.0 - g.scroll.Y
		// section titles
		margin, _, _, _ := g.gridMetrics()
		for i, c := range g.cols {
			if i >= len(g.sectionY) {
				break
			}
			ty := g.sectionY[i]
			Text(dst, c.Title, Bold, 20, margin, ty, p.Text, AlignLeft)
			tw := TextWidth(c.Title, Bold, 20)
			if c.ComingSoon && (c.PuzzleCount() == 0 || len(c.Subcategories) > 0) {
				lbl := "Coming soon"
				if c.PuzzleCount() > 0 {
					lbl = "More coming soon"
				}
				pw := TextWidth(lbl, Bold, 11.5) + 20
				FillRR(dst, Rect{margin + tw + 12, ty - 12, pw, 24}, 12, p.Accent)
				Text(dst, lbl, Bold, 11.5, margin+tw+12+pw/2, ty, p.AccentText, AlignCenter)
				tw += pw + 12
			}
			Text(dst, c.Subtitle, Regular, 14, margin+tw+14, ty+1, p.TextDim, AlignLeft)
		}
		_ = y
		for _, c := range g.cards {
			if c.r.Bottom() < 80 || c.r.Y > SH {
				continue
			}
			if c.cat != nil {
				drawCategoryCard(dst, c.r, c.cat)
			} else if c.soon {
				g.drawSoonCard(dst, c.r, c.col)
			} else {
				g.drawCard(dst, c)
			}
		}
	}
	// header bar
	FillRect(dst, Rect{0, 0, SW, 90}, alpha(p.Bg, 0.96))
	drawHeader(dst, g.title, g.sub)
	for _, w := range g.widgets {
		w.Draw(dst)
	}
	g.scroll.DrawBar(dst, Rect{0, 90, SW - 4, SH - 90})
	if g.selected != nil {
		g.drawModal(dst)
	}
}

func completedModes(id string) []bool {
	out := make([]bool, len(modes))
	for i, m := range modes {
		_, out[i] = settings.Completed[id+"|"+m.Key]
	}
	return out
}

func (g *GalleryScene) drawCard(dst *ebiten.Image, c cardInfo) {
	p := P()
	hv := hoverAnim["card_"+c.p.ID]
	r := c.r.Move(0, -hv*4)
	Shadow(dst, r, 18, 1+hv)
	FillRR(dst, r, 18, p.Surface)
	if hv > 0.01 {
		StrokeRR(dst, r, 18, 2, alpha(p.Accent, hv))
	} else {
		StrokeRR(dst, r, 18, 1, p.Border)
	}
	img := Rect{r.X + 10, r.Y + 10, r.W - 20, r.W*0.5625 - 10}
	DrawImageRR(dst, c.p.Thumb(), img, 12, 1)
	Text(dst, c.p.Title, Bold, 15, r.X+16, img.Bottom()+24, p.Text, AlignLeft)
	// completion dots
	done := completedModes(c.p.ID)
	n := 0
	for i, d := range done {
		cx := r.Right() - 20 - float64(len(done)-1-i)*14
		if d {
			FillCircle(dst, cx, img.Bottom()+44, 5, p.Teal)
			n++
		} else {
			StrokeCircle(dst, cx, img.Bottom()+44, 4.5, 1.5, p.TextDim)
		}
	}
	saved := 0
	for i := range modes {
		if g.has(c.p.ID, i) {
			saved++
		}
	}
	if saved > 0 {
		pill := "In progress"
		pw := TextWidth(pill, Bold, 11.5) + 20
		FillRR(dst, Rect{img.X + 10, img.Y + 10, pw, 24}, 12, alpha(p.Accent, 0.92))
		Text(dst, pill, Bold, 11.5, img.X+10+pw/2, img.Y+22, p.AccentText, AlignCenter)
	}
	sub := "Not started yet"
	if saved > 0 && n == 0 {
		sub = "In progress"
	}
	if n == len(modes) {
		sub = "All levels complete!"
	} else if n > 0 {
		sub = fmt.Sprintf("%d of %d levels complete", n, len(modes))
	}
	Text(dst, sub, Regular, 12, r.X+16, img.Bottom()+44, p.TextDim, AlignLeft)
}

func (g *GalleryScene) drawSoonCard(dst *ebiten.Image, r Rect, col *Collection) {
	p := P()
	FillRR(dst, r, 18, alpha(p.Accent, 0.07))
	// dashed border
	steps := 60
	per := 2 * (r.W + r.H)
	for i := 0; i < steps; i += 2 {
		a := per * float64(i) / float64(steps)
		b := per * float64(i+1) / float64(steps)
		x0, y0 := perim(r, a)
		x1, y1 := perim(r, b)
		Line(dst, x0, y0, x1, y1, 2, alpha(p.Accent, 0.6))
	}
	badge, title, sub := "+5", "More puzzles coming soon!", "5 new puzzles every month"
	if col != nil && col.ComingSoon {
		badge, title, sub = col.SoonBadge, col.Title+" coming soon!", col.SoonText
		if len(col.Puzzles) > 0 {
			title = "More " + strings.ToLower(col.Title) + " coming soon!"
		}
		if badge == "" {
			badge = "+"
		}
		if col.SoonTitle != "" {
			title = col.SoonTitle
		}
	}
	cy := r.Y + r.H*0.42
	FillCircle(dst, r.CX(), cy-14, 26, p.Accent)
	Text(dst, badge, Bold, 18, r.CX(), cy-14, p.AccentText, AlignCenter)
	Text(dst, title, Bold, 15, r.CX(), cy+36, p.Text, AlignCenter)
	if sub != "" {
		Text(dst, sub, Regular, 12.5, r.CX(), cy+58, p.TextDim, AlignCenter)
	}
}

func perim(r Rect, d float64) (float64, float64) {
	in := r.Inset(6)
	if d < in.W {
		return in.X + d, in.Y
	}
	d -= in.W
	if d < in.H {
		return in.Right(), in.Y + d
	}
	d -= in.H
	if d < in.W {
		return in.Right() - d, in.Bottom()
	}
	d -= in.W
	return in.X, in.Bottom() - math.Min(d, in.H)
}

// ---- difficulty modal

func (g *GalleryScene) modalRect() Rect {
	w := math.Min(860, SW-80)
	h := 470.0
	return Rect{SW/2 - w/2, SH/2 - h/2, w, h}
}

func (g *GalleryScene) updateModal() {
	r := g.modalRect()
	pz := g.selected
	closeFn := func() { g.selected = nil; g.confirm = -1 }
	g.widgets = append(g.widgets, &Button{ID: "md_close", R: Rect{r.Right() - 54, r.Y + 14, 40, 40}, Label: "X", Kind: BtnGhost, Size: 16, OnClick: closeFn})
	lx := r.X + r.W*0.5 + 10
	lw := r.Right() - 28 - lx
	y := r.Y + 96
	aspect := pz.Aspect()
	for i := range modes {
		mi := i
		rowR := Rect{lx, y, lw, 74}
		saved := g.has(pz.ID, mi)
		bw := 118.0
		if g.confirm == mi {
			g.widgets = append(g.widgets,
				&Button{ID: fmt.Sprintf("md_yes%d", mi), R: Rect{rowR.Right() - 2*bw - 8, rowR.Y + 16, bw, 42}, Label: "Start over", Kind: BtnDanger, Size: 13.5, OnClick: func() {
					deleteSave(saveKey(pz.ID, mi))
					game.SetScene(NewPlayScene(pz, mi, nil))
				}},
				&Button{ID: fmt.Sprintf("md_no%d", mi), R: Rect{rowR.Right() - bw, rowR.Y + 16, bw, 42}, Label: "Cancel", Kind: BtnSecondary, Size: 13.5, OnClick: func() { g.confirm = -1 }},
			)
		} else if saved {
			g.widgets = append(g.widgets,
				&Button{ID: fmt.Sprintf("md_cont%d", mi), R: Rect{rowR.Right() - 2*bw - 8, rowR.Y + 16, bw, 42}, Label: "Continue", Kind: BtnPrimary, Size: 13.5, OnClick: func() {
					if sv, err := readSave(saveKey(pz.ID, mi)); err == nil {
						game.SetScene(NewPlayScene(pz, mi, sv))
					}
				}},
				&Button{ID: fmt.Sprintf("md_re%d", mi), R: Rect{rowR.Right() - bw, rowR.Y + 16, bw, 42}, Label: "Restart", Kind: BtnSecondary, Size: 13.5, OnClick: func() { g.confirm = mi }},
			)
		} else {
			g.widgets = append(g.widgets, &Button{ID: fmt.Sprintf("md_play%d", mi), R: Rect{rowR.Right() - bw, rowR.Y + 16, bw, 42}, Label: "Play", Kind: BtnPrimary, Size: 14, OnClick: func() {
				game.SetScene(NewPlayScene(pz, mi, nil))
			}})
		}
		_ = aspect
		y += 84
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		closeFn()
	}
	for _, w := range g.widgets {
		w.Update()
	}
	if in.Pressed && !in.consumed && !r.Contains(in.MX, in.MY) {
		closeFn()
	}
}

func (g *GalleryScene) drawModal(dst *ebiten.Image) {
	p := P()
	FillRect(dst, Rect{0, 0, SW, SH}, alpha(p.Overlay, 0.5))
	r := g.modalRect()
	Panel(dst, r, 24)
	pz := g.selected
	img := Rect{r.X + 24, r.Y + 24, r.W*0.5 - 34, (r.W*0.5 - 34) * 0.5625}
	DrawImageRR(dst, pz.Thumb(), img, 16, 1)
	Text(dst, pz.Title, Bold, 22, r.X+28, img.Bottom()+36, p.Text, AlignLeft)
	colTitle := ""
	if c := findCollection(pz.Collection); c != nil {
		colTitle = c.Title
	}
	Text(dst, colTitle, Regular, 14, r.X+28, img.Bottom()+64, p.TextDim, AlignLeft)
	TextWrapped(dst, "Pick a difficulty. Your progress saves automatically, so you can stop any time and come back later.", Regular, 13.5, r.X+28, img.Bottom()+84, img.W-10, 21, p.TextDim, AlignLeft)

	lx := r.X + r.W*0.5 + 10
	lw := r.Right() - 28 - lx
	Text(dst, "Choose difficulty", Bold, 18, lx, r.Y+60, p.Text, AlignLeft)
	y := r.Y + 96
	aspect := pz.Aspect()
	accents := []color.RGBA{p.Teal, hex("#4E9F3D"), p.Accent, p.Plum}
	for i, m := range modes {
		rowR := Rect{lx, y, lw, 74}
		FillRR(dst, rowR, 14, p.Surface2)
		FillRR(dst, Rect{rowR.X, rowR.Y, 6, rowR.H}, 3, accents[i])
		n := modePieceCount(i, aspect)
		Text(dst, m.Name, Bold, 16, rowR.X+22, rowR.Y+26, p.Text, AlignLeft)
		sub := fmt.Sprintf("%d pieces", n)
		if best, ok := settings.Completed[pz.ID+"|"+m.Key]; ok {
			sub += "  -  best " + formatDuration(float64(best))
		}
		if pct, ok := g.saved[saveKey(pz.ID, i)]; ok {
			sub = fmt.Sprintf("%d pieces  -  %d%% done", n, pct)
		}
		Text(dst, sub, Regular, 13, rowR.X+22, rowR.Y+50, p.TextDim, AlignLeft)
		if g.confirm == i {
			Text(dst, "Lose progress?", Medium, 12.5, rowR.Right()-252, rowR.Y+8+0, p.Danger, AlignRight)
		}
		y += 84
	}
	for _, w := range g.widgets {
		w.Draw(dst)
	}
}
