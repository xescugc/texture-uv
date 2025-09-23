package uv

import (
	"fmt"
	"image"
	"image/color"
)

type UVMap map[color.Color]image.Point

// NewSource returns a source image(UV) that then can be used to Apply a Lookup.
// It expects an Overlay image and a Map image.
func NewSource(o, m image.Image) (image.Image, error) {
	uvm := make(UVMap)

	mbounds := m.Bounds()
	mwidth, mheight := mbounds.Max.X, mbounds.Max.Y

	for x := range mwidth {
		for y := range mheight {
			c := m.At(x, y)
			if _, ok := uvm[c]; ok {
				return nil, fmt.Errorf("The color %s is already present in the Map", c)
			}
			_, _, _, a := c.RGBA()
			if a > 0 {
				uvm[c] = image.Pt(x, y)
			}
		}
	}

	obounds := o.Bounds()
	owidth, oheight := obounds.Max.X, obounds.Max.Y

	srcimg := image.NewNRGBA(obounds)
	for x := range owidth {
		for y := range oheight {
			c := o.At(x, y)
			_, _, _, a := c.RGBA()
			if a > 0 {
				uv, ok := uvm[c]
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
	swidth, sheight := sbounds.Max.X, sbounds.Max.Y

	img := image.NewNRGBA(sbounds)
	for x := range swidth {
		for y := range sheight {
			c := s.At(x, y)
			nc := c.(color.NRGBA)
			if nc.A > 0 {
				lc := l.At(int(nc.R), int(nc.G))
				_, _, _, a := lc.RGBA()
				if a > 0 {
					img.Set(x, y, lc)
				}
			}
		}
	}

	return img
}
