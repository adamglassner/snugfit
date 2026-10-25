package main

import (
	"image"
	"image/draw"
	"math"
	"math/rand"
	"runtime"
	"sync"

	xvector "golang.org/x/image/vector"
)

type Mode struct {
	Key    string
	Name   string
	Target int
}

var modes = []Mode{
	{"beginner", "Beginner", 100},
	{"easy", "Easy", 250},
	{"medium", "Medium", 500},
	{"advanced", "Advanced", 1000},
}

// gridFor picks columns/rows so the piece count is as close to target as
// possible while keeping pieces close to square.
func gridFor(target int, aspect float64) (int, int) {
	bestC, bestR, best := 0, 0, math.Inf(1)
	for c := 2; c <= 120; c++ {
		for _, r := range []int{int(math.Floor(float64(target) / float64(c))), int(math.Ceil(float64(target) / float64(c)))} {
			if r < 2 {
				continue
			}
			n := c * r
			pa := aspect * float64(r) / float64(c) // piece width / height
			var score float64
			if n == target && pa >= 0.85 && pa <= 1.2 {
				score = math.Abs(math.Log(pa)) * 0.1
			} else {
				score = math.Abs(float64(n-target))/float64(target)*25 + math.Abs(math.Log(pa))*4
				if n < target {
					score += 1.5
				}
			}
			if score < best {
				best, bestC, bestR = score, c, r
			}
		}
	}
	return bestC, bestR
}

func modePieceCount(mode int, aspect float64) int {
	c, r := gridFor(modes[mode].Target, aspect)
	return c * r
}

type pt struct{ X, Y float64 }
type cubic [4]pt
type curve []cubic

// tab template in (t along edge, n normal) units
var tabTemplate = []pt{
	{0, 0}, {0.35, 0}, // line
	{0.37, 0}, {0.39, 0.02}, {0.38, 0.06},
	{0.35, 0.13}, {0.38, 0.23}, {0.5, 0.23},
	{0.62, 0.23}, {0.65, 0.13}, {0.62, 0.06},
	{0.61, 0.02}, {0.63, 0}, {0.65, 0},
	{1, 0}, // line
}

func makeEdge(p0, p1 pt, normScale float64, bulge float64, rng *rand.Rand) curve {
	dx, dy := p1.X-p0.X, p1.Y-p0.Y
	L := math.Hypot(dx, dy)
	ux, uy := dx/L, dy/L
	// normal pointing toward +y for horizontal edges, +x for vertical (bulge sign decides direction)
	nx, ny := -uy, ux
	if nx < 0 || ny < 0 {
		nx, ny = -nx, -ny
	}
	s := 0.88 + rng.Float64()*0.2
	dt := (rng.Float64() - 0.5) * 0.1
	jit := make([]float64, len(tabTemplate))
	for i := range jit {
		jit[i] = (rng.Float64() - 0.5) * 0.03
	}
	wave1 := (rng.Float64() - 0.5) * 0.04
	wave2 := (rng.Float64() - 0.5) * 0.04
	m := func(i int) pt {
		t, n := tabTemplate[i].X, tabTemplate[i].Y
		if i > 0 && i < len(tabTemplate)-1 {
			t = 0.5 + (t-0.5)*s + dt + jit[i]*0.5
			n = (n + jit[i]*0.4) * s
		}
		n *= bulge
		return pt{p0.X + ux*t*L + nx*n*normScale, p0.Y + uy*t*L + ny*n*normScale}
	}
	off := func(t, n float64) pt {
		return pt{p0.X + ux*t*L + nx*n*normScale, p0.Y + uy*t*L + ny*n*normScale}
	}
	a, b := m(0), m(1)
	c := curve{}
	// gently waved straight sections
	c = append(c, cubic{a, off(0.12, wave1), off(0.24, wave1), b})
	for i := 1; i+3 < len(tabTemplate)-1; i += 3 {
		c = append(c, cubic{m(i), m(i + 1), m(i + 2), m(i + 3)})
	}
	e := m(len(tabTemplate) - 1)
	c = append(c, cubic{m(len(tabTemplate) - 2), off(0.76, wave2), off(0.88, wave2), e})
	return c
}

func straightEdge(p0, p1 pt) curve {
	return curve{{p0, {p0.X + (p1.X-p0.X)/3, p0.Y + (p1.Y-p0.Y)/3}, {p0.X + 2*(p1.X-p0.X)/3, p0.Y + 2*(p1.Y-p0.Y)/3}, p1}}
}

func (c curve) reversed() curve {
	out := make(curve, len(c))
	for i, s := range c {
		out[len(c)-1-i] = cubic{s[3], s[2], s[1], s[0]}
	}
	return out
}

type Geometry struct {
	Cols, Rows int
	BW, BH     float64
	PW, PH     float64
	H          [][]curve // (rows+1) x cols
	V          [][]curve // rows x (cols+1)
	Margin     float64
}

func newGeometry(bw, bh float64, cols, rows int, seed int64) *Geometry {
	g := &Geometry{Cols: cols, Rows: rows, BW: bw, BH: bh, PW: bw / float64(cols), PH: bh / float64(rows)}
	rng := rand.New(rand.NewSource(seed))
	ns := math.Min(g.PW, g.PH)
	g.Margin = math.Ceil(ns*0.28) + 2
	g.H = make([][]curve, rows+1)
	for r := 0; r <= rows; r++ {
		g.H[r] = make([]curve, cols)
		for c := 0; c < cols; c++ {
			p0 := pt{float64(c) * g.PW, float64(r) * g.PH}
			p1 := pt{float64(c+1) * g.PW, float64(r) * g.PH}
			if r == 0 || r == rows {
				g.H[r][c] = straightEdge(p0, p1)
			} else {
				b := 1.0
				if rng.Intn(2) == 0 {
					b = -1
				}
				g.H[r][c] = makeEdge(p0, p1, ns, b, rng)
			}
		}
	}
	g.V = make([][]curve, rows)
	for r := 0; r < rows; r++ {
		g.V[r] = make([]curve, cols+1)
		for c := 0; c <= cols; c++ {
			p0 := pt{float64(c) * g.PW, float64(r) * g.PH}
			p1 := pt{float64(c) * g.PW, float64(r+1) * g.PH}
			if c == 0 || c == cols {
				g.V[r][c] = straightEdge(p0, p1)
			} else {
				b := 1.0
				if rng.Intn(2) == 0 {
					b = -1
				}
				g.V[r][c] = makeEdge(p0, p1, ns, b, rng)
			}
		}
	}
	return g
}

// Outline returns the closed outline of piece (c, r) in board coordinates.
func (g *Geometry) Outline(c, r int) curve {
	var out curve
	out = append(out, g.H[r][c]...)
	out = append(out, g.V[r][c+1]...)
	out = append(out, g.H[r+1][c].reversed()...)
	out = append(out, g.V[r][c].reversed()...)
	return out
}

// PieceRect returns the integer image rect (board coords) that contains piece (c,r) incl. tabs.
func (g *Geometry) PieceRect(c, r int) image.Rectangle {
	x0 := int(math.Floor(float64(c)*g.PW - g.Margin))
	y0 := int(math.Floor(float64(r)*g.PH - g.Margin))
	x1 := int(math.Ceil(float64(c+1)*g.PW + g.Margin))
	y1 := int(math.Ceil(float64(r+1)*g.PH + g.Margin))
	return image.Rect(x0, y0, x1, y1)
}

type PieceRaster struct {
	Index int
	Rect  image.Rectangle // board coords
	RGBA  *image.RGBA
	Mask  []uint8
}

func rasterPiece(g *Geometry, src *image.RGBA, c, r int) PieceRaster {
	rect := g.PieceRect(c, r)
	w, h := rect.Dx(), rect.Dy()
	z := xvector.NewRasterizer(w, h)
	ol := g.Outline(c, r)
	ox, oy := float64(rect.Min.X), float64(rect.Min.Y)
	f := func(p pt) (float32, float32) { return float32(p.X - ox), float32(p.Y - oy) }
	x, y := f(ol[0][0])
	z.MoveTo(x, y)
	for _, s := range ol {
		bx, by := f(s[1])
		cx, cy := f(s[2])
		dx, dy := f(s[3])
		z.CubeTo(bx, by, cx, cy, dx, dy)
	}
	z.ClosePath()
	alphaImg := image.NewAlpha(image.Rect(0, 0, w, h))
	z.Draw(alphaImg, alphaImg.Bounds(), image.Opaque, image.Point{})
	mask := alphaImg.Pix
	stride := alphaImg.Stride

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	o := int(math.Max(1, math.Round(math.Min(g.PW, g.PH)/40)))
	sb := src.Bounds()
	get := func(x, y int) int {
		if x < 0 || y < 0 || x >= w || y >= h {
			return 0
		}
		return int(mask[y*stride+x])
	}
	for yy := 0; yy < h; yy++ {
		sy := rect.Min.Y + yy
		if sy < sb.Min.Y {
			sy = sb.Min.Y
		} else if sy >= sb.Max.Y {
			sy = sb.Max.Y - 1
		}
		for xx := 0; xx < w; xx++ {
			m := int(mask[yy*stride+xx])
			if m == 0 {
				continue
			}
			sx := rect.Min.X + xx
			if sx < sb.Min.X {
				sx = sb.Min.X
			} else if sx >= sb.Max.X {
				sx = sb.Max.X - 1
			}
			si := src.PixOffset(sx, sy)
			hi := m - get(xx-o, yy-o)
			lo := m - get(xx+o, yy+o)
			if hi < 0 {
				hi = 0
			}
			if lo < 0 {
				lo = 0
			}
			k := 1 + 0.32*float64(hi)/255 - 0.42*float64(lo)/255
			add := 22 * float64(hi) / 255
			cr := float64(src.Pix[si])*k + add
			cg := float64(src.Pix[si+1])*k + add
			cbv := float64(src.Pix[si+2])*k + add
			fm := float64(m) / 255
			di := yy*out.Stride + xx*4
			out.Pix[di] = uint8(clamp(cr, 0, 255) * fm)
			out.Pix[di+1] = uint8(clamp(cg, 0, 255) * fm)
			out.Pix[di+2] = uint8(clamp(cbv, 0, 255) * fm)
			out.Pix[di+3] = uint8(m)
		}
	}
	mk := make([]uint8, w*h)
	for yy := 0; yy < h; yy++ {
		copy(mk[yy*w:(yy+1)*w], mask[yy*stride:yy*stride+w])
	}
	return PieceRaster{Index: r*g.Cols + c, Rect: rect, RGBA: out, Mask: mk}
}

func toRGBA(img image.Image) *image.RGBA {
	if r, ok := img.(*image.RGBA); ok {
		return r
	}
	b := img.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), img, b.Min, draw.Src)
	return out
}

// generatePieces rasterizes all pieces in parallel, streaming results.
func generatePieces(g *Geometry, src *image.RGBA, out chan<- PieceRaster) {
	n := g.Cols * g.Rows
	jobs := make(chan int, n)
	for i := 0; i < n; i++ {
		jobs <- i
	}
	close(jobs)
	workers := runtime.NumCPU()
	if workers > 8 {
		workers = 8
	}
	if workers < 1 {
		workers = 1
	}
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				out <- rasterPiece(g, src, i%g.Cols, i/g.Cols)
			}
		}()
	}
	wg.Wait()
	close(out)
}
