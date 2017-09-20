package captcha

import (
	"bytes"
	"image/color"
	"testing"
)

func TestNewCaptcha(t *testing.T) {
	New(36, 12)
	data, err := New(150, 50)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	data.WriteTo(buf)
}

func TestNewCaptchaOptions(t *testing.T) {
	New(100, 34, func(options *Options) {
		options.BackgroundColor = color.Opaque
		options.CharPreset = "1234567890"
		options.CurveNumber = 0
		options.TextLength = 6
	})
}

func TestCovNilFontError(t *testing.T) {
	temp := ttfFont
	ttfFont = nil

	_, err := New(150, 50)
	if err == nil {
		t.Fatal("Expect to get nil font error")
	}

	ttfFont = temp
}
