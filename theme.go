package main

import "image/color"

type Palette struct {
	Bg, Surface, Surface2, Text, TextDim, Border, Shadow color.RGBA
	Accent, AccentHover, AccentText, Teal, Plum, Danger      color.RGBA
	Overlay, Track, Gold                                     color.RGBA
}

func hex(s string) color.RGBA {
	var v [3]uint8
	for i := 0; i < 3; i++ {
		v[i] = hexByte(s[1+i*2], s[2+i*2])
	}
	return color.RGBA{v[0], v[1], v[2], 255}
}

func hexByte(a, b byte) uint8 {
	h := func(c byte) uint8 {
		switch {
		case c >= '0' && c <= '9':
			return c - '0'
		case c >= 'a' && c <= 'f':
			return c - 'a' + 10
		case c >= 'A' && c <= 'F':
			return c - 'A' + 10
		}
		return 0
	}
	return h(a)<<4 | h(b)
}

func alpha(c color.RGBA, a float64) color.RGBA {
	// returns premultiplied color with the given opacity
	f := a * float64(c.A) / 255
	return color.RGBA{uint8(float64(c.R) * f), uint8(float64(c.G) * f), uint8(float64(c.B) * f), uint8(255 * f)}
}

func lerpColor(a, b color.RGBA, t float64) color.RGBA {
	l := func(x, y uint8) uint8 { return uint8(float64(x) + (float64(y)-float64(x))*t) }
	return color.RGBA{l(a.R, b.R), l(a.G, b.G), l(a.B, b.B), l(a.A, b.A)}
}

var lightPalette = Palette{
	Bg: hex("#FBF5EC"), Surface: hex("#FFFFFF"), Surface2: hex("#F3EADC"), Text: hex("#2E2A33"), TextDim: hex("#7A7280"),
	Border: hex("#E6DCCB"), Shadow: hex("#5A4630"), Accent: hex("#F07A2A"), AccentHover: hex("#E06818"), AccentText: hex("#FFFFFF"),
	Teal: hex("#2C9A92"), Plum: hex("#5B3F8C"), Danger: hex("#D64545"), Overlay: hex("#2E2A33"), Track: hex("#E9DFCF"), Gold: hex("#E8B031"),
}

var darkPalette = Palette{
	Bg: hex("#17161F"), Surface: hex("#23222E"), Surface2: hex("#2D2C3B"), Text: hex("#F4EEE6"), TextDim: hex("#A59FB0"),
	Border: hex("#3A3849"), Shadow: hex("#000000"), Accent: hex("#FF8E45"), AccentHover: hex("#FF7A2A"), AccentText: hex("#1A1216"),
	Teal: hex("#3DBDB3"), Plum: hex("#9B7FD1"), Danger: hex("#F06464"), Overlay: hex("#000000"), Track: hex("#3A3849"), Gold: hex("#F2C14E"),
}

func P() *Palette {
	if settings.Theme == "dark" {
		return &darkPalette
	}
	return &lightPalette
}

type TableColor struct {
	Name  string
	Light color.RGBA
	Dark  color.RGBA
}

var tableColors = []TableColor{
	{"Theme", hex("#EFE5D6"), hex("#1F1E29")},
	{"Felt green", hex("#2F5D4A"), hex("#22433A")},
	{"Ocean blue", hex("#2A4A6B"), hex("#1C3049")},
	{"Walnut", hex("#6B4A33"), hex("#3F2C20")},
	{"Plum", hex("#4C3A66"), hex("#30253F")},
	{"Slate", hex("#5A6270"), hex("#2E333C")},
}

func tableColor() color.RGBA {
	i := settings.TableColor
	if i < 0 || i >= len(tableColors) {
		i = 0
	}
	if settings.Theme == "dark" {
		return tableColors[i].Dark
	}
	return tableColors[i].Light
}

func isDarkTable() bool {
	c := tableColor()
	return int(c.R)+int(c.G)+int(c.B) < 3*110
}
