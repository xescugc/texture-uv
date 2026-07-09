package uv_test

import (
	"image"
	"image/color"
	_ "image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xescugc/texture-uv/uv"
)

func TestNewSource(t *testing.T) {
	fo, err := os.Open("../testdata/overlay.character_walk.png")
	require.NoError(t, err)
	defer fo.Close()

	fm, err := os.Open("../testdata/map.character.png")
	require.NoError(t, err)
	defer fm.Close()

	oimg, _, err := image.Decode(fo)
	require.NoError(t, err)
	mimg, _, err := image.Decode(fm)
	require.NoError(t, err)

	src, err := uv.NewSource(oimg, mimg)
	require.NoError(t, err)

	fs, err := os.Open("../testdata/source.character_walk.png")
	require.NoError(t, err)
	defer fs.Close()

	simg, _, err := image.Decode(fs)
	require.NoError(t, err)

	assert.Equal(t, simg, src)
}

func TestApply(t *testing.T) {
	fs, err := os.Open("../testdata/source.character_walk.png")
	require.NoError(t, err)
	defer fs.Close()

	fl, err := os.Open("../testdata/lookup.character_basic.png")
	require.NoError(t, err)
	defer fl.Close()

	simg, _, err := image.Decode(fs)
	require.NoError(t, err)
	limg, _, err := image.Decode(fl)
	require.NoError(t, err)

	img := uv.Apply(simg, limg)

	fcw, err := os.Open("../testdata/character_walk.png")
	require.NoError(t, err)
	defer fcw.Close()

	cwimg, _, err := image.Decode(fcw)
	require.NoError(t, err)

	assert.Equal(t, img, cwimg)
}

func TestNewSource_DuplicateColorInMap(t *testing.T) {
	// Create a 2x1 map with the same color in both pixels
	m := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	c := color.NRGBA{R: 100, G: 50, B: 25, A: 255}
	m.Set(0, 0, c)
	m.Set(1, 0, c)

	o := image.NewNRGBA(image.Rect(0, 0, 1, 1))

	_, err := uv.NewSource(o, m)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already present in the Map")
}

func TestNewSource_ColorNotFoundInMap(t *testing.T) {
	// Map has one color, overlay has a different one
	m := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	m.Set(0, 0, color.NRGBA{R: 100, G: 50, B: 25, A: 255})

	o := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	o.Set(0, 0, color.NRGBA{R: 200, G: 100, B: 50, A: 255})

	_, err := uv.NewSource(o, m)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "was not found on the Map")
}

func TestNewSource_MapTooLarge(t *testing.T) {
	m := image.NewNRGBA(image.Rect(0, 0, 256, 256))
	o := image.NewNRGBA(image.Rect(0, 0, 1, 1))

	_, err := uv.NewSource(o, m)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceed the maximum of 255x255")
}

func TestNewSource_TransparentPixelsSkipped(t *testing.T) {
	// Map with one opaque and one transparent pixel
	m := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	m.Set(0, 0, color.NRGBA{R: 100, G: 50, B: 25, A: 255})
	m.Set(1, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 0}) // transparent

	// Overlay uses only the opaque color
	o := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	o.Set(0, 0, color.NRGBA{R: 100, G: 50, B: 25, A: 255})

	src, err := uv.NewSource(o, m)
	require.NoError(t, err)

	expected := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	expected.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	assert.Equal(t, expected, src)
}

func TestNewSource_OverlayTransparentPixelsSkipped(t *testing.T) {
	m := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	m.Set(0, 0, color.NRGBA{R: 100, G: 50, B: 25, A: 255})

	// Overlay with a transparent pixel — should not error even though
	// the transparent color isn't in the map
	o := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	o.Set(0, 0, color.NRGBA{R: 100, G: 50, B: 25, A: 255})
	o.Set(1, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 0}) // transparent

	src, err := uv.NewSource(o, m)
	require.NoError(t, err)

	expected := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	expected.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	// (1,0) stays zero-value (transparent)
	assert.Equal(t, expected, src)
}

func TestNewSource_SmallImage(t *testing.T) {
	// 3x1 map with 3 distinct colors
	m := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	c0 := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	c1 := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	c2 := color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	m.Set(0, 0, c0)
	m.Set(1, 0, c1)
	m.Set(2, 0, c2)

	// 3x1 overlay using those colors
	o := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	o.Set(0, 0, c2)
	o.Set(1, 0, c0)
	o.Set(2, 0, c1)

	src, err := uv.NewSource(o, m)
	require.NoError(t, err)

	expected := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	expected.Set(0, 0, color.NRGBA{R: 2, G: 0, B: 0, A: 255}) // c2 is at (2,0)
	expected.Set(1, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255}) // c0 is at (0,0)
	expected.Set(2, 0, color.NRGBA{R: 1, G: 0, B: 0, A: 255}) // c1 is at (1,0)
	assert.Equal(t, expected, src)
}

func TestApply_TransparentSourcePixelsSkipped(t *testing.T) {
	// Source with one opaque and one transparent pixel
	s := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	s.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	s.Set(1, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 0}) // transparent

	// Lookup with a color at (0,0)
	l := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	l.Set(0, 0, color.NRGBA{R: 42, G: 42, B: 42, A: 255})

	img := uv.Apply(s, l)

	expected := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	expected.Set(0, 0, color.NRGBA{R: 42, G: 42, B: 42, A: 255})
	// (1,0) stays transparent
	assert.Equal(t, expected, img)
}

func TestValidate_NoMismatches(t *testing.T) {
	m := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	c0 := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	c1 := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	m.Set(0, 0, c0)
	m.Set(1, 0, c1)

	o := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	o.Set(0, 0, c0)
	o.Set(1, 0, c1)

	mismatches, err := uv.Validate(o, m)
	require.NoError(t, err)
	assert.Empty(t, mismatches)
}

func TestValidate_WithMismatches(t *testing.T) {
	m := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	m.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})

	o := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	o.Set(0, 0, color.NRGBA{R: 0, G: 255, B: 0, A: 255}) // not in map
	o.Set(1, 0, color.NRGBA{R: 0, G: 0, B: 255, A: 255}) // not in map

	mismatches, err := uv.Validate(o, m)
	require.NoError(t, err)
	require.Len(t, mismatches, 2)
	assert.Equal(t, 0, mismatches[0].X)
	assert.Equal(t, 0, mismatches[0].Y)
	assert.Equal(t, 1, mismatches[1].X)
	assert.Equal(t, 0, mismatches[1].Y)
}

func TestValidate_DuplicateColorInMap(t *testing.T) {
	m := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	c := color.NRGBA{R: 100, G: 50, B: 25, A: 255}
	m.Set(0, 0, c)
	m.Set(1, 0, c)

	o := image.NewNRGBA(image.Rect(0, 0, 1, 1))

	_, err := uv.Validate(o, m)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already present in the Map")
}

func TestValidate_TransparentPixelsSkipped(t *testing.T) {
	m := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	m.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})

	o := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	o.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	o.Set(1, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 0}) // transparent — should be skipped

	mismatches, err := uv.Validate(o, m)
	require.NoError(t, err)
	assert.Empty(t, mismatches)
}

func TestValidate_RealImages(t *testing.T) {
	fo, err := os.Open("../testdata/overlay.character_walk.png")
	require.NoError(t, err)
	defer fo.Close()

	fm, err := os.Open("../testdata/map.character.png")
	require.NoError(t, err)
	defer fm.Close()

	oimg, _, err := image.Decode(fo)
	require.NoError(t, err)
	mimg, _, err := image.Decode(fm)
	require.NoError(t, err)

	mismatches, err := uv.Validate(oimg, mimg)
	require.NoError(t, err)
	assert.Empty(t, mismatches)
}

func TestDiff_IdenticalOverlays(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	img.Set(1, 0, color.NRGBA{R: 0, G: 255, B: 0, A: 255})
	img.Set(0, 1, color.NRGBA{R: 0, G: 0, B: 255, A: 255})
	img.Set(1, 1, color.NRGBA{R: 0, G: 0, B: 0, A: 0})

	diff := uv.Diff(img, img)

	expected := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	assert.Equal(t, expected, diff)
}

func TestDiff_SinglePixelDifference(t *testing.T) {
	o1 := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	o1.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	o1.Set(1, 0, color.NRGBA{R: 0, G: 255, B: 0, A: 255})

	o2 := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	o2.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255}) // same
	o2.Set(1, 0, color.NRGBA{R: 0, G: 0, B: 255, A: 255}) // different

	diff := uv.Diff(o1, o2)

	expected := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	expected.Set(1, 0, color.NRGBA{R: 255, G: 0, B: 255, A: 255}) // magenta
	assert.Equal(t, expected, diff)
}

func TestDiff_DifferentSizes(t *testing.T) {
	o1 := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	o1.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})

	o2 := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	o2.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255}) // same
	o2.Set(1, 0, color.NRGBA{R: 0, G: 255, B: 0, A: 255}) // extra
	o2.Set(0, 1, color.NRGBA{R: 0, G: 0, B: 255, A: 255}) // extra
	o2.Set(1, 1, color.NRGBA{R: 0, G: 0, B: 0, A: 0})     // both transparent

	diff := uv.Diff(o1, o2)

	// Output covers union (2x2)
	assert.Equal(t, image.Rect(0, 0, 2, 2), diff.Bounds())

	magenta := color.NRGBA{R: 255, G: 0, B: 255, A: 255}
	transparent := color.NRGBA{}

	nrgba := diff.(*image.NRGBA)
	assert.Equal(t, transparent, nrgba.NRGBAAt(0, 0)) // same pixel
	assert.Equal(t, magenta, nrgba.NRGBAAt(1, 0))     // extra in o2
	assert.Equal(t, magenta, nrgba.NRGBAAt(0, 1))     // extra in o2
	assert.Equal(t, transparent, nrgba.NRGBAAt(1, 1)) // both transparent
}

func TestDiff_BothTransparent(t *testing.T) {
	o1 := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	o2 := image.NewNRGBA(image.Rect(0, 0, 1, 1))

	diff := uv.Diff(o1, o2)

	expected := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	assert.Equal(t, expected, diff)
}

func TestApply_TransparentLookupPixelsSkipped(t *testing.T) {
	// Source pointing to (0,0) in lookup
	s := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	s.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

	// Lookup with a transparent pixel at (0,0)
	l := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	l.Set(0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 0})

	img := uv.Apply(s, l)

	expected := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	// stays transparent since lookup pixel is transparent
	assert.Equal(t, expected, img)
}
