package uv

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
)

// ColorMapping describes how sprite colors map to UV regions per column.
type ColorMapping struct {
	FrameSize [2]int                  `json:"frame_size"`
	Columns   map[string]ColumnConfig `json:"columns"`
}

// ColumnConfig holds the groups for a single column (viewing direction).
type ColumnConfig struct {
	Groups []Group `json:"groups"`
}

// Group maps a set of sprite colors to a UV region in the map image.
type Group struct {
	UV     [4]int   `json:"uv"`
	Colors [][3]int `json:"colors"`
}

// NewOverlay generates an overlay image from a sprite, a map image, and a color mapping.
// For each group in each column, it finds the bounding box of matching pixels in the frame
// and maps them 1:1 to the UV region specified in the group.
func NewOverlay(sprite, mapImg image.Image, mapping ColorMapping) (image.Image, error) {
	type point struct{ x, y int }
	type groupColorSet struct {
		group    Group
		colorSet map[colorKey]struct{}
	}

	sb := sprite.Bounds()
	frameW := mapping.FrameSize[0]
	frameH := mapping.FrameSize[1]

	numRows := sb.Dy() / frameH

	out := image.NewNRGBA(sb)

	for colStr, colCfg := range mapping.Columns {
		col, err := strconv.Atoi(colStr)
		if err != nil {
			return nil, fmt.Errorf("invalid column key %q: %w", colStr, err)
		}
		groupSets := make([]groupColorSet, len(colCfg.Groups))
		for i, group := range colCfg.Groups {
			cs := make(map[colorKey]struct{}, len(group.Colors))
			for _, c := range group.Colors {
				ck := toColorKey(color.NRGBA{
					R: uint8(c[0]),
					G: uint8(c[1]),
					B: uint8(c[2]),
					A: 255,
				})
				cs[ck] = struct{}{}
			}
			groupSets[i] = groupColorSet{group: group, colorSet: cs}
		}

		for row := 0; row < numRows; row++ {
			frameX := sb.Min.X + col*frameW
			frameY := sb.Min.Y + row*frameH

			for _, gs := range groupSets {
				group := gs.group
				colorSet := gs.colorSet

				// Pass 1: find matching pixels and bounding box
				var matches []point
				bbMinX, bbMinY := frameX+frameW, frameY+frameH

				for x := frameX; x < frameX+frameW; x++ {
					for y := frameY; y < frameY+frameH; y++ {
						c := sprite.At(x, y)
						_, _, _, a := c.RGBA()
						if a == 0 {
							continue
						}
						ck := toColorKey(c)
						if _, ok := colorSet[ck]; ok {
							matches = append(matches, point{x, y})
							if x < bbMinX {
								bbMinX = x
							}
							if y < bbMinY {
								bbMinY = y
							}
						}
					}
				}

				if len(matches) == 0 {
					continue
				}

				// Pass 2: map matching pixels to UV region
				uvX, uvY, uvW, uvH := group.UV[0], group.UV[1], group.UV[2], group.UV[3]

				for _, p := range matches {
					rx := p.x - bbMinX
					ry := p.y - bbMinY
					if rx >= uvW || ry >= uvH {
						return nil, fmt.Errorf(
							"bounding box exceeds UV region: pixel offset (%d,%d) outside UV size (%d,%d) in column %d, row %d",
							rx, ry, uvW, uvH, col, row,
						)
					}
					out.Set(p.x, p.y, mapImg.At(uvX+rx, uvY+ry))
				}
			}
		}
	}

	return out, nil
}
