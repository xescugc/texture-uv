package uv_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xescugc/texture-uv/uv"
)

func TestNewMap_BasicTemplate(t *testing.T) {
	tmpl := image.NewNRGBA(image.Rect(0, 0, 3, 3))
	tmpl.Set(0, 0, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	tmpl.Set(2, 1, color.NRGBA{R: 100, G: 100, B: 100, A: 255})
	tmpl.Set(1, 2, color.NRGBA{R: 50, G: 50, B: 50, A: 255})

	out, err := uv.NewMap(tmpl)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	assert.Equal(t, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, nrgba.NRGBAAt(0, 0))
	assert.Equal(t, color.NRGBA{R: 2, G: 1, B: 0, A: 255}, nrgba.NRGBAAt(2, 1))
	assert.Equal(t, color.NRGBA{R: 1, G: 2, B: 0, A: 255}, nrgba.NRGBAAt(1, 2))
	// Transparent pixels stay transparent
	assert.Equal(t, color.NRGBA{}, nrgba.NRGBAAt(1, 0))
	assert.Equal(t, color.NRGBA{}, nrgba.NRGBAAt(0, 1))
}

func TestNewMap_AllOpaque(t *testing.T) {
	tmpl := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			tmpl.Set(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}

	out, err := uv.NewMap(tmpl)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	// Every pixel should have a unique color
	seen := make(map[color.NRGBA]bool)
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			c := nrgba.NRGBAAt(x, y)
			assert.Equal(t, color.NRGBA{R: uint8(x), G: uint8(y), B: 0, A: 255}, c)
			assert.False(t, seen[c], "duplicate color at (%d,%d)", x, y)
			seen[c] = true
		}
	}
}

func TestNewMap_AllTransparent(t *testing.T) {
	tmpl := image.NewNRGBA(image.Rect(0, 0, 3, 3))

	out, err := uv.NewMap(tmpl)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			assert.Equal(t, color.NRGBA{}, nrgba.NRGBAAt(x, y))
		}
	}
}

func TestNewMap_TemplateTooLarge(t *testing.T) {
	_, err := uv.NewMap(image.NewNRGBA(image.Rect(0, 0, 256, 256)))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceed the maximum of 255x255")

	_, err = uv.NewMap(image.NewNRGBA(image.Rect(0, 0, 256, 1)))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceed the maximum of 255x255")

	_, err = uv.NewMap(image.NewNRGBA(image.Rect(0, 0, 1, 256)))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceed the maximum of 255x255")
}

func TestNewMap_MaxSize(t *testing.T) {
	tmpl := image.NewNRGBA(image.Rect(0, 0, 255, 255))
	tmpl.Set(0, 0, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	tmpl.Set(254, 254, color.NRGBA{R: 255, G: 255, B: 255, A: 255})

	out, err := uv.NewMap(tmpl)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	assert.Equal(t, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, nrgba.NRGBAAt(0, 0))
	assert.Equal(t, color.NRGBA{R: 254, G: 254, B: 0, A: 255}, nrgba.NRGBAAt(254, 254))
}

func TestNewMap_NonZeroOrigin(t *testing.T) {
	tmpl := image.NewNRGBA(image.Rect(5, 5, 10, 10))
	tmpl.Set(5, 5, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	tmpl.Set(9, 9, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	tmpl.Set(7, 6, color.NRGBA{R: 255, G: 255, B: 255, A: 255})

	out, err := uv.NewMap(tmpl)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	// Coordinates should be relative to Min
	assert.Equal(t, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, nrgba.NRGBAAt(5, 5))
	assert.Equal(t, color.NRGBA{R: 4, G: 4, B: 0, A: 255}, nrgba.NRGBAAt(9, 9))
	assert.Equal(t, color.NRGBA{R: 2, G: 1, B: 0, A: 255}, nrgba.NRGBAAt(7, 6))
}

func TestNewMap_SemiTransparent(t *testing.T) {
	tmpl := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	tmpl.Set(0, 0, color.NRGBA{R: 255, G: 255, B: 255, A: 128}) // semi-transparent
	tmpl.Set(1, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 1})   // barely visible

	out, err := uv.NewMap(tmpl)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	// Semi-transparent pixels should be treated as opaque
	assert.Equal(t, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, nrgba.NRGBAAt(0, 0))
	assert.Equal(t, color.NRGBA{R: 1, G: 1, B: 0, A: 255}, nrgba.NRGBAAt(1, 1))
	// Fully transparent stays transparent
	assert.Equal(t, color.NRGBA{}, nrgba.NRGBAAt(1, 0))
	assert.Equal(t, color.NRGBA{}, nrgba.NRGBAAt(0, 1))
}
