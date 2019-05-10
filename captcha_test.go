package captcha

import (
	"bytes"
	"errors"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"math/rand"
	"os"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

func TestNewCaptcha(t *testing.T) {
	data, err := New(150, 50)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = data.WriteImage(buf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSmallCaptcha(t *testing.T) {
	_, err := New(36, 12)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEncodeJPG(t *testing.T) {
	data, err := New(150, 50)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = data.WriteJPG(buf, &jpeg.Options{Quality: 70})
	if err != nil {
		t.Fatal(err)
	}
}

func TestEncodeGIF(t *testing.T) {
	data, err := New(150, 50)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = data.WriteGIF(buf, new(gif.Options))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCaptchaOptions(t *testing.T) {
	_, err := New(100, 34, func(options *Options) {
		options.BackgroundColor = color.Opaque
		options.CharPreset = "1234567890"
		options.CurveNumber = 0
		options.TextLength = 6
		options.Palette = palette.WebSafe
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewMathExpr(100, 34, func(options *Options) {
		options.BackgroundColor = color.Black
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewMathExpr(t *testing.T) {
	_, err := NewMathExpr(150, 50)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCovInternalFontErr(t *testing.T) {
	if ttfFont := MustLoadFont(ttf); ttfFont == nil {
		t.Fatal("Fail to load internal font")
	}
}

type errReader struct{}

func (errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("")
}

func TestCovReaderErr(t *testing.T) {
	_, err := LoadFontFromReader(errReader{})
	if err == nil {
		t.Fatal("Expect to get io.Reader error")
	}
}

func TestLoadFont(t *testing.T) {
	if _, err := LoadFont(goregular.TTF); err != nil {
		t.Fatal("Fail to load go font")
	}

	if _, err := LoadFont([]byte("invalid")); err == nil {
		t.Fatal("LoadFont incorrectly parse an invalid font")
	}
}

func TestLoadFontFromReader(t *testing.T) {
	file, err := os.Open("./fonts/Comismsh.ttf")
	if err != nil {
		t.Fatal("Fail to load test file")
	}

	if _, err := LoadFontFromReader(file); err != nil {
		t.Fatal("Fail to load font from io.Reader")
	}
}

func TestMaxColor(t *testing.T) {
	result := maxColor()
	if result != 0 {
		t.Fatalf("Expect max color to be 0, got %v", result)
	}

	result = maxColor(1)
	if result != 1 {
		t.Fatalf("Expect max color to be 1, got %v", result)
	}

	result = maxColor(52428, 65535)
	if result != 255 {
		t.Fatalf("Expect max color to be 255, got %v", result)
	}

	var rng = rand.New(rand.NewSource(0))
	for i := 0; i < 10; i++ {
		result = maxColor(rng.Uint32(), rng.Uint32(), rng.Uint32())
		if result > 255 {
			t.Fatalf("Number out of range: %v", result)
		}
	}
}

func TestMinColor(t *testing.T) {
	result := minColor()
	if result != 255 {
		t.Fatalf("Expect min color to be 255, got %v", result)
	}

	result = minColor(1)
	if result != 1 {
		t.Fatalf("Expect min color to be 1, got %v", result)
	}

	result = minColor(52428, 65535)
	if result != 204 {
		t.Fatalf("Expect min color to be 1, got %v", result)
	}

	var rng = rand.New(rand.NewSource(0))
	for i := 0; i < 10; i++ {
		result = minColor(rng.Uint32(), rng.Uint32(), rng.Uint32())
		if result > 255 {
			t.Fatalf("Number out of range: %v", result)
		}
	}
}
