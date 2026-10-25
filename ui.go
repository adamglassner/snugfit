package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ---- input

type Input struct {
	MX, MY            float64 // logical
	Pressed, Released bool
	Down              bool
	RPressed, RDown   bool
	MPressed, MDown   bool
	Wheel             float64
	consumed          bool
	Tick              int64
}

var in Input

// fakeQueue lets automated tests feed mouse input one frame at a time.
var fakeQueue []Input

func (i *Input) poll() {
	if len(fakeQueue) > 0 {
		f := fakeQueue[0]
		fakeQueue = fakeQueue[1:]
		f.Tick = i.Tick + 1
		*i = f
		return
	}
	x, y := ebiten.CursorPosition()
	i.MX, i.MY = float64(x)/S, float64(y)/S
	i.Pressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	i.Released = inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
	i.Down = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	i.RPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	i.RDown = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	i.MPressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle)
	i.MDown = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	_, wy := ebiten.Wheel()
	i.Wheel = wy
	i.consumed = false
	i.Tick++
}

// Click returns true once if a press happened inside r and wasn't consumed.
func (i *Input) Click(r Rect) bool {
	if i.Pressed && !i.consumed && r.Contains(i.MX, i.MY) {
		i.consumed = true
		return true
	}
	return false
}

func (i *Input) Hover(r Rect) bool { return r.Contains(i.MX, i.MY) }

var hoverAnim = map[string]float64{}
var hoverWas = map[string]bool{}

func animHover(id string, hovering bool) float64 {
	v := hoverAnim[id]
	t := 0.0
	if hovering {
		t = 1
	}
	v = approach(v, t, 0.25)
	hoverAnim[id] = v
	if hovering && !hoverWas[id] {
		playSFX("hover")
	}
	hoverWas[id] = hovering
	return v
}

// ---- widgets

type Widget interface {
	Update()
	Draw(dst *ebiten.Image)
}

type ButtonKind int

const (
	BtnPrimary ButtonKind = iota
	BtnSecondary
	BtnGhost
	BtnDanger
	BtnTeal
)

type Button struct {
	ID       string
	R        Rect
	Label    string
	Sub      string
	Kind     ButtonKind
	Size     float64
	Disabled bool
	Selected bool
	Badge    bool
	Align    Align
	OnClick  func()
	hover    float64
}

func (b *Button) Update() {
	h := !b.Disabled && in.Hover(b.R)
	b.hover = animHover(b.ID, h)
	if h {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	}
	if !b.Disabled && in.Click(b.R) {
		playSFX("click")
		if b.OnClick != nil {
			b.OnClick()
		}
	}
}

func (b *Button) Draw(dst *ebiten.Image) {
	p := P()
	r := b.R
	size := b.Size
	if size == 0 {
		size = 15
	}
	var bg, fg, border color.RGBA
	switch b.Kind {
	case BtnPrimary:
		bg = lerpColor(p.Accent, p.AccentHover, b.hover)
		fg = p.AccentText
	case BtnTeal:
		bg = lerpColor(p.Teal, lerpColor(p.Teal, color.RGBA{0, 0, 0, 255}, 0.15), b.hover)
		fg = color.RGBA{255, 255, 255, 255}
	case BtnDanger:
		bg = lerpColor(p.Danger, lerpColor(p.Danger, color.RGBA{0, 0, 0, 255}, 0.15), b.hover)
		fg = color.RGBA{255, 255, 255, 255}
	case BtnSecondary:
		bg = lerpColor(p.Surface, p.Surface2, b.hover)
		fg = p.Text
		border = p.Border
	case BtnGhost:
		bg = alpha(p.Surface2, b.hover)
		fg = p.Text
	}
	if b.Selected {
		bg = alpha(p.Accent, 0.16)
		border = p.Accent
		fg = p.Text
	}
	if b.Disabled {
		bg = alpha(bg, 0.45)
		fg = alpha(fg, 0.5)
	}
	rad := math.Min(12, r.H/2)
	if b.Kind == BtnPrimary || b.Kind == BtnTeal || b.Kind == BtnDanger {
		if !b.Disabled {
			Shadow(dst, r.Move(0, -b.hover), rad, 0.6+b.hover*0.6)
		}
	}
	rr := r.Move(0, -b.hover*1.5)
	FillRR(dst, rr, rad, bg)
	if border.A > 0 {
		StrokeRR(dst, rr, rad, 1.2, border)
	}
	wgt := Medium
	if b.Kind == BtnPrimary || b.Kind == BtnTeal || b.Kind == BtnDanger {
		wgt = Bold
	}
	cy := rr.CY()
	if b.Sub != "" {
		cy -= size * 0.55
	}
	switch b.Align {
	case AlignLeft:
		Text(dst, b.Label, wgt, size, rr.X+18, cy, fg, AlignLeft)
		if b.Sub != "" {
			Text(dst, b.Sub, Regular, size*0.78, rr.X+18, cy+size*1.25, alpha(fg, 0.7), AlignLeft)
		}
	default:
		Text(dst, b.Label, wgt, size, rr.CX(), cy, fg, AlignCenter)
		if b.Sub != "" {
			Text(dst, b.Sub, Regular, size*0.78, rr.CX(), cy+size*1.25, alpha(fg, 0.7), AlignCenter)
		}
	}
	if b.Badge {
		FillCircle(dst, rr.Right()-8, rr.Y+8, 5, p.Danger)
	}
}

type Slider struct {
	ID       string
	R        Rect
	Value    *float64
	OnChange func()
	drag     bool
}

var activeDrag string

func (s *Slider) Update() {
	hit := Rect{s.R.X - 10, s.R.Y - 12, s.R.W + 20, s.R.H + 24}
	h := in.Hover(hit)
	animHover(s.ID, h)
	if h {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	}
	if in.Pressed && h && !in.consumed {
		in.consumed = true
		activeDrag = s.ID
	}
	if activeDrag == s.ID {
		if in.Down {
			v := clamp((in.MX-s.R.X)/s.R.W, 0, 1)
			if v != *s.Value {
				*s.Value = v
				if s.OnChange != nil {
					s.OnChange()
				}
			}
		} else {
			activeDrag = ""
			playSFX("click")
		}
	}
}

func (s *Slider) Draw(dst *ebiten.Image) {
	p := P()
	v := *s.Value
	cy := s.R.CY()
	FillRR(dst, Rect{s.R.X, cy - 4, s.R.W, 8}, 4, p.Track)
	FillRR(dst, Rect{s.R.X, cy - 4, s.R.W * v, 8}, 4, p.Accent)
	kx := s.R.X + s.R.W*v
	hv := hoverAnim[s.ID]
	FillCircle(dst, kx, cy+1.5, 11+hv*2, alpha(p.Shadow, 0.18))
	FillCircle(dst, kx, cy, 11+hv*2, p.Surface)
	StrokeCircle(dst, kx, cy, 11+hv*2, 2.5, p.Accent)
}

type Toggle struct {
	ID       string
	R        Rect // pill rect (approx 46x26)
	Value    *bool
	OnChange func()
	anim     float64
}

func (t *Toggle) Update() {
	h := in.Hover(t.R.Inset(-6))
	animHover(t.ID, h)
	if h {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	}
	if in.Click(t.R.Inset(-6)) {
		*t.Value = !*t.Value
		playSFX("click")
		if t.OnChange != nil {
			t.OnChange()
		}
	}
	tg := 0.0
	if *t.Value {
		tg = 1
	}
	k := "tg_" + t.ID
	hoverAnim[k] = approach(hoverAnim[k], tg, 0.3)
	t.anim = hoverAnim[k]
}

func (t *Toggle) Draw(dst *ebiten.Image) {
	p := P()
	bg := lerpColor(p.Track, p.Teal, t.anim)
	FillRR(dst, t.R, t.R.H/2, bg)
	kx := t.R.X + t.R.H/2 + (t.R.W-t.R.H)*t.anim
	FillCircle(dst, kx, t.R.CY()+1, t.R.H/2-3, alpha(color.RGBA{0, 0, 0, 255}, 0.15))
	FillCircle(dst, kx, t.R.CY(), t.R.H/2-3, color.RGBA{255, 255, 255, 255})
}

type Segmented struct {
	ID       string
	R        Rect
	Options  []string
	Selected int
	OnSelect func(i int)
}

func (s *Segmented) Update() {
	n := len(s.Options)
	w := s.R.W / float64(n)
	for i := 0; i < n; i++ {
		r := Rect{s.R.X + float64(i)*w, s.R.Y, w, s.R.H}
		h := in.Hover(r)
		animHover(s.ID+string(rune('0'+i)), h)
		if h {
			ebiten.SetCursorShape(ebiten.CursorShapePointer)
		}
		if in.Click(r) && i != s.Selected {
			playSFX("click")
			s.Selected = i
			if s.OnSelect != nil {
				s.OnSelect(i)
			}
		}
	}
}

func (s *Segmented) Draw(dst *ebiten.Image) {
	p := P()
	FillRR(dst, s.R, 12, p.Surface2)
	n := len(s.Options)
	w := s.R.W / float64(n)
	for i, o := range s.Options {
		r := Rect{s.R.X + float64(i)*w, s.R.Y, w, s.R.H}.Inset(4)
		if i == s.Selected {
			Shadow(dst, r, 9, 0.5)
			FillRR(dst, r, 9, p.Surface)
			Text(dst, o, Bold, 14, r.CX(), r.CY(), p.Text, AlignCenter)
		} else {
			hv := hoverAnim[s.ID+string(rune('0'+i))]
			if hv > 0.01 {
				FillRR(dst, r, 9, alpha(p.Surface, hv*0.5))
			}
			Text(dst, o, Medium, 14, r.CX(), r.CY(), p.TextDim, AlignCenter)
		}
	}
}

// Panel draws a card surface.
func Panel(dst *ebiten.Image, r Rect, rad float64) {
	p := P()
	Shadow(dst, r, rad, 1)
	FillRR(dst, r, rad, p.Surface)
	StrokeRR(dst, r, rad, 1, p.Border)
}

// ScrollState is a simple vertical scroller.
type ScrollState struct {
	Y, Target, Max float64
}

func (s *ScrollState) Update(area Rect) {
	if in.Hover(area) && in.Wheel != 0 {
		s.Target -= in.Wheel * 60
	}
	s.Target = clamp(s.Target, 0, math.Max(0, s.Max))
	s.Y = approach(s.Y, s.Target, 0.3)
	if math.Abs(s.Y-s.Target) < 0.3 {
		s.Y = s.Target
	}
}

func (s *ScrollState) DrawBar(dst *ebiten.Image, area Rect) {
	if s.Max <= 0 {
		return
	}
	p := P()
	total := area.H + s.Max
	h := math.Max(40, area.H*area.H/total)
	y := area.Y + (area.H-h)*(s.Y/s.Max)
	FillRR(dst, Rect{area.Right() - 6, y, 4, h}, 2, alpha(p.TextDim, 0.5))
}
