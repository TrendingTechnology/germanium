package main

import (
	"image"
	"image/color"
	"image/draw"
)

var (
	pw = 60
	ph = 60
	wh = 20 * 3
	lw = 40

	radius = 10

	background = color.RGBA{170, 170, 255, 255}
	dracula    = color.RGBA{40, 42, 54, 255}

	close   = color.RGBA{255, 95, 86, 255}
	minimum = color.RGBA{255, 189, 46, 255}
	maximum = color.RGBA{39, 201, 63, 255}
)

func NewPanels(w, h int) (*Rect, *Rect, *Rect) {
	// base panel
	base := NewRect(0, 0, w, h, background)
	base.fillColor()

	// window panel
	window := NewRect(pw, ph, w-pw, h-ph, dracula)
	window.fillColor()
	base.drawOver(window.img)

	// window control bar
	if !opts.NoWindowAccessBar {
		wc := NewRect(pw, ph, w-pw, ph+wh, dracula)
		wc.fillColor()

		// control buttons
		buttons := []color.RGBA{close, minimum, maximum}
		for i, b := range buttons {
			wc.drawCircle(image.Point{X: pw + (i * 30) + 20, Y: ph + 10*2}, radius, b)
		}
		base.drawOver(wc.img)
	} else {
		wh = 0
	}

	// round corner
	round := NewRect(pw-radius, ph-radius, w-pw+radius, h-ph+radius, dracula)
	corners := []image.Point{
		image.Point{pw, ph},
		image.Point{w - pw, ph},
		image.Point{pw, h - ph},
		image.Point{w - pw, h - ph},
	}
	for _, c := range corners {
		round.drawCircle(c, radius, round.color)
	}
	base.drawOver(round.img)

	aroundbars := []*Rect{
		NewRect(pw-radius, ph, pw, h-ph, dracula),
		NewRect(pw, ph-radius, w-pw, ph, dracula),
		NewRect(w-pw, ph, w-pw+radius, h-ph, dracula),
		NewRect(pw, h-ph, w-pw, h-ph+radius, dracula),
	}
	for _, ab := range aroundbars {
		ab.fillColor()
		base.drawOver(ab.img)
	}

	// editor panel
	editor := NewRect(pw, ph+wh, w-pw, h-ph, dracula)
	editor.fillColor()
	base.drawOver(editor.img)

	// line number panel
	line := NewRect(pw, ph+wh, pw+lw, h-ph, dracula)
	line.fillColor()
	base.drawOver(line.img)

	return base, editor, line
}

type Rect struct {
	img   *image.RGBA
	color color.RGBA
}

func NewRect(sx, sy, ex, ey int, c color.RGBA) *Rect {
	rect := image.Rect(sx, sy, ex, ey)
	return &Rect{img: image.NewRGBA(rect), color: c}
}

func (r *Rect) fillColor() {
	for x := r.img.Rect.Min.X; x < r.img.Rect.Max.X; x++ {
		for y := r.img.Rect.Min.Y; y < r.img.Rect.Max.Y; y++ {
			r.img.Set(x, y, r.color)
		}
	}
}

// drawOver draw image over r.img
func (r *Rect) drawOver(img *image.RGBA) {
	draw.Draw(r.img, r.img.Bounds(), img, image.Point{0, 0}, draw.Over)
}

// drawCircle draw circle over r.img
// http://dencha.ojaru.jp/programs_07/pg_graphic_09a1.html section7
func (r *Rect) drawCircle(center image.Point, radius int, c color.RGBA) {
	var cx, cy, d, dh, dd int
	d = 1 - radius
	dh = 3
	dd = 5 - 2*radius
	cy = radius

	for cx = 0; cx <= cy; cx++ {
		if d < 0 {
			d += dh
			dh += 2
			dd += 2
		} else {
			d += dd
			dh += 2
			dd += 4
			cy--
		}

		r.img.Set(center.X+cy, center.Y+cx, c) // 0-45
		r.img.Set(center.X+cx, center.Y+cy, c) // 45-90
		r.img.Set(center.X-cx, center.Y+cy, c) // 90-135
		r.img.Set(center.X-cy, center.Y+cx, c) // 135-180
		r.img.Set(center.X-cy, center.Y-cx, c) // 180-225
		r.img.Set(center.X-cx, center.Y-cy, c) // 225-270
		r.img.Set(center.X+cx, center.Y-cy, c) // 270-315
		r.img.Set(center.X+cy, center.Y-cx, c) // 315-360

		// draw line same y position
		for x := center.X - cy; x <= center.X+cy; x++ {
			r.img.Set(x, center.Y+cx, c)
			r.img.Set(x, center.Y-cx, c)
		}
		for x := center.X - cx; x <= center.X+cx; x++ {
			r.img.Set(x, center.Y+cy, c)
			r.img.Set(x, center.Y-cy, c)
		}
	}
}
