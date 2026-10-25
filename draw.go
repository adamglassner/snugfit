package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// S is the device scale factor: UI is laid out in logical units and multiplied by S.
var S = 1.0

var whiteImg = func() *ebiten.Image {
	img := ebiten.NewImage(3, 3)
	img.Fill(color.White)
	return img
}()
var whiteSub = whiteImg.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

type Rect struct{ X, Y, W, H float64 }

func (r Rect) Contains(x, y float64) bool {
	return x >= r.X && y >= r.Y && x < r.X+r.W && y < r.Y+r.H
}
func (r Rect) Inset(d float64) Rect   { return Rect{r.X + d, r.Y + d, r.W - 2*d, r.H - 2*d} }
func (r Rect) CX() float64            { return r.X + r.W/2 }
func (r Rect) CY() float64            { return r.Y + r.H/2 }
func (r Rect) Right() float64         { return r.X + r.W }
func (r Rect) Bottom() float64        { return r.Y + r.H }
func (r Rect) Move(dx, dy float64) Rect { return Rect{r.X + dx, r.Y + dy, r.W, r.H} }

func rrPath(x, y, w, h, r float64) *vector.Path {
	p := &vector.Path{}
	r = math.Max(0, math.Min(r, math.Min(w, h)/2))
	fx, fy, fw, fh, fr := float32(x), float32(y), float32(w), float32(h), float32(r)
	if r <= 0.01 {
		p.MoveTo(fx, fy)
		p.LineTo(fx+fw, fy)
		p.LineTo(fx+fw, fy+fh)
		p.LineTo(fx, fy+fh)
		p.Close()
		return p
	}
	p.MoveTo(fx+fr, fy)
	p.LineTo(fx+fw-fr, fy)
	p.Arc(fx+fw-fr, fy+fr, fr, -math.Pi/2, 0, vector.Clockwise)
	p.LineTo(fx+fw, fy+fh-fr)
	p.Arc(fx+fw-fr, fy+fh-fr, fr, 0, math.Pi/2, vector.Clockwise)
	p.LineTo(fx+fr, fy+fh)
	p.Arc(fx+fr, fy+fh-fr, fr, math.Pi/2, math.Pi, vector.Clockwise)
	p.LineTo(fx, fy+fr)
	p.Arc(fx+fr, fy+fr, fr, math.Pi, 3*math.Pi/2, vector.Clockwise)
	p.Close()
	return p
}

func fillPath(dst *ebiten.Image, p *vector.Path, c color.RGBA) {
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	cr, cg, cb, ca := float32(c.R)/255, float32(c.G)/255, float32(c.B)/255, float32(c.A)/255
	for i := range vs {
		vs[i].SrcX, vs[i].SrcY = 1, 1
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = cr, cg, cb, ca
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true, FillRule: ebiten.FillRuleNonZero}
	dst.DrawTriangles(vs, is, whiteSub, op)
}

func strokePath(dst *ebiten.Image, p *vector.Path, width float64, c color.RGBA) {
	sop := &vector.StrokeOptions{Width: float32(width), LineJoin: vector.LineJoinRound, LineCap: vector.LineCapRound}
	vs, is := p.AppendVerticesAndIndicesForStroke(nil, nil, sop)
	cr, cg, cb, ca := float32(c.R)/255, float32(c.G)/255, float32(c.B)/255, float32(c.A)/255
	for i := range vs {
		vs[i].SrcX, vs[i].SrcY = 1, 1
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = cr, cg, cb, ca
	}
	dst.DrawTriangles(vs, is, whiteSub, &ebiten.DrawTrianglesOptions{AntiAlias: true})
}

// Logical-unit helpers (multiplied by S)

func FillRR(dst *ebiten.Image, r Rect, rad float64, c color.RGBA) {
	fillPath(dst, rrPath(r.X*S, r.Y*S, r.W*S, r.H*S, rad*S), c)
}

func StrokeRR(dst *ebiten.Image, r Rect, rad, width float64, c color.RGBA) {
	strokePath(dst, rrPath(r.X*S, r.Y*S, r.W*S, r.H*S, rad*S), width*S, c)
}

func FillRect(dst *ebiten.Image, r Rect, c color.RGBA) {
	vector.DrawFilledRect(dst, float32(r.X*S), float32(r.Y*S), float32(r.W*S), float32(r.H*S), c, false)
}

func FillCircle(dst *ebiten.Image, x, y, r float64, c color.RGBA) {
	vector.DrawFilledCircle(dst, float32(x*S), float32(y*S), float32(r*S), c, true)
}

func StrokeCircle(dst *ebiten.Image, x, y, r, w float64, c color.RGBA) {
	vector.StrokeCircle(dst, float32(x*S), float32(y*S), float32(r*S), float32(w*S), c, true)
}

func Line(dst *ebiten.Image, x0, y0, x1, y1, w float64, c color.RGBA) {
	vector.StrokeLine(dst, float32(x0*S), float32(y0*S), float32(x1*S), float32(y1*S), float32(w*S), c, true)
}

// Shadow draws a soft drop shadow under a rounded rect.
func Shadow(dst *ebiten.Image, r Rect, rad float64, strength float64) {
	sc := P().Shadow
	base := 0.05
	if settings.Theme == "dark" {
		base = 0.12
	}
	for i := 4; i >= 1; i-- {
		d := float64(i) * 2.5
		FillRR(dst, Rect{r.X - d*0.5, r.Y + d*0.6, r.W + d, r.H + d}, rad+d, alpha(sc, base*strength))
	}
}

// DrawImageRR draws img into rect r (logical) with rounded corners, cover-fitting.
func DrawImageRR(dst *ebiten.Image, img *ebiten.Image, r Rect, rad float64, opacity float64) {
	if img == nil {
		return
	}
	b := img.Bounds()
	iw, ih := float64(b.Dx()), float64(b.Dy())
	// cover fit
	sx0, sy0, sw, sh := 0.0, 0.0, iw, ih
	ra := r.W / r.H
	if iw/ih > ra {
		sw = ih * ra
		sx0 = (iw - sw) / 2
	} else {
		sh = iw / ra
		sy0 = (ih - sh) / 2
	}
	p := rrPath(r.X*S, r.Y*S, r.W*S, r.H*S, rad*S)
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	o := float32(opacity)
	for i := range vs {
		u := (float64(vs[i].DstX)/S - r.X) / r.W
		v := (float64(vs[i].DstY)/S - r.Y) / r.H
		vs[i].SrcX = float32(float64(b.Min.X) + sx0 + u*sw)
		vs[i].SrcY = float32(float64(b.Min.Y) + sy0 + v*sh)
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = o, o, o, o
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true, Filter: ebiten.FilterLinear}
	dst.DrawTriangles(vs, is, img, op)
}

// DrawImageFit draws an image scaled to fit inside r (logical), keeping aspect.
func DrawImageFit(dst *ebiten.Image, img *ebiten.Image, r Rect, opacity float64) Rect {
	b := img.Bounds()
	iw, ih := float64(b.Dx()), float64(b.Dy())
	sc := math.Min(r.W/iw, r.H/ih)
	w, h := iw*sc, ih*sc
	x, y := r.X+(r.W-w)/2, r.Y+(r.H-h)/2
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Scale(sc*S, sc*S)
	op.GeoM.Translate(x*S, y*S)
	op.ColorScale.ScaleAlpha(float32(opacity))
	dst.DrawImage(img, op)
	return Rect{x, y, w, h}
}

// ---- text

type Weight int

const (
	Regular Weight = iota
	Medium
	Bold
)

var fontSources [3]*text.GoTextFaceSource

type faceKey struct {
	w    Weight
	size float64
	s    float64
}

var faceCache = map[faceKey]*text.GoTextFace{}

func Face(w Weight, size float64) *text.GoTextFace {
	k := faceKey{w, size, S}
	if f, ok := faceCache[k]; ok {
		return f
	}
	f := &text.GoTextFace{Source: fontSources[w], Size: size * S}
	faceCache[k] = f
	return f
}

type Align int

const (
	AlignCenter Align = iota
	AlignLeft
	AlignRight
)

// Text draws a single line. y is the vertical center of the line.
func Text(dst *ebiten.Image, s string, w Weight, size float64, x, y float64, c color.RGBA, al Align) {
	f := Face(w, size)
	op := &text.DrawOptions{}
	op.GeoM.Translate(x*S, y*S)
	switch al {
	case AlignCenter:
		op.PrimaryAlign = text.AlignCenter
	case AlignRight:
		op.PrimaryAlign = text.AlignEnd
	}
	op.SecondaryAlign = text.AlignCenter
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, s, f, op)
}

func TextWidth(s string, w Weight, size float64) float64 {
	wd, _ := text.Measure(s, Face(w, size), 0)
	return wd / S
}

// WrapText splits s into lines that fit maxW.
func WrapText(s string, w Weight, size, maxW float64) []string {
	var lines []string
	words := splitWords(s)
	cur := ""
	for _, wd := range words {
		try := wd
		if cur != "" {
			try = cur + " " + wd
		}
		if TextWidth(try, w, size) > maxW && cur != "" {
			lines = append(lines, cur)
			cur = wd
		} else {
			cur = try
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

func splitWords(s string) []string {
	var out []string
	cur := []rune{}
	for _, r := range s {
		if r == ' ' || r == '\n' || r == '\t' {
			if len(cur) > 0 {
				out = append(out, string(cur))
				cur = cur[:0]
			}
			continue
		}
		cur = append(cur, r)
	}
	if len(cur) > 0 {
		out = append(out, string(cur))
	}
	return out
}

func TextWrapped(dst *ebiten.Image, s string, w Weight, size, x, y, maxW, lineH float64, c color.RGBA, al Align) float64 {
	lines := WrapText(s, w, size, maxW)
	for i, l := range lines {
		lx := x
		if al == AlignCenter {
			lx = x + maxW/2
		}
		Text(dst, l, w, size, lx, y+float64(i)*lineH+lineH/2, c, al)
	}
	return float64(len(lines)) * lineH
}

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }

func approach(cur, target, speed float64) float64 {
	return cur + (target-cur)*clamp(speed, 0, 1)
}

func rectToImage(r Rect) image.Rectangle {
	return image.Rect(int(r.X*S), int(r.Y*S), int(math.Ceil(r.Right()*S)), int(math.Ceil(r.Bottom()*S)))
}
