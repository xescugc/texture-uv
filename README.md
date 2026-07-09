# Texture UV

Generates a Texture Coordinate file (UV) named Source from a Map and an Overlay that then can be used (the UV) to 
apply a Texture/Lookup to it like a Skin on a character

## How to use it

For the latest version you can download it [here](https://github.com/xescugc/texture-uv/releases)

You need to have 2 things to start with:
* A **Overlay**: Which uses the colors of the unique colors of the Map to overlay the animation
<img src="testdata/overlay.character_walk.png" width=200>

* A **Map**: Which is what it would hold the texture in 2D
<img src="testdata/map.character.png" width=200>

Then you can generate the **Source** with 

```
texture-uv new-source ./testdata/overlay.character_walk.png ./testdata/map.character.png -o ./testdata/source.character_walk.png
```

Which will generate the **Source** (UV coordinates) that we'll be able to apply any texture after

<img src="testdata/source.character_walk.png" width=200>

Then to generate a new animation with a texture you need to use the **Lookup**(which is the texture or skin) and the **Source**

<img src="testdata/lookup.character_basic.png" width=200>

```
texture-uv apply ./testdata/source.character_walk.png ./testdata/lookup.character_basic.png -o ./testdata/character_walk.png
```

<img src="testdata/character_walk.png" width=200>

Now you can use a different **Lookup** to generate a different animation

<img src="testdata/lookup.character_helmet.png" width=200>
<img src="testdata/character_helmet_walk.png" width=200>

## Import Package (If Golang)

You can also use it as an import package, the `uv` has:
* `uv.NewSource(o, m image.Image) image.Image`
* `uv.Apply(s, l image.Image) image.Image` 

Which can then be used dynamically from the code if you want to dynamically generate the images.

### GPU Shader (Ebiten)

The `ebitenuv` sub-module provides a GPU-based alternative to `uv.Apply` using an Ebiten Kage shader. Instead of generating N spritesheets in memory, you keep one Source texture and N small Lookup textures, composited on the GPU per frame.

```go
import "github.com/xescugc/texture-uv/ebitenuv"

// In Ebiten game's Draw method:
func (g *Game) Draw(screen *ebiten.Image) {
    opts := &ebiten.DrawRectShaderOptions{}
    opts.GeoM.Translate(100, 50)
    ebitenuv.DrawWithLookup(screen, g.source, g.currentLookup, opts)
}
```

Install with:
```
go get github.com/xescugc/texture-uv/ebitenuv
```

### `go:generate`

To not have this process be manual, you can run the `new-source` from a `//go:generate` 

## TODO

* Make it so there is a CLI cmd to read a directory and from a naming convention it automatically generates the **Source** files

## Inspiration

This lib is 100% inspired in "Pixel Art Animation. Reinvented - Astortion Devlog" video from "aarthificial" which I recommend to watch to understand
even more how this works

[![Pixel Art Animation. Reinvented - Astortion Devlog](https://img.youtube.com/vi/HsOKwUwL1bE/0.jpg)](https://www.youtube.com/watch?v=HsOKwUwL1bE)
