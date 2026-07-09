//go:build ebitenuv_test

package ebitenuv

import (
	"errors"
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type testGame struct {
	fn  func()
	ran bool
}

func (g *testGame) Update() error {
	if !g.ran {
		g.ran = true
		g.fn()
	}
	return ebiten.Termination
}

func (g *testGame) Draw(screen *ebiten.Image) {}

func (g *testGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func runGameFunc(fn func()) {
	g := &testGame{fn: fn}
	if err := ebiten.RunGame(g); err != nil && !errors.Is(err, ebiten.Termination) {
		panic(err)
	}
}

func TestShaderCompiles(t *testing.T) {
	runGameFunc(func() {
		s := shader()
		if s == nil {
			t.Fatal("shader() returned nil")
		}
	})
}

func TestDrawWithLookup_MatchesCPU(t *testing.T) {
	runGameFunc(func() {
		const size = 4

		// Build a source image where R=x, G=y (matching uv.NewSource encoding)
		srcImg := image.NewNRGBA(image.Rect(0, 0, size, size))
		for x := 0; x < size; x++ {
			for y := 0; y < size; y++ {
				srcImg.Set(x, y, color.NRGBA{uint8(x), uint8(y), 0, 255})
			}
		}

		// Build a lookup where each pixel has a distinct color
		lookupImg := image.NewNRGBA(image.Rect(0, 0, size, size))
		for x := 0; x < size; x++ {
			for y := 0; y < size; y++ {
				lookupImg.Set(x, y, color.NRGBA{
					uint8(x * 60),
					uint8(y * 60),
					uint8((x + y) * 30),
					255,
				})
			}
		}

		source := ebiten.NewImageFromImage(srcImg)
		lookup := ebiten.NewImageFromImage(lookupImg)
		dst := ebiten.NewImage(size, size)

		DrawWithLookup(dst, source, lookup, nil)

		// Compare output against expected CPU result
		for x := 0; x < size; x++ {
			for y := 0; y < size; y++ {
				got := dst.At(x, y)
				want := lookupImg.At(x, y)

				gr, gg, gb, ga := got.RGBA()
				wr, wg, wb, wa := want.RGBA()

				// Allow ±1 tolerance per channel (GPU float precision)
				if diff(gr, wr) > 1 || diff(gg, wg) > 1 || diff(gb, wb) > 1 || diff(ga, wa) > 1 {
					t.Errorf("pixel (%d,%d): got %v, want %v", x, y, got, want)
				}
			}
		}
	})
}

func diff(a, b uint32) uint32 {
	return uint32(math.Abs(float64(a>>8) - float64(b>>8)))
}
