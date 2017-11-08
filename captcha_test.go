package captcha

import (
	"bytes"
	"errors"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"image/gif"
	"image/jpeg"
	"math/rand"
	"os"
	"testing"
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
	err = data.WriteGIF(buf, &gif.Options{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCaptchaOptions(t *testing.T) {
	New(100, 34, func(options *Options) {
		options.BackgroundColor = color.Opaque
		options.CharPreset = "1234567890"
		options.CurveNumber = 0
		options.TextLength = 6
	})

	NewMathExpr(100, 34, func(options *Options) {
		options.BackgroundColor = color.Black
	})
}

func TestNewMathExpr(t *testing.T) {
	_, err := NewMathExpr(150, 50)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCovNilFontError(t *testing.T) {
	temp := ttfFont
	ttfFont = nil

	_, err := New(150, 50)
	if err == nil {
		t.Fatal("Expect to get nil font error")
	}

	_, err = NewMathExpr(150, 50)
	if err == nil {
		t.Fatal("Expect to get nil font error")
	}

	ttfFont = temp
}

type errReader struct{}

func (errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("")
}

func TestCovReaderErr(t *testing.T) {
	err := LoadFontFromReader(errReader{})
	if err == nil {
		t.Fatal("Expect to get io.Reader error")
	}
}

func TestLoadFont(t *testing.T) {
	err := LoadFont(goregular.TTF)
	if err != nil {
		t.Fatal("Fail to load go font")
	}

	err = LoadFont([]byte("invalid"))
	if err == nil {
		t.Fatal("LoadFont incorrectly parse an invalid font")
	}
}

func TestLoadFontFromReader(t *testing.T) {
	file, err := os.Open("./fonts/Comismsh.ttf")
	if err != nil {
		t.Fatal("Fail to load test file")
	}

	if err = LoadFontFromReader(file); err != nil {
		t.Fatal("Fail to load font from io.Reader")
	}
}

func TestMaxColor(t *testing.T) {
	var result uint32
	result = maxColor()
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
	var result uint32
	result = minColor()
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
