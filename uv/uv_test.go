package uv_test

import (
	"image"
	_ "image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xescugc/golang-texture-uv/uv"
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
