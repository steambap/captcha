package captcha

import (
	"image/color"
	"testing"
)

func TestHSVAInterface(t *testing.T) {
	var _ color.Color = hsva{}
}
