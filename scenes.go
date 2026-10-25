package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ================= Saved puzzles =================

type SavesScene struct {
	widgets []Widget
	saves   []SaveInfo
	scroll  ScrollState
	confirm string
}

func NewSavesScene() *SavesScene { return &SavesScene{saves: listSaves()} }

func (s *SavesScene) listRect() Rect {
	w := math.Min(900, SW-80)
	return Rect{SW/2 - w/2, 110, w, SH - 130}
}

func (s *SavesScene) Update() {
	s.widgets = s.widgets[:0]
	header(&s.widgets, "", func() { game.SetScene(NewMenuScene()) })
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		game.SetScene(NewMenuScene())
	}
	lr := s.listRect()
	rowH := 116.0
	s.scroll.Max = float64(len(s.saves))*(rowH+14) - lr.H
	s.scroll.Update(lr)
	y := lr.Y - s.scroll.Y
	for i := range s.saves {
		si := s.saves[i]
		r := Rect{lr.X, y, lr.W, rowH}
		y += rowH + 14
		if r.Bottom() < 90 || r.Y > SH {
			continue
		}
		bw := 120.0
		if s.confirm == si.Key {
			s.widgets = append(s.widgets,
				&Button{ID: "sv_yes" + si.Key, R: Rect{r.Right() - 2*bw - 30, r.CY() - 21, bw, 42}, Label: "Delete", Kind: BtnDanger, Size: 14, OnClick: func() {
					deleteSave(si.Key)
					s.saves = listSaves()
					s.confirm = ""
				}},
				&Button{ID: "sv_no" + si.Key, R: Rect{r.Right() - bw - 20, r.CY() - 21, bw, 42}, Label: "Keep", Kind: BtnSecondary, Size: 14, OnClick: func() { s.confirm = "" }},
			)
		} else {
			s.widgets = append(s.widgets,
				&Button{ID: "sv_go" + si.Key, R: Rect{r.Right() - 2*bw - 30, r.CY() - 21, bw, 42}, Label: "Continue", Kind: BtnPrimary, Size: 14, OnClick: func() {
					if sv, err := readSave(si.Key); err == nil {
						game.SetScene(NewPlayScene(si.Puzzle, sv.Mode, sv))
					}
				}},
				&Button{ID: "sv_del" + si.Key, R: Rect{r.Right() - bw - 20, r.CY() - 21, bw, 42}, Label: "Delete", Kind: BtnSecondary, Size: 14, OnClick: func() { s.confirm = si.Key }},
			)
		}
	}
	for _, w := range s.widgets {
		w.Update()
	}
}

func (s *SavesScene) Draw(dst *ebiten.Image) {
	p := P()
	dst.Fill(p.Bg)
	drawFloatPieces(dst)
	lr := s.listRect()
	if len(s.saves) == 0 {
		r := Rect{SW/2 - 260, SH/2 - 110, 520, 200}
		Panel(dst, r, 22)
		Text(dst, "No saved puzzles yet", Bold, 22, r.CX(), r.Y+70, p.Text, AlignCenter)
		TextWrapped(dst, "Start any puzzle and your progress is saved automatically. It will show up here so you can finish it later.", Regular, 14, r.X+40, r.Y+100, r.W-80, 22, p.TextDim, AlignCenter)
	}
	y := lr.Y - s.scroll.Y
	rowH := 116.0
	for _, si := range s.saves {
		r := Rect{lr.X, y, lr.W, rowH}
		y += rowH + 14
		if r.Bottom() < 90 || r.Y > SH {
			continue
		}
		Panel(dst, r, 18)
		img := Rect{r.X + 12, r.Y + 12, (rowH - 24) * 16 / 9, rowH - 24}
		DrawImageRR(dst, si.Puzzle.Thumb(), img, 12, 1)
		tx := img.Right() + 22
		g := si.Game
		Text(dst, si.Puzzle.Title, Bold, 17, tx, r.Y+30, p.Text, AlignLeft)
		Text(dst, fmt.Sprintf("%s  -  %d pieces  -  played %s", modes[g.Mode].Name, g.Total, timeAgo(g.Updated)), Regular, 13, tx, r.Y+56, p.TextDim, AlignLeft)
		bw := math.Max(80, r.Right()-300-tx)
		bar := Rect{tx, r.Y + 80, bw, 8}
		FillRR(dst, bar, 4, p.Track)
		FillRR(dst, Rect{bar.X, bar.Y, math.Max(8, bar.W*si.Percent), 8}, 4, p.Teal)
		Text(dst, fmt.Sprintf("%d%%  -  %s", int(si.Percent*100), formatDuration(g.Elapsed)), Medium, 12.5, tx, r.Y+100, p.TextDim, AlignLeft)
	}
	FillRect(dst, Rect{0, 0, SW, 96}, alpha(p.Bg, 0.96))
	drawHeader(dst, "Saved puzzles", "Pick up right where you left off")
	for _, w := range s.widgets {
		w.Draw(dst)
	}
	s.scroll.DrawBar(dst, lr)
}

// ================= Options =================

type OptionsScene struct {
	widgets []Widget
	back    func()
	scroll  ScrollState
	musicV  float64
	sfxV    float64
	contentH float64
}

func NewOptionsScene(back func()) *OptionsScene {
	o := &OptionsScene{back: back, musicV: settings.MusicVolume, sfxV: settings.SFXVolume}
	return o
}

func (o *OptionsScene) goBack() {
	saveSettings()
	if o.back != nil {
		o.back()
	} else {
		game.SetScene(NewMenuScene())
	}
}

func (o *OptionsScene) col() Rect {
	w := math.Min(760, SW-80)
	return Rect{SW/2 - w/2, 110 - o.scroll.Y, w, 0}
}

// section rects (computed from a fixed layout)
func (o *OptionsScene) sections() (Rect, Rect, Rect) {
	c := o.col()
	audio := Rect{c.X, c.Y, c.W, 250}
	disp := Rect{c.X, audio.Bottom() + 20, c.W, 290}
	upd := Rect{c.X, disp.Bottom() + 20, c.W, 210}
	o.contentH = upd.Bottom() + o.scroll.Y + 30
	return audio, disp, upd
}

func (o *OptionsScene) Update() {
	o.widgets = o.widgets[:0]
	header(&o.widgets, "", o.goBack)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		o.goBack()
		return
	}
	a, d, u := o.sections()
	o.scroll.Max = o.contentH - SH
	o.scroll.Update(Rect{0, 90, SW, SH - 90})
	sl := a.X + 230
	sw := a.W - 230 - 190
	o.widgets = append(o.widgets,
		&Slider{ID: "mvol", R: Rect{sl, a.Y + 92, sw, 20}, Value: &settings.MusicVolume, OnChange: func() { settings.MusicMuted = false }},
		&Toggle{ID: "mmute", R: Rect{a.Right() - 74, a.Y + 89, 48, 26}, Value: boolPtrInv(&settings.MusicMuted), OnChange: func() { settings.MusicMuted = !*musicOn; saveSettings() }},
		&Slider{ID: "svol", R: Rect{sl, a.Y + 152, sw, 20}, Value: &settings.SFXVolume, OnChange: func() { settings.SFXMuted = false }},
		&Toggle{ID: "smute", R: Rect{a.Right() - 74, a.Y + 149, 48, 26}, Value: boolPtrInvS(&settings.SFXMuted), OnChange: func() { settings.SFXMuted = !*sfxOn; saveSettings(); playSFX("snap") }},
		&Button{ID: "nexttrack", R: Rect{sl, a.Y + 190, 160, 40}, Label: "Next song", Kind: BtnSecondary, Size: 13.5, OnClick: func() {
			startTrack(musicIdx + 1)
		}},
	)
	themeSel := 0
	if settings.Theme == "dark" {
		themeSel = 1
	}
	o.widgets = append(o.widgets, &Segmented{ID: "theme", R: Rect{sl, d.Y + 74, 280, 46}, Options: []string{"Light", "Dark"}, Selected: themeSel, OnSelect: func(i int) {
		if i == 1 {
			settings.Theme = "dark"
		} else {
			settings.Theme = "light"
		}
		saveSettings()
	}})
	// table swatches
	for i := range tableColors {
		ti := i
		r := Rect{sl + float64(i)*52, d.Y + 140, 40, 40}
		o.widgets = append(o.widgets, &swatch{id: fmt.Sprintf("sw%d", i), r: r, idx: ti})
	}
	fs := settings.Fullscreen
	o.widgets = append(o.widgets,
		&Toggle{ID: "fs", R: Rect{sl, d.Y + 202, 48, 26}, Value: &settings.Fullscreen, OnChange: func() { ebiten.SetFullscreen(settings.Fullscreen); saveSettings() }},
		&Toggle{ID: "ghost", R: Rect{sl, d.Y + 246, 48, 26}, Value: &settings.GhostDefault, OnChange: saveSettings},
	)
	_ = fs
	// updates
	uc := updateChecker
	if uc.state == updDownloading {
		// no buttons
	} else if uc.available() {
		o.widgets = append(o.widgets, &Button{ID: "updinstall", R: Rect{u.X + 28, u.Y + 140, 240, 46}, Label: "Download & install v" + uc.latest.Version, Kind: BtnPrimary, Size: 14, OnClick: func() {
			if ps, ok := o.backScene(); ok {
				ps.saveNow(false)
			}
			uc.download()
		}})
		o.widgets = append(o.widgets, &Button{ID: "updnotes", R: Rect{u.X + 282, u.Y + 140, 160, 46}, Label: "Release notes", Kind: BtnSecondary, Size: 14, OnClick: func() {
			game.SetScene(NewReleaseNotesScene(uc.latest, func() { game.SetScene(o) }))
		}})
	} else {
		o.widgets = append(o.widgets, &Button{ID: "updcheck", R: Rect{u.X + 28, u.Y + 140, 200, 46}, Label: "Check for updates", Kind: BtnTeal, Size: 14, Disabled: uc.state == updChecking, OnClick: func() { uc.check(true) }})
	}
	o.widgets = append(o.widgets, &Toggle{ID: "autoupd", R: Rect{u.Right() - 74, u.Y + 150, 48, 26}, Value: &settings.AutoCheckUpdates, OnChange: saveSettings})
	musicOnV = !settings.MusicMuted
	sfxOnV = !settings.SFXMuted
	for _, w := range o.widgets {
		w.Update()
	}
	if activeDrag == "" && (o.musicV != settings.MusicVolume || o.sfxV != settings.SFXVolume) {
		if o.sfxV != settings.SFXVolume {
			playSFX("snap")
		}
		o.musicV, o.sfxV = settings.MusicVolume, settings.SFXVolume
		saveSettings()
	}
}

func (o *OptionsScene) backScene() (*PlayScene, bool) { return nil, false }

var musicOnV, sfxOnV bool
var musicOn = &musicOnV
var sfxOn = &sfxOnV

func boolPtrInv(_ *bool) *bool  { musicOnV = !settings.MusicMuted; return musicOn }
func boolPtrInvS(_ *bool) *bool { sfxOnV = !settings.SFXMuted; return sfxOn }

type swatch struct {
	id  string
	r   Rect
	idx int
}

func (s *swatch) Update() {
	h := in.Hover(s.r)
	animHover(s.id, h)
	if h {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	}
	if in.Click(s.r) {
		playSFX("click")
		settings.TableColor = s.idx
		saveSettings()
	}
}

func (s *swatch) Draw(dst *ebiten.Image) {
	p := P()
	tc := tableColors[s.idx]
	c := tc.Light
	if settings.Theme == "dark" {
		c = tc.Dark
	}
	hv := hoverAnim[s.id]
	r := s.r.Inset(-hv * 2)
	FillRR(dst, r, 10, c)
	StrokeRR(dst, r, 10, 1, p.Border)
	if settings.TableColor == s.idx {
		StrokeRR(dst, r.Inset(-4), 13, 2.5, p.Accent)
	}
}

func optSection(dst *ebiten.Image, r Rect, title string) {
	p := P()
	Panel(dst, r, 20)
	Text(dst, title, Bold, 18, r.X+28, r.Y+38, p.Text, AlignLeft)
}

func (o *OptionsScene) Draw(dst *ebiten.Image) {
	p := P()
	dst.Fill(p.Bg)
	drawFloatPieces(dst)
	a, d, u := o.sections()
	optSection(dst, a, "Sound")
	Text(dst, "Music volume", Medium, 15, a.X+28, a.Y+102, p.Text, AlignLeft)
	Text(dst, fmt.Sprintf("%d%%", int(settings.MusicVolume*100)), Medium, 13, a.Right()-112, a.Y+102, p.TextDim, AlignRight)
	Text(dst, "Sound effects", Medium, 15, a.X+28, a.Y+162, p.Text, AlignLeft)
	Text(dst, fmt.Sprintf("%d%%", int(settings.SFXVolume*100)), Medium, 13, a.Right()-112, a.Y+162, p.TextDim, AlignRight)
	Text(dst, "On", Regular, 12, a.Right()-50, a.Y+70, p.TextDim, AlignCenter)
	if n := musicTrackName(); n != "" {
		Text(dst, "Now playing: "+n, Regular, 13, a.X+430, a.Y+210, p.TextDim, AlignLeft)
	}
	Text(dst, "Background music", Regular, 12.5, a.X+28, a.Y+210, p.TextDim, AlignLeft)

	optSection(dst, d, "Display")
	Text(dst, "Theme", Medium, 15, d.X+28, d.Y+97, p.Text, AlignLeft)
	Text(dst, "Table color", Medium, 15, d.X+28, d.Y+160, p.Text, AlignLeft)
	Text(dst, tableColors[clampInt(settings.TableColor, 0, len(tableColors)-1)].Name, Regular, 13, d.X+230+float64(len(tableColors))*52+8, d.Y+160, p.TextDim, AlignLeft)
	Text(dst, "Fullscreen", Medium, 15, d.X+28, d.Y+215, p.Text, AlignLeft)
	Text(dst, "F11 also toggles fullscreen", Regular, 12.5, d.X+300, d.Y+215, p.TextDim, AlignLeft)
	Text(dst, "Ghost image", Medium, 15, d.X+28, d.Y+259, p.Text, AlignLeft)
	Text(dst, "Show a faint picture on the board in new puzzles", Regular, 12.5, d.X+300, d.Y+259, p.TextDim, AlignLeft)

	optSection(dst, u, "Updates")
	uc := updateChecker
	Text(dst, "You have version "+Version, Medium, 15, u.X+28, u.Y+78, p.Text, AlignLeft)
	Text(dst, uc.statusText(), Regular, 13.5, u.X+28, u.Y+106, uc.statusColor(), AlignLeft)
	Text(dst, "Check automatically", Regular, 13, u.Right()-90, u.Y+163, p.TextDim, AlignRight)
	if uc.state == updDownloading {
		bar := Rect{u.X + 28, u.Y + 152, u.W - 300, 10}
		FillRR(dst, bar, 5, p.Track)
		FillRR(dst, Rect{bar.X, bar.Y, math.Max(10, bar.W*uc.progress), 10}, 5, p.Accent)
	}

	FillRect(dst, Rect{0, 0, SW, 96}, alpha(p.Bg, 0.96))
	drawHeader(dst, "Options", "Sound, themes and updates")
	for _, w := range o.widgets {
		w.Draw(dst)
	}
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ================= Changelog / release notes =================

type mdLine struct {
	kind int // 0 text, 1 h1, 2 h2, 3 h3, 4 bullet, 5 blank
	text string
}

func parseMD(md string) []mdLine {
	var out []mdLine
	for _, l := range strings.Split(strings.ReplaceAll(md, "\r", ""), "\n") {
		t := strings.TrimSpace(l)
		switch {
		case t == "":
			out = append(out, mdLine{5, ""})
		case strings.HasPrefix(t, "### "):
			out = append(out, mdLine{3, strings.TrimPrefix(t, "### ")})
		case strings.HasPrefix(t, "## "):
			out = append(out, mdLine{2, strings.TrimPrefix(t, "## ")})
		case strings.HasPrefix(t, "# "):
			out = append(out, mdLine{1, strings.TrimPrefix(t, "# ")})
		case strings.HasPrefix(t, "- ") || strings.HasPrefix(t, "* "):
			out = append(out, mdLine{4, strings.TrimSpace(t[2:])})
		default:
			out = append(out, mdLine{0, t})
		}
	}
	return out
}

func stripMD(s string) string {
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "`", "")
	return s
}

type ChangelogScene struct {
	widgets []Widget
	lines   []mdLine
	scroll  ScrollState
	title   string
	sub     string
	back    func()
}

func NewChangelogScene() *ChangelogScene {
	return &ChangelogScene{lines: parseMD(changelogMD), title: "What's new", sub: "Version history and release notes", back: func() { game.SetScene(NewMenuScene()) }}
}

func NewReleaseNotesScene(r *ReleaseInfo, back func()) *ChangelogScene {
	body := r.Notes
	if strings.TrimSpace(body) == "" {
		body = "No release notes were provided for this version."
	}
	return &ChangelogScene{lines: parseMD(body), title: "Version " + r.Version, sub: "Release notes from GitHub", back: back}
}

func (c *ChangelogScene) area() Rect {
	w := math.Min(820, SW-80)
	return Rect{SW/2 - w/2, 110, w, SH - 130}
}

func (c *ChangelogScene) Update() {
	c.widgets = c.widgets[:0]
	header(&c.widgets, "", c.back)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		c.back()
		return
	}
	ar := c.area()
	c.scroll.Max = c.measure(ar) - ar.H + 60
	c.scroll.Update(ar)
	for _, w := range c.widgets {
		w.Update()
	}
}

func (c *ChangelogScene) render(dst *ebiten.Image, ar Rect, draw bool) float64 {
	p := P()
	x := ar.X + 40
	w := ar.W - 80
	y := ar.Y + 30 - c.scroll.Y
	start := y
	for _, l := range c.lines {
		t := stripMD(l.text)
		switch l.kind {
		case 1:
			// document title is shown in the header already
		case 2:
			y += 10
			if draw {
				FillRR(dst, Rect{x - 14, y, 5, 26}, 2.5, p.Accent)
				Text(dst, t, Bold, 21, x, y+13, p.Text, AlignLeft)
			}
			y += 40
		case 3:
			if draw {
				Text(dst, t, Bold, 15, x, y+11, p.Teal, AlignLeft)
			}
			y += 28
		case 4:
			lines := WrapText(t, Regular, 14.5, w-26)
			for i, ln := range lines {
				if draw {
					if i == 0 {
						FillCircle(dst, x+6, y+12, 3, p.Accent)
					}
					Text(dst, ln, Regular, 14.5, x+22, y+12, p.Text, AlignLeft)
				}
				y += 24
			}
			y += 3
		case 5:
			y += 8
		default:
			lines := WrapText(t, Regular, 14.5, w)
			for _, ln := range lines {
				if draw {
					Text(dst, ln, Regular, 14.5, x, y+12, p.TextDim, AlignLeft)
				}
				y += 24
			}
		}
	}
	return y - start
}

func (c *ChangelogScene) measure(ar Rect) float64 { return c.render(nil, ar, false) }

func (c *ChangelogScene) Draw(dst *ebiten.Image) {
	p := P()
	dst.Fill(p.Bg)
	drawFloatPieces(dst)
	ar := c.area()
	Panel(dst, Rect{ar.X, ar.Y, ar.W, math.Max(ar.H, 100)}, 22)
	sub := dst.SubImage(rectToImage(Rect{ar.X, ar.Y + 2, ar.W, ar.H - 4})).(*ebiten.Image)
	c.render(sub, ar, true)
	c.scroll.DrawBar(dst, ar.Inset(6))
	FillRect(dst, Rect{0, 0, SW, 96}, alpha(p.Bg, 0.96))
	drawHeader(dst, c.title, c.sub)
	for _, w := range c.widgets {
		w.Draw(dst)
	}
}
