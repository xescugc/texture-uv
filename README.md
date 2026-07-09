# Texture UV

[![Go Reference](https://pkg.go.dev/badge/github.com/xescugc/texture-uv.svg)](https://pkg.go.dev/github.com/xescugc/texture-uv)

Texture UV implements the palette-swapping technique from [aarthificial's "Pixel Art Animation. Reinvented"](https://www.youtube.com/watch?v=HsOKwUwL1bE). It generates UV-coordinate images (Sources) from sprite Overlays and color Maps, then applies any Lookup texture to produce recolored animations — like skins on a character.

## Concepts

The system revolves around four image types:

| Image | Description |
|-------|-------------|
| **Map** | A small image (max 255x255) where each pixel has a unique color. Each color represents a "slot" in the texture. |
| **Overlay** | A spritesheet where each pixel is painted with colors from the Map, indicating which Map slot that pixel belongs to. |
| **Source** | A UV-coordinate image generated from an Overlay + Map. Each pixel stores the Map coordinates in its R (X) and G (Y) channels. This is the reusable artifact. |
| **Lookup** | A texture/skin image with the same dimensions as the Map. Pixel colors in the Lookup are mapped onto the Source to produce the final output. |

### Workflow

```
Map + Overlay  ──▶  Source (UV coords)
                        │
Source + Lookup  ──▶  Final spritesheet
```

1. Create a **Map** (manually or via `new-map`)
2. Paint an **Overlay** spritesheet using Map colors
3. Generate a **Source** from the Overlay + Map (one-time)
4. Apply any **Lookup** to the Source to produce a recolored spritesheet

You can swap Lookups freely to produce different skins without regenerating the Source.

## Installation

Download prebuilt binaries from the [releases page](https://github.com/xescugc/texture-uv/releases), or install with Go:

```
go install github.com/xescugc/texture-uv@latest
```

## CLI Usage

### `new-map`

Generates a Map image from a template. Every non-transparent pixel in the template gets a unique color based on its position (R=x, G=y, B=0, A=255).

```
texture-uv new-map --from-template template.png -o map.character.png
```

| Flag | Description |
|------|-------------|
| `--from-template` | Path to the template/mask image (required) |
| `-o, --output` | Output path for the generated Map (required) |

### `new-source`

Creates a Source (UV coordinate image) from an Overlay and a Map.

```
texture-uv new-source overlay.character_walk.png map.character.png -o source.character_walk.png
```

| Argument | Description |
|----------|-------------|
| 1st | Path to the Overlay image |
| 2nd | Path to the Map image |

| Flag | Description |
|------|-------------|
| `-o, --output` | Output path for the generated Source (required) |

### `apply`

Applies a Lookup texture to a Source, producing the final recolored spritesheet.

```
texture-uv apply source.character_walk.png lookup.character_basic.png -o character_walk.png
```

| Argument | Description |
|----------|-------------|
| 1st | Path to the Source image |
| 2nd | Path to the Lookup image |

| Flag | Description |
|------|-------------|
| `-o, --output` | Output path for the result image (required) |

### `batch-apply`

Applies multiple Lookup textures from a directory to the same Source.

```
texture-uv batch-apply source.character_walk.png ./lookups/ -o ./output/
```

| Argument | Description |
|----------|-------------|
| 1st | Path to the Source image |
| 2nd | Directory containing Lookup PNG files |

| Flag | Description |
|------|-------------|
| `-o, --output` | Output directory (required, created if missing) |

### `preview`

Applies a Lookup to a Source and opens the result in your system's default image viewer.

```
texture-uv preview source.character_walk.png lookup.character_basic.png
```

| Argument | Description |
|----------|-------------|
| 1st | Path to the Source image |
| 2nd | Path to the Lookup image |

### `generate`

Scans a directory for Overlay/Map pairs using a naming convention and generates Source files. Pairs are matched by name: `overlay.<name>.png` pairs with `map.<name>.png` to produce `source.<name>.png`. Source files are skipped if they are newer than both their Overlay and Map.

```
texture-uv generate ./assets/
```

| Argument | Description |
|----------|-------------|
| 1st | Directory to scan (defaults to `.`) |

### `validate`

Checks that every non-transparent pixel in an Overlay has a matching color in the Map.

```
texture-uv validate overlay.character_walk.png map.character.png
```

| Argument | Description |
|----------|-------------|
| 1st | Path to the Overlay image |
| 2nd | Path to the Map image |

### `diff`

Compares two Overlay images pixel-by-pixel and outputs an image highlighting differences in magenta.

```
texture-uv diff overlay_v1.png overlay_v2.png -o diff.png
```

| Argument | Description |
|----------|-------------|
| 1st | Path to the first Overlay image |
| 2nd | Path to the second Overlay image |

| Flag | Description |
|------|-------------|
| `-o, --output` | Output path for the diff image (required) |

### `new-overlay`

Generates an Overlay from a sprite image, a color mapping JSON file, and a Map image. The color mapping defines how sprite pixel colors correspond to UV regions in the Map.

```
texture-uv new-overlay sprite.png color_mapping.json map.character.png -o overlay.character_walk.png
```

| Argument | Description |
|----------|-------------|
| 1st | Path to the sprite image |
| 2nd | Path to the color mapping JSON file |
| 3rd | Path to the Map image |

| Flag | Description |
|------|-------------|
| `-o, --output` | Output path for the generated Overlay (required) |

The color mapping JSON has the following structure:

```json
{
  "frame_size": [32, 32],
  "columns": {
    "0": {
      "groups": [
        {
          "uv": [0, 0, 16, 16],
          "colors": [[255, 0, 0], [0, 255, 0]]
        }
      ]
    }
  }
}
```

- `frame_size`: Width and height of each frame in the spritesheet
- `columns`: Keyed by column index (as a string). Each column represents a viewing direction or variant.
- `groups[].uv`: `[x, y, width, height]` region in the Map image
- `groups[].colors`: List of `[R, G, B]` colors to match in the sprite

## Go Library (`uv` package)

```
go get github.com/xescugc/texture-uv
```

```go
import "github.com/xescugc/texture-uv/uv"
```

### `uv.NewMap`

```go
func NewMap(template image.Image) (image.Image, error)
```

Generates a Map image from a template. Each non-transparent pixel gets a unique color based on its position (R=x, G=y, B=0, A=255). The template must not exceed 255x255 pixels.

### `uv.NewSource`

```go
func NewSource(o, m image.Image) (image.Image, error)
```

Creates a Source (UV coordinate) image from an Overlay (`o`) and a Map (`m`). Each non-transparent Overlay pixel is looked up in the Map by color, and the corresponding Map position is encoded as R=x, G=y in the output. Returns an error if the Map exceeds 255x255, contains duplicate colors, or if an Overlay color is not found in the Map.

### `uv.Apply`

```go
func Apply(s, l image.Image) image.Image
```

Applies a Lookup (`l`) to a Source (`s`). For each non-transparent Source pixel, reads the R and G channels as UV coordinates, samples the Lookup at those coordinates, and writes the result. Transparent Source pixels and transparent Lookup pixels remain transparent.

### `uv.NewOverlay`

```go
func NewOverlay(sprite, mapImg image.Image, mapping ColorMapping) (image.Image, error)
```

Generates an Overlay from a sprite image, a Map image, and a `ColorMapping`. For each group in each column, it finds matching pixels in the sprite frame by color, computes their bounding box, and maps them 1:1 into the UV region specified by the group.

### `uv.Validate`

```go
func Validate(o, m image.Image) ([]Mismatch, error)
```

Checks that every non-transparent pixel in the Overlay (`o`) has a matching color in the Map (`m`). Returns a list of `Mismatch` values for pixels with no match, and an error if the Map contains duplicate colors.

### `uv.Diff`

```go
func Diff(o1, o2 image.Image) image.Image
```

Compares two Overlay images pixel-by-pixel. Pixels that differ are highlighted in magenta; identical or both-transparent pixels are transparent. The output is sized to the union of both bounds.

### Types

```go
// Mismatch represents a non-transparent overlay pixel whose color is not present in the map.
type Mismatch struct {
    X, Y       int
    R, G, B, A uint32
}

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
    UV     [4]int   `json:"uv"`     // [x, y, width, height]
    Colors [][3]int `json:"colors"` // [[R, G, B], ...]
}
```

## GPU Shader (`ebitenuv` package)

The `ebitenuv` sub-module provides a GPU-based alternative to `uv.Apply` for [Ebiten](https://ebitengine.org/) games. Instead of generating N spritesheets in memory (one per skin), you keep a single Source texture and N small Lookup textures, composited on the GPU each frame via a Kage shader.

```
go get github.com/xescugc/texture-uv/ebitenuv
```

```go
import "github.com/xescugc/texture-uv/ebitenuv"
```

### `ebitenuv.DrawWithLookup`

```go
func DrawWithLookup(dst, source, lookup *ebiten.Image, opts *ebiten.DrawRectShaderOptions)
```

Draws the Source image onto `dst`, replacing UV-encoded pixels with colors sampled from the Lookup texture. This replicates the logic of `uv.Apply` but runs entirely on the GPU. The shader is compiled lazily on first call via `sync.Once`.

- `dst`: Destination image (typically the screen)
- `source`: The Source (UV coordinate) image as an `*ebiten.Image`
- `lookup`: The Lookup (skin/texture) image as an `*ebiten.Image`
- `opts`: Draw options for positioning via `GeoM`, blending, etc. Pass `nil` for defaults.

### `ebitenuv.SourceFromImage`

```go
func SourceFromImage(img image.Image) *ebiten.Image
```

Converts a standard `image.Image` to an `*ebiten.Image` suitable for use as the `source` parameter in `DrawWithLookup`.

### `ebitenuv.LookupFromImage`

```go
func LookupFromImage(img image.Image) *ebiten.Image
```

Converts a standard `image.Image` to an `*ebiten.Image` suitable for use as the `lookup` parameter in `DrawWithLookup`.

### Example

```go
package main

import (
    "image/png"
    "os"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/xescugc/texture-uv/ebitenuv"
)

type Game struct {
    source        *ebiten.Image
    currentLookup *ebiten.Image
}

func (g *Game) Draw(screen *ebiten.Image) {
    opts := &ebiten.DrawRectShaderOptions{}
    opts.GeoM.Translate(100, 50)
    ebitenuv.DrawWithLookup(screen, g.source, g.currentLookup, opts)
}

func (g *Game) Update() error { return nil }
func (g *Game) Layout(w, h int) (int, int) { return w, h }

func main() {
    // Load source and lookup as standard images, then convert
    f, _ := os.Open("source.character_walk.png")
    srcImg, _ := png.Decode(f)
    f.Close()

    f, _ = os.Open("lookup.character_basic.png")
    lookupImg, _ := png.Decode(f)
    f.Close()

    game := &Game{
        source:        ebitenuv.SourceFromImage(srcImg),
        currentLookup: ebitenuv.LookupFromImage(lookupImg),
    }
    ebiten.RunGame(game)
}
```

### `go:generate`

To automate Source generation, you can use `//go:generate` directives:

```go
//go:generate texture-uv new-source overlay.character_walk.png map.character.png -o source.character_walk.png
```

## Visual Example

**Overlay** (uses Map colors to mark regions):

<img src="testdata/overlay.character_walk.png" width=200>

**Map** (each pixel has a unique color):

<img src="testdata/map.character.png" width=200>

**Source** (generated UV coordinates — R,G encode position):

<img src="testdata/source.character_walk.png" width=200>

**Lookup** (the skin/texture to apply):

<img src="testdata/lookup.character_basic.png" width=200>

**Result** (`apply` Source + Lookup):

<img src="testdata/character_walk.png" width=200>

**Different Lookup, different result:**

<img src="testdata/lookup.character_helmet.png" width=200>
<img src="testdata/character_helmet_walk.png" width=200>

## Inspiration

This library is 100% inspired by the "Pixel Art Animation. Reinvented - Astortion Devlog" video from aarthificial, which explains the technique in detail:

[![Pixel Art Animation. Reinvented - Astortion Devlog](https://img.youtube.com/vi/HsOKwUwL1bE/0.jpg)](https://www.youtube.com/watch?v=HsOKwUwL1bE)
