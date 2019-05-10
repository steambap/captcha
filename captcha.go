// Package captcha provides an easy to use, unopinionated API for captcha generation
package captcha

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const charPreset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// nolint: gochecknoglobals
var (
	rng     = rand.New(rand.NewSource(time.Now().UnixNano()))
	ttfFont *truetype.Font
)

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
	// FontDPI controls DPI (dots per inch) of font.
	// The default is 72.0.
	FontDPI float64
	// FontScale controls the scale of font.
	// The default is 1.0.
	FontScale float64
	// Noise controls the number of noise drawn.
	// A noise dot is drawn for every 28 pixel by default.
	// The default is 1.0.
	Noise float64
	// Palette is the set of colors to chose from
	Palette color.Palette

	width  int
	height int
}

func newDefaultOption(width, height int) *Options {
	return &Options{
		BackgroundColor: color.Transparent,
		CharPreset:      charPreset,
		TextLength:      4,
		CurveNumber:     2,
		FontDPI:         72.0,
		FontScale:       1.0,
		Noise:           1.0,
		Palette:         []color.Color{},
		width:           width,
		height:          height,
	}
}

// SetOption is a function that can be used to modify default options.
type SetOption func(*Options)

// Data is the result of captcha generation.
// It has a `Text` field and a private `img` field that will
// be used in `WriteImage` receiver.
type Data struct {
	// Text is captcha solution.
	Text string

	img *image.NRGBA
}

// WriteImage encodes image data and writes to an io.Writer.
// It returns possible error from PNG encoding.
func (data *Data) WriteImage(w io.Writer) error {
	return png.Encode(w, data.img)
}

// WriteJPG encodes image data in JPEG format and writes to an io.Writer.
// It returns possible error from JPEG encoding.
func (data *Data) WriteJPG(w io.Writer, o *jpeg.Options) error {
	return jpeg.Encode(w, data.img, o)
}

// WriteGIF encodes image data in GIF format and writes to an io.Writer.
// It returns possible error from GIF encoding.
func (data *Data) WriteGIF(w io.Writer, o *gif.Options) error {
	return gif.Encode(w, data.img, o)
}

// nolint: gochecknoinits
func init() {
	ttfFont, _ = freetype.ParseFont(ttf)
}

// LoadFont let you load an external font.
func LoadFont(fontData []byte) error {
	var err error
	ttfFont, err = freetype.ParseFont(fontData)
	return err
}

// LoadFontFromReader load an external font from an io.Reader interface.
func LoadFontFromReader(reader io.Reader) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return err
	}

	return LoadFont(buf.Bytes())
}

// New creates a new captcha.
// It returns captcha data and any freetype drawing error encountered.
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

// NewMathExpr creates a new captcha.
// It will generate a image with a math expression like `1 + 2`.
func NewMathExpr(width int, height int, option ...SetOption) (*Data, error) {
	options := newDefaultOption(width, height)
	for _, setOption := range option {
		setOption(options)
	}

	text, equation := randomEquation()
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{options.BackgroundColor}, image.ZP, draw.Src)
	drawNoise(img, options)
	drawCurves(img, options)
	err := drawText(equation, img, options)
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
	noiseCount := (opts.width * opts.height) / int(28.0/opts.Noise)
	for i := 0; i < noiseCount; i++ {
		x := rng.Intn(opts.width)
		y := rng.Intn(opts.height)
		img.Set(x, y, randomColor())
	}
}

func randomColor() color.RGBA {
	red := rng.Intn(256)
	green := rng.Intn(256)
	blue := rng.Intn(256)

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
	yFlip := 1.0
	if rng.Intn(2) == 0 {
		yFlip = -1.0
	}
	curveColor := randomColorFromOptions(opts)

	for x1 := xStart; x1 <= xEnd; x1++ {
		y := math.Sin(math.Pi*angle*float64(x1)/float64(opts.width)) * curveHeight * yFlip
		img.Set(x1, int(y)+yStart, curveColor)
	}
}

func drawText(text string, img *image.NRGBA, opts *Options) error { // nolint: interfacer
	ctx := freetype.NewContext()
	ctx.SetDPI(opts.FontDPI)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetHinting(font.HintingFull)
	ctx.SetFont(ttfFont)

	fontSpacing := opts.width / len(text)
	fontOffset := rng.Intn(fontSpacing / 2)

	for idx, char := range text {
		fontScale := 0.8 + rng.Float64()*0.4
		fontSize := float64(opts.height) / fontScale * opts.FontScale
		ctx.SetFontSize(fontSize)
		ctx.SetSrc(image.NewUniform(randomColorFromOptions(opts)))
		x := fontSpacing*idx + fontOffset
		y := opts.height/6 + rng.Intn(opts.height/3) + int(fontSize/2)
		pt := freetype.Pt(x, y)
		if _, err := ctx.DrawString(string(char), pt); err != nil {
			return err
		}
	}

	return nil
}

func randomColorFromOptions(opts *Options) color.Color {
	length := len(opts.Palette)
	if length == 0 {
		return randomInvertColor(opts.BackgroundColor)
	}

	return opts.Palette[rng.Intn(length)]
}

func randomInvertColor(base color.Color) color.Color {
	baseLightness := getLightness(base)
	var value float64
	if baseLightness >= 0.5 {
		value = baseLightness - 0.3 - rng.Float64()*0.2
	} else {
		value = baseLightness + 0.3 + rng.Float64()*0.2
	}
	hue := float64(rng.Intn(361)) / 360
	saturation := 0.6 + rng.Float64()*0.2

	return hsva{h: hue, s: saturation, v: value, a: 255}
}

func getLightness(colour color.Color) float64 {
	r, g, b, a := colour.RGBA()
	// transparent
	if a == 0 {
		return 1.0
	}
	max := maxColor(r, g, b)
	min := minColor(r, g, b)

	l := (float64(max) + float64(min)) / (2 * 255)

	return l
}

func maxColor(numList ...uint32) (max uint32) {
	for _, num := range numList {
		colorVal := num & 255
		if colorVal > max {
			max = colorVal
		}
	}

	return max
}

func minColor(numList ...uint32) (min uint32) {
	min = 255
	for _, num := range numList {
		colorVal := num & 255
		if colorVal < min {
			min = colorVal
		}
	}

	return min
}

func randomEquation() (text string, equation string) {
	left := 1 + rng.Intn(9)
	right := 1 + rng.Intn(9)
	text = strconv.Itoa(left + right)
	equation = strconv.Itoa(left) + "+" + strconv.Itoa(right)

	return text, equation
}
