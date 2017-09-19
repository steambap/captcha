package captcha

import (
	"bytes"
	"testing"
)

func TestNewCaptcha(t *testing.T) {
	data, err := New(150, 50)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	data.WriteTo(buf)
}
