package captcha

import (
	"fmt"
	"image/color"
	"testing"
)

func TestHSVAInterface(_ *testing.T) {
	var _ color.Color = hsva{}
}

func TestConversionGray(t *testing.T) {
	var white1 color.Color = hsva{h: 0.0, s: 0.0, v: 1.0, a: uint8(255)}
	var white2 color.Color = color.RGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(255)}

	if err := eq(white1, white2); err != nil {
		t.Fatal(err)
	}

	var black1 color.Color = hsva{h: 0.0, s: 0.0, v: 0.0, a: uint8(255)}
	var black2 color.Color = color.RGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}

	if err := eq(black1, black2); err != nil {
		t.Fatal(err)
	}

	var gray1 color.Color = hsva{h: 1.0, s: 0.0, v: 0.8, a: uint8(255)}
	var gray2 color.Color = color.RGBA{R: uint8(204), G: uint8(204), B: uint8(204), A: uint8(255)}

	if err := eq(gray1, gray2); err != nil {
		t.Fatal(err)
	}
}

func TestConversionRGB(t *testing.T) {
	var yellow1 color.Color = hsva{h: 60.0 / 360.0, s: 0.6, v: 1.0, a: uint8(255)}
	var yellow2 color.Color = color.RGBA{R: uint8(255), G: uint8(255), B: uint8(102), A: uint8(255)}

	if err := eq(yellow1, yellow2); err != nil {
		t.Fatal(err)
	}

	var green1 color.Color = hsva{h: 120.0 / 360.0, s: 1.0, v: 0.8, a: uint8(255)}
	var green2 color.Color = color.RGBA{R: uint8(0), G: uint8(204), B: uint8(0), A: uint8(255)}

	if err := eq(green1, green2); err != nil {
		t.Fatal(err)
	}

	var teal1 color.Color = hsva{h: 180.0 / 360.0, s: 0.5, v: 0.8, a: uint8(255)}
	var teal2 color.Color = color.RGBA{R: uint8(102), G: uint8(204), B: uint8(204), A: uint8(255)}

	if err := eq(teal1, teal2); err != nil {
		t.Fatal(err)
	}

	var blue1 color.Color = hsva{h: 240.0 / 360.0, s: 0.75, v: 0.8, a: uint8(255)}
	var blue2 color.Color = color.RGBA{R: uint8(51), G: uint8(51), B: uint8(204), A: uint8(255)}

	if err := eq(blue1, blue2); err != nil {
		t.Fatal(err)
	}

	var pink1 color.Color = hsva{h: 300.0 / 360.0, s: 0.2, v: 1.0, a: uint8(255)}
	var pink2 color.Color = color.RGBA{R: uint8(255), G: uint8(204), B: uint8(255), A: uint8(255)}

	if err := eq(pink1, pink2); err != nil {
		t.Fatal(err)
	}
}

func eq(c0, c1 color.Color) error {
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()
	if r0 != r1 || g0 != g1 || b0 != b1 || a0 != a1 {
		return fmt.Errorf("got  0x%04x 0x%04x 0x%04x 0x%04x\nwant 0x%04x 0x%04x 0x%04x 0x%04x",
			r0, g0, b0, a0, r1, g1, b1, a1)
	}
	return nil
}
