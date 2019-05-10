package captcha

import "math"

type hsva struct {
	h, s, v float64
	a       uint8
}

// https://gist.github.com/mjackson/5311256
func (c hsva) RGBA() (r, g, b, a uint32) {
	var i = math.Floor(c.h * 6)
	var f = c.h*6 - i
	var p = c.v * (1.0 - c.s)
	var q = c.v * (1.0 - f*c.s)
	var t = c.v * (1 - (1-f)*c.s)

	var red, green, blue float64
	switch int(i) % 6 {
	case 0:
		red, green, blue = c.v, t, p
	case 1:
		red, green, blue = q, c.v, p
	case 2:
		red, green, blue = p, c.v, t
	case 3:
		red, green, blue = p, q, c.v
	case 4:
		red, green, blue = t, p, c.v
	case 5:
		red, green, blue = c.v, p, q
	}

	r = uint32(red * 255)
	r |= r << 8
	g = uint32(green * 255)
	g |= g << 8
	b = uint32(blue * 255)
	b |= b << 8
	a = uint32(c.a)
	a |= a << 8

	return r, g, b, a
}
