package captcha

import (
	"io"
	"image/color"
	"image"
	"image/png"
	"math/rand"
	"time"
)

const charPreset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

type Options struct {
	BackgroundColor color.RGBA
	CharPreset string
	TxtLength int
}

func newDefaultOption() *Options {
	return &Options{
		CharPreset: charPreset,
		TxtLength: 4,
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

func New(width int, height int, option... Option) *Data {
	options := newDefaultOption()
	for _, setOption := range option {
		setOption(options)
	}

	text := randomText(options)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	return &Data{Text: text, img: img}
}

func randomText(opts *Options) (text string) {
	n := len(opts.CharPreset)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i :=0; i < opts.TxtLength; i++ {
		text += string(opts.CharPreset[r.Intn(n)])
	}

	return text
}
