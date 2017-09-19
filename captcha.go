// Package captcha provides a simple API for captcha generation
package captcha

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"math/rand"
	"time"
)

const charPreset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
var ttfFont *truetype.Font

// Options manage captcha generation details.
type Options struct {
	// BackgroundColor is captcha image's background color.
	// It defaults to color.Transparent.
	BackgroundColor color.Color
	// CharPreset decides what text will be on captcha image.
	// It defaults to digit 0-9 and all English alphabet.
	// ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789
	CharPreset string
	// TextLength is the length of captcha text.
	// It defaults to 4.
	TextLength int
	// CurveNumber is the number of curves to draw on captcha image.
	// It defaults to 2.
	CurveNumber int

	width       int
	height      int
}

func newDefaultOption(width, height int) *Options {
	return &Options{
		BackgroundColor: color.Transparent,
		CharPreset:      charPreset,
		TextLength:      4,
		CurveNumber:     2,
		width:           width,
		height:          height,
	}
}

// SetOption is a function that can be used to modify default options.
type SetOption func(*Options)

// Data is the result of captcha generation.
// It has a `Text` field and a private `img` field that will
// be used in `WriteTo` receiver
type Data struct {
	// Text is captcha solution
	Text string

	img *image.NRGBA
}

// WriteTo encodes image data and writes to an io.Writer.
// It returns possible error from PNG encoding
func (data *Data) WriteTo(w io.Writer) error {
	return png.Encode(w, data.img)
}

func init() {
	var err error
	ttfFont, err = freetype.ParseFont(TTF)
	if err != nil {
		panic(err)
	}
}

// New creates a new captcha.
// It returns captcha data and any freetype drawing error encountered
func New(width int, height int, option ...SetOption) (*Data, error) {
	options := newDefaultOption(width, height)
	for _, setOption := range option {
		setOption(options)
	}

	text := randomText(options)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{options.BackgroundColor}, image.ZP, draw.Src)
	drawNoise(img, options)
	drawCurves(img, options)
	err := drawText(text, img, options)
	if err != nil {
		return nil, err
	}

	return &Data{Text: text, img: img}, nil
}

func randomText(opts *Options) (text string) {
	n := len(opts.CharPreset)
	for i := 0; i < opts.TextLength; i++ {
		text += string(opts.CharPreset[rng.Intn(n)])
	}

	return text
}

func drawNoise(img *image.NRGBA, opts *Options) {
	noiseCount := (opts.width * opts.height) / 28
	for i := 0; i < noiseCount; i++ {
		x := rng.Intn(opts.width)
		y := rng.Intn(opts.height)
		img.Set(x, y, randomColor())
	}
}

func randomColor() color.RGBA {
	red := rng.Intn(255)
	green := rng.Intn(255)
	blue := rng.Intn(255)

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

func drawCurves(img *image.NRGBA, opts *Options) {
	for i := 0; i < opts.CurveNumber; i++ {
		drawSineCurve(img, opts)
	}
}

// Ideally we want to draw bezier curves
// For now sine curves will do the job
func drawSineCurve(img *image.NRGBA, opts *Options) {
	var xStart, xEnd int
	if opts.width <= 40 {
		xStart, xEnd = 1, opts.width-1
	} else {
		xStart = rng.Intn(opts.width/10) + 1
		xEnd = opts.width - rng.Intn(opts.width/10) - 1
	}
	curveHeight := float64(rng.Intn(opts.height/6) + opts.height/6)
	yStart := rng.Intn(opts.height*2/3) + opts.height/6
	angle := 1.0 + rng.Float64()
	flip := rng.Intn(2) == 0
	yFlip := 1.0
	if flip {
		yFlip = -1.0
	}
	curveColor := randomDarkGray()

	for x1 := xStart; x1 <= xEnd; x1++ {
		y := math.Sin(math.Pi*angle*float64(x1)/float64(opts.width)) * curveHeight * yFlip
		img.Set(x1, int(y)+yStart, curveColor)
	}
}

func randomDarkGray() color.Gray {
	gray := rng.Intn(128) + 20

	return color.Gray{Y: uint8(gray)}
}

func drawText(text string, img *image.NRGBA, opts *Options) error {
	ctx := freetype.NewContext()
	ctx.SetDPI(92.0)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetHinting(font.HintingFull)
	ctx.SetFont(ttfFont)

	fontSpacing := opts.width / len(text)

	for idx, char := range text {
		fontScale := 1 + rng.Float64() * 0.5
		fontSize := float64(opts.height) / fontScale
		ctx.SetFontSize(fontSize)
		ctx.SetSrc(image.NewUniform(randomDarkGray()))
		x := fontSpacing*idx + fontSpacing/int(fontSize)
		y := opts.height/6 + rng.Intn(opts.height/3) + int(fontSize/2)
		pt := freetype.Pt(x, y)
		if _, err := ctx.DrawString(string(char), pt); err != nil {
			return err
		}
	}

	return nil
}
