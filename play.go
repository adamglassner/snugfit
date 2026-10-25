package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const hudH = 64.0

type Piece struct {
	C, R           int
	Img            *ebiten.Image
	Mask           []uint8
	IX, IY, IW, IH int
	X, Y           float64
	AX, AY         float64 // snap animation offset
	Group          int
	Locked         bool
	Edge           bool
}

type loadResult struct {
	rgba *image.RGBA
	geo  *Geometry
	err  error
}

type particle struct {
	x, y, vx, vy, rot, vr, size float64
	c                           color.RGBA
	life                        float64
}

type PlayScene struct {
	def  *PuzzleDef
	mode int
	seed int64
	geo  *Geometry
	src  *ebiten.Image

	pieces []*Piece
	groups map[int][]int
	order  []int

	// loading
	loading   bool
	loadErr   string
	resCh     chan loadResult
	rasterCh  chan PieceRaster
	received  int
	total     int
	pending   *SaveGame
	loadStart time.Time

	camX, camY, zoom float64
	camTX, camTY     float64
	zoomT            float64
	camAnim          bool

	dragging            bool
	dragGroup           int
	lastWX, lastWY      float64
	panning             bool
	panSX, panSY        float64
	panCX, panCY        float64
	spaceHeld           bool

	elapsed  float64
	placed   int
	snapDist float64

	showPreview, showGhost, edgesOnly bool

	paused   bool
	confirm  string
	complete bool
	compT    float64
	hidePanel bool
	confetti []particle
	bestTime bool

	dirty    bool
	saveT    float64
	toast    string
	toastT   float64
	hintT    float64
	widgets  []Widget
	hoverHUD bool
}

func NewPlayScene(def *PuzzleDef, mode int, save *SaveGame) *PlayScene {
	s := &PlayScene{def: def, mode: mode, pending: save, groups: map[int][]int{}}
	s.seed = time.Now().UnixNano()
	s.showGhost = settings.GhostDefault
	if save != nil {
		s.seed = save.Seed
		s.elapsed = save.Elapsed
		s.showGhost = save.Ghost
	}
	s.startLoading()
	settings.LastSave = saveKey(def.ID, mode)
	saveSettings()
	return s
}

func (s *PlayScene) startLoading() {
	s.loading = true
	s.loadStart = time.Now()
	s.resCh = make(chan loadResult, 1)
	s.rasterCh = make(chan PieceRaster, 256)
	def, mode, seed, save := s.def, s.mode, s.seed, s.pending
	go func() {
		img, err := loadImageFile(assetPath("puzzles", def.Image))
		if err != nil {
			s.resCh <- loadResult{err: err}
			return
		}
		rgba := toRGBA(img)
		b := rgba.Bounds()
		aspect := float64(b.Dx()) / float64(b.Dy())
		cols, rows := gridFor(modes[mode].Target, aspect)
		if save != nil && save.Cols > 1 && save.Rows > 1 && len(save.Pieces) == save.Cols*save.Rows {
			cols, rows = save.Cols, save.Rows
		}
		geo := newGeometry(float64(b.Dx()), float64(b.Dy()), cols, rows, seed)
		s.resCh <- loadResult{rgba: rgba, geo: geo}
		generatePieces(geo, rgba, s.rasterCh)
	}()
}

func (s *PlayScene) updateLoading() {
	if s.geo == nil {
		select {
		case r := <-s.resCh:
			if r.err != nil {
				s.loadErr = "Could not open the puzzle picture: " + r.err.Error()
				return
			}
			s.geo = r.geo
			s.src = ebiten.NewImageFromImage(r.rgba)
			s.total = s.geo.Cols * s.geo.Rows
			s.pieces = make([]*Piece, s.total)
			s.snapDist = math.Min(s.geo.PW, s.geo.PH) * 0.24
		default:
		}
		return
	}
	start := time.Now()
	for time.Since(start) < 12*time.Millisecond {
		select {
		case pr, ok := <-s.rasterCh:
			if !ok {
				if s.received >= s.total {
					s.finishLoading()
				}
				return
			}
			c, r := pr.Index%s.geo.Cols, pr.Index/s.geo.Cols
			p := &Piece{C: c, R: r, Img: ebiten.NewImageFromImage(pr.RGBA), Mask: pr.Mask,
				IX: pr.Rect.Min.X, IY: pr.Rect.Min.Y, IW: pr.Rect.Dx(), IH: pr.Rect.Dy(),
				Group: pr.Index, Edge: c == 0 || r == 0 || c == s.geo.Cols-1 || r == s.geo.Rows-1}
			p.X, p.Y = float64(c)*s.geo.PW, float64(r)*s.geo.PH
			s.pieces[pr.Index] = p
			s.received++
		default:
			return
		}
	}
}

func (s *PlayScene) finishLoading() {
	s.loading = false
	s.groups = map[int][]int{}
	if sv := s.pending; sv != nil && len(sv.Pieces) == s.total {
		for i, sp := range sv.Pieces {
			p := s.pieces[i]
			p.X, p.Y, p.Group, p.Locked = sp.X, sp.Y, sp.Group, sp.Locked
		}
		for i, p := range s.pieces {
			if p.Group < 0 || p.Group >= s.total {
				p.Group = i
			}
			s.groups[p.Group] = append(s.groups[p.Group], i)
		}
	} else {
		for i, p := range s.pieces {
			p.Group = i
			s.groups[i] = []int{i}
		}
		all := make([]int, s.total)
		for i := range all {
			all[i] = i
		}
		s.scatter(all, false)
	}
	s.order = make([]int, 0, s.total)
	for i := range s.pieces {
		s.order = append(s.order, i)
	}
	rand.Shuffle(len(s.order), func(i, j int) { s.order[i], s.order[j] = s.order[j], s.order[i] })
	s.countPlaced()
	if sv := s.pending; sv != nil && sv.Zoom > 0 {
		s.camX, s.camY, s.zoom = sv.CamX, sv.CamY, sv.Zoom*S
		s.camTX, s.camTY, s.zoomT = s.camX, s.camY, s.zoom
	} else {
		s.fit(false)
	}
	s.pending = nil
	s.hintT = 14
	playSFX("whoosh")
}

// scatter places the given pieces in a ring of slots around the board.
func (s *PlayScene) scatter(indices []int, edgesFirst bool) {
	g := s.geo
	sx, sy := g.PW*1.38, g.PH*1.38
	gap := math.Max(g.PW, g.PH) * 0.9
	type slot struct{ x, y, d float64 }
	// obstacles: bounding boxes of multi-piece loose groups
	var obst []Rect
	for _, gl := range s.groups {
		if len(gl) < 2 || s.pieces[gl[0]].Locked {
			continue
		}
		minX, minY, maxX, maxY := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
		for _, i := range gl {
			p := s.pieces[i]
			minX, minY = math.Min(minX, p.X), math.Min(minY, p.Y)
			maxX, maxY = math.Max(maxX, p.X+g.PW), math.Max(maxY, p.Y+g.PH)
		}
		obst = append(obst, Rect{minX - g.PW*0.3, minY - g.PH*0.3, maxX - minX + g.PW*0.6, maxY - minY + g.PH*0.6})
	}
	board := Rect{-gap, -gap, g.BW + 2*gap, g.BH + 2*gap}
	var slots []slot
	for k := 2; k < 400; k++ {
		slots = slots[:0]
		ex, ey := float64(k)*sx, float64(k)*sy*0.8
		for y := -ey; y <= g.BH+ey; y += sy {
			for x := -ex; x <= g.BW+ex; x += sx {
				r := Rect{x, y, sx, sy}
				if rectOverlap(r, board) {
					continue
				}
				bad := false
				for _, o := range obst {
					if rectOverlap(r, o) {
						bad = true
						break
					}
				}
				if bad {
					continue
				}
				cx, cy := x+sx/2, y+sy/2
				dx := math.Max(0, math.Max(-cx, cx-g.BW)) / g.BW
				dy := math.Max(0, math.Max(-cy, cy-g.BH)) / g.BH
				slots = append(slots, slot{cx, cy, math.Hypot(dx*1.0, dy*1.6) + rand.Float64()*0.02})
			}
		}
		if len(slots) >= len(indices) {
			break
		}
	}
	sort.Slice(slots, func(i, j int) bool { return slots[i].d < slots[j].d })
	if len(slots) > len(indices) {
		slots = slots[:len(indices)]
	}
	ord := append([]int(nil), indices...)
	rand.Shuffle(len(ord), func(i, j int) { ord[i], ord[j] = ord[j], ord[i] })
	sl := append([]slot(nil), slots...)
	if edgesFirst {
		sort.SliceStable(ord, func(i, j int) bool { return s.pieces[ord[i]].Edge && !s.pieces[ord[j]].Edge })
		// slots already sorted nearest-first; shuffle within edge/non-edge halves for variety
		ne := 0
		for _, i := range ord {
			if s.pieces[i].Edge {
				ne++
			}
		}
		if ne < len(sl) {
			rand.Shuffle(ne, func(i, j int) { sl[i], sl[j] = sl[j], sl[i] })
			rest := sl[ne:]
			rand.Shuffle(len(rest), func(i, j int) { rest[i], rest[j] = rest[j], rest[i] })
		}
	} else {
		rand.Shuffle(len(sl), func(i, j int) { sl[i], sl[j] = sl[j], sl[i] })
	}
	for k, i := range ord {
		if k >= len(sl) {
			break
		}
		p := s.pieces[i]
		nx := sl[k].x - g.PW/2 + (rand.Float64()-0.5)*g.PW*0.18
		ny := sl[k].y - g.PH/2 + (rand.Float64()-0.5)*g.PH*0.18
		if !s.loading {
			p.AX += p.X - nx
			p.AY += p.Y - ny
		}
		p.X, p.Y = nx, ny
	}
}

func rectOverlap(a, b Rect) bool {
	return a.X < b.Right() && a.Right() > b.X && a.Y < b.Bottom() && a.Bottom() > b.Y
}

func (s *PlayScene) countPlaced() {
	n := 0
	for _, p := range s.pieces {
		if p.Locked {
			n++
		}
	}
	s.placed = n
}

// ---- camera

func (s *PlayScene) worldBounds() (float64, float64, float64, float64) {
	minX, minY, maxX, maxY := 0.0, 0.0, s.geo.BW, s.geo.BH
	for _, p := range s.pieces {
		if s.edgesOnly && !p.Locked && !s.groupHasEdge(p.Group) {
			continue
		}
		minX = math.Min(minX, p.X)
		minY = math.Min(minY, p.Y)
		maxX = math.Max(maxX, p.X+s.geo.PW)
		maxY = math.Max(maxY, p.Y+s.geo.PH)
	}
	pad := math.Max(s.geo.PW, s.geo.PH) * 0.6
	return minX - pad, minY - pad, maxX + pad, maxY + pad
}

func (s *PlayScene) fit(animate bool) {
	x0, y0, x1, y1 := s.worldBounds()
	area := Rect{0, hudH, SW, SH - hudH}
	z := math.Min(area.W*S/(x1-x0), area.H*S/(y1-y0))
	cx, cy := (x0+x1)/2, (y0+y1)/2
	tx := cx - (area.CX()*S)/z
	ty := cy - (area.CY()*S)/z
	if animate {
		s.camTX, s.camTY, s.zoomT, s.camAnim = tx, ty, z, true
	} else {
		s.camX, s.camY, s.zoom = tx, ty, z
		s.camTX, s.camTY, s.zoomT = tx, ty, z
	}
}

func (s *PlayScene) zoomLimits() (float64, float64) {
	maxZ := 420 * S / math.Min(s.geo.PW, s.geo.PH)
	minZ := math.Min(SW*S/(s.geo.BW*4), maxZ/8)
	return minZ, maxZ
}

func (s *PlayScene) zoomAt(factor, mx, my float64) {
	minZ, maxZ := s.zoomLimits()
	nz := clamp(s.zoom*factor, minZ, maxZ)
	wx := s.camX + mx*S/s.zoom
	wy := s.camY + my*S/s.zoom
	s.zoom = nz
	s.camX = wx - mx*S/nz
	s.camY = wy - my*S/nz
	s.camAnim = false
}

func (s *PlayScene) mouseWorld() (float64, float64) {
	return s.camX + in.MX*S/s.zoom, s.camY + in.MY*S/s.zoom
}

// ---- groups

func (s *PlayScene) groupHasEdge(g int) bool {
	for _, i := range s.groups[g] {
		if s.pieces[i].Edge {
			return true
		}
	}
	return false
}

func (s *PlayScene) visible(p *Piece) bool {
	if p.Locked || !s.edgesOnly {
		return true
	}
	return s.groupHasEdge(p.Group)
}

func (s *PlayScene) pick(wx, wy float64) int {
	for k := len(s.order) - 1; k >= 0; k-- {
		i := s.order[k]
		p := s.pieces[i]
		if p.Locked || !s.visible(p) {
			continue
		}
		lx := int(wx - (p.X + float64(p.IX) - float64(p.C)*s.geo.PW))
		ly := int(wy - (p.Y + float64(p.IY) - float64(p.R)*s.geo.PH))
		if lx < 0 || ly < 0 || lx >= p.IW || ly >= p.IH {
			continue
		}
		if p.Mask[ly*p.IW+lx] > 60 {
			return i
		}
	}
	return -1
}

func (s *PlayScene) bringToTop(g int) {
	in := map[int]bool{}
	for _, i := range s.groups[g] {
		in[i] = true
	}
	out := s.order[:0]
	var top []int
	for _, i := range s.order {
		if in[i] {
			top = append(top, i)
		} else {
			out = append(out, i)
		}
	}
	s.order = append(out, top...)
}

func (s *PlayScene) moveGroup(g int, dx, dy float64, anim bool) {
	for _, i := range s.groups[g] {
		p := s.pieces[i]
		p.X += dx
		p.Y += dy
		if anim {
			p.AX -= dx
			p.AY -= dy
		}
	}
}

func (s *PlayScene) merge(into, from int) int {
	locked := s.pieces[s.groups[into][0]].Locked
	for _, i := range s.groups[from] {
		s.pieces[i].Group = into
		if locked {
			s.pieces[i].Locked = true
		}
	}
	s.groups[into] = append(s.groups[into], s.groups[from]...)
	delete(s.groups, from)
	return into
}

var dirs = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func (s *PlayScene) drop(g int) {
	connected := false
	for {
		found := false
		for _, i := range s.groups[g] {
			p := s.pieces[i]
			for _, d := range dirs {
				c2, r2 := p.C+d[0], p.R+d[1]
				if c2 < 0 || r2 < 0 || c2 >= s.geo.Cols || r2 >= s.geo.Rows {
					continue
				}
				q := s.pieces[r2*s.geo.Cols+c2]
				if q.Group == g || !s.visible(q) {
					continue
				}
				ex := p.X + float64(d[0])*s.geo.PW
				ey := p.Y + float64(d[1])*s.geo.PH
				if math.Hypot(q.X-ex, q.Y-ey) < s.snapDist {
					s.moveGroup(g, q.X-ex, q.Y-ey, true)
					g = s.merge(q.Group, g)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			break
		}
		connected = true
	}
	locked := s.pieces[s.groups[g][0]].Locked
	if !locked {
		for _, i := range s.groups[g] {
			p := s.pieces[i]
			tx, ty := float64(p.C)*s.geo.PW, float64(p.R)*s.geo.PH
			if math.Hypot(p.X-tx, p.Y-ty) < s.snapDist {
				s.moveGroup(g, tx-p.X, ty-p.Y, true)
				for _, j := range s.groups[g] {
					s.pieces[j].Locked = true
				}
				locked = true
				break
			}
		}
		if locked {
			playSFX("place")
		}
	} else if connected {
		playSFX("place")
	}
	if !locked {
		if connected {
			playSFX("snap")
		} else {
			playSFX("drop")
		}
	}
	// if everything is one loose group, put it on the board
	if len(s.groups[g]) == s.total && !s.pieces[s.groups[g][0]].Locked {
		p := s.pieces[s.groups[g][0]]
		s.moveGroup(g, float64(p.C)*s.geo.PW-p.X, float64(p.R)*s.geo.PH-p.Y, true)
		for _, j := range s.groups[g] {
			s.pieces[j].Locked = true
		}
	}
	s.countPlaced()
	s.dirty = true
	if s.placed == s.total {
		s.onComplete()
	}
}

func (s *PlayScene) onComplete() {
	s.complete = true
	s.compT = 0
	s.edgesOnly = false
	key := s.def.ID + "|" + modes[s.mode].Key
	best, ok := settings.Completed[key]
	if !ok || int64(s.elapsed) < best {
		settings.Completed[key] = int64(s.elapsed)
		s.bestTime = ok
	}
	deleteSave(saveKey(s.def.ID, s.mode))
	saveSettings()
	playSFX("complete")
	s.fit(true)
	cols := []color.RGBA{hex("#F07A2A"), hex("#2C9A92"), hex("#F2C14E"), hex("#9B7FD1"), hex("#E84A5F"), hex("#4CC9F0")}
	for i := 0; i < 260; i++ {
		s.confetti = append(s.confetti, particle{
			x: rand.Float64() * SW, y: -rand.Float64()*SH*0.6 - 20,
			vx: (rand.Float64() - 0.5) * 3, vy: 2 + rand.Float64()*3,
			rot: rand.Float64() * 6, vr: (rand.Float64() - 0.5) * 0.3, size: 6 + rand.Float64()*8,
			c: cols[rand.Intn(len(cols))], life: 1,
		})
	}
}

// ---- saving

func (s *PlayScene) buildSave() *SaveGame {
	sv := &SaveGame{PuzzleID: s.def.ID, Mode: s.mode, Cols: s.geo.Cols, Rows: s.geo.Rows, Seed: s.seed,
		Elapsed: s.elapsed, Placed: s.placed, Total: s.total, CamX: s.camX, CamY: s.camY, Zoom: s.zoom / S, Ghost: s.showGhost}
	sv.Pieces = make([]SavedPiece, s.total)
	for i, p := range s.pieces {
		sv.Pieces[i] = SavedPiece{X: p.X, Y: p.Y, Group: p.Group, Locked: p.Locked}
	}
	return sv
}

func (s *PlayScene) saveNow(showToast bool) {
	if s.loading || s.complete || s.geo == nil {
		return
	}
	if err := writeSave(s.buildSave()); err == nil {
		s.dirty = false
		if showToast {
			s.showToast("Progress saved")
		}
	} else {
		s.showToast("Could not save: " + err.Error())
	}
}

func (s *PlayScene) showToast(t string) { s.toast, s.toastT = t, 2.4 }

func (s *PlayScene) exitToMenu() {
	s.saveNow(false)
	game.SetScene(NewMenuScene())
}

// ---- update

func (s *PlayScene) Update() {
	dt := 1.0 / 60
	if s.loadErr != "" {
		s.widgets = []Widget{&Button{ID: "lerr", R: Rect{SW/2 - 90, SH/2 + 40, 180, 46}, Label: "Back to menu", OnClick: func() { game.SetScene(NewMenuScene()) }}}
		for _, w := range s.widgets {
			w.Update()
		}
		return
	}
	if s.loading {
		s.updateLoading()
		return
	}
	if s.toastT > 0 {
		s.toastT -= dt
	}
	if s.hintT > 0 {
		s.hintT -= dt
	}
	// snap animation decay
	for _, p := range s.pieces {
		if p.AX != 0 || p.AY != 0 {
			p.AX *= 0.72
			p.AY *= 0.72
			if math.Abs(p.AX) < 0.3 && math.Abs(p.AY) < 0.3 {
				p.AX, p.AY = 0, 0
			}
		}
	}
	if s.camAnim {
		s.camX = approach(s.camX, s.camTX, 0.15)
		s.camY = approach(s.camY, s.camTY, 0.15)
		s.zoom = approach(s.zoom, s.zoomT, 0.15)
		if math.Abs(s.zoom-s.zoomT) < 0.0005*s.zoomT && math.Abs(s.camX-s.camTX) < 0.5 {
			s.camAnim = false
		}
	}

	s.buildHUD()
	for _, w := range s.widgets {
		w.Update()
	}

	if s.complete {
		s.compT += dt
		for i := range s.confetti {
			c := &s.confetti[i]
			c.x += c.vx
			c.y += c.vy
			c.vy += 0.03
			c.vx *= 0.995
			c.rot += c.vr
			if c.y > SH+20 {
				c.life = 0
			}
		}
		s.handleCamera(true)
		return
	}
	if s.paused {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			if s.confirm != "" {
				s.confirm = ""
			} else {
				s.paused = false
			}
		}
		return
	}
	if ebiten.IsFocused() {
		s.elapsed += dt
	}
	s.handleKeys()
	s.handleCamera(false)
	s.handlePieces()
	s.saveT += dt
	if s.dirty && s.saveT > 20 {
		s.saveT = 0
		s.saveNow(false)
	}
}

func (s *PlayScene) handleKeys() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.paused = true
		s.saveNow(false)
		playSFX("click")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		s.showGhost = !s.showGhost
		s.dirty = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		s.showPreview = !s.showPreview
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		s.edgesOnly = !s.edgesOnly
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		s.fit(true)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		s.tidy()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		s.saveNow(true)
	}
	pan := 14 * S / s.zoom
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		s.camX -= pan
		s.camAnim = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		s.camX += pan
		s.camAnim = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		s.camY -= pan
		s.camAnim = false
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		s.camY += pan
		s.camAnim = false
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) || inpututil.IsKeyJustPressed(ebiten.KeyKPAdd) {
		s.zoomAt(1.25, SW/2, (SH+hudH)/2)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) || inpututil.IsKeyJustPressed(ebiten.KeyKPSubtract) {
		s.zoomAt(0.8, SW/2, (SH+hudH)/2)
	}
	s.spaceHeld = ebiten.IsKeyPressed(ebiten.KeySpace)
}

func (s *PlayScene) handleCamera(viewOnly bool) {
	if in.Wheel != 0 && !s.hoverHUD {
		s.zoomAt(math.Pow(1.15, in.Wheel), in.MX, in.MY)
	}
	startPan := in.RPressed || in.MPressed || (in.Pressed && (s.spaceHeld || viewOnly) && !in.consumed && in.MY > hudH)
	if startPan && !s.dragging {
		s.panning = true
		s.panSX, s.panSY = in.MX, in.MY
		s.panCX, s.panCY = s.camX, s.camY
		s.camAnim = false
		if in.Pressed {
			in.consumed = true
		}
	}
	if s.panning {
		ebiten.SetCursorShape(ebiten.CursorShapeMove)
		if in.RDown || in.MDown || in.Down {
			s.camX = s.panCX - (in.MX-s.panSX)*S/s.zoom
			s.camY = s.panCY - (in.MY-s.panSY)*S/s.zoom
		} else {
			s.panning = false
		}
	}
}

func (s *PlayScene) handlePieces() {
	wx, wy := s.mouseWorld()
	if s.dragging {
		ebiten.SetCursorShape(ebiten.CursorShapeMove)
		if in.Down {
			s.moveGroup(s.dragGroup, wx-s.lastWX, wy-s.lastWY, false)
			s.lastWX, s.lastWY = wx, wy
			// edge auto-scroll
			edge := 24.0
			sp := 10 * S / s.zoom
			if in.MX < edge {
				s.camX -= sp
				s.moveGroup(s.dragGroup, -sp, 0, false)
			} else if in.MX > SW-edge {
				s.camX += sp
				s.moveGroup(s.dragGroup, sp, 0, false)
			}
			if in.MY < hudH+edge && in.MY > hudH-10 {
				s.camY -= sp
				s.moveGroup(s.dragGroup, 0, -sp, false)
			} else if in.MY > SH-edge {
				s.camY += sp
				s.moveGroup(s.dragGroup, 0, sp, false)
			}
			s.lastWX, s.lastWY = s.mouseWorld()
		} else {
			s.dragging = false
			s.drop(s.dragGroup)
		}
		return
	}
	if s.panning || s.hoverHUD {
		return
	}
	hit := s.pick(wx, wy)
	if hit >= 0 {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	}
	if in.Pressed && !in.consumed {
		in.consumed = true
		if hit >= 0 {
			s.dragging = true
			s.dragGroup = s.pieces[hit].Group
			s.lastWX, s.lastWY = wx, wy
			s.bringToTop(s.dragGroup)
			playSFX("pick")
			s.hintT = math.Min(s.hintT, 3)
		} else {
			// drag empty table = pan
			s.panning = true
			s.panSX, s.panSY = in.MX, in.MY
			s.panCX, s.panCY = s.camX, s.camY
			s.camAnim = false
		}
	}
}

func (s *PlayScene) tidy() {
	var loose []int
	for _, gl := range s.groups {
		if len(gl) == 1 && !s.pieces[gl[0]].Locked {
			loose = append(loose, gl[0])
		}
	}
	if len(loose) == 0 {
		return
	}
	s.scatter(loose, true)
	playSFX("whoosh")
	s.dirty = true
	s.showToast("Loose pieces tidied - edge pieces are closest to the board")
}

func (s *PlayScene) restart() {
	deleteSave(saveKey(s.def.ID, s.mode))
	game.SetScene(NewPlayScene(s.def, s.mode, nil))
}

// ---- HUD

func (s *PlayScene) buildHUD() {
	s.widgets = s.widgets[:0]
	s.hoverHUD = in.MY < hudH
	if s.paused {
		s.buildPause()
		return
	}
	s.widgets = append(s.widgets, &Button{ID: "menu", R: Rect{12, 12, 100, 40}, Label: "< Menu", Kind: BtnGhost, OnClick: func() {
		s.paused = true
		s.saveNow(false)
	}})
	type tb struct {
		id, label string
		on        bool
		fn        func()
	}
	items := []tb{
		{"preview", "Preview", s.showPreview, func() { s.showPreview = !s.showPreview }},
		{"ghost", "Ghost", s.showGhost, func() { s.showGhost = !s.showGhost; s.dirty = true }},
		{"edges", "Edges only", s.edgesOnly, func() { s.edgesOnly = !s.edgesOnly }},
		{"tidy", "Tidy", false, s.tidy},
		{"fit", "Fit view", false, func() { s.fit(true) }},
	}
	if s.complete {
		items = items[:2]
	}
	x := SW - 12
	for i := len(items) - 1; i >= 0; i-- {
		it := items[i]
		w := TextWidth(it.label, Medium, 14) + 30
		x -= w
		b := &Button{ID: "tb_" + it.id, R: Rect{x, 12, w, 40}, Label: it.label, Kind: BtnSecondary, Size: 14, Selected: it.on, OnClick: it.fn}
		s.widgets = append(s.widgets, b)
		x -= 8
	}
	// zoom buttons
	zr := Rect{SW - 60, SH - 120, 44, 44}
	s.widgets = append(s.widgets,
		&Button{ID: "zin", R: zr, Label: "+", Kind: BtnSecondary, Size: 22, OnClick: func() { s.zoomAt(1.25, SW/2, (SH+hudH)/2) }},
		&Button{ID: "zout", R: zr.Move(0, 52), Label: "-", Kind: BtnSecondary, Size: 22, OnClick: func() { s.zoomAt(0.8, SW/2, (SH+hudH)/2) }},
	)
	if in.Hover(Rect{zr.X - 6, zr.Y - 6, 56, 108}) {
		s.hoverHUD = true
	}
	if s.showPreview {
		pr := s.previewRect()
		if in.Hover(pr) {
			s.hoverHUD = true
		}
	}
	if s.complete && !s.hidePanel && s.compT > 0.6 {
		s.buildComplete()
	} else if s.complete && s.hidePanel {
		s.widgets = append(s.widgets, &Button{ID: "showpanel", R: Rect{SW/2 - 90, SH - 70, 180, 46}, Label: "Show results", Kind: BtnPrimary, OnClick: func() { s.hidePanel = false }})
	}
}

func (s *PlayScene) previewRect() Rect {
	w := math.Min(360, SW*0.3)
	h := w * s.geo.BH / s.geo.BW
	return Rect{SW - w - 16, hudH + 12, w, h}
}

func (s *PlayScene) buildPause() {
	pw, ph := 380.0, 380.0
	r := Rect{SW/2 - pw/2, SH/2 - ph/2, pw, ph}
	bx := r.X + 40
	bw := pw - 80
	if s.confirm == "restart" {
		s.widgets = append(s.widgets,
			&Button{ID: "c_yes", R: Rect{bx, r.Y + 210, bw, 48}, Label: "Yes, start over", Kind: BtnDanger, OnClick: s.restart},
			&Button{ID: "c_no", R: Rect{bx, r.Y + 270, bw, 48}, Label: "Cancel", Kind: BtnSecondary, OnClick: func() { s.confirm = "" }},
		)
		return
	}
	y := r.Y + 110
	add := func(id, label string, kind ButtonKind, fn func()) {
		s.widgets = append(s.widgets, &Button{ID: id, R: Rect{bx, y, bw, 48}, Label: label, Kind: kind, OnClick: fn})
		y += 58
	}
	add("p_resume", "Resume", BtnPrimary, func() { s.paused = false })
	add("p_opts", "Options", BtnSecondary, func() {
		game.SetScene(NewOptionsScene(func() { s.paused = true; game.SetScene(s) }))
	})
	add("p_restart", "Restart puzzle", BtnSecondary, func() { s.confirm = "restart" })
	add("p_exit", "Save & exit to menu", BtnSecondary, s.exitToMenu)
}

func (s *PlayScene) buildComplete() {
	pw, ph := 460.0, 360.0
	r := Rect{SW/2 - pw/2, SH/2 - ph/2, pw, ph}
	bw := (pw - 100) / 2
	s.widgets = append(s.widgets,
		&Button{ID: "cp_next", R: Rect{r.X + 40, r.Bottom() - 128, pw - 80, 48}, Label: "Play another puzzle", Kind: BtnPrimary, OnClick: func() { game.SetScene(NewGalleryScene()) }},
		&Button{ID: "cp_admire", R: Rect{r.X + 40, r.Bottom() - 68, bw, 44}, Label: "Admire it", Kind: BtnSecondary, OnClick: func() { s.hidePanel = true }},
		&Button{ID: "cp_menu", R: Rect{r.X + 60 + bw, r.Bottom() - 68, bw, 44}, Label: "Main menu", Kind: BtnSecondary, OnClick: func() { game.SetScene(NewMenuScene()) }},
	)
}

// ---- draw

func (s *PlayScene) w2s(wx, wy float64) (float64, float64) {
	return (wx - s.camX) * s.zoom, (wy - s.camY) * s.zoom
}

func (s *PlayScene) Draw(dst *ebiten.Image) {
	p := P()
	dst.Fill(tableColor())
	if s.loadErr != "" {
		Text(dst, s.loadErr, Medium, 16, SW/2, SH/2-10, p.Text, AlignCenter)
		for _, w := range s.widgets {
			w.Draw(dst)
		}
		return
	}
	if s.loading {
		s.drawLoading(dst)
		return
	}
	g := s.geo
	// board
	bx, by := s.w2s(0, 0)
	bw, bh := g.BW*s.zoom, g.BH*s.zoom
	boardCol := lerpColor(tableColor(), color.RGBA{0, 0, 0, 255}, 0.12)
	if !isDarkTable() {
		boardCol = lerpColor(tableColor(), color.RGBA{0, 0, 0, 255}, 0.07)
	}
	pad := 6 * S
	fillPath(dst, rrPath(bx-pad, by-pad, bw+2*pad, bh+2*pad, 10*S), alpha(color.RGBA{0, 0, 0, 255}, 0.10))
	fillPath(dst, rrPath(bx, by, bw, bh, 2*S), boardCol)
	if s.showGhost || s.complete {
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		op.GeoM.Scale(s.zoom, s.zoom)
		op.GeoM.Translate(bx, by)
		a := float32(0.24)
		if s.complete {
			a = 1
		}
		op.ColorScale.ScaleAlpha(a)
		dst.DrawImage(s.src, op)
	}
	if !s.complete || s.compT < 1.2 {
		s.drawPieces(dst, true)
		s.drawPieces(dst, false)
	}
	if s.complete && s.compT < 1.2 {
		// fade in the clean picture over the pieces
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		op.GeoM.Scale(s.zoom, s.zoom)
		op.GeoM.Translate(bx, by)
		op.ColorScale.ScaleAlpha(float32(clamp(s.compT/1.2, 0, 1)))
		dst.DrawImage(s.src, op)
	}
	s.drawHUD(dst)
}

func (s *PlayScene) drawPieces(dst *ebiten.Image, locked bool) {
	g := s.geo
	sw, sh := float64(dst.Bounds().Dx()), float64(dst.Bounds().Dy())
	filter := ebiten.FilterLinear
	shadowOff := math.Max(2, math.Min(6, s.zoom*math.Min(g.PW, g.PH)*0.04))
	for _, i := range s.order {
		pc := s.pieces[i]
		if pc.Locked != locked || !s.visible(pc) {
			continue
		}
		wx := pc.X + pc.AX + float64(pc.IX) - float64(pc.C)*g.PW
		wy := pc.Y + pc.AY + float64(pc.IY) - float64(pc.R)*g.PH
		x, y := s.w2s(wx, wy)
		w, h := float64(pc.IW)*s.zoom, float64(pc.IH)*s.zoom
		if x > sw || y > sh || x+w < 0 || y+h < 0 {
			continue
		}
		lifted := s.dragging && pc.Group == s.dragGroup
		if !locked {
			so := shadowOff
			sa := float32(0.32)
			if lifted {
				so *= 2.6
				sa = 0.28
			}
			op := &ebiten.DrawImageOptions{Filter: filter}
			op.GeoM.Scale(s.zoom, s.zoom)
			op.GeoM.Translate(x+so*0.6, y+so)
			op.ColorScale.Scale(0, 0, 0, sa)
			dst.DrawImage(pc.Img, op)
		}
		op := &ebiten.DrawImageOptions{Filter: filter}
		op.GeoM.Scale(s.zoom, s.zoom)
		if lifted {
			op.GeoM.Translate(x-1*S, y-2*S)
		} else {
			op.GeoM.Translate(x, y)
		}
		dst.DrawImage(pc.Img, op)
	}
}

func (s *PlayScene) drawLoading(dst *ebiten.Image) {
	p := P()
	dst.Fill(p.Bg)
	r := Rect{SW/2 - 230, SH/2 - 150, 460, 300}
	Panel(dst, r, 22)
	DrawImageRR(dst, s.def.Thumb(), Rect{r.X + 30, r.Y + 30, r.W - 60, 140}, 12, 1)
	Text(dst, s.def.Title, Bold, 20, r.CX(), r.Y+198, p.Text, AlignCenter)
	frac := 0.0
	label := "Opening picture..."
	if s.total > 0 {
		frac = float64(s.received) / float64(s.total)
		label = fmt.Sprintf("Cutting %d pieces...  %d%%", s.total, int(frac*100))
	}
	Text(dst, label, Regular, 14, r.CX(), r.Y+228, p.TextDim, AlignCenter)
	bar := Rect{r.X + 40, r.Y + 252, r.W - 80, 10}
	FillRR(dst, bar, 5, p.Track)
	FillRR(dst, Rect{bar.X, bar.Y, math.Max(10, bar.W*frac), bar.H}, 5, p.Accent)
}

func (s *PlayScene) drawHUD(dst *ebiten.Image) {
	p := P()
	// top bar
	FillRect(dst, Rect{0, 0, SW, hudH}, alpha(p.Surface, 0.94))
	FillRect(dst, Rect{0, hudH - 1, SW, 1}, p.Border)
	tx := 128.0
	Text(dst, s.def.Title, Bold, 16, tx, 24, p.Text, AlignLeft)
	info := fmt.Sprintf("%s  -  %d pieces", modes[s.mode].Name, s.total)
	Text(dst, info, Regular, 12.5, tx, 44, p.TextDim, AlignLeft)
	// progress
	px := tx + math.Max(TextWidth(s.def.Title, Bold, 16), TextWidth(info, Regular, 12.5)) + 36
	prog := float64(s.placed) / float64(s.total)
	maxPX := SW - 560
	if px < maxPX {
		bw := math.Min(220, maxPX-px)
		if bw > 80 {
			Text(dst, fmt.Sprintf("%d / %d placed", s.placed, s.total), Medium, 13, px, 22, p.Text, AlignLeft)
			Text(dst, formatDuration(s.elapsed), Medium, 13, px+bw, 22, p.TextDim, AlignRight)
			FillRR(dst, Rect{px, 36, bw, 8}, 4, p.Track)
			FillRR(dst, Rect{px, 36, math.Max(8, bw*prog), 8}, 4, p.Teal)
		}
	}
	for _, w := range s.widgets {
		if s.paused || (s.complete && !s.hidePanel && s.compT > 0.6) {
			break
		}
		w.Draw(dst)
	}
	if s.showPreview {
		pr := s.previewRect()
		Panel(dst, pr.Inset(-6), 14)
		DrawImageRR(dst, s.src, pr, 9, 1)
	}
	if s.hintT > 0 && !s.complete {
		a := clamp(s.hintT, 0, 1)
		msg := "Drag pieces to connect them  -  drag the table or right-click to pan  -  scroll to zoom  -  Esc for menu"
		w := TextWidth(msg, Medium, 13.5) + 40
		r := Rect{SW/2 - w/2, SH - 64, w, 40}
		FillRR(dst, r, 20, alpha(p.Overlay, 0.78*a))
		Text(dst, msg, Medium, 13.5, r.CX(), r.CY(), alpha(color.RGBA{255, 255, 255, 255}, a), AlignCenter)
	}
	if s.toastT > 0 {
		a := clamp(s.toastT/0.4, 0, 1)
		w := TextWidth(s.toast, Medium, 14) + 40
		r := Rect{SW/2 - w/2, hudH + 16, w, 40}
		FillRR(dst, r, 20, alpha(p.Teal, a))
		Text(dst, s.toast, Medium, 14, r.CX(), r.CY(), alpha(color.RGBA{255, 255, 255, 255}, a), AlignCenter)
	}
	if s.paused {
		s.drawPause(dst)
	}
	if s.complete {
		s.drawComplete(dst)
	}
}

func (s *PlayScene) drawPause(dst *ebiten.Image) {
	p := P()
	FillRect(dst, Rect{0, 0, SW, SH}, alpha(p.Overlay, 0.45))
	pw, ph := 380.0, 380.0
	r := Rect{SW/2 - pw/2, SH/2 - ph/2, pw, ph}
	Panel(dst, r, 22)
	if s.confirm == "restart" {
		Text(dst, "Start over?", Bold, 24, r.CX(), r.Y+60, p.Text, AlignCenter)
		TextWrapped(dst, "Your progress on this puzzle will be lost and the pieces will be mixed up again.", Regular, 14, r.X+40, r.Y+100, r.W-80, 22, p.TextDim, AlignCenter)
	} else {
		Text(dst, "Paused", Bold, 26, r.CX(), r.Y+52, p.Text, AlignCenter)
		Text(dst, fmt.Sprintf("%d of %d pieces placed  -  %s", s.placed, s.total, formatDuration(s.elapsed)), Regular, 13.5, r.CX(), r.Y+84, p.TextDim, AlignCenter)
	}
	for _, w := range s.widgets {
		w.Draw(dst)
	}
}

func (s *PlayScene) drawComplete(dst *ebiten.Image) {
	p := P()
	for _, c := range s.confetti {
		if c.life <= 0 {
			continue
		}
		w := c.size * math.Abs(math.Cos(c.rot))
		FillRect(dst, Rect{c.x - w/2, c.y - c.size*0.3, math.Max(1, w), c.size * 0.6}, c.c)
	}
	if s.hidePanel || s.compT < 0.6 {
		if s.hidePanel {
			for _, w := range s.widgets {
				if b, ok := w.(*Button); ok && b.ID == "showpanel" {
					b.Draw(dst)
				}
			}
		}
		return
	}
	a := clamp((s.compT-0.6)/0.4, 0, 1)
	FillRect(dst, Rect{0, 0, SW, SH}, alpha(p.Overlay, 0.35*a))
	pw, ph := 460.0, 360.0
	r := Rect{SW/2 - pw/2, SH/2 - ph/2 + (1-a)*30, pw, ph}
	Panel(dst, r, 24)
	Text(dst, "Puzzle complete!", Bold, 30, r.CX(), r.Y+58, p.Text, AlignCenter)
	Text(dst, s.def.Title, Medium, 16, r.CX(), r.Y+98, p.Accent, AlignCenter)
	Text(dst, fmt.Sprintf("%s  -  %d pieces", modes[s.mode].Name, s.total), Regular, 14, r.CX(), r.Y+124, p.TextDim, AlignCenter)
	Text(dst, "Time  "+formatDuration(s.elapsed), Bold, 20, r.CX(), r.Y+168, p.Text, AlignCenter)
	if s.bestTime {
		Text(dst, "New best time!", Medium, 13, r.CX(), r.Y+196, p.Teal, AlignCenter)
	}
	for _, w := range s.widgets {
		w.Draw(dst)
	}
}
