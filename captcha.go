package captcha

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
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

type Options struct {
	BackgroundColor color.Color
	CharPreset      string
	TxtLength       int
	width           int
	height          int
}

func newDefaultOption(width, height int) *Options {
	return &Options{
		BackgroundColor: color.Transparent,
		CharPreset:      charPreset,
		TxtLength:       4,
		width:           width,
		height:          height,
	}
}

type SetOption func(*Options)

type Data struct {
	Text string

	img *image.NRGBA
}

func (data *Data) WriteTo(w io.Writer) error {
	return png.Encode(w, data.img)
}

func init() {
	var err error
	ttfFont, err = freetype.ParseFont(goregular.TTF)
	if err != nil {
		panic(err)
	}
}

func New(width int, height int, option ...SetOption) (*Data, error) {
	options := newDefaultOption(width, height)
	for _, setOption := range option {
		setOption(options)
	}

	text := randomText(options)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{options.BackgroundColor}, image.ZP, draw.Src)
	drawNoise(img, options)
	drawLine(img, options)
	err := drawText(text, img, options)
	if err != nil {
		return nil, err
	}

	return &Data{Text: text, img: img}, nil
}

func randomText(opts *Options) (text string) {
	n := len(opts.CharPreset)
	for i := 0; i < opts.TxtLength; i++ {
		text += string(opts.CharPreset[rng.Intn(n)])
	}

	return text
}

func drawNoise(img *image.NRGBA, opts *Options) {
	noiseCount := (opts.width * opts.height) / 18
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

func drawLine(img *image.NRGBA, opts *Options) {
	for i := 0; i < 3; i++ {
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
	ctx.SetDPI(72.0)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetHinting(font.HintingFull)
	ctx.SetFont(ttfFont)

	fontSpacing := opts.width / len(text)

	for idx, char := range text {
		fontScale := 1 + float64(rng.Intn(7))/float64(9)
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
