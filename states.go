package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// stateTint gives each state card a gentle color from the app palette.
func stateTint(i int) color.RGBA {
	p := P()
	cols := []color.RGBA{p.Teal, p.Accent, p.Plum, hex("#4E9F3D"), hex("#2F6FB0"), hex("#C0504D")}
	return cols[i%len(cols)]
}

func drawSoonPill(dst *ebiten.Image, x, y float64, label string, align Align) {
	p := P()
	pw := TextWidth(label, Bold, 11.5) + 20
	if align == AlignCenter {
		x -= pw / 2
	} else if align == AlignRight {
		x -= pw
	}
	FillRR(dst, Rect{x, y - 12, pw, 24}, 12, p.Accent)
	Text(dst, label, Bold, 11.5, x+pw/2, y, p.AccentText, AlignCenter)
}

// drawCategoryCard draws the single card that opens a category of categories.
func drawCategoryCard(dst *ebiten.Image, r Rect, c *Collection) {
	p := P()
	hv := hoverAnim["cat_"+c.ID]
	r = r.Move(0, -hv*4)
	Shadow(dst, r, 18, 1+hv)
	FillRR(dst, r, 18, p.Surface)
	if hv > 0.01 {
		StrokeRR(dst, r, 18, 2, alpha(p.Accent, hv))
	} else {
		StrokeRR(dst, r, 18, 1, p.Border)
	}
	img := Rect{r.X + 10, r.Y + 10, r.W - 20, r.W*0.5625 - 10}
	FillRR(dst, img, 12, p.Teal)
	// a mosaic of state silhouettes
	n := len(c.Subcategories)
	gcols, grows := 10, (n+9)/10
	cw := img.W / float64(gcols)
	chh := img.H / float64(grows)
	s := math.Min(cw, chh) * 0.86
	for i, sc := range c.Subcategories {
		ic := sc.IconImage()
		if ic == nil {
			continue
		}
		cx := img.X + cw*(float64(i%gcols)+0.5)
		cy := img.Y + chh*(float64(i/gcols)+0.5)
		b := ic.Bounds()
		k := s / float64(b.Dx())
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		op.GeoM.Scale(k*S, k*S)
		op.GeoM.Translate((cx-s/2)*S, (cy-s/2)*S)
		a := float32(0.92)
		op.ColorScale.Scale(a, a, a*0.96, a)
		dst.DrawImage(ic, op)
	}
	if c.PuzzleCount() == 0 {
		drawSoonPill(dst, img.X+10, img.Y+22, "Coming soon", AlignLeft)
	}
	Text(dst, c.Title, Bold, 15, r.X+16, img.Bottom()+24, p.Text, AlignLeft)
	sub := fmt.Sprintf("%d states", n)
	if pc := c.PuzzleCount(); pc > 0 {
		sub = fmt.Sprintf("%d states  -  %d puzzles", n, pc)
	}
	Text(dst, sub, Regular, 12, r.X+16, img.Bottom()+44, p.TextDim, AlignLeft)
	Text(dst, "Open >", Medium, 12.5, r.Right()-18, img.Bottom()+44, p.Accent, AlignRight)
}

// ================= Category screen (e.g. the 50 states) =================

type CategoryScene struct {
	cat     *Collection
	widgets []Widget
	scroll  ScrollState
	cards   []Rect
	toast   string
	toastT  float64
}

func NewCategoryScene(c *Collection) *CategoryScene { return &CategoryScene{cat: c} }

func (s *CategoryScene) metrics() (float64, float64, float64, float64, int) {
	margin := math.Max(40, SW*0.05)
	gap := 16.0
	cols := int((SW - 2*margin + gap) / (176 + gap))
	if cols < 3 {
		cols = 3
	}
	if cols > 8 {
		cols = 8
	}
	cw := (SW - 2*margin - gap*float64(cols-1)) / float64(cols)
	return margin, gap, cw, cw * 1.08, cols
}

func (s *CategoryScene) layout() float64 {
	s.cards = s.cards[:0]
	margin, gap, cw, ch, cols := s.metrics()
	y0 := 170.0 - s.scroll.Y
	for i := range s.cat.Subcategories {
		r := Rect{margin + float64(i%cols)*(cw+gap), y0 + float64(i/cols)*(ch+gap), cw, ch}
		s.cards = append(s.cards, r)
	}
	rows := (len(s.cat.Subcategories) + cols - 1) / cols
	return 170 + float64(rows)*(ch+gap) + 20
}

func (s *CategoryScene) goBack() { game.SetScene(NewGalleryScene()) }

func (s *CategoryScene) Update() {
	s.widgets = s.widgets[:0]
	header(&s.widgets, "", s.goBack)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.goBack()
		return
	}
	h := s.layout()
	s.scroll.Max = h - SH
	s.scroll.Update(Rect{0, 90, SW, SH - 90})
	if s.toastT > 0 {
		s.toastT -= 1.0 / 60
	}
	for _, w := range s.widgets {
		w.Update()
	}
	for i, r := range s.cards {
		sc := s.cat.Subcategories[i]
		hov := in.Hover(r) && in.MY > 140
		animHover("st_"+sc.ID, hov)
		if hov {
			ebiten.SetCursorShape(ebiten.CursorShapePointer)
		}
		if hov && in.Click(r) {
			playSFX("click")
			if len(sc.Puzzles) > 0 {
				game.SetScene(NewSubcategoryGallery(s.cat, sc))
			} else {
				s.toast = sc.Title + " puzzles are coming soon!"
				s.toastT = 2.2
			}
		}
	}
}

func (s *CategoryScene) Draw(dst *ebiten.Image) {
	p := P()
	dst.Fill(p.Bg)
	drawFloatPieces(dst)
	s.layout()
	for i, r := range s.cards {
		if r.Bottom() < 130 || r.Y > SH {
			continue
		}
		sc := s.cat.Subcategories[i]
		hv := hoverAnim["st_"+sc.ID]
		rr := r.Move(0, -hv*3)
		Shadow(dst, rr, 16, 0.8+hv)
		FillRR(dst, rr, 16, p.Surface)
		if hv > 0.01 {
			StrokeRR(dst, rr, 16, 2, alpha(p.Accent, hv))
		} else {
			StrokeRR(dst, rr, 16, 1, p.Border)
		}
		tint := stateTint(i)
		art := Rect{rr.X + 10, rr.Y + 10, rr.W - 20, rr.H - 74}
		FillRR(dst, art, 12, lerpColor(p.Surface2, tint, 0.12))
		if ic := sc.IconImage(); ic != nil {
			sz := math.Min(art.W, art.H) * 0.84
			b := ic.Bounds()
			k := sz / float64(b.Dx())
			op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
			op.GeoM.Scale(k*S, k*S)
			op.GeoM.Translate((art.CX()-sz/2)*S, (art.CY()-sz/2)*S)
			op.ColorScale.ScaleWithColor(tint)
			dst.DrawImage(ic, op)
		}
		Text(dst, sc.Abbr, Bold, 11, art.Right()-10, art.Y+14, alpha(tint, 0.9), AlignRight)
		name := sc.Title
		size := 14.5
		for size > 11 && TextWidth(name, Bold, size) > rr.W-24 {
			size -= 0.5
		}
		Text(dst, name, Bold, size, rr.CX(), art.Bottom()+20, p.Text, AlignCenter)
		if len(sc.Puzzles) == 0 {
			drawSoonPill(dst, rr.CX(), art.Bottom()+46, "Coming soon", AlignCenter)
		} else {
			Text(dst, fmt.Sprintf("%d puzzles", len(sc.Puzzles)), Medium, 12.5, rr.CX(), art.Bottom()+46, p.Teal, AlignCenter)
		}
	}
	FillRect(dst, Rect{0, 0, SW, 140}, alpha(p.Bg, 0.96))
	drawHeader(dst, s.cat.Title, s.cat.Subtitle)
	margin, _, _, _, _ := s.metrics()
	if s.cat.PuzzleCount() == 0 {
		drawSoonPill(dst, margin, 116, "Coming soon", AlignLeft)
		Text(dst, "Every state is getting its own puzzles. Check back for updates!", Regular, 14, margin+120, 116, p.TextDim, AlignLeft)
	}
	for _, w := range s.widgets {
		w.Draw(dst)
	}
	s.scroll.DrawBar(dst, Rect{0, 140, SW - 4, SH - 140})
	if s.toastT > 0 {
		a := clamp(s.toastT/0.4, 0, 1)
		w := TextWidth(s.toast, Medium, 14) + 40
		r := Rect{SW/2 - w/2, SH - 70, w, 42}
		FillRR(dst, r, 21, alpha(p.Accent, a))
		Text(dst, s.toast, Medium, 14, r.CX(), r.CY(), alpha(p.AccentText, a), AlignCenter)
	}
}
