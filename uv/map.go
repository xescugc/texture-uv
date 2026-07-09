package uv

import (
	"fmt"
	"image"
	"image/color"
)

// NewMap generates a Map image from a template/mask image.
// Each non-transparent pixel gets a unique color based on its position: R=x, G=y, B=0, A=255.
// The template must not exceed 255x255 pixels.
func NewMap(template image.Image) (image.Image, error) {
	bounds := template.Bounds()

	if bounds.Dx() > 255 || bounds.Dy() > 255 {
		return nil, fmt.Errorf("template dimensions (%dx%d) exceed the maximum of 255x255", bounds.Dx(), bounds.Dy())
	}

	out := image.NewNRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			_, _, _, a := template.At(x, y).RGBA()
			if a > 0 {
				out.Set(x, y, color.NRGBA{
					R: uint8(x - bounds.Min.X),
					G: uint8(y - bounds.Min.Y),
					B: 0,
					A: 255,
				})
			}
		}
	}

	return out, nil
}
