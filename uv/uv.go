package uv

import (
	"fmt"
	"image"
	"image/color"
)

// colorKey is a normalized color representation used as a map key,
// so that colors from different image types (RGBA, NRGBA, etc.)
// can be compared correctly.
type colorKey struct {
	r, g, b, a uint32
}

func toColorKey(c color.Color) colorKey {
	r, g, b, a := c.RGBA()
	return colorKey{r, g, b, a}
}

// NewSource returns a source image(UV) that then can be used to Apply a Lookup.
// It expects an Overlay image and a Map image.
// The Map image must not exceed 255x255 pixels.
func NewSource(o, m image.Image) (image.Image, error) {
	uvm := make(map[colorKey]image.Point)

	mbounds := m.Bounds()

	if mbounds.Dx() > 255 || mbounds.Dy() > 255 {
		return nil, fmt.Errorf("Map image dimensions (%dx%d) exceed the maximum of 255x255", mbounds.Dx(), mbounds.Dy())
	}

	for x := mbounds.Min.X; x < mbounds.Max.X; x++ {
		for y := mbounds.Min.Y; y < mbounds.Max.Y; y++ {
			c := m.At(x, y)
			ck := toColorKey(c)
			if _, ok := uvm[ck]; ok {
				return nil, fmt.Errorf("The color %s is already present in the Map", c)
			}
			_, _, _, a := c.RGBA()
			if a > 0 {
				uvm[ck] = image.Pt(x-mbounds.Min.X, y-mbounds.Min.Y)
			}
		}
	}

	obounds := o.Bounds()

	srcimg := image.NewNRGBA(obounds)
	for x := obounds.Min.X; x < obounds.Max.X; x++ {
		for y := obounds.Min.Y; y < obounds.Max.Y; y++ {
			c := o.At(x, y)
			_, _, _, a := c.RGBA()
			if a > 0 {
				ck := toColorKey(c)
				uv, ok := uvm[ck]
				if !ok {
					return nil, fmt.Errorf("The color %s was not found on the Map", c)
				}
				srcimg.Set(x, y, color.NRGBA{uint8(uv.X), uint8(uv.Y), 0, 255})
			}
		}
	}

	return srcimg, nil
}

// Apply applies to the Source(s) the Lookup(l)
func Apply(s, l image.Image) image.Image {
	sbounds := s.Bounds()

	img := image.NewNRGBA(sbounds)
	for x := sbounds.Min.X; x < sbounds.Max.X; x++ {
		for y := sbounds.Min.Y; y < sbounds.Max.Y; y++ {
			c := s.At(x, y)
			r, g, _, a := c.RGBA()
			if a > 0 {
				// r and g are pre-scaled to 0-65535, convert back to 0-255
				uvX := int(r >> 8)
				uvY := int(g >> 8)
				lc := l.At(uvX, uvY)
				_, _, _, la := lc.RGBA()
				if la > 0 {
					img.Set(x, y, lc)
				}
			}
		}
	}

	return img
}
