package ebitenuv

import (
	"image"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	shaderOnce sync.Once
	uvShader   *ebiten.Shader
)

func shader() *ebiten.Shader {
	shaderOnce.Do(func() {
		s, err := ebiten.NewShader([]byte(shaderSource))
		if err != nil {
			panic("ebitenuv: failed to compile shader: " + err.Error())
		}
		uvShader = s
	})
	return uvShader
}

// DrawWithLookup draws the source image onto dst, replacing UV-encoded pixels
// with colors sampled from the lookup texture. This replicates the logic of
// uv.Apply but runs entirely on the GPU.
func DrawWithLookup(dst, source, lookup *ebiten.Image, opts *ebiten.DrawRectShaderOptions) {
	if opts == nil {
		opts = &ebiten.DrawRectShaderOptions{}
	}
	bounds := source.Bounds()
	opts.Images[0] = source
	opts.Images[1] = lookup
	dst.DrawRectShader(bounds.Dx(), bounds.Dy(), shader(), opts)
}

// SourceFromImage converts a standard image.Image to an *ebiten.Image suitable
// for use as the source parameter in DrawWithLookup.
func SourceFromImage(img image.Image) *ebiten.Image {
	return ebiten.NewImageFromImage(img)
}

// LookupFromImage converts a standard image.Image to an *ebiten.Image suitable
// for use as the lookup parameter in DrawWithLookup.
func LookupFromImage(img image.Image) *ebiten.Image {
	return ebiten.NewImageFromImage(img)
}
