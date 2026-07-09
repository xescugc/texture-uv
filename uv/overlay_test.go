package uv_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xescugc/texture-uv/uv"
)

func TestNewOverlay_SingleGroupSingleFrame(t *testing.T) {
	// 4x4 sprite with a 2x2 block of color at (1,1)
	sprite := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	c := color.NRGBA{R: 46, G: 34, B: 47, A: 255}
	sprite.Set(1, 1, c)
	sprite.Set(2, 1, c)
	sprite.Set(1, 2, c)
	sprite.Set(2, 2, c)

	// Map image: 8x8 with distinct colors in the UV region (0,0)-(2,2)
	mapImg := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	mapImg.Set(0, 0, color.NRGBA{R: 10, G: 20, B: 30, A: 255})
	mapImg.Set(1, 0, color.NRGBA{R: 40, G: 50, B: 60, A: 255})
	mapImg.Set(0, 1, color.NRGBA{R: 70, G: 80, B: 90, A: 255})
	mapImg.Set(1, 1, color.NRGBA{R: 100, G: 110, B: 120, A: 255})

	mapping := uv.ColorMapping{
		FrameSize: [2]int{4, 4},
		Columns: map[string]uv.ColumnConfig{
			"0": {
				Groups: []uv.Group{
					{UV: [4]int{0, 0, 2, 2}, Colors: [][3]int{{46, 34, 47}}},
				},
			},
		},
	}

	out, err := uv.NewOverlay(sprite, mapImg, mapping)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	assert.Equal(t, color.NRGBA{R: 10, G: 20, B: 30, A: 255}, nrgba.NRGBAAt(1, 1))
	assert.Equal(t, color.NRGBA{R: 40, G: 50, B: 60, A: 255}, nrgba.NRGBAAt(2, 1))
	assert.Equal(t, color.NRGBA{R: 70, G: 80, B: 90, A: 255}, nrgba.NRGBAAt(1, 2))
	assert.Equal(t, color.NRGBA{R: 100, G: 110, B: 120, A: 255}, nrgba.NRGBAAt(2, 2))
	// Unmapped pixels stay transparent
	assert.Equal(t, color.NRGBA{}, nrgba.NRGBAAt(0, 0))
}

func TestNewOverlay_MultipleGroups(t *testing.T) {
	sprite := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	c1 := color.NRGBA{R: 100, G: 0, B: 0, A: 255}
	c2 := color.NRGBA{R: 0, G: 100, B: 0, A: 255}
	sprite.Set(0, 0, c1)
	sprite.Set(2, 2, c2)

	mapImg := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	mapImg.Set(0, 0, color.NRGBA{R: 11, G: 22, B: 33, A: 255})
	mapImg.Set(4, 4, color.NRGBA{R: 44, G: 55, B: 66, A: 255})

	mapping := uv.ColorMapping{
		FrameSize: [2]int{4, 4},
		Columns: map[string]uv.ColumnConfig{
			"0": {
				Groups: []uv.Group{
					{UV: [4]int{0, 0, 1, 1}, Colors: [][3]int{{100, 0, 0}}},
					{UV: [4]int{4, 4, 1, 1}, Colors: [][3]int{{0, 100, 0}}},
				},
			},
		},
	}

	out, err := uv.NewOverlay(sprite, mapImg, mapping)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	assert.Equal(t, color.NRGBA{R: 11, G: 22, B: 33, A: 255}, nrgba.NRGBAAt(0, 0))
	assert.Equal(t, color.NRGBA{R: 44, G: 55, B: 66, A: 255}, nrgba.NRGBAAt(2, 2))
}

func TestNewOverlay_MultipleColumns(t *testing.T) {
	// 8x4 sprite: 2 columns of 4x4 frames
	sprite := image.NewNRGBA(image.Rect(0, 0, 8, 4))
	c := color.NRGBA{R: 50, G: 50, B: 50, A: 255}
	sprite.Set(1, 1, c) // column 0
	sprite.Set(5, 1, c) // column 1

	mapImg := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	mapImg.Set(0, 0, color.NRGBA{R: 10, G: 10, B: 10, A: 255})
	mapImg.Set(2, 0, color.NRGBA{R: 20, G: 20, B: 20, A: 255})

	mapping := uv.ColorMapping{
		FrameSize: [2]int{4, 4},
		Columns: map[string]uv.ColumnConfig{
			"0": {
				Groups: []uv.Group{
					{UV: [4]int{0, 0, 1, 1}, Colors: [][3]int{{50, 50, 50}}},
				},
			},
			"1": {
				Groups: []uv.Group{
					{UV: [4]int{2, 0, 1, 1}, Colors: [][3]int{{50, 50, 50}}},
				},
			},
		},
	}

	out, err := uv.NewOverlay(sprite, mapImg, mapping)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	assert.Equal(t, color.NRGBA{R: 10, G: 10, B: 10, A: 255}, nrgba.NRGBAAt(1, 1))
	assert.Equal(t, color.NRGBA{R: 20, G: 20, B: 20, A: 255}, nrgba.NRGBAAt(5, 1))
}

func TestNewOverlay_MultipleFrames(t *testing.T) {
	// 4x8 sprite: 1 column, 2 rows
	sprite := image.NewNRGBA(image.Rect(0, 0, 4, 8))
	c := color.NRGBA{R: 60, G: 60, B: 60, A: 255}
	sprite.Set(0, 0, c) // frame row 0
	sprite.Set(1, 5, c) // frame row 1, at (1, 1) within frame

	mapImg := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	mapImg.Set(0, 0, color.NRGBA{R: 11, G: 11, B: 11, A: 255})

	mapping := uv.ColorMapping{
		FrameSize: [2]int{4, 4},
		Columns: map[string]uv.ColumnConfig{
			"0": {
				Groups: []uv.Group{
					{UV: [4]int{0, 0, 2, 2}, Colors: [][3]int{{60, 60, 60}}},
				},
			},
		},
	}

	out, err := uv.NewOverlay(sprite, mapImg, mapping)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	// Frame 0: single pixel at (0,0), bb is (0,0)-(0,0), maps to map(0,0)
	assert.Equal(t, color.NRGBA{R: 11, G: 11, B: 11, A: 255}, nrgba.NRGBAAt(0, 0))
	// Frame 1: single pixel at (1,5), bb is (1,5)-(1,5), maps to map(0,0)
	assert.Equal(t, color.NRGBA{R: 11, G: 11, B: 11, A: 255}, nrgba.NRGBAAt(1, 5))
}

func TestNewOverlay_TransparentPixelsSkipped(t *testing.T) {
	sprite := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	// Same RGB as the group color but alpha=0
	sprite.Set(0, 0, color.NRGBA{R: 50, G: 50, B: 50, A: 0})
	// Opaque pixel with group color
	sprite.Set(1, 1, color.NRGBA{R: 50, G: 50, B: 50, A: 255})

	mapImg := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	mapImg.Set(0, 0, color.NRGBA{R: 99, G: 99, B: 99, A: 255})

	mapping := uv.ColorMapping{
		FrameSize: [2]int{4, 4},
		Columns: map[string]uv.ColumnConfig{
			"0": {
				Groups: []uv.Group{
					{UV: [4]int{0, 0, 1, 1}, Colors: [][3]int{{50, 50, 50}}},
				},
			},
		},
	}

	out, err := uv.NewOverlay(sprite, mapImg, mapping)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	// Transparent pixel should NOT be in the bounding box
	assert.Equal(t, color.NRGBA{}, nrgba.NRGBAAt(0, 0))
	assert.Equal(t, color.NRGBA{R: 99, G: 99, B: 99, A: 255}, nrgba.NRGBAAt(1, 1))
}

func TestNewOverlay_UnmappedColorsStayTransparent(t *testing.T) {
	sprite := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	sprite.Set(0, 0, color.NRGBA{R: 50, G: 50, B: 50, A: 255})   // mapped
	sprite.Set(1, 0, color.NRGBA{R: 200, G: 200, B: 200, A: 255}) // NOT mapped

	mapImg := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	mapImg.Set(0, 0, color.NRGBA{R: 99, G: 99, B: 99, A: 255})

	mapping := uv.ColorMapping{
		FrameSize: [2]int{4, 4},
		Columns: map[string]uv.ColumnConfig{
			"0": {
				Groups: []uv.Group{
					{UV: [4]int{0, 0, 1, 1}, Colors: [][3]int{{50, 50, 50}}},
				},
			},
		},
	}

	out, err := uv.NewOverlay(sprite, mapImg, mapping)
	require.NoError(t, err)

	nrgba := out.(*image.NRGBA)
	assert.Equal(t, color.NRGBA{R: 99, G: 99, B: 99, A: 255}, nrgba.NRGBAAt(0, 0))
	assert.Equal(t, color.NRGBA{}, nrgba.NRGBAAt(1, 0)) // unmapped → transparent
}

func TestNewOverlay_BoundingBoxExceedsUV(t *testing.T) {
	sprite := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	c := color.NRGBA{R: 50, G: 50, B: 50, A: 255}
	// 3 pixels wide bounding box but UV is only 2 wide
	sprite.Set(0, 0, c)
	sprite.Set(2, 0, c)

	mapImg := image.NewNRGBA(image.Rect(0, 0, 8, 8))

	mapping := uv.ColorMapping{
		FrameSize: [2]int{4, 4},
		Columns: map[string]uv.ColumnConfig{
			"0": {
				Groups: []uv.Group{
					{UV: [4]int{0, 0, 2, 1}, Colors: [][3]int{{50, 50, 50}}},
				},
			},
		},
	}

	_, err := uv.NewOverlay(sprite, mapImg, mapping)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bounding box exceeds UV region")
}

func TestNewOverlay_NoMatchingPixels(t *testing.T) {
	sprite := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	// No pixels match the group color

	mapImg := image.NewNRGBA(image.Rect(0, 0, 8, 8))

	mapping := uv.ColorMapping{
		FrameSize: [2]int{4, 4},
		Columns: map[string]uv.ColumnConfig{
			"0": {
				Groups: []uv.Group{
					{UV: [4]int{0, 0, 2, 2}, Colors: [][3]int{{50, 50, 50}}},
				},
			},
		},
	}

	out, err := uv.NewOverlay(sprite, mapImg, mapping)
	require.NoError(t, err)

	// All pixels should be transparent
	expected := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	assert.Equal(t, expected, out)
}
