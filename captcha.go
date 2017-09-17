package captcha

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math/rand"
	"time"
)

const charPreset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

type Options struct {
	BackgroundColor color.RGBA
	CharPreset      string
	TxtLength       int
	width           int
	height          int
}

func newDefaultOption(width, height int) *Options {
	return &Options{
		CharPreset: charPreset,
		TxtLength:  4,
		width:      width,
		height:     height,
	}
}

type Option func(*Options)

type Data struct {
	Text string

	img *image.NRGBA
}

func (data *Data) WriteTo(w io.Writer) error {
	return png.Encode(w, data.img)
}

func New(width int, height int, option ...Option) *Data {
	options := newDefaultOption(width, height)
	for _, setOption := range option {
		setOption(options)
	}

	text := randomText(options)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	drawNoise(img, options)

	return &Data{Text: text, img: img}
}

func randomText(opts *Options) (text string) {
	n := len(opts.CharPreset)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < opts.TxtLength; i++ {
		text += string(opts.CharPreset[rng.Intn(n)])
	}

	return text
}

func drawNoise(img *image.NRGBA, opts *Options) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	noiseCount := (opts.width * opts.height) / 18
	for i := 0; i < noiseCount; i++ {
		x := rng.Intn(opts.width)
		y := rng.Intn(opts.height)
		img.Set(x, y, randomColor())
	}
}

func randomColor() color.RGBA {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	red := rng.Intn(255)
	green := rng.Intn(255)
	blue := rng.Intn(255)

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}
