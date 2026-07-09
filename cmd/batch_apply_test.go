package cmd

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writePNG(t *testing.T, path string, img image.Image) {
	t.Helper()
	f, err := os.Create(path)
	require.NoError(t, err)
	require.NoError(t, png.Encode(f, img))
	require.NoError(t, f.Close())
}

func makeTestImage(w, h int, c color.Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func TestBatchApply_MultipleLookups(t *testing.T) {
	srcDir := t.TempDir()
	lookupsDir := t.TempDir()
	outputDir := filepath.Join(t.TempDir(), "out")

	sourceImg := makeTestImage(2, 2, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	writePNG(t, filepath.Join(srcDir, "source.png"), sourceImg)

	lookup1 := makeTestImage(2, 2, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	lookup2 := makeTestImage(2, 2, color.NRGBA{R: 0, G: 255, B: 0, A: 255})
	writePNG(t, filepath.Join(lookupsDir, "skin1.png"), lookup1)
	writePNG(t, filepath.Join(lookupsDir, "skin2.png"), lookup2)

	cmd := batchApplyCmd
	err := cmd.Run(context.Background(), []string{"batch-apply", "-o", outputDir, filepath.Join(srcDir, "source.png"), lookupsDir})
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(outputDir, "skin1.png"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(outputDir, "skin2.png"))
	assert.NoError(t, err)
}

func TestBatchApply_EmptyDir(t *testing.T) {
	srcDir := t.TempDir()
	lookupsDir := t.TempDir()
	outputDir := filepath.Join(t.TempDir(), "out")

	sourceImg := makeTestImage(2, 2, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	writePNG(t, filepath.Join(srcDir, "source.png"), sourceImg)

	cmd := batchApplyCmd
	err := cmd.Run(context.Background(), []string{"batch-apply", "-o", outputDir, filepath.Join(srcDir, "source.png"), lookupsDir})
	require.NoError(t, err)

	entries, err := os.ReadDir(outputDir)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestBatchApply_SkipsSubdirectories(t *testing.T) {
	srcDir := t.TempDir()
	lookupsDir := t.TempDir()
	outputDir := filepath.Join(t.TempDir(), "out")

	sourceImg := makeTestImage(2, 2, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	writePNG(t, filepath.Join(srcDir, "source.png"), sourceImg)

	require.NoError(t, os.Mkdir(filepath.Join(lookupsDir, "subdir"), 0755))

	cmd := batchApplyCmd
	err := cmd.Run(context.Background(), []string{"batch-apply", "-o", outputDir, filepath.Join(srcDir, "source.png"), lookupsDir})
	require.NoError(t, err)

	entries, err := os.ReadDir(outputDir)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestBatchApply_NonPNGFilesIgnored(t *testing.T) {
	srcDir := t.TempDir()
	lookupsDir := t.TempDir()
	outputDir := filepath.Join(t.TempDir(), "out")

	sourceImg := makeTestImage(2, 2, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	writePNG(t, filepath.Join(srcDir, "source.png"), sourceImg)

	require.NoError(t, os.WriteFile(filepath.Join(lookupsDir, "readme.txt"), []byte("hello"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(lookupsDir, "data.json"), []byte("{}"), 0644))

	cmd := batchApplyCmd
	err := cmd.Run(context.Background(), []string{"batch-apply", "-o", outputDir, filepath.Join(srcDir, "source.png"), lookupsDir})
	require.NoError(t, err)

	entries, err := os.ReadDir(outputDir)
	require.NoError(t, err)
	assert.Empty(t, entries)
}
